package common

import (
	"context"
)

// KVStore is a simple key value store interface.
type KVStore[T any] interface {
	// GetItem returns the value associated with a given key.
	GetItem(ctx context.Context, key string) (*T, error)
	// UpdateItem updates the value for the given key.
	UpdateItem(ctx context.Context, key string, value *T) error
}

// KVStoreVersioned extends KVStore with version information
type KVStoreVersioned[T any] interface {
	KVStore[T] // Embeds KVStore interface

	// GetItemWithVersion returns the value associated with a given key and version
	GetItemWithVersion(ctx context.Context, key string) (*T, int, error)
	// UpdateItem updates the value for the given key with version
	UpdateItemWithVersion(ctx context.Context, key string, value *T, expectedVersion int) error
}
