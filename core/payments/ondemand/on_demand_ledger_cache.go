package ondemand

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Stores a collection of OnDemandLedgers in an LRU cache
//
// The OnDemandLedgers created and stored in this cache are backed by DynamoDB, so that on-demand payment usage is
// persistent.
type OnDemandLedgerCache struct {
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
	minNumSymbols uint32
	// protects concurrent access to the ledgers cache during ledger creation
	//
	// The lru.Cache object itself is threadsafe, as are the OnDemandLedger values contained in the cache. This lock
	// is to make sure that only one caller is constructing a new OnDemandLedger at a time for a specific account.
	// Otherwise, it would be possible for two separate callers to get a cache miss for the same account, create the
	// new object for the same account key, and try to add them to the cache.
	ledgerCreationLock *common.IndexLock
	// monitors the PaymentVault for changes, and updates cached ledgers accordingly
	vaultMonitor *OnDemandVaultMonitor
}

func NewOnDemandLedgerCache(
	ctx context.Context,
	logger logging.Logger,
	maxLedgers int,
	paymentVault payments.PaymentVault,
	updateInterval time.Duration,
	dynamoClient *dynamodb.Client,
	onDemandTableName string,
) (*OnDemandLedgerCache, error) {
	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(accountAddress gethcommon.Address, _ *OnDemandLedger) {
			logger.Infof("evicted account %s from LRU on-demand ledger cache", accountAddress.Hex())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	if paymentVault == nil {
		return nil, errors.New("payment vault must be non-nil")
	}

	if dynamoClient == nil {
		return nil, errors.New("dynamo client must be non-nil")
	}

	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	ledgerCache := &OnDemandLedgerCache{
		cache:              cache,
		paymentVault:       paymentVault,
		dynamoClient:       dynamoClient,
		onDemandTableName:  onDemandTableName,
		pricePerSymbol:     new(big.Int).SetUint64(pricePerSymbol),
		minNumSymbols:      minNumSymbols,
		ledgerCreationLock: common.NewIndexLock(256),
	}

	// Create the vault monitor with callback functions
	ledgerCache.vaultMonitor, err = NewOnDemandVaultMonitor(
		ctx,
		logger,
		paymentVault,
		updateInterval,
		// relatively arbitrary value. much higher than account number in practice, but much lower than what the RPC
		// could actually handle. Since the "sweet spot" is really wide, hardcode this instead of spending time wiring
		// in a config value
		1024,
		ledgerCache.GetAccountsToUpdate,
		ledgerCache.UpdateTotalDeposit,
	)
	if err != nil {
		return nil, fmt.Errorf("new on-demand vault monitor: %w", err)
	}

	return ledgerCache, nil
}

// Retrieves an existing OnDemandLedger for the given account, or creates a new one if it doesn't exist
//
// Note: there exists a potential race condition with the access pattern of this method:
// 1. A ledger is retrieved from the cache
// 2. A large amount of activity (or a small configured cache size) causes the ledger to be evicted from the cache
// before the ledger operation has been completed
// 3. A different caller tries to retrieve the ledger for that account, gets a cache miss, and constructs a new instance
//
// With this sequence of events, there could be multiple existing ledger instances for the same account. The
// underlying cumulative payment store isn't designed to function with multiple instantiated ledger structs, so the
// operation of one instance would overwrite the operation of the other. Practically, this would mean that the user
// would get one free dispersal. The multiple instance problem would resolve itself after a single operation, since
// the LRU cache can only maintain a single instance, and the other instance would be destroyed.
//
// It is very unlikely for this race condition to take place if the cache has been configured with a sane size. Given
// the low probability of the occurrence, and the low severity of the race condition, we are not addressing it right
// now to avoid the complexity of the potential workarounds.
func (c *OnDemandLedgerCache) GetOrCreate(ctx context.Context, accountID gethcommon.Address) (*OnDemandLedger, error) {
	// Fast path: check if ledger already exists in cache
	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	// Slow path: acquire per-account lock and check again
	accountIndex := binary.BigEndian.Uint64(accountID.Bytes()[:8])
	c.ledgerCreationLock.Lock(accountIndex)
	defer c.ledgerCreationLock.Unlock(accountIndex)

	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	totalDeposit, err := c.paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get total deposit for account %v: %w", accountID.Hex(), err)
	}

	cumulativePaymentStore, err := NewCumulativePaymentStore(c.dynamoClient, c.onDemandTableName, accountID)
	if err != nil {
		return nil, fmt.Errorf("new cumulative payment store: %w", err)
	}

	newLedger, err := OnDemandLedgerFromStore(
		ctx,
		totalDeposit,
		c.pricePerSymbol,
		c.minNumSymbols,
		cumulativePaymentStore,
	)
	if err != nil {
		return nil, fmt.Errorf("create ledger from store: %w", err)
	}

	c.cache.Add(accountID, newLedger)
	return newLedger, nil
}

// Returns all accounts currently being tracked in the cache
//
// This method is used to determine which values need to be fetched from the PaymentVault, when periodically
// checking for updates.
func (c *OnDemandLedgerCache) GetAccountsToUpdate() []gethcommon.Address {
	return c.cache.Keys()
}

// Updates the total deposit for an account
func (c *OnDemandLedgerCache) UpdateTotalDeposit(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
	ledger, exists := c.cache.Get(accountID)
	if !exists {
		// Account was evicted from cache, nothing to update
		return nil
	}

	currentDeposit := ledger.GetTotalDeposits()
	if currentDeposit.Cmp(newTotalDeposit) != 0 {
		return ledger.UpdateTotalDeposits(newTotalDeposit)
	}
	return nil
}
