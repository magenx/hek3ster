package hetzner

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/pkg/version"
)

// Client wraps the Hetzner Cloud client
type Client struct {
	hcloud *hcloud.Client
	token  string
}

// NewClient creates a new Hetzner client
func NewClient(token string) *Client {
	opts := []hcloud.ClientOption{
		hcloud.WithToken(token),
		hcloud.WithApplication("hek3ster", version.Get()),
	}

	return &Client{
		hcloud: hcloud.NewClient(opts...),
		token:  token,
	}
}

// GetLocations returns all available locations
func (c *Client) GetLocations(ctx context.Context) ([]*hcloud.Location, error) {
	locations, err := c.hcloud.Location.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	return locations, nil
}

// GetServerTypes returns all available server types
func (c *Client) GetServerTypes(ctx context.Context) ([]*hcloud.ServerType, error) {
	serverTypes, err := c.hcloud.ServerType.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server types: %w", err)
	}
	return serverTypes, nil
}

// GetServerType returns a specific server type by name
func (c *Client) GetServerType(ctx context.Context, name string) (*hcloud.ServerType, error) {
	serverType, _, err := c.hcloud.ServerType.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server type %s: %w", name, err)
	}
	if serverType == nil {
		return nil, fmt.Errorf("server type %s not found", name)
	}
	return serverType, nil
}

// GetLocation returns a specific location by name
func (c *Client) GetLocation(ctx context.Context, name string) (*hcloud.Location, error) {
	location, _, err := c.hcloud.Location.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location %s: %w", name, err)
	}
	if location == nil {
		return nil, fmt.Errorf("location %s not found", name)
	}
	return location, nil
}

// GetImage returns a specific image by name or ID
func (c *Client) GetImage(ctx context.Context, nameOrID string) (*hcloud.Image, error) {
	image, _, err := c.hcloud.Image.GetByName(ctx, nameOrID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image %s: %w", nameOrID, err)
	}
	if image == nil {
		return nil, fmt.Errorf("image %s not found", nameOrID)
	}
	return image, nil
}

// ListServers returns all servers matching the label selector
func (c *Client) ListServers(ctx context.Context, opts hcloud.ServerListOpts) ([]*hcloud.Server, error) {
	servers, err := c.hcloud.Server.AllWithOpts(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}
	return servers, nil
}

// GetServer returns a specific server by name
func (c *Client) GetServer(ctx context.Context, name string) (*hcloud.Server, error) {
	server, _, err := c.hcloud.Server.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server %s: %w", name, err)
	}
	return server, nil
}

// CreateServer creates a new server
func (c *Client) CreateServer(ctx context.Context, opts hcloud.ServerCreateOpts) (*hcloud.Server, error) {
	result, _, err := c.hcloud.Server.Create(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	// Wait for the action to complete
	if result.Action != nil {
		if err := c.waitForAction(ctx, result.Action); err != nil {
			return nil, fmt.Errorf("server creation action failed: %w", err)
		}
	}

	// Wait for any next actions
	if len(result.NextActions) > 0 {
		if err := c.waitForActions(ctx, result.NextActions); err != nil {
			return nil, fmt.Errorf("server creation next actions failed: %w", err)
		}
	}

	return result.Server, nil
}

// DeleteServer deletes a server
func (c *Client) DeleteServer(ctx context.Context, server *hcloud.Server) error {
	result, _, err := c.hcloud.Server.DeleteWithResult(ctx, server)
	if err != nil {
		return fmt.Errorf("failed to delete server %s: %w", server.Name, err)
	}

	// Wait for the delete action to complete
	if result.Action != nil {
		if err := c.waitForAction(ctx, result.Action); err != nil {
			return fmt.Errorf("server deletion action failed: %w", err)
		}
	}

	return nil
}

// CreateNetwork creates a new network
func (c *Client) CreateNetwork(ctx context.Context, opts hcloud.NetworkCreateOpts) (*hcloud.Network, error) {
	network, _, err := c.hcloud.Network.Create(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}
	return network, nil
}

// GetNetwork returns a specific network by name
func (c *Client) GetNetwork(ctx context.Context, name string) (*hcloud.Network, error) {
	network, _, err := c.hcloud.Network.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch network %s: %w", name, err)
	}
	return network, nil
}

// AddRouteToNetwork adds a route to a network
func (c *Client) AddRouteToNetwork(ctx context.Context, network *hcloud.Network, opts hcloud.NetworkAddRouteOpts) error {
	action, _, err := c.hcloud.Network.AddRoute(ctx, network, opts)
	if err != nil {
		return fmt.Errorf("failed to add route to network %s: %w", network.Name, err)
	}

	// Wait for the action to complete
	if err := c.waitForAction(ctx, action); err != nil {
		return fmt.Errorf("network add route action failed: %w", err)
	}

	return nil
}

