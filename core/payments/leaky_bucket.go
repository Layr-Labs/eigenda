package payments

import (
	"fmt"
	"time"
)

// This struct implements the [leaky bucket](https://en.wikipedia.org/wiki/Leaky_bucket) algorithm as a meter.
//
// Symbols "leak out" of the bucket at a constant rate, creating capacity for new symbols. The bucket can be "filled"
// with additional symbols if there is enough available capacity.
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
// NOTE: This struct doesn't do any synchronization! The caller is responsible for making sure that only one goroutine
// is using it at a time.
type LeakyBucket struct {
	// Defines whether we should err on the side of permitting more or less throughput
	biasBehavior BiasBehavior

	// Defines different ways that overfilling the bucket should be handled
	overfillBehavior OverfillBehavior

	// The total number of symbols that fit in the bucket
	bucketCapacity int64

	// The number of symbols that leak out of the bucket each second
	symbolsPerSecondLeakRate int64

	// The number of symbols currently in the bucket
	currentFillLevel int64

	// The time at which the previous leak calculation was made
	previousLeakTime time.Time

	// The number of symbols which leaked in the "partial second" of the previous leak calculation.
	//
	// Since the leaky bucket uses integers instead of floats, leak math isn't straight forward. It's easy to calculate
	// the number of symbols that leak in a full second, since leak rate is defined in terms of symbols / second. But
	// determining how many symbols leak in a number of nanoseconds requires making a rounding choice. Leak calculation
	// N needs to take the partialSecondLeakage of calculation N-1 into account, so that the precisely correct number
	// of symbols are leaked for each full second.
	//
	//   (pipe characters represent second boundaries)
	//
	//   Previous leak (N-1)                      Current Leak (N)
	//            ↓                                      ↓
	//       |----*----------|----------------|----------*-----|
	//       ↑____↑
	//    previousPartialSecond
	previousPartialSecondLeakage int64
}

// Creates a new instance of the leaky bucket algorithm
func NewLeakyBucket(
	// how fast symbols leak out of the bucket
	symbolsPerSecondLeakRate int64,
	// bucketCapacityDuration * symbolsPerSecondLeakRate becomes the bucket capacity
	bucketCapacityDuration time.Duration,
	// whether to err on the side of permitting more or less throughput
	biasBehavior BiasBehavior,
	// how to handle overfilling the bucket
	overfillBehavior OverfillBehavior,
	// the current time, when this is being constructed
	now time.Time,
) (*LeakyBucket, error) {
	if symbolsPerSecondLeakRate <= 0 {
		return nil, fmt.Errorf("symbolsPerSecondLeakRate must be > 0, got %d", symbolsPerSecondLeakRate)
	}

	bucketCapacity := symbolsPerSecondLeakRate * bucketCapacityDuration.Nanoseconds() / 1e9

	if bucketCapacity <= 0 {
		return nil, fmt.Errorf("bucket capacity must be > 0, got %d (from leak rate %d symbols/sec * duration %s)",
			bucketCapacity, symbolsPerSecondLeakRate, bucketCapacityDuration)
	}

	var currentFillLevel int64
	switch biasBehavior {
	case BiasPermitMore:
		// starting with a fill level of 0 means the bucket starts out with available capacity
		currentFillLevel = 0
	case BiasPermitLess:
		// starting with a full bucket means some time must elapse to allow leakage before the bucket can be used
		currentFillLevel = bucketCapacity
	default:
		return nil, fmt.Errorf("unknown bias behavior %s", biasBehavior)
	}

	return &LeakyBucket{
		biasBehavior:                 biasBehavior,
		overfillBehavior:             overfillBehavior,
		bucketCapacity:               bucketCapacity,
		symbolsPerSecondLeakRate:     symbolsPerSecondLeakRate,
		currentFillLevel:             currentFillLevel,
		previousLeakTime:             now,
		previousPartialSecondLeakage: 0,
	}, nil
}

// Fill the bucket with a number of symbols.
//
// - Returns nil if the leaky bucket has enough capacity to accept the fill.
// - Returns an InsufficientReservationCapacityError if bucket lacks capacity to permit the fill.
// - Returns a TimeMovedBackwardError if input time is before previous leak time.
// - Returns a generic error for all other modes of failure.
//
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed fill doesn't count against the meter.
func (lb *LeakyBucket) Fill(now time.Time, symbolCount int64) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	// this is how full the bucket would be, if the fill were to be accepted
	newFillLevel := lb.currentFillLevel + symbolCount

	// if newFillLevel is <= the total bucket capacity, no further checks are required
	if newFillLevel <= lb.bucketCapacity {
		lb.currentFillLevel = newFillLevel
		return nil
	}

	// this fill would result in the bucket being overfilled, so we check the overfill behavior to decide what to do
	switch lb.overfillBehavior {
	case OverfillNotPermitted:
		return &InsufficientReservationCapacityError{symbolCount}
	case OverfillOncePermitted:
		zeroCapacityAvailable := lb.currentFillLevel >= lb.bucketCapacity

		// if there is no available capacity whatsoever, dispersal is never permitted, no matter the overfill behavior
		if zeroCapacityAvailable {
			return &InsufficientReservationCapacityError{symbolCount}
		}

		lb.currentFillLevel = newFillLevel
		return nil
	default:
		return fmt.Errorf("unknown overfill behavior %s", lb.overfillBehavior)
	}
}

