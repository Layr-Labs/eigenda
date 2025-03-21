package test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

// A list of builders for various stores to be tested.
var storeBuilders = []func(logger logging.Logger, path string) (kvstore.Store[[]byte], error){
	//func(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
	//	return mapstore.NewStore(), nil
	//},
	//func(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
	//	return leveldb.NewStore(logger, path, nil)
	//},
	//func(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
	//	config := tablestore.DefaultMapStoreConfig()
	//	config.Schema = []string{"test"}
	//	tableStore, err := tablestore.Start(logger, config)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return NewTableAsAStore(tableStore)
	//},
	//func(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
	//	config := tablestore.DefaultLevelDBConfig(path)
	//	config.Schema = []string{"test"}
	//	tableStore, err := tablestore.Start(logger, config)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return NewTableAsAStore(tableStore)
	//},
	func(logger logging.Logger, path string) (kvstore.Store[[]byte], error) {
		config := tablestore.DefaultLotusDBConfig(path)
		config.Schema = []string{"test"}
		tableStore, err := tablestore.Start(logger, config)
		if err != nil {
			return nil, err
		}
		return NewTableAsAStore(tableStore)
	},
}

func randomOperationsTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

	expectedData := make(map[string][]byte)

	for i := 0; i < 1000; i++ {
		choice := rand.Float64()
		if len(expectedData) == 0 || choice < 0.50 {
			// Write a random value.

			key := tu.RandomBytes(32)
			value := tu.RandomBytes(32)

			err := store.Put(key, value)
			require.NoError(t, err)

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
			require.NoError(t, err)
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
			require.NoError(t, err)
		} else {
			// Drop a non-existent value.

			key := tu.RandomBytes(32)
			err := store.Delete(key)
			require.Nil(t, err)
		}

		if i%10 == 0 {
			// Every so often, check that the store matches the expected data.
			for key, expectedValue := range expectedData {
				value, err := store.Get([]byte(key))
				require.NoError(t, err)
				require.Equal(t, expectedValue, value)
			}

			// Try and get a value that isn't in the store.
			key := tu.RandomBytes(32)
			value, err := store.Get(key)
			require.Equal(t, kvstore.ErrNotFound, err)
			require.Nil(t, value)
		}
	}

	err := store.Shutdown()
	require.NoError(t, err)
	err = store.Destroy()
	require.NoError(t, err)
}

func TestRandomOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		randomOperationsTest(t, store)
	}
}

func writeBatchTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

	rand := random.NewTestRandom()

	var err error

	expectedData := make(map[string][]byte)
	batch := store.NewBatch()

	for i := 0; i < 1000; i++ {
		// Write a random value.
		key := []byte(rand.String(32))

		var value []byte
		//if i%50 == 0 { // TODO
		//	// nil values are interpreted as empty slices.
		//	value = nil
		//} else {
		value = tu.RandomBytes(32)
		//}

		batch.Put(key, value)

		if value == nil {
			expectedData[string(key)] = []byte{}
		} else {
			expectedData[string(key)] = value
		}

		if true { //i%10 == 0 {
			// Every so often, apply the batch and check that the store matches the expected data.

			err = batch.Apply()
			require.NoError(t, err)

			for key, expectedValue := range expectedData {
				value, err = store.Get([]byte(key))
				require.NoError(t, err)
				require.Equal(t, expectedValue, value)
			}

			// Try and get a value that isn't in the store.
			key = tu.RandomBytes(32)
			value, err = store.Get(key)
			require.Equal(t, kvstore.ErrNotFound, err)
			require.Nil(t, value)

			batch = store.NewBatch()
		}
	}

	err = batch.Apply()
	require.NoError(t, err)

	err = store.Shutdown()
	require.NoError(t, err)
	err = store.Destroy()
	require.NoError(t, err)
}

func TestWriteBatch(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		writeBatchTest(t, store)
	}
}

func deleteBatchTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

	expectedData := make(map[string][]byte)

	// Add some data to the store.
	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(32)
		value := tu.RandomBytes(32)

		err := store.Put(key, value)
		require.NoError(t, err)

		expectedData[string(key)] = value
	}

	batch := store.NewBatch()

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
	require.NoError(t, err)

	// Check that the store matches the expected data.
	for key, expectedValue := range expectedData {
		value, err := store.Get([]byte(key))
		require.NoError(t, err)
		require.Equal(t, expectedValue, value)
	}

	// Try and get a value that isn't in the store.
	key := tu.RandomBytes(32)
	value, err := store.Get(key)
	require.Equal(t, kvstore.ErrNotFound, err)
	require.Nil(t, value)

	err = store.Shutdown()
	require.NoError(t, err)
	err = store.Destroy()
	require.NoError(t, err)
}

func TestDeleteBatch(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		deleteBatchTest(t, store)
	}
}

func iterationTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

	expectedData := make(map[string][]byte)

	// Insert some data into the store.
	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(32)
		value := tu.RandomBytes(32)

		err := store.Put(key, value)
		require.NoError(t, err)

		expectedData[string(key)] = value
	}

	// Iterate over the store and check that the data matches the expected data.
	foundKeys := make(map[string]bool)

	iterator, err := store.NewIterator(nil)
	require.NoError(t, err)

	for iterator.Next() {
		key := string(iterator.Key())
		value := iterator.Value()

		expectedValue, ok := expectedData[key]
		require.True(t, ok)
		require.Equal(t, expectedValue, value)

		foundKeys[key] = true
	}
	require.Equal(t, len(expectedData), len(foundKeys))

	iterator.Release()
	err = store.Destroy()
	require.NoError(t, err)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		iterationTest(t, store)
	}
}

func iterationWithPrefixTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

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
		require.NoError(t, err)
	}

	// Iterate over the store with prefixA and check that the data matches the expected data.
	foundKeysA := make(map[string]bool)
	iteratorA, err := store.NewIterator(prefixA)
	require.NoError(t, err)

	index := 0

	for iteratorA.Next() {
		index++

		key := string(iteratorA.Key())
		value := iteratorA.Value()

		expectedValue, ok := expectedDataA[key]
		require.True(t, ok)
		require.Equal(t, expectedValue, value)

		foundKeysA[key] = true
	}
	require.Equal(t, len(expectedDataA), len(foundKeysA))
	iteratorA.Release()

	// Iterate over the store with prefixB and check that the data matches the expected data.
	foundKeysB := make(map[string]bool)
	iteratorB, err := store.NewIterator(prefixB)

	require.NoError(t, err)

	for iteratorB.Next() {
		key := string(iteratorB.Key())
		value := iteratorB.Value()

		expectedValue, ok := expectedDataB[key]
		require.True(t, ok)
		require.Equal(t, expectedValue, value)

		foundKeysB[key] = true
	}
	require.Equal(t, len(expectedDataB), len(foundKeysB))
	iteratorB.Release()

	err = store.Destroy()
	require.NoError(t, err)
}

func TestIterationWithPrefix(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		iterationWithPrefixTest(t, store)
	}
}

func putNilTest(t *testing.T, store kvstore.Store[[]byte]) {
	tu.InitializeRandom()

	key := tu.RandomBytes(32)

	err := store.Put(key, nil)
	require.NoError(t, err)

	value, err := store.Get(key)
	require.NoError(t, err)
	require.Equal(t, []byte{}, value)

	err = store.Destroy()
	require.NoError(t, err)
}

func TestPutNil(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range storeBuilders {
		dbPath := t.TempDir()
		store, err := builder(logger, dbPath)
		require.NoError(t, err)
		putNilTest(t, store)
	}
}
