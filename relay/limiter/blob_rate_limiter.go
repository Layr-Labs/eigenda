package limiter

import (
	"fmt"
	"golang.org/x/time/rate"
	"golang.org/x/tools/container/intsets"
	"sync/atomic"
	"time"
)

// TODO test

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
	operationsInFlight atomic.Int64
}

func NewBlobRateLimiter(config *Config) *BlobRateLimiter {
	globalGetBlobOpLimiter := rate.NewLimiter(rate.Limit(config.MaxGetBlobOpsPerSecond), 1)

	// Burst size is set to MaxInt. This is safe, as the requested size is always a size we've
	// determined by reading the blob metadata, which is guaranteed to respect maximum blob size.
	globalGetBlobBandwidthLimiter := rate.NewLimiter(rate.Limit(config.MaxGetBlobBytesPerSecond), intsets.MaxInt)

	return &BlobRateLimiter{
		config:           config,
		opLimiter:        globalGetBlobOpLimiter,
		bandwidthLimiter: globalGetBlobBandwidthLimiter,
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

	countInFlight := l.operationsInFlight.Add(1)
	if countInFlight > int64(l.config.MaxConcurrentGetBlobOps) {
		l.operationsInFlight.Add(-1)
		return fmt.Errorf("global concurrent request limit exceeded for getBlob operations, try again later")
	}

	allowed := l.opLimiter.AllowN(now, 1)

	if !allowed {
		l.operationsInFlight.Add(-1)
		return fmt.Errorf("global rate limit exceeded for getBlob operations, try again later")
	}
	return nil
}

// FinishGetBlobOperation should be called exactly once for each time BeginGetBlobOperation is called and
// returns nil.
func (l *BlobRateLimiter) FinishGetBlobOperation() {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return
	}

	l.operationsInFlight.Add(-1)
}

// RequestGetBlobBandwidth should be called when a GetBlob is about to start downloading blob data
// from S3. It returns an error if there is insufficient bandwidth available. If it returns nil, the
// operation should proceed.
func (l *BlobRateLimiter) RequestGetBlobBandwidth(now time.Time, bytes uint32) error {
	if l == nil {
		// If the rate limiter is nil, do not enforce rate limits.
		return nil
	}

	allowed := l.bandwidthLimiter.AllowN(now, int(bytes))
	if !allowed {
		return fmt.Errorf("global rate limit exceeded for getBlob bandwidth, try again later")
	}
	return nil
}
