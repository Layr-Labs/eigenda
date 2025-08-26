package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ChurnerConfig defines configuration for the churner container
type ChurnerConfig struct {
	// Enable churner service
	Enabled bool

	// Log format (text/json)
	LogFormat string

	// Network configuration
	Hostname string
	GRPCPort string

	// EigenDA configuration
	EigenDADirectory       string
	OperatorStateRetriever string
	ServiceManager         string

	// Chain configuration
	ChainRPC   string
	PrivateKey string

	// Graph configuration
	GraphURL            string
	IndexerPullInterval string

	// Metrics configuration
	EnableMetrics   bool
	MetricsHTTPPort string

	// Churner configuration
	ChurnApprovalInterval string

	// Container image
	Image string
}

// ChurnerContainer represents a running churner container
type ChurnerContainer struct {
	testcontainers.Container
	config  ChurnerConfig
	url     string
	logPath string // Path to log file on host
}

// DefaultChurnerConfig returns a default churner configuration
func DefaultChurnerConfig() ChurnerConfig {
	return ChurnerConfig{
		Enabled:                true,
		LogFormat:              "text",
		Hostname:               "0.0.0.0",
		GRPCPort:               "32001",
		EigenDADirectory:       "", // Will be populated from contract deployment
		OperatorStateRetriever: "", // Will be populated from contract deployment
		ServiceManager:         "", // Will be populated from contract deployment
		ChainRPC:               "", // Will be populated from Anvil
		PrivateKey:             "", // Will be populated from deployer key
		GraphURL:               "", // Will be populated from GraphNode if enabled
		IndexerPullInterval:    "1s",
		EnableMetrics:          true,
		MetricsHTTPPort:        "9095",
		ChurnApprovalInterval:  "900s",
		Image:                  "ghcr.io/layr-labs/eigenda/churner:dev",
	}
}

// NewChurnerContainerWithNetwork creates and starts a new churner container with a custom network
func NewChurnerContainerWithNetwork(ctx context.Context, config ChurnerConfig, network *testcontainers.DockerNetwork) (*ChurnerContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", "churner-logs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "churner.log")

	// Build environment variables
	env := buildChurnerEnv(config)

	// Configure container request with network
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{config.GRPCPort + "/tcp"},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {"churner"},
		},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: logDir,
				},
				Target: "/logs",
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.GRPCPort+"/tcp")).WithStartupTimeout(30*time.Second),
			wait.ForLog("churner server listening at").WithStartupTimeout(30*time.Second),
		),
		Name:            "eigenda-churner",
		AlwaysPullImage: false, // Use local image if available
	}

	// Add metrics port if enabled
	if config.EnableMetrics {
		req.ExposedPorts = append(req.ExposedPorts, config.MetricsHTTPPort+"/tcp")
	}

	// Create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start churner container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, nat.Port(config.GRPCPort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the container host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	churnerURL := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	fmt.Printf("Churner logs will be available at: %s\n", logPath)

	return &ChurnerContainer{
		Container: container,
		config:    config,
		url:       churnerURL,
		logPath:   logPath,
	}, nil
}

// URL returns the churner service URL
func (c *ChurnerContainer) URL() string {
	return c.url
}

// InternalURL returns the churner service URL for internal network communication
func (c *ChurnerContainer) InternalURL() string {
	return fmt.Sprintf("churner:%s", c.config.GRPCPort)
}

// Config returns the churner configuration
func (c *ChurnerContainer) Config() ChurnerConfig {
	return c.config
}

// LogPath returns the path to the churner log file on the host
func (c *ChurnerContainer) LogPath() string {
	return c.logPath
}

// buildChurnerEnv builds environment variables for the churner container
func buildChurnerEnv(config ChurnerConfig) map[string]string {
	// Strip 0x prefix from private key if present
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")

	env := map[string]string{
		"CHURNER_LOG_FORMAT":                  config.LogFormat,
		"CHURNER_LOG_PATH":                    "/logs/churner.log", // Log to a file we can access
		"CHURNER_LOG_LEVEL":                   "debug",             // Enable debug logging
		"CHURNER_HOSTNAME":                    config.Hostname,
		"CHURNER_GRPC_PORT":                   config.GRPCPort,
		"CHURNER_EIGENDA_DIRECTORY":           config.EigenDADirectory,
		"CHURNER_BLS_OPERATOR_STATE_RETRIVER": config.OperatorStateRetriever,
		"CHURNER_EIGENDA_SERVICE_MANAGER":     config.ServiceManager,
		"CHURNER_CHAIN_RPC":                   config.ChainRPC,
		"CHURNER_PRIVATE_KEY":                 privateKey,
		"CHURNER_INDEXER_PULL_INTERVAL":       config.IndexerPullInterval,
		"CHURNER_ENABLE_METRICS":              fmt.Sprintf("%t", config.EnableMetrics),
		"CHURNER_METRICS_HTTP_PORT":           config.MetricsHTTPPort,
		"CHURNER_CHURN_APPROVAL_INTERVAL":     config.ChurnApprovalInterval,
	}

	// Add graph URL if provided
	if config.GraphURL != "" {
		env["CHURNER_GRAPH_URL"] = config.GraphURL
	}

	return env
}
