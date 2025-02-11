package litt

import "time"

// Table is a key-value store with a namespace that does not overlap with other tables.
// Values may be written to the table, but once written, they may not be changed or deleted (except via TTL).
type Table interface {
	// Name returns the name of the table.
	Name() string

	// Put stores a value in the database. May not be used to overwrite an existing value.
	// Note that when this method returns, data written may not be crash durable on disk
	// (although the write does have atomicity). In order to ensure crash durability, call Flush().
	Put(key []byte, value []byte) error

	// Get retrieves a value from the database. Returns an error if the value does not exist.
	Get(key []byte) ([]byte, error)

	// Flush ensures that all data written to the database is crash durable on disk. When this method returns,
	// all data written by Put() operations is guaranteed to be crash durable. Put() operations called synchronously
	// with this method may not be crash durable after this method returns.
	//
	// Note that data flushed at the same time is not atomic. If the process crashes mid-flush, some data
	// being flushed may become persistent, while some may not. Each individual key-value pair is atomic
	// in the event of a crash, though.
	Flush() error

	// SetTTL sets the time to live for data in this table. This TTL is immediately applied to data already in
	// the table. Note that deletion is lazy. That is, when the data expires, it may not be deleted immediately.
	SetTTL(ttl time.Duration)
}
