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
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/urfave/cli"
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

	metrics := encoder.NewMetrics(config.MetricsConfig.HTTPPort, logger)
	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)
	}

	if config.EncoderVersion == V2 {
		// Do not load G2 points for v2
		prover, err := prover.NewProver(&config.EncoderConfig, false)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		// Create a new s3 client
		s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		// Create a new blob store
		blobStoreBucketName := config.BlobStoreConfig.BucketName
		logger.Info("Blob store", "bucket", blobStoreBucketName)
		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, s3Client, logger)

		// Create a new chunk store writer
		chunkStoreBucketName := config.ChunkStoreConfig.BucketName
		logger.Info("Chunk store writer", "bucket", blobStoreBucketName)
		chunkWriter := chunkstore.NewChunkWriter(logger, s3Client, chunkStoreBucketName, 1) // TODO: make fragment size configurable?

		server := encoder.NewEncoderServerV2(
			*config.ServerConfig,
			blobStore,
			chunkWriter,
			logger,
			prover,
			metrics,
		)

		return server.Start()
	}

	prover, err := prover.NewProver(&config.EncoderConfig, true)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(*config.ServerConfig, logger, prover, metrics)

	return server.Start()

}
