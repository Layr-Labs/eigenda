package payments

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// TODO: consider overflows in this. Are we being _too_ careful?

// TODO: write unit tests

// TODO: where do we check whether a dispersal fits within the correct time? it seems like we might need to return
// the time to put into the payment header when doing the debit function....

// TODO: consider extracting leaky bucket out of here after all. As I go, there is more and more stuff that isn't related
// to leaky bucket being added.

// TODO: at what point in the process will be construct the payment header? I don't want to do it in this class,
// since then it won't be reusable: only the client should be creating the payment header, everyone else should
// be extracting data, and verifying that the dispersal is permitted.

// This struct implements the [leaky bucket](https://en.wikipedia.org/wiki/Leaky_bucket) algorithm as a meter.
//
// Symbols "leak out" of the bucket at a constant rate, creating capacity for new symbols. The bucket can be "filled"
// with additional symbols if there is enough available capacity.
//
// The standard golang golang.org/x/time/rate.Limiter is not suitable for our use-case, for the following reasons:
//
//  1. The Limiter doesn't support the concept of overfilling the bucket. We require the concept of overfill, for cases
//     where a bucket size might be too small to fit the largest permissible blob size. We don't want to prevent users
//     with a small reservation size from submitting large blobs.
//  2. The Limiter uses floating point math. Though it would *probably* be ok to use floats, it makes the distributed
//     system harder to reason about. What level of error accumulation would we see with frequent updates? Under
//     what conditions would it be possible for the client and server representations of a given leaky bucket to
//     diverge, and what impact would that have on our assumptions? These questions can be avoided entirely by using
//     an integer based implementation.
//
// NOTE: methods on this struct should not be called from separate goroutines: it's not threadsafe.
type ReservationLedger struct {
	config ReservationLedgerConfig

	getNow func() time.Time

	leakyBucket *LeakyBucket
}

// Creates a new reservation ledger, which represents the reservation of a single user with a leaky bucket
func NewReservationLedger(
	config ReservationLedgerConfig,
	getNow func() time.Time,
) (*ReservationLedger, error) {
	leakyBucket, err := NewLeakyBucket(config, getNow())
	if err != nil {
		return nil, fmt.Errorf("new leaky bucket: %w", err)
	}

	return &ReservationLedger{
		getNow:      getNow,
		config:      config,
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
func (rl *ReservationLedger) Debit(symbolCount int64, quorums []core.QuorumID) (*core.PaymentMetadata, error) {
	err := rl.config.reservation.CheckQuorumsPermitted(quorums)
	if err != nil {
		return nil, fmt.Errorf("check quorums permitted: %w", err)
	}

	now := rl.getNow()

	err = rl.config.reservation.CheckTime(now)
	if err != nil {
		
	}

	err = rl.leakyBucket.Fill(now, symbolCount)
	if err != nil {
		return nil, fmt.Errorf("fill leaky bucket: %w", err)
	}

	paymentMetadata, err := core.NewPaymentMetadata(rl.config.reservation.accountID, now, nil)
	if err != nil {
		return nil, fmt.Errorf("new payment metadata: %w", err)
	}

	return paymentMetadata, nil
}

// Credit the reservation with a number of symbols. This method "undoes" a previous debit, following a failed dispersal.
func (rl *ReservationLedger) RevertDebit(symbolCount int64) error {
	err := rl.leakyBucket.RevertFill(rl.getNow(), symbolCount)
	if err != nil {
		return fmt.Errorf("revert fill: %w", err)
	}

	return nil
}
