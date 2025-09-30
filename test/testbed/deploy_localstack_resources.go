package testbed

import (
	"context"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DeployResourcesConfig holds configuration for deploying AWS resources
type DeployResourcesConfig struct {
	LocalStackEndpoint  string
	MetadataTableName   string
	BucketTableName     string
	BucketName          string // Optional: S3 bucket name, defaults to "test-eigenda-blobstore"
	V2MetadataTableName string
	V2PaymentPrefix     string         // Optional: prefix for v2 payment tables, defaults to "e2e_v2_"
	Region              string         // Optional: AWS region, defaults to "us-east-1"
	AccessKey           string         // Optional: AWS access key, defaults to "localstack"
	SecretAccessKey     string         // Optional: AWS secret key, defaults to "localstack"
	Logger              logging.Logger // Optional: logger for output messages
}

// DeployResources creates AWS resources (S3 buckets and DynamoDB tables) on LocalStack
func DeployResources(ctx context.Context, config DeployResourcesConfig) error {
	// Use a default logger if none provided
	logger := config.Logger
	if logger == nil {
		loggerConfig := &common.LoggerConfig{
			Format:       common.TextLogFormat,
			OutputWriter: os.Stdout,
		}
		var err error
		logger, err = common.NewLogger(loggerConfig)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}
	}

	// Add component to logger
	logger = logger.With("component", "DeployResources")

	// Set defaults
	if config.Region == "" {
		config.Region = "us-east-1"
	}
	if config.AccessKey == "" {
		config.AccessKey = "localstack"
	}
	if config.SecretAccessKey == "" {
		config.SecretAccessKey = "localstack"
	}
	if config.V2PaymentPrefix == "" {
		config.V2PaymentPrefix = "e2e_v2_"
	}
	if config.BucketName == "" {
		config.BucketName = "test-eigenda-blobstore"
	}

	// Create AWS client config
	cfg := aws.ClientConfig{
		Region:          config.Region,
		AccessKey:       config.AccessKey,
		SecretAccessKey: config.SecretAccessKey,
		EndpointURL:     config.LocalStackEndpoint,
	}

	// Create S3 bucket
	if err := createS3Bucket(ctx, cfg, config.BucketName, logger); err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	// Create metadata table
	if config.MetadataTableName != "" {
		_, err := test_utils.CreateTable(ctx, cfg, config.MetadataTableName,
			blobstore.GenerateTableSchema(config.MetadataTableName, 10, 10))
		if err != nil {
			return fmt.Errorf("failed to create metadata table %s: %w", config.MetadataTableName, err)
		}
		logger.Info("Created metadata table", "table", config.MetadataTableName)
	}

	// Create bucket table
	if config.BucketTableName != "" {
		_, err := test_utils.CreateTable(ctx, cfg, config.BucketTableName,
			store.GenerateTableSchema(10, 10, config.BucketTableName))
		if err != nil {
			return fmt.Errorf("failed to create bucket table %s: %w", config.BucketTableName, err)
		}
		logger.Info("Created bucket table", "table", config.BucketTableName)
	}

	// Create v2 tables if specified
	if config.V2MetadataTableName != "" {
		logger.Info("Creating v2 tables")

		// Create v2 metadata table
		_, err := test_utils.CreateTable(ctx, cfg, config.V2MetadataTableName,
			blobstorev2.GenerateTableSchema(config.V2MetadataTableName, 10, 10))
		if err != nil {
			return fmt.Errorf("failed to create v2 metadata table %s: %w", config.V2MetadataTableName, err)
		}
		logger.Info("Created v2 metadata table", "table", config.V2MetadataTableName)

		// Create payment related tables
		if err := createPaymentTables(cfg, config.V2PaymentPrefix, logger); err != nil {
			return fmt.Errorf("failed to create payment tables: %w", err)
		}
	}

	return nil
}

// createS3Bucket creates the S3 bucket using the AWS SDK
func createS3Bucket(ctx context.Context, cfg aws.ClientConfig, bucketName string, logger logging.Logger) error {

	// Create AWS SDK config with custom endpoint resolver
	customResolver := awssdk.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (awssdk.Endpoint, error) {
			if cfg.EndpointURL != "" {
				return awssdk.Endpoint{
					PartitionID:   "aws",
					URL:           cfg.EndpointURL,
					SigningRegion: cfg.Region,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to its default resolution
			return awssdk.Endpoint{}, &awssdk.EndpointNotFoundError{}
		})

	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")),
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with path-style addressing for LocalStack
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Check if bucket already exists
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucketName,
	})

	if err == nil {
		logger.Info("Bucket already exists", "bucket", bucketName)
		return nil
	}

	// Create the bucket
	createBucketConfig := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	// Only add LocationConstraint for non us-east-1 regions
	if cfg.Region != "us-east-1" {
		createBucketConfig.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(cfg.Region),
		}
	}

	_, err = s3Client.CreateBucket(ctx, createBucketConfig)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	logger.Info("Created S3 bucket", "bucket", bucketName)
	return nil
}

// createPaymentTables creates the payment-related tables
func createPaymentTables(cfg aws.ClientConfig, prefix string, logger logging.Logger) error {
	// Create reservation table
	if err := meterer.CreateReservationTable(cfg, prefix+"reservation"); err != nil {
		return fmt.Errorf("failed to create reservation table: %w", err)
	}
	logger.Info("Created reservation table", "table", prefix+"reservation")

	// Create on-demand table
	if err := meterer.CreateOnDemandTable(cfg, prefix+"ondemand"); err != nil {
		return fmt.Errorf("failed to create on-demand table: %w", err)
	}
	logger.Info("Created on-demand table", "table", prefix+"ondemand")

	// Create global reservation table
	if err := meterer.CreateGlobalReservationTable(cfg, prefix+"global_reservation"); err != nil {
		return fmt.Errorf("failed to create global reservation table: %w", err)
	}
	logger.Info("Created global reservation table", "table", prefix+"global_reservation")

	return nil
}

// DeployResourcesWithContainer is a convenience function that uses a LocalStackContainer
func DeployResourcesWithContainer(ctx context.Context,
	container *LocalStackContainer, config DeployResourcesConfig) error {
	// Override the endpoint with the container's endpoint
	config.LocalStackEndpoint = container.Endpoint()
	return DeployResources(ctx, config)
}
