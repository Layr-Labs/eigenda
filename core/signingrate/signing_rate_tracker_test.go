package signingrate

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

// Do a dump of a tracker and validate the contents.
func validateTrackerDump(
	t *testing.T,
	now time.Time,
	expectedBuckets []*SigningRateBucket,
	tracker SigningRateTracker,
	timeSpan time.Duration,
	dumpStart time.Time,
) {
	gcThreshold := now.Add(-timeSpan)
	cutoffTime := gcThreshold
	if dumpStart.After(gcThreshold) {
		cutoffTime = dumpStart
	}

	// Request all available buckets that are still before the cutoff time.
	dumpedBuckets, err := tracker.GetSigningRateDump(dumpStart)
	require.NoError(t, err)

	if len(dumpedBuckets) == 0 {
		// It is ok to return zero dumped buckets iff no data has been added yet.
		require.Equal(t, 0, len(expectedBuckets[0].signingRateInfo))
		return
	}

	// We shouldn't see any buckets that end before the cutoff time.
	for _, bucket := range dumpedBuckets {
		require.True(t, bucket.GetEndTimestamp() >= uint64(cutoffTime.Unix()))
	}

	// Find the index of the first expected bucket that ends after the cutoff time. This should align
	// with the first bucket in dumpedBuckets.
	indexOffset := 0
	for expectedBuckets[indexOffset].endTimestamp.Unix() <= cutoffTime.Unix() {
		indexOffset++
	}

	expectedDumpSize := len(expectedBuckets) - indexOffset
	require.Equal(t, expectedDumpSize, len(dumpedBuckets))

	// For each remaining bucket, the expected bucket should exactly match the dumped bucket.
	for index := 0; index < len(expectedBuckets)-indexOffset; index++ {
		expectedBucket := expectedBuckets[index+indexOffset]
		dumpedBucket := dumpedBuckets[index]

		require.Equal(t, int(uint64(expectedBucket.startTimestamp.Unix())), int(dumpedBucket.GetStartTimestamp()))
		require.Equal(t, uint64(expectedBucket.endTimestamp.Unix()), dumpedBucket.GetEndTimestamp())
		for _, quorumInfo := range dumpedBucket.GetQuorumSigningRates() {
			quorumID := core.QuorumID(quorumInfo.GetQuorumId())
			for _, signingRate := range quorumInfo.GetValidatorSigningRates() {
				validatorID := core.OperatorID(signingRate.GetId())
				expectedSigningRate := expectedBucket.signingRateInfo[quorumID][validatorID]
				require.True(t, areSigningRatesEqual(expectedSigningRate, signingRate))
			}
		}
	}
}

