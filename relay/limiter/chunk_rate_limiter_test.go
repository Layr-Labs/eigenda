package limiter

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"math"
	"testing"
	"time"
)

func TestConcurrentGetChunksOperations(t *testing.T) {
	tu.InitializeRandom()

	concurrencyLimit := 1 + rand.Intn(10)

	config := defaultConfig()
	config.MaxConcurrentGetChunkOps = concurrencyLimit
	config.MaxConcurrentGetChunkOpsClient = math.MaxInt32
	config.GetChunkOpsBurstiness = math.MaxInt32
	config.GetChunkOpsBurstinessClient = math.MaxInt32

	userID := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// We should be able to start this many operations concurrently
	for i := 0; i < concurrencyLimit; i++ {
		err := limiter.BeginGetChunkOperation(now, userID)
		require.NoError(t, err)
	}

	// Starting one more operation should fail due to the concurrency limit
	err := limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)

	// Finish an operation. This should permit exactly one more operation to start
	limiter.FinishGetChunkOperation(userID)
	err = limiter.BeginGetChunkOperation(now, userID)
	require.NoError(t, err)
	err = limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)
}

func TestGetChunksRateLimit(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetChunkOpsPerSecond = float64(2 + rand.Intn(10))
	config.GetChunkOpsBurstiness = int(config.MaxGetChunkOpsPerSecond) + rand.Intn(10)
	config.GetChunkOpsBurstinessClient = math.MaxInt32
	config.MaxConcurrentGetChunkOps = 1

	userID := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// Without advancing time, we should be able to perform a number of operations equal to the burstiness limit.
	for i := 0; i < config.GetChunkOpsBurstiness; i++ {
		err := limiter.BeginGetChunkOperation(now, userID)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID)
	}

	// We are now at the rate limit, and should not be able to start another operation.
	err := limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)

	// Advance time by one second. We should now be able to perform a number of operations equal to the rate limit.
	now = now.Add(time.Second)
	for i := 0; i < int(config.MaxGetChunkOpsPerSecond); i++ {
		err = limiter.BeginGetChunkOperation(now, userID)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID)
	}

	// We are now at the rate limit, and should not be able to start another operation.
	err = limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)

	// Advance time by one second.
	// Intentionally do not finish the operation. We are attempting to see what happens when an operation fails
	// due to the limit on parallel operations.
	now = now.Add(time.Second)
	err = limiter.BeginGetChunkOperation(now, userID)
	require.NoError(t, err)

	// This operation will fail due to the concurrency limit. It should not affect the rate limit.
	err = limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)

	// Finish the operation that was started in the previous second. This should permit the next operation to start.
	limiter.FinishGetChunkOperation(userID)

	// Verify that we have the expected number of available tokens.
	for i := 0; i < int(config.MaxGetChunkOpsPerSecond)-1; i++ {
		err = limiter.BeginGetChunkOperation(now, userID)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID)
	}

	// We are now at the rate limit, and should not be able to start another operation.
	err = limiter.BeginGetChunkOperation(now, userID)
	require.Error(t, err)
}

func TestGetChunksBandwidthLimit(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetChunkBytesPerSecond = float64(1024 + rand.Intn(1024*1024))
	config.GetChunkBytesBurstiness = int(config.MaxGetBlobBytesPerSecond) + rand.Intn(1024*1024)
	config.GetChunkBytesBurstinessClient = math.MaxInt32

	userID := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// "register" the user ID
	err := limiter.BeginGetChunkOperation(now, userID)
	require.NoError(t, err)
	limiter.FinishGetChunkOperation(userID)

	// Without advancing time, we should be able to utilize a number of bytes equal to the burstiness limit.
	bytesRemaining := config.GetChunkBytesBurstiness
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetChunkBandwidth(now, userID, bytesToRequest)
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err = limiter.RequestGetChunkBandwidth(now, userID, 1)
	require.Error(t, err)

	// Advance time by one second. We should gain a number of tokens equal to the rate limit.
	now = now.Add(time.Second)
	bytesRemaining = int(config.MaxGetChunkBytesPerSecond)
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetChunkBandwidth(now, userID, bytesToRequest)
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err = limiter.RequestGetChunkBandwidth(now, userID, 1)
	require.Error(t, err)
}

func TestPerClientConcurrencyLimit(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxConcurrentGetChunkOpsClient = 1 + rand.Intn(10)
	config.MaxConcurrentGetChunkOps = 2 * config.MaxConcurrentGetChunkOpsClient
	config.GetChunkOpsBurstinessClient = math.MaxInt32
	config.GetChunkOpsBurstiness = math.MaxInt32

	userID1 := tu.RandomString(64)
	userID2 := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// Start the maximum permitted number of operations for user 1
	for i := 0; i < config.MaxConcurrentGetChunkOpsClient; i++ {
		err := limiter.BeginGetChunkOperation(now, userID1)
		require.NoError(t, err)
	}

	// Starting another operation for user 1 should fail due to the concurrency limit
	err := limiter.BeginGetChunkOperation(now, userID1)
	require.Error(t, err)

	// The failure to start the operation for client 1 should not use up any of the global concurrency slots.
	// To verify this, allow the maximum number of operations for client 2 to start.
	for i := 0; i < config.MaxConcurrentGetChunkOpsClient; i++ {
		err := limiter.BeginGetChunkOperation(now, userID2)
		require.NoError(t, err)
	}

	// Starting another operation for client 2 should fail due to the concurrency limit
	err = limiter.BeginGetChunkOperation(now, userID2)
	require.Error(t, err)

	// Ending an operation from client 2 should not affect the concurrency limit for client 1.
	limiter.FinishGetChunkOperation(userID2)
	err = limiter.BeginGetChunkOperation(now, userID1)
	require.Error(t, err)

	// Ending an operation from client 1 should permit another operation for client 1 to start.
	limiter.FinishGetChunkOperation(userID1)
	err = limiter.BeginGetChunkOperation(now, userID1)
	require.NoError(t, err)
}

