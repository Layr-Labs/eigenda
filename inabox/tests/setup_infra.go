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

// SetupGlobalInfrastructure creates the shared infrastructure that persists across all tests.
// This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
// This should be called once in TestMain.
func SetupGlobalInfrastructure(templateName, testName string, inMemoryBlobStore bool, logger logging.Logger) (*InfrastructureHarness, error) {
	infra := &InfrastructureHarness{
		TemplateName:        templateName,
		TestName:            testName,
		InMemoryBlobStore:   inMemoryBlobStore,
		Logger:              logger,
		MetadataTableName:   "test-BlobMetadata",
		BucketTableName:     "test-BucketStore",
		MetadataTableNameV2: "test-BlobMetadata-v2",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	infra.Cancel = cancel
	defer func() {
		// Don't cancel here - the context is needed for the lifetime of the infra
		// It will be cancelled in TeardownGlobalInfrastructure
	}()

	rootPath := "../../"

	// Create test directory if needed
	if infra.TestName == "" {
		var err error
		infra.TestName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create test directory: %w", err)
		}
	}

	infra.TestConfig = deploy.ReadTestConfig(infra.TestName, rootPath)

	if infra.TestConfig.Environment.IsLocal() {
		if err := setupLocalInfrastructure(ctx, infra); err != nil {
			cancel()
			return nil, err
		}
	}

	return infra, nil
}

func setupLocalInfrastructure(ctx context.Context, infra *InfrastructureHarness) error {
	// Create a shared Docker network for all containers
	var err error
	infra.ChainDockerNetwork, err = network.New(context.Background(),
		network.WithDriver("bridge"),
		network.WithAttachable())
	if err != nil {
		return fmt.Errorf("failed to create docker network: %w", err)
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
			return fmt.Errorf("failed to start localstack: %w", err)
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
			return fmt.Errorf("failed to deploy resources: %w", err)
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
		return fmt.Errorf("failed to start anvil: %w", err)
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
			return fmt.Errorf("failed to start graph node: %w", err)
		}
	}

	// Deploy contracts
	infra.Logger.Info("Deploying experiment")
	err = infra.TestConfig.DeployExperiment()
	if err != nil {
		return fmt.Errorf("failed to deploy experiment: %w", err)
	}

	// Start churner server
	infra.Logger.Info("Starting churner server")
	err = StartChurnerForInfrastructure(infra)
	if err != nil {
		return fmt.Errorf("failed to start churner server: %w", err)
	}
	infra.Logger.Info("Churner server started", "port", "32002")

	infra.Logger.Info("Starting binaries")
	infra.TestConfig.StartBinaries(true) // true = for tests, will skip churner

	return nil
}

// TeardownGlobalInfrastructure cleans up all global infrastructure
func TeardownGlobalInfrastructure(infra *InfrastructureHarness) {
	if infra == nil || infra.TestConfig == nil || !infra.TestConfig.Environment.IsLocal() {
		return
	}

	infra.Logger.Info("Tearing down global infrastructure")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if infra.Cancel != nil {
		infra.Cancel()
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
