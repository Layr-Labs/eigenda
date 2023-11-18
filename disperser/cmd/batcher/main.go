package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shurcooL/graphql"

	"github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/core/thegraph"

	inmemstore "github.com/Layr-Labs/eigenda/indexer/inmem"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/eth"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/Layr-Labs/eigenda/disperser/cmd/batcher/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "batcher"
	app.Usage = "EigenDA Batcher"
	app.Description = "Service for creating a batch from queued blobs, distributing coded chunks to nodes, and confirming onchain"

	app.Action = RunBatcher
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunBatcher(ctx *cli.Context) error {
	config := NewConfig(ctx)

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	bucketName := config.BlobstoreConfig.BucketName
	s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
	if err != nil {
		return err
	}
	logger.Info("Initialized S3 client", "bucket", bucketName)

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	dispatcher := dispatcher.NewDispatcher(&dispatcher.Config{
		Timeout: config.TimeoutConfig.AttestationTimeout,
	}, logger)
	agg := core.NewStdSignatureAggregator(logger)
	asgn := &core.StdAssignmentCoordinator{}

	client, err := geth.NewClient(config.EthClientConfig, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}
	rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURL)
	if err != nil {
		return err
	}
	tx, err := coreeth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}
	confirmer, err := eth.NewBatchConfirmer(tx, config.TimeoutConfig.ChainWriteTimeout)
	if err != nil {
		return err
	}

	blockStaleMeasure, err := tx.GetBlockStaleMeasure(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
	}
	storeDurationBlocks, err := tx.GetStoreDurationBlocks(context.Background())
	if err != nil || storeDurationBlocks == 0 {
		return fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
	}
	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, time.Duration((storeDurationBlocks+blockStaleMeasure)*12)*time.Second)
	queue := blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)

	cs := coreeth.NewChainState(tx, client)

	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")
		querier := graphql.NewClient(config.GraphUrl, nil)
		logger.Info("Connecting to subgraph", "url", config.GraphUrl)
		ics = thegraph.NewIndexedChainState(cs, querier, logger)
	} else {
		logger.Info("Using built-in indexer")

		store := inmemstore.NewHeaderStore()

		ics, err = indexer.NewIndexedChainState(&config.IndexerConfig, gethcommon.HexToAddress(config.EigenDAServiceManagerAddr), cs, store, client, rpcClient, logger)
		if err != nil {
			return err
		}
	}

	metrics := batcher.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	if len(config.BatcherConfig.EncoderSocket) == 0 {
		return fmt.Errorf("encoder socket must be specified")
	}
	encoderClient, err := encoder.NewEncoderClient(config.BatcherConfig.EncoderSocket, config.TimeoutConfig.EncodingTimeout)
	if err != nil {
		return err
	}
	finalizer := batcher.NewFinalizer(config.TimeoutConfig.ChainReadTimeout, config.BatcherConfig.FinalizerInterval, queue, client, rpcClient, logger)
	batcher, err := batcher.NewBatcher(config.BatcherConfig, config.TimeoutConfig, queue, dispatcher, confirmer, ics, asgn, encoderClient, agg, client, finalizer, logger, metrics)
	if err != nil {
		return err
	}

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Batcher", "socket", httpSocket)
	}

	err = batcher.Start(context.Background())
	if err != nil {
		return err
	}

	return nil

}
