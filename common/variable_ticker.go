package common

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Any frequency at or below this value will be interpreted as a frequency of 0 Hz. Needed to avoid overflow.
// The factor of 2 is to take care of floating point precision issues.
const MinimumFrequency = float64(time.Nanosecond/math.MaxInt64) * 2.0

// Any frequency above this value will be interpreted as a frequency of MaximumFrequency Hz. Needed to avoid overflow.
const MaximumFrequency = float64(time.Nanosecond)

// VariableTicker behaves like a ticker with a frequency that can be changed at runtime.
type VariableTicker struct {
	ctx   context.Context
	close context.CancelFunc

	// The target frequency for the ticker, in HZ.
	targetFrequency float64

	// If the current frequency is not equal to the target frequency, the frequency will move towards the
	// target frequency at this rate per second. If zero, the  ticker will immediately adopt its target frequency.
	acceleration float64

	// The current frequency of the ticker, in HZ.
	currentFrequency float64

	// Matches currentFrequency. currentFrequency is the "source of truth", but we cache the period to avoid
	// recomputing it each tick.
	currentPeriod time.Duration

	// The previous frequency held by this ticker the last time its configuration was changed.
	anchorFrequency float64

	// The time at which the ticker last had its configuration changed.
	anchorTime time.Time

	// The channel that produces an output every time the ticker ticks.
	tickChan chan struct{}

	// The channel used to send control messages to main ticker loop.
	controlChan chan any
}

// frequencyUpdate is a control message to update the target frequency of the ticker.
type frequencyUpdate struct {
	// The target frequency to move towards.
	targetFrequency float64
}

// accelerationUpdate is a control message to update the acceleration of the ticker.
type accelerationUpdate struct {
	// The acceleration to apply to the ticker.
	acceleration float64
}

// NewVariableTickerWithPeriod creates a new VariableTicker given a target period.
func NewVariableTickerWithPeriod(ctx context.Context, period time.Duration) (*VariableTicker, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be positive, got %v", period)
	}
	frequency := float64(time.Second) / float64(period)
	return NewVariableTickerWithFrequency(ctx, frequency)
}

// NewVariableTickerWithFrequency creates a new VariableTicker given a target frequency.
func NewVariableTickerWithFrequency(ctx context.Context, frequency float64) (*VariableTicker, error) {
	if frequency < 0 {
		return nil, fmt.Errorf("frequency must be non-negative, got %v", frequency)
	}

	ctx, cancel := context.WithCancel(ctx)

	currentPeriod := time.Duration(0)
	if frequency > 0 {
		currentPeriod = time.Duration(float64(time.Second) / frequency)
	}

	ticker := &VariableTicker{
		ctx:              ctx,
		close:            cancel,
		acceleration:     0.0,
		currentFrequency: frequency,
		currentPeriod:    currentPeriod,
		targetFrequency:  frequency,
		tickChan:         make(chan struct{}),
		controlChan:      make(chan any, 2),
	}

	go ticker.run()

	return ticker, nil
}

// SetTargetPeriod sets the target period for the ticker. If acceleration is non-zero, the ticker will
// move towards the target period at the rate of acceleration per second. If acceleration is zero,
// the ticker will immediately adopt the target period.
func (t *VariableTicker) SetTargetPeriod(period time.Duration) error {
	if period <= 0 {
		return fmt.Errorf("invalid period %v, period must be positive", period)
	}
	frequency := float64(time.Second) / float64(period)
	return t.SetTargetFrequency(frequency)
}

func (t *VariableTicker) SetTargetFrequency(frequency float64) error {
	if frequency < 0 {
		return fmt.Errorf("invalid frequency %v, frequency must be non-negative", frequency)
	}

	if frequency < MinimumFrequency {
		frequency = 0.0
	}
	if frequency > MaximumFrequency {
		frequency = MaximumFrequency
	}

	t.controlChan <- &frequencyUpdate{
		targetFrequency: frequency,
	}

	return nil
}

// SetAcceleration sets the acceleration for the frequency of the ticker, in HZ/second (i.e. 1/s/s).
func (t *VariableTicker) SetAcceleration(acceleration float64) error {
	if acceleration < 0 {
		return fmt.Errorf("invalid acceleration %v, acceleration must be non-negative", acceleration)
	}

	t.controlChan <- &accelerationUpdate{
		acceleration: acceleration,
	}

	return nil
}

// Tick returns a channel that produces an output every time the ticker ticks.
func (t *VariableTicker) Tick() <-chan struct{} {
	return t.tickChan
}

// Close stops the ticker and releases any resources it holds.
func (t *VariableTicker) Close() {
	t.close()
}

