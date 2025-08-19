package reservation

import (
	"errors"
	"fmt"
	"time"
)

// This struct implements the [leaky bucket](https://en.wikipedia.org/wiki/Leaky_bucket) algorithm as a meter.
//
// Symbols "leak out" of the bucket at a constant rate, creating capacity for new symbols. The bucket can be "filled"
// with additional symbols if there is enough available capacity.
//
// The standard golang golang.org/x/time/rate.Limiter is not suitable for our use-case, since the Limiter doesn't
// support the concept of overfilling the bucket. We require the concept of overfill, for cases where a bucket size
// might be too small to fit the largest permissible blob size. We don't want to prevent users with a small reservation
// size from submitting large blobs.
//
// NOTE: This struct doesn't do any synchronization! The caller is responsible for making sure that only one goroutine
// is using it at a time.
type LeakyBucket struct {
	// Defines different ways that overfilling the bucket should be handled
	overfillBehavior OverfillBehavior

	// The total number of symbols that fit in the bucket
	bucketCapacity float64

	// The number of symbols that leak out of the bucket each second, as determined by the reservation.
	symbolsPerSecondLeakRate float64

	// The number of symbols currently in the bucket
	currentFillLevel float64

	// The time at which the previous leak calculation was made
	previousLeakTime time.Time
}

// Creates a new instance of the [LeakyBucket] algorithm
func NewLeakyBucket(
	// how fast symbols leak out of the bucket
	symbolsPerSecondLeakRate uint64,
	// bucketCapacityDuration * symbolsPerSecondLeakRate becomes the bucket capacity
	bucketCapacityDuration time.Duration,
	// whether the bucket should start full or empty
	startFull bool,
	// how to handle overfilling the bucket
	overfillBehavior OverfillBehavior,
	// the current time, when this is being constructed
	now time.Time,
) (*LeakyBucket, error) {
	if symbolsPerSecondLeakRate == 0 {
		return nil, errors.New("symbolsPerSecondLeakRate must be > 0")
	}

	bucketCapacity := float64(symbolsPerSecondLeakRate) * bucketCapacityDuration.Seconds()
	if bucketCapacity <= 0 {
		return nil, fmt.Errorf("bucket capacity must be > 0 (from leak rate %d symbols/sec * duration %s)",
			symbolsPerSecondLeakRate, bucketCapacityDuration)
	}

	currentFillLevel := float64(0)
	if startFull {
		// starting with a full bucket means some time must elapse to allow leakage before the bucket can be used
		currentFillLevel = bucketCapacity
	}

	return &LeakyBucket{
		overfillBehavior:         overfillBehavior,
		bucketCapacity:           bucketCapacity,
		symbolsPerSecondLeakRate: float64(symbolsPerSecondLeakRate),
		currentFillLevel:         currentFillLevel,
		previousLeakTime:         now,
	}, nil
}

// Fill the bucket with a number of symbols.
//
// Use a time source that includes monotonic time for best results.
//
// - Returns (true, nil) if the leaky bucket has enough capacity to accept the fill.
// - Returns (false, nil) if bucket lacks capacity to permit the fill.
// - Returns (false, error) for actual errors:
//   - [TimeMovedBackwardError] if input time is before previous leak time (only possible if monotonic time isn't used).
//   - Generic error for all other modes of failure.
//
// If the bucket doesn't have enough capacity to accommodate the fill, symbolCount IS NOT added to the bucket, i.e. a
// failed fill doesn't count against the meter.
func (lb *LeakyBucket) Fill(now time.Time, symbolCount uint32) (bool, error) {
	if symbolCount == 0 {
		return false, errors.New("symbolCount must be > 0")
	}

	err := lb.leak(now)
	if err != nil {
		return false, fmt.Errorf("leak: %w", err)
	}

	// this is how full the bucket would be, if the fill were to be accepted
	newFillLevel := lb.currentFillLevel + float64(symbolCount)

	// if newFillLevel is <= the total bucket capacity, no further checks are required
	if newFillLevel <= lb.bucketCapacity {
		lb.currentFillLevel = newFillLevel
		return true, nil
	}

	// this fill would result in the bucket being overfilled, so we check the overfill behavior to decide what to do
	switch lb.overfillBehavior {
	case OverfillNotPermitted:
		return false, nil
	case OverfillOncePermitted:
		bucketFull := lb.currentFillLevel >= lb.bucketCapacity

		// if there is no available capacity whatsoever, dispersal is never permitted, no matter the overfill behavior
		if bucketFull {
			return false, nil
		}

		lb.currentFillLevel = newFillLevel
		return true, nil
	default:
		panic(fmt.Sprintf("unknown overfill behavior %s", lb.overfillBehavior))
	}
}

// Reverts a previous fill, i.e. removes the number of symbols that got added to the bucket
//
// Use a time source that includes monotonic time for best results.
//
// - Returns [TimeMovedBackwardError] if input time is before previous leak time (only possible if monotonic time
// isn't used).
// - Returns a generic error for all other modes of failure.
//
// The input time should be the most up-to-date time, NOT the time of the original fill.
func (lb *LeakyBucket) RevertFill(now time.Time, symbolCount uint32) error {
	if symbolCount == 0 {
		return errors.New("symbolCount must be > 0")
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	lb.currentFillLevel = lb.currentFillLevel - float64(symbolCount)

	// Ensure fill level doesn't go negative
	if lb.currentFillLevel < 0 {
		lb.currentFillLevel = 0
	}

	return nil
}

// Lets the correct number of symbols leak out of the bucket, based on when we last leaked
//
// Returns [TimeMovedBackwardError] if input time is before previous leak time.
func (lb *LeakyBucket) leak(now time.Time) error {
	elapsed := now.Sub(lb.previousLeakTime)

	if elapsed < 0 {
		// This can only happen if the user passes in time instances without monotonic timestamps
		return &TimeMovedBackwardError{
			PreviousTime: lb.previousLeakTime,
			CurrentTime:  now,
		}
	}

	if elapsed == 0 {
		// nothing leaks if no time has passed
		return nil
	}

	leakage := elapsed.Seconds() * lb.symbolsPerSecondLeakRate
	lb.currentFillLevel = lb.currentFillLevel - leakage

	if lb.currentFillLevel < 0 {
		lb.currentFillLevel = 0
	}

	lb.previousLeakTime = now
	return nil
}
