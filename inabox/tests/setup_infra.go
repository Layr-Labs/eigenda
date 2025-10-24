package integration

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
	TemplateName string
	TestName     string
	Logger       logging.Logger

	// V2 disperser related configuration
	V2MetadataTableName string
	BlobStoreBucketName string

	// Number of relay instances to start, if not specified, no relays will be started.
	RelayCount int
}

// SetupInfrastructure creates the shared infrastructure that persists across all tests.
// This includes containers for Anvil, LocalStack, GraphNode, and the Churner server.
func SetupInfrastructure(ctx context.Context, config *InfrastructureConfig) (*InfrastructureHarness, error) {
	if config.V2MetadataTableName == "" {
		config.V2MetadataTableName = "test-BlobMetadata-v2"
	}
	if config.BlobStoreBucketName == "" {
		config.BlobStoreBucketName = "test-eigenda-blobstore"
	}

	logger := config.Logger

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
		SharedNetwork: sharedDockerNetwork,
		TestConfig:    testConfig,
		TemplateName:  config.TemplateName,
		TestName:      testName,
		Logger:        config.Logger,
		Cancel:        infraCancel,
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

	// Setup a shared localstack container
	sharedLocalStack, err := testbed.NewLocalStackContainerWithOptions(
		infraCtx,
		testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       infra.LocalStackPort,
			Logger:         logger,
		})
	if err != nil {
		setupErr = fmt.Errorf("failed to create shared localstack container: %w", err)
		return nil, setupErr
	}

	// Setup Disperser Harness (V2) (DynamoDB tables, S3 buckets, relays, encoders)
	disperserHarnessConfig := &DisperserHarnessConfig{
		TestConfig: testConfig,
		TestName:   testName,
		RelayCount: config.RelayCount,

		// LocalStack resources
		BlobStoreBucketName: config.BlobStoreBucketName,
		V2MetadataTableName: config.V2MetadataTableName,
	}
	disperserHarness, err := SetupDisperserHarness(
		infraCtx,
		logger,
		infra.ChainHarness.EthClient,
		sharedLocalStack,
		*disperserHarnessConfig,
	)
	if err != nil {
		setupErr = fmt.Errorf("failed to setup disperser harness: %w", err)
		return nil, setupErr
	}
	infra.DisperserHarness = *disperserHarness

	// Start remaining binaries (disperser, batcher, etc.)
	// TODO(dmanc): Once all of these components are migrated to goroutines, we can remove this.
	if testConfig != nil {
		config.Logger.Info("Starting remaining binaries")
		// Get encoder addresses, using empty string if instances are nil
		// TODO(dmanc): This is a hack to get the tests to pass when using in-memory blob store, we should refactor this.
		encoderV1Address := "" // V1 disperser no longer supported
		encoderV2Address := ""
		if disperserHarness.EncoderV2Instance != nil {
			encoderV2Address = disperserHarness.EncoderV2Instance.URL
		}
		err := testConfig.GenerateAllVariables(
			encoderV1Address,
			encoderV2Address,
		)
		if err != nil {
			return nil, fmt.Errorf("could not generate environment variables: %w", err)
		}
		testConfig.StartBinaries(true) // true = for tests, will skip churner and operators
	}

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

	// Stop test binaries
	infra.Logger.Info("Stopping binaries")
	infra.TestConfig.StopBinaries()

	// Stop operator goroutines using the harness cleanup
	infra.OperatorHarness.Cleanup(cleanupCtx, infra.Logger)

	// Clean up disperser harness (graph node and localstack)
	infra.DisperserHarness.Cleanup(cleanupCtx, infra.Logger)

	// Clean up localstack
	if err := infra.SharedLocalStack.Terminate(cleanupCtx); err != nil {
		infra.Logger.Error("Failed to terminate localstack container", "error", err)
	}

	// Clean up chain harness (churner and anvil)
	infra.ChainHarness.Cleanup(cleanupCtx, infra.Logger)

	// Clean up the shared Docker network last since multiple harnesses use it
	if infra.SharedNetwork != nil {
		infra.Logger.Info("Removing shared Docker network")
		_ = infra.SharedNetwork.Remove(cleanupCtx)
	}

	infra.Logger.Info("Teardown completed")
}
