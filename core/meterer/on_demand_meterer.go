package meterer

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// OnDemandMeterer handles global throughput rate limiting for on-demand payments.
// It ensures that the global maximum throughput is observed across all on-demand dispersals.
//
// This struct is safe for use by multiple goroutines.
type OnDemandMeterer struct {
	limiter *rate.Limiter
	getNow  func() time.Time
}

// Creates a new OnDemandMeterer with the specified rate limiting parameters.
func NewOnDemandMeterer(
	globalSymbolsPerSecond uint64,
	globalRatePeriodInterval uint64,
	getNow func() time.Time,
) *OnDemandMeterer {
	burstSize := int(globalSymbolsPerSecond * globalRatePeriodInterval)
	limiter := rate.NewLimiter(rate.Limit(globalSymbolsPerSecond), burstSize)

	return &OnDemandMeterer{
		limiter: limiter,
		getNow:  getNow,
	}
}

// Reserves tokens for a dispersal with the given number of symbols.
//
// Returns a reservation that can be cancelled if the dispersal is not performed (e.g., if payment verification fails).
// The reservation will automatically take effect if not cancelled.
//
// This method only succeeds if tokens are immediately available (no queueing/waiting). If a reservation is returned,
// it is safe to proceed with dispersal without checking the delay.
func (m *OnDemandMeterer) MeterDispersal(symbolCount uint32) (*rate.Reservation, error) {
	reservation := m.limiter.ReserveN(m.getNow(), int(symbolCount))

	if !reservation.OK() || reservation.Delay() > 0 {
		reservation.Cancel()
		return nil, fmt.Errorf("global rate limit exceeded: cannot reserve %d symbols", symbolCount)
	}

	return reservation, nil
}

// Cancels a reservation obtained by MeterDispersal, returning tokens to the rate limiter.
// This should be called when a reserved dispersal will not be performed (e.g., payment verification failed).
func (m *OnDemandMeterer) CancelDispersal(reservation *rate.Reservation) {
	if reservation != nil {
		reservation.Cancel()
	}
}
