package reservation

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewReservationVaultMonitorInvalidInterval(t *testing.T) {
	t.Run("zero interval", func(t *testing.T) {
		monitor, err := NewReservationVaultMonitor(
			context.Background(),
			testutils.GetLogger(),
			vault.NewTestPaymentVault(),
			0, // zero interval
			1024,
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *Reservation) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})

	t.Run("negative interval", func(t *testing.T) {
		monitor, err := NewReservationVaultMonitor(
			context.Background(),
			testutils.GetLogger(),
			vault.NewTestPaymentVault(),
			-time.Second, // negative interval
			1024,
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *Reservation) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})
}

func TestReservationVaultMonitor(t *testing.T) {
	testTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

	ctx, cancel := context.WithCancel(context.Background())
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
		testVault.SetReservation(addr, &bindings.IPaymentVaultReservation{
			SymbolsPerSecond: uint64(100 + i*10),
			StartTimestamp:   uint64(testTime.Unix()),
			EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
			QuorumNumbers:    []byte{0},
			QuorumSplits:     []byte{100},
		})
	}

	capturedUpdates := make(map[gethcommon.Address]*Reservation)
	updateReservation := func(accountID gethcommon.Address, newReservation *Reservation) error {
		capturedUpdates[accountID] = newReservation
		return nil
	}

	monitor, err := NewReservationVaultMonitor(
		ctx,
		testutils.GetLogger(),
		testVault,
		updateInterval,
		2, // Small batch size to force multiple batches
		func() []gethcommon.Address { return accounts },
		updateReservation,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	time.Sleep(updateInterval * 10)

	require.Equal(t, len(accounts), len(capturedUpdates))
	for i, addr := range accounts {
		reservation, ok := capturedUpdates[addr]
		require.True(t, ok, "account %s should have been updated", addr.Hex())
		require.NotNil(t, reservation)
		require.Equal(t, uint64(100+i*10), reservation.symbolsPerSecond)
	}

	// update one of the reservations
	testAccount := accounts[2]
	testVault.SetReservation(testAccount, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 999, // Changed
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	// Clear captured updates to verify new updates
	capturedUpdates = make(map[gethcommon.Address]*Reservation)

	// Wait for the monitor to fetch the updated reservation
	time.Sleep(updateInterval * 10)

	// Verify all accounts were updated again
	require.Equal(t, len(accounts), len(capturedUpdates))

	// Check that the specific account was updated correctly
	updatedReservation, ok := capturedUpdates[testAccount]
	require.True(t, ok, "account %s should have been updated", testAccount.Hex())
	require.NotNil(t, updatedReservation)
	require.Equal(t, uint64(999), updatedReservation.symbolsPerSecond)

	// Other accounts should remain unchanged
	for i, addr := range accounts {
		if addr != testAccount {
			reservation, ok := capturedUpdates[addr]
			require.True(t, ok, "account %s should have been updated", addr.Hex())
			require.NotNil(t, reservation)
			require.Equal(t, uint64(100+i*10), reservation.symbolsPerSecond)
		}
	}
}

func TestReservationVaultMonitorNoBatching(t *testing.T) {
	testTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updateInterval := time.Millisecond

	// Create multiple accounts to verify they're all fetched in a single batch
	accounts := []gethcommon.Address{
		gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
	}

	testVault := vault.NewTestPaymentVault()
	for i, addr := range accounts {
		testVault.SetReservation(addr, &bindings.IPaymentVaultReservation{
			SymbolsPerSecond: uint64(200 + i*20),
			StartTimestamp:   uint64(testTime.Unix()),
			EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
			QuorumNumbers:    []byte{0},
			QuorumSplits:     []byte{100},
		})
	}

	capturedUpdates := make(map[gethcommon.Address]*Reservation)
	updateReservation := func(accountID gethcommon.Address, newReservation *Reservation) error {
		capturedUpdates[accountID] = newReservation
		return nil
	}

	monitor, err := NewReservationVaultMonitor(
		ctx,
		testutils.GetLogger(),
		testVault,
		updateInterval,
		0, // Batch size 0 means no batching - all accounts in one call
		func() []gethcommon.Address { return accounts },
		updateReservation,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	// Wait for updates
	time.Sleep(updateInterval * 10)

	// Verify all accounts were updated correctly with no batching
	require.Equal(t, len(accounts), len(capturedUpdates))
	for i, addr := range accounts {
		reservation, ok := capturedUpdates[addr]
		require.True(t, ok, "account %s should have been updated", addr.Hex())
		require.NotNil(t, reservation)
		require.Equal(t, uint64(200+i*20), reservation.symbolsPerSecond)
	}
}
