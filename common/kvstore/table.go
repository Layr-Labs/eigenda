package kvstore

import (
	"errors"
	"time"
)

// ErrTableNotFound is returned when a table is not found.
var ErrTableNotFound = errors.New("table not found")

// Table can be used to operate on data in a specific table in a TableStore.
type Table interface {
	Store

	// Name returns the name of the table.
	Name() string

	// TableKey creates a new key scoped to this table that can be used for TableStoreBatch
	// operations that modify this table. Using keys in TableStore batches that are not created using this method
	// has undefined behavior. Use of this method in a TableStoreBatch is not optional.
	TableKey(key []byte) []byte

	// PutWithTTL adds a key-value pair to the store that expires after a specified duration.
	// Key is eventually deleted after the TTL elapses.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithTTL(key []byte, value []byte, ttl time.Duration) error

	// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
	// Key is eventually deleted after the expiry time.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithExpiration(key []byte, value []byte, expiryTime time.Time) error

	// NewTTLBatch creates a new TTLBatch that can be used to perform multiple operations atomically.
	// Use this instead of NewBatch to create a batch that supports TTL/expiration.
	NewTTLBatch() TTLBatch
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

	// GetTable gets the table with the given name. If the table does not exist, it is first created.
	// Returns ErrTableNotFound if the table does not exist and cannot be created.
	GetTable(name string) (Table, error)

	// GetTables returns a list of all tables in the store in no particular order.
	GetTables() []Table

	// NewBatch creates a new batch that can be used to perform multiple operations across tables atomically.
	NewBatch() TTLBatch

	// Shutdown shuts down the store, flushing any remaining data to disk.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	Destroy() error
}
