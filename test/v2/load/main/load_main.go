package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/Layr-Labs/eigenda/test/v2/load"
)

func main() {

	cfg, err := config.Bootstrap(
		load.DefaultTrafficGeneratorConfig,
		"TRAFFIC_GENERATOR_SIGNER_PRIVATE_KEY_HEX",
		"TRAFFIC_GENERATOR_RPC_URLS",
	)
	if err != nil {
		panic(fmt.Errorf("failed to bootstrap config: %w", err))
	}

	loggerConfig := common.DefaultTextLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}

	var metrics *client.TestClientMetrics
	if !cfg.Environment.DisableMetrics {
		metrics = client.NewTestClientMetrics(logger, cfg.Environment.MetricsPort)
		metrics.Start()
	}

	testClient, err := client.NewTestClient(context.Background(), logger, metrics, &cfg.Environment)
	if err != nil {
		panic(fmt.Errorf("failed to create test client: %w", err))
	}

	generator, err := load.NewLoadGenerator(&cfg.Load, testClient)
	if err != nil {
		panic(fmt.Errorf("failed to create load generator: %w", err))
	}

	signals := make(chan os.Signal)
	go func() {
		<-signals
		generator.Stop()
	}()

	generator.Start(true)
}
