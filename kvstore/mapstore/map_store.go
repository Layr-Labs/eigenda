package mapstore

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var _ kvstore.Store = &mapStore{}

// mapStore is a simple in-memory implementation of KVStore.
type mapStore struct {
	data      map[string][]byte
	destroyed bool
}

// NewStore creates a new mapStore.
func NewStore() kvstore.Store {
	return &mapStore{
		data: make(map[string][]byte),
	}
}

// Put adds a key-value pair to the store.
func (store *mapStore) Put(key []byte, value []byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

// Delete removes a key-value pair from the store.
func (store *mapStore) Delete(key []byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	stringifiedKey := string(key)
	delete(store.data, stringifiedKey)
	return nil
}

// DeleteBatch removes multiple key-value pairs from the store.
func (store *mapStore) DeleteBatch(keys [][]byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	for _, key := range keys {
		err := store.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteBatch adds multiple key-value pairs to the store.
func (store *mapStore) WriteBatch(keys, values [][]byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	if len(keys) != len(values) {
		return fmt.Errorf("keys and values slices must have the same length")
	}

	for i, key := range keys {
		err := store.Put(key, values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// NewIterator creates a new iterator for the store. Only keys with the given prefix are returned.
// WARNING: this implementation does not take a snapshot of the store, and so the iterator may return a
// consistent view of the store if there is concurrent modification.
func (store *mapStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	if store.destroyed {
		return nil, fmt.Errorf("mapStore is destroyed")
	}

	// This will not be implemented until we upgrade to go 1.23, which standardizes iterators as part of the language.

	return nil, nil
}

// Get retrieves data from the mapStore. Returns nil if the data is not found.
func (store *mapStore) Get(key []byte) ([]byte, error) {
	if store.destroyed {
		return nil, fmt.Errorf("mapStore is destroyed")
	}

	stringifiedKey := string(key)

	data, ok := store.data[stringifiedKey]

	if !ok {
		return nil, kvstore.ErrNotFound
	}

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	return dataCopy, nil // TODO test that it is safe to modify the returned data
}

// Shutdown stops the mapStore and releases any resources it holds. Does not delete any on-disk data.
func (store *mapStore) Shutdown() error {
	return store.Destroy()
}

// Destroy permanently stops the mapStore and deletes all data (including data on disk).
func (store *mapStore) Destroy() error {
	store.data = nil
	store.destroyed = true
	return nil
}