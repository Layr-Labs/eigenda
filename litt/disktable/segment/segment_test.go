package segment

import (
	"bytes"
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/stretchr/testify/require"
)

// countFilesInDirectory returns the number of files in the given directory.
func countFilesInDirectory(t *testing.T, directory string) int {
	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	return len(files)
}

func TestWriteAndReadSegmentSingleShard(t *testing.T) {
	rand := random.NewTestRandom()
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(1000, 2000)
	keys := make([][]byte, valueCount)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		key := rand.PrintableVariableBytes(1, 100)
		keys[i] = key
		values[i] = rand.PrintableVariableBytes(1, 100)
	}

	// a map from keys to values
	expectedValues := make(map[string][]byte)

	// a map from keys to addresses
	addressMap := make(map[string]types.Address)

	expectedLargestShardSize := uint64(0)

	salt := rand.Uint32()
	seg, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		1,
		salt,
		false)

	require.NoError(t, err)

	// Write values to the segment.
	for i := 0; i < int(valueCount); i++ {
		key := keys[i]
		value := values[i]
		expectedValues[string(key)] = value

		expectedLargestShardSize += uint64(len(value)) + 4 /* uint32 length */

		largestShardSize, err := seg.Write(&types.KVPair{Key: key, Value: value})
		require.NoError(t, err)
		require.Equal(t, expectedLargestShardSize, largestShardSize)

		// Occasionally flush the segment to disk.
		if rand.BoolWithProbability(0.25) {
			flushedKeys, err := seg.Flush()()
			require.NoError(t, err)
			for _, flushedKey := range flushedKeys {
				addressMap[string(flushedKey.Key)] = flushedKey.Address
			}

			// after flushing, the address map should be the same size as the expected values map
			require.Equal(t, len(expectedValues), len(addressMap))
		}

		// Occasionally scan all addresses and values in the segment.
		if rand.BoolWithProbability(0.1) {
			flushedKeys, err := seg.Flush()()
			require.NoError(t, err)
			for _, flushedKey := range flushedKeys {
				addressMap[string(flushedKey.Key)] = flushedKey.Address
			}

			// after flushing, the address map should be the same size as the expected values map
			require.Equal(t, len(expectedValues), len(addressMap))

			for k, addr := range addressMap {
				readValue, err := seg.Read([]byte(k), addr)
				require.NoError(t, err)
				require.Equal(t, expectedValues[k], readValue)
			}
		}
	}

	// Seal the segment and read all keys and values.
	require.False(t, seg.IsSealed())
	sealTime := rand.Time()
	flushedKeys, err := seg.Seal(sealTime)
	require.NoError(t, err)
	require.True(t, seg.IsSealed())

	for _, flushedKey := range flushedKeys {
		addressMap[string(flushedKey.Key)] = flushedKey.Address
	}

	// after flushing, the address map should be the same size as the expected values map
	require.Equal(t, len(expectedValues), len(addressMap))

	require.Equal(t, sealTime.UnixNano(), seg.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment, err := seg.GetKeys()
	require.NoError(t, err)
	for i, ka := range keysFromSegment {
		require.Equal(t, ka.Key, keys[i])
	}

	// Reopen the segment and read all keys and values.
	seg2, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		1,
		salt,
		false)
	require.NoError(t, err)
	require.True(t, seg2.IsSealed())

	require.Equal(t, sealTime.UnixNano(), seg2.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg2.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment2, err := seg2.GetKeys()
	require.NoError(t, err)
	require.Equal(t, keysFromSegment, keysFromSegment2)

	// delete the segment
	require.Equal(t, 3, countFilesInDirectory(t, directory))

	err = seg.delete()
	require.NoError(t, err)

	require.Equal(t, 0, countFilesInDirectory(t, directory))
}

