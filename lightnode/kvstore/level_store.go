package kvstore

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"time"
)

var _ KVStore = &LevelStore{}

// TODO add timeouts for operations maybe, see node/store.go

// LevelStore implements KVStore using LevelDB.
type LevelStore struct {
	db   *leveldb.DB
	path string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewLevelStore creates a new LevelStore.
func NewLevelStore(logger logging.Logger, path string) (KVStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelStore{
		path:   path,
		db:     db,
		logger: logger,
	}, nil
}

// All "regular" keys in the database are prefixed with this string. This is an internal implementation detail, and
// external users of this store do not need to take this into consideration. This prefix is used to help implement
// the TTL feature.
const keyPrefix string = "k"

// All keys used to describe expiration times are prefixed with this string. This is an internal implementation detail,
// and external users of this store do not need to take this into consideration. This prefix is used to help implement
// the TTL feature.
const expirationPrefix string = "e"

// Put stores a data in the store.
func (store *LevelStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	if ttl == 0 {
		return store.db.Put(key, value, nil)
	} else {
		// TODO
		return nil
	}
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *LevelStore) Get(key []byte) ([]byte, error) {
	if store.shutdown {
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
func (store *LevelStore) Drop(key []byte) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	return store.db.Delete(key, nil)
}

// BatchUpdate performs a batch of Put and Drop operations.
func (store *LevelStore) BatchUpdate(operations []*BatchOperation) error {
	// TODO implement a real batch update

	for _, operation := range operations {
		if operation.Value == nil {
			err := store.Drop(operation.Key)
			if err != nil {
				return err
			}
		} else {
			err := store.Put(operation.Key, operation.Value, operation.TTL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Shutdown shuts down the store.
func (store *LevelStore) Shutdown() error {
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
func (store *LevelStore) Destroy() error {
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

// IsShutDown returns true if the store is shut down.
func (store *LevelStore) IsShutDown() bool {
	return store.shutdown
}
