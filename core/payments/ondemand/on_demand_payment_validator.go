package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

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
	paymentVaultState *OnDemandPaymentVaultState

	// Background update configuration
	updateInterval time.Duration
	cancelFunc     context.CancelFunc
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
	paymentVaultParams PaymentVaultParams,
	// provides access to on-demand payment state and update notifications
	paymentVaultState *OnDemandPaymentVaultState,
	dynamoClient *dynamodb.Client,
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string,
	// interval for checking for payment updates
	updateInterval time.Duration,
) (*OnDemandPaymentValidator, error) {
	if paymentVaultState == nil {
		return nil, errors.New("paymentVaultState cannot be nil")
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

	return &OnDemandPaymentValidator{
		logger:             logger,
		ledgers:            cache,
		paymentVaultParams: paymentVaultParams,
		paymentVaultState:  paymentVaultState,
		updateInterval:     updateInterval,
		dynamoClient:       dynamoClient,
		onDemandTableName:  onDemandTableName,
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

	onDemandPayment, err := pv.paymentVaultState.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment for account %v: %w", accountID.Hex(), err)
	}

	cumulativePaymentStore, err := NewCumulativePaymentStore(pv.dynamoClient, pv.onDemandTableName, accountID)
	if err != nil {
		return nil, fmt.Errorf("new cumulative payment store: %w", err)
	}

	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		onDemandPayment.CumulativePayment,
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

// Start starts the background update thread
func (pv *OnDemandPaymentValidator) Start(ctx context.Context) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	pv.cancelFunc = cancel

	go pv.runUpdateLoop(ctxWithCancel)
}

// Stop stops the background update thread
func (pv *OnDemandPaymentValidator) Stop() {
	if pv.cancelFunc != nil {
		pv.cancelFunc()
	}
}

// Runs the background update loop, to periodically consume updates made to the PaymentVault
//
// TODO(litt3): Replace periodic polling with event-driven updates from PaymentVault contract
func (pv *OnDemandPaymentValidator) runUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(pv.updateInterval)
	defer ticker.Stop()

	pv.logger.Info("Starting OnDemandPaymentValidator background update thread", "updateInterval", pv.updateInterval)

	for {
		select {
		case <-ticker.C:
			if err := pv.performUpdates(ctx); err != nil {
				pv.logger.Error("perform on-demand payment updates", "error", err)
			}
		case <-ctx.Done():
			pv.logger.Info("OnDemandPaymentValidator background update thread stopped")
			return
		}
	}
}

// performUpdates fetches and applies updates that have been made to the payment vault
func (pv *OnDemandPaymentValidator) performUpdates(ctx context.Context) error {
	updates, err := pv.paymentVaultState.RefreshOnDemandPayments(ctx)
	if err != nil {
		return fmt.Errorf("refresh on-demand payments: %w", err)
	}

	var result *multierror.Error
	for _, update := range updates {
		ledger, exists := pv.ledgers.Get(update.AccountAddress)
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
