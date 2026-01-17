package cluster

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/magenx/hek3ster/internal/config"
	"github.com/magenx/hek3ster/internal/util"
	"github.com/magenx/hek3ster/pkg/hetzner"
)

// NetworkResourceManager handles creation of firewalls and load balancers
type NetworkResourceManager struct {
	Config        *config.Main
	HetznerClient *hetzner.Client
	ctx           context.Context
}

// NewNetworkResourceManager creates a new network resource manager
func NewNetworkResourceManager(cfg *config.Main, hetznerClient *hetzner.Client) *NetworkResourceManager {
	return &NetworkResourceManager{
		Config:        cfg,
		HetznerClient: hetznerClient,
		ctx:           context.Background(),
	}
}

// CreateAPILoadBalancer creates a load balancer for the Kubernetes API
func (n *NetworkResourceManager) CreateAPILoadBalancer(masterServers []*hcloud.Server, location string, network *hcloud.Network) (*hcloud.LoadBalancer, error) {
	util.LogInfo("Creating load balancer for Kubernetes API", "load balancer")

	lbName := fmt.Sprintf("%s-api-lb", n.Config.ClusterName)

	// Check if load balancer already exists
	existingLB, err := n.HetznerClient.GetLoadBalancer(n.ctx, lbName)
	if err == nil && existingLB != nil {
		util.LogInfo("API load balancer already exists, using existing load balancer", "load balancer")
		return existingLB, nil
	}

	// Determine if API load balancer should be attached to network
	// Only attach if private network is enabled and network object is provided
	shouldAttachToNetwork := n.Config.Networking.PrivateNetwork.Enabled && network != nil

	// Create load balancer without targets first to avoid API validation issues
	// Targets will be added after creation
	opts := hcloud.LoadBalancerCreateOpts{
		Name:             lbName,
		LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"}, // Smallest LB type
		Location:         &hcloud.Location{Name: location},
		Labels: map[string]string{
			"cluster": n.Config.ClusterName,
			"role":    "api-lb",
			"managed": "hek3ster",
		},
		Services: []hcloud.LoadBalancerCreateOptsService{
			{
				Protocol:        hcloud.LoadBalancerServiceProtocolTCP,
				ListenPort:      hcloud.Ptr(6443),
				DestinationPort: hcloud.Ptr(6443),
				HealthCheck: &hcloud.LoadBalancerCreateOptsServiceHealthCheck{
					Protocol: hcloud.LoadBalancerServiceProtocolTCP,
					Port:     hcloud.Ptr(6443),
					Interval: hcloud.Duration(15 * time.Second),
					Timeout:  hcloud.Duration(10 * time.Second),
					Retries:  hcloud.Ptr(3),
				},
			},
		},
		PublicInterface: hcloud.Ptr(true),
	}

	// Attach to network if private network is enabled
	if shouldAttachToNetwork {
		opts.Network = network
	}

	lb, err := n.HetznerClient.CreateLoadBalancer(n.ctx, opts)

	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	util.LogSuccess(fmt.Sprintf("Load balancer created: %s (IP: %s)", lbName, lb.PublicNet.IPv4.IP.String()), "load balancer")

	// Default retry configuration for network attachment verification
	const (
		maxNetworkAttachmentRetries = 5
		initialRetryDelay           = 2 * time.Second
		stabilizationDelay          = 5 * time.Second // Additional delay after network attachment is verified
	)

	// If load balancer was created with a network attachment, refresh its data from the API
	// to ensure the network attachment is fully reflected before adding targets
	if shouldAttachToNetwork {
		util.LogInfo("Refreshing load balancer data to verify network attachment", "load balancer")

		// Retry with exponential backoff to ensure network attachment is complete
		retryDelay := initialRetryDelay

		for attempt := 0; attempt < maxNetworkAttachmentRetries; attempt++ {
			if attempt > 0 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
			}

			refreshedLB, err := n.HetznerClient.GetLoadBalancer(n.ctx, lbName)
			if err != nil {
				return nil, fmt.Errorf("failed to refresh load balancer: %w", err)
			}
			if refreshedLB == nil {
				return nil, fmt.Errorf("load balancer %s not found after creation", lbName)
			}

			// Check if network is attached (PrivateNet should be populated)
			if len(refreshedLB.PrivateNet) > 0 {
				lb = refreshedLB
				util.LogInfo("Load balancer network attachment verified", "load balancer")
				// Add stabilization delay to ensure network attachment is fully operational
				// in Hetzner's backend before attempting to add targets
				util.LogInfo("Waiting for network attachment to stabilize", "load balancer")
				time.Sleep(stabilizationDelay)
				break
			}

			if attempt == maxNetworkAttachmentRetries-1 {
				return nil, fmt.Errorf("load balancer network attachment not complete after %d attempts", maxNetworkAttachmentRetries)
			}

			util.LogInfo(fmt.Sprintf("Network attachment in progress, retrying... (attempt %d/%d)", attempt+2, maxNetworkAttachmentRetries), "load balancer")
		}
	} else {
		// Add a stabilization delay to ensure the load balancer is fully operational
		// in Hetzner's backend before attempting to add server targets
		// This prevents "cloud target was not found" errors
		util.LogInfo("Waiting for load balancer to stabilize before adding targets", "load balancer")
		time.Sleep(stabilizationDelay)
	}

	// Add master servers as targets using label selector
	// This is more reliable than adding individual server targets
	util.LogInfo("Adding master servers as targets to API load balancer", "load balancer")
	labelSelector := fmt.Sprintf("role=master,cluster=%s", n.Config.ClusterName)

	// Use private IP if network is attached, otherwise use public IP
	usePrivateIP := shouldAttachToNetwork

	err = n.HetznerClient.AddLabelSelectorTargetToLoadBalancer(n.ctx, lb, hcloud.LoadBalancerAddLabelSelectorTargetOpts{
		Selector:     labelSelector,
		UsePrivateIP: hcloud.Ptr(usePrivateIP),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add master servers as targets: %w", err)
	}

	util.LogSuccess("Master servers added as targets to API load balancer", "load balancer")

	return lb, nil
}

