package benchmark

// Database encapsulates a database (e.g. LittDB). The goal of this layer of indirection is to support
// benchmarking different database solutions.
type Database interface {

	// Write a value to the DB. Value does not have to be persistent when this method returns.
	// Value must be visible to future Read() calls when this method returns.
	Write(key []byte, value []byte) error

	// Flush all previously written values out to disk. Values must be persistent when this method returns.
	// For DB implementations that don't support a flush pattern, this wrapper may want to accumulate values
	// in batches when Write() is called, and commit those batches when Flush() is called.
	Flush() error

	// Read returns a value previously written to the DB. Does not return an error if the value is not present.
	Read(key []byte) (value []byte, exists bool, err error)

	// Close cleanly shuts down the DB.
	Close() error
}
