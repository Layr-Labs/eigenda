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
