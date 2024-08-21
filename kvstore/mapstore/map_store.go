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

func (store *mapStore) Put(key []byte, value []byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

func (store *mapStore) Delete(key []byte) error {
	if store.destroyed {
		return fmt.Errorf("mapStore is destroyed")
	}

	stringifiedKey := string(key)
	delete(store.data, stringifiedKey)
	return nil
}

func (store *mapStore) DeleteBatch(keys [][]byte) error {
	for _, key := range keys {
		err := store.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func (store *mapStore) NewIterator(prefix []byte) iterator.Iterator {
	//TODO implement me
	// TODO unit test
	panic("implement me")
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

// IsShutDown returns true if the mapStore is shut down.
func (store *mapStore) IsShutDown() bool {
	return store.destroyed // TODO necessary?
}
