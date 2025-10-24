package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/testcontainers/testcontainers-go"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	Network        *testcontainers.DockerNetwork
	TestConfig     *deploy.Config
	TestName       string
	LocalStackPort string

	// S3 bucket name for blob storage
	S3BucketName string

	// V1 metadata table name
	MetadataTableName string

	// V2 metadata table name
	MetadataTableNameV2 string

	// DynamoDB table name for on-demand payments, currently used by the controller.
	OnDemandTableName string

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int

	// OperatorStateSubgraphURL is the URL for the operator state subgraph
	OperatorStateSubgraphURL string
}

// DisperserHarness is the harness for spinning up the disperser infrastructure as goroutines.
// It will only support V2 components of the disperser.
// TODO: Add api server
type DisperserHarness struct {
	// LocalStack infrastructure for blobstore and metadata store
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadata   string
		BlobMetadataV2 string
	}
	S3Buckets struct {
		BlobStore string
	}

	// Relay
	RelayServers []*relay.Server

	// Encoder V2
	EncoderServerV2 *encoder.EncoderServerV2

	// Controller components
	// TODO: Refactor into a single struct for controller components
	EncodingManager  *controller.EncodingManager
	Dispatcher       *controller.Dispatcher
	ControllerServer *server.Server
}

// setupLocalStackResources initializes LocalStack and deploys AWS resources
func setupLocalStackResources(
	ctx context.Context,
	logger logging.Logger,
	config DisperserHarnessConfig,
) (*testbed.LocalStackContainer, error) {
	logger.Info("Setting up LocalStack for blob store")
	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(
		ctx,
		testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       config.LocalStackPort,
			Logger:         logger,
			Network:        config.Network,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to start localstack: %w", err)
	}

	// Deploy AWS resources (DynamoDB tables and S3 buckets)
	logger.Info("Deploying AWS resources in LocalStack")
	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  localstackContainer.Endpoint(),
		BlobStoreBucketName: config.S3BucketName,
		V1MetadataTableName: config.MetadataTableName,
		V2MetadataTableName: config.MetadataTableNameV2,
		AWSConfig:           localstackContainer.GetAWSClientConfig(),
		Logger:              logger,
	}
	if err := testbed.DeployResources(ctx, deployConfig); err != nil {
		return nil, fmt.Errorf("failed to deploy resources: %w", err)
	}
	logger.Info("AWS resources deployed successfully")

	return localstackContainer, nil
}

