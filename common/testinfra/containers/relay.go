package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RelayConfig defines configuration for an EigenDA relay node
type RelayConfig struct {
	// Relay identification
	ID       int    // Relay index (0-based)
	Hostname string // Hostname for the relay

	// Network ports
	GRPCPort         string
	InternalGRPCPort string
	MetricsPort      string

	// AWS configuration
	AWSEndpointURL     string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Storage configuration
	BucketName        string
	MetadataTableName string

	// Relay keys (which relay indices this relay serves)
	RelayKeys []int

	// Chain configuration
	ChainRPC string
	GraphURL string

	// Contract addresses
	BLSOperatorStateRetriever string
	EigenDAServiceManager     string
	EigenDADirectory          string

	// Cache configurations
	MetadataCacheSize        int
	BlobCacheBytes           uint64
	ChunkCacheBytes          int64
	AuthenticationKeyMaxSize int

	// Rate limiting configurations
	MaxGRPCMessageSize             int
	MaxConcurrentGetBlobOps        int
	MaxConcurrentGetChunkOps       int
	MaxConcurrentGetChunkOpsClient int
	MaxKeysPerGetChunksRequest     int

	// GetBlob rate limiting
	MaxGetBlobOpsPerSecond   float64
	GetBlobOpsBurstiness     int
	MaxGetBlobBytesPerSecond float64
	GetBlobBytesBurstiness   int

	// GetChunk rate limiting (global)
	MaxGetChunkOpsPerSecond   float64
	GetChunkOpsBurstiness     int
	MaxGetChunkBytesPerSecond float64
	GetChunkBytesBurstiness   int

	// GetChunk rate limiting (per client)
	MaxGetChunkOpsPerSecondClient   float64
	GetChunkOpsBurstinessClient     int
	MaxGetChunkBytesPerSecondClient float64
	GetChunkBytesBurstinessClient   int

	// Timeouts
	AuthenticationTimeout          time.Duration
	GetChunksTimeout               time.Duration
	GetBlobTimeout                 time.Duration
	InternalGetMetadataTimeout     time.Duration
	InternalGetBlobTimeout         time.Duration
	InternalGetProofsTimeout       time.Duration
	InternalGetCoefficientsTimeout time.Duration
	OnchainStateRefreshInterval    time.Duration

	// Feature flags
	AuthenticationDisabled bool
	EnableMetrics          bool

	// Logging
	LogFormat string
	LogLevel  string

	// Container image
	Image string
}

// RelayContainer represents a running relay container
type RelayContainer struct {
	testcontainers.Container
	config  RelayConfig
	logPath string
}

