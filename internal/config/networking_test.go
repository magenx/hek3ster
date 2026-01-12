package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestPublicNetworkIPv4_UnmarshalYAML_BooleanFormat(t *testing.T) {
	yamlData := `
public_network:
  ipv4: true
  ipv6: false
`
	var config struct {
		PublicNetwork PublicNetwork `yaml:"public_network"`
	}

	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML with boolean format: %v", err)
	}

	if config.PublicNetwork.IPv4 == nil {
		t.Fatal("Expected IPv4 to be set, got nil")
	}
	if !config.PublicNetwork.IPv4.Enabled {
		t.Error("Expected IPv4.Enabled to be true, got false")
	}

	if config.PublicNetwork.IPv6 == nil {
		t.Fatal("Expected IPv6 to be set, got nil")
	}
	if config.PublicNetwork.IPv6.Enabled {
		t.Error("Expected IPv6.Enabled to be false, got true")
	}
}

func TestPublicNetworkIPv4_UnmarshalYAML_ObjectFormat(t *testing.T) {
	yamlData := `
public_network:
  ipv4:
    enabled: true
  ipv6:
    enabled: false
`
	var config struct {
		PublicNetwork PublicNetwork `yaml:"public_network"`
	}

	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML with object format: %v", err)
	}

	if config.PublicNetwork.IPv4 == nil {
		t.Fatal("Expected IPv4 to be set, got nil")
	}
	if !config.PublicNetwork.IPv4.Enabled {
		t.Error("Expected IPv4.Enabled to be true, got false")
	}

	if config.PublicNetwork.IPv6 == nil {
		t.Fatal("Expected IPv6 to be set, got nil")
	}
	if config.PublicNetwork.IPv6.Enabled {
		t.Error("Expected IPv6.Enabled to be false, got true")
	}
}

func TestPublicNetworkIPv4_UnmarshalYAML_MixedFormats(t *testing.T) {
	// Test that we can have one as boolean and one as object
	yamlData := `
public_network:
  ipv4: true
  ipv6:
    enabled: false
`
	var config struct {
		PublicNetwork PublicNetwork `yaml:"public_network"`
	}

	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML with mixed formats: %v", err)
	}

	if config.PublicNetwork.IPv4 == nil {
		t.Fatal("Expected IPv4 to be set, got nil")
	}
	if !config.PublicNetwork.IPv4.Enabled {
		t.Error("Expected IPv4.Enabled to be true, got false")
	}

	if config.PublicNetwork.IPv6 == nil {
		t.Fatal("Expected IPv6 to be set, got nil")
	}
	if config.PublicNetwork.IPv6.Enabled {
		t.Error("Expected IPv6.Enabled to be false, got true")
	}
}

func TestPublicNetwork_UseLocalFirewall(t *testing.T) {
	yamlData := `
public_network:
  ipv4: true
  ipv6: false
  use_local_firewall: true
  hetzner_ips_query_server_url: "http://example.com/ips"
`
	type TestConfig struct {
		PublicNetwork PublicNetwork `yaml:"public_network"`
	}

	var config TestConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if config.PublicNetwork.IPv4 == nil {
		t.Fatal("IPv4 should not be nil")
	}
	if !config.PublicNetwork.IPv4.Enabled {
		t.Error("IPv4 should be enabled")
	}

	if config.PublicNetwork.IPv6 == nil {
		t.Fatal("IPv6 should not be nil")
	}
	if config.PublicNetwork.IPv6.Enabled {
		t.Error("IPv6 should be disabled")
	}

	if !config.PublicNetwork.UseLocalFirewall {
		t.Error("UseLocalFirewall should be true")
	}

	if config.PublicNetwork.HetznerIPsQueryServerURL != "http://example.com/ips" {
		t.Errorf("Expected HetznerIPsQueryServerURL to be 'http://example.com/ips', got '%s'", config.PublicNetwork.HetznerIPsQueryServerURL)
	}
}
