package ephemeral

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
)

// Implements the CumulativePaymentStore interface, by storing values in memory
//
// This struct supports decrementing cumulative payments.
//
// As specified by the interface contract, this struct is goroutine safe
type EphemeralCumulativePaymentStore struct {
	cumulativePayment *big.Int
	lock              sync.Mutex
}

var _ ondemand.CumulativePaymentStore = (*EphemeralCumulativePaymentStore)(nil)

// Constructs a new in-memory cumulative payment store
func NewEphemeralCumulativePaymentStore() *EphemeralCumulativePaymentStore {
	return &EphemeralCumulativePaymentStore{
		cumulativePayment: big.NewInt(0),
	}
}

// Atomically increments the cumulative payment by the given amount.
func (e *EphemeralCumulativePaymentStore) AddCumulativePayment(
	_ context.Context,
	amount *big.Int,
	maxCumulativePayment *big.Int,
) (*big.Int, error) {
	if amount == nil {
		return nil, errors.New("amount cannot be nil")
	}
	if maxCumulativePayment == nil {
		return nil, errors.New("maxCumulativePayment cannot be nil")
	}
	if maxCumulativePayment.Sign() < 0 {
		return nil, fmt.Errorf("maxCumulativePayment cannot be negative: received %s", maxCumulativePayment.String())
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	newCumulativePayment := new(big.Int).Add(e.cumulativePayment, amount)

	if newCumulativePayment.Sign() < 0 {
		return nil, fmt.Errorf("operation would result in negative cumulative payment: current=%s, addition amount=%s",
			e.cumulativePayment.String(), amount.String())
	}

	if newCumulativePayment.Cmp(maxCumulativePayment) > 0 {
		return nil, &ondemand.InsufficientFundsError{
			CurrentCumulativePayment: e.cumulativePayment,
			MaxCumulativePayment:     maxCumulativePayment,
			BlobCost:                 amount,
		}
	}

	e.cumulativePayment.Set(newCumulativePayment)

	// Return the copy we made, so the caller can't modify the internal cumulativePayment value
	return newCumulativePayment, nil
}