// DefaultRelayConfig returns a default relay configuration
func DefaultRelayConfig(id int) RelayConfig {
	// Calculate base ports for this relay
	// Each relay gets a range of 2 ports for different services
	// Starting at 34000 to avoid conflicts with other services
	basePort := 34000 + (id * 2)

	return RelayConfig{
		ID:               id,
		Hostname:         "0.0.0.0",
		GRPCPort:         fmt.Sprintf("%d", basePort),
		InternalGRPCPort: fmt.Sprintf("%d", basePort),
		MetricsPort:      fmt.Sprintf("%d", basePort+1),

		// AWS defaults (will be overridden by LocalStack if enabled)
		AWSRegion: "us-east-1",

		// Storage defaults
		BucketName:        "test-eigenda-blobstore",
		MetadataTableName: "test-BlobMetadata-v2",

		// Relay keys - by default each relay serves its own index
		RelayKeys: []int{id},

		// Cache configurations
		MetadataCacheSize:        1048576,    // 1 MiB items
		BlobCacheBytes:           1073741824, // 1 GiB
		ChunkCacheBytes:          1073741824, // 1 GiB
		AuthenticationKeyMaxSize: 1048576,

		// Rate limiting configurations
		MaxGRPCMessageSize:             4194304, // 4 MiB
		MaxConcurrentGetBlobOps:        1024,
		MaxConcurrentGetChunkOps:       1024,
		MaxConcurrentGetChunkOpsClient: 1,
		MaxKeysPerGetChunksRequest:     1024,

		// GetBlob rate limiting
		MaxGetBlobOpsPerSecond:   1024,
		GetBlobOpsBurstiness:     1024,
		MaxGetBlobBytesPerSecond: 20971520, // 20 MiB/s
		GetBlobBytesBurstiness:   20971520, // 20 MiB

		// GetChunk rate limiting (global)
		MaxGetChunkOpsPerSecond:   1024,
		GetChunkOpsBurstiness:     1024,
		MaxGetChunkBytesPerSecond: 83886080,  // 80 MiB/s
		GetChunkBytesBurstiness:   838860800, // 800 MiB

		// GetChunk rate limiting (per client)
		MaxGetChunkOpsPerSecondClient:   8,
		GetChunkOpsBurstinessClient:     8,
		MaxGetChunkBytesPerSecondClient: 41943040,  // 40 MiB/s
		GetChunkBytesBurstinessClient:   419430400, // 400 MiB

		// Timeouts
		AuthenticationTimeout:          0,
		GetChunksTimeout:               20 * time.Second,
		GetBlobTimeout:                 20 * time.Second,
		InternalGetMetadataTimeout:     5 * time.Second,
		InternalGetBlobTimeout:         20 * time.Second,
		InternalGetProofsTimeout:       5 * time.Second,
		InternalGetCoefficientsTimeout: 20 * time.Second,
		OnchainStateRefreshInterval:    1 * time.Hour,

		// Feature flags
		AuthenticationDisabled: false,
		EnableMetrics:          true,

		// Logging
		LogFormat: "text",
		LogLevel:  "debug",

		// Container image
		Image: "ghcr.io/layr-labs/eigenda/relay:dev",
	}
}

// NewRelayContainerWithNetwork creates and starts a new relay container with a custom network
func NewRelayContainerWithNetwork(ctx context.Context, config RelayConfig, nw *testcontainers.DockerNetwork) (*RelayContainer, error) {
	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", fmt.Sprintf("relay-%d-logs-*", config.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "relay.log")

	// Build environment variables
	env := buildRelayEnv(config)

	// Debug log critical environment variables
	fmt.Printf("DEBUG: Relay %d env - RELAY_GRPC_PORT=%s, RELAY_EIGENDA_SERVICE_MANAGER=%s\n",
		config.ID, env["RELAY_GRPC_PORT"], env["RELAY_EIGEN_DA_SERVICE_MANAGER"])

	// Configure container request
	req := testcontainers.ContainerRequest{
		Image: config.Image,
		Env:   env,
		ExposedPorts: []string{
			config.InternalGRPCPort + "/tcp",
			config.MetricsPort + "/tcp",
		},
		Networks:       []string{},
		NetworkAliases: map[string][]string{},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: logDir,
				},
				Target: "/logs",
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.InternalGRPCPort+"/tcp")).WithStartupTimeout(60*time.Second),
			wait.ForLog("GRPC Listening").WithStartupTimeout(90*time.Second),
		),
		Name:            fmt.Sprintf("eigenda-relay-%d", config.ID),
		AlwaysPullImage: false, // Use local image if available
	}

	// Add port bindings when hostname is 0.0.0.0 or using localhost domain
	if config.Hostname == "0.0.0.0" || strings.HasSuffix(config.Hostname, ".localtest.me") {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port(config.InternalGRPCPort + "/tcp"): []nat.PortBinding{{HostPort: config.GRPCPort}},
				nat.Port(config.MetricsPort + "/tcp"):      []nat.PortBinding{{HostPort: config.MetricsPort}},
			}
		}
	}

	// Add network configuration if provided
	if nw != nil {
		req.Networks = []string{nw.Name}
		// Use localhost domain which resolves to 127.0.0.1
		relayHostname := fmt.Sprintf("relay-%d.localtest.me", config.ID)
		req.NetworkAliases = map[string][]string{
			nw.Name: {relayHostname, config.Hostname, fmt.Sprintf("relay-%d", config.ID), fmt.Sprintf("relay%d", config.ID)},
		}
	}

	// Create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start relay container %d: %w", config.ID, err)
	}

	fmt.Printf("Relay %d logs will be available at: %s\n", config.ID, logPath)

	return &RelayContainer{
		Container: container,
		config:    config,
		logPath:   logPath,
	}, nil
}

