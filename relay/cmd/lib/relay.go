package lib

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/relay"
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

	// Create relay server dependencies
	deps, err := relay.NewServerDependencies(
		ctx,
		relay.ServerDependenciesConfig{
			AWSConfig:                  config.AWS,
			MetadataTableName:          config.MetadataTableName,
			BucketName:                 config.BucketName,
			OperatorStateRetrieverAddr: config.OperatorStateRetrieverAddr,
			ServiceManagerAddr:         config.EigenDAServiceManagerAddr,
			ChainStateConfig:           config.ChainStateConfig,
			EthClientConfig:            config.EthClientConfig,
			LoggerConfig:               config.Log,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create relay server dependencies: %w", err)
	}

	// Create listener on configured port
	addr := fmt.Sprintf("0.0.0.0:%d", config.RelayConfig.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	// Create and start the relay instance
	instance, err := relay.NewInstanceWithDependencies(ctx, &config.RelayConfig, deps, listener)
	if err != nil {
		return fmt.Errorf("failed to create relay instance: %w", err)
	}

	deps.Logger.Info("Relay server started successfully", "port", instance.Port)

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	deps.Logger.Info("Received shutdown signal, stopping relay server", "signal", sig)

	// Gracefully stop the server
	if err := instance.Stop(); err != nil {
		deps.Logger.Warn("Error stopping relay server", "error", err)
		return fmt.Errorf("error stopping relay server: %w", err)
	}

	return nil
}
