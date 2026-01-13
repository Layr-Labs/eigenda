package meterer

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/stretchr/testify/require"
)

var startTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

func TestMeterDispersal(t *testing.T) {
	ctx := t.Context()
	timeSource := func() time.Time { return startTime }

	paymentVault := vault.NewTestPaymentVault()
	// bucket capacity is 100*10 = 1000 symbols
	paymentVault.SetGlobalSymbolsPerSecond(100)
	paymentVault.SetGlobalRatePeriodInterval(10)
	paymentVault.SetMinNumSymbols(100)

	meterer, err := NewOnDemandMeterer(ctx, paymentVault, timeSource, nil, 1.0)
	require.NoError(t, err)

	// blob larger than minNumSymbols
	reservation, err := meterer.MeterDispersal(850)
	require.NoError(t, err)
	require.NotNil(t, reservation)

	// blob below minNumSymbols - should meter minNumSymbols (100)
	reservation, err = meterer.MeterDispersal(50)
	require.NoError(t, err)
	require.NotNil(t, reservation)

	// blob below minNumSymbols - should meter minNumSymbols (100), but we've exhausted capacity
	reservation, err = meterer.MeterDispersal(1)
	require.Error(t, err, "should have exceeded available meter capacity")
	require.Nil(t, reservation)
}

func TestCancelDispersal(t *testing.T) {
	ctx := t.Context()
	timeSource := func() time.Time { return startTime }
	paymentVault := vault.NewTestPaymentVault()
	paymentVault.SetGlobalSymbolsPerSecond(100)
	paymentVault.SetGlobalRatePeriodInterval(10)
	paymentVault.SetMinNumSymbols(100)

	meterer, err := NewOnDemandMeterer(ctx, paymentVault, timeSource, nil, 1.0)
	require.NoError(t, err)

	reservation, err := meterer.MeterDispersal(500)
	require.NoError(t, err)
	require.NotNil(t, reservation)

	// don't panic
	meterer.CancelDispersal(reservation)
}

func TestRefreshUpdatesLimits(t *testing.T) {
	ctx := context.Background()
	timeSource := func() time.Time { return startTime }

	paymentVault := vault.NewTestPaymentVault()
	paymentVault.SetGlobalSymbolsPerSecond(100)
	paymentVault.SetGlobalRatePeriodInterval(10)
	paymentVault.SetMinNumSymbols(1)

	meterer, err := NewOnDemandMeterer(ctx, paymentVault, timeSource, nil, 1.0)
	require.NoError(t, err)

	// Exhaust initial capacity (100 * 10 = 1000)
	_, err = meterer.MeterDispersal(1000)
	require.NoError(t, err)
	_, err = meterer.MeterDispersal(1)
	require.Error(t, err, "expected exhaustion at initial capacity")

	// Increase on-chain limit and refresh; capacity should expand
	paymentVault.SetGlobalSymbolsPerSecond(200) // new capacity: 2000
	err = meterer.Refresh(ctx)
	require.NoError(t, err)

	_, err = meterer.MeterDispersal(1000) // should consume remaining new capacity
	require.NoError(t, err)
	_, err = meterer.MeterDispersal(1)
	require.Error(t, err, "expected exhaustion after consuming expanded capacity")

	// Refresh with unchanged params should be a no-op
	err = meterer.Refresh(ctx)
	require.NoError(t, err)
}