// Validate information in the signing rate tracker against expected information.
func validateTracker(
	t *testing.T,
	now time.Time,
	expectedBuckets []*SigningRateBucket,
	validatorIDs []core.OperatorID,
	tracker SigningRateTracker,
	timeSpan time.Duration,
	rand *random.TestRandom,
	empty bool,
) {

	err := tracker.Flush()
	require.NoError(t, err)

	// Check the start timestamp of the last bucket.
	if empty {
		// We should get a zero timestamp if no data has been added yet.
		timestamp, err := tracker.GetLastBucketStartTime()
		require.NoError(t, err)
		require.True(t, timestamp.IsZero())
	} else {
		expectedTimestamp := expectedBuckets[len(expectedBuckets)-1].startTimestamp
		actualTimestamp, err := tracker.GetLastBucketStartTime()
		require.NoError(t, err)
		require.Equal(t, expectedTimestamp, actualTimestamp)
	}

	// Dump entire tracker.
	validateTrackerDump(t, now, expectedBuckets, tracker, timeSpan, time.Time{})

	// Choose a random cutoff time within the last timeSpan.
	cutoffTime := now.Add(-time.Duration(rand.Float64Range(0, float64(timeSpan))))
	validateTrackerDump(t, now, expectedBuckets, tracker, timeSpan, cutoffTime)

	// For a random validator and a random time span, verify reported validator signing rates.
	validatorIndex := rand.Intn(len(validatorIDs))
	validatorID := validatorIDs[validatorIndex]
	startTime := now.Add(-time.Duration(rand.Float64Range(0, float64(timeSpan))))
	// intentionally allow endTime to be after now
	endTime := startTime.Add(time.Duration(rand.Float64Range(0, float64(timeSpan))))

	expectedSigningRate := &validator.ValidatorSigningRate{
		Id: validatorID[:],
	}
	for _, bucket := range expectedBuckets {
		if bucket.endTimestamp.Before(startTime) {
			// This bucket is entirely before the requested time range.
			continue
		}
		if bucket.startTimestamp.After(endTime) || bucket.startTimestamp.Equal(endTime) {
			// This bucket is entirely after the requested time range.
			break
		}
		expectedSigningRate.SignedBatches += bucket.signingRateInfo[0][validatorID].GetSignedBatches()
		expectedSigningRate.SignedBytes += bucket.signingRateInfo[0][validatorID].GetSignedBytes()
		expectedSigningRate.UnsignedBatches += bucket.signingRateInfo[0][validatorID].GetUnsignedBatches()
		expectedSigningRate.UnsignedBytes += bucket.signingRateInfo[0][validatorID].GetUnsignedBytes()
		expectedSigningRate.SigningLatency += bucket.signingRateInfo[0][validatorID].GetSigningLatency()
	}

	reportedSigningRate, err := tracker.GetValidatorSigningRate(0, validatorID, startTime, endTime)
	require.NoError(t, err)

	require.True(t, areSigningRatesEqual(expectedSigningRate, reportedSigningRate))
}

// Copy recent updates into the clone and validate that it matches the original.
func validateTrackerClone(
	t *testing.T,
	now time.Time,
	expectedBuckets []*SigningRateBucket,
	validatorIDs []core.OperatorID,
	tracker SigningRateTracker,
	trackerClone SigningRateTracker,
	timeSpan time.Duration,
	rand *random.TestRandom,
	empty bool,
) {

	err := tracker.Flush()
	require.NoError(t, err)
	err = trackerClone.Flush()
	require.NoError(t, err)

	// Only request data from the clone starting at the last bucket start time it knows about.
	dumpStartTimestamp, err := trackerClone.GetLastBucketStartTime()
	require.NoError(t, err)

	dump, err := tracker.GetSigningRateDump(dumpStartTimestamp)
	require.NoError(t, err)
	for _, dumpedBucket := range dump {
		trackerClone.UpdateLastBucket(now, dumpedBucket)
	}

	validateTracker(t, now, expectedBuckets, validatorIDs, trackerClone, timeSpan, rand, empty)

	// The clone should never mark buckets as needing flushing.
	buckets, err := trackerClone.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Equal(t, 0, len(buckets))
}

