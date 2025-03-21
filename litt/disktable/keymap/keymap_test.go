package keymap

import (
	"os"
	"path"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var builders = []keymapBuilder{
	buildMemKeymap,
	buildLevelDBKeymap,
}

type keymapBuilder func(logger logging.Logger, path string) (Keymap, error)

func buildMemKeymap(logger logging.Logger, path string) (Keymap, error) {
	return NewMemKeymap(logger, true), nil
}

func buildLevelDBKeymap(logger logging.Logger, path string) (Keymap, error) {
	return NewLevelDBKeymap(logger, path, true)
}

func testBasicBehavior(t *testing.T, keymap Keymap) {
	rand := random.NewTestRandom()

	expected := make(map[string]types.Address)

	operations := 1000
	for i := 0; i < operations; i++ {
		choice := rand.Float64()
		if choice < 0.5 {
			// Write a random value
			key := []byte(rand.String(32))
			address := types.Address(rand.Uint64())

			err := keymap.Put([]*types.KAPair{{Key: key, Address: address}})
			require.NoError(t, err)
			expected[string(key)] = address
		} else if choice < 0.75 {
			// Delete a few random values
			numberToDelete := rand.Int32Range(1, 10)
			numberToDelete = min(numberToDelete, int32(len(expected)))
			keysToDelete := make([]*types.KAPair, 0, numberToDelete)
			for key := range expected {
				if numberToDelete == int32(len(keysToDelete)) {
					break
				}
				keysToDelete = append(keysToDelete, &types.KAPair{Key: []byte(key)})
				numberToDelete--
			}

			err := keymap.Delete(keysToDelete)
			require.NoError(t, err)
			for _, key := range keysToDelete {
				delete(expected, string(key.Key))
			}
		} else {
			// Write a batch of random values
			numberToWrite := rand.Int32Range(1, 10)
			pairs := make([]*types.KAPair, numberToWrite)
			for i := 0; i < int(numberToWrite); i++ {
				key := []byte(rand.String(32))
				address := types.Address(rand.Uint64())
				pairs[i] = &types.KAPair{Key: key, Address: address}
				expected[string(key)] = address
			}
			err := keymap.Put(pairs)
			require.NoError(t, err)
		}

		// Every once in a while, verify that the keymap is correct
		if rand.BoolWithProbability(0.1) {
			for key, expectedAddress := range expected {
				address, ok, err := keymap.Get([]byte(key))
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, expectedAddress, address)
			}
		}
	}

	for key, expectedAddress := range expected {
		address, ok, err := keymap.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, expectedAddress, address)
	}

	err := keymap.Destroy()
	require.NoError(t, err)
}

func TestBasicBehavior(t *testing.T) {
	testDir := t.TempDir()
	dbDir := path.Join(testDir, "keymap")

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	for _, builder := range builders {
		keymap, err := builder(logger, dbDir)
		if err != nil {
			t.Fatalf("Failed to create keymap: %v", err)
		}
		testBasicBehavior(t, keymap)

		// verify that test dir is empty (destroy should have deleted everything)
		_, err = os.Stat(dbDir)
		if err != nil {
			if !os.IsNotExist(err) {
				require.NoError(t, err)
			}

			// Directory doesn't exist. We are good.
		} else {
			// Directory exists. Make sure it's emtpy.
			entries, err := os.ReadDir(dbDir)
			require.NoError(t, err)
			require.Empty(t, entries)
		}
	}
}

