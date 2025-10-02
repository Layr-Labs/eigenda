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
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	proverv2 "github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	TestConfig          *deploy.Config
	TestName            string
	LocalStackPort      string
	V2MetadataTableName string
	BlobStoreBucketName string // S3 bucket name for blob storage

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int
}

// TODO: Add api server, controller, batcher
type DisperserHarness struct {
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetaV2 string
	}
	S3Buckets struct {
		BlobStore string
	}
	RelayInstances    []*RelayInstance
	EncoderV2Instance *EncoderV2Instance
}

// setupLocalStackResources initializes LocalStack and deploys AWS resources
func setupV2LocalStackResources(
	ctx context.Context,
	logger logging.Logger,
	localstack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*testbed.LocalStackContainer, error) {
	// Deploy AWS resources (DynamoDB tables and S3 buckets)
	logger.Info("Deploying AWS resources in LocalStack")
	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  localstack.Endpoint(),
		V2MetadataTableName: config.V2MetadataTableName,
		BlobStoreBucketName: config.BlobStoreBucketName,
		AWSConfig:           localstack.GetAWSClientConfig(),
		Logger:              logger,
	}
	if err := testbed.DeployResources(ctx, deployConfig); err != nil {
		return nil, fmt.Errorf("failed to deploy resources: %w", err)
	}
	logger.Info("AWS resources deployed successfully")

	return localstack, nil
}

// setupDisperserKeypairAndRegistrations generates disperser keypair and performs registrations
func setupDisperserKeypairAndRegistrations(
	logger logging.Logger,
	ethClient common.EthClient,
	config DisperserHarnessConfig,
) error {
	if config.TestConfig == nil {
		return nil
	}

	logger.Info("Attempting to generate disperser keypair with LocalStack running")
	if err := config.TestConfig.GenerateDisperserKeypair(); err != nil {
		return fmt.Errorf("failed to generate disperser keypair: %w", err)
	}

	// Register disperser keypair on chain
	if config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		config.TestConfig.PerformDisperserRegistrations(ethClient)
	}

	return nil
}

// SetupDisperserHarness creates and initializes the disperser infrastructure
// (LocalStack, DynamoDB tables, S3 buckets, relays)
func SetupDisperserHarness(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	localstack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*DisperserHarness, error) {
	// Check if localstack resources are empty
	if config.V2MetadataTableName == "" || config.BlobStoreBucketName == "" {
		return nil, fmt.Errorf("missing name for localstack resources")
	}

	harness := &DisperserHarness{
		RelayInstances: make([]*RelayInstance, 0),
	}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetaV2 = config.V2MetadataTableName
	harness.S3Buckets.BlobStore = config.BlobStoreBucketName

	// Setup LocalStack if not using in-memory blob store
	localstack, err := setupV2LocalStackResources(ctx, logger, localstack, config)
	if err != nil {
		return nil, err
	}
	harness.LocalStack = localstack

	// Generate disperser keypair and perform registrations
	if err := setupDisperserKeypairAndRegistrations(logger, ethClient, config); err != nil {
		return nil, err
	}

	// Start relay goroutines if relay count is specified
	if config.RelayCount > 0 {
		if err := startRelays(ctx, logger, ethClient, harness, config); err != nil {
			return nil, fmt.Errorf("failed to start relays: %w", err)
		}
	} else {
		logger.Warn("Relay count is not specified, skipping relay setup")
	}

	// Start encoder v2 instance
	logger.Info("Starting encoder v2 instance")
	encoderInstance, err := startEncoderV2(ctx, logger, harness, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start encoder v2: %w", err)
	}
	harness.EncoderV2Instance = encoderInstance

	return harness, nil
}

// RelayInstance holds the state for a single relay
// TODO(dmanc): This (or something similar) should live in the relay package instead of here.
type RelayInstance struct {
	Server   *relay.Server
	Listener net.Listener
	Port     string
	URL      string
	Logger   logging.Logger
}

// EncoderV2Instance holds the state for a single encoder v2
type EncoderV2Instance struct {
	Server   *encoder.EncoderServerV2
	Listener net.Listener
	Port     string
	URL      string
	Logger   logging.Logger
}

