package testbed_test

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// TestAnvilBlockNumberReachesFive verifies that Anvil will eventually reach block 5
func TestAnvilBlockNumberReachesFive(t *testing.T) {
	ctx := t.Context()
	logger := test.GetLogger()

	// Start Anvil container with 1 second block time
	anvil, err := testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		Logger:         logger,
		BlockTime:      1, // 1 second block intervals
	})
	require.NoError(t, err)
	defer func() {
		_ = anvil.Terminate(ctx)
	}()

	// Connect to Anvil RPC
	client, err := ethclient.Dial(anvil.RpcURL())
	require.NoError(t, err)
	defer client.Close()

	// Assert that block number eventually reaches at least 5
	require.Eventually(t, func() bool {
		blockNum, err := client.BlockNumber(ctx)
		if err != nil {
			logger.Warn("Failed to get block number", "error", err)
			return false
		}
		logger.Debug("Current block number", "block", blockNum)
		return blockNum >= 5
	}, 10*time.Second, 500*time.Millisecond, "Block number should reach at least 5 within 10 seconds")

	logger.Info("Successfully reached block 5")
}
