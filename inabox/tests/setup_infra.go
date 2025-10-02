package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
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
	RootPath            string
	MetadataTableName   string
	BucketTableName     string
	S3BucketName        string
	MetadataTableNameV2 string

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int

	// The following fields are temporary, to be able to test different payments configurations. They will be removed
	// once legacy payments are removed.
	UserReservationSymbolsPerSecond uint64
	ClientLedgerMode                clientledger.ClientLedgerMode
	ControllerUseNewPayments        bool
}

// SetupInfrastructure creates the shared infrastructure that persists across all tests.
// This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
func SetupInfrastructure(ctx context.Context, config *InfrastructureConfig) (*InfrastructureHarness, error) {
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
	testConfig.ClientLedgerMode = config.ClientLedgerMode
	testConfig.UseControllerMediatedPayments = config.ControllerUseNewPayments

	// Create a long-lived context for the infrastructure lifecycle
	infraCtx, infraCancel := context.WithCancel(ctx)

	// Ensure we cancel the context if we return an error
	var setupErr error
	defer func() {
		if setupErr != nil && infraCancel != nil {
			infraCancel()
		}
	}()

	// Create shared Docker network, primarily for Anvil and Graph Node
	sharedDockerNetwork, err := network.New(
		infraCtx,
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
		Cancel:            infraCancel,
	}

	// Setup Chain Harness first (Anvil, Graph Node, Contracts, Churner)
	chainHarnessConfig := &ChainHarnessConfig{
		TestConfig: testConfig,
		TestName:   testName,
		Logger:     logger,
		Network:    sharedDockerNetwork,
	}
	chainHarness, err := SetupChainHarness(infraCtx, chainHarnessConfig)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup chain harness: %w", err)
		return nil, setupErr
	}
	infra.ChainHarness = *chainHarness

	// Setup Disperser Harness second (LocalStack, DynamoDB tables, S3 buckets, relays)
	disperserHarnessConfig := &DisperserHarnessConfig{
		Logger:              logger,
		Network:             sharedDockerNetwork,
		TestConfig:          testConfig,
		TestName:            testName,
		InMemoryBlobStore:   config.InMemoryBlobStore,
		LocalStackPort:      infra.LocalStackPort,
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		S3BucketName:        config.S3BucketName,
		MetadataTableNameV2: config.MetadataTableNameV2,
		EthClient:           infra.ChainHarness.EthClient,
		RelayCount:          config.RelayCount,
	}
	disperserHarness, err := SetupDisperserHarness(infraCtx, *disperserHarnessConfig)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup disperser harness: %w", err)
		return nil, setupErr
	}
	infra.DisperserHarness = *disperserHarness

	// Setup Operator Harness third (requires chain and disperser to be ready)
	operatorHarnessConfig := &OperatorHarnessConfig{
		TestConfig: testConfig,
		TestName:   testName,
		Logger:     logger,
	}
	operatorHarness, err := SetupOperatorHarness(infraCtx, operatorHarnessConfig, &infra.ChainHarness)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup operator harness: %w", err)
		return nil, setupErr
	}
	infra.OperatorHarness = *operatorHarness

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
