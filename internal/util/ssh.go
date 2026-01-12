package util

import (
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	DefaultConnectTimeout = 15 * time.Second
	DefaultCommandTimeout = 30 * time.Second
	DefaultMaxAttempts    = 20
	DefaultRetryDelay     = 5 * time.Second
	CloudInitWaitTimeout  = 6 * time.Minute // Allows for 5 min wait in script + overhead
)

// cloudInitWaitScript contains the cloud-init wait script embedded from templates/cloud_init_wait_script.sh
//
//go:embed templates/cloud_init_wait_script.sh
var cloudInitWaitScript string

// SSH represents an SSH client wrapper
type SSH struct {
	privateKeyPath string
	publicKeyPath  string
	bastionHost    string // Bastion/jump host for ProxyJump
	bastionPort    int    // Bastion SSH port
}

// NewSSH creates a new SSH client
func NewSSH(privateKeyPath, publicKeyPath string) *SSH {
	return &SSH{
		privateKeyPath: privateKeyPath,
		publicKeyPath:  publicKeyPath,
	}
}

// SetBastion configures a bastion/jump host for SSH connections
func (s *SSH) SetBastion(host string, port int) {
	s.bastionHost = host
	s.bastionPort = port
}

// CalculateFingerprint calculates the MD5 fingerprint of a public key
// NOTE: MD5 is used here for compatibility with legacy SSH implementations
// and Hetzner Cloud API expectations. While MD5 is cryptographically weak,
// it's acceptable for this non-security-critical display/comparison purpose.
// Modern SSH implementations prefer SHA256 fingerprints.
func CalculateFingerprint(publicKeyPath string) (string, error) {
	data, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key: %w", err)
	}

	// Parse the public key
	parts := strings.Fields(string(data))
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid public key format")
	}

	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	// MD5 hash for legacy compatibility (not for security purposes)
	hash := md5.Sum(keyData)
	fingerprint := fmt.Sprintf("%x", hash)

	// Format as colon-separated pairs
	var result []string
	for i := 0; i < len(fingerprint); i += 2 {
		result = append(result, fingerprint[i:i+2])
	}

	return strings.Join(result, ":"), nil
}

// getSSHConfig returns SSH client configuration
func (s *SSH) getSSHConfig(useAgent bool) (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	if useAgent {
		// Try to use SSH agent
		if agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			authMethods = append(authMethods, ssh.PublicKeysCallback(agent.NewClient(agentConn).Signers))
		}
	}

	// Read private key
	key, err := os.ReadFile(s.privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	authMethods = append(authMethods, ssh.PublicKeys(signer))

	config := &ssh.ClientConfig{
		User: "root",
		Auth: authMethods,
		// WARNING: InsecureIgnoreHostKey disables SSH host key verification.
		// This is used for cluster provisioning where we're connecting to newly
		// created servers with unknown host keys. This is a security trade-off
		// for automation convenience.
		//
		// SECURITY CONSIDERATIONS:
		// - Only use this for initial cluster provisioning
		// - Ensure connection is over a trusted network or VPN
		// - Consider implementing host key storage and verification for production
		// - Be aware this is vulnerable to MITM attacks during initial connection
		//
		// For enhanced security, consider:
		// 1. Using ssh.FixedHostKey() with known host keys
		// 2. Implementing a known_hosts file system
		// 3. Using certificate-based authentication
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         DefaultConnectTimeout,
	}

	return config, nil
}

// getClient returns an SSH client, optionally through a bastion host
func (s *SSH) getClient(config *ssh.ClientConfig, host string, port int) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	// If bastion host is configured, connect through it
	if s.bastionHost != "" {
		// Connect to bastion
		bastionAddr := fmt.Sprintf("%s:%d", s.bastionHost, s.bastionPort)
		bastionClient, err := ssh.Dial("tcp", bastionAddr, config)
		if err != nil {
			return nil, fmt.Errorf("failed to dial bastion %s: %w", bastionAddr, err)
		}

		// Connect to target through bastion
		conn, err := bastionClient.Dial("tcp", addr)
		if err != nil {
			bastionClient.Close()
			return nil, fmt.Errorf("failed to dial %s through bastion: %w", addr, err)
		}

		// Create SSH connection over the tunneled connection
		ncc, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			conn.Close()
			bastionClient.Close()
			return nil, fmt.Errorf("failed to create client connection through bastion: %w", err)
		}

		// Note: The bastion client connection is embedded in the returned SSH client
		// and will be cleaned up when the client is closed
		return ssh.NewClient(ncc, chans, reqs), nil
	}

	// Direct connection (no bastion)
	return ssh.Dial("tcp", addr, config)
}

