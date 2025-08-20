package deployment

import (
	"context"
	"fmt"
	"strings"

	eigendaaws "github.com/Layr-Labs/eigenda/common/aws"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// LocalStackDeploymentConfig holds configuration for AWS resource deployment in LocalStack
type LocalStackDeploymentConfig struct {
	BucketName          string `json:"bucket_name"`
	MetadataTableName   string `json:"metadata_table_name"`
	BucketTableName     string `json:"bucket_table_name"`
	V2MetadataTableName string `json:"v2_metadata_table_name"`
	V2PaymentPrefix     string `json:"v2_payment_prefix"`
	CreateV2Resources   bool   `json:"create_v2_resources"`
}

// DefaultLocalStackDeploymentConfig returns a default configuration for AWS resource deployment
func DefaultLocalStackDeploymentConfig() LocalStackDeploymentConfig {
	return LocalStackDeploymentConfig{
		BucketName:          "test-eigenda-blobstore",
		MetadataTableName:   "test-BlobMetadata",
		BucketTableName:     "test-BucketStore",
		V2MetadataTableName: "test-BlobMetadataV2",
		V2PaymentPrefix:     "test_v2_",
		CreateV2Resources:   true,
	}
}

// DeployLocalStackResources creates S3 buckets and DynamoDB tables in LocalStack
// This function replaces the DeployResources function from inabox/deploy/localstack.go
func DeployLocalStackResources(ctx context.Context, ls *containers.LocalStackContainer, cfg LocalStackDeploymentConfig) error {
	if ls == nil {
		return fmt.Errorf("localstack container is not initialized")
	}

	// Create AWS config for LocalStack (using same credentials as original script)
	awsConfig := eigendaaws.ClientConfig{
		Region:          ls.Region(),
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     ls.Endpoint(),
	}

	// Create S3 bucket
	if err := createS3Bucket(ctx, cfg.BucketName, awsConfig); err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	// Create DynamoDB tables
	if err := createLocalStackDynamoDBTables(ctx, cfg, awsConfig); err != nil {
		return fmt.Errorf("failed to create DynamoDB tables: %w", err)
	}

	return nil
}

// createS3Bucket creates an S3 bucket in LocalStack using the AWS SDK
func createS3Bucket(ctx context.Context, bucketName string, awsConfig eigendaaws.ClientConfig) error {
	fmt.Printf("Creating S3 bucket: %s\n", bucketName)

	// Create AWS SDK config with custom endpoint resolver for LocalStack
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AccessKey,
			awsConfig.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to create AWS config: %w", err)
	}

	// Create S3 client with LocalStack-specific configuration
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Use path-style addressing for LocalStack compatibility
		o.UsePathStyle = true
		// Set custom endpoint
		o.BaseEndpoint = &awsConfig.EndpointURL
		// Disable SSL verification for LocalStack
		o.EndpointOptions.DisableHTTPS = strings.HasPrefix(awsConfig.EndpointURL, "http://")
	})

	// Check if bucket already exists
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err == nil {
		fmt.Printf("Bucket %s already exists\n", bucketName)
		return nil
	}

	// Create bucket
	input := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}
	_, err = s3Client.CreateBucket(ctx, input)
	if err != nil {
		// Check if it's a BucketAlreadyExists error - that's okay
		if strings.Contains(err.Error(), "BucketAlreadyExists") ||
			strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			fmt.Printf("Bucket %s already exists\n", bucketName)
			return nil
		}
		return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	fmt.Printf("Successfully created S3 bucket: %s\n", bucketName)
	return nil
}

// createLocalStackDynamoDBTables creates all required DynamoDB tables in LocalStack
func createLocalStackDynamoDBTables(ctx context.Context, cfg LocalStackDeploymentConfig, awsConfig eigendaaws.ClientConfig) error {
	// Create metadata table for v1 blob storage
	fmt.Println("Creating v1 metadata table:", cfg.MetadataTableName)
	if _, err := test_utils.CreateTable(ctx, awsConfig, cfg.MetadataTableName, blobstore.GenerateTableSchema(cfg.MetadataTableName, 10, 10)); err != nil {
		return fmt.Errorf("failed to create metadata table: %w", err)
	}

	// Create bucket table for general storage
	fmt.Println("Creating bucket table:", cfg.BucketTableName)
	if _, err := test_utils.CreateTable(ctx, awsConfig, cfg.BucketTableName, store.GenerateTableSchema(10, 10, cfg.BucketTableName)); err != nil {
		return fmt.Errorf("failed to create bucket table: %w", err)
	}

	// Create v2 resources if enabled
	if cfg.CreateV2Resources && cfg.V2MetadataTableName != "" {
		fmt.Println("Creating v2 resources...")

		// Create v2 metadata table
		fmt.Println("Creating v2 metadata table:", cfg.V2MetadataTableName)
		if _, err := test_utils.CreateTable(ctx, awsConfig, cfg.V2MetadataTableName, blobstorev2.GenerateTableSchema(cfg.V2MetadataTableName, 10, 10)); err != nil {
			return fmt.Errorf("failed to create v2 metadata table: %w", err)
		}

		// Create payment-related tables
		if err := createPaymentTables(awsConfig, cfg.V2PaymentPrefix); err != nil {
			return fmt.Errorf("failed to create payment tables: %w", err)
		}
	}

	fmt.Println("Successfully created all DynamoDB tables")
	return nil
}

// createPaymentTables creates payment-related DynamoDB tables
func createPaymentTables(awsConfig eigendaaws.ClientConfig, paymentPrefix string) error {
	tables := []struct {
		name   string
		create func(eigendaaws.ClientConfig, string) error
	}{
		{"reservation", meterer.CreateReservationTable},
		{"ondemand", meterer.CreateOnDemandTable},
		{"global_reservation", meterer.CreateGlobalReservationTable},
	}

	for _, table := range tables {
		tableName := paymentPrefix + table.name
		fmt.Printf("Creating payment table: %s\n", tableName)
		if err := table.create(awsConfig, tableName); err != nil {
			return fmt.Errorf("failed to create %s table: %w", tableName, err)
		}
	}

	return nil
}

// GetLocalStackAWSClientConfig returns an AWS client configuration for the LocalStack instance
// This is a convenience function for tests that need to create AWS clients
func GetLocalStackAWSClientConfig(ls *containers.LocalStackContainer) eigendaaws.ClientConfig {
	return eigendaaws.ClientConfig{
		Region:          ls.Region(),
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     ls.Endpoint(),
	}
}