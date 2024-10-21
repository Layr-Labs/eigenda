package leveldb

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ kvstore.Store[[]byte] = &levelDBStore{}

// levelDBStore implements kvstore.Store interfaces with levelDB as the backend engine.
type levelDBStore struct {
	db   *leveldb.DB
	path string

	logger logging.Logger

	shutdown bool
}

// NewStore returns a new levelDBStore built using LevelDB.
func NewStore(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
	levelDB, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, err
	}

	return &levelDBStore{
		db:     levelDB,
		path:   path,
		logger: logger,
	}, nil
}

// Put stores a data in the store.
func (store *levelDBStore) Put(key []byte, value []byte) error {
	if value == nil {
		value = []byte{}
	}
	return store.db.Put(key, value, nil)
}

// Get retrieves data from the store. Returns kvstore.ErrNotFound if the data is not found.
func (store *levelDBStore) Get(key []byte) ([]byte, error) {
	data, err := store.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, kvstore.ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

// NewIterator creates a new iterator. Only keys prefixed with the given prefix will be iterated.
func (store *levelDBStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	return store.db.NewIterator(util.BytesPrefix(prefix), nil), nil
}

// Delete deletes data from the store.
func (store *levelDBStore) Delete(key []byte) error {
	return store.db.Delete(key, nil)
}

// DeleteBatch deletes multiple key-value pairs from the store.
func (store *levelDBStore) DeleteBatch(keys [][]byte) error {
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(key)
	}
	return store.db.Write(batch, nil)
}

// WriteBatch adds multiple key-value pairs to the store.
func (store *levelDBStore) WriteBatch(keys [][]byte, values [][]byte) error {
	batch := new(leveldb.Batch)
	for i, key := range keys {
		batch.Put(key, values[i])
	}
	return store.db.Write(batch, nil)
}

// NewBatch creates a new batch for the store.
func (store *levelDBStore) NewBatch() kvstore.Batch[[]byte] {
	return &levelDBBatch{
		store: store,
		batch: new(leveldb.Batch),
	}
}

type levelDBBatch struct {
	store *levelDBStore
	batch *leveldb.Batch
}

func (m *levelDBBatch) Put(key []byte, value []byte) {
	if value == nil {
		value = []byte{}
	}
	m.batch.Put(key, value)
}

func (m *levelDBBatch) Delete(key []byte) {
	m.batch.Delete(key)
}

func (m *levelDBBatch) Apply() error {
	return m.store.db.Write(m.batch, nil)
}

// Size returns the number of operations in the batch.
func (m *levelDBBatch) Size() uint32 {
	return uint32(m.batch.Len())
}

// Shutdown shuts down the store.
//
// Warning: it is not thread safe to call this method concurrently with other methods on this class,
// or while there exist unclosed iterators.
func (store *levelDBStore) Shutdown() error {
	err := store.db.Close()
	if err != nil {
		return err
	}

	store.shutdown = true
	return nil
}

// Destroy destroys the store.
//
// Warning: it is not thread safe to call this method concurrently with other methods on this class,
// or while there exist unclosed iterators.
func (store *levelDBStore) Destroy() error {
	if !store.shutdown {
		err := store.Shutdown()
		if err != nil {
			return err
		}
	}

	store.logger.Info(fmt.Sprintf("destroying LevelDB store at path: %s", store.path))
	err := os.RemoveAll(store.path)
	if err != nil {
		return err
	}
	return nil
}
