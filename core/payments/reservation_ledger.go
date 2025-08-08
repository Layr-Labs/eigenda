package payments

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

// TODO: consider extracting leaky bucket out of here after all. As I go, there is more and more stuff that isn't related
// to leaky bucket being added.

// TODO: at what point in the process will be construct the payment header? I don't want to do it in this class,
// since then it won't be reusable: only the client should be creating the payment header, everyone else should
// be extracting data, and verifying that the dispersal is permitted.


type ReservationLedger struct {
	config ReservationLedgerConfig

	lock        *semaphore.Weighted
	leakyBucket *LeakyBucket
}

// Creates a new reservation ledger, which represents the reservation of a single user with a leaky bucket
func NewReservationLedger(
	config ReservationLedgerConfig,
	now time.Time,
) (*ReservationLedger, error) {
	leakyBucket, err := NewLeakyBucket(config, now)
	if err != nil {
		return nil, fmt.Errorf("new leaky bucket: %w", err)
	}

	return &ReservationLedger{
		config:      config,
		lock:        semaphore.NewWeighted(1),
		leakyBucket: leakyBucket,
	}, nil
}

// TODO: consider whether the concept of a debit slip makes sense

// Debit the reservation with a number of symbols.
//
// Algorithmically, that means adding a number of symbols to the leaky bucket.
//
// Returns nil if the leaky bucket has enough capacity to accept the fill. Returns an
// InsufficientReservationCapacityError if bucket lacks capacity to permit the fill.
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
		// TODO: error here
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
