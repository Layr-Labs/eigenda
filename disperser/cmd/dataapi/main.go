package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

// @title			EigenDA Data Access API
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
	config := NewConfig(ctx)

	logger, err := logging.GetLogger(config.LoggerConfig)
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

	promApi, err := prometheus.NewApi(config.PrometheusConfig)
	if err != nil {
		return err
	}

	client, err := geth.NewClient(config.EthClientConfig, logger)
	if err != nil {
		return err
	}

	tx, err := coreeth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}

	var (
		promClient        = dataapi.NewPrometheusClient(promApi, config.PrometheusConfig.Cluster)
		blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, 0)
		sharedStorage     = blobstore.NewSharedStorage(config.BlobstoreConfig.BucketName, s3Client, blobMetadataStore, logger)
		subgraphApi       = subgraph.NewApi(config.SubgraphApiBatchMetadataAddr, config.SubgraphApiOperatorStateAddr)
		subgraphClient    = dataapi.NewSubgraphClient(subgraphApi, logger)
		chainState        = coreeth.NewChainState(tx, client)
		metrics           = dataapi.NewMetrics(config.MetricsConfig.HTTPPort, logger)
		server            = dataapi.NewServer(
			dataapi.Config{
				ServerMode:   config.ServerMode,
				SocketAddr:   config.SocketAddr,
				AllowOrigins: config.AllowOrigins,
			},
			sharedStorage,
			promClient,
			subgraphClient,
			tx,
			chainState,
			logger,
			metrics,
		)
	)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Data Access API", "socket", httpSocket)
	}

	return server.Start()
}
