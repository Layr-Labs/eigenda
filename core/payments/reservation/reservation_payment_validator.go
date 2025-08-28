package reservation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	lru "github.com/hashicorp/golang-lru/v2"
)

// ReservationPaymentValidator validates reservation payments for multiple accounts
type ReservationPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked.
	//
	// New ReservationLedger entries are added to this cache as Debit requests are received from new accounts. Least
	// recently used ReservationLedger entries are removed if the cache gets above the configured size.
	//
	// Since the LeakyBuckets the underlie the reservation ledgers are only in memory, evicting a ReservationLedger
	// from this cache can result in data loss! If a ledger is deleted and then re-added, the new instance is
	// instantiated *empty*, as if the user has had no recent dispersals. If the evicted ledger *wasn't* actually empty,
	// that means the user is not being limited to the capacity of their reservation. This errs on the side of
	// permitting more throughput, since being more permissive is preferable to the alternative of *not* providing the
	// amount of throughput guaranteed by a reservation.
	//
	// IMPORTANT: If the cache size is configured to be too small and there is a lot of churn, then
	// dishonest clients may be able to utilize more than their allotted reservations! Be sure to configure a large
	// enough cache. If any ReservationLedgers are removed from this cache with non-empty buckets, the occurrence will
	// be logged as an error. If such error logs are observed, this cache size must be increased.
	ledgers *lru.Cache[gethcommon.Address, *ReservationLedger]
	// protects concurrent access to the ledgers cache during ledger creation
	ledgerCreationLock sync.Mutex
	// Provides access to the values stored in the PaymentVault contract and update notifications
	paymentVaultState *ReservationPaymentVaultState

	// Background update configuration
	updateInterval time.Duration
	cancelFunc     context.CancelFunc
	// source of current time for the leaky bucket algorithm
	timeSource func() time.Time

	overfillBehavior     OverfillBehavior
	bucketCapacityPeriod time.Duration
}

// NewReservationPaymentValidator creates a new ReservationPaymentValidator with specified cache size
func NewReservationPaymentValidator(
	logger logging.Logger,
	// the maximum number of ReservationLedger entries to be kept in the LRU cache
	maxLedgers int,
	// provides access to reservation payment state and update notifications
	paymentVaultState *ReservationPaymentVaultState,
	// source of current time for the leaky bucket algorithm
	timeSource func() time.Time,
	// how to handle requests that would overfill the bucket
	overfillBehavior OverfillBehavior,
	// duration used to calculate bucket capacity
	bucketCapacityPeriod time.Duration,
	// interval for checking for payment updates
	updateInterval time.Duration,
) (*ReservationPaymentValidator, error) {
	if paymentVaultState == nil {
		return nil, errors.New("paymentVaultState cannot be nil")
	}

	if updateInterval <= 0 {
		return nil, errors.New("updateInterval must be > 0")
	}

	if bucketCapacityPeriod <= 0 {
		return nil, errors.New("bucketCapacityPeriod must be > 0")
	}

	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(key gethcommon.Address, reservationLedger *ReservationLedger) {
			// TODO: not threadsafe, fix it
			// TODO: would it make sense to double the cache size if this happens??
			if reservationLedger.leakyBucket.currentFillLevel > 0 {
				logger.Errorf("evicted account %s from LRU reservation ledger cache, but the underlying leaky bucket "+
					"wasn't empty! You must increase the ReservationPaymentValidator LRU cache size", key.Hex())
			} else {
				logger.Infof("evicted account %s from LRU reservation ledger cache", key.Hex())
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	return &ReservationPaymentValidator{
		logger:               logger,
		ledgers:              cache,
		paymentVaultState:    paymentVaultState,
		updateInterval:       updateInterval,
		timeSource:           timeSource,
		overfillBehavior:     overfillBehavior,
		bucketCapacityPeriod: bucketCapacityPeriod,
	}, nil
}

// Debit validates a reservation payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (pv *ReservationPaymentValidator) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) error {
	ledger, err := pv.getOrCreateLedger(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get or create reservation ledger: %w", err)
	}

	now := pv.timeSource()
	success, err := ledger.Debit(now, dispersalTime, symbolCount, quorumNumbers)
	if err != nil {
		return fmt.Errorf("debit reservation payment: %w", err)
	}

	if !success {
		return fmt.Errorf("reservation debit failed: insufficient capacity")
	}

	return nil
}

// getOrCreateLedger gets an existing reservation ledger or creates a new one if it doesn't exist
func (pv *ReservationPaymentValidator) getOrCreateLedger(
	ctx context.Context,
	accountID gethcommon.Address,
) (*ReservationLedger, error) {
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

	// Fetch on-chain reservation for account
	reservedPayment, err := pv.paymentVaultState.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reserved payment for account %v: %w", accountID.Hex(), err)
	}

	reservationObj, err := NewReservation(
		reservedPayment.SymbolsPerSecond,
		time.Unix(int64(reservedPayment.StartTimestamp), 0),
		time.Unix(int64(reservedPayment.EndTimestamp), 0),
		reservedPayment.QuorumNumbers)
	if err != nil {
		return nil, fmt.Errorf("new reservation: %w", err)
	}

	reservationLedgerConfig, err := NewReservationLedgerConfig(
		*reservationObj,
		// start empty, to err on the side of permitting more throughput instead of less
		false,
		pv.overfillBehavior,
		pv.bucketCapacityPeriod,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	now := pv.timeSource()
	newLedger, err := NewReservationLedger(*reservationLedgerConfig, now)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	pv.ledgers.Add(accountID, newLedger)
	return newLedger, nil
}

// Start starts the background update thread
func (pv *ReservationPaymentValidator) Start(ctx context.Context) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	pv.cancelFunc = cancel

	go pv.runUpdateLoop(ctxWithCancel)
}

// Stop stops the background update thread
func (pv *ReservationPaymentValidator) Stop() {
	if pv.cancelFunc != nil {
		pv.cancelFunc()
	}
}

// Runs the background update loop, to periodically consume updates made to the PaymentVault
//
// TODO(litt3): Replace periodic polling with event-driven updates from PaymentVault contract
func (pv *ReservationPaymentValidator) runUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(pv.updateInterval)
	defer ticker.Stop()

	pv.logger.Info("Starting ReservationPaymentValidator background update thread", "updateInterval", pv.updateInterval)

	for {
		select {
		case <-ticker.C:
			if err := pv.performUpdates(ctx); err != nil {
				pv.logger.Error("Failed to perform reservation payment updates", "error", err)
			}
		case <-ctx.Done():
			pv.logger.Info("ReservationPaymentValidator background update thread stopped")
			return
		}
	}
}

// performUpdates fetches and applies updates immediately as they are discovered
func (pv *ReservationPaymentValidator) performUpdates(ctx context.Context) error {
	updates, err := pv.paymentVaultState.RefreshReservedPayments(ctx)
	if err != nil {
		return fmt.Errorf("refresh reserved payments: %w", err)
	}

	now := pv.timeSource()
	var result *multierror.Error
	for _, update := range updates {
		ledger, exists := pv.ledgers.Get(update.AccountAddress)
		if !exists {
			// if we aren't already tracking the account, there's nothing to do. we'll start tracking it if the
			// account ever makes a reservation dispersal
			continue
		}

		err := ledger.UpdateReservation(update.NewReservation, now)
		if err != nil {
			result = multierror.Append(
				result, fmt.Errorf("update reservation for account %v: %w", update.AccountAddress.Hex(), err))
			continue
		}
	}

	if err := result.ErrorOrNil(); err != nil {
		return fmt.Errorf("update reservations: %w", err)
	}
	return nil
}
