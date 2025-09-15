package signingrate

import (
	"context"
	"time"
)

// Responsible for storing historical signing rate information in a manner that is restart/crash safe.
type SigningRateStorage interface {
	// Store one or more buckets. If a bucket with the same start time already exists, it will be overwritten.
	StoreBuckets(ctx context.Context, buckets []*SigningRateBucket) error

	// Load all buckets with data starting at or after startTimestamp.
	LoadBuckets(ctx context.Context, startTimestamp time.Time) ([]*SigningRateBucket, error)
}
