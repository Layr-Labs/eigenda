package reservationvalidation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Validates reservation payments for multiple accounts
type ReservationPaymentValidator struct {
	logger logging.Logger
	// A cache of the ledgers being tracked
	ledgerCache *ReservationLedgerCache
	timeSource  func() time.Time
	metrics     *ReservationValidatorMetrics
}

func NewReservationPaymentValidator(
	ctx context.Context,
	logger logging.Logger,
	config ReservationLedgerCacheConfig,
	// provides access to payment vault contract
	paymentVault payments.PaymentVault,
	// source of current time for the leaky bucket algorithm
	timeSource func() time.Time,
	validatorMetrics *ReservationValidatorMetrics,
	cacheMetrics *ReservationCacheMetrics,
) (*ReservationPaymentValidator, error) {

	ledgerCache, err := NewReservationLedgerCache(
		ctx,
		logger,
		config,
		paymentVault,
		timeSource,
		cacheMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger cache: %w", err)
	}

	return &ReservationPaymentValidator{
		logger:      logger,
		ledgerCache: ledgerCache,
		timeSource:  timeSource,
		metrics:     validatorMetrics,
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

	if err == nil {
		if success {
			pv.metrics.RecordSuccess(symbolCount)
		} else {
			pv.metrics.IncrementInsufficientBandwidth()
		}
		return success, nil
	}

	var quorumNotPermittedErr *reservation.QuorumNotPermittedError
	if errors.As(err, &quorumNotPermittedErr) {
		pv.metrics.IncrementQuorumNotPermitted()
		return false, err
	}

	var timeOutOfRangeErr *reservation.TimeOutOfRangeError
	if errors.As(err, &timeOutOfRangeErr) {
		pv.metrics.IncrementTimeOutOfRange()
		return false, err
	}

	var timeMovedBackwardErr *reservation.TimeMovedBackwardError
	if errors.As(err, &timeMovedBackwardErr) {
		pv.metrics.IncrementTimeMovedBackward()
		return false, err
	}

	pv.metrics.IncrementUnexpectedErrors()
	return false, err
}