// setupDisperserKeypairAndRegistrations generates disperser keypair and performs registrations
func setupDisperserKeypairAndRegistrations(
	logger logging.Logger,
	ethClient common.EthClient,
	config DisperserHarnessConfig) error {
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
	config DisperserHarnessConfig,
) (*DisperserHarness, error) {
	harness := &DisperserHarness{
		RelayServers: make([]*relay.Server, 0),
	}

	if config.OperatorStateSubgraphURL == "" {
		return nil, fmt.Errorf("operator state subgraph URL is required")
	}

	// Set default values if not provided
	if config.LocalStackPort == "" {
		config.LocalStackPort = "4570"
	}
	if config.S3BucketName == "" {
		config.S3BucketName = "test-eigenda-blobstore"
	}
	if config.MetadataTableName == "" {
		config.MetadataTableName = "test-BlobMetadata"
	}
	if config.MetadataTableNameV2 == "" {
		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
	}
	if config.OnDemandTableName == "" {
		config.OnDemandTableName = "e2e_v2_ondemand"
	}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetadata = config.MetadataTableName
	harness.DynamoDBTables.BlobMetadataV2 = config.MetadataTableNameV2
	harness.S3Buckets.BlobStore = config.S3BucketName

	localstack, err := setupLocalStackResources(ctx, logger, config)
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

	// Start encoder v2 goroutine
	if err := startEncoderV2(ctx, harness, config); err != nil {
		return nil, fmt.Errorf("failed to start encoder v2: %w", err)
	}

	// Start controller goroutine
	if err := startController(ctx, ethClient, config.OperatorStateSubgraphURL, harness, config); err != nil {
		return nil, fmt.Errorf("failed to start controller: %w", err)
	}

	// Start remaining binaries (disperser, batcher, etc.)
	if config.TestConfig != nil {
		logger.Info("Starting remaining binaries")
		err := config.TestConfig.GenerateAllVariables("", "")
		if err != nil {
			return nil, fmt.Errorf("could not generate environment variables: %w", err)
		}

		// Start binaries for tests, will skip churner, operators, encoder v2, controller, and relays
		config.TestConfig.StartBinaries(true)
	}

	return harness, nil
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
		instance, err := startRelayWithListener(ctx, ethClient, i, listener, harness, config)
		if err != nil {
			// Clean up any relays we started and all remaining listeners
			stopAllRelays(harness.RelayServers, logger)
			for j := i; j < len(listeners); j++ {
				err := listeners[j].Close()
				if err != nil {
					logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to start relay %d (%s): %w", i, actualURLs[i], err)
		}
		harness.RelayServers = append(harness.RelayServers, instance)
		logger.Info("Started relay", "index", i, "url", actualURLs[i])
	}

	return nil
}

// Cleanup releases resources held by the DisperserHarness (excluding shared network)
func (dh *DisperserHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	// Stop encoder v2 server
	if dh.EncoderServerV2 != nil {
		logger.Info("Stopping encoder v2 server")
		dh.EncoderServerV2.Close()
	}

	// Stop controller components
	if dh.ControllerServer != nil {
		logger.Info("Stopping controller gRPC server")
		dh.ControllerServer.Stop()
	}

	// Note: EncodingManager and Dispatcher don't have explicit Stop methods in the current implementation
	// They will be cleaned up when the context is cancelled or the process exits

	// Stop relay goroutines
	if len(dh.RelayServers) > 0 {
		logger.Info("Stopping relay goroutines")
		stopAllRelays(dh.RelayServers, logger)
	}
}

// startRelayWithListener starts a single relay with the given index and pre-created listener
func startRelayWithListener(
	ctx context.Context,
	ethClient common.EthClient,
	relayIndex int,
	listener net.Listener,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) (*relay.Server, error) {
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

	// Create relay logger config for file output
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := harness.LocalStack.GetAWSClientConfig()

	// Create logger
	logger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig, logger)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create S3 client
	s3Client, err := s3.NewClient(ctx, awsConfig, logger)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.MetadataTableNameV2)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "relay",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create blob store and chunk reader
	blobStore := blobstore.NewBlobStore(config.S3BucketName, s3Client, logger)
	chunkReader := chunkstore.NewChunkReader(logger, s3Client, config.S3BucketName)

	// Create eth writer
	tx, err := eth.NewWriter(
		logger,
		ethClient,
		config.TestConfig.EigenDA.OperatorStateRetriever,
		config.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create eth writer: %w", err)
	}

	// Create chain state
	cs := eth.NewChainState(tx, ethClient)
	ics := thegraph.MakeIndexedChainState(thegraph.Config{}, cs, logger)

	// Create relay test configuration
	relayConfig := relay.NewTestConfig(relayIndex)

	// Create server
	server, err := relay.NewServer(
		ctx,
		metricsRegistry,
		logger,
		relayConfig,
		metadataStore,
		blobStore,
		chunkReader,
		tx,
		ics,
		listener,
	)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create relay server: %w", err)
	}

	// Start server in background
	go func() {
		logger.Info("Starting relay server", "address", listener.Addr().String(), "logFile", logFilePath)
		if err := server.Start(ctx); err != nil {
			logger.Error("Relay server failed", "error", err)
		}
	}()

	// TODO(dmanc): Replace with proper health check endpoint
	logger.Info("Relay server started successfully", "port", listener.Addr().(*net.TCPAddr).Port, "logFile", logFilePath)

	return server, nil
}

