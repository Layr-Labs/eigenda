package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/testcontainers/testcontainers-go"
)

// DisperserV1HarnessConfig contains the configuration for setting up the v1 disperser harness
type DisperserV1HarnessConfig struct {
	Logger              logging.Logger
	Network             *testcontainers.DockerNetwork
	TestConfig          *deploy.Config
	TestName            string
	InMemoryBlobStore   bool
	BucketTableName     string
	V1MetadataTableName string
	BlobStoreBucketName string
	EthClient           common.EthClient
}

// DisperserV1Harness is a simpler harness that only uses v1 encoder (no v2, no relays)
type DisperserV1Harness struct {
	LocalStack     *testbed.LocalStackContainer
	DynamoDBTables struct {
		BlobMetadataV1 string
	}
	S3Buckets struct {
		BlobStore string
	}
	EncoderV1Instance *EncoderV1Instance
}

// EncoderV1Instance holds the state for a single encoder v1
type EncoderV1Instance struct {
	Server   *encoder.EncoderServer
	Listener net.Listener
	Port     string
	URL      string
	Logger   logging.Logger
}

// SetupDisperserV1Harness creates and initializes the v1 disperser infrastructure
// (LocalStack, DynamoDB tables, S3 buckets, encoder v1 goroutine).
func SetupDisperserV1Harness(
	ctx context.Context, localstack *testbed.LocalStackContainer, config DisperserV1HarnessConfig,
) (*DisperserV1Harness, error) {
	// Check if localstack resources are empty
	if config.V1MetadataTableName == "" || config.BucketTableName == "" || config.BlobStoreBucketName == "" {
		return nil, fmt.Errorf("missing name for localstack resources")
	}

	harness := &DisperserV1Harness{}

	// Populate the harness tables and buckets metadata
	harness.DynamoDBTables.BlobMetadataV1 = config.V1MetadataTableName
	harness.S3Buckets.BlobStore = config.BlobStoreBucketName

	// Setup LocalStack if not using in-memory blob store
	if !config.InMemoryBlobStore {
		// Setup LocalStack resources (reuse the same function from main harness)
		harnessConfig := DisperserV1HarnessConfig{
			Logger:              config.Logger,
			Network:             config.Network,
			TestConfig:          config.TestConfig,
			TestName:            config.TestName,
			V1MetadataTableName: config.V1MetadataTableName,
			BucketTableName:     config.BucketTableName,
			BlobStoreBucketName: config.BlobStoreBucketName,
			EthClient:           config.EthClient,
		}

		localstack, err := setupV1LocalStackResources(ctx, localstack, harnessConfig)
		if err != nil {
			return nil, err
		}
		harness.LocalStack = localstack

		// Start encoder v1 instance as a goroutine
		config.Logger.Info("Starting encoder v1 instance")
		encoderInstance, err := startEncoderV1(ctx, harness, config)
		if err != nil {
			return nil, fmt.Errorf("failed to start encoder v1: %w", err)
		}
		harness.EncoderV1Instance = encoderInstance
	} else {
		config.Logger.Info("Using in-memory blob store, skipping LocalStack setup")
	}

	return harness, nil
}

// setupLocalStackResources initializes LocalStack and deploys AWS resources
func setupV1LocalStackResources(
	ctx context.Context, localstack *testbed.LocalStackContainer, config DisperserV1HarnessConfig,
) (*testbed.LocalStackContainer, error) {
	// Deploy AWS resources (DynamoDB tables and S3 buckets)
	config.Logger.Info("Deploying AWS resources in LocalStack")
	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  localstack.Endpoint(),
		V1MetadataTableName: config.V1MetadataTableName,
		BucketTableName:     config.BucketTableName,
		BlobStoreBucketName: config.BlobStoreBucketName,
		AWSConfig:           localstack.GetAWSClientConfig(),
		Logger:              config.Logger,
	}
	if err := testbed.DeployResources(ctx, deployConfig); err != nil {
		return nil, fmt.Errorf("failed to deploy resources: %w", err)
	}
	config.Logger.Info("AWS resources deployed successfully")

	return localstack, nil
}

// startEncoderV1 starts a single encoder v1 instance
func startEncoderV1(
	ctx context.Context,
	harness *DisperserV1Harness,
	config DisperserV1HarnessConfig,
) (*EncoderV1Instance, error) {
	config.Logger.Info("Starting encoder v1 instance")

	// Get SRS paths using the same function as operator setup
	g1Path, g2Path, err := getSRSPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Pre-create listener with port 0 (OS assigns port)
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener for encoder v1: %w", err)
	}

	// Extract the actual port assigned by the OS
	actualPort := listener.Addr().(*net.TCPAddr).Port
	encoderURL := fmt.Sprintf("0.0.0.0:%d", actualPort)
	port := fmt.Sprintf("%d", actualPort)

	config.Logger.Info("Created listener for encoder v1", "assigned_port", actualPort)

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/encoder_v1.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to open encoder v1 log file: %w", err)
	}

	// Create encoder logger
	loggerConfig := common.LoggerConfig{
		Format:       common.TextLogFormat,
		OutputWriter: io.MultiWriter(os.Stdout, logFile),
		HandlerOpts: logging.SLoggerOptions{
			Level:     slog.LevelDebug,
			NoColor:   true,
			AddSource: true,
		},
	}

	encoderLogger, err := common.NewLogger(&loggerConfig)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create encoder v1 logger: %w", err)
	}

	// Create prover with dynamically determined paths
	// TODO(dmanc): Make these configurable
	kzgConfig := &kzg.KzgConfig{
		SRSOrder:        10000,
		SRSNumberToLoad: 10000,
		G1Path:          g1Path,
		G2Path:          g2Path,
		CacheDir:        fmt.Sprintf("testdata/%s/cache/encoder", config.TestName),
		NumWorker:       1,
		Verbose:         false,
		LoadG2Points:    true, // v1 encoder needs G2 points
	}
	encodingConfig := &encoding.Config{
		BackendType: encoding.GnarkBackend,
		GPUEnable:   false,
		NumWorker:   1,
	}

	proverV1, err := prover.NewProver(kzgConfig, encodingConfig)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to create prover v1: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(metricsRegistry, "9100", encoderLogger)
	grpcMetrics := grpcprom.NewServerMetrics()
	metricsRegistry.MustRegister(grpcMetrics)

	// Create encoder server configuration
	serverConfig := encoder.ServerConfig{
		GrpcPort:              port,
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}

	// Create encoder server
	server := encoder.NewEncoderServer(
		serverConfig,
		encoderLogger,
		proverV1,
		metrics,
		grpcMetrics,
	)

	// Start the encoder server in a goroutine using the pre-created listener
	go func() {
		encoderLogger.Info("Starting encoder v1 server with listener", "port", port)
		if err := server.StartWithListener(listener); err != nil {
			encoderLogger.Error("Encoder v1 server failed", "error", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	encoderLogger.Info("Encoder v1 server started successfully", "port", port, "logFile", logFilePath)

	return &EncoderV1Instance{
		Server:   server,
		Listener: listener,
		Port:     port,
		URL:      encoderURL,
		Logger:   encoderLogger,
	}, nil
}

// Cleanup releases resources held by the DisperserV1Harness
func (dh *DisperserV1Harness) Cleanup(ctx context.Context, logger logging.Logger) {
	// Stop encoder v1 instance
	if dh.EncoderV1Instance != nil {
		logger.Info("Stopping encoder v1 instance")
		dh.EncoderV1Instance.Logger.Info("Stopping encoder v1")
		// TODO: Add Close method to encoder v1 server
	}
}
