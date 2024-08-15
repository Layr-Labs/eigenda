package kvstore

//
//import (
//	"errors"
//	"fmt"
//	"github.com/Layr-Labs/eigensdk-go/logging"
//	"github.com/cockroachdb/pebble"
//	"os"
//	"time"
//)
//
//var _ KVStore = &PebbleStore{}
//
//type PebbleStore struct {
//	db   *pebble.DB
//	path string
//
//	logger logging.Logger
//
//	shutdown  bool
//	destroyed bool
//}
//
//func NewPebbleStore(logger logging.Logger, path string) (KVStore, error) {
//	db, err := pebble.Open(path, &pebble.Options{
//		DisableWAL: true,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return &PebbleStore{
//		db:     db,
//		path:   path,
//		logger: logger,
//	}, nil
//}
//
//func (store *PebbleStore) Put(key []byte, value []byte, ttl time.Duration) error {
//	if store.shutdown {
//		return fmt.Errorf("store is offline")
//	}
//
//	return store.db.Set(key, value, &pebble.WriteOptions{
//		Sync: false,
//	})
//}
//
//func (store *PebbleStore) Get(key []byte) ([]byte, error) {
//	if store.shutdown {
//		return nil, fmt.Errorf("store is offline")
//	}
//
//	data, closer, err := store.db.Get(key)
//
//	if errors.Is(err, pebble.ErrNotFound) {
//		return nil, nil
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	dataCopy := make([]byte, len(data))
//	copy(dataCopy, data)
//
//	err = closer.Close()
//	if err != nil {
//		return nil, err
//	}
//
//	return dataCopy, nil
//}
//
//func (store *PebbleStore) Drop(key []byte) error {
//	if store.shutdown {
//		return fmt.Errorf("store is offline")
//	}
//
//	return store.db.Delete(key, &pebble.WriteOptions{
//		Sync: false,
//	})
//}
//
//func (store *PebbleStore) BatchUpdate(operations []*BatchOperation) error {
//	if store.shutdown {
//		return fmt.Errorf("store is offline")
//	}
//
//	batch := store.db.NewIndexedBatch()
//
//	for _, operation := range operations {
//		if operation.Value == nil {
//			err := batch.Delete(operation.Key, nil)
//			if err != nil {
//				return err
//			}
//		} else {
//			err := batch.Set(operation.Key, operation.Value, nil)
//			if err != nil {
//				return err
//			}
//		}
//
//	}
//	return store.db.Apply(batch, &pebble.WriteOptions{
//		Sync: false,
//	})
//}
//
//func (store *PebbleStore) Shutdown() error {
//	if store.shutdown {
//		return nil
//	}
//	store.shutdown = true
//
//	return store.db.Close()
//}
//
//func (store *PebbleStore) Destroy() error {
//	if store.destroyed {
//		return nil
//	}
//
//	if !store.shutdown {
//		err := store.Shutdown()
//		if err != nil {
//			return err
//		}
//	}
//
//	store.logger.Info(fmt.Sprintf("destroying Pebble store at path: %s", store.path))
//	err := os.RemoveAll(store.path)
//	if err != nil {
//		return err
//	}
//	store.destroyed = true
//	return nil
//}
//
//func (store *PebbleStore) IsShutDown() bool {
//	return store.shutdown
//}
