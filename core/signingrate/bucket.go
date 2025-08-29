package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// ValidatorStatus represents whether a validator is considered to be "up" or "down". This is from the perspective
// of the signing rate tracker, and may not reflect the actual status of the validator.
type ValidatorStatus bool

const (
	// The status of a validator that we consider to be "up".
	Up ValidatorStatus = true
	// The status of a validator that we consider to be "down".
	Down ValidatorStatus = false
)

// TODO unit test
// TODO: periodic logging of validators with bad signing rates

// A Bucket for tracking signing rates. A bucket holds information about signing rates for a specific time interval.
// Roughly correlates to the validator.SigningRateBucket protobuf message, but also includes extra data structures
// to help with tracking state while the bucket is being written to.
type Bucket struct {
	logger logging.Logger

	// The timestamp when data could have first been added to this bucket.
	startTimestamp time.Time

	// The timestamp when the last data could have been added to this bucket.
	endTimestamp time.Time

	// Signing rate info for each validator in this bucket. If a validator is not present in this map, it can be
	// assumed that the validator has been Down for the entire duration of the bucket.
	validatorInfo map[core.OperatorID]*SigningRate

	// For each validator, whether we think that validator is Up or Down.
	// If a validator is not present in this map, it can be assumed that it has been "down" for the entire duration
	// of the bucket. This field is ignored when a bucket is immutable, since it's only
	// a tool for tracking state while the bucket is being written to.
	statusMap map[core.OperatorID]ValidatorStatus

	// For each validator, the timestamp when its entry in the statusMap last changed. If a validator is not present
	// in this map, it can be assumed that the validator has been down for the entire duration of the bucket.
	lastStatusUpdate map[core.OperatorID]time.Time

	// A cached protobuf representation of this bucket. Set to nil whenever the bucket is modified.
	cachedProtobuf *validator.SigningRateBucket
}

// Create a new empty Bucket. If provided with validator IDs, then those validators will be initialized as Up.
// Any validator not provided to the constructor will be assumed to be Down until it is reported as Up.
func NewBucket(
	logger logging.Logger,
	startTime time.Time,
	span time.Duration,
	onlineValidators ...core.OperatorID,
) *Bucket {

	validatorInfo := make(map[core.OperatorID]*SigningRate)
	statusMap := make(map[core.OperatorID]ValidatorStatus)
	lastStatusUpdate := make(map[core.OperatorID]time.Time)
	bucket := &Bucket{
		logger:           logger,
		startTimestamp:   startTime,
		endTimestamp:     startTime.Add(span),
		validatorInfo:    validatorInfo,
		statusMap:        statusMap,
		lastStatusUpdate: lastStatusUpdate,
	}

	for _, id := range onlineValidators {
		bucket.newValidator(id, startTime, Up)
	}

	return bucket
}

// Parse a Bucket from its protobuf representation.
func NewBucketFromProto(
	logger logging.Logger,
	pb *validator.SigningRateBucket,
) *Bucket {

	startTime := time.Unix(int64(pb.GetStartTimestamp()), 0)
	endTime := time.Unix(int64(pb.GetEndTimestamp()), 0)

	validatorInfo := make(map[core.OperatorID]*SigningRate)
	statusMap := make(map[core.OperatorID]ValidatorStatus)
	lastStatusUpdate := make(map[core.OperatorID]time.Time)

	for _, info := range pb.GetValidatorSigningRates() {
		var id core.OperatorID
		copy(id[:], info.GetId())
		validatorInfo[id] = NewSigningRateFromProtobuf(info)
		lastStatusUpdate[id] = startTime
	}

	return &Bucket{
		logger:           logger,
		startTimestamp:   startTime,
		endTimestamp:     endTime,
		validatorInfo:    validatorInfo,
		statusMap:        statusMap,
		lastStatusUpdate: lastStatusUpdate,
	}
}

