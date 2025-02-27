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
	"github.com/dgraph-io/badger/v4"
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

	reportInterval := units.MiB * 100
	interval := 0
	previousIntervalDataWritten := uint64(0)
	previousIntervalTimestamp := time.Now()

	for dataWritten < totalDataToWrite {
		key := rand.Bytes(32)
		value := rand.Bytes(int(dataSize))

		err := write(key, value)
		require.NoError(t, err)

		dataWritten += dataSize

		// Do some console logging
		newInterval := int(dataWritten / uint64(reportInterval))
		if newInterval > interval {
			interval = newInterval

			timeSinceLastInterval := time.Since(previousIntervalTimestamp)
			previousIntervalTimestamp = time.Now()
			dataSinceLastInterval := dataWritten - previousIntervalDataWritten
			previousIntervalDataWritten = dataWritten

			mbPerSecondOverLastInterval := float64(dataSinceLastInterval) / (units.MiB) /
				(float64(timeSinceLastInterval.Nanoseconds()) / float64(time.Second))

			timeSinceStart := time.Since(start)
			mbPerSecondTotal := float64(dataWritten) / (units.MiB) /
				(float64(timeSinceStart.Nanoseconds()) / float64(time.Second))

			mbWritten := float64(dataWritten) / units.MiB

			completionPercentage := int(float64(dataWritten) / float64(totalDataToWrite) * 100)

			fmt.Printf("%d%%: wrote %d MiB. %d mb/s during recent period, %d mb/s overall.\r",
				completionPercentage, int(mbWritten), int(mbPerSecondOverLastInterval), int(mbPerSecondTotal))
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

	err = db.Start()
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

	err = db.Start()
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

func TestBadgerDBWrite(t *testing.T) {
	directory := "./test-data"
	opts := badger.DefaultOptions(directory)

	opts.Logger = nil
	db, err := badger.Open(opts)
	require.NoError(t, err)

	batch := db.NewWriteBatch()
	objectsInBatch := 0

	writeFunction := func(key []byte, value []byte) error {
		err = batch.Set(key, value)
		if err != nil {
			return err
		}
		objectsInBatch++

		if objectsInBatch >= batchSize {
			err = batch.Flush()
			if err != nil {
				return err
			}
			batch = db.NewWriteBatch()
			objectsInBatch = 0
		}

		return nil
	}

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)

	err = db.Close()
	require.NoError(t, err)
}
