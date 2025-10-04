package integration

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// Generates a random ECDSA private key and returns it as a hex string (without the "0x" prefix)
func GenerateRandomPrivateKeyHex(t *testing.T) string {
	t.Helper()
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	return gethcommon.Bytes2Hex(crypto.FromECDSA(privateKey))
}

// Derives the Ethereum address (account ID) from a private key hex string
func GetAccountIDFromPrivateKeyHex(t *testing.T, privateKeyHex string) gethcommon.Address {
	t.Helper()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
