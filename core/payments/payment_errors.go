package payments

import (
	"fmt"

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
