package test

import (
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/db"
	"github.com/Layr-Labs/eigenda/db/leveldb"
	"github.com/Layr-Labs/eigenda/db/memdb"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

func randomOperationsTest(t *testing.T, store db.DB) {
	tu.InitializeRandom()

	// Delete the database directory, just in case it was left over from a previous run.
	err := os.RemoveAll(dbPath)
	assert.NoError(t, err)

	expectedData := make(map[string][]byte)

	for i := 0; i < 1000; i++ {

		choice := rand.Float64()
		if len(expectedData) == 0 || choice < 0.50 {
			// Write a random value.

			key := tu.RandomBytes(32)
			value := tu.RandomBytes(32)

			err = store.Put(key, value)
			assert.NoError(t, err)

			expectedData[string(key)] = value
		} else if choice < 0.75 {
			// Modify a random value.

			var key string
			for k := range expectedData {
				key = k
				break
			}
			value := tu.RandomBytes(32)
			err = store.Put([]byte(key), value)
			assert.NoError(t, err)
			expectedData[key] = value
		} else if choice < 0.90 {
			// Drop a random value.

			var key string
			for k := range expectedData {
				key = k
				break
			}
			delete(expectedData, key)
			err = store.Delete([]byte(key))
			assert.NoError(t, err)
		} else {
			// Drop a non-existent value.

			key := tu.RandomBytes(32)
			err = store.Delete(key)
			assert.Nil(t, err)
		}

		if i%10 == 0 {
			// Every so often, check that the store matches the expected data.
			for key, expectedValue := range expectedData {
				value, err := store.Get([]byte(key))
				assert.NoError(t, err)
				assert.Equal(t, expectedValue, value)
			}

			// Try and get a value that isn't in the store.
			key := tu.RandomBytes(32)
			value, err := store.Get(key)
			assert.Equal(t, db.ErrNotFound, err)
			assert.Nil(t, value)
		}
	}

	err = store.Shutdown()
	assert.NoError(t, err)
	err = store.Destroy()
	assert.NoError(t, err)
}

var dbPath = "test-store"

func verifyDBIsDeleted(t *testing.T) {
	_, err := os.Stat(dbPath)
	assert.True(t, os.IsNotExist(err))
}

func TestRandomOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)
	var store db.DB

	// In memory store

	randomOperationsTest(t, memdb.NewInMemoryStore())
	//randomOperationsTest(t, ThreadSafeWrapper(NewInMemoryStore()))
	//randomOperationsTest(t, BatchingWrapper(NewInMemoryStore(), 32*5))

	// LevelDB store

	store, err = leveldb.NewLevelDBStore(logger, dbPath)
	assert.NoError(t, err)
	randomOperationsTest(t, store)
	verifyDBIsDeleted(t)

	//store, err = NewLevelStore(logger, dbPath)
	//store = ThreadSafeWrapper(store)
	//assert.NoError(t, err)
	//randomOperationsTest(t, store)
	//verifyDBIsDeleted(t)
	//
	//store, err = NewLevelStore(logger, dbPath)
	//store = BatchingWrapper(store, 32*5)
	//assert.NoError(t, err)
	//randomOperationsTest(t, store)
	//verifyDBIsDeleted(t)
}

