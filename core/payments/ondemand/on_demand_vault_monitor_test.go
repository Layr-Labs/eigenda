package ondemand

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/test"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandVaultMonitorInvalidInterval(t *testing.T) {
	ctx := t.Context()
	t.Run("zero interval", func(t *testing.T) {
		monitor, err := NewOnDemandVaultMonitor(
			ctx,
			test.GetLogger(),
			vault.NewTestPaymentVault(),
			0, // zero interval
			1024,
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *big.Int) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})

	t.Run("negative interval", func(t *testing.T) {
		monitor, err := NewOnDemandVaultMonitor(
			ctx,
			test.GetLogger(),
			vault.NewTestPaymentVault(),
			-time.Second, // negative interval
			1024,
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *big.Int) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})
}

func TestOnDemandVaultMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	updateInterval := time.Millisecond

	accounts := []gethcommon.Address{
		gethcommon.HexToAddress("0x1111111111111111111111111111111111111111"),
		gethcommon.HexToAddress("0x2222222222222222222222222222222222222222"),
		gethcommon.HexToAddress("0x3333333333333333333333333333333333333333"),
		gethcommon.HexToAddress("0x4444444444444444444444444444444444444444"),
		gethcommon.HexToAddress("0x5555555555555555555555555555555555555555"),
	}

	testVault := vault.NewTestPaymentVault()
	for i, addr := range accounts {
		testVault.SetTotalDeposit(addr, big.NewInt(int64(1000+i*100)))
	}

	capturedUpdates := make(map[gethcommon.Address]*big.Int)
	updateTotalDeposit := func(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
		capturedUpdates[accountID] = newTotalDeposit
		return nil
	}

	monitor, err := NewOnDemandVaultMonitor(
		ctx,
		test.GetLogger(),
		testVault,
		updateInterval,
		2, // Small batch size to force multiple batches
		func() []gethcommon.Address { return accounts },
		updateTotalDeposit,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	test.AssertEventuallyEquals(t, len(accounts), func() int {
		return len(capturedUpdates)
	}, time.Second)

	for i, addr := range accounts {
		deposit, ok := capturedUpdates[addr]
		require.True(t, ok, "account %s should have been updated", addr.Hex())
		require.NotNil(t, deposit)
		require.Equal(t, big.NewInt(int64(1000+i*100)), deposit)
	}

	// update one of the deposits
	testAccount := accounts[2]
	testVault.SetTotalDeposit(testAccount, big.NewInt(9999)) // Changed

	// Clear captured updates to verify new updates
	capturedUpdates = make(map[gethcommon.Address]*big.Int)

	// Wait for the monitor to fetch the updated deposits
	test.AssertEventuallyEquals(t, len(accounts), func() int {
		return len(capturedUpdates)
	}, time.Second)

	// Check that the specific account was updated correctly
	updatedDeposit, ok := capturedUpdates[testAccount]
	require.True(t, ok, "account %s should have been updated", testAccount.Hex())
	require.NotNil(t, updatedDeposit)
	require.Equal(t, big.NewInt(9999), updatedDeposit)

	// Other accounts should remain unchanged
	for i, addr := range accounts {
		if addr != testAccount {
			deposit, ok := capturedUpdates[addr]
			require.True(t, ok, "account %s should have been updated", addr.Hex())
			require.NotNil(t, deposit)
			require.Equal(t, big.NewInt(int64(1000+i*100)), deposit)
		}
	}
}

func TestOnDemandVaultMonitorNoBatching(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	updateInterval := time.Millisecond

	// Create multiple accounts to verify they're all fetched in a single batch
	accounts := []gethcommon.Address{
		gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
	}

	testVault := vault.NewTestPaymentVault()
	for i, addr := range accounts {
		testVault.SetTotalDeposit(addr, big.NewInt(int64(2000+i*200)))
	}

	capturedUpdates := make(map[gethcommon.Address]*big.Int)
	updateTotalDeposit := func(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
		capturedUpdates[accountID] = newTotalDeposit
		return nil
	}

	monitor, err := NewOnDemandVaultMonitor(
		ctx,
		test.GetLogger(),
		testVault,
		updateInterval,
		0, // Batch size 0 means no batching - all accounts in one call
		func() []gethcommon.Address { return accounts },
		updateTotalDeposit,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	// Wait for updates
	test.AssertEventuallyEquals(t, len(accounts), func() int {
		return len(capturedUpdates)
	}, time.Second)
	for i, addr := range accounts {
		deposit, ok := capturedUpdates[addr]
		require.True(t, ok, "account %s should have been updated", addr.Hex())
		require.NotNil(t, deposit)
		require.Equal(t, big.NewInt(int64(2000+i*200)), deposit)
	}
}
