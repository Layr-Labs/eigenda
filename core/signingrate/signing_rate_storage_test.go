package signingrate

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

// simulateRandomSigningRateActivity simulates random signing activity on the given tracker. Does not attempt to
// advance time.
func simulateRandomSigningRateActivity(
	rand *random.TestRandom,
	tracker SigningRateTracker,
	quorumCount core.QuorumID,
	validatorIDs []core.OperatorID,
	iterations int,
) {

	for i := 0; i < iterations; i++ {
		quorum := core.QuorumID(rand.Intn(int(quorumCount)))
		validator := validatorIDs[rand.Intn(len(validatorIDs))]
		batchSize := uint64(rand.Intn(10) + 1)
		signingLatency := time.Duration(rand.Intn(1000)) * time.Millisecond

		if rand.Bool() {
			tracker.ReportSuccess(quorum, validator, batchSize, signingLatency)
		} else {
			tracker.ReportFailure(quorum, validator, batchSize)
		}
	}
}

// Compare two ValidatorSigningRate slices for equality.
func areValidatorSigningRatesEqual(
	expected []*validator.SigningRateBucket,
	actual []*validator.SigningRateBucket,
) bool {

	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		expectedBucket := expected[i]
		actualBucket := actual[i]

		if expectedBucket.GetStartTimestamp() != actualBucket.GetStartTimestamp() ||
			expectedBucket.GetEndTimestamp() != actualBucket.GetEndTimestamp() {
			return false
		}

		if len(expectedBucket.GetQuorumSigningRates()) != len(actualBucket.GetQuorumSigningRates()) {
			return false
		}

		for j := range expectedBucket.GetQuorumSigningRates() {
			expectedQuorum := expectedBucket.GetQuorumSigningRates()[j]
			actualQuorum := actualBucket.GetQuorumSigningRates()[j]

			if expectedQuorum.GetQuorumId() != actualQuorum.GetQuorumId() {
				return false
			}

			if len(expectedQuorum.GetValidatorSigningRates()) != len(actualQuorum.GetValidatorSigningRates()) {
				return false
			}

			for k := range expectedQuorum.GetValidatorSigningRates() {
				expectedValidator := expectedQuorum.GetValidatorSigningRates()[k]
				actualValidator := actualQuorum.GetValidatorSigningRates()[k]

				if !areSigningRatesEqual(expectedValidator, actualValidator) {
					return false
				}
			}
		}
	}

	return true
}

// Setting up local stack is slow. Cram a bunch of test cases into this one test to avoid this cost.
func TestSigningRateStorage(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	cleanup, err := test.DeployDynamoLocalstack(t.Context())
	require.NoError(t, err)
	defer cleanup()

	dynamoClient, err := test.GetDynamoClient()
	require.NoError(t, err)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	tableName := "TestSigningRateStorage"
	storage, err := NewDynamoSigningRateStorage(t.Context(), logger, dynamoClient, tableName)
	require.NoError(t, err)

	validatorCount := rand.Intn(10) + 5
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
	}

	quorumCount := core.QuorumID(rand.Intn(5) + 3)

	now := rand.Time()
	timePointer := atomic.Pointer[time.Time]{}
	timePointer.Store(&now)
	timeSource := func() time.Time {
		return *timePointer.Load()
	}

	// Use a signing rate tracker as a "source of truth". This data structure is validated in its own unit tests,
	// so trust it here.
	tracker, err := NewSigningRateTracker(logger, time.Hour*100, time.Minute*10, timeSource)
	require.NoError(t, err)

	// Check query behavior when there are no buckets.
	buckets, err := storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, buckets, 0)

	// Add a single bucket and check it can be retrieved.
	simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)

	unflushedBuckets, err := tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 1)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err := tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 1)
	actualBuckets, err := storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))

	// Add several more buckets.
	for i := 0; i < 5; i++ {
		now = now.Add(time.Minute * 10)
		timePointer.Store(&now)
		simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)
	}

	unflushedBuckets, err = tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 5)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err = tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 6)
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))

	// Query for a subset of the data.

	// Fetch data starting exactly at the start of a bucket.
	targetIndex := len(expectedBuckets) / 2
	startTimestamp := expectedBuckets[targetIndex].GetStartTimestamp()

	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(int64(startTimestamp), 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets[targetIndex:], actualBuckets))

	// If we subtract one second from the starting timestamp, we should snag the previous bucket as well.
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(int64(startTimestamp)-1, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets[targetIndex-1:], actualBuckets))

	// Modify the last bucket and ensure it gets overwritten correctly.
	// Note that we are not advancing time, so this activity goes into the last bucket.
	simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)

	unflushedBuckets, err = tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 1)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err = tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 6)
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)

	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))
}