func TestWriteAndReadSegmentMultiShard(t *testing.T) {
	rand := random.NewTestRandom()
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(1000, 2000)
	shardCount := rand.Uint32Range(2, 32)
	keys := make([][]byte, valueCount)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		key := rand.PrintableVariableBytes(1, 100)
		keys[i] = key
		values[i] = rand.PrintableVariableBytes(1, 100)
	}

	// a map from keys to values
	expectedValues := make(map[string][]byte)

	// a map from keys to addresses
	addressMap := make(map[string]types.Address)

	salt := rand.Uint32()
	seg, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		shardCount,
		salt,
		false)

	require.NoError(t, err)

	// Write values to the segment.
	for i := 0; i < int(valueCount); i++ {
		key := keys[i]
		value := values[i]
		expectedValues[string(key)] = value

		largestShardSize, err := seg.Write(&types.KVPair{Key: key, Value: value})
		require.NoError(t, err)
		require.True(t, largestShardSize >= uint64(len(value)+4))

		// Occasionally flush the segment to disk.
		if rand.BoolWithProbability(0.25) {
			flushedKeys, err := seg.Flush()()
			require.NoError(t, err)
			for _, flushedKey := range flushedKeys {
				addressMap[string(flushedKey.Key)] = flushedKey.Address
			}

			// after flushing, the address map should be the same size as the expected values map
			require.Equal(t, len(expectedValues), len(addressMap))
		}

		// Occasionally scan all addresses and values in the segment.
		if rand.BoolWithProbability(0.1) {
			flushedKeys, err := seg.Flush()()
			require.NoError(t, err)
			for _, flushedKey := range flushedKeys {
				addressMap[string(flushedKey.Key)] = flushedKey.Address
			}

			// after flushing, the address map should be the same size as the expected values map
			require.Equal(t, len(expectedValues), len(addressMap))

			for k, addr := range addressMap {
				readValue, err := seg.Read([]byte(k), addr)
				require.NoError(t, err)
				require.Equal(t, expectedValues[k], readValue)
			}
		}
	}

	// Seal the segment and read all keys and values.
	require.False(t, seg.IsSealed())
	sealTime := rand.Time()
	flushedKeys, err := seg.Seal(sealTime)
	require.NoError(t, err)
	require.True(t, seg.IsSealed())

	for _, flushedKey := range flushedKeys {
		addressMap[string(flushedKey.Key)] = flushedKey.Address
	}

	// after flushing, the address map should be the same size as the expected values map
	require.Equal(t, len(expectedValues), len(addressMap))

	require.Equal(t, sealTime.UnixNano(), seg.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment, err := seg.GetKeys()
	require.NoError(t, err)
	// Sort keys. With more than one shard, keys may have random order.
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	sort.Slice(keysFromSegment, func(i, j int) bool {
		return bytes.Compare(keysFromSegment[i].Key, keysFromSegment[j].Key) < 0
	})
	for i, ka := range keysFromSegment {
		require.Equal(t, ka.Key, keys[i])
	}

	// Reopen the segment and read all keys and values.
	seg2, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		1,
		salt,
		false)
	require.NoError(t, err)
	require.True(t, seg2.IsSealed())

	require.Equal(t, sealTime.UnixNano(), seg2.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg2.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment2, err := seg2.GetKeys()
	sort.Slice(keysFromSegment2, func(i, j int) bool {
		return bytes.Compare(keysFromSegment2[i].Key, keysFromSegment2[j].Key) < 0
	})
	require.NoError(t, err)
	require.Equal(t, keysFromSegment, keysFromSegment2)

	// delete the segment
	require.Equal(t, int(2+shardCount), countFilesInDirectory(t, directory))

	err = seg.delete()
	require.NoError(t, err)

	require.Equal(t, 0, countFilesInDirectory(t, directory))
}

// Tests writing and reading, but allocates more shards than values written to force some shards to be empty.
func TestWriteAndReadColdShard(t *testing.T) {
	rand := random.NewTestRandom()
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	shardCount := rand.Uint32Range(2, 32)
	valueCount := shardCount * 2
	keys := make([][]byte, valueCount)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		key := rand.PrintableVariableBytes(1, 100)
		keys[i] = key
		values[i] = rand.PrintableVariableBytes(1, 100)
	}

	// a map from keys to values
	expectedValues := make(map[string][]byte)

	// a map from keys to addresses
	addressMap := make(map[string]types.Address)

	salt := rand.Uint32()
	seg, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		shardCount,
		salt,
		false)

	require.NoError(t, err)

	// Write values to the segment.
	for i := 0; i < int(valueCount); i++ {
		key := keys[i]
		value := values[i]
		expectedValues[string(key)] = value

		largestShardSize, err := seg.Write(&types.KVPair{Key: key, Value: value})
		require.NoError(t, err)
		require.True(t, largestShardSize >= uint64(len(value)+4))
	}

	// Seal the segment and read all keys and values.
	require.False(t, seg.IsSealed())
	sealTime := rand.Time()
	flushedKeys, err := seg.Seal(sealTime)
	require.NoError(t, err)
	require.True(t, seg.IsSealed())

	for _, flushedKey := range flushedKeys {
		addressMap[string(flushedKey.Key)] = flushedKey.Address
	}

	// after flushing, the address map should be the same size as the expected values map
	require.Equal(t, len(expectedValues), len(addressMap))

	require.Equal(t, sealTime.UnixNano(), seg.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment, err := seg.GetKeys()
	require.NoError(t, err)
	// Sort keys. With more than one shard, keys may have random order.
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	sort.Slice(keysFromSegment, func(i, j int) bool {
		return bytes.Compare(keysFromSegment[i].Key, keysFromSegment[j].Key) < 0
	})
	for i, ka := range keysFromSegment {
		require.Equal(t, ka.Key, keys[i])
	}

	// Reopen the segment and read all keys and values.
	seg2, err := NewSegment(
		context.Background(),
		logger,
		util.NewDBPanic(logger),
		index,
		[]string{directory},
		time.Now(),
		1,
		salt,
		false)
	require.NoError(t, err)
	require.True(t, seg2.IsSealed())

	require.Equal(t, sealTime.UnixNano(), seg2.GetSealTime().UnixNano())

	for k, addr := range addressMap {
		readValue, err := seg2.Read([]byte(k), addr)
		require.NoError(t, err)
		require.Equal(t, expectedValues[k], readValue)
	}

	keysFromSegment2, err := seg2.GetKeys()
	sort.Slice(keysFromSegment2, func(i, j int) bool {
		return bytes.Compare(keysFromSegment2[i].Key, keysFromSegment2[j].Key) < 0
	})
	require.NoError(t, err)
	require.Equal(t, keysFromSegment, keysFromSegment2)

	// delete the segment
	require.Equal(t, int(2+shardCount), countFilesInDirectory(t, directory))

	err = seg.delete()
	require.NoError(t, err)

	require.Equal(t, 0, countFilesInDirectory(t, directory))
}
