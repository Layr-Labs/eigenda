package storeutil

import (
	"github.com/Layr-Labs/eigenda/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"time"
)

// TTLStore adds a time-to-live (TTL) capability to the store.
type TTLStore struct {
	store kvstore.Store
}

// TTLWrapper extends the given store with TTL capabilities.
func TTLWrapper(store kvstore.Store) *TTLStore {
	return &TTLStore{
		store: store,
	}
}

// PutWithTTL adds a key-value pair to the store with a TTL. Key is eventually deleted after the time to live expires.`
func (store *TTLStore) PutWithTTL(key []byte, value []byte, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) Put(key []byte, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) Delete(key []byte) error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) DeleteBatch(keys [][]byte) error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) WriteBatch(keys, values [][]byte) error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) NewIterator(prefix []byte) iterator.Iterator {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (store *TTLStore) Destroy() error {
	//TODO implement me
	panic("implement me")
}
