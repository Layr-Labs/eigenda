package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

// MineAnvilBlocks mines the specified number of blocks in Anvil.
func MineAnvilBlocks(t *testing.T, rpcClient common.RPCEthClient, numBlocks int) {
	t.Helper()
	for i := 0; i < numBlocks; i++ {
		err := rpcClient.CallContext(t.Context(), nil, "evm_mine")
		require.NoError(t, err)
	}
}

// getSRSPaths returns the correct paths to SRS files based on the source file location.
// This uses runtime.Caller to determine where this file is located and calculates
// the relative path to the resources/srs directory from there.
func getSRSPaths() (g1Path, g2Path, g2TrailingPath string, err error) {
	// Get the path of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", "", "", fmt.Errorf("failed to get caller information")
	}

	// We need to go up 2 directories from tests/ to get to inabox/, then up one more to get to the project root
	// From project root, resources/srs is the target
	testDir := filepath.Dir(filename)
	inaboxDir := filepath.Dir(testDir)
	projectRoot := filepath.Dir(inaboxDir)

	g1Path = filepath.Join(projectRoot, "resources", "srs", "g1.point")
	g2Path = filepath.Join(projectRoot, "resources", "srs", "g2.point")
	g2TrailingPath = filepath.Join(projectRoot, "resources", "srs", "g2.trailing.point")

	return g1Path, g2Path, g2TrailingPath, nil
}
