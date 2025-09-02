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
// Returns (true, nil) if the reservation has enough capacity to perform the debit.
// Returns (false, nil) if the bucket lacks capacity to permit the fill.
// Returns (false, error) if an error occurs. Possible errors include:
//   - [QuorumNotPermittedError]: one or more of the requested quorums are not permitted by the reservation
//   - [TimeOutOfRangeError]: the dispersal time is outside the reservation's valid time range
//   - [TimeMovedBackwardError]: current time is before a previously observed time (only possible if input time
//     instances don't included monotonic timestamps)
//   - Generic errors for all other unexpected behavior
//
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
) (bool, error) {

	err := rl.config.reservation.CheckQuorumsPermitted(quorums)
	if err != nil {
		return false, fmt.Errorf("check quorums permitted: %w", err)
	}

	err = rl.config.reservation.CheckTime(dispersalTime)
	if err != nil {
		return false, fmt.Errorf("check time: %w", err)
	}

	rl.lock.Lock()
	defer rl.lock.Unlock()

	success, err := rl.leakyBucket.Fill(now, symbolCount)
	if err != nil {
		return false, fmt.Errorf("fill: %w", err)
	}

	return success, nil
}

// Credit the reservation with a number of symbols. This method "undoes" a previous debit, following a failed dispersal.
//
// Note that this method doesn't reset the state of the ledger to be the same as when the debit was made: it just
// "refunds" the amount of symbols that were originally debited. Since the leaky bucket backing the reservation can't
// get emptier than "empty", it may be the case that only a portion of the debit is reverted, with the final capacity
// being clamped to 0.
func (rl *ReservationLedger) RevertDebit(now time.Time, symbolCount uint32) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	err := rl.leakyBucket.RevertFill(now, symbolCount)
	if err != nil {
		return fmt.Errorf("revert fill: %w", err)
	}

	return nil
}

// Checks if the underlying leaky bucket is empty.
func (rl *ReservationLedger) IsBucketEmpty(now time.Time) (bool, error) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	fillLevel, err := rl.leakyBucket.CheckFillLevel(now)
	if err != nil {
		return false, fmt.Errorf("check fill level: %w", err)
	}

	return fillLevel <= 0, nil
}

// UpdateReservation updates the reservation parameters and recreates the leaky bucket
//
// This method completely replaces the current reservation with a new one. The leaky bucket
// is recreated with the new parameters, but preserves the current bucket state by starting
// with the same fill level as the previous bucket had at the update time.
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

	// Create new config with the updated reservation
	newConfig := ReservationLedgerConfig{
		reservation:            *newReservation,
		startFull:              rl.config.startFull,
		overfillBehavior:       rl.config.overfillBehavior,
		bucketCapacityDuration: rl.config.bucketCapacityDuration,
	}

	// Get current bucket state to preserve fill level during transition
	// Note: We can't directly call leak() as it's private, so we'll use the current state
	// This means there might be some imprecision if significant time has passed since the last operation
	oldCapacity := rl.leakyBucket.bucketCapacity
	currentFillLevel := rl.leakyBucket.currentFillLevel

	// Create new leaky bucket
	newLeakyBucket, err := NewLeakyBucket(
		newConfig.reservation.symbolsPerSecond,
		newConfig.bucketCapacityDuration,
		false, // start empty - we'll set the fill level below
		newConfig.overfillBehavior,
		now,
	)
	if err != nil {
		return fmt.Errorf("create new leaky bucket: %w", err)
	}

	// Preserve the fill level proportionally if the bucket capacity changed
	newCapacity := newLeakyBucket.bucketCapacity
	var newFillLevel float64
	if oldCapacity > 0 {
		// Scale the current fill level proportionally to the new capacity
		// This maintains the same relative "fullness" of the bucket
		newFillLevel = (currentFillLevel * newCapacity) / oldCapacity
		if newFillLevel > newCapacity {
			newFillLevel = newCapacity
		}
	}

	// Set the new fill level directly in the bucket
	newLeakyBucket.currentFillLevel = newFillLevel

	// Update the ledger with new config and bucket
	rl.config = newConfig
	rl.leakyBucket = newLeakyBucket

	return nil
}
