package reservation

import (
	"errors"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ReservationUpdate represents an update to reservation parameters for a specific account
type ReservationUpdate struct {
	AccountAddress gethcommon.Address
	NewReservation *Reservation
}

// NewReservationUpdate creates a new ReservationUpdate with proper validation
//
// Returns an error if:
// - accountAddress is the zero address
// - newReservation is nil
func NewReservationUpdate(accountAddress gethcommon.Address, newReservation *Reservation) (*ReservationUpdate, error) {
	if accountAddress == (gethcommon.Address{}) {
		return nil, errors.New("accountAddress cannot be zero address")
	}

	if newReservation == nil {
		return nil, errors.New("newReservation cannot be nil")
	}

	return &ReservationUpdate{
		AccountAddress: accountAddress,
		NewReservation: newReservation,
	}, nil
}
