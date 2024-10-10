package test

import (
	"context"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/common/kvstore/ttl"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

// A list of builders for various stores to be tested.
var storeBuilders = []func(logger logging.Logger, path string) (kvstore.Store, error){
	func(logger logging.Logger, path string) (kvstore.Store, error) {
		return mapstore.NewStore(), nil
	},

	func(logger logging.Logger, path string) (kvstore.Store, error) {
		return ttl.TTLWrapper(context.Background(), logger, mapstore.NewStore(), 0), nil
	},
	func(logger logging.Logger, path string) (kvstore.Store, error) {
		return leveldb.NewStore(logger, path)
	},
	func(logger logging.Logger, path string) (kvstore.Store, error) {
		store, err := leveldb.NewStore(logger, path)
		if err != nil {
			return nil, err
		}
		return ttl.TTLWrapper(context.Background(), logger, store, 0), nil
	},
	func(logger logging.Logger, path string) (kvstore.Store, error) {
		tableStore, err := tablestore.MapStore.Start(logger, path, "test")
		if err != nil {
			return nil, err
		}
		store, err := tableStore.GetTable("test")
		if err != nil {
			return nil, err
		}
		return store, nil
	},
	func(logger logging.Logger, path string) (kvstore.Store, error) {
		tableStore, err := tablestore.LevelDB.Start(logger, path, "test")
		if err != nil {
			return nil, err
		}
		store, err := tableStore.GetTable("test")
		if err != nil {
			return nil, err
		}
		return store, nil
	},
}

var dbPath = "test-store"

func deleteDBDirectory(t *testing.T) {
	err := os.RemoveAll(dbPath)
	assert.NoError(t, err)
}

func verifyDBIsDeleted(t *testing.T) {
	_, err := os.Stat(dbPath)
	assert.True(t, os.IsNotExist(err))
}

func randomOperationsTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	expectedData := make(map[string][]byte)

	for i := 0; i < 1000; i++ {

		choice := rand.Float64()
		if len(expectedData) == 0 || choice < 0.50 {
			// Write a random value.

			key := tu.RandomBytes(32)
			value := tu.RandomBytes(32)

			err := store.Put(key, value)
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
			err := store.Put([]byte(key), value)
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
			err := store.Delete([]byte(key))
			assert.NoError(t, err)
		} else {
			// Drop a non-existent value.

			key := tu.RandomBytes(32)
			err := store.Delete(key)
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
			assert.Equal(t, kvstore.ErrNotFound, err)
			assert.Nil(t, value)
		}
	}

	err := store.Shutdown()
	assert.NoError(t, err)
	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestRandomOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		randomOperationsTest(t, store)
	}
}

func writeBatchTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	var err error

	expectedData := make(map[string][]byte)
	batch := store.NewBatch()

	for i := 0; i < 1000; i++ {
		// Write a random value.
		key := tu.RandomBytes(32)

		var value []byte
		if i%50 == 0 {
			// nil values are interpreted as empty slices.
			value = nil
		} else {
			value = tu.RandomBytes(32)
		}

		batch.Put(key, value)

		if value == nil {
			expectedData[string(key)] = []byte{}
		} else {
			expectedData[string(key)] = value
		}

		if i%10 == 0 {
			// Every so often, apply the batch and check that the store matches the expected data.

			err = batch.Apply()
			assert.NoError(t, err)

			for key, expectedValue := range expectedData {
				value, err = store.Get([]byte(key))
				assert.NoError(t, err)
				assert.Equal(t, expectedValue, value)
			}

			// Try and get a value that isn't in the store.
			key = tu.RandomBytes(32)
			value, err = store.Get(key)
			assert.Equal(t, kvstore.ErrNotFound, err)
			assert.Nil(t, value)
		}
	}

	err = store.Shutdown()
	assert.NoError(t, err)
	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestWriteBatch(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		writeBatchTest(t, store)
	}
}

func deleteBatchTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	expectedData := make(map[string][]byte)

	batch := store.NewBatch()

	// Add some data to the store.
	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(32)
		value := tu.RandomBytes(32)

		err := store.Put(key, value)
		assert.NoError(t, err)

		expectedData[string(key)] = value
	}

	// Delete some of the data.
	for key := range expectedData {
		choice := rand.Float64()
		if choice < 0.5 {
			batch.Delete([]byte(key))
			delete(expectedData, key)
		} else if choice < 0.75 {
			// Delete a non-existent key.
			batch.Delete(tu.RandomBytes(32))
		}
	}

	err := batch.Apply()
	assert.NoError(t, err)

	// Check that the store matches the expected data.
	for key, expectedValue := range expectedData {
		value, err := store.Get([]byte(key))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Try and get a value that isn't in the store.
	key := tu.RandomBytes(32)
	value, err := store.Get(key)
	assert.Equal(t, kvstore.ErrNotFound, err)
	assert.Nil(t, value)

	err = store.Shutdown()
	assert.NoError(t, err)
	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func TestDeleteBatch(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		deleteBatchTest(t, store)
	}
}

func iterationTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	expectedData := make(map[string][]byte)

	// Insert some data into the store.
	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(32)
		value := tu.RandomBytes(32)

		err := store.Put(key, value)
		assert.NoError(t, err)

		expectedData[string(key)] = value
	}

	// Iterate over the store and check that the data matches the expected data.
	foundKeys := make(map[string]bool)

	iterator, err := store.NewIterator(nil)
	assert.NoError(t, err)
	defer iterator.Release()

	for iterator.Next() {
		key := string(iterator.Key())
		value := iterator.Value()

		expectedValue, ok := expectedData[key]
		assert.True(t, ok)
		assert.Equal(t, expectedValue, value)

		foundKeys[key] = true
	}
	assert.Equal(t, len(expectedData), len(foundKeys))

	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		iterationTest(t, store)
	}
}

func iterationWithPrefixTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	prefixA := tu.RandomBytes(8)
	prefixB := tu.RandomBytes(8)

	expectedDataA := make(map[string][]byte)
	expectedDataB := make(map[string][]byte)

	// Insert some data into the store.
	for i := 0; i < 1000; i++ {
		choice := rand.Float64()

		var key []byte
		value := tu.RandomBytes(32)

		if choice < 0.5 {
			key = append(prefixA, tu.RandomBytes(24)...)
			expectedDataA[string(key)] = value
		} else {
			key = append(prefixB, tu.RandomBytes(24)...)
			expectedDataB[string(key)] = value
		}

		err := store.Put(key, value)
		assert.NoError(t, err)
	}

	// Iterate over the store with prefixA and check that the data matches the expected data.
	foundKeysA := make(map[string]bool)
	iteratorA, err := store.NewIterator(prefixA)
	defer iteratorA.Release()
	assert.NoError(t, err)

	index := 0

	for iteratorA.Next() {
		index++

		key := string(iteratorA.Key())
		value := iteratorA.Value()

		expectedValue, ok := expectedDataA[key]
		assert.True(t, ok)
		assert.Equal(t, expectedValue, value)

		foundKeysA[key] = true
	}
	assert.Equal(t, len(expectedDataA), len(foundKeysA))

	// Iterate over the store with prefixB and check that the data matches the expected data.
	foundKeysB := make(map[string]bool)
	iteratorB, err := store.NewIterator(prefixB)
	defer iteratorB.Release()

	assert.NoError(t, err)

	for iteratorB.Next() {
		key := string(iteratorB.Key())
		value := iteratorB.Value()

		expectedValue, ok := expectedDataB[key]
		assert.True(t, ok)
		assert.Equal(t, expectedValue, value)

		foundKeysB[key] = true
	}
	assert.Equal(t, len(expectedDataB), len(foundKeysB))

	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestIterationWithPrefix(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		iterationWithPrefixTest(t, store)
	}
}

func putNilTest(t *testing.T, store kvstore.Store) {
	tu.InitializeRandom()
	deleteDBDirectory(t)

	key := tu.RandomBytes(32)

	err := store.Put(key, nil)
	assert.NoError(t, err)

	value, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, value)

	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestPutNil(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	for _, builder := range storeBuilders {
		store, err := builder(logger, dbPath)
		assert.NoError(t, err)
		putNilTest(t, store)
	}
}
