package config

import (
	"os"
	"strings"
	"testing"
)

func TestValidateWorkerPools_WithAutoscaling(t *testing.T) {
	// Test case: Worker pool with autoscaling enabled should NOT require instance_count
	phpPoolName := "php"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:          &phpPoolName,
					InstanceType:  "cpx32",
					InstanceCount: 0, // This should be ignored when autoscaling is enabled
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: 1,
						MaxInstances: 3,
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should not have any errors about instance_count
	for _, err := range validator.GetErrors() {
		if strings.Contains(err, "instance_count") {
			t.Errorf("Expected no instance_count validation error when autoscaling is enabled, got: %s", err)
		}
	}
}

func TestValidateWorkerPools_WithoutAutoscaling(t *testing.T) {
	// Test case: Worker pool without autoscaling must have instance_count >= 1
	mariadbPoolName := "mariadb"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:          &mariadbPoolName,
					InstanceType:  "cpx32",
					InstanceCount: 0, // This should trigger an error
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should have an error about instance_count
	foundError := false
	for _, err := range validator.GetErrors() {
		if strings.Contains(err, "instance_count must be at least 1") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Error("Expected instance_count validation error when autoscaling is not enabled")
	}
}

func TestValidateWorkerPools_AutoscalingMinInstances(t *testing.T) {
	// Test case: Autoscaling min_instances can be 0
	poolName := "test"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:         &poolName,
					InstanceType: "cpx32",
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: 0, // Valid: can start with 0 nodes and scale up
						MaxInstances: 3,
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should NOT have any errors - min_instances: 0 is valid
	if len(validator.GetErrors()) > 0 {
		t.Errorf("Expected no validation errors for min_instances: 0, got: %v", validator.GetErrors())
	}
}

func TestValidateWorkerPools_AutoscalingMaxGreaterThanMin(t *testing.T) {
	// Test case: Autoscaling max_instances must be greater than min_instances
	poolName := "test"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:         &poolName,
					InstanceType: "cpx32",
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: 3,
						MaxInstances: 3, // Invalid: should be greater than min
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should have an error about max_instances
	foundError := false
	for _, err := range validator.GetErrors() {
		if strings.Contains(err, "max_instances must be greater than min_instances") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Error("Expected max_instances validation error")
	}
}

func TestValidateWorkerPools_AutoscalingNegativeMinInstances(t *testing.T) {
	// Test case: Autoscaling min_instances cannot be negative
	poolName := "test"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:         &poolName,
					InstanceType: "cpx32",
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: -1, // Invalid: cannot be negative
						MaxInstances: 3,
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should have an error about negative min_instances
	foundError := false
	for _, err := range validator.GetErrors() {
		if strings.Contains(err, "min_instances cannot be negative") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Error("Expected min_instances negative validation error")
	}
}

func TestValidateWorkerPools_MixedAutoscalingAndStatic(t *testing.T) {
	// Test case: Mix of autoscaling and static pools
	mariadbPoolName := "mariadb"
	phpPoolName := "php"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:          &mariadbPoolName,
					InstanceType:  "cpx32",
					InstanceCount: 1, // Static pool with explicit count
				},
				Location: "nbg1",
			},
			{
				NodePool: NodePool{
					Name:         &phpPoolName,
					InstanceType: "cpx32",
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: 1,
						MaxInstances: 3,
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should have no errors
	if len(validator.GetErrors()) > 0 {
		t.Errorf("Expected no errors for valid mixed configuration, got: %v", validator.GetErrors())
	}
}

func TestValidateWorkerPools_ValidAutoscaling(t *testing.T) {
	// Test case: Valid autoscaling configuration
	poolName := "php"
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		WorkerNodePools: []WorkerNodePool{
			{
				NodePool: NodePool{
					Name:         &poolName,
					InstanceType: "cpx32",
					Autoscaling: &Autoscaling{
						Enabled:      true,
						MinInstances: 1,
						MaxInstances: 5,
					},
				},
				Location: "nbg1",
			},
		},
	}

	validator := NewValidator(config)
	validator.validateWorkerPools()

	// Should have no errors
	if len(validator.GetErrors()) > 0 {
		t.Errorf("Expected no errors for valid autoscaling configuration, got: %v", validator.GetErrors())
	}
}

func TestValidateSSHKeys_TildeExpansion(t *testing.T) {
	// Test case: SSH key paths with tilde should be properly expanded before validation
	// Create temporary SSH keys for testing
	tmpDir := t.TempDir()
	privateKeyPath := tmpDir + "/id_test"
	publicKeyPath := tmpDir + "/id_test.pub"

	// Create dummy SSH key files
	if err := createDummyFile(privateKeyPath); err != nil {
		t.Fatalf("Failed to create temporary private key: %v", err)
	}
	if err := createDummyFile(publicKeyPath); err != nil {
		t.Fatalf("Failed to create temporary public key: %v", err)
	}

	// Test with absolute paths first - should work
	config := &Main{
		ClusterName: "test-cluster",
		K3sVersion:  "v1.32.0+k3s1",
		Networking: Networking{
			SSH: SSH{
				PrivateKeyPath: privateKeyPath,
				PublicKeyPath:  publicKeyPath,
				Port:           22,
			},
		},
		MastersPool: MasterNodePool{
			NodePool: NodePool{
				InstanceType:  "cpx11",
				InstanceCount: 1,
			},
			Locations: []string{"nbg1"},
		},
	}

	validator := NewValidator(config)
	validator.validateSSHKeys()

	// Should have no errors for valid absolute paths
	sshErrors := filterSSHKeyErrors(validator.GetErrors())
	if len(sshErrors) > 0 {
		t.Errorf("Expected no SSH key validation errors for absolute paths, got: %v", sshErrors)
	}
}

// Helper function to create a dummy file for testing
func createDummyFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("dummy content")
	return err
}

// Helper function to filter SSH key related errors
func filterSSHKeyErrors(errors []string) []string {
	var sshErrors []string
	for _, err := range errors {
		if strings.Contains(err, "SSH") && (strings.Contains(err, "key") || strings.Contains(err, "path")) {
			sshErrors = append(sshErrors, err)
		}
	}
	return sshErrors
}
