package ondemand

import (
	"errors"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TotalDepositUpdate represents an update to total deposits for a specific account
type TotalDepositUpdate struct {
	AccountAddress  gethcommon.Address
	NewTotalDeposit *big.Int
}

// NewTotalDepositUpdate creates a new TotalDepositUpdate with proper validation
//
// Returns an error if:
// - accountAddress is the zero address
// - newTotalDeposit is nil
// - newTotalDeposit is negative
func NewTotalDepositUpdate(accountAddress gethcommon.Address, newTotalDeposit *big.Int) (*TotalDepositUpdate, error) {
	if accountAddress == (gethcommon.Address{}) {
		return nil, errors.New("accountAddress cannot be zero address")
	}

	if newTotalDeposit == nil {
		return nil, errors.New("newTotalDeposit cannot be nil")
	}

	if newTotalDeposit.Sign() < 0 {
		return nil, errors.New("newTotalDeposit cannot be negative")
	}

	return &TotalDepositUpdate{
		AccountAddress:  accountAddress,
		NewTotalDeposit: new(big.Int).Set(newTotalDeposit), // Create a copy to avoid shared references
	}, nil
}
