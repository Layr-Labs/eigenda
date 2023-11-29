package elasticcache_test

import (
	"context"
	"log"
	"testing"
	"time"

	elasticCache "github.com/Layr-Labs/eigenda/common/aws/elasticcache"
	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient(t *testing.T) {
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

	// Set up Redis client
	cfg := elasticCache.RedisClientConfig{
		EndpointURL: "localhost",
		Port:        "6379",
	}

	logger := &cmock.Logger{}
	client, err := elasticCache.NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}
	assert.NoError(t, err, "NewClient should not return an error")
	assert.NotNil(t, client, "RedisClient should not be nil")

	// Test Set method
	key := "testKey"
	value := "testValue"
	_, err = client.Set(context.Background(), key, value, 0) // 0 expiration means no expiration
	assert.NoError(t, err, "Set should not return an error")

	// Test Get method
	getCmd := client.Get(context.Background(), key)
	getResult, err := getCmd.Result()
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, value, getResult, "Get should return the value that was set")

	// Test AcquireLock and ReleaseLock methods
	lockKey := "testLockKey"
	lockValue := "uniqueLockValue"
	assert.True(t, client.AcquireLock(lockKey, lockValue, time.Second*10), "AcquireLock should return true")
	assert.NoError(t, client.ReleaseLock(lockKey, lockValue), "ReleaseLock should not return an error")
}
