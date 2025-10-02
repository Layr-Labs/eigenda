package reservation

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments"
)

// Tracks usage of a single account reservation
//
// This struct is goroutine safe.
type ReservationLedger struct {
	config ReservationLedgerConfig

	// synchronizes access to the underlying leaky bucket algorithm
	lock sync.Mutex

	// an instance of the algorithm which tracks reservation usage
	leakyBucket *ratelimit.LeakyBucket
}

// Creates a new reservation ledger, which represents the reservation of a single user with a [LeakyBucket]
func NewReservationLedger(
	config ReservationLedgerConfig,
	// now should be from a source that includes monotonic timestamp for best results
	now time.Time,
) (*ReservationLedger, error) {
	leakyBucket, err := ratelimit.NewLeakyBucket(
		float64(config.reservation.symbolsPerSecond),
		config.bucketCapacityDuration,
		config.startFull,
		config.overfillBehavior,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("new leaky bucket: %w", err)
	}

	return &ReservationLedger{
		config:      config,
		leakyBucket: leakyBucket,
	}, nil
}

// Debit the reservation with a number of symbols.
//
// Returns (true, remainingCapacity, nil) if the reservation has enough capacity to perform the debit.
// Returns (false, remainingCapacity, nil) if the bucket lacks capacity to permit the fill.
// Returns (false, 0, error) if an error occurs. Possible errors include:
//   - [QuorumNotPermittedError]: one or more of the requested quorums are not permitted by the reservation
//   - [TimeOutOfRangeError]: the dispersal time is outside the reservation's valid time range
//   - [TimeMovedBackwardError]: current time is before a previously observed time (only possible if input time
//     instances don't included monotonic timestamps)
//   - Generic errors for all other unexpected behavior
//
// The remainingCapacity is the amount of space left in the bucket after the operation (in symbols).
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed debit doesn't count against the meter.
func (rl *ReservationLedger) Debit(
	// now should be from a source that includes monotonic timestamp for best results.
	// This a local time from the perspective of the entity that owns this ledger instance, to be used with the local
	// leaky bucket: it should NOT be sourced from the PaymentHeader
	now time.Time,
	// the timestamp included, or planned to be included, in the PaymentHeader
	dispersalTime time.Time,
	// the number of symbols to debit
	symbolCount uint32,
	// the quorums being dispersed to
	quorums []core.QuorumID,
) (bool, float64, error) {

	err := rl.config.reservation.CheckQuorumsPermitted(quorums)
	if err != nil {
		return false, 0, fmt.Errorf("check quorums permitted: %w", err)
	}

	err = rl.config.reservation.CheckTime(dispersalTime)
	if err != nil {
		return false, 0, fmt.Errorf("check time: %w", err)
	}

	billableSymbols := payments.CalculateBillableSymbols(symbolCount, rl.config.minNumSymbols)

	rl.lock.Lock()
	defer rl.lock.Unlock()

	success, err := rl.leakyBucket.Fill(now, float64(billableSymbols))
	if err != nil {
		return false, 0, fmt.Errorf("fill: %w", err)
	}

	remainingCapacity := rl.leakyBucket.GetRemainingCapacity()

	return success, remainingCapacity, nil
}

// Credit the reservation with a number of symbols. This method "undoes" a previous debit, following a failed dispersal.
//
// Note that this method doesn't reset the state of the ledger to be the same as when the debit was made: it just
// "refunds" the amount of symbols that were originally debited. Since the leaky bucket backing the reservation can't
// get emptier than "empty", it may be the case that only a portion of the debit is reverted, with the final capacity
// being clamped to 0.
//
// Returns the remaining capacity in the bucket after the revert operation.
func (rl *ReservationLedger) RevertDebit(now time.Time, symbolCount uint32) (float64, error) {
	billableSymbols := payments.CalculateBillableSymbols(symbolCount, rl.config.minNumSymbols)

	rl.lock.Lock()
	defer rl.lock.Unlock()

	err := rl.leakyBucket.RevertFill(now, float64(billableSymbols))
	if err != nil {
		return 0, fmt.Errorf("revert fill: %w", err)
	}

	remainingCapacity := rl.leakyBucket.GetRemainingCapacity()

	return remainingCapacity, nil
}

// Checks if the underlying leaky bucket is empty.
//
// This method cannot be used as an oracle to determine whether the bucket will be empty at some point in the future:
// it causes the ledger to update it's internal state, so only an *honest* representation of "now" should be provided.
func (rl *ReservationLedger) IsBucketEmpty(now time.Time) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	// Intentionally ignore the error here, can only happen if time moved backwards.
	fillLevel, _ := rl.leakyBucket.GetFillLevel(now)

	return fillLevel <= 0
}

// UpdateReservation updates the reservation parameters and recreates the leaky bucket, if necessary
//
// This method replaces the current reservation with a new one if the new reservation differs from the old.
//
// When an update occurs, the leaky bucket is recreated with the new parameters, but the old bucket
// state is preserved by starting the new bucket with the same fill level as the old.
//
// Returns an error if:
//   - newReservation is nil
//   - the new reservation configuration is invalid
//   - there's an error creating the new leaky bucket
func (rl *ReservationLedger) UpdateReservation(newReservation *Reservation, now time.Time) error {
	if newReservation == nil {
		return fmt.Errorf("newReservation cannot be nil")
	}

	rl.lock.Lock()
	defer rl.lock.Unlock()

	if rl.config.reservation.Equal(newReservation) {
		// if the reservation didn't change, there isn't anything to do
		return nil
	}

	// Create new config with the updated reservation
	newConfig, err := NewReservationLedgerConfig(
		*newReservation,
		rl.config.minNumSymbols,
		rl.config.startFull,
		rl.config.overfillBehavior,
		rl.config.bucketCapacityDuration)
	if err != nil {
		return fmt.Errorf("new reservation ledger config: %w", err)
	}
	rl.config = *newConfig

	err = rl.leakyBucket.Reconfigure(
		float64(newConfig.reservation.symbolsPerSecond),
		newConfig.bucketCapacityDuration,
		newConfig.overfillBehavior,
		now)
	if err != nil {
		return fmt.Errorf("reconfigure leaky bucket: %w", err)
	}

	return nil
}

// Returns the total bucket capacity in symbols
func (rl *ReservationLedger) GetBucketCapacity() float64 {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	return rl.leakyBucket.GetCapacity()
}
