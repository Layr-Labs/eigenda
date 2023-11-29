package store_test

import (
	"context"
	"log"
	"os"
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

var pool *dockertest.Pool
var resource *dockertest.Resource

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "latest",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "", HostPort: "6379"}},
		},
	})
	if err != nil {
		log.Fatalf("Could not start Redis container: %v", err)
	}

	// Wait for Redis to be ready
	if err := pool.Retry(func() error {
		// Perform a health check...
		return nil // return nil if healthy
	}); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Run tests
	code := m.Run()

	// Teardown: Stop and remove the Redis container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Redis container: %v", err)
	}

	os.Exit(code)
}

func TestRedisStore(t *testing.T) {

	// Set up the Redis client to point to your local Redis server
	clientConfig := elasticcache.RedisClientConfig{
		EndpointURL: "localhost",
		Port:        "6379",
	}

	redisClient, err := elasticcache.NewClient(clientConfig, &cmock.Logger{}) // Assuming logger can be nil
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	redisStore := store.NewRedisStore[common.RateBucketParams](redisClient)

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
