package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type BucketStore = common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe]

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

// RateLimiting Algorithm defined as per this specification: https://docs.google.com/document/d/1KbGTRfLYX0buD0baCCnh97L01QeHwqsxDeGS1IrBANc
// Gist is to use running time.
// Bucket is initialized with current time minus the global bucket size
// If BucketLevel + Required Capacity (calclulated as blobSize/rate) is less than now then the request can be allowed
// Versioning makes it optimisically safe
// Using DynamodB Atomic ADD operation to update BucketLevel as opposed to GET and then SET
// TODO: Add a Flag to switch between the two implementations
func (d *rateLimiter) AllowRequestConcurrencySafeVersion(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	var bs common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe] = d.bucketStore.(common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe])

	// Fetch the current bucket parameters
	bucketParams, _, err := bs.GetItemWithVersion(ctx, requesterID)
	if err != nil {
		return false, fmt.Errorf("error fetching bucket params: %v", err)
	}

	// Initialize bucket parameters if they don't exist
	if bucketParams == nil {
		bucketParams = &common.RateBucketParamsConcurrencySafe{
			BucketLevels: make([]uint64, len(d.globalRateParams.BucketSizes)),
		}

		now := time.Now().UTC().UnixMilli()
		// Initialize BucketLevels to their maximum sizes
		for i, size := range d.globalRateParams.BucketSizes {

			// Initialize BucketLevel to the current time minus the bucket size
			// Converting to Uint64 to avoid overflow
			bucketParams.BucketLevels[i] = uint64(now - size.Milliseconds())
		}
		// Optionally persist the initial state ??
		// TODO: Double Check if these need to be set before checking BucketSize
	}

	// Conservative check: Ensure that the request can potentially be allowed
	allowed := false
	reachedMaxBucketSize := false
	now := time.Now().UTC().UnixMilli()
	bucketLevels := bucketParams.BucketLevels
	for i, bucketSize := range d.globalRateParams.BucketSizes {
		// BucketLevel is initialized as the current time minus the bucket size, so if the bucket level is greater than
		// the current time + required capacity it will overflow bucketsize as a result the request should not be allowed
		calclulatedSize := uint64(now - bucketSize.Milliseconds())
		if bucketLevels[i] < calclulatedSize {
			reachedMaxBucketSize = true
			bucketLevels[i] = calclulatedSize
			//break
		}

		// Calculate the required capacity for new request
		requiredCapacity := uint64(float64(blobSize) / float64(rate) * float64(time.Second.Milliseconds()))
		// Since bucketsize is subtracted from now which is assigned to bucketLevel[i] above
		// Now check if the required capacity + bucketLevel[i] is less than now
		// If it is less than now then the request can be allowed
		totalCapacity := bucketLevels[i] + requiredCapacity
		allowed = totalCapacity <= uint64(now)

		if allowed || d.globalRateParams.CountFailed {
			if !reachedMaxBucketSize {

				// Update Bucket Level by required capacity as it is still less than now
				updateBuilder := expression.Add(
					expression.Name(fmt.Sprintf("BucketLevels[%d]", i)),
					expression.Value(requiredCapacity),
				).Add(
					expression.Name("Version"),
					expression.Value(1),
				)
				// TODO: Specifically Handle Case when this fails for same document being updated by multiple requests
				// Handle this:error StatusCode: 400, RequestID: e5764d48-e0fd-46bc-8aca-42de3d38c055, api error ValidationException:
				// Invalid UpdateExpression: Two document paths overlap with each other; must remove or rewrite one of these paths; path one: [Version], path two: [Version]
				// Note this will result in a call to UpdateItem for each bucket level
				// This should be revisited if needed later as it may result in too many calls to dynamodb
				// Other alternative is to make all updates to each level as an array and then call UpdateItem once but that does not make it concurrency safe
				err = bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
				if err != nil {
					return false, fmt.Errorf("failed to update bucket level %d atomically: %v", i, err)
				}
				return allowed, nil
			} else {

				// Set
				updateBuilder := expression.Set(
					expression.Name(fmt.Sprintf("BucketLevels[%d]", i)),
					expression.Value(totalCapacity),
				).Add(
					expression.Name("Version"),
					expression.Value(1),
				)
				// TODO: Specifically Handle Case when this fails for same document being updated by multiple requests
				// Handle this:error StatusCode: 400, RequestID: e5764d48-e0fd-46bc-8aca-42de3d38c055, api error ValidationException:
				// Invalid UpdateExpression: Two document paths overlap with each other; must remove or rewrite one of these paths; path one: [Version], path two: [Version]
				// Note this will result in a call to UpdateItem for each bucket level
				// This should be revisited if needed later as it may result in too many calls to dynamodb
				err = bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
				if err != nil {
					return false, fmt.Errorf("failed to update bucket level %d atomically: %v", i, err)
				}
			}
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
