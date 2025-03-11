package dbbench

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"github.com/docker/go-units"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/stretchr/testify/require"
)

// writer is a function that writes a key-value pair to a database.
type writer func(key []byte, value []byte) error

// reader is a function that reads a key-value pair from a database.
type reader func(key []byte) ([]byte, error)

const totalToWrite = 10 * units.TiB
const dataSize = 1 * units.MiB
const batchSize = 100
const parallelWriters = 2
const readBytesPerSecond = 1 * units.MiB
const readerCount = 1
const TTL = 2 * time.Hour

// given a seed, deterministically generate a key/value pair
func generateKVPair(seed int64) ([]byte, []byte) {
	rand := random.NewTestRandomNoPrint(seed)
	key := []byte(rand.String(32))
	value := rand.Bytes(dataSize)
	return key, value
}

// Describes data with an expiration time. In order to determine the key and the expected value, use the seed with
// the generateKVPair() function.
type dataWithExpiration struct {
	seed       int64
	expiration time.Time
}

// randomUnexpiredSeed returns a random seed from the unexpiredData map. If no unexpired data is present, it returns
// false.
func randomUnexpiredSeed(
	lock *sync.RWMutex,
	unexpiredData map[int64]struct{}) (int64, bool) {
	lock.RLock()
	defer lock.RUnlock()
	for seed := range unexpiredData {
		return seed, true
	}
	return 0, false
}

// TODO separate out the benchmark framework from the DB implementations, make this into a struct, not a mega function

