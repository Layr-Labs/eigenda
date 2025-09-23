package meterer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var startTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

func TestMeterDispersal(t *testing.T) {
	timeSource := func() time.Time { return startTime }
	meterer := NewOnDemandMeterer(100, 10, timeSource, nil)

	reservation, err := meterer.MeterDispersal(500)
	require.NoError(t, err)
	require.NotNil(t, reservation)
	require.True(t, reservation.OK())
}

func TestCancelDispersal(t *testing.T) {
	timeSource := func() time.Time { return startTime }
	meterer := NewOnDemandMeterer(100, 10, timeSource, nil)

	reservation, err := meterer.MeterDispersal(500)
	require.NoError(t, err)
	require.NotNil(t, reservation)

	// don't panic
	meterer.CancelDispersal(reservation)
}
