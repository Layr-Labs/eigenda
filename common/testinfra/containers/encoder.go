package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// EncoderConfig defines configuration for the encoder container
type EncoderConfig struct {
	// Enable encoder service
	Enabled bool

	// Log configuration
	LogFormat string
	LogLevel  string

	// Network configuration
	GRPCPort string

	// AWS configuration (for S3 storage if using encoder v2)
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSEndpointURL     string
	S3BucketName       string // Only for encoder v2

	// Metrics configuration
	EnableMetrics   bool
	MetricsHTTPPort string

	// SRS configuration
	G1Path         string
	G2Path         string
	G2PowerOf2Path string
	SRSOrder       string
	SRSLoad        string
	CachePath      string

	// Performance configuration
	Verbose               string
	NumWorkers            string
	MaxConcurrentRequests string
	RequestPoolSize       string
	RequestQueueSize      string

	// Encoder version (1 or 2)
	EncoderVersion string

	// Container image
	Image string
}

// EncoderContainer represents a running encoder container
type EncoderContainer struct {
	testcontainers.Container
	config  EncoderConfig
	url     string
	logPath string // Path to log file on host
}

// DefaultEncoderV1Config returns a default encoder v1 configuration
func DefaultEncoderV1Config() EncoderConfig {
	return EncoderConfig{
		Enabled:               true,
		LogFormat:             "text",
		LogLevel:              "debug",
		GRPCPort:              "34000",
		AWSRegion:             "",
		AWSAccessKeyID:        "",
		AWSSecretAccessKey:    "",
		AWSEndpointURL:        "",
		S3BucketName:          "",
		EnableMetrics:         true,
		MetricsHTTPPort:       "9095",
		G1Path:                "",
		G2Path:                "",
		G2PowerOf2Path:        "",
		SRSOrder:              "10000",
		SRSLoad:               "10000",
		CachePath:             "",
		Verbose:               "true",
		NumWorkers:            "8",
		MaxConcurrentRequests: "16",
		RequestPoolSize:       "32",
		RequestQueueSize:      "32",
		EncoderVersion:        "1",
		Image:                 "ghcr.io/layr-labs/eigenda/encoder:dev",
	}
}

// DefaultEncoderV2Config returns a default encoder v2 configuration
func DefaultEncoderV2Config() EncoderConfig {
	config := DefaultEncoderV1Config()
	config.EncoderVersion = "2"
	config.S3BucketName = "test-eigenda-blobstore"
	return config
}

// NewEncoderContainerWithNetwork creates and starts a new encoder container with a custom network
func NewEncoderContainerWithNetwork(ctx context.Context, config EncoderConfig, network *testcontainers.DockerNetwork) (*EncoderContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", "encoder-logs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "encoder.log")

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

	// Mount KZG resources if paths are provided
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
	env := buildEncoderEnv(containerConfig)

	// Determine network alias and container name based on encoder version
	networkAlias := "encoder-v1"
	containerName := "eigenda-encoder-v1"
	if config.EncoderVersion == "2" {
		networkAlias = "encoder-v2"
		containerName = "eigenda-encoder-v2"
	}

	// Configure container request with network
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{config.GRPCPort + "/tcp"},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {networkAlias},
		},
		Mounts: mounts,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.GRPCPort+"/tcp")).WithStartupTimeout(30*time.Second),
			wait.ForLog("GRPC Listening").WithStartupTimeout(30*time.Second),
		),
		Name:            containerName,
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
		return nil, fmt.Errorf("failed to start encoder container: %w", err)
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

	encoderURL := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	fmt.Printf("Encoder logs will be available at: %s\n", logPath)

	return &EncoderContainer{
		Container: container,
		config:    config,
		url:       encoderURL,
		logPath:   logPath,
	}, nil
}

// URL returns the encoder service URL (gRPC endpoint)
func (c *EncoderContainer) URL() string {
	return c.url
}

// InternalURL returns the encoder service URL for internal network communication
func (c *EncoderContainer) InternalURL() string {
	networkAlias := "encoder-v1"
	if c.config.EncoderVersion == "2" {
		networkAlias = "encoder-v2"
	}
	return fmt.Sprintf("%s:%s", networkAlias, c.config.GRPCPort)
}

// Config returns the encoder configuration
func (c *EncoderContainer) Config() EncoderConfig {
	return c.config
}

// LogPath returns the path to the encoder log file on the host
func (c *EncoderContainer) LogPath() string {
	return c.logPath
}

// buildEncoderEnv builds environment variables for the encoder container
func buildEncoderEnv(config EncoderConfig) map[string]string {
	env := map[string]string{
		"DISPERSER_ENCODER_LOG_FORMAT":              config.LogFormat,
		"DISPERSER_ENCODER_LOG_LEVEL":               config.LogLevel,
		"DISPERSER_ENCODER_LOG_PATH":                "/logs/encoder.log", // Log to a file we can access
		"DISPERSER_ENCODER_GRPC_PORT":               config.GRPCPort,
		"DISPERSER_ENCODER_ENABLE_METRICS":          fmt.Sprintf("%t", config.EnableMetrics),
		"DISPERSER_ENCODER_METRICS_HTTP_PORT":       config.MetricsHTTPPort,
		"DISPERSER_ENCODER_SRS_ORDER":               config.SRSOrder,
		"DISPERSER_ENCODER_SRS_LOAD":                config.SRSLoad,
		"DISPERSER_ENCODER_VERBOSE":                 config.Verbose,
		"DISPERSER_ENCODER_NUM_WORKERS":             config.NumWorkers,
		"DISPERSER_ENCODER_MAX_CONCURRENT_REQUESTS": config.MaxConcurrentRequests,
		"DISPERSER_ENCODER_REQUEST_POOL_SIZE":       config.RequestPoolSize,
		"DISPERSER_ENCODER_REQUEST_QUEUE_SIZE":      config.RequestQueueSize,
	}

	// Add optional configurations if provided
	if config.G1Path != "" {
		env["DISPERSER_ENCODER_G1_PATH"] = config.G1Path
	}
	if config.G2Path != "" {
		env["DISPERSER_ENCODER_G2_PATH"] = config.G2Path
	}
	if config.G2PowerOf2Path != "" {
		env["DISPERSER_ENCODER_G2_POWER_OF_2_PATH"] = config.G2PowerOf2Path
	}
	if config.CachePath != "" {
		env["DISPERSER_ENCODER_CACHE_PATH"] = config.CachePath
	}
	if config.AWSRegion != "" {
		env["DISPERSER_ENCODER_AWS_REGION"] = config.AWSRegion
	}
	if config.AWSAccessKeyID != "" {
		env["DISPERSER_ENCODER_AWS_ACCESS_KEY_ID"] = config.AWSAccessKeyID
	}
	if config.AWSSecretAccessKey != "" {
		env["DISPERSER_ENCODER_AWS_SECRET_ACCESS_KEY"] = config.AWSSecretAccessKey
	}
	if config.AWSEndpointURL != "" {
		env["DISPERSER_ENCODER_AWS_ENDPOINT_URL"] = config.AWSEndpointURL
	}

	// Add encoder version if specified
	if config.EncoderVersion != "" {
		env["DISPERSER_ENCODER_ENCODER_VERSION"] = config.EncoderVersion
	}

	// Add S3 bucket name for encoder v2
	if config.EncoderVersion == "2" && config.S3BucketName != "" {
		env["DISPERSER_ENCODER_S3_BUCKET_NAME"] = config.S3BucketName
	}

	return env
}
