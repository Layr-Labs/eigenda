package payments

import (
	"fmt"
	"math"
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
type LeakyBucket struct {
	timeSource func() time.Time

	overfillBehavior OverfillBehavior

	maxWaitTime time.Duration

	bucketCapacity   uint32
	currentFillLevel uint32

	symbolsPerSecondLeakRate uint32

	previousUpdateTime       time.Time
	previousSubSecondLeakage uint32

	bias BiasBehavior
}

func NewLeakyBucket(
	timeSource func() time.Time,
	bucketCapacity uint32,
	bias BiasBehavior,
	symbolsPerSecondLeakRate uint32,
	overfillBehavior OverfillBehavior,
	maxWaitTime time.Duration,
) (*LeakyBucket, error) {
	var currentFillLevel uint32
	switch bias {
	case BiasPermitMore:
		currentFillLevel = 0
	case BiasPermitLess:
		currentFillLevel = bucketCapacity
	default:
		return nil, fmt.Errorf("unknown bias type %s", bias)
	}

	return &LeakyBucket{
		timeSource:               timeSource,
		overfillBehavior:         overfillBehavior,
		maxWaitTime:              maxWaitTime,
		bucketCapacity:           bucketCapacity,
		currentFillLevel:         currentFillLevel,
		symbolsPerSecondLeakRate: symbolsPerSecondLeakRate,
		previousUpdateTime:       timeSource(),
		bias:                     bias,
	}, nil
}

func (lb *LeakyBucket) Fill(symbols uint32) error {
	err := lb.leak()
	if err != nil {
		return fmt.Errorf("leak: %w", err)
	}

	newFillLevel := lb.currentFillLevel + symbols
	if newFillLevel < lb.bucketCapacity {
		lb.currentFillLevel = newFillLevel
		return nil
	}

	nonZeroCapacityAvailable := lb.currentFillLevel < lb.bucketCapacity

	if nonZeroCapacityAvailable && lb.overfillBehavior == OverfillOncePermitted {
		lb.currentFillLevel = newFillLevel
		return nil
	}
	// TODO: not done yet. keep going with the logic here
}

func (lb *LeakyBucket) computeDispersalTime(symbols uint32) (*time.Time, error) {
	var targetLevel uint32
	switch lb.overfillBehavior {
	case OverfillNotPermitted:
		// TODO; consider what if this is < 0
		targetLevel = lb.bucketCapacity - symbols
	case OverfillOncePermitted:
		targetLevel = lb.bucketCapacity - 1
	default:
		return nil, fmt.Errorf("unrecognized overfill behavior %s", lb.overfillBehavior)
	}

	// todo: note assumption that this method is only called if you actually need to wait
	amountToLeak := lb.currentFillLevel - targetLevel

	// the amount of leakage that will happen in the remainder of the current second
	subSecondRemainderToLeak := lb.symbolsPerSecondLeakRate - lb.previousSubSecondLeakage

	nanosRemainingInCurrentSecond := int64(1e9 - lb.previousUpdateTime.Nanosecond())

	if amountToLeak < subSecondRemainderToLeak {
		// We need to wait for a fraction of the remainder of the current second
		// Calculate how many nanos we need to wait based on the proportion of amountToLeak to symbolsPerSecondLeakRate

		// The formula is:
		//   time_to_wait = amountToLeak / symbolsPerSecondLeakRate  (in seconds)
		//
		//   To convert to nanoseconds:
		//   time_to_wait_nanos = (amountToLeak / symbolsPerSecondLeakRate) * 1e9
		//
		//   To avoid floating point division, we rearrange:
		//   time_to_wait_nanos = (amountToLeak * 1e9) / symbolsPerSecondLeakRate

		switch lb.bias {
		case BiasPermitMore:
			// Round down wait time (wait less, permit sooner)
			nanosToWait := int64(amountToLeak) * 1e9 / int64(lb.symbolsPerSecondLeakRate)
			dispersalTime := lb.previousUpdateTime.Add(time.Duration(nanosToWait))
			return &dispersalTime, nil
		case BiasPermitLess:
			// Round up wait time (wait more, be conservative)
			nanosToWait := (int64(amountToLeak)*1e9 + int64(lb.symbolsPerSecondLeakRate) - 1) / int64(lb.symbolsPerSecondLeakRate)
			dispersalTime := lb.previousUpdateTime.Add(time.Duration(nanosToWait))
			return &dispersalTime, nil
		default:
			return nil, fmt.Errorf("unknown bias: %s", lb.bias)
		}
	} else if amountToLeak == subSecondRemainderToLeak {
		// The amount to wait is exactly the sub-second time remainder
		// We need to wait until the next whole second
		dispersalTime := lb.previousUpdateTime.Add(time.Duration(nanosRemainingInCurrentSecond))
		return &dispersalTime, nil
	}

	// else we have to wait for the rest of the current second to elapse, then some number of seconds, then potentially
	// some number of nanos

	// First, subtract what will leak in the remainder of the current second
	amountToLeak -= subSecondRemainderToLeak

	// Calculate how many full seconds we need to wait
	fullSecondsToWait := amountToLeak / lb.symbolsPerSecondLeakRate

	// the amount of symbols that need to leak in the final fractional second
	nanosToLeakInFinalSecond := amountToLeak % lb.symbolsPerSecondLeakRate

	// Calculate nanos for the final fractional second
	var nanosToLeakRemainder int64
	if nanosToLeakInFinalSecond > 0 {
		// Apply bias for the partial second calculation
		switch lb.bias {
		case BiasPermitMore:
			// Round down wait time
			nanosToLeakRemainder = (int64(nanosToLeakInFinalSecond) * 1e9) / int64(lb.symbolsPerSecondLeakRate)
		case BiasPermitLess:
			// Round up wait time
			nanosToLeakRemainder = ((int64(nanosToLeakInFinalSecond) * 1e9) + int64(lb.symbolsPerSecondLeakRate) - 1) / int64(lb.symbolsPerSecondLeakRate)
		default:
			return nil, fmt.Errorf("unknown bias: %s", lb.bias)
		}
	}

	totalNanosToWait := nanosRemainingInCurrentSecond + int64(fullSecondsToWait)*1e9 + nanosToLeakRemainder
	dispersalTime := lb.previousUpdateTime.Add(time.Duration(totalNanosToWait))
	return &dispersalTime, nil
}

func (lb *LeakyBucket) leak() error {
	currentTime := lb.timeSource()
	defer func() {
		lb.previousUpdateTime = currentTime
	}()

	fullSecondLeakage, err := lb.computeFullSecondLeakage(currentTime.Unix())
	if err != nil {
		return fmt.Errorf("compute full second leakage: %w", err)
	}
	subSecondLeakage, err := lb.computeSubSecondLeakage(int64(currentTime.Nanosecond()))
	if err != nil {
		return fmt.Errorf("compute sub second leakage")
	}
	defer func() {
		lb.previousSubSecondLeakage = subSecondLeakage
	}()

	actualLeakage := fullSecondLeakage + subSecondLeakage - lb.previousSubSecondLeakage

	if actualLeakage > lb.currentFillLevel {
		lb.currentFillLevel = 0
		return nil
	}

	lb.currentFillLevel = lb.currentFillLevel - actualLeakage

	return nil
}

func (lb *LeakyBucket) computeSubSecondLeakage(nanos int64) (uint32, error) {
	// Check that nanos is less than 1 second
	if nanos >= 1e9 || nanos < 0 {
		return 0, fmt.Errorf("nanos must be between 0 and 1e9, got %d", nanos)
	}

	// Since nanos < 1e9 and symbolsPerSecondLeakRate is uint32,
	// the product will fit in int64 (max value would be (1e9-1) * (2^32-1) < 2^63)
	product := nanos * int64(lb.symbolsPerSecondLeakRate)

	switch lb.bias {
	case BiasPermitMore:
		// Round up when bias is to permit more (more leakage = more capacity freed up)
		// Add (1e9 - 1) before dividing to round up
		// The result after division will always fit in uint32 since we're dividing by 1e9
		return uint32((product + 1e9 - 1) / 1e9), nil
	case BiasPermitLess:
		// Round down when bias is to permit less (less leakage = less capacity freed up)
		return uint32(product / 1e9), nil
	default:
		return 0, fmt.Errorf("unknown bias: %s", lb.bias)
	}
}

func (lb *LeakyBucket) computeFullSecondLeakage(seconds int64) (uint32, error) {
	if seconds < 0 {
		return 0, fmt.Errorf("seconds must be >= 0, got %d", seconds)
	}

	if seconds < lb.previousUpdateTime.Unix() {
		return 0, fmt.Errorf("current time %d is before previous time %d", seconds, lb.previousUpdateTime.Unix())
	}

	secondsSinceLastUpdate := seconds - lb.previousUpdateTime.Unix()

	// Check if secondsDiff * symbolsPerSecondLeakRate would overflow uint32
	if uint64(secondsSinceLastUpdate) > math.MaxUint32/uint64(lb.symbolsPerSecondLeakRate) {
		// TODO; double check this overflow logic
		return 0, fmt.Errorf("TODO, error message")
	}

	fullSecondLeakage := uint32(secondsSinceLastUpdate) * lb.symbolsPerSecondLeakRate
	return fullSecondLeakage, nil
}