// stopAllRelays stops all relay servers
func stopAllRelays(servers []*relay.Server, logger logging.Logger) {
	for i, server := range servers {
		if server == nil {
			continue
		}
		logger.Info("Stopping relay", "index", i)
		if err := server.Stop(); err != nil {
			logger.Warn("Error stopping relay server", "index", i, "error", err)
		}
	}
}

// startEncoderV2 starts the encoder v2 server as a goroutine
func startEncoderV2(
	ctx context.Context,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) error {
	if config.TestConfig == nil {
		return fmt.Errorf("test config is required to start encoder v2")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/enc1.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open encoder log file: %w", err)
	}

	// Create encoder logger config for file output
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	// Create logger
	encoderLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := harness.LocalStack.GetAWSClientConfig()

	// Create S3 client
	s3Client, err := s3.NewClient(ctx, awsConfig, encoderLogger)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create encoder metrics
	encoderMetrics := encoder.NewMetrics(metricsRegistry, "9099", encoderLogger)
	grpcMetrics := grpcprom.NewServerMetrics()
	metricsRegistry.MustRegister(grpcMetrics)

	// Start metrics server
	encoderMetrics.Start(ctx)

	// Get SRS paths using the utility function
	g1Path, _, err := getSRSPaths()
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Construct cache directory path from g1Path
	srsDir := filepath.Dir(g1Path)
	cacheDir := filepath.Join(srsDir, "SRSTables")

	// Create prover
	kzgConfig := prover.KzgConfig{
		G1Path:          g1Path,
		CacheDir:        cacheDir,
		SRSNumberToLoad: 10000,
		NumWorker:       1,
	}

	encodingConfig := &encoding.Config{
		BackendType: encoding.GnarkBackend,
		GPUEnable:   false,
		NumWorker:   1,
	}

	proverV2, err := prover.NewProver(encoderLogger, &kzgConfig, encodingConfig)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create prover: %w", err)
	}

	// Create blob store
	blobStore := blobstore.NewBlobStore(config.S3BucketName, s3Client, encoderLogger)

	// Create chunk writer
	const DefaultFragmentSizeBytes = 4 * 1024 * 1024
	chunkWriter := chunkstore.NewChunkWriter(encoderLogger, s3Client, config.S3BucketName, DefaultFragmentSizeBytes)

	// Create encoder server config
	serverConfig := encoder.ServerConfig{
		GrpcPort:              "34001",
		MaxConcurrentRequests: 16,
		RequestQueueSize:      32,
		PreventReencoding:     true,
		Backend:               "gnark",
		GPUEnable:             false,
	}

	// Create encoder server
	encoderServerV2 := encoder.NewEncoderServerV2(
		serverConfig,
		blobStore,
		chunkWriter,
		encoderLogger,
		proverV2,
		encoderMetrics,
		grpcMetrics,
	)

	// Pre-create listener
	listener, err := net.Listen("tcp", "0.0.0.0:34001")
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create listener for encoder v2: %w", err)
	}

	// Start encoder server in background
	go func() {
		encoderLogger.Info("Starting encoder v2 server", "address", listener.Addr().String(), "logFile", logFilePath)
		if err := encoderServerV2.StartWithListener(listener); err != nil {
			encoderLogger.Error("Encoder v2 server failed", "error", err)
		}
	}()

	// Store encoder in harness
	harness.EncoderServerV2 = encoderServerV2

	encoderLogger.Info("Encoder v2 server started successfully", "logFile", logFilePath)

	return nil
}

