package payments

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	accountID     = common.HexToAddress("0x1234567890123456789012345678901234567890")
	testStartTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
)

func TestReservationOnly(t *testing.T) {
	t.Run("insufficient capacity error", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime },
		)
		require.NoError(t, err)
		require.NotNil(t, clientLedger)

		// first dispersal is permitted, even though it overfills bucket
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1000, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.False(t, paymentMetadata.IsOnDemand())
		require.Equal(t, big.NewInt(0), paymentMetadata.CumulativePayment)
		require.Equal(t, accountID, paymentMetadata.AccountID)

		// any additional symbols aren't permitted
		paymentMetadata, err = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.Error(t, err, "should be over capacity")
		require.Nil(t, paymentMetadata)
	})

	t.Run("time moved backward error", func(t *testing.T) {
		currentTime := testStartTime
		getNow := func() time.Time {
			return currentTime
		}

		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			getNow,
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		// First debit to establish a time baseline
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.NotNil(t, paymentMetadata)
		require.NoError(t, err)

		// Move time backward
		currentTime = testStartTime.Add(-time.Minute)

		paymentMetadata, err = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.Error(t, err, "time moved backward should cause error")
		require.Nil(t, paymentMetadata)
	})

	t.Run("quorum not permitted panic", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{99})
		})
	})

	t.Run("time out of range panic", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime.Add(2 * time.Hour) },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		}, "expired reservation should cause fatal panic")
	})
}

func TestOnDemandOnly(t *testing.T) {
	t.Run("successful debit with cumulative payment", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		paymentMetadata, err := clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.True(t, paymentMetadata.IsOnDemand())
		// 100 symbols * 10 wei per symbol = 1000 wei
		require.Equal(t, big.NewInt(1000), paymentMetadata.CumulativePayment)
		require.Equal(t, accountID, paymentMetadata.AccountID)
	})

	t.Run("insufficient funds panic", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1001, []core.QuorumID{0, 1})
		}, "insufficient funds should cause fatal panic in on-demand only mode")
	})
}

func TestReservationAndOnDemand(t *testing.T) {
	t.Run("fallback to on-demand", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		// First debit uses all reservation capacity
		_, err = clientLedger.Debit(context.Background(), 1000, []core.QuorumID{0, 1})
		require.NoError(t, err)

		// Second debit should fallback to on-demand
		paymentMetadata, err := clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.True(t, paymentMetadata.IsOnDemand())
		// 100 symbols * 10 wei per symbol = 1000 wei
		require.Equal(t, big.NewInt(1000), paymentMetadata.CumulativePayment)
		require.Equal(t, accountID, paymentMetadata.AccountID)
	})

	t.Run("time moved backward error", func(t *testing.T) {
		currentTime := testStartTime
		getNow := func() time.Time {
			return currentTime
		}

		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			getNow,
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		// First debit to establish a time baseline
		_, err = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.NoError(t, err)

		// Move time backward
		currentTime = testStartTime.Add(-time.Minute)

		paymentMetadata, err := clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.Error(t, err, "time moved backward should cause retriable error")
		require.Nil(t, paymentMetadata)
	})

	t.Run("insufficient funds error from on-demand", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		// First debit uses all reservation capacity
		_, err = clientLedger.Debit(context.Background(), 1000, []core.QuorumID{0, 1})
		require.NoError(t, err)

		// Second debit should fallback to on-demand but fail due to insufficient funds
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1001, []core.QuorumID{0, 1})
		require.Error(t, err, "insufficient funds in on-demand should cause retriable error in combined mode")
		require.Nil(t, paymentMetadata)
	})

	t.Run("fatal errors cause panic", func(t *testing.T) {
		clientLedger, err := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)
		require.NoError(t, err)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{99})
		}, "forbidden quorum should cause fatal panic")
	})
}

func buildReservationLedger(t *testing.T) *reservation.ReservationLedger {
	res, err := reservation.NewReservation(
		10, testStartTime.Add(-time.Hour), testStartTime.Add(time.Hour), []core.QuorumID{0, 1})
	require.NotNil(t, res)
	require.NoError(t, err)

	reservationLedgerConfig, err := reservation.NewReservationLedgerConfig(
		*res, false, reservation.OverfillOncePermitted, time.Minute)
	require.NotNil(t, reservationLedgerConfig)
	require.NoError(t, err)

	reservationLedger, err := reservation.NewReservationLedger(*reservationLedgerConfig, testStartTime)
	require.NotNil(t, reservationLedger)
	require.NoError(t, err)

	return reservationLedger
}

func buildOnDemandLedger(t *testing.T) *ondemand.OnDemandLedger {
	onDemandLedger, err := ondemand.OnDemandLedgerFromValue(
		big.NewInt(10000),
		big.NewInt(10),
		10,
		big.NewInt(0),
	)
	require.NoError(t, err)
	require.NotNil(t, onDemandLedger)

	return onDemandLedger
}
