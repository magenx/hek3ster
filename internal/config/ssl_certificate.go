package config

// SSLCertificate represents SSL certificate configuration for managed certificates
type SSLCertificate struct {
	Enabled bool   `yaml:"enabled,omitempty"`
	Name    string `yaml:"name,omitempty"`   // Certificate name in Hetzner
	Domain  string `yaml:"domain,omitempty"` // Domain for the certificate
}

// SetDefaults sets default values for SSL certificate configuration
func (s *SSLCertificate) SetDefaults() {
	// No specific defaults needed - all fields are optional
	// and will be derived from other config values if not set
}
