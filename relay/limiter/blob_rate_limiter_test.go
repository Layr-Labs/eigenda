package limiter

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func defaultConfig() *Config {
	return &Config{
		MaxGetBlobOpsPerSecond:          1024,
		GetBlobOpsBurstiness:            1024,
		MaxGetBlobBytesPerSecond:        20 * 1024 * 1024,
		GetBlobBytesBurstiness:          20 * 1024 * 1024,
		MaxConcurrentGetBlobOps:         1024,
		MaxGetChunkOpsPerSecond:         1024,
		GetChunkOpsBurstiness:           1024,
		MaxGetChunkBytesPerSecond:       20 * 1024 * 1024,
		GetChunkBytesBurstiness:         20 * 1024 * 1024,
		MaxConcurrentGetChunkOps:        1024,
		MaxGetChunkOpsPerSecondClient:   8,
		GetChunkOpsBurstinessClient:     8,
		MaxGetChunkBytesPerSecondClient: 2 * 1024 * 1024,
		GetChunkBytesBurstinessClient:   2 * 1024 * 1024,
		MaxConcurrentGetChunkOpsClient:  1,
	}
}

func TestConcurrentBlobOperations(t *testing.T) {
	tu.InitializeRandom()

	concurrencyLimit := 1 + rand.Intn(10)

	config := defaultConfig()
	config.MaxConcurrentGetBlobOps = concurrencyLimit
	// Make the burstiness limit high enough that we won't be rate limited
	config.GetBlobOpsBurstiness = concurrencyLimit * 100

	limiter := NewBlobRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// We should be able to start this many operations concurrently
	for i := 0; i < concurrencyLimit; i++ {
		err := limiter.BeginGetBlobOperation(now)
		require.NoError(t, err)
	}

	// Starting one more operation should fail due to the concurrency limit
	err := limiter.BeginGetBlobOperation(now)
	require.Error(t, err)

	// Finish an operation. This should permit exactly one more operation to start
	limiter.FinishGetBlobOperation()
	err = limiter.BeginGetBlobOperation(now)
	require.NoError(t, err)
	err = limiter.BeginGetBlobOperation(now)
	require.Error(t, err)
}

func TestGetBlobOpRateLimit(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetBlobOpsPerSecond = float64(2 + rand.Intn(10))
	config.GetBlobOpsBurstiness = int(config.MaxGetBlobOpsPerSecond) + rand.Intn(10)
	config.MaxConcurrentGetBlobOps = 1

	limiter := NewBlobRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// Without advancing time, we should be able to perform a number of operations equal to the burstiness limit.
	for i := 0; i < config.GetBlobOpsBurstiness; i++ {
		err := limiter.BeginGetBlobOperation(now)
		require.NoError(t, err)
		limiter.FinishGetBlobOperation()
	}

	// We are not at the rate limit, and should be able to start another operation.
	err := limiter.BeginGetBlobOperation(now)
	require.Error(t, err)

	// Advance time by one second. We should gain a number of tokens equal to the rate limit.
	now = now.Add(time.Second)
	for i := 0; i < int(config.MaxGetBlobOpsPerSecond); i++ {
		err = limiter.BeginGetBlobOperation(now)
		require.NoError(t, err)
		limiter.FinishGetBlobOperation()
	}

	// We have once again hit the rate limit. We should not be able to start another operation.
	err = limiter.BeginGetBlobOperation(now)
	require.Error(t, err)

	// Advance time by another second. We should gain another number of tokens equal to the rate limit.
	// Intentionally do not finish the next operation. We are attempting to get a failure by exceeding
	// the max concurrent operations limit.
	now = now.Add(time.Second)
	err = limiter.BeginGetBlobOperation(now)
	require.NoError(t, err)

	// This operation should fail since we have limited concurrent operations to 1. It should not count
	// against the rate limit.
	err = limiter.BeginGetBlobOperation(now)
	require.Error(t, err)

	// "finish" the prior operation. Verify that we have all expected tokens available.
	limiter.FinishGetBlobOperation()
	for i := 0; i < int(config.MaxGetBlobOpsPerSecond)-1; i++ {
		err = limiter.BeginGetBlobOperation(now)
		require.NoError(t, err)
		limiter.FinishGetBlobOperation()
	}

	// We should now be at the rate limit. We should not be able to start another operation.
	err = limiter.BeginGetBlobOperation(now)
	require.Error(t, err)
}

func TestGetBlobBandwidthLimit(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetBlobBytesPerSecond = float64(1024 + rand.Intn(1024*1024))
	config.GetBlobBytesBurstiness = int(config.MaxGetBlobBytesPerSecond) + rand.Intn(1024*1024)

	limiter := NewBlobRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// Without advancing time, we should be able to utilize a number of bytes equal to the burstiness limit.
	bytesRemaining := config.GetBlobBytesBurstiness
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err := limiter.RequestGetBlobBandwidth(now, uint32(bytesToRequest))
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err := limiter.RequestGetBlobBandwidth(now, 1)
	require.Error(t, err)

	// Advance time by one second. We should gain a number of tokens equal to the rate limit.
	now = now.Add(time.Second)
	bytesRemaining = int(config.MaxGetBlobBytesPerSecond)
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetBlobBandwidth(now, uint32(bytesToRequest))
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err = limiter.RequestGetBlobBandwidth(now, 1)
	require.Error(t, err)
}
