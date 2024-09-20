package kvstore

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Key is a key in a TableStore. It is a combination of a table name and a key within that table.
type Key interface {
	// GetKeyBytes returns the key within the table, interpreted as a byte slice.
	GetKeyBytes() []byte

	// GetKeyString returns the key within the table, interpreted as a string. Calling this
	// method on keys that do not represent a string may return odd results.
	GetKeyString() string

	// GetKeyUint32 returns the key within the table, interpreted as a uint32. Calling this
	// method on keys that do not represent a uint32 may return odd results.
	GetKeyUint32() uint32

	// GetKeyUint64 returns the key within the table, interpreted as a uint64. Calling this
	// method on keys that do not represent a uint64 may return odd results.
	GetKeyUint64() uint64

	// GetRawBytes gets the representation of the key as used internally by the store.
	GetRawBytes() []byte
}

// KeyBuilder is used to create new keys in a specific table.
type KeyBuilder interface {
	// Key creates a new key in a specific table using the given key bytes.
	Key(key []byte) Key

	// StringKey creates a new key in a specific table using the given key string.
	StringKey(key string) Key

	// Uint32Key creates a new key in a specific table using the given uint32 as a key.
	Uint32Key(key uint32) Key

	// Uint64Key creates a new key in a specific table using the given uint64 as a key.
	Uint64Key(key uint64) Key
}

// TableStore implements a key-value store, with the addition of the abstraction of tables.
// A "table" in this context is a disjoint keyspace. Keys in one table to not collide with keys in another table,
// and keys within a particular table can be iterated over efficiently.
//
// A TableStore is only required to support a maximum of 2^32-X unique, where X is a small integer number of tables
// reserved for internal use by the table store. The exact value of X is implementation dependent.
//
// Implementations of this interface are expected to be thread-safe, except where noted.
type TableStore interface {

	// GetOrCreateTable creates a new table with the given name if one does not exist
	// and returns a key builder for that table.
	//
	// WARNING: this method is not thread safe with respect to any other methods in this interface.
	GetOrCreateTable(name string) (KeyBuilder, error)

	// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
	//
	// WARNING: this method is not thread safe with respect to any other methods in this interface.
	DropTable(name string) error

	// GetKeyBuilder returns a key builder for the table with the given name,
	// returning an error if the table does not exist.
	GetKeyBuilder(name string) (KeyBuilder, error)

	// GetMaxTableCount returns the maximum number of tables that can be created in the store
	// (excluding internal tables utilized by the store).
	GetMaxTableCount() uint32

	// GetCurrentTableCount returns the current number of tables in the store
	// (excluding internal tables utilized by the store).
	GetCurrentTableCount() uint32

	// GetTables returns a list of all tables in the store (excluding internal tables utilized by the store).
	GetTables() []string

	// Put stores the given key / value pair in the database, overwriting any existing value for that key.
	Put(key Key, value []byte) error

	// Get retrieves the value for the given key. Returns a ErrNotFound error if the key does not exist.
	// The value returned is safe to modify.
	Get(key Key) ([]byte, error)

	// Delete removes the key from the database. Does not return an error if the key does not exist.
	Delete(key Key) error

	// DeleteBatch atomically removes a list of keys from the database.
	DeleteBatch(keys []Key) error

	// WriteBatch atomically writes a list of key / value pairs to the database. The key at index i in the keys slice
	// corresponds to the value at index i in the values slice.
	WriteBatch(keys []Key, values [][]byte) error

	// NewIterator returns an iterator that can be used to iterate over a subset of the keys in the database.
	// Only keys with the given key's table with prefix matching the key will be iterated. The iterator must be closed
	// by calling Release() when done. The iterator will return keys in lexicographically sorted order. The iterator
	// walks over a consistent snapshot of the database, so it will not see any writes that occur after the iterator
	// is created.
	NewIterator(prefix Key) (iterator.Iterator, error)

	// Shutdown shuts down the store, flushing any remaining data to disk.
	//
	// Warning: it is not thread safe to call this method concurrently with other methods in this interface,
	// or while there exist unclosed iterators.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	//
	// Warning: it is not thread safe to call this method concurrently with other methods in this interface,
	// or while there exist unclosed iterators.
	Destroy() error
}
