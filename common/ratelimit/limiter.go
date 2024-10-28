package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type BucketStore = common.KVStore[common.RateBucketParams]

type rateLimiter struct {
	globalRateParams common.GlobalRateParams
	bucketStore      BucketStore

	logger logging.Logger

	// Prometheus metrics
	bucketLevels *prometheus.GaugeVec
}

func NewRateLimiter(reg prometheus.Registerer, rateParams common.GlobalRateParams, bucketStore BucketStore, logger logging.Logger) common.RateLimiter {
	return &rateLimiter{
		globalRateParams: rateParams,
		bucketStore:      bucketStore,
		logger:           logger.With("component", "RateLimiter"),
		bucketLevels: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "rate_limiter_bucket_levels",
			Help: "Current level of each bucket for rate limiting",
		}, []string{"requester_id", "requester_name", "bucket_index"}),
	}
}

// AllowRequest checks whether the request should be allowed. If the request is allowed, the function returns true.
// If the request is not allowed, the function returns false and the RequestParams of the request that was not allowed.
// In order to for the request to be allowed, all of the requests represented by the RequestParams slice must be allowed.
// Each RequestParams object represents a single request. Each request is subjected to the same GlobalRateParams, but the
// individual parameters of the request can differ.
//
// If CountFailed is set to true in the GlobalRateParams, AllowRequest will count failed requests towards the rate limit.
// If CountFailed is set to false, the rate limiter will stop processing requests as soon as it encounters a request that
// is not allowed.
func (d *rateLimiter) AllowRequest(ctx context.Context, params []common.RequestParams) (bool, *common.RequestParams, error) {

	updatedBucketParams := make([]*common.RateBucketParams, len(params))

	allowed := true

	var limitedParam *common.RequestParams

	for i, param := range params {
		allowedForParam, bucketParams := d.checkAllowed(ctx, param)
		updatedBucketParams[i] = bucketParams
		if !allowedForParam {
			allowed = false
			limitedParam = &param

			if !d.globalRateParams.CountFailed {
				break
			}
		}
	}

	if allowed || d.globalRateParams.CountFailed {
		err := d.updateBucketParams(ctx, params, updatedBucketParams)
		if err != nil {
			return false, nil, err
		}
	}

	return allowed, limitedParam, nil

}

func (d *rateLimiter) updateBucketParams(ctx context.Context, params []common.RequestParams, updatedBucketParams []*common.RateBucketParams) error {
	for i, param := range params {
		err := d.bucketStore.UpdateItem(ctx, param.RequesterID, updatedBucketParams[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *rateLimiter) checkAllowed(ctx context.Context, params common.RequestParams) (bool, *common.RateBucketParams) {

	bucketParams, err := d.bucketStore.GetItem(ctx, params.RequesterID)
	if err != nil {

		bucketLevels := make([]time.Duration, len(d.globalRateParams.BucketSizes))
		copy(bucketLevels, d.globalRateParams.BucketSizes)

		bucketParams = &common.RateBucketParams{
			BucketLevels:    bucketLevels,
			LastRequestTime: time.Now().UTC(),
		}
	}

	bucketLevels := make([]time.Duration, len(d.globalRateParams.BucketSizes))

	// Check whether the request is allowed based on the rate

	// Get interval since last request
	interval := time.Since(bucketParams.LastRequestTime)
	lastRequestTime := time.Now().UTC()

	// Calculate updated bucket levels
	allowed := true
	for i, size := range d.globalRateParams.BucketSizes {

		// Determine bucket deduction
		deduction := time.Microsecond * time.Duration(1e6*float32(params.BlobSize)/float32(params.Rate)/d.globalRateParams.Multipliers[i])

		// Update the bucket level
		bucketLevels[i] = getBucketLevel(bucketParams.BucketLevels[i], size, interval, deduction)
		allowed = allowed && bucketLevels[i] > 0

		d.logger.Debug("Bucket level updated", "key", params.RequesterID, "name", params.RequesterName, "prevLevel", bucketParams.BucketLevels[i], "level", bucketLevels[i], "size", size, "interval", interval, "deduction", deduction, "allowed", allowed)

		// Update metrics only if the requester name is provided. We're making
		// an assumption that the requester name is only provided for authenticated
		// requests so it should limit the cardinality of the requester_id label.
		if params.RequesterName != "" {
			d.bucketLevels.With(prometheus.Labels{
				"requester_id":   params.RequesterID,
				"requester_name": params.RequesterName,
				"bucket_index":   strconv.Itoa(i),
			}).Set(float64(bucketLevels[i]))
		}
	}

	bucketParams = &common.RateBucketParams{
		LastRequestTime: lastRequestTime,
		BucketLevels:    bucketLevels,
	}

	return allowed, bucketParams

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
