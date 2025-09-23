package ratelimit

import (
	"errors"
	"fmt"
	"time"
)

// TimeMovedBackwardError indicates a timestamp was observed that is before a previously observed timestamp.
type TimeMovedBackwardError struct {
	PreviousTime time.Time
	CurrentTime  time.Time
}

func (e *TimeMovedBackwardError) Error() string {
	return fmt.Sprintf("time moved backward: previous=%v, current=%v", e.PreviousTime, e.CurrentTime)
}

// This struct implements the [leaky bucket](https://en.wikipedia.org/wiki/Leaky_bucket) algorithm as a meter.
//
// A leaky bucket is a metaphor for rate limiting. The bucket has a fixed capacity, and it leaks at a constant rate.
// When work is done, the bucket is "filled" with an amount of "water" proportional to the work done.
// Water "leaks out" of the bucket at a constant rate, creating capacity for new work.
//
// The standard golang golang.org/x/time/rate.Limiter is not suitable for some use cases, since the Limiter doesn't
// support the concept of overfilling the bucket. We require the concept of overfill, for cases where a bucket size
// might be too small to fit the largest permissible quantity of work.
//
// NOTE: This struct doesn't do any synchronization! The caller is responsible for making sure that only one goroutine
// is using it at a time.
type LeakyBucket struct {
	// Defines different ways that overfilling the bucket should be handled
	overfillBehavior OverfillBehavior

	// The total quantity of "water" that fit in the bucket
	bucketCapacity float64

	// The quantity of "water" that leaks out of the bucket each second, as determined by the configuration.
	leakRate float64

	// The amount of "water" currently in the bucket
	currentFillLevel float64

	// The time at which the previous leak calculation was made
	previousLeakTime time.Time
}

// Creates a new instance of the LeakyBucket algorithm
func NewLeakyBucket(
	// how fast "water" leaks out of the bucket
	leakRate uint64,
	// bucketCapacityDuration * leakRate becomes the bucket capacity
	bucketCapacityDuration time.Duration,
	// whether the bucket should start full or empty
	startFull bool,
	// how to handle overfilling the bucket
	overfillBehavior OverfillBehavior,
	// the current time, when this is being constructed
	now time.Time,
) (*LeakyBucket, error) {
	if leakRate == 0 {
		return nil, errors.New("leakRate must be > 0")
	}

	bucketCapacity := float64(leakRate) * bucketCapacityDuration.Seconds()
	if bucketCapacity <= 0 {
		return nil, fmt.Errorf("bucket capacity must be > 0 (from leak rate %d * duration %s)",
			leakRate, bucketCapacityDuration)
	}

	currentFillLevel := float64(0)
	if startFull {
		// starting with a full bucket means some time must elapse to allow leakage before the bucket can be used
		currentFillLevel = bucketCapacity
	}

	return &LeakyBucket{
		overfillBehavior: overfillBehavior,
		bucketCapacity:   bucketCapacity,
		leakRate:         float64(leakRate),
		currentFillLevel: currentFillLevel,
		previousLeakTime: now,
	}, nil
}

// Fill the bucket with "water", symbolizing work being done.
//
// Use a time source that includes monotonic time for best results.
//
// - Returns (true, nil) if the leaky bucket has enough capacity to accept the fill.
// - Returns (false, nil) if bucket lacks capacity to permit the fill.
// - Returns (false, error) for actual errors:
//   - [TimeMovedBackwardError] if input time is before previous leak time (only possible if monotonic time isn't used).
//   - Generic error for all other modes of failure.
//
// If the bucket doesn't have enough capacity to accommodate the fill, "water" IS NOT added to the bucket, i.e. a
// failed fill doesn't count against the meter.
func (lb *LeakyBucket) Fill(now time.Time, quantity float64) (bool, error) {
	if quantity <= 0 {
		return false, errors.New("quantity must be > 0")
	}

	err := lb.leak(now)
	if err != nil {
		return false, fmt.Errorf("leak: %w", err)
	}

	// this is how full the bucket would be, if the fill were to be accepted
	newFillLevel := lb.currentFillLevel + quantity

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

// Gets the current fill level of the bucket
//
// Use a time source that includes monotonic time for best results.
func (lb *LeakyBucket) CheckFillLevel(now time.Time) float64 {
	// even if there is an error, we still want to just return whatever the current fill level is
	_ = lb.leak(now)

	return lb.currentFillLevel
}

// Overrides the current fill level of the bucket, setting it to the specified value.
func (lb *LeakyBucket) SetFillLevel(now time.Time, fillLevel float64) error {
	if fillLevel < 0 {
		return errors.New("fill level must be >= 0")
	}

	if fillLevel > lb.bucketCapacity && lb.overfillBehavior == OverfillNotPermitted {
		return fmt.Errorf("fill level %f exceeds bucket capacity %f, but overfilling is not permitted",
			fillLevel, lb.bucketCapacity)
	}

	lb.previousLeakTime = now
	lb.currentFillLevel = fillLevel
	return nil
}

// Reverts a previous fill, i.e. removes a quantity of "water" that got added to the bucket
//
// Use a time source that includes monotonic time for best results.
//
// - Returns [TimeMovedBackwardError] if input time is before previous leak time (only possible if monotonic time
// isn't used).
// - Returns a generic error for all other modes of failure.
//
// The input time should be the most up-to-date time, NOT the time of the original fill.
func (lb *LeakyBucket) RevertFill(now time.Time, quantity float64) error {
	if quantity <= 0 {
		return errors.New("quantity must be > 0")
	}

	err := lb.leak(now)
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	lb.currentFillLevel = lb.currentFillLevel - quantity

	// Ensure fill level doesn't go negative
	if lb.currentFillLevel < 0 {
		lb.currentFillLevel = 0
	}

	return nil
}

// Lets the correct quantity of "water" leak out of the bucket, based on when we last leaked
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

	leakage := elapsed.Seconds() * lb.leakRate
	lb.currentFillLevel = lb.currentFillLevel - leakage

	if lb.currentFillLevel < 0 {
		lb.currentFillLevel = 0
	}

	lb.previousLeakTime = now
	return nil
}

// Gets the amount of capacity available in the bucket, i.e. how much "water" must be added to make the bucket
// exactly full.
//
// May be negative if the bucket is currently overfilled
func (lb *LeakyBucket) GetRemainingCapacity() float64 {
	return lb.bucketCapacity - lb.currentFillLevel
}

// Gets the total capacity of the bucket.
func (lb *LeakyBucket) GetBucketCapacity() float64 {
	return lb.bucketCapacity
}
