package cloudinit

import (
	"strings"
	"testing"

	"github.com/magenx/hek3ster/internal/config"
)

// TestConfigMergingLogic tests that pool-specific settings override root settings
func TestConfigMergingLogic(t *testing.T) {
	// Root config
	rootConfig := &config.Main{
		AdditionalPackages:        []string{"htop", "vim"},
		AdditionalPreK3sCommands:  []string{"apt update"},
		AdditionalPostK3sCommands: []string{"apt autoremove -y"},
	}

	// Pool with overrides
	poolWithOverrides := &config.NodePool{
		AdditionalPackages:        []string{"curl"},
		AdditionalPreK3sCommands:  []string{"echo 'pool pre'"},
		AdditionalPostK3sCommands: []string{"echo 'pool post'"},
	}

	// Pool without overrides (should use root)
	poolWithoutOverrides := &config.NodePool{
		AdditionalPackages:        nil,
		AdditionalPreK3sCommands:  nil,
		AdditionalPostK3sCommands: nil,
	}

	t.Run("pool with overrides", func(t *testing.T) {
		packages := rootConfig.AdditionalPackages
		if poolWithOverrides.AdditionalPackages != nil {
			packages = poolWithOverrides.AdditionalPackages
		}

		preCommands := rootConfig.AdditionalPreK3sCommands
		if poolWithOverrides.AdditionalPreK3sCommands != nil {
			preCommands = poolWithOverrides.AdditionalPreK3sCommands
		}

		postCommands := rootConfig.AdditionalPostK3sCommands
		if poolWithOverrides.AdditionalPostK3sCommands != nil {
			postCommands = poolWithOverrides.AdditionalPostK3sCommands
		}

		// Verify pool overrides are used
		if len(packages) != 1 || packages[0] != "curl" {
			t.Errorf("Expected pool packages [curl], got %v", packages)
		}
		if len(preCommands) != 1 || preCommands[0] != "echo 'pool pre'" {
			t.Errorf("Expected pool pre commands [echo 'pool pre'], got %v", preCommands)
		}
		if len(postCommands) != 1 || postCommands[0] != "echo 'pool post'" {
			t.Errorf("Expected pool post commands [echo 'pool post'], got %v", postCommands)
		}

		// Generate cloud-init with pool overrides
		cfg := &Config{
			SSHPort:                   22,
			Packages:                  packages,
			AdditionalPreK3sCommands:  preCommands,
			AdditionalPostK3sCommands: postCommands,
		}
		generator := NewGenerator(cfg)
		cloudInit, err := generator.Generate()
		if err != nil {
			t.Fatalf("Failed to generate cloud-init: %v", err)
		}

		// Verify pool-specific settings are in the cloud-init
		if !strings.Contains(cloudInit, "'curl'") {
			t.Error("Cloud-init should contain pool-specific package 'curl'")
		}
		if !strings.Contains(cloudInit, "echo 'pool pre'") {
			t.Error("Cloud-init should contain pool-specific pre command")
		}
		if !strings.Contains(cloudInit, "echo 'pool post'") {
			t.Error("Cloud-init should contain pool-specific post command")
		}
		// Verify root settings are NOT in the cloud-init
		if strings.Contains(cloudInit, "'htop'") {
			t.Error("Cloud-init should not contain root package 'htop' when pool overrides")
		}
	})

	t.Run("pool without overrides uses root settings", func(t *testing.T) {
		packages := rootConfig.AdditionalPackages
		if poolWithoutOverrides.AdditionalPackages != nil {
			packages = poolWithoutOverrides.AdditionalPackages
		}

		preCommands := rootConfig.AdditionalPreK3sCommands
		if poolWithoutOverrides.AdditionalPreK3sCommands != nil {
			preCommands = poolWithoutOverrides.AdditionalPreK3sCommands
		}

		postCommands := rootConfig.AdditionalPostK3sCommands
		if poolWithoutOverrides.AdditionalPostK3sCommands != nil {
			postCommands = poolWithoutOverrides.AdditionalPostK3sCommands
		}

		// Verify root settings are used
		if len(packages) != 2 || packages[0] != "htop" || packages[1] != "vim" {
			t.Errorf("Expected root packages [htop vim], got %v", packages)
		}
		if len(preCommands) != 1 || preCommands[0] != "apt update" {
			t.Errorf("Expected root pre commands [apt update], got %v", preCommands)
		}
		if len(postCommands) != 1 || postCommands[0] != "apt autoremove -y" {
			t.Errorf("Expected root post commands [apt autoremove -y], got %v", postCommands)
		}

		// Generate cloud-init with root settings
		cfg := &Config{
			SSHPort:                   22,
			Packages:                  packages,
			AdditionalPreK3sCommands:  preCommands,
			AdditionalPostK3sCommands: postCommands,
		}
		generator := NewGenerator(cfg)
		cloudInit, err := generator.Generate()
		if err != nil {
			t.Fatalf("Failed to generate cloud-init: %v", err)
		}

		// Verify root settings are in the cloud-init
		if !strings.Contains(cloudInit, "'htop'") {
			t.Error("Cloud-init should contain root package 'htop'")
		}
		if !strings.Contains(cloudInit, "'vim'") {
			t.Error("Cloud-init should contain root package 'vim'")
		}
		if !strings.Contains(cloudInit, "apt update") {
			t.Error("Cloud-init should contain root pre command 'apt update'")
		}
		if !strings.Contains(cloudInit, "apt autoremove -y") {
			t.Error("Cloud-init should contain root post command 'apt autoremove -y'")
		}
	})
}
