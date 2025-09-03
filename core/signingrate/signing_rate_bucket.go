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

	// Signing rate info for each validator in this bucket. If a validator is not present in this map, it can be
	// assumed that the validator has been Down for the entire duration of the bucket.
	validatorInfo map[core.OperatorID]*validator.ValidatorSigningRate

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

	validatorInfo := make(map[core.OperatorID]*validator.ValidatorSigningRate)
	bucket := &SigningRateBucket{
		startTimestamp: startTimestamp,
		endTimestamp:   endTimestamp,
		validatorInfo:  validatorInfo,
	}

	return bucket, nil
}

// Parse a SigningRateBucket from its protobuf representation.
func NewBucketFromProto(pb *validator.SigningRateBucket) *SigningRateBucket {

	startTime := time.Unix(int64(pb.GetStartTimestamp()), 0)
	endTime := time.Unix(int64(pb.GetEndTimestamp()), 0)

	validatorInfo := make(map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, info := range pb.GetValidatorSigningRates() {
		var id core.OperatorID
		copy(id[:], info.GetId())
		validatorInfo[id] = info
	}

	return &SigningRateBucket{
		startTimestamp: startTime,
		endTimestamp:   endTime,
		validatorInfo:  validatorInfo,
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

	validatorSigningRates := make([]*validator.ValidatorSigningRate, 0, len(b.validatorInfo))
	for _, info := range b.validatorInfo {
		validatorSigningRates = append(validatorSigningRates, cloneValidatorSigningRate(info))
	}

	// Sort for deterministic output. Not strictly necessary, but sometimes nice to have.
	sortValidatorSigningRate(validatorSigningRates)

	b.cachedProtobuf = &validator.SigningRateBucket{
		StartTimestamp:        start,
		EndTimestamp:          end,
		ValidatorSigningRates: validatorSigningRates,
	}
	return b.cachedProtobuf
}

// Report that a validator has successfully signed a batch of the given size.
//
// If the validator was previously Down, it will be marked as Up.
func (b *SigningRateBucket) ReportSuccess(
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	info := b.getValidator(id)

	info.SignedBatches += 1
	info.SignedBytes += batchSize
	info.SigningLatency += uint64(signingLatency.Nanoseconds())

	b.cachedProtobuf = nil
}

// Report that a validator has failed to sign a batch of the given size.
//
// If the validator was previously Up, it will be marked as Down.
func (b *SigningRateBucket) ReportFailure(id core.OperatorID, batchSize uint64) {
	info := b.getValidator(id)

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

// Get the signing rate info for a validator, creating a new entry if necessary. Is not a deep copy.
func (b *SigningRateBucket) getValidator(id core.OperatorID) *validator.ValidatorSigningRate {
	info, exists := b.validatorInfo[id]
	if !exists {
		info = b.newValidator(id)
	}
	return info
}

// Get the signing rate info for a validator if it is registered, or nil if it is not. Is not a deep copy.
func (b *SigningRateBucket) getValidatorIfExists(
	id core.OperatorID,
) (signingRate *validator.ValidatorSigningRate, exists bool) {

	signingRate, exists = b.validatorInfo[id]
	return signingRate, exists
}

// Add a new validator to the set of validators tracked by this bucket.
func (b *SigningRateBucket) newValidator(id core.OperatorID) *validator.ValidatorSigningRate {
	signingRate := &validator.ValidatorSigningRate{
		Id: id[:],
	}
	b.validatorInfo[id] = signingRate
	return signingRate
}