// Config returns the relay configuration
func (c *RelayContainer) Config() RelayConfig {
	return c.config
}

// LogPath returns the path to the relay log file on the host
func (c *RelayContainer) LogPath() string {
	return c.logPath
}

// GetGRPCAddress returns the full gRPC address for this relay (hostname:port)
func (c *RelayContainer) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%s", c.config.Hostname, c.config.GRPCPort)
}

// GetInternalGRPCAddress returns the internal gRPC address for Docker network communication
func (c *RelayContainer) GetInternalGRPCAddress() string {
	// Use the container alias for internal communication
	return fmt.Sprintf("relay%d:%s", c.config.ID, c.config.InternalGRPCPort)
}

// buildRelayEnv builds environment variables for the relay container
func buildRelayEnv(config RelayConfig) map[string]string {
	// Format relay keys as comma-separated string
	relayKeys := make([]string, len(config.RelayKeys))
	for i, key := range config.RelayKeys {
		relayKeys[i] = fmt.Sprintf("%d", key)
	}
	relayKeysStr := strings.Join(relayKeys, ",")

	// Ensure AWS region has a default value
	awsRegion := config.AWSRegion
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	env := map[string]string{
		// Network configuration
		"RELAY_GRPC_PORT": config.GRPCPort,

		// AWS configuration (RELAY_ prefix is required for the relay binary)
		"RELAY_AWS_ENDPOINT_URL":      config.AWSEndpointURL,
		"RELAY_AWS_REGION":            awsRegion,
		"RELAY_AWS_ACCESS_KEY_ID":     config.AWSAccessKeyID,
		"RELAY_AWS_SECRET_ACCESS_KEY": config.AWSSecretAccessKey,

		// Storage configuration
		"RELAY_BUCKET_NAME":         config.BucketName,
		"RELAY_METADATA_TABLE_NAME": config.MetadataTableName,
		"RELAY_RELAY_KEYS":          relayKeysStr,

		// Chain configuration
		"RELAY_CHAIN_RPC":   config.ChainRPC,
		"RELAY_GRAPH_URL":   config.GraphURL,
		"RELAY_PRIVATE_KEY": "123", // Dummy private key as it's not used but required by the binary

		// Contract addresses (using _ADDR suffix as expected by relay flags)
		"RELAY_BLS_OPERATOR_STATE_RETRIEVER_ADDR": config.BLSOperatorStateRetriever,
		"RELAY_EIGEN_DA_SERVICE_MANAGER_ADDR":     config.EigenDAServiceManager,

		// Cache configurations
		"RELAY_METADATA_CACHE_SIZE":           fmt.Sprintf("%d", config.MetadataCacheSize),
		"RELAY_BLOB_CACHE_BYTES":              fmt.Sprintf("%d", config.BlobCacheBytes),
		"RELAY_CHUNK_CACHE_BYTES":             fmt.Sprintf("%d", config.ChunkCacheBytes),
		"RELAY_AUTHENTICATION_KEY_CACHE_SIZE": fmt.Sprintf("%d", config.AuthenticationKeyMaxSize),

		// Rate limiting configurations
		"RELAY_MAX_GRPC_MESSAGE_SIZE":               fmt.Sprintf("%d", config.MaxGRPCMessageSize),
		"RELAY_MAX_CONCURRENT_GET_BLOB_OPS":         fmt.Sprintf("%d", config.MaxConcurrentGetBlobOps),
		"RELAY_MAX_CONCURRENT_GET_CHUNK_OPS":        fmt.Sprintf("%d", config.MaxConcurrentGetChunkOps),
		"RELAY_MAX_CONCURRENT_GET_CHUNK_OPS_CLIENT": fmt.Sprintf("%d", config.MaxConcurrentGetChunkOpsClient),
		"RELAY_MAX_KEYS_PER_GET_CHUNKS_REQUEST":     fmt.Sprintf("%d", config.MaxKeysPerGetChunksRequest),

		// GetBlob rate limiting
		"RELAY_MAX_GET_BLOB_OPS_PER_SECOND":   fmt.Sprintf("%f", config.MaxGetBlobOpsPerSecond),
		"RELAY_GET_BLOB_OPS_BURSTINESS":       fmt.Sprintf("%d", config.GetBlobOpsBurstiness),
		"RELAY_MAX_GET_BLOB_BYTES_PER_SECOND": fmt.Sprintf("%f", config.MaxGetBlobBytesPerSecond),
		"RELAY_GET_BLOB_BYTES_BURSTINESS":     fmt.Sprintf("%d", config.GetBlobBytesBurstiness),

		// GetChunk rate limiting (global)
		"RELAY_MAX_GET_CHUNK_OPS_PER_SECOND":   fmt.Sprintf("%f", config.MaxGetChunkOpsPerSecond),
		"RELAY_GET_CHUNK_OPS_BURSTINESS":       fmt.Sprintf("%d", config.GetChunkOpsBurstiness),
		"RELAY_MAX_GET_CHUNK_BYTES_PER_SECOND": fmt.Sprintf("%f", config.MaxGetChunkBytesPerSecond),
		"RELAY_GET_CHUNK_BYTES_BURSTINESS":     fmt.Sprintf("%d", config.GetChunkBytesBurstiness),

		// GetChunk rate limiting (per client)
		"RELAY_MAX_GET_CHUNK_OPS_PER_SECOND_CLIENT":   fmt.Sprintf("%f", config.MaxGetChunkOpsPerSecondClient),
		"RELAY_GET_CHUNK_OPS_BURSTINESS_CLIENT":       fmt.Sprintf("%d", config.GetChunkOpsBurstinessClient),
		"RELAY_MAX_GET_CHUNK_BYTES_PER_SECOND_CLIENT": fmt.Sprintf("%f", config.MaxGetChunkBytesPerSecondClient),
		"RELAY_GET_CHUNK_BYTES_BURSTINESS_CLIENT":     fmt.Sprintf("%d", config.GetChunkBytesBurstinessClient),

		// Timeouts
		"RELAY_AUTHENTICATION_TIMEOUT":            fmt.Sprintf("%s", config.AuthenticationTimeout),
		"RELAY_GET_CHUNKS_TIMEOUT":                fmt.Sprintf("%s", config.GetChunksTimeout),
		"RELAY_GET_BLOB_TIMEOUT":                  fmt.Sprintf("%s", config.GetBlobTimeout),
		"RELAY_INTERNAL_GET_METADATA_TIMEOUT":     fmt.Sprintf("%s", config.InternalGetMetadataTimeout),
		"RELAY_INTERNAL_GET_BLOB_TIMEOUT":         fmt.Sprintf("%s", config.InternalGetBlobTimeout),
		"RELAY_INTERNAL_GET_PROOFS_TIMEOUT":       fmt.Sprintf("%s", config.InternalGetProofsTimeout),
		"RELAY_INTERNAL_GET_COEFFICIENTS_TIMEOUT": fmt.Sprintf("%s", config.InternalGetCoefficientsTimeout),
		"RELAY_ONCHAIN_STATE_REFRESH_INTERVAL":    fmt.Sprintf("%s", config.OnchainStateRefreshInterval),

		// Feature flags
		"RELAY_AUTHENTICATION_DISABLED": fmt.Sprintf("%t", config.AuthenticationDisabled),
		"RELAY_ENABLE_METRICS":          fmt.Sprintf("%t", config.EnableMetrics),
		"RELAY_METRICS_PORT":            config.MetricsPort,

		// Logging
		"RELAY_LOG_FORMAT": config.LogFormat,
		"RELAY_LOG_LEVEL":  config.LogLevel,
		"RELAY_LOG_PATH":   filepath.Join("/logs", "relay.log"),

		// Metadata and blob max concurrency
		"RELAY_METADATA_MAX_CONCURRENCY": "32",
		"RELAY_BLOB_MAX_CONCURRENCY":     "32",
		"RELAY_CHUNK_MAX_CONCURRENCY":    "32",
	}

	// Add EigenDA directory if configured (takes precedence over individual contracts)
	if config.EigenDADirectory != "" {
		env["RELAY_EIGENDA_DIRECTORY"] = config.EigenDADirectory
	}

	return env
}
