package ratelimit_test

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
	"github.com/Layr-Labs/eigenda/common/mock"
	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
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
	bucketTableName  = "BucketStore"
)

func makeTestRatelimiter() (common.RateLimiter, error) {

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{time.Second, time.Minute},
		Multipliers: []float32{1, 1},
	}
	bucketStoreSize := 1000

	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](bucketStoreSize)
	if err != nil {
		return nil, err
	}

	ratelimiter, err := ratelimit.NewRateLimiter(globalParams, bucketStore, []string{"testRetriever2"}, &mock.Logger{})
	if err != nil {
		return nil, err
	}

	return ratelimiter, nil

}

func makeTestRatelimiterDynamodB() (common.RateLimiter, error) {

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
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
		panic("failed to create dynamodb table: " + err.Error())
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		panic("failed to create dynamodb client: " + err.Error())
	}

	dynamoParamStore = store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, bucketTableName)

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{time.Second, time.Minute},
		Multipliers: []float32{1, 1},
	}

	bucketStore := store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, bucketTableName)
	if err != nil {
		return nil, err
	}

	ratelimiter, err := ratelimit.NewRateLimiter(globalParams, bucketStore, []string{}, &mock.Logger{})
	if err != nil {
		return nil, err
	}

	return ratelimiter, nil

}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestRatelimit(t *testing.T) {

	ratelimiter, err := makeTestRatelimiter()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever"

	for i := 0; i < 10; i++ {
		allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, true, allow)
	}

	allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)
}

func TestRatelimitAllowList(t *testing.T) {

	ratelimiter, err := makeTestRatelimiter()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever2"

	for i := 0; i < 10; i++ {
		allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, true, allow)
	}

	allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)
}

func TestRatelimitDynamodBStore(t *testing.T) {

	ratelimiter, err := makeTestRatelimiterDynamodB()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever"

	for i := 0; i < 10; i++ {
		fmt.Println("i", i)
		allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, true, allow)
	}

	allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 1000, 100)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	teardown()
}