// This function performs a number of random operations on a tracker, and verifies that it provides the expected
// information. It periodically clones the data to a "follower" tracker, and verifies that both trackers provide
// the same information.
func randomOperationsTest(
	t *testing.T,
	tracker SigningRateTracker,
	trackerClone SigningRateTracker,
	timeSpan time.Duration,
	bucketSpan time.Duration,
	timePointer *atomic.Pointer[time.Time],
) {
	rand := random.NewTestRandom()

	validatorCount := rand.IntRange(1, 10)
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
	}

	testSpan := timeSpan * 2
	totalBuckets := int(testSpan / bucketSpan)

	expectedBuckets := make([]*SigningRateBucket, 0, totalBuckets)

	// Each iteration, step forward in time by exactly one second.
	startTime := rand.Time()
	endTime := startTime.Add(testSpan)
	currentTime := startTime
	bucket, err := NewSigningRateBucket(startTime, bucketSpan)
	require.NoError(t, err)
	expectedBuckets = append(expectedBuckets, bucket)

	// verify before we've added any data
	validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, true)
	validateTrackerClone(t, currentTime, expectedBuckets, validatorIDs, tracker, trackerClone, timeSpan, rand, true)

	for currentTime.Before(endTime) {
		batchSize := rand.Uint64Range(1, 1000)
		validatorIndex := rand.Intn(validatorCount)
		validatorID := validatorIDs[validatorIndex]

		expectedBucket := expectedBuckets[len(expectedBuckets)-1]
		if !expectedBucket.Contains(currentTime) {
			// We've moved into a new bucket.
			expectedBucket, err = NewSigningRateBucket(currentTime, bucketSpan)
			require.NoError(t, err)
			expectedBuckets = append(expectedBuckets, expectedBucket)
		}

		// TODO use more than just quorum 0

		if rand.Bool() {
			latency := rand.DurationRange(time.Second, time.Hour)
			tracker.ReportSuccess(0, validatorID, batchSize, latency)
			expectedBucket.ReportSuccess(0, validatorID, batchSize, latency)
		} else {
			tracker.ReportFailure(0, validatorID, batchSize)
			expectedBucket.ReportFailure(0, validatorID, batchSize)
		}

		// On average, validate once per bucket.
		if rand.Float64() < 1.0/(bucketSpan.Seconds()) {
			validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, false)
			validateTrackerClone(
				t, currentTime, expectedBuckets, validatorIDs, tracker, trackerClone, timeSpan, rand, false)
		}

		nextTime := currentTime.Add(time.Second)
		if !nextTime.Before(endTime) {
			// Do one last validation at the end of the test.
			validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, false)
			validateTrackerClone(
				t, currentTime, expectedBuckets, validatorIDs, tracker, trackerClone, timeSpan, rand, false)
		}

		// There should be one unflushed bucket.
		buckets, err := tracker.GetUnflushedBuckets()
		require.NoError(t, err)
		require.Equal(t, 1, len(buckets))
		// Asking for unflushed buckets again should return none, since the first call marks them as flushed.
		buckets, err = tracker.GetUnflushedBuckets()
		require.NoError(t, err)
		require.Equal(t, 0, len(buckets))

		currentTime = nextTime
		timePointer.Store(&currentTime)
	}
}

func TestRandomOperations(t *testing.T) {
	t.Parallel()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	// The size of each bucket
	bucketSpan := time.Minute
	// The amount of time the tracker remembers data for
	timeSpan := bucketSpan * 100

	t.Run("signingRateTracker", func(t *testing.T) {
		t.Parallel()

		currentTime := &atomic.Pointer[time.Time]{}
		timeSource := func() time.Time {
			return *currentTime.Load()
		}

		tracker, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, timeSource)
		require.NoError(t, err)
		defer tracker.Close()

		trackerClone, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, timeSource)
		require.NoError(t, err)
		defer trackerClone.Close()

		randomOperationsTest(t, tracker, trackerClone, timeSpan, bucketSpan, currentTime)
	})

	t.Run("threadsafeSigningRateTracker", func(t *testing.T) {
		t.Parallel()

		currentTime := &atomic.Pointer[time.Time]{}
		timeSource := func() time.Time {
			return *currentTime.Load()
		}

		tracker, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, timeSource)
		require.NoError(t, err)
		tracker = NewThreadsafeSigningRateTracker(tracker)
		defer tracker.Close()

		trackerClone, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, nil)
		require.NoError(t, err)
		trackerClone = NewThreadsafeSigningRateTracker(trackerClone)
		defer trackerClone.Close()

		randomOperationsTest(t, tracker, trackerClone, timeSpan, bucketSpan, currentTime)
	})

}

