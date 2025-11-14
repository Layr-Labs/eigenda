package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	commonpprof "github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/disperser/cmd/encoder/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg/prover"
	proverv2 "github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
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

	logger, err := common.NewLogger(&config.LoggerConfig)
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

	// Start pprof server if enabled (works for both v1 and v2)
	pprofProfiler := commonpprof.NewPprofProfiler(config.ServerConfig.PprofHttpPort, logger)
	if config.ServerConfig.EnablePprof {
		go pprofProfiler.Start()
		logger.Info("Enabled pprof for encoder server", "port", config.ServerConfig.PprofHttpPort)
	}

	backendType, err := encoding.ParseBackendType(config.ServerConfig.Backend)
	if err != nil {
		return err
	}

	// Set the encoding config
	encodingConfig := &encoding.Config{
		BackendType:                           backendType,
		GPUEnable:                             config.ServerConfig.GPUEnable,
		GPUConcurrentFrameGenerationDangerous: int64(config.ServerConfig.MaxConcurrentRequestsDangerous),
		NumWorker:                             config.EncoderConfig.NumWorker,
	}

	// Read the GRPC port from flags
	grpcPort := ctx.GlobalString(flags.GrpcPortFlag.Name)

	// Create listener
	addr := fmt.Sprintf("0.0.0.0:%s", grpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("Failed to close listener", "error", err)
		}
	}()

	if config.EncoderVersion == V2 {
		// We no longer load the G2 points in V2 because the KZG commitments are computed
		// on the API server side.
		config.EncoderConfig.LoadG2Points = false
		prover, err := proverv2.NewProver(logger, proverv2.KzgConfigFromV1Config(&config.EncoderConfig), encodingConfig)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		// Create object storage client (supports both S3 and OCI)
		objectStorageClient, err := blobstore.CreateObjectStorageClient(
			context.Background(), config.BlobStoreConfig, config.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		blobStoreBucketName := config.BlobStoreConfig.BucketName
		if blobStoreBucketName == "" {
			return fmt.Errorf("blob store bucket name is required")
		}

		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, objectStorageClient, logger)
		logger.Info("Blob store", "bucket", blobStoreBucketName, "backend", config.BlobStoreConfig.Backend)

		chunkStoreBucketName := config.ChunkStoreConfig.BucketName
		chunkWriter := chunkstore.NewChunkWriter(
			objectStorageClient,
			chunkStoreBucketName)
		logger.Info("Chunk store writer", "bucket", chunkStoreBucketName, "backend", config.ChunkStoreConfig.Backend)

		server := encoder.NewEncoderServerV2(
			*config.ServerConfig,
			blobStore,
			chunkWriter,
			logger,
			prover,
			metrics,
			grpcMetrics,
		)

		logger.Info("Starting encoder v2 server", "address", listener.Addr().String())

		//nolint:wrapcheck
		return server.StartWithListener(listener)
	}

	config.EncoderConfig.LoadG2Points = true
	prover, err := prover.NewProver(&config.EncoderConfig, encodingConfig)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(*config.ServerConfig, logger, prover, metrics, grpcMetrics)

	logger.Info("Starting encoder v1 server", "address", listener.Addr().String())

	//nolint:wrapcheck
	return server.StartWithListener(listener)
}
