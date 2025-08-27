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

// DisperserConfig defines configuration for the disperser container
type DisperserConfig struct {
	// Enable disperser service
	Enabled bool

	// Version of disperser (1 or 2)
	Version int

	// Log format (text/json)
	LogFormat string

	// S3 and DynamoDB configuration
	S3BucketName        string
	DynamoDBTableName   string
	RateBucketTableName string
	RateBucketStoreSize string

	// Network configuration
	GRPCPort string

	// Metrics configuration
	EnableMetrics   bool
	MetricsHTTPPort string

	// Chain configuration
	ChainRPC         string
	PrivateKey       string
	NumConfirmations string

	// Rate limiting configuration
	RegisteredQuorumID    string
	TotalUnauthByteRate   string
	PerUserUnauthByteRate string
	TotalUnauthBlobRate   string
	PerUserUnauthBlobRate string
	EnableRatelimiter     bool
	RetrievalBlobRate     string
	RetrievalByteRate     string
	BucketSizes           string
	BucketMultipliers     string
	CountFailed           bool

	// EigenDA configuration
	EigenDADirectory       string
	OperatorStateRetriever string
	ServiceManager         string

	// V2 specific configuration
	EnablePaymentMeterer  bool
	ReservedOnly          bool
	ReservationsTableName string
	OnDemandTableName     string
	GlobalRateTableName   string

	// Encoder configuration
	EncoderAddress string

	// KZG parameters (required for v2)
	G1Path           string
	G2Path           string
	G2PowerOf2Path   string
	CachePath        string
	SRSOrder         string
	SRSLoad          string

	// Container image
	Image string

	// AWS configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSEndpointURL     string
}

// DisperserContainer represents a running disperser container
type DisperserContainer struct {
	testcontainers.Container
	config  DisperserConfig
	url     string
	logPath string // Path to log file on host
}

// DefaultDisperserConfig returns a default disperser configuration
func DefaultDisperserConfig(version int) DisperserConfig {
	config := DisperserConfig{
		Enabled:                true,
		Version:                version,
		LogFormat:              "text",
		S3BucketName:           "test-eigenda-blobstore",
		RateBucketStoreSize:    "100000",
		GRPCPort:               "32003",
		EnableMetrics:          true,
		MetricsHTTPPort:        "9093",
		ChainRPC:               "", // Will be populated from Anvil
		PrivateKey:             "123",
		NumConfirmations:       "0",
		RegisteredQuorumID:     "0,1",
		TotalUnauthByteRate:    "10000000,10000000",
		PerUserUnauthByteRate:  "32000,32000",
		TotalUnauthBlobRate:    "10,10",
		PerUserUnauthBlobRate:  "2,2",
		EnableRatelimiter:      true,
		RetrievalBlobRate:      "4",
		RetrievalByteRate:      "10000000",
		BucketSizes:            "5s",
		BucketMultipliers:      "1",
		CountFailed:            true,
		EigenDADirectory:       "", // Will be populated from contract deployment
		OperatorStateRetriever: "", // Will be populated from contract deployment
		ServiceManager:         "", // Will be populated from contract deployment
		Image:                  "ghcr.io/layr-labs/eigenda/apiserver:dev",
		AWSRegion:              "us-east-1",
		AWSAccessKeyID:         "localstack",
		AWSSecretAccessKey:     "localstack",
		AWSEndpointURL:         "", // Will be populated from LocalStack
		EncoderAddress:         "", // Will be populated from encoder
	}

	// V1 vs V2 specific defaults
	if version == 1 {
		config.DynamoDBTableName = "test-BlobMetadata"
		config.RateBucketTableName = ""
	} else if version == 2 {
		config.DynamoDBTableName = "test-BlobMetadata-v2"
		config.RateBucketTableName = ""
		config.EnablePaymentMeterer = true
		config.ReservedOnly = false
		config.ReservationsTableName = "e2e_v2_reservation"
		config.OnDemandTableName = "e2e_v2_ondemand"
		config.GlobalRateTableName = "e2e_v2_global_reservation"
	}

	return config
}

