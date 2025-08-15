package ondemand

import (
	"context"
	"math/big"
)

// Knows how to store and retrieve a cumulative payment value, stored as a *big.Int representing a number of wei
type CumulativePaymentStore interface {
	// Gets the stored cumulative payment in wei
	GetCumulativePayment(ctx context.Context) (*big.Int, error)

	// Sets the cumulative payment in wei. Overwrites the previously stored value
	//
	// Returns an error if newCumulativePayment param is nil or < 0
	SetCumulativePayment(ctx context.Context, newCumulativePayment *big.Int) error
}
