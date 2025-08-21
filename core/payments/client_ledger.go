package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TODO: write unit tests

// TODO: work out how to fit metrics into this

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
	logger    logging.Logger
	accountID gethcommon.Address

	// Though it would theoretically be possible to infer mode of operation based on on-chain state, it's important
	// that this is directly configurable by the user, to ensure that reality matches intention.
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
	// The account that this client ledger is for.
	accountID gethcommon.Address,
	clientLedgerMode ClientLedgerMode,
	reservationLedger *reservation.ReservationLedger,
	onDemandLedger *ondemand.OnDemandLedger,
	// Should be a timesource which includes monotonic timestamps, for best results. Otherwise, reservation payments
	// may occasionally fail due to NTP adjustments
	getNow func() time.Time,
) (*ClientLedger, error) {

	switch clientLedgerMode {
	case ClientLedgerModeReservationOnly:
		if reservationLedger == nil || onDemandLedger != nil {
			panic(fmt.Sprintf("in %s mode, expected reservation ledger to be non-nil and on-demand ledger to be nil",
				ClientLedgerModeReservationOnly))
		}
	case ClientLedgerModeOnDemandOnly:
		if onDemandLedger == nil || reservationLedger != nil {
			panic(fmt.Sprintf("in %s mode, expected on-demand ledger to be non-nil and reservation ledger to be nil",
				ClientLedgerModeOnDemandOnly))
		}
	case ClientLedgerModeReservationAndOnDemand:
		if reservationLedger == nil || onDemandLedger == nil {
			panic(fmt.Sprintf("in %s mode, expected reservation and on-demand ledgers to be non-nil",
				ClientLedgerModeReservationAndOnDemand))
		}
	default:
		panic(fmt.Sprintf("unknown clientLedgerMode %s", clientLedgerMode))
	}

	clientLedger := &ClientLedger{
		logger:            logger,
		accountID:         accountID,
		clientLedgerMode:  clientLedgerMode,
		reservationLedger: reservationLedger,
		onDemandLedger:    onDemandLedger,
		getNow:            getNow,
	}

	return clientLedger, nil
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

// Used ClientLedger instances where only reservation payments are configured.
func (cl *ClientLedger) debitReservationOnly(
	now time.Time,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	// As the client, "now" and the dispersal time are the same. The client is responsible for populating the
	// dispersal time when constructing the payment header (below), and it does so with its conception of "now"
	success, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, err
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
	if err != nil {
		panic(fmt.Sprintf("new payment metadata: %w", err))
	}
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
	if err != nil {
		panic(fmt.Sprintf("on-demand debit failed. reservations aren't configured, and the ledger won't become "+
			"aware of new on-chain deposits without a restart: %v", err))
	}

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
	if err != nil {
		panic(fmt.Sprintf("new payment metadata: %w", err))
	}
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
	// dispersal time when constructing the payment header (below), and it does so with its conception of "now"
	success, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, err
		}

		// all other modes of failure are fatal
		panic(fmt.Sprintf("reservation debit failed: %v", err))
	}

	if success {
		paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
		if err != nil {
			panic(fmt.Sprintf("new payment metadata: %w", err))
		}
		return paymentMetadata, nil
	}

	cl.logger.Infof("Reservation lacks capacity for blob with %d symbols (%d bytes). Falling back to on-demand.",
		blobLengthSymbols, blobLengthSymbols*encoding.BYTES_PER_SYMBOL)

	cumulativePayment, err := cl.onDemandLedger.Debit(ctx, blobLengthSymbols, quorums)
	if err != nil {
		var InsufficientFundsError *ondemand.InsufficientFundsError
		if errors.As(err, &InsufficientFundsError) {
			// don't panic, since future dispersals could still use the reservation, once more capacity is available
			return nil, err
		}

		// everything else is a more serious problem, which requires human intervention
		panic(fmt.Sprintf("on-demand debit failed: %v", err))
	}

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
	if err != nil {
		panic(fmt.Sprintf("new payment metadata: %w", err))
	}
	return paymentMetadata, nil
}

// Undoes a previous debit.
//
// This should be called in cases where the client does accounting for a blob, but then the dispersal fails before
// the being accounted for by the disperser.
func (cl *ClientLedger) revertDebit(
	ctx context.Context,
	paymentMetadata *core.PaymentMetadata,
	blobSymbolCount uint32,
) error {
	if paymentMetadata.IsOnDemand() {
		if cl.onDemandLedger == nil {
			return fmt.Errorf("unable to revert on demand payment with nil onDemandLedger")
		}

		err := cl.onDemandLedger.RevertDebit(ctx, blobSymbolCount)
		if err != nil {
			return fmt.Errorf("revert debit: %w", err)
		}
	} else {
		if cl.reservationLedger == nil {
			return fmt.Errorf("unable to revert reservation payment with nil reservationLedger")
		}

		err := cl.reservationLedger.RevertDebit(cl.getNow(), blobSymbolCount)
		if err != nil {
			return fmt.Errorf("revert reservation debit: %w", err)
		}
	}

	return nil
}

func (cl *ClientLedger) DispersalSent(
	ctx context.Context,
	paymentMetadata *core.PaymentMetadata,
	symbolCount uint32,
	success bool,
) error {
	if success {
		return nil
	}

	// If the dispersal wasn't a success, that means that the disperser didn't charge the client for it, so the local
	// ledger should "refund" itself the cost of the dispersal
	err := cl.revertDebit(ctx, paymentMetadata, symbolCount)
	if err != nil {
		return fmt.Errorf("revert debit: %w", err)
	}

	return nil
}
