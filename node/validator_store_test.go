package node

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

func TestRandomInsertions(t *testing.T) {

	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config := &Config{
		GetChunksHotCacheReadLimitMB:  units.GiB,
		GetChunksHotBurstLimitMB:      units.GiB,
		GetChunksColdCacheReadLimitMB: units.GiB,
		GetChunksColdBurstLimitMB:     units.GiB,
		LittDBStoragePaths:            []string{testDir},
	}

	store, err := NewValidatorStore(logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// A map from bundle key to bundle bytes
	expectedData := make(map[string][]byte)

	// Write data to the store
	for i := 0; i < iterations; i++ {
		bundleCount := rand.Int32Range(1, 10)
		bundles := make([]*BundleToStore, 0, bundleCount)
		for j := 0; j < int(bundleCount); j++ {
			bundleKey := rand.PrintableBytes(32)
			bundleBytes := rand.PrintableVariableBytes(1, 64)
			bundles = append(bundles, &BundleToStore{
				BundleKey:   bundleKey,
				BundleBytes: bundleBytes,
			})
			expectedData[string(bundleKey)] = bundleBytes
		}

		_, err = store.StoreBatch(bundles)
		require.NoError(t, err)
	}

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	err = store.Stop()
	require.NoError(t, err)
}

func TestRestart(t *testing.T) {

	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config := &Config{
		GetChunksHotCacheReadLimitMB:  units.GiB,
		GetChunksHotBurstLimitMB:      units.GiB,
		GetChunksColdCacheReadLimitMB: units.GiB,
		GetChunksColdBurstLimitMB:     units.GiB,
		LittDBStoragePaths:            []string{testDir},
	}

	store, err := NewValidatorStore(logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// A map from bundle key to bundle bytes
	expectedData := make(map[string][]byte)

	// Write data to the store
	for i := 0; i < iterations; i++ {
		bundleCount := rand.Int32Range(1, 10)
		bundles := make([]*BundleToStore, 0, bundleCount)
		for j := 0; j < int(bundleCount); j++ {
			bundleKey := rand.PrintableBytes(32)
			bundleBytes := rand.PrintableVariableBytes(1, 64)
			bundles = append(bundles, &BundleToStore{
				BundleKey:   bundleKey,
				BundleBytes: bundleBytes,
			})
			expectedData[string(bundleKey)] = bundleBytes
		}

		_, err = store.StoreBatch(bundles)
		require.NoError(t, err)
	}

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	err = store.Stop()
	require.NoError(t, err)

	// Restart the store
	store, err = NewValidatorStore(logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	err = store.Stop()
	require.NoError(t, err)
}

func TestDoubleInsertionLittDB(t *testing.T) {
	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config := &Config{
		LittDBDoubleWriteProtection:   true, // causes littDB to throw if data is written twice
		GetChunksHotCacheReadLimitMB:  units.GiB,
		GetChunksHotBurstLimitMB:      units.GiB,
		GetChunksColdCacheReadLimitMB: units.GiB,
		GetChunksColdBurstLimitMB:     units.GiB,
		LittDBStoragePaths:            []string{testDir},
	}

	store, err := NewValidatorStore(logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// A map from bundle key to bundle bytes
	expectedData := make(map[string][]byte)

	// Write data to the store
	for i := 0; i < iterations; i++ {
		bundleCount := rand.Int32Range(1, 10)
		bundles := make([]*BundleToStore, 0, bundleCount)
		for j := 0; j < int(bundleCount); j++ {
			bundleKey := rand.PrintableBytes(32)
			bundleBytes := rand.PrintableVariableBytes(1, 64)
			bundles = append(bundles, &BundleToStore{
				BundleKey:   bundleKey,
				BundleBytes: bundleBytes,
			})
			expectedData[string(bundleKey)] = bundleBytes
		}

		_, err = store.StoreBatch(bundles)
		require.NoError(t, err)
	}

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// Attempt to insert the same data again
	for bundleKey, bundleBytes := range expectedData {
		bundles := make([]*BundleToStore, 0, 1)
		bundles = append(bundles, &BundleToStore{
			BundleKey:   []byte(bundleKey),
			BundleBytes: bundleBytes[:],
		})
		_, err = store.StoreBatch(bundles)
		require.NoError(t, err)
	}

	// Restart and try again.
	err = store.Stop()
	require.NoError(t, err)

	store, err = NewValidatorStore(logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// Attempt to insert the same data again
	for bundleKey, bundleBytes := range expectedData {
		bundles := make([]*BundleToStore, 0, 1)
		bundles = append(bundles, &BundleToStore{
			BundleKey:   []byte(bundleKey),
			BundleBytes: bundleBytes[:],
		})
		_, err = store.StoreBatch(bundles)
		require.NoError(t, err)
	}

	err = store.Stop()
	require.NoError(t, err)
}
