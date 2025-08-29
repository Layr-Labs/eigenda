package ondemand

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandPaymentValidator validates on-demand payments for multiple accounts
type OnDemandPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked
	ledgerCache *OnDemandLedgerCache
	// Provides access to the values stored in the PaymentVault contract and update notifications
	vaultMonitor *OnDemandVaultMonitor
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator with specified cache size
func NewOnDemandPaymentValidator(
	ctx context.Context,
	logger logging.Logger,
	// the maximum number of OnDemandLedger entries to be kept in the LRU cache
	maxLedgers int,
	// provides access to payment vault contract
	paymentVault payments.PaymentVault,
	dynamoClient *dynamodb.Client,
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string,
	// interval for checking for payment updates
	updateInterval time.Duration,
) (*OnDemandPaymentValidator, error) {
	if paymentVault == nil {
		return nil, errors.New("paymentVault cannot be nil")
	}

	if updateInterval <= 0 {
		return nil, errors.New("updateInterval must be > 0")
	}

	ledgerCache, err := NewOnDemandLedgerCache(logger, maxLedgers, paymentVault, dynamoClient, onDemandTableName)
	if err != nil {
		return nil, fmt.Errorf("new on-demand ledger cache: %w", err)
	}

	if dynamoClient == nil {
		return nil, errors.New("dynamo client cannot be nil")
	}

	if onDemandTableName == "" {
		return nil, errors.New("on demand table name cannot be empty")
	}

	vaultMonitor := NewOnDemandVaultMonitor(ctx, logger, paymentVault, ledgerCache, updateInterval)

	return &OnDemandPaymentValidator{
		logger:       logger,
		ledgerCache:  ledgerCache,
		vaultMonitor: vaultMonitor,
	}, nil
}

// Debit validates an on-demand payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (pv *OnDemandPaymentValidator) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
) error {
	ledger, err := pv.ledgerCache.GetOrCreate(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get or create ledger: %w", err)
	}

	_, err = ledger.Debit(ctx, symbolCount, quorumNumbers)
	if err != nil {
		return fmt.Errorf("debit on-demand payment: %w", err)
	}

	return nil
}

// Stop stops the background vault monitoring thread
func (pv *OnDemandPaymentValidator) Stop() {
	if pv.vaultMonitor != nil {
		pv.vaultMonitor.Stop()
	}
}
