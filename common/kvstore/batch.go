package kvstore

// Batch is a collection of operations that can be applied atomically to a TableStore.
type Batch[T any] interface {
	// TODO test nil value case
	// Put stores the given key / value pair in the batch, overwriting any existing value for that key.
	// If nil is passed as the value, a byte slice of length 0 will be stored.
	Put(key T, value []byte)
	// Delete removes the key from the batch.
	Delete(key T)
	// Apply atomically writes all the key / value pairs in the batch to the database.
	Apply() error
	// Size returns the number of operations in the batch.
	Size() uint32
}

// BatchOperator is an interface for creating new batches.
type BatchOperator[T any] interface {
	// NewBatch creates a new batch that can be used to perform multiple operations atomically.
	NewBatch() Batch[T]
}