// CreateClusterFirewall creates a firewall for the cluster
func (n *NetworkResourceManager) CreateClusterFirewall(network *hcloud.Network) (*hcloud.Firewall, error) {
	util.LogInfo("Creating firewall for cluster", "firewall")

	fwName := fmt.Sprintf("%s-firewall", n.Config.ClusterName)

	// Define firewall rules
	var rules []hcloud.FirewallRule

	// Allow SSH from configured networks
	if len(n.Config.Networking.AllowedNetworks.SSH) > 0 {
		for _, cidr := range n.Config.Networking.AllowedNetworks.SSH {
			rules = append(rules, hcloud.FirewallRule{
				Direction:   hcloud.FirewallRuleDirectionIn,
				SourceIPs:   []net.IPNet{parseCIDR(cidr)},
				Protocol:    hcloud.FirewallRuleProtocolTCP,
				Port:        hcloud.Ptr(fmt.Sprintf("%d", n.Config.Networking.SSH.Port)),
				Description: hcloud.Ptr("SSH access"),
			})
		}
	}

	// Allow Kubernetes API from configured networks
	if len(n.Config.Networking.AllowedNetworks.API) > 0 {
		for _, cidr := range n.Config.Networking.AllowedNetworks.API {
			rules = append(rules, hcloud.FirewallRule{
				Direction:   hcloud.FirewallRuleDirectionIn,
				SourceIPs:   []net.IPNet{parseCIDR(cidr)},
				Protocol:    hcloud.FirewallRuleProtocolTCP,
				Port:        hcloud.Ptr("6443"),
				Description: hcloud.Ptr("Kubernetes API access"),
			})
		}
	}

	// Allow all traffic within private network if enabled
	if n.Config.Networking.PrivateNetwork.Enabled {
		_, privateNet, err := net.ParseCIDR(n.Config.Networking.PrivateNetwork.Subnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private network subnet: %w", err)
		}
		rules = append(rules, hcloud.FirewallRule{
			Direction:   hcloud.FirewallRuleDirectionIn,
			SourceIPs:   []net.IPNet{*privateNet},
			Protocol:    hcloud.FirewallRuleProtocolTCP,
			Port:        hcloud.Ptr("1-65535"),
			Description: hcloud.Ptr("Allow all TCP within cluster network"),
		})
		rules = append(rules, hcloud.FirewallRule{
			Direction:   hcloud.FirewallRuleDirectionIn,
			SourceIPs:   []net.IPNet{*privateNet},
			Protocol:    hcloud.FirewallRuleProtocolUDP,
			Port:        hcloud.Ptr("1-65535"),
			Description: hcloud.Ptr("Allow all UDP within cluster network"),
		})
		rules = append(rules, hcloud.FirewallRule{
			Direction:   hcloud.FirewallRuleDirectionIn,
			SourceIPs:   []net.IPNet{*privateNet},
			Protocol:    hcloud.FirewallRuleProtocolICMP,
			Description: hcloud.Ptr("Allow ICMP within cluster network"),
		})
	}

	// Create firewall
	fw, err := n.HetznerClient.CreateFirewall(n.ctx, hcloud.FirewallCreateOpts{
		Name: fwName,
		Labels: map[string]string{
			"cluster": n.Config.ClusterName,
		},
		Rules: rules,
		ApplyTo: []hcloud.FirewallResource{
			{
				Type: hcloud.FirewallResourceTypeLabelSelector,
				LabelSelector: &hcloud.FirewallResourceLabelSelector{
					Selector: fmt.Sprintf("cluster=%s", n.Config.ClusterName),
				},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create firewall: %w", err)
	}

	util.LogSuccess(fmt.Sprintf("Firewall created: %s with %d rule(s)", fwName, len(rules)), "firewall")

	return fw, nil
}

// CreateGlobalLoadBalancer creates a global load balancer for application traffic
// It attaches to specified worker pools via label selector and configures services
// If SSL certificate is enabled and HTTPS services are configured, it will attach the certificate
func (n *NetworkResourceManager) CreateGlobalLoadBalancer(network *hcloud.Network, location string, certificate *hcloud.Certificate) (*hcloud.LoadBalancer, error) {
	// Skip if load balancer is not enabled
	if !n.Config.LoadBalancer.Enabled {
		util.LogInfo("Global load balancer is disabled, skipping creation", "load balancer")
		return nil, nil
	}

	util.LogInfo("Creating global load balancer for application traffic", "load balancer")

	// Determine load balancer name
	lbName := fmt.Sprintf("%s-global-lb", n.Config.ClusterName)
	if n.Config.LoadBalancer.Name != nil {
		lbName = *n.Config.LoadBalancer.Name
	}

	// Check if load balancer already exists
	existingLB, err := n.HetznerClient.GetLoadBalancer(n.ctx, lbName)
	if err == nil && existingLB != nil {
		util.LogInfo("Global load balancer already exists, using existing load balancer", "load balancer")
		return existingLB, nil
	}

	// Build services configuration
	var services []hcloud.LoadBalancerCreateOptsService
	for _, svc := range n.Config.LoadBalancer.Services {
		service := hcloud.LoadBalancerCreateOptsService{
			Protocol:        hcloud.LoadBalancerServiceProtocol(svc.Protocol),
			ListenPort:      hcloud.Ptr(svc.ListenPort),
			DestinationPort: hcloud.Ptr(svc.DestinationPort),
			Proxyprotocol:   hcloud.Ptr(svc.ProxyProtocol),
		}

		// Add HTTP configuration for HTTPS services with certificate
		if strings.ToLower(svc.Protocol) == "https" && certificate != nil {
			httpConfig := &hcloud.LoadBalancerCreateOptsServiceHTTP{
				Certificates: []*hcloud.Certificate{certificate},
			}
			service.HTTP = httpConfig
			util.LogInfo(fmt.Sprintf("Attaching SSL certificate to HTTPS service on port %d", svc.ListenPort), "load balancer")
		}

		// Add health check if configured
		if svc.HealthCheck != nil {
			healthCheck := &hcloud.LoadBalancerCreateOptsServiceHealthCheck{
				Protocol: hcloud.LoadBalancerServiceProtocol(svc.HealthCheck.Protocol),
				Port:     hcloud.Ptr(svc.HealthCheck.Port),
				Interval: hcloud.Duration(time.Duration(svc.HealthCheck.Interval) * time.Second),
				Timeout:  hcloud.Duration(time.Duration(svc.HealthCheck.Timeout) * time.Second),
				Retries:  hcloud.Ptr(svc.HealthCheck.Retries),
			}

			// Add HTTP-specific health check settings
			if svc.HealthCheck.HTTP != nil {
				healthCheck.HTTP = &hcloud.LoadBalancerCreateOptsServiceHealthCheckHTTP{
					Domain:      hcloud.Ptr(svc.HealthCheck.HTTP.Domain),
					Path:        hcloud.Ptr(svc.HealthCheck.HTTP.Path),
					StatusCodes: svc.HealthCheck.HTTP.StatusCodes,
					TLS:         hcloud.Ptr(svc.HealthCheck.HTTP.TLS),
				}
			}

			service.HealthCheck = healthCheck
		}

		services = append(services, service)
	}

	// Determine if load balancer should be attached to network
	// This must be calculated BEFORE determining usePrivateIP
	shouldAttachToNetwork := n.Config.LoadBalancer.AttachToNetwork && n.Config.Networking.PrivateNetwork.Enabled && network != nil

	// Determine if we should use private IP for targets
	// Private IPs can only be used if the load balancer is attached to the network
	usePrivateIP := false
	if n.Config.LoadBalancer.UsePrivateIP != nil {
		// If explicitly set, use that value (but it must match network attachment)
		usePrivateIP = *n.Config.LoadBalancer.UsePrivateIP && shouldAttachToNetwork
	} else if shouldAttachToNetwork {
		// If not explicitly set and load balancer is attached to network, default to true
		usePrivateIP = true
	}

	// Determine location for load balancer
	lbLocation := location
	if n.Config.LoadBalancer.Location != "" {
		lbLocation = n.Config.LoadBalancer.Location
	}

	// Build create options without targets (to avoid API validation issues)
	// Targets will be added after creation
	opts := hcloud.LoadBalancerCreateOpts{
		Name:             lbName,
		LoadBalancerType: &hcloud.LoadBalancerType{Name: n.Config.LoadBalancer.Type},
		Location:         &hcloud.Location{Name: lbLocation},
		Labels: map[string]string{
			"cluster": n.Config.ClusterName,
			"role":    "global-lb",
			"managed": "hek3ster",
		},
		Algorithm: &hcloud.LoadBalancerAlgorithm{
			Type: hcloud.LoadBalancerAlgorithmType(n.Config.LoadBalancer.Algorithm.Type),
		},
		Services:        services,
		PublicInterface: hcloud.Ptr(true),
	}

	// Attach to network if configured and private network is enabled
	if shouldAttachToNetwork {
		opts.Network = network
	}

	// Create load balancer without targets first
	lb, err := n.HetznerClient.CreateLoadBalancer(n.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Default retry configuration for network attachment verification
	const (
		maxNetworkAttachmentRetries = 5
		initialRetryDelay           = 2 * time.Second
		stabilizationDelay          = 5 * time.Second // Additional delay after network attachment is verified
	)

	util.LogSuccess(fmt.Sprintf("Global load balancer created: %s (IP: %s)", lbName, lb.PublicNet.IPv4.IP.String()), "load balancer")

	// If load balancer was created with a network attachment, refresh its data from the API
	// to ensure the network attachment is fully reflected before adding targets
	if shouldAttachToNetwork {
		util.LogInfo("Refreshing load balancer data to verify network attachment", "load balancer")

		// Retry with exponential backoff to ensure network attachment is complete
		retryDelay := initialRetryDelay

		for attempt := 0; attempt < maxNetworkAttachmentRetries; attempt++ {
			if attempt > 0 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
			}

			refreshedLB, err := n.HetznerClient.GetLoadBalancer(n.ctx, lbName)
			if err != nil {
				return nil, fmt.Errorf("failed to refresh load balancer: %w", err)
			}
			if refreshedLB == nil {
				return nil, fmt.Errorf("load balancer %s not found after creation", lbName)
			}

			// Check if network is attached (PrivateNet should be populated)
			if len(refreshedLB.PrivateNet) > 0 {
				lb = refreshedLB
				util.LogInfo("Load balancer network attachment verified", "load balancer")
				// Add stabilization delay to ensure network attachment is fully operational
				// in Hetzner's backend before attempting to add targets
				util.LogInfo("Waiting for network attachment to stabilize", "load balancer")
				time.Sleep(stabilizationDelay)
				break
			}

			if attempt == maxNetworkAttachmentRetries-1 {
				return nil, fmt.Errorf("load balancer network attachment not complete after %d attempts", maxNetworkAttachmentRetries)
			}

			util.LogInfo(fmt.Sprintf("Network attachment in progress, retrying... (attempt %d/%d)", attempt+2, maxNetworkAttachmentRetries), "load balancer")
		}
	}

	// Add label selector targets after creation
	util.LogInfo("Adding targets to load balancer", "load balancer")

	if len(n.Config.LoadBalancer.TargetPools) > 0 {
		// Add a separate label selector target for each pool
		for _, poolName := range n.Config.LoadBalancer.TargetPools {
			labelSelector := fmt.Sprintf("pool=%s", poolName)
			err = n.HetznerClient.AddLabelSelectorTargetToLoadBalancer(n.ctx, lb, hcloud.LoadBalancerAddLabelSelectorTargetOpts{
				Selector:     labelSelector,
				UsePrivateIP: hcloud.Ptr(usePrivateIP),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to add targets to load balancer for pool %s: %w", poolName, err)
			}
		}
	} else {
		// Default to all worker nodes in the cluster
		labelSelector := fmt.Sprintf("role=worker,cluster=%s", n.Config.ClusterName)
		err = n.HetznerClient.AddLabelSelectorTargetToLoadBalancer(n.ctx, lb, hcloud.LoadBalancerAddLabelSelectorTargetOpts{
			Selector:     labelSelector,
			UsePrivateIP: hcloud.Ptr(usePrivateIP),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add targets to load balancer: %w", err)
		}
	}

	util.LogSuccess("Targets added to global load balancer", "load balancer")

	return lb, nil
}

// CreateDNSZone creates a DNS zone in Hetzner for the domain
func (n *NetworkResourceManager) CreateDNSZone() (*hcloud.Zone, error) {
	// Skip if DNS zone is not enabled
	if !n.Config.DNSZone.Enabled {
		util.LogInfo("DNS zone creation is disabled, skipping", "dns")
		return nil, nil
	}

	// Skip if domain is not set
	if n.Config.Domain == "" {
		util.LogInfo("Domain is not set, skipping DNS zone creation", "dns")
		return nil, nil
	}

	util.LogInfo(fmt.Sprintf("Creating DNS zone for domain: %s", n.Config.Domain), "dns")

	// Determine zone name (use override if provided, otherwise use domain)
	zoneName := n.Config.Domain
	if n.Config.DNSZone.Name != "" {
		zoneName = n.Config.DNSZone.Name
	}

	// Check if zone already exists
	existingZone, err := n.HetznerClient.GetZone(n.ctx, zoneName)
	if err == nil && existingZone != nil {
		util.LogInfo(fmt.Sprintf("DNS zone already exists: %s", zoneName), "dns")
		return existingZone, nil
	}

	// Create DNS zone
	zone, err := n.HetznerClient.CreateZone(n.ctx, hcloud.ZoneCreateOpts{
		Name: zoneName,
		Mode: hcloud.ZoneModePrimary,
		TTL:  hcloud.Ptr(n.Config.DNSZone.TTL),
		Labels: map[string]string{
			"cluster": n.Config.ClusterName,
			"managed": "hek3ster",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS zone: %w", err)
	}

	util.LogSuccess(fmt.Sprintf("DNS zone created: %s", zoneName), "dns")

	// Display nameservers information
	if len(zone.AuthoritativeNameservers.Assigned) > 0 {
		util.LogInfo("DNS zone nameservers:", "dns")
		for _, ns := range zone.AuthoritativeNameservers.Assigned {
			util.LogInfo(fmt.Sprintf("  - %s", ns), "dns")
		}
		util.LogInfo(fmt.Sprintf("Update your domain registrar to use these nameservers for domain: %s", n.Config.Domain), "dns")
	}

	return zone, nil
}

// parseCIDR parses a CIDR string and returns net.IPNet
func parseCIDR(cidr string) net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		// Return a restrictive default that won't match anything if parsing fails
		// This is safer than allowing 0.0.0.0/0
		_, ipnet, _ = net.ParseCIDR("127.0.0.1/32")
		util.LogWarning(fmt.Sprintf("Failed to parse CIDR %s, using restrictive default", cidr), "firewall")
	}
	return *ipnet
}

// CreateSSLCertificate creates a managed SSL certificate for the domain
// The certificate will cover both the root domain and wildcard subdomain
func (n *NetworkResourceManager) CreateSSLCertificate() (*hcloud.Certificate, error) {
	// Skip if SSL certificate is not enabled
	if !n.Config.SSLCertificate.Enabled {
		util.LogInfo("SSL certificate creation is disabled, skipping", "ssl")
		return nil, nil
	}

	// Skip if domain is not set
	if n.Config.Domain == "" {
		util.LogInfo("Domain is not set, skipping SSL certificate creation", "ssl")
		return nil, nil
	}

	util.LogInfo(fmt.Sprintf("Creating managed SSL certificate for domain: %s", n.Config.Domain), "ssl")

	// Determine certificate name (use override if provided, otherwise use domain)
	certName := n.Config.Domain
	if n.Config.SSLCertificate.Name != "" {
		certName = n.Config.SSLCertificate.Name
	}

	// Determine domain for certificate (use override if provided, otherwise use domain from config)
	certDomain := n.Config.Domain
	if n.Config.SSLCertificate.Domain != "" {
		certDomain = n.Config.SSLCertificate.Domain
	}

	// Check if certificate already exists
	existingCert, err := n.HetznerClient.GetCertificate(n.ctx, certName)
	if err == nil && existingCert != nil {
		util.LogInfo(fmt.Sprintf("SSL certificate already exists: %s", certName), "ssl")
		return existingCert, nil
	}

	// Create managed certificate with root domain and wildcard
	// This allows the certificate to be used for both example.com and *.example.com
	domainNames := []string{
		certDomain,
		fmt.Sprintf("*.%s", certDomain),
	}

	cert, err := n.HetznerClient.CreateManagedCertificate(n.ctx, hcloud.CertificateCreateOpts{
		Name:        certName,
		Type:        hcloud.CertificateTypeManaged,
		DomainNames: domainNames,
		Labels: map[string]string{
			"cluster": n.Config.ClusterName,
			"managed": "hek3ster",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create SSL certificate: %w", err)
	}

	util.LogSuccess(fmt.Sprintf("SSL certificate created: %s (covers: %s)", certName, strings.Join(domainNames, ", ")), "ssl")
	util.LogInfo("Certificate validation will happen in the background (may take up to 5 minutes)", "ssl")
	util.LogInfo("The certificate will be automatically validated via DNS records in your DNS zone", "ssl")

	return cert, nil
}
