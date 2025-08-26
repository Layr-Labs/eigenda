package reservation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// ReservationLedgers manages and validates reservation payments for multiple accounts
type ReservationLedgers struct {
	logger       logging.Logger
	ledgers      *lru.Cache[gethcommon.Address, *ReservationLedger]
	onChainState meterer.OnchainPayment

	timeSource func() time.Time

	// Configuration for constructing underlying ReservationLedger instances
	startFull            bool
	overfillBehavior     OverfillBehavior
	bucketCapacityPeriod time.Duration
}

// NewReservationPaymentValidator creates a new ReservationLedgers with specified cache size and on-chain reader
func NewReservationPaymentValidator(
	logger logging.Logger,
	maxLedgers int,
	// expected to be initialized and have its background update thread started
	onChainState meterer.OnchainPayment,
	timeSource func() time.Time,
	startFull bool,
	overfillBehavior OverfillBehavior,
	bucketCapacityPeriod time.Duration,
) (*ReservationLedgers, error) {
	if onChainState == nil {
		return nil, errors.New("onChainState cannot be nil")
	}
	if maxLedgers <= 0 {
		return nil, errors.New("maxLedgers must be > 0")
	}
	if bucketCapacityPeriod <= 0 {
		return nil, errors.New("bucketCapacityPeriod must be > 0")
	}

	cache, err := lru.NewWithEvict(
		maxLedgers,
		func(key gethcommon.Address, _ *ReservationLedger) {
			logger.Infof("evicted account %s from LRU reservation ledger cache", key.Hex())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new LRU cache with evict: %w", err)
	}

	return &ReservationLedgers{
		logger:               logger,
		ledgers:              cache,
		onChainState:         onChainState,
		timeSource:           timeSource,
		startFull:            startFull,
		overfillBehavior:     overfillBehavior,
		bucketCapacityPeriod: bucketCapacityPeriod,
	}, nil
}

// Debit validates a reservation payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (rl *ReservationLedgers) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) error {
	ledger, err := rl.getOrCreateLedger(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get or create reservation ledger: %w", err)
	}

	now := rl.timeSource()
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
func (rl *ReservationLedgers) getOrCreateLedger(
	ctx context.Context,
	accountID gethcommon.Address,
) (*ReservationLedger, error) {
	if ledger, exists := rl.ledgers.Get(accountID); exists {
		return ledger, nil
	}

	// Fetch on-chain reservation for account
	reservedPayment, err := rl.onChainState.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reserved payment for account %v: %w", accountID.Hex(), err)
	}

	// Build Reservation object
	startTime := time.Unix(int64(reservedPayment.StartTimestamp), 0)
	endTime := time.Unix(int64(reservedPayment.EndTimestamp), 0)

	reservationObj, err := NewReservation(
		reservedPayment.SymbolsPerSecond,
		startTime,
		endTime,
		reservedPayment.QuorumNumbers)
	if err != nil {
		return nil, fmt.Errorf("new reservation: %w", err)
	}

	reservationLedgerConfig, err := NewReservationLedgerConfig(
		*reservationObj,
		rl.startFull,
		rl.overfillBehavior,
		rl.bucketCapacityPeriod,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	now := rl.timeSource()
	newLedger, err := NewReservationLedger(*reservationLedgerConfig, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation ledger: %w", err)
	}

	rl.ledgers.Add(accountID, newLedger)
	return newLedger, nil
}
