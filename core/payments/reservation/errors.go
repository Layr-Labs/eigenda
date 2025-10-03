package reservation

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

// QuorumNotPermittedError indicates that a requested quorum is not permitted by the reservation.
type QuorumNotPermittedError struct {
	Quorum           core.QuorumID
	PermittedQuorums []core.QuorumID
}

func (e *QuorumNotPermittedError) Error() string {
	return fmt.Sprintf("quorum %d not in permitted set %v", e.Quorum, e.PermittedQuorums)
}

// TimeOutOfRangeError indicates the dispersal time is outside the reservation's valid time range.
type TimeOutOfRangeError struct {
	DispersalTime        time.Time
	ReservationStartTime time.Time
	ReservationEndTime   time.Time
}

func (e *TimeOutOfRangeError) Error() string {
	return fmt.Sprintf("dispersal time %s is outside permitted range [%s, %s]",
		e.DispersalTime.Format(time.RFC3339),
		e.ReservationStartTime.Format(time.RFC3339),
		e.ReservationEndTime.Format(time.RFC3339))
}
