package reservation

import (
	"fmt"
	"time"
)

// TODO: docs
// TODO: consider config pointers throughout

type ReservationLedgerConfig struct {
	reservation       Reservation
	startFull         bool
	overfillBehavior OverfillBehavior
	// bucketCapacity is how much time worth of reservations should the capacity be
	bucketCapacityDuration time.Duration
}

func NewReservationLedgerConfig(
	reservation Reservation,
	startFull bool,
	overfillBehavior OverfillBehavior,
	bucketCapacityDuration time.Duration,
) (*ReservationLedgerConfig, error) {
	if bucketCapacityDuration <= 0 {
		return nil, fmt.Errorf("bucket capacity duration must be > 0, got %d", bucketCapacityDuration)
	}

	return &ReservationLedgerConfig{
		reservation:            reservation,
		startFull:              startFull,
		overfillBehavior:       overfillBehavior,
		bucketCapacityDuration: bucketCapacityDuration,
	}, nil
}
