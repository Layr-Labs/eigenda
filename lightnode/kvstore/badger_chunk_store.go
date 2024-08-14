package kvstore

// go get github.com/dgraph-io/badger/v4

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	badger "github.com/dgraph-io/badger/v4"
	"os"
	"time"
)

var _ KVStore = &BadgerKVStore{}

// TODO we need to periodically tell levelDB to do background merges

// BadgerKVStore implements KVStore using GopherDB.
type BadgerKVStore struct {
	db   *badger.DB
	path string

	logger logging.Logger

	shutdown  bool
	destroyed bool
}

// NewBadgerKVStore creates a new BadgerKVStore.
func NewBadgerKVStore(logger logging.Logger, path string) (*BadgerKVStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	return &BadgerKVStore{
		db:     db,
		path:   path,
		logger: logger,
	}, nil
}

// Put stores a data in the store.
func (store *BadgerKVStore) Put(key []byte, value []byte, ttl time.Duration) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	// TODO a new transaction for each put may not be the most efficient way to do this

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
func (store *BadgerKVStore) Get(key []byte) ([]byte, error) {
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
func (store *BadgerKVStore) Drop(key []byte) error {
	if store.shutdown {
		return fmt.Errorf("store is offline")
	}

	return store.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Shutdown shuts down the store.
func (store *BadgerKVStore) Shutdown() error {
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
func (store *BadgerKVStore) Destroy() error {
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
