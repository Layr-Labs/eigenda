package common

import (
	"context"
	"fmt"
	"time"
)

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

	// The previous period held by this ticker the last time its configuration was changed.
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
	// The target period to move towards.
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
	return NewVariableTickerWithFrequency(ctx, frequency), nil
}

// NewVariableTickerWithFrequency creates a new VariableTicker given a target frequency.
func NewVariableTickerWithFrequency(ctx context.Context, frequency float64) *VariableTicker {
	ctx, cancel := context.WithCancel(ctx)

	ticker := &VariableTicker{
		ctx:              ctx,
		close:            cancel,
		acceleration:     0.0,
		currentFrequency: frequency,
		currentPeriod:    time.Duration(1.0 / frequency),
		targetFrequency:  frequency,
		tickChan:         make(chan struct{}),
		controlChan:      make(chan any, 2),
	}

	go ticker.run()

	return ticker
}

// SetTargetPeriod sets the target period for the ticker. If acceleration is non-zero, the ticker will
// move towards the target period at the rate of acceleration per second. If acceleration is zero,
// the ticker will immediately adopt the target period.
func (t *VariableTicker) SetTargetPeriod(period time.Duration) error {
	if period <= 0 {
		return fmt.Errorf("invalid period %v, period must be positive", period)
	}

	frequency := float64(time.Second) / float64(period)

	t.controlChan <- &frequencyUpdate{
		targetFrequency: frequency,
	}

	return nil
}

func (t *VariableTicker) SetTargetFrequency(frequency float64) {
	t.controlChan <- &frequencyUpdate{
		targetFrequency: frequency,
	}
}

// SetAcceleration sets the acceleration for the frequency of the ticker, in HZ/second (i.e. 1/s/s).
func (t *VariableTicker) SetAcceleration(acceleration float64) {
	t.controlChan <- &accelerationUpdate{
		acceleration: acceleration,
	}
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
		// Check for control messages to update the ticker's configuration.
		select {
		case msg := <-t.controlChan:
			t.handleControlMessage(msg)
		default:
			// No control message received, continue with the current configuration.
		}

		if t.currentFrequency == 0 {
			// If the current frequency is zero, never tick.
			continue
		}

		// Send a tick.
		select {
		case t.tickChan <- struct{}{}:
		case <-t.ctx.Done():
			return
		}

		// Wait until it is time to send the next tick.
		sleepTime := t.computeSleepTime()
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

// computeSleepTime calculates the time to sleep until the next tick based on the current and target periods.
func (t *VariableTicker) computeSleepTime() time.Duration {
	if t.currentFrequency == t.targetFrequency {
		return t.currentPeriod
	}

	elapsedSinceAnchorTime := time.Since(t.anchorTime)

	var currentFrequency float64

	if t.acceleration == 0 {
		// Acceleration zero is defined as infinite acceleration. Immediately adopt the target frequency.
		currentFrequency = t.targetFrequency
	} else if t.currentFrequency < t.targetFrequency {
		// We are below the target frequency.
		currentFrequency = t.anchorFrequency + (t.acceleration * elapsedSinceAnchorTime.Seconds())
		if currentFrequency > t.targetFrequency {
			currentFrequency = t.targetFrequency
		}
	} else {
		// We are above the target frequency.
		currentFrequency = t.targetFrequency - (t.acceleration * elapsedSinceAnchorTime.Seconds())
		if currentFrequency < t.targetFrequency {
			currentFrequency = t.targetFrequency
		}
	}

	t.currentPeriod = time.Duration(1.0 / currentFrequency)
	return t.currentPeriod
}
