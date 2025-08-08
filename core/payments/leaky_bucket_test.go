package payments

import (
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

var testStartTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

// Tests the NewLeakyBucket constructor
func TestNewLeakyBucket(t *testing.T) {
	t.Run("create with valid parameters", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(10, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.NotNil(t, leakyBucket)
		require.NoError(t, err)
	})

	t.Run("create with invalid leak rate", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(-10, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "negative leak rate should cause error")

		leakyBucket, err = NewLeakyBucket(0, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "zero leak rate should cause error")
	})

	t.Run("create with invalid bucket size duration", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(10, -10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "negative bucket duration should cause error")

		leakyBucket, err = NewLeakyBucket(10, 0, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "zero bucket duration should cause error")
	})
}

// Test the Fill method
func TestFill(t *testing.T) {
	// verify that overfill logic is working as expected
	t.Run("test overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(11, 10*time.Second, BiasPermitMore, OverfillOncePermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity+10)
		require.NoError(t, err)
		require.Equal(t, leakyBucket.bucketCapacity+10, leakyBucket.currentFillLevel, "first overfill should succeed")

		// no time elapses, so bucket is still overfilled
		err = leakyBucket.Fill(testStartTime, 1)
		require.Error(t, err)
		var insufficientErr *InsufficientReservationCapacityError
		require.True(t, errors.As(err, &insufficientErr), "overfill should fail, if bucket is overfilled")

		// let some time elapse, so there is a little bit of available capacity
		err = leakyBucket.Fill(testStartTime.Add(time.Second), leakyBucket.bucketCapacity+10)
		require.NoError(t, err, "any available capacity should permit overfill")
	})

	// make sure filling without overfill works
	t.Run("valid non-overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity-10)
		require.NoError(t, err)

		require.Equal(t, leakyBucket.bucketCapacity-10, leakyBucket.currentFillLevel)
	})

	// test edge case of filling directly to capacity
	t.Run("fill to exact capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		err = leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity)
		require.NoError(t, err)
		require.Equal(t, leakyBucket.bucketCapacity, leakyBucket.currentFillLevel)
	})

	t.Run("fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.Fill(testStartTime, 0)
		require.Error(t, err, "zero fill should not be permitted")

		err = leakyBucket.Fill(testStartTime, -10)
		require.Error(t, err, "negative fill should not be permitted")

		require.Equal(t, int64(0), leakyBucket.currentFillLevel, "nothing should have been added to the bucket")
	})

	// tests that waiting a really long time leaks the bucket empty, and that filling after that behaves as expected
	t.Run("large idle leakage to empty", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		// wait longer than the bucket duration
		futureTime := testStartTime.Add(15 * time.Second)

		err = leakyBucket.Fill(futureTime, 50)
		require.NoError(t, err)

		require.Equal(t, int64(50), leakyBucket.currentFillLevel, "bucket should leak empty, then be filled")
	})
}

// Tests that revert fill works
func TestRevertFill(t *testing.T) {
	t.Run("valid revert fill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.Equal(t, int64(500), leakyBucket.currentFillLevel)

		err = leakyBucket.RevertFill(testStartTime, 200)
		require.NoError(t, err)

		require.Equal(t, int64(300), leakyBucket.currentFillLevel)
	})

	t.Run("revert fill resulting in 0 capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.Equal(t, int64(500), leakyBucket.currentFillLevel)

		// revert fill with greater than the amount in the bucket
		err = leakyBucket.RevertFill(testStartTime, 600)
		require.NoError(t, err)

		require.Equal(t, int64(0), leakyBucket.currentFillLevel, "revert fill should clamp to 0")
	})

	t.Run("revert fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.RevertFill(testStartTime, 0)
		require.Error(t, err, "revert fill with 0 symbols should cause an error")

		err = leakyBucket.RevertFill(testStartTime, -10)
		require.Error(t, err, "revert fill with negative symbols should cause an error")

		require.Equal(t, int64(0), leakyBucket.currentFillLevel)
	})
}

