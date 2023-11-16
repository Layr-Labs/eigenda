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

	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/blobstore"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/ory/dockertest/v3"
)

var (
	logger         = &cmock.Logger{}
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
	s3Client   = cmock.NewS3Client()
	bucketName = "test-eigenda-blobstore"
	blobHash   = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	blobSize   = uint(len(blob.Data))

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	localStackPort     string

	dynamoClient      *dynamodb.Client
	blobMetadataStore *blobstore.BlobMetadataStore
	metadataTableName = "test-BlobMetadata"
	sharedStorage     *blobstore.SharedBlobStore
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {
	localStackPort = "4569"
	pool, resource, err := deploy.StartDockertestWithLocalstackContainer(localStackPort)
	dockertestPool = pool
	dockertestResource = resource
	if err != nil {
		teardown()
		panic("failed to start localstack container")
	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}

	_, err = test_utils.CreateTableIfNotExists(context.Background(), cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName, time.Hour)
	sharedStorage = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)
}

func teardown() {
	deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
}
