package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// The ClientLedger manages payment state for a single account. It is only used by *clients*, not by the disperser
// or validator nodes.
//
// The ClientLedger aggressively triggers panics for errors that indicate no future payments will succeed. A client
// is only useful if it can disperse blobs, and blobs can only be dispersed with a functioning payment mechanism.
//
// TODO(litt3): Currently, the client ledger has no mechanism to observe the following changes that may occur in the
// PaymentVault:
//
// 1. Reservation expiration
// 2. Reservation update
// 3. OnDemand deposit
//
// It is the responsibility of the user to restart the client if such a change occurs. A mechanism should be implemented
// to remove this burden from the user.
type ClientLedger struct {
	logger             logging.Logger
	accountantMetricer metrics.AccountantMetricer
	accountID          gethcommon.Address

	// Though it would theoretically be possible to infer mode of operation based on on-chain state, it's important
	// that this be directly configurable by the user, to ensure that reality matches intention.
	//
	// Consider, for example, if a user intends to operate with a reservation covering the majority of dispersals,
	// with an on-demand balance as a backup. If there is a configuration issue which prevents the reservation from
	// being used, the client could mistakenly burn through all backup funds before becoming aware of the
	// misconfiguration. In such cases, it's better to fail early, to bring the misconfiguration to the attention of the
	// user as soon as possible.
	clientLedgerMode ClientLedgerMode

	reservationLedger *reservation.ReservationLedger
	onDemandLedger    *ondemand.OnDemandLedger
	getNow            func() time.Time
}

// Creates a ClientLedger, which is responsible for managing payments for a single client.
func NewClientLedger(
	logger logging.Logger,
	accountantMetricer metrics.AccountantMetricer,
	// The account that this client ledger is for
	accountID gethcommon.Address,
	clientLedgerMode ClientLedgerMode,
	// may be nil if clientLedgerMode is configured to not use reservations
	reservationLedger *reservation.ReservationLedger,
	// may be nil if clientLedgerMode is configured to not use on-demand payments
	onDemandLedger *ondemand.OnDemandLedger,
	// Should be a timesource which includes monotonic timestamps, for best results. Otherwise, reservation payments
	// may occasionally fail due to NTP adjustments
	getNow func() time.Time,
) *ClientLedger {
	if accountantMetricer == nil {
		accountantMetricer = metrics.NoopAccountantMetrics
	}

	enforce.NotEquals(accountID, gethcommon.Address{}, "account ID cannot be zero address")

	switch clientLedgerMode {
	case ClientLedgerModeReservationOnly:
		enforce.NotNil(reservationLedger,
			"in %s mode, reservation ledger must be non-nil", ClientLedgerModeReservationOnly)
		enforce.Nil(onDemandLedger, "in %s mode, on-demand ledger must be nil", ClientLedgerModeReservationOnly)
	case ClientLedgerModeOnDemandOnly:
		enforce.NotNil(onDemandLedger, "in %s mode, on-demand ledger must be non-nil", ClientLedgerModeOnDemandOnly)
		enforce.Nil(reservationLedger, "in %s mode, reservation ledger must be nil", ClientLedgerModeOnDemandOnly)
	case ClientLedgerModeReservationAndOnDemand:
		enforce.NotNil(reservationLedger, "in %s mode, reservation ledger must be non-nil",
			ClientLedgerModeReservationAndOnDemand)
		enforce.NotNil(onDemandLedger, "in %s mode, on-demand ledger must be non-nil",
			ClientLedgerModeReservationAndOnDemand)
	default:
		panic(fmt.Sprintf("unknown clientLedgerMode %s", clientLedgerMode))
	}

	enforce.True(getNow != nil, "getNow function must not be nil")

	clientLedger := &ClientLedger{
		logger:             logger,
		accountantMetricer: accountantMetricer,
		accountID:          accountID,
		clientLedgerMode:   clientLedgerMode,
		reservationLedger:  reservationLedger,
		onDemandLedger:     onDemandLedger,
		getNow:             getNow,
	}

	return clientLedger
}

// Accepts parameters describing the aspects of a blob dispersal that are relevant for accounting. Attempts to use the
// configured payment method(s) to account for the blob.
//
// Returns a PaymentMetadata if the blob was successfully accounted for. This PaymentMetadata contains the
// information necessary to craft the dispersal message, and implicitly describes the payment mechanism being used.
//
// Returns an error for payment failures that could conceivably be resolved by retrying. Panics for all other failure
// modes, since inability to pay for dispersals requires intervention.
func (cl *ClientLedger) Debit(
	ctx context.Context,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	now := cl.getNow()

	// the handle methods in this switch contain some duplicate logic, but trying to generalize these operations
	// incurs a high complexity cost: the same underlying function calls are being made, but logging + error behavior
	// differs, depending on the specific mode of operation.
	switch cl.clientLedgerMode {
	case ClientLedgerModeReservationOnly:
		return cl.debitReservationOnly(now, blobLengthSymbols, quorums)
	case ClientLedgerModeOnDemandOnly:
		return cl.debitOnDemandOnly(ctx, now, blobLengthSymbols, quorums)
	case ClientLedgerModeReservationAndOnDemand:
		return cl.debitReservationOrOnDemand(ctx, now, blobLengthSymbols, quorums)
	default:
		panic(fmt.Sprintf("unknown clientLedgerMode %s", cl.clientLedgerMode))
	}
}

