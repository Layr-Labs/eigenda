package kvstore

import "errors"

// TableKey is a key scoped to a particular table. It can be used to perform batch operations that modify multiple
// table keys atomically.
type TableKey []byte

// Table can be used to operate on data in a specific table in a TableStore.
type Table interface {
	// Store permits access to the table as if it were a store.
	Store

	// Name returns the name of the table.
	Name() string

	// TableKey creates a new key scoped to this table that can be used for batch operations that modify this table.
	TableKey(key []byte) TableKey
}

// ErrTableLimitExceeded is returned when the maximum number of tables has been reached.
var ErrTableLimitExceeded = errors.New("table limit exceeded")

// ErrTableNotFound is returned when a table is not found.
var ErrTableNotFound = errors.New("table not found")

// TableStoreBuilder is used to create a new TableStore instance. It can be used to add and remove
// tables from the store. Once the TableStore is created, the TableStoreBuilder instance should not be used again,
// and no tables may be added or removed from the store. If table modifications are required on an existing
// TableStore, it should first be shut down and a new TableStoreBuilder should be created.
type TableStoreBuilder interface {

	// CreateTable creates a new table with the given name. If a table with the given name already exists,
	// this method becomes a no-op. Returns ErrTableLimitExceeded if the maximum number of tables has been reached.
	CreateTable(name string) error

	// DropTable deletes the table with the given name. If the table does not exist, this method becomes a no-op.
	DropTable(name string) error

	// GetMaxTableCount returns the maximum number of tables that can be created in the store
	// (excluding internal tables utilized by the store).
	GetMaxTableCount() uint32

	// GetTableCount returns the current number of tables in the store
	// (excluding internal tables utilized by the store).
	GetTableCount() uint32

	// GetTableNames returns a list of the names of all tables in the store, in no particular order.
	GetTableNames() []string

	// Build creates a new TableStore instance with the specified tables. After this method is called,
	// the TableStoreBuilder should not be used again.
	Build() (TableStore, error)

	// Shutdown shuts down the store, flushing any remaining data to disk.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	Destroy() error
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
	NewBatch() Batch[TableKey]

	// Shutdown shuts down the store, flushing any remaining data to disk.
	Shutdown() error

	// Destroy shuts down and permanently deletes all data in the store.
	Destroy() error
}
