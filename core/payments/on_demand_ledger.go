package payments

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"golang.org/x/sync/semaphore"
)

// TODO: we need to keep track of how many in flight dispersals there are, and not let that number exceed a certain
// value. The account ledger will need to check whether the on demand ledger is available before trying to debit,
// and do a wait if it isn't. We also need to consider how to "time out" an old request that was made to the disperser
// which was never responded to. We can't wait forever, eventually we need to declare a dispersal "failed", and move on

type OnDemandLedger struct {
	config            OnDemandLedgerConfig
	lock              *semaphore.Weighted
	cumulativePayment *big.Int
}

func NewOnDemandLedger(
	config OnDemandLedgerConfig,
) (*OnDemandLedger, error) {
	// TODO: get this from the disperser
	cumulativePayment := big.NewInt(0)

	return &OnDemandLedger{
		config:            config,
		lock:              semaphore.NewWeighted(1), // Binary semaphore acts as a mutex
		cumulativePayment: cumulativePayment,
	}, nil
}

// TODO: reconsider int64
func (odl *OnDemandLedger) Debit(
	ctx context.Context,
	now time.Time,
	symbolCount int64,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	if symbolCount <= 0 {
		return nil, fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	err := checkForOnDemandSupport(quorums)
	if err != nil {
		return nil, fmt.Errorf("check for on demand support: %w", err)
	}

	blobCost, err := odl.computeBlobCost(symbolCount)
	if err != nil {
		return nil, fmt.Errorf("compute blob cost: %w", err)
	}

	if err := odl.lock.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("acquire lock: %w", err)
	}
	defer odl.lock.Release(1)

	newCumulativePayment := new(big.Int).Add(odl.cumulativePayment, big.NewInt(blobCost))

	if newCumulativePayment.Cmp(odl.config.totalDeposits) > 0 {
		// TODO: make a specific error type with this, with appropriate fields
		return nil, fmt.Errorf("insufficient on-demand funds")
	}

	paymentMetadata, err := core.NewPaymentMetadata(odl.config.accountID, now, newCumulativePayment)
	if err != nil {
		return nil, fmt.Errorf("new payment metadata: %w", err)
	}

	odl.cumulativePayment = newCumulativePayment

	return paymentMetadata, nil
}

// RevertDebit reverts a previous debit operation, following a failed dispersal.
func (odl *OnDemandLedger) RevertDebit(ctx context.Context, symbolCount int64) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	blobCost, err := odl.computeBlobCost(symbolCount)
	if err != nil {
		return fmt.Errorf("compute blob cost: %w", err)
	}

	// Acquire the semaphore with context timeout
	if err := odl.lock.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer odl.lock.Release(1)

	newCumulativePayment := new(big.Int).Sub(odl.cumulativePayment, big.NewInt(blobCost))

	if newCumulativePayment.Sign() < 0 {
		return fmt.Errorf("cannot revert debit: would result in negative cumulative payment")
	}

	odl.cumulativePayment = newCumulativePayment

	return nil
}

func checkForOnDemandSupport(quorumsToCheck []core.QuorumID) error {
	for _, quorum := range quorumsToCheck {
		if quorum == 0 || quorum == 1 {
			continue
		}

		return fmt.Errorf("only quorums 0 and 1 are supported for on demand payments, got %d", quorum)
	}

	return nil
}

func (odl *OnDemandLedger) computeBlobCost(symbolCount int64) (int64, error) {
	if symbolCount <= 0 {
		return 0, fmt.Errorf("symbol count must be > 0, got %d", symbolCount)
	}

	billableSymbols := symbolCount
	if symbolCount < int64(odl.config.minNumSymbols) {
		billableSymbols = int64(odl.config.minNumSymbols)
	}

	return billableSymbols * int64(odl.config.pricePerSymbol), nil
}
