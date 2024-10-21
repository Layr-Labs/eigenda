package kvstore

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// ErrNotFound is returned when a key is not found in the database.
var ErrNotFound = errors.New("not found")

// Store implements a key-value store. May be backed by a database like LevelDB.
// The generic type K is the type of the keys in the store.
//
// Implementations of this interface are expected to be thread-safe.
type Store[K any] interface {

	// Put stores the given key / value pair in the database, overwriting any existing value for that key.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	Put(k K, value []byte) error

	// Get retrieves the value for the given key. Returns a ErrNotFound error if the key does not exist.
	// The value returned is safe to modify.
	Get(k K) ([]byte, error)

	// Delete removes the key from the database. Does not return an error if the key does not exist.
	Delete(k K) error

	// NewBatch creates a new batch that can be used to perform multiple operations atomically.
	NewBatch() Batch[K]

	// NewIterator returns an iterator that can be used to iterate over a subset of the keys in the database.
	// Only keys with the given prefix will be iterated. The iterator must be closed by calling Release() when done.
	// The iterator will return keys in lexicographically sorted order. The iterator walks over a consistent snapshot
	// of the database, so it will not see any writes that occur after the iterator is created.
	NewIterator(prefix K) (iterator.Iterator, error)

	// Shutdown shuts down the store, flushing any remaining data to disk.
	//
	// Warning: it is not thread safe to call this method concurrently with other methods on this class,
	// or while there exist unclosed iterators.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	//
	// Warning: it is not thread safe to call this method concurrently with other methods on this class,
	// or while there exist unclosed iterators.
	Destroy() error
}
