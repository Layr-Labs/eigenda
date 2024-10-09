package kvstore

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// ErrNotFound is returned when a key is not found in the database.
var ErrNotFound = errors.New("not found")

// StoreBatch is a collection of operations that can be applied atomically to a Store.
type StoreBatch Batch[[]byte]

// Store implements a key-value store. May be backed by a database like LevelDB.
//
// Implementations of this interface are expected to be thread-safe.
type Store interface {

	// Put stores the given key / value pair in the database, overwriting any existing value for that key.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	Put(key []byte, value []byte) error

	// Get retrieves the value for the given key. Returns a ErrNotFound error if the key does not exist.
	// The value returned is safe to modify.
	Get(key []byte) ([]byte, error)

	// Delete removes the key from the database. Does not return an error if the key does not exist.
	Delete(key []byte) error

	// NewBatch creates a new batch that can be used to perform multiple operations atomically.
	NewBatch() StoreBatch

	// NewIterator returns an iterator that can be used to iterate over a subset of the keys in the database.
	// Only keys with the given prefix will be iterated. The iterator must be closed by calling Release() when done.
	// The iterator will return keys in lexicographically sorted order. The iterator walks over a consistent snapshot
	// of the database, so it will not see any writes that occur after the iterator is created.
	NewIterator(prefix []byte) (iterator.Iterator, error)

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
