package payments

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
)

type LeakyBucket struct {
	biasBehavior      BiasBehavior
	overdraftBehavior OverdraftBehavior

	// The total number of symbols that fit in the bucket
	bucketCapacity int64

	// The number of symbols that leak out of the bucket each second
	//
	// The leak rate is a uint64 on-chain, so that's what we accept in the constructor. But it's converted to an int64
	// under the hood to match other int types in the struct, making syntax less verbose.
	//
	// The assumption that symbolsPerSecondLeakRate fits in an int64 should be fine, as long as we are careful not to
	// give out any reservations that exceed 590 exabytes.
	symbolsPerSecondLeakRate int64

	// The number of symbols currently in the bucket
	currentFillLevel int64

	// The time at which the previous leak calculation was made.
	previousLeakTime time.Time

	// The number of symbols which leaked in the "partial second" of the previous leak calculation. A "partial" second
	// is `epochNanoTime % 1e9`.
	//
	// Since the leaky bucket uses integers instead of floats, leak math isn't straight forward. It's easy to calculate
	// the number of symbols that leak in a full second, since leak rate is defined in terms of symbols / second. But
	// determining how many symbols leak in a number of nanoseconds requires making a rounding choice. Leak calculation
	// N needs to take the partialSecondLeakage of calculation N-1 into account, so that the precisely correct number
	// of symbols are leaked for each full second.
	//
	// It would be possible to recalculate this value for N-1 when doing the calculation for N, but storing this value
	// as a member variable keeps things simple and avoids re-doing the math.
	previousPartialSecondLeakage int64
}

func NewLeakyBucket(
	// todo: note that we assume this is well formed. use proper constructors
	config ReservationLedgerConfig,
	now time.Time,
) (*LeakyBucket, error) {
	bucketCapacity := int64(float64(config.reservation.symbolsPerSecond) * config.bucketCapacityDuration.Seconds())

	var currentFillLevel int64
	switch config.biasBehavior {
	case BiasPermitMore:
		currentFillLevel = 0
	case BiasPermitLess:
		currentFillLevel = bucketCapacity
	default:
		return nil, fmt.Errorf("unknown bias behavior %s", config.biasBehavior)
	}

	return &LeakyBucket{
		biasBehavior:                 config.biasBehavior,
		overdraftBehavior:            config.overdraftBehavior,
		bucketCapacity:               bucketCapacity,
		symbolsPerSecondLeakRate:     config.reservation.symbolsPerSecond,
		currentFillLevel:             currentFillLevel,
		previousLeakTime:             now,
		previousPartialSecondLeakage: 0,
	}, nil
}

// Debit the reservation with a number of symbols.
//
// Algorithmically, that means adding a number of symbols to the leaky bucket.
//
// Returns nil if the leaky bucket has enough capacity to accept the fill. Returns an
// InsufficientReservationCapacityError if bucket lacks capacity to permit the fill.
//
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed debit doesn't count against the meter.
func (lb *LeakyBucket) Fill(now time.Time, symbolCount int64) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	// this is how full the bucket would be, if the fill were to be accepted
	newFillLevel, err := common.SafeAddInt64(lb.currentFillLevel, symbolCount)
	if err != nil {
		return fmt.Errorf("safe add to compute newFillLevel: %w", err)
	}

	// if newFillLevel is less than the total bucket capacity, no further checks are required
	if newFillLevel < lb.bucketCapacity {
		lb.currentFillLevel = newFillLevel
		return nil
	}

	// this fill would result in the bucket being overfilled, so we check the overfill behavior to decide what to do
	switch lb.overdraftBehavior {
	case OverdraftNotPermitted:
		return &InsufficientReservationCapacityError{symbolCount}
	case OverdraftOncePermitted:
		zeroCapacityAvailable := lb.currentFillLevel >= lb.bucketCapacity

		// if there is no available capacity whatsoever, dispersal is never permitted, no matter the overfill behavior
		if zeroCapacityAvailable {
			return &InsufficientReservationCapacityError{symbolCount}
		}

		lb.currentFillLevel = newFillLevel
		return nil
	default:
		return fmt.Errorf("unknown overfill behavior %s", lb.overdraftBehavior)
	}
}

// TODO: doc
func (lb *LeakyBucket) RevertFill(now time.Time, symbolCount int64) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	newFillLevel, err := common.SafeSubtractInt64(lb.currentFillLevel, symbolCount)
	if err != nil {
		return fmt.Errorf("safe subtract to compute newFillLevel: %w", err)
	}

	lb.currentFillLevel = newFillLevel
	return nil
}

