package signingrate

import (
	"fmt"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ SigningRateTracker = (*signingRateTracker)(nil)

// A standard implementation of the SigningRateTracker interface. Is not thread safe on its own.
type signingRateTracker struct {
	logger logging.Logger

	// Signing data storage, split up into buckets for each time interval. Buckets are stored in chronological order.
	buckets *common.RandomAccessDeque[*SigningRateBucket]

	// Buckets that have not yet been flushed to storage. Keyed by the bucket's start time.
	unflushedBuckets map[time.Time]*SigningRateBucket

	// The length of time to keep loaded in memory.
	timeSpan time.Duration

	// The duration of each bucket. Buckets loaded from storage may have different spans, but new buckets will
	// always have this span.
	bucketSpan time.Duration

	// A function that returns the current time.
	timeSource func() time.Time
}

// Create a new SigningRateTracker.
func NewSigningRateTracker(
	logger logging.Logger,
	// The amount of time to keep in memory. Queries are only supported for this timeSpan.
	timeSpan time.Duration,
	// The duration of each bucket
	bucketSpan time.Duration,
	timeSource func() time.Time,
) (SigningRateTracker, error) {

	if timeSpan.Seconds() < 1 {
		return nil, fmt.Errorf("time span must be at least one second, got %s", timeSpan)
	}
	if bucketSpan.Seconds() < 1 {
		return nil, fmt.Errorf("bucket span must be at least one second, got %s", bucketSpan)
	}

	store := &signingRateTracker{
		logger:           logger,
		buckets:          common.NewRandomAccessDeque[*SigningRateBucket](0),
		timeSpan:         timeSpan,
		bucketSpan:       bucketSpan,
		unflushedBuckets: make(map[time.Time]*SigningRateBucket),
		timeSource:       timeSource,
	}

	return store, nil
}

// Report that a validator has successfully signed a batch of the given size.
func (s *signingRateTracker) ReportSuccess(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {
	now := s.timeSource()

	bucket := s.getMutableBucket(now)
	bucket.ReportSuccess(quorum, validatorID, batchSize, signingLatency)
	s.markUnflushed(bucket)

	s.garbageCollectBuckets(now)
}

// Report that a validator has failed to sign a batch of the given size.
func (s *signingRateTracker) ReportFailure(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	batchSize uint64,
) {
	now := s.timeSource()

	bucket := s.getMutableBucket(now)
	bucket.ReportFailure(quorum, validatorID, batchSize)
	s.markUnflushed(bucket)

	s.garbageCollectBuckets(now)
}

func (s *signingRateTracker) GetValidatorSigningRate(
	quorum core.QuorumID,
	validatorID core.OperatorID,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {

	if !endTime.After(startTime) {
		return nil, fmt.Errorf("end time %v is not after start time %v", endTime, startTime)
	}

	if s.buckets.Size() == 0 {
		// Special case: no data available.
		return &validator.ValidatorSigningRate{
			ValidatorId: validatorID[:],
		}, nil
	}

	comparator := func(timestamp time.Time, bucket *SigningRateBucket) int {
		unixTimestamp := timestamp.Unix()

		if unixTimestamp < bucket.startTimestamp.Unix() {
			return -1
		} else if unixTimestamp >= bucket.endTimestamp.Unix() {
			// unixTimestamp == bucket.endTimestamp.Unix(), then timestamp is "after" the bucket since end is exclusive
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
		ValidatorId: validatorID[:],
	}

	iterator, err := s.buckets.IteratorFrom(startIndex)
	enforce.NilError(err, "should be impossible with a valid index")
	for _, bucket := range iterator {
		if !bucket.startTimestamp.Before(endTime) {
			break
		}

		signingRate, exists := bucket.getValidatorIfExists(quorum, validatorID)
		if !exists {
			// No info for validator during this bucket, skip it.
			continue
		}

		totalSigningRate.SignedBatches += signingRate.GetSignedBatches()
		totalSigningRate.UnsignedBatches += signingRate.GetUnsignedBatches()
		totalSigningRate.SignedBytes += signingRate.GetSignedBytes()
		totalSigningRate.UnsignedBytes += signingRate.GetUnsignedBytes()
		totalSigningRate.SigningLatency += signingRate.GetSigningLatency()
	}

	return totalSigningRate, nil
}

func (s *signingRateTracker) GetSigningRateDump(
	startTime time.Time,
) ([]*validator.SigningRateBucket, error) {

	buckets := make([]*validator.SigningRateBucket, 0, s.buckets.Size())

	// Iterate backwards. In general, dump requests will only be used to fetch recent data, so
	// we should optimize the case where we are requesting a few buckets from the end of the deque.
	// Worst case scenario, we iterate the entire deque. If we do that, we are about to transmit the contents
	// of the deque over a network connection. And so in that case, the cost of iteration doesn't really matter.
	for _, bucket := range s.buckets.ReverseIterator() {
		if !bucket.EndTimestamp().After(startTime) {
			// This bucket is too old, skip it and stop iterating.
			break
		}
		buckets = append(buckets, bucket.ToProtobuf())
	}

	// We iterated in reverse, so reverse again to get chronological ordering.
	slices.Reverse(buckets)

	return buckets, nil
}

func (s *signingRateTracker) GetUnflushedBuckets() ([]*validator.SigningRateBucket, error) {
	buckets := make([]*validator.SigningRateBucket, 0, len(s.unflushedBuckets))

	for _, bucket := range s.unflushedBuckets {
		proto := bucket.ToProtobuf()
		buckets = append(buckets, proto)
	}
	s.unflushedBuckets = make(map[time.Time]*SigningRateBucket)

	sortValidatorSigningRateBuckets(buckets)

	return buckets, nil
}

func (s *signingRateTracker) UpdateLastBucket(bucket *validator.SigningRateBucket) {
	convertedBucket := NewBucketFromProto(bucket)

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

	if convertedBucket.startTimestamp.Before(previousBucket.startTimestamp) {
		// This method should not be used to add buckets out of order.
		// But no need to crash if it happens, just ignore the request.
		s.logger.Errorf(
			"Attempted to add bucket with start time %v after last bucket with start time %v, ignoring",
			convertedBucket.startTimestamp, previousBucket.startTimestamp)
		return
	}

	// Add the new bucket to the end of the list.
	s.buckets.PushBack(convertedBucket)

	s.garbageCollectBuckets(s.timeSource())
}

func (s *signingRateTracker) GetLastBucketStartTime() (time.Time, error) {
	if s.buckets.Size() == 0 {
		return time.Time{}, nil
	}
	bucket, err := s.buckets.PeekBack()
	enforce.NilError(err, "should be impossible with a non-empty deque")
	return bucket.startTimestamp, nil
}

func (s *signingRateTracker) Flush() error {
	// Intentional no-op, as this implementation is synchronous.
	return nil
}

// Get the bucket that is currently being written to. This is always the latest bucket.
func (s *signingRateTracker) getMutableBucket(now time.Time) *SigningRateBucket {

	if s.buckets.Size() == 0 {
		// Create the first bucket.
		newBucket, err := NewSigningRateBucket(now, s.bucketSpan)
		enforce.NilError(err, "should be impossible with a valid bucket span")
		s.buckets.PushBack(newBucket)
	}

	bucket, err := s.buckets.PeekBack()
	enforce.NilError(err, "should be impossible with a non-empty deque")

	if !bucket.Contains(now) {
		// The current bucket's time span has elapsed, create a new bucket.

		bucket, err = NewSigningRateBucket(now, s.bucketSpan)
		enforce.NilError(err, "should be impossible with a valid bucket span")
		s.buckets.PushBack(bucket)

		// Now is a good time to do garbage collection. As long as bucket size remains fixed, we should be removing
		// one bucket for each new bucket we add once we reach steady state.
		s.garbageCollectBuckets(now)
	}

	return bucket
}

// Remove old buckets that are outside the configured timeSpan.
func (s *signingRateTracker) garbageCollectBuckets(now time.Time) {
	cutoff := now.Add(-s.timeSpan)

	for s.buckets.Size() > 0 {
		bucket, err := s.buckets.PeekFront()
		enforce.NilError(err, "should be impossible with a non-empty deque")

		if cutoff.Before(bucket.EndTimestamp()) {
			// This bucket is new enough, so all later buckets will also be new enough.
			break
		}

		// This bucket is too old, remove it.
		_, err = s.buckets.PopFront()
		enforce.NilError(err, "should be impossible with a non-empty deque")
	}
}

// Mark a bucket as needing to be flushed to storage.
func (s *signingRateTracker) markUnflushed(bucket *SigningRateBucket) {
	s.unflushedBuckets[bucket.startTimestamp] = bucket
}
