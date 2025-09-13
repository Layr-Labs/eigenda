package blobstore_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	awsmock "github.com/Layr-Labs/eigenda/common/aws/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/google/uuid"
)

var (
	logger         = test.GetLogger()
	securityParams = []*core.SecurityParam{{
		QuorumID:           1,
		AdversaryThreshold: 80,
		QuorumRate:         32000,
	},
	}
	blob = &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: []byte("test"),
	}
	s3Client   = awsmock.NewS3Client()
	bucketName = "test-eigenda-blobstore"
	blobHash   = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	blobSize   = uint(len(blob.Data))

	localstackContainer *testbed.LocalStackContainer

	deployLocalStack bool
	localstackPort   = "4569"

	dynamoClient      dynamodb.Client
	blobMetadataStore *blobstore.BlobMetadataStore
	sharedStorage     *blobstore.SharedBlobStore

	UUID                    = uuid.New()
	metadataTableName       = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	shadowMetadataTableName = fmt.Sprintf("test-BlobMetadata-Shadow-%v", UUID)
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(_ *testing.M) {
	ctx := context.Background()

	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localstackPort,
			Services:       []string{"s3", "dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			teardown()
			logger.Fatal("Failed to start localstack container:", err)
		}

	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
	}

	_, err := test_utils.CreateTable(ctx, cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		logger.Fatal("Failed to create dynamodb table:", err)
	}

	if shadowMetadataTableName != "" {
		_, err = test_utils.CreateTable(ctx, cfg, shadowMetadataTableName,
			blobstore.GenerateTableSchema(shadowMetadataTableName, 10, 10))
		if err != nil {
			teardown()
			logger.Fatal("Failed to create shadow dynamodb table:", err)
		}
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create dynamodb client:", err)
	}

	blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName, time.Hour)
	sharedStorage = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)
}

func teardown() {
	if deployLocalStack && localstackContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = localstackContainer.Terminate(ctx)
	}
}
