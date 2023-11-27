package store

import (
	"context"
	"encoding/json"

	"github.com/Layr-Labs/eigenda/common"
	elasticCache "github.com/Layr-Labs/eigenda/common/aws/elasticcache"
)

type RedisStore[T any] struct {
	client    *elasticCache.RedisClient
	lockKey   string
	lockValue string // Unique value for the lock, e.g., UUID or server identifier
}

func NewRedisStore[T any](client *elasticCache.RedisClient, lockKey string, lockValue string) common.KVStore[T] {
	return &RedisStore[T]{client: client,
		lockKey:   lockKey,
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

	_, err = s.client.Set(ctx, key, jsonData, s.lockKey, s.lockValue, 0) // 0 means no expiration
	return err
}
