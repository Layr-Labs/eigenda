package dbbench

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/cockroachdb/pebble"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

type writer func(key []byte, value []byte) error

const totalToWrite = 50 * units.GiB

// const totalToWrite = 1024 * units.GiB * 10
const dataSize = 1 * units.MiB
const batchSize = 100

const littDBShards = 1

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

func TestLevelDBNoCompactionWrite(t *testing.T) {
	directory := "./test-data"
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	config := tablestore.DefaultLevelDBConfig(directory)
	config.DisableCompaction = true
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

	config, err := littbuilder.DefaultConfig(directory)
	require.NoError(t, err)
	config.ShardingFactor = littDBShards

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

func TestLittDBWithGCWrite(t *testing.T) {
	directory := "./test-data"

	config, err := littbuilder.DefaultConfig(directory)
	require.NoError(t, err)
	config.ShardingFactor = littDBShards
	config.TTL = 2 * time.Hour

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

	config, err := littbuilder.DefaultConfig(directory)
	require.NoError(t, err)
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

func TestBadgerDBWithGCWrite(t *testing.T) {
	directory := "./test-data"
	opts := badger.DefaultOptions(directory)
	opts.Compression = options.None
	opts.CompactL0OnClose = true

	db, err := badger.Open(opts)
	require.NoError(t, err)

	transaction := db.NewTransaction(true)
	objectsInBatch := 0

	ttl := 5 * time.Minute

	keys := make([][]byte, 0)
	writeFunction := func(key []byte, value []byte) error {
		keys = append(keys, key)

		entry := badger.NewEntry(key, value).WithTTL(ttl)
		err = transaction.SetEntry(entry)

		if err != nil {
			return err
		}
		objectsInBatch++

		if objectsInBatch >= batchSize {
			err = transaction.Commit()

			if err != nil {
				return err
			}
			transaction = db.NewTransaction(true)
			objectsInBatch = 0
		}

		return nil
	}

	alive := atomic.Bool{}
	alive.Store(true)
	compactionDone := make(chan struct{})
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer func() {
			compactionDone <- struct{}{}
		}()
		for alive.Load() {
			<-ticker.C

			fmt.Printf("\nRunning GC\n")
			startTime := time.Now()

			gcIterations := 0
			for {
				gcIterations++
				err = db.RunValueLogGC(0.125)
				if err != nil {
					if !strings.Contains(err.Error(), "Value log GC attempt didn't result in any cleanup") {
						fmt.Printf("\nError running GC: %v\n", err)
					}
					break
				}
			}

			err = db.Flatten(1)
			if err != nil {
				fmt.Printf("\nError flattening DB: %v\n", err)
			}

			fmt.Printf("\nGC took %v, did %d iterations\n", time.Since(startTime), gcIterations)
		}
	}()

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)
	alive.Store(false)
	<-compactionDone

	fmt.Printf("doing some final compaction to see what happens\n")
	fmt.Printf("First, let's sleep for a little while (2 minutes)")
	time.Sleep(2 * time.Minute)
	fmt.Printf("Now, let's run the compaction\n")
	iterations := 0
	for {
		iterations++
		err = db.RunValueLogGC(0.125)
		if err != nil {
			if !strings.Contains(err.Error(), "Value log GC attempt didn't result in any cleanup") {
				fmt.Printf("\nError running GC: %v\n", err)
			}
			break
		}
	}
	fmt.Printf("Compaction took %d iterations\n", iterations)
	fmt.Printf("Now, let's flatten the DB\n")
	err = db.Flatten(1)
	if err != nil {
		fmt.Printf("\nError flattening DB: %v\n", err)
	}

	fmt.Printf("checking to see what keys are still present. Based on timing, all keys should be expired.\n")
	keysPresent := 0
	keysMissing := 0
	err = db.View(func(txn *badger.Txn) error {
		for _, key := range keys {
			_, err = txn.Get(key)
			if err == nil {
				keysPresent++
			} else if errors.Is(badger.ErrKeyNotFound, err) {
				keysMissing++
			} else {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	fmt.Printf("Keys present: %d, keys missing: %d\n", keysPresent, keysMissing)

	err = db.Close()
	require.NoError(t, err)
}

func TestPebbleDBWrite(t *testing.T) {
	directory := "./test-data"

	options := &pebble.Options{}

	db, err := pebble.Open(directory, options)
	require.NoError(t, err)

	batch := db.NewBatch()
	objectsInBatch := 0

	writeFunction := func(key []byte, value []byte) error {

		err = batch.Set(key, value, nil)
		if err != nil {
			return err
		}
		objectsInBatch++

		if objectsInBatch >= batchSize {
			err = batch.Commit(nil)
			if err != nil {
				return err
			}
			batch = db.NewBatch()
			objectsInBatch = 0
		}

		return nil
	}

	runWriteBenchmark(t, writeFunction, totalToWrite, dataSize)

	err = db.Close()
	require.NoError(t, err)
}
