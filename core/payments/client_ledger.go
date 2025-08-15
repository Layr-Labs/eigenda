package payments

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/semaphore"
)

// TODO: write unit tests

// TODO: work out how to fit metrics into this
type ClientLedger struct {
	// TODO: add logger

	accountID gethcommon.Address

	getNow func() time.Time

	reservationLedger *reservation.ReservationLedger

	onDemandLedger *ondemand.OnDemandLedger

	alive atomic.Bool

	inflightOnDemandDispersals *semaphore.Weighted
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
		accountID:                  accountID,
		getNow:                     getNow,
		reservationLedger:          reservationLedger,
		onDemandLedger:             onDemandLedger,
		inflightOnDemandDispersals: semaphore.NewWeighted(1),
	}
	clientLedger.alive.Store(true)

	return clientLedger, nil
}

func (cl *ClientLedger) Debit(
	ctx context.Context,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	if !cl.alive.Load() {
		// TODO: make special error type, which causes the whole client to crash. cannot continue without a ledger
		return nil, fmt.Errorf("ledger is not alive")
	}

	now := cl.getNow()

	if cl.reservationLedger != nil {
		err := cl.reservationLedger.CheckInvariants(quorums, now)
		if err != nil {
			return nil, fmt.Errorf("") // TODO make this a good error. make sure this causes client to shut down
		}

		success, err := cl.reservationLedger.Debit(now, blobLengthSymbols)
		if err != nil {
			// TODO: make this a type of error which causes the client to shut down
			cl.alive.Store(false)
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
		// If not successful, continue to on-demand
		// todo: add info log
	}

	if cl.onDemandLedger != nil {
		if err := cl.inflightOnDemandDispersals.Acquire(ctx, 1); err != nil {
			return nil, fmt.Errorf("acquire inflight dispersal slot: %w", err)
		}

		cumulativePayment, err := cl.onDemandLedger.Debit(ctx, blobLengthSymbols, quorums)
		if err == nil {
			// Success - blob accounted for via on-demand
			paymentMetadata, err := core.NewPaymentMetadata(cl.accountID, now, cumulativePayment)
			if err != nil {
				return nil, fmt.Errorf("new payment metadata: %w", err)
			}
			return paymentMetadata, nil
		} else {
			// TODO: make this a type of error which causes the client to shut down
			// actually, not all errors.... failure to acquire semaphore shouldn't do that
			cl.alive.Store(false)
			return nil, fmt.Errorf("something unexpected happened, shut down")
		}
	}

	return nil, fmt.Errorf("TODO: make a REALLY good error here, with all sorts of juicy details")
}

// TODO: doc
func (cl *ClientLedger) revertDebit(
	ctx context.Context,
	paymentMetadata *core.PaymentMetadata,
	blobSymbolCount uint32,
) error {
	if !cl.alive.Load() {
		// TODO: make special error type, which causes the whole client to crash. cannot continue without a ledger
		return fmt.Errorf("ledger is not alive")
	}

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
	if paymentMetadata.IsOnDemand() {
		cl.inflightOnDemandDispersals.Release(1)
	}

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
