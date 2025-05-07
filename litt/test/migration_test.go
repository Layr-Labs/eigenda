package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/stretchr/testify/require"
)

// This file contains tests for data migrations (i.e. when the on-disk format of the data changes).

// Enable and run this "test" to generate data for a migration test at the current version.
func TestGenerateData(t *testing.T) {
	t.Skip() // comment out this line to generate data

	version := segment.CurrentSerializationVersion
	dataDir := fmt.Sprintf("testdata/v%d", version)

	exists, err := util.Exists(dataDir)
	require.NoError(t, err)
	if exists {
		fmt.Printf("deleting existing data at %s\n", dataDir)
		err = os.RemoveAll(dataDir)
		require.NoError(t, err)
	}

	fmt.Printf("generating migration test data at %s\n", dataDir)

	err = os.MkdirAll(dataDir, 0777)
	require.NoError(t, err)

	config, err := litt.DefaultConfig(dataDir)
	require.NoError(t, err)
	config.DoubleWriteProtection = true
	config.Fsync = false
	config.ShardingFactor = 4
	config.TargetSegmentFileSize = 100

	db, err := littbuilder.NewDB(config)
	require.NoError(t, err)

	table, err := db.GetTable("test")
	require.NoError(t, err)

	for key, value := range migrationData {
		err = table.Put([]byte(key), []byte(value))
		require.NoError(t, err)
	}

	// verify the data in the table
	for key, value := range migrationData {
		v, exists, err := table.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, value, string(v))
	}

	// Shut the DB down.
	err = db.Close()
	require.NoError(t, err)
}

// A "control experiment" for migration. Write the expected data and reload it with the same version.
// This is a sanity check on the migration test.
func TestMigrationControl(t *testing.T) {
	dataDir := t.TempDir()

	// Write the data to a table at the old version.
	config, err := litt.DefaultConfig(dataDir)
	require.NoError(t, err)
	config.DoubleWriteProtection = true
	config.Fsync = false

	db, err := littbuilder.NewDB(config)
	require.NoError(t, err)

	table, err := db.GetTable("test")
	require.NoError(t, err)

	for key, value := range migrationData {
		err = table.Put([]byte(key), []byte(value))
		require.NoError(t, err)
	}

	// verify the data in the table
	for key, value := range migrationData {
		v, exists, err := table.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, value, string(v))
	}

	// Shut the DB down.
	err = db.Close()
	require.NoError(t, err)

	// Reload the DB.
	config, err = litt.DefaultConfig(dataDir)
	require.NoError(t, err)
	config.DoubleWriteProtection = true
	config.Fsync = false
	db, err = littbuilder.NewDB(config)
	require.NoError(t, err)

	table, err = db.GetTable("test")
	require.NoError(t, err)

	// verify the data in the table
	for key, value := range migrationData {
		v, exists, err := table.Get([]byte(key))
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, value, string(v))
	}

	// This is just a sanity check test, delete the table.
	err = db.Destroy()
	require.NoError(t, err)
}
