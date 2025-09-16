package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

var _ SigningRateTracker = (*noOpSigningRateTracker)(nil)

// A no-op implementation of the SigningRateTracker interface, for unit tests.
type noOpSigningRateTracker struct {
}

// Create a new no-op SigningRateTracker.
func NewNoOpSigningRateTracker() SigningRateTracker {
	return &noOpSigningRateTracker{}
}

func (n *noOpSigningRateTracker) GetValidatorSigningRate(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {
	return &validator.ValidatorSigningRate{
		ValidatorId:     validatorID[:],
		SignedBatches:   0,
		UnsignedBatches: 0,
		SignedBytes:     0,
		UnsignedBytes:   0,
		SigningLatency:  0,
	}, nil
}

func (n *noOpSigningRateTracker) GetSigningRateDump(startTime time.Time) ([]*validator.SigningRateBucket, error) {
	return make([]*validator.SigningRateBucket, 0), nil
}

func (n *noOpSigningRateTracker) GetUnflushedBuckets() ([]*validator.SigningRateBucket, error) {
	return make([]*validator.SigningRateBucket, 0), nil
}

func (n *noOpSigningRateTracker) ReportSuccess(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	// no-op
}

func (n *noOpSigningRateTracker) ReportFailure(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	batchSize uint64,
) {

	// no-op
}

func (n *noOpSigningRateTracker) UpdateLastBucket(bucket *validator.SigningRateBucket) {
	// no-op
}

func (n *noOpSigningRateTracker) GetLastBucketStartTime() (time.Time, error) {
	return time.Time{}, nil
}

func (n *noOpSigningRateTracker) Flush() error {
	return nil
}

func (n *noOpSigningRateTracker) Close() {
	// no-op
}
