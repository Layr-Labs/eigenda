package elasticcache_test

import (
	"context"
	"testing"
	"time"

	elasticCache "github.com/Layr-Labs/eigenda/common/aws/elasticcache"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient(t *testing.T) {
	// Set up the Redis client
	cfg := elasticCache.RedisClientConfig{
		EndpointURL: "localhost",
		Port:        "6379", // Assuming Redis is running on the default port
	}

	client, err := elasticCache.NewClient(cfg, nil)
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
