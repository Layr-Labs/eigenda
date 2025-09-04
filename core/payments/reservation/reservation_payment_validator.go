package reservation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Validates reservation payments for multiple accounts
type ReservationPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked
	ledgerCache *ReservationLedgerCache
	timeSource   func() time.Time
}

func NewReservationPaymentValidator(
	ctx context.Context,
	logger logging.Logger,
	// the maximum number of ReservationLedger entries to be kept in the LRU cache
	maxLedgers int,
	// provides access to payment vault contract
	paymentVault payments.PaymentVault,
	// source of current time for the leaky bucket algorithm
	timeSource func() time.Time,
	// how to handle requests that would overfill the bucket
	overfillBehavior OverfillBehavior,
	// duration used to calculate bucket capacity
	bucketCapacityPeriod time.Duration,
	// interval for checking for payment updates
	updateInterval time.Duration,
) (*ReservationPaymentValidator, error) {
	if paymentVault == nil {
		return nil, errors.New("paymentVault cannot be nil")
	}

	if timeSource == nil {
		return nil, errors.New("timeSource cannot be nil")
	}

	if updateInterval <= 0 {
		return nil, errors.New("updateInterval must be > 0")
	}

	if bucketCapacityPeriod <= 0 {
		return nil, errors.New("bucketCapacityPeriod must be > 0")
	}

	ledgerCache, err := NewReservationLedgerCache(
		ctx,
		logger,
		maxLedgers,
		paymentVault,
		timeSource,
		overfillBehavior,
		bucketCapacityPeriod,
		updateInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger cache: %w", err)
	}

	return &ReservationPaymentValidator{
		logger:      logger,
		ledgerCache: ledgerCache,
		timeSource:  timeSource,
	}, nil
}

// Validates a reservation payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (pv *ReservationPaymentValidator) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) error {
	ledger, err := pv.ledgerCache.GetOrCreate(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get or create ledger: %w", err)
	}

	now := pv.timeSource()
	success, _, err := ledger.Debit(now, dispersalTime, symbolCount, quorumNumbers)
	if err != nil {
		return fmt.Errorf("debit reservation payment: %w", err)
	}

	if !success {
		return fmt.Errorf("reservation debit failed: insufficient capacity")
	}

	return nil
}

