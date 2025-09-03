package reservation

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// Tracks usage of a single account reservation
//
// This struct is goroutine safe.
type ReservationLedger struct {
	config ReservationLedgerConfig

	// synchronizes access to the underlying leaky bucket algorithm
	lock sync.Mutex

	// an instance of the algorithm which tracks reservation usage
	leakyBucket *LeakyBucket
}

// Creates a new reservation ledger, which represents the reservation of a single user with a [LeakyBucket]
func NewReservationLedger(
	config ReservationLedgerConfig,
	// now should be from a source that includes monotonic timestamp for best results
	now time.Time,
) (*ReservationLedger, error) {
	leakyBucket, err := NewLeakyBucket(
		config.reservation.symbolsPerSecond,
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

	rl.lock.Lock()
	defer rl.lock.Unlock()

	success, err := rl.leakyBucket.Fill(now, symbolCount)
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
	rl.lock.Lock()
	defer rl.lock.Unlock()

	err := rl.leakyBucket.RevertFill(now, symbolCount)
	if err != nil {
		return 0, fmt.Errorf("revert fill: %w", err)
	}

	remainingCapacity := rl.leakyBucket.GetRemainingCapacity()

	return remainingCapacity, nil
}

// Checks if the underlying leaky bucket is empty.
func (rl *ReservationLedger) IsBucketEmpty(now time.Time) (bool, error) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	return rl.leakyBucket.CheckFillLevel(now) <= 0, nil
}

// UpdateReservation updates the reservation parameters and recreates the leaky bucket
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
	newConfig := ReservationLedgerConfig{
		reservation:            *newReservation,
		startFull:              rl.config.startFull,
		overfillBehavior:       rl.config.overfillBehavior,
		bucketCapacityDuration: rl.config.bucketCapacityDuration,
	}

	previousFillLevel := rl.leakyBucket.CheckFillLevel(now)

	newLeakyBucket, err := NewLeakyBucket(
		newConfig.reservation.symbolsPerSecond,
		newConfig.bucketCapacityDuration,
		false, // fill level is explicitly set below
		newConfig.overfillBehavior,
		now,
	)
	if err != nil {
		return fmt.Errorf("new leaky bucket: %w", err)
	}

	newLeakyBucket.currentFillLevel = previousFillLevel

	rl.config = newConfig
	rl.leakyBucket = newLeakyBucket

	return nil
}

// Returns the total bucket capacity in symbols
func (rl *ReservationLedger) GetBucketCapacity() float64 {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	return rl.leakyBucket.bucketCapacity
}
