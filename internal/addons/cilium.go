package addons

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/internal/config"
	"github.com/magenx/hek3ster/internal/util"
)

//go:embed templates/cilium_values.yaml.tmpl
var ciliumValuesTemplate string

// CiliumInstaller handles Cilium CNI installation
type CiliumInstaller struct {
	Config    *config.Main
	SSHClient *util.SSH
	ctx       context.Context
}

// NewCiliumInstaller creates a new Cilium installer
func NewCiliumInstaller(cfg *config.Main, sshClient *util.SSH) *CiliumInstaller {
	return &CiliumInstaller{
		Config:    cfg,
		SSHClient: sshClient,
		ctx:       context.Background(),
	}
}

// Install installs Cilium CNI via Helm
func (c *CiliumInstaller) Install(firstMaster *hcloud.Server, masterSSHIP string) error {
	if c.Config.Networking.CNI.Mode != "cilium" {
		return nil // Not using Cilium, skip installation
	}

	if c.Config.Networking.CNI.Cilium == nil {
		return fmt.Errorf("Cilium configuration is missing")
	}

	// Check if Cilium is already installed by checking helm status
	checkCmd := "helm status cilium -n kube-system > /dev/null 2>&1"
	_, err := c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, checkCmd, c.Config.Networking.SSH.UseAgent)
	if err == nil {
		// Cilium is already installed
		util.LogInfo("Cilium CNI already installed, skipping installation", "addons")
		return nil
	}

	util.LogInfo("Installing Cilium CNI", "cilium")

	// Step 1: Add Cilium Helm repository
	if err := c.addHelmRepo(masterSSHIP); err != nil {
		return fmt.Errorf("failed to add Cilium Helm repository: %w", err)
	}

	// Step 2: Install Cilium using Helm
	if err := c.installCiliumHelm(masterSSHIP); err != nil {
		return fmt.Errorf("failed to install Cilium: %w", err)
	}

	// Step 3: Wait for Cilium to be ready
	if err := c.waitForCiliumReady(masterSSHIP); err != nil {
		return fmt.Errorf("failed to verify Cilium status: %w", err)
	}

	util.LogSuccess("Cilium CNI installed successfully", "cilium")
	return nil
}

// addHelmRepo adds the Cilium Helm repository
func (c *CiliumInstaller) addHelmRepo(masterSSHIP string) error {
	cmd := "helm repo add cilium https://helm.cilium.io/"
	_, err := c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, cmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		// Repo might already exist, try updating
		updateCmd := "helm repo update"
		_, updateErr := c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, updateCmd, c.Config.Networking.SSH.UseAgent)
		if updateErr != nil {
			return fmt.Errorf("failed to add or update Helm repository: %w", err)
		}
	}
	return nil
}

// installCiliumHelm installs Cilium using Helm
func (c *CiliumInstaller) installCiliumHelm(masterSSHIP string) error {
	ciliumConfig := c.Config.Networking.CNI.Cilium

	// Check if custom Helm values file is provided
	if ciliumConfig.HelmValuesPath != "" {
		// Use custom values file
		return c.installWithCustomValues(masterSSHIP, ciliumConfig.HelmValuesPath)
	}

	// Generate values from template
	return c.installWithGeneratedValues(masterSSHIP)
}

// installWithCustomValues installs Cilium with a custom values file
func (c *CiliumInstaller) installWithCustomValues(masterSSHIP, valuesPath string) error {
	ciliumConfig := c.Config.Networking.CNI.Cilium

	// Read custom values file
	if !util.FileExists(valuesPath) {
		return fmt.Errorf("custom Helm values file not found: %s", valuesPath)
	}

	valuesContent, err := os.ReadFile(valuesPath)
	if err != nil {
		return fmt.Errorf("failed to read custom values file: %w", err)
	}

	// Sanitize version to prevent command injection
	version := sanitizeVersion(ciliumConfig.HelmChartVersion)

	// Upload values file to master node
	remoteValuesPath := "/tmp/cilium_values.yaml"
	uploadCmd := fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", remoteValuesPath, string(valuesContent))
	_, err = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, uploadCmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		return fmt.Errorf("failed to upload values file: %w", err)
	}

	// Install Cilium with custom values - remote path is safe since it's hardcoded
	installCmd := fmt.Sprintf(
		"helm upgrade --install --version %s --namespace kube-system --values %s cilium cilium/cilium",
		version,
		remoteValuesPath,
	)

	_, err = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, installCmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		return fmt.Errorf("helm install failed: %w", err)
	}

	// Clean up temporary file
	cleanupCmd := fmt.Sprintf("rm -f %s", remoteValuesPath)
	_, _ = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, cleanupCmd, c.Config.Networking.SSH.UseAgent)

	return nil
}

