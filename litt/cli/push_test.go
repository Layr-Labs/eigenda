package main

import (
	"path"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/stretchr/testify/require"
)

func pushTest(
	t *testing.T,
	sourceDirs uint64,
	destDirs uint64,
	deleteAfterTransfer bool,
	verbose bool,
) {

	//logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	//require.NoError(t, err)

	rand := random.NewTestRandom()
	sourceRoot := t.TempDir()
	destRoot := t.TempDir()

	sourceDirList := make([]string, 0, sourceDirs)
	destDirList := make([]string, 0, destDirs)

	for i := uint64(0); i < sourceDirs; i++ {
		sourceDirList = append(sourceDirList, path.Join(sourceRoot, rand.String(32)))
	}
	for i := uint64(0); i < destDirs; i++ {
		destDirList = append(destDirList, path.Join(destRoot, rand.String(32)))
	}

	tableCount := rand.Uint64Range(2, 4)
	tableNames := make([]string, 0, tableCount)
	for i := uint64(0); i < tableCount; i++ {
		tableNames = append(tableNames, rand.String(32))
	}

	shardingFactor := sourceDirs + rand.Uint64Range(0, 4)

	config, err := litt.DefaultConfig(sourceDirList...)
	require.NoError(t, err)
	config.DoubleWriteProtection = true
	config.ShardingFactor = uint32(shardingFactor)
	config.Fsync = false

	db, err := littbuilder.NewDB(config)
	require.NoError(t, err)

	expectedData := make(map[string] /*table*/ map[string] /*value*/ []byte)
	for _, tableName := range tableNames {
		expectedData[tableName] = make(map[string][]byte)
	}

	// Insert data into the tables.
	keyCount := uint64(1024)
	for i := uint64(0); i < keyCount; i++ {
		tableIndex := rand.Uint64Range(0, tableCount)
		table, err := db.GetTable(tableNames[tableIndex])
		require.NoError(t, err)
		key := rand.PrintableBytes(32)
		value := rand.PrintableVariableBytes(10, 100)

		expectedData[table.Name()][string(key)] = value
		err = table.Put(key, value)
		require.NoError(t, err, "failed to put key %s in table %s", key, table.Name())
	}

	// Flush all tables.
	for _, tableName := range tableNames {
		table, err := db.GetTable(tableName)
		require.NoError(t, err)
		err = table.Flush()
		require.NoError(t, err, "failed to flush table %s", table.Name())
	}

	// Verify the data in the DB.
	for tableName := range expectedData {
		table, err := db.GetTable(tableName)
		require.NoError(t, err, "failed to get table %s", tableName)
		for key := range expectedData[tableName] {
			value, ok, err := table.Get([]byte(key))
			require.NoError(t, err, "failed to get key %s in table %s", key, tableName)
			require.True(t, ok, "key %s not found in table %s", key, tableName)
			require.Equal(t, expectedData[tableName][key], value,
				"value for key %s in table %s does not match expected value", key, tableName)
		}
	}

	// Verify expected directories.
	for _, sourceDir := range sourceDirList {
		// We should see each source dir.
		exists, err := util.Exists(sourceDir)
		require.NoError(t, err)
		require.True(t, exists, "source directory %s does not exist", sourceDir)
	}
	for _, destDir := range destDirList {
		// We should not see dest dirs yet.
		exists, err := util.Exists(destDir)
		require.NoError(t, err)
		require.False(t, exists, "destination directory %s exists", destDir)
	}

	// Start a container that is running an SSH server. The push() command will communicate with this server.
	container := util.SetupSSHTestContainer(t)
	defer func() { _ = container.Cleanup() }()

	//// pushing with the DB still open should fail.
	//err = push()
	//require.Error(t, err)
	//
	//// None of the source dirs should have been deleted.
	//for _, sourceDir := range sourceDirList {
	//	// We should see each source dir.
	//	exists, err := util.Exists(sourceDir)
	//	require.NoError(t, err)
	//	require.True(t, exists, "source directory %s does not exist", sourceDir)
	//}
	//
	//// The failed rebase should not have changed the data in the DB.
	//for tableName := range expectedData {
	//	table, err := db.GetTable(tableName)
	//	require.NoError(t, err, "failed to get table %s", tableName)
	//	for key := range expectedData[tableName] {
	//		value, ok, err := table.Get([]byte(key))
	//		require.NoError(t, err, "failed to get key %s in table %s", key, tableName)
	//		require.True(t, ok, "key %s not found in table %s", key, tableName)
	//		require.Equal(t, expectedData[tableName][key], value,
	//			"value for key %s in table %s does not match expected value", key, tableName)
	//	}
	//}
	//
	//// Shut down the DB and rebase it.
	//err = db.Close()
	//require.NoError(t, err, "failed to close DB")
	//
	//err = rebase(logger, sourceDirList, destDirList, shallow, preserveOriginal, false, verbose)
	//require.NoError(t, err, "failed to rebase DB")
	//
	//// Verify the new directories.
	//for _, sourceDir := range sourceDirList {
	//	exists, err := util.Exists(sourceDir)
	//	require.NoError(t, err)
	//
	//	if preserveOriginal {
	//		// We should see each source dir if preserveOriginal is true.
	//		require.True(t, exists, "source directory %s does not exist", sourceDir)
	//	} else {
	//		// If we aren't preserving the original, then a source directory should only exist if it overlaps.
	//		if _, ok := destDirSet[sourceDir]; !ok {
	//			require.False(t, exists, "source directory %s exists but should not", sourceDir)
	//		} else {
	//			require.True(t, exists, "source directory %s does not exist but should", sourceDir)
	//		}
	//	}
	//}
	//for _, destDir := range destDirList {
	//	// We should see all destination dirs.
	//	exists, err := util.Exists(destDir)
	//	require.NoError(t, err)
	//	require.True(t, exists, "destination directory %s does not exist", destDir)
	//}
	//
	//// Reopen the DB at the new destination directories.
	//config.Paths = destDirList
	//db, err = littbuilder.NewDB(config)
	//require.NoError(t, err, "failed to open DB after rebase")
	//
	//// Verify the data in the DB.
	//for tableName := range expectedData {
	//	table, err := db.GetTable(tableName)
	//	require.NoError(t, err, "failed to get table %s", tableName)
	//	for key := range expectedData[tableName] {
	//		value, ok, err := table.Get([]byte(key))
	//		require.NoError(t, err, "failed to get key %s in table %s", key, tableName)
	//		require.True(t, ok, "key %s not found in table %s", key, tableName)
	//		require.Equal(t, expectedData[tableName][key], value,
	//			"value for key %s in table %s does not match expected value", key, tableName)
	//	}
	//}
	//
	//err = db.Close()
	//require.NoError(t, err, "failed to close DB after rebase")
}
