package store_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	logger = logging.NewNoopLogger()

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource

	deployLocalStack bool
	localStackPort   = "4566"

	dynamoClient     dynamodb.Client
	dynamoParamStore common.KVStore[common.RateBucketParams]
	bucketTableName  = "BucketStore"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	_, err := test_utils.CreateTable(context.Background(), cfg, bucketTableName, store.GenerateTableSchema(10, 10, bucketTableName))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}

	dynamoParamStore = store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, bucketTableName)
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestDynamoBucketStore(t *testing.T) {
	ctx := context.Background()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now().UTC(),
	}

	p2, err := dynamoParamStore.GetItem(ctx, "testRetriever")
	assert.Error(t, err)
	assert.Nil(t, p2)

	err = dynamoParamStore.UpdateItem(ctx, "testRetriever", p)
	assert.NoError(t, err)

	p2, err = dynamoParamStore.GetItem(ctx, "testRetriever")

	assert.NoError(t, err)
	assert.Equal(t, p, p2)
}
