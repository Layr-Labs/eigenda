package kvstore

import (
	"sync"
	"time"
)

var _ KVStore = &threadSafeKVStore{}

type threadSafeKVStore struct {
	store KVStore
	lock  sync.RWMutex
}

// ThreadSafeWrapper creates returns a wrapper around a KVStore that makes store access thread safe.
func ThreadSafeWrapper(store KVStore) KVStore {
	return &threadSafeKVStore{
		store: store,
		lock:  sync.RWMutex{},
	}
}

// Put stores a data in the store.
func (store *threadSafeKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Put(key, value, ttl)
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *threadSafeKVStore) Get(key []byte) ([]byte, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	return store.store.Get(key)
}

// Drop deletes data from the store.
func (store *threadSafeKVStore) Drop(key []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Drop(key)
}

// Shutdown shuts down the store.
func (store *threadSafeKVStore) Shutdown() error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Shutdown()
}

// Destroy destroys the store.
func (store *threadSafeKVStore) Destroy() error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Destroy()
}
