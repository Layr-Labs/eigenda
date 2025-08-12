package reservation

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
		leakyBucket, err := NewLeakyBucket(0, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
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
	// verify logic with overfill enabled
	t.Run("test overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(11, 10*time.Second, BiasPermitMore, OverfillOncePermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, uint32(leakyBucket.bucketCapacity+10))
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity+10, leakyBucket.currentFillLevel, "first overfill should succeed")

		// no time elapses, so bucket is still over capacity
		success, err = leakyBucket.Fill(testStartTime, 1)
		require.NoError(t, err)
		require.False(t, success, "overfill should fail, if bucket is already over capacity")

		// let some time elapse, so there is a little bit of available capacity
		success, err = leakyBucket.Fill(testStartTime.Add(time.Second), uint32(leakyBucket.bucketCapacity+10))
		require.NoError(t, err)
		require.True(t, success, "any available capacity should permit overfill")
	})

	// make sure filling without overfill enabled works
	t.Run("non-overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, uint32(leakyBucket.bucketCapacity-10))
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity-10, leakyBucket.currentFillLevel)

		success, err = leakyBucket.Fill(testStartTime, 11)
		require.NoError(t, err)
		require.False(t, success, "if no overfill is enabled, any amount of overfill should fail")
		require.Equal(t, leakyBucket.bucketCapacity-10, leakyBucket.currentFillLevel)
	})

	t.Run("fill to exact capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		success, err := leakyBucket.Fill(testStartTime, uint32(leakyBucket.bucketCapacity))
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity, leakyBucket.currentFillLevel)
	})

	t.Run("fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 0)
		require.Error(t, err, "zero fill should not be permitted")
		require.False(t, success)

		require.Equal(t, uint64(0), leakyBucket.currentFillLevel, "nothing should have been added to the bucket")
	})

	// tests that waiting a really long time leaks the bucket empty, and that filling after that behaves as expected
	t.Run("large idle leakage to empty", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		// wait longer than the bucket duration
		success, err := leakyBucket.Fill(testStartTime.Add(15*time.Second), 50)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, uint64(50), leakyBucket.currentFillLevel, "bucket should leak empty, then be filled")
	})
}

// Tests that revert fill works
func TestRevertFill(t *testing.T) {
	t.Run("valid revert fill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, uint64(500), leakyBucket.currentFillLevel)

		err = leakyBucket.RevertFill(testStartTime, 200)
		require.NoError(t, err)

		require.Equal(t, uint64(300), leakyBucket.currentFillLevel)
	})

	t.Run("revert fill resulting in 0 capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, uint64(500), leakyBucket.currentFillLevel)

		// revert fill with greater than the amount in the bucket
		err = leakyBucket.RevertFill(testStartTime, 600)
		require.NoError(t, err)

		require.Equal(t, uint64(0), leakyBucket.currentFillLevel, "revert fill should clamp to 0")
	})

	t.Run("revert fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.RevertFill(testStartTime, 0)
		require.Error(t, err, "revert fill with 0 symbols should cause an error")

		require.Equal(t, uint64(0), leakyBucket.currentFillLevel)
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
	leakRate := uint64(5)

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

		success, err := leakyBucket.Fill(workingTime, 1)
		require.NoError(t, err)
		require.True(t, success)
	}

	// round time off to a full second, so we can predict what the expected leak should be.
	// it's easy math if the total test time is a round second
	workingTime = workingTime.Add(time.Duration(1e9-workingTime.Nanosecond()) * time.Nanosecond)
	require.Equal(t, 0, workingTime.Nanosecond(), "bug in test logic: workingTime should be a round second")

	// do a final fill, to bring the total test time to the round second value
	success, err := leakyBucket.Fill(workingTime, 1)
	require.NoError(t, err)
	require.True(t, success)

	// compute how much should have leaked throughout the test duration
	timeDelta := workingTime.Sub(testStartTime)
	expectedLeak := uint64(timeDelta.Seconds()) * leakRate

	// original fill, minus what we expected to leak, plus what we filled during iteration (and the final fill)
	expectedFill := halfFull - expectedLeak + uint64(iterations+1)

	require.Equal(t, expectedFill, leakyBucket.currentFillLevel, "fill level didn't match expected")
}

// Tests that time going backwards throws the right error
func TestTimeRegression(t *testing.T) {
	leakyBucket, err := NewLeakyBucket(100, 10*time.Second, BiasPermitMore, OverfillNotPermitted, testStartTime)
	require.NoError(t, err)

	success, err := leakyBucket.Fill(testStartTime.Add(5*time.Second), 100)
	require.NoError(t, err)
	require.True(t, success)

	success, err = leakyBucket.Fill(testStartTime.Add(3*time.Second), 50)
	require.Error(t, err)
	require.False(t, success)
	require.True(t, errors.Is(err, ErrTimeMovedBackward))

	err = leakyBucket.RevertFill(testStartTime.Add(2*time.Second), 50)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrTimeMovedBackward))
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
		require.Equal(t, uint64(98), leakyBucket.currentFillLevel)
	})

	t.Run("BiasPermitLess rounds down", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(7, 10*time.Second, BiasPermitLess, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		leakyBucket.currentFillLevel = 100

		// 25% of 7, rounded down, is 1
		err = leakyBucket.leak(testStartTime.Add(250 * time.Millisecond))
		require.NoError(t, err)
		require.Equal(t, uint64(99), leakyBucket.currentFillLevel)
	})
}