func unflushedBucketsTest(
	t *testing.T,
	tracker SigningRateTracker,
	timeSpan time.Duration,
	bucketSpan time.Duration,
	timePointer *atomic.Pointer[time.Time],
) {
	rand := random.NewTestRandom()

	validatorCount := rand.IntRange(1, 10)
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
	}

	testSpan := timeSpan * 2
	totalBuckets := int(testSpan / bucketSpan)

	expectedBuckets := make([]*SigningRateBucket, 0, totalBuckets)

	// Each iteration, step forward in time by exactly one second.
	startTime := rand.Time()
	endTime := startTime.Add(testSpan)
	currentTime := startTime
	bucket, err := NewSigningRateBucket(startTime, bucketSpan)
	require.NoError(t, err)
	expectedBuckets = append(expectedBuckets, bucket)

	// verify before we've added any data
	validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, true)

	for currentTime.Before(endTime) {
		batchSize := rand.Uint64Range(1, 1000)
		validatorIndex := rand.Intn(validatorCount)
		validatorID := validatorIDs[validatorIndex]

		expectedBucket := expectedBuckets[len(expectedBuckets)-1]
		if !expectedBucket.Contains(currentTime) {
			// We've moved into a new bucket.
			expectedBucket, err = NewSigningRateBucket(currentTime, bucketSpan)
			require.NoError(t, err)
			expectedBuckets = append(expectedBuckets, expectedBucket)
		}

		if rand.Bool() {
			latency := rand.DurationRange(time.Second, time.Hour)
			tracker.ReportSuccess(0, validatorID, batchSize, latency)
			expectedBucket.ReportSuccess(0, validatorID, batchSize, latency)
		} else {
			tracker.ReportFailure(0, validatorID, batchSize)
			expectedBucket.ReportFailure(0, validatorID, batchSize)
		}

		// On average, validate once per bucket.
		if rand.Float64() < 1.0/(bucketSpan.Seconds()) {
			validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, false)
		}

		nextTime := currentTime.Add(time.Second)
		if !nextTime.Before(endTime) {
			// Do one last validation at the end of the test.
			validateTracker(t, currentTime, expectedBuckets, validatorIDs, tracker, timeSpan, rand, false)
		}

		// Unlike TestRandomOperations, wait until the end of the test to look at unflushed buckets.

		currentTime = nextTime
		timePointer.Store(&currentTime)
	}

	err = tracker.Flush()
	require.NoError(t, err)

	// Get unflushed buckets. This should exactly match expectedBuckets
	// (i.e. it should have all data written during this test).
	unflushedBuckets, err := tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Equal(t, len(expectedBuckets), len(unflushedBuckets))
	for i, bucket := range unflushedBuckets {
		expectedBucket := expectedBuckets[i]
		require.Equal(t, int(uint64(expectedBucket.startTimestamp.Unix())), int(bucket.GetStartTimestamp()))
		require.Equal(t, uint64(expectedBucket.endTimestamp.Unix()), bucket.GetEndTimestamp())
		for _, quorumInfo := range bucket.GetQuorumSigningRates() {
			quorumID := core.QuorumID(quorumInfo.GetQuorumId())
			for _, signingRate := range quorumInfo.GetValidatorSigningRates() {
				validatorID := core.OperatorID(signingRate.GetId())
				expectedSigningRate := expectedBucket.signingRateInfo[quorumID][validatorID]
				require.True(t, areSigningRatesEqual(expectedSigningRate, signingRate))
			}
		}
	}

	// There should no longer be any unflushed buckets.
	unflushedBuckets, err = tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Equal(t, 0, len(unflushedBuckets))
}

// Perform a bunch of random operations. At the end, request the unflushed buckets. We should see all data in the
// proper order.
func TestUnflushedBuckets(t *testing.T) {
	t.Parallel()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	// The size of each bucket
	bucketSpan := time.Minute
	// The amount of time the tracker remembers data for
	timeSpan := bucketSpan * 100

	t.Run("signingRateTracker", func(t *testing.T) {
		t.Parallel()

		currentTime := &atomic.Pointer[time.Time]{}
		timeSource := func() time.Time {
			return *currentTime.Load()
		}

		tracker, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, timeSource)
		require.NoError(t, err)
		defer tracker.Close()

		unflushedBucketsTest(t, tracker, timeSpan, bucketSpan, currentTime)
	})

	t.Run("threadsafeSigningRateTracker", func(t *testing.T) {
		t.Parallel()

		currentTime := &atomic.Pointer[time.Time]{}
		timeSource := func() time.Time {
			return *currentTime.Load()
		}

		tracker, err := NewSigningRateTracker(logger, timeSpan, bucketSpan, timeSource)
		require.NoError(t, err)
		tracker = NewThreadsafeSigningRateTracker(tracker)
		defer tracker.Close()

		unflushedBucketsTest(t, tracker, timeSpan, bucketSpan, currentTime)
	})
}
