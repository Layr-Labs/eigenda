package kvstore

// Batch is a collection of operations that can be applied atomically to a TableStore.
type Batch[T any] interface {
	// Put stores the given key / value pair in the batch, overwriting any existing value for that key.
	Put(key T, value []byte)
	// Delete removes the key from the batch. Does not return an error if the key does not exist.
	Delete(key T)
	// Apply atomically writes all the key / value pairs in the batch to the database.
	Apply() error
}

// BatchOperator is an interface for creating new batches.
type BatchOperator[T any] interface {
	// NewBatch creates a new batch that can be used to perform multiple operations atomically.
	NewBatch() Batch[T]
}
