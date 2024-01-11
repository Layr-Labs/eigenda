package ratelimit

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/common"
)

type BucketStore = common.KVStore[common.RateBucketParams]

type rateLimiter struct {
	globalRateParams common.GlobalRateParams
	bucketStore      BucketStore

	logger common.Logger
}

func NewRateLimiter(rateParams common.GlobalRateParams, bucketStore BucketStore, logger common.Logger) common.RateLimiter {
	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		logger:           logger,
	}
}

// Checks whether a request from the given requesterID is allowed
func (d *rateLimiter) AllowRequest(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	// Retrieve bucket params for the requester ID
	// This will be from dynamo for Disperser and from local storage for DA node
	bucketParams, err := d.bucketStore.GetItem(ctx, requesterID)
	if err != nil {

		bucketLevels := make([]time.Duration, len(d.globalRateParams.BucketSizes))
		copy(bucketLevels, d.globalRateParams.BucketSizes)

		bucketParams = &common.RateBucketParams{
			BucketLevels:    bucketLevels,
			LastRequestTime: time.Now().UTC(),
		}
	}

	// Check whether the request is allowed based on the rate

	// Get interval since last request
	interval := time.Since(bucketParams.LastRequestTime)
	bucketParams.LastRequestTime = time.Now().UTC()

	// Calculate updated bucket levels
	allowed := true
	for i, size := range d.globalRateParams.BucketSizes {

		// Determine bucket deduction
		deduction := time.Microsecond * time.Duration(1e6*float32(blobSize)/float32(rate)/d.globalRateParams.Multipliers[i])

		// Update the bucket level
		bucketParams.BucketLevels[i] = getBucketLevel(bucketParams.BucketLevels[i], size, interval, deduction)
		allowed = allowed && bucketParams.BucketLevels[i] > 0
	}

	// Update the bucket based on blob size and current rate
	if allowed || d.globalRateParams.CountFailed {
		// Update bucket params
		err := d.bucketStore.UpdateItem(ctx, requesterID, bucketParams)
		if err != nil {
			return allowed, err
		}

	}

	return allowed, nil

	// (DA Node) Store the rate params and account ID along with the blob
}

func getBucketLevel(bucketLevel, bucketSize, interval, deduction time.Duration) time.Duration {

	newLevel := bucketLevel + interval - deduction
	if newLevel < 0 {
		newLevel = 0
	}
	if newLevel > bucketSize {
		newLevel = bucketSize
	}

	return newLevel

}
