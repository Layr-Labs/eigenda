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
	timeSource  func() time.Time
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
	if timeSource == nil {
		return nil, errors.New("timeSource cannot be nil")
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
//
// Returns (true, nil) if the reservation has enough capacity to perform the debit.
// Returns (false, nil) if the bucket lacks capacity to permit the dispersal.
// Returns (false, error) if an error occurs during validation.
func (pv *ReservationPaymentValidator) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) (bool, error) {
	ledger, err := pv.ledgerCache.GetOrCreate(ctx, accountID)
	if err != nil {
		return false, fmt.Errorf("get or create ledger: %w", err)
	}

	now := pv.timeSource()
	success, _, err := ledger.Debit(now, dispersalTime, symbolCount, quorumNumbers)
	if err != nil {
		return false, fmt.Errorf("debit reservation payment: %w", err)
	}

	return success, nil
}
