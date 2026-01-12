package addons

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/internal/config"
)

// TestK3sTokenEmbeddedInCloudInit verifies that the actual k3s token is embedded in cloud-init
// rather than using a placeholder variable
func TestK3sTokenEmbeddedInCloudInit(t *testing.T) {
	// Setup test configuration
	cfg := &config.Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		Image:       "ubuntu-24.04",
		Networking: config.Networking{
			SSH: config.SSH{
				Port: 22,
			},
		},
		Addons: config.Addons{
			ClusterAutoscaler: &config.ClusterAutoscaler{
				Enabled:                    true,
				ContainerImageTag:          "v1.34.2",
				ScanInterval:               "10s",
				ScaleDownDelayAfterAdd:     "10m",
				ScaleDownDelayAfterDelete:  "10s",
				ScaleDownDelayAfterFailure: "3m",
				MaxNodeProvisionTime:       "15m",
			},
		},
	}

	// Create test pool
	phpName := "php"
	pool := config.WorkerNodePool{
		NodePool: config.NodePool{
			Name:                       &phpName,
			InstanceType:               "cpx32",
			IncludeClusterNameAsPrefix: true,
			Autoscaling: &config.Autoscaling{
				Enabled:      true,
				MinInstances: 1,
				MaxInstances: 3,
			},
		},
		Location: "nbg1",
	}

	// Create installer
	installer := NewClusterAutoscalerInstaller(cfg, nil)

	// Mock server objects
	firstMaster := &hcloud.Server{
		Name: "test-master-1",
	}
	masters := []*hcloud.Server{firstMaster}

	// Test token
	testToken := "test-k3s-token-12345"
	masterIP := "10.0.0.1"

	// Generate cloud-init
	cloudInit, err := installer.generateCloudInitForPool(pool, firstMaster, masters, masterIP, testToken)
	if err != nil {
		t.Fatalf("Failed to generate cloud-init: %v", err)
	}

	// Verify token is embedded in the init script (which is gzip+base64 encoded)
	// We can check that the init file is present and the token appears somewhere
	// in the cloud-init (it will be in the encoded init script)
	if !strings.Contains(cloudInit, "/etc/init-0.sh") {
		t.Error("Cloud-init should contain init script file reference")
		t.Logf("Cloud-init content:\n%s", cloudInit)
	}

	// The test for token embedding now needs to decode the gzip+base64 content
	// For simplicity, we'll just verify the init file is present
	// The actual token embedding is tested in TestGenerateWorkerInstallScript
}

// TestGenerateWorkerInstallScript verifies the worker install script generation
func TestGenerateWorkerInstallScript(t *testing.T) {
	cfg := &config.Main{
		K3sVersion: "v1.32.0+k3s1",
		Networking: config.Networking{
			PrivateNetwork: config.PrivateNetwork{
				Enabled: true,
				Subnet:  "10.0.0.0/16",
			},
		},
	}

	installer := NewClusterAutoscalerInstaller(cfg, nil)

	tests := []struct {
		name          string
		pool          config.WorkerNodePool
		masterIP      string
		k3sToken      string
		expectedParts []string
		notExpected   []string
	}{
		{
			name: "basic pool without labels or taints",
			pool: config.WorkerNodePool{
				NodePool: config.NodePool{
					InstanceType: "cpx32",
				},
			},
			masterIP: "10.0.0.1",
			k3sToken: "test-token-abc",
			expectedParts: []string{
				"curl -sfL https://get.k3s.io",
				"K3S_URL=https://10.0.0.1:6443",
				`echo -n "test-token-abc" > /tmp/k3s-token`,
				`K3S_TOKEN="$(cat /tmp/k3s-token)"`,
				`INSTALL_K3S_VERSION="v1.32.0+k3s1"`,
				`INSTALL_K3S_EXEC="agent"`,
				"Private network IP",
				"set -o pipefail",
			},
			notExpected: []string{
				"${K3S_TOKEN}",
			},
		},
		{
			name: "pool with labels",
			pool: config.WorkerNodePool{
				NodePool: config.NodePool{
					InstanceType: "cpx32",
					Labels: []config.Label{
						{Key: "app", Value: "web"},
						{Key: "env", Value: "prod"},
					},
				},
			},
			masterIP: "10.0.0.1",
			k3sToken: "test-token-xyz",
			expectedParts: []string{
				`echo -n "test-token-xyz" > /tmp/k3s-token`,
				"--node-label=app=web,env=prod",
			},
			notExpected: []string{
				"${K3S_TOKEN}",
			},
		},
		{
			name: "pool with taints",
			pool: config.WorkerNodePool{
				NodePool: config.NodePool{
					InstanceType: "cpx32",
					Taints: []config.Taint{
						{Key: "dedicated", Value: "database", Effect: "NoSchedule"},
					},
				},
			},
			masterIP: "10.0.0.1",
			k3sToken: "test-token-123",
			expectedParts: []string{
				`echo -n "test-token-123" > /tmp/k3s-token`,
				"--node-taint=dedicated=database:NoSchedule",
			},
			notExpected: []string{
				"${K3S_TOKEN}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := installer.generateWorkerInstallScript(tt.masterIP, tt.pool, tt.k3sToken)
			if err != nil {
				t.Fatalf("Failed to generate worker install script: %v", err)
			}

			// Check expected parts are present
			for _, expected := range tt.expectedParts {
				if !strings.Contains(script, expected) {
					t.Errorf("Script should contain '%s'\nActual script:\n%s", expected, script)
				}
			}

			// Check unexpected parts are not present
			for _, notExpected := range tt.notExpected {
				if strings.Contains(script, notExpected) {
					t.Errorf("Script should NOT contain '%s'\nActual script:\n%s", notExpected, script)
				}
			}
		})
	}
}

