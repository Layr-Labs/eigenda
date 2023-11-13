package common

import "context"

// KVStore is a simple key value store interface.
type KVStore[T any] interface {
	// GetItem returns the value associated with a given key.
	GetItem(ctx context.Context, key string) (*T, error)
	// UpdateItem updates the value for the given key.
	UpdateItem(ctx context.Context, key string, value *T) error
}
