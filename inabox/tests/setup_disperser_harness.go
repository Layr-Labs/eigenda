package integration

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/testcontainers/testcontainers-go"
)

// DisperserHarnessConfig contains the configuration for setting up the disperser harness
type DisperserHarnessConfig struct {
	Logger              logging.Logger
	Network             *testcontainers.DockerNetwork
	TestConfig          *deploy.Config
	InMemoryBlobStore   bool
	LocalStackPort      string
	MetadataTableName   string
	BucketTableName     string
	MetadataTableNameV2 string
}

// TODO: Add encoder, api server, relay, controller, batcher
type DisperserHarness struct {
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadataV1 string
		BlobMetaV2     string
	}
	S3Buckets struct {
		BlobStore string
	}
}

// setupLocalStackResources initializes LocalStack and deploys AWS resources
func setupLocalStackResources(
	ctx context.Context, config *DisperserHarnessConfig,
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
		LocalStackEndpoint:  fmt.Sprintf("http://0.0.0.0:%s", config.LocalStackPort),
		MetadataTableName:   config.MetadataTableName,
		BucketTableName:     config.BucketTableName,
		V2MetadataTableName: config.MetadataTableNameV2,
		Logger:              config.Logger,
	}
	if err := testbed.DeployResources(ctx, deployConfig); err != nil {
		return nil, fmt.Errorf("failed to deploy resources: %w", err)
	}
	config.Logger.Info("AWS resources deployed successfully")

	return localstackContainer, nil
}

// setupDisperserKeypairAndRegistrations generates disperser keypair and performs registrations
func setupDisperserKeypairAndRegistrations(config *DisperserHarnessConfig) error {
	if config.TestConfig == nil {
		return nil
	}

	config.Logger.Info("Attempting to generate disperser keypair with LocalStack running")
	if err := config.TestConfig.GenerateDisperserKeypair(); err != nil {
		return fmt.Errorf("failed to generate disperser keypair: %w", err)
	}

	// Register blob versions, relays, and disperser keypair
	if config.TestConfig.EigenDA.Deployer != "" && config.TestConfig.IsEigenDADeployed() {
		config.TestConfig.PerformDisperserRegistrations()
	}

	return nil
}

// SetupDisperserHarness creates and initializes the disperser infrastructure (LocalStack, DynamoDB tables, S3 buckets)
func SetupDisperserHarness(ctx context.Context, config *DisperserHarnessConfig) (*DisperserHarness, error) {
	harness := &DisperserHarness{}

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
	if config.MetadataTableNameV2 == "" {
		config.MetadataTableNameV2 = "test-BlobMetadata-v2"
	}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetadataV1 = config.MetadataTableName
	harness.DynamoDBTables.BlobMetaV2 = config.MetadataTableNameV2
	harness.S3Buckets.BlobStore = config.BucketTableName

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
	} else {
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

// Cleanup releases resources held by the DisperserHarness (excluding shared network)
func (dh *DisperserHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	if dh.LocalStack != nil {
		logger.Info("Stopping localstack container")
		if err := dh.LocalStack.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate localstack container", "error", err)
		}
	}
}
