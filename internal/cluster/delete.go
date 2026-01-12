package cluster

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/internal/config"
	"github.com/magenx/hek3ster/internal/util"
	"github.com/magenx/hek3ster/pkg/hetzner"
)

const (
	// HCloudNodeGroupLabel is the label key used by the cluster autoscaler to identify node groups
	HCloudNodeGroupLabel = "hcloud/node-group"
)

// Deleter handles cluster deletion
type Deleter struct {
	Config        *config.Main
	HetznerClient *hetzner.Client
	Force         bool
	ctx           context.Context
}

// NewDeleter creates a new cluster deleter
func NewDeleter(cfg *config.Main, hetznerClient *hetzner.Client, force bool) *Deleter {
	return &Deleter{
		Config:        cfg,
		HetznerClient: hetznerClient,
		Force:         force,
		ctx:           context.Background(),
	}
}

// Run executes the cluster deletion process
func (d *Deleter) Run() error {
	util.LogInfo("Starting cluster deletion", d.Config.ClusterName)

	// Confirm deletion if not forced
	if !d.Force {
		// Request cluster name confirmation
		if err := d.requestClusterNameConfirmation(); err != nil {
			return err
		}

		// Check protection against deletion
		if d.Config.ProtectAgainstDeletion {
			util.LogError("Cluster cannot be deleted. If you are sure about this, disable the protection by setting `protect_against_deletion` to `false` in the config file. Aborting deletion.", "")
			return fmt.Errorf("cluster is protected against deletion")
		}
	}

	// Track errors during deletion
	var deletionErrors []string

	// Find all resources with cluster label
	clusterLabel := fmt.Sprintf("cluster=%s", d.Config.ClusterName)

	// Step 1: Delete servers first
	spinner := util.NewSpinner("Finding and deleting servers", "servers")
	spinner.Start()
	servers, err := d.HetznerClient.ListServers(d.ctx, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: clusterLabel,
		},
	})
	if err != nil {
		spinner.Stop(true)
		return fmt.Errorf("failed to list servers: %w", err)
	}

	// Also find servers from autoscaling-enabled worker node pools
	// These are created by the cluster autoscaler and have the HCloudNodeGroupLabel
	autoscaledServers, err := d.findAutoscaledPoolServers()
	if err != nil {
		spinner.Stop(true)
		return fmt.Errorf("failed to find autoscaled pool servers: %w", err)
	}

	// Merge servers, avoiding duplicates
	serverMap := make(map[int64]*hcloud.Server)
	for _, server := range servers {
		serverMap[server.ID] = server
	}
	for _, server := range autoscaledServers {
		serverMap[server.ID] = server
	}

	// Convert map back to slice
	allServers := make([]*hcloud.Server, 0, len(serverMap))
	for _, server := range serverMap {
		allServers = append(allServers, server)
	}

	spinner.Stop(true)

	// Log found servers
	if len(allServers) == 0 {
		util.LogInfo("No servers found", "servers")
	} else {
		util.LogInfo(fmt.Sprintf("Found %d server(s):", len(allServers)), "servers")
		for _, server := range allServers {
			fmt.Printf("  - %s\n", server.Name)
		}

		// Delete servers with individual logging
		for _, server := range allServers {
			util.LogInfo(fmt.Sprintf("Deleting server: %s", server.Name), "servers")
			if err := d.HetznerClient.DeleteServer(d.ctx, server); err != nil {
				errMsg := fmt.Sprintf("Failed to delete server %s: %v", server.Name, err)
				util.LogError(errMsg, "servers")
				deletionErrors = append(deletionErrors, errMsg)
			} else {
				util.LogSuccess(fmt.Sprintf("Deleted server: %s", server.Name), "servers")
			}
		}
	}

	util.LogSuccess(fmt.Sprintf("Completed deletion of %d server(s)", len(allServers)), "servers")

	// Step 2: Delete load balancers
	util.LogInfo("Finding and deleting load balancers", "load balancer")

	// Delete API load balancer
	apiLbName := fmt.Sprintf("%s-api-lb", d.Config.ClusterName)
	apiLb, err := d.HetznerClient.GetLoadBalancer(d.ctx, apiLbName)
	if err == nil && apiLb != nil {
		if err := d.HetznerClient.DeleteLoadBalancer(d.ctx, apiLb); err != nil {
			errMsg := fmt.Sprintf("Failed to delete API load balancer: %v", err)
			util.LogError(errMsg, "load balancer")
			deletionErrors = append(deletionErrors, errMsg)
		} else {
			util.LogSuccess("API load balancer deleted", "load balancer")
		}
	}

	// Delete global load balancer (if it was created)
	globalLbName := fmt.Sprintf("%s-global-lb", d.Config.ClusterName)
	if d.Config.LoadBalancer.Name != nil {
		globalLbName = *d.Config.LoadBalancer.Name
	}
	globalLb, err := d.HetznerClient.GetLoadBalancer(d.ctx, globalLbName)
	if err == nil && globalLb != nil {
		if err := d.HetznerClient.DeleteLoadBalancer(d.ctx, globalLb); err != nil {
			errMsg := fmt.Sprintf("Failed to delete global load balancer: %v", err)
			util.LogError(errMsg, "load balancer")
			deletionErrors = append(deletionErrors, errMsg)
		} else {
			util.LogSuccess("Global load balancer deleted", "load balancer")
		}
	}

	// Step 3: Delete network (after servers and load balancers that might be using it)
	if d.Config.Networking.PrivateNetwork.Enabled {
		util.LogInfo("Finding and deleting network", "network")
		// Use cluster name as network name (matching creation logic)
		networkName := d.Config.ClusterName
		network, err := d.HetznerClient.GetNetwork(d.ctx, networkName)
		if err == nil && network != nil {
			if err := d.HetznerClient.DeleteNetwork(d.ctx, network); err != nil {
				errMsg := fmt.Sprintf("Failed to delete network: %v", err)
				util.LogError(errMsg, "network")
				deletionErrors = append(deletionErrors, errMsg)
			} else {
				util.LogSuccess("Network deleted", "network")
			}
		}
	}

	// Step 4: Delete firewalls (after all other resources, before SSH key)
	// Skip if local firewall is enabled
	if !d.Config.Networking.PrivateNetwork.Enabled && d.Config.Networking.PublicNetwork.UseLocalFirewall {
		util.LogInfo("Local firewall was used, skipping Hetzner cloud firewall deletion", "firewall")
	} else {
		util.LogInfo("Finding and deleting firewalls", "firewall")
		firewallName := fmt.Sprintf("%s-firewall", d.Config.ClusterName)
		firewall, err := d.HetznerClient.GetFirewall(d.ctx, firewallName)
		if err == nil && firewall != nil {
			if err := d.HetznerClient.DeleteFirewall(d.ctx, firewall); err != nil {
				errMsg := fmt.Sprintf("Failed to delete firewall: %v", err)
				util.LogError(errMsg, "firewall")
				deletionErrors = append(deletionErrors, errMsg)
			} else {
				util.LogSuccess("Firewall deleted", "firewall")
			}
		}
	}

	// Step 5: Delete SSH key (last, no dependencies)
	util.LogInfo("Finding and deleting SSH key", "ssh key")
	keyName := fmt.Sprintf("%s-ssh-key", d.Config.ClusterName)
	sshKey, err := d.HetznerClient.GetSSHKey(d.ctx, keyName)
	if err == nil && sshKey != nil {
		if err := d.HetznerClient.DeleteSSHKey(d.ctx, sshKey); err != nil {
			errMsg := fmt.Sprintf("Failed to delete SSH key: %v", err)
			util.LogError(errMsg, "ssh key")
			deletionErrors = append(deletionErrors, errMsg)
		} else {
			util.LogSuccess("SSH key deleted", "ssh key")
		}
	}

	fmt.Println()

	// Report final status based on whether errors occurred
	if len(deletionErrors) > 0 {
		util.LogError("Cluster deletion completed with errors!", d.Config.ClusterName)
		util.LogWarning("The following resources failed to delete:", "")
		for _, errMsg := range deletionErrors {
			fmt.Printf("  - %s\n", errMsg)
		}
		return fmt.Errorf("cluster deletion completed with %d error(s)", len(deletionErrors))
	}

	// Clean up kubeconfig file if it exists
	d.cleanupKubeconfig()

	util.LogSuccess("Cluster deletion completed successfully!", d.Config.ClusterName)
	return nil
}

