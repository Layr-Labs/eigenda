package ratelimit

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func makeTestRatelimiter(t *testing.T) (common.RateLimiter, error) {
	t.Helper()

	logger := test.GetLogger()

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{time.Second, time.Minute},
		Multipliers: []float32{1, 1},
	}
	bucketStoreSize := 1000

	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](bucketStoreSize)
	if err != nil {
		return nil, err
	}

	ratelimiter := NewRateLimiter(prometheus.NewRegistry(), globalParams, bucketStore, logger)

	return ratelimiter, nil

}

func TestRatelimit(t *testing.T) {
	ctx := t.Context()

	ratelimiter, err := makeTestRatelimiter(t)
	require.NoError(t, err)

	retrieverID := "testRetriever"

	params := []common.RequestParams{
		{
			RequesterID: retrieverID,
			BlobSize:    10,
			Rate:        100,
		},
	}

	for i := 0; i < 10; i++ {
		allow, _, err := ratelimiter.AllowRequest(ctx, params)
		require.NoError(t, err)
		require.Equal(t, true, allow)
	}

	allow, _, err := ratelimiter.AllowRequest(ctx, params)
	require.NoError(t, err)
	require.Equal(t, false, allow)
}
