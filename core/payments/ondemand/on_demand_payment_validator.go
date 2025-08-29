package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
	//
	// The lru.Cache object itself is threadsafe, as are the OnDemandLedger values contained in the cache. This lock
	// is to make sure that only one caller is constructing a new OnDemandLedger at a time. Otherwise, it would be
	// possible for two separate callers to get a cache miss, create the new object for the same account key, and try
	// to add them to the cache.
	ledgerCreationLock sync.Mutex
	// global payment parameters from the PaymentVault
	//
	// TODO(litt3): there is currently no consideration for updates that may be made to these parameters. The strategy
	// used by the old metering logic wasn't actually safe: updates to the global payment params must be made
	// deterministically based on RBN. This logic should be implemented before updates to these parameters are made
	// on-chain.
	paymentVaultParams PaymentVaultParams
	// Provides access to the values stored in the PaymentVault contract and update notifications
	onDemandPaymentVault *OnDemandVaultMonitor
	// the underlying dynamo client, which is used by all OnDemandLedger instances created by this struct
	dynamoClient *dynamodb.Client
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string
}

// NewOnDemandPaymentValidator creates a new OnDemandPaymentValidator with specified cache size
func NewOnDemandPaymentValidator(
	ctx context.Context,
	logger logging.Logger,
	// the maximum number of OnDemandLedger entries to be kept in the LRU cache
	maxLedgers int,
	paymentVaultParams PaymentVaultParams,
	// provides access to payment vault contract
	paymentVault *vault.PaymentVault,
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

	validator := &OnDemandPaymentValidator{
		logger:             logger,
		ledgers:            cache,
		paymentVaultParams: paymentVaultParams,
		dynamoClient:       dynamoClient,
		onDemandTableName:  onDemandTableName,
	}

	// Create callback for this validator instance
	callback := validator.CreateUpdateCallback()

	// Create vault state with the callback - it will auto-start background monitoring
	onDemandPaymentVault := NewOnDemandVaultMonitor(ctx, logger, paymentVault, callback, updateInterval)
	validator.onDemandPaymentVault = onDemandPaymentVault

	return validator, nil
}

// CreateUpdateCallback creates a callback function for handling payment updates
// This method is exported so that external code can create vault states with the proper callback
func (pv *OnDemandPaymentValidator) CreateUpdateCallback() onDemandPaymentUpdateCallback {
	return func(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
		ledger, exists := pv.ledgers.Get(accountID)
		if !exists {
			// if we aren't already tracking the account, there's nothing to do. we'll start tracking it if the
			// account ever makes an on-demand dispersal
			return nil
		}

		err := ledger.UpdateTotalDeposits(newTotalDeposit)
		if err != nil {
			return fmt.Errorf("update total deposit for account %v: %w", accountID.Hex(), err)
		}
		return nil
	}
}

// Debit validates an on-demand payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (pv *OnDemandPaymentValidator) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
) error {
	ledger, err := pv.getOrCreateLedger(ctx, accountID)
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
func (pv *OnDemandPaymentValidator) getOrCreateLedger(
	ctx context.Context,
	accountID gethcommon.Address,
) (*OnDemandLedger, error) {
	// Fast path: check if ledger already exists in cache
	if ledger, exists := pv.ledgers.Get(accountID); exists {
		return ledger, nil
	}

	// Slow path: acquire lock and check again
	pv.ledgerCreationLock.Lock()
	defer pv.ledgerCreationLock.Unlock()

	if ledger, exists := pv.ledgers.Get(accountID); exists {
		return ledger, nil
	}

	onDemandPayment, err := pv.onDemandPaymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment for account %v: %w", accountID.Hex(), err)
	}

	cumulativePaymentStore, err := NewCumulativePaymentStore(pv.dynamoClient, pv.onDemandTableName, accountID)
	if err != nil {
		return nil, fmt.Errorf("new cumulative payment store: %w", err)
	}

	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		onDemandPayment,
		big.NewInt(int64(pv.paymentVaultParams.PricePerSymbol)),
		pv.paymentVaultParams.MinNumSymbols,
		cumulativePaymentStore,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-demand ledger: %w", err)
	}

	pv.ledgers.Add(accountID, newLedger)
	return newLedger, nil
}
