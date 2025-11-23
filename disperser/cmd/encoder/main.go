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
func run(ctx context.Context) error {
	rootCfg, err := config.Bootstrap(DefaultRootEncoderConfig)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}
	encoderConfig := rootCfg.Config
	// Ensure we don't accidentally use rootCfg after this point.
	rootCfg = nil

	loggerConfig := common.DefaultLoggerConfig()
	loggerConfig.Format = common.LogFormat(encoderConfig.LogOutputType)
	loggerConfig.HandlerOpts.NoColor = !encoderConfig.LogColor

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	reg := prometheus.NewRegistry()
	metrics := encoder.NewMetrics(reg, encoderConfig.MetricsConfig.HTTPPort, logger)
	grpcMetrics := grpcprom.NewServerMetrics()
	if encoderConfig.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", encoderConfig.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)

		reg.MustRegister(grpcMetrics)
	}

	// Start pprof server if enabled (works for both v1 and v2)
	pprofProfiler := commonpprof.NewPprofProfiler(encoderConfig.ServerConfig.PprofHttpPort, logger)
	if encoderConfig.ServerConfig.EnablePprof {
		go pprofProfiler.Start()
		logger.Info("Enabled pprof for encoder server", "port", encoderConfig.ServerConfig.PprofHttpPort)
	}

	backendType, err := encoding.ParseBackendType(encoderConfig.ServerConfig.Backend)
	if err != nil {
		return err
	}

	// Set the encoding config
	encodingConfig := &encoding.Config{
		BackendType:                           backendType,
		GPUEnable:                             encoderConfig.ServerConfig.GPUEnable,
		GPUConcurrentFrameGenerationDangerous: int64(encoderConfig.ServerConfig.MaxConcurrentRequestsDangerous),
		NumWorker:                             encoderConfig.KzgConfig.NumWorker,
	}

	// Create listener
	addr := fmt.Sprintf("0.0.0.0:%s", encoderConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("Failed to close listener", "error", err)
		}
	}()

	if encoderConfig.EncoderVersion == 2 {
		// We no longer load the G2 points in V2 because the KZG commitments are computed
		// on the API server side.
		encoderConfig.KzgConfig.LoadG2Points = false
		prover, err := proverv2.NewProver(logger, proverv2.KzgConfigFromV1Config(&encoderConfig.KzgConfig), encodingConfig)
		if err != nil {
			return fmt.Errorf("failed to create encoder: %w", err)
		}

		// Create object storage client (supports both S3 and OCI)
		objectStorageClient, err := blobstore.CreateObjectStorageClient(
			context.Background(), encoderConfig.BlobStoreConfig, encoderConfig.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		blobStoreBucketName := encoderConfig.BlobStoreConfig.BucketName
		if blobStoreBucketName == "" {
			return fmt.Errorf("blob store bucket name is required")
		}

		blobStore := blobstorev2.NewBlobStore(blobStoreBucketName, objectStorageClient, logger)
		logger.Info("Blob store", "bucket", blobStoreBucketName, "backend", encoderConfig.BlobStoreConfig.Backend)

		chunkStoreBucketName := encoderConfig.ChunkStoreConfig.BucketName
		chunkWriter := chunkstore.NewChunkWriter(
			objectStorageClient,
			chunkStoreBucketName)
		logger.Info("Chunk store writer", "bucket", chunkStoreBucketName, "backend", encoderConfig.ChunkStoreConfig.Backend)

		server := encoder.NewEncoderServerV2(
			encoderConfig.ServerConfig,
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

	encoderConfig.KzgConfig.LoadG2Points = true
	prover, err := prover.NewProver(&encoderConfig.KzgConfig, encodingConfig)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	server := encoder.NewEncoderServer(encoderConfig.ServerConfig, logger, prover, metrics, grpcMetrics)

	logger.Info("Starting encoder v1 server", "address", listener.Addr().String())

	//nolint:wrapcheck
	return server.StartWithListener(listener)
}
