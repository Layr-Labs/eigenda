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
	grpccontroller "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	awss3 "github.com/Layr-Labs/eigenda/common/s3/aws"
	"github.com/Layr-Labs/eigenda/core"
	authv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	Network        *testcontainers.DockerNetwork
	TestConfig     *deploy.Config
	TestName       string
	LocalStackPort string

	// LocalStack resources for blobstore and metadata store
	MetadataTableName string
	BucketTableName   string

	// S3 bucket name for blob storage
	S3BucketName string

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
type DisperserHarness struct {
	// LocalStack infrastructure for blobstore and metadata store
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadataV1 string
		BlobMetadataV2 string
	}
	S3Buckets struct {
		BlobStore string
	}

	// Relay
	RelayServers []*relay.Server

	// Encoder
	EncoderServer *encoder.EncoderServerV2

	// API Server
	APIServer        *apiserver.DispersalServerV2
	APIServerAddress string

	// Controller components
	// TODO: Refactor into a single struct for controller components
	EncodingManager  *controller.EncodingManager
	Dispatcher       *controller.Dispatcher
	ControllerServer *server.Server
}

// TODO: Consider refactoring these component structs into the underlying packages (relay, encoder, controller,
// apiserver). This would reduce maintenance burden on tests - if the production code changes, the component structs
// would be updated alongside it. Currently these exist here because production code runs each service as a separate
// binary, while the test harness runs them as goroutines and needs to return/track the created objects.

// RelayComponents contains the components created by startRelays
type RelayComponents struct {
	Servers []*relay.Server
}

// EncoderComponents contains the components created by startEncoder
type EncoderComponents struct {
	Server  *encoder.EncoderServerV2
	Address string
}

// ControllerComponents contains the components created by startController
type ControllerComponents struct {
	EncodingManager  *controller.EncodingManager
	Dispatcher       *controller.Dispatcher
	ControllerServer *server.Server
	Address          string
}

// APIServerComponents contains the components created by startAPIServer
type APIServerComponents struct {
	Server  *apiserver.DispersalServerV2
	Address string
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
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
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
	if config.OnDemandTableName == "" {
		config.OnDemandTableName = "e2e_v2_ondemand"
	}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetadataV1 = config.MetadataTableName
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
		relayComponents, err := startRelays(ctx, logger, ethClient, harness.LocalStack, config)
		if err != nil {
			return nil, fmt.Errorf("failed to start relays: %w", err)
		}
		harness.RelayServers = relayComponents.Servers
	} else {
		logger.Warn("Relay count is not specified, skipping relay setup")
	}

	// Start encoder goroutine
	encoderComponents, err := startEncoder(ctx, harness.LocalStack, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start encoder: %w", err)
	}
	harness.EncoderServer = encoderComponents.Server
	encoderAddress := encoderComponents.Address

	// Start controller goroutine
	controllerComponents, err := startController(
		ctx,
		ethClient,
		config.OperatorStateSubgraphURL,
		encoderAddress,
		harness.LocalStack,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start controller: %w", err)
	}
	harness.EncodingManager = controllerComponents.EncodingManager
	harness.Dispatcher = controllerComponents.Dispatcher
	harness.ControllerServer = controllerComponents.ControllerServer

	// Start API server goroutine
	apiServerComponents, err := startAPIServer(
		ctx,
		ethClient,
		controllerComponents.Address,
		harness.LocalStack,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start API server: %w", err)
	}
	harness.APIServer = apiServerComponents.Server
	harness.APIServerAddress = apiServerComponents.Address

	// Start remaining binaries (disperser, batcher, etc.)
	if config.TestConfig != nil {
		logger.Info("Starting remaining binaries")
		err := config.TestConfig.GenerateAllVariables()
		if err != nil {
			return nil, fmt.Errorf("could not generate environment variables: %w", err)
		}

		// Start binaries for tests, will skip churner, operators, encoder, controller, and relays
		config.TestConfig.StartBinaries(true)
	}

	return harness, nil
}

