package kvstore

import (
	"fmt"
	"time"
)

var _ KVStore = &InMemoryStore{}

// TODO create priority queue for TTL

// InMemoryStore is a simple in-memory implementation of KVStore.
type InMemoryStore struct {
	data      map[string][]byte
	destroyed bool
}

// NewInMemoryStore creates a new InMemoryStore.
func NewInMemoryStore() KVStore {
	return &InMemoryStore{
		data: make(map[string][]byte),
	}
}

// Put stores a data in the store.
func (store *InMemoryStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.destroyed {
		return fmt.Errorf("store is destroyed")
	}

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *InMemoryStore) Get(key []byte) ([]byte, error) {
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

// BatchUpdate performs a batch of Put and Drop operations.
func (store *InMemoryStore) BatchUpdate(puts []PutOperation, drops []DropOperation) error {
	for _, put := range puts {
		err := store.Put(put.Key, put.Value, put.TTL)
		if err != nil {
			return err
		}
	}

	for _, drop := range drops {
		err := store.Drop(drop)
		if err != nil {
			return err
		}
	}

	return nil
}

// Drop removes data from the store.
func (store *InMemoryStore) Drop(key []byte) error {
	if store.destroyed {
		return fmt.Errorf("store is destroyed")
	}

	stringifiedKey := string(key)
	delete(store.data, stringifiedKey)
	return nil
}

// Shutdown stops the store and releases any resources it holds. Does not delete any on-disk data.
func (store *InMemoryStore) Shutdown() error {
	return store.Destroy()
}

// Destroy permanently stops the store and deletes all data (including data on disk).
func (store *InMemoryStore) Destroy() error {
	store.data = nil
	store.destroyed = true
	return nil
}

// IsShutDown returns true if the store is shut down.
func (store *InMemoryStore) IsShutDown() bool {
	return store.destroyed
}
