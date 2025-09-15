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

	quorumCount := core.QuorumID(5)
	for quorum := core.QuorumID(0); quorum < quorumCount; quorum++ {
		bucket.signingRateInfo[quorum] = make(map[core.OperatorID]*validator.ValidatorSigningRate)
		for _, validatorID := range validatorIDs {
			bucket.signingRateInfo[quorum][validatorID] = &validator.ValidatorSigningRate{
				Id:             validatorID[:],
				SignedBatches:  rand.Uint64(),
				SignedBytes:    rand.Uint64(),
				UnsignedBytes:  rand.Uint64(),
				SigningLatency: rand.Uint64(),
			}
		}
	}

	// Convert the entire bucket to a protobuf
	pb := bucket.ToProtobuf()
	require.Equal(t, uint64(bucket.startTimestamp.Unix()), pb.GetStartTimestamp())
	require.Equal(t, uint64(bucket.endTimestamp.Unix()), pb.GetEndTimestamp())
	for _, quorumInfo := range pb.GetQuorumSigningRates() {
		quorumID := core.QuorumID(quorumInfo.GetQuorumId())
		for index, actualSigningRate := range quorumInfo.GetValidatorSigningRates() {
			expected := bucket.signingRateInfo[quorumID][validatorIDs[index]]

			require.True(t, areSigningRatesEqual(expected, actualSigningRate))
			require.True(t, expected != actualSigningRate, "Expected a deep copy of the signing rate info")
		}
	}

	// Getting the protobuf again should yield the same object (cached)
	pb2 := bucket.ToProtobuf()
	require.True(t, pb == pb2, "Expected the cached protobuf to be returned")

	// Convert protobuf back into a bucket
	bucket2 := NewBucketFromProto(pb)
	require.Equal(t, bucket.startTimestamp.Unix(), bucket2.startTimestamp.Unix())
	require.Equal(t, bucket.endTimestamp.Unix(), bucket2.endTimestamp.Unix())
	for quorum := core.QuorumID(0); quorum < quorumCount; quorum++ {
		for id, info := range bucket.signingRateInfo[quorum] {
			info2, exists := bucket2.signingRateInfo[quorum][id]
			require.True(t, exists, "Validator ID missing in converted bucket")
			require.True(t, areSigningRatesEqual(info, info2))
			require.True(t, info != info2, "Expected a deep copy of the signing rate info")
		}
	}

	// Perform updates. This should clear the cached protobuf.
	bucket.ReportSuccess(0, validatorIDs[0], 0, 0)
	pb3 := bucket.ToProtobuf()
	require.True(t, pb3 != pb, "Expected a new protobuf to be generated after the bucket was modified")
	pb4 := bucket.ToProtobuf()
	require.True(t, pb3 == pb4, "Expected the cached protobuf to be returned")
	bucket.ReportFailure(1, validatorIDs[0], 0)
	pb5 := bucket.ToProtobuf()
	require.True(t, pb5 != pb4, "Expected a new protobuf to be generated after the bucket was modified")
	pb6 := bucket.ToProtobuf()
	require.True(t, pb5 == pb6, "Expected the cached protobuf to be returned")
}

func TestReporting(t *testing.T) {
	rand := random.NewTestRandom()

	expectedSuccesses := make(map[core.QuorumID]map[core.OperatorID]uint64)
	expectedFailures := make(map[core.QuorumID]map[core.OperatorID]uint64)
	expectedSuccessBytes := make(map[core.QuorumID]map[core.OperatorID]uint64)
	expectedFailureBytes := make(map[core.QuorumID]map[core.OperatorID]uint64)
	expectedLatency := make(map[core.QuorumID]map[core.OperatorID]uint64)

	quorumCount := core.QuorumID(5)
	for quorum := core.QuorumID(0); quorum < quorumCount; quorum++ {
		expectedSuccesses[quorum] = make(map[core.OperatorID]uint64)
		expectedFailures[quorum] = make(map[core.OperatorID]uint64)
		expectedSuccessBytes[quorum] = make(map[core.OperatorID]uint64)
		expectedFailureBytes[quorum] = make(map[core.OperatorID]uint64)
		expectedLatency[quorum] = make(map[core.OperatorID]uint64)
	}

	validatorCount := rand.IntRange(1, 10)

	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))

		for quorum := core.QuorumID(0); quorum < quorumCount; quorum++ {
			expectedSuccesses[quorum][validatorIDs[i]] = 0
			expectedFailures[quorum][validatorIDs[i]] = 0
			expectedSuccessBytes[quorum][validatorIDs[i]] = 0
			expectedFailureBytes[quorum][validatorIDs[i]] = 0
			expectedLatency[quorum][validatorIDs[i]] = 0
		}
	}

	span := rand.DurationRange(time.Second, time.Hour)
	bucket, err := NewSigningRateBucket(rand.Time(), span)
	require.NoError(t, err)

	// Simulate a bunch of random reports.
	for i := 0; i < 10_000; i++ {
		batchSize := rand.Uint64Range(1, 1000)
		validatorIndex := rand.Intn(validatorCount)
		validatorID := validatorIDs[validatorIndex]
		quorum := core.QuorumID(rand.Intn((int)(quorumCount)))

		if rand.Bool() {
			latency := rand.DurationRange(time.Second, time.Hour)
			bucket.ReportSuccess(quorum, validatorID, batchSize, latency)

			expectedSuccesses[quorum][validatorID] += 1
			expectedSuccessBytes[quorum][validatorID] += batchSize
			expectedLatency[quorum][validatorID] += uint64(latency.Nanoseconds())
		} else {
			bucket.ReportFailure(quorum, validatorID, batchSize)

			expectedFailures[quorum][validatorID] += 1
			expectedFailureBytes[quorum][validatorID] += batchSize
		}
	}

	// Verify the results.
	for quorum := core.QuorumID(0); quorum < quorumCount; quorum++ {
		for _, validatorID := range validatorIDs {
			signingRate := bucket.getValidator(quorum, validatorID)
			require.Equal(t, expectedSuccesses[quorum][validatorID], signingRate.GetSignedBatches())
			require.Equal(t, expectedSuccessBytes[quorum][validatorID], signingRate.GetSignedBytes())
			require.Equal(t, expectedFailures[quorum][validatorID], signingRate.GetUnsignedBatches())
			require.Equal(t, expectedFailureBytes[quorum][validatorID], signingRate.GetUnsignedBytes())
			require.Equal(t, expectedLatency[quorum][validatorID], signingRate.GetSigningLatency())
		}
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
