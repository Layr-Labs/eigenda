package ejector

import (
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/test/random"
	geth "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TODO test cases:
//  - error checking if ejection is in progress
//  - rate limiting
//     - overfill
//     - start empty/full
//     - leak rate
//  - error checking if validator is in any quorum

// For a target trigger time, determine if it is time to trigger. Time to trigger is defined as the first
// timestamp that appears after the target time (which means that the previous time is before the target time).
func isTriggerTime(now time.Time, previousTime time.Time, target time.Time) bool {
	return now.After(target) && previousTime.Before(target)
}

func TestStandardEjection(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 3)
}

func TestConstructorBlacklist(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true

	// Blacklist B and C, so only A should be ejected.
	blacklist := []geth.Address{
		validatorB,
		validatorC,
	}

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			manager.BeginEjection(validatorB, nil)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			manager.BeginEjection(validatorC, nil)
		}
		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}

		// Neither B nor C should ever have their ejections started or finalized, since they are blacklisted.
		_, started := ejectionTransactor.inProgressEjections[validatorB]
		require.False(t, started)
		_, finalized := ejectionTransactor.completedEjections[validatorB]
		require.False(t, finalized)
		_, started = ejectionTransactor.inProgressEjections[validatorC]
		require.False(t, started)
		_, finalized = ejectionTransactor.completedEjections[validatorC]
		require.False(t, finalized)

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 1)
}

func TestEjectionAlreadyInProgress(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	// Mark the ejection for validator B as already in progress. If the ejection manager tries to start it again,
	// the mock transactor will raise an error.
	ejectionTransactor.inProgressEjections[validatorB] = struct{}{}

	// Verify that the mock transactor will raise an error if asked to start an ejection that is already in progress.
	err := ejectionTransactor.StartEjection(t.Context(), validatorB)
	require.Error(t, err)

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)
			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 3)
}

func TestMinimumTimeBetweenEjections(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Simulate an ejection for B before we start the main loop. The retry delay is configured to be > 10 minutes,
	// so the ejection scheduled below should be skipped
	manager.BeginEjection(validatorB, nil)
	currentTime = currentTime.Add(5 * time.Minute)
	manager.FinalizeEjections()
	// Put B "back into" a quorum so that it is eligible for ejection again.
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	delete(ejectionTransactor.completedEjections, validatorB)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			manager.BeginEjection(validatorB, nil)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		// Due to timing, the ejection for B should never be started in this loop.
		_, started := ejectionTransactor.inProgressEjections[validatorB]
		require.False(t, started)
		_, finalized := ejectionTransactor.completedEjections[validatorB]
		require.False(t, finalized)

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 2)

	// Fast-forward time so that B's retry delay has passed, then try again
	currentTime = currentTime.Add(10 * time.Minute)
	manager.BeginEjection(validatorB, nil)
	currentTime = currentTime.Add(10 * time.Minute)
	manager.FinalizeEjections()

	require.Len(t, ejectionTransactor.completedEjections, 3)

	// Fast-forward time again. Ensure that the ejection manager has forgotten about all prior ejections.
	// If we don't, it's a memory leak.
	currentTime = currentTime.Add(2 * time.Hour)
	manager.FinalizeEjections()
	require.Equal(t, 0, len(manager.(*ejectionManager).recentEjectionTimes))
}

func TestEjectedBySomebodyElse(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)

			// Before running the iteration that would otherwise eject B, simulate what happens if some other entity
			// finalizes the ejection before us.
			delete(ejectionTransactor.inProgressEjections, validatorB)
			ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = false

			// Asking the transactor to finalize an ejection for B should now return an error, since the
			// ejection is no longer in progress due to being finalized by somebody else. The mock
			// transactor keeps track of completed ejections, so we can verify that the mock transactor
			// will work as expected if the ejection manager tries to finalize the ejection incorrectly.
			err := ejectionTransactor.CompleteEjection(t.Context(), validatorB)
			require.Error(t, err)

		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 2)

	// Being ejected by somebody else shouldn't have been counted as a failed ejection attempt.
	require.Equal(t, uint32(0), manager.(*ejectionManager).failedEjectionAttempts[validatorB])
}

func TestCancellation(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := rand.DurationRange(0, time.Minute)
	retryAttempts := uint32(3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	ejectionTimeA := currentTime.Add(time.Minute)

	// Make a bunch of attempts at ejecting B. The first 3 will be cancelled, the last should not be attempted.
	ejectionTimeB1 := currentTime.Add(time.Minute)
	ejectionTimeB2 := ejectionTimeB1.Add(5 * time.Minute)
	ejectionTimeB3 := ejectionTimeB2.Add(5 * time.Minute)
	ejectionTimeB4 := ejectionTimeB3.Add(5 * time.Minute)

	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB1 time.Time
	var expectedFinalizeTimeB2 time.Time
	var expectedFinalizeTimeB3 time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB1) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB1 = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB2) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB2 = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB3) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB3 = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB4) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			// Since we've failed 3 times already, B should be in the blacklist. The ejection should not be started.
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB1) {
			// Simulate the ejection being cancelled.
			delete(ejectionTransactor.inProgressEjections, validatorB)

			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB2) {
			// Simulate the ejection being cancelled.
			delete(ejectionTransactor.inProgressEjections, validatorB)

			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB3) {
			// Simulate the ejection being cancelled.
			delete(ejectionTransactor.inProgressEjections, validatorB)

			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB1) {
			// Ejection was cancelled, so it shouldn't be finalized.
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB2) {
			// Ejection was cancelled, so it shouldn't be finalized.
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB3) {
			// Ejection was cancelled, so it shouldn't be finalized.
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 2)

	// B should be on the blacklist now.
	_, blacklisted := manager.(*ejectionManager).ejectionBlacklist[validatorB]
	require.True(t, blacklisted)
}

