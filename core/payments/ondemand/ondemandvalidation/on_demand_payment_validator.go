package ondemandvalidation

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandPaymentValidator validates on-demand payments for multiple accounts
type OnDemandPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked
	ledgerCache *OnDemandLedgerCache
	metrics     *OnDemandValidatorMetrics
}

func NewOnDemandPaymentValidator(
	ctx context.Context,
	logger logging.Logger,
	config OnDemandLedgerCacheConfig,
	// provides access to payment vault contract
	paymentVault payments.PaymentVault,
	dynamoClient *dynamodb.Client,
	validatorMetrics *OnDemandValidatorMetrics,
	cacheMetrics *OnDemandCacheMetrics,
) (*OnDemandPaymentValidator, error) {
	ledgerCache, err := NewOnDemandLedgerCache(ctx, logger, config, paymentVault, dynamoClient, cacheMetrics)
	if err != nil {
		return nil, fmt.Errorf("new on-demand ledger cache: %w", err)
	}

	return &OnDemandPaymentValidator{
		logger:      logger,
		ledgerCache: ledgerCache,
		metrics:     validatorMetrics,
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
	if err == nil {
		pv.metrics.RecordSuccess(accountID.Hex(), symbolCount)
		return nil
	}

	var insufficientFundsErr *ondemand.InsufficientFundsError
	if errors.As(err, &insufficientFundsErr) {
		pv.metrics.IncrementInsufficientFunds()
		return err
	}

	var quorumNotSupportedErr *ondemand.QuorumNotSupportedError
	if errors.As(err, &quorumNotSupportedErr) {
		pv.metrics.IncrementQuorumNotSupported()
		return err
	}

	pv.metrics.IncrementUnexpectedErrors()
	return err
}
