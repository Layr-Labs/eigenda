package reservation

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Stores a collection of ReservationLedgers in an LRU cache
type ReservationLedgerCache struct {
	logger logging.Logger
	// A cache of the ledgers being tracked.
	//
	// Least recently used ReservationLedger entries are removed if the cache gets above the configured size. Since
	// the LeakyBuckets that underlie the reservation ledgers are only in memory, evicting a ReservationLedger
	// from this cache can result in data loss! If a ledger is deleted and then re-added, the new instance is
	// instantiated *empty*, as if the user has had no recent dispersals. This errs on the side of permitting more
	// throughput, since being more permissive is preferable to the alternative of *not* providing the amount of
	// throughput guaranteed by a reservation.
	//
	// IMPORTANT: If the cache size is configured to be too small and there is a lot of churn, then
	// dishonest clients may be able to utilize more than their allotted reservations! Be sure to configure a large
	// enough cache. If any ReservationLedgers are removed from this cache with non-empty buckets, the occurrence will
	// be logged as an error. If such error logs are observed, this cache size must be increased.
	cache *lru.Cache[gethcommon.Address, *ReservationLedger]
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

	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(accountAddress gethcommon.Address, reservationLedger *ReservationLedger) {
			isEmpty, err := reservationLedger.IsBucketEmpty(timeSource())
			if err != nil {
				logger.Errorf("failed to check if bucket is empty for account %s: %v", accountAddress.Hex(), err)
			}

			if !isEmpty {
				logger.Errorf("evicted account %s from LRU reservation ledger cache, but the underlying leaky bucket "+
					"wasn't empty! You must increase the ReservationLedgerCache LRU cache size", accountAddress.Hex())
				return
			}

			logger.Infof("evicted account %s from LRU reservation ledger cache", accountAddress.Hex())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	ledgerCache := &ReservationLedgerCache{
		logger:               logger,
		cache:                cache,
		paymentVault:         paymentVault,
		timeSource:           timeSource,
		overfillBehavior:     overfillBehavior,
		bucketCapacityPeriod: bucketCapacityPeriod,
		minNumSymbols:        minNumSymbols,
		ledgerCreationLock:   common.NewIndexLock(256),
	}

	// Create the vault monitor with callback functions
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
	accountIndex := binary.BigEndian.Uint64(accountID.Bytes()[:8])
	c.ledgerCreationLock.Lock(accountIndex)
	defer c.ledgerCreationLock.Unlock(accountIndex)

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
