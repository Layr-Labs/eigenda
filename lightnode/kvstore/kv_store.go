package kvstore

import "time"

// TODO: write up the consistency/concurrency model for this interface.

// KVStore is an interface for a key-value store. It is used by a light node to store data.
type KVStore interface {

	// Put stores data for later retrieval. The data will be available via Get until at least when
	// the TTL (time to live) expires.
	Put(key []byte, value []byte, ttl time.Duration) error

	// Get retrieves data that was previously stored with StoreData. Returns an error if the data
	// is unavailable for any reason.
	Get(key []byte) ([]byte, error)

	// Drop removes data from the store. This is a no-op if the data does not exist.
	Drop(key []byte) error

	// Shutdown stops the store and releases any resources it holds. Does not delete any on-disk data.
	// Calling shutdown on a store that has already been shut down or destroyed is a no-op.
	Shutdown() error

	// Destroy permanently stops the store and deletes all data (including data on disk).
	// Calling destroy on a store that has already been destroyed is a no-op. It is not
	// necessary to call Shutdown before calling Destroy.
	Destroy() error
}

// Operations to potentially support in the future:
// - iteration over all keys
// - information about the size of the store (entry count, number of bytes, etc)
// - methods to support automated migration from any implementation to any other implementation
