package payments

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
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

	reservationLedger *ReservationLedger
	onDemandLedger    *OnDemandLedger

	queue  chan coretypes.Blob
	status atomic.Value // stores LedgerStatus
}

func NewAccountLedger(
	// TODO: may be nil if no reservation exists
	reservationConfig *ReservationLedgerConfig,
	getNow func() time.Time,
) (*AccountLedger, error) {
	var reservationLedger *ReservationLedger
	if reservationConfig != nil {
		var err error
		reservationLedger, err = NewReservationLedger(*reservationConfig, getNow)
		if err != nil {
			return nil, fmt.Errorf("new reservation ledger: %w", err)
		}
	}

	// TODO: we probably want to add an on demand ledger config, and if that is nonnil, we init the on demand ledger

	onDemandLedger, err := NewOnDemandLedger()
	if err != nil {
		return nil, fmt.Errorf("new on demand ledger: %w", err)
	}

	accountLedger := &AccountLedger{
		reservationLedger: reservationLedger,
		onDemandLedger:    onDemandLedger,
		queue:             make(chan coretypes.Blob, 100), // buffer size of 100. TODO add to config
	}
	accountLedger.status.Store(LedgerStatusAlive)

	go accountLedger.ProcessBlobQueue()

	return accountLedger, nil
}

func (al *AccountLedger) Stop() {
	al.status.Store(LedgerStatusDead)
	close(al.queue)
}

func (al *AccountLedger) EnqueueBlobForAccounting(ctx context.Context, blob coretypes.Blob) error {
	if al.status.Load().(LedgerStatus) != LedgerStatusAlive {
		// TODO: make special error type, which causes the whole client to crash. cannot continue without a ledger
		return fmt.Errorf("ledger is not alive")
	}

	select {
	case al.queue <- blob:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("failed to enqueue blob: %w", ctx.Err())
	}
}

// ProcessBlobQueue pops blobs off the queue and processes them
func (al *AccountLedger) ProcessBlobQueue() {
	for blob := range al.queue {
		err := al.processBlob(blob)

		if err != nil {
			// TODO: debug log
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// TODO: doc, also better method name
func (al *AccountLedger) processBlob(blob coretypes.Blob) error {
	if al.reservationLedger != nil {
		// TODO: need to get quorums from blob
		var quorums []core.QuorumID
		// TODO: use this payment metadata
		_, err := al.reservationLedger.Debit(int64(blob.BlobLengthSymbols()), quorums)

		switch err.(type) {
		case nil:
			// Success - blob accounted for
			return nil
		case *InsufficientReservationCapacityError:
			// todo: add info log, then continue to on-demand
		default:
			al.Stop()
			// return nil, since we don't want to do the sleep in the ProcessBlobQueue loop
			return nil
		}
	}

	// try on-demand if configured to do so

	// TODO: make a REALLY good error here, with all sorts of juicy details
}