// DeleteNetwork deletes a network
func (c *Client) DeleteNetwork(ctx context.Context, network *hcloud.Network) error {
	_, err := c.hcloud.Network.Delete(ctx, network)
	if err != nil {
		return fmt.Errorf("failed to delete network %s: %w", network.Name, err)
	}

	// Wait for network to actually be deleted
	// Poll with timeout to prevent infinite loops on persistent API issues
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(2 * time.Minute)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timeout waiting for network %s to be deleted", network.Name)
		case <-ticker.C:
			// Check if network still exists
			net, _, err := c.hcloud.Network.GetByID(ctx, network.ID)
			if err != nil {
				// Only treat 'not found' errors as successful deletion
				if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
					return nil
				}
				// For other errors (e.g., transient network issues, API throttling),
				// continue polling rather than failing the deletion immediately.
				// The timeout will prevent infinite loops on persistent API issues.
			}
			if net == nil {
				// Network is deleted
				return nil
			}
		}
	}
}

// CreateSSHKey creates a new SSH key
func (c *Client) CreateSSHKey(ctx context.Context, opts hcloud.SSHKeyCreateOpts) (*hcloud.SSHKey, error) {
	sshKey, _, err := c.hcloud.SSHKey.Create(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH key: %w", err)
	}
	return sshKey, nil
}

// GetSSHKey returns a specific SSH key by name
func (c *Client) GetSSHKey(ctx context.Context, name string) (*hcloud.SSHKey, error) {
	sshKey, _, err := c.hcloud.SSHKey.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SSH key %s: %w", name, err)
	}
	return sshKey, nil
}

// DeleteSSHKey deletes an SSH key
func (c *Client) DeleteSSHKey(ctx context.Context, sshKey *hcloud.SSHKey) error {
	_, err := c.hcloud.SSHKey.Delete(ctx, sshKey)
	if err != nil {
		return fmt.Errorf("failed to delete SSH key %s: %w", sshKey.Name, err)
	}
	return nil
}

// CreateFirewall creates a new firewall
func (c *Client) CreateFirewall(ctx context.Context, opts hcloud.FirewallCreateOpts) (*hcloud.Firewall, error) {
	result, _, err := c.hcloud.Firewall.Create(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create firewall: %w", err)
	}

	// Wait for any actions to complete
	if len(result.Actions) > 0 {
		if err := c.waitForActions(ctx, result.Actions); err != nil {
			return nil, fmt.Errorf("firewall creation action failed: %w", err)
		}
	}

	return result.Firewall, nil
}

// GetFirewall returns a specific firewall by name
func (c *Client) GetFirewall(ctx context.Context, name string) (*hcloud.Firewall, error) {
	firewall, _, err := c.hcloud.Firewall.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch firewall %s: %w", name, err)
	}
	return firewall, nil
}

// DeleteFirewall deletes a firewall
func (c *Client) DeleteFirewall(ctx context.Context, firewall *hcloud.Firewall) error {
	_, err := c.hcloud.Firewall.Delete(ctx, firewall)
	if err != nil {
		return fmt.Errorf("failed to delete firewall %s: %w", firewall.Name, err)
	}

	// Wait for firewall to actually be deleted
	// Poll with timeout to prevent infinite loops on persistent API issues
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(2 * time.Minute)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timeout waiting for firewall %s to be deleted", firewall.Name)
		case <-ticker.C:
			// Check if firewall still exists
			fw, _, err := c.hcloud.Firewall.GetByID(ctx, firewall.ID)
			if err != nil {
				// Only treat 'not found' errors as successful deletion
				if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
					return nil
				}
				// For other errors (e.g., transient network issues, API throttling),
				// continue polling rather than failing the deletion immediately.
				// The timeout will prevent infinite loops on persistent API issues.
			}
			if fw == nil {
				// Firewall is deleted
				return nil
			}
		}
	}
}

// CreateLoadBalancer creates a new load balancer
func (c *Client) CreateLoadBalancer(ctx context.Context, opts hcloud.LoadBalancerCreateOpts) (*hcloud.LoadBalancer, error) {
	result, _, err := c.hcloud.LoadBalancer.Create(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Wait for the action to complete
	if err := c.waitForAction(ctx, result.Action); err != nil {
		return nil, fmt.Errorf("load balancer creation action failed: %w", err)
	}

	return result.LoadBalancer, nil
}

// GetLoadBalancer returns a specific load balancer by name
func (c *Client) GetLoadBalancer(ctx context.Context, name string) (*hcloud.LoadBalancer, error) {
	lb, _, err := c.hcloud.LoadBalancer.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch load balancer %s: %w", name, err)
	}
	return lb, nil
}

