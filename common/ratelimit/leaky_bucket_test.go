package ratelimit

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

var testStartTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

func TestNewLeakyBucket(t *testing.T) {
	t.Run("create with valid parameters", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(10, 10*time.Second, true, OverfillNotPermitted, testStartTime)
		require.NotNil(t, leakyBucket)
		require.NoError(t, err)
	})

	t.Run("create with invalid leak rate", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(0, 10*time.Second, true, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "zero leak rate should cause error")
	})

	t.Run("create with invalid bucket size duration", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(10, -10*time.Second, true, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "negative bucket duration should cause error")

		leakyBucket, err = NewLeakyBucket(10, 0, true, OverfillNotPermitted, testStartTime)
		require.Nil(t, leakyBucket)
		require.Error(t, err, "zero bucket duration should cause error")
	})
}

func TestFill(t *testing.T) {
	t.Run("test overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(11, 10*time.Second, false, OverfillOncePermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity+10)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity+10, leakyBucket.currentFillLevel, "first overfill should succeed")

		// no time elapses, so bucket is still over capacity
		success, err = leakyBucket.Fill(testStartTime, 1)
		require.NoError(t, err)
		require.False(t, success, "overfill should fail, if bucket is already over capacity")

		// let some time elapse, so there is a little bit of available capacity
		success, err = leakyBucket.Fill(testStartTime.Add(time.Second), leakyBucket.bucketCapacity+10)
		require.NoError(t, err)
		require.True(t, success, "any available capacity should permit overfill")
	})

	t.Run("non-overfill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity-10)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity-10, leakyBucket.currentFillLevel)

		success, err = leakyBucket.Fill(testStartTime, 11)
		require.NoError(t, err)
		require.False(t, success, "if no overfill is enabled, any amount of overfill should fail")
		require.Equal(t, leakyBucket.bucketCapacity-10, leakyBucket.currentFillLevel)
	})

	t.Run("fill to exact capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		success, err := leakyBucket.Fill(testStartTime, leakyBucket.bucketCapacity)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, leakyBucket.bucketCapacity, leakyBucket.currentFillLevel)
	})

	t.Run("fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 0)
		require.Error(t, err, "zero fill should not be permitted")
		require.False(t, success)

		require.Equal(t, float64(0), leakyBucket.currentFillLevel, "nothing should have been added to the bucket")
	})

	// tests that waiting a really long time leaks the bucket empty, and that filling after that behaves as expected
	t.Run("large idle leakage to empty", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, true, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)

		// wait longer than the bucket duration
		success, err := leakyBucket.Fill(testStartTime.Add(15*time.Second), 50)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, float64(50), leakyBucket.currentFillLevel, "bucket should leak empty, then be filled")
	})
}

func TestRevertFill(t *testing.T) {
	t.Run("valid revert fill", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, float64(500), leakyBucket.currentFillLevel)

		err = leakyBucket.RevertFill(testStartTime, 200)
		require.NoError(t, err)

		require.Equal(t, float64(300), leakyBucket.currentFillLevel)
	})

	t.Run("revert fill resulting in 0 capacity", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		success, err := leakyBucket.Fill(testStartTime, 500)
		require.NoError(t, err)
		require.True(t, success)
		require.Equal(t, float64(500), leakyBucket.currentFillLevel)

		// revert fill with greater than the amount in the bucket
		err = leakyBucket.RevertFill(testStartTime, 600)
		require.NoError(t, err)

		require.Equal(t, float64(0), leakyBucket.currentFillLevel, "revert fill should clamp to 0")
	})

	t.Run("revert fill with invalid symbol count", func(t *testing.T) {
		leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
		require.NoError(t, err)
		require.NotNil(t, leakyBucket)

		err = leakyBucket.RevertFill(testStartTime, 0)
		require.Error(t, err, "revert fill with 0 symbols should cause an error")

		require.Equal(t, float64(0), leakyBucket.currentFillLevel)
	})
}

