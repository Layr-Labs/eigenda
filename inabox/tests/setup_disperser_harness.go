package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/testcontainers/testcontainers-go"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	Network        *testcontainers.DockerNetwork
	TestConfig     *deploy.Config
	TestName       string
	LocalStackPort string

	// LocalStack resources for blobstore and metadata store
	MetadataTableName   string
	BucketTableName     string
	S3BucketName        string // S3 bucket name for blob storage
	MetadataTableNameV2 string

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int

	// OperatorStateSubgraphURL is the URL for the operator state subgraph
	OperatorStateSubgraphURL string
}

// TODO: Add encoder, api server, batcher
type DisperserHarness struct {
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadataV1 string
		BlobMetaV2     string
	}
	S3Buckets struct {
		BlobStore string
	}
	RelayServers []*relay.Server
	Controller   *controller.Controller
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
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		BlobStoreBucketName: config.S3BucketName,
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

	// Start controller goroutine
	if err := startController(ctx, ethClient, config.OperatorStateSubgraphURL, harness, config); err != nil {
		return nil, fmt.Errorf("failed to start controller: %w", err)
	}

	// Start remaining binaries (disperser, encoder, batcher, etc.)
	if config.TestConfig != nil {
		logger.Info("Starting remaining binaries")
		err := config.TestConfig.GenerateAllVariables()
		if err != nil {
			return nil, fmt.Errorf("could not generate environment variables: %w", err)
		}
		config.TestConfig.StartBinaries(true) // true = for tests, will skip churner and operators
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
	// Stop controller
	if dh.Controller != nil {
		logger.Info("Stopping controller")
		stopController(dh.Controller, logger)
	}

	// Stop relay goroutines
	if len(dh.RelayServers) > 0 {
		logger.Info("Stopping relay goroutines")
		stopAllRelays(dh.RelayServers, logger)
	}

	if dh.LocalStack != nil {
		logger.Info("Stopping localstack container")
		if err := dh.LocalStack.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate localstack container", "error", err)
		}
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

// startController starts the controller as a singleton goroutine
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
	chainStateConfig := thegraph.DefaultConfig()
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

	// Create encoder client
	encoderClient, err := encoder.NewEncoderClientV2(encodingManagerConfig.EncoderAddress)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create encoder client: %w", err)
	}

	// Create signature aggregator
	sigAgg, err := core.NewStdSignatureAggregator(controllerLogger, chainReader)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create signature aggregator: %w", err)
	}

	// Create indexed chain state
	chainState := eth.NewChainState(chainReader, ethClient)
	ics := thegraph.MakeIndexedChainState(chainStateConfig, chainState, controllerLogger)

	// Create node client manager
	nodeClientManager, err := controller.NewNodeClientManager(dispatcherConfig.NodeClientCacheSize, requestSigner, controllerLogger)
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

	// Create controller
	controllerInstance, err := controller.NewController(
		controllerLogger,
		metadataStore,
		chainReader,
		encoderClient,
		ics,
		batchMetadataManager,
		sigAgg,
		nodeClientManager,
		metricsRegistry,
		encodingManagerConfig,
		dispatcherConfig,
		nil, // No server config for tests
		nil, // No payment authorization handler for tests
	)
	if err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to create controller: %w", err)
	}

	// Start controller
	if err := controllerInstance.Start(ctx, controller.StartOptions{}); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("failed to start controller: %w", err)
	}

	// Store controller instance in harness
	harness.Controller = controllerInstance

	controllerLogger.Info("Controller started successfully", "logFile", logFilePath)

	return nil
}

// stopController stops the controller
func stopController(instance *controller.Controller, logger logging.Logger) {
	if instance == nil {
		return
	}

	// Note: EncodingManager and Dispatcher don't have explicit Stop methods in the current implementation
	logger.Info("Controller stopped")
}
