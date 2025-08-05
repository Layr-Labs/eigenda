package payments

import (
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
)

// TODO: consider overflows in this. Are we being _too_ careful?

// This struct implements the [leaky bucket](https://en.wikipedia.org/wiki/Leaky_bucket) algorithm as a meter.
//
// Units "leak out" of the bucket at a constant rate, creating capacity for new units. The bucket can be "filled"
// with additional units if there is enough available capacity.
//
// The standard golang golang.org/x/time/rate.Limiter is not suitable for our use-case, for the following reasons:
//
//  1. The Limiter doesn't support the concept of overfilling the bucket. We require the concept of overfill, for cases
//     where a bucket size might be too small to fit the largest permissible blob size. We don't want to prevent users
//     with a small reservation size from submitting large blobs.
//  2. The Limiter uses floating point math. Though it would *probably* be ok to use floats, it makes the distributed
//     system harder to reason about. What level of error accumulation would we see with frequent updates? Under
//     what conditions would it be possible for the client and server representations of a given leaky bucket to
//     diverge, and what impact would that have on our assumptions? These questions can be avoided entirely by using
//     an integer based implementation.
//
// NOTE: methods on this struct should not be called from separate goroutines: it's not threadsafe.
type LeakyBucket struct {
	timeSource func() time.Time

	// Describes alternate strategies for how the leaky bucket responds when a Fill is attempted that exceeds the
	// available bucket capacity
	overdraftBehavior OverdraftBehavior

	// The total number of units that fit in the bucket
	bucketCapacity int64

	// The number of units that leak out of the bucket each second
	//
	// The leak rate is a uint64 on-chain, so that's what we accept in the constructor. But it's converted to an int64
	// under the hood to match other int types in the struct, making syntax less verbose.
	//
	// The assumption that unitsPerSecondLeakRate fits in an int64 should be fine, as long as we are careful not to
	// give out any reservations that exceed 590 exabytes.
	unitsPerSecondLeakRate int64

	// The type of biasBehavior to display, when confronted with a decision of whether to err on the side of permitting
	// more or less traffic.
	biasBehavior BiasBehavior

	// The number of units currently in the bucket
	currentFillLevel int64

	// The time at which the previous leak calculation was made.
	previousLeakTime time.Time

	// The number of units which leaked in the "partial second" of the previous leak calculation. A "partial" second
	// is `epochNanoTime % 1e9`.
	//
	// Since the leaky bucket uses integers instead of floats, leak math isn't straight forward. It's easy to calculate
	// the number of units that leak in a full second, since leak rate is defined in terms of units / second. But
	// determining how many units leak in a number of nanoseconds requires making a rounding choice. Leak calculation
	// N needs to take the partialSecondLeakage of calculation N-1 into account, so that the precisely correct number
	// of units are leaked for each full second.
	//
	// It would be possible to recalculate this value for N-1 when doing the calculation for N, but storing this value
	// as a member variable keeps things simple and avoids re-doing the math.
	previousPartialSecondLeakage int64
}

// Creates a new leaky bucket, which represents the reservation of a single user
func NewLeakyBucket(
	timeSource func() time.Time,
	config *ReservationConfig,
) (*LeakyBucket, error) {
	if config.bucketCapacityDuration <= 0 {
		return nil, fmt.Errorf("bucket capacity must be > 0, got %d", config.bucketCapacityDuration)
	}

	bucketCapacity := int64(float64(config.symbolsPerSecond) * config.bucketCapacityDuration.Seconds())

	if config.symbolsPerSecond == 0 {
		return nil, errors.New("symbolsPerSecond must be > 0")
	}

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
		timeSource:                   timeSource,
		overdraftBehavior:            config.overdraftBehavior,
		bucketCapacity:               bucketCapacity,
		unitsPerSecondLeakRate:       int64(config.symbolsPerSecond),
		biasBehavior:                 config.biasBehavior,
		currentFillLevel:             currentFillLevel,
		previousLeakTime:             timeSource(),
		previousPartialSecondLeakage: 0,
	}, nil
}

// Add a number of units to the leaky bucket.
//
// Returns nil if the bucket has enough capacity to accept the fill. Returns an InsufficientReservationCapacityError
// if bucket lacks capacity to permit the fill.
//
// If the bucket doesn't have enough capacity to accommodate the fill, unitCount IS NOT added to the bucket, i.e. a
// failed fill doesn't count against the meter.
//
// TODO: consider whether we should return the available capacity from this method?
func (lb *LeakyBucket) Fill(unitCount int64) error {
	if unitCount <= 0 {
		return fmt.Errorf("unitCount must be > 0, got %d", unitCount)
	}

	// leak the correct number of units, based on how long it's been since the last leak calculation
	err := lb.leak()
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	// this is how full the bucket would be, if the fill were to be accepted
	newFillLevel, err := common.SafeAddInt64(lb.currentFillLevel, unitCount)
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
		return &InsufficientReservationCapacityError{unitCount}
	case OverdraftOncePermitted:
		zeroCapacityAvailable := lb.currentFillLevel >= lb.bucketCapacity

		// if there is no available capacity whatsoever, dispersal is never permitted, no matter the overfill behavior
		if zeroCapacityAvailable {
			return &InsufficientReservationCapacityError{unitCount}
		}

		lb.currentFillLevel = newFillLevel
		return nil
	default:
		return fmt.Errorf("unknown overfill behavior %s", lb.overdraftBehavior)
	}
}

// Lets the correct number of units leak out of the bucket, based on when we last leaked
//
// Returns an error if any of the calculations fail, which should not happen during normal usage.
func (lb *LeakyBucket) leak() error {
	currentTime := lb.timeSource()
	defer func() {
		lb.previousLeakTime = currentTime
	}()

	fullSecondLeakage, err := lb.computeFullSecondLeakage(currentTime.Unix())
	if err != nil {
		return fmt.Errorf("compute full second leakage: %w", err)
	}

	// We need to correct the full-second leakage value: the previous leak calculation already let some units from a
	// partial second period leak out, and those units shouldn't be allowed to leak twice
	//
	// This value can be negative if the previous leak calculation was within the same second as this calculation.
	correctedFullSecondLeakage, err := common.SafeSubtractInt64(fullSecondLeakage, lb.previousPartialSecondLeakage)
	if err != nil {
		return fmt.Errorf("safe subtract to compute correctedFullSecondLeakage: %w", err)
	}

	partialSecondLeakage, err := lb.computePartialSecondLeakage(currentTime.Nanosecond())
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

// Accepts the current number of seconds since epoch. Returns the number of units that should leak from the bucket,
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

	fullSecondLeakage, err := common.SafeMultiplyInt64(secondsSinceLastUpdate, lb.unitsPerSecondLeakRate)
	if err != nil {
		return 0, fmt.Errorf("safe multiply to compute fullSecondLeakage: %w", err)
	}
	return fullSecondLeakage, nil
}

func (lb *LeakyBucket) computePartialSecondLeakage(nanos int) (int64, error) {
	if nanos >= 1e9 || nanos < 0 {
		return 0, fmt.Errorf("nanos must be between 0 and 1e9, got %d", nanos)
	}

	product, err := common.SafeMultiplyInt64(int64(nanos), lb.unitsPerSecondLeakRate)
	if err != nil {
		return 0, fmt.Errorf("safe multiply to compute nanos * unitsPerSecondLeakRate: %w", err)
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
