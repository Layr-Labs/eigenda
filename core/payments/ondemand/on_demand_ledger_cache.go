package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Stores a collection of OnDemandLedgers in an LRU cache
type OnDemandLedgerCache struct {
	logger logging.Logger
	// A cache of the ledgers being tracked.
	//
	// Least recently used OnDemandLedger entries are removed if the cache gets above the configured size. Since
	// on-demand payment data is stored in a persistent way, deleting an OnDemandLedger from memory doesn't result in
	// data loss: it just means that a new OnDemandLedger object will need to be constructed if needed in the future.
	cache *lru.Cache[gethcommon.Address, *OnDemandLedger]
	// can access state of the PaymentVault contract
	paymentVault payments.PaymentVault
	// the underlying dynamo client, which is used by all OnDemandLedger instances created by this struct
	dynamoClient *dynamodb.Client
	// the name of the dynamo table where on-demand payment information is stored
	onDemandTableName string
	// price per symbol in wei, from the PaymentVault
	pricePerSymbol *big.Int
	// minimum number of symbols to bill, from the PaymentVault
	minNumSymbols uint64
	// protects concurrent access to the ledgers cache during ledger creation
	//
	// The lru.Cache object itself is threadsafe, as are the OnDemandLedger values contained in the cache. This lock
	// is to make sure that only one caller is constructing a new OnDemandLedger at a time. Otherwise, it would be
	// possible for two separate callers to get a cache miss, create the new object for the same account key, and try
	// to add them to the cache.
	//
	// This lock is intentionally more restrictive than it needs to be, for the sake of simplicity. It could be
	// converted to an account-based lock instead of a global creation lock, if it ever becomes a bottleneck.
	ledgerCreationLock sync.Mutex
}

var _ UpdatableOnDemandLedgers = &OnDemandLedgerCache{}

func NewOnDemandLedgerCache(
	ctx context.Context,
	logger logging.Logger,
	maxLedgers int,
	paymentVault payments.PaymentVault,
	dynamoClient *dynamodb.Client,
	onDemandTableName string,
) (*OnDemandLedgerCache, error) {
	if paymentVault == nil {
		return nil, errors.New("payment vault must be non-nil")
	}

	if dynamoClient == nil {
		return nil, errors.New("dynamo client must be non-nil")
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

	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	return &OnDemandLedgerCache{
		logger:            logger,
		cache:             cache,
		paymentVault:      paymentVault,
		dynamoClient:      dynamoClient,
		onDemandTableName: onDemandTableName,
		pricePerSymbol:    new(big.Int).SetUint64(pricePerSymbol),
		minNumSymbols:     minNumSymbols,
	}, nil
}

// GetOrCreate retrieves an existing OnDemandLedger for the given account, or creates a new one if it doesn't exist
func (c *OnDemandLedgerCache) GetOrCreate(ctx context.Context, accountID gethcommon.Address) (*OnDemandLedger, error) {
	// Fast path: check if ledger already exists in cache
	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	// Slow path: acquire lock and check again
	c.ledgerCreationLock.Lock()
	defer c.ledgerCreationLock.Unlock()

	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	onDemandPayment, err := c.paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment for account %v: %w", accountID.Hex(), err)
	}

	cumulativePaymentStore, err := NewCumulativePaymentStore(c.dynamoClient, c.onDemandTableName, accountID)
	if err != nil {
		return nil, fmt.Errorf("new cumulative payment store: %w", err)
	}

	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		onDemandPayment,
		c.pricePerSymbol,
		c.minNumSymbols,
		cumulativePaymentStore,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-demand ledger: %w", err)
	}

	c.cache.Add(accountID, newLedger)
	return newLedger, nil
}

// GetAccountsToUpdate returns all accounts currently being tracked in the cache
func (c *OnDemandLedgerCache) GetAccountsToUpdate() []gethcommon.Address {
	return c.cache.Keys()
}

// UpdateTotalDeposit updates the total deposit for an account if different from current value
func (c *OnDemandLedgerCache) UpdateTotalDeposit(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
	ledger, exists := c.cache.Get(accountID)
	if !exists {
		// Account was evicted from cache, nothing to update
		return nil
	}

	return ledger.UpdateTotalDeposits(newTotalDeposit)
}
