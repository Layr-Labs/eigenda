package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"golang.org/x/time/rate"
)



// Payment errors
var (
	// ErrExceedsMaxWaitTime is returned when the wait time for a debit exceeds the configured maxWaitTime
	ErrExceedsMaxWaitTime = errors.New("wait time exceeds maximum allowed")
	// ErrInsufficientCapacity is returned when there isn't enough capacity for the requested debit
	ErrInsufficientCapacity = errors.New("insufficient reservation capacity")
)

type ReservationLedger struct {
	timeSource func() time.Time

	// A "token bucket" which is used to track reservation usage.
	limiter *rate.Limiter


}

func NewReservationLedger(
	reservation *core.ReservedPayment,
	initBehavior InitBehavior,
	overdraftBehavior OverdraftBehavior,
	maxWaitTime time.Duration,
	// bucketCapacity is how much time worth of reservations should the capacity be
	bucketCapacity time.Duration,
) (*ReservationLedger, error) {
	bucketSize := int(float64(reservation.SymbolsPerSecond) * bucketCapacity.Seconds())
	limiter := rate.NewLimiter(rate.Limit(reservation.SymbolsPerSecond), bucketSize)

	// If init behavior is empty, consume all tokens to start with empty bucket
	if initBehavior == InitEmpty {
		// Use ReserveN to consume all tokens in the bucket
		// This ensures the bucket starts empty even if AllowN would fail
		now := timeSource()

		// There isn't any standard way to start with an empty token bucket. To work around this, we simply take out a
		// reservation for the entire bucket capacity.
		reservation := limiter.ReserveN(now, bucketSize)
		if !reservation.OK() {
			return nil, fmt.Errorf("Drain bucket as prescribed by %v config. This should not be possible", InitEmpty)
		}
	}

	return &ReservationLedger{
		timeSource:        timeSource,
		limiter:           limiter,
		maxWaitTime:       maxWaitTime,
		overdraftBehavior: overdraftBehavior,
	}, nil
}

// // A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
// // A Reservation may be canceled, which may enable the Limiter to permit additional events.
// type Reservation struct {
// 	ok        bool
// 	lim       *Limiter
// 	tokens    int
// 	timeToAct time.Time
// 	// This is the Limit at reservation time, it can change later.
// 	limit Limit
// }

func (rl *ReservationLedger) Debit(
	ctx context.Context,
	symbols uint64,
) error {

	// TODO: make sure symbols doesn't exceed max value of int. also make sure it isn't negative

	now := rl.timeSource()

	reservation := rl.limiter.ReserveN(now, int(symbols))


	
	// // Handle OverdraftNotPermitted case
	// if rl.overdraftBehavior == OverdraftNotPermitted {
	// 	// Try to allow immediately without waiting
	// 	if rl.limiter.AllowN(now, int(symbols)) {
	// 		return nil
	// 	}
		
	// 	// Not enough tokens immediately available, check if we should wait
	// 	reservation := rl.limiter.ReserveN(now, int(symbols))
	// 	if !reservation.OK() {
	// 		return ErrInsufficientCapacity
	// 	}
		
	// 	waitDuration := reservation.DelayFrom(now)
		
	// 	// If wait time exceeds max, cancel and return error
	// 	if waitDuration > rl.maxWaitTime {
	// 		reservation.CancelAt(now)
	// 		return fmt.Errorf("%w: would need to wait %v", ErrExceedsMaxWaitTime, waitDuration)
	// 	}
		
	// 	// Wait for the reservation to be ready
	// 	if waitDuration > 0 {
	// 		timer := time.NewTimer(waitDuration)
	// 		defer timer.Stop()
			
	// 		select {
	// 		case <-ctx.Done():
	// 			reservation.CancelAt(rl.timeSource())
	// 			return ctx.Err()
	// 		case <-timer.C:
	// 			// Reservation is now ready
	// 			return nil
	// 		}
	// 	}
		
	// 	return nil
	// }
	
	// // Handle OverdraftOncePermitted case
	// if rl.overdraftBehavior == OverdraftOncePermitted {
	// 	// Check if there are ANY tokens available (even 1)
	// 	tokens := rl.limiter.Tokens()
		
	// 	if tokens > 0 {
	// 		// We have at least some tokens, allow the overdraft
	// 		reservation := rl.limiter.ReserveN(now, int(symbols))
	// 		if !reservation.OK() {
	// 			return ErrInsufficientCapacity
	// 		}
			
	// 		// The reservation will go into overdraft, but that's allowed in this mode
	// 		// No need to wait since we're overdrafting
	// 		return nil
	// 	}
		
	// 	// No tokens available at all, need to wait
	// 	reservation := rl.limiter.ReserveN(now, int(symbols))
	// 	if !reservation.OK() {
	// 		return ErrInsufficientCapacity
	// 	}
		
	// 	waitDuration := reservation.DelayFrom(now)
		
	// 	// Check if wait time exceeds max
	// 	if waitDuration > rl.maxWaitTime {
	// 		reservation.CancelAt(now)
	// 		return fmt.Errorf("%w: would need to wait %v", ErrExceedsMaxWaitTime, waitDuration)
	// 	}
		
	// 	// Wait for tokens to become available
	// 	if waitDuration > 0 {
	// 		timer := time.NewTimer(waitDuration)
	// 		defer timer.Stop()
			
	// 		select {
	// 		case <-ctx.Done():
	// 			reservation.CancelAt(rl.timeSource())
	// 			return ctx.Err()
	// 		case <-timer.C:
	// 			// Reservation is now ready
	// 			return nil
	// 		}
	// 	}
		
	// 	return nil
	// }
	
	// // Unknown overdraft behavior
	// return fmt.Errorf("unknown overdraft behavior: %v", rl.overdraftBehavior)
}
