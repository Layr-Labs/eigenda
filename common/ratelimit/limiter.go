package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/cenkalti/backoff/v4"
)

type BucketStoreConcurrencySafe = common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe]
type BucketStore = common.KVStore[common.RateBucketParams]

type rateLimiter struct {
	globalRateParams common.GlobalRateParams

	bucketStore      interface{}
	logger           common.Logger
	retryWithBackOff backoff.BackOff
}

// NewRateLimiter is a factory function to create new RateLimiter
func NewRateLimiter(rateParams common.GlobalRateParams, bucketStore interface{}, logger common.Logger) (common.RateLimiter, error) {

	retryWithBackOff := backoff.NewExponentialBackOff()
	retryWithBackOff.MaxElapsedTime = 5 * time.Minute
	retryWithBackOff.MaxInterval = 15 * time.Second
	retryWithBackOff.InitialInterval = 1 * time.Second

	switch bs := bucketStore.(type) {
	case common.KVStore[common.RateBucketParams]:
		return newRateLimiter(rateParams, bs, logger, retryWithBackOff), nil
	case common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe]:
		return newVersionedConcurrencySafeRateLimiter(rateParams, bs, logger, retryWithBackOff), nil
	default:
		return nil, errors.New("unsupported bucketStore type")
	}
}

// Specific constructors for each type of RateLimiter
func newRateLimiter(rateParams common.GlobalRateParams, bucketStore common.KVStore[common.RateBucketParams], logger common.Logger, retryWithBackOff backoff.BackOff) common.RateLimiter {
	// initialize and return a simple rate limiter
	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		logger:           logger,
		retryWithBackOff: retryWithBackOff,
	}
}

func newVersionedConcurrencySafeRateLimiter(rateParams common.GlobalRateParams, bucketStore common.KVStoreVersioned[common.RateBucketParamsConcurrencySafe], logger common.Logger, retryWithBackOff backoff.BackOff) common.RateLimiter {
	// initialize and return a versioned, concurrency-safe rate limiter
	// initialize and return a simple rate limiter
	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		logger:           logger,
		retryWithBackOff: retryWithBackOff,
	}
}

// Checks whether a request from the given requesterID is allowed
func (d *rateLimiter) AllowRequest(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	var bs BucketStore = d.bucketStore.(BucketStore)

	// Retrieve bucket params for the requester ID
	// This will be from dynamo for Disperser and from local storage for DA node
	bucketParams, err := bs.GetItem(ctx, requesterID)
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

		prevLevel := bucketParams.BucketLevels[i]

		// Update the bucket level
		bucketParams.BucketLevels[i] = getBucketLevel(bucketParams.BucketLevels[i], size, interval, deduction)
		allowed = allowed && bucketParams.BucketLevels[i] > 0

		d.logger.Debug("Bucket level", "key", requesterID, "prevLevel", prevLevel, "level", bucketParams.BucketLevels[i], "size", size, "interval", interval, "deduction", deduction, "allowed", allowed)
	}

	// Update the bucket based on blob size and current rate
	if allowed || d.globalRateParams.CountFailed {
		// Update bucket params
		err := bs.UpdateItem(ctx, requesterID, bucketParams)
		if err != nil {
			return allowed, err
		}

	}

	return allowed, nil

	// (DA Node) Store the rate params and account ID along with the blob
}

// RateLimiting Algorithm defined as per this specification: https://docs.google.com/document/d/1KbGTRfLYX0buD0baCCnh97L01QeHwqsxDeGS1IrBANc
// Gist is to use running time.
// Bucket is initialized with current time minus the global bucket size
// If BucketLevel + Required Capacity (calclulated as blobSize/rate) is less than now then the request can be allowed
// Versioning makes it optimisically safe
// Using DynamodB Atomic ADD operation to update BucketLevel as opposed to GET and then SET
// TODO: Add a Flag to switch between the two implementations
func (d *rateLimiter) AllowRequestConcurrencySafeVersion(ctx context.Context, requesterID common.RequesterID, blobSize uint, rate common.RateParam) (bool, error) {
	fmt.Println("AllowRequestConcurrencySafeVersion")
	var bs BucketStoreConcurrencySafe = d.bucketStore.(BucketStoreConcurrencySafe)

	if bs == nil {
		return false, errors.New("bucketStore must be KVStoreVersioned")
	}
	// Fetch the current bucket parameters
	bucketParams, _, err := bs.GetItemWithVersion(ctx, requesterID)
	if err != nil {
		fmt.Printf("err %v\n", err)
		fmt.Printf("Item not found for requesterID %v\n", requesterID)
	}
	fmt.Printf("bucketParams %v\n", bucketParams)
	// Initialize bucket parameters if they don't exist
	if bucketParams == nil {
		fmt.Printf("Initialize BucketParams\n")
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
		// Save BucketParams for Requester
		err := bs.UpdateItem(ctx, requesterID, bucketParams)

		if err != nil {
			return false, fmt.Errorf("failed to save bucket params dynamodB: %v", err)
		}
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
		requiredCapacity := uint64(time.Duration(uint32(blobSize) / rate))
		// Since bucketsize is subtracted from now which is assigned to bucketLevel[i] above
		// Now check if the required capacity + bucketLevel[i] is less than now
		// If it is less than now then the request can be allowed
		totalCapacity := bucketLevels[i] + requiredCapacity
		allowed = totalCapacity < uint64(now)
		fmt.Printf("allowed %v\n", allowed)

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
				// Note this will result in a call to UpdateItem for each bucket level
				// This should be revisited if needed later as it may result in too many calls to dynamodb
				// Other alternative is to make all updates to each level as an array and then call UpdateItem once but that does not make it concurrency safe
				err = bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
				if err != nil {
					return false, fmt.Errorf("failed to update bucket level %d atomically: %v with required capacity %v", i, err, requiredCapacity)
				}

				operationUpdateItemWithExpression := func() error {
					err := bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
					if err != nil {
						// Return transient error to trigger a retry
						if isConcurrentUpdateError(err) {
							return err
						}
					}
					return nil // No error, no retry
				}

				// Retry with exponential backoff
				err := backoff.Retry(operationUpdateItemWithExpression, d.retryWithBackOff)
				if err != nil {
					return false, fmt.Errorf("failed after retries: %w", err)
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
				// Note this will result in a call to UpdateItem for each bucket level
				// This should be revisited if needed later as it may result in too many calls to dynamodb
				err = bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
				if err != nil {
					return false, fmt.Errorf("failed to update bucket level %d atomically: %v", i, err)
				}

				operationUpdateItemWithExpression := func() error {
					err := bs.UpdateItemWithExpression(ctx, requesterID, &updateBuilder)
					if err != nil {
						// Return transient error to trigger a retry
						if isConcurrentUpdateError(err) {
							return err
						}
					}
					// For other errors, do not retry
					return nil // No error, no retry
				}

				// Retry with exponential backoff
				err := backoff.Retry(operationUpdateItemWithExpression, d.retryWithBackOff)
				if err != nil {
					return false, fmt.Errorf("failed after retries: %w", err)
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

func isConcurrentUpdateError(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		// Check if the error code is ValidationException, which is used for a variety of input errors
		if awsErr.Code() == "ValidationException" {
			// Check if the message contains the specific error detail
			return strings.Contains(awsErr.Message(), "Two document paths overlap with each other")
		}
	}
	return false
}
