package segment

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"testing"
)

// countFilesInDirectory returns the number of files in the given directory.
func countFilesInDirectory(t *testing.T, directory string) int {
	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	return len(files)
}

func TestWriteAndReadSegment(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(100, 200)
	values := make([][]byte, valueCount)
	keys := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		values[i] = rand.VariableBytes(1, 100)
		keys[i] = rand.VariableBytes(1, 100)
	}

	addressMap := make(map[Address][]byte)

	seg, err := NewSegment(logger, index, directory, math.MaxUint32)
	require.NoError(t, err)

	// Write values to the segment.
	for i := 0; i < int(valueCount); i++ {
		key := keys[i]
		value := values[i]

		address, ok, err := seg.Write(key, value)
		require.NoError(t, err)
		require.True(t, ok)
		addressMap[address] = value

		// Occasionally flush the segment to disk.
		if rand.BoolWithProbability(0.25) {
			err := seg.Flush()
			require.NoError(t, err)
		}

		// Occasionally scan all addresses and values in the segment. Some will be on disk, some will be in memory.
		if rand.BoolWithProbability(0.1) {
			for addr, val := range addressMap {
				readValue, err := seg.Read(addr)
				require.NoError(t, err)
				require.Equal(t, val, readValue)
			}
		}
	}

	// Seal the segment and read all keys and values.
	require.False(t, seg.IsSealed())
	sealTime := rand.Time()
	err = seg.Seal(sealTime)
	require.NoError(t, err)
	require.True(t, seg.IsSealed())

	require.Equal(t, sealTime.UnixNano(), seg.GetSealTime().UnixNano())

	for addr, val := range addressMap {
		readValue, err := seg.Read(addr)
		require.NoError(t, err)
		require.Equal(t, val, readValue)
	}

	keysFromSegment, err := seg.GetKeys()
	require.NoError(t, err)
	require.Equal(t, keys, keysFromSegment)

	expectedSize := uint64(0)
	for _, value := range values {
		expectedSize += uint64(len(value)) + 4
	}
	require.Equal(t, expectedSize, seg.CurrentSize())

	// Reopen the segment and read all keys and values.
	seg2, err := NewSegment(logger, index, directory, math.MaxUint32)
	require.NoError(t, err)
	require.True(t, seg2.IsSealed())

	require.Equal(t, sealTime.UnixNano(), seg2.GetSealTime().UnixNano())

	for addr, val := range addressMap {
		readValue, err := seg2.Read(addr)
		require.NoError(t, err)
		require.Equal(t, val, readValue)
	}

	keysFromSegment2, err := seg2.GetKeys()
	require.NoError(t, err)
	require.Equal(t, keys, keysFromSegment2)

	require.Equal(t, expectedSize, seg2.CurrentSize())

	// delete the segment
	require.Equal(t, 3, countFilesInDirectory(t, directory))

	err = seg.delete()
	require.NoError(t, err)

	require.Equal(t, 0, countFilesInDirectory(t, directory))
}

func TestWriteToFullSegment(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	sizeOfAllButLastValue := uint64(0)
	sizeOfLastValue := uint64(0)

	index := rand.Uint32()
	valueCount := rand.Int32Range(100, 200)
	values := make([][]byte, valueCount)
	keys := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		values[i] = rand.VariableBytes(1, 100)
		keys[i] = rand.VariableBytes(1, 100)

		if i < int(valueCount)-1 {
			sizeOfAllButLastValue += uint64(len(values[i])) + 4
		} else {
			sizeOfLastValue = uint64(len(values[i])) + 4
		}
	}

	addressMap := make(map[Address][]byte)

	capacity := rand.Uint32Range(uint32(sizeOfAllButLastValue), uint32(sizeOfAllButLastValue+sizeOfLastValue))
	seg, err := NewSegment(logger, index, directory, capacity)
	require.NoError(t, err)

	// Write the values. All but the last one should fit.
	for i := 0; i < int(valueCount)-1; i++ {
		key := keys[i]
		value := values[i]

		address, ok, err := seg.Write(key, value)
		require.NoError(t, err)
		require.True(t, ok)
		addressMap[address] = value
	}

	// The last value should not fit.
	key := keys[int(valueCount)-1]
	value := values[int(valueCount)-1]
	_, ok, err := seg.Write(key, value)
	require.NoError(t, err)
	require.False(t, ok)

	// Read the values back. All but the last one should be there.
	require.False(t, seg.IsSealed())
	sealTime := rand.Time()
	err = seg.Seal(sealTime)
	require.NoError(t, err)

	for addr, val := range addressMap {
		readValue, err := seg.Read(addr)
		require.NoError(t, err)
		require.Equal(t, val, readValue)
	}

	keysFromSegment, err := seg.GetKeys()
	require.NoError(t, err)
	require.Equal(t, keys[:len(keys)-1], keysFromSegment)

	// Measure the value file size on disk. It should exactly match sizeOfAllButLastValue.
	require.Equal(t, sizeOfAllButLastValue, seg.CurrentSize())
	bytesFromFile, err := os.ReadFile(seg.values.path())
	require.NoError(t, err)
	require.Equal(t, int(sizeOfAllButLastValue), len(bytesFromFile))

	// Delete the segment.
	require.Equal(t, 3, countFilesInDirectory(t, directory))
	err = seg.delete()
	require.NoError(t, err)
	require.Equal(t, 0, countFilesInDirectory(t, directory))
}

// TODO reservation tests
// TODO future cody: start here, then finish garbage collection on the segment manager, then test segment manager
