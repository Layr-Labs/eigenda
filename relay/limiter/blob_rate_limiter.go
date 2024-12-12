package limiter

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/relay/metrics"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// BlobRateLimiter enforces rate limits on GetBlob operations.
type BlobRateLimiter struct {

	// config is the rate limit configuration.
	config *Config

	// opLimiter enforces rate limits on the maximum rate of GetBlob operations
	opLimiter *rate.Limiter

	// bandwidthLimiter enforces rate limits on the maximum bandwidth consumed by GetBlob operations. Only the size
	// of the blob data is considered, not the size of the entire response.
	bandwidthLimiter *rate.Limiter

	// operationsInFlight is the number of GetBlob operations currently in flight.
	operationsInFlight int

	// Encapsulates relay metrics.
	relayMetrics *metrics.RelayMetrics

	// this lock is used to provide thread safety
	lock sync.Mutex
}

// NewBlobRateLimiter creates a new BlobRateLimiter.
func NewBlobRateLimiter(config *Config, relayMetrics *metrics.RelayMetrics) *BlobRateLimiter {
	globalGetBlobOpLimiter := rate.NewLimiter(
		rate.Limit(config.MaxGetBlobOpsPerSecond),
		config.GetBlobOpsBurstiness)

	globalGetBlobBandwidthLimiter := rate.NewLimiter(
		rate.Limit(config.MaxGetBlobBytesPerSecond),
		config.GetBlobBytesBurstiness)

	return &BlobRateLimiter{
		config:           config,
		opLimiter:        globalGetBlobOpLimiter,
		bandwidthLimiter: globalGetBlobBandwidthLimiter,
		relayMetrics:     relayMetrics,
	}
}

// BeginGetBlobOperation should be called when a GetBlob operation is about to begin. If it returns an error,
// the operation should not be performed. If it does not return an error, FinishGetBlobOperation should be
// called when the operation completes.
func (l *BlobRateLimiter) BeginGetBlobOperation(now time.Time) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	if l.operationsInFlight >= l.config.MaxConcurrentGetBlobOps {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportBlobRateLimited("global concurrency")
		}
		return fmt.Errorf("global concurrent request limit %d exceeded for getBlob operations, try again later",
			l.config.MaxConcurrentGetBlobOps)
	}
	if l.opLimiter.TokensAt(now) < 1 {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportBlobRateLimited("global rate")
		}
		return fmt.Errorf("global rate limit %0.1fhz exceeded for getBlob operations, try again later",
			l.config.MaxGetBlobOpsPerSecond)
	}

	l.operationsInFlight++
	l.opLimiter.AllowN(now, 1)

	return nil
}

// FinishGetBlobOperation should be called exactly once for each time BeginGetBlobOperation is called and
// returns nil.
func (l *BlobRateLimiter) FinishGetBlobOperation() {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	l.operationsInFlight--
}

// RequestGetBlobBandwidth should be called when a GetBlob is about to start downloading blob data
// from S3. It returns an error if there is insufficient bandwidth available. If it returns nil, the
// operation should proceed.
func (l *BlobRateLimiter) RequestGetBlobBandwidth(now time.Time, bytes uint32) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

	// no locking needed, the only thing we touch here is the bandwidthLimiter, which is inherently thread-safe

	allowed := l.bandwidthLimiter.AllowN(now, int(bytes))
	if !allowed {
		if l.relayMetrics != nil {
			l.relayMetrics.ReportBlobRateLimited("global bandwidth")
		}
		return fmt.Errorf("global rate limit %dMib/s exceeded for getBlob bandwidth, try again later",
			int(l.config.MaxGetBlobBytesPerSecond/1024/1024))
	}
	return nil
}