func TestOpLimitPerClient(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetChunkOpsPerSecondClient = float64(2 + rand.Intn(10))
	config.GetChunkOpsBurstinessClient = int(config.MaxGetChunkOpsPerSecondClient) + rand.Intn(10)
	config.GetChunkOpsBurstiness = math.MaxInt32

	userID1 := tu.RandomString(64)
	userID2 := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// Without advancing time, we should be able to perform a number of operations equal to the burstiness limit.
	for i := 0; i < config.GetChunkOpsBurstinessClient; i++ {
		err := limiter.BeginGetChunkOperation(now, userID1)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID1)
	}

	// We are not at the rate limit, and should be able to start another operation.
	err := limiter.BeginGetChunkOperation(now, userID1)
	require.Error(t, err)

	// Client 2 should not be rate limited based on actions by client 1.
	for i := 0; i < config.GetChunkOpsBurstinessClient; i++ {
		err := limiter.BeginGetChunkOperation(now, userID2)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID2)
	}

	// Client 2 should now have exhausted its burstiness limit.
	err = limiter.BeginGetChunkOperation(now, userID2)
	require.Error(t, err)

	// Advancing time by a second should permit more operations.
	now = now.Add(time.Second)
	for i := 0; i < int(config.MaxGetChunkOpsPerSecondClient); i++ {
		err = limiter.BeginGetChunkOperation(now, userID1)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID1)
		err = limiter.BeginGetChunkOperation(now, userID2)
		require.NoError(t, err)
		limiter.FinishGetChunkOperation(userID2)
	}

	// No more operations should be permitted for either client.
	err = limiter.BeginGetChunkOperation(now, userID1)
	require.Error(t, err)
	err = limiter.BeginGetChunkOperation(now, userID2)
	require.Error(t, err)
}

func TestBandwidthLimitPerClient(t *testing.T) {
	tu.InitializeRandom()

	config := defaultConfig()
	config.MaxGetChunkBytesPerSecondClient = float64(1024 + rand.Intn(1024*1024))
	config.GetChunkBytesBurstinessClient = int(config.MaxGetBlobBytesPerSecond) + rand.Intn(1024*1024)
	config.GetChunkBytesBurstiness = math.MaxInt32
	config.GetChunkOpsBurstiness = math.MaxInt32
	config.GetChunkOpsBurstinessClient = math.MaxInt32

	userID1 := tu.RandomString(64)
	userID2 := tu.RandomString(64)

	limiter := NewChunkRateLimiter(config, nil)

	// time starts at current time, but advances manually afterward
	now := time.Now()

	// "register" the user IDs
	err := limiter.BeginGetChunkOperation(now, userID1)
	require.NoError(t, err)
	limiter.FinishGetChunkOperation(userID1)
	err = limiter.BeginGetChunkOperation(now, userID2)
	require.NoError(t, err)
	limiter.FinishGetChunkOperation(userID2)

	// Request maximum possible bandwidth for client 1
	bytesRemaining := config.GetChunkBytesBurstinessClient
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetChunkBandwidth(now, userID1, bytesToRequest)
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err = limiter.RequestGetChunkBandwidth(now, userID1, 1)
	require.Error(t, err)

	// User 2 should have its full bandwidth allowance available
	bytesRemaining = config.GetChunkBytesBurstinessClient
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetChunkBandwidth(now, userID2, bytesToRequest)
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// Requesting one more byte should fail due to the bandwidth limit
	err = limiter.RequestGetChunkBandwidth(now, userID2, 1)
	require.Error(t, err)

	// Advance time by one second. We should gain a number of tokens equal to the rate limit.
	now = now.Add(time.Second)
	bytesRemaining = int(config.MaxGetChunkBytesPerSecondClient)
	for bytesRemaining > 0 {
		bytesToRequest := 1 + rand.Intn(bytesRemaining)
		err = limiter.RequestGetChunkBandwidth(now, userID1, bytesToRequest)
		require.NoError(t, err)
		err = limiter.RequestGetChunkBandwidth(now, userID2, bytesToRequest)
		require.NoError(t, err)
		bytesRemaining -= bytesToRequest
	}

	// All bandwidth should now be exhausted for both clients
	err = limiter.RequestGetChunkBandwidth(now, userID1, 1)
	require.Error(t, err)
	err = limiter.RequestGetChunkBandwidth(now, userID2, 1)
	require.Error(t, err)
}