// startController starts the controller components (encoding manager and dispatcher)
func startController(
	ctx context.Context,
	ethClient common.EthClient,
	operatorStateSubgraphURL string,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) error {
	if config.TestConfig == nil {
		return fmt.Errorf("test config is required to start controller")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/controller.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open controller log file: %w", err)
	}

	// Create controller logger config for file output
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	// Create logger
	controllerLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := harness.LocalStack.GetAWSClientConfig()

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig, controllerLogger)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Get available relays from config
	availableRelays := make([]corev2.RelayKey, config.RelayCount)
	for i := range config.RelayCount {
		availableRelays[i] = corev2.RelayKey(i)
	}

	requestSigner, err := clients.NewDispersalRequestSigner(
		ctx,
		clients.DispersalRequestSignerConfig{
			Region:   awsConfig.Region,
			Endpoint: awsConfig.EndpointURL,
			KeyID:    config.TestConfig.DisperserKMSKeyID,
		})
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create dispersal request signer: %w", err)
	}

	// Build encoding manager configs
	encodingManagerConfig := controller.DefaultEncodingManagerConfig()
	encodingManagerConfig.NumRelayAssignment = uint16(config.RelayCount)
	encodingManagerConfig.AvailableRelays = availableRelays
	encodingManagerConfig.EncoderAddress = "localhost:34001"

	// Build dispatcher configs
	dispatcherConfig := controller.DefaultDispatcherConfig()
	dispatcherConfig.FinalizationBlockDelay = 5
	dispatcherConfig.BatchMetadataUpdatePeriod = 100 * time.Millisecond

	// Chain state config
	chainStateConfig := thegraph.Config{
		PullInterval: 100 * time.Millisecond,
		MaxRetries:   5,
	}
	chainStateConfig.Endpoint = operatorStateSubgraphURL

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, controllerLogger, config.MetadataTableNameV2)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "controller",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create chain reader
	chainReader, err := eth.NewReader(
		controllerLogger,
		ethClient,
		config.TestConfig.EigenDA.OperatorStateRetriever,
		config.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Create heartbeat channel
	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)

	// Create encoder client
	encoderClient, err := encoder.NewEncoderClientV2(encodingManagerConfig.EncoderAddress)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create encoder client: %w", err)
	}

	// Create encoding manager with workerpool and blob set
	encodingPool := workerpool.New(encodingManagerConfig.NumConcurrentRequests)
	encodingManagerBlobSet := controller.NewBlobSet()
	encodingManager, err := controller.NewEncodingManager(
		encodingManagerConfig,
		metadataStore,
		encodingPool,
		encoderClient,
		chainReader,
		controllerLogger,
		metricsRegistry,
		encodingManagerBlobSet,
		controllerLivenessChan,
	)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create encoding manager: %w", err)
	}

	// Create signature aggregator
	sigAgg, err := core.NewStdSignatureAggregator(controllerLogger, chainReader)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create signature aggregator: %w", err)
	}

	// Create dispatcher pool
	dispatcherPool := workerpool.New(dispatcherConfig.NumConcurrentRequests)

	// Create indexed chain state
	chainState := eth.NewChainState(chainReader, ethClient)
	ics := thegraph.MakeIndexedChainState(chainStateConfig, chainState, controllerLogger)

	// Create node client manager
	nodeClientManager, err := controller.NewNodeClientManager(
		dispatcherConfig.NodeClientCacheSize,
		requestSigner,
		controllerLogger,
	)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create node client manager: %w", err)
	}

	// Create batch metadata manager
	batchMetadataManager, err := metadata.NewBatchMetadataManager(
		ctx,
		controllerLogger,
		ethClient,
		ics,
		gethcommon.HexToAddress(config.TestConfig.EigenDA.RegistryCoordinator),
		dispatcherConfig.BatchMetadataUpdatePeriod,
		dispatcherConfig.FinalizationBlockDelay,
	)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create batch metadata manager: %w", err)
	}

	// Create beforeDispatch callback to remove blobs from encoding manager's set
	beforeDispatch := func(blobKey corev2.BlobKey) error {
		encodingManagerBlobSet.RemoveBlob(blobKey)
		return nil
	}
	dispatcherBlobSet := controller.NewBlobSet()

	// Create dispatcher
	dispatcher, err := controller.NewDispatcher(
		dispatcherConfig,
		metadataStore,
		dispatcherPool,
		ics,
		batchMetadataManager,
		sigAgg,
		nodeClientManager,
		controllerLogger,
		metricsRegistry,
		beforeDispatch,
		dispatcherBlobSet,
		controllerLivenessChan,
	)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create dispatcher: %w", err)
	}

	// Recover state before starting
	if err := controller.RecoverState(ctx, metadataStore, controllerLogger); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to recover state: %w", err)
	}

	// Start encoding manager
	if err := encodingManager.Start(ctx); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to start encoding manager: %w", err)
	}

	// Start dispatcher
	if err := dispatcher.Start(ctx); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to start dispatcher: %w", err)
	}

	// Store components in harness
	harness.EncodingManager = encodingManager
	harness.Dispatcher = dispatcher

	// Build and start gRPC server if payments are enabled
	if config.TestConfig.UseControllerMediatedPayments {
		controllerLogger.Info("UseControllerMediatedPayments enabled - starting gRPC server")

		// Create contract directory
		contractDirectory, err := directory.NewContractDirectory(
			ctx,
			controllerLogger,
			ethClient,
			gethcommon.HexToAddress(config.TestConfig.EigenDA.EigenDADirectory),
		)
		if err != nil {
			_ = logFile.Close()
			return fmt.Errorf("failed to create contract directory: %w", err)
		}

		// Build payment authorization handler
		paymentAuthConfig := controller.DefaultPaymentAuthorizationConfig()
		paymentAuthConfig.OnDemandConfig.OnDemandTableName = config.OnDemandTableName
		paymentAuthConfig.OnDemandConfig.UpdateInterval = 1 * time.Second
		paymentAuthConfig.ReservationConfig.UpdateInterval = 1 * time.Second

		paymentAuthorizationHandler, err := controller.BuildPaymentAuthorizationHandler(
			ctx,
			controllerLogger,
			*paymentAuthConfig,
			contractDirectory,
			ethClient,
			dynamoClient.GetAwsClient(),
			metricsRegistry,
		)
		if err != nil {
			_ = logFile.Close()
			return fmt.Errorf("failed to build payment authorization handler: %w", err)
		}

		// Create server config
		grpcServerConfig, err := common.NewGRPCServerConfig(
			true,
			30000, // TODO(dmanc): inject listener instead
			1024*1024,
			5*time.Minute,
			5*time.Minute,
			3*time.Minute,
		)
		if err != nil {
			_ = logFile.Close()
			return fmt.Errorf("failed to create gRPC server config: %w", err)
		}

		serverConfig, err := server.NewConfig(
			grpcServerConfig,
			true, // EnablePaymentAuthentication
		)
		if err != nil {
			_ = logFile.Close()
			return fmt.Errorf("failed to create server config: %w", err)
		}

		// Create and start gRPC server
		grpcServer, err := server.NewServer(
			ctx,
			serverConfig,
			controllerLogger,
			metricsRegistry,
			paymentAuthorizationHandler,
		)
		if err != nil {
			_ = logFile.Close()
			return fmt.Errorf("failed to create gRPC server: %w", err)
		}

		go func() {
			controllerLogger.Info("Starting controller gRPC server", "port", serverConfig.GrpcPort)
			if err := grpcServer.Start(); err != nil {
				controllerLogger.Error("gRPC server failed", "error", err)
			}
		}()

		harness.ControllerServer = grpcServer
		controllerLogger.Info("Controller gRPC server started successfully")
	} else {
		controllerLogger.Info("UseControllerMediatedPayments disabled - controller will not have server")
	}

	controllerLogger.Info("Controller components started successfully", "logFile", logFilePath)

	return nil
}
