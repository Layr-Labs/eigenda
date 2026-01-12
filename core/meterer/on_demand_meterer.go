package meterer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core/payments"
)

// OnDemandMeterer handles global throughput rate limiting for on-demand payments.
// It ensures that the global maximum throughput is observed across all on-demand dispersals.
//
// This struct is safe for use by multiple goroutines.
type OnDemandMeterer struct {
	mu            sync.RWMutex
	bucket        *ratelimit.LeakyBucket
	getNow        func() time.Time
	metrics       *OnDemandMetererMetrics
	minNumSymbols uint32
	paymentVault  payments.PaymentVault
	fuzzFactor    float64

	// cached on-chain params for change detection
	globalSymbolsPerSecond   uint64
	globalRatePeriodInterval uint64
}

type bucketParams struct {
	leakRate     float64
	capacity     time.Duration
	minSymbols   uint32
	rawSymbolsPS uint64
	rawPeriod    uint64
}

// OnDemandReservation captures a bucket fill that can be reverted.
type OnDemandReservation struct {
	quantity float64
}

// Creates a new OnDemandMeterer with the specified rate limiting parameters.
func NewOnDemandMeterer(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	getNow func() time.Time,
	metrics *OnDemandMetererMetrics,
	fuzzFactor float64,
) (*OnDemandMeterer, error) {
	if fuzzFactor <= 0 {
		return nil, fmt.Errorf("fuzz factor must be > 0: got %f", fuzzFactor)
	}

	params, err := buildBucket(ctx, paymentVault, fuzzFactor)
	if err != nil {
		return nil, err
	}

	startTime := getNow()

	bucket, err := ratelimit.NewLeakyBucket(
		params.leakRate,
		params.capacity,
		false, /* start empty so capacity represents available tokens */
		ratelimit.OverfillNotPermitted,
		startTime,
	)
	if err != nil {
		return nil, fmt.Errorf("create leaky bucket: %w", err)
	}

	return &OnDemandMeterer{
		mu:            sync.RWMutex{},
		bucket:        bucket,
		getNow:        getNow,
		metrics:       metrics,
		minNumSymbols: params.minSymbols,
		paymentVault:  paymentVault,
		fuzzFactor:    fuzzFactor,

		globalSymbolsPerSecond:   params.rawSymbolsPS,
		globalRatePeriodInterval: params.rawPeriod,
	}, nil
}

// Reserves tokens for a dispersal with the given number of symbols.
//
// The actual number of tokens reserved is the billable symbols (applying the minNumSymbols threshold),
// not the raw symbol count.
//
// Returns a reservation that can be cancelled if the dispersal is not performed (e.g., if payment verification fails).
// The reservation will automatically take effect if not cancelled.
//
// This method only succeeds if tokens are immediately available (no queueing/waiting). If a reservation is returned,
// it is safe to proceed with dispersal without checking the delay.
func (m *OnDemandMeterer) MeterDispersal(symbolCount uint32) (*OnDemandReservation, error) {
	now := m.getNow()

	m.mu.RLock()
	billableSymbols := payments.CalculateBillableSymbols(symbolCount, m.minNumSymbols)
	ok, err := m.bucket.Fill(now, float64(billableSymbols))
	m.mu.RUnlock()

	if err != nil {
		return nil, fmt.Errorf("fill leaky bucket: %w", err)
	}

	if !ok {
		m.metrics.RecordGlobalMeterExhaustion(billableSymbols)
		return nil, fmt.Errorf("global rate limit exceeded: cannot reserve %d symbols", billableSymbols)
	}

	m.metrics.RecordGlobalMeterThroughput(billableSymbols)
	return &OnDemandReservation{quantity: float64(billableSymbols)}, nil
}

// Cancels a reservation obtained by MeterDispersal, returning tokens to the rate limiter.
// This should be called when a reserved dispersal will not be performed (e.g., payment verification failed).
//
// Input reservation must be non-nil, otherwise this will panic
func (m *OnDemandMeterer) CancelDispersal(reservation *OnDemandReservation) {
	if reservation == nil {
		return
	}

	now := m.getNow()

	m.mu.Lock()
	_ = m.bucket.RevertFill(now, reservation.quantity)
	m.mu.Unlock()
}

// Refresh updates the limiter parameters from the PaymentVault to track any on-chain changes.
func (m *OnDemandMeterer) Refresh(ctx context.Context) error {
	params, err := buildBucket(ctx, m.paymentVault, m.fuzzFactor)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if params.rawSymbolsPS == m.globalSymbolsPerSecond &&
		params.rawPeriod == m.globalRatePeriodInterval &&
		params.minSymbols == m.minNumSymbols {
		return nil
	}

	if err := m.bucket.Reconfigure(
		params.leakRate,
		params.capacity,
		ratelimit.OverfillNotPermitted,
		m.getNow(),
	); err != nil {
		return fmt.Errorf("reconfigure leaky bucket: %w", err)
	}
	m.minNumSymbols = params.minSymbols
	m.globalSymbolsPerSecond = params.rawSymbolsPS
	m.globalRatePeriodInterval = params.rawPeriod
	return nil
}

func buildBucket(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	fuzzFactor float64,
) (*bucketParams, error) {
	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global symbols per second: %w", err)
	}

	globalRatePeriodInterval, err := paymentVault.GetGlobalRatePeriodInterval(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global rate period interval: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	effectiveSymbolsPerSecond := float64(globalSymbolsPerSecond) * fuzzFactor
	if effectiveSymbolsPerSecond < 1 {
		effectiveSymbolsPerSecond = 1
	}

	capacityDuration := time.Duration(globalRatePeriodInterval) * time.Second
	return &bucketParams{
		leakRate:     effectiveSymbolsPerSecond,
		capacity:     capacityDuration,
		minSymbols:   minNumSymbols,
		rawSymbolsPS: globalSymbolsPerSecond,
		rawPeriod:    globalRatePeriodInterval,
	}, nil
}
