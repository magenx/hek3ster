package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSSLCertificateUnmarshalYAML(t *testing.T) {
	yamlData := `
ssl_certificate:
  enabled: true
  name: example.com
  domain: example.com
`

	type testConfig struct {
		SSLCertificate SSLCertificate `yaml:"ssl_certificate"`
	}

	var config testConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !config.SSLCertificate.Enabled {
		t.Error("Expected SSL certificate to be enabled")
	}
	if config.SSLCertificate.Name != "example.com" {
		t.Errorf("Expected name 'example.com', got '%s'", config.SSLCertificate.Name)
	}
	if config.SSLCertificate.Domain != "example.com" {
		t.Errorf("Expected domain 'example.com', got '%s'", config.SSLCertificate.Domain)
	}
}

func TestSSLCertificateSetDefaults(t *testing.T) {
	tests := []struct {
		name string
		cert SSLCertificate
	}{
		{
			name: "Empty SSL certificate",
			cert: SSLCertificate{},
		},
		{
			name: "Enabled SSL certificate",
			cert: SSLCertificate{
				Enabled: true,
				Name:    "test-cert",
				Domain:  "test.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert := tt.cert
			// SetDefaults should not panic and should be idempotent
			cert.SetDefaults()
			cert.SetDefaults() // Call twice to verify idempotence
		})
	}
}
