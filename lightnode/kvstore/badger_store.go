package kvstore

// go get github.com/dgraph-io/badger/v4

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/dgraph-io/badger/v4"
	boptions "github.com/dgraph-io/badger/v4/options"
	"os"
	"time"
)

var _ KVStore = &BadgerStore{}

// TODO we need to periodically tell levelDB to do background merges

// BadgerStore implements KVStore using GopherDB.
type BadgerStore struct {
	db   *badger.DB
	path string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewBadgerStore creates a new BadgerStore.
func NewBadgerStore(logger logging.Logger, path string) (KVStore, error) {

	options := badger.DefaultOptions(path)
	options.Compression = boptions.None

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &BadgerStore{
		db:     db,
		path:   path,
		logger: logger,
	}, nil
}

// Put stores a data in the store.
func (store *BadgerStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	return store.db.Update(func(txn *badger.Txn) error {
		if ttl == 0 {
			return txn.Set(key, value)
		} else {
			entry := badger.NewEntry(key, value).WithTTL(ttl) // TODO test ttl
			return txn.SetEntry(entry)
		}
	})
}

// Get retrieves data from the store. Returns nil if the data is not found.
func (store *BadgerStore) Get(key []byte) ([]byte, error) {
	var value []byte

	err := store.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})

	if err == nil {
		return value, nil
	}

	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	return nil, err
}

// Drop deletes data from the store.
func (store *BadgerStore) Drop(key []byte) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	return store.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// BatchUpdate performs a batch of Put and Drop operations.
func (store *BadgerStore) BatchUpdate(operations []*BatchOperation) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	return store.db.Update(func(txn *badger.Txn) error {
		for _, operation := range operations {
			if operation.Value == nil {
				err := txn.Delete(operation.Key)
				if err != nil {
					return err
				}
			} else if operation.TTL == 0 {
				err := txn.Set(operation.Key, operation.Value)
				if err != nil {
					return err
				}
			} else {
				entry := badger.NewEntry(operation.Key, operation.Value).WithTTL(operation.TTL)
				err := txn.SetEntry(entry)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Shutdown shuts down the store.
func (store *BadgerStore) Shutdown() error {
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
func (store *BadgerStore) Destroy() error {
	if store.destroyed {
		return nil
	}

	if !store.shutdown {
		err := store.Shutdown()
		if err != nil {
			return err
		}
	}

	store.logger.Info(fmt.Sprintf("destroying BadgerDB store at path: %s", store.path))
	err := os.RemoveAll(store.path)
	if err != nil {
		return err
	}
	store.destroyed = true
	return nil
}

// IsShutDown returns true if the store is shut down.
func (store *BadgerStore) IsShutDown() bool {
	return store.shutdown
}
