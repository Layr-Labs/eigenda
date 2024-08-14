package kvstore

import "time"

var _ KVStore = &GopherKVStore{}

// GopherKVStore implements KVStore using GopherDB.
type GopherKVStore struct {
}

func (store *GopherKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (store *GopherKVStore) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (store *GopherKVStore) Drop(key []byte) error {
	//TODO implement me
	panic("implement me")
}

func (store *GopherKVStore) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (store *GopherKVStore) Destroy() error {
	//TODO implement me
	panic("implement me")
}
