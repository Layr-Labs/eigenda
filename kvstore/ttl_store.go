package kvstore

import "time"

// TTLStore is a store that supports key-value pairs with time-to-live (TTL) or expiration time.
type TTLStore interface {
	Store

	// PutWithTTL adds a key-value pair to the store that expires after a specified duration.
	// Key is eventually deleted after the TTL elapses.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithTTL(key []byte, value []byte, ttl time.Duration) error

	// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
	// Key is eventually deleted after the expiry time.
	//
	// Warning: updating the value of a key with a ttl/expiration has undefined behavior. Support for this pattern
	// may be implemented in the future if a use case is identified.
	PutWithExpiration(key []byte, value []byte, expiryTime time.Time) error
}
