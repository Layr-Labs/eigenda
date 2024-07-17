package store

import "context"


type Stats struct {
	Entries int
	Reads   int
}

type Store interface {
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, value []byte) (key []byte, err error)
	// Stats returns the current usage metrics of the key-value data store.
	Stats() *Stats
}
