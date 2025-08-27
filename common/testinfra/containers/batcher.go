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

// BatcherConfig defines configuration for the batcher container
type BatcherConfig struct {
	// Enable batcher service
	Enabled bool

	// Log configuration
	LogFormat string
	LogLevel  string

	// AWS configuration
	S3BucketName       string
	DynamoDBTableName  string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSEndpointURL     string

	// Batch configuration
	PullInterval    string
	BatchSizeLimit  string
	SRSOrder        string
	UseGraph        bool
	GraphURL        string
	GraphBackoff    string
	GraphMaxRetries string

	// Encoder configuration
	EncoderAddress            string
	EncodingTimeout           string
	EncodingRequestQueueSize  string
	NumConnections            string
	MaxNumRetriesPerBlob      string
	TargetNumChunks           string
	MaxBlobsToFetchFromStore  string
	EnableGnarkBundleEncoding bool

	// Metrics configuration
	EnableMetrics   bool
	MetricsHTTPPort string

	// Chain configuration
	ChainRPC                    string
	ChainRPCFallback            string
	PrivateKey                  string
	NumConfirmations            string
	NumRetries                  string
	ChainReadTimeout            string
	ChainWriteTimeout           string
	ChainStateTimeout           string
	TransactionBroadcastTimeout string

	// EigenDA configuration
	EigenDADirectory       string
	OperatorStateRetriever string
	ServiceManager         string
	IndexerPullInterval    string
	IndexerDataDir         string

	// Finalizer configuration
	FinalizerInterval      string
	FinalizerPoolSize      string
	FinalizationBlockDelay string

	// Attestation configuration
	AttestationTimeout        string
	BatchAttestationTimeout   string
	MaxNodeConnections        string
	MaxNumRetriesPerDispersal string

	// Fragment configuration
	FragmentPrefixChars         string
	FragmentParallelismFactor   string
	FragmentParallelismConstant string
	FragmentReadTimeout         string
	FragmentWriteTimeout        string

	// KMS configuration
	KMSKeyID      string
	KMSKeyRegion  string
	KMSKeyDisable bool

	// Container image
	Image string
}

// BatcherContainer represents a running batcher container
type BatcherContainer struct {
	testcontainers.Container
	config  BatcherConfig
	url     string
	logPath string // Path to log file on host
}

// DefaultBatcherConfig returns a default batcher configuration
func DefaultBatcherConfig() BatcherConfig {
	return BatcherConfig{
		Enabled:                     true,
		LogFormat:                   "text",
		LogLevel:                    "debug",
		S3BucketName:                "test-eigenda-blobstore",
		DynamoDBTableName:           "test-BlobMetadata",
		AWSRegion:                   "us-east-1",
		AWSAccessKeyID:              "", // Will be populated from LocalStack or AWS config
		AWSSecretAccessKey:          "", // Will be populated from LocalStack or AWS config
		AWSEndpointURL:              "", // Will be populated from LocalStack if enabled
		PullInterval:                "5s",
		BatchSizeLimit:              "10240", // 10 GiB
		SRSOrder:                    "300000",
		UseGraph:                    true,
		GraphURL:                    "", // Will be populated from GraphNode if enabled
		GraphBackoff:                "100ms",
		GraphMaxRetries:             "5",
		EncoderAddress:              "", // Will be populated when encoder is deployed
		EncodingTimeout:             "10m",
		EncodingRequestQueueSize:    "500",
		NumConnections:              "256",
		MaxNumRetriesPerBlob:        "2",
		TargetNumChunks:             "8192",
		MaxBlobsToFetchFromStore:    "100",
		EnableGnarkBundleEncoding:   true,
		EnableMetrics:               true,
		MetricsHTTPPort:             "9094",
		ChainRPC:                    "", // Will be populated from Anvil
		ChainRPCFallback:            "",
		PrivateKey:                  "", // Will be populated from deployer key
		NumConfirmations:            "0",
		NumRetries:                  "3",
		ChainReadTimeout:            "12s",
		ChainWriteTimeout:           "12s",
		ChainStateTimeout:           "12s",
		TransactionBroadcastTimeout: "12s",
		EigenDADirectory:            "", // Will be populated from contract deployment
		OperatorStateRetriever:      "", // Will be populated from contract deployment
		ServiceManager:              "", // Will be populated from contract deployment
		IndexerPullInterval:         "1s",
		IndexerDataDir:              "./data/test-indexer",
		FinalizerInterval:           "6m",
		FinalizerPoolSize:           "5",
		FinalizationBlockDelay:      "0",
		AttestationTimeout:          "10m",
		BatchAttestationTimeout:     "11m",
		MaxNodeConnections:          "100",
		MaxNumRetriesPerDispersal:   "2",
		FragmentPrefixChars:         "3",
		FragmentParallelismFactor:   "8",
		FragmentParallelismConstant: "0",
		FragmentReadTimeout:         "30s",
		FragmentWriteTimeout:        "30s",
		KMSKeyID:                    "",
		KMSKeyRegion:                "",
		KMSKeyDisable:               true,
		Image:                       "ghcr.io/layr-labs/eigenda/batcher:dev",
	}
}

