package kvstore

// TableKey is a key that can be used for batch operations that modify a table.
type TableKey []byte

// Table can be used to operate on data in a specific table in a TableStore.
type Table interface {
	Store

	// Name returns the name of the table.
	Name() string

	// TableKey creates a new key that can be used for batch operations that modify this table.
	TableKey(key []byte) TableKey
}

// Batch is a collection of operations that can be applied atomically to a TableStore.
type Batch interface {
	// Put stores the given key / value pair in the batch, overwriting any existing value for that key.
	Put(key TableKey, value []byte)
	// Delete removes the key from the batch. Does not return an error if the key does not exist.
	Delete(key TableKey)
	// Apply atomically writes all the key / value pairs in the batch to the database.
	Apply() error
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
	//
	// WARNING: this method is not thread safe with respect to any other methods in this interface or
	// any methods on any Table objects associated with this store.
	GetTable(name string) (Table, error)

	// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
	//
	// WARNING: this method is not thread safe with respect to any other methods in this interface or
	// any methods on any Table objects associated with this store.
	DropTable(name string) error

	// GetMaxTableCount returns the maximum number of tables that can be created in the store
	// (excluding internal tables utilized by the store).
	GetMaxTableCount() uint32

	// GetTableCount returns the current number of tables in the store
	// (excluding internal tables utilized by the store).
	GetTableCount() uint32

	// GetTables returns a list of all tables in the store (excluding internal tables utilized by the store).
	GetTables() []Table

	// NewBatch creates a new batch that can be used to perform multiple operations atomically.
	NewBatch() Batch

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
