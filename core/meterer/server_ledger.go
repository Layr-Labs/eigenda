package meterer

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ServerLedger manages global rate limiting for all accounts
// It only tracks global on-demand bin usage and delegates account-specific
// validation to individual ServerAccountLedger instances
type ServerLedger struct {
	// Global dependencies
	meteringStore MeteringStore
	config        Config
	logger        logging.Logger

	// Cache of individual account ledgers
	accounts      map[gethcommon.Address]*ServerAccountLedger
	accountsMutex sync.RWMutex
}

// NewServerLedger creates a global accounts ledger manager
func NewServerLedger(
	config Config,
	meteringStore MeteringStore,
	logger logging.Logger,
) *ServerLedger {
	return &ServerLedger{
		meteringStore: meteringStore,
		config:        config,
		logger:        logger,
		accounts:      make(map[gethcommon.Address]*ServerAccountLedger),
	}
}

// GetAccountLedger retrieves or creates a ServerAccountLedger for a specific account
func (sal *ServerLedger) GetAccountLedger(
	ctx context.Context,
	accountID gethcommon.Address,
	chainPaymentState OnchainPayment,
) (*ServerAccountLedger, error) {
	// Try to get existing account ledger with read lock
	sal.accountsMutex.RLock()
	if ledger, exists := sal.accounts[accountID]; exists {
		sal.accountsMutex.RUnlock()
		return ledger, nil
	}
	sal.accountsMutex.RUnlock()

	// Create new account ledger with write lock
	sal.accountsMutex.Lock()
	defer sal.accountsMutex.Unlock()

	// Double-check in case another goroutine created it
	if ledger, exists := sal.accounts[accountID]; exists {
		return ledger, nil
	}

	// Create new account-specific ledger
	ledger, err := NewServerAccountLedger(
		ctx,
		accountID,
		chainPaymentState,
		sal.meteringStore,
		sal.config,
		sal.logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create account ledger for %s: %w", accountID.Hex(), err)
	}

	sal.accounts[accountID] = ledger
	sal.logger.Debug("Created new account ledger", "accountID", accountID.Hex())

	return ledger, nil
}

// ValidateGlobalOnDemandUsage validates global rate limits for on-demand usage
// This is the only global validation performed at this level
func (sal *ServerLedger) ValidateGlobalOnDemandUsage(
	ctx context.Context,
	timestampNs int64,
	numSymbols uint64,
	params *PaymentVaultParams,
) (uint64, error) {
	// Calculate the period for global bin
	period := payment_logic.GetReservationPeriodByNanosecond(timestampNs, 3600) // 1 hour window

	// Update global bin usage
	newGlobalUsage, err := sal.meteringStore.UpdateGlobalBin(ctx, period, numSymbols)
	if err != nil {
		return 0, fmt.Errorf("failed to update global bin: %w", err)
	}

	// Get global rate limit from first available quorum config
	var globalLimit uint64
	for _, quorumConfig := range params.QuorumPaymentConfigs {
		globalLimit = payment_logic.GetBinLimit(quorumConfig.OnDemandSymbolsPerSecond, 3600) // 1 hour window
		break
	}

	// Check global rate limit
	if newGlobalUsage > globalLimit {
		return newGlobalUsage, fmt.Errorf("global rate limit exceeded: %d > %d", newGlobalUsage, globalLimit)
	}

	return newGlobalUsage, nil
}

// RemoveAccount removes an account ledger from the cache
func (sal *ServerLedger) RemoveAccount(accountID gethcommon.Address) {
	sal.accountsMutex.Lock()
	defer sal.accountsMutex.Unlock()
	delete(sal.accounts, accountID)
	sal.logger.Debug("Removed account ledger", "accountID", accountID.Hex())
}

// GetAccountCount returns the number of cached account ledgers
func (sal *ServerLedger) GetAccountCount() int {
	sal.accountsMutex.RLock()
	defer sal.accountsMutex.RUnlock()
	return len(sal.accounts)
}
