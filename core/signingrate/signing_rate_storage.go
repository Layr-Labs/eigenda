package signingrate

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Responsible for storing historical signing rate information in a manner that is restart/crash safe.
type SigningRateStorage interface {
	// Store one or more buckets. If a bucket with the same start time already exists, it will be overwritten.
	StoreBuckets(ctx context.Context, buckets []*validator.SigningRateBucket) error

	// Load all buckets that contain any data from after the provided startTimestamp. A bucket is returned
	// even if it also has some data that is before the startTimestamp, so long as it also contains data after it.
	LoadBuckets(ctx context.Context, startTimestamp time.Time) ([]*validator.SigningRateBucket, error)
}
