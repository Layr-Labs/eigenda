package config

// ServerConfig ... Config for the proxy HTTP server
type ServerConfig struct {
	DisperseToV2 bool
	Host         string
	Port         int
}
