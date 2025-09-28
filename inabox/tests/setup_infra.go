package integration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/testcontainers/testcontainers-go/network"
)

// InfrastructureConfig contains the configuration for setting up the infrastructure
type InfrastructureConfig struct {
	TemplateName        string
	TestName            string
	InMemoryBlobStore   bool
	Logger              logging.Logger
	MetadataTableName   string
	BucketTableName     string
	MetadataTableNameV2 string
}

// SetupGlobalInfrastructure creates the shared infrastructure that persists across all tests.
// This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
// This should be called once in TestMain.
func SetupGlobalInfrastructure(config *InfrastructureConfig) (*InfrastructureHarness, error) {
	if config.MetadataTableName == "" {
		config.MetadataTableName = "test-BlobMetadata"
	}
	if config.BucketTableName == "" {
		config.BucketTableName = "test-BucketStore"
	}
	if config.MetadataTableNameV2 == "" {
		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
	}

	// Create a timeout context for setup operations only
	setupCtx, setupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer setupCancel()

	rootPath := "../../"

	// Create test directory if needed
	testName := config.TestName
	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(config.TemplateName, rootPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create test directory: %w", err)
		}
	}

	testConfig := deploy.ReadTestConfig(testName, rootPath)

	// Create a long-lived context for the infrastructure lifecycle
	infraCtx, infraCancel := context.WithCancel(context.Background())

	if testConfig.Environment.IsLocal() {
		return setupLocalInfrastructure(setupCtx, infraCtx, infraCancel, config, testName, testConfig)
	}

	// For non-local environments, just return a minimal harness
	return &InfrastructureHarness{
		TemplateName:        config.TemplateName,
		TestName:            testName,
		InMemoryBlobStore:   config.InMemoryBlobStore,
		Logger:              config.Logger,
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		MetadataTableNameV2: config.MetadataTableNameV2,
		TestConfig:          testConfig,
		Ctx:                 infraCtx,
		Cancel:              infraCancel,
	}, nil
}

