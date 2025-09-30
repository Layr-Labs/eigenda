package integration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/testcontainers/testcontainers-go/network"
)

// InfrastructureConfig contains the configuration for setting up the infrastructure
type InfrastructureConfig struct {
	TemplateName                    string
	TestName                        string
	InMemoryBlobStore               bool
	Logger                          logging.Logger
	RootPath                        string
	MetadataTableName               string
	BucketTableName                 string
	MetadataTableNameV2             string
	UserReservationSymbolsPerSecond uint64
	UserOnDemandDeposit             uint64
	ReservationPeriodInterval       uint64
	ClientLedgerMode                clientledger.ClientLedgerMode
	ControllerUseNewPayments        bool
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create test directory if needed
	testName := config.TestName
	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(config.TemplateName, config.RootPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create test directory: %w", err)
		}
	}

	testConfig := deploy.ReadTestConfig(testName, config.RootPath)
	testConfig.UserReservationSymbolsPerSecond = config.UserReservationSymbolsPerSecond
	testConfig.UserOnDemandDeposit = config.UserOnDemandDeposit
	testConfig.ReservationPeriodInterval = config.ReservationPeriodInterval
	testConfig.ClientLedgerMode = config.ClientLedgerMode
	testConfig.UseControllerMediatedPayments = config.ControllerUseNewPayments

	if testConfig.Environment.IsLocal() {
		return setupLocalInfrastructure(ctx, config, testName, testConfig)
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
	}, nil
}

func setupLocalInfrastructure(
	ctx context.Context,
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
	}

	// Create a shared Docker network for all containers
	var err error
	infra.ChainDockerNetwork, err = network.New(context.Background(),
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
		infra.LocalstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
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
		err = testbed.DeployResources(ctx, deployConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy resources: %w", err)
		}
	} else {
		infra.Logger.Info("Using in-memory Blob Store")
	}

	// Setup Anvil
	infra.Logger.Info("Starting anvil")
	infra.AnvilContainer, err = testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         infra.Logger,
		Network:        infra.ChainDockerNetwork,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil: %w", err)
	}

	// Setup Graph Node if needed
	deployer, ok := infra.TestConfig.GetDeployer(infra.TestConfig.EigenDA.Deployer)
	if ok && deployer.DeploySubgraphs {
		infra.Logger.Info("Starting graph node")
		anvilInternalEndpoint := infra.AnvilContainer.InternalEndpoint()
		infra.GraphNodeContainer, err = testbed.NewGraphNodeContainerWithOptions(
			context.Background(), testbed.GraphNodeOptions{
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

	// Start churner server
	infra.Logger.Info("Starting churner server")
	err = StartChurnerForInfrastructure(infra)
	if err != nil {
		return nil, fmt.Errorf("failed to start churner server: %w", err)
	}
	infra.Logger.Info("Churner server started", "port", "32002")

	infra.Logger.Info("Starting binaries")
	infra.TestConfig.StartBinaries(true) // true = for tests, will skip churner

	return infra, nil
}

// TeardownGlobalInfrastructure cleans up all global infrastructure
func TeardownGlobalInfrastructure(infra *InfrastructureHarness) {
	if infra == nil || infra.TestConfig == nil || !infra.TestConfig.Environment.IsLocal() {
		return
	}

	infra.Logger.Info("Tearing down global infrastructure")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

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
		if err := infra.AnvilContainer.Terminate(ctx); err != nil {
			infra.Logger.Warn("Failed to terminate anvil container", "error", err)
		}
	}

	if infra.GraphNodeContainer != nil {
		infra.Logger.Info("Stopping graph node")
		_ = infra.GraphNodeContainer.Terminate(context.Background())
	}

	if infra.ChainDockerNetwork != nil {
		infra.Logger.Info("Removing Docker network")
		_ = infra.ChainDockerNetwork.Remove(context.Background())
	}

	if infra.LocalstackContainer != nil {
		infra.Logger.Info("Stopping localstack container")
		if err := infra.LocalstackContainer.Terminate(ctx); err != nil {
			infra.Logger.Warn("Failed to terminate localstack container", "error", err)
		}
	}

	infra.Logger.Info("Teardown completed")
}
