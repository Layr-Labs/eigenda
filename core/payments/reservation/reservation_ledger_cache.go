package reservation

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	// maxReservationLRUCacheSize is the maximum number of reservation ledgers that can be stored in the cache.
	// Set to 2^16 = 65,536 entries.
	//
	// To do some napkin math: each cache entry is <500 bytes in size, so 65k cache entries would have a memory
	// footprint <33MiB. This isn't a catastrophic amount of memory, and 65k active reservation users is absurdly high.
	maxReservationLRUCacheSize = 65536
)

// Stores a collection of ReservationLedgers in an LRU cache
type ReservationLedgerCache struct {
	logger logging.Logger
	// A cache of the ledgers being tracked.
	//
	// Least recently used ReservationLedger entries are removed if the cache gets above the configured size.
	//
	// The LeakyBuckets that underlie the reservation ledgers are *only* in memory. This means that evicting a ledger
	// prematurely from the cache (when the LeakyBucket isn't empty) results in information loss! If the prematurely
	// evicted ledger were to be reinstantiated, it would start with an *empty* bucket, potentially permitting more
	// throughput than it should (assuming a malicious client).
	//
	// The solution to prevent this from happening is that we will detect when a ledger is evicted prematurely, and
	// automatically resize the cache in response. This prevents the cache from getting into a thrashy state, where
	// many ledgers are being evicted prematurely and then reinstantiated.
	cache *lru.Cache[gethcommon.Address, *ReservationLedger]
	// current maximum number of ledgers the cache can hold (will be dynamically increased if premature evictions are
	// observed)
	maxLedgers int
	// can access state of the PaymentVault contract
	paymentVault payments.PaymentVault
	// source of current time for the leaky bucket algorithm
	timeSource func() time.Time
	// how to handle requests that would overfill the bucket
	overfillBehavior OverfillBehavior
	// duration used to calculate bucket capacity
	bucketCapacityPeriod time.Duration
	// minimum number of symbols to bill, from the PaymentVault
	minNumSymbols uint32
	// protects concurrent access to the ledgers cache during ledger creation
	//
	// The lru.Cache object itself is threadsafe, as are the ReservationLedger values contained in the cache. This lock
	// is to make sure that only one caller is constructing a new ReservationLedger at a time for a specific account.
	// Otherwise, it would be possible for two separate callers to get a cache miss for the same account, create the
	// new object for the same account key, and try to add them to the cache.
	ledgerCreationLock *common.IndexLock
	// protects the cache eviction process, ensures that only one eviction can be processed at a time and preventing
	// race conditions during cache resizing
	evictionLock sync.Mutex
	// monitors the PaymentVault for changes, and updates cached ledgers accordingly
	vaultMonitor *ReservationVaultMonitor
}

