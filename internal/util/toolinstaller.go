package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// ToolInstaller handles detection and installation of required tools
type ToolInstaller struct {
	kubectlVersion string
	helmVersion    string
}

// NewToolInstaller creates a new tool installer with versions
func NewToolInstaller(k3sVersion string) (*ToolInstaller, error) {
	kubectlVersion := extractKubectlVersionFromK3s(k3sVersion)
	if kubectlVersion == "" {
		return nil, fmt.Errorf("unable to extract kubectl version from k3s version: %s", k3sVersion)
	}

	return &ToolInstaller{
		kubectlVersion: kubectlVersion,
		helmVersion:    "", // Helm version determined by get-helm script
	}, nil
}

// extractKubectlVersionFromK3s extracts the Kubernetes version from k3s version
// Example: v1.35.0+k3s1 -> v1.35.0
func extractKubectlVersionFromK3s(k3sVersion string) string {
	if k3sVersion == "" {
		return ""
	}

	// Match pattern vX.Y.Z+k3sN and extract vX.Y.Z
	re := regexp.MustCompile(`^(v?\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(k3sVersion)

	if len(matches) > 1 {
		version := matches[1]
		// Ensure it has the 'v' prefix
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
		return version
	}

	return ""
}

// SetHelmVersion allows setting a custom helm version
func (t *ToolInstaller) SetHelmVersion(version string) {
	t.helmVersion = version
}

// GetKubectlVersion returns the kubectl version that will be installed
func (t *ToolInstaller) GetKubectlVersion() string {
	return t.kubectlVersion
}

// GetHelmVersion returns the helm version that will be installed
func (t *ToolInstaller) GetHelmVersion() string {
	return t.helmVersion
}

// IsKubectlInstalled checks if kubectl is available
func (t *ToolInstaller) IsKubectlInstalled() bool {
	_, err := exec.LookPath("kubectl")
	return err == nil
}

// IsHelmInstalled checks if helm is available
func (t *ToolInstaller) IsHelmInstalled() bool {
	_, err := exec.LookPath("helm")
	return err == nil
}

// IsKubectlAIInstalled checks if kubectl-ai is available
func (t *ToolInstaller) IsKubectlAIInstalled() bool {
	_, err := exec.LookPath("kubectl-ai")
	return err == nil
}

// InstallKubectl installs kubectl globally using direct download
func (t *ToolInstaller) InstallKubectl() error {
	fmt.Printf("Installing kubectl %s globally...\n", t.kubectlVersion)

	osName := runtime.GOOS
	if osName != "linux" && osName != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", osName)
	}

	// Validate architecture (supported: amd64, arm64)
	arch := runtime.GOARCH
	if arch != "amd64" && arch != "arm64" {
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	fmt.Printf("Downloading kubectl %s for %s/%s...\n", t.kubectlVersion, osName, arch)

	// Download kubectl binary
	kubectlURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/%s/%s/kubectl", t.kubectlVersion, osName, arch)
	if err := t.runCommand("curl", "-LO", kubectlURL); err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}

	// Download checksum
	checksumURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/%s/%s/kubectl.sha256", t.kubectlVersion, osName, arch)
	if err := t.runCommand("curl", "-LO", checksumURL); err != nil {
		return fmt.Errorf("failed to download kubectl checksum: %w", err)
	}

	// Verify checksum (command differs between Linux and macOS)
	fmt.Println("Verifying kubectl checksum")
	var checksumCmd string
	if osName == "linux" {
		checksumCmd = "echo \"$(cat kubectl.sha256)  kubectl\" | sha256sum --check"
	} else if osName == "darwin" {
		checksumCmd = "echo \"$(cat kubectl.sha256)  kubectl\" | shasum -a 256 --check"
	} else {
		// Should never reach here due to OS validation above
		return fmt.Errorf("unsupported operating system: %s", osName)
	}

	cmd := exec.Command("bash", "-c", checksumCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Clean up downloaded files
		os.Remove("kubectl")
		os.Remove("kubectl.sha256")
		return fmt.Errorf("kubectl checksum verification failed: %w\nOutput: %s", err, string(output))
	}

	// Make executable and move to /usr/local/bin
	if err := os.Chmod("kubectl", 0755); err != nil {
		os.Remove("kubectl")
		os.Remove("kubectl.sha256")
		return fmt.Errorf("failed to make kubectl executable: %w", err)
	}

	if err := t.runCommand("sudo", "mv", "kubectl", "/usr/local/bin/kubectl"); err != nil {
		os.Remove("kubectl")
		os.Remove("kubectl.sha256")
		return fmt.Errorf("failed to move kubectl to /usr/local/bin: %w", err)
	}

	// Clean up checksum file
	os.Remove("kubectl.sha256")

	fmt.Println("✓ kubectl installed successfully to /usr/local/bin/kubectl")
	return nil
}

// InstallHelm installs helm globally using the official get-helm-4 script
// The script automatically detects OS and architecture
func (t *ToolInstaller) InstallHelm() error {
	fmt.Println("Installing helm globally")

	osName := runtime.GOOS
	if osName != "linux" && osName != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", osName)
	}

	fmt.Println("Downloading helm installation script")

	// Download the get-helm-4 script which auto-detects OS and architecture
	scriptURL := "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4"
	if err := t.runCommand("curl", "-fsSL", "-o", "get_helm.sh", scriptURL); err != nil {
		return fmt.Errorf("failed to download helm install script: %w", err)
	}

	// Make the script executable
	if err := os.Chmod("get_helm.sh", 0700); err != nil {
		os.Remove("get_helm.sh")
		return fmt.Errorf("failed to make helm script executable: %w", err)
	}

	// Execute the script
	fmt.Println("Running helm installation script")
	cmd := exec.Command("./get_helm.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove("get_helm.sh")
		return fmt.Errorf("failed to run helm install script: %w", err)
	}

	// Clean up the script
	os.Remove("get_helm.sh")

	fmt.Println("✓ helm installed successfully")
	return nil
}

// InstallKubectlAI installs kubectl-ai globally using the official installation script
// The script automatically detects OS and architecture
func (t *ToolInstaller) InstallKubectlAI() error {
	fmt.Println("Installing kubectl-ai globally")

	osName := runtime.GOOS
	if osName != "linux" && osName != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", osName)
	}

	fmt.Println("Downloading kubectl-ai installation script")

	// Download the kubectl-ai installation script
	scriptURL := "https://raw.githubusercontent.com/GoogleCloudPlatform/kubectl-ai/main/install.sh"
	if err := t.runCommand("curl", "-fsSL", "-o", "install_kubectl_ai.sh", scriptURL); err != nil {
		return fmt.Errorf("failed to download kubectl-ai install script: %w", err)
	}

	// Make the script executable
	if err := os.Chmod("install_kubectl_ai.sh", 0700); err != nil {
		os.Remove("install_kubectl_ai.sh")
		return fmt.Errorf("failed to make kubectl-ai script executable: %w", err)
	}

	// Execute the script
	fmt.Println("Running kubectl-ai installation script")
	cmd := exec.Command("bash", "install_kubectl_ai.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove("install_kubectl_ai.sh")
		return fmt.Errorf("failed to run kubectl-ai install script: %w", err)
	}

	// Clean up the script
	os.Remove("install_kubectl_ai.sh")

	fmt.Println("✓ kubectl-ai installed successfully")
	return nil
}

// EnsureToolsInstalled checks and installs kubectl, helm, and kubectl-ai if needed
func (t *ToolInstaller) EnsureToolsInstalled() error {
	var errors []string

	if !t.IsKubectlInstalled() {
		if err := t.InstallKubectl(); err != nil {
			errors = append(errors, fmt.Sprintf("kubectl: %v", err))
		}
	} else {
		fmt.Println("✓ kubectl is already installed")
	}

	if !t.IsHelmInstalled() {
		if err := t.InstallHelm(); err != nil {
			errors = append(errors, fmt.Sprintf("helm: %v", err))
		}
	} else {
		fmt.Println("✓ helm is already installed")
	}

	if !t.IsKubectlAIInstalled() {
		if err := t.InstallKubectlAI(); err != nil {
			errors = append(errors, fmt.Sprintf("kubectl-ai: %v", err))
		}
	} else {
		fmt.Println("✓ kubectl-ai is already installed")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to install tools:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// commandExists checks if a command is available in PATH
func (t *ToolInstaller) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// runCommand executes a command and shows output
func (t *ToolInstaller) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Print output if there's any
	if stdout.Len() > 0 {
		fmt.Print(stdout.String())
	}
	if stderr.Len() > 0 {
		fmt.Fprint(os.Stderr, stderr.String())
	}

	return err
}
