package cluster

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/internal/config"
)

// GetServerIP returns the appropriate IP address for a server based on network configuration
// It prefers private IP if private networking is enabled and available, otherwise uses public IPv4
func GetServerIP(server *hcloud.Server, cfg *config.Main) (string, error) {
	// Prefer private IP if private networking is enabled and available
	if cfg.Networking.PrivateNetwork.Enabled && len(server.PrivateNet) > 0 {
		return server.PrivateNet[0].IP.String(), nil
	}

	// Fall back to public IPv4
	if server.PublicNet.IPv4.IP != nil {
		return server.PublicNet.IPv4.IP.String(), nil
	}

	return "", fmt.Errorf("server %s has no accessible IP address (private networking disabled or unavailable, and no public IPv4)", server.Name)
}

// GetServerPublicIP returns the public IPv4 address of a server for external access (e.g., kubeconfig)
func GetServerPublicIP(server *hcloud.Server) (string, error) {
	if server.PublicNet.IPv4.IP != nil {
		return server.PublicNet.IPv4.IP.String(), nil
	}
	return "", fmt.Errorf("server %s has no public IPv4 address", server.Name)
}

// GetServerSSHIP returns the appropriate IP address for SSH connections from external machines
// It always prefers public IP for SSH access, as the control machine may not have access to private IPs.
// The fallback to private IP is only for edge cases where servers are created without public IPs
// (e.g., when public_network.ipv4 is explicitly disabled in the configuration).
func GetServerSSHIP(server *hcloud.Server) (string, error) {
	// Always prefer public IP for SSH from external machines
	if server.PublicNet.IPv4.IP != nil {
		return server.PublicNet.IPv4.IP.String(), nil
	}

	// Fall back to private IP only if no public IP is available
	// Note: SSH will likely fail from external machines in this case unless VPN/bastion is used
	if len(server.PrivateNet) > 0 {
		return server.PrivateNet[0].IP.String(), nil
	}

	return "", fmt.Errorf("server %s has no accessible IP address for SSH", server.Name)
}

// GenerateTLSSans generates TLS SAN (Subject Alternative Name) flags for k3s installation
// This ensures the k3s API server certificate includes all necessary IP addresses and hostnames
func GenerateTLSSans(cfg *config.Main, masters []*hcloud.Server, firstMaster *hcloud.Server, apiLoadBalancer *hcloud.LoadBalancer) (string, error) {
	// Use a map to collect unique SANs while building the list
	uniqueSans := make(map[string]bool)

	// Add first master's API server IP (private IP if available, otherwise public)
	apiServerIP, err := GetServerIP(firstMaster, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get API server IP: %w", err)
	}
	uniqueSans[fmt.Sprintf("--tls-san=%s", apiServerIP)] = true

	// Always add localhost
	uniqueSans["--tls-san=127.0.0.1"] = true

	// Add API load balancer IP if configured and created
	if apiLoadBalancer != nil && apiLoadBalancer.PublicNet.IPv4.IP != nil {
		lbIP := apiLoadBalancer.PublicNet.IPv4.IP.String()
		uniqueSans[fmt.Sprintf("--tls-san=%s", lbIP)] = true
	}

	// Add API server hostname if configured
	if cfg.APIServerHostname != "" {
		uniqueSans[fmt.Sprintf("--tls-san=%s", cfg.APIServerHostname)] = true
	}

	// Add all master IPs (both private and public)
	for _, master := range masters {
		// Add private IP
		if len(master.PrivateNet) > 0 {
			privateIP := master.PrivateNet[0].IP.String()
			uniqueSans[fmt.Sprintf("--tls-san=%s", privateIP)] = true
		}

		// Add public IP
		if master.PublicNet.IPv4.IP != nil {
			publicIP := master.PublicNet.IPv4.IP.String()
			uniqueSans[fmt.Sprintf("--tls-san=%s", publicIP)] = true
		}
	}

	// Convert map to sorted slice
	sortedSans := make([]string, 0, len(uniqueSans))
	for san := range uniqueSans {
		sortedSans = append(sortedSans, san)
	}

	// Sort for deterministic output
	sort.Strings(sortedSans)

	return strings.Join(sortedSans, " "), nil
}
