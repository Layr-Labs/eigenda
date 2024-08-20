package leveldb

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ kvstore.Store = &Store{}

// Store implements kvstore.Store interfaces with levelDB as the backend engine.
type Store struct { // TODO carefully consider naming
	database *leveldb.DB // TODO rename
	path     string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewStore returns a new Store built using LevelDB.
func NewStore(logger logging.Logger, path string) (*Store, error) { // TODO return type
	levelDB, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, err
	}

	return &Store{
		database: levelDB,
		logger:   logger,
	}, nil
}

func (store *Store) Put(key []byte, value []byte) error {
	return store.database.Put(key, value, nil)
}

func (store *Store) Get(key []byte) ([]byte, error) {
	data, err := store.database.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, kvstore.ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (store *Store) NewIterator(prefix []byte) iterator.Iterator {
	return store.database.NewIterator(util.BytesPrefix(prefix), nil)
}

func (store *Store) Delete(key []byte) error {
	err := store.database.Delete(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return kvstore.ErrNotFound
		}
		return nil
	}
	return nil
}

func (store *Store) DeleteBatch(keys [][]byte) error {
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(key)
	}
	return store.database.Write(batch, nil)
}

func (store *Store) WriteBatch(keys, values [][]byte) error {
	batch := new(leveldb.Batch)
	for i, key := range keys {
		batch.Put(key, values[i])
	}
	return store.database.Write(batch, nil)
}

// Shutdown shuts down the store.
func (store *Store) Shutdown() error {
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
func (store *Store) Destroy() error {
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
