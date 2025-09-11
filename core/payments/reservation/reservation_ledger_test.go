package reservation

import (
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

func TestDebit(t *testing.T) {
	t.Run("successful debit", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)
		dispersalTime := startTime.Add(time.Hour)

		success, remainingCapacity, err := ledger.Debit(
			startTime,
			dispersalTime,
			50,
			[]core.QuorumID{0},
		)
		require.NoError(t, err)
		require.True(t, success)
		require.Greater(t, remainingCapacity, float64(0))
	})

	t.Run("invalid quorum", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)
		dispersalTime := startTime.Add(time.Hour)

		success, _, err := ledger.Debit(
			startTime,
			dispersalTime,
			50,
			[]core.QuorumID{0, 1, 5}, // quorum 5 not permitted
		)
		require.Error(t, err)
		require.False(t, success)

		var quorumNotPermittedError *QuorumNotPermittedError
		require.True(t, errors.As(err, &quorumNotPermittedError))
	})

	t.Run("invalid dispersal time", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)

		// before reservation start
		success, _, err := ledger.Debit(
			startTime,
			startTime.Add(-time.Hour),
			50,
			[]core.QuorumID{0},
		)
		require.Error(t, err)
		require.False(t, success)
		var timeOutOfRangeError *TimeOutOfRangeError
		require.True(t, errors.As(err, &timeOutOfRangeError))

		// after reservation end
		success, _, err = ledger.Debit(
			startTime,
			startTime.Add(25*time.Hour),
			50,
			[]core.QuorumID{0},
		)
		require.Error(t, err)
		require.False(t, success)
		require.True(t, errors.As(err, &timeOutOfRangeError))
	})

	t.Run("minimum symbols applied", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)
		dispersalTime := startTime.Add(time.Hour)

		// debit 5 symbols, but minNumSymbols is 10
		success, remainingCapacity, err := ledger.Debit(
			startTime,
			dispersalTime,
			5,
			[]core.QuorumID{0},
		)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, float64(990), remainingCapacity)
	})
}

func TestRevertDebit(t *testing.T) {
	t.Run("successful revert", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)
		dispersalTime := startTime.Add(time.Hour)

		// debit first
		success, _, err := ledger.Debit(
			startTime,
			dispersalTime,
			100,
			[]core.QuorumID{0},
		)
		require.NoError(t, err)
		require.True(t, success)

		// revert the debit
		remainingCapacity, err := ledger.RevertDebit(startTime, 50)
		require.NoError(t, err)
		require.Equal(t, float64(950), remainingCapacity)
	})

	t.Run("minimum symbols applied", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		ledger := createTestLedger(t, 100, false, startTime)
		dispersalTime := startTime.Add(time.Hour)

		// debit 5 (charged 10 due to minimum)
		success, _, err := ledger.Debit(
			startTime,
			dispersalTime,
			5,
			[]core.QuorumID{0},
		)
		require.NoError(t, err)
		require.True(t, success)

		// revert 5 (should revert 10 due to minimum)
		remainingCapacity, err := ledger.RevertDebit(startTime, 5)
		require.NoError(t, err)
		require.Equal(t, float64(1000), remainingCapacity)
	})
}

func TestIsBucketEmpty(t *testing.T) {
	startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
	ledger := createTestLedger(t, 100, false, startTime)
	dispersalTime := startTime.Add(time.Hour)

	// initially empty
	isEmpty, err := ledger.IsBucketEmpty(startTime)
	require.NoError(t, err)
	require.True(t, isEmpty)

	// after debit, not empty
	success, _, err := ledger.Debit(
		startTime,
		dispersalTime,
		100,
		[]core.QuorumID{0},
	)
	require.NoError(t, err)
	require.True(t, success)

	isEmpty, err = ledger.IsBucketEmpty(startTime)
	require.NoError(t, err)
	require.False(t, isEmpty)
}

func TestUpdateReservation(t *testing.T) {
	startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
	ledger := createTestLedger(t, 100, false, startTime)
	dispersalTime := startTime.Add(time.Hour)

	// debit 500 symbols to establish a fill level
	success, remainingCapacity, err := ledger.Debit(
		startTime,
		dispersalTime,
		500,
		[]core.QuorumID{0, 1},
	)
	require.NoError(t, err)
	require.True(t, success)
	// 1000 - 500 = 500
	require.Equal(t, float64(500), remainingCapacity)

	// update with identical reservation
	endTime := startTime.Add(24 * time.Hour)
	identicalReservation, err := NewReservation(100, startTime, endTime, []core.QuorumID{0, 1})
	require.NoError(t, err)
	err = ledger.UpdateReservation(identicalReservation, startTime)
	require.NoError(t, err)

	// totalCapacity should remain the same
	totalCapacity := ledger.GetBucketCapacity()
	require.Equal(t, float64(1000), totalCapacity)
	// verify fill level was preserved by doing another debit (100 symbols)
	success, remainingCapacity, err = ledger.Debit(
		startTime,
		dispersalTime,
		100,
		[]core.QuorumID{0},
	)
	require.NoError(t, err)
	require.True(t, success)
	// 1000 - 500 - 100 = 400
	require.Equal(t, float64(400), remainingCapacity)

	// update all fields
	newStartTime := startTime.Add(-time.Hour)
	newEndTime := startTime.Add(48 * time.Hour)
	newReservation, err := NewReservation(200, newStartTime, newEndTime, []core.QuorumID{0}) // only quorum 0 now
	require.NoError(t, err)
	err = ledger.UpdateReservation(newReservation, startTime)
	require.NoError(t, err)

	// verify new total capacity (200 * 10 = 2000)
	totalCapacity = ledger.GetBucketCapacity()
	require.Equal(t, float64(2000), totalCapacity)
	// verify fill level was preserved by doing another debit (100 symbols)
	success, remainingCapacity, err = ledger.Debit(
		startTime,
		dispersalTime,
		100,
		[]core.QuorumID{0},
	)
	require.NoError(t, err)
	require.True(t, success)
	// 2000 - 500 - 100 - 100 = 1300
	require.Equal(t, float64(1300), remainingCapacity)

	// verify new quorum restrictions are enforced
	success, _, err = ledger.Debit(
		startTime,
		dispersalTime,
		50,
		[]core.QuorumID{1}, // quorum 1 no longer permitted
	)
	require.Error(t, err)
	require.False(t, success)
	var quorumNotPermittedError *QuorumNotPermittedError
	require.True(t, errors.As(err, &quorumNotPermittedError))

	// verify new time window is enforced
	lateDispersalTime := startTime.Add(30 * time.Hour)
	success, _, err = ledger.Debit(
		startTime,
		lateDispersalTime, // within new 48 hour window
		50,
		[]core.QuorumID{0},
	)
	require.NoError(t, err)
	require.True(t, success)

	// update with nil reservation
	err = ledger.UpdateReservation(nil, startTime)
	require.Error(t, err)
}

func createTestLedger(t *testing.T, symbolsPerSecond uint64, startFull bool, startTime time.Time) *ReservationLedger {
	endTime := startTime.Add(24 * time.Hour)
	permittedQuorums := []core.QuorumID{0, 1}

	reservation, err := NewReservation(symbolsPerSecond, startTime, endTime, permittedQuorums)
	require.NoError(t, err)

	config, err := NewReservationLedgerConfig(
		*reservation,
		10, // minNumSymbols
		startFull,
		OverfillOncePermitted,
		10*time.Second,
	)
	require.NoError(t, err)

	ledger, err := NewReservationLedger(*config, startTime)
	require.NoError(t, err)
	require.NotNil(t, ledger)

	return ledger
}
