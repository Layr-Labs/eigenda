package ondemand

import (
	"math/big"
)

// Knows how to store and retrieve a cumulative payment value, stored as a *big.Int representing a number of wei
type CumulativePaymentStore interface {
	// Gets the stored cumulative payment in wei
	GetCumulativePayment() (*big.Int, error)

	// Sets the cumulative payment in wei. Overwrites the previously stored value
	SetCumulativePayment(newCumulativePayment *big.Int) error
}
