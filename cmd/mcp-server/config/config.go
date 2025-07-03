package config

import (
	"fmt"
	"os"
)

// Config holds the configuration for the MCP server
type Config struct {
	// Temporal Cloud API configuration
	CloudAPIKey string

	// Temporal namespace configuration
	Namespace         string
	NamespaceAPIKey   string
	NamespaceTLSCert  string
	NamespaceTLSKey   string

	// MCP server configuration
	ServerName string
	ServerVersion string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	config := &Config{
		CloudAPIKey:       os.Getenv("TEMPORAL_CLOUD_API_KEY"),
		Namespace:         os.Getenv("TEMPORAL_CLOUD_NAMESPACE"),
		NamespaceAPIKey:   os.Getenv("TEMPORAL_CLOUD_NAMESPACE_API_KEY"),
		NamespaceTLSCert:  os.Getenv("TEMPORAL_CLOUD_NAMESPACE_TLS_CERT"),
		NamespaceTLSKey:   os.Getenv("TEMPORAL_CLOUD_NAMESPACE_TLS_KEY"),
		ServerName:        getEnvOrDefault("MCP_SERVER_NAME", "temporal-cloud-mcp-server"),
		ServerVersion:     getEnvOrDefault("MCP_SERVER_VERSION", "1.0.0"),
	}

	// Validate required configuration
	if config.CloudAPIKey == "" {
		return nil, fmt.Errorf("TEMPORAL_CLOUD_API_KEY is required")
	}

	return config, nil
}

// HasNamespaceAuth returns true if namespace authentication is configured
func (c *Config) HasNamespaceAuth() bool {
	return c.NamespaceAPIKey != "" || (c.NamespaceTLSCert != "" && c.NamespaceTLSKey != "")
}

// HasmTLSAuth returns true if mTLS authentication is configured
func (c *Config) HasmTLSAuth() bool {
	return c.NamespaceTLSCert != "" && c.NamespaceTLSKey != ""
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}