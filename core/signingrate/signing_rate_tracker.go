package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
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
	GetSigningRateDump(startTime time.Time) ([]*controller.SigningRateBucket, error)

	// Returns a list of buckets that have not yet been flushed to persistent storage.
	// Buckets are in chronological order. Allows for an external process to periodically
	// flush data in this tracker to persistent storage.
	//
	// Returned data threadsafe to read, but should not be modified.
	GetUnflushedBuckets() ([]*controller.SigningRateBucket, error)

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
		timeout bool,
	)

	// Update a bucket, overwriting an existing bucket with the same start time if it is present. Should
	// only be used to update the last bucket in the store. Data is ignored if the bucket won't be the
	// last bucket.
	//
	// The intended use of this method is to set up a SigningRateTracker that mirrors a remote SigningRateTracker.
	// The remote tracker is the source of truth, and this local tracker is just a cache. Periodically, get data
	// from the remote tracker using GetSigningRateDump(), and then insert the data returned into this tracker using
	// UpdateLastBucket().
	//
	// This operation doesn't mark a bucket as unflushed. A bucket is only marked as unflushed when it is modified,
	// not when it is provided whole-sale from an external source.
	UpdateLastBucket(now time.Time, bucket *controller.SigningRateBucket)

	// Get the start time of the last bucket in the store. If the store is empty, returns the zero time.
	// Useful for determining how much data to request from a remote store when mirroring.
	GetLastBucketStartTime() (time.Time, error)

	// Several methods on this interface may asynchronously modify internal state. This method blocks
	// until all previously queued modifications have been applied.
	Flush() error

	// Close the store and free any associated resources.
	Close()
}
