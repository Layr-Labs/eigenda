package reservation

import (
	"fmt"
	"time"
)

// Configuration for a [ReservationLedger], which manages the reservation of a single account
type ReservationLedgerConfig struct {
	// Contains the parameters of the reservation that the [ReservationLedger] is responsible for
	reservation Reservation
	// Whether the underlying reservation [LeakyBucket] should start full or empty.
	// This asymmetric approach is necessary to handle restart scenarios correctly for different entities.
	//
	// Validators and Dispersers should start empty:
	// - An empty bucket allows dispersals immediately upon construction, up to available capacity
	// - This errs on the side of permitting more throughput from clients
	// - Critical for restart scenarios: if a validator had an empty bucket before restart and initialized
	//   with a full bucket, it would incorrectly deny all dispersals until leakage provides capacity,
	//   unfairly blocking honest clients from their reserved throughput
	//
	// Clients should start full:
	// - A full bucket requires waiting for leakage before dispersals can be made
	// - This errs on the side of underutilization
	// - Critical for restart scenarios: if a client had a full bucket before restart and initialized
	//   with an empty bucket, it would incorrectly believe it has full capacity to disperse
	// - Without this protection, clients with recurring problems that restart rapidly could over-utilize
	//   their reservation so severely that validators would begin rejecting dispersals
	startFull bool
	// Controls how the [LeakyBucket] handles dispersals that would exceed bucket capacity.
	//
	// Background: Small reservations may have bucket capacities smaller than the maximum blob size.
	// Without overfill support, users with small reservations would be unable to disperse large blobs,
	// even though their average rate permits it over time.
	//
	// This configuration parameter exists just in case we want to limit the cases that overfill is permitted in the
	// future, but the current intention is for all entities to run with [OverfillBehavior] == [OverfillOncePermitted]
	overfillBehavior OverfillBehavior
	// Determines the maximum burst capacity of the [LeakyBucket].
	//
	// The actual bucket capacity in symbols = symbolsPerSecond * bucketCapacityDuration
	//
	// This duration will be different for different parties, even for a given reservation. Clients will be configured
	// to have a smaller bucket size than dispersers and validators, to account for latency in the dispersal process.
	bucketCapacityDuration time.Duration
}

func NewReservationLedgerConfig(
	reservation Reservation,
	startFull bool,
	overfillBehavior OverfillBehavior,
	bucketCapacityDuration time.Duration,
) (*ReservationLedgerConfig, error) {
	if bucketCapacityDuration <= 0 {
		return nil, fmt.Errorf("bucket capacity duration must be > 0, got %v", bucketCapacityDuration)
	}

	return &ReservationLedgerConfig{
		reservation:            reservation,
		startFull:              startFull,
		overfillBehavior:       overfillBehavior,
		bucketCapacityDuration: bucketCapacityDuration,
	}, nil
}
