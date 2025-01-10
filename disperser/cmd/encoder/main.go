package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/disperser/cmd/encoder/flags"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

const (
	// DefaultFragmentSizeBytes represents the size of each fragment in bytes (4MB)
	DefaultFragmentSizeBytes = 4 * 1024 * 1024
)

var (
	// Version is the version of the binary.
	Version   string
	GitCommit string
	GitDate   string
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "encoder"
	app.Usage = "EigenDA Encoder"
	app.Description = "Service for encoding blobs"

	app.Action = RunEncoderServer
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunEncoderServer(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	reg := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(reg, config.MetricsConfig.HTTPPort, logger)
	grpcMetrics := grpcprom.NewServerMetrics()
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)

		reg.MustRegister(grpcMetrics)
	}

	backendType, err := encoding.ParseBackendType(config.ServerConfig.Backend)
	if err != nil {
		return err
	}

	// Set the encoding config
	encodingConfig := &encoding.Config{
		BackendType: backendType,
		GPUEnable:   config.ServerConfig.GPUEnable,
		NumWorker:   config.EncoderConfig.NumWorker,
	}

	if config.EncoderVersion == V2 {
		// We no longer load the G2 points in V2 because the KZG commitments are computed
		// on the API server side.
		config.EncoderConfig.LoadG2Points = false
		prover, err := prover.NewProver(&config.EncoderConfig, encodingConfig)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		blobStoreBucketName := config.BlobStoreConfig.BucketName
		if blobStoreBucketName == "" {
			return fmt.Errorf("blob store bucket name is required")
		}

		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, s3Client, logger)
		logger.Info("Blob store", "bucket", blobStoreBucketName)

		chunkStoreBucketName := config.ChunkStoreConfig.BucketName
		chunkWriter := chunkstore.NewChunkWriter(logger, s3Client, chunkStoreBucketName, DefaultFragmentSizeBytes)
		logger.Info("Chunk store writer", "bucket", blobStoreBucketName)

		server := encoder.NewEncoderServerV2(
			*config.ServerConfig,
			blobStore,
			chunkWriter,
			logger,
			prover,
			metrics,
			grpcMetrics,
		)

		return server.Start()
	}

	config.EncoderConfig.LoadG2Points = true
	prover, err := prover.NewProver(&config.EncoderConfig, encodingConfig)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(*config.ServerConfig, logger, prover, metrics, grpcMetrics)

	return server.Start()
}
