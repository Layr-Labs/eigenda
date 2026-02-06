package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/api"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// Version information set via ldflags during build
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := run(ctx)
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	// Bootstrap config using documented config framework
	cfg, err := config.Bootstrap(chainstate.DefaultRootIndexerConfig, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}

	// handle help flag correctly
	if cfg == nil {
		return nil
	}

	secretConfig := cfg.Secret
	indexerConfig := cfg.Config
	cfg = nil // Safety: discard root config to prevent accidental secret logging

	// Create logger
	logger, err := common.NewLogger(&indexerConfig.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("Starting ChainState Indexer",
		"version", Version,
		"gitCommit", GitCommit,
		"gitDate", GitDate)

	// Create Ethereum client
	if len(secretConfig.EthRpcUrls) == 0 {
		return fmt.Errorf("no Ethereum RPC URLs configured")
	}

	ethClient, err := ethclient.Dial(secretConfig.EthRpcUrls[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	defer ethClient.Close()

	// TODO(iquidus): dont leak api key in logs
	logger.Info("Connected to Ethereum RPC", "url", secretConfig.EthRpcUrls[0])

	// Create indexer
	indexer, err := chainstate.NewIndexer(ctx, indexerConfig, ethClient, logger)
	if err != nil {
		return fmt.Errorf("failed to create indexer: %w", err)
	}

	// Start indexer
	if err := indexer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start indexer: %w", err)
	}

	logger.Info("Indexer started successfully")

	// Create and start API server
	apiServer := api.NewServer(indexerConfig, indexer.GetStore(), logger)

	// Start API server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := apiServer.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	logger.Info("API server started", "port", indexerConfig.HTTPPort)

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("Received shutdown signal, stopping relay server")
	case err := <-errChan:
		logger.Error("Relay server failed", "error", err)
		return fmt.Errorf("relay server failed: %w", err)
	}

	// Give some time for graceful shutdown
	logger.Info("Shutdown complete")

	return nil
}
