package reservation

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"golang.org/x/sync/semaphore"
)

// TODO: write unit tests

// TODO: where do we check whether a dispersal fits within the correct time? it seems like we might need to return
// the time to put into the payment header when doing the debit function....

// TODO: at what point in the process will be construct the payment header? I don't want to do it in this class,
// since then it won't be reusable: only the client should be creating the payment header, everyone else should
// be extracting data, and verifying that the dispersal is permitted.

// Keeps track of the state of a given reservation
//
// This is a goroutine safe wrapper around the LeakyBucket algorithm.
type ReservationLedger struct {
	config ReservationLedgerConfig

	// synchronizes access to the underlying leaky bucket algorithm
	// this is a semaphore instead of a lock, for the sake of fairness: goroutines should acquire the lock in the
	// order that requests arrive
	lock *semaphore.Weighted

	// an instance of the algorithm which tracks reservation usage
	leakyBucket *LeakyBucket
}

// Creates a new reservation ledger, which represents the reservation of a single user with a leaky bucket
func NewReservationLedger(
	config ReservationLedgerConfig,
	now time.Time,
) (*ReservationLedger, error) {
	leakyBucket, err := NewLeakyBucket(
		config.reservation.symbolsPerSecond,
		config.bucketCapacityDuration,
		config.biasBehavior,
		config.overfillBehavior,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("new leaky bucket: %w", err)
	}

	return &ReservationLedger{
		config:      config,
		lock:        semaphore.NewWeighted(1),
		leakyBucket: leakyBucket,
	}, nil
}

// Debit the reservation with a number of symbols.
//
// Algorithmically, that means adding a number of symbols to the leaky bucket.
//
// Returns nil if the leaky bucket has enough capacity to accept the fill. Returns an
// InsufficientReservationCapacityError if bucket lacks capacity to permit the fill.

// - Returns nil if the reservation has enough capacity to perform the debit.
// - Returns an InsufficientReservationCapacityError if bucket lacks capacity to perform the debit
// - Returns a TimeMovedBackwardError if input time is before previous leak time.
// - Returns a generic error for all other modes of failure.
//
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed debit doesn't count against the meter.
func (rl *ReservationLedger) Debit(
	ctx context.Context,
	now time.Time,
	symbolCount int64,
	quorums []core.QuorumID,
) error {
	err := rl.config.reservation.CheckQuorumsPermitted(quorums)
	if err != nil {
		return fmt.Errorf("check quorums permitted: %w", err)
	}

	err = rl.config.reservation.CheckTime(now)
	if err != nil {
		return fmt.Errorf("check time: %w", err)
	}

	if err := rl.lock.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer rl.lock.Release(1)

	err = rl.leakyBucket.Fill(now, symbolCount)
	if err != nil {
		return fmt.Errorf("fill leaky bucket: %w", err)
	}

	return nil
}

// Credit the reservation with a number of symbols. This method "undoes" a previous debit, following a failed dispersal.
func (rl *ReservationLedger) RevertDebit(ctx context.Context, now time.Time, symbolCount int64) error {
	if err := rl.lock.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer rl.lock.Release(1)

	err := rl.leakyBucket.RevertFill(now, symbolCount)
	if err != nil {
		return fmt.Errorf("revert fill: %w", err)
	}

	return nil
}
