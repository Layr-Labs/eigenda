package meterer

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"golang.org/x/time/rate"
)

// OnDemandMeterer handles global throughput rate limiting for on-demand payments.
// It ensures that the global maximum throughput is observed across all on-demand dispersals.
//
// This struct is safe for use by multiple goroutines.
type OnDemandMeterer struct {
	limiter       *rate.Limiter
	getNow        func() time.Time
	metrics       *OnDemandMetererMetrics
	minNumSymbols atomic.Uint32
}

// OnDemandMetererConfig configures how the meterer refreshes on-chain limits and applies fuzz.
type OnDemandMetererConfig struct {
	// RefreshInterval controls how often on-chain limits are fetched.
	RefreshInterval time.Duration
	// FuzzFactor is a multiplier applied on top of the on-chain limit to tolerate small drifts.
	FuzzFactor float64
}

const (
	// DefaultOnDemandMetererRefreshInterval is how often the meterer re-reads on-chain limits.
	DefaultOnDemandMetererRefreshInterval = 5 * time.Minute
	// DefaultOnDemandMetererFuzzFactor allows 10% additional throughput above on-chain limit.
	DefaultOnDemandMetererFuzzFactor = 1.1
)

// WithDefaults fills unset fields with sane defaults.
func (c OnDemandMetererConfig) WithDefaults() OnDemandMetererConfig {
	if c.RefreshInterval <= 0 {
		c.RefreshInterval = DefaultOnDemandMetererRefreshInterval
	}
	if c.FuzzFactor == 0 {
		c.FuzzFactor = DefaultOnDemandMetererFuzzFactor
	}
	return c
}

// Verify ensures config values are valid.
func (c OnDemandMetererConfig) Verify() error {
	if c.RefreshInterval <= 0 {
		return fmt.Errorf("refresh interval must be positive")
	}
	if c.FuzzFactor <= 0 {
		return fmt.Errorf("fuzz factor must be positive")
	}
	return nil
}

type limiterParams struct {
	limit         rate.Limit
	burst         int
	minNumSymbols uint32
}

// buildLimiterParams loads current on-chain limits and applies fuzz to produce limiter settings.
func buildLimiterParams(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	fuzzFactor float64,
) (limiterParams, error) {
	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return limiterParams{}, fmt.Errorf("get global symbols per second: %w", err)
	}

	globalRatePeriodInterval, err := paymentVault.GetGlobalRatePeriodInterval(ctx)
	if err != nil {
		return limiterParams{}, fmt.Errorf("get global rate period interval: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return limiterParams{}, fmt.Errorf("get min num symbols: %w", err)
	}

	limit := rate.Limit(float64(globalSymbolsPerSecond) * fuzzFactor)
	burst := int(math.Ceil(float64(globalSymbolsPerSecond*globalRatePeriodInterval) * fuzzFactor))
	if burst < 1 {
		burst = 1
	}

	return limiterParams{
		limit:         limit,
		burst:         burst,
		minNumSymbols: minNumSymbols,
	}, nil
}

// Creates a new OnDemandMeterer with the specified rate limiting parameters.
func NewOnDemandMeterer(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	getNow func() time.Time,
	metrics *OnDemandMetererMetrics,
	config OnDemandMetererConfig,
) (*OnDemandMeterer, error) {
	config = config.WithDefaults()
	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid on-demand meterer config: %w", err)
	}

	params, err := buildLimiterParams(ctx, paymentVault, config.FuzzFactor)
	if err != nil {
		return nil, err
	}

	limiter := rate.NewLimiter(params.limit, params.burst)

	m := &OnDemandMeterer{
		limiter: limiter,
		getNow:  getNow,
		metrics: metrics,
	}
	m.minNumSymbols.Store(params.minNumSymbols)

	go m.refreshLoop(ctx, paymentVault, config)

	return m, nil
}

func (m *OnDemandMeterer) refreshLoop(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	config OnDemandMetererConfig,
) {
	ticker := time.NewTicker(config.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			params, err := buildLimiterParams(ctx, paymentVault, config.FuzzFactor)
			if err != nil {
				// Keep existing values on error; next tick may succeed.
				continue
			}
			m.applyLimiterParams(params)
		}
	}
}

func (m *OnDemandMeterer) applyLimiterParams(params limiterParams) {
	m.limiter.SetLimit(params.limit)
	m.limiter.SetBurst(params.burst)
	m.minNumSymbols.Store(params.minNumSymbols)
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
func (m *OnDemandMeterer) MeterDispersal(symbolCount uint32) (*rate.Reservation, error) {
	now := m.getNow()

	billableSymbols := payments.CalculateBillableSymbols(symbolCount, m.minNumSymbols.Load())
	reservation := m.limiter.ReserveN(now, int(billableSymbols))

	if !reservation.OK() || reservation.DelayFrom(now) > 0 {
		reservation.Cancel()
		m.metrics.RecordGlobalMeterExhaustion(billableSymbols)
		return nil, fmt.Errorf("global rate limit exceeded: cannot reserve %d symbols", billableSymbols)
	}

	m.metrics.RecordGlobalMeterThroughput(billableSymbols)
	return reservation, nil
}

// Cancels a reservation obtained by MeterDispersal, returning tokens to the rate limiter.
// This should be called when a reserved dispersal will not be performed (e.g., payment verification failed).
//
// Input reservation must be non-nil, otherwise this will panic
func (m *OnDemandMeterer) CancelDispersal(reservation *rate.Reservation) {
	reservation.Cancel()
}
