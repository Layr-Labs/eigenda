package ondemand

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
)

// TODO: we need to keep track of how many in flight dispersals there are, and not let that number exceed a certain
// value. The client ledger will need to check whether the on demand ledger is available before trying to debit,
// and do a wait if it isn't. We also need to consider how to "time out" an old request that was made to the disperser
// which was never responded to. We can't wait forever, eventually we need to declare a dispersal "failed", and move on

type OnDemandLedger struct {
	config OnDemandLedgerConfig
	// synchronizes access to the cumulative payment store
	lock                   sync.Mutex
	cumulativePaymentStore CumulativePaymentStore
}

func NewOnDemandLedger(
	config OnDemandLedgerConfig,
	cumulativePaymentStore CumulativePaymentStore,
) (*OnDemandLedger, error) {
	return &OnDemandLedger{
		config:                 config,
		cumulativePaymentStore: cumulativePaymentStore,
	}, nil
}

func (odl *OnDemandLedger) Debit(
	symbolCount uint32,
	quorums []core.QuorumID,
) (*big.Int, error) {
	if symbolCount == 0 {
		return nil, errors.New("symbolCount must be > 0")
	}

	err := checkForOnDemandSupport(quorums)
	if err != nil {
		return nil, fmt.Errorf("check for on demand support: %w", err)
	}

	blobCost := odl.computeBlobCost(symbolCount)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	currentCumulativePayment, err := odl.cumulativePaymentStore.GetCumulativePayment()
	if err != nil {
		return nil, fmt.Errorf("get cumulative payment: %w", err)
	}

	newCumulativePayment := new(big.Int).Add(currentCumulativePayment, blobCost)

	if newCumulativePayment.Cmp(odl.config.totalDeposits) > 0 {
		// TODO: make a specific error type with this, with appropriate fields
		return nil, fmt.Errorf("insufficient on-demand funds")
	}

	if err := odl.cumulativePaymentStore.SetCumulativePayment(newCumulativePayment); err != nil {
		return nil, fmt.Errorf("set cumulative payment: %w", err)
	}

	return newCumulativePayment, nil
}

// RevertDebit reverts a previous debit operation, following a failed dispersal.
func (odl *OnDemandLedger) RevertDebit(symbolCount uint32) error {
	if symbolCount == 0 {
		return errors.New("symbolCount must be > 0")
	}

	blobCost := odl.computeBlobCost(symbolCount)

	odl.lock.Lock()
	defer odl.lock.Unlock()

	currentCumulativePayment, err := odl.cumulativePaymentStore.GetCumulativePayment()
	if err != nil {
		return fmt.Errorf("get cumulative payment: %w", err)
	}

	newCumulativePayment := new(big.Int).Sub(currentCumulativePayment, blobCost)

	if newCumulativePayment.Sign() < 0 {
		return fmt.Errorf("cannot revert debit: would result in negative cumulative payment")
	}

	if err := odl.cumulativePaymentStore.SetCumulativePayment(newCumulativePayment); err != nil {
		return fmt.Errorf("set cumulative payment: %w", err)
	}

	return nil
}

func checkForOnDemandSupport(quorumsToCheck []core.QuorumID) error {
	for _, quorum := range quorumsToCheck {
		if quorum == 0 || quorum == 1 {
			continue
		}

		return fmt.Errorf("%w: quorum %d not in supported set [0, 1]", ErrQuorumNotSupported, quorum)
	}

	return nil
}

func (odl *OnDemandLedger) computeBlobCost(symbolCount uint32) *big.Int {
	symbolCountBig := big.NewInt(int64(symbolCount))

	billableSymbols := symbolCountBig
	if symbolCountBig.Cmp(odl.config.minNumSymbols) < 0 {
		billableSymbols = odl.config.minNumSymbols
	}

	return new(big.Int).Mul(billableSymbols, odl.config.pricePerSymbol)
}
