package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
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

// SetupInfrastructure creates the shared infrastructure that persists across all tests.
// This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
func SetupInfrastructure(config *InfrastructureConfig) (*InfrastructureHarness, error) {
	if config.MetadataTableName == "" {
		config.MetadataTableName = "test-BlobMetadata"
	}
	if config.BucketTableName == "" {
		config.BucketTableName = "test-BucketStore"
	}
	if config.MetadataTableNameV2 == "" {
		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
	}

	logger := config.Logger

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

	// Ensure we cancel the context if we return an error
	var setupErr error
	defer func() {
		if setupErr != nil && infraCancel != nil {
			infraCancel()
		}
	}()

	// Create a shared Docker network for all containers
	sharedDockerNetwork, err := network.New(
		setupCtx,
		network.WithDriver("bridge"),
		network.WithAttachable())
	if err != nil {
		setupErr = fmt.Errorf("failed to create docker network: %w", err)
		return nil, setupErr
	}
	logger.Info("Created Docker network", "name", sharedDockerNetwork.Name)

	// Create infrastructure harness early so we can populate it incrementally
	infra := &InfrastructureHarness{
		SharedNetwork:     sharedDockerNetwork,
		TestConfig:        testConfig,
		TemplateName:      config.TemplateName,
		TestName:          testName,
		InMemoryBlobStore: config.InMemoryBlobStore,
		LocalStackPort:    "4570",
		Logger:            config.Logger,
		Ctx:               infraCtx,
		Cancel:            infraCancel,
	}

	// Setup Chain Harness first (Anvil, Graph Node, contracts, Churner)
	chainHarnessConfig := &ChainHarnessConfig{
		TestConfig: testConfig,
		TestName:   testName,
		Logger:     logger,
		Network:    sharedDockerNetwork,
	}

	chainHarness, err := SetupChainHarness(setupCtx, chainHarnessConfig)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup chain harness: %w", err)
		return nil, setupErr
	}
	infra.ChainHarness = *chainHarness

	// Setup Operator Harness second (requires chain to be ready)
	operatorHarnessConfig := &OperatorHarnessConfig{
		TestConfig:   testConfig,
		TestName:     testName,
		Logger:       logger,
		ChainHarness: &infra.ChainHarness,
		Ctx:          infraCtx,
	}
	operatorHarness, err := SetupOperatorHarness(setupCtx, operatorHarnessConfig)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup operator harness: %w", err)
		return nil, setupErr
	}
	infra.OperatorHarness = *operatorHarness

	// Setup Disperser Harness third (LocalStack, DynamoDB tables, S3 buckets)
	disperserHarnessConfig := &DisperserHarnessConfig{
		Logger:              logger,
		Network:             sharedDockerNetwork,
		TestConfig:          testConfig,
		InMemoryBlobStore:   config.InMemoryBlobStore,
		LocalStackPort:      infra.LocalStackPort,
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		MetadataTableNameV2: config.MetadataTableNameV2,
	}

	disperserHarness, err := SetupDisperserHarness(setupCtx, disperserHarnessConfig)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup disperser harness: %w", err)
		return nil, setupErr
	}
	infra.DisperserHarness = *disperserHarness

	return infra, nil
}

// TeardownGlobalInfrastructure cleans up all global infrastructure
func TeardownInfrastructure(infra *InfrastructureHarness) {
	infra.Logger.Info("Tearing down global infrastructure")

	// Cancel the infrastructure context to signal all components to shut down
	if infra.Cancel != nil {
		infra.Logger.Info("Cancelling infrastructure context")
		infra.Cancel()
	}

	// Create a separate timeout context for cleanup operations
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cleanupCancel()

	// Stop operator goroutines using the harness cleanup
	infra.OperatorHarness.Cleanup(cleanupCtx, infra.Logger)

	// Stop test binaries
	infra.Logger.Info("Stopping binaries")
	infra.TestConfig.StopBinaries()

	// Clean up disperser harness (graph node and localstack)
	infra.DisperserHarness.Cleanup(cleanupCtx, infra.Logger)

	// Clean up chain harness (churner and anvil)
	infra.ChainHarness.Cleanup(cleanupCtx, infra.Logger)

	// Clean up the shared Docker network last since multiple harnesses use it
	if infra.SharedNetwork != nil {
		infra.Logger.Info("Removing shared Docker network")
		_ = infra.SharedNetwork.Remove(cleanupCtx)
	}

	infra.Logger.Info("Teardown completed")
}
