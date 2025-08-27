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

// A Bucket for tracking signing rates. A bucket holds information about signing rates for a specific time interval.
type Bucket struct {
	logger logging.Logger

	// The timestamp when data could have first been added to this bucket.
	startTimestamp time.Time

	// The timestamp when the last data could have been added to this bucket.
	endTimestamp time.Time

	// Signing rate info for each validator in this bucket. If a validator is not present in this map, it can be
	// assumed that the validator has been Down for the entire duration of the bucket.
	validatorInfo map[core.OperatorID]*validator.ValidatorSigningRate

	// For each validator, whether we think that validator is Up or Down.
	// If a validator is not present in this map, it can be assumed that it has been "down" for the entire duration
	// of the bucket.
	statusMap map[core.OperatorID]ValidatorStatus

	// For each validator, the timestamp when its entry in the statusMap last changed. If a validator is not present
	// in this map, it can be assumed that the validator has been down for the entire duration of the bucket.
	lastStatusUpdate map[core.OperatorID]time.Time
}

// Create a new empty Bucket. If provided with validator IDs, then those validators will be initialized as Up.
// Any validator not provided to the constructor will be assumed to be Down until it is reported as Up.
func NewBucket(
	logger logging.Logger,
	startTime time.Time,
	onlineValidators ...core.OperatorID,
) *Bucket {

	validatorInfo := make(map[core.OperatorID]*validator.ValidatorSigningRate)
	statusMap := make(map[core.OperatorID]ValidatorStatus)
	lastStatusUpdate := make(map[core.OperatorID]time.Time)
	bucket := &Bucket{
		logger:           logger,
		startTimestamp:   startTime,
		endTimestamp:     startTime,
		validatorInfo:    validatorInfo,
		statusMap:        statusMap,
		lastStatusUpdate: lastStatusUpdate,
	}

	for _, id := range onlineValidators {
		bucket.newValidator(id, startTime, Up)
	}

	return bucket
}

// Convert this Bucket to its protobuf representation.
//
// If now is nil, then the EndTimestamp will be set to the last time that data was added to this bucket.
// If now is non-nil, then the EndTimestamp will be set to now. In general, a non-nil value for now should
// be provided when getting information about a bucket that is currently being written to, and nil should
// be provided when getting information about bucket in the past that is no longer being written to.
func (b *Bucket) ToProto(now *time.Time) *validator.SigningRateBucket {
	start := uint64(b.startTimestamp.Unix())
	var end uint64
	if now != nil {
		end = uint64(now.Unix())
	} else {
		end = uint64(b.endTimestamp.Unix())
	}

	validatorSigningRates := make([]*validator.ValidatorSigningRate, 0, len(b.validatorInfo))
	for _, info := range b.validatorInfo {
		validatorSigningRates = append(validatorSigningRates, info)
	}

	return &validator.SigningRateBucket{
		StartTimestamp:        start,
		EndTimestamp:          end,
		ValidatorSigningRates: validatorSigningRates,
	}
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

	info.SignedBatches += 1
	info.SignedBytes += batchSize
	info.SigningLatency += uint64(signingLatency.Nanoseconds())

	b.endTimestamp = now
}

// Report that a validator has failed to sign a batch of the given size.
//
// If the validator was previously Up, it will be marked as Down.
func (b *Bucket) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {

	info := b.getValidator(id, now)

	b.handleStatusChange(now, id, Down)

	info.UnsignedBatches += 1
	info.UnsignedBytes += batchSize

	b.endTimestamp = now
}

// Get the signing rate info for a validator, creating a new entry if necessary.
func (b *Bucket) getValidator(id core.OperatorID, now time.Time) *validator.ValidatorSigningRate {
	info, exists := b.validatorInfo[id]
	if !exists {
		info = b.newValidator(id, now, Down)
	}
	return info
}

// Add a new validator to the set of validators tracked by this bucket.
func (b *Bucket) newValidator(
	id core.OperatorID,
	now time.Time,
	status ValidatorStatus) *validator.ValidatorSigningRate {

	signingRate := &validator.ValidatorSigningRate{
		Id: id[:],
	}
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
	info.Uptime += uint64(nanosecondsUp)
}
