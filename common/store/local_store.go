package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

type lockInfo struct {
	mu     sync.Mutex
	locked bool
}

// Global map for locks
var globalLocks sync.Map

type localParamStore[T any] struct {
	cache     *lru.Cache[string, T]
	lockValue string // global lock value
}

func NewLocalParamStore[T any](size int, lockValue string) (common.LockableKVStore[T], error) {
	cache, err := lru.New[string, T](size)
	if err != nil {
		return nil, err
	}

	return &localParamStore[T]{
		cache:     cache,
		lockValue: lockValue,
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

// AcquireLock implementation for RedisStore
func (s *localParamStore[T]) AcquireLock(key string, expiration time.Duration) bool {

	val, _ := globalLocks.LoadOrStore(key, &lockInfo{})
	lock := val.(*lockInfo)

	lock.mu.Lock()
	defer lock.mu.Unlock()

	if !lock.locked {
		lock.locked = true
		return true
	}
	return false
}

// ReleaseLock implementation for RedisStore
func (s *localParamStore[T]) ReleaseLock(key string) error {
	val, ok := globalLocks.Load(key)
	if !ok {
		return errors.New("no lock found for the given key")
	}

	lock := val.(*lockInfo)

	lock.mu.Lock()
	defer lock.mu.Unlock()

	if lock.locked {
		lock.locked = false
		return nil
	}
	return errors.New("lock not acquired or already released")
}
