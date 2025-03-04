package litt

// DB is a highly specialized key-value store. It is intentionally very feature poor, sacrificing
// unnecessary features for simplicity, high performance, and low memory usage.
//
// Litt: slang, a synonym for "cool" or "awesome". e.g. "Man, that database is litt, bro!".
//
// Supported features:
// - writing values
// - reading values
// - TTLs and automatic (lazy) deletion of expired values
// - tables with non-overlapping namespaces
// - thread safety: all methods are safe to call concurrently, and all modifications are atomic
//
// Unsupported features:
// - mutating existing values (once a value is written, it cannot be changed)
// - deleting values (values only leave the DB when they expire via a TTL)
// - transactions (individual operations are atomic, but there is no way to group operations atomically)
// - fine granularity for TTL (all data in the same table must have the same TTL)
type DB interface {
	// GetTable gets a table by name, creating one if it does not exist.
	//
	// The first time a table is fetched (either a new table or an existing one loaded from disk), its TTL is always
	// set to 0 (i.e. it has no TTL). If you want to set a TTL, you must call Table.SetTTL() to do so. This is
	// necessary after each time the database is started/restarted.
	GetTable(name string) (Table, error)

	// DropTable deletes a table and all of its data.
	//
	// Note that it is NOT thread safe to drop a table concurrently with any operation that accesses the table.
	// The table returned by GetTable() before DropTable() is called must not be used once DropTable() is called.
	DropTable(name string) error

	// Start starts the database. This method must be called before any other method is called.
	Start() error

	// Stop stops the database. This method must be called when the database is no longer needed.
	// Stop ensures that all non-flushed data is crash durable on disk before returning. Calls to
	// Put() concurrent with Stop() may not be crash durable after Stop() returns.¬
	Stop() error

	// Destroy deletes all data in the database.
	Destroy() error
}
