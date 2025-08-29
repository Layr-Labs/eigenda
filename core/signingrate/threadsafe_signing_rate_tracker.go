package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

var _ SigningRateTracker = (*threadsafeSigningRateTracker)(nil)

// A thread-safe wrapper around a SigningRateTracker.
type threadsafeSigningRateTracker struct {
	base SigningRateTracker
}

func (t threadsafeSigningRateTracker) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {
	//TODO implement me
	panic("implement me")
}

func (t threadsafeSigningRateTracker) GetSigningRateDump(
	startTime time.Time,
	now time.Time,
) []*validator.SigningRateBucket {
	//TODO implement me
	panic("implement me")
}

func (t threadsafeSigningRateTracker) GetUnflushedBuckets() []*validator.SigningRateBucket {
	//TODO implement me
	panic("implement me")
}

func (t threadsafeSigningRateTracker) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	//TODO implement me
	panic("implement me")
}

func (t threadsafeSigningRateTracker) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	//TODO implement me
	panic("implement me")
}

func (t threadsafeSigningRateTracker) Close() {
	//TODO implement me
	panic("implement me")
}
