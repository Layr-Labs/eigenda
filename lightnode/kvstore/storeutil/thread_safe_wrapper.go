package storeutil

import (
	"github.com/Layr-Labs/eigenda/lightnode/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"sync"
)

var _ kvstore.Store = &threadSafeStore{}

// threadSafeStore is a wrapper around a KVStore that makes store access thread safe.
type threadSafeStore struct {
	store kvstore.Store
	lock  sync.RWMutex
}

// ThreadSafeWrapper creates returns a wrapper around a KVStore that makes store access thread safe.
func ThreadSafeWrapper(store kvstore.Store) kvstore.Store {
	return &threadSafeStore{
		store: store,
		lock:  sync.RWMutex{},
	}
}

// Put stores a data in the store.
func (store *threadSafeStore) Put(key []byte, value []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Put(key, value)
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *threadSafeStore) Get(key []byte) ([]byte, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	return store.store.Get(key)
}

// Delete deletes data from the store.
func (store *threadSafeStore) Delete(key []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.Delete(key)
}

// DeleteBatch deletes multiple key-value pairs from the store.
func (store *threadSafeStore) DeleteBatch(keys [][]byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.DeleteBatch(keys)
}

// WriteBatch adds multiple key-value pairs to the store.
func (store *threadSafeStore) WriteBatch(keys, values [][]byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.store.WriteBatch(keys, values)
}

// NewIterator creates a new iterator.
func (store *threadSafeStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	return store.store.NewIterator(prefix)
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
