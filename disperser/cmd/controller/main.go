package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gammazero/workerpool"
	"github.com/urfave/cli"
)

var (
	version   string
	gitCommit string
	gitDate   string
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

func RunController(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}
	chainReader, err := eth.NewReader(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}

	blobMetadataStore := blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		config.DynamoDBTableName,
	)

	encoderClient, err := encoder.NewEncoderClientV2(config.EncodingManagerConfig.EncoderAddress)
	if err != nil {
		return fmt.Errorf("failed to create encoder client: %v", err)
	}
	encodingPool := workerpool.New(config.NumConcurrentEncodingRequests)
	encodingManager, err := controller.NewEncodingManager(
		config.EncodingManagerConfig,
		blobMetadataStore,
		encodingPool,
		encoderClient,
		chainReader,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create encoding manager: %v", err)
	}

	sigAgg, err := core.NewStdSignatureAggregator(logger, chainReader)
	if err != nil {
		return fmt.Errorf("failed to create signature aggregator: %v", err)
	}
	dispatcherPool := workerpool.New(config.NumConcurrentDispersalRequests)
	chainState := eth.NewChainState(chainReader, gethClient)
	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")

		logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
		ics = thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)
	} else {
		logger.Info("Using built-in indexer")
		rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURLs[0])
		if err != nil {
			return err
		}
		idx, err := indexer.CreateNewIndexer(
			&config.IndexerConfig,
			gethClient,
			rpcClient,
			config.EigenDAServiceManagerAddr,
			logger,
		)
		if err != nil {
			return err
		}
		ics, err = indexer.NewIndexedChainState(chainState, idx)
		if err != nil {
			return err
		}
	}
	nodeClientManager, err := controller.NewNodeClientManager(config.NodeClientCacheSize, logger)
	if err != nil {
		return fmt.Errorf("failed to create node client manager: %v", err)
	}
	dispatcher, err := controller.NewDispatcher(
		config.DispatcherConfig,
		blobMetadataStore,
		dispatcherPool,
		ics,
		sigAgg,
		nodeClientManager,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create dispatcher: %v", err)
	}

	c := context.Background()
	err = encodingManager.Start(c)
	if err != nil {
		return fmt.Errorf("failed to start encoding manager: %v", err)
	}

	err = dispatcher.Start(c)
	if err != nil {
		return fmt.Errorf("failed to start dispatcher: %v", err)
	}

	return nil
}
