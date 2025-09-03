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
// [ReservationLedger.CheckInvariants] should be called prior to calling [ReservationLedger.Debit], to make sure the
// dispersal is permitted under the parameters of the reservation. If [ReservationLedger.CheckInvariants] succeeds,
// then [ReservationLedger.Debit] is called to make sure the dispersal doesn't exceed reservation capacity.
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
