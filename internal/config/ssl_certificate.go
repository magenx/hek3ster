package config

// SSLCertificate represents SSL certificate configuration for managed certificates
type SSLCertificate struct {
	Enabled bool   `yaml:"enabled,omitempty"`
	Name    string `yaml:"name,omitempty"`   // Certificate name in Hetzner
	Domain  string `yaml:"domain,omitempty"` // Domain for the certificate
}

// SetDefaults sets default values for SSL certificate configuration
func (s *SSLCertificate) SetDefaults() {
	// No defaults are set here - name and domain are derived from the main
	// configuration's domain field when creating the certificate if not explicitly set
}
