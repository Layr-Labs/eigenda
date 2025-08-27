package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	lru "github.com/hashicorp/golang-lru/v2"
)

// OnDemandPaymentValidator validates on-demand payments for multiple accounts
type OnDemandPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked.
	//
	// New OnDemandLedger entries are added to this cache as Debit requests are received from new accounts. Least
	// recently used OnDemandLedger entries are removed if the cache gets above the configured size. Since on-demand
	// payment data is stored in a persistent way, deleting an OnDemandLedger from memory doesn't result in data loss:
	// it just means that a new OnDemandLedger object will need to be constructed if any future Debits must be handled
	// from that account.
	ledgers *lru.Cache[gethcommon.Address, *OnDemandLedger]
	// protects concurrent access to the ledgers cache during ledger creation
	ledgerCreationLock sync.Mutex
	// Provides access to the values stored in the PaymentVault contract.
	//
	// The state of this object is updated on a background thread.
	onChainState meterer.OnchainPayment
	// the underlying dynamo client, which is used by all OnDemandLedger instances created by this struct
	dynamoClient *dynamodb.Client
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator with specified cache size
func NewOnDemandPaymentValidator(
	logger logging.Logger,
	// the maximum number of OnDemandLedger entries to be kept in the LRU cache
	maxLedgers int,
	// expected to be initialized and have its background update thread started
	onChainState meterer.OnchainPayment,
	dynamoClient *dynamodb.Client,
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string,
) (*OnDemandPaymentValidator, error) {
	if onChainState == nil {
		return nil, errors.New("onChainState cannot be nil")
	}
	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(accountAddress gethcommon.Address, _ *OnDemandLedger) {
			logger.Infof("evicted account %s from LRU on-demand ledger cache", accountAddress.Hex())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	if dynamoClient == nil {
		return nil, errors.New("dynamo client cannot be nil")
	}

	if onDemandTableName == "" {
		return nil, errors.New("on demand table name cannot be empty")
	}

	return &OnDemandPaymentValidator{
		logger:            logger,
		ledgers:           cache,
		onChainState:      onChainState,
		dynamoClient:      dynamoClient,
		onDemandTableName: onDemandTableName,
	}, nil
}

// Debit validates an on-demand payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (odl *OnDemandPaymentValidator) Debit(
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

// getOrCreateLedger gets an existing on-demand ledger from the cache, or creates a new one if it doesn't exist
func (odl *OnDemandPaymentValidator) getOrCreateLedger(
	ctx context.Context,
	accountID gethcommon.Address,
) (*OnDemandLedger, error) {
	// Fast path: check if ledger already exists in cache
	if ledger, exists := odl.ledgers.Get(accountID); exists {
		return ledger, nil
	}

	// Slow path: acquire lock and check again
	odl.ledgerCreationLock.Lock()
	defer odl.ledgerCreationLock.Unlock()

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
//
// Will attempt to make all updates, even if one update fails. A multierror is returned, describing any/all errors
// that occurred during the updates.
func (odl *OnDemandPaymentValidator) UpdateTotalDeposits(updates []TotalDepositUpdate) error {
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
				result, fmt.Errorf("update total deposit for account %v: %w", update.AccountAddress.Hex(), err))
			continue
		}
	}

	if err := result.ErrorOrNil(); err != nil {
		return fmt.Errorf("update total deposits: %w", err)
	}
	return nil
}