// NewDisperserContainerWithNetwork creates and starts a new disperser container with a custom network
func NewDisperserContainerWithNetwork(ctx context.Context, config DisperserConfig, network *testcontainers.DockerNetwork) (*DisperserContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", fmt.Sprintf("disperser-v%d-logs-*", config.Version))
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "disperser.log")

	// Create a copy of config to modify paths for container
	containerConfig := config

	// Prepare mounts
	mounts := testcontainers.ContainerMounts{
		{
			Source: testcontainers.GenericBindMountSource{
				HostPath: logDir,
			},
			Target: "/logs",
		},
	}

	// Mount KZG resources if paths are provided (required for v2)
	if config.G1Path != "" {
		// Extract directory from G1 path to mount the entire KZG directory
		kzgDir := filepath.Dir(config.G1Path)
		mounts = append(mounts, testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: kzgDir,
			},
			Target:   "/kzg",
			ReadOnly: true,
		})
		// Update config copy to use mounted paths
		g1Filename := filepath.Base(config.G1Path)
		containerConfig.G1Path = "/kzg/" + g1Filename

		if config.G2Path != "" {
			g2Filename := filepath.Base(config.G2Path)
			containerConfig.G2Path = "/kzg/" + g2Filename
		}

		if config.G2PowerOf2Path != "" {
			g2PowerOf2Filename := filepath.Base(config.G2PowerOf2Path)
			containerConfig.G2PowerOf2Path = "/kzg/" + g2PowerOf2Filename
		}
	}

	// Mount cache directory if provided
	if config.CachePath != "" {
		mounts = append(mounts, testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: config.CachePath,
			},
			Target: "/cache",
		})
		// Update config copy to use mounted path
		containerConfig.CachePath = "/cache"
	}

	// Build environment variables with updated paths
	env := buildDisperserEnv(containerConfig)

	// Configure container request with network
	alias := fmt.Sprintf("disperser-v%d", config.Version)
	if config.Version == 1 {
		alias = "dis0"
	} else if config.Version == 2 {
		alias = "dis1"
	}

	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{config.GRPCPort + "/tcp"},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {alias},
		},
		Mounts: mounts,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.GRPCPort+"/tcp")).WithStartupTimeout(60*time.Second),
			wait.ForLog("GRPC Listening").WithStartupTimeout(60*time.Second),
		),
		Name:            fmt.Sprintf("eigenda-disperser-v%d", config.Version),
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
		return nil, fmt.Errorf("failed to start disperser v%d container: %w", config.Version, err)
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

	disperserURL := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	fmt.Printf("Disperser v%d logs will be available at: %s\n", config.Version, logPath)

	return &DisperserContainer{
		Container: container,
		config:    config,
		url:       disperserURL,
		logPath:   logPath,
	}, nil
}

// URL returns the disperser service URL
func (c *DisperserContainer) URL() string {
	return c.url
}

// InternalURL returns the disperser service URL for internal network communication
func (c *DisperserContainer) InternalURL() string {
	if c.config.Version == 1 {
		return fmt.Sprintf("dis0:%s", c.config.GRPCPort)
	}
	return fmt.Sprintf("dis1:%s", c.config.GRPCPort)
}

// Config returns the disperser configuration
func (c *DisperserContainer) Config() DisperserConfig {
	return c.config
}

// LogPath returns the path to the disperser log file on the host
func (c *DisperserContainer) LogPath() string {
	return c.logPath
}

// Version returns the disperser version (1 or 2)
func (c *DisperserContainer) Version() int {
	return c.config.Version
}

