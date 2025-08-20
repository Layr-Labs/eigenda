package ondemand

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandLedgers manages and validates on-demand payments for multiple accounts
type OnDemandLedgers struct {
	// map from account ID to OnDemandLedger
	ledgers map[gethcommon.Address]*OnDemandLedger

	// lock protects concurrent access to the ledgers map
	lock sync.Mutex
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator
func NewOnDemandPaymentValidator() *OnDemandLedgers {
	return &OnDemandLedgers{
		ledgers: make(map[gethcommon.Address]*OnDemandLedger),
	}
}

// Debit validates an on-demand payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (odl *OnDemandLedgers) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
) error {
	ledger, err := odl.getOrCreateLedger(accountID)
	if err != nil {
		return fmt.Errorf("get or create on-demand ledger: %w", err)
	}

	// For on-demand payments, call Debit with actual values
	_, err = ledger.Debit(ctx, symbolCount, quorumNumbers)
	if err != nil {
		return fmt.Errorf("debit on-demand payment: %w", err)
	}

	// TODO: Consider in what cases we should remove the ledger from the map
	// Possible cases:
	// - Account has exhausted all funds
	// - Account has been inactive for a certain period
	// - Explicit cleanup request

	return nil
}

// getOrCreateLedger gets an existing on-demand ledger or creates a new one if it doesn't exist
func (odl *OnDemandLedgers) getOrCreateLedger(accountID gethcommon.Address) (*OnDemandLedger, error) {
	odl.lock.Lock()
	defer odl.lock.Unlock()

	if ledger, exists := odl.ledgers[accountID]; exists {
		return ledger, nil
	}

	// TODO: These are placeholder values - need to get actual values from chain or config
	totalDeposits := big.NewInt(1000000000000000000) // 1 ETH placeholder
	pricePerSymbol := big.NewInt(100000000000000)    // 0.0001 ETH placeholder
	minNumSymbols := uint64(100)                     // 100 symbols placeholder

	// TODO: Need to provide actual CumulativePaymentStore implementation
	// For now, using nil as placeholder - this will need to be fixed
	newLedger, err := NewOnDemandLedger(
		totalDeposits,
		pricePerSymbol,
		minNumSymbols,
		nil, // placeholder for CumulativePaymentStore
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-demand ledger: %w", err)
	}

	odl.ledgers[accountID] = newLedger
	return newLedger, nil
}