// Run executes a command on a remote host via SSH
func (s *SSH) Run(ctx context.Context, host string, port int, command string, useAgent bool) (string, error) {
	config, err := s.getSSHConfig(useAgent)
	if err != nil {
		return "", err
	}

	// Connect to the remote host (possibly through bastion)
	client, err := s.getClient(config, host, port)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Execute the command with timeout
	type result struct {
		output string
		err    error
	}

	resultChan := make(chan result, 1)

	go func() {
		output, err := session.CombinedOutput(command)
		resultChan <- result{output: string(output), err: err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultChan:
		return res.output, res.err
	}
}

// RunWithOutput executes a command and streams output
func (s *SSH) RunWithOutput(ctx context.Context, host string, port int, command string, useAgent bool, prefix string) error {
	config, err := s.getSSHConfig(useAgent)
	if err != nil {
		return err
	}

	// Connect to the remote host (possibly through bastion)
	client, err := s.getClient(config, host, port)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Setup pipes
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := session.Start(command); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output
	done := make(chan error, 1)
	go func() {
		done <- session.Wait()
	}()

	// Read output
	go s.streamOutput(stdout, prefix, false)
	go s.streamOutput(stderr, prefix, true)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// streamOutput reads from a reader and prints with prefix
func (s *SSH) streamOutput(reader io.Reader, prefix string, isError bool) {
	buf := make([]byte, 4096)
	leftover := ""

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			text := leftover + string(buf[:n])
			lines := strings.Split(text, "\n")

			for i := 0; i < len(lines)-1; i++ {
				line := strings.TrimRight(lines[i], "\r")
				if line != "" {
					if isError {
						fmt.Fprintf(os.Stderr, "[%s] %s\n", prefix, line)
					} else {
						fmt.Printf("[%s] %s\n", prefix, line)
					}
				}
			}

			leftover = lines[len(lines)-1]
		}

		if err != nil {
			if err != io.EOF && leftover != "" {
				leftover = strings.TrimRight(leftover, "\r")
				if leftover != "" {
					if isError {
						fmt.Fprintf(os.Stderr, "[%s] %s\n", prefix, leftover)
					} else {
						fmt.Printf("[%s] %s\n", prefix, leftover)
					}
				}
			}
			break
		}
	}
}

// WaitForInstance waits for an instance to be ready by running a test command
func (s *SSH) WaitForInstance(ctx context.Context, host string, port int, testCommand string, expectedResult string, useAgent bool, maxAttempts int) error {
	if maxAttempts == 0 {
		maxAttempts = DefaultMaxAttempts
	}

	for i := 0; i < maxAttempts; i++ {
		// Create context with timeout for this attempt
		attemptCtx, cancel := context.WithTimeout(ctx, DefaultCommandTimeout)

		result, err := s.Run(attemptCtx, host, port, testCommand, useAgent)
		cancel()

		if err == nil {
			result = strings.TrimSpace(result)
			if result == expectedResult {
				return nil
			}
		}

		if i < maxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(DefaultRetryDelay):
				// Continue to next attempt
			}
		}
	}

	return fmt.Errorf("instance not ready after %d attempts", maxAttempts)
}

// WaitForCloudInit waits for cloud-init to complete on a remote host
func (s *SSH) WaitForCloudInit(ctx context.Context, host string, port int, useAgent bool) error {
	// Execute the wait script with a timeout that allows for the script's 5 min wait + overhead
	waitCtx, cancel := context.WithTimeout(ctx, CloudInitWaitTimeout)
	defer cancel()

	_, err := s.Run(waitCtx, host, port, cloudInitWaitScript, useAgent)
	if err != nil {
		return fmt.Errorf("cloud-init did not complete: %w", err)
	}

	return nil
}

// CopyFile copies a file to the remote host via SCP
func (s *SSH) CopyFile(ctx context.Context, host string, port int, localPath string, remotePath string, useAgent bool) error {
	config, err := s.getSSHConfig(useAgent)
	if err != nil {
		return err
	}

	// Connect to the remote host (possibly through bastion)
	client, err := s.getClient(config, host, port)
	if err != nil {
		return err
	}
	defer client.Close()

	// Read local file
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Create remote file
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Write file content via SSH
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := session.Start(fmt.Sprintf("cat > %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if _, err := stdin.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	stdin.Close()

	return session.Wait()
}