// Lets the correct number of symbols leak out of the bucket, based on when we last leaked
//
// Returns an error if any of the calculations fail, which should not happen during normal usage.
func (lb *LeakyBucket) leak(now time.Time) error {
	defer func() {
		lb.previousLeakTime = now
	}()

	fullSecondLeakage, err := lb.computeFullSecondLeakage(now.Unix())
	if err != nil {
		return fmt.Errorf("compute full second leakage: %w", err)
	}

	// We need to correct the full-second leakage value: the previous leak calculation already let some symbols from a
	// partial second period leak out, and those symbols shouldn't be allowed to leak twice
	//
	// This value can be negative if the previous leak calculation was within the same second as this calculation.
	correctedFullSecondLeakage, err := common.SafeSubtractInt64(fullSecondLeakage, lb.previousPartialSecondLeakage)
	if err != nil {
		return fmt.Errorf("safe subtract to compute correctedFullSecondLeakage: %w", err)
	}

	partialSecondLeakage, err := lb.computePartialSecondLeakage(now.Nanosecond())
	if err != nil {
		return fmt.Errorf("compute partial second leakage: %w", err)
	}
	defer func() {
		lb.previousPartialSecondLeakage = partialSecondLeakage
	}()

	actualLeakage, err := common.SafeAddInt64(correctedFullSecondLeakage, partialSecondLeakage)
	if err != nil {
		return fmt.Errorf("safe add to compute actualLeakage: %w", err)
	}

	// don't let the bucket leak past empty
	if actualLeakage >= lb.currentFillLevel {
		lb.currentFillLevel = 0
		return nil
	}

	newFillLevel, err := common.SafeSubtractInt64(lb.currentFillLevel, actualLeakage)
	if err != nil {
		return fmt.Errorf("safe subtract to update currentFillLevel: %w", err)
	}

	lb.currentFillLevel = newFillLevel
	return nil
}

// Accepts the current number of seconds since epoch. Returns the number of symbols that should leak from the bucket,
// based on when we last leaked.
//
// Since this method only takes full seconds into consideration, the returned value must be used carefully. See leak()
// for details.
//
// Returns an error if the leakage calculation fails, which should not happen during normal usage.
func (lb *LeakyBucket) computeFullSecondLeakage(epochSeconds int64) (int64, error) {
	if epochSeconds < 0 {
		return 0, fmt.Errorf("epochSeconds must be >= 0, got %d", epochSeconds)
	}

	// epoch seconds should never go backwards, but could be the same
	if epochSeconds < lb.previousLeakTime.Unix() {
		return 0, fmt.Errorf("current time %s (%d) is before previous time %s (%d)",
			time.Unix(epochSeconds, 0).UTC().Format(time.RFC3339),
			epochSeconds,
			lb.previousLeakTime.UTC().Format(time.RFC3339),
			lb.previousLeakTime.Unix())
	}

	secondsSinceLastUpdate, err := common.SafeSubtractInt64(epochSeconds, lb.previousLeakTime.Unix())
	if err != nil {
		return 0, fmt.Errorf("safe subtract to compute secondsSinceLastUpdate: %w", err)
	}

	fullSecondLeakage, err := common.SafeMultiplyInt64(secondsSinceLastUpdate, lb.symbolsPerSecondLeakRate)
	if err != nil {
		return 0, fmt.Errorf("safe multiply to compute fullSecondLeakage: %w", err)
	}
	return fullSecondLeakage, nil
}

// TODO: doc
func (lb *LeakyBucket) computePartialSecondLeakage(nanos int) (int64, error) {
	if nanos >= 1e9 || nanos < 0 {
		return 0, fmt.Errorf("nanos must be between 0 and 1e9, got %d", nanos)
	}

	product, err := common.SafeMultiplyInt64(int64(nanos), lb.symbolsPerSecondLeakRate)
	if err != nil {
		return 0, fmt.Errorf("safe multiply to compute nanos * symbolsPerSecondLeakRate: %w", err)
	}

	switch lb.biasBehavior {
	case BiasPermitMore:
		// Round up, to permit more (more leakage = more capacity freed up)
		// Add (1e9 - 1) before dividing to round up
		sum, err := common.SafeAddInt64(product, 1e9-1)
		if err != nil {
			return 0, fmt.Errorf("safe add to compute rounding sum: %w", err)
		}
		return sum / 1e9, nil
	case BiasPermitLess:
		// Round down, to permit less (less leakage = less capacity freed up)
		return product / 1e9, nil
	default:
		return 0, fmt.Errorf("unknown bias: %s", lb.biasBehavior)
	}
}
