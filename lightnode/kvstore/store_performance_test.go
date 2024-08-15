package kvstore

import (
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func simpleWritingBenchmark(b *testing.B, store KVStore) {
	keySize := 8
	valueSize := 1024

	// reuse the byte arrays for all operations to avoid the overhead of generating random data.
	baseKey := tu.RandomBytes(keySize)
	baseValue := tu.RandomBytes(valueSize)

	bytesToWrite := 1 * 1024 * 1024 * 1024 // 1GB
	keysToWrite := bytesToWrite / valueSize

	for i := 0; i < b.N; i++ {

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

		for j := 0; j < keysToWrite; j++ {
			err := store.Put(baseKey, baseValue, 0)
			assert.NoError(b, err)
		}
	}

	//fmt.Printf("Wrote %d bytes\n", bytesToWrite)
	//fmt.Printf("Value size: %d bytes\n", valueSize)
	//fmt.Printf("Wrote %d keys\n", keysToWrite)
	//fmt.Printf("Completed in %f s\n", float64(b.Elapsed())/1000/1000/1000)

	err := store.Destroy()
	assert.NoError(b, err)
}

//func BenchmarkWritingInMemory(b *testing.B) {
//	store := NewInMemoryStore()
//	simpleWritingBenchmark(b, store)
//}

func BenchmarkWritingLevelDB(b *testing.B) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewLevelStore(logger, "testdb")
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBatchedLevelDB(b *testing.B) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewLevelStore(logger, "testdb")
	store = BatchingWrapper(store, 1024*32)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBadgerDB(b *testing.B) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewBadgerStore(logger, "testdb")
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}

func BenchmarkWritingBatchedBadgerDB(b *testing.B) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := NewBadgerStore(logger, "testdb")
	store = BatchingWrapper(store, 1024*1024)
	assert.NoError(b, err)
	simpleWritingBenchmark(b, store)
}
