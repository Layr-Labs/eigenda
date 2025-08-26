package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	lru "github.com/hashicorp/golang-lru/v2"
)

// OnDemandLedgers manages and validates on-demand payments for multiple accounts
type OnDemandLedgers struct {
	logger       logging.Logger
	ledgers      *lru.Cache[gethcommon.Address, *OnDemandLedger]
	onChainState meterer.OnchainPayment

	dynamoClient      *dynamodb.Client
	onDemandTableName string
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator with specified cache size
func NewOnDemandPaymentValidator(
	logger logging.Logger,
	maxLedgers int,
	// expected to be initialized and have its background update thread started
	onChainState meterer.OnchainPayment,
	dynamoClient *dynamodb.Client,
	onDemandTableName string,
) (*OnDemandLedgers, error) {
	if onChainState == nil {
		return nil, errors.New("onChainState cannot be nil")
	}
	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(key gethcommon.Address, _ *OnDemandLedger) {
			logger.Infof("evicted account %s from LRU on-demand ledger cache", key.Hex())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	if dynamoClient == nil {
		return nil, errors.New("dynamo client cannot be nil")
	}

	if onDemandTableName == "" {
		return nil, errors.New("on demand table name cannot be nil")
	}

	return &OnDemandLedgers{
		logger:            logger,
		ledgers:           cache,
		onChainState:      onChainState,
		dynamoClient:      dynamoClient,
		onDemandTableName: onDemandTableName,
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

	onDemandPayment, err := odl.onChainState.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment for account %v: %w", accountID.Hex(), err)
	}

	pricePerSymbol := big.NewInt(int64(odl.onChainState.GetPricePerSymbol()))
	minNumSymbols := odl.onChainState.GetMinNumSymbols()

	cumulativePaymentStore, err := NewCumulativePaymentStore(odl.dynamoClient, odl.onDemandTableName, accountID)
	if err != nil {
		return nil, fmt.Errorf("new cumulative payment store: %w", err)
	}

	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		onDemandPayment.CumulativePayment,
		pricePerSymbol,
		minNumSymbols,
		cumulativePaymentStore,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-demand ledger: %w", err)
	}

	odl.ledgers.Add(accountID, newLedger)
	return newLedger, nil
}

// UpdateTotalDeposits updates the total deposits for multiple accounts.
func (odl *OnDemandLedgers) UpdateTotalDeposits(updates []TotalDepositUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	var result *multierror.Error
	for _, update := range updates {
		ledger, exists := odl.ledgers.Get(update.AccountAddress)
		if !exists {
			// if we aren't already tracking the account, there's nothing to do. we'll start tracking it if the
			// account ever makes an on-demand dispersal
			continue
		}

		err := ledger.UpdateTotalDeposits(update.NewTotalDeposit)
		if err != nil {
			result = multierror.Append(
				result, fmt.Errorf("failed to update deposits for account %v: %w", update.AccountAddress.Hex(), err))
			continue
		}
	}

	return fmt.Errorf("update total deposits: %w", result.ErrorOrNil())
}
