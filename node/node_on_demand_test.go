package node

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/stretchr/testify/require"
)

var testStartTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

func TestNodeOnDemandMeteringPaths(t *testing.T) {
	ctx := context.Background()
	pv := vault.NewTestPaymentVault()
	// Reduce capacity so we can exhaust quickly: 10 * 1 = 10 symbols
	pv.SetGlobalSymbolsPerSecond(10)
	pv.SetGlobalRatePeriodInterval(1)
	pv.SetMinNumSymbols(1)

	timeSource := func() time.Time { return testStartTime }
	m, err := meterer.NewOnDemandMeterer(ctx, pv, timeSource, nil, 1.0)
	require.NoError(t, err)

	n := &Node{}
	n.SetOnDemandMeterer(m)

	// Success path: reserve within capacity
	res, err := n.MeterOnDemandDispersal(5)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Cancel should be safe even when reservation is nil
	n.CancelOnDemandDispersal(nil)
	n.CancelOnDemandDispersal(res)

	// Consume remaining capacity then verify exhaustion
	res, err = n.MeterOnDemandDispersal(10)
	require.NoError(t, err)

	_, err = n.MeterOnDemandDispersal(1)
	require.Error(t, err)
}
