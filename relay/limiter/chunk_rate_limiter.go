package limiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/relay/metrics"
	"golang.org/x/time/rate"
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

	// perClientOpLimiter enforces per-client rate limits on the maximum rate of GetChunk operations
	perClientOpLimiter map[string]*rate.Limiter

	// perClientBandwidthLimiter enforces per-client rate limits on the maximum bandwidth consumed by
	// GetChunk operations.
	perClientBandwidthLimiter map[string]*rate.Limiter

	// perClientOperationsInFlight is the number of GetChunk operations currently in flight for each client.
	perClientOperationsInFlight map[string]int

	// Encapsulates relay metrics.
	relayMetrics *metrics.RelayMetrics

	// this lock is used to provide thread safety
	lock sync.Mutex
}

// NewChunkRateLimiter creates a new ChunkRateLimiter.
func NewChunkRateLimiter(
	config *Config,
	relayMetrics *metrics.RelayMetrics) *ChunkRateLimiter {

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
		relayMetrics:                relayMetrics,
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
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("global concurrency")
		}
		return fmt.Errorf(
			"global concurrent request limit %d exceeded for GetChunks operations, try again later",
			l.config.MaxConcurrentGetChunkOps)
	}
	if l.globalOpLimiter.TokensAt(now) < 1 {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("global rate")
		}
		return fmt.Errorf("global rate limit %0.1fhz exceeded for GetChunks operations, try again later",
			l.config.MaxGetChunkOpsPerSecond)
	}
	if l.perClientOperationsInFlight[requesterID] >= l.config.MaxConcurrentGetChunkOpsClient {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("client concurrency")
		}
		return fmt.Errorf("client concurrent request limit %d exceeded for GetChunks",
			l.config.MaxConcurrentGetChunkOpsClient)
	}
	if l.perClientOpLimiter[requesterID].TokensAt(now) < 1 {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("client rate")
		}
		return fmt.Errorf("client rate limit %0.1fhz exceeded for GetChunks, try again later",
			l.config.MaxGetChunkOpsPerSecondClient)
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
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("global bandwidth")
		}
		return fmt.Errorf("global rate limit %dMiB exceeded for GetChunk bandwidth, try again later",
			int(l.config.MaxGetChunkBytesPerSecond/1024/1024))
	}

	limiter, ok := l.perClientBandwidthLimiter[requesterID]
	if !ok {
		return fmt.Errorf("internal error, unable to find bandwidth limiter for client ID %s", requesterID)
	}
	allowed = limiter.AllowN(now, bytes)
	if !allowed {
		l.globalBandwidthLimiter.AllowN(now, -bytes)
		if l.relayMetrics != nil {
			l.relayMetrics.ReportChunkRateLimited("client bandwidth")
		}
		return fmt.Errorf("client rate limit %dMiB exceeded for GetChunk bandwidth, try again later",
			int(l.config.MaxGetChunkBytesPerSecondClient/1024/1024))
	}

	return nil
}