// requestClusterNameConfirmation prompts the user to confirm deletion by typing the cluster name
func (d *Deleter) requestClusterNameConfirmation() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Please enter the cluster name to confirm that you want to delete it: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		if input == "" {
			util.LogError("Input cannot be empty. Please enter the cluster name.", "")
			continue
		}

		if input != d.Config.ClusterName {
			util.LogError(fmt.Sprintf("Cluster name '%s' does not match expected '%s'. Aborting deletion.", input, d.Config.ClusterName), "")
			return fmt.Errorf("cluster name confirmation failed")
		}

		break
	}

	return nil
}

// cleanupKubeconfig removes the kubeconfig file if it exists
func (d *Deleter) cleanupKubeconfig() {
	kubeconfigPath, err := config.ExpandPath(d.Config.KubeconfigPath)
	if err != nil {
		util.LogWarning(fmt.Sprintf("Failed to expand kubeconfig path: %v", err), "kubeconfig")
		return
	}

	if _, err := os.Stat(kubeconfigPath); err == nil {
		if err := os.Remove(kubeconfigPath); err != nil {
			util.LogWarning(fmt.Sprintf("Failed to delete kubeconfig file: %v", err), "kubeconfig")
		} else {
			util.LogSuccess("Kubeconfig file deleted", "kubeconfig")
		}
	}
}

// findAutoscaledPoolServers finds servers created by the cluster autoscaler
// These servers have the HCloudNodeGroupLabel label instead of the cluster label
func (d *Deleter) findAutoscaledPoolServers() ([]*hcloud.Server, error) {
	var allServers []*hcloud.Server

	// Iterate through all worker node pools
	for _, pool := range d.Config.WorkerNodePools {
		// Only process autoscaling-enabled pools
		if !pool.AutoscalingEnabled() {
			continue
		}

		// Build the node pool name (must match the name used by cluster autoscaler)
		poolName := pool.BuildNodePoolName(d.Config.ClusterName)

		// Search for servers with the HCloudNodeGroupLabel
		nodeGroupLabel := fmt.Sprintf("%s=%s", HCloudNodeGroupLabel, poolName)
		servers, err := d.HetznerClient.ListServers(d.ctx, hcloud.ServerListOpts{
			ListOpts: hcloud.ListOpts{
				LabelSelector: nodeGroupLabel,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list servers for node group %s: %w", poolName, err)
		}

		allServers = append(allServers, servers...)
	}

	return allServers, nil
}
