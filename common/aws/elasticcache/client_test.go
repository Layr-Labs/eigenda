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

	// Test setting a value
	key := "testKey"
	value := "testValue"
	_, err = client.Set(context.Background(), key, value, 10*time.Second).Result()
	assert.NoError(t, err, "Set should not return an error")

	// Test getting the value
	stringCmd := client.Get(context.Background(), key)
	result, err := stringCmd.Result()
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, value, result, "Get should return the value that was set")
}
