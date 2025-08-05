package payments

import (
	"time"
)

// TODO: docs

type ReservationConfig struct {
	symbolsPerSecond  uint64
	biasBehavior      BiasBehavior
	overdraftBehavior OverdraftBehavior
	// bucketCapacity is how much time worth of reservations should the capacity be
	bucketCapacityDuration time.Duration
}
