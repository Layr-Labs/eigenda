package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/testcontainers/testcontainers-go"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	Logger              logging.Logger
	Network             *testcontainers.DockerNetwork
	TestConfig          *deploy.Config
	TestName            string
	InMemoryBlobStore   bool
	LocalStackPort      string
	MetadataTableName   string
	BucketTableName     string
	S3BucketName        string // S3 bucket name for blob storage
	MetadataTableNameV2 string
	EthClient           common.EthClient

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int
}

// TODO: Add encoder, api server, controller, batcher
type DisperserHarness struct {
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadataV1 string
		BlobMetaV2     string
	}
	S3Buckets struct {
		BlobStore string
	}
	RelayInstances []*RelayInstance
}

// setupLocalStackResources initializes LocalStack and deploys AWS resources
func setupLocalStackResources(
	ctx context.Context, config DisperserHarnessConfig,
) (*testbed.LocalStackContainer, error) {
	config.Logger.Info("Setting up LocalStack for blob store")
	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(
		ctx,
		testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       config.LocalStackPort,
			Logger:         config.Logger,
			Network:        config.Network,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to start localstack: %w", err)
	}

	// Deploy AWS resources (DynamoDB tables and S3 buckets)
	config.Logger.Info("Deploying AWS resources in LocalStack")
	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  localstackContainer.Endpoint(),
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		BlobStoreBucketName: config.S3BucketName,
		V2MetadataTableName: config.MetadataTableNameV2,
		AWSConfig:           localstackContainer.GetAWSClientConfig(),
		Logger:              config.Logger,
	}
	if err := testbed.DeployResources(ctx, deployConfig); err != nil {
		return nil, fmt.Errorf("failed to deploy resources: %w", err)
	}
	config.Logger.Info("AWS resources deployed successfully")

	return localstackContainer, nil
}

// setupDisperserKeypairAndRegistrations generates disperser keypair and performs registrations
func setupDisperserKeypairAndRegistrations(config DisperserHarnessConfig) error {
	if config.TestConfig == nil {
		return nil
	}

	config.Logger.Info("Attempting to generate disperser keypair with LocalStack running")
	if err := config.TestConfig.GenerateDisperserKeypair(); err != nil {
		return fmt.Errorf("failed to generate disperser keypair: %w", err)
	}

	// Register disperser keypair on chain
	if config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		config.TestConfig.PerformDisperserRegistrations(config.EthClient)
	}

	return nil
}

// SetupDisperserHarness creates and initializes the disperser infrastructure
// (LocalStack, DynamoDB tables, S3 buckets, relays)
func SetupDisperserHarness(ctx context.Context, config DisperserHarnessConfig) (*DisperserHarness, error) {
	harness := &DisperserHarness{
		RelayInstances: make([]*RelayInstance, 0),
	}

	// Set default values if not provided
	if config.LocalStackPort == "" {
		config.LocalStackPort = "4570"
	}
	if config.MetadataTableName == "" {
		config.MetadataTableName = "test-BlobMetadata"
	}
	if config.BucketTableName == "" {
		config.BucketTableName = "test-BucketStore"
	}
	if config.S3BucketName == "" {
		config.S3BucketName = "test-eigenda-blobstore"
	}
	if config.MetadataTableNameV2 == "" {
		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
	}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetadataV1 = config.MetadataTableName
	harness.DynamoDBTables.BlobMetaV2 = config.MetadataTableNameV2
	harness.S3Buckets.BlobStore = config.S3BucketName

	// Setup LocalStack if not using in-memory blob store
	if !config.InMemoryBlobStore {
		localstack, err := setupLocalStackResources(ctx, config)
		if err != nil {
			return nil, err
		}
		harness.LocalStack = localstack

		// Generate disperser keypair and perform registrations
		if err := setupDisperserKeypairAndRegistrations(config); err != nil {
			return nil, err
		}

		// Start relay goroutines if relay count is specified
		if config.RelayCount > 0 {
			if err := startRelays(ctx, harness, config); err != nil {
				return nil, fmt.Errorf("failed to start relays: %w", err)
			}
		} else {
			config.Logger.Warn("Relay count is not specified, skipping relay setup")
		}
	} else {
		// TODO(dmanc): Do the relays even work when not using S3 as the blob store?
		config.Logger.Info("Using in-memory blob store, skipping LocalStack setup")
	}

	// Start remaining binaries (disperser, encoder, batcher, etc.)
	if config.TestConfig != nil {
		config.Logger.Info("Starting remaining binaries")
		err := config.TestConfig.GenerateAllVariables()
		if err != nil {
			return nil, fmt.Errorf("could not generate environment variables: %w", err)
		}
		config.TestConfig.StartBinaries(true) // true = for tests, will skip churner and operators
	}

	return harness, nil
}

