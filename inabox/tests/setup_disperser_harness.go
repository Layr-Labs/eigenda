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
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
	RelayInstances []*relay.Instance
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
		RelayInstances: make([]*relay.Instance, 0),
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

// startRelays starts all relay goroutines
func startRelays(ctx context.Context, harness *DisperserHarness, config DisperserHarnessConfig) error {
	config.Logger.Info("Pre-creating listeners for relay goroutines", "count", config.RelayCount)

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
					config.Logger.Warn("Failed to close listener for relay", "index", j, "error", err)
				}
			}
			return fmt.Errorf("failed to create listener for relay %d: %w", i, err)
		}
		listeners[i] = listener

		// Extract the actual port assigned by the OS
		actualPort := listener.Addr().(*net.TCPAddr).Port
		actualURLs[i] = fmt.Sprintf("0.0.0.0:%d", actualPort)

		config.Logger.Info("Created listener for relay", "index", i, "assigned_port", actualPort)
	}

	// Now that we have all the actual URLs, register them on-chain
	if config.TestConfig != nil && config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		config.Logger.Info("Registering relay URLs with actual ports", "urls", actualURLs)
		config.TestConfig.RegisterRelays(config.EthClient, actualURLs, config.EthClient.GetAccountAddress())
	}

	// Now start each relay with its pre-created listener
	for i, listener := range listeners {
		instance, err := startRelayWithListener(ctx, i, listener, harness, config)
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

// startRelayWithListener starts a single relay with the given index and pre-created listener
func startRelayWithListener(
	ctx context.Context,
	relayIndex int,
	listener net.Listener,
	harness *DisperserHarness,
	config DisperserHarnessConfig,
) (*relay.Instance, error) {
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

	// Create eth client config with RPC URL from deployer
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:    []string{config.TestConfig.Deployers[0].RPC},
		NumRetries: 2,
	}

	// Create relay server dependencies
	deps, err := relay.NewServerDependencies(
		ctx,
		relay.ServerDependenciesConfig{
			AWSConfig:                  awsConfig,
			MetadataTableName:          config.MetadataTableNameV2,
			BucketName:                 config.S3BucketName,
			OperatorStateRetrieverAddr: config.TestConfig.EigenDA.OperatorStateRetriever,
			ServiceManagerAddr:         config.TestConfig.EigenDA.ServiceManager,
			ChainStateConfig:           thegraph.Config{},
			EthClientConfig:            ethClientConfig,
			LoggerConfig:               loggerConfig,
		},
	)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create relay server dependencies: %w", err)
	}

	// Create relay test configuration (GRPCPort is ignored since listener is pre-created)
	relayConfig := relay.NewTestConfig(relayIndex)

	// Create and start the relay instance using the helper
	instance, err := relay.NewInstanceWithDependencies(ctx, relayConfig, deps, listener)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to create relay instance: %w", err)
	}

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	deps.Logger.Info("Relay server started successfully", "port", port, "logFile", logFilePath)

	return instance, nil
}

// stopAllRelays stops all relay instances
func stopAllRelays(instances []*relay.Instance, logger logging.Logger) {
	for i, instance := range instances {
		if instance == nil {
			continue
		}
		logger.Info("Stopping relay", "index", i, "url", instance.URL)
		if err := instance.Stop(); err != nil {
			logger.Warn("Error stopping relay instance", "index", i, "error", err)
		}
	}
}
