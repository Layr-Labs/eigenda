package store

import (
	"context"
	"encoding/json"

	"github.com/Layr-Labs/eigenda/common"
	commoncache "github.com/Layr-Labs/eigenda/common/aws/elasticcache"
)

type RedisStore[T any] struct {
	client *commoncache.RedisClient
}

func NewRedisStore[T any](client *commoncache.RedisClient) common.KVStore[T] {
	return &RedisStore[T]{client: client}
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

	return s.client.Set(ctx, key, jsonData, 0).Err() // 0 means no expiration
}
