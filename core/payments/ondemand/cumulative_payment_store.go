package ondemand

import (
	"math/big"
)

// Knows how to store and retrieve cumulative payment values, stored as *big.Ints
type CumulativePaymentStore interface {
	// Gets the stored cumulative payment
	GetCumulativePayment() (*big.Int, error)

	// Sets the cumulative payment. Overwrites the previously stored value
	SetCumulativePayment(newCumulativePayment *big.Int) error
}
