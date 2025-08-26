package ondemand

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// OnDemandLedgers manages and validates on-demand payments for multiple accounts
type OnDemandLedgers struct {
	logger  logging.Logger
	ledgers *lru.Cache[gethcommon.Address, *OnDemandLedger]
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator with specified cache size
func NewOnDemandPaymentValidator(logger logging.Logger, maxLedgers int) (*OnDemandLedgers, error) {
	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(key gethcommon.Address, _ *OnDemandLedger) {
			logger.Infof("purged account %v from LRU on-demand ledger cache", key)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	return &OnDemandLedgers{
		logger:  logger,
		ledgers: cache,
	}, nil
}

// Debit validates an on-demand payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (odl *OnDemandLedgers) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
) error {
	ledger, err := odl.getOrCreateLedger(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get or create on-demand ledger: %w", err)
	}

	_, err = ledger.Debit(ctx, symbolCount, quorumNumbers)
	if err != nil {
		return fmt.Errorf("debit on-demand payment: %w", err)
	}

	return nil
}

// getOrCreateLedger gets an existing on-demand ledger or creates a new one if it doesn't exist
func (odl *OnDemandLedgers) getOrCreateLedger(
	ctx context.Context,
	accountID gethcommon.Address,
) (*OnDemandLedger, error) {
	if ledger, exists := odl.ledgers.Get(accountID); exists {
		return ledger, nil
	}

	// TODO: These are placeholder values - need to get actual values from chain or config
	totalDeposits := big.NewInt(1000000000000000000) // 1 ETH placeholder
	pricePerSymbol := big.NewInt(100000000000000)    // 0.0001 ETH placeholder
	minNumSymbols := uint64(100)                     // 100 symbols placeholder

	// TODO: Need to provide actual CumulativePaymentStore implementation
	// For now, using nil as placeholder - this will need to be fixed
	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		totalDeposits,
		pricePerSymbol,
		minNumSymbols,
		nil, // placeholder for CumulativePaymentStore
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-demand ledger: %w", err)
	}

	odl.ledgers.Add(accountID, newLedger)
	return newLedger, nil
}
