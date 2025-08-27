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

func TestClientLedgerConstructor(t *testing.T) {
	t.Run("zero address panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				common.Address{}, // zero address
				ClientLedgerModeReservationOnly,
				buildReservationLedger(t),
				nil,
				func() time.Time { return testStartTime },
			)
		}, "zero address should cause panic")
	})

	t.Run("nil getNow panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeReservationOnly,
				buildReservationLedger(t),
				nil,
				nil, // nil getNow
			)
		}, "nil getNow should cause panic")
	})

	t.Run("invalid mode panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerMode("invalid_mode"),
				buildReservationLedger(t),
				nil,
				func() time.Time { return testStartTime },
			)
		}, "invalid mode should cause panic")
	})

	t.Run("reservation-only mode with nil reservation ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeReservationOnly,
				nil, // nil reservation ledger
				nil,
				func() time.Time { return testStartTime },
			)
		}, "reservation-only mode with nil reservation ledger should cause panic")
	})

	t.Run("reservation-only mode with non-nil on-demand ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeReservationOnly,
				buildReservationLedger(t),
				buildOnDemandLedger(t), // should be nil
				func() time.Time { return testStartTime },
			)
		}, "reservation-only mode with non-nil on-demand ledger should cause panic")
	})

	t.Run("on-demand-only mode with nil on-demand ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeOnDemandOnly,
				nil,
				nil, // nil on-demand ledger
				func() time.Time { return testStartTime },
			)
		}, "on-demand-only mode with nil on-demand ledger should cause panic")
	})

	t.Run("on-demand-only mode with non-nil reservation ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeOnDemandOnly,
				buildReservationLedger(t), // should be nil
				buildOnDemandLedger(t),
				func() time.Time { return testStartTime },
			)
		}, "on-demand-only mode with non-nil reservation ledger should cause panic")
	})

	t.Run("reservation-and-on-demand mode with nil reservation ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeReservationAndOnDemand,
				nil, // nil reservation ledger
				buildOnDemandLedger(t),
				func() time.Time { return testStartTime },
			)
		}, "reservation-and-on-demand mode with nil reservation ledger should cause panic")
	})

	t.Run("reservation-and-on-demand mode with nil on-demand ledger panic", func(t *testing.T) {
		require.Panics(t, func() {
			NewClientLedger(
				testutils.GetLogger(),
				nil,
				accountID,
				ClientLedgerModeReservationAndOnDemand,
				buildReservationLedger(t),
				nil, // nil on-demand ledger
				func() time.Time { return testStartTime },
			)
		}, "reservation-and-on-demand mode with nil on-demand ledger should cause panic")
	})
}

func TestReservationOnly(t *testing.T) {
	t.Run("insufficient capacity error", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime },
		)
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

		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			getNow,
		)
		require.NotNil(t, clientLedger)

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
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{99})
		})
	})

	t.Run("time out of range panic", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime.Add(2 * time.Hour) },
		)
		require.NotNil(t, clientLedger)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		}, "expired reservation should cause fatal panic")
	})
}

func TestOnDemandOnly(t *testing.T) {
	t.Run("successful debit with cumulative payment", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		paymentMetadata, err := clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.True(t, paymentMetadata.IsOnDemand())
		// 100 symbols * 10 wei per symbol = 1000 wei
		require.Equal(t, big.NewInt(1000), paymentMetadata.CumulativePayment)
		require.Equal(t, accountID, paymentMetadata.AccountID)
	})

	t.Run("insufficient funds panic", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1001, []core.QuorumID{0, 1})
		}, "insufficient funds should cause fatal panic in on-demand only mode")
	})

	t.Run("fatal errors cause panic", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{99})
		}, "forbidden quorum should cause fatal panic")
	})
}

func TestReservationAndOnDemand(t *testing.T) {
	t.Run("fallback to on-demand", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		// First debit uses all reservation capacity
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1000, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.False(t, paymentMetadata.IsOnDemand())

		// Second debit should fallback to on-demand
		paymentMetadata, err = clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
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

		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			getNow,
		)
		require.NotNil(t, clientLedger)

		// First debit to establish a time baseline
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.False(t, paymentMetadata.IsOnDemand())

		// Move time backward
		currentTime = testStartTime.Add(-time.Minute)

		paymentMetadata, err = clientLedger.Debit(context.Background(), 1, []core.QuorumID{0, 1})
		require.Error(t, err, "time moved backward should cause retriable error")
		require.Nil(t, paymentMetadata)
	})

	t.Run("insufficient funds error from on-demand", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		// First debit uses all reservation capacity
		paymentMetadata, err := clientLedger.Debit(context.Background(), 1000, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.False(t, paymentMetadata.IsOnDemand())

		// Second debit should fallback to on-demand but fails due to insufficient funds
		paymentMetadata, err = clientLedger.Debit(context.Background(), 1001, []core.QuorumID{0, 1})
		require.Error(t, err, "insufficient funds in on-demand should cause retriable error in combined mode")
		require.Nil(t, paymentMetadata)
	})

	t.Run("fatal errors cause panic", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationAndOnDemand,
			buildReservationLedger(t),
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		require.Panics(t, func() {
			_, _ = clientLedger.Debit(context.Background(), 1, []core.QuorumID{99})
		}, "forbidden quorum should cause fatal panic")
	})
}

func TestRevertDebit(t *testing.T) {
	t.Run("successful reservation revert", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeReservationOnly,
			buildReservationLedger(t),
			nil,
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		paymentMetadata, err := clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.False(t, paymentMetadata.IsOnDemand())

		err = clientLedger.RevertDebit(context.Background(), paymentMetadata, 100)
		require.NoError(t, err)
	})

	t.Run("successful on-demand revert", func(t *testing.T) {
		clientLedger := NewClientLedger(
			testutils.GetLogger(),
			nil,
			accountID,
			ClientLedgerModeOnDemandOnly,
			nil,
			buildOnDemandLedger(t),
			func() time.Time { return testStartTime },
		)
		require.NotNil(t, clientLedger)

		paymentMetadata, err := clientLedger.Debit(context.Background(), 100, []core.QuorumID{0, 1})
		require.NoError(t, err)
		require.NotNil(t, paymentMetadata)
		require.True(t, paymentMetadata.IsOnDemand())

		err = clientLedger.RevertDebit(context.Background(), paymentMetadata, 100)
		require.NoError(t, err)
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
