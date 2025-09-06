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
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *Reservation) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})
}

func TestReservationVaultMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updateInterval := time.Millisecond
	address := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetReservation(address, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 100,
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	var capturedAccountID gethcommon.Address
	var capturedReservation *Reservation
	updateReservation := func(accountID gethcommon.Address, newReservation *Reservation) error {
		capturedAccountID = accountID
		capturedReservation = newReservation
		return nil
	}

	monitor, err := NewReservationVaultMonitor(
		ctx,
		testutils.GetLogger(),
		testVault,
		updateInterval,
		func() []gethcommon.Address { return []gethcommon.Address{address} },
		updateReservation,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	time.Sleep(updateInterval * 10)
	require.Equal(t, address, capturedAccountID)
	require.NotNil(t, capturedReservation)

	// Update the reservation
	testVault.SetReservation(address, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 200, // Changed
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})
	capturedReservation = nil

	time.Sleep(updateInterval * 10)
	require.NotNil(t, capturedReservation, "update should have been observed")
	require.Equal(t, uint64(200), capturedReservation.symbolsPerSecond)
}