func setupLocalInfrastructure(
	setupCtx context.Context,
	infraCtx context.Context,
	infraCancel context.CancelFunc,
	config *InfrastructureConfig,
	testName string,
	testConfig *deploy.Config,
) (*InfrastructureHarness, error) {
	infra := &InfrastructureHarness{
		TemplateName:        config.TemplateName,
		TestName:            testName,
		InMemoryBlobStore:   config.InMemoryBlobStore,
		Logger:              config.Logger,
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		MetadataTableNameV2: config.MetadataTableNameV2,
		TestConfig:          testConfig,
		Ctx:                 infraCtx,
		Cancel:              infraCancel,
	}

	// Create a shared Docker network for all containers
	var err error
	infra.ChainDockerNetwork, err = network.New(
		setupCtx,
		network.WithDriver("bridge"),
		network.WithAttachable())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}
	infra.Logger.Info("Created Docker network", "name", infra.ChainDockerNetwork.Name)

	// Setup blob store
	if !infra.InMemoryBlobStore {
		infra.Logger.Info("Using shared Blob Store")
		infra.LocalStackPort = "4570"
		infra.LocalstackContainer, err = testbed.NewLocalStackContainerWithOptions(
			setupCtx,
			testbed.LocalStackOptions{
				ExposeHostPort: true,
				HostPort:       infra.LocalStackPort,
				Logger:         infra.Logger,
				Network:        infra.ChainDockerNetwork,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to start localstack: %w", err)
		}

		deployConfig := testbed.DeployResourcesConfig{
			LocalStackEndpoint:  fmt.Sprintf("http://0.0.0.0:%s", infra.LocalStackPort),
			MetadataTableName:   infra.MetadataTableName,
			BucketTableName:     infra.BucketTableName,
			V2MetadataTableName: infra.MetadataTableNameV2,
			Logger:              infra.Logger,
		}
		err = testbed.DeployResources(setupCtx, deployConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy resources: %w", err)
		}
	} else {
		infra.Logger.Info("Using in-memory Blob Store")
	}

	// Setup Anvil
	infra.Logger.Info("Starting anvil")
	infra.AnvilContainer, err = testbed.NewAnvilContainerWithOptions(
		setupCtx,
		testbed.AnvilOptions{
			ExposeHostPort: true,
			HostPort:       "8545",
			Logger:         infra.Logger,
			Network:        infra.ChainDockerNetwork,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil: %w", err)
	}

	// Get Anvil RPC URL from the running container
	anvilRPC := infra.AnvilContainer.RpcURL()

	// Setup Graph Node if needed
	deployer, ok := infra.TestConfig.GetDeployer(infra.TestConfig.EigenDA.Deployer)
	if ok && deployer.DeploySubgraphs {
		infra.Logger.Info("Starting graph node")
		anvilInternalEndpoint := infra.AnvilContainer.InternalEndpoint()
		infra.GraphNodeContainer, err = testbed.NewGraphNodeContainerWithOptions(
			setupCtx,
			testbed.GraphNodeOptions{
				PostgresDB:     "graph-node",
				PostgresUser:   "graph-node",
				PostgresPass:   "let-me-in",
				EthereumRPC:    anvilInternalEndpoint,
				ExposeHostPort: true,
				HostHTTPPort:   "8000",
				HostWSPort:     "8001",
				HostAdminPort:  "8020",
				HostIPFSPort:   "5001",
				Logger:         infra.Logger,
				Network:        infra.ChainDockerNetwork,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to start graph node: %w", err)
		}
	}

	// Deploy contracts
	infra.Logger.Info("Deploying experiment")
	err = infra.TestConfig.DeployExperiment()
	if err != nil {
		return nil, fmt.Errorf("failed to deploy experiment: %w", err)
	}

	// Start churner goroutine
	infra.Logger.Info("Starting churner server")
	churnerURL, err := StartChurnerForInfrastructure(infra, anvilRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to start churner server: %w", err)
	}
	infra.Logger.Info("Churner server started", "address", churnerURL)

	// Start operator goroutines
	infra.Logger.Info("Starting operator goroutines")
	err = StartOperatorsForInfrastructure(infra, anvilRPC, churnerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to start operator goroutines: %w", err)
	}

	// Start remaining binaries (disperser, encoder, batcher, etc.)
	infra.Logger.Info("Starting remaining binaries")
	infra.TestConfig.StartBinaries(true) // true = for tests, will skip churner and operators

	return infra, nil
}

// TeardownGlobalInfrastructure cleans up all global infrastructure
func TeardownGlobalInfrastructure(infra *InfrastructureHarness) {
	infra.Logger.Info("Tearing down global infrastructure")

	// Cancel the infrastructure context to signal all components to shut down
	if infra.Cancel != nil {
		infra.Logger.Info("Cancelling infrastructure context")
		infra.Cancel()
	}

	// Create a separate timeout context for cleanup operations
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cleanupCancel()

	// Stop operator goroutines
	if len(infra.OperatorInstances) > 0 {
		infra.Logger.Info("Stopping operator goroutines")
		StopAllOperators(infra)
	}

	infra.Logger.Info("Stopping binaries")
	infra.TestConfig.StopBinaries()

	if infra.ChurnerServer != nil {
		infra.Logger.Info("Stopping churner server")
		infra.ChurnerServer.GracefulStop()
		if infra.ChurnerListener != nil {
			_ = infra.ChurnerListener.Close()
		}
	}

	if infra.AnvilContainer != nil {
		infra.Logger.Info("Stopping anvil")
		if err := infra.AnvilContainer.Terminate(cleanupCtx); err != nil {
			infra.Logger.Warn("Failed to terminate anvil container", "error", err)
		}
	}

	if infra.GraphNodeContainer != nil {
		infra.Logger.Info("Stopping graph node")
		_ = infra.GraphNodeContainer.Terminate(cleanupCtx)
	}

	if infra.ChainDockerNetwork != nil {
		infra.Logger.Info("Removing Docker network")
		_ = infra.ChainDockerNetwork.Remove(cleanupCtx)
	}

	if infra.LocalstackContainer != nil {
		infra.Logger.Info("Stopping localstack container")
		if err := infra.LocalstackContainer.Terminate(cleanupCtx); err != nil {
			infra.Logger.Warn("Failed to terminate localstack container", "error", err)
		}
	}

	infra.Logger.Info("Teardown completed")
}
