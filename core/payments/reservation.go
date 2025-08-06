package payments

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TODO: write docs somewhere about implications of adding a new reservation, or an old reservation expiring, while
// a client is running. I think the correct thing would be to just restart the client if a new reservation is made...

// Represents a reservation for a single user.
//
// TODO(litt3): I opted to duplicate the preexisting `ReservedPayment` struct, rather than using the old one. There
// are nontrivial changes I wanted to make, and making those changes in a way that's compatible with the preexisting
// usages was going to be messy. Instead, `ReservedPayment` can just be removed, when we remove the deprecated payment
// system.
type Reservation struct {
	accountID gethcommon.Address

	// The number of symbols / second that the holder of this reservation is entitled to disperse
	//
	// The leak rate is a uint64 on-chain, so that's what we accept in the constructor. But it's converted to an int64
	// under the hood to match other int types in reservation tracking, making syntax less verbose.
	//
	// The assumption that symbolsPerSecondLeakRate fits in an int64 should be fine, as long as we are careful not to
	// give out any reservations that exceed ~590 exabytes/second.
	symbolsPerSecond int64

	// The time at which the reservation becomes active
	startTime time.Time

	// The time at which the reservation expires
	endTime time.Time

	// The quorums that the holder of this reservation is entitled to disperse to
	permittedQuorumIDs map[core.QuorumID]bool
}

// TODO doc
func NewReservation(
	accountID gethcommon.Address,
	symbolsPerSecond uint64,
	startTime time.Time,
	endTime time.Time,
	permittedQuorumIDs []core.QuorumID,
) (*Reservation, error) {
	if accountID == (gethcommon.Address{}) {
		return nil, fmt.Errorf("account ID cannot be zero address")
	}

	if symbolsPerSecond <= 0 {
		return nil, fmt.Errorf("reservation must have >0 symbols per second, got %d", symbolsPerSecond)
	}

	if symbolsPerSecond > math.MaxInt64 {
		return nil, fmt.Errorf("symbolsPerSecond must be < math.MaxInt64 (got %d). Technically, anything up to "+
			"math.MaxUint64 is permitted on-chain, but practical limits are put in place to simplify implementation.",
			symbolsPerSecond)
	}

	if startTime == endTime || endTime.Before(startTime) {
		return nil, fmt.Errorf("start time (%v) must be before end time (%v)", startTime, endTime)
	}

	permittedQuorumIDsLen := len(permittedQuorumIDs)
	if permittedQuorumIDsLen == 0 {
		return nil, errors.New("reservation must permit at least one quorum")
	}

	permittedQuorumIDSet := make(map[core.QuorumID]bool, permittedQuorumIDsLen)
	for _, quorumID := range permittedQuorumIDs {
		permittedQuorumIDSet[quorumID] = true
	}

	return &Reservation{
		accountID:          accountID,
		symbolsPerSecond:   int64(symbolsPerSecond),
		startTime:          startTime,
		endTime:            endTime,
		permittedQuorumIDs: permittedQuorumIDSet,
	}, nil
}

// Checks whether an input list of quorums are all permitted by the reservation.
//
// Returns nil if all input quorums are permitted, otherwise an error.
func (r *Reservation) CheckQuorumsPermitted(quorums []core.QuorumID) error {
	for _, quorum := range quorums {
		if !r.permittedQuorumIDs[quorum] {
			return fmt.Errorf("quorum %v not permitted", quorum)
		}
	}

	return nil
}

// TODO doc
func (r *Reservation) CheckTime(timeToCheck time.Time) error {
	if timeToCheck.Before(r.startTime) {
		return fmt.Errorf("timeToCheck %s is before reservation start time %s",
			timeToCheck.Format(time.RFC3339), r.startTime.Format(time.RFC3339))
	}

	if timeToCheck.After(r.endTime) {
		return fmt.Errorf("timeToCheck %s is after reservation end time %s",
			timeToCheck.Format(time.RFC3339), r.endTime.Format(time.RFC3339))
	}

	return nil
}
