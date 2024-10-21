package mapstore

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sort"
	"sync"
)

var _ kvstore.Store[[]byte] = &mapStore{}

// mapStore is a simple in-memory implementation of KVStore. Designed more as a correctness test than a
// production implementation -- there are things that may not be performant with this implementation.
type mapStore struct {
	data map[string][]byte
	lock sync.RWMutex
}

// NewStore creates a new mapStore.
func NewStore() kvstore.Store[[]byte] {
	return &mapStore{
		data: make(map[string][]byte),
	}
}

// Put adds a key-value pair to the store.
func (store *mapStore) Put(key []byte, value []byte) error {
	if value == nil {
		value = []byte{}
	}

	store.lock.Lock()
	defer store.lock.Unlock()

	stringifiedKey := string(key)

	store.data[stringifiedKey] = value
	return nil
}

// Delete removes a key-value pair from the store.
func (store *mapStore) Delete(key []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	stringifiedKey := string(key)
	delete(store.data, stringifiedKey)
	return nil
}

// DeleteBatch removes multiple key-value pairs from the store.
func (store *mapStore) DeleteBatch(keys [][]byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	for _, key := range keys {
		stringifiedKey := string(key)
		delete(store.data, stringifiedKey)
	}
	return nil
}

// WriteBatch adds multiple key-value pairs to the store.
func (store *mapStore) WriteBatch(keys, values [][]byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	if len(keys) != len(values) {
		return fmt.Errorf("keys and values slices must have the same length")
	}

	for i, key := range keys {
		stringifiedKey := string(key)
		store.data[stringifiedKey] = values[i]
	}
	return nil
}

// NewBatch creates a new batch for the store.
func (store *mapStore) NewBatch() kvstore.Batch[[]byte] {
	return &batch{
		store:  store,
		keys:   make([][]byte, 0),
		values: make([][]byte, 0),
	}
}

// batch is a collection of operations that can be applied atomically to a mapStore.
type batch struct {
	store  *mapStore
	keys   [][]byte
	values [][]byte
}

// Put stores the given key / value pair in the batch, overwriting any existing value for that key.
func (m *batch) Put(key []byte, value []byte) {
	if value == nil {
		value = []byte{}
	}
	m.keys = append(m.keys, key)
	m.values = append(m.values, value)
}

// Delete removes the key from the batch. Does not return an error if the key does not exist.
func (m *batch) Delete(key []byte) {
	m.keys = append(m.keys, key)
	m.values = append(m.values, nil)
}

// Apply atomically writes & deletes all the key / value pairs in the batch to the database.
func (m *batch) Apply() error {
	for i, key := range m.keys {
		if m.values[i] == nil {
			err := m.store.Delete(key)
			if err != nil {
				return err
			}
		} else {
			err := m.store.Put(key, m.values[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Size returns the number of operations in the batch.
func (m *batch) Size() uint32 {
	return uint32(len(m.keys))
}

type mapIterator struct {
	keys         []string
	values       map[string][]byte
	currentIndex int
}

func (it *mapIterator) First() bool {
	if len(it.keys) == 0 {
		return false
	}
	it.currentIndex = 0
	return true
}

func (it *mapIterator) Last() bool {
	if len(it.keys) == 0 {
		return false
	}
	it.currentIndex = len(it.keys) - 1
	return true
}

func (it *mapIterator) Seek(key []byte) bool {
	// Not efficient. But then again, nothing is efficient in this iterator implementation.
	for i, k := range it.keys {
		if k == string(key) {
			it.currentIndex = i
			return true
		}
	}
	return false
}

func (it *mapIterator) Next() bool {
	if it.currentIndex == len(it.keys)-1 {
		return false
	}
	it.currentIndex++
	return true
}

func (it *mapIterator) Prev() bool {
	if it.currentIndex == 0 {
		return false
	}
	it.currentIndex--
	return true
}

func (it *mapIterator) Release() {
	// no op
}

func (it *mapIterator) SetReleaser(releaser util.Releaser) {
	// no op
}

func (it *mapIterator) Valid() bool {
	return true
}

func (it *mapIterator) Error() error {
	return nil
}

func (it *mapIterator) Key() []byte {
	return []byte(it.keys[it.currentIndex])
}

func (it *mapIterator) Value() []byte {
	return it.values[it.keys[it.currentIndex]]
}

// NewIterator creates a new iterator for the store. Only keys with the given prefix are returned.
// WARNING: this implementation does a full copy to return the iterator. This is not efficient.
// Not for production use.
func (store *mapStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	mapCopy := make(map[string][]byte)
	keys := make([]string, 0, len(store.data))

	for k, v := range store.data {
		keyBytes := []byte(k)

		if len(keyBytes) < len(prefix) {
			// Key is shorter than the prefix, so it can't have the prefix
			continue
		}

		// check if the key has the prefix
		matchesPrefix := true
		for i, b := range prefix {
			if keyBytes[i] != b {
				matchesPrefix = false
				break
			}
		}
		if !matchesPrefix {
			continue
		}

		mapCopy[k] = v
		keys = append(keys, k)
	}

	// Iterator must walk over keys in lexicographical order
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return &mapIterator{
		keys:         keys,
		values:       mapCopy,
		currentIndex: -1,
	}, nil
}

// Get retrieves data from the mapStore. Returns nil if the data is not found.
func (store *mapStore) Get(key []byte) ([]byte, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	data, ok := store.data[string(key)]

	if !ok {
		return nil, kvstore.ErrNotFound
	}

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	return dataCopy, nil
}

// Shutdown stops the mapStore and releases any resources it holds. Does not delete any on-disk data.
func (store *mapStore) Shutdown() error {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.data = nil
	return nil
}

// Destroy permanently stops the mapStore and deletes all data (including data on disk).
func (store *mapStore) Destroy() error {
	return store.Shutdown()
}