// DeleteLoadBalancer deletes a load balancer
func (c *Client) DeleteLoadBalancer(ctx context.Context, lb *hcloud.LoadBalancer) error {
	_, err := c.hcloud.LoadBalancer.Delete(ctx, lb)
	if err != nil {
		return fmt.Errorf("failed to delete load balancer %s: %w", lb.Name, err)
	}

	// Wait for load balancer to actually be deleted
	// Poll with timeout to prevent infinite loops on persistent API issues
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(2 * time.Minute)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timeout waiting for load balancer %s to be deleted", lb.Name)
		case <-ticker.C:
			// Check if load balancer still exists
			loadBalancer, _, err := c.hcloud.LoadBalancer.GetByID(ctx, lb.ID)
			if err != nil {
				// Only treat 'not found' errors as successful deletion
				if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
					return nil
				}
				// For other errors (e.g., transient network issues, API throttling),
				// continue polling rather than failing the deletion immediately.
				// The timeout will prevent infinite loops on persistent API issues.
			}
			if loadBalancer == nil {
				// Load balancer is deleted
				return nil
			}
		}
	}
}

// AddServiceToLoadBalancer adds a service to an existing load balancer
func (c *Client) AddServiceToLoadBalancer(ctx context.Context, lb *hcloud.LoadBalancer, opts hcloud.LoadBalancerAddServiceOpts) error {
	action, _, err := c.hcloud.LoadBalancer.AddService(ctx, lb, opts)
	if err != nil {
		return fmt.Errorf("failed to add service to load balancer %s: %w", lb.Name, err)
	}

	// Wait for the action to complete
	if err := c.waitForAction(ctx, action); err != nil {
		return fmt.Errorf("add service action failed: %w", err)
	}

	return nil
}

// AddLabelSelectorTargetToLoadBalancer adds a label selector target to an existing load balancer
func (c *Client) AddLabelSelectorTargetToLoadBalancer(ctx context.Context, lb *hcloud.LoadBalancer, opts hcloud.LoadBalancerAddLabelSelectorTargetOpts) error {
	action, _, err := c.hcloud.LoadBalancer.AddLabelSelectorTarget(ctx, lb, opts)
	if err != nil {
		return fmt.Errorf("failed to add label selector target to load balancer %s: %w", lb.Name, err)
	}

	// Wait for the action to complete
	if err := c.waitForAction(ctx, action); err != nil {
		return fmt.Errorf("add label selector target action failed: %w", err)
	}

	return nil
}

// AddServerTargetToLoadBalancer adds a server target to an existing load balancer
func (c *Client) AddServerTargetToLoadBalancer(ctx context.Context, lb *hcloud.LoadBalancer, opts hcloud.LoadBalancerAddServerTargetOpts) error {
	action, _, err := c.hcloud.LoadBalancer.AddServerTarget(ctx, lb, opts)
	if err != nil {
		return fmt.Errorf("failed to add server target to load balancer %s: %w", lb.Name, err)
	}

	// Wait for the action to complete
	if err := c.waitForAction(ctx, action); err != nil {
		return fmt.Errorf("add server target action failed: %w", err)
	}

	return nil
}

// waitForAction waits for a single action to complete
func (c *Client) waitForAction(ctx context.Context, action *hcloud.Action) error {
	if action == nil {
		return nil
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			a, _, err := c.hcloud.Action.GetByID(ctx, action.ID)
			if err != nil {
				return fmt.Errorf("failed to get action status: %w", err)
			}

			if a.Status == hcloud.ActionStatusSuccess {
				return nil
			}
			if a.Status == hcloud.ActionStatusError {
				return fmt.Errorf("action failed: %s", a.ErrorMessage)
			}
		}
	}
}

// waitForActions waits for multiple actions to complete
func (c *Client) waitForActions(ctx context.Context, actions []*hcloud.Action) error {
	for _, action := range actions {
		if err := c.waitForAction(ctx, action); err != nil {
			return err
		}
	}
	return nil
}

// WaitForServerStatus waits for a server to reach the specified status
func (c *Client) WaitForServerStatus(ctx context.Context, server *hcloud.Server, targetStatus hcloud.ServerStatus, timeout time.Duration) error {
	if server == nil {
		return fmt.Errorf("server is nil")
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout waiting for server %s to reach status %s", server.Name, targetStatus)
		case <-ticker.C:
			// Refresh server status
			srv, _, err := c.hcloud.Server.GetByID(ctx, server.ID)
			if err != nil {
				return fmt.Errorf("failed to get server status: %w", err)
			}
			if srv == nil {
				return fmt.Errorf("server %s not found", server.Name)
			}

			if srv.Status == targetStatus {
				return nil
			}
		}
	}
}

// GetHCloudClient returns the underlying hcloud client for advanced operations
func (c *Client) GetHCloudClient() *hcloud.Client {
	return c.hcloud
}
