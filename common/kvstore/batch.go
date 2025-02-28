package kvstore

import "time"

// Batch is a collection of key / value pairs that will be written atomically to a database.
// Although it is thread safe to modify different batches in parallel or to modify a batch while
// the store is being modified, it is not thread safe to concurrently modify the same batch.
type Batch[K any] interface {
	// Put stores the given key / value pair in the batch, overwriting any existing value for that key.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	Put(key K, value []byte)
	// Delete removes the key from the batch.
	Delete(key K)
	// Apply atomically writes all the key / value pairs in the batch to the database.
	Apply() error
	// Size returns the number of operations in the batch.
	Size() uint32
}

// TTLBatch is a collection of key / value pairs that will be written atomically to a database with
// time-to-live (TTL) or expiration times. Although it is thread safe to modify different batches in
// parallel or to modify a batch while the store is being modified, it is not thread safe to concurrently
// modify the same batch.
type TTLBatch[K any] interface {
	Batch[K]
	// PutWithTTL stores the given key / value pair in the batch with a time-to-live (TTL) or expiration time.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	PutWithTTL(key K, value []byte, ttl time.Duration)
	// PutWithExpiration stores the given key / value pair in the batch with an expiration time.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	PutWithExpiration(key K, value []byte, expiryTime time.Time)
}
