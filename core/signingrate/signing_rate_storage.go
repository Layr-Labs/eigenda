package signingrate

import "time"

// Responsible for storing historical signing rate information in a manner that is restart/crash safe.
type SigningRateStorage interface {
	// Store one or more buckets. If a bucket with the same start time already exists, it will be overwritten.
	StoreBuckets(buckets []*SigningRateBucket) error

	// Load all buckets with data starting at or after startTimestamp.
	LoadBuckets(startTimestamp time.Time) ([]*SigningRateBucket, error)
}
