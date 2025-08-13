package ondemand

import (
	"math/big"
)

type CumulativePaymentStore interface {
	GetCumulativePayment() (*big.Int, error)
	SetCumulativePayment(newCumulativePayment *big.Int) error
}
