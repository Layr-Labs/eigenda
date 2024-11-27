package metrics

import (
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"time"
)

// RunningAverage keeps track of the average of a series of values over a given time window.
type RunningAverage struct {
	maxAge  time.Duration
	sum     float64
	count   int
	entries queues.Queue
}

// NewRunningAverage creates a new RunningAverage with the given time window.
func NewRunningAverage(maxAge time.Duration) *RunningAverage {
	return &RunningAverage{
		maxAge:  maxAge,
		entries: linkedlistqueue.New(),
	}
}

type runningAverageEntry struct {
	value float64
	time  time.Time
}

// Update adds a new value to the RunningAverage and returns the new average.
func (a *RunningAverage) Update(now time.Time, value float64) float64 {
	a.count++
	a.sum += value
	a.entries.Enqueue(&runningAverageEntry{value: value, time: now})
	return a.GetAverage(now)
}

// GetAverage returns the current average of the RunningAverage.
func (a *RunningAverage) GetAverage(now time.Time) float64 {
	a.cleanup(now)
	if a.count == 0 {
		return 0
	}
	return a.sum / float64(a.count)
}

// cleanup removes old entries from the RunningAverage.
func (a *RunningAverage) cleanup(now time.Time) {

	for {
		v, ok := a.entries.Peek()
		if !ok {
			break
		}
		entry := v.(*runningAverageEntry)

		if now.Sub(entry.time) <= a.maxAge {
			break
		}

		a.entries.Dequeue()
		a.sum -= entry.value
		a.count--
	}

	if a.count == 0 {
		// clear away any cruft from accumulated floating point errors if we have no entries
		a.sum = 0
	}
}
