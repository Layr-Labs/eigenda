package store

import (
	"context"

	"github.com/Layr-Labs/eigenda-proxy/common"
)

type Store interface {
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte, domain common.DomainType) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, value []byte) (key []byte, err error)
	// Entries returns the number of values in the store; only used for memstore
	Stats() *common.Stats
}
