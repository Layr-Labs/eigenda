package wrappers

// A wrapper around a database, enables a database to be used by the benchmark engine.
type DatabaseWrapper interface {
	// BuildThreadLocalWrapper creates a thread-local wrapper for the database.
	BuildThreadLocalWrapper() (ThreadLocalDatabaseWrapper, error)
}

// ThreadLocalDatabaseWrapper is a wrapper around a database that is used on a single thread. Useful for DBs that
// want to do things like batching.
type ThreadLocalDatabaseWrapper interface {
	// Insert a key-value pair into the database.
	Put(key, value []byte) error

	// Get a value by its key from the database.
	Get(key []byte) (value []byte, exists bool, err error)

	// Flush data out to the database.
	Flush() error
}