// Reverts a previous fill, i.e. removes the number of symbols that got added to the bucket
//
// - Returns a TimeMovedBackwardError if input time is before previous leak time.
// - Returns a generic error for all other modes of failure.
//
// The input time should be the most up-to-date time, NOT the time of the original fill.
func (lb *LeakyBucket) RevertFill(now time.Time, symbolCount int64) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	newFillLevel := lb.currentFillLevel - symbolCount

	// don't let the bucket get emptier than "totally empty"
	if newFillLevel < 0 {
		lb.currentFillLevel = 0
		return nil
	}

	lb.currentFillLevel = newFillLevel
	return nil
}

// Lets the correct number of symbols leak out of the bucket, based on when we last leaked
//
// - Returns a TimeMovedBackwardError if input time is before previous leak time.
// - Returns a generic error if any of the calculations fail, which should not happen during normal usage.
//
//	Previous leak (N-1)                      Current Leak (N)
//	         ↓                                      ↓
//	    |----*----------|----------------|----------*-----|
//	         ↑______________________________________↑
//	   leaks the correct number of symbols for this time span
func (lb *LeakyBucket) leak(now time.Time) error {
	if now.Before(lb.previousLeakTime) {
		return &TimeMovedBackwardError{PreviousTime: lb.previousLeakTime, CurrentTime: now}
	}

	defer func() {
		lb.previousLeakTime = now
	}()

	fullSecondLeakage, err := lb.computeFullSecondLeakage(now.Unix())
	if err != nil {
		return fmt.Errorf("compute full second leakage: %w", err)
	}

	// We need to correct the full-second leakage value: the previous leak calculation already let some symbols from a
	// partial second period leak out, and those symbols shouldn't leak twice
	//
	// This value can be negative if the previous leak calculation was within the same second as this calculation,
	// since in that case fullSecondLeakage would be 0.
	//
	// Previous leak (N-1)                    Current Leak (N)
	//        ↓                                      ↓
	//   |----*----------|----------------|----------*-----|
	//        ↑___________________________↑
	//   this is the corrected full second leakage span
	correctedFullSecondLeakage := fullSecondLeakage - lb.previousPartialSecondLeakage

	partialSecondLeakage, err := lb.computePartialSecondLeakage(int64(now.Nanosecond()))
	if err != nil {
		return fmt.Errorf("compute partial second leakage: %w", err)
	}
	lb.previousPartialSecondLeakage = partialSecondLeakage

	actualLeakage := correctedFullSecondLeakage + partialSecondLeakage

	newFillLevel := lb.currentFillLevel - actualLeakage

	// don't let the bucket get emptier than "totally empty"
	if newFillLevel < 0 {
		lb.currentFillLevel = 0
		return nil
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
//
//	 Previous leak (N-1)                      Current Leak (N)
//	        ↓                                      ↓
//	   |----*----------|----------------|----------*-----|
//	   ↑________________________________↑
//	computes the correct number of symbols for this time span
func (lb *LeakyBucket) computeFullSecondLeakage(epochSeconds int64) (int64, error) {
	secondsSinceLastUpdate := epochSeconds - lb.previousLeakTime.Unix()
	fullSecondLeakage := secondsSinceLastUpdate * lb.symbolsPerSecondLeakRate
	return fullSecondLeakage, nil
}

// Accepts a number of nanoseconds, which represent a fraction of a single second.
//
// Computes the number of symbols which leak out in the given fractional second. Since this deals with integers,
// the configured bias determines which direction we round in.
//
//	Previous leak (N-1)                      Current Leak (N)
//	       ↓                                      ↓
//	  |----*----------|----------------|----------*-----|
//	                                   ↑__________↑
//	            returns the correct number of symbols for this time span
func (lb *LeakyBucket) computePartialSecondLeakage(nanos int64) (int64, error) {
	switch lb.biasBehavior {
	case BiasPermitMore:
		// Round up, to permit more (more leakage = more capacity freed up)
		// Add (1e9 - 1) before dividing to round up
		return (nanos*lb.symbolsPerSecondLeakRate + 1e9 - 1) / 1e9, nil
	case BiasPermitLess:
		// Round down, to permit less (less leakage = less capacity freed up)
		return nanos * lb.symbolsPerSecondLeakRate / 1e9, nil
	default:
		return 0, fmt.Errorf("unknown bias: %s", lb.biasBehavior)
	}
}
