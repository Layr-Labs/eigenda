package reservation

import (
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// Represents a reservation for a single user.
//
// TODO(litt3): I opted to duplicate the preexisting `ReservedPayment` struct, rather than using the old one. There
// are nontrivial changes I wanted to make, and making those changes in a way that's compatible with the preexisting
// usages was going to be messy. Instead, `ReservedPayment` can just be removed, when we remove the deprecated payment
// system.
type Reservation struct {
	// The number of symbols / second that the holder of this reservation is entitled to disperse
	symbolsPerSecond uint64

	// The time at which the reservation becomes active
	startTime time.Time

	// The time at which the reservation expires
	endTime time.Time

	// The quorums that the holder of this reservation is entitled to disperse to
	permittedQuorumIDs map[core.QuorumID]bool
}

// Create a representation of a single account Reservation.
func NewReservation(
	symbolsPerSecond uint64,
	startTime time.Time,
	endTime time.Time,
	permittedQuorumIDs []core.QuorumID,
) (*Reservation, error) {
	if symbolsPerSecond == 0 {
		return nil, errors.New("reservation must have >0 symbols per second")
	}

	if startTime.Equal(endTime) || endTime.Before(startTime) {
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
		symbolsPerSecond:   symbolsPerSecond,
		startTime:          startTime,
		endTime:            endTime,
		permittedQuorumIDs: permittedQuorumIDSet,
	}, nil
}

// Checks whether an input list of quorums are all permitted by the reservation.
//
// Returns nil if all input quorums are permitted, otherwise returns ErrQuorumNotPermitted.
func (r *Reservation) CheckQuorumsPermitted(quorums []core.QuorumID) error {
	for _, quorum := range quorums {
		if !r.permittedQuorumIDs[quorum] {
			permittedQuorums := make([]core.QuorumID, 0, len(r.permittedQuorumIDs))
			for quorumID := range r.permittedQuorumIDs {
				permittedQuorums = append(permittedQuorums, quorumID)
			}
			return fmt.Errorf("%w: quorum %d not in permitted set %v", ErrQuorumNotPermitted, quorum, permittedQuorums)
		}
	}

	return nil
}

// Verifies that the given time falls within the reservation's valid time range.
//
// Returns ErrTimeOutOfRange if the time is outside the valid range.
func (r *Reservation) CheckTime(timeToCheck time.Time) error {
	if timeToCheck.Before(r.startTime) || timeToCheck.After(r.endTime) {
		return fmt.Errorf("%w: dispersal time %s is outside range [%s, %s]",
			ErrTimeOutOfRange,
			timeToCheck.Format(time.RFC3339),
			r.startTime.Format(time.RFC3339),
			r.endTime.Format(time.RFC3339))
	}

	return nil
}
