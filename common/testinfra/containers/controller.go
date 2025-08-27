package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ControllerConfig defines configuration for the controller container
type ControllerConfig struct {
	// Enable controller service
	Enabled bool

	// Log configuration
	LogFormat string
	LogLevel  string

	// Network configuration
	MetricsPort string

	// AWS configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSEndpointURL     string

	// DynamoDB configuration
	DynamoDBTableName string

	// Encoder configuration
	EncoderAddress string

	// Chain configuration
	ChainRPC                  string
	ChainPrivateKey           string // Private key for signing transactions
	EigenDAServiceManager     string
	BLSOperatorStateRetriever string
	NumConfirmations          uint
	FinalizationBlockDelay    uint

	// Graph configuration
	UseGraph bool
	GraphURL string

	// Indexer configuration (if not using graph)
	IndexerPullInterval string

	// Dispatcher configuration
	DispatcherPullInterval string
	AvailableRelays        string

	// Encoding configuration
	EncodingPullInterval string

	// Attestation configuration
	AttestationTimeout      string
	BatchAttestationTimeout string

	// Disperser configuration
	DisperserKMSKeyID                   string
	DisperserStoreChunksSigningDisabled bool

	// Performance configuration
	NumConcurrentEncodingRequests  int
	NumConcurrentDispersalRequests int

	// Container image
	Image string
}

// ControllerContainer represents a running controller container
type ControllerContainer struct {
	testcontainers.Container
	config     ControllerConfig
	metricsURL string
	logPath    string // Path to log file on host
}

// DefaultControllerConfig returns a default controller configuration
func DefaultControllerConfig() ControllerConfig {
	return ControllerConfig{
		Enabled:                             true,
		LogFormat:                           "text",
		LogLevel:                            "debug",
		MetricsPort:                         "9100",
		AWSRegion:                           "us-east-1",
		AWSAccessKeyID:                      "localstack",
		AWSSecretAccessKey:                  "localstack",
		AWSEndpointURL:                      "",
		DynamoDBTableName:                   "test-BlobMetadata-v2",
		EncoderAddress:                      "encoder-v2:34001", // Default to v2 encoder
		ChainRPC:                            "",
		ChainPrivateKey:                     "", // Will be set from deployment
		EigenDAServiceManager:               "",
		BLSOperatorStateRetriever:           "",
		NumConfirmations:                    0,
		FinalizationBlockDelay:              5,
		UseGraph:                            false,
		GraphURL:                            "",
		IndexerPullInterval:                 "1s",
		DispatcherPullInterval:              "3s",
		AvailableRelays:                     "0,1,2,3",
		EncodingPullInterval:                "1s",
		AttestationTimeout:                  "5s",
		BatchAttestationTimeout:             "6s",
		DisperserKMSKeyID:                   "",
		DisperserStoreChunksSigningDisabled: false,
		NumConcurrentEncodingRequests:       16,
		NumConcurrentDispersalRequests:      8,
		Image:                               "ghcr.io/layr-labs/eigenda/controller:dev",
	}
}

// NewControllerContainerWithNetwork creates and starts a new controller container with a custom network
func NewControllerContainerWithNetwork(ctx context.Context, config ControllerConfig, network *testcontainers.DockerNetwork) (*ControllerContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", "controller-logs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "controller.log")

	// Prepare mounts
	mounts := testcontainers.ContainerMounts{
		{
			Source: testcontainers.GenericBindMountSource{
				HostPath: logDir,
			},
			Target: "/logs",
		},
	}

	// Build environment variables
	env := buildControllerEnv(config)

	// Configure container request with network
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{config.MetricsPort + "/tcp"},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {"controller"},
		},
		// Add host access configuration for accessing operators on host ports
		HostAccessPorts: []int{
			32011, 32012, 32013, 32014, 32015, 32016, // operator-0 ports
			32021, 32022, 32023, 32024, 32025, 32026, // operator-1 ports
			32031, 32032, 32033, 32034, 32035, 32036, // operator-2 ports
			32041, 32042, 32043, 32044, 32045, 32046, // operator-3 ports
		},
		Mounts: mounts,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.MetricsPort+"/tcp")).WithStartupTimeout(30*time.Second),
			wait.ForLog("no blobs to encode").WithStartupTimeout(30*time.Second),
		),
		Name:            "eigenda-controller",
		AlwaysPullImage: false, // Use local image if available
	}

	// Add HostConfigModifier to set up ExtraHosts for operator localhost domains
	// This maps operator-{i}.localtest.me to host-gateway, allowing the controller
	// to reach operators running on the host through the localhost domain
	req.HostConfigModifier = func(hc *container.HostConfig) {
		// Add entries for each operator's localhost domain
		// host-gateway is a special Docker hostname that resolves to the host machine
		for i := 0; i < 4; i++ {
			operatorHost := fmt.Sprintf("operator-%d.localtest.me:host-gateway", i)
			hc.ExtraHosts = append(hc.ExtraHosts, operatorHost)
		}
	}

	// Create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start controller container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, nat.Port(config.MetricsPort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the container host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	metricsURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	fmt.Printf("Controller logs will be available at: %s\n", logPath)

	return &ControllerContainer{
		Container:  container,
		config:     config,
		metricsURL: metricsURL,
		logPath:    logPath,
	}, nil
}

