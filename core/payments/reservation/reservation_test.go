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

		err = reservation.CheckQuorumsPermitted([]core.QuorumID{0, 1, 3})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrQuorumNotPermitted))
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

		err = reservation.CheckTime(startTime.Add(-time.Minute))
		require.Error(t, err, "time before start time should fail")
		require.True(t, errors.Is(err, ErrTimeOutOfRange))
	})

	t.Run("late time", func(t *testing.T) {
		startTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)
		permittedQuorums := []core.QuorumID{0, 1}

		reservation, err := NewReservation(100, startTime, endTime, permittedQuorums)
		require.NotNil(t, reservation)
		require.NoError(t, err)

		err = reservation.CheckTime(endTime.Add(time.Minute))
		require.Error(t, err, "time after end time should fail")
		require.True(t, errors.Is(err, ErrTimeOutOfRange))
	})
}
