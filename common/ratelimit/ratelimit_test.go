package ratelimit_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
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

	dynamoClient                    *dynamodb.Client
	dynamoParamStore                common.KVStoreVersioned[common.RateBucketParams]
	dynamoParamStoreConccurencySafe common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe]
	bucketTableName                 = "BucketStoreRateLimit"
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

	ratelimiter, err := ratelimit.NewRateLimiter(globalParams, bucketStore, &mock.Logger{})
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

	ratelimiter, err := ratelimit.NewRateLimiter(globalParams, bucketStore, &mock.Logger{})
	if err != nil {
		return nil, err
	}

	return ratelimiter, nil

}

func makeTestRatelimiterDynamodBConccurencySafe() (common.RateLimiter, error) {

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

	dynamoParamStoreConccurencySafe = store.NewDynamoParamStore[common.RateBucketParamsConcurrencySafe](dynamoClient, bucketTableName)

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{time.Second, time.Second},
		Multipliers: []float32{1, 1},
	}
	fmt.Printf("DynamoStore %v\n", dynamoParamStoreConccurencySafe)

	bucketStore := store.NewDynamoParamStore[common.RateBucketParamsConcurrencySafe](dynamoClient, bucketTableName)

	fmt.Printf("bucketStore %v\n", bucketStore)

	ratelimiter, err := ratelimit.NewRateLimiter(globalParams, bucketStore, &mock.Logger{})
	if err != nil {
		fmt.Printf("err %v\n", err)
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

func TestRatelimitDynamodBStore(t *testing.T) {

	ratelimiter, err := makeTestRatelimiterDynamodB()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever"

	for i := 0; i < 10; i++ {
		allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, true, allow)
	}

	allow, err := ratelimiter.AllowRequest(ctx, retreiverID, 1000, 100)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	teardown()
}

func TestRatelimitDynamodBStoreConccurencySafe(t *testing.T) {

	ratelimiter, err := makeTestRatelimiterDynamodBConccurencySafe()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever2"

	// Allow small blob at rate of 1
	allow, err := ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small blob at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small blob at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 100, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small of volume 1000 blob at rate of 10 because calculate rate is 100
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 1000, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Larger Blob of volume 10000 fails at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10000, 10)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	// // Disperse a smaller Blob should succeed
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 100000, 10000)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 1000000000, 10000)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	teardown()
}

func TestRatelimitDynamodBStoreConccurencySafeMultipleSerialRequests(t *testing.T) {

	ratelimiter, err := makeTestRatelimiterDynamodBConccurencySafe()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever2"

	// Allow small blob at rate of 1
	allow, err := ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small blob at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small blob at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 100, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Allow small of volume 1000 blob at rate of 10 because calculate rate is 100
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 1000, 10)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	// Larger Blob of volume 10000 fails at rate of 10
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 10000, 10)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	// // Disperse a smaller Blob should succeed
	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 100000, 10000)
	assert.NoError(t, err)
	assert.Equal(t, true, allow)

	allow, err = ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retreiverID, 1000000000, 10000)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)

	teardown()
}

func TestRatelimitDynamodBStoreConcurrencySafeMultipleParallelRequests(t *testing.T) {
	ratelimiter, err := makeTestRatelimiterDynamodBConccurencySafe()
	assert.NoError(t, err)

	ctx := context.Background()
	retrieverID := "testRetriever2"

	var wg sync.WaitGroup
	numGoroutines := 3 // Number of concurrent requests
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Simulate concurrent requests with varying blob sizes and rates
			blobSize := uint(1000 + rand.Intn(1000)) // Random blob size between 1000 and 2000
			rate := uint32(10 + rand.Intn(5))        // Random rate between 10 and 15

			allow, err := ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retrieverID, blobSize, rate)
			assert.NoError(t, err)
			assert.True(t, allow)

			fmt.Printf("Goroutine %d completed with allow: %v\n", goroutineID, allow)
		}(i)

		// Delay between starting each goroutine
		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait() // Wait for all goroutines to finish

	teardown()
}

func TestRatelimitDynamodBStoreConcurrencySafeMultipleRequestsFailedWithMoreTimeBetweenRequests(t *testing.T) {
	ratelimiter, err := makeTestRatelimiterDynamodBConccurencySafe()
	assert.NoError(t, err)

	ctx := context.Background()
	retrieverID := "testRetriever2"
	blobSizes := []uint{1000, 10000, 100000}

	var wg sync.WaitGroup
	numGoroutines := 3 // Number of concurrent requests
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Simulate concurrent requests with varying blob sizes and rates
			blobSize := uint(1000 + blobSizes[i]) // Random blob size between 1000 and 2000
			rate := uint32(10 + rand.Intn(5))     // Random rate between 10 and 15

			allow, err := ratelimiter.AllowRequestConcurrencySafeVersion(ctx, retrieverID, blobSize, rate)
			assert.NoError(t, err)

			if i == 2 {
				assert.False(t, allow)
			} else {
				assert.True(t, allow)
			}
			fmt.Printf("Goroutine %d completed with allow: %v\n", goroutineID, allow)
		}(i)

		// Delay between starting each goroutine
		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait() // Wait for all goroutines to finish

	teardown()
}
