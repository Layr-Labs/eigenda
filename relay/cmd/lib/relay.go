package lib

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	blobstorefactory "github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

// RunRelay is the entrypoint for the relay.
func RunRelay(cliCtx *cli.Context) error {
	config, err := NewConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to create relay config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create logger
	logger, err := common.NewLogger(&config.Log)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Create eth client
	ethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return fmt.Errorf("failed to create eth client: %w", err)
	}

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(config.AWS, logger)
	if err != nil {
		return fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create object storage client (supports both S3 and OCI)
	blobStoreConfig := blobstorefactory.Config{
		BucketName: config.BucketName,
		Backend:    blobstorefactory.ObjectStorageBackend(config.ObjectStorageBackend),
	}
	objectStorageClient, err := blobstorefactory.CreateObjectStorageClient(
		ctx, blobStoreConfig, config.AWS, logger)
	if err != nil {
		return fmt.Errorf("failed to create object storage client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.MetadataTableName)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "relay",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create blob store and chunk reader
	blobStore := blobstore.NewBlobStore(config.BucketName, objectStorageClient, logger)
	chunkReader := chunkstore.NewChunkReader(logger, objectStorageClient, config.BucketName)

	// Create eth writer
	tx, err := eth.NewWriter(logger, ethClient, config.OperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return fmt.Errorf("failed to create eth writer: %w", err)
	}

	// Create chain state
	cs := eth.NewChainState(tx, ethClient)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, cs, logger)

	// Create listener
	addr := fmt.Sprintf("0.0.0.0:%d", config.RelayConfig.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	// Create server
	server, err := relay.NewServer(
		ctx,
		metricsRegistry,
		logger,
		&config.RelayConfig,
		metadataStore,
		blobStore,
		chunkReader,
		tx,
		ics,
		listener,
	)
	if err != nil {
		_ = listener.Close()
		return fmt.Errorf("failed to create relay server: %w", err)
	}

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Starting relay server", "address", listener.Addr().String())
		if err := server.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal, stopping relay server", "signal", sig)
	case err := <-errChan:
		logger.Error("Relay server failed", "error", err)
		return fmt.Errorf("relay server failed: %w", err)
	}

	// Gracefully stop the server
	if err := server.Stop(); err != nil {
		logger.Warn("Error stopping relay server", "error", err)
		return fmt.Errorf("error stopping relay server: %w", err)
	}

	return nil
}
