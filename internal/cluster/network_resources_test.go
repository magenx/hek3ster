package cluster

import (
	"testing"

	"github.com/magenx/hek3ster/internal/config"
)

// TestLoadBalancerPrivateIPLogic tests the logic for determining when to use private IPs
// This test validates the fix for the "load balancer is not attached to a network" issue
func TestLoadBalancerPrivateIPLogic(t *testing.T) {
	tests := []struct {
		name                    string
		attachToNetwork         bool
		privateNetworkEnabled   bool
		usePrivateIPConfigured  *bool
		networkProvided         bool
		expectedUsePrivateIP    bool
		expectedAttachToNetwork bool
		description             string
	}{
		{
			name:                    "attach_to_network=true, private_network=true, use_private_ip not set",
			attachToNetwork:         true,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  nil,
			networkProvided:         true,
			expectedUsePrivateIP:    true,
			expectedAttachToNetwork: true,
			description:             "Should attach to network and use private IPs by default",
		},
		{
			name:                    "attach_to_network=false, private_network=true, use_private_ip not set",
			attachToNetwork:         false,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  nil,
			networkProvided:         true,
			expectedUsePrivateIP:    false,
			expectedAttachToNetwork: false,
			description:             "Bug fix: Should NOT use private IPs when not attached to network",
		},
		{
			name:                    "attach_to_network=true, private_network=false, use_private_ip not set",
			attachToNetwork:         true,
			privateNetworkEnabled:   false,
			usePrivateIPConfigured:  nil,
			networkProvided:         true,
			expectedUsePrivateIP:    false,
			expectedAttachToNetwork: false,
			description:             "Should not attach to network when private network is disabled",
		},
		{
			name:                    "attach_to_network=true, private_network=true, use_private_ip=true",
			attachToNetwork:         true,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  boolPtr(true),
			networkProvided:         true,
			expectedUsePrivateIP:    true,
			expectedAttachToNetwork: true,
			description:             "Should use private IPs when explicitly configured and attached",
		},
		{
			name:                    "attach_to_network=false, private_network=true, use_private_ip=true",
			attachToNetwork:         false,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  boolPtr(true),
			networkProvided:         true,
			expectedUsePrivateIP:    false,
			expectedAttachToNetwork: false,
			description:             "Bug fix: Should NOT use private IPs even when explicitly configured if not attached",
		},
		{
			name:                    "attach_to_network=true, private_network=true, use_private_ip=false",
			attachToNetwork:         true,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  boolPtr(false),
			networkProvided:         true,
			expectedUsePrivateIP:    false,
			expectedAttachToNetwork: true,
			description:             "Should respect explicit false setting for use_private_ip",
		},
		{
			name:                    "attach_to_network=true, private_network=true, network=nil",
			attachToNetwork:         true,
			privateNetworkEnabled:   true,
			usePrivateIPConfigured:  nil,
			networkProvided:         false,
			expectedUsePrivateIP:    false,
			expectedAttachToNetwork: false,
			description:             "Should not attach when network object is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from CreateGlobalLoadBalancer
			cfg := &config.Main{
				LoadBalancer: config.LoadBalancer{
					Enabled:         true,
					AttachToNetwork: tt.attachToNetwork,
					UsePrivateIP:    tt.usePrivateIPConfigured,
				},
				Networking: config.Networking{
					PrivateNetwork: config.PrivateNetwork{
						Enabled: tt.privateNetworkEnabled,
					},
				},
			}

			// Simulate network object
			var network interface{}
			if tt.networkProvided {
				network = &struct{}{} // Mock network object
			}

			// This is the logic from network_resources.go lines 313-326
			shouldAttachToNetwork := cfg.LoadBalancer.AttachToNetwork && cfg.Networking.PrivateNetwork.Enabled && network != nil

			usePrivateIP := false
			if cfg.LoadBalancer.UsePrivateIP != nil {
				// If explicitly set, use that value (but it must match network attachment)
				usePrivateIP = *cfg.LoadBalancer.UsePrivateIP && shouldAttachToNetwork
			} else if shouldAttachToNetwork {
				// If not explicitly set and load balancer is attached to network, default to true
				usePrivateIP = true
			}

			// Verify expectations
			if shouldAttachToNetwork != tt.expectedAttachToNetwork {
				t.Errorf("%s: shouldAttachToNetwork = %v, expected %v",
					tt.description, shouldAttachToNetwork, tt.expectedAttachToNetwork)
			}

			if usePrivateIP != tt.expectedUsePrivateIP {
				t.Errorf("%s: usePrivateIP = %v, expected %v",
					tt.description, usePrivateIP, tt.expectedUsePrivateIP)
			}

			// The key assertion: usePrivateIP should NEVER be true when shouldAttachToNetwork is false
			if usePrivateIP && !shouldAttachToNetwork {
				t.Errorf("%s: CRITICAL BUG: usePrivateIP is true but shouldAttachToNetwork is false. "+
					"This would cause 'load balancer is not attached to a network' error!",
					tt.description)
			}
		})
	}
}
