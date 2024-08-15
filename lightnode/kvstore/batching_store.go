package kvstore

import (
	"fmt"
	"time"
)

var _ KVStore = &batchingStore{}

// batchingStore is a wrapper around a KVStore that batches Put and Drop operations.
// This may (significantly) improve performance for stores where batching is helpful.
// Note that if the store is not properly shutdown, some data may be lost.
type batchingStore struct {
	cache map[string]*BatchOperation // TODO is there a way not to have to stringify?
	store KVStore

	maxCacheSize uint
	cacheSize    uint
}

// TODO should the key size be included in the cache size? probably it should

// BatchingWrapper creates a wrapper around a KVStore that batches Put and Drop operations. The maxCacheSize parameter
// the number of bytes that can be stored in the cache before it is flushed to the underlying store.
func BatchingWrapper(store KVStore, maxCacheSize uint) KVStore {
	return &batchingStore{
		cache:        make(map[string]*BatchOperation),
		store:        store,
		maxCacheSize: maxCacheSize,
	}
}

// Put stores a data in the store.
func (store *batchingStore) Put(key []byte, value []byte, ttl time.Duration) error {

	stringifiedKey := string(key)

	previousCacheEntry := store.cache[stringifiedKey]
	if previousCacheEntry != nil {
		store.cacheSize -= uint(len(previousCacheEntry.Value))
	}

	store.cache[stringifiedKey] = &BatchOperation{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	store.cacheSize += uint(len(value))

	if store.maxCacheSize >= store.maxCacheSize {
		return store.flushCache()
	}

	return nil
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *batchingStore) Get(key []byte) ([]byte, error) {
	if store.store.IsShutDown() {
		return nil, fmt.Errorf("store is offline")
	}

	operation, ok := store.cache[string(key)]

	if !ok {
		// Entry is not in the write cache. Fetch from store.
		return store.store.Get(key)
	}

	if operation == nil {
		// key is dropped from cache
		return nil, nil
	}

	return operation.Value, nil
}

// Drop deletes data from the store.
func (store *batchingStore) Drop(key []byte) error {
	if store.store.IsShutDown() {
		return fmt.Errorf("store is offline")
	}

	stringifiedKey := string(key)

	previousCacheEntry := store.cache[stringifiedKey]
	if previousCacheEntry != nil {
		store.cacheSize -= uint(len(previousCacheEntry.Value))
	}

	store.cache[string(key)] = nil
	return nil
}

// BatchUpdate performs a batch of Put and Drop operations.
func (store *batchingStore) BatchUpdate(operations []*BatchOperation) error {
	if store.store.IsShutDown() {
		return fmt.Errorf("store is offline")
	}

	for _, operation := range operations {
		stringifiedKey := string(operation.Key)

		previousCacheEntry, present := store.cache[stringifiedKey]
		if present && previousCacheEntry != nil {
			store.cacheSize -= uint(len(previousCacheEntry.Value))
		}

		store.cache[stringifiedKey] = operation

		if operation.Value != nil {
			store.cacheSize += uint(len(operation.Value))
		}
	}

	if store.cacheSize >= store.maxCacheSize {
		return store.flushCache()
	}

	return nil
}

// Shutdown shuts down the store.
func (store *batchingStore) Shutdown() error {
	if store.store.IsShutDown() {
		return nil
	}

	// TODO write tests for shutdown behavior
	err := store.flushCache()
	if err != nil {
		return err
	}
	return store.store.Shutdown()
}

// Destroy destroys the store.
func (store *batchingStore) Destroy() error {
	err := store.Shutdown()
	if err != nil {
		return err
	}
	return store.store.Destroy()
}

// flushCache flushes the cache to the underlying store.
func (store *batchingStore) flushCache() error {
	for key, operation := range store.cache {
		if operation == nil {
			err := store.store.Drop([]byte(key))
			if err != nil {
				return err
			}
		} else {
			err := store.store.Put(operation.Key, operation.Value, operation.TTL)
			if err != nil {
				return err
			}
		}
	}

	store.cacheSize = 0
	store.cache = make(map[string]*BatchOperation)

	return nil
}

// IsShutDown returns true if the store is shut down.
func (store *batchingStore) IsShutDown() bool {
	return store.store.IsShutDown()
}
