package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

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
	GetSigningRateDump(startTime time.Time, now time.Time) ([]*validator.SigningRateBucket, error)

	// Returns a list of buckets that have not yet been flushed to persistent storage.
	// Buckets are in chronological order.
	//
	// Returned data threadsafe to read, but should not be modified.
	GetUnflushedBuckets(time time.Time) ([]*validator.SigningRateBucket, error)

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

	// Update a bucket, overwriting an existing bucket with the same start time if it is present. Should
	// only be used to update the last bucket in the store. Data is ignored if the bucket won't be the
	// last bucket.
	//
	// This operation doesn't mark a bucket as unflushed. A bucket is only marked as unflushed when it is modified,
	// not when it is provided whole-sale from an external source.
	UpdateLastBucket(now time.Time, bucket *validator.SigningRateBucket)

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

	// Metrics about signing rates.
	metrics *SigningRateMetrics
}

// Create a new SigningRateTracker.
//
//   - signingRateDatabase: The database to use for storing historical signing rate information.
//   - timespan: The amount of time to keep in memory. Queries are only supported for this timespan.
//   - bucketSpan: The duration of each bucket.
func NewSigningRateTracker(
	logger logging.Logger,
	timespan time.Duration,
	bucketSpan time.Duration,
	registry *prometheus.Registry,
) (SigningRateTracker, error) {

	store := &signingRateTracker{
		logger:     logger,
		buckets:    common.NewRandomAccessDeque[*Bucket](0),
		timespan:   timespan,
		bucketSpan: bucketSpan,
		metrics:    NewSigningRateMetrics(registry),
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

	bucket := s.getMutableBucket(now)
	bucket.ReportSuccess(now, id, batchSize, signingLatency)
	s.markUnflushed(bucket)
	s.metrics.ReportSuccess(id, batchSize, signingLatency)
}

// Report that a validator has failed to sign a batch of the given size.
func (s *signingRateTracker) ReportFailure(now time.Time, id core.OperatorID, batchSize uint64) {
	bucket := s.getMutableBucket(now)
	bucket.ReportFailure(now, id, batchSize)
	s.markUnflushed(bucket)
	s.metrics.ReportFailure(id, batchSize)
}

func (s *signingRateTracker) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {

	comparator := func(timestamp time.Time, bucket *Bucket) int {
		if bucket.startTimestamp.Before(timestamp) {
			return -1
		} else if bucket.startTimestamp.After(timestamp) {
			return 1
		}
		return 0
	}

	startIndex, exact := common.BinarySearchInOrderedDeque(s.buckets, startTime, comparator)

	if !exact && startIndex > 0 {
		// We didn't find the bucket with the exact start time, so round backwards to the previous bucket.
		startIndex--
	}

	totalSigningRate := &validator.ValidatorSigningRate{
		Id: operatorID,
	}

	iterator, err := s.buckets.IteratorFrom(startIndex)
	enforce.NilError(err, "should be impossible with a valid index")
	for _, bucket := range iterator {
		if bucket.startTimestamp.After(endTime) {
			break
		}

		signingRate, exists := bucket.getValidatorIfExists(core.OperatorID(operatorID))
		if !exists {
			// No info for validator during this bucket, skip it.
			continue
		}

		totalSigningRate.SignedBatches += signingRate.SignedBatches()
		totalSigningRate.UnsignedBatches += signingRate.UnsignedBatches()
		totalSigningRate.SignedBytes += signingRate.SignedBytes()
		totalSigningRate.UnsignedBytes += signingRate.UnsignedBytes()
		totalSigningRate.SigningLatency += signingRate.SigningLatency()
	}

	return totalSigningRate, nil
}

func (s *signingRateTracker) GetSigningRateDump(
	startTime time.Time,
	now time.Time,
) ([]*validator.SigningRateBucket, error) {

	buckets := make([]*validator.SigningRateBucket, 0, s.buckets.Size())

	// Iterate backwards. In general, dump requests will only be used to fetch recent data, so
	// we should optimize the case where we are requesting a few buckets from the end of the deque.
	for _, bucket := range s.buckets.ReverseIterator() {
		if bucket.EndTimestamp().Before(startTime) {
			// This bucket is too old, skip it and stop iterating.
			break
		}
		buckets = append(buckets, bucket.ToProtobuf(now))
	}

	// We iterated in reverse, so reverse again to get chronological ordering.
	for i, j := 0, len(buckets)-1; i < j; i, j = i+1, j-1 {
		buckets[i], buckets[j] = buckets[j], buckets[i]
	}

	return buckets, nil
}

func (s *signingRateTracker) GetUnflushedBuckets(now time.Time) ([]*validator.SigningRateBucket, error) {
	buckets := make([]*validator.SigningRateBucket, 0, len(s.unflushedBuckets))

	for _, bucket := range s.unflushedBuckets {
		proto := bucket.ToProtobuf(now)
		buckets = append(buckets, proto)
	}

	sortValidatorSigningRateBuckets(buckets)

	return buckets, nil
}

func (s *signingRateTracker) UpdateLastBucket(now time.Time, bucket *validator.SigningRateBucket) {
	convertedBucket := NewBucketFromProto(s.logger, bucket)

	if s.buckets.Size() == 0 {
		s.buckets.PushBack(convertedBucket)
		return
	}

	previousBucket, err := s.buckets.PeekBack()
	enforce.NilError(err, "should be impossible with a non-empty deque")

	if previousBucket.startTimestamp.Equal(convertedBucket.startTimestamp) {
		// We have a bucket with the same start time, replace it.
		_, err := s.buckets.SetFromBack(0, convertedBucket)
		enforce.NilError(err, "should be impossible with a valid index")
		return
	}

	if previousBucket.startTimestamp.Before(convertedBucket.startTimestamp) {
		// This method should not be used to add buckets out of order.
		// In theory, if the controller loses a large amount of history (i.e. hours), it could try to
		// send out of date old buckets to fill in the gap. Scream about this in the logs, but
		// no need to bring things crashing down over it.
		s.logger.Errorf(
			"Attempted to add bucket with start time %v after last bucket with start time %v, ignoring",
			convertedBucket.startTimestamp, previousBucket.startTimestamp)
		return
	}

	// Add the new bucket to the end of the list.
	s.buckets.PushBack(convertedBucket)

	// Now is as good a time as any to do garbage collection.
	s.garbageCollectBuckets(now)
}

// Get the bucket that is currently being written to. This is always the latest bucket.
func (s *signingRateTracker) getMutableBucket(now time.Time) *Bucket {

	if s.buckets.Size() == 0 {
		// Create the first bucket.
		newBucket := NewBucket(s.logger, now, s.bucketSpan)
		s.buckets.PushBack(newBucket)
	}

	bucket, err := s.buckets.PeekBack()
	enforce.NilError(err, "should be impossible with a non-empty deque")

	if now.After(bucket.EndTimestamp()) {
		// The current bucket's time span has elapsed, create a new bucket.

		bucket = NewBucket(s.logger, now, s.bucketSpan)
		s.buckets.PushBack(bucket)

		// Now is a good time to do garbage collection. As long as bucket size remains fixed, we should be removing
		// one bucket for each new bucket we add once we reach steady state.
		s.garbageCollectBuckets(now)
	}

	return bucket
}

// Remove old buckets that are outside the configured timespan.
func (s *signingRateTracker) garbageCollectBuckets(now time.Time) {
	cutoff := now.Add(-s.timespan)

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
