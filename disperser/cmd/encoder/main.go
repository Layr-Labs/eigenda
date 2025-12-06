package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	commonpprof "github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg/prover"
	proverv2 "github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Version is the version of the binary.
	Version   string
	GitCommit string
	GitDate   string
)

func main() {
	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	// Block forever, the encoder runs as a server.
	select {}
}

// Run the encoder. This method is split from main() so we only have to use log.Fatalf() once.
func run(_ context.Context) error {
	config, err := config.Bootstrap(encoder.DefaultEncoderConfig, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}

	loggerConfig := common.DefaultLoggerConfig()
	loggerConfig.Format = config.LogFormat
	loggerConfig.HandlerOpts.NoColor = !config.LogColor
	level, err := common.StringToLogLevel(config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	loggerConfig.HandlerOpts.Level = level

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	reg := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(reg, config.MetricsPort, logger)
	grpcMetrics := grpcprom.NewServerMetrics()
	if config.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)

		reg.MustRegister(grpcMetrics)
	}

	// Start pprof server if enabled (works for both v1 and v2)
	pprofProfiler := commonpprof.NewPprofProfiler(config.Server.PprofHttpPort, logger)
	if config.Server.EnablePprof {
		go pprofProfiler.Start()
		logger.Info("Enabled pprof for encoder server", "port", config.Server.PprofHttpPort)
	}

	backendType, err := encoding.ParseBackendType(config.Server.Backend)
	if err != nil {
		return err
	}

	// Set the encoding config
	encodingConfig := &encoding.Config{
		BackendType:                           backendType,
		GPUEnable:                             config.Server.GPUEnable,
		GPUConcurrentFrameGenerationDangerous: int64(config.Server.MaxConcurrentRequestsDangerous),
		NumWorker:                             config.Kzg.NumWorker,
	}

	// Create listener
	addr := fmt.Sprintf("0.0.0.0:%s", config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("Failed to close listener", "error", err)
		}
	}()

	if config.Version == encoder.V2 {
		// We no longer load the G2 points in V2 because the KZG commitments are computed
		// on the API server side.
		config.Kzg.LoadG2Points = false
		prover, err := proverv2.NewProver(logger, proverv2.KzgConfigFromV1Config(&config.Kzg), encodingConfig)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		// Create object storage client (supports both S3 and OCI)
		objectStorageClient, err := blobstore.CreateObjectStorageClient(
			context.Background(), config.BlobStore, config.Aws, logger)
		if err != nil {
			return err
		}

		blobStoreBucketName := config.BlobStore.BucketName
		if blobStoreBucketName == "" {
			return fmt.Errorf("blob store bucket name is required")
		}

		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, objectStorageClient, logger)
		logger.Info("Blob store", "bucket", blobStoreBucketName, "backend", config.BlobStore.Backend)

		chunkStoreBucketName := config.ChunkStore.BucketName
		chunkWriter := chunkstore.NewChunkWriter(
			objectStorageClient,
			chunkStoreBucketName)
		logger.Info("Chunk store writer", "bucket", chunkStoreBucketName, "backend", config.ChunkStore.Backend)

		server := encoder.NewEncoderServerV2(
			config.Server,
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

	config.Kzg.LoadG2Points = true
	prover, err := prover.NewProver(&config.Kzg, encodingConfig)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(config.Server, logger, prover, metrics, grpcMetrics)

	logger.Info("Starting encoder v1 server", "address", listener.Addr().String())

	//nolint:wrapcheck
	return server.StartWithListener(listener)
}
