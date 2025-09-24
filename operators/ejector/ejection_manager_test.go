package ejector

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/test/random"
	geth "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TODO test cases:
//  - blacklisted in constructor
//  - ejection already in progress
//  - too recent of an attempt
//  - ejection already in progress
//  - error checking if ejection is in progress
//  - rate limiting
//     - overfill
//     - start empty/full
//     - leak rate
//  - validator leaves quorum prior to being ejected
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