// Convert this Bucket to its protobuf representation.
//
// If now is nil, then the EndTimestamp will be set to the last time that data was added to this bucket.
// If now is non-nil, then the EndTimestamp will be set to now. In general, a non-nil value for now should
// be provided when getting information about a bucket that is currently being written to, and nil should
// be provided when getting information about bucket in the past that is no longer being written to.
func (b *Bucket) ToProtobuf(now time.Time) *validator.SigningRateBucket {
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
		validatorSigningRates = append(validatorSigningRates, info.ToProtobuf())
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
func (b *Bucket) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {

	info := b.getValidator(id, now)

	b.handleStatusChange(now, id, Up)

	info.IncrementSignedBatches()
	info.AddSignedBytes(batchSize)
	info.AddSigningLatency(uint64(signingLatency.Nanoseconds()))

	b.endTimestamp = now
	b.cachedProtobuf = nil
}

// Report that a validator has failed to sign a batch of the given size.
//
// If the validator was previously Up, it will be marked as Down.
func (b *Bucket) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {

	info := b.getValidator(id, now)

	b.handleStatusChange(now, id, Down)

	info.IncrementSignedBatches()
	info.AddUnsignedBytes(batchSize)

	b.endTimestamp = now
	b.cachedProtobuf = nil
}

// Get the start timestamp of this bucket.
func (b *Bucket) StartTimestamp() time.Time {
	return b.startTimestamp
}

// Get the end timestamp of this bucket.
func (b *Bucket) EndTimestamp() time.Time {
	return b.endTimestamp
}

// Fetch a list of validator IDs that are currently considered to be Up.
func (b *Bucket) GetOnlineValidators() []core.OperatorID {
	ids := make([]core.OperatorID, 0, len(b.statusMap))
	for id, status := range b.statusMap {
		if status == Up {
			ids = append(ids, id)
		}
	}
	return ids
}

// Get the signing rate info for a validator, creating a new entry if necessary.
func (b *Bucket) getValidator(id core.OperatorID, now time.Time) *SigningRate {
	info, exists := b.validatorInfo[id]
	if !exists {
		info = b.newValidator(id, now, Down)
	}
	return info
}

// Get the signing rate info for a validator if it is registered, or nil if it is not.
func (b *Bucket) getValidatorIfExists(id core.OperatorID) (signingRate *SigningRate, exists bool) {
	signingRate, exists = b.validatorInfo[id]
	return signingRate, exists
}

// Add a new validator to the set of validators tracked by this bucket.
func (b *Bucket) newValidator(
	id core.OperatorID,
	now time.Time,
	status ValidatorStatus) *SigningRate {

	signingRate := NewSigningRate(id)
	b.validatorInfo[id] = signingRate
	b.statusMap[id] = status
	b.lastStatusUpdate[id] = now

	return signingRate
}

// Handle a change in the status of a validator. This is a no-op if the status has not changed.
func (b *Bucket) handleStatusChange(now time.Time, id core.OperatorID, newStatus ValidatorStatus) {
	currentStatus, ok := b.statusMap[id]
	enforce.True(ok, "validator %v not found in status map", id)

	if currentStatus == newStatus {
		// no change in status
		return
	}

	defer func() {
		b.lastStatusUpdate[id] = now
	}()

	if newStatus == Up {
		// Transition from Down to Up
		b.statusMap[id] = Up

		// No need to record the amount of time the validator spent Down, since we only track Uptime.
		// It's simple enough to infer Downtime from Uptime and the total duration of the bucket.
		return
	}

	// Transition from Up to Down
	b.statusMap[id] = Down

	// Record the amount of time the validator spent Up.
	lastUpdateTime, ok := b.lastStatusUpdate[id]
	enforce.True(ok, "validator %v not found in last status map", id)
	if now.Before(lastUpdateTime) {
		b.logger.Errorf("clock went backwards: now=%v, lastUpdateTime=%v", now, lastUpdateTime)
		return
	}
	nanosecondsUp := now.Sub(lastUpdateTime).Nanoseconds()
	info, ok := b.validatorInfo[id]
	enforce.True(ok, "validator %v not found in validator info map", id)
	info.AddUptime(uint64(nanosecondsUp))
}
