package store_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/elasticcache"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {

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

	// Set up the Redis client to point to your local Redis server
	clientConfig := elasticcache.RedisClientConfig{
		EndpointURL: "localhost",
		Port:        "6379",
	}
	redisClient, err := elasticcache.NewClient(clientConfig, nil) // Assuming logger can be nil
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	redisStore := store.NewRedisStore[common.RateBucketParams](redisClient)

	// Run your tests here
	// Example: Test Set and Get
	ctx := context.Background()
	testKey := "testKey"
	testValue := common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now().UTC(),
	}

	err = redisStore.UpdateItem(ctx, testKey, &testValue)
	assert.NoError(t, err, "UpdateItem should not return an error")

	result, err := redisStore.GetItem(ctx, testKey)
	assert.NoError(t, err, "GetItem should not return an error")
	assert.Equal(t, testValue, *result, "GetItem should return the value that was set")
}
