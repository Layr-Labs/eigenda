package clientledger

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments"
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

	reservationMonitor *reservation.ReservationVaultMonitor
	onDemandMonitor    *ondemand.OnDemandVaultMonitor
}

// Creates a ClientLedger, which is responsible for managing payments for a single client.
func NewClientLedger(
	ctx context.Context,
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
	// provides access to payment vault contract
	paymentVault payments.PaymentVault,
	// interval for checking for PaymentVault updates
	updateInterval time.Duration,
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
	if paymentVault == nil {
		panic("payment vault must not be nil")
	}

	clientLedger := &ClientLedger{
		logger:             logger,
		accountantMetricer: accountantMetricer,
		accountID:          accountID,
		clientLedgerMode:   clientLedgerMode,
		reservationLedger:  reservationLedger,
		onDemandLedger:     onDemandLedger,
		getNow:             getNow,
	}

	var err error
	if clientLedger.reservationLedger != nil {
		clientLedger.reservationMonitor, err = reservation.NewReservationVaultMonitor(
			ctx,
			logger,
			paymentVault,
			updateInterval,
			clientLedger.GetAccountsToUpdate,
			clientLedger.UpdateReservation)
		enforce.NilError(err, "new reservation vault monitor")
	}

	if clientLedger.onDemandLedger != nil {
		clientLedger.onDemandMonitor, err = ondemand.NewOnDemandVaultMonitor(
			ctx,
			logger,
			paymentVault,
			updateInterval,
			clientLedger.GetAccountsToUpdate,
			clientLedger.UpdateTotalDeposit)
		enforce.NilError(err, "new on demand vault monitor")
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
	success, remainingCapacity, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, fmt.Errorf("debit reservation: %w", err)
		}

		// all other modes of failure are fatal
		panic(fmt.Sprintf("reservation debit failed: %v", err))
	}

	cl.accountantMetricer.RecordReservationPayment(remainingCapacity)

	if !success {
		return nil, fmt.Errorf(
			"reservation lacks capacity for blob with %d symbols (%d bytes), "+
				"and no on-demand fallback is configured",
			blobLengthSymbols, blobLengthSymbols*encoding.BYTES_PER_SYMBOL)
	}

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
	enforce.NilError(err, "new payment metadata")
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
		"aware of new on-chain deposits without a restart")

	paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
	enforce.NilError(err, "new payment metadata")

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
	success, remainingCapacity, err := cl.reservationLedger.Debit(now, now, blobLengthSymbols, quorums)
	if err != nil {
		var timeMovedBackwardErr *reservation.TimeMovedBackwardError
		if errors.As(err, &timeMovedBackwardErr) {
			// this is the only class of error that can be returned from Debit where trying again might help
			return nil, fmt.Errorf("debit reservation: %w", err)
		}

		// all other modes of failure are fatal
		panic(fmt.Sprintf("reservation debit failed: %v", err))
	}

	cl.accountantMetricer.RecordReservationPayment(remainingCapacity)

	if success {
		paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
		enforce.NilError(err, "new payment metadata")
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
	enforce.NilError(err, "new payment metadata")

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

		remainingCapacity, err := cl.reservationLedger.RevertDebit(cl.getNow(), blobSymbolCount)
		if err != nil {
			return fmt.Errorf("revert reservation debit: %w", err)
		}

		cl.accountantMetricer.RecordReservationPayment(remainingCapacity)
	}

	return nil
}

// Returns the single account being tracked by this client ledger
func (cl *ClientLedger) GetAccountsToUpdate() []gethcommon.Address {
	return []gethcommon.Address{cl.accountID}
}

// Updates the reservation for the client's account
func (cl *ClientLedger) UpdateReservation(accountID gethcommon.Address, newReservation *reservation.Reservation) error {
	if accountID != cl.accountID {
		panic(fmt.Sprintf(
			"attempted to update reservation for the wrong account. Received account: %s, actual account: %s",
			accountID, cl.accountID))
	}

	err := cl.reservationLedger.UpdateReservation(newReservation, cl.getNow())
	if err != nil {
		return fmt.Errorf("update reservation: %w", err)
	}

	cl.accountantMetricer.RecordReservationBucketCapacity(cl.reservationLedger.GetBucketCapacity())

	return nil
}

// Updates the total deposit for the client's account
func (cl *ClientLedger) UpdateTotalDeposit(accountID gethcommon.Address, newTotalDeposit *big.Int) error {
	if accountID != cl.accountID {
		panic(fmt.Sprintf(
			"attempted to update total deposit for the wrong account. Received account: %s, actual account: %s",
			accountID, cl.accountID))
	}

	err := cl.onDemandLedger.UpdateTotalDeposits(newTotalDeposit)
	if err != nil {
		return fmt.Errorf("update total deposits: %w", err)
	}

	cl.accountantMetricer.RecordOnDemandTotalDeposits(cl.accountID.Hex(), newTotalDeposit)

	return nil
}
