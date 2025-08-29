package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// TODO don't use time.Now()

// Tracks signing rates for validators and serves queries about signing rates.
type SigningRateTracker interface {

	// Get the signing rate for a validator over the specified time range. Start time is rounded forwards/backwards
	// to the nearest bucket boundaries.
	//
	// Returned data threadsafe to read, but should not be modified.
	GetValidatorSigningRate(
		operatorID []byte,
		startTime time.Time,
		endTime time.Time,
	) (*validator.ValidatorSigningRate, error)

	// Extract all signing rate data currently tracked by the store starting at a given timestamp.
	// Data is returned in chronological order.
	//
	// Returned data threadsafe to read, but should not be modified.
	GetSigningRateDump(startTime time.Time, now time.Time) []*validator.SigningRateBucket

	// Returns a list of buckets that have not yet been flushed to persistent storage.
	// Buckets are in chronological order.
	//
	// Returned data threadsafe to read, but should not be modified.
	GetUnflushedBuckets() []*validator.SigningRateBucket

	// Report that a validator has successfully signed a batch of the given size.
	ReportSuccess(
		now time.Time,
		id core.OperatorID,
		batchSize uint64,
		signingLatency time.Duration,
	)

	// Report that a validator has failed to sign a batch of the given size.
	ReportFailure(
		now time.Time,
		id core.OperatorID,
		batchSize uint64,
	)

	// Close the store and free any associated resources.
	Close()
}

var _ SigningRateTracker = (*signingRateTracker)(nil)

// A standard implementation of the SigningRateTracker interface. Is not thread safe on its own.
type signingRateTracker struct {
	logger logging.Logger

	// Signing data storage, split up into buckets for each time interval. Buckets are stored in chronological order.
	buckets *common.RandomAccessDeque[*Bucket]

	// Buckets that have not yet been flushed to storage. Keyed by the bucket's start time.
	unflushedBuckets map[time.Time]*Bucket

	// The length of time to keep loaded in memory.
	timespan time.Duration

	// The duration of each bucket. Buckets loaded from storage may have different spans, but new buckets will
	// always have this span.
	bucketSpan time.Duration
}

// Create a new SigningRateTracker.
//
//   - signingRateDatabase: The database to use for storing historical signing rate information.
//   - timespan: The amount of time to keep in memory. Queries are only supported for this timespan.
//   - bucketSpan: The duration of each bucket.
//   - flushPeriod: How often to flush in-memory data to the database. If the process is shut down/crashes, any data
//     not yet flushed to the database may be lost.
func NewSigningRateTracker(
	logger logging.Logger,
	timespan time.Duration,
	bucketSpan time.Duration,
	buckets []*Bucket,
) (SigningRateTracker, error) {

	store := &signingRateTracker{
		logger:     logger,
		buckets:    common.NewRandomAccessDeque[*Bucket](0),
		timespan:   timespan,
		bucketSpan: bucketSpan,
	}

	// Load old buckets.
	for _, bucket := range buckets {
		store.buckets.PushBack(bucket)
	}

	return store, nil
}

func (s *signingRateTracker) Close() {
	// This implementation has no resources that the garbage collector won't clean up, so nothing to do here.
}

// Report that a validator has successfully signed a batch of the given size.
func (s *signingRateTracker) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {

	bucket := s.getMutableBucket()
	bucket.ReportSuccess(now, id, batchSize, signingLatency)
	s.markUnflushed(bucket)
}

// Report that a validator has failed to sign a batch of the given size.
func (s *signingRateTracker) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	bucket := s.getMutableBucket()
	bucket.ReportFailure(now, id, batchSize)
	s.markUnflushed(bucket)
}

func (s *signingRateTracker) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {
	return nil, nil // TODO
}

func (s *signingRateTracker) GetSigningRateDump(startTime time.Time, now time.Time) []*validator.SigningRateBucket {
	buckets := make([]*validator.SigningRateBucket, 0, s.buckets.Size())

	for i := uint64(0); i < s.buckets.Size(); i++ {
		bucket, err := s.buckets.Get(i)
		enforce.NilError(err, "should be impossible with valid index")
		proto := bucket.ToProtobuf(now)
		buckets = append(buckets, proto)
	}

	// No need to sort, s.buckets is always in chronological order.

	return buckets
}

func (s *signingRateTracker) GetUnflushedBuckets() []*validator.SigningRateBucket {
	buckets := make([]*validator.SigningRateBucket, 0, len(s.unflushedBuckets))

	for _, bucket := range s.unflushedBuckets {
		proto := bucket.ToProtobuf(time.Now())
		buckets = append(buckets, proto)
	}

	sortValidatorSigningRateBuckets(buckets)

	return buckets
}

// Get the bucket that is currently being written to. This is always the latest bucket.
func (s *signingRateTracker) getMutableBucket() *Bucket {

	now := time.Now()

	if s.buckets.Size() == 0 {
		// Create the first bucket.
		newBucket := NewBucket(s.logger, now, s.bucketSpan)
		s.buckets.PushBack(newBucket)
	}

	bucket, err := s.buckets.PeekBack()
	enforce.NilError(err, "should be impossible with a non-empty deque")

	if now.After(bucket.EndTimestamp()) {
		// The current bucket's time span has elapsed, create a new bucket.

		bucket = NewBucket(s.logger, now, s.bucketSpan, bucket.GetOnlineValidators()...)
		s.buckets.PushBack(bucket)

		// Now is a good time to do garbage collection. As long as bucket size remains fixed, we should be removing
		// one bucket for each new bucket we add once we reach steady state.
		s.garbageCollectBuckets()
	}

	return bucket
}

// Remove old buckets that are outside the configured timespan.
func (s *signingRateTracker) garbageCollectBuckets() {
	cutoff := time.Now().Add(-s.timespan)

	for s.buckets.Size() > 0 {
		bucket, err := s.buckets.PeekFront()
		enforce.NilError(err, "should be impossible with a non-empty deque")

		if bucket.EndTimestamp().After(cutoff) {
			// This bucket is new enough, so all later buckets will also be new enough.
			break
		}

		// This bucket is too old, remove it.
		_, err = s.buckets.PopFront()
		enforce.NilError(err, "should be impossible with a non-empty deque")
	}
}

// Mark a bucket as needing to be flushed to storage.
func (s *signingRateTracker) markUnflushed(bucket *Bucket) {
	s.unflushedBuckets[bucket.startTimestamp] = bucket
}
