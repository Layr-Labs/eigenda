package replay

import (
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func TestTooOldRequest(t *testing.T) {
	rand := random.NewTestRandom()

	now := rand.Time()
	timeSource := func() time.Time {
		return now
	}

	maxTimeInPast := time.Duration(rand.Intn(5)+1) * time.Minute
	maxTimeInFuture := time.Duration(rand.Intn(5)+1) * time.Minute

	rGuard := NewReplayGuardian(timeSource, maxTimeInPast, maxTimeInFuture)

	requestAge := maxTimeInPast + 1
	requestTime := now.Add(-requestAge)

	err := rGuard.VerifyRequest(rand.Bytes(32), requestTime)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "too far in the past"))

	// Verify that nothing has been added to the observedHashes set.
	g := rGuard.(*replayGuardian)
	require.Zero(t, len(g.observedHashes))
	require.Zero(t, g.expirationQueue.Size())
}

func TestTooFarInFutureRequest(t *testing.T) {
	rand := random.NewTestRandom()

	now := rand.Time()
	timeSource := func() time.Time {
		return now
	}

	maxTimeInPast := time.Duration(rand.Intn(5)+1) * time.Minute
	maxTimeInFuture := time.Duration(rand.Intn(5)+1) * time.Minute

	rGuard := NewReplayGuardian(timeSource, maxTimeInPast, maxTimeInFuture)

	requestTimeInFuture := maxTimeInFuture + 1
	requestTime := now.Add(requestTimeInFuture)

	err := rGuard.VerifyRequest(rand.Bytes(32), requestTime)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "too far in the future"))

	// Verify that nothing has been added to the observedHashes set.
	g := rGuard.(*replayGuardian)
	require.Zero(t, len(g.observedHashes))
	require.Zero(t, g.expirationQueue.Size())
}

func TestDuplicateRequests(t *testing.T) {
	rand := random.NewTestRandom()

	now := rand.Time()
	timeSource := func() time.Time {
		return now
	}

	maxTimeInPast := time.Duration(rand.Intn(5)+1) * time.Minute
	maxTimeInFuture := time.Duration(rand.Intn(5)+1) * time.Minute

	rGuard := NewReplayGuardian(timeSource, maxTimeInPast, maxTimeInFuture)
	submittedHashes := make(map[string]struct{})

	for i := 0; i < 5; i++ {
		now = rand.TimeInRange(now, now.Add(10*time.Second))

		// Submit a new request
		earliestLegalTime := now.Add(-maxTimeInPast)
		latestLegalTime := now.Add(maxTimeInFuture)

		hash := rand.Bytes(32)
		var requestTime time.Time

		choice := rand.Float64()
		if choice < 0.05 {
			// once in a while, choose a time that is the maximum time in the past
			requestTime = earliestLegalTime
		} else if choice < 0.1 {
			// once in a while, choose a time that is the maximum time in the future
			requestTime = latestLegalTime
		} else {
			// choose a time that is within the legal range
			requestTime = rand.TimeInRange(earliestLegalTime, latestLegalTime)
		}

		err := rGuard.VerifyRequest(hash, requestTime)
		require.NoError(t, err)
		submittedHashes[string(hash)] = struct{}{}

		if rand.Float64() < 0.01 {
			// Once in a while, scan through the submitted hashes and verify that they are all rejected.
			for submittedHash := range submittedHashes {
				err = rGuard.VerifyRequest([]byte(submittedHash), requestTime)
				require.Error(t, err)
			}
		}
	}

	// Move time forward a long time in order to prune all the hashes. Submit a single request to trigger cleanup.
	now = now.Add(maxTimeInPast + maxTimeInFuture + 1)

	err := rGuard.VerifyRequest(rand.Bytes(32), now)
	require.NoError(t, err)

	// Only the most recent hash should be in the observedHashes set.
	g := rGuard.(*replayGuardian)
	require.Equal(t, 1, len(g.observedHashes))
	require.Equal(t, 1, g.expirationQueue.Size())
}
