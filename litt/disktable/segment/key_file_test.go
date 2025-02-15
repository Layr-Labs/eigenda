package segment

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWriteThenRead(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	assert.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()

	keyCount := rand.Intn(100) + 1
	keys := make([][]byte, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = rand.VariableBytes(1, 100)
	}

	file, err := newKeyFile(logger, index, directory, false)
	require.NoError(t, err)

	for _, key := range keys {
		err := file.write(key)
		require.NoError(t, err)
	}

	// Reading the file prior to sealing it is forbidden.
	_, err = file.readKeys()
	require.Error(t, err)

	err = file.seal()
	require.NoError(t, err)

	// Reading the file after sealing it is allowed.
	readKeys, err := file.readKeys()
	require.NoError(t, err)

	for i, key := range keys {
		assert.Equal(t, key, readKeys[i])
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

func TestReadingTruncatedFile(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	assert.NoError(t, err)
	directory := t.TempDir()

	index := rand.Uint32()

	keyCount := rand.Intn(100) + 1
	keys := make([][]byte, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = rand.VariableBytes(1, 100)
	}

	file, err := newKeyFile(logger, index, directory, false)
	require.NoError(t, err)

	for _, key := range keys {
		err := file.write(key)
		require.NoError(t, err)
	}

	err = file.seal()
	require.NoError(t, err)

	// Truncate the file. Chop off some bytes from the last key, but do not corrupt the length prefix.
	lastKeyLength := len(keys[keyCount-1])

	filePath := file.path()

	originalBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

	bytesToRemove := rand.Intn(lastKeyLength-1) + 1
	bytes := originalBytes[:len(originalBytes)-bytesToRemove]

	err = os.WriteFile(filePath, bytes, 0644)
	require.NoError(t, err)

	// We should be able to read the keys up to the point where the file was truncated.
	readKeys, err := file.readKeys()
	require.NoError(t, err)

	require.Equal(t, keyCount-1, len(readKeys))
	for i, key := range keys[:keyCount-1] {
		assert.Equal(t, key, readKeys[i])
	}

	// Truncate the file. This time, chop off some of the length prefix of the last key.
	prefixBytesToRemove := rand.Intn(2) + 1
	bytes = originalBytes[:len(originalBytes)-prefixBytesToRemove]

	err = os.WriteFile(filePath, bytes, 0644)
	require.NoError(t, err)

	// We should not be able to read the keys if the length prefix is truncated.
	keys, err = file.readKeys()
	require.NoError(t, err)

	require.Equal(t, keyCount-1, len(keys))
	for i, key := range keys[:keyCount-1] {
		assert.Equal(t, key, keys[i])
	}

	// delete the file
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = file.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}
