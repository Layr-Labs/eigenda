package payments

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
)

// InsufficientReservationCapacityError is returned when the leaky bucket doesn't have enough capacity to accommodate
// a requested dispersal.
type InsufficientReservationCapacityError struct {
	// The number of symbols that were requested to be dispersed
	RequestedSymbols int64
}

// Implements the error interface
func (e *InsufficientReservationCapacityError) Error() string {
	return fmt.Sprintf("insufficient reservation capacity to disperse %d symbols (%d bytes)",
		e.RequestedSymbols, e.RequestedSymbols*encoding.BYTES_PER_SYMBOL)
}

// InvalidReservationPeriod is returned when attempting to use a reservation outside its valid time window
type InvalidReservationPeriod struct {
	ReservationStartTime time.Time
	ReservationEndTime   time.Time
	TimeAttempted        time.Time
}

func (e *InvalidReservationPeriod) Error() string {
	if e.TimeAttempted.Before(e.ReservationStartTime) {
		return fmt.Sprintf("reservation not yet active: valid from %s to %s, attempted at %s",
			e.ReservationStartTime.Format(time.RFC3339),
			e.ReservationEndTime.Format(time.RFC3339),
			e.TimeAttempted.Format(time.RFC3339))
	}
	return fmt.Sprintf("reservation expired: valid from %s to %s, attempted at %s",
		e.ReservationStartTime.Format(time.RFC3339),
		e.ReservationEndTime.Format(time.RFC3339),
		e.TimeAttempted.Format(time.RFC3339))
}
