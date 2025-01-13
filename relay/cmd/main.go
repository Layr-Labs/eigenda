package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common/geth"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/cmd/flags"
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
	app.Name = "relay"
	app.Usage = "EigenDA Relay"
	app.Description = "EigenDA relay for serving blobs and chunks data"

	app.Action = RunRelay
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}
	select {}
}

// RunRelay is the entrypoint for the relay.
func RunRelay(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create relay config: %w", err)
	}

	logger, err := common.NewLogger(config.Log)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	dynamoClient, err := dynamodb.NewClient(config.AWS, logger)
	if err != nil {
		return fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	s3Client, err := s3.NewClient(context.Background(), config.AWS, logger)
	if err != nil {
		return fmt.Errorf("failed to create s3 client: %w", err)
	}

	metadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.MetadataTableName)
	blobStore := blobstore.NewBlobStore(config.BucketName, s3Client, logger)
	chunkReader := chunkstore.NewChunkReader(logger, s3Client, config.BucketName)
	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return fmt.Errorf("failed to create eth client: %w", err)
	}

	tx, err := coreeth.NewWriter(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return fmt.Errorf("failed to create eth writer: %w", err)
	}

	cs := coreeth.NewChainState(tx, client)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, cs, logger)

	server, err := relay.NewServer(
		context.Background(),
		logger,
		&config.RelayConfig,
		metadataStore,
		blobStore,
		chunkReader,
		tx,
		ics,
	)
	if err != nil {
		return fmt.Errorf("failed to create relay server: %w", err)
	}

	err = server.Start(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start relay server: %w", err)
	}

	return nil
}
