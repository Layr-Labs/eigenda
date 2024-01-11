package ratelimit

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
)

type BucketStore = common.KVStoreVersioned[common.RateBucketParams]

type rateLimiter struct {
	globalRateParams common.GlobalRateParams

	bucketStore interface{}
	allowlist   []string

	logger common.Logger
}

// Note: This could instead be a factor of RateLimiter with different bucket types
func NewRateLimiter(rateParams common.GlobalRateParams, bucketStore interface{}, allowlist []string, logger common.Logger) (common.RateLimiter, error) {

	if _, isKVStore := bucketStore.(common.KVStore[common.RateBucketParams]); !isKVStore {
		if _, isKVStoreVersioned := bucketStore.(common.KVStoreVersioned[common.RateBucketParams]); !isKVStoreVersioned {
			return nil, errors.New("bucketStore must be either KVStore or KVStoreVersioned")
		}
	}

	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		allowlist:        allowlist,
		logger:           logger,
	}, nil
}

func (d *rateLimiter) AllowRequest(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	// TODO: temporary allowlist that unconditionally allows request
	// for testing purposes only
	for _, id := range d.allowlist {
		if strings.Contains(requesterID, id) {
			return true, nil
		}
	}

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

// func (d *RateLimiter) AllowRequestRunningTimeWorkingDontDelete(ctx context.Context, requesterID RequesterID, blobSize uint, rate RateParam) (bool, error) {
// 	var bs common.KVStoreVersioned[common.RateBucketParams] = d.bucketStore.(common.KVStoreVersioned[common.RateBucketParams])
// 	var version int = 0

// 	now := time.Now().UTC()

// 	// Calculate the required capacity for new request
// 	requiredCapacity := time.Duration(float64(blobSize)/float64(rate)) * time.Second

// 	// Fetch the current bucket parameters
// 	bucketParams, err := bs.GetItemWithVersion(ctx, requesterID)
// 	if err != nil {
// 		return false, fmt.Errorf("error fetching bucket params: %v", err)
// 	}

// 	// Initialize bucket parameters if they don't exist
// 	if bucketParams == nil {
// 		bucketParams = &common.RateBucketParams{
// 			BucketLevels:    make([]time.Duration, len(d.GlobalParams.BucketSizes)),
// 			LastRequestTime: now,
// 		}
// 		// Initialize BucketLevels to their maximum sizes
// 		for i, size := range d.GlobalParams.BucketSizes {
// 			bucketParams.BucketLevels[i] = size
// 		}
// 		// Optionally persist the initial state ??
// 	}

// 	// Conservative check: Ensure that the request can potentially be allowed
// 	canBeAllowed := true
// 	for i, bucketSize := range d.GlobalParams.BucketSizes {
// 		if bucketLevel[i]+requiredCapacity < bucketSize {
// 			canBeAllowed = false
// 			break
// 		}
// 	}

// 	if !canBeAllowed {
// 		return false, nil
// 	}

// 	// Build the update expression for each bucket level
// 	updateBuilder := expression.UpdateBuilder{}
// 	for i := range d.GlobalParams.BucketSizes {
// 		bucketLevelKey := fmt.Sprintf("BucketLevels[%d]", i)

// 		if updateKeySeparately {
// 			updateBuilder = updateBuilder.Add(
// 				expression.Name(bucketLevelKey),
// 				expression.Value(-requiredCapacity.Milliseconds()),
// 			)
// 			updateBuilder = updateBuilder.Set(
// 				expression.Name("LastRequestTime"),
// 				expression.Value(now.UnixMilli()),
// 			)

// 			err = d.Store.UpdateItemWithExpression(ctx, requesterID, expr)
// 			if err != nil {
// 				return false, fmt.Errorf("failed to update bucket level %d atomically: %v", i, err)
// 			}
// 			updateBuilder = expression.UpdateBuilder{} // Reset builder for next iteration
// 		} else {
// 			updateBuilder = updateBuilder.Add(
// 				expression.Name(bucketLevelKey),
// 				expression.Value(-requiredCapacity.Milliseconds()),
// 			)
// 		}
// 	}

// 	if !updateKeySeparately {
// 		// Update last request time
// 		updateBuilder = updateBuilder.Set(
// 			expression.Name("LastRequestTime"),
// 			expression.Value(now.UnixMilli()),
// 		)
// 		err = d.Store.UpdateItemWithExpression(ctx, requesterID, expr)
// 		if err != nil {
// 			return false, fmt.Errorf("failed to update bucket levels atomically: %v", err)
// 		}
// 	}

// 	// Return true, assuming the request is allowed after the update
// 	return true, nil
// }

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
