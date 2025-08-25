package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
)

// Keeps track of the cumulative payment state for on-demand dispersals for a single account.
//
// On-demand payments use a cumulative payment system where, each time a dispersal is made, we keep track of the total
// amount paid by the account for that and all previous dispersals. The cumulative payment is chosen by the client
// based on the state of its local accounting, and the chosen value can be verified by checking:
// 1. that the claimed value is <= the total deposits belonging to the account in the PaymentVault contract
// 2. that the value has increased by at least the cost of the dispersal from the previously observed value
//
// The cost of each dispersal is calculated by multiplying the number of symbols (with a minimum of minNumSymbols) by
// the pricePerSymbol.
//
// On-demand payments are currently limited to quorums 0 (ETH) and 1 (EIGEN) and can only be used when dispersing
// through the EigenDA disperser.
//
// This is a goroutine safe struct.
type OnDemandLedger struct {
	// total deposits available for the account in wei
	totalDeposits *big.Int
	// price per symbol in wei
	pricePerSymbol *big.Int
	// minimum number of symbols to bill
	minNumSymbols uint64

	// an optional store to back the cumulative payment for this account
	//
	// if non-nil, the new cumulative payment value will be stored here after each debit
	cumulativePaymentStore *CumulativePaymentStore

	// the latest cumulative payment for the account
	cumulativePayment *big.Int
	// used to synchronize computation and optional storing of the cumulativePayment
	lock sync.Mutex
}

// Creates a new OnDemandLedger, backed by a CumulativePaymentStore
//
// The CumulativePaymentStore is used in this constructor to get the current cumulative payment value. After each
// debit, the latest cumulative payment will be stored in the CumulativePaymentStore.
func OnDemandLedgerFromStore(
	ctx context.Context,
	totalDeposits *big.Int,
	pricePerSymbol *big.Int,
	minNumSymbols uint64,
	cumulativePaymentStore *CumulativePaymentStore,
) (*OnDemandLedger, error) {
	if cumulativePaymentStore == nil {
		return nil, errors.New("cumulativePaymentStore cannot be nil")
	}

	cumulativePayment, err := cumulativePaymentStore.GetCumulativePayment(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cumulative payment from store: %w", err)
	}

	return newOnDemandLedger(totalDeposits, pricePerSymbol, minNumSymbols, cumulativePaymentStore, cumulativePayment)
}

// Creates a new OnDemandLedger, which *isn't* backed by a CumulativePayment store: the only representation of the
// cumulative payment is in memory.
func OnDemandLedgerFromValue(
	totalDeposits *big.Int,
	pricePerSymbol *big.Int,
	minNumSymbols uint64,
	cumulativePayment *big.Int,
) (*OnDemandLedger, error) {
	return newOnDemandLedger(totalDeposits, pricePerSymbol, minNumSymbols, nil, cumulativePayment)
}

// Creates an OnDemandLedger, checking all input parameters
func newOnDemandLedger(
	totalDeposits *big.Int,
	pricePerSymbol *big.Int,
	minNumSymbols uint64,
	cumulativePaymentStore *CumulativePaymentStore,
	cumulativePayment *big.Int,
) (*OnDemandLedger, error) {
	if totalDeposits == nil {
		return nil, errors.New("totalDeposits cannot be nil")
	}
	if totalDeposits.Sign() < 0 {
		return nil, errors.New("totalDeposits cannot be negative")
	}

	if pricePerSymbol == nil {
		return nil, errors.New("pricePerSymbol cannot be nil")
	}
	if pricePerSymbol.Sign() < 0 {
		return nil, errors.New("pricePerSymbol cannot be negative")
	}

	if cumulativePayment == nil {
		return nil, errors.New("cumulativePayment cannot be nil")
	}
	if cumulativePayment.Sign() < 0 {
		return nil, errors.New("cumulativePayment cannot be negative")
	}
	if cumulativePayment.Cmp(totalDeposits) > 0 {
		return nil, errors.New("cumulativePayment cannot exceed totalDeposits")
	}

	return &OnDemandLedger{
		totalDeposits:          new(big.Int).Set(totalDeposits),
		pricePerSymbol:         new(big.Int).Set(pricePerSymbol),
		minNumSymbols:          minNumSymbols,
		cumulativePaymentStore: cumulativePaymentStore,
		cumulativePayment:      new(big.Int).Set(cumulativePayment),
	}, nil
}