// RelayInstance holds the state for a single relay
type RelayInstance struct {
	Server   *relay.Server
	Listener net.Listener
	Port     string
	URL      string
	Logger   logging.Logger
}

// startRelays starts all relay goroutines
func startRelays(ctx context.Context, harness *DisperserHarness, config DisperserHarnessConfig) error {
	config.Logger.Info("Pre-creating listeners for relay goroutines", "count", config.RelayCount)

	// Pre-create all listeners with port 0 (OS assigns ports)
	listeners := make([]net.Listener, config.RelayCount)
	actualURLs := make([]string, config.RelayCount)

	for i := range config.RelayCount {
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			// Clean up any listeners we created before failing
			for j := range i {
				err := listeners[j].Close()
				if err != nil {
					config.Logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to create listener for relay %d: %w", i, err)
		}
		listeners[i] = listener

		// Extract the actual port assigned by the OS
		actualPort := listener.Addr().(*net.TCPAddr).Port
		actualURLs[i] = fmt.Sprintf("localhost:%d", actualPort)

		config.Logger.Info("Created listener for relay", "index", i, "assigned_port", actualPort)
	}

	// Now that we have all the actual URLs, register them on-chain
	if config.TestConfig != nil && config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		config.Logger.Info("Registering relay URLs with actual ports", "urls", actualURLs)
		config.TestConfig.RegisterRelays(config.EthClient, actualURLs, config.EthClient.GetAccountAddress())
	}

	// Now start each relay with its pre-created listener
	for i, listener := range listeners {
		instance, err := startRelayWithListener(ctx, i, actualURLs[i], listener, harness, config)
		if err != nil {
			// Clean up any relays we started and all remaining listeners
			stopAllRelays(harness.RelayInstances, config.Logger)
			for j := i; j < len(listeners); j++ {
				err := listeners[j].Close()
				if err != nil {
					config.Logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to start relay %d (%s): %w", i, actualURLs[i], err)
		}
		harness.RelayInstances = append(harness.RelayInstances, instance)
		config.Logger.Info("Started relay", "index", i, "url", actualURLs[i])
	}

	return nil
}

// Cleanup releases resources held by the DisperserHarness (excluding shared network)
func (dh *DisperserHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	// Stop relay goroutines
	if len(dh.RelayInstances) > 0 {
		logger.Info("Stopping relay goroutines")
		stopAllRelays(dh.RelayInstances, logger)
	}

	if dh.LocalStack != nil {
		logger.Info("Stopping localstack container")
		if err := dh.LocalStack.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate localstack container", "error", err)
		}
	}
}

// startRelayWithListener starts a single relay with the given index, URL, and pre-created listener
func startRelayWithListener(
	ctx context.Context,
	relayIndex int,
	relayURL string,
	listener net.Listener,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) (*RelayInstance, error) {
	// Extract port from the listener's address
	port := fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/relay_%d.log", logsDir, relayIndex)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open relay log file: %w", err)
	}

	// Create relay logger
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	relayLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := harness.LocalStack.GetAWSClientConfig()
	dynamoClient, err := dynamodb.NewClient(awsConfig, relayLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	s3Client, err := s3.NewClient(ctx, awsConfig, relayLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, relayLogger, config.MetadataTableNameV2)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "relay",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create blob store and chunk reader
	blobStore := blobstore.NewBlobStore(harness.S3Buckets.BlobStore, s3Client, relayLogger)
	chunkReader := chunkstore.NewChunkReader(relayLogger, s3Client, harness.S3Buckets.BlobStore)

	// Create eth writer and chain state
	tx, err := coreeth.NewWriter(
		relayLogger,
		config.EthClient,
		config.TestConfig.EigenDA.OperatorStateRetriever,
		config.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth writer: %w", err)
	}

	cs := coreeth.NewChainState(tx, config.EthClient)
	ics := thegraph.MakeIndexedChainState(thegraph.Config{}, cs, relayLogger)

	// Create relay configuration
	relayConfig := &relay.Config{
		RelayKeys:                  []v2.RelayKey{v2.RelayKey(relayIndex)}, // Serve data for any shard
		GRPCPort:                   mustParsePort(port),
		MaxGRPCMessageSize:         1024 * 1024 * 300,
		MetadataCacheSize:          1024 * 1024,
		MetadataMaxConcurrency:     32,
		BlobCacheBytes:             32 * 1024 * 1024,
		BlobMaxConcurrency:         32,
		ChunkCacheBytes:            32 * 1024 * 1024,
		ChunkMaxConcurrency:        32,
		MaxKeysPerGetChunksRequest: 1024,
		RateLimits: limiter.Config{
			MaxGetBlobOpsPerSecond:          1024,
			GetBlobOpsBurstiness:            1024,
			MaxGetBlobBytesPerSecond:        20 * 1024 * 1024,
			GetBlobBytesBurstiness:          20 * 1024 * 1024,
			MaxConcurrentGetBlobOps:         1024,
			MaxGetChunkOpsPerSecond:         1024,
			GetChunkOpsBurstiness:           1024,
			MaxGetChunkBytesPerSecond:       20 * 1024 * 1024,
			GetChunkBytesBurstiness:         20 * 1024 * 1024,
			MaxConcurrentGetChunkOps:        1024,
			MaxGetChunkOpsPerSecondClient:   8,
			GetChunkOpsBurstinessClient:     8,
			MaxGetChunkBytesPerSecondClient: 2 * 1024 * 1024,
			GetChunkBytesBurstinessClient:   2 * 1024 * 1024,
			MaxConcurrentGetChunkOpsClient:  1,
		},
		AuthenticationKeyCacheSize:   1024,
		AuthenticationDisabled:       true, // Disabled for testing
		GetChunksRequestMaxPastAge:   5 * time.Minute,
		GetChunksRequestMaxFutureAge: 1 * time.Minute,
		Timeouts: relay.TimeoutConfig{
			GetChunksTimeout:               20 * time.Second,
			GetBlobTimeout:                 20 * time.Second,
			InternalGetMetadataTimeout:     5 * time.Second,
			InternalGetBlobTimeout:         20 * time.Second,
			InternalGetProofsTimeout:       5 * time.Second,
			InternalGetCoefficientsTimeout: 20 * time.Second,
		},
		OnchainStateRefreshInterval: 10 * time.Second,
		MetricsPort:                 9100 + relayIndex,
		EnableMetrics:               true,
		EnablePprof:                 false,
		PprofHttpPort:               0,
		MaxConnectionAge:            0,
		MaxConnectionAgeGrace:       5 * time.Second,
		MaxIdleConnectionAge:        30 * time.Second,
	}

	// Create relay server
	server, err := relay.NewServer(
		ctx,
		metricsRegistry,
		relayLogger,
		relayConfig,
		metadataStore,
		blobStore,
		chunkReader,
		tx,
		ics,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay server: %w", err)
	}

	// Start the relay server in a goroutine using the pre-created listener
	go func() {
		relayLogger.Info("Starting relay server with listener", "port", port)
		if err := server.StartWithListener(ctx, listener); err != nil {
			relayLogger.Error("Relay server failed", "error", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	relayLogger.Info("Relay server started successfully", "port", port, "logFile", logFilePath)

	return &RelayInstance{
		Server:   server,
		Listener: listener,
		Port:     port,
		URL:      relayURL,
		Logger:   relayLogger,
	}, nil
}

// stopAllRelays stops all relay instances
func stopAllRelays(instances []*RelayInstance, logger logging.Logger) {
	for i, instance := range instances {
		if instance == nil {
			continue
		}
		logger.Info("Stopping relay", "index", i, "url", instance.URL)
		instance.Logger.Info("Stopping relay")
		// TODO: Add graceful shutdown of relay server once it's implemented
		// For now, the context cancellation in TeardownInfrastructure will stop the server
	}
}

// mustParsePort parses a port string to an int, panicking on error
func mustParsePort(portStr string) int {
	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		panic(fmt.Sprintf("invalid port: %s", portStr))
	}
	return port
}
