package kvstore

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func simpleWritingBenchmark(b *testing.B, store KVStore) {
	keySize := 8
	valueSize := 1024

	if store == nil {
		panic("store is nil") // todo
	}

	// reuse the byte arrays for all operations to avoid the overhead of generating random data.
	baseKey := tu.RandomBytes(keySize)
	baseValue := tu.RandomBytes(valueSize)

	bytesToWrite := 1 * 1024 * 1024 * 1024 // 1 GB
	keysToWrite := bytesToWrite / valueSize

	for i := 0; i < keysToWrite; i++ {

		//if i%1000 == 0 {
		//	fmt.Printf("i: %d\n", i) // TODO
		//}

		// Change a few bytes in the key to avoid collisions. Change a few bytes in the value to
		// avoid the DB taking shortcuts (since we aren't using random data for the sake of benchmark performance).

		baseKey[0] ^= byte(i)
		baseKey[1] ^= byte(i >> 8)
		baseKey[2] ^= byte(i >> 16)
		baseKey[3] ^= byte(i >> 24)

		baseValue[0] ^= byte(i)
		baseValue[1] ^= byte(i >> 8)
		baseValue[2] ^= byte(i >> 16)
		baseValue[3] ^= byte(i >> 24)

		err := store.Put(baseKey, baseValue, 0)
		assert.NoError(b, err)
	}

	err := store.Destroy()
	assert.NoError(b, err)

	_, err = os.Stat(dbPath)
	assert.True(b, os.IsNotExist(err))
}

//func BenchmarkWritingInMemory(b *testing.B) {
//	fmt.Println("-------------------------------------------------- BenchmarkWritingInMemory") // TODO
//
//	store := NewInMemoryStore()
//	simpleWritingBenchmark(b, store)
//}

func BenchmarkWritingLevelDB(b *testing.B) {
	fmt.Println("-------------------------------------------------- BenchmarkWritingLevelDB") // TODO

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewLevelStore(logger, dbPath)
	assert.NoError(b, err)

	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBatchedLevelDB(b *testing.B) {

	fmt.Println("-------------------------------------------------- BenchmarkWritingBatchedLevelDB") // TODO

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewLevelStore(logger, dbPath)
	store = BatchingWrapper(store, 1024*32)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBadgerDB(b *testing.B) {

	fmt.Println("-------------------------------------------------- BenchmarkWritingBadgerDB") // TODO

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewBadgerStore(logger, dbPath)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBatchedBadgerDB(b *testing.B) {

	fmt.Println("-------------------------------------------------- BenchmarkWritingBatchedBadgerDB") // TODO

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewBadgerStore(logger, dbPath)
	store = BatchingWrapper(store, 1024*1024)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingPebble(b *testing.B) {

	fmt.Println("-------------------------------------------------- BenchmarkWritingPebbleDB") // TODO

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewPebbleStore(logger, dbPath)
	store = BatchingWrapper(store, 1024*1024*32)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}
