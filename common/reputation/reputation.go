package reputation

import (
	"math"
	"time"
)

// Reputation tracks the reliability of an entity using exponential moving average.
//
// Each time an interaction succeeds or fails, the reputation score moves toward 1.0 (perfect)
// or 0.0 (completely unreliable).
//
// The update rates control how quickly the score adapts. A higher rate means recent outcomes
// matter more. A lower rate means the score is more stable and takes longer to change.
//
// Forgiveness increases low scores toward a neutral point over time, even without new
// interactions. This prevents entities from being permanently penalized based on old
// failures when we lack recent information.
//
// This structure is NOT goroutine safe.
type Reputation struct {
	config        ReputationConfig
	score         float64
	lastUpdatedAt time.Time
}

// Creates a new reputation tracker starting at the neutral forgiveness target.
func NewReputation(config ReputationConfig, now time.Time) *Reputation {
	return &Reputation{
		config:        config,
		score:         config.ForgivenessTarget,
		lastUpdatedAt: now,
	}
}

// Updates the reputation after a successful interaction.
// Moves the score toward 1.0 based on the configured success update rate.
func (r *Reputation) Success() {
	r.score = (1-r.config.SuccessUpdateRate)*r.score + r.config.SuccessUpdateRate
}

// Failure updates the reputation after a failed interaction.
// Moves the score toward 0.0 based on the configured failure update rate.
func (r *Reputation) Failure() {
	r.score = (1 - r.config.FailureUpdateRate) * r.score
}

// Returns the current reputation score.
func (r *Reputation) Score() float64 {
	return r.score
}

// Forgive applies time-based drift toward the neutral forgiveness target.
// Only increases scores that are below the target - scores above the target are unchanged.
// This should be called on all reputations before making selection decisions,
// so that entities which haven't been interacted with recently have their
// low scores recover rather than remaining at their old values.
func (r *Reputation) Forgive(now time.Time) {
	if r.lastUpdatedAt.IsZero() {
		r.lastUpdatedAt = now
		return
	}

	// Only forgive if score is below the forgiveness target
	if r.score >= r.config.ForgivenessTarget {
		return
	}

	elapsed := now.Sub(r.lastUpdatedAt).Seconds()
	if elapsed <= 0 {
		return
	}

	forgivenessRate := math.Log(2) / r.config.ForgivenessHalfLife.Seconds()
	forgivenessFraction := 1 - math.Exp(-forgivenessRate*elapsed)

	r.score = (1-forgivenessFraction)*r.score + forgivenessFraction*r.config.ForgivenessTarget
	r.lastUpdatedAt = now
}
