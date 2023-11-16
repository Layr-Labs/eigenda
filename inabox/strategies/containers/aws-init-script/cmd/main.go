package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Layr-Labs/eigenda/common/aws"
	dynamodbutils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	s3utils "github.com/Layr-Labs/eigenda/common/aws/s3/utils"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/disperser/blobstore"
	"github.com/urfave/cli/v2"
)

var (
	localstackFlagName = "localstack-port"

	metadataTableName = "test-BlobMetadata"
	bucketTableName   = "test-BucketStore"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  localstackFlagName,
				Value: "4570",
				Usage: "path to the config file",
			},
		},
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	err := DeployResources(ctx.String(localstackFlagName), metadataTableName, bucketTableName)
	if err != nil {
		return err
	}

	// Indicates that the script is done and other docker compose services can start
	startHealthcheckServer()

	return nil
}

func startHealthcheckServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Println("TCP server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		conn.Close() // Close the connection immediately for the health check purpose
	}
}

func DeployResources(localStackPort, metadataTableName, bucketTableName string) error {
	endpoint := fmt.Sprintf("http://localstack:%s", localStackPort)

	s3utils.CheckOrCreateBucket(
		"localstack",
		"localstack",
		"test-eigenda-blobstore",
		"us-east-1",
		endpoint,
	)

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     endpoint,
	}

	_, err := dynamodbutils.CreateTableIfNotExists(context.Background(), cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		return err
	}

	_, err = dynamodbutils.CreateTableIfNotExists(context.Background(), cfg, bucketTableName, store.GenerateTableSchema(10, 10, bucketTableName))
	return err
}
