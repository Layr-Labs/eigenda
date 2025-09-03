package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// A SigningRateBucket for tracking signing rates. A bucket holds information about signing rates for a specific time interval.
// Roughly correlates to the validator.SigningRateBucket protobuf message, but also includes extra data structures
// to help with tracking state while the bucket is being written to.
type SigningRateBucket struct {
	logger logging.Logger

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

// Create a new empty SigningRateBucket. If provided with validator IDs, then those validators will be initialized as Up.
// Any validator not provided to the constructor will be assumed to be Down until it is reported as Up.
func NewBucket(
	logger logging.Logger,
	startTime time.Time,
	span time.Duration,
) *SigningRateBucket {

	validatorInfo := make(map[core.OperatorID]*validator.ValidatorSigningRate)
	bucket := &SigningRateBucket{
		logger:         logger,
		startTimestamp: startTime,
		endTimestamp:   startTime.Add(span),
		validatorInfo:  validatorInfo,
	}

	return bucket
}

// Parse a SigningRateBucket from its protobuf representation.
func NewBucketFromProto(
	logger logging.Logger,
	pb *validator.SigningRateBucket,
) *SigningRateBucket {

	startTime := time.Unix(int64(pb.GetStartTimestamp()), 0)
	endTime := time.Unix(int64(pb.GetEndTimestamp()), 0)

	validatorInfo := make(map[core.OperatorID]*validator.ValidatorSigningRate)

	for _, info := range pb.GetValidatorSigningRates() {
		var id core.OperatorID
		copy(id[:], info.GetId())
		validatorInfo[id] = info
	}

	return &SigningRateBucket{
		logger:         logger,
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
func (b *SigningRateBucket) ToProtobuf(now time.Time) *validator.SigningRateBucket {
	if b.cachedProtobuf != nil {
		return b.cachedProtobuf
	}

	start := uint64(b.startTimestamp.Unix())

	var end uint64
	if b.endTimestamp.IsZero() {
		// The end timestamp of this bucket is not yet set. Use the current time as the end timestamp.
		end = uint64(now.Unix())
	} else {
		// The end timestamp of this bucket is set. Use it.
		end = uint64(b.endTimestamp.Unix())
	}

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
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	info := b.getValidator(id)

	info.SignedBatches += 1
	info.SignedBytes += batchSize
	info.SigningLatency += uint64(signingLatency.Nanoseconds())

	b.endTimestamp = now
	b.cachedProtobuf = nil
}

// Report that a validator has failed to sign a batch of the given size.
//
// If the validator was previously Up, it will be marked as Down.
func (b *SigningRateBucket) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	info := b.getValidator(id)

	info.SignedBatches += 1
	info.UnsignedBytes += batchSize

	b.endTimestamp = now
	b.cachedProtobuf = nil
}

// Get the start timestamp of this bucket.
func (b *SigningRateBucket) StartTimestamp() time.Time {
	return b.startTimestamp
}

// Get the end timestamp of this bucket.
func (b *SigningRateBucket) EndTimestamp() time.Time {
	return b.endTimestamp
}

// Get the signing rate info for a validator, creating a new entry if necessary.
func (b *SigningRateBucket) getValidator(id core.OperatorID) *validator.ValidatorSigningRate {
	info, exists := b.validatorInfo[id]
	if !exists {
		info = b.newValidator(id)
	}
	return info
}

// Get the signing rate info for a validator if it is registered, or nil if it is not.
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
