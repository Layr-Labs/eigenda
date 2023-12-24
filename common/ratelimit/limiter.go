package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/common"
)

type BucketStore = common.KVStoreVersioned[common.RateBucketParams]

type rateLimiter struct {
	globalRateParams common.GlobalRateParams

	bucketStore interface{}

	logger common.Logger
}

func NewRateLimiter(rateParams common.GlobalRateParams, bucketStore interface{}, logger common.Logger) (common.RateLimiter, error) {

	if _, isKVStore := bucketStore.(common.KVStore[common.RateBucketParams]); !isKVStore {
		if _, isKVStoreVersioned := bucketStore.(common.KVStoreVersioned[common.RateBucketParams]); !isKVStoreVersioned {
			return nil, errors.New("bucketStore must be either KVStore or KVStoreVersioned")
		}
	}

	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		logger:           logger,
	}, nil
}

func (d *rateLimiter) AllowRequest(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	// Retrieve bucket params for the requester ID
	var bucketParams *common.RateBucketParams
	var err error
	var version int = 0

	// Determine the type of bucketStore and use the appropriate GetItem method
	switch bs := d.bucketStore.(type) {
	case common.KVStoreVersioned[common.RateBucketParams]:
		bucketParams, version, err = bs.GetItemWithVersion(ctx, requesterID)
	case common.KVStore[common.RateBucketParams]:
		bucketParams, err = bs.GetItem(ctx, requesterID)
	default:
		return false, errors.New("unknown bucketStore type")
	}

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
		switch bs := d.bucketStore.(type) {
		case common.KVStoreVersioned[common.RateBucketParams]:
			// Use UpdateItemWithVersion for KVStoreVersioned
			err = bs.UpdateItemWithVersion(ctx, requesterID, bucketParams, version)
		case common.KVStore[common.RateBucketParams]:
			// Use UpdateItem for KVStore
			err = bs.UpdateItem(ctx, requesterID, bucketParams)
		default:
			return false, errors.New("unknown bucketStore type")
		}

		if err != nil {
			return allowed, err
		}
	}

	return allowed, nil
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
