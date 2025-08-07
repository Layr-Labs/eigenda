package payments

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

type LedgerStatus string

const (
	LedgerStatusAlive LedgerStatus = "alive"
	LedgerStatusDead  LedgerStatus = "dead"
)

// TODO: write unit tests

// TODO: work out how to fit metrics into this

// TODO: this is a replacement for the accountant, bu it will not be the replacement for the meterer.
// this struct will need to construct payment headers, wait for the individual ledgers to be available with timeouts, etc.
// the disperser and validator nodes don't need to do any of that.
type AccountLedger struct {
	// TODO: add logger

	getNow func() time.Time

	reservationLedger *ReservationLedger

	onDemandLedger *OnDemandLedger

	status atomic.Value
}

func NewAccountLedger(
	// TODO: may be nil if no reservation exists
	reservationLedgerConfig *ReservationLedgerConfig,
	// TODO: may be nil if no on demand payments are enabled
	onDemandLedgerConfig *OnDemandLedgerConfig,
	getNow func() time.Time,
) (*AccountLedger, error) {
	var reservationLedger *ReservationLedger
	if reservationLedgerConfig != nil {
		var err error
		reservationLedger, err = NewReservationLedger(*reservationLedgerConfig, getNow())
		if err != nil {
			return nil, fmt.Errorf("new reservation ledger: %w", err)
		}
	}

	var onDemandLedger *OnDemandLedger
	if onDemandLedgerConfig != nil {
		var err error
		onDemandLedger, err = NewOnDemandLedger(*onDemandLedgerConfig)
		if err != nil {
			return nil, fmt.Errorf("new on demand ledger: %w", err)
		}
	}

	accountLedger := &AccountLedger{
		getNow:            getNow,
		reservationLedger: reservationLedger,
		onDemandLedger:    onDemandLedger,
	}
	accountLedger.status.Store(LedgerStatusAlive)

	return accountLedger, nil
}

// TODO: consider timeouts

// TODO: doc, also better method name
func (al *AccountLedger) Debit(
	ctx context.Context,
	blobLengthSymbols uint32,
	quorums []core.QuorumID,
) (*core.PaymentMetadata, error) {
	if al.status.Load().(LedgerStatus) != LedgerStatusAlive {
		// TODO: make special error type, which causes the whole client to crash. cannot continue without a ledger
		return nil, fmt.Errorf("ledger is not alive")
	}

	now := al.getNow()

	if al.reservationLedger != nil {
		paymentMetadata, err := al.reservationLedger.Debit(ctx, now, int64(blobLengthSymbols), quorums)

		switch err.(type) {
		case nil:
			// Success - blob accounted for
			return paymentMetadata, nil
		case *InsufficientReservationCapacityError:
			// todo: add info log, then continue to on-demand
		default:
			// TODO: make this a type of error which causes the client to shut down
			al.status.Store(LedgerStatusDead)
			return nil, fmt.Errorf("something unexpected happened, shut down")
		}
	}

	if al.onDemandLedger != nil {
		paymentMetadata, err := al.onDemandLedger.Debit(ctx, now, int64(blobLengthSymbols), quorums)
		switch err.(type) {
		case nil:
			// Success - blob accounted for
			return paymentMetadata, nil
		default:
			// TODO: make this a type of error which causes the client to shut down
			al.status.Store(LedgerStatusDead)
			return nil, fmt.Errorf("something unexpected happened, shut down")
		}
	}

	return nil, fmt.Errorf("TODO: make a REALLY good error here, with all sorts of juicy details")
}

// TODO: doc
func (al *AccountLedger) RevertDebit(
	ctx context.Context,
	paymentMetadata *core.PaymentMetadata,
	blobSymbolCount uint32,
) error {
	if al.status.Load().(LedgerStatus) != LedgerStatusAlive {
		// TODO: make special error type, which causes the whole client to crash. cannot continue without a ledger
		return fmt.Errorf("ledger is not alive")
	}

	if paymentMetadata.IsOnDemand() {
		if al.onDemandLedger == nil {
			return fmt.Errorf("unable to revert on demand payment with nil onDemandLedger")
		}

		err := al.onDemandLedger.RevertDebit(ctx, int64(blobSymbolCount))
		if err != nil {
			return fmt.Errorf("revert debit: %w", err)
		}
	} else {
		if al.reservationLedger == nil {
			return fmt.Errorf("unable to revert reservation payment with nil reservationLedger")
		}

		err := al.reservationLedger.RevertDebit(ctx, al.getNow(), int64(blobSymbolCount))
		if err != nil {
			return fmt.Errorf("revert reservation debit: %w", err)
		}
	}

	return nil
}