func TestLeak(t *testing.T) {
	t.Run("leak with 'permitMore' bias", func(t *testing.T) {
		leakTest(t, BiasPermitMore)
	})

	t.Run("leak with 'permitLess' bias", func(t *testing.T) {
		leakTest(t, BiasPermitLess)
	})
}

// this function does many leaks, and makes sure the end values match up with expected values
func leakTest(t *testing.T, bias BiasBehavior) {
	leakRate := int64(5)

	// This test uses a large capacity, to make sure that none of the fills or leaks are bumping up against the
	// limits of the bucket
	leakyBucket, err := NewLeakyBucket(leakRate, 10*time.Hour, bias, OverfillNotPermitted, testStartTime)
	require.NotNil(t, leakyBucket)
	require.NoError(t, err)

	// We set the bucket fill to half way, so we're far away from both full and empty
	halfFull := leakyBucket.bucketCapacity / 2
	leakyBucket.currentFillLevel = halfFull

	testRandom := random.NewTestRandom()
	iterations := 1000

	workingTime := testStartTime
	for range iterations {
		// randomly advance between 1 nanosecond and 2 seconds for each iteration
		workingTime = workingTime.Add(time.Duration(testRandom.Intn(2_000_000_000) + 1))

		err := leakyBucket.Fill(workingTime, 1)
		require.NoError(t, err)
	}

	// round time off to a full second, so we can predict what the expected leak should be
	// it's easy math if the total test time is a round second
	workingTime = workingTime.Add(time.Duration(1e9-workingTime.Nanosecond()) * time.Nanosecond)
	require.Equal(t, 0, workingTime.Nanosecond(), "bug in test logic: workingTime should be a round second")

	// do a final fill, to bring the total test time to the round second value
	err = leakyBucket.Fill(workingTime, 1)
	require.NoError(t, err)

	// compute how much should have leaked throughout the test duration
	timeDelta := workingTime.Sub(testStartTime)
	expectedLeak := int64(timeDelta.Seconds()) * leakRate

	// original fill, minus what we expected to leak, plus what we filled during iteration (and the final fill)
	expectedFill := halfFull - expectedLeak + int64(iterations+1)

	require.Equal(t, expectedFill, leakyBucket.currentFillLevel, "fill level didn't match expected")
}

// Tests that time going backwards throws the right error
func TestTimeRegression(t *testing.T) {
	leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
	require.NoError(t, err)

	err = leakyBucket.Fill(testStartTime.Add(5*time.Second), 100)
	require.NoError(t, err)

	err = leakyBucket.Fill(testStartTime.Add(3*time.Second), 50)
	require.Error(t, err)
	var timeErr *TimeMovedBackwardError
	require.True(t, errors.As(err, &timeErr))

	err = leakyBucket.RevertFill(testStartTime.Add(2*time.Second), 50)
	require.Error(t, err)
	require.True(t, errors.As(err, &timeErr))
}

// Directly meddles with the leak function, to do a sanity check that rounding is happening as expected, based on the
// configured bias
func TestPartialSecondRoundingDifference(t *testing.T) {
	t.Run("BiasPermitMore rounds up", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(7, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		leakyBucket.currentFillLevel = 100

		// 25% of 7, rounded up, is 2
		err = leakyBucket.leak(testStartTime.Add(250 * time.Millisecond))
		require.NoError(t, err)
		require.Equal(t, int64(98), leakyBucket.currentFillLevel)
	})

	t.Run("BiasPermitLess rounds down", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(7, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		leakyBucket.currentFillLevel = 100

		// 25% of 7, rounded down, is 1
		err = leakyBucket.leak(testStartTime.Add(250 * time.Millisecond))
		require.NoError(t, err)
		require.Equal(t, int64(99), leakyBucket.currentFillLevel)
	})
}
