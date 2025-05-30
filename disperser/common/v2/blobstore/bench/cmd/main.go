package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore/bench"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	logger logging.Logger
)

// Command line flags
var (
	storeType      = flag.String("store", "dynamodb", "Store type: dynamodb or postgresql")
	operations     = flag.String("ops", "UpdateBlobStatus:200:30s", "Operations to benchmark (format: op1:rate:duration,op2:rate:duration)")
	warmupTime     = flag.Duration("warmup", 10*time.Second, "Warmup duration")
	reportInterval = flag.Duration("report", 5*time.Second, "Reporting interval")
	workers        = flag.Int("workers", 10, "Number of workers per operation")

	// DynamoDB flags
	dynamoTable    = flag.String("dynamo-table", "test-metadata", "DynamoDB table name")
	dynamoRegion   = flag.String("dynamo-region", "us-east-1", "AWS region")
	dynamoEndpoint = flag.String("dynamo-endpoint", "", "DynamoDB endpoint (for local testing)")

	// PostgreSQL flags
	pgHost     = flag.String("pg-host", "localhost", "PostgreSQL host")
	pgPort     = flag.Int("pg-port", 5432, "PostgreSQL port")
	pgUser     = flag.String("pg-user", "postgres", "PostgreSQL username")
	pgPassword = flag.String("pg-password", "", "PostgreSQL password")
	pgDatabase = flag.String("pg-database", "eigenda_benchmark", "PostgreSQL database")
	pgSSLMode  = flag.String("pg-sslmode", "disable", "PostgreSQL SSL mode")

	// Logging
	logLevel = flag.String("log-level", "info", "Log level: debug, info, warn, error")
)

func main() {
	flag.Parse()

	// Setup logger
	loggerConfig := common.DefaultTextLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Parse operations
	opConfigs, err := parseOperations(*operations)
	if err != nil {
		logger.Fatal("Failed to parse operations", "error", err)
	}

	// Create benchmark config
	benchConfig := bench.BenchmarkConfig{
		StoreType:      bench.StoreType(*storeType),
		Operations:     opConfigs,
		WarmupTime:     *warmupTime,
		ReportInterval: *reportInterval,
	}

	// Setup store-specific configuration
	switch benchConfig.StoreType {
	case bench.DynamoDB:
		clientConfig := commonaws.ClientConfig{
			Region:      *dynamoRegion,
			EndpointURL: *dynamoEndpoint, // Only set for local testing
		}
		// Only use local credentials when endpoint is specified (local testing)
		if *dynamoEndpoint != "" {
			clientConfig.AccessKey = "local"
			clientConfig.SecretAccessKey = "local"
		}
		dynamoClient, err := commondynamodb.NewClient(clientConfig, logger)
		if err != nil {
			logger.Fatal("Failed to create DynamoDB client", "error", err)
		}
		benchConfig.DynamoDBConfig = &blobstore.DynamoDBConfig{
			Client:    dynamoClient,
			TableName: *dynamoTable,
			Region:    *dynamoRegion,
			Endpoint:  *dynamoEndpoint,
		}
	case bench.PostgreSQL:
		benchConfig.PostgresConfig = &blobstore.PostgresConfig{
			Host:     *pgHost,
			Port:     *pgPort,
			Username: *pgUser,
			Password: *pgPassword,
			Database: *pgDatabase,
			SSLMode:  *pgSSLMode,
		}
	default:
		logger.Fatal("Invalid store type", "store", *storeType)
	}

	// Create metadata store
	store, err := createMetadataStore(benchConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create metadata store", "error", err)
	}

	// Create and run benchmark
	runner := bench.NewBenchmarkRunner(benchConfig, store, logger)

	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		logger.Fatal("Benchmark failed", "error", err)
	}
}

// parseOperations parses the operations string into OperationConfig slice
func parseOperations(opsStr string) ([]bench.OperationConfig, error) {
	if opsStr == "" {
		return nil, fmt.Errorf("no operations specified")
	}

	var configs []bench.OperationConfig
	ops := strings.Split(opsStr, ",")

	for _, op := range ops {
		parts := strings.Split(op, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid operation format: %s (expected op:rate:duration)", op)
		}

		opType := bench.OperationType(parts[0])

		var rate int
		_, err := fmt.Sscanf(parts[1], "%d", &rate)
		if err != nil {
			return nil, fmt.Errorf("invalid rate for %s: %s", parts[0], parts[1])
		}

		duration, err := time.ParseDuration(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid duration for %s: %s", parts[0], parts[2])
		}

		configs = append(configs, bench.OperationConfig{
			Type:       opType,
			RatePerSec: rate,
			Duration:   duration,
			Workers:    *workers,
		})
	}

	return configs, nil
}

// createMetadataStore creates the appropriate metadata store based on configuration
func createMetadataStore(config bench.BenchmarkConfig, logger logging.Logger) (blobstore.MetadataStore, error) {
	switch config.StoreType {
	case bench.DynamoDB:
		return createDynamoDBStore(config.DynamoDBConfig, logger)
	case bench.PostgreSQL:
		return createPostgreSQLStore(config.PostgresConfig, logger)
	default:
		return nil, fmt.Errorf("unsupported store type: %s", config.StoreType)
	}
}

// createDynamoDBStore creates a DynamoDB metadata store
func createDynamoDBStore(cfg *blobstore.DynamoDBConfig, logger logging.Logger) (blobstore.MetadataStore, error) {
	// Create store
	store := blobstore.NewDynamoDBBlobMetadataStore(logger, *cfg)

	// Optionally create table if it doesn't exist (for local testing)
	if cfg.Endpoint != "" {
		if err := createDynamoDBTableIfNotExists(cfg.Client, cfg.TableName); err != nil {
			return nil, fmt.Errorf("failed to create DynamoDB table: %w", err)
		}
	}

	return store, nil
}

// createPostgreSQLStore creates a PostgreSQL metadata store
func createPostgreSQLStore(cfg *blobstore.PostgresConfig, logger logging.Logger) (blobstore.MetadataStore, error) {
	store, err := blobstore.NewPostgresBlobMetadataStore(logger, *cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL store: %w", err)
	}

	// Initialize tables
	if err := store.ResetTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL tables: %w", err)
	}

	return store, nil
}

// createDynamoDBTableIfNotExists creates the DynamoDB table if it doesn't exist
func createDynamoDBTableIfNotExists(client commondynamodb.Client, tableName string) error {
	ctx := context.Background()

	// Check if table exists
	err := client.TableExists(ctx, tableName)
	if err == nil {
		// Table exists
		return nil
	}

	// Create table using test_utils
	input := blobstore.GenerateTableSchema(tableName, 100, 100)
	createConfig := commonaws.ClientConfig{
		Region:      *dynamoRegion,
		EndpointURL: *dynamoEndpoint,
	}
	// Only use local credentials when endpoint is specified (local testing)
	if *dynamoEndpoint != "" {
		createConfig.AccessKey = "local"
		createConfig.SecretAccessKey = "local"
	}
	_, err = test_utils.CreateTable(ctx, createConfig, tableName, input)
	return err
}
