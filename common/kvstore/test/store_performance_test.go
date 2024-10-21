package test

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func writeThenReadBenchmark(b *testing.B, store kvstore.Store[[]byte]) {
	keySize := 8
	valueSize := 1024

	bytesToWrite := 1 * 1024 * 1024 * 1024 // 1 GB
	keysToWrite := bytesToWrite / valueSize

	keysToRead := keysToWrite / 10

	// reuse the byte arrays for all operations to avoid the overhead of generating random data.
	baseKey := tu.RandomBytes(keySize)
	baseValue := tu.RandomBytes(valueSize)

	start := time.Now()

	for i := 0; i < keysToWrite; i++ {

		if i%1000 == 0 {
			fmt.Printf("writing %d\n", i)
		}

		// Change a few bytes in the key to avoid collisions. Change a few bytes in the value to
		// avoid the DB taking shortcuts (since we aren't using random data for the sake of benchmark performance).

		key := make([]byte, keySize)
		copy(key, baseKey)

		key[0] ^= byte(i)
		key[1] ^= byte(i >> 8)
		key[2] ^= byte(i >> 16)
		key[3] ^= byte(i >> 24)

		baseValue[0] ^= byte(i)
		baseValue[1] ^= byte(i >> 8)
		baseValue[2] ^= byte(i >> 16)
		baseValue[3] ^= byte(i >> 24)

		err := store.Put(key, baseValue)
		assert.NoError(b, err)
	}

	doneWriting := time.Now()

	for i := 0; i < keysToRead; i++ {

		if i%1000 == 0 {
			fmt.Printf("reading %d\n", i)
		}

		keyIndex := rand.Intn(keysToWrite)

		key := make([]byte, keySize)
		copy(key, baseKey)

		key[0] ^= byte(keyIndex)
		key[1] ^= byte(keyIndex >> 8)
		key[2] ^= byte(keyIndex >> 16)
		key[3] ^= byte(keyIndex >> 24)

		_, err := store.Get(key)
		assert.NoError(b, err)
	}

	doneReading := time.Now()

	err := store.Destroy()
	assert.NoError(b, err)

	doneDestroying := time.Now()

	timeToWrite := doneWriting.Sub(start)
	timeToRead := doneReading.Sub(doneWriting)
	timeToDestroy := doneDestroying.Sub(doneReading)

	timeToWriteSeconds := float64(timeToWrite) / float64(time.Second)
	timeToReadSeconds := float64(timeToRead) / float64(time.Second)
	timeToDestroySeconds := float64(timeToDestroy) / float64(time.Second)

	mbWrittenPerSecond := float64(bytesToWrite) / timeToWriteSeconds / float64(1024*1024)
	mbReadPerSecond := float64(keysToRead*valueSize) / timeToReadSeconds / float64(1024*1024)

	fmt.Printf("Bytes written: %d\n", bytesToWrite)
	fmt.Printf("Time to write: %.1fs\n", timeToWriteSeconds)
	fmt.Printf("Time to read: %.1fs\n", timeToReadSeconds)
	fmt.Printf("Time to destroy: %.1fs\n", timeToDestroySeconds)
	fmt.Printf("Write speed: %.1f MB/s\n", mbWrittenPerSecond)
	fmt.Printf("Read speed: %.1f MB/s\n", mbReadPerSecond)
}

func BenchmarkWritingInMemory(b *testing.B) {
	store := mapstore.NewStore()
	writeThenReadBenchmark(b, store)
}

func BenchmarkLevelDB(b *testing.B) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(b, err)

	store, err := leveldb.NewStore(logger, dbPath)
	assert.NoError(b, err)

	writeThenReadBenchmark(b, store)
}
