package ondemand

import (
	"context"
	"math/big"
)

// Implements the CumulativePaymentStore interface, by storing values in memory
//
// This struct supports decrementing cumulative payments.
//
// NOTE: This struct doesn't do any synchronization! The caller is responsible for making sure that only one goroutine
// is using it at a time.
type EphemeralCumulativePaymentStore struct {
	cumulativePayment *big.Int
}

var _ CumulativePaymentStore = (*EphemeralCumulativePaymentStore)(nil)

// Constructs a new in-memory cumulative payment store
func NewEphemeralCumulativePaymentStore() *EphemeralCumulativePaymentStore {
	return &EphemeralCumulativePaymentStore{
		cumulativePayment: big.NewInt(0),
	}
}

// Gets the stored cumulative payment in wei
func (e *EphemeralCumulativePaymentStore) GetCumulativePayment(_ context.Context) (*big.Int, error) {
	return new(big.Int).Set(e.cumulativePayment), nil
}

// Sets the cumulative payment in wei, overwriting the previous value
func (e *EphemeralCumulativePaymentStore) SetCumulativePayment(_ context.Context, newCumulativePayment *big.Int) error {
	e.cumulativePayment.Set(newCumulativePayment)
	return nil
}
