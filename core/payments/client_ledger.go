package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TODO: write unit tests

// TODO: work out how to fit metrics into this
type ClientLedger struct {
	// TODO: add logger

	accountID gethcommon.Address

	getNow func() time.Time

	reservationLedger *reservation.ReservationLedger

	onDemandLedger *ondemand.OnDemandLedger
}

func NewClientLedger(
	accountID gethcommon.Address,
	// may be nil if no reservation exists
	reservationLedger *reservation.ReservationLedger,
	// may be nil if no on demand payments are enabled
	onDemandLedger *ondemand.OnDemandLedger,
	getNow func() time.Time,
) (*ClientLedger, error) {

	clientLedger := &ClientLedger{
		accountID:         accountID,
		getNow:            getNow,
		reservationLedger: reservationLedger,
		onDemandLedger:    onDemandLedger,
	}

	return clientLedger, nil
}

func (cl *ClientLedger) Debit(
	ctx context.Context,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	now := cl.getNow()

	if cl.reservationLedger != nil {
		err := cl.reservationLedger.CheckInvariants(quorums, now)
		if err != nil {
			// TODO: add panic text here, make sure to include error
			panic("")
		}

		success, err := cl.reservationLedger.Debit(now, blobLengthSymbols)
		if err != nil {

			// TODO: check if this is a recoverable error. recoverable errors are any structured errors that debit may return
			// if the error isn't recoverable, panic. if error is recoverable, log it and continue on (don't return)
			return nil, fmt.Errorf("reservation debit error: %w", err)
		}

		if success {
			// Success - blob accounted for via reservation
			paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, nil)
			if err != nil {
				return nil, fmt.Errorf("new payment metadata: %w", err)
			}
			return paymentMetadata, nil
		}
		// todo: add info log, saying reservation payment failed
	}

	if cl.onDemandLedger != nil {
		cumulativePayment, err := cl.onDemandLedger.Debit(ctx, blobLengthSymbols, quorums)
		if err == nil {
			// Success - blob accounted for via on-demand
			paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
			if err != nil {
				return nil, fmt.Errorf("new payment metadata: %w", err)
			}
			return paymentMetadata, nil
		} else {
			// TODO: check if this is a recoverable error. recoverable errors are any structured errors that debit may return
			// if the error isn't recoverable, panic
			return nil, fmt.Errorf("something unexpected happened, shut down")
		}
	}

	return nil, fmt.Errorf("")
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
