package node_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/require"
)

func TestDeleteV1Data_NonExistentDirectory(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	dbPath := t.TempDir()
	// Don't create the chunk subdirectory - it should not exist

	err = node.DeleteV1Data(logger, dbPath)
	require.NoError(t, err)
}

func TestDeleteV1Data_FileInsteadOfDirectory(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	dbPath := t.TempDir()
	v1DataPath := filepath.Join(dbPath, node.V1ChunkSubdir)

	// Create a file (not a directory) at the v1 data path
	err = os.WriteFile(v1DataPath, []byte("not a directory"), 0644)
	require.NoError(t, err)

	err = node.DeleteV1Data(logger, dbPath)
	require.Error(t, err, "should return error when path is a file instead of directory")
}

func TestDeleteV1Data_NestedDirectories(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	dbPath := t.TempDir()
	v1DataPath := filepath.Join(dbPath, node.V1ChunkSubdir)

	// Create nested directory structure
	nestedPath := filepath.Join(v1DataPath, "subdir1", "subdir2")
	err = os.MkdirAll(nestedPath, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(v1DataPath, "file1.db"), []byte("data1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(nestedPath, "file2.db"), []byte("data2"), 0644)
	require.NoError(t, err)

	err = node.DeleteV1Data(logger, dbPath)
	require.NoError(t, err)

	_, err = os.Stat(v1DataPath)
	require.True(t, os.IsNotExist(err), "v1 data directory should not exist after deletion")
}
