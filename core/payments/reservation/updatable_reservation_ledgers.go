package reservation

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Defines methods for updating reservation ledgers when changes are made to the parameters in the PaymentVault
type UpdatableReservationLedgers interface {
	// Returns the accounts included in this interface which should be updated with PaymentVault changes
	GetAccountsToUpdate() []gethcommon.Address

	// Updates the reservation for an account
	UpdateReservation(accountID gethcommon.Address, newReservation *Reservation) error
}
