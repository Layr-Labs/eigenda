package leveldb

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/db"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ db.DB = &LevelDBStore{}

// LevelDBStore implements db.DB interfaces with levelDB as the backend engine.
type LevelDBStore struct { // TODO carefully consider naming
	database *leveldb.DB // TODO rename
	path     string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewLevelDBStore returns a new DB built using LevelDB.
func NewLevelDBStore(logger logging.Logger, path string) (*LevelDBStore, error) { // TODO return type
	levelDB, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		database: levelDB,
		logger:   logger,
	}, nil
}

func (store *LevelDBStore) Put(key []byte, value []byte) error {
	return store.database.Put(key, value, nil)
}

func (store *LevelDBStore) Get(key []byte) ([]byte, error) {
	data, err := store.database.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, db.ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (store *LevelDBStore) NewIterator(prefix []byte) iterator.Iterator {
	return store.database.NewIterator(util.BytesPrefix(prefix), nil)
}

func (store *LevelDBStore) Delete(key []byte) error {
	err := store.database.Delete(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return db.ErrNotFound
		}
		return nil
	}
	return nil
}

func (store *LevelDBStore) DeleteBatch(keys [][]byte) error {
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(key)
	}
	return store.database.Write(batch, nil)
}

func (store *LevelDBStore) WriteBatch(keys, values [][]byte) error {
	batch := new(leveldb.Batch)
	for i, key := range keys {
		batch.Put(key, values[i])
	}
	return store.database.Write(batch, nil)
}

// Shutdown shuts down the store.
func (store *LevelDBStore) Shutdown() error {
	if store.shutdown {
		return nil
	}

	err := store.database.Close()
	if err != nil {
		return err
	}

	store.shutdown = true
	return nil
}

// Destroy destroys the store.
func (store *LevelDBStore) Destroy() error {
	if store.destroyed {
		return nil
	}

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
	store.destroyed = true
	return nil
}