func NewReservationLedgerCache(
	ctx context.Context,
	logger logging.Logger,
	maxLedgers int,
	paymentVault payments.PaymentVault,
	timeSource func() time.Time,
	overfillBehavior OverfillBehavior,
	bucketCapacityPeriod time.Duration,
	updateInterval time.Duration,
) (*ReservationLedgerCache, error) {
	if paymentVault == nil {
		return nil, errors.New("payment vault must be non-nil")
	}

	if timeSource == nil {
		return nil, errors.New("time source must be non-nil")
	}

	if bucketCapacityPeriod <= 0 {
		return nil, errors.New("bucket capacity period must be > 0")
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	ledgerCache := &ReservationLedgerCache{
		logger:               logger,
		maxLedgers:           maxLedgers,
		paymentVault:         paymentVault,
		timeSource:           timeSource,
		overfillBehavior:     overfillBehavior,
		bucketCapacityPeriod: bucketCapacityPeriod,
		minNumSymbols:        minNumSymbols,
		ledgerCreationLock:   common.NewIndexLock(256),
	}

	ledgerCache.cache, err = lru.NewWithEvict(maxLedgers, ledgerCache.handleEviction)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	ledgerCache.vaultMonitor, err = NewReservationVaultMonitor(
		ctx,
		logger,
		paymentVault,
		updateInterval,
		// relatively arbitrary value. much higher than account number in practice, but much lower than what the RPC
		// could actually handle. Since the "sweet spot" is really wide, hardcode this instead of spending time wiring
		// in a config value
		1024,
		ledgerCache.GetAccountsToUpdate,
		ledgerCache.UpdateReservation,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation vault monitor: %w", err)
	}

	return ledgerCache, nil
}

// GetOrCreate retrieves an existing ReservationLedger for the given account, or creates a new one if it doesn't exist
func (c *ReservationLedgerCache) GetOrCreate(
	ctx context.Context,
	accountID gethcommon.Address,
) (*ReservationLedger, error) {
	// Fast path: check if ledger already exists in cache
	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	// Slow path: acquire per-account lock and check again
	defer c.acquireLedgerLock(accountID)()

	if ledger, exists := c.cache.Get(accountID); exists {
		return ledger, nil
	}

	reservationData, err := c.paymentVault.GetReservation(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reservation for account %v: %w", accountID.Hex(), err)
	}

	if reservationData == nil {
		return nil, fmt.Errorf("no reservation found for account %v", accountID.Hex())
	}

	reservationObj, err := FromContractStruct(reservationData)
	if err != nil {
		return nil, fmt.Errorf("from contract struct: %w", err)
	}

	reservationLedgerConfig, err := NewReservationLedgerConfig(
		*reservationObj,
		c.minNumSymbols,
		// start empty, to err on the side of permitting more throughput instead of less
		false,
		c.overfillBehavior,
		c.bucketCapacityPeriod,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	now := c.timeSource()
	newLedger, err := NewReservationLedger(*reservationLedgerConfig, now)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	c.cache.Add(accountID, newLedger)
	return newLedger, nil
}

// GetAccountsToUpdate returns all accounts currently being tracked in the cache
func (c *ReservationLedgerCache) GetAccountsToUpdate() []gethcommon.Address {
	return c.cache.Keys()
}

// UpdateReservation updates the reservation for an account if different from current value
// If newReservation is nil, the account is removed from the cache
func (c *ReservationLedgerCache) UpdateReservation(accountID gethcommon.Address, newReservation *Reservation) error {
	ledger, exists := c.cache.Get(accountID)
	if !exists {
		// Account was evicted from cache or never existed, nothing to update
		return nil
	}

	if newReservation == nil {
		c.cache.Remove(accountID)
		c.logger.Debugf("Removed account %s from cache due to nil reservation", accountID.Hex())
		return nil
	}

	now := c.timeSource()
	return ledger.UpdateReservation(newReservation, now)
}

// Called when an item is evicted from the LRU cache.
//
// If the evicted ledger has a non-empty bucket, it resizes the cache and re-adds the ledger.
func (c *ReservationLedgerCache) handleEviction(
	accountID gethcommon.Address,
	reservationLedger *ReservationLedger,
) {
	c.evictionLock.Lock()
	defer c.evictionLock.Unlock()

	if reservationLedger.IsBucketEmpty(c.timeSource()) {
		c.logger.Debugf("evicted account %s from LRU reservation ledger cache", accountID.Hex())
		return
	}

	// The bucket is not empty!!! This was a premature eviction: we must resize the cache

	newSize := c.maxLedgers * 2
	if newSize > maxReservationLRUCacheSize {
		c.logger.Errorf(
			"Cannot resize LRU reservation ledger cache beyond maximum size of %d entries. Current size: %d",
			maxReservationLRUCacheSize, c.maxLedgers)
		// We've hit the maximum cache size - still evict the entry but don't resize
		return
	}

	c.logger.Infof("Resizing LRU reservation ledger cache from %d to %d entries.", c.maxLedgers, newSize)

	c.maxLedgers = newSize
	c.cache.Resize(c.maxLedgers)

	// Don't bother checking if another routine already re-created this ledger. Even if another routine *did* create
	// a new instance, it's reasonable to preference the old instance over the new. There may be some small discrepancy
	// here, but there would be no feasible way for a malicious client to exploit this. In the worst case, the leaky
	// bucket will be slightly less filled than it ought to have been. Since it's incredibly unlikely to happen in the
	// first place, it's not worth contorting the design to address.
	c.cache.Add(accountID, reservationLedger)
}

// Acquires the per-account lock for the given account address and returns a function that should be called to release
// the lock via defer
func (c *ReservationLedgerCache) acquireLedgerLock(accountID gethcommon.Address) func() {
	accountIndex := binary.BigEndian.Uint64(accountID.Bytes()[:8])
	c.ledgerCreationLock.Lock(accountIndex)
	return func() {
		c.ledgerCreationLock.Unlock(accountIndex)
	}
}
