package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func makeTestRatelimiter() (common.RateLimiter, error) {

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{time.Second, time.Minute},
		Multipliers: []float32{1, 1},
	}
	bucketStoreSize := 1000

	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](bucketStoreSize)
	if err != nil {
		return nil, err
	}

	ratelimiter := ratelimit.NewRateLimiter(prometheus.NewRegistry(), globalParams, bucketStore, logging.NewNoopLogger())

	return ratelimiter, nil

}

func TestRatelimit(t *testing.T) {

	ratelimiter, err := makeTestRatelimiter()
	assert.NoError(t, err)

	ctx := context.Background()

	retreiverID := "testRetriever"

	params := []common.RequestParams{
		{
			RequesterID: retreiverID,
			BlobSize:    10,
			Rate:        100,
		},
	}

	for i := 0; i < 10; i++ {
		allow, _, err := ratelimiter.AllowRequest(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, true, allow)
	}

	allow, _, err := ratelimiter.AllowRequest(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, false, allow)
}
