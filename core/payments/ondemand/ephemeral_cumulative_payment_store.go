package ondemand

import (
	"errors"
	"math/big"
)

type EphemeralCumulativePaymentStore struct {
	cumulativePayment *big.Int
}

var _ CumulativePaymentStore = (*EphemeralCumulativePaymentStore)(nil)

func NewEphemeralCumulativePaymentStore() *EphemeralCumulativePaymentStore {
	return &EphemeralCumulativePaymentStore{
		cumulativePayment: big.NewInt(0),
	}
}

func (e *EphemeralCumulativePaymentStore) GetCumulativePayment() (*big.Int, error) {
	if e.cumulativePayment == nil {
		return nil, errors.New("underlying cumulative payment is nil")
	}
	return new(big.Int).Set(e.cumulativePayment), nil
}

func (e *EphemeralCumulativePaymentStore) SetCumulativePayment(newCumulativePayment *big.Int) error {
	if e.cumulativePayment == nil {
		return errors.New("newCumulativePayment is nil")
	}

	e.cumulativePayment.Set(newCumulativePayment)
	return nil
}
