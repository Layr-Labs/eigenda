package kvstore

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"time"
)

var _ KVStore = &LevelKVStore{}

// LevelKVStore implements KVStore using LevelDB.
type LevelKVStore struct {
	db   *leveldb.DB
	path string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewLevelKVStore creates a new LevelKVStore.
func NewLevelKVStore(logger logging.Logger, path string) (*LevelKVStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelKVStore{
		path:   path,
		db:     db,
		logger: logger,
	}, nil
}

// Put stores a data in the store.
func (store *LevelKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.shutdown || store.destroyed {
		return fmt.Errorf("store is offline")
	}
	// TODO improve performance by buffering and writing in larger batches
	// TODO could this pattern be encapsulated in a helper class?

	return store.db.Put(key, value, nil)
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *LevelKVStore) Get(key []byte) ([]byte, error) {
	if store.shutdown || store.destroyed {
		return nil, fmt.Errorf("store is offline")
	}

	value, err := store.db.Get(key, nil)

	if err == nil {
		return value, nil
	}

	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, nil
	}

	return nil, err

}

// Drop deletes data from the store.
func (store *LevelKVStore) Drop(key []byte) error {
	if store.shutdown || store.destroyed {
		return fmt.Errorf("store is offline")
	}

	return store.db.Delete(key, nil)
}

// Shutdown shuts down the store.
func (store *LevelKVStore) Shutdown() error {
	if store.shutdown {
		return nil
	}

	err := store.db.Close()
	if err != nil {
		return err
	}

	store.shutdown = true
	return nil
}

// Destroy destroys the store.
func (store *LevelKVStore) Destroy() error {
	if store.destroyed {
		return nil
	}

	if !store.shutdown {
		err := store.Shutdown()
		if err != nil {
			return err
		}
	}

	store.logger.Info("destroying LevelDB store at path: %s", store.path)
	err := os.RemoveAll(store.path)
	if err != nil {
		return err
	}
	store.destroyed = true
	return nil
}
