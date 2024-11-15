package limiter

import (
	"fmt"
	"golang.org/x/time/rate"
	"golang.org/x/tools/container/intsets"
	"sync/atomic"
	"time"
)

// TODO test

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
	globalOperationsInFlight atomic.Int64

	// per-client limiters

	// Note: in its current form, these expose a DOS vector, since an attacker can create many clients IDs
	// and force these maps to become arbitrarily large. This will be remedied when authentication
	// is implemented, as only authentication will happen prior to rate limiting.

	// perClientOpLimiter enforces per-client rate limits on the maximum rate of GetChunk operations
	perClientOpLimiter map[string]*rate.Limiter

	// perClientBandwidthLimiter enforces per-client rate limits on the maximum bandwidth consumed by
	// GetChunk operations.
	perClientBandwidthLimiter map[string]*rate.Limiter

	// perClientOperationsInFlight is the number of GetChunk operations currently in flight for each client.
	perClientOperationsInFlight map[string]*atomic.Int64
}

func NewChunkRateLimiter(config *Config) *ChunkRateLimiter {

	globalOpLimiter := rate.NewLimiter(rate.Limit(config.MaxGetChunkOpsPerSecond), 1)
	globalBandwidthLimiter := rate.NewLimiter(rate.Limit(config.MaxGetChunkBytesPerSecond), intsets.MaxInt)

	return &ChunkRateLimiter{
		config:                      config,
		globalOpLimiter:             globalOpLimiter,
		globalBandwidthLimiter:      globalBandwidthLimiter,
		globalOperationsInFlight:    atomic.Int64{},
		perClientOpLimiter:          make(map[string]*rate.Limiter),
		perClientBandwidthLimiter:   make(map[string]*rate.Limiter),
		perClientOperationsInFlight: make(map[string]*atomic.Int64),
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

	countInFlight := l.globalOperationsInFlight.Add(1)
	if countInFlight > int64(l.config.MaxConcurrentGetChunkOps) {
		l.globalOperationsInFlight.Add(-1)
		return fmt.Errorf("global concurrent request limit exceeded for GetChunks operations, try again later")
	}

	allowed := l.globalOpLimiter.AllowN(now, 1)
	if !allowed {
		l.globalOperationsInFlight.Add(-1)
		return fmt.Errorf("global rate limit exceeded for GetChunks operations, try again later")
	}

	clientInFlightCounter, ok := l.perClientOperationsInFlight[requesterID]
	if !ok {
		// This is the first time we've seen this client ID.

		l.perClientOperationsInFlight[requesterID] = &atomic.Int64{}
		clientInFlightCounter = l.perClientOperationsInFlight[requesterID]

		l.perClientBandwidthLimiter[requesterID] = rate.NewLimiter(
			rate.Limit(l.config.MaxGetChunkBytesPerSecond), intsets.MaxInt)
	}

	countInFlight = clientInFlightCounter.Add(1)
	if countInFlight > int64(l.config.MaxConcurrentGetChunkOpsClient) {
		l.globalOperationsInFlight.Add(-1)
		clientInFlightCounter.Add(-1)
		return fmt.Errorf("client concurrent request limit exceeded for GetChunks")
	}

	allowed = l.perClientOpLimiter[requesterID].AllowN(now, 1)
	if !allowed {
		l.globalOperationsInFlight.Add(-1)
		clientInFlightCounter.Add(-1)
		return fmt.Errorf("client rate limit exceeded for GetChunks, try again later")
	}

	return nil
}

// FinishGetChunkOperation should be called when a GetChunk operation completes.
func (l *ChunkRateLimiter) FinishGetChunkOperation(requesterID string) {
	if l == nil {
		return
	}

	l.globalOperationsInFlight.Add(-1)
	l.perClientOperationsInFlight[requesterID].Add(-1)
}

// RequestGetChunkBandwidth should be called when a GetChunk is about to start downloading chunk data.
func (l *ChunkRateLimiter) RequestGetChunkBandwidth(now time.Time, requesterID string, bytes int) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

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
