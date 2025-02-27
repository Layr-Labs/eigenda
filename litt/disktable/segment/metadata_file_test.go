package segment

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUnsealedSerialization(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	timestamp := rand.Uint64()
	m := &metadataFile{
		index:                index,
		serializationVersion: currentSerializationVersion,
		sealed:               false,
		timestamp:            timestamp,
		parentDirectory:      directory,
	}
	err := m.write()
	require.NoError(t, err)

	deserialized, err := newMetadataFile(index, m.parentDirectory)
	require.NoError(t, err)
	require.Equal(t, *m, *deserialized)

	// delete the file
	filePath := m.path()
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = m.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}

func TestSealedSerialization(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	timestamp := rand.Uint64()
	m := &metadataFile{
		index:                index,
		serializationVersion: currentSerializationVersion,
		sealed:               true,
		timestamp:            timestamp,
		parentDirectory:      directory,
	}
	err := m.write()
	require.NoError(t, err)

	deserialized, err := newMetadataFile(index, m.parentDirectory)
	require.NoError(t, err)
	require.Equal(t, *m, *deserialized)

	// delete the file
	filePath := m.path()
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = m.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}

func TestFreshFileSerialization(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	m, err := newMetadataFile(index, directory)
	require.NoError(t, err)

	require.Equal(t, index, m.index)
	require.Equal(t, currentSerializationVersion, m.serializationVersion)
	require.False(t, m.sealed)
	require.Zero(t, m.timestamp)
	require.Equal(t, directory, m.parentDirectory)

	deserialized, err := newMetadataFile(index, m.parentDirectory)
	require.NoError(t, err)
	require.Equal(t, *m, *deserialized)

	// delete the file
	filePath := m.path()
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = m.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}

func TestSealing(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	m, err := newMetadataFile(index, directory)
	require.NoError(t, err)

	// seal the file
	sealTime := rand.Time()
	err = m.seal(sealTime)
	require.NoError(t, err)

	require.Equal(t, index, m.index)
	require.Equal(t, currentSerializationVersion, m.serializationVersion)
	require.True(t, m.sealed)
	require.Equal(t, uint64(sealTime.UnixNano()), m.timestamp)
	require.Equal(t, directory, m.parentDirectory)

	// load the file
	deserialized, err := newMetadataFile(index, m.parentDirectory)
	require.NoError(t, err)
	require.Equal(t, *m, *deserialized)

	// delete the file
	filePath := m.path()
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	err = m.delete()
	require.NoError(t, err)

	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err))
}