func TestEjectionInProgressError(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			// Make IsEjectionInProgress return an error for A.
			ejectionTransactor.isEjectionInProgressErrors[validatorA] = errors.New("intentional error")

			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)

			// Since there was an error checking if the ejection was in progress,
			// the ejection should not have been started.
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)

			// Allow ejection checks to proceed normally again.
			delete(ejectionTransactor.isEjectionInProgressErrors, validatorA)
			manager.BeginEjection(validatorA, nil)

			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)

			// When checking if the ejection is in progress, return an error for A.
			ejectionTransactor.isEjectionInProgressErrors[validatorA] = errors.New("intentional error")

		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			// Since there was an error checking if the ejection was in progress,
			// the ejection should not have been finalized for A.
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 2)

	// Failures to call eth transactions should not be treated as failed ejection attempts for the purposes
	// of blacklisting.
	require.Equal(t, uint32(0), manager.(*ejectionManager).failedEjectionAttempts[validatorA])
}

func TestStartEjectionError(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)

			// Simulate a failure when calling StartEjection for A.
			ejectionTransactor.startEjectionErrors[validatorA] = errors.New("intentional error")

			manager.BeginEjection(validatorA, nil)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)

			// Allow the second attempt to proceed normally.
			delete(ejectionTransactor.startEjectionErrors, validatorA)

			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 3)
}

func TestIsValidatorPresentInAnyQuorumError(t *testing.T) {
	rand := random.NewTestRandom()

	logger := common.TestLogger(t)

	start := rand.Time()
	currentTime := start
	previousTime := currentTime

	timeSource := func() time.Time {
		return currentTime
	}

	validatorA := rand.Address()
	validatorB := rand.Address()
	validatorC := rand.Address()

	ejectionTransactor := newMockEjectionTransactor()
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorA] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorB] = true
	ejectionTransactor.isValidatorPresentInAnyQuorumResponses[validatorC] = true

	ejectionDelay := time.Minute + rand.DurationRange(0, time.Minute)
	retryDelay := 10*time.Minute + rand.DurationRange(0, time.Minute)
	retryAttempts := rand.Uint32Range(1, 3)
	maxEjectionRate := 1.00
	bucketDuration := time.Hour
	startBucketFull := true
	var blacklist []geth.Address

	manager, err := NewEjectionManager(
		t.Context(),
		logger,
		timeSource,
		ejectionTransactor,
		ejectionDelay,
		retryDelay,
		retryAttempts,
		maxEjectionRate,
		bucketDuration,
		startBucketFull,
		blacklist)
	require.NoError(t, err)

	// Eject A and B at the same time. Eject C a bit later.
	ejectionTimeA := currentTime.Add(time.Minute)
	ejectionTimeB := currentTime.Add(time.Minute)
	ejectionTimeC := currentTime.Add(2 * time.Minute)

	var expectedFinalizeTimeA time.Time
	var expectedFinalizeTimeB time.Time
	var expectedFinalizeTimeC time.Time

	// Step forward in time in ~5 second increments, checking the state of ejections along the way.
	endTime := start.Add(30 * time.Minute)
	for currentTime.Before(endTime) {

		// Start ejections when ready.
		if isTriggerTime(currentTime, previousTime, ejectionTimeA) {
			_, started := ejectionTransactor.inProgressEjections[validatorA]
			require.False(t, started)
			manager.BeginEjection(validatorA, nil)
			expectedFinalizeTimeA = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorA]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeB) {
			_, started := ejectionTransactor.inProgressEjections[validatorB]
			require.False(t, started)
			manager.BeginEjection(validatorB, nil)
			expectedFinalizeTimeB = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorB]
			require.True(t, started)
		}
		if isTriggerTime(currentTime, previousTime, ejectionTimeC) {
			_, started := ejectionTransactor.inProgressEjections[validatorC]
			require.False(t, started)
			manager.BeginEjection(validatorC, nil)

			// Ejecting twice shouldn't harm anything. It will log, but otherwise be a no-op.
			manager.BeginEjection(validatorC, nil)

			expectedFinalizeTimeC = currentTime.Add(ejectionDelay)
			_, started = ejectionTransactor.inProgressEjections[validatorC]
			require.True(t, started)
		}

		// If right before the expected finalize time, ejection should not yet be finalized.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {

			// Simulate the ejection being cancelled, but IsValidatorPresentInAnyQuorum returning an error.
			delete(ejectionTransactor.inProgressEjections, validatorB)
			ejectionTransactor.isValidatorPresentInAnyQuorumErrors[validatorB] = errors.New("intentional error")

			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.False(t, finalized)
		}

		// Call this each iteration. Most of the time it won't do anything, but when the time is right it will finalize
		// ejections that are ready.
		manager.FinalizeEjections()

		// Once finalize is called, verify that the ejection has been completed if it is the expected time.
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeA) {
			_, finalized := ejectionTransactor.completedEjections[validatorA]
			require.True(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeB) {
			_, finalized := ejectionTransactor.completedEjections[validatorB]
			require.False(t, finalized)
		}
		if isTriggerTime(currentTime, previousTime, expectedFinalizeTimeC) {
			_, finalized := ejectionTransactor.completedEjections[validatorC]
			require.True(t, finalized)
		}

		previousTime = currentTime
		currentTime = currentTime.Add(rand.DurationRange(time.Second, 5*time.Second))
	}

	// Sanity check: we should see all three ejections completed. This is more a verification that the unit
	// test itself worked as expected, rather than a test of the ejection manager.
	require.Len(t, ejectionTransactor.completedEjections, 2)
}
