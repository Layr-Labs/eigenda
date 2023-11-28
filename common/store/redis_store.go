package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	elasticCache "github.com/Layr-Labs/eigenda/common/aws/elasticcache"
)

type RedisStore[T any] struct {
	client    *elasticCache.RedisClient
	lockValue string // global lock value
}

func NewRedisStore[T any](client *elasticCache.RedisClient, lockValue string) common.LockableKVStore[T] {

	// LockValue is an identifier for the application initializing this store
	return &RedisStore[T]{
		client:    client,
		lockValue: lockValue,
	}
}

func (s *RedisStore[T]) GetItem(ctx context.Context, key string) (*T, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var item T
	err = json.Unmarshal([]byte(val), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *RedisStore[T]) UpdateItem(ctx context.Context, key string, value *T) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.client.Set(ctx, key, jsonData, 0) // 0 means no expiration
	return err
}

// AcquireLock implementation for RedisStore
func (s *RedisStore[T]) AcquireLock(key string, expiration time.Duration) bool {
	return s.client.AcquireLock(key, s.lockValue, expiration)
}

// ReleaseLock implementation for RedisStore
func (s *RedisStore[T]) ReleaseLock(key string) error {
	return s.client.ReleaseLock(key, s.lockValue)
}