// MetricsURL returns the controller metrics URL
func (c *ControllerContainer) MetricsURL() string {
	return c.metricsURL
}

// InternalMetricsURL returns the controller metrics URL for internal network communication
func (c *ControllerContainer) InternalMetricsURL() string {
	return fmt.Sprintf("http://controller:%s", c.config.MetricsPort)
}

// Config returns the controller configuration
func (c *ControllerContainer) Config() ControllerConfig {
	return c.config
}

// LogPath returns the path to the controller log file on the host
func (c *ControllerContainer) LogPath() string {
	return c.logPath
}

// buildControllerEnv builds environment variables for the controller container
func buildControllerEnv(config ControllerConfig) map[string]string {
	env := map[string]string{
		"CONTROLLER_LOG_FORMAT":                              config.LogFormat,
		"CONTROLLER_LOG_LEVEL":                               config.LogLevel,
		"CONTROLLER_LOG_PATH":                                "/logs/controller.log",
		"CONTROLLER_METRICS_PORT":                            config.MetricsPort,
		"CONTROLLER_AWS_REGION":                              config.AWSRegion,
		"CONTROLLER_AWS_ACCESS_KEY_ID":                       config.AWSAccessKeyID,
		"CONTROLLER_AWS_SECRET_ACCESS_KEY":                   config.AWSSecretAccessKey,
		"CONTROLLER_DYNAMODB_TABLE_NAME":                     config.DynamoDBTableName,
		"CONTROLLER_ENCODER_ADDRESS":                         config.EncoderAddress,
		"CONTROLLER_CHAIN_RPC":                               config.ChainRPC,
		"CONTROLLER_PRIVATE_KEY":                             config.ChainPrivateKey,
		"CONTROLLER_EIGENDA_SERVICE_MANAGER":                 config.EigenDAServiceManager,
		"CONTROLLER_BLS_OPERATOR_STATE_RETRIVER":             config.BLSOperatorStateRetriever,
		"CONTROLLER_NUM_CONFIRMATIONS":                       fmt.Sprintf("%d", config.NumConfirmations),
		"CONTROLLER_FINALIZATION_BLOCK_DELAY":                fmt.Sprintf("%d", config.FinalizationBlockDelay),
		"CONTROLLER_USE_GRAPH":                               fmt.Sprintf("%t", config.UseGraph),
		"CONTROLLER_INDEXER_PULL_INTERVAL":                   config.IndexerPullInterval,
		"CONTROLLER_DISPATCHER_PULL_INTERVAL":                config.DispatcherPullInterval,
		"CONTROLLER_AVAILABLE_RELAYS":                        config.AvailableRelays,
		"CONTROLLER_ENCODING_PULL_INTERVAL":                  config.EncodingPullInterval,
		"CONTROLLER_ATTESTATION_TIMEOUT":                     config.AttestationTimeout,
		"CONTROLLER_BATCH_ATTESTATION_TIMEOUT":               config.BatchAttestationTimeout,
		"CONTROLLER_DISPERSER_STORE_CHUNKS_SIGNING_DISABLED": fmt.Sprintf("%t", config.DisperserStoreChunksSigningDisabled),
		"CONTROLLER_NUM_CONCURRENT_ENCODING_REQUESTS":        fmt.Sprintf("%d", config.NumConcurrentEncodingRequests),
		"CONTROLLER_NUM_CONCURRENT_DISPERSAL_REQUESTS":       fmt.Sprintf("%d", config.NumConcurrentDispersalRequests),
	}

	// Add optional configurations if provided
	if config.AWSEndpointURL != "" {
		env["CONTROLLER_AWS_ENDPOINT_URL"] = config.AWSEndpointURL
	}
	if config.GraphURL != "" {
		env["CONTROLLER_GRAPH_URL"] = config.GraphURL
	}
	if config.DisperserKMSKeyID != "" {
		env["CONTROLLER_DISPERSER_KMS_KEY_ID"] = config.DisperserKMSKeyID
	}
	// If private key is empty, set a dummy value to satisfy the flag requirement
	if config.ChainPrivateKey == "" {
		env["CONTROLLER_PRIVATE_KEY"] = "0x0000000000000000000000000000000000000000000000000000000000000001"
	}

	return env
}
