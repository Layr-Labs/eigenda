package payments

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// TODO: this is a replacement for the accountant and the meterer
type AccountLedger struct {
	leakyBucket *LeakyBucket
}

func NewAccountLedger(
	timeSource func() time.Time,
	// TODO: may be nil if no reservation exists
	reservationConfig *ReservationConfig,
) (*AccountLedger, error) {
	var leakyBucket *LeakyBucket
	if reservationConfig != nil {
		var err error
		leakyBucket, err = NewLeakyBucket(timeSource, reservationConfig)
		if err != nil {
			return nil, fmt.Errorf("new leaky bucket: %w", err)
		}
	}

	return &AccountLedger{
		leakyBucket: leakyBucket,
	}, nil
}

func (al *AccountLedger) Debit(blob coretypes.Blob) error {
	if al.leakyBucket != nil {
		err := al.leakyBucket.Fill(int64(blob.BlobLengthSymbols()))

		if err != nil {
			return nil
		}
	}

	// try on-demand if configured to do so

	// TODO: make a REALLY good error here, with all sorts of juicy details
	return fmt.Errorf("")
}
