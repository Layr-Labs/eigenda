package reservation

import (
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

func TestNewReservation(t *testing.T) {
	t.Run("create with valid parameters", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)
	})

	t.Run("create with invalid parameters", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(0, startTime, endTime, permittedQuorums)
		require.Nil(t, reservation)
		require.Error(t, err, "zero symbols per second should error")

		reservation, err = NewReservation(100, startTime, startTime, permittedQuorums)
		require.Nil(t, reservation)
		require.Error(t, err, "startTime == endTime should error")

		reservation, err = NewReservation(100, endTime, startTime, permittedQuorums)
		require.Nil(t, reservation)
		require.Error(t, err, "endTime < startTime should error")

		reservation, err = NewReservation(100, startTime, endTime, []core.QuorumID{})
		require.Nil(t, reservation)
		require.Error(t, err, "no permitted quorums should error")
	})
}

func TestCheckQuorumsPermitted(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		err = reservation.CheckQuorumsPermitted(permittedQuorums)
		require.NoError(t, err)
	})

	t.Run("invalid quorum", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		var quorumNotPermittedError *QuorumNotPermittedError

		err = reservation.CheckQuorumsPermitted([]core.QuorumID{0, 1, 3})
		require.Error(t, err)
		require.True(t, errors.As(err, &quorumNotPermittedError))
	})
}

func TestCheckTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		err = reservation.CheckTime(startTime.Add(time.Minute))
		require.NoError(t, err)
	})

	t.Run("early time", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		var timeOutOfRangeError *TimeOutOfRangeError

		err = reservation.CheckTime(startTime.Add(-time.Minute))
		require.Error(t, err, "time before start time should fail")
		require.True(t, errors.As(err, &timeOutOfRangeError))
	})

	t.Run("late time", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		var timeOutOfRangeError *TimeOutOfRangeError

		err = reservation.CheckTime(endTime.Add(time.Minute))
		require.Error(t, err, "time after end time should fail")
		require.True(t, errors.As(err, &timeOutOfRangeError))
	})
}

func TestEqual(t *testing.T) {
	startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(time.Hour)
	quorums := []core.QuorumID{0, 1}

	// equal reservations
	r1, err := NewReservation(100, startTime, endTime, quorums)
	require.NoError(t, err)
	r2, err := NewReservation(100, startTime, endTime, quorums)
	require.NoError(t, err)
	require.True(t, r1.Equal(r2))

	// nil comparison
	require.False(t, r1.Equal(nil))

	// different symbols per second
	r3, err := NewReservation(200, startTime, endTime, quorums)
	require.NoError(t, err)
	require.False(t, r1.Equal(r3))

	// different start time
	r4, err := NewReservation(100, startTime.Add(time.Second), endTime, quorums)
	require.NoError(t, err)
	require.False(t, r1.Equal(r4))

	// different end time
	r5, err := NewReservation(100, startTime, endTime.Add(time.Second), quorums)
	require.NoError(t, err)
	require.False(t, r1.Equal(r5))

	// different number of quorums
	r6, err := NewReservation(100, startTime, endTime, []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	require.False(t, r1.Equal(r6))

	// different quorum IDs (same length, different values)
	r7, err := NewReservation(100, startTime, endTime, []core.QuorumID{0, 3})
	require.NoError(t, err)
	require.False(t, r1.Equal(r7))
}
