package kvstore

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"time"
)

// ErrTableNotFound is returned when a table is not found.
var ErrTableNotFound = errors.New("table not found")

// TableStore implements a key-value store, with the addition of the abstraction of tables.
// A "table" in this context is a disjoint keyspace. Keys in one table to not collide with keys in another table,
// and keys within a particular table can be iterated over efficiently.
//
// A TableStore is only required to support a maximum of 2^32-X unique, where X is a small integer number of tables
// reserved for internal use by the table store. The exact value of X is implementation dependent.
//
// Implementations of this interface are expected to be thread-safe, except where noted.
type TableStore interface {
	Store[Key]

	// GetKeyBuilder gets the key builder for a particular table. Returns ErrTableNotFound if the table does not exist.
	// The returned KeyBuilder can be used to interact with the table.
	//
	// Warning: Do not use key builders created by one TableStore instance with another TableStore instance.
	// This may result in odd and undefined behavior.
	GetKeyBuilder(name string) (KeyBuilder, error)

	// GetKeyBuilders returns all key builders in the store.
	GetKeyBuilders() []KeyBuilder

	// GetTables returns a list of the table names currently in the store.
	GetTables() []string

	// PutWithTTL adds a key-value pair to the store that expires after a specified duration.
	// Key is eventually deleted after the TTL elapses.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithTTL(key Key, value []byte, ttl time.Duration) error

	// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
	// Key is eventually deleted after the expiry time.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithExpiration(key Key, value []byte, expiryTime time.Time) error

	// NewTTLBatch creates a new TTLBatch that can be used to perform multiple operations atomically.
	// Use this instead of NewBatch to create a batch that supports TTL/expiration.
	NewTTLBatch() TTLBatch[Key]

	// NewTableIterator returns an iterator that can be used to iterate over all keys in a table.
	// Equivalent to NewIterator(keyBuilder.Key([]byte{})).
	NewTableIterator(keyBuilder KeyBuilder) (iterator.Iterator, error)
}
