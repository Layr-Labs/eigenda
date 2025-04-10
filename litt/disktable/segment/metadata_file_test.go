package segment

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func TestUnsealedSerialization(t *testing.T) {
	t.Parallel()
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	shardingFactor := rand.Uint32()
	salt := rand.Uint32()
	timestamp := rand.Uint64()
	m := &metadataFile{
		index:                index,
		serializationVersion: currentSerializationVersion,
		shardingFactor:       shardingFactor,
		salt:                 salt,
		timestamp:            timestamp,
		sealed:               false,
		parentDirectory:      directory,
	}
	err := m.write()
	require.NoError(t, err)

	deserialized, err := loadMetadataFile(index, []string{m.parentDirectory})
	require.NoError(t, err)
	require.Equal(t, *m, *deserialized)

	reportedSize := m.Size()
	stat, err := os.Stat(m.path())
	require.NoError(t, err)
	actualSize := uint64(stat.Size())
	require.Equal(t, actualSize, reportedSize)

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
	t.Parallel()
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	shardingFactor := rand.Uint32()
	salt := rand.Uint32()
	timestamp := rand.Uint64()
	m := &metadataFile{
		index:                index,
		serializationVersion: currentSerializationVersion,
		shardingFactor:       shardingFactor,
		salt:                 salt,
		timestamp:            timestamp,
		sealed:               true,
		parentDirectory:      directory,
	}
	err := m.write()
	require.NoError(t, err)

	reportedSize := m.Size()
	stat, err := os.Stat(m.path())
	require.NoError(t, err)
	actualSize := uint64(stat.Size())
	require.Equal(t, actualSize, reportedSize)

	deserialized, err := loadMetadataFile(index, []string{m.parentDirectory})
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
	t.Parallel()
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	m, err := createMetadataFile(index, 1234, 5678, directory)
	require.NoError(t, err)

	require.Equal(t, index, m.index)
	require.Equal(t, currentSerializationVersion, m.serializationVersion)
	require.False(t, m.sealed)
	require.Zero(t, m.timestamp)
	require.Equal(t, directory, m.parentDirectory)

	reportedSize := m.Size()
	stat, err := os.Stat(m.path())
	require.NoError(t, err)
	actualSize := uint64(stat.Size())
	require.Equal(t, actualSize, reportedSize)

	deserialized, err := loadMetadataFile(index, []string{m.parentDirectory})
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
	t.Parallel()
	rand := random.NewTestRandom()
	directory := t.TempDir()

	index := rand.Uint32()
	m, err := createMetadataFile(index, 1234, 5678, directory)
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
	deserialized, err := loadMetadataFile(index, []string{m.parentDirectory})
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