// run produces ticks at the configured rate.
func (t *VariableTicker) run() {
	timer := time.NewTimer(t.currentPeriod)
	defer timer.Stop()

	for {
		t.computePeriod()
		if t.currentPeriod == 0 {
			// Period==0 is overloaded, and is used as a proxy for an infinitely long period (i.e. a frequency of 0).
			// In that case, do not tick.
			//
			// Only internal logic can set the period to 0. A user is unable to directly set the period to 0,
			// since if we interpret a period of 0 literally, it would require us to tick infinitely fast,

			select {
			case msg := <-t.controlChan:
				// check to see if we have a control message to process.
				t.handleControlMessage(msg)
			default:
				// to avoid busy waiting.
				time.Sleep(time.Millisecond)
			}
			continue
		}

		// Send a tick. Also listen for control messages.
		startOfTick := time.Now()
		var tickSent bool
		for !tickSent {
			select {
			case msg := <-t.controlChan:
				t.handleControlMessage(msg)
			case t.tickChan <- struct{}{}:
				tickSent = true
			case <-t.ctx.Done():
				return
			}
		}

		elapsed := time.Since(startOfTick)
		sleepTime := t.currentPeriod - elapsed
		if sleepTime < 0 {
			// If ticks are requested less often than the configured frequency, no need to sleep.
			continue
		}

		timer.Reset(sleepTime)
		select {
		case <-timer.C:
		case <-t.ctx.Done():
			return
		}
	}
}

// handleControlMessage processes control messages that update the ticker's configuration.
func (t *VariableTicker) handleControlMessage(msg any) {
	targetFrequency := t.targetFrequency
	acceleration := t.acceleration

	switch m := msg.(type) {
	case *frequencyUpdate:
		targetFrequency = m.targetFrequency
	case *accelerationUpdate:
		acceleration = m.acceleration
	default:
		// This should not be possible.
		panic(fmt.Sprintf("invalid control message type: %T", msg))
	}

	t.anchorTime = time.Now()
	t.anchorFrequency = t.currentFrequency
	t.targetFrequency = targetFrequency
	t.acceleration = acceleration
}

// computePeriod updates the current period based on configured frequency and acceleration
func (t *VariableTicker) computePeriod() {
	if t.currentFrequency == t.targetFrequency {
		// shortcut, don't recompute period if the period is already correct
		return
	}

	elapsedSinceAnchorTime := time.Since(t.anchorTime)

	if t.acceleration == 0 {
		// Acceleration zero is defined as infinite acceleration. Immediately adopt the target frequency.
		t.currentFrequency = t.targetFrequency
	} else if t.currentFrequency < t.targetFrequency {
		// We are below the target frequency.
		t.currentFrequency = t.anchorFrequency + (t.acceleration * elapsedSinceAnchorTime.Seconds())
		if t.currentFrequency > t.targetFrequency {
			// If we over shoot, adopt the target frequency.
			t.currentFrequency = t.targetFrequency
		} else {
			// When speeding up, substitute the current frequency with the inflection frequency.
			// This is to avoid sleeping for a very long time when starting from a low frequency.
			t.currentFrequency = t.computeInflectionFrequency()
		}
	} else {
		// We are above the target frequency.
		t.currentFrequency = t.anchorFrequency - (t.acceleration * elapsedSinceAnchorTime.Seconds())
		if t.currentFrequency < t.targetFrequency {
			// If we over shoot, adopt the target frequency.
			t.currentFrequency = t.targetFrequency
		}
	}

	if t.currentFrequency == 0 {
		t.currentPeriod = 0
	} else {
		t.currentPeriod = time.Duration(float64(time.Second) / t.currentFrequency)
	}
}

// computeInflectionFrequency handles an edge case when starting from a very low frequency. Suppose we start at 0.0 hz
// and are accelerating. At the moment we start accelerating, the frequency is zero and the period is infinite (1/0=∞),
// which is obviously not what we want. The "inflection frequency" is an adjusted frequency that will cause us to sleep
// for a more reasonable time. Specifically, it causes us to sleep long enough so that at the moment we wake up,
// the frequency at the moment we wake up will produce a period equal to the time we just slept.
func (t *VariableTicker) computeInflectionFrequency() float64 {
	// T0 = the current time, at this time we have frequency F0 and period P0=1/F0
	// T1 = the time at which we would wake up if we sleep for a period we calculate the period using F0
	//
	// T0                                      Ti                                      T1
	// |---------------------------------------|---------------------------------------|
	// <-----------------Pi-------------------->
	//
	// Ti = the inflection time, i.e. the time we want to wake up at
	// Pi = inflection period
	//
	// The goal is that at time Ti, if we use the inflection frequency Fi, we will find that we have a period of Pi.
	//
	// A = acceleration
	//
	// a) Pi = (Ti - T0) / 2
	// b) Fi = F0 + A * Pi
	// c) Pi = 1 / Fi
	//
	// Combine equations b and c:
	// d) Pi = 1 / (F0 + A * Pi)
	//
	// Plug equation d into an algebraic solver:
	// https://www.wolframalpha.com/input?i=solve+for+x+in+%28x+%3D+1%2F%28f+%2B+x+*+a%29%29
	// Variable substitution done since WolframAlpha gets confused by multi-character variables.
	// e) Pi = (sqrt(4A + F0^2) - F0) / 2A
	//
	// Combine equations c and e (i.e. invert the period to get the frequency):
	// f) Fi = 2A / (sqrt(4A + F0^2) - F0)
	return (2 * t.acceleration) / (math.Sqrt(4*t.acceleration+math.Pow(t.currentFrequency, 2)) - t.currentFrequency)
}