// sanitizeVersion sanitizes version string to prevent command injection
func sanitizeVersion(version string) string {
	// Remove any potentially dangerous characters, allow only alphanumeric, dots, and dashes
	version = strings.TrimSpace(version)
	var result strings.Builder
	for _, r := range version {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == 'v' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// installWithGeneratedValues installs Cilium with generated values from template
func (c *CiliumInstaller) installWithGeneratedValues(masterSSHIP string) error {
	ciliumConfig := c.Config.Networking.CNI.Cilium

	// Generate Helm values from template
	valuesContent, err := c.generateHelmValues()
	if err != nil {
		return fmt.Errorf("failed to generate Helm values: %w", err)
	}

	// Upload generated values to master node
	remoteValuesPath := "/tmp/cilium_values.yaml"
	uploadCmd := fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", remoteValuesPath, valuesContent)
	_, err = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, uploadCmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		return fmt.Errorf("failed to upload generated values: %w", err)
	}

	// Sanitize version to prevent command injection
	version := sanitizeVersion(ciliumConfig.HelmChartVersion)

	// Install Cilium with generated values - remote path is safe since it's hardcoded
	installCmd := fmt.Sprintf(
		"helm upgrade --install --version %s --namespace kube-system --values %s cilium cilium/cilium",
		version,
		remoteValuesPath,
	)

	_, err = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, installCmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		return fmt.Errorf("helm install failed: %w", err)
	}

	// Clean up temporary file
	cleanupCmd := fmt.Sprintf("rm -f %s", remoteValuesPath)
	_, _ = c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, cleanupCmd, c.Config.Networking.SSH.UseAgent)

	return nil
}

// generateHelmValues generates Helm values from template
func (c *CiliumInstaller) generateHelmValues() (string, error) {
	ciliumConfig := c.Config.Networking.CNI.Cilium

	// Ensure defaults are set
	ciliumConfig.SetDefaults()

	// Determine encryption enabled - Cilium has its own encryption settings
	// Enable encryption if EncryptionType is configured
	encryptionEnabled := ciliumConfig.EncryptionType != ""

	// Build Hubble metrics array
	hubbleMetrics := c.buildHubbleMetrics(ciliumConfig.HubbleMetrics)

	// Safe dereferencing of pointers with defaults
	hubbleEnabled := ciliumConfig.HubbleEnabled != nil && *ciliumConfig.HubbleEnabled
	hubbleRelayEnabled := ciliumConfig.HubbleRelayEnabled != nil && *ciliumConfig.HubbleRelayEnabled
	hubbleUIEnabled := ciliumConfig.HubbleUIEnabled != nil && *ciliumConfig.HubbleUIEnabled

	// Prepare template data
	data := map[string]interface{}{
		"encryption_enabled":      encryptionEnabled,
		"encryption_type":         ciliumConfig.EncryptionType,
		"routing_mode":            ciliumConfig.RoutingMode,
		"tunnel_protocol":         ciliumConfig.TunnelProtocol,
		"hubble_enabled":          hubbleEnabled,
		"hubble_metrics_enabled":  len(hubbleMetrics) > 0,
		"hubble_metrics":          hubbleMetrics,
		"hubble_relay_enabled":    hubbleRelayEnabled,
		"hubble_ui_enabled":       hubbleUIEnabled,
		"k8s_service_host":        ciliumConfig.K8sServiceHost,
		"k8s_service_port":        ciliumConfig.K8sServicePort,
		"operator_replicas":       ciliumConfig.OperatorReplicas,
		"operator_memory_request": ciliumConfig.OperatorMemoryRequest,
		"agent_memory_request":    ciliumConfig.AgentMemoryRequest,
		"egress_gateway_enabled":  ciliumConfig.EgressGatewayEnabled,
	}

	// Parse and execute template
	tmpl, err := template.New("cilium_values").Parse(ciliumValuesTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// buildHubbleMetrics builds Hubble metrics array from slice
func (c *CiliumInstaller) buildHubbleMetrics(customMetrics []string) string {
	var metrics []string
	if len(customMetrics) > 0 {
		// Use custom metrics
		metrics = customMetrics
	} else {
		// Default metrics
		metrics = []string{"dns", "drop", "tcp", "flow", "port-distribution", "icmp", "http"}
	}
	
	// Serialize to JSON format
	jsonBytes, err := json.Marshal(metrics)
	if err != nil {
		// Fallback to default if serialization fails
		return `["dns", "drop", "tcp", "flow", "port-distribution", "icmp", "http"]`
	}
	return string(jsonBytes)
}

// waitForCiliumReady waits for Cilium to be ready
func (c *CiliumInstaller) waitForCiliumReady(masterSSHIP string) error {
	checkCmd := "kubectl -n kube-system rollout status ds cilium --timeout=300s"
	_, err := c.SSHClient.Run(c.ctx, masterSSHIP, c.Config.Networking.SSH.Port, checkCmd, c.Config.Networking.SSH.UseAgent)
	if err != nil {
		return fmt.Errorf("Cilium rollout check failed: %w", err)
	}
	return nil
}
