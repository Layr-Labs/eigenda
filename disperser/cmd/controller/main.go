package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var (
	version   string
	gitCommit string
	gitDate   string

	controllerMaxStallDuration = 240 * time.Second
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "controller"
	app.Usage = "EigenDA Controller"
	app.Description = "EigenDA control plane for encoding and dispatching blobs"

	app.Action = RunController
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}
	select {}
}

func RunController(cliCtx *cli.Context) error {
	config, err := NewConfig(cliCtx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Reset readiness probe upon start-up
	if err := os.Remove(config.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to clean up readiness file", "error", err, "path", config.ControllerReadinessProbePath)
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create DynamoDB client: %w", err)
	}
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	contractDirectory, err := directory.NewContractDirectory(
		context.Background(),
		logger,
		gethClient,
		gethcommon.HexToAddress(config.EigenDAContractDirectoryAddress))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	operatorStateRetrieverAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.OperatorStateRetriever)
	if err != nil {
		return fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
	}
	serviceManagerAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.ServiceManager)
	if err != nil {
		return fmt.Errorf("failed to get ServiceManager address: %w", err)
	}
	registryCoordinatorAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.RegistryCoordinator)
	if err != nil {
		return fmt.Errorf("failed to get registry coordinator address: %w", err)
	}

	metricsRegistry := prometheus.NewRegistry()
	metricsRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metricsRegistry.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", config.MetricsPort)
	addr := fmt.Sprintf(":%d", config.MetricsPort)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		metricsRegistry,
		promhttp.HandlerOpts{},
	))
	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if !config.UseGraph {
		return fmt.Errorf("built-in indexer is deprecated and will be removed soon, please use UseGraph=true")
	}
	logger.Info("Using graph node")
	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)

	var requestSigner clients.DispersalRequestSigner
	if config.DisperserStoreChunksSigningDisabled {
		logger.Warn("StoreChunks() signing is disabled")
	} else {
		requestSigner, err = clients.NewDispersalRequestSigner(
			context.Background(),
			clients.DispersalRequestSignerConfig{
				Region:   config.AwsClientConfig.Region,
				Endpoint: config.AwsClientConfig.EndpointURL,
				KeyID:    config.DisperserKMSKeyID,
			})
		if err != nil {
			return fmt.Errorf("failed to create request signer: %v", err)
		}
	}

	ctx := context.Background()
	_, err = controller.StartController(
		ctx,
		logger,
		gethClient,
		dynamoClient,
		config.DynamoDBTableName,
		metricsRegistry,
		requestSigner,
		operatorStateRetrieverAddress,
		serviceManagerAddress,
		registryCoordinatorAddress,
		config.ChainStateConfig.Endpoint,
		&config.EncodingManagerConfig,
		&config.DispatcherConfig,

		// Chain state config
		config.ChainStateConfig,

		// Optional components
		metricsServer,
		config.ControllerReadinessProbePath,
		healthcheck.HeartbeatMonitorConfig{
			FilePath:         config.ControllerHealthProbePath,
			MaxStallDuration: controllerMaxStallDuration,
		},

		// Server config
		&config.ServerConfig,
		&config.OnDemandConfig,
		&config.ReservationConfig,
		contractDirectory,
	)
	if err != nil {
		return fmt.Errorf("failed to start controller: %w", err)
	}

	return nil
}