func TestLeak(t *testing.T) {
	leakRate := float64(5)

	// This test uses a large capacity, to make sure that none of the fills or leaks are bumping up against the
	// limits of the bucket
	leakyBucket, err := NewLeakyBucket(leakRate, 10*time.Hour, true, OverfillNotPermitted, testStartTime)
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

	// compute how much should have leaked throughout the test duration
	timeDelta := workingTime.Sub(testStartTime)
	expectedLeak := timeDelta.Seconds() * leakRate

	// original fill, minus what we expected to leak, plus what we filled during iteration
	expectedFill := halfFull - expectedLeak + float64(iterations)

	require.InDelta(t, expectedFill, leakyBucket.currentFillLevel, 0.0001, "fill level didn't match expected")
}

func TestTimeRegression(t *testing.T) {
	leakyBucket, err := NewLeakyBucket(100, 10*time.Second, false, OverfillNotPermitted, testStartTime)
	require.NoError(t, err)

	success, err := leakyBucket.Fill(testStartTime.Add(5*time.Second), 100)
	require.NoError(t, err)
	require.True(t, success)

	var timeMovedBackwardError *TimeMovedBackwardError

	success, err = leakyBucket.Fill(testStartTime.Add(3*time.Second), 50)
	require.Error(t, err)
	require.False(t, success)
	require.ErrorAs(t, err, &timeMovedBackwardError)

	err = leakyBucket.RevertFill(testStartTime.Add(2*time.Second), 50)
	require.Error(t, err)
	require.ErrorAs(t, err, &timeMovedBackwardError)
}

func TestReconfigure(t *testing.T) {
	rand := random.NewTestRandom()

	leakyBucket, err := NewLeakyBucket(1, 11*time.Second, false, OverfillOncePermitted, testStartTime)
	require.NoError(t, err)
	require.NotNil(t, leakyBucket)

	now := rand.Time()

	// Leak a few times, do not advance time. All should pass.
	for i := 1; i <= 6; i++ {
		success, err := leakyBucket.Fill(now, 2)
		require.NoError(t, err)
		require.True(t, success)
	}

	// We are currently overfilled, so we should be unable to fill any more.
	success, err := leakyBucket.Fill(now, 1)
	require.NoError(t, err)
	require.False(t, success, "overfill should not be permitted when already overfilled")

	fillLevel, err := leakyBucket.GetFillLevel(now)
	require.NoError(t, err)
	require.Equal(t, 12.0, fillLevel)

	// Advance time by 5 seconds, should leak 5 symbols.
	now = now.Add(5 * time.Second)

	// At this point in time, the expected fill level is 7.
	// Resize the leak rate to 2 symbols per second, and bucket duration to 1 second.
	// Resulting bucket size is 2 symbols, so we should be overfilled.
	err = leakyBucket.Reconfigure(2, 1*time.Second, OverfillNotPermitted, now)
	require.NoError(t, err)

	fillLevel, err = leakyBucket.GetFillLevel(now)
	require.NoError(t, err)
	require.Equal(t, 7.0, fillLevel, "fill level should be unchanged by reconfigure")

	// Wait 3 seconds, should leak 5 symbols, for a resulting fill level of 1.
	now = now.Add(3 * time.Second)
	fillLevel, err = leakyBucket.GetFillLevel(now)
	require.NoError(t, err)
	require.Equal(t, 1.0, fillLevel, "fill level should be 1 after leaking")

	// We toggled off overfill, so we should not be able to fill beyond capacity.
	success, err = leakyBucket.Fill(now, 2)
	require.NoError(t, err)
	require.False(t, success, "overfill should not be permitted")

	// Now, increase the bucket size to 10 symbols, and enable overfill once again.
	err = leakyBucket.Reconfigure(2, 5*time.Second, OverfillOncePermitted, now)
	require.NoError(t, err)

	fillLevel, err = leakyBucket.GetFillLevel(now)
	require.NoError(t, err)
	require.Equal(t, 1.0, fillLevel, "fill level should be unchanged by reconfigure")

	// We should be able to fill up to 10 symbols now.
	success, err = leakyBucket.Fill(now, 9)
	require.NoError(t, err)
	require.True(t, success, "fill within capacity should be permitted")

	fillLevel, err = leakyBucket.GetFillLevel(now)
	require.NoError(t, err)
	require.Equal(t, 10.0, fillLevel, "fill level should be 10 after fill")

	// Let a little drain away to verify that we can overfill again.
	now = now.Add(1 * time.Second)
	success, err = leakyBucket.Fill(now, 100)
	require.NoError(t, err)
	require.True(t, success, "overfill should be permitted again")
}
