package signingrate

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

// A SigningRateBucket for tracking signing rates. A bucket holds information about signing rates for a specific time
// interval. Roughly correlates to the validator.SigningRateBucket protobuf message, but also includes extra data
// structures to help with tracking state while the bucket is being written to.
type SigningRateBucket struct {
	// The timestamp when data could have first been added to this bucket.
	startTimestamp time.Time

	// The timestamp when the last data could have been added to this bucket.
	endTimestamp time.Time

	// The signing rate information for the time period covered by this bucket.
	signingRateInfo map[core.QuorumID]map[core.OperatorID]*validator.ValidatorSigningRate

	// A cached protobuf representation of this bucket. Set to nil whenever the bucket is modified.
	cachedProtobuf *validator.SigningRateBucket
}

// Create a new empty SigningRateBucket.
func NewSigningRateBucket(startTime time.Time, span time.Duration) (*SigningRateBucket, error) {
	startTimestamp, err := bucketStartTimestamp(span, startTime)
	if err != nil {
		return nil, fmt.Errorf("error creating signing rate bucket: %w", err)
	}

	endTimestamp, err := bucketEndTimestamp(span, startTime)
	if err != nil {
		return nil, fmt.Errorf("error creating signing rate bucket: %w", err)
	}

	validatorInfo := make(map[core.QuorumID]map[core.OperatorID]*validator.ValidatorSigningRate)
	bucket := &SigningRateBucket{
		startTimestamp:  startTimestamp,
		endTimestamp:    endTimestamp,
		signingRateInfo: validatorInfo,
	}

	return bucket, nil
}

// Parse a SigningRateBucket from its protobuf representation.
func NewBucketFromProto(pb *validator.SigningRateBucket) *SigningRateBucket {

	startTime := time.Unix(int64(pb.GetStartTimestamp()), 0)
	endTime := time.Unix(int64(pb.GetEndTimestamp()), 0)

	signingRateInfo := make(map[core.QuorumID]map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, quorumInfo := range pb.GetQuorumSigningRates() {
		quorumID := core.QuorumID(quorumInfo.GetQuorumId())
		signingRateInfo[quorumID] = make(map[core.OperatorID]*validator.ValidatorSigningRate)

		for _, validatorInfo := range quorumInfo.GetValidatorSigningRates() {
			validatorID := core.OperatorID{}
			copy(validatorID[:], validatorInfo.GetId())

			signingRateInfo[quorumID][validatorID] = cloneValidatorSigningRate(validatorInfo)
		}
	}

	return &SigningRateBucket{
		startTimestamp:  startTime,
		endTimestamp:    endTime,
		signingRateInfo: signingRateInfo,
	}
}

// Convert this SigningRateBucket to its protobuf representation.
//
// If now is nil, then the EndTimestamp will be set to the last time that data was added to this bucket.
// If now is non-nil, then the EndTimestamp will be set to now. In general, a non-nil value for now should
// be provided when getting information about a bucket that is currently being written to, and nil should
// be provided when getting information about bucket in the past that is no longer being written to.
//
// The resulting protobuf is a deep copy, and is therefore threadsafe to use concurrently with this SigningRateBucket.
func (b *SigningRateBucket) ToProtobuf() *validator.SigningRateBucket {
	if b.cachedProtobuf != nil {
		return b.cachedProtobuf
	}

	start := uint64(b.startTimestamp.Unix())
	end := uint64(b.endTimestamp.Unix())

	quorumSigningRates := make([]*validator.QuorumSigningRate, len(b.signingRateInfo))

	for quorumID, quorumInfo := range b.signingRateInfo {
		validatorSigningRates := make([]*validator.ValidatorSigningRate, 0, len(b.signingRateInfo))

		for _, validatorInfo := range quorumInfo {
			validatorSigningRates = append(validatorSigningRates, cloneValidatorSigningRate(validatorInfo))
		}
		sortValidatorSigningRates(validatorSigningRates)

		quorumSigningRates = append(quorumSigningRates,
			&validator.QuorumSigningRate{
				QuorumId:              uint64(quorumID),
				ValidatorSigningRates: validatorSigningRates,
			})
	}
	sortQuorumSigningRates(quorumSigningRates)

	b.cachedProtobuf = &validator.SigningRateBucket{
		StartTimestamp:     start,
		EndTimestamp:       end,
		QuorumSigningRates: quorumSigningRates,
	}
	return b.cachedProtobuf
}

// Report that a validator has successfully signed a batch of the given size.
//
// If the validator was previously Down, it will be marked as Up.
func (b *SigningRateBucket) ReportSuccess(
	quorum core.QuorumID,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	info := b.getValidator(quorum, id)

	info.SignedBatches += 1
	info.SignedBytes += batchSize
	info.SigningLatency += uint64(signingLatency.Nanoseconds())

	b.cachedProtobuf = nil
}

// Report that a validator has failed to sign a batch of the given size.
//
// If the validator was previously Up, it will be marked as Down.
func (b *SigningRateBucket) ReportFailure(quorum core.QuorumID, id core.OperatorID, batchSize uint64) {
	info := b.getValidator(quorum, id)

	info.UnsignedBatches += 1
	info.UnsignedBytes += batchSize

	b.cachedProtobuf = nil
}

// Get the start timestamp of this bucket (inclusive).
func (b *SigningRateBucket) StartTimestamp() time.Time {
	return b.startTimestamp
}

// Get the end timestamp of this bucket (exclusive).
func (b *SigningRateBucket) EndTimestamp() time.Time {
	return b.endTimestamp
}

// Returns true if the given time is contained within this bucket. Start time is inclusive, end time is exclusive.
func (b *SigningRateBucket) Contains(t time.Time) bool {
	return !t.Before(b.startTimestamp) && t.Before(b.endTimestamp)
}

// Get the signing rate info for a validator in a particular quorum, creating a new entry if necessary. Is not a deep copy.
func (b *SigningRateBucket) getValidator(quorum core.QuorumID, id core.OperatorID) *validator.ValidatorSigningRate {
	quorumSigningRate, exists := b.signingRateInfo[quorum]
	if !exists {
		quorumSigningRate = make(map[core.OperatorID]*validator.ValidatorSigningRate)
		b.signingRateInfo[quorum] = quorumSigningRate
	}

	validatorSigningRate, exists := quorumSigningRate[id]
	if !exists {
		validatorSigningRate = &validator.ValidatorSigningRate{
			Id: id[:],
		}
		quorumSigningRate[id] = validatorSigningRate
	}
	return validatorSigningRate
}

// Get the signing rate info for a validator in a quorum if it is registered, or nil if it is not. Is not a deep copy.
func (b *SigningRateBucket) getValidatorIfExists(
	quorum core.QuorumID,
	id core.OperatorID,
) (signingRate *validator.ValidatorSigningRate, exists bool) {

	quorumSigningRate, exists := b.signingRateInfo[quorum]
	if !exists {
		return nil, false
	}

	signingRate, exists = quorumSigningRate[id]
	return signingRate, exists
}