// startRelays starts all relay goroutines
func startRelays(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	localStack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*RelayComponents, error) {
	logger.Info("Pre-creating listeners for relay goroutines", "count", config.RelayCount)

	// Pre-create all listeners with port 0 (OS assigns ports)
	listeners := make([]net.Listener, config.RelayCount)
	assignedURLs := make([]string, config.RelayCount)

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
			return nil, fmt.Errorf("failed to create listener for relay %d: %w", i, err)
		}
		listeners[i] = listener

		// Extract the port assigned by the OS
		assignedPort := listener.Addr().(*net.TCPAddr).Port
		assignedURLs[i] = fmt.Sprintf("0.0.0.0:%d", assignedPort)

		logger.Info("Created listener for relay", "index", i, "assigned_port", assignedPort)
	}

	// Now that we have all the assigned URLs, register them on-chain
	if config.TestConfig != nil && config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		logger.Info("Registering relay URLs with assigned ports", "urls", assignedURLs)
		config.TestConfig.RegisterRelays(ethClient, assignedURLs, ethClient.GetAccountAddress())
	}

	// Now start each relay with its pre-created listener
	relayServers := make([]*relay.Server, 0, config.RelayCount)
	for i, listener := range listeners {
		instance, err := startRelayWithListener(ctx, ethClient, i, listener, localStack, config)
		if err != nil {
			// Clean up any relays we started and all remaining listeners
			stopAllRelays(relayServers, logger)
			for j := i; j < len(listeners); j++ {
				err := listeners[j].Close()
				if err != nil {
					logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return nil, fmt.Errorf("failed to start relay %d (%s): %w", i, assignedURLs[i], err)
		}
		relayServers = append(relayServers, instance)
		logger.Info("Started relay", "index", i, "url", assignedURLs[i])
	}

	return &RelayComponents{
		Servers: relayServers,
	}, nil
}

// Cleanup releases resources held by the DisperserHarness (excluding shared network)
func (dh *DisperserHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	// Stop encoder server
	if dh.EncoderServer != nil {
		logger.Info("Stopping encoder server")
		dh.EncoderServer.Close()
	}

	// Stop API server
	if dh.APIServer != nil {
		logger.Info("Stopping API server")
		if err := dh.APIServer.Stop(); err != nil {
			logger.Error("Failed to stop API server", "error", err)
		}
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

	// Clean up LocalStack
	if dh.LocalStack != nil {
		logger.Info("Terminating LocalStack container")
		if err := dh.LocalStack.Terminate(ctx); err != nil {
			logger.Error("Failed to terminate LocalStack container", "error", err)
		}
	}
}

// startRelayWithListener starts a single relay with the given index and pre-created listener
func startRelayWithListener(
	ctx context.Context,
	ethClient common.EthClient,
	relayIndex int,
	listener net.Listener,
	localStack *testbed.LocalStackContainer,
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
	defer func() {
		if err != nil {
			_ = logFile.Close()
		}
	}()

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
	awsConfig := localStack.GetAWSClientConfig()

	// Create logger
	logger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create S3 client
	s3Client, err := awss3.NewAwsS3Client(
		ctx,
		logger,
		awsConfig.EndpointURL,
		awsConfig.Region,
		awsConfig.FragmentParallelismFactor,
		awsConfig.FragmentParallelismConstant,
		awsConfig.AccessKey,
		awsConfig.SecretAccessKey,
	)
	if err != nil {
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

// startEncoder starts the encoder server as a goroutine and returns the encoder components
func startEncoder(
	ctx context.Context,
	localStack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*EncoderComponents, error) {
	if config.TestConfig == nil {
		return nil, fmt.Errorf("test config is required to start encoder")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/enc1.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open encoder log file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = logFile.Close()
		}
	}()

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
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := localStack.GetAWSClientConfig()

	// Create S3 client
	s3Client, err := awss3.NewAwsS3Client(
		ctx,
		encoderLogger,
		awsConfig.EndpointURL,
		awsConfig.Region,
		awsConfig.FragmentParallelismFactor,
		awsConfig.FragmentParallelismConstant,
		awsConfig.AccessKey,
		awsConfig.SecretAccessKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
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
	g1Path, _, _, err := getSRSPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to determine SRS file paths: %w", err)
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

	prover, err := prover.NewProver(encoderLogger, &kzgConfig, encodingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create prover: %w", err)
	}

	// Create blob store
	blobStore := blobstore.NewBlobStore(config.S3BucketName, s3Client, encoderLogger)

	// Create chunk writer
	const DefaultFragmentSizeBytes = 4 * 1024 * 1024
	chunkWriter := chunkstore.NewChunkWriter(encoderLogger, s3Client, config.S3BucketName, DefaultFragmentSizeBytes)

	// Create encoder server config
	serverConfig := encoder.ServerConfig{
		MaxConcurrentRequestsDangerous: 16,
		RequestQueueSize:               32,
		PreventReencoding:              true,
		Backend:                        "gnark",
		GPUEnable:                      false,
	}

	// Create encoder server
	encoderServer := encoder.NewEncoderServerV2(
		serverConfig,
		blobStore,
		chunkWriter,
		encoderLogger,
		prover,
		encoderMetrics,
		grpcMetrics,
	)

	// Pre-create listener with port 0 (OS assigns random port)
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener for encoder: %w", err)
	}

	// Extract the port assigned by the OS
	assignedPort := listener.Addr().(*net.TCPAddr).Port
	assignedAddress := fmt.Sprintf("localhost:%d", assignedPort)

	encoderLogger.Info("Created listener for encoder", "assigned_port", assignedPort, "address", assignedAddress)

	// Start encoder server in background
	go func() {
		encoderLogger.Info("Starting encoder server", "address", listener.Addr().String(), "logFile", logFilePath)
		if err := encoderServer.StartWithListener(listener); err != nil {
			encoderLogger.Error("Encoder server failed", "error", err)
		}
	}()

	encoderLogger.Info("Encoder server started successfully", "address", assignedAddress, "logFile", logFilePath)

	return &EncoderComponents{
		Server:  encoderServer,
		Address: assignedAddress,
	}, nil
}

// startController starts the controller components (encoding manager and dispatcher)
// and returns the controller components
func startController(
	ctx context.Context,
	ethClient common.EthClient,
	operatorStateSubgraphURL string,
	encoderAddress string,
	localStack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*ControllerComponents, error) {
	if config.TestConfig == nil {
		return nil, fmt.Errorf("test config is required to start controller")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/controller.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open controller log file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = logFile.Close()
		}
	}()

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
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := localStack.GetAWSClientConfig()

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig, controllerLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
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
		return nil, fmt.Errorf("failed to create dispersal request signer: %w", err)
	}

	// Build encoding manager configs
	encodingManagerConfig := controller.DefaultEncodingManagerConfig()
	encodingManagerConfig.NumRelayAssignment = uint16(config.RelayCount)
	encodingManagerConfig.AvailableRelays = availableRelays
	encodingManagerConfig.EncoderAddress = encoderAddress

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
		return nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Create heartbeat channel
	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)

	// Create encoder client
	encoderClient, err := encoder.NewEncoderClientV2(encodingManagerConfig.EncoderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder client: %w", err)
	}

	// Create encoding manager with workerpool and blob set
	encodingPool := workerpool.New(encodingManagerConfig.NumConcurrentRequests)
	encodingManagerBlobSet := controller.NewBlobSet()
	encodingManager, err := controller.NewEncodingManager(
		encodingManagerConfig,
		time.Now,
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
		return nil, fmt.Errorf("failed to create encoding manager: %w", err)
	}

	// Create signature aggregator
	sigAgg, err := core.NewStdSignatureAggregator(controllerLogger, chainReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature aggregator: %w", err)
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
		return nil, fmt.Errorf("failed to create node client manager: %w", err)
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
		return nil, fmt.Errorf("failed to create batch metadata manager: %w", err)
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
		return nil, fmt.Errorf("failed to create dispatcher: %w", err)
	}

	// Recover state before starting
	if err := controller.RecoverState(ctx, metadataStore, controllerLogger); err != nil {
		return nil, fmt.Errorf("failed to recover state: %w", err)
	}

	// Start encoding manager
	if err := encodingManager.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start encoding manager: %w", err)
	}

	// Start dispatcher
	if err := dispatcher.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start dispatcher: %w", err)
	}

	// Build and start gRPC server if payments are enabled
	var controllerServer *server.Server
	var controllerAddress string
	//nolint:nestif // Complex nested block is temporary until old payment system is removed
	if config.TestConfig.UseNewPayments {
		controllerLogger.Info("UseNewPayments enabled - starting gRPC server")

		// Create contract directory
		contractDirectory, err := directory.NewContractDirectory(
			ctx,
			controllerLogger,
			ethClient,
			gethcommon.HexToAddress(config.TestConfig.EigenDA.EigenDADirectory),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create contract directory: %w", err)
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
			return nil, fmt.Errorf("failed to build payment authorization handler: %w", err)
		}

		// Pre-create listener with port 0 (OS assigns random port)
		listener, err := net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			return nil, fmt.Errorf("failed to create listener for controller: %w", err)
		}
		defer func() {
			if err != nil {
				_ = listener.Close()
			}
		}()

		// Extract the port assigned by the OS
		assignedPort := listener.Addr().(*net.TCPAddr).Port
		controllerLogger.Info("Created listener for controller", "assigned_port", assignedPort)

		// Create server config
		grpcServerConfig, err := common.NewGRPCServerConfig(
			true,
			uint16(assignedPort),
			1024*1024,
			5*time.Minute,
			5*time.Minute,
			3*time.Minute,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC server config: %w", err)
		}

		serverConfig, err := server.NewConfig(
			grpcServerConfig,
			true, // EnablePaymentAuthentication
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create server config: %w", err)
		}

		// Create and start gRPC server
		grpcServer, err := server.NewServer(
			ctx,
			serverConfig,
			controllerLogger,
			metricsRegistry,
			paymentAuthorizationHandler,
			listener,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC server: %w", err)
		}

		go func() {
			controllerLogger.Info("Starting controller gRPC server", "address", listener.Addr().String())
			if err := grpcServer.Start(); err != nil {
				controllerLogger.Error("gRPC server failed", "error", err)
			}
		}()

		controllerServer = grpcServer
		controllerAddress = fmt.Sprintf("localhost:%d", assignedPort)
		controllerLogger.Info("Controller gRPC server started successfully", "address", controllerAddress)
	} else {
		// When server is disabled, use empty address
		controllerAddress = ""
		controllerLogger.Info("UseNewPayments disabled - controller will not have server")
	}

	controllerLogger.Info("Controller components started successfully",
		"address", controllerAddress, "logFile", logFilePath)

	return &ControllerComponents{
		EncodingManager:  encodingManager,
		Dispatcher:       dispatcher,
		ControllerServer: controllerServer,
		Address:          controllerAddress,
	}, nil
}

