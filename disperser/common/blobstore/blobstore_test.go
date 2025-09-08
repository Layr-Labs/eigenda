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
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/google/uuid"
)

var (
	logger         = testutils.GetLogger()
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

	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		cfg := testbed.DefaultLocalStackConfig()
		cfg.Services = []string{"s3", "dynamodb"}
		cfg.Port = localstackPort
		cfg.Host = "0.0.0.0"

		localstackContainer, err = testbed.NewLocalStackContainer(context.Background(), cfg)
		if err != nil {
			teardown()
			panic("failed to start localstack container: " + err.Error())
		}

	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
	}

	_, err := test_utils.CreateTable(context.Background(), cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	if shadowMetadataTableName != "" {
		_, err = test_utils.CreateTable(context.Background(), cfg, shadowMetadataTableName, blobstore.GenerateTableSchema(shadowMetadataTableName, 10, 10))
		if err != nil {
			teardown()
			panic("failed to create shadow dynamodb table: " + err.Error())
		}
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}

	blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName, time.Hour)
	sharedStorage = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)
}

func teardown() {
	if deployLocalStack && localstackContainer != nil {
		_ = localstackContainer.Terminate(context.Background())
	}
}
