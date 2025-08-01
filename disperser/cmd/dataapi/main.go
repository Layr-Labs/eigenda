package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/geth"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	dataapiprometheus "github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	serverv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

// @title			EigenDA Data Access API V1
// @description	This is the EigenDA Data Access API server.
// @version		1
// @Schemes		https http
func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "data-access-api"
	app.Usage = "EigenDA Data Access API"
	app.Description = "Service that provides access to data blobs."

	app.Action = RunDataApi
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunDataApi(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return err
	}

	s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	promApi, err := dataapiprometheus.NewApi(config.PrometheusConfig)
	if err != nil {
		return err
	}

	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return err
	}

	tx, err := coreeth.NewReader(logger, client, config.EigenDADirectory, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}

	var (
		reg               = prometheus.NewRegistry()
		promClient        = dataapi.NewPrometheusClient(promApi, config.PrometheusConfig.Cluster)
		subgraphApi       = subgraph.NewApi(config.SubgraphApiBatchMetadataAddr, config.SubgraphApiOperatorStateAddr, config.SubgraphApiPaymentsAddr)
		subgraphClient    = dataapi.NewSubgraphClient(subgraphApi, logger)
		chainState        = coreeth.NewChainState(tx, client)
		indexedChainState = thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)
	)

	if config.ServerVersion == 2 {
		baseBlobMetadataStorev2 := blobstorev2.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName)
		blobMetadataStorev2 := blobstorev2.NewInstrumentedMetadataStore(baseBlobMetadataStorev2, blobstorev2.InstrumentedMetadataStoreConfig{
			ServiceName: "dataapi",
			Registry:    reg,
			Backend:     blobstorev2.BackendDynamoDB,
		})

		// Register reservation collector
		reservationCollector := serverv2.NewReservationExpirationCollector(subgraphClient, logger)
		reg.MustRegister(reservationCollector)

		metrics := dataapi.NewMetrics(config.ServerVersion, reg, blobMetadataStorev2, config.MetricsConfig.HTTPPort, logger)
		serverv2, err := serverv2.NewServerV2(
			dataapi.Config{
				ServerMode:         config.ServerMode,
				SocketAddr:         config.SocketAddr,
				AllowOrigins:       config.AllowOrigins,
				DisperserHostname:  config.DisperserHostname,
				ChurnerHostname:    config.ChurnerHostname,
				BatcherHealthEndpt: config.BatcherHealthEndpt,
			},
			blobMetadataStorev2,
			promClient,
			subgraphClient,
			tx,
			chainState,
			indexedChainState,
			logger,
			metrics,
		)
		if err != nil {
			return fmt.Errorf("failed to create v2 server: %w", err)
		}

		// Enable Metrics Block
		if config.MetricsConfig.EnableMetrics {
			httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
			metrics.Start(context.Background())
			logger.Info("Enabled metrics for Data Access API", "socket", httpSocket)
		}

		return runServer(serverv2, logger)
	}

	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, 0)
	sharedStorage := blobstore.NewSharedStorage(config.BlobstoreConfig.BucketName, s3Client, blobMetadataStore, logger)
	metrics := dataapi.NewMetrics(config.ServerVersion, reg, blobMetadataStore, config.MetricsConfig.HTTPPort, logger)

	server, err := dataapi.NewServer(
		dataapi.Config{
			ServerMode:         config.ServerMode,
			SocketAddr:         config.SocketAddr,
			AllowOrigins:       config.AllowOrigins,
			DisperserHostname:  config.DisperserHostname,
			ChurnerHostname:    config.ChurnerHostname,
			BatcherHealthEndpt: config.BatcherHealthEndpt,
		},
		sharedStorage,
		promClient,
		subgraphClient,
		tx,
		chainState,
		indexedChainState,
		logger,
		metrics,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create v1 server: %w", err)
	}

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Data Access API", "socket", httpSocket)
	}

	return runServer(server, logger)
}

func runServer[T dataapi.ServerInterface](server T, logger logging.Logger) error {
	// Setup channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	// catch SIGINT (Ctrl+C) and SIGTERM (e.g., from `kill`)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a separate goroutine so that it doesn't block.
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block until a signal is received.
	<-quit
	logger.Info("Shutting down server...")
	err := server.Shutdown()

	if err != nil {
		logger.Errorf("Failed to shutdown server: %v", err)
	}

	return err
}
