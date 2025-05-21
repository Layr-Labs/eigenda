package benchmark

import (
	"context"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

func TestTrackerDeterminism(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	config := DefaultBenchmarkConfig()
	config.RandomPoolSize = units.MiB
	config.CohortSize = rand.Uint64Range(10, 20)
	config.MetadataDirectory = directory
	config.Seed = rand.Int63()
	config.ValueSizeMB = 1.0 / 1024 // 1kb

	// Generate enough data to fill 10ish cohorts.
	keyCount := 10*config.CohortSize + rand.Uint64Range(0, 10)

	dataTracker, err := NewDataTracker(context.Background(), config)
	require.NoError(t, err)

	// map from indices to keys
	expectedKeys := make(map[uint64][]byte)

	// map from indices to values
	expectedValues := make(map[uint64][]byte)

	// Get a bunch of values.
	for i := uint64(0); i < keyCount; i++ {
		writeInfo := dataTracker.GetWriteInfo()
		require.Equal(t, i, writeInfo.Index)
		require.Equal(t, 32, len(writeInfo.Key))
		require.Equal(t, units.KiB, len(writeInfo.Value))

		expectedKeys[i] = writeInfo.Key
		expectedValues[i] = writeInfo.Value
	}

	dataTracker.Close()

	// Rebuild the tracker at genesis. We should get the same sequence of keys and values.
	err = os.RemoveAll(directory)
	require.NoError(t, err)
	err = os.MkdirAll(directory, os.ModePerm)
	require.NoError(t, err)

	dataTracker, err = NewDataTracker(context.Background(), config)
	require.NoError(t, err)

	for i := uint64(0); i < keyCount; i++ {
		writeInfo := dataTracker.GetWriteInfo()
		require.Equal(t, i, writeInfo.Index)
		require.Equal(t, 32, len(writeInfo.Key))
		require.Equal(t, units.KiB, len(writeInfo.Value))
		require.Equal(t, expectedKeys[i], writeInfo.Key)
		require.Equal(t, expectedValues[i], writeInfo.Value)
	}

	dataTracker.Close()

	err = os.RemoveAll(directory)
	require.NoError(t, err)
}

func TestTrackerRestart(t *testing.T) {
	rand := random.NewTestRandom()
	directory := t.TempDir()

	config := DefaultBenchmarkConfig()
	config.RandomPoolSize = units.MiB
	config.CohortSize = rand.Uint64Range(10, 20)
	config.MetadataDirectory = directory
	config.Seed = rand.Int63()
	config.ValueSizeMB = 1.0 / 1024 // 1kb

	// Generate enough data to fill 10ish cohorts.
	keyCount := 10*config.CohortSize + rand.Uint64Range(0, 10)

	dataTracker, err := NewDataTracker(context.Background(), config)
	require.NoError(t, err)

	indexSet := make(map[uint64]struct{})

	// Generate a bunch of values.
	for i := uint64(0); i < keyCount; i++ {
		writeInfo := dataTracker.GetWriteInfo()
		require.Equal(t, i, writeInfo.Index)
		require.Equal(t, 32, len(writeInfo.Key))
		require.Equal(t, units.KiB, len(writeInfo.Value))

		indexSet[writeInfo.Index] = struct{}{}
	}

	// All indices should be unique.
	require.Equal(t, keyCount, uint64(len(indexSet)))

	// Restart.
	dataTracker.Close()
	dataTracker, err = NewDataTracker(context.Background(), config)
	require.NoError(t, err)

	// Generate more values.
	for i := uint64(0); i < keyCount; i++ {
		writeInfo := dataTracker.GetWriteInfo()
		indexSet[writeInfo.Index] = struct{}{}
	}

	// If we aren't reusing indices after the restart, then the set should now be equal to 2*keyCount.
	require.Equal(t, 2*keyCount, uint64(len(indexSet)))

	dataTracker.Close()

	err = os.RemoveAll(directory)
	require.NoError(t, err)
}

func TestTrackReads(t *testing.T) {

}

// TODO: test read expirations somehow
