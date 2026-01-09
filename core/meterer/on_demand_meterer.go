package meterer

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"golang.org/x/time/rate"
)

// OnDemandMeterer handles global throughput rate limiting for on-demand payments.
// It ensures that the global maximum throughput is observed across all on-demand dispersals.
//
// This struct is safe for use by multiple goroutines.
type OnDemandMeterer struct {
	mu            sync.RWMutex
	limiter       *rate.Limiter
	getNow        func() time.Time
	metrics       *OnDemandMetererMetrics
	minNumSymbols uint32
	paymentVault  payments.PaymentVault
	fuzzFactor    float64
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

	limiter, minNumSymbols, err := buildLimiter(ctx, paymentVault, fuzzFactor)
	if err != nil {
		return nil, err
	}

	return &OnDemandMeterer{
		mu:            sync.RWMutex{},
		limiter:       limiter,
		getNow:        getNow,
		metrics:       metrics,
		minNumSymbols: minNumSymbols,
		paymentVault:  paymentVault,
		fuzzFactor:    fuzzFactor,
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
func (m *OnDemandMeterer) MeterDispersal(symbolCount uint32) (*rate.Reservation, error) {
	now := m.getNow()

	m.mu.RLock()
	billableSymbols := payments.CalculateBillableSymbols(symbolCount, m.minNumSymbols)
	reservation := m.limiter.ReserveN(now, int(billableSymbols))
	m.mu.RUnlock()

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

// Refresh updates the limiter parameters from the PaymentVault to track any on-chain changes.
func (m *OnDemandMeterer) Refresh(ctx context.Context) error {
	limiter, minNumSymbols, err := buildLimiter(ctx, m.paymentVault, m.fuzzFactor)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.limiter = limiter
	m.minNumSymbols = minNumSymbols
	return nil
}

func buildLimiter(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	fuzzFactor float64,
) (*rate.Limiter, uint32, error) {
	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get global symbols per second: %w", err)
	}

	globalRatePeriodInterval, err := paymentVault.GetGlobalRatePeriodInterval(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get global rate period interval: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get min num symbols: %w", err)
	}

	effectiveSymbolsPerSecond := float64(globalSymbolsPerSecond) * fuzzFactor
	if effectiveSymbolsPerSecond < 1 {
		effectiveSymbolsPerSecond = 1
	}

	burstSize := int(math.Ceil(effectiveSymbolsPerSecond * float64(globalRatePeriodInterval)))
	if burstSize < 1 {
		burstSize = 1
	}

	limiter := rate.NewLimiter(rate.Limit(effectiveSymbolsPerSecond), burstSize)
	return limiter, minNumSymbols, nil
}
