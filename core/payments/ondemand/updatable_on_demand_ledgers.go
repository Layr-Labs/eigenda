package ondemand

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// UpdatableOnDemandLedgers defines the interface for a collection of OnDemandLedgers that need to be updated when
// there are changes to the state of the payment vault.
type UpdatableOnDemandLedgers interface {
	// Returns the accounts included in this interface which should be updated with PaymentVault changes
	GetAccountsToUpdate() []gethcommon.Address

	// Updates the total deposit for an account
	UpdateTotalDeposit(accountID gethcommon.Address, newTotalDeposit *big.Int) error
}
