package kvstore

import "time"

var _ KVStore = &LevelKVStore{}

// LevelKVStore implements KVStore using LevelDB.
type LevelKVStore struct {
}

func (store *LevelKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (store *LevelKVStore) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
