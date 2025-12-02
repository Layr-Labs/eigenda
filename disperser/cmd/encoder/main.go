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
	cfg, err := config.Bootstrap(encoder.DefaultEncoderConfig)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}

	loggerConfig := common.DefaultLoggerConfig()
	loggerConfig.Format = common.LogFormat(cfg.LogFormat)
	loggerConfig.HandlerOpts.NoColor = !cfg.LogColor
	level, err := common.StringToLogLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	loggerConfig.HandlerOpts.Level = level

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	reg := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(reg, cfg.MetricsPort, logger)
	grpcMetrics := grpcprom.NewServerMetrics()
	if cfg.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", cfg.MetricsPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)

		reg.MustRegister(grpcMetrics)
	}

	// Start pprof server if enabled (works for both v1 and v2)
	pprofProfiler := commonpprof.NewPprofProfiler(cfg.Server.PprofHttpPort, logger)
	if cfg.Server.EnablePprof {
		go pprofProfiler.Start()
		logger.Info("Enabled pprof for encoder server", "port", cfg.Server.PprofHttpPort)
	}

	backendType, err := encoding.ParseBackendType(cfg.Server.Backend)
	if err != nil {
		return err
	}

	// Set the encoding config
	encodingConfig := &encoding.Config{
		BackendType:                           backendType,
		GPUEnable:                             cfg.Server.GPUEnable,
		GPUConcurrentFrameGenerationDangerous: int64(cfg.Server.MaxConcurrentRequestsDangerous),
		NumWorker:                             cfg.Kzg.NumWorker,
	}

	// Create listener
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("Failed to close listener", "error", err)
		}
	}()

	if cfg.EncoderVersion == 2 {
		// We no longer load the G2 points in V2 because the KZG commitments are computed
		// on the API server side.
		cfg.Kzg.LoadG2Points = false
		prover, err := proverv2.NewProver(logger, proverv2.KzgConfigFromV1Config(&cfg.Kzg), encodingConfig)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		// Create object storage client (supports both S3 and OCI)
		objectStorageClient, err := blobstore.CreateObjectStorageClient(
			context.Background(), cfg.BlobStore, cfg.Aws, logger)
		if err != nil {
			return err
		}

		blobStoreBucketName := cfg.BlobStore.BucketName
		if blobStoreBucketName == "" {
			return fmt.Errorf("blob store bucket name is required")
		}

		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, objectStorageClient, logger)
		logger.Info("Blob store", "bucket", blobStoreBucketName, "backend", cfg.BlobStore.Backend)

		chunkStoreBucketName := cfg.ChunkStore.BucketName
		chunkWriter := chunkstore.NewChunkWriter(
			objectStorageClient,
			chunkStoreBucketName)
		logger.Info("Chunk store writer", "bucket", chunkStoreBucketName, "backend", cfg.ChunkStore.Backend)

		server := encoder.NewEncoderServerV2(
			cfg.Server,
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

	cfg.Kzg.LoadG2Points = true
	prover, err := prover.NewProver(&cfg.Kzg, encodingConfig)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(cfg.Server, logger, prover, metrics, grpcMetrics)

	logger.Info("Starting encoder v1 server", "address", listener.Addr().String())

	//nolint:wrapcheck
	return server.StartWithListener(listener)
}
