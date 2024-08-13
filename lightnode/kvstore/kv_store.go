package kvstore

import "time"

// KVStore is an interface for a key-value store. It is used by a light node to store data.
type KVStore interface {

	// Put stores data for later retrieval. The data will be available via Get until at least when
	// the TTL (time to live) expires.
	Put(key []byte, value []byte, ttl time.Duration) error

	// Get retrieves data that was previously stored with StoreData. Returns an error if the data
	// is unavailable for any reason.
	Get(key []byte) ([]byte, error)
}
