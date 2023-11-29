package common

import (
	"context"
	"time"
)

// KVStore is a simple key value store interface.
type KVStore[T any] interface {
	// GetItem returns the value associated with a given key.
	GetItem(ctx context.Context, key string) (*T, error)
	// UpdateItem updates the value for the given key.
	UpdateItem(ctx context.Context, key string, value *T) error
}

// LockableKVStore extends KVStore with lock and unlock capabilities.
type LockableKVStore[T any] interface {
	KVStore[T] // Embedding KVStore

	// AcquireLock tries to acquire a lock and returns true if successful.
	AcquireLock(lockKey string, expiration time.Duration) bool
	// ReleaseLock releases the acquired lock.
	ReleaseLock(lockKey string) error
}
