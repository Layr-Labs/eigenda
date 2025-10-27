package integration

// // InfrastructureConfig contains the configuration for setting up the infrastructure
// type InfrastructureConfig struct {
// 	TemplateName        string
// 	TestName            string
// 	Logger              logging.Logger
// 	RootPath            string
// 	MetadataTableName   string
// 	BucketTableName     string
// 	S3BucketName        string
// 	MetadataTableNameV2 string

// 	// Number of relay instances to start, if not specified, no relays will be started.
// 	RelayCount int

// 	// DisableDisperser disables the disperser deployment when set to true. This is useful for
// 	// tests that do not require the disperser infrastructure to be deployed (e.g. testing graph
// 	// node with operator registration)
// 	DisableDisperser bool

// 	// The following field is temporary, to be able to test different payments configurations. It will be removed
// 	// once legacy payments are removed.
// 	ControllerUseNewPayments bool
// }

// // SetupInfrastructure creates the shared infrastructure that persists across all tests.
// // This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
// func SetupInfrastructure(ctx context.Context, config *InfrastructureConfig) (*InfrastructureHarness, error) {
// 	var err error
// 	var infra *InfrastructureHarness
// 	if config.MetadataTableName == "" {
// 		config.MetadataTableName = "test-BlobMetadata"
// 	}
// 	if config.BucketTableName == "" {
// 		config.BucketTableName = "test-BucketStore"
// 	}
// 	if config.MetadataTableNameV2 == "" {
// 		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
// 	}

// 	logger := config.Logger

// 	// Create test directory if needed
// 	testName := config.TestName
// 	if testName == "" {
// 		testName, err = deploy.CreateNewTestDirectory(config.TemplateName, config.RootPath)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create test directory: %w", err)
// 		}
// 	}

// 	testConfig := deploy.ReadTestConfig(testName, config.RootPath)
// 	testConfig.UseControllerMediatedPayments = config.ControllerUseNewPayments

// 	// Create a long-lived context for the infrastructure lifecycle
// 	infraCtx, infraCancel := context.WithCancel(ctx)

// 	// Ensure we cancel the context if we return an error
// 	defer func() {
// 		if err != nil {
// 			infraCancel()
// 		}
// 	}()

// 	// Create shared Docker network, primarily for Anvil and Graph Node
// 	sharedDockerNetwork, err := network.New(
// 		infraCtx,
// 		network.WithDriver("bridge"),
// 		network.WithAttachable())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create docker network: %w", err)
// 	}
// 	logger.Info("Created Docker network", "name", sharedDockerNetwork.Name)

// 	// Create infrastructure harness early so we can populate it incrementally
// 	infra = &InfrastructureHarness{
// 		SharedNetwork:  sharedDockerNetwork,
// 		TestConfig:     testConfig,
// 		TemplateName:   config.TemplateName,
// 		TestName:       testName,
// 		LocalStackPort: "4570",
// 		Logger:         config.Logger,
// 		Cancel:         infraCancel,
// 	}

// 	// Setup Chain Harness first (Anvil, Graph Node, Contracts, Churner)
// 	chainHarnessConfig := &ChainHarnessConfig{
// 		TestConfig: testConfig,
// 		TestName:   testName,
// 		Logger:     logger,
// 		Network:    sharedDockerNetwork,
// 	}
// 	chainHarness, err := SetupChainHarness(infraCtx, chainHarnessConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to setup chain harness: %w", err)
// 	}
// 	infra.ChainHarness = *chainHarness

// 	// // Setup Disperser Harness second (LocalStack, DynamoDB tables, S3 buckets, relays)
// 	// if !config.DisableDisperser {
// 	// 	disperserHarnessConfig := &DisperserHarnessConfig{
// 	// 		Network:             sharedDockerNetwork,
// 	// 		TestConfig:          testConfig,
// 	// 		TestName:            testName,
// 	// 		LocalStackPort:      infra.LocalStackPort,
// 	// 		MetadataTableName:   config.MetadataTableName,
// 	// 		BucketTableName:     config.BucketTableName,
// 	// 		S3BucketName:        config.S3BucketName,
// 	// 		MetadataTableNameV2: config.MetadataTableNameV2,
// 	// 		RelayCount:          config.RelayCount,
// 	// 	}
// 	// 	disperserHarness, err := SetupDisperserHarness(
// 	// 		infraCtx,
// 	// 		logger,
// 	// 		infra.ChainHarness.EthClient,
// 	// 		*disperserHarnessConfig,
// 	// 	)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to setup disperser harness: %w", err)
// 	// 	}
// 	// 	infra.DisperserHarness = *disperserHarness
// 	// } else {
// 	// 	logger.Info("Disperser deployment disabled, skipping disperser harness setup")
// 	// }

// 	// Setup Operator Harness third (requires chain and disperser to be ready)
// 	operatorHarnessConfig := &OperatorHarnessConfig{
// 		TestConfig: testConfig,
// 		TestName:   testName,
// 	}
// 	operatorHarness, err := SetupOperatorHarness(infraCtx, logger, &infra.ChainHarness, operatorHarnessConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to setup operator harness: %w", err)
// 	}
// 	infra.OperatorHarness = *operatorHarness

// 	return infra, nil
// }

// // TeardownGlobalInfrastructure cleans up all global infrastructure
// func TeardownInfrastructure(infra *InfrastructureHarness) {
// 	infra.Logger.Info("Tearing down global infrastructure")

// 	// Cancel the infrastructure context to signal all components to shut down
// 	if infra.Cancel != nil {
// 		infra.Logger.Info("Cancelling infrastructure context")
// 		infra.Cancel()
// 	}

// 	// Create a separate timeout context for cleanup operations
// 	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	defer cleanupCancel()

// 	// Stop operator goroutines using the harness cleanup
// 	infra.OperatorHarness.Cleanup(infra.Logger)

// 	// Stop test binaries
// 	infra.Logger.Info("Stopping binaries")
// 	infra.TestConfig.StopBinaries()

// 	// Clean up disperser harness
// 	// infra.DisperserHarness.Cleanup(cleanupCtx, infra.Logger)

// 	// Clean up chain harness (churner and anvil)
// 	infra.ChainHarness.Cleanup(cleanupCtx, infra.Logger)

// 	// Clean up the shared Docker network last since multiple harnesses use it
// 	if infra.SharedNetwork != nil {
// 		infra.Logger.Info("Removing shared Docker network")
// 		_ = infra.SharedNetwork.Remove(cleanupCtx)
// 	}

// 	infra.Logger.Info("Teardown completed")
// }
