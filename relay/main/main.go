package main

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
)

// main is the entrypoint for the relay.
func main() {

	config, err := relay.LoadConfigWithViper()
	if err != nil {
		panic(fmt.Sprintf("fatal error loading config: %s", err))
	}

	logger, err := common.NewLogger(config.Log)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating logger: %s", err))
	}
	logger.Info(fmt.Sprintf("Relay configuration: %#v", config))

	dynamoClient, err := dynamodb.NewClient(config.AWS, logger)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating dynamodb client: %s", err))
	}

	s3Client, err := s3.NewClient(context.Background(), config.AWS, logger)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating s3 client: %s", err))
	}

	metadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.MetadataTableName)
	blobStore := blobstore.NewBlobStore(config.BucketName, s3Client, logger)
	chunkReader := chunkstore.NewChunkReader(logger, s3Client, config.BucketName)

	server, err := relay.NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		chunkReader)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating relay server: %s", err))
	}

	err = server.Start()
	if err != nil {
		panic(fmt.Sprintf("fatal error starting relay server: %s", err))
	}

	// Block forever.
	select {}
}
