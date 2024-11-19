package limiter

import (
	"fmt"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// ChunkRateLimiter enforces rate limits on GetChunk operations.
type ChunkRateLimiter struct {

	// config is the rate limit configuration.
	config *Config

	// global limiters

	// globalOpLimiter enforces global rate limits on the maximum rate of GetChunk operations
	globalOpLimiter *rate.Limiter

	// globalBandwidthLimiter enforces global rate limits on the maximum bandwidth consumed by GetChunk operations.
	globalBandwidthLimiter *rate.Limiter

	// globalOperationsInFlight is the number of GetChunk operations currently in flight.
	globalOperationsInFlight int

	// per-client limiters

	// Note: in its current form, these expose a DOS vector, since an attacker can create many clients IDs
	// and force these maps to become arbitrarily large. This will be remedied when auth
	// is implemented, as only auth will happen prior to rate limiting.

	// perClientOpLimiter enforces per-client rate limits on the maximum rate of GetChunk operations
	perClientOpLimiter map[string]*rate.Limiter

	// perClientBandwidthLimiter enforces per-client rate limits on the maximum bandwidth consumed by
	// GetChunk operations.
	perClientBandwidthLimiter map[string]*rate.Limiter

	// perClientOperationsInFlight is the number of GetChunk operations currently in flight for each client.
	perClientOperationsInFlight map[string]int

	// this lock is used to provide thread safety
	lock sync.Mutex
}

// NewChunkRateLimiter creates a new ChunkRateLimiter.
func NewChunkRateLimiter(config *Config) *ChunkRateLimiter {

	globalOpLimiter := rate.NewLimiter(rate.Limit(
		config.MaxGetChunkOpsPerSecond),
		config.GetChunkOpsBurstiness)

	globalBandwidthLimiter := rate.NewLimiter(rate.Limit(
		config.MaxGetChunkBytesPerSecond),
		config.GetChunkBytesBurstiness)

	return &ChunkRateLimiter{
		config:                      config,
		globalOpLimiter:             globalOpLimiter,
		globalBandwidthLimiter:      globalBandwidthLimiter,
		perClientOpLimiter:          make(map[string]*rate.Limiter),
		perClientBandwidthLimiter:   make(map[string]*rate.Limiter),
		perClientOperationsInFlight: make(map[string]int),
	}
}

// BeginGetChunkOperation should be called when a GetChunk operation is about to begin. If it returns an error,
// the operation should not be performed. If it does not return an error, FinishGetChunkOperation should be
// called when the operation completes.
func (l *ChunkRateLimiter) BeginGetChunkOperation(
	now time.Time,
	requesterID string) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.perClientOperationsInFlight[requesterID]
	if !ok {
		// This is the first time we've seen this client ID.
		l.perClientOperationsInFlight[requesterID] = 0

		l.perClientOpLimiter[requesterID] = rate.NewLimiter(
			rate.Limit(l.config.MaxGetChunkOpsPerSecondClient),
			l.config.GetChunkOpsBurstinessClient)

		l.perClientBandwidthLimiter[requesterID] = rate.NewLimiter(
			rate.Limit(l.config.MaxGetChunkBytesPerSecondClient),
			l.config.GetChunkBytesBurstinessClient)
	}

	if l.globalOperationsInFlight >= l.config.MaxConcurrentGetChunkOps {
		return fmt.Errorf("global concurrent request limit exceeded for GetChunks operations, try again later")
	}
	if l.globalOpLimiter.TokensAt(now) < 1 {
		return fmt.Errorf("global rate limit exceeded for GetChunks operations, try again later")
	}
	if l.perClientOperationsInFlight[requesterID] >= l.config.MaxConcurrentGetChunkOpsClient {
		return fmt.Errorf("client concurrent request limit exceeded for GetChunks")
	}
	if l.perClientOpLimiter[requesterID].TokensAt(now) < 1 {
		return fmt.Errorf("client rate limit exceeded for GetChunks, try again later")
	}

	l.globalOperationsInFlight++
	l.perClientOperationsInFlight[requesterID]++
	l.globalOpLimiter.AllowN(now, 1)
	l.perClientOpLimiter[requesterID].AllowN(now, 1)

	return nil
}

// FinishGetChunkOperation should be called when a GetChunk operation completes.
func (l *ChunkRateLimiter) FinishGetChunkOperation(requesterID string) {
	if l == nil {
		return
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	l.globalOperationsInFlight--
	l.perClientOperationsInFlight[requesterID]--
}

// RequestGetChunkBandwidth should be called when a GetChunk is about to start downloading chunk data.
func (l *ChunkRateLimiter) RequestGetChunkBandwidth(now time.Time, requesterID string, bytes int) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

	// no lock needed here, as the bandwidth limiters themselves are thread-safe

	allowed := l.globalBandwidthLimiter.AllowN(now, bytes)
	if !allowed {
		return fmt.Errorf("global rate limit exceeded for GetChunk bandwidth, try again later")
	}

	allowed = l.perClientBandwidthLimiter[requesterID].AllowN(now, bytes)
	if !allowed {
		l.globalBandwidthLimiter.AllowN(now, -bytes)
		return fmt.Errorf("client rate limit exceeded for GetChunk bandwidth, try again later")
	}

	return nil
}
