package signingrate

import (
	"bytes"
	"sort"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

// Validate information in the signing rate tracker against expected information.
func validateTracker(
	t *testing.T,
	now time.Time,
	expectedBuckets []*SigningRateBucket,
	tracker SigningRateTracker,
	timeSpan time.Duration,
) {

	cutoffTime := now.Add(-timeSpan)

	// Request all available buckets that are still before the cutoff time.
	dumpedBuckets, err := tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)

	if len(dumpedBuckets) == 0 {
		// It is ok to return zero dumped buckets iff no data has been added yet.
		require.Equal(t, 0, len(expectedBuckets[0].validatorInfo))
		return
	}

	// We shouldn't see any buckets that end before the cutoff time.
	for _, bucket := range dumpedBuckets {
		require.True(t, bucket.GetEndTimestamp() >= uint64(cutoffTime.Unix()))
	}

	// Find the index of the first expected bucket that ends after the cutoff time. This should align
	// with the first bucket in dumpedBuckets.
	index := 0
	for ; expectedBuckets[index].endTimestamp.Unix() < cutoffTime.Unix(); index++ {
	}

	expectedDumpSize := len(expectedBuckets) - index
	require.Equal(t, expectedDumpSize, len(dumpedBuckets))

	// For each remaining bucket, the expected bucket should exactly match the dumped bucket.
	for ; index < len(expectedBuckets); index++ {
		expectedBucket := expectedBuckets[index]
		dumpedBucket := dumpedBuckets[index]

		require.Equal(t, uint64(expectedBucket.startTimestamp.Unix()), dumpedBucket.GetStartTimestamp())
		require.Equal(t, uint64(expectedBucket.endTimestamp.Unix()), dumpedBucket.GetEndTimestamp())
		for _, signingRate := range dumpedBucket.GetValidatorSigningRates() {
			validatorID := core.OperatorID(signingRate.GetId())
			expectedSigningRate := expectedBucket.validatorInfo[validatorID]
			require.True(t, areSigningRatesEqual(expectedSigningRate, signingRate))
		}
	}

}

func randomOperationsTest(
	t *testing.T,
	tracker SigningRateTracker,
	bucketSpan time.Duration,
	timeSpan time.Duration,
) {
	defer tracker.Close()
	rand := random.NewTestRandom()

	validatorCount := rand.IntRange(1, 10)
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
	}

	// Sort validator IDs. This is the expected ordering within the protobuf.
	sort.Slice(validatorIDs, func(i, j int) bool {
		return bytes.Compare(validatorIDs[i][:], validatorIDs[j][:]) < 0
	})

	testSpan := timeSpan * 2
	totalBuckets := int(testSpan / bucketSpan)

	expectedBuckets := make([]*SigningRateBucket, 0, totalBuckets)

	// Each iteration, step forward in time by exactly one second.
	startTime := rand.Time()
	endTime := startTime.Add(testSpan)
	currentTime := startTime
	expectedBuckets = append(expectedBuckets, NewSigningRateBucket(startTime, bucketSpan))

	// verify before we've added any data
	validateTracker(t, currentTime, expectedBuckets, tracker, timeSpan)

	for currentTime.Before(endTime) {
		batchSize := rand.Uint64Range(1, 1000)
		validatorIndex := rand.Intn(validatorCount)
		validatorID := validatorIDs[validatorIndex]

		expectedBucket := expectedBuckets[len(expectedBuckets)-1]
		if !currentTime.Before(expectedBucket.endTimestamp) {
			// We've moved into a new bucket.

			newBucketStart := expectedBucket.endTimestamp
			expectedBucket = NewSigningRateBucket(newBucketStart, bucketSpan)
			expectedBuckets = append(expectedBuckets, expectedBucket)
		}

		if rand.Bool() {
			latency := rand.DurationRange(time.Second, time.Hour)
			tracker.ReportSuccess(currentTime, validatorID, batchSize, latency)
			expectedBucket.ReportSuccess(validatorID, batchSize, latency)
		} else {
			tracker.ReportFailure(currentTime, validatorID, batchSize)
			expectedBucket.ReportFailure(validatorID, batchSize)
		}

		// On average, validate once per bucket.
		if rand.Float64() < 1.0/(bucketSpan.Seconds()) {
			validateTracker(t, currentTime, expectedBuckets, tracker, timeSpan)
		}

		currentTime = currentTime.Add(time.Second)
		if currentTime.After(endTime) {
			// Do one last validation at the end of the test.
			validateTracker(t, currentTime, expectedBuckets, tracker, timeSpan)
		}
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
		tracker := NewSigningRateTracker(logger, timeSpan, bucketSpan, nil)
		randomOperationsTest(t, tracker, timeSpan, bucketSpan)
	})

	// TODO: we will need a flush operation to make this work...
	//t.Run("threadsafeSigningRateTracker", func(t *testing.T) {
	//	t.Parallel()
	//	tracker := NewSigningRateTracker(logger, timeSpan, bucketSpan, nil)
	//	tracker = NewThreadsafeSigningRateTracker(tracker)
	//	randomOperationsTest(t, tracker, timeSpan, bucketSpan)
	//})

}
