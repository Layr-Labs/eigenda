package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Tracks signing rates for validators and serves queries about signing rates.
type SigningRateTracker interface {
	// Close the store and free any associated resources.
	Close()

	// Get the signing rate for a validator over the specified time range. Start time is rounded forwards/backwards
	// to the nearest bucket boundaries. Returned data is immutable.
	GetValidatorSigningRate(
		operatorID []byte,
		startTime time.Time,
		endTime time.Time,
	) (*validator.ValidatorSigningRate, error)

	// Extract all signing rate data currently tracked by the store. Data is returned in chronological order.
	// Returned data is immutable.
	GetSigningRateDump() ([]*validator.SigningRateBucket, error)

	// Report that a validator has successfully signed a batch of the given size.
	ReportSuccess(
		now time.Time,
		id core.OperatorID,
		batchSize uint64,
		signingLatency time.Duration,
	)

	// Report that a validator has failed to sign a batch of the given size.
	ReportFailure(now time.Time, id core.OperatorID, batchSize uint64)
}

// A standard implementation of the SigningRateTracker interface. Is not thread safe on its own.
type signingRateTracker struct {
	logger logging.Logger
}

// Create a new SigningRateTracker.
//
//   - signingRateDatabase: The database to use for storing historical signing rate information.
//   - timespan: The amount of time to keep in memory. Queries are only supported for this timespan.
//   - bucketSpan: The duration of each bucket.
//   - flushPeriod: How often to flush in-memory data to the database. If the process is shut down/crashes, any data
//     not yet flushed to the database may be lost.
func NewSigningRateTracker(
	signingRateDatabase SigningRateStorage,
	timespan time.Duration,
	bucketSpan time.Duration,
	flushPeriod time.Duration,
) (SigningRateTracker, error) {

	store := &signingRateTracker{}

	return store, nil
}

func (s *signingRateTracker) Close() {
	// TODO
}

// Get the signing rate for a validator over the specified time range. Start time is rounded forwards/backwards
// to the nearest bucket boundaries. Returned data is immutable.
func (s *signingRateTracker) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {
	return nil, nil // TODO
}

// Extract all signing rate data currently tracked by the store. Data is returned in chronological order.
// Returned data is immutable.
func (s *signingRateTracker) GetSigningRateDump() ([]*validator.SigningRateBucket, error) {
	return nil, nil // TODO
}

// Report that a validator has successfully signed a batch of the given size.
func (s *signingRateTracker) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	s.getMutableBucket().ReportSuccess(now, id, batchSize, signingLatency)
}

// Report that a validator has failed to sign a batch of the given size.
func (s *signingRateTracker) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	s.getMutableBucket().ReportFailure(now, id, batchSize)
}

// Get the bucket that is currently being written to. This is always the latest bucket.
func (s *signingRateTracker) getMutableBucket() *Bucket {
	return nil // TODO
}