// NewBatcherContainerWithNetwork creates and starts a new batcher container with a custom network
func NewBatcherContainerWithNetwork(ctx context.Context, config BatcherConfig, network *testcontainers.DockerNetwork) (*BatcherContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", "batcher-logs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "batcher.log")

	// Build environment variables
	env := buildBatcherEnv(config)

	// Configure container request with network
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {"batcher"},
		},
		// Add host access configuration for accessing operators on host ports
		HostAccessPorts: []int{
			32011, 32012, 32013, 32014, 32015, 32016, // operator-0 ports
			32021, 32022, 32023, 32024, 32025, 32026, // operator-1 ports
			32031, 32032, 32033, 32034, 32035, 32036, // operator-2 ports
			32041, 32042, 32043, 32044, 32045, 32046, // operator-3 ports
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
			wait.ForLog("starting metrics server").WithStartupTimeout(60 * time.Second),
		),
		Name:            "eigenda-batcher",
		AlwaysPullImage: false, // Use local image if available
	}

	// Add metrics port if enabled
	if config.EnableMetrics {
		req.ExposedPorts = append(req.ExposedPorts, config.MetricsHTTPPort+"/tcp")
		req.WaitingFor = wait.ForAll(
			wait.ForListeningPort(nat.Port(config.MetricsHTTPPort+"/tcp")).WithStartupTimeout(60*time.Second),
			wait.ForLog("starting metrics server").WithStartupTimeout(60*time.Second),
		)
	}

	// Add HostConfigModifier to set up ExtraHosts for operator localhost domains
	// This maps operator-{i}.localhost to host-gateway, allowing the batcher
	// to reach operators running on the host through the localhost domain
	req.HostConfigModifier = func(hc *container.HostConfig) {
		// Add entries for each operator's localhost domain
		// host-gateway is a special Docker hostname that resolves to the host machine
		for i := 0; i < 4; i++ {
			operatorHost := fmt.Sprintf("operator-%d.localhost:host-gateway", i)
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
		return nil, fmt.Errorf("failed to start batcher container: %w", err)
	}

	// Get the container host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	batcherURL := host
	if config.EnableMetrics {
		// Get the mapped port for metrics
		mappedPort, err := container.MappedPort(ctx, nat.Port(config.MetricsHTTPPort))
		if err != nil {
			return nil, fmt.Errorf("failed to get mapped port: %w", err)
		}
		batcherURL = fmt.Sprintf("%s:%s", host, mappedPort.Port())
	}

	fmt.Printf("Batcher logs will be available at: %s\n", logPath)

	return &BatcherContainer{
		Container: container,
		config:    config,
		url:       batcherURL,
		logPath:   logPath,
	}, nil
}

// URL returns the batcher service URL (metrics endpoint if enabled)
func (c *BatcherContainer) URL() string {
	return c.url
}

// InternalURL returns the batcher service URL for internal network communication (metrics endpoint if enabled)
func (c *BatcherContainer) InternalURL() string {
	if c.config.EnableMetrics {
		return fmt.Sprintf("batcher:%s", c.config.MetricsHTTPPort)
	}
	return "batcher"
}

// Config returns the batcher configuration
func (c *BatcherContainer) Config() BatcherConfig {
	return c.config
}

// LogPath returns the path to the batcher log file on the host
func (c *BatcherContainer) LogPath() string {
	return c.logPath
}

