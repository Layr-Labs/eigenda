package littdb

// LittDB is a highly specialized key-value store. It is intentionally very feature poor, sacrificing
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
type LittDB interface {

	// Put stores a value in the database. May not be used to overwrite an existing value.
	// Note that when this method returns, data written may not be crash durable on disk
	// (although the write does have atomicity). In order to ensure crash durability, call Flush().
	Put(table string, key []byte, value []byte) error

	// Get retrieves a value from the database. Returns an error if the value does not exist.
	Get(table string, key []byte) ([]byte, error)

	// Flush ensures that all data written to the database is crash durable on disk. When this method returns,
	// all data written by Put() operations is guaranteed to be crash durable. Put() operations called synchronously
	// with this method may not be crash durable after this method returns.
	Flush() error

	// Start starts the database. This method must be called before any other method is called.
	Start() error

	// Stop stops the database. This method must be called when the database is no longer needed.
	// Stop ensures that all non-flushed data is crash durable on disk before returning. Calls to
	// Put() concurrent with Stop() may not be crash durable after Stop() returns.¬
	Stop() error
}