//func batchOperationsTest(t *testing.T, store KVStore) {
//	tu.InitializeRandom()
//
//	var err error
//
//	expectedData := make(map[string][]byte)
//
//	var operations []*BatchOperation
//
//	for i := 0; i < 11; i++ { // TODO 1000
//
//		choice := rand.Float64()
//		if len(expectedData) == 0 || choice < 0.66 {
//			// Write a random value.
//
//			key := tu.RandomBytes(32)
//			value := tu.RandomBytes(32)
//
//			operations = append(operations, &BatchOperation{
//				Key:   key,
//				Value: value,
//				TTL:   0,
//			})
//
//			expectedData[string(key)] = value
//		} else if choice < 0.90 {
//			// Drop a random value.
//
//			var key string
//			for k := range expectedData {
//				key = k
//			}
//			delete(expectedData, key)
//
//			operations = append(operations, &BatchOperation{
//				Key: []byte(key),
//			})
//		} else {
//			// Drop a non-existent value.
//
//			key := tu.RandomBytes(32)
//			operations = append(operations, &BatchOperation{
//				Key: key,
//			})
//		}
//
//		if i%10 == 0 {
//			// Every so often, apply the batch and check that the store matches the expected data.
//
//			err := store.BatchUpdate(operations)
//			assert.NoError(t, err)
//
//			operations = nil
//
//			for key, expectedValue := range expectedData {
//				value, err := store.Get([]byte(key))
//				assert.NoError(t, err)
//				assert.Equal(t, expectedValue, value)
//			}
//
//			// Try and get a value that isn't in the store.
//			key := tu.RandomBytes(32)
//			value, err := store.Get(key)
//			assert.NoError(t, err)
//			assert.Nil(t, value)
//		}
//	}
//
//	err = store.Shutdown()
//	assert.NoError(t, err)
//	err = store.Destroy()
//	assert.NoError(t, err)
//}
//
//func TestBatchOperations(t *testing.T) {
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//	var store KVStore
//
//	// In memory store
//
//	batchOperationsTest(t, NewInMemoryStore())
//	batchOperationsTest(t, ThreadSafeWrapper(NewInMemoryStore()))
//	batchOperationsTest(t, BatchingWrapper(NewInMemoryStore(), 32*5))
//
//	// LevelDB store
//
//	store, err = NewLevelStore(logger, dbPath)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewLevelStore(logger, dbPath)
//	store = ThreadSafeWrapper(store)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewLevelStore(logger, dbPath)
//	store = BatchingWrapper(store, 32*5)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//
//	// BadgerDB store
//
//	//store, err = NewBadgerStore(logger, dbPath)
//	//assert.NoError(t, err)
//	//batchOperationsTest(t, store)
//	//verifyDBIsDeleted(t)
//	//
//	//store, err = NewBadgerStore(logger, dbPath)
//	//store = ThreadSafeWrapper(store)
//	//assert.NoError(t, err)
//	//batchOperationsTest(t, store)
//	//verifyDBIsDeleted(t)
//	//
//	//store, err = NewBadgerStore(logger, dbPath)
//	//store = BatchingWrapper(store, 32*5)
//	//assert.NoError(t, err)
//	//batchOperationsTest(t, store)
//	//verifyDBIsDeleted(t)
//
//	// Pebble store
//
//	store, err = NewPebbleStore(logger, dbPath)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewPebbleStore(logger, dbPath)
//	store = ThreadSafeWrapper(store)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewPebbleStore(logger, dbPath)
//	store = BatchingWrapper(store, 32*5)
//	assert.NoError(t, err)
//	batchOperationsTest(t, store)
//	verifyDBIsDeleted(t)
//}
//
//func operationsOnShutdownStoreTest(t *testing.T, store KVStore) {
//	err := store.Shutdown()
//	assert.NoError(t, err)
//
//	err = store.Put([]byte("key"), []byte("value"), 0)
//	assert.Error(t, err)
//
//	_, err = store.Get([]byte("key"))
//	assert.Error(t, err)
//
//	err = store.Drop([]byte("key"))
//	assert.Error(t, err)
//
//	err = store.Shutdown()
//	assert.NoError(t, err)
//
//	err = store.Destroy()
//	assert.NoError(t, err)
//}
//
//func TestOperationsOnShutdownStore(t *testing.T) {
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//	var store KVStore
//
//	// In memory store
//
//	operationsOnShutdownStoreTest(t, NewInMemoryStore())
//	operationsOnShutdownStoreTest(t, ThreadSafeWrapper(NewInMemoryStore()))
//	operationsOnShutdownStoreTest(t, BatchingWrapper(NewInMemoryStore(), 32*5))
//
//	// LevelDB store
//
//	store, err = NewLevelStore(logger, dbPath)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewLevelStore(logger, dbPath)
//	store = ThreadSafeWrapper(store)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewLevelStore(logger, dbPath)
//	store = BatchingWrapper(store, 32*5)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//
//	// BadgerDB store
//	//store, err = NewBadgerStore(logger, dbPath)
//	//assert.NoError(t, err)
//	//operationsOnShutdownStoreTest(t, store)
//	//verifyDBIsDeleted(t)
//	//
//	//store, err = NewBadgerStore(logger, dbPath)
//	//store = ThreadSafeWrapper(store)
//	//assert.NoError(t, err)
//	//operationsOnShutdownStoreTest(t, store)
//	//verifyDBIsDeleted(t)
//	//
//	//store, err = NewBadgerStore(logger, dbPath)
//	//store = BatchingWrapper(store, 32*5)
//	//assert.NoError(t, err)
//	//operationsOnShutdownStoreTest(t, store)
//	//verifyDBIsDeleted(t)
//
//	// Pebble store
//
//	store, err = NewPebbleStore(logger, dbPath)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewPebbleStore(logger, dbPath)
//	store = ThreadSafeWrapper(store)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//
//	store, err = NewPebbleStore(logger, dbPath)
//	store = BatchingWrapper(store, 32*5)
//	assert.NoError(t, err)
//	operationsOnShutdownStoreTest(t, store)
//	verifyDBIsDeleted(t)
//}