// Debit the on-demand account with the cost of a dispersal, based on the number of symbols.
//
// Returns (cumulativePayment, nil) if the account has sufficient funds to perform the debit.
// The returned cumulativePayment represents the new total amount spent from this account, including this blob.
//
// Returns (nil, error) if an error occurs. Possible errors include:
//   - [QuorumNotSupportedError]: requested quorums are not supported for on-demand payments
//   - [InsufficientFundsError]: the debit would exceed the total deposits available
//   - Generic errors for all other unexpected behavior
//
// If the account doesn't have sufficient funds to accommodate the debit, the cumulative payment
// IS NOT updated, i.e. a failed debit doesn't modify the payment state.
func (odl *OnDemandLedger) Debit(
	ctx context.Context,
	symbolCount uint32,
	quorums []core.QuorumID,
) (*big.Int, error) {
	if symbolCount == 0 {
		return nil, errors.New("symbolCount must be > 0")
	}

	err := checkForOnDemandSupport(quorums)
	if err != nil {
		return nil, err
	}

	blobCost := odl.computeCost(symbolCount)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	newCumulativePayment := new(big.Int).Add(odl.cumulativePayment, blobCost)
	if newCumulativePayment.Cmp(odl.totalDeposits) > 0 {
		return nil, &InsufficientFundsError{
			CurrentCumulativePayment: odl.cumulativePayment,
			TotalDeposits:            odl.totalDeposits,
			BlobCost:                 blobCost,
		}
	}

	if odl.cumulativePaymentStore != nil {
		err := odl.cumulativePaymentStore.StoreCumulativePayment(ctx, newCumulativePayment)
		if err != nil {
			return nil, fmt.Errorf("store cumulative payment: %w", err)
		}
	}

	odl.cumulativePayment.Set(newCumulativePayment)

	return newCumulativePayment, nil
}

// RevertDebit reverts a previous debit operation, following a failed dispersal.
//
// Returns the new cumulative payment amount after the revert.
func (odl *OnDemandLedger) RevertDebit(ctx context.Context, symbolCount uint32) (*big.Int, error) {
	if symbolCount == 0 {
		return nil, errors.New("symbolCount must be > 0")
	}

	blobCost := odl.computeCost(symbolCount)
	blobCost.Neg(blobCost)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	newCumulativePayment := new(big.Int).Add(odl.cumulativePayment, blobCost)
	if newCumulativePayment.Sign() < 0 {
		return nil, fmt.Errorf("operation would result in negative cumulative payment: current=%s, addition amount=%s",
			odl.cumulativePayment.String(), blobCost.String())
	}

	if odl.cumulativePaymentStore != nil {
		err := odl.cumulativePaymentStore.StoreCumulativePayment(ctx, newCumulativePayment)
		if err != nil {
			return nil, fmt.Errorf("store cumulative payment: %w", err)
		}
	}

	odl.cumulativePayment.Set(newCumulativePayment)

	return newCumulativePayment, nil
}

// Checks whether all input quorum IDs are supported for on demand payments
//
// Returns an error if any input quorum isn't supported, otherwise nil
func checkForOnDemandSupport(quorumsToCheck []core.QuorumID) error {
	for _, quorum := range quorumsToCheck {
		if quorum == 0 || quorum == 1 {
			continue
		}

		return &QuorumNotSupportedError{
			RequestedQuorum:  quorum,
			SupportedQuorums: []core.QuorumID{0, 1},
		}
	}

	return nil
}

// Computes the on demand cost of a number of symbols
func (odl *OnDemandLedger) computeCost(symbolCount uint32) *big.Int {
	billableSymbols := uint64(symbolCount)
	if billableSymbols < odl.minNumSymbols {
		billableSymbols = odl.minNumSymbols
	}

	return new(big.Int).Mul(big.NewInt(int64(billableSymbols)), odl.pricePerSymbol)
}
