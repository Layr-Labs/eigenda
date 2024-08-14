package kvstore

import (
	"sync"
	"time"
)

var _ KVStore = &threadSafeStore{}

type threadSafeStore struct {
	store KVStore
	lock  sync.RWMutex
}

// ThreadSafeWrapper creates returns a wrapper around a KVStore that makes store access thread safe.
func ThreadSafeWrapper(store KVStore) KVStore {
	return &threadSafeStore{
		store: store,
		lock:  sync.RWMutex{},
	}
}

// Put stores a data in the store.
func (store *threadSafeStore) Put(key []byte, value []byte, ttl time.Duration) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Put(key, value, ttl)
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *threadSafeStore) Get(key []byte) ([]byte, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	return store.store.Get(key)
}

// Drop deletes data from the store.
func (store *threadSafeStore) Drop(key []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Drop(key)
}

// BatchUpdate performs a batch of Put and Drop operations.
func (store *threadSafeStore) BatchUpdate(puts []PutOperation, drops []DropOperation) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.BatchUpdate(puts, drops)
}

// Shutdown shuts down the store.
func (store *threadSafeStore) Shutdown() error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Shutdown()
}

// Destroy destroys the store.
func (store *threadSafeStore) Destroy() error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Destroy()
}

// IsShutDown returns true if the store is shut down.
func (store *threadSafeStore) IsShutDown() bool {
	return store.store.IsShutDown()
}