// buildBatcherEnv builds environment variables for the batcher container
func buildBatcherEnv(config BatcherConfig) map[string]string {
	// Strip 0x prefix from private key if present
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")

	env := map[string]string{
		"BATCHER_LOG_FORMAT":                    config.LogFormat,
		"BATCHER_LOG_LEVEL":                     config.LogLevel,
		"BATCHER_LOG_PATH":                      "/logs/batcher.log", // Log to a file we can access
		"BATCHER_S3_BUCKET_NAME":                config.S3BucketName,
		"BATCHER_DYNAMODB_TABLE_NAME":           config.DynamoDBTableName,
		"BATCHER_PULL_INTERVAL":                 config.PullInterval,
		"BATCHER_ENCODER_ADDRESS":               config.EncoderAddress,
		"BATCHER_ENABLE_METRICS":                fmt.Sprintf("%t", config.EnableMetrics),
		"BATCHER_METRICS_HTTP_PORT":             config.MetricsHTTPPort,
		"BATCHER_BATCH_SIZE_LIMIT":              config.BatchSizeLimit,
		"BATCHER_USE_GRAPH":                     fmt.Sprintf("%t", config.UseGraph),
		"BATCHER_SRS_ORDER":                     config.SRSOrder,
		"BATCHER_INDEXER_DATA_DIR":              config.IndexerDataDir,
		"BATCHER_ENCODING_TIMEOUT":              config.EncodingTimeout,
		"BATCHER_ATTESTATION_TIMEOUT":           config.AttestationTimeout,
		"BATCHER_BATCH_ATTESTATION_TIMEOUT":     config.BatchAttestationTimeout,
		"BATCHER_CHAIN_READ_TIMEOUT":            config.ChainReadTimeout,
		"BATCHER_CHAIN_WRITE_TIMEOUT":           config.ChainWriteTimeout,
		"BATCHER_CHAIN_STATE_TIMEOUT":           config.ChainStateTimeout,
		"BATCHER_TRANSACTION_BROADCAST_TIMEOUT": config.TransactionBroadcastTimeout,
		"BATCHER_NUM_CONNECTIONS":               config.NumConnections,
		"BATCHER_FINALIZER_INTERVAL":            config.FinalizerInterval,
		"BATCHER_FINALIZER_POOL_SIZE":           config.FinalizerPoolSize,
		"BATCHER_ENCODING_REQUEST_QUEUE_SIZE":   config.EncodingRequestQueueSize,
		"BATCHER_MAX_NUM_RETRIES_PER_BLOB":      config.MaxNumRetriesPerBlob,
		"BATCHER_TARGET_NUM_CHUNKS":             config.TargetNumChunks,
		"BATCHER_MAX_BLOBS_TO_FETCH_FROM_STORE": config.MaxBlobsToFetchFromStore,
		"BATCHER_FINALIZATION_BLOCK_DELAY":      config.FinalizationBlockDelay,
		"BATCHER_MAX_NODE_CONNECTIONS":          config.MaxNodeConnections,
		"BATCHER_MAX_NUM_RETRIES_PER_DISPERSAL": config.MaxNumRetriesPerDispersal,
		"BATCHER_ENABLE_GNARK_BUNDLE_ENCODING":  fmt.Sprintf("%t", config.EnableGnarkBundleEncoding),
		"BATCHER_BLS_OPERATOR_STATE_RETRIVER":   config.OperatorStateRetriever,
		"BATCHER_EIGENDA_SERVICE_MANAGER":       config.ServiceManager,
		"BATCHER_EIGENDA_DIRECTORY":             config.EigenDADirectory,
		"BATCHER_CHAIN_RPC":                     config.ChainRPC,
		"BATCHER_PRIVATE_KEY":                   privateKey,
		"BATCHER_NUM_CONFIRMATIONS":             config.NumConfirmations,
		"BATCHER_NUM_RETRIES":                   config.NumRetries,
		"BATCHER_INDEXER_PULL_INTERVAL":         config.IndexerPullInterval,
		"BATCHER_FRAGMENT_PREFIX_CHARS":         config.FragmentPrefixChars,
		"BATCHER_FRAGMENT_PARALLELISM_FACTOR":   config.FragmentParallelismFactor,
		"BATCHER_FRAGMENT_PARALLELISM_CONSTANT": config.FragmentParallelismConstant,
		"BATCHER_FRAGMENT_READ_TIMEOUT":         config.FragmentReadTimeout,
		"BATCHER_FRAGMENT_WRITE_TIMEOUT":        config.FragmentWriteTimeout,
		"BATCHER_KMS_KEY_DISABLE":               fmt.Sprintf("%t", config.KMSKeyDisable),
	}

	// Add optional configurations if provided
	if config.AWSRegion != "" {
		env["BATCHER_AWS_REGION"] = config.AWSRegion
	}
	if config.AWSAccessKeyID != "" {
		env["BATCHER_AWS_ACCESS_KEY_ID"] = config.AWSAccessKeyID
	}
	if config.AWSSecretAccessKey != "" {
		env["BATCHER_AWS_SECRET_ACCESS_KEY"] = config.AWSSecretAccessKey
	}
	if config.AWSEndpointURL != "" {
		env["BATCHER_AWS_ENDPOINT_URL"] = config.AWSEndpointURL
	}
	if config.GraphURL != "" {
		env["BATCHER_GRAPH_URL"] = config.GraphURL
	}
	if config.GraphBackoff != "" {
		env["BATCHER_GRAPH_BACKOFF"] = config.GraphBackoff
	}
	if config.GraphMaxRetries != "" {
		env["BATCHER_GRAPH_MAX_RETRIES"] = config.GraphMaxRetries
	}
	if config.ChainRPCFallback != "" {
		env["BATCHER_CHAIN_RPC_FALLBACK"] = config.ChainRPCFallback
	}
	if config.KMSKeyID != "" {
		env["BATCHER_KMS_KEY_ID"] = config.KMSKeyID
	}
	if config.KMSKeyRegion != "" {
		env["BATCHER_KMS_KEY_REGION"] = config.KMSKeyRegion
	}

	return env
}
