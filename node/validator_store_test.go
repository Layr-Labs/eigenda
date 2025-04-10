package node

import (
	"context"
	"os"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func randomInsertionsTest(t *testing.T, config *Config) {
	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config.DbPath = testDir

	store, err := NewValidatorStore(context.Background(), logger, config, time.Now, 2*time.Hour, nil)
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

		batchHeaderHash := rand.PrintableBytes(32)

		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
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

func TestRandomInsertions(t *testing.T) {
	t.Run("levelDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled:             false,
			ExpirationPollIntervalSec: 1,
		}
		randomInsertionsTest(t, config)
	})
	t.Run("littDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled: true,
		}
		randomInsertionsTest(t, config)
	})
}

func restartTest(t *testing.T, config *Config) {
	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config.DbPath = testDir

	store, err := NewValidatorStore(context.Background(), logger, config, time.Now, 2*time.Hour, nil)
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

		batchHeaderHash := rand.PrintableBytes(32)

		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
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
	store, err = NewValidatorStore(context.Background(), logger, config, time.Now, 2*time.Hour, nil)
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

func TestRestart(t *testing.T) {
	t.Run("levelDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled:             false,
			ExpirationPollIntervalSec: 1,
		}
		restartTest(t, config)
	})
	t.Run("littDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled: true,
		}
		restartTest(t, config)
	})
}

func doubleInsertionTest(t *testing.T, config *Config) {
	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config.DbPath = testDir

	store, err := NewValidatorStore(context.Background(), logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// A map from bundle key to bundle bytes
	expectedData := make(map[string][]byte)

	batchHeaderHashes := make([][]byte, 0, iterations)

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

		batchHeaderHash := rand.PrintableBytes(32)
		batchHeaderHashes = append(batchHeaderHashes, batchHeaderHash)

		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
		require.NoError(t, err)
	}

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// Attempt to insert data with the same batch header hash
	for _, batchHeaderHash := range batchHeaderHashes {
		bundles := make([]*BundleToStore, 0, 1)
		bundleKey := rand.PrintableBytes(32)
		bundleBytes := rand.PrintableVariableBytes(1, 64)
		bundles = append(bundles, &BundleToStore{
			BundleKey:   bundleKey,
			BundleBytes: bundleBytes,
		})
		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
		require.Error(t, err)
	}

	// Restarting should not permit double insertion.
	err = store.Stop()
	require.NoError(t, err)

	store, err = NewValidatorStore(context.Background(), logger, config, time.Now, 2*time.Hour, nil)
	require.NoError(t, err)

	// Read data back from the store
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// Attempt to insert data with the same batch header hash
	for _, batchHeaderHash := range batchHeaderHashes {
		bundles := make([]*BundleToStore, 0, 1)
		bundleKey := rand.PrintableBytes(32)
		bundleBytes := rand.PrintableVariableBytes(1, 64)
		bundles = append(bundles, &BundleToStore{
			BundleKey:   bundleKey,
			BundleBytes: bundleBytes,
		})
		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
		require.Error(t, err)
	}

	err = store.Stop()
	require.NoError(t, err)
}

func TestDoubleInsertion(t *testing.T) {
	t.Run("levelDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled:             false,
			ExpirationPollIntervalSec: 1,
		}
		doubleInsertionTest(t, config)
	})
	t.Run("littDB", func(t *testing.T) {
		config := &Config{
			LittDBEnabled: true,
		}
		doubleInsertionTest(t, config)
	})
}

func TestMigration(t *testing.T) {
	rand := random.NewTestRandom()
	testDir := t.TempDir()

	iterations := 10

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	config := &Config{
		LittDBEnabled:             false,
		ExpirationPollIntervalSec: 1,
		DbPath:                    testDir,
	}

	ttl := 2 * time.Hour

	store, err := NewValidatorStore(context.Background(), logger, config, time.Now, ttl, nil)
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

		batchHeaderHash := rand.PrintableBytes(32)

		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
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

	// Restart the store in migration mode.
	config.LittDBEnabled = true

	store, err = NewValidatorStore(context.Background(), logger, config, time.Now, ttl, nil)
	require.NoError(t, err)

	// Write some new data.
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

		batchHeaderHash := rand.PrintableBytes(32)

		_, _, err = store.StoreBatch(batchHeaderHash, bundles)
		require.NoError(t, err)
	}

	// Verify all data is present. Some data will be in littDB, other data will be in levelDB.
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// Restart the store while still in migration mode.
	// This invokes a different code pathway than the initial migration.
	err = store.Stop()
	require.NoError(t, err)

	timeSourceDelta := atomic.Uint64{}
	timeSource := func() time.Time {
		return time.Now().Add(time.Duration(timeSourceDelta.Load()) * time.Second)
	}

	store, err = NewValidatorStore(context.Background(), logger, config, timeSource, ttl, nil)
	require.NoError(t, err)

	// Verify all data is present. Some data will be in littDB, other data will be in levelDB.
	for bundleKey, expectedBundleBytes := range expectedData {
		bundleBytes, err := store.GetBundleData([]byte(bundleKey))
		require.NoError(t, err)
		require.Equal(t, expectedBundleBytes, bundleBytes)
	}

	// At this point in time, the levelDB directory should still be present.
	_, err = os.Stat(path.Join(testDir, LevelDBPath))
	require.NoError(t, err)

	// Simulate time moving forward by a few hours and manually trigger the step to clean up levelDB data.
	timeSourceDelta.Store(3 * 60 * 60)
	store.(*validatorStore).finalizeMigration(context.Background())

	// The levelDB directory should now be gone.
	_, err = os.Stat(path.Join(testDir, LevelDBPath))
	require.True(t, os.IsNotExist(err))

	err = store.Stop()
	require.NoError(t, err)
}
