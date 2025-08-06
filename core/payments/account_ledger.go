package payments

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// TODO: write unit tests

// TODO: work out how to fit metrics into this

// TODO: this is a replacement for the accountant and the meterer

// TODO: this ledger will need a queue that can be pushed on to from different go routines. Elements will then be
// popped off of the queue, and handled via "Debit"
type AccountLedger struct {
	reservationLedger *ReservationLedger
}

func NewAccountLedger(
	timeSource func() time.Time,
	// TODO: may be nil if no reservation exists
	reservationConfig *ReservationConfig,
) (*AccountLedger, error) {
	var leakyBucket *ReservationLedger
	if reservationConfig != nil {
		var err error
		leakyBucket, err = NewReservationLedger(timeSource, reservationConfig)
		if err != nil {
			return nil, fmt.Errorf("new leaky bucket: %w", err)
		}
	}

	return &AccountLedger{
		reservationLedger: leakyBucket,
	}, nil
}

func (al *AccountLedger) Debit(blob coretypes.Blob) error {
	if al.reservationLedger != nil {
		err := al.reservationLedger.Debit(int64(blob.BlobLengthSymbols()))

		if err != nil {
			return nil
		}
	}

	// try on-demand if configured to do so

	// TODO: make a REALLY good error here, with all sorts of juicy details
	return fmt.Errorf("")
}
