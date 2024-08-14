package kvstore

import (
	"fmt"
	"time"
)

var _ KVStore = &InMemoryKVStore{}

// TODO create priority queue for TTL

// InMemoryKVStore is a simple in-memory implementation of KVStore.
type InMemoryKVStore struct {
	data      map[string][]byte
	destroyed bool
}

// NewInMemoryChunkStore creates a new InMemoryKVStore.
func NewInMemoryChunkStore() *InMemoryKVStore {
	return &InMemoryKVStore{
		data: make(map[string][]byte),
	}
}

// Put stores a data in the store.
func (store *InMemoryKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.destroyed {
		return fmt.Errorf("store is destroyed")
	}

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *InMemoryKVStore) Get(key []byte) ([]byte, error) {
	if store.destroyed {
		return nil, fmt.Errorf("store is destroyed")
	}

	stringifiedKey := string(key)

	data, ok := store.data[stringifiedKey]

	if !ok {
		return nil, nil
	}

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	return dataCopy, nil // TODO test that it is safe to modify the returned data
}

// Drop removes data from the store.
func (store *InMemoryKVStore) Drop(key []byte) error {
	if store.destroyed {
		return fmt.Errorf("store is destroyed")
	}

	stringifiedKey := string(key)
	delete(store.data, stringifiedKey)
	return nil
}

// Shutdown stops the store and releases any resources it holds. Does not delete any on-disk data.
func (store *InMemoryKVStore) Shutdown() error {
	return store.Destroy()
}

// Destroy permanently stops the store and deletes all data (including data on disk).
func (store *InMemoryKVStore) Destroy() error {
	store.data = nil
	store.destroyed = true
	return nil
}
