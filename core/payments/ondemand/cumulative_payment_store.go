package ondemand

import (
	"context"
	"math/big"
)

// Stores a cumulative payment, and can atomically add values to it while remaining within specified bounds.
//
// The cumulative payment is a *big.Int representing a number of wei
type CumulativePaymentStore interface {
	// Atomically adds a number of wei to the cumulative payment.
	//
	// May optionally support subtraction via negative amount parameter, depending on the implementation.
	//
	// Returns the new cumulative payment value after the addition.
	// Returns An an [ondemand.InsufficientFundsError] if the addition would cause the cumulative payment to exceed
	// the input maxCumulativePayment. In this case, the underlying cumulative payment value is not modified.
	AddCumulativePayment(
		ctx context.Context,
		// the amount to add to the cumulative payment, in wei
		amount *big.Int,
		// the maximum value for the cumulative payment, after the addition. Should be determined by the `TotalDeposits`
		// for the account in the payment vault
		maxCumulativePayment *big.Int,
	) (*big.Int, error)
}
