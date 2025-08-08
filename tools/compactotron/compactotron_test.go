package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestCompaction(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	rand := random.NewTestRandom()

	testDir := t.TempDir()
	source := path.Join(testDir, "source")
	destination := path.Join(testDir, "destination")

	db, err := leveldb.OpenFile(source, nil)
	require.NoError(t, err)

	fmt.Printf("writing values into original table\n")
	expectedValues := make(map[string][]byte)
	for i := 0; i < 1024; i++ {
		key := rand.String(32)
		value := rand.PrintableBytes(units.MiB)

		expectedValues[key] = value
		err = db.Put([]byte(key), value, nil)
	}

	err = db.Close()
	require.NoError(t, err)

	fmt.Printf("doing migration\n")
	err = CompactLevelDB(source, destination)
	require.NoError(t, err)

	fmt.Printf("opening compacted table and comparing it to the original\n")
	db, err = leveldb.OpenFile(destination, nil)
	require.NoError(t, err)

	for key, expectedValue := range expectedValues {
		actualValue, err := db.Get([]byte(key), nil)
		require.NoError(t, err)
		require.Equal(t, expectedValue, actualValue, "value for key %s does not match", key)
	}

	err = db.Close()
	require.NoError(t, err)

	fmt.Printf("opening original table to check if it is still intact\n")
	db, err = leveldb.OpenFile(source, nil)
	require.NoError(t, err)

	for key, expectedValue := range expectedValues {
		actualValue, err := db.Get([]byte(key), nil)
		require.NoError(t, err)
		require.Equal(t, expectedValue, actualValue, "value for key %s does not match in original table", key)
	}

	err = db.Close()
	require.NoError(t, err)
}
