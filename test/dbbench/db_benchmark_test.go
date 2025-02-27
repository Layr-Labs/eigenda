package dbbench

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

type writer func(key []byte, value []byte) error

const totalToWrite = 10 * units.GiB
const dataSize = 1 * units.MiB
const batchSize = 100

// runWriteBenchmark runs a simple benchmark. Its goal is to write a ton of data to the database as fast as possible.
func runWriteBenchmark(
	t *testing.T,
	write writer,
	totalDataToWrite uint64,
	dataSize uint64) {

	rand := random.NewTestRandom()

	start := time.Now()
	dataWritten := uint64(0)

	reportInterval := units.MiB * 10
	interval := 0

	for dataWritten < totalDataToWrite {
		key := rand.Bytes(32)
		value := rand.Bytes(int(dataSize))

		err := write(key, value)
		require.NoError(t, err)

		dataWritten += dataSize

		newInterval := int(dataWritten / uint64(reportInterval))
		if newInterval > interval {
			fmt.Printf("wrote %d MiB\r", dataWritten/units.MiB)
			interval = newInterval
		}
	}

	fmt.Printf("\n")

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Write benchmark took %v\n", elapsed)
	mbPerSecond := float64(totalDataToWrite) / (units.MiB) / (float64(elapsed.Nanoseconds()) / float64(time.Second))
	fmt.Printf("Write benchmark speed: %.2f MB/s\n", mbPerSecond)
}

func TestLevelDBWrite(t *testing.T) {
	directory := "./test-data"
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	config := tablestore.DefaultLevelDBConfig(directory)
	config.Schema = []string{"test"}
	db, err := tablestore.Start(logger, config)
	require.NoError(t, err)

	keyBuilder, err := db.GetKeyBuilder("test")
	require.NoError(t, err)

	batch := db.NewBatch()

	writeFunction := func(key []byte, value []byte) error {

		batch.Put(keyBuilder.Key(key), value)
		if batch.Size() >= batchSize {
			err = batch.Apply()
			if err != nil {
				return err
			}
			batch = db.NewBatch()
		}

		return nil
	}

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)

	err = db.Shutdown()
	require.NoError(t, err)
}

func TestLittDBWrite(t *testing.T) {
	directory := "./test-data"

	config := littbuilder.DefaultConfig(directory)

	db, err := config.Build(context.Background())
	require.NoError(t, err)

	table, err := db.GetTable("test")
	require.NoError(t, err)

	unflushedCount := uint32(0)

	writeFunction := func(key []byte, value []byte) error {
		err = table.Put(key, value)
		if err != nil {
			return err
		}

		unflushedCount++
		if unflushedCount >= batchSize {
			err = table.Flush()
			if err != nil {
				return err
			}
			unflushedCount = 0
		}

		return nil
	}

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)

	err = db.Stop()
	require.NoError(t, err)
}

func TestMemKeymapLittDBWrite(t *testing.T) {
	directory := "./test-data"

	config := littbuilder.DefaultConfig(directory)
	config.KeyMapType = littbuilder.MemKeyMap

	db, err := config.Build(context.Background())
	require.NoError(t, err)

	table, err := db.GetTable("test")
	require.NoError(t, err)

	unflushedCount := uint32(0)

	writeFunction := func(key []byte, value []byte) error {
		err = table.Put(key, value)
		if err != nil {
			return err
		}

		unflushedCount++
		if unflushedCount >= batchSize {
			err = table.Flush()
			if err != nil {
				return err
			}
			unflushedCount = 0
		}

		return nil
	}

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)

	err = db.Stop()
	require.NoError(t, err)
}