// runBenchmark runs a simple benchmark. Its goal is to write a ton of data to the database as fast as possible.
func runBenchmark(write writer, read reader) {

	rand := random.NewTestRandom()

	fmt.Printf("Starting benchmark\n")
	fmt.Printf("Writing %d bytes\n", totalToWrite)
	fmt.Printf("Data size: %d bytes\n", dataSize)
	fmt.Printf("Batch size: %d\n", batchSize)
	fmt.Printf("Read bytes per second: %d\n", readBytesPerSecond)
	fmt.Printf("Reader count: %d\n", readerCount)
	fmt.Printf("TTL: %v\n", TTL)
	fmt.Printf("Write parallelism: %d\n", parallelWriters)

	start := time.Now()
	dataWritten := uint64(0)
	alive := atomic.Bool{}
	alive.Store(true)
	defer alive.Store(false)

	// data in the database that is expected to be present when the reader threads read it. This map implements
	// a set containing the seed values for all key-value pairs that are expected to be present in the database.
	unexpiredData := map[int64]struct{}{}
	// manually track expiration to maintain the unexpiredData map, which is needed by the reader threads.
	expirationQueue := linkedlistqueue.New()

	// protects access to unexpiredData and expirationQueue
	lock := &sync.RWMutex{}

	// The DB is permitted to buffer up a certain number of reads. Don't add a value to the unexpiredData map until
	// we are certain that the DB has flushed the batch.
	possiblyUnflushedData := linkedlistqueue.New()

	// variables for generating reports
	reportInterval := units.MiB * 100
	interval := 0
	previousIntervalDataWritten := uint64(0)
	previousIntervalTimestamp := time.Now()

	// Set up a goroutine to handle removal of elements from the unexpiredData map.
	gcTicker := time.NewTicker(1 * time.Second)
	gcDone := make(chan struct{})
	go func() {
		for alive.Load() {
			<-gcTicker.C

			lock.Lock()
			for {
				next, ok := expirationQueue.Peek()
				if !ok {
					break
				}

				data := next.(dataWithExpiration)
				if time.Now().After(data.expiration) {
					expirationQueue.Dequeue()
					delete(unexpiredData, data.seed)
				} else {
					break
				}
			}
			lock.Unlock()
		}
		gcDone <- struct{}{}
	}()
	defer func() {
		<-gcDone
	}()

	readRatePerGoroutine := readBytesPerSecond / readerCount
	readsPerSecond := readRatePerGoroutine / dataSize
	readerDoneChannels := make([]chan struct{}, readerCount)
	totalReadsPerformed := atomic.Uint64{}
	totalNanosecondsSpentOnReads := atomic.Uint64{}

	longestRead := atomic.Pointer[time.Duration]{}
	longestReadLock := &sync.Mutex{}
	defaultLongestRead := time.Duration(0)
	longestRead.Store(&defaultLongestRead)

	// Set up goroutines to read data from the database
	for i := 0; i < readerCount; i++ {
		readerDoneChannels[i] = make(chan struct{})
		readerTicker := time.NewTicker(time.Second / time.Duration(readsPerSecond))
		readerIndex := i
		go func() {
			for alive.Load() {
				<-readerTicker.C
				seed, ok := randomUnexpiredSeed(lock, unexpiredData)
				if !ok {
					continue
				}
				key, expectedValue := generateKVPair(seed)
				readStart := time.Now()
				value, err := read(key)
				readLatency := time.Since(readStart)
				if err != nil {
					panic(fmt.Errorf("error reading key %v: %v", key, err))
				}
				if !bytes.Equal(expectedValue, value) {
					panic(fmt.Errorf("expected value %v, got %v", expectedValue, value))
				}
				totalReadsPerformed.Add(1)
				totalNanosecondsSpentOnReads.Add(uint64(readLatency.Nanoseconds()))

				longestReadLock.Lock()
				if readLatency > *longestRead.Load() {
					longestRead.Store(&readLatency)
				}
				longestReadLock.Unlock()
			}
			readerDoneChannels[readerIndex] <- struct{}{}
		}()
	}
	defer func() {
		for _, done := range readerDoneChannels {
			<-done
		}
	}()

	longestWrite := atomic.Pointer[time.Duration]{}
	longestWriteLock := &sync.Mutex{}
	defaultLongestWrite := time.Duration(0)
	longestWrite.Store(&defaultLongestWrite)

	// Write data to the database
	for dataWritten < totalToWrite {
		seed := rand.Int63()
		key, value := generateKVPair(seed)

		writeStart := time.Now()

		err := write(key, value)
		if err != nil {
			panic(fmt.Errorf("error writing key %v: %v", key, err))
		}

		writeLatency := time.Since(writeStart)
		longestWriteLock.Lock()
		if writeLatency > *longestWrite.Load() {
			longestWrite.Store(&writeLatency)
		}
		longestWriteLock.Unlock()

		dataWritten += dataSize

		possiblyUnflushedData.Enqueue(seed)
		lock.Lock()
		// Subtract 10 minutes from the actual expiration time to avoid race conditions with the reader threads.
		// This means that the reader threads will stop making attempts to read a key/value pair 10 minutes before
		// that key/value pair is actually scheduled to be deleted.
		expirationQueue.Enqueue(dataWithExpiration{seed: seed, expiration: time.Now().Add(TTL).Add(-10 * time.Minute)})

		if possiblyUnflushedData.Size() > batchSize*parallelWriters {
			// Data that has had a number of writes afterward that exceeds the maximum batchSize
			// it is guaranteed to be flushed if the DB respects batch sizes.
			next, _ := possiblyUnflushedData.Dequeue()
			unexpiredData[next.(int64)] = struct{}{}
		}
		lock.Unlock()

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

			completionPercentage := int(float64(dataWritten) / float64(totalToWrite) * 100)

			averageReadLatencyNanoseconds := uint64(0)
			if totalReadsPerformed.Load() > 0 {
				averageReadLatencyNanoseconds = totalNanosecondsSpentOnReads.Load() / totalReadsPerformed.Load()
			}
			averageReadLatencyMs := float64(averageReadLatencyNanoseconds) / float64(time.Millisecond)

			averageReadThroughputMBPerSecond := float64(totalReadsPerformed.Load()*dataSize) / units.MiB /
				timeSinceStart.Seconds()

			fmt.Printf("%d%%: wrote %d MiB. %0.1f mb/s recently, %0.1f mb/s overall. "+
				"Average rLat is %0.3fms, average rThrpt is %0.1f mb/s. Max read: %s, Max write: %s%s\r",
				completionPercentage, int(mbWritten), mbPerSecondOverLastInterval, mbPerSecondTotal,
				averageReadLatencyMs, averageReadThroughputMBPerSecond, *longestRead.Load(), *longestWrite.Load(),
				strings.Repeat(" ", 10))
		}
	}

	fmt.Printf("\n")

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Write benchmark took %v\n", elapsed)
	mbPerSecond := float64(totalToWrite) / (units.MiB) / (float64(elapsed.Nanoseconds()) / float64(time.Second))
	fmt.Printf("Write benchmark speed: %.2f MB/s\n", mbPerSecond)
}

