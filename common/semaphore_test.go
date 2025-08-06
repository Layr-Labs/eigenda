package common_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func TestZeroTokens(t *testing.T) {
	t.Parallel()

	_, err := common.NewSemaphore(0)
	require.Error(t, err)
}

func TestRequestTooManyTokens(t *testing.T) {
	t.Parallel()

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

func TestReleaseTooManyTokens(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// release more tokens than the maximum capacity of the semaphore
	err = semaphore.Release(tokenCount + 1)
	require.Error(t, err)

	// release a legal number of tokens, but tokens that were never acquired
	err = semaphore.Release(tokenCount / 2)
	require.Error(t, err)
}

func TestExclusion(t *testing.T) {
	t.Parallel()

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
	time.Sleep(10 * time.Millisecond)
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

// Similar to TestExclusion, but with multiple releases prior to the goroutine acquiring tokens.
func TestMultiExclusion(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// Acquire more than half of all tokens.
	err = semaphore.Acquire(t.Context(), tokenCount)
	require.NoError(t, err)

	// Create a goroutine that tries to acquire the same number of tokens.
	// The semaphore should prevent this.
	acquiredTokens := atomic.Bool{}
	acquiredTokensChan := make(chan struct{})
	go func() {
		err := semaphore.Acquire(t.Context(), tokenCount)
		require.NoError(t, err)
		acquiredTokens.Store(true)
		acquiredTokensChan <- struct{}{}
	}()

	// Release tokens gradually.
	tokensToRelease := tokenCount
	for tokensToRelease > 0 {
		time.Sleep(time.Millisecond)
		require.False(t, acquiredTokens.Load())

		releaseCount := tokensToRelease / 2
		if releaseCount == 0 {
			// the last release will round 0.5 down to 0, so we need to ensure at least 1 token is released
			releaseCount = 1
		}
		tokensToRelease -= releaseCount
		err = semaphore.Release(releaseCount)
		require.NoError(t, err)
	}

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

// We want to make sure that the closing of the semaphore does not cause deadlocks.
func TestClosingSemaphore(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)

	err = semaphore.Acquire(t.Context(), tokenCount)
	require.NoError(t, err)

	// Create a goroutine that tries to acquire the same number of tokens.
	// The semaphore should prevent this. As soon as the semaphore is closed, the goroutine should return an error.
	unblocked := atomic.Bool{}
	unblockedChan := make(chan struct{})
	go func() {
		err := semaphore.Acquire(t.Context(), tokenCount)
		require.Error(t, err)
		unblocked.Store(true)
		unblockedChan <- struct{}{}
	}()

	// Wait for a little while. Goroutine should still be blocked.
	time.Sleep(10 * time.Millisecond)
	require.False(t, unblocked.Load())

	// Close the semaphore. Should unblock the goroutine.
	semaphore.Close()

	select {
	case <-unblockedChan:
		require.True(t, unblocked.Load())
	case <-time.After(10000 * time.Millisecond):
		require.Fail(t, "goroutine did not unblock in time")
	}

	// Closing multiple times shouldn't cause any issues.
	semaphore.Close()

	// Additional calls to Acquire and Release after closing should return errors.
	err = semaphore.Acquire(t.Context(), 1)
	require.Error(t, err)
	err = semaphore.Release(0)
	require.Error(t, err)
}

// See what happens when the context is cancelled prior to the Acquire call.
func TestContextImmediatelyCancelled(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// Acquire more than half of all tokens.
	err = semaphore.Acquire(t.Context(), tokenCount/2+1)
	require.NoError(t, err)

	// Create a context that is already canceled.
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	// This should return an error immediately without blocking.
	err = semaphore.Acquire(ctx, tokenCount/2+1)
	require.Error(t, err)
}

// See what happens when the context is cancelled after the Acquire call, but before the request is put on the channel
// to the control loop.
func TestContextCancelledBeforeSending(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	// Acquire all tokens.
	err = semaphore.Acquire(t.Context(), tokenCount)
	require.NoError(t, err)

	// Fill up the channel. After this operation, Acquire() should get blocked on insertion into the channel.
	// The channel has capacity 64, and the control loop will buffer one request. So we need to add 65 requests.
	for i := 0; i < 65; i++ {
		go func() {
			// this will block until we close the semaphore
			_ = semaphore.Acquire(t.Context(), 1)
		}()
	}

	// Wait a little while to give the goroutines time to block. This test won't fail if not all goroutines block,
	// but we may not exercise the code path we want to test.
	time.Sleep(50 * time.Millisecond)

	// Create a context that we will eventually cancel.
	ctx, cancel := context.WithCancel(t.Context())

	// Submit a request with that context.
	unblocked := atomic.Bool{}
	unblockedChan := make(chan struct{})
	go func() {
		err := semaphore.Acquire(ctx, 1)
		require.Error(t, err)
		unblocked.Store(true)
		unblockedChan <- struct{}{}
	}()

	// Wait a little while. Goroutine should still be blocked.
	time.Sleep(10 * time.Millisecond)
	require.False(t, unblocked.Load())

	// Cancel the context. This should unblock the goroutine and return an error.
	cancel()

	select {
	case <-unblockedChan:
		require.True(t, unblocked.Load())
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "goroutine did not unblock in time")
	}
}

// See what happens when the context is cancelled after the request is sent to the control loop, but before the
// control loop responds to the request.
func TestContextCancelledBeforeResponse(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	tokenCount := rand.Uint64()/2 + 1

	semaphore, err := common.NewSemaphore(tokenCount)
	require.NoError(t, err)
	defer semaphore.Close()

	err = semaphore.Acquire(t.Context(), tokenCount)
	require.NoError(t, err)

	// Create a context that we will eventually cancel.
	ctx, cancel := context.WithCancel(t.Context())

	// Create a goroutine that tries to acquire the same number of tokens.
	// This should block until the context is cancelled.
	unblocked := atomic.Bool{}
	unblockedChan := make(chan struct{})
	go func() {
		err := semaphore.Acquire(ctx, tokenCount)
		require.Error(t, err)
		unblocked.Store(true)
		unblockedChan <- struct{}{}
	}()

	// Wait a little while. Goroutine should still be blocked. We want the request to be sent to the control loop
	// during this time. If that doesn't happen, this test will still pass, but we won't exercise the code path we
	// want to test.
	time.Sleep(50 * time.Millisecond)
	require.False(t, unblocked.Load())

	// Cancel the context. This should unblock the goroutine.
	cancel()

	select {
	case <-unblockedChan:
		require.True(t, unblocked.Load())
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "goroutine did not unblock in time")
	}
}
