package leveldb

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var ErrNotFound = errors.New("not found")

// This is an implementation of node.DB interfaces with levelDB as the backend engine.
type LevelDBStore struct {
	*leveldb.DB
}

func NewLevelDBStore(path string) (*LevelDBStore, error) {
	handle, err := leveldb.OpenFile(path, nil)
	return &LevelDBStore{handle}, err
}

func (d *LevelDBStore) Put(key []byte, value []byte) error {
	return d.DB.Put(key, value, nil)
}

func (d *LevelDBStore) Get(key []byte) ([]byte, error) {
	data, err := d.DB.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (d *LevelDBStore) NewIterator(prefix []byte) iterator.Iterator {
	return d.DB.NewIterator(util.BytesPrefix(prefix), nil)
}

func (d *LevelDBStore) Delete(key []byte) error {
	return d.DB.Delete(key, nil)
}

func (d *LevelDBStore) DeleteBatch(keys [][]byte) error {
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(key)
	}
	return d.DB.Write(batch, nil)
}

func (d *LevelDBStore) WriteBatch(keys, values [][]byte) error {
	batch := new(leveldb.Batch)
	for i, key := range keys {
		batch.Put(key, values[i])
	}
	return d.DB.Write(batch, nil)
}
