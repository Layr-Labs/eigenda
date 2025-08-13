package ondemand

import (
	"math/big"
)

// Implements the CumulativePaymentStore interface, by storing values in memory
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

// Gets the stored cumulative payment.
func (e *EphemeralCumulativePaymentStore) GetCumulativePayment() (*big.Int, error) {
	return new(big.Int).Set(e.cumulativePayment), nil
}

// Sets the cumulative payment, overwriting the previous value
func (e *EphemeralCumulativePaymentStore) SetCumulativePayment(newCumulativePayment *big.Int) error {
	e.cumulativePayment.Set(newCumulativePayment)
	return nil
}
