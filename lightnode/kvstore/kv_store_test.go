package kvstore

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func randomOperationsTest(t *testing.T, store KVStore) {
	tu.InitializeRandom()

	var err error

	expectedData := make(map[string][]byte)

	for i := 0; i < 1000; i++ {

		choice := rand.Float64()
		if len(expectedData) == 0 || choice < 0.66 {
			// Write a random value.

			key := tu.RandomBytes(32)
			value := tu.RandomBytes(32)

			err = store.Put(key, value, 0)
			assert.NoError(t, err)

			expectedData[string(key)] = value
		} else if choice < 0.90 {
			// Drop a random value.

			var key string
			for k := range expectedData {
				key = k
			}
			delete(expectedData, key)
			err = store.Drop([]byte(key))
			assert.NoError(t, err)
		} else {
			// Drop a non-existent value.

			key := tu.RandomBytes(32)
			err = store.Drop(key)
			assert.NoError(t, err)
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
			assert.NoError(t, err)
			assert.Nil(t, value)
		}
	}

	err = store.Shutdown()
	assert.NoError(t, err)
	err = store.Destroy()
	assert.NoError(t, err)
}

func TestRandomOperations(t *testing.T) {
	randomOperationsTest(t, NewInMemoryChunkStore())
}

func operationsOnShutdownStoreTest(t *testing.T, store KVStore) {
	err := store.Shutdown()
	assert.NoError(t, err)

	err = store.Put([]byte("key"), []byte("value"), 0)
	assert.Error(t, err)

	_, err = store.Get([]byte("key"))
	assert.Error(t, err)

	err = store.Drop([]byte("key"))
	assert.Error(t, err)

	err = store.Shutdown()
	assert.NoError(t, err)

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestOperationsOnShutdownStore(t *testing.T) {
	operationsOnShutdownStoreTest(t, NewInMemoryChunkStore())
}
