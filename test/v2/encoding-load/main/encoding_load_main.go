package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	encodingload "github.com/Layr-Labs/eigenda/test/v2/encoding-load"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("Expected 2 args, got %d. Usage: %s <encoding_load_file>.\n"+
			"If '-' is passed in lieu of a config file, the config file path is read from the environment variable "+
			"$GENERATOR_ENCODING_LOAD.\n",
			len(os.Args), os.Args[0]))
	}

	loadFile := os.Args[1]
	if loadFile == "-" {
		loadFile = os.Getenv("GENERATOR_ENCODING_LOAD")
		if loadFile == "" {
			panic("$GENERATOR_ENCODING_LOAD not set")
		}
	}

	// Initialize logger
	loggerConfig := common.DefaultConsoleLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}

	// Initialize metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create the encoder client v2
	encoderURL := os.Getenv("ENCODER_URL")
	if encoderURL == "" {
		encoderURL = "localhost:8090" // Default encoder URL
		logger.Info("Using default encoder URL", "url", encoderURL)
	} else {
		logger.Info("Using encoder URL from environment", "url", encoderURL)
	}

	// Create the encoder client v2
	encoderClientV2, err := encoder.NewEncoderClientV2(encoderURL)
	if err != nil {
		panic(fmt.Errorf("failed to create encoder client v2: %w", err))
	}

	// Test the encoder client connection
	logger.Info("Successfully created encoder client v2", "client", encoderClientV2)

	// Read the config file first
	config, err := encodingload.ReadConfigFile(loadFile)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	// Configure S3 bucket
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "test-eigenda-blobstore" // Default S3 bucket
		logger.Info("Using default S3 bucket", "bucket", s3Bucket)
	} else {
		logger.Info("Using S3 bucket from environment", "bucket", s3Bucket)
	}
	os.Setenv("DISPERSER_SERVER_S3_BUCKET_NAME", s3Bucket)

	// Create a blob store for storing blobs
	awsConfig := &aws.ClientConfig{
		Region: os.Getenv("AWS_REGION"),
	}
	if awsConfig.Region == "" {
		awsConfig.Region = "us-east-1" // Default region
		logger.Info("Using default AWS region", "region", awsConfig.Region)
	}

	s3Client, err := s3.NewClient(context.Background(), *awsConfig, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create S3 client: %w", err))
	}

	blobStore := blobstorev2.NewBlobStore(s3Bucket, s3Client, logger)

	// Create the encoding load generator
	generator := encodingload.NewEncodingLoadGenerator(config, logger, metricsRegistry, encoderClientV2, blobStore)

	// Set up signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		logger.Info("Shutting down encoding load generator...")
		generator.Stop()
	}()

	logger.Info("Starting encoding load generator...")
	generator.Start(true)
}
