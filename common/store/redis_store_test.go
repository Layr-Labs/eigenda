package store_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/elasticcache"
	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {
	// Start Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %v", err)
	}

	// Start Redis container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "latest",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "", HostPort: "6379"}},
		},
	})
	if err != nil {
		t.Fatalf("Could not start Redis container: %v", err)
	}

	// Delay cleanup until after all tests have run
	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge Redis container: %v", err)
		}
	})

	// Wait for Redis to be ready
	if err := pool.Retry(func() error {
		// Perform a health check...
		return nil // return nil if healthy
	}); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Set up the Redis client to point to your local Redis server
	clientConfig := elasticcache.RedisClientConfig{
		EndpointURL: "localhost",
		Port:        "6379",
	}

	redisClient, err := elasticcache.NewClient(clientConfig, &cmock.Logger{}) // Assuming logger can be nil
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	redisStore := store.NewRedisStore[common.RateBucketParams](redisClient, "testKey")

	// Test Update and Get Item
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

func TestRedisStoreAcquireAndReleaseLock(t *testing.T) {

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

	redisStore := store.NewRedisStore[common.RateBucketParams](redisClient, "testKey")

	// Acquire and Release Lock
	testKey := "testKey"
	locked := redisStore.AcquireLock(testKey, 0)
	assert.True(t, locked)

	err = redisStore.ReleaseLock(testKey)
	assert.NoError(t, err, "Release should not return an error")
}
