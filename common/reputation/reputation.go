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
// Forgiveness increases low scores toward a neutral point over time.
//
// This structure is NOT goroutine safe.
type Reputation struct {
	config                  ReputationConfig
	score                   float64
	previousForgivenessTime time.Time
}

// Creates a new reputation tracker starting at the neutral forgiveness target.
func NewReputation(config ReputationConfig, now time.Time) *Reputation {
	return &Reputation{
		config:                  config,
		score:                   config.ForgivenessTarget,
		previousForgivenessTime: now,
	}
}

// Updates the reputation after a successful interaction.
// Moves the score toward 1.0 based on the configured success update rate.
// Applies forgiveness before updating the score.
func (r *Reputation) ReportSuccess(now time.Time) {
	r.forgive(now)
	r.score = (1-r.config.SuccessUpdateRate)*r.score + r.config.SuccessUpdateRate
}

// Updates the reputation after a failed interaction.
// Moves the score toward 0.0 based on the configured failure update rate.
// Applies forgiveness before updating the score.
func (r *Reputation) ReportFailure(now time.Time) {
	r.forgive(now)
	r.score = (1 - r.config.FailureUpdateRate) * r.score
}

// Returns the current reputation score.
// Applies forgiveness before returning the score.
func (r *Reputation) Score(now time.Time) float64 {
	r.forgive(now)
	return r.score
}

// Applies time-based drift toward the neutral forgiveness target.
// Only increases scores that are below the target - scores above the target are unchanged.
func (r *Reputation) forgive(now time.Time) {
	if r.previousForgivenessTime.IsZero() {
		r.previousForgivenessTime = now
		return
	}

	// Only forgive if score is below the forgiveness target
	if r.score >= r.config.ForgivenessTarget {
		return
	}

	elapsed := now.Sub(r.previousForgivenessTime).Seconds()
	if elapsed <= 0 {
		return
	}

	forgivenessRate := math.Log(2) / r.config.ForgivenessHalfLife.Seconds()
	forgivenessFraction := 1 - math.Exp(-forgivenessRate*elapsed)

	r.score = (1-forgivenessFraction)*r.score + forgivenessFraction*r.config.ForgivenessTarget
	r.previousForgivenessTime = now
}