// startRelays starts all relay goroutines
func startRelays(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) error {
	logger.Info("Pre-creating listeners for relay goroutines", "count", config.RelayCount)

	// Pre-create all listeners with port 0 (OS assigns ports)
	listeners := make([]net.Listener, config.RelayCount)
	actualURLs := make([]string, config.RelayCount)

	for i := range config.RelayCount {
		listener, err := net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			// Clean up any listeners we created before failing
			for j := range i {
				err := listeners[j].Close()
				if err != nil {
					logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to create listener for relay %d: %w", i, err)
		}
		listeners[i] = listener

		// Extract the actual port assigned by the OS
		actualPort := listener.Addr().(*net.TCPAddr).Port
		actualURLs[i] = fmt.Sprintf("0.0.0.0:%d", actualPort)

		logger.Info("Created listener for relay", "index", i, "assigned_port", actualPort)
	}

	// Now that we have all the actual URLs, register them on-chain
	if config.TestConfig != nil && config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		logger.Info("Registering relay URLs with actual ports", "urls", actualURLs)
		config.TestConfig.RegisterRelays(ethClient, actualURLs, ethClient.GetAccountAddress())
	}

	// Now start each relay with its pre-created listener
	for i, listener := range listeners {
		instance, err := startRelayWithListener(ctx, logger, ethClient, i, actualURLs[i], listener, harness, config)
		if err != nil {
			// Clean up any relays we started and all remaining listeners
			stopAllRelays(harness.RelayInstances, logger)
			for j := i; j < len(listeners); j++ {
				err := listeners[j].Close()
				if err != nil {
					logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to start relay %d (%s): %w", i, actualURLs[i], err)
		}
		harness.RelayInstances = append(harness.RelayInstances, instance)
		logger.Info("Started relay", "index", i, "url", actualURLs[i])
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

	// Stop encoder v2 instance
	if dh.EncoderV2Instance != nil {
		logger.Info("Stopping encoder v2 instance")
		dh.EncoderV2Instance.Logger.Info("Stopping encoder v2")
		dh.EncoderV2Instance.Server.Close()
	}
}

// startRelayWithListener starts a single relay with the given index, URL, and pre-created listener
func startRelayWithListener(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	relayIndex int,
	relayURL string,
	listener net.Listener,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) (*RelayInstance, error) {
	// Extract port from the listener's address
	port := fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)

	// Create logs directory
	// TODO(dmanc): If possible we should have a centralized place for creating loggers and injecting them into the config.
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
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, relayLogger, config.V2MetadataTableName)
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
		ethClient,
		config.TestConfig.EigenDA.OperatorStateRetriever,
		config.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth writer: %w", err)
	}

	cs := coreeth.NewChainState(tx, ethClient)
	ics := thegraph.MakeIndexedChainState(thegraph.Config{}, cs, relayLogger)

	// Create relay configuration
	// TODO(dmanc): In addition to loggers, we should have a centralized place for
	// setting up configuration and injecting it into the harness config.
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

// startEncoderV2 starts a single encoder v2 instance
func startEncoderV2(
	ctx context.Context,
	logger logging.Logger,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) (*EncoderV2Instance, error) {
	logger.Info("Starting encoder v2 instance")

	// Get SRS paths using the same function as operator setup
	g1Path, _, err := getSRSPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Pre-create listener with port 0 (OS assigns port)
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener for encoder v2: %w", err)
	}

	// Extract the actual port assigned by the OS
	actualPort := listener.Addr().(*net.TCPAddr).Port
	encoderURL := fmt.Sprintf("0.0.0.0:%d", actualPort)
	port := fmt.Sprintf("%d", actualPort)

	logger.Info("Created listener for encoder v2", "assigned_port", actualPort)

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/encoder_v2.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to open encoder v2 log file: %w", err)
	}

	// Create encoder logger
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	encoderLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create encoder v2 logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := harness.LocalStack.GetAWSClientConfig()
	s3Client, err := s3.NewClient(ctx, awsConfig, encoderLogger)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create blob store and chunk writer
	blobStore := blobstore.NewBlobStore(harness.S3Buckets.BlobStore, s3Client, encoderLogger)
	chunkWriter := chunkstore.NewChunkWriter(encoderLogger, s3Client, harness.S3Buckets.BlobStore, 4*1024*1024)

	// Create prover with dynamically determined paths
	// TODO(dmanc): Make these configurable
	kzgConfig := &proverv2.KzgConfig{
		SRSNumberToLoad: 10000,
		G1Path:          g1Path,
		LoadG2Points:    false, // Encoder doesn't need G2 points
		PreloadEncoder:  false,
		CacheDir:        fmt.Sprintf("testdata/%s/cache/encoder", config.TestName),
		NumWorker:       1,
		Verbose:         false,
	}
	encodingConfig := &encoding.Config{
		BackendType: encoding.GnarkBackend,
		GPUEnable:   false,
		NumWorker:   1,
	}

	prover, err := proverv2.NewProver(kzgConfig, encodingConfig)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create prover: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(metricsRegistry, "9200", encoderLogger)
	grpcMetrics := grpcprom.NewServerMetrics()
	metricsRegistry.MustRegister(grpcMetrics)

	// Create encoder server configuration
	serverConfig := encoder.ServerConfig{
		GrpcPort:              port,
		MaxConcurrentRequests: 32,
		RequestQueueSize:      100,
		PreventReencoding:     false,
		Backend:               "gnark",
		GPUEnable:             false,
	}

	// Create encoder server
	server := encoder.NewEncoderServerV2(
		serverConfig,
		blobStore,
		chunkWriter,
		encoderLogger,
		prover,
		metrics,
		grpcMetrics,
	)

	// Start the encoder server in a goroutine using the pre-created listener
	go func() {
		encoderLogger.Info("Starting encoder v2 server with listener", "port", port)
		if err := server.StartWithListener(listener); err != nil {
			encoderLogger.Error("Encoder v2 server failed", "error", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	encoderLogger.Info("Encoder v2 server started successfully", "port", port, "logFile", logFilePath)

	return &EncoderV2Instance{
		Server:   server,
		Listener: listener,
		Port:     port,
		URL:      encoderURL,
		Logger:   encoderLogger,
	}, nil
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