// buildDisperserEnv builds environment variables for the disperser container
func buildDisperserEnv(config DisperserConfig) map[string]string {
	// Strip 0x prefix from private key if present
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")

	env := map[string]string{
		"DISPERSER_SERVER_LOG_FORMAT":                  config.LogFormat,
		"DISPERSER_SERVER_LOG_PATH":                    "/logs/disperser.log", // Log to a file we can access
		"DISPERSER_SERVER_LOG_LEVEL":                   "debug",               // Enable debug logging
		"DISPERSER_SERVER_S3_BUCKET_NAME":              config.S3BucketName,
		"DISPERSER_SERVER_DYNAMODB_TABLE_NAME":         config.DynamoDBTableName,
		"DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME":      config.RateBucketTableName,
		"DISPERSER_SERVER_RATE_BUCKET_STORE_SIZE":      config.RateBucketStoreSize,
		"DISPERSER_SERVER_GRPC_PORT":                   config.GRPCPort,
		"DISPERSER_SERVER_ENABLE_METRICS":              fmt.Sprintf("%t", config.EnableMetrics),
		"DISPERSER_SERVER_METRICS_HTTP_PORT":           config.MetricsHTTPPort,
		"DISPERSER_SERVER_CHAIN_RPC":                   config.ChainRPC,
		"DISPERSER_SERVER_PRIVATE_KEY":                 privateKey,
		"DISPERSER_SERVER_NUM_CONFIRMATIONS":           config.NumConfirmations,
		"DISPERSER_SERVER_REGISTERED_QUORUM_ID":        config.RegisteredQuorumID,
		"DISPERSER_SERVER_TOTAL_UNAUTH_BYTE_RATE":      config.TotalUnauthByteRate,
		"DISPERSER_SERVER_PER_USER_UNAUTH_BYTE_RATE":   config.PerUserUnauthByteRate,
		"DISPERSER_SERVER_TOTAL_UNAUTH_BLOB_RATE":      config.TotalUnauthBlobRate,
		"DISPERSER_SERVER_PER_USER_UNAUTH_BLOB_RATE":   config.PerUserUnauthBlobRate,
		"DISPERSER_SERVER_ENABLE_RATELIMITER":          fmt.Sprintf("%t", config.EnableRatelimiter),
		"DISPERSER_SERVER_RETRIEVAL_BLOB_RATE":         config.RetrievalBlobRate,
		"DISPERSER_SERVER_RETRIEVAL_BYTE_RATE":         config.RetrievalByteRate,
		"DISPERSER_SERVER_BUCKET_SIZES":                config.BucketSizes,
		"DISPERSER_SERVER_BUCKET_MULTIPLIERS":          config.BucketMultipliers,
		"DISPERSER_SERVER_COUNT_FAILED":                fmt.Sprintf("%t", config.CountFailed),
		"DISPERSER_SERVER_EIGENDA_DIRECTORY":           config.EigenDADirectory,
		"DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER": config.OperatorStateRetriever,
		"DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER":     config.ServiceManager,
		"DISPERSER_SERVER_AWS_REGION":                  config.AWSRegion,
		"DISPERSER_SERVER_AWS_ACCESS_KEY_ID":           config.AWSAccessKeyID,
		"DISPERSER_SERVER_AWS_SECRET_ACCESS_KEY":       config.AWSSecretAccessKey,
	}

	// Add AWS endpoint URL if provided
	if config.AWSEndpointURL != "" {
		env["DISPERSER_SERVER_AWS_ENDPOINT_URL"] = config.AWSEndpointURL
	}

	// Add encoder address if provided
	if config.EncoderAddress != "" {
		env["DISPERSER_SERVER_ENCODER_ADDRESS"] = config.EncoderAddress
	}

	// Add version-specific environment variables
	if config.Version == 2 {
		env["DISPERSER_SERVER_DISPERSER_VERSION"] = "2"
		env["DISPERSER_SERVER_ENABLE_PAYMENT_METERER"] = fmt.Sprintf("%t", config.EnablePaymentMeterer)
		env["DISPERSER_SERVER_RESERVED_ONLY"] = fmt.Sprintf("%t", config.ReservedOnly)
		env["DISPERSER_SERVER_RESERVATIONS_TABLE_NAME"] = config.ReservationsTableName
		env["DISPERSER_SERVER_ON_DEMAND_TABLE_NAME"] = config.OnDemandTableName
		env["DISPERSER_SERVER_GLOBAL_RATE_TABLE_NAME"] = config.GlobalRateTableName
	}

	// Add KZG parameters if provided (required for v2)
	if config.G1Path != "" {
		env["DISPERSER_SERVER_G1_PATH"] = config.G1Path
	}
	if config.G2Path != "" {
		env["DISPERSER_SERVER_G2_PATH"] = config.G2Path
	}
	if config.G2PowerOf2Path != "" {
		env["DISPERSER_SERVER_G2_POWER_OF_2_PATH"] = config.G2PowerOf2Path
	}
	if config.CachePath != "" {
		env["DISPERSER_SERVER_CACHE_PATH"] = config.CachePath
	}
	if config.SRSOrder != "" {
		env["DISPERSER_SERVER_SRS_ORDER"] = config.SRSOrder
	}
	if config.SRSLoad != "" {
		env["DISPERSER_SERVER_SRS_LOAD"] = config.SRSLoad
	}

	return env
}