func TestRestart(t *testing.T) {
	rand := random.NewTestRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	dbDir := path.Join(testDir, "keymap")

	keymap, err := NewLevelDBKeymap(logger, dbDir, true)
	require.NoError(t, err)

	expected := make(map[string]types.Address)

	operations := 1000
	for i := 0; i < operations; i++ {
		choice := rand.Float64()
		if choice < 0.5 {
			// Write a random value
			key := []byte(rand.String(32))
			address := types.Address(rand.Uint64())

			err := keymap.Put([]*types.KAPair{{Key: key, Address: address}})
			require.NoError(t, err)
			expected[string(key)] = address
		} else if choice < 0.75 {
			// Delete a few random values
			numberToDelete := rand.Int32Range(1, 10)
			numberToDelete = min(numberToDelete, int32(len(expected)))
			keysToDelete := make([]*types.KAPair, 0, numberToDelete)
			for key := range expected {
				if numberToDelete == int32(len(keysToDelete)) {
					break
				}
				keysToDelete = append(keysToDelete, &types.KAPair{Key: []byte(key)})
				numberToDelete--
			}

			err := keymap.Delete(keysToDelete)
			require.NoError(t, err)
			for _, key := range keysToDelete {
				delete(expected, string(key.Key))
			}
		} else {
			// Write a batch of random values
			numberToWrite := rand.Int32Range(1, 10)
			pairs := make([]*types.KAPair, numberToWrite)
			for i := 0; i < int(numberToWrite); i++ {
				key := []byte(rand.String(32))
				address := types.Address(rand.Uint64())
				pairs[i] = &types.KAPair{Key: key, Address: address}
				expected[string(key)] = address
			}
			err := keymap.Put(pairs)
			require.NoError(t, err)
		}

		// Every once in a while, verify that the keymap is correct
		if rand.BoolWithProbability(0.1) {
			for key, expectedAddress := range expected {
				address, ok, err := keymap.Get([]byte(key))
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, expectedAddress, address)
			}
		}
	}

	for key, expectedAddress := range expected {
		address, ok, err := keymap.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, expectedAddress, address)
	}

	// Shut down the keymap and reload it
	err = keymap.Stop()
	require.NoError(t, err)

	keymap, err = NewLevelDBKeymap(logger, dbDir, true)
	require.NoError(t, err)

	// Expected data should be present
	for key, expectedAddress := range expected {
		address, ok, err := keymap.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, expectedAddress, address)
	}

	for i := 0; i < operations; i++ {
		choice := rand.Float64()
		if choice < 0.5 {
			// Write a random value
			key := []byte(rand.String(32))
			address := types.Address(rand.Uint64())

			err := keymap.Put([]*types.KAPair{{Key: key, Address: address}})
			require.NoError(t, err)
			expected[string(key)] = address
		} else if choice < 0.75 {
			// Delete a few random values
			numberToDelete := rand.Int32Range(1, 10)
			numberToDelete = min(numberToDelete, int32(len(expected)))
			keysToDelete := make([]*types.KAPair, 0, numberToDelete)
			for key := range expected {
				if numberToDelete == int32(len(keysToDelete)) {
					break
				}
				keysToDelete = append(keysToDelete, &types.KAPair{Key: []byte(key)})
				numberToDelete--
			}

			err := keymap.Delete(keysToDelete)
			require.NoError(t, err)
			for _, key := range keysToDelete {
				delete(expected, string(key.Key))
			}
		} else {
			// Write a batch of random values
			numberToWrite := rand.Int32Range(1, 10)
			pairs := make([]*types.KAPair, numberToWrite)
			for i := 0; i < int(numberToWrite); i++ {
				key := []byte(rand.String(32))
				address := types.Address(rand.Uint64())
				pairs[i] = &types.KAPair{Key: key, Address: address}
				expected[string(key)] = address
			}
			err := keymap.Put(pairs)
			require.NoError(t, err)
		}

		// Every once in a while, verify that the keymap is correct
		if rand.BoolWithProbability(0.1) {
			for key, expectedAddress := range expected {
				address, ok, err := keymap.Get([]byte(key))
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, expectedAddress, address)
			}
		}
	}

	for key, expectedAddress := range expected {
		address, ok, err := keymap.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, expectedAddress, address)
	}

	err = keymap.Destroy()
	require.NoError(t, err)
}
