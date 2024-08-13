package kvstore

import (
	"fmt"
	"sync"
	"time"
)

var _ KVStore = &InMemoryKVStore{}

// TODO create priority queue for TTL

// InMemoryKVStore is a simple in-memory implementation of KVStore.
type InMemoryKVStore struct {
	data map[string][]byte
	lock sync.RWMutex
}

// NewInMemoryChunkStore creates a new InMemoryKVStore.
func NewInMemoryChunkStore() *InMemoryKVStore {
	return &InMemoryKVStore{
		data: make(map[string][]byte),
	}
}

// Put stores a data in the store.
func (store *InMemoryKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

// Get retrieves data from the store.
func (store *InMemoryKVStore) Get(key []byte) ([]byte, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	stringifiedKey := string(key)

	data, ok := store.data[stringifiedKey]

	if !ok {
		return nil, fmt.Errorf("data not found for key %s", stringifiedKey)
	}

	return data, nil
}
