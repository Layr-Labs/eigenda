package segment

import (
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/stretchr/testify/require"
)

// countFilesInDirectory returns the number of files in the given directory.
func countFilesInDirectory(t *testing.T, directory string) int {
	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	return len(files)
}

func TestWriteAndReadSegment(t *testing.T) {
	rand := random.NewTestRandom()
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(100, 200)
	keys := make([]*types.KAPair, valueCount)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		key := rand.VariableBytes(1, 100)
		keys[i] = &types.KAPair{Key: key}
		values[i] = rand.VariableBytes(1, 100)
	}

	addressMap := make(map[types.Address][]byte)

	seg, err := NewSegment(logger, index, directory, time.Now(), false)
	require.NoError(t, err)

	// Write values to the segment.
	for i := 0; i < int(valueCount); i++ {
		key := keys[i]
		value := values[i]

		key.Address, err = seg.Write(key.Key, value)
		require.NoError(t, err)
		addressMap[key.Address] = value

		// Occasionally flush the segment to disk.
		if rand.BoolWithProbability(0.25) {
			err := seg.Flush()
			require.NoError(t, err)
		}

		// Occasionally scan all addresses and values in the segment.
		if rand.BoolWithProbability(0.1) {
			err := seg.Flush()
			require.NoError(t, err)
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
	seg2, err := NewSegment(logger, index, directory, time.Now(), false)
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
