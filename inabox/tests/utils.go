package integration_test

import (
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
