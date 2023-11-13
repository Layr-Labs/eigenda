package node

import (
	"encoding/binary"

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

// ToByteArray converts an uint64 into byte array in big endian.
func ToByteArray(i uint64) []byte {
	arr := make([]byte, 8)
	binary.BigEndian.PutUint64(arr[0:8], uint64(i))
	return arr
}

// ToUint64 converts a byte array into an uint64, assuming big endian.
func ToUint64(arr []byte) uint64 {
	i := binary.BigEndian.Uint64(arr)
	return i
}
