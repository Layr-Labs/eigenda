package elasticcache

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/go-redis/redis/v8"
)

type RedisClientConfig struct {
	EndpointURL string
	Port        string
}

type RedisClient struct {
	redisClient *redis.Client
	logger      common.Logger // Ensure common.Logger is imported correctly
}

func NewClient(cfg RedisClientConfig, logger common.Logger) (*RedisClient, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.EndpointURL + ":" + cfg.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test the Redis connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err // Return the error instead of logging and exiting
	}
	logger.Info("Redis connection successful")

	return &RedisClient{redisClient: redisClient, logger: logger}, nil
}

// Get retrieves a value from Redis
func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.redisClient.Get(ctx, key)
}

// Set sets a value in Redis
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (*redis.StatusCmd, error) {

	return c.redisClient.Set(ctx, key, value, expiration), nil
}

func (c *RedisClient) AcquireLock(lockKey string, lockValue string, expiration time.Duration) bool {
	result, err := c.redisClient.SetNX(context.Background(), lockKey, lockValue, expiration).Result()
	if err != nil {
		// Handle error
		return false
	}
	return result
}

func (c *RedisClient) ReleaseLock(lockKey string, lockValue string) error {
	script := `
    if redis.call("get",KEYS[1]) == ARGV[1] then
        return redis.call("del",KEYS[1])
    else
        return 0
    end
    `
	_, err := c.redisClient.Eval(context.Background(), script, []string{lockKey}, lockValue).Result()
	if err != nil {
		return err
	}

	return nil
}
