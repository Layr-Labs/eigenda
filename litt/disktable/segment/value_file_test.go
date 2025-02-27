package segment

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/stretchr/testify/require"
)

func TestWriteThenReadValues(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(100, 200)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		values[i] = rand.VariableBytes(1, 100)
	}

	addressMap := make(map[types.Address][]byte)

	file, err := newValueFile(logger, index, directory, false)
	require.NoError(t, err)

	for _, value := range values {
		address, err := file.write(value)
		require.NoError(t, err)
		addressMap[address] = value

		// Occasionally flush the file to disk.
		if rand.BoolWithProbability(0.25) {
			err := file.flush()
			require.NoError(t, err)
		}

		// Occasionally scan all addresses and values in the file.
		if rand.BoolWithProbability(0.1) {
			err = file.flush()
			require.NoError(t, err)
			for key, val := range addressMap {
				readValue, err := file.read(key)
				require.NoError(t, err)
				require.Equal(t, val, readValue)
			}
		}
	}

	// Seal the file and read all values.
	err = file.seal()
	require.NoError(t, err)
	for key, val := range addressMap {
		readValue, err := file.read(key)
		require.NoError(t, err)
		require.Equal(t, val, readValue)
	}

	// delete the file
	filePath := file.path()
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = file.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}

func TestReadingTruncatedValueFile(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()
	valueCount := rand.Int32Range(100, 200)
	values := make([][]byte, valueCount)
	for i := 0; i < int(valueCount); i++ {
		values[i] = rand.VariableBytes(1, 100)
	}

	addressMap := make(map[types.Address][]byte)

	file, err := newValueFile(logger, index, directory, false)
	require.NoError(t, err)

	var lastAddress types.Address
	for _, value := range values {
		address, err := file.write(value)
		require.NoError(t, err)
		addressMap[address] = value
		lastAddress = address
	}

	err = file.seal()
	require.NoError(t, err)

	// Truncate the file. Chop off some bytes from the last value, but do not corrupt the length prefix.
	lastValueLength := len(values[valueCount-1])

	filePath := file.path()

	originalBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

	bytesToRemove := rand.Int32Range(1, int32(lastValueLength)+1)
	bytes := originalBytes[:len(originalBytes)-int(bytesToRemove)]

	err = os.WriteFile(filePath, bytes, 0644)
	require.NoError(t, err)

	file, err = newValueFile(logger, index, directory, true)
	require.NoError(t, err)

	// We should be able to read all values except for the last one.
	for key, val := range addressMap {
		if key == lastAddress {
			_, err := file.read(key)
			require.Error(t, err)
		} else {
			readValue, err := file.read(key)
			require.NoError(t, err)
			require.Equal(t, val, readValue)
		}
	}

	// Truncate the file. Corrupt the length prefix of the last value.
	prefixBytesToRemove := rand.Int32Range(1, 4)
	bytes = originalBytes[:len(originalBytes)-int(prefixBytesToRemove)]

	err = os.WriteFile(filePath, bytes, 0644)
	require.NoError(t, err)

	file, err = newValueFile(logger, index, directory, true)
	require.NoError(t, err)

	// We should be able to read all values except for the last one.
	for key, val := range addressMap {
		if key == lastAddress {
			_, err := file.read(key)
			require.Error(t, err)
		} else {
			readValue, err := file.read(key)
			require.NoError(t, err)
			require.Equal(t, val, readValue)
		}
	}

	// delete the file
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = file.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}