func TestLevelDB(t *testing.T) {
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

	batch := db.NewTTLBatch()

	writeLimiter := make(chan struct{}, parallelWriters)

	writeFunction := func(key []byte, value []byte) error {
		batch.PutWithTTL(keyBuilder.Key(key), value, TTL)
		if batch.Size() >= batchSize {

			writeLimiter <- struct{}{}

			batchToFlush := batch
			go func() {
				err = batchToFlush.Apply()
				if err != nil {
					panic(err)
				}
				<-writeLimiter
			}()

			batch = db.NewTTLBatch()
		}

		return nil
	}

	readFunction := func(key []byte) ([]byte, error) {
		return db.Get(keyBuilder.Key(key))
	}

	runBenchmark(writeFunction, readFunction)

	err = db.Shutdown()
	require.NoError(t, err)
}

func TestLittDB(t *testing.T) {
	directory := "./test-data"

	config, err := littbuilder.DefaultConfig(directory)
	require.NoError(t, err)
	config.ShardingFactor = 1
	config.TTL = TTL
	config.KeyMapType = keymap.MemKeyMapType

	db, err := config.Build(context.Background())
	require.NoError(t, err)

	table, err := db.GetTable("test")
	require.NoError(t, err)

	unflushedCount := uint32(0)

	writeLimiter := make(chan struct{}, parallelWriters)

	writeFunction := func(key []byte, value []byte) error {
		err = table.Put(key, value)
		if err != nil {
			return err
		}

		unflushedCount++
		if unflushedCount >= batchSize {

			writeLimiter <- struct{}{}
			go func() {
				err = table.Flush()
				if err != nil {
					panic(err)
				}
				<-writeLimiter
			}()

			unflushedCount = 0
		}

		return nil
	}

	readFunction := func(key []byte) ([]byte, error) {
		value, ok, err := table.Get(key)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("key not found")
		}
		return value, nil
	}

	runBenchmark(writeFunction, readFunction)

	err = db.Stop()
	require.NoError(t, err)
}

func TestBadgerDB(t *testing.T) {
	directory := "./test-data"
	opts := badger.DefaultOptions(directory)
	opts.Compression = options.None
	opts.Logger = nil

	opts.ValueThreshold = 0

	opts.BaseTableSize = 10 * units.KiB
	opts.TableSizeMultiplier = 2

	opts.BaseLevelSize = 10 * units.KiB
	opts.LevelSizeMultiplier = 2

	opts.MemTableSize = 10 * units.KiB

	opts.NumMemtables = 1
	opts.NumLevelZeroTables = 1

	opts.MaxLevels = 10

	db, err := badger.Open(opts)
	require.NoError(t, err)

	transaction := db.NewTransaction(true)
	objectsInBatch := 0

	ttl := TTL

	writeLimiter := make(chan struct{}, parallelWriters)

	keys := make([][]byte, 0)
	writeFunction := func(key []byte, value []byte) error {
		keys = append(keys, key)

		entry := badger.NewEntry(key, value).WithTTL(ttl)
		err = transaction.SetEntry(entry)

		if err != nil {
			if strings.Contains(err.Error(), "Txn is too big to fit into one request") {
				err = transaction.Commit()
				if err != nil {
					return err
				}
				transaction = db.NewTransaction(true)
				objectsInBatch = 0
				err = transaction.SetEntry(entry)
			} else {
				return err
			}
			return err
		}
		objectsInBatch++

		if objectsInBatch >= batchSize {

			writeLimiter <- struct{}{}
			transactionToCommit := transaction

			go func() {
				err = transactionToCommit.Commit()
				if err != nil {
					panic(err)
				}
				<-writeLimiter
			}()

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

			err = db.Flatten(8)
			if err != nil {
				fmt.Printf("\nError flattening DB: %v\n", err)
			}

			gcIterations := 0
			for alive.Load() {
				gcIterations++
				err = db.RunValueLogGC(0.125)
				if err != nil {
					if !strings.Contains(err.Error(), "Value log GC attempt didn't result in any cleanup") {
						fmt.Printf("\nError running GC: %v\n", err)
					}
					break
				}
			}

			fmt.Printf("\nGC took %v, did %d iterations\n", time.Since(startTime), gcIterations)
		}
	}()

	readFunction := func(key []byte) ([]byte, error) {
		txn := db.NewTransaction(false)
		defer txn.Discard()
		item, err := txn.Get(key)
		if err != nil {
			return nil, err
		}
		value, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	runBenchmark(writeFunction, readFunction)
	alive.Store(false)
	<-compactionDone

	err = db.Close()
	require.NoError(t, err)
}
