package ondemand_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/test"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandVaultMonitorInvalidInterval(t *testing.T) {
	ctx := t.Context()

	t.Run("zero interval", func(t *testing.T) {
		monitor, err := ondemand.NewOnDemandVaultMonitor(
			ctx,
			test.GetLogger(),
			vault.NewTestPaymentVault(),
			0, // zero interval
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *big.Int) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})

	t.Run("negative interval", func(t *testing.T) {
		monitor, err := ondemand.NewOnDemandVaultMonitor(
			ctx,
			test.GetLogger(),
			vault.NewTestPaymentVault(),
			-time.Second, // negative interval
			func() []gethcommon.Address { return nil },
			func(gethcommon.Address, *big.Int) error { return nil },
		)
		require.Error(t, err)
		require.Nil(t, monitor)
	})
}

// tests basic vault monitor behavior
func TestOnDemandVaultMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	logger := test.GetLogger()
	updateInterval := 1 * time.Millisecond
	address := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetTotalDeposit(address, big.NewInt(5000))

	var capturedAccountID gethcommon.Address
	var capturedDeposit *big.Int
	updateTotalDeposit := func(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
		capturedAccountID = accountID
		capturedDeposit = newTotalDeposit
		return nil
	}

	monitor, err := ondemand.NewOnDemandVaultMonitor(
		ctx,
		logger,
		testVault,
		updateInterval,
		func() []gethcommon.Address { return []gethcommon.Address{address} },
		updateTotalDeposit,
	)
	require.NoError(t, err)
	require.NotNil(t, monitor)

	time.Sleep(updateInterval * 10)
	require.Equal(t, address, capturedAccountID)
	require.Equal(t, big.NewInt(5000), capturedDeposit)

	testVault.SetTotalDeposit(address, big.NewInt(5001))
	time.Sleep(updateInterval * 10)
	require.Equal(t, big.NewInt(5001), capturedDeposit, "update should have been observed")
}
