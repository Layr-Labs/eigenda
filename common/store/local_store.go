package store

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

type localParamStore[T any] struct {
	cache *lru.Cache[string, T]
}

func NewLocalParamStore[T any](size int) (common.KVStore[T], error) {
	cache, err := lru.New[string, T](size)
	if err != nil {
		return nil, err
	}

	return &localParamStore[T]{
		cache: cache,
	}, nil
}

func (s *localParamStore[T]) GetItem(ctx context.Context, key string) (*T, error) {

	obj, ok := s.cache.Get(key)
	if !ok {
		return nil, errors.New("error retrieving key")
	}

	return &obj, nil

}

func (s *localParamStore[T]) UpdateItem(ctx context.Context, key string, params *T) error {

	s.cache.Add(key, *params)

	return nil
}