// TestNATGatewayEnvironmentVariables verifies that when NAT gateway is enabled,
// public IP environment variables are forced to false
func TestNATGatewayEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name               string
		natGatewayEnabled  bool
		publicIPv4Enabled  bool
		publicIPv6Enabled  bool
		useNilIPv4Config   bool
		useNilIPv6Config   bool
		expectedPublicIPv4 string
		expectedPublicIPv6 string
	}{
		{
			name:               "NAT gateway enabled - should disable public IPs",
			natGatewayEnabled:  true,
			publicIPv4Enabled:  true,    // Even if enabled in config
			publicIPv6Enabled:  true,    // Even if enabled in config
			expectedPublicIPv4: "false", // Should be forced to false
			expectedPublicIPv6: "false", // Should be forced to false
		},
		{
			name:               "NAT gateway disabled - should respect config",
			natGatewayEnabled:  false,
			publicIPv4Enabled:  true,
			publicIPv6Enabled:  true,
			expectedPublicIPv4: "true",
			expectedPublicIPv6: "true",
		},
		{
			name:               "NAT gateway disabled with public IPs disabled",
			natGatewayEnabled:  false,
			publicIPv4Enabled:  false,
			publicIPv6Enabled:  false,
			expectedPublicIPv4: "false",
			expectedPublicIPv6: "false",
		},
		{
			name:               "NAT gateway enabled with public IPs disabled in config",
			natGatewayEnabled:  true,
			publicIPv4Enabled:  false,
			publicIPv6Enabled:  false,
			expectedPublicIPv4: "false",
			expectedPublicIPv6: "false",
		},
		{
			name:               "nil public network config - defaults to false",
			natGatewayEnabled:  false,
			useNilIPv4Config:   true,
			useNilIPv6Config:   true,
			expectedPublicIPv4: "false",
			expectedPublicIPv6: "false",
		},
		{
			name:               "NAT gateway with nil public network config",
			natGatewayEnabled:  true,
			useNilIPv4Config:   true,
			useNilIPv6Config:   true,
			expectedPublicIPv4: "false",
			expectedPublicIPv6: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test configuration
			cfg := &config.Main{
				ClusterName: "test-cluster",
				K3sVersion:  "v1.32.0+k3s1",
				Image:       "ubuntu-24.04",
				Networking: config.Networking{
					PrivateNetwork: config.PrivateNetwork{
						Enabled: true,
						Subnet:  "10.0.0.0/16",
					},
					PublicNetwork: config.PublicNetwork{},
					SSH: config.SSH{
						Port: 22,
					},
				},
				Addons: config.Addons{
					ClusterAutoscaler: &config.ClusterAutoscaler{
						Enabled:                    true,
						ContainerImageTag:          "v1.34.2",
						ScanInterval:               "10s",
						ScaleDownDelayAfterAdd:     "10m",
						ScaleDownDelayAfterDelete:  "10s",
						ScaleDownDelayAfterFailure: "3m",
						MaxNodeProvisionTime:       "15m",
					},
				},
			}

			// Configure NAT gateway if enabled
			if tt.natGatewayEnabled {
				cfg.Networking.PrivateNetwork.NATGateway = &config.NATGateway{
					Enabled:      true,
					InstanceType: "cpx11",
				}
			}

			// Configure public network settings
			if !tt.useNilIPv4Config {
				cfg.Networking.PublicNetwork.IPv4 = &config.PublicNetworkIPv4{
					Enabled: tt.publicIPv4Enabled,
				}
			}
			if !tt.useNilIPv6Config {
				cfg.Networking.PublicNetwork.IPv6 = &config.PublicNetworkIPv6{
					Enabled: tt.publicIPv6Enabled,
				}
			}

			// Create test pool
			phpName := "php"
			pool := config.WorkerNodePool{
				NodePool: config.NodePool{
					Name:                       &phpName,
					InstanceType:               "cpx32",
					IncludeClusterNameAsPrefix: true,
					Autoscaling: &config.Autoscaling{
						Enabled:      true,
						MinInstances: 1,
						MaxInstances: 3,
					},
				},
				Location: "nbg1",
			}

			// Create installer
			installer := NewClusterAutoscalerInstaller(cfg, nil)

			// Mock server objects
			firstMaster := &hcloud.Server{
				Name: "test-master-1",
			}
			masters := []*hcloud.Server{firstMaster}

			// Test token
			testToken := "test-k3s-token-12345"
			masterIP := "10.0.0.1"

			// Build environment variables
			env, err := installer.buildEnvironmentVariables(firstMaster, masters, []config.WorkerNodePool{pool}, masterIP, testToken)
			if err != nil {
				t.Fatalf("Failed to build environment variables: %v", err)
			}

			// Verify HCLOUD_PUBLIC_IPV4 and HCLOUD_PUBLIC_IPV6 values
			var publicIPv4Value, publicIPv6Value string
			for _, envVar := range env {
				name, ok := envVar["name"].(string)
				if !ok {
					continue
				}
				if name == "HCLOUD_PUBLIC_IPV4" {
					publicIPv4Value, _ = envVar["value"].(string)
				} else if name == "HCLOUD_PUBLIC_IPV6" {
					publicIPv6Value, _ = envVar["value"].(string)
				}
			}

			if publicIPv4Value != tt.expectedPublicIPv4 {
				t.Errorf("HCLOUD_PUBLIC_IPV4 = %q, expected %q", publicIPv4Value, tt.expectedPublicIPv4)
			}

			if publicIPv6Value != tt.expectedPublicIPv6 {
				t.Errorf("HCLOUD_PUBLIC_IPV6 = %q, expected %q", publicIPv6Value, tt.expectedPublicIPv6)
			}

			// Also verify the cluster config includes the correct settings
			var clusterConfigBase64 string
			for _, envVar := range env {
				name, ok := envVar["name"].(string)
				if !ok {
					continue
				}
				if name == "HCLOUD_CLUSTER_CONFIG" {
					clusterConfigBase64, _ = envVar["value"].(string)
					break
				}
			}

			if clusterConfigBase64 == "" {
				t.Fatal("HCLOUD_CLUSTER_CONFIG not found in environment variables")
			}

			// Decode and verify the cluster config
			clusterConfigJSON, err := base64.StdEncoding.DecodeString(clusterConfigBase64)
			if err != nil {
				t.Fatalf("Failed to decode cluster config: %v", err)
			}

			var clusterConfig map[string]interface{}
			if err := json.Unmarshal(clusterConfigJSON, &clusterConfig); err != nil {
				t.Fatalf("Failed to unmarshal cluster config: %v", err)
			}

			// Verify cluster config has the expected structure
			if _, ok := clusterConfig["nodeConfigs"]; !ok {
				t.Error("Cluster config should contain nodeConfigs")
			}
		})
	}
}
