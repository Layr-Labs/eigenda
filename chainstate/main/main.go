package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/service"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
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

	// Derive a cancellable context for the indexer so that, even if shutdown is
	// triggered by an API server error rather than a signal, we can stop the
	// indexer and let its persister flush a final state save before exiting.
	ctx, cancelIndexer := context.WithCancel(ctx)
	defer cancelIndexer()

	// Create logger
	logger, err := common.NewLogger(&indexerConfig.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("Starting ChainState Indexer",
		"version", Version,
		"gitCommit", GitCommit,
		"gitDate", GitDate)

	// Create the service (dials RPC, builds the indexer and API server).
	svc, err := service.New(ctx, indexerConfig, secretConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Start indexing and the API server.
	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	logger.Info("Indexer and API server started", "port", indexerConfig.HTTPPort)

	// Wait for either a shutdown signal (via ctx) or a fatal API server error.
	var runErr error
	select {
	case <-ctx.Done():
		logger.Info("Received shutdown signal, stopping indexer")
	case err := <-svc.Errors():
		logger.Error("API server failed", "error", err)
		runErr = fmt.Errorf("API server failed: %w", err)
	}

	// Cancel the indexer context (a no-op if a signal already cancelled it) and
	// block until its background goroutines stop, ensuring the persister's final
	// state save completes before we exit on every shutdown path.
	cancelIndexer()
	svc.Wait()
	logger.Info("Shutdown complete")

	return runErr
}
