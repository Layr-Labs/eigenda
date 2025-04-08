package config

import (
	"slices"
)

// ServerConfig ... Config for the proxy HTTP server
type ServerConfig struct {
	Host string
	Port int
	// EnabledAPIs contains the list of API types that are enabled.
	// When empty (default), no special API endpoints are registered.
	// Example: If it contains "admin", administrative endpoints like /admin/eigenda-dispersal-backend will be available.
	EnabledAPIs []string
}

// IsAPIEnabled checks if a specific API type is enabled
func (c *ServerConfig) IsAPIEnabled(apiType string) bool {
	return slices.Contains(c.EnabledAPIs, apiType)
}
