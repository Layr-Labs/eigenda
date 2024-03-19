package node

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// DB is an interface to access the local database, such as leveldb, rocksdb.
type DB interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	DeleteBatch(keys [][]byte) error
	WriteBatch(keys, values [][]byte) error
	NewIterator(prefix []byte) iterator.Iterator
}
