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
	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	logger = &cmock.Logger{}

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource

	deployLocalStack bool
	localStackPort   = "4566"

	dynamoClient     *dynamodb.Client
	dynamoParamStore common.KVStoreVersioned[common.RateBucketParams]
	bucketTableName  = "BucketStoreVersioned"
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

func TestDynamoBucketStoreVersioned(t *testing.T) {
	ctx := context.Background()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now().UTC(),
	}

	p2, version, err := dynamoParamStore.GetItemWithVersion(ctx, "testRetriever")
	assert.Error(t, err)
	assert.Nil(t, p2)
	assert.Equal(t, 0, version)

	err = dynamoParamStore.UpdateItemWithVersion(ctx, "testRetriever", p, version)
	assert.NoError(t, err)

	p2, version, err = dynamoParamStore.GetItemWithVersion(ctx, "testRetriever")

	assert.NoError(t, err)
	assert.Equal(t, p, p2)
	assert.Equal(t, 1, version)
}

func TestUpsertMultipleUpdateAsSeparateOperationWithExpression(t *testing.T) {

	ctx := context.Background()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{30 * time.Second, 30 * time.Second},
		LastRequestTime: time.Now().UTC(),
	}
	err := dynamoParamStore.UpdateItemWithVersion(ctx, "testRetriever2", p, 0)
	assert.NoError(t, err)
	assert.NoError(t, err)

	// Retrieve and check the initial item
	p2, _, err := dynamoParamStore.GetItemWithVersion(ctx, "testRetriever2")

	assert.NoError(t, err)
	for i := 0; i < len(p2.BucketLevels); i++ {
		delta := uint64(100 * time.Second)
		// Create a new UpdateBuilder for each attribute update
		// Ideally This should be an ADD Operation but DynamoDB only supports numeric types (like integers or floating-point numbers) directly
		// Chain VersionName in the update
		updateBuilder := expression.Add(
			expression.Name(fmt.Sprintf("BucketLevels[%d]", i)),
			expression.Value(delta),
		).Add(
			expression.Name("Version"),
			expression.Value(1),
		)

		err := dynamoParamStore.UpdateItemWithExpression(ctx, "testRetriever2", &updateBuilder)
		assert.NoError(t, err)
	}

	p3, version, err := dynamoParamStore.GetItemWithVersion(ctx, "testRetriever2")

	// Validate that the item was updated
	for i := 0; i < len(p2.BucketLevels); i++ {
		p.BucketLevels[i] += 100 * time.Second
		fmt.Printf("p3.BucketLevels[%d]: %v\n", i, p3.BucketLevels[i])
		assert.Equal(t, p.BucketLevels[i], p3.BucketLevels[i])
	}

	assert.NoError(t, err)
	assert.Equal(t, 3, version)
}