// startAPIServer starts the API server as a goroutine and returns the API server components
func startAPIServer(
	ctx context.Context,
	ethClient common.EthClient,
	controllerAddress string,
	localStack *testbed.LocalStackContainer,
	config DisperserHarnessConfig,
) (*APIServerComponents, error) {
	if config.TestConfig == nil {
		return nil, fmt.Errorf("test config is required to start API server")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/apiserver.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open API server log file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = logFile.Close()
		}
	}()

	// Create API server logger config for file output
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
	apiServerLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create AWS clients using LocalStack container's configuration
	awsConfig := localStack.GetAWSClientConfig()

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig, apiServerLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create S3 client
	s3Client, err := awss3.NewAwsS3Client(
		ctx,
		apiServerLogger,
		awsConfig.EndpointURL,
		awsConfig.Region,
		awsConfig.FragmentParallelismFactor,
		awsConfig.FragmentParallelismConstant,
		awsConfig.AccessKey,
		awsConfig.SecretAccessKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, apiServerLogger, config.MetadataTableNameV2)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "apiserver",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create blob store
	blobStore := blobstore.NewBlobStore(config.S3BucketName, s3Client, apiServerLogger)

	// Create committer
	g1Path, g2Path, g2TrailingPath, err := getSRSPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	committerConfig := committer.Config{
		SRSNumberToLoad:   10000,
		G1SRSPath:         g1Path,
		G2SRSPath:         g2Path,
		G2TrailingSRSPath: g2TrailingPath,
	}

	kzgCommitter, err := committer.NewFromConfig(committerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create committer: %w", err)
	}

	// Create chain reader
	chainReader, err := eth.NewReader(
		apiServerLogger,
		ethClient,
		config.TestConfig.EigenDA.OperatorStateRetriever,
		config.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Create blob request authenticator
	authenticator := authv2.NewPaymentStateAuthenticator(
		5*time.Minute, // AuthPmtStateRequestMaxPastAge
		5*time.Minute, // AuthPmtStateRequestMaxFutureAge
	)

	// Create meterer
	// Note: The meterer is always created to serve GetPaymentState calls, even when using
	// controller-mediated payments. The UseNewPayments flag controls which
	// payment system is used for authorization during dispersals, but doesn't affect
	// whether the meterer is available for querying payment state.
	apiServerLogger.Info("Creating meterer")

	mtConfig := meterer.Config{
		ChainReadTimeout: 5 * time.Second,
		UpdateInterval:   1 * time.Second, // Match deploy config for tests
	}

	paymentChainState, err := meterer.NewOnchainPaymentState(ctx, chainReader, apiServerLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create onchain payment state: %w", err)
	}
	if err := paymentChainState.RefreshOnchainPaymentState(ctx); err != nil {
		return nil, fmt.Errorf("failed to make initial query to the on-chain state: %w", err)
	}

	// Use the standard v2 payment table prefix
	const v2PaymentPrefix = "e2e_v2_"
	meteringStore, err := meterer.NewDynamoDBMeteringStore(
		awsConfig,
		v2PaymentPrefix+"reservation",        // ReservationsTableName
		v2PaymentPrefix+"ondemand",           // OnDemandTableName
		v2PaymentPrefix+"global_reservation", // GlobalRateTableName
		apiServerLogger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create offchain store: %w", err)
	}

	mt := meterer.NewMeterer(
		mtConfig,
		paymentChainState,
		meteringStore,
		apiServerLogger,
	)
	mt.Start(ctx)

	// Pre-create listener with port 0 (OS assigns random port)
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener for API server: %w", err)
	}
	defer func() {
		if err != nil {
			_ = listener.Close()
		}
	}()

	// Extract the port assigned by the OS
	assignedPort := listener.Addr().(*net.TCPAddr).Port
	apiServerLogger.Info("Created listener for API server", "assigned_port", assignedPort)

	// Create server config
	serverConfig := disperser.ServerConfig{
		GrpcPort:              fmt.Sprintf("%d", assignedPort),
		GrpcTimeout:           10 * time.Second,
		MaxConnectionAge:      5 * time.Minute,
		MaxConnectionAgeGrace: 30 * time.Second,
		MaxIdleConnectionAge:  1 * time.Minute,
	}

	metricsConfig := disperser.MetricsConfig{
		HTTPPort:      "9100",
		EnableMetrics: true,
	}

	// Max number of symbols per blob (based on typical config)
	const maxNumSymbolsPerBlob = 16 * 1024 * 1024

	// Onchain state refresh interval
	onchainStateRefreshInterval := 1 * time.Second

	// Create controller client if using new payments
	var controllerConnection *grpc.ClientConn
	var controllerClient grpccontroller.ControllerServiceClient
	if config.TestConfig.UseNewPayments {
		if controllerAddress == "" {
			return nil, fmt.Errorf("controller address is empty but UseNewPayments is true")
		}
		connection, err := grpc.NewClient(
			controllerAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("create controller connection: %w", err)
		}
		controllerConnection = connection
		controllerClient = grpccontroller.NewControllerServiceClient(connection)
	}

	// Create API server
	// Note: meterer is nil when using controller-mediated payments, otherwise it's the legacy meterer
	apiServer, err := apiserver.NewDispersalServerV2(
		serverConfig,
		time.Now,
		blobStore,
		metadataStore,
		chainReader,
		mt,
		authenticator,
		kzgCommitter,
		maxNumSymbolsPerBlob,
		onchainStateRefreshInterval,
		45*time.Second, // maxDispersalAge
		45*time.Second, // maxFutureDispersalTime
		apiServerLogger,
		metricsRegistry,
		metricsConfig,
		false,                            // ReservedOnly
		config.TestConfig.UseNewPayments, // useControllerMediatedPayments
		controllerConnection,
		controllerClient,
		listener,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API server: %w", err)
	}

	// Start API server in background
	go func() {
		apiServerLogger.Info("Starting API server", "address", listener.Addr().String(), "logFile", logFilePath)
		if err := apiServer.Start(ctx); err != nil {
			apiServerLogger.Error("API server failed", "error", err)
		}
	}()

	actualAddress := fmt.Sprintf("localhost:%d", assignedPort)
	apiServerLogger.Info("API server started successfully", "address", actualAddress, "logFile", logFilePath)

	return &APIServerComponents{
		Server:  apiServer,
		Address: actualAddress,
	}, nil
}