// Used by ClientLedger instances where only reservation payments are configured.
func (cl *ClientLedger) debitReservationOnly(
	now time.Time,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	// As the client, "now" and the dispersal time are the same. The client is responsible for populating the
	// dispersal time when constructing the payment header, and it does so with its conception of "now"
	success, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, fmt.Errorf("debit reservation: %w", err)
		}

		// all other modes of failure are fatal
		panic(fmt.Sprintf("reservation debit failed: %v", err))
	}

	if !success {
		return nil, fmt.Errorf(
			"reservation lacks capacity for blob with %d symbols (%d bytes), "+
				"and no on-demand fallback is configured",
			blobLengthSymbols, blobLengthSymbols*encoding.BYTES_PER_SYMBOL)
	}

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
	enforce.NilError(err, "new payment metadata: %v", err)
	return paymentMetadata, nil
}

// Used by ClientLedger instances where only on-demand payments are configured.
func (cl *ClientLedger) debitOnDemandOnly(
	ctx context.Context,
	now time.Time,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	cumulativePayment, err := cl.onDemandLedger.Debit(ctx, blobLengthSymbols, quorums)
	enforce.NilError(err, "on-demand debit failed. reservations aren't configured, and the ledger won't become "+
		"aware of new on-chain deposits without a restart: %v", err)

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
	enforce.NilError(err, "new payment metadata: %v", err)

	cl.accountantMetricer.RecordCumulativePayment(cl.accountID.Hex(), cumulativePayment)

	return paymentMetadata, nil
}

// Used by ClientLedger instances where both reservation and on-demand payments are configured.
//
// First tries to pay for a dispersal with the reservation, and falls back to on-demand if the reservation
// lacks capacity.
func (cl *ClientLedger) debitReservationOrOnDemand(
	ctx context.Context,
	now time.Time,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	// As the client, "now" and the dispersal time are the same. The client is responsible for populating the
	// dispersal time when constructing the payment header, and it does so with its conception of "now"
	success, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, fmt.Errorf("debit reservation: %w", err)
		}

		// all other modes of failure are fatal
		panic(fmt.Sprintf("reservation debit failed: %v", err))
	}

	if success {
		paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
		enforce.NilError(err, "new payment metadata: %v", err)
		return paymentMetadata, nil
	}

	cl.logger.Infof("Reservation lacks capacity for blob with %d symbols (%d bytes). Falling back to on-demand.",
		blobLengthSymbols, blobLengthSymbols*encoding.BYTES_PER_SYMBOL)

	cumulativePayment, err := cl.onDemandLedger.Debit(ctx, blobLengthSymbols, quorums)
	if err != nil {
		var InsufficientFundsError *ondemand.InsufficientFundsError
		if errors.As(err, &InsufficientFundsError) {
			// don't panic, since future dispersals could still use the reservation, once more capacity is available
			return nil, fmt.Errorf("debit on-demand: %w", err)
		}

		// everything else is a more serious problem, which requires human intervention
		panic(fmt.Sprintf("on-demand debit failed: %v", err))
	}

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
	enforce.NilError(err, "new payment metadata: %v", err)

	cl.accountantMetricer.RecordCumulativePayment(cl.accountID.Hex(), cumulativePayment)

	return paymentMetadata, nil
}

// RevertDebit undoes a previous debit.
//
// This should be called in cases where the client does accounting for a blob, but then the dispersal fails before
// being accounted for by the disperser.
func (cl *ClientLedger) RevertDebit(
	ctx context.Context,
	paymentMetadata *core.PaymentMetadata,
	blobSymbolCount uint32,
) error {
	if paymentMetadata.IsOnDemand() {
		enforce.NotNil(cl.onDemandLedger, "payment metadata is for an on-demand payment, but OnDemandLedger is nil")

		newCumulativePayment, err := cl.onDemandLedger.RevertDebit(ctx, blobSymbolCount)
		if err != nil {
			return fmt.Errorf("revert on-demand debit: %w", err)
		}

		cl.accountantMetricer.RecordCumulativePayment(cl.accountID.Hex(), newCumulativePayment)
	} else {
		enforce.NotNil(cl.reservationLedger,
			"payment metadata is for a reservation payment, but ReservationLedger is nil")

		err := cl.reservationLedger.RevertDebit(cl.getNow(), blobSymbolCount)
		if err != nil {
			return fmt.Errorf("revert reservation debit: %w", err)
		}
	}

	return nil
}
