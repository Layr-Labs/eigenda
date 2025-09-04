package signingrate

import (
	"bytes"
	"sort"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

// Returns true if two validator.ValidatorSigningRate messages are equal.
func areSigningRatesEqual(a *validator.ValidatorSigningRate, b *validator.ValidatorSigningRate) bool {
	if a == nil || b == nil {
		return a == b
	}
	if !bytes.Equal(a.GetId(), b.GetId()) {
		return false
	}
	if a.GetSignedBatches() != b.GetSignedBatches() {
		return false
	}
	if a.GetSignedBytes() != b.GetSignedBytes() {
		return false
	}
	if a.GetUnsignedBytes() != b.GetUnsignedBytes() {
		return false
	}
	if a.GetSigningLatency() != b.GetSigningLatency() {
		return false
	}
	return true
}

func TestProtoConversion(t *testing.T) {
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

	span := rand.DurationRange(time.Second, time.Hour)
	bucket, err := NewSigningRateBucket(rand.Time(), span)
	require.NoError(t, err)

	for _, validatorID := range validatorIDs {
		bucket.validatorInfo[validatorID] = &validator.ValidatorSigningRate{
			Id:             validatorID[:],
			SignedBatches:  rand.Uint64(),
			SignedBytes:    rand.Uint64(),
			UnsignedBytes:  rand.Uint64(),
			SigningLatency: rand.Uint64(),
		}
	}

	// Convert the entire bucket to a protobuf
	pb := bucket.ToProtobuf()
	require.Equal(t, uint64(bucket.startTimestamp.Unix()), pb.GetStartTimestamp())
	require.Equal(t, uint64(bucket.endTimestamp.Unix()), pb.GetEndTimestamp())
	for index := range pb.GetValidatorSigningRates() {
		expected := bucket.validatorInfo[validatorIDs[index]]
		actual := pb.GetValidatorSigningRates()[index]
		require.True(t, areSigningRatesEqual(expected, actual))
		require.True(t, expected != actual, "Expected a deep copy of the signing rate info")
	}

	// Getting the protobuf again should yield the same object (cached)
	pb2 := bucket.ToProtobuf()
	require.True(t, pb == pb2, "Expected the cached protobuf to be returned")

	// Convert protobuf back into a bucket
	bucket2 := NewBucketFromProto(pb)
	require.Equal(t, bucket.startTimestamp.Unix(), bucket2.startTimestamp.Unix())
	require.Equal(t, bucket.endTimestamp.Unix(), bucket2.endTimestamp.Unix())
	for id, info := range bucket.validatorInfo {
		info2, exists := bucket2.validatorInfo[id]
		require.True(t, exists, "Validator ID missing in converted bucket")
		require.True(t, areSigningRatesEqual(info, info2))
		require.True(t, info != info2, "Expected a deep copy of the signing rate info")
	}

	// Perform updates. This should clear the cached protobuf.
	bucket.ReportSuccess(validatorIDs[0], 0, 0)
	pb3 := bucket.ToProtobuf()
	require.True(t, pb3 != pb, "Expected a new protobuf to be generated after the bucket was modified")
	pb4 := bucket.ToProtobuf()
	require.True(t, pb3 == pb4, "Expected the cached protobuf to be returned")
	bucket.ReportFailure(validatorIDs[0], 0)
	pb5 := bucket.ToProtobuf()
	require.True(t, pb5 != pb4, "Expected a new protobuf to be generated after the bucket was modified")
	pb6 := bucket.ToProtobuf()
	require.True(t, pb5 == pb6, "Expected the cached protobuf to be returned")
}

func TestReporting(t *testing.T) {
	rand := random.NewTestRandom()

	expectedSuccesses := make(map[core.OperatorID]uint64)
	expectedFailures := make(map[core.OperatorID]uint64)
	expectedSuccessBytes := make(map[core.OperatorID]uint64)
	expectedFailureBytes := make(map[core.OperatorID]uint64)
	expectedLatency := make(map[core.OperatorID]uint64)

	validatorCount := rand.IntRange(1, 10)
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
		expectedSuccesses[validatorIDs[i]] = 0
		expectedFailures[validatorIDs[i]] = 0
		expectedSuccessBytes[validatorIDs[i]] = 0
		expectedFailureBytes[validatorIDs[i]] = 0
		expectedLatency[validatorIDs[i]] = 0
	}

	span := rand.DurationRange(time.Second, time.Hour)
	bucket, err := NewSigningRateBucket(rand.Time(), span)
	require.NoError(t, err)

	// Simulate a bunch of random reports.
	for i := 0; i < 10_000; i++ {
		batchSize := rand.Uint64Range(1, 1000)
		validatorIndex := rand.Intn(validatorCount)
		validatorID := validatorIDs[validatorIndex]

		if rand.Bool() {
			latency := rand.DurationRange(time.Second, time.Hour)
			bucket.ReportSuccess(validatorID, batchSize, latency)

			expectedSuccesses[validatorID] += 1
			expectedSuccessBytes[validatorID] += batchSize
			expectedLatency[validatorID] += uint64(latency.Nanoseconds())
		} else {
			bucket.ReportFailure(validatorID, batchSize)

			expectedFailures[validatorID] += 1
			expectedFailureBytes[validatorID] += batchSize
		}
	}

	// Verify the results.
	for _, validatorID := range validatorIDs {
		signingRate := bucket.getValidator(validatorID)
		require.Equal(t, expectedSuccesses[validatorID], signingRate.GetSignedBatches())
		require.Equal(t, expectedSuccessBytes[validatorID], signingRate.GetSignedBytes())
		require.Equal(t, expectedFailures[validatorID], signingRate.GetUnsignedBatches())
		require.Equal(t, expectedFailureBytes[validatorID], signingRate.GetUnsignedBytes())
		require.Equal(t, expectedLatency[validatorID], signingRate.GetSigningLatency())
	}
}

func TestCloneValidatorSigningRate(t *testing.T) {
	rand := random.NewTestRandom()

	signingRate := &validator.ValidatorSigningRate{
		Id:             rand.Bytes(32),
		SignedBatches:  rand.Uint64(),
		SignedBytes:    rand.Uint64(),
		UnsignedBytes:  rand.Uint64(),
		SigningLatency: rand.Uint64(),
	}

	clone := cloneValidatorSigningRate(signingRate)
	require.True(t, areSigningRatesEqual(signingRate, clone))
}
