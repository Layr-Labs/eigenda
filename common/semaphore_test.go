package common_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TODO
// - mutual exclusion
// - ordering
// - context is cancelled
// - semaphore is closed
// - semaphore with 0 capacity
// - asking for more tokens than available
// - releasing more tokens than acquired

func TestZeroTokens(t *testing.T) {
	_, err := common.NewSemaphore(0)
	require.Error(t, err)
}

func TestRequestTooManyTokens(t *testing.T) {
	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// We should be able to acquire the exact number of tokens available.
	err = semaphore.Acquire(t.Context(), tokenCount)
	require.NoError(t, err)
	err = semaphore.Release(tokenCount)
	require.NoError(t, err)

	// We should not be able to acquire more tokens than available.
	err = semaphore.Acquire(t.Context(), tokenCount+1)
	require.Error(t, err)
}

func TestExclusion(t *testing.T) {
	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// Acquire more than half of all tokens.
	err = semaphore.Acquire(t.Context(), tokenCount/2+1)
	require.NoError(t, err)

	// Create a goroutine that tries to acquire the same number of tokens.
	// The semaphore should prevent this.
	acquiredTokens := atomic.Bool{}
	acquiredTokensChan := make(chan struct{})
	go func() {
		err := semaphore.Acquire(t.Context(), tokenCount/2+1)
		require.NoError(t, err)
		acquiredTokens.Store(true)
		acquiredTokensChan <- struct{}{}
	}()

	// Wait for a little while. Goroutine should not be able to acquire the tokens.
	time.Sleep(50 * time.Millisecond)
	require.False(t, acquiredTokens.Load())

	// Release the tokens we acquired.
	err = semaphore.Release(tokenCount/2 + 1)
	require.NoError(t, err)

	// Now the goroutine should be able to acquire the tokens.
	select {
	case <-acquiredTokensChan:
		require.True(t, acquiredTokens.Load())
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "goroutine did not acquire tokens in time")
	}

	// Release the tokens acquired by the goroutine.
	err = semaphore.Release(tokenCount/2 + 1)
	require.NoError(t, err)

}
