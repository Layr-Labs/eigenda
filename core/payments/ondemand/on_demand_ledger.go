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
// amount paid by the account for that and all previous dispersals. The cumulative payment is chosen by the dispersing
// client based on the state of its local accounting, and the chosen value can be verified by all other parties by
// checking:
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
	config OnDemandLedgerConfig
	// synchronizes access to the cumulative payment store
	lock sync.Mutex
	// stores the cumulative payment for this ledger in wei
	cumulativePaymentStore CumulativePaymentStore
}

// Constructs a new OnDemandLedger, backed by the input CumulativePaymentStore
func NewOnDemandLedger(
	config OnDemandLedgerConfig,
	cumulativePaymentStore CumulativePaymentStore,
) (*OnDemandLedger, error) {
	return &OnDemandLedger{
		config:                 config,
		cumulativePaymentStore: cumulativePaymentStore,
	}, nil
}

// Debit the on-demand account with the cost of a dispersal, based on the number of symbols.
//
// Returns (cumulativePayment, nil) if the account has sufficient funds to perform the debit.
// The returned cumulativePayment represents the new total amount spent from this account.
//
// Returns (nil, error) if an error occurs. Possible errors include:
//   - ErrQuorumNotSupported: requested quorums are not supported for on-demand payments
//   - ErrInsufficientFunds: the debit would exceed the total deposits available
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
		return nil, fmt.Errorf("%w: %s", ErrQuorumNotSupported, err.Error())
	}

	blobCost := odl.computeCost(symbolCount)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	currentCumulativePayment, err := odl.cumulativePaymentStore.GetCumulativePayment(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cumulative payment: %w", err)
	}

	newCumulativePayment := new(big.Int).Add(currentCumulativePayment, blobCost)

	if newCumulativePayment.Cmp(odl.config.totalDeposits) > 0 {
		return nil, fmt.Errorf(
			"%w: current cumulative payment: %s wei, total deposits: %s wei, blob cost: %s wei",
			ErrInsufficientFunds,
			currentCumulativePayment.String(),
			odl.config.totalDeposits.String(),
			blobCost.String())
	}

	if err := odl.cumulativePaymentStore.SetCumulativePayment(ctx, newCumulativePayment); err != nil {
		return nil, fmt.Errorf("set cumulative payment: %w", err)
	}

	return newCumulativePayment, nil
}

// RevertDebit reverts a previous debit operation, following a failed dispersal.
//
// Note: this method will only succeed if the underlying CumulativePaymentStore supports decrementing the
// cumulative payment.
func (odl *OnDemandLedger) RevertDebit(ctx context.Context, symbolCount uint32) error {
	if symbolCount == 0 {
		return errors.New("symbolCount must be > 0")
	}

	blobCost := odl.computeCost(symbolCount)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	currentCumulativePayment, err := odl.cumulativePaymentStore.GetCumulativePayment(ctx)
	if err != nil {
		return fmt.Errorf("get cumulative payment: %w", err)
	}

	newCumulativePayment := new(big.Int).Sub(currentCumulativePayment, blobCost)

	if newCumulativePayment.Sign() < 0 {
		return fmt.Errorf("cannot revert debit: would result in negative cumulative payment")
	}

	if err := odl.cumulativePaymentStore.SetCumulativePayment(ctx, newCumulativePayment); err != nil {
		return fmt.Errorf("set cumulative payment: %w", err)
	}

	return nil
}

// Checks whether all input quorum IDs are supported for on demand payments
//
// Returns an error if any input quorum isn't supported, otherwise nil
func checkForOnDemandSupport(quorumsToCheck []core.QuorumID) error {
	for _, quorum := range quorumsToCheck {
		if quorum == 0 || quorum == 1 {
			continue
		}

		return fmt.Errorf("quorum %d not in supported set [0, 1]", quorum)
	}

	return nil
}

// Computes the on demand cost of a number of symbols
func (odl *OnDemandLedger) computeCost(symbolCount uint32) *big.Int {
	symbolCountBig := big.NewInt(int64(symbolCount))

	billableSymbols := symbolCountBig
	if symbolCountBig.Cmp(odl.config.minNumSymbols) < 0 {
		billableSymbols = odl.config.minNumSymbols
	}

	return new(big.Int).Mul(billableSymbols, odl.config.pricePerSymbol)
}
