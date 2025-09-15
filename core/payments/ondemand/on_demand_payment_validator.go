package ondemand

import (
	"context"
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
}

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
	ledgerCache, err := NewOnDemandLedgerCache(
		ctx, logger, maxLedgers, paymentVault, updateInterval, dynamoClient, onDemandTableName)
	if err != nil {
		return nil, fmt.Errorf("new on-demand ledger cache: %w", err)
	}

	return &OnDemandPaymentValidator{
		logger:      logger,
		ledgerCache: ledgerCache,
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
