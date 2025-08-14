package reservation

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// Keeps track of the state of a given reservation
//
// This is a goroutine safe wrapper around the LeakyBucket algorithm.
type ReservationLedger struct {
	config ReservationLedgerConfig

	// synchronizes access to the underlying leaky bucket algorithm
	lock sync.Mutex

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
		leakyBucket: leakyBucket,
	}, nil
}

// Debit the reservation with a number of symbols.
//
// Algorithmically, that means adding a number of symbols to the leaky bucket.
//
// Returns (true, nil) if the reservation has enough capacity to perform the debit.
// Returns (false, nil) if the bucket lacks capacity to permit the fill.
// Returns (false, error) if an error occurs. Possible errors include:
//   - ErrQuorumNotPermitted: requested quorums are not permitted by the reservation
//   - ErrTimeOutOfRange: dispersal time is outside the reservation's valid time range
//   - ErrLockAcquisition: failed to acquire the internal reservation lock
//   - ErrTimeMovedBackward: current time is before a previously observed time
//   - Generic errors for all other unexpected behavior
//
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed debit doesn't count against the meter.
func (rl *ReservationLedger) Debit(
	now time.Time,
	symbolCount uint32,
	quorums []core.QuorumID,
) (bool, error) {
	err := rl.config.reservation.CheckQuorumsPermitted(quorums)
	if err != nil {
		return false, fmt.Errorf("check quorums permitted: %w", err)
	}

	err = rl.config.reservation.CheckTime(now)
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
func (rl *ReservationLedger) RevertDebit(now time.Time, symbolCount uint32) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	err := rl.leakyBucket.RevertFill(now, symbolCount)
	if err != nil {
		return fmt.Errorf("revert fill: %w", err)
	}

	return nil
}
