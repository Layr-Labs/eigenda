package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

var _ SigningRateStore = (*ThreadsafeSigningRateStore)(nil)

// A thread-safe wrapper around a SigningRateStore.
type ThreadsafeSigningRateStore struct {
	base SigningRateStore
}

// Create a new ThreadsafeSigningRateStore that wraps the provided base SigningRateStore.
func NewThreadsafeSigningRateStore(base SigningRateStore) SigningRateStore {

	store := &ThreadsafeSigningRateStore{
		base: base,
	}

	return store
}

func (t ThreadsafeSigningRateStore) Close() {
	//TODO implement me
	panic("implement me")
}

func (t ThreadsafeSigningRateStore) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {

	//TODO implement me
	panic("implement me")
}

func (t ThreadsafeSigningRateStore) GetSigningRateDump() ([]*validator.SigningRateBucket, error) {
	//TODO implement me
	panic("implement me")
}

func (t ThreadsafeSigningRateStore) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	//TODO implement me
	panic("implement me")
}

func (t ThreadsafeSigningRateStore) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	//TODO implement me
	panic("implement me")
}
