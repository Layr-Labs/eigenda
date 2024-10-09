package kvstore

import "errors"

// ErrTableLimitExceeded is returned when the maximum number of tables has been reached.
var ErrTableLimitExceeded = errors.New("table limit exceeded")

// ErrTableNotFound is returned when a table is not found.
var ErrTableNotFound = errors.New("table not found")

// Table can be used to operate on data in a specific table in a TableStore.
type Table interface {
	// Store permits access to the table as if it were a store.
	Store

	// Name returns the name of the table.
	Name() string

	// TableKey creates a new key scoped to this table that can be used for batch operations that modify this table.
	TableKey(key []byte) TableKey
}

// TableKey is a key scoped to a particular table. It can be used to perform batch operations that modify multiple
// table keys atomically.
type TableKey []byte

// TableBatch is a collection of operations that can be applied atomically to a TableStore.
type TableBatch Batch[TableKey]

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
	NewBatch() TableBatch

	// Shutdown shuts down the store, flushing any remaining data to disk.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	Destroy() error
}
