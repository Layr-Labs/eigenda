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
type Reputation struct {
	config        ReputationConfig
	Reputation    float64
	LastUpdatedAt time.Time
}

// Creates a new reputation tracker starting at the neutral forgiveness target.
func NewReputation(config ReputationConfig, now time.Time) *Reputation {
	return &Reputation{
		config:        config,
		Reputation:    config.ForgivenessTarget,
		LastUpdatedAt: now,
	}
}

// Updates the reputation after a successful interaction.
// Moves the score toward 1.0 based on the configured success update rate.
func (r *Reputation) Success() {
	r.Reputation = (1-r.config.SuccessUpdateRate)*r.Reputation + r.config.SuccessUpdateRate
}

// Failure updates the reputation after a failed interaction.
// Moves the score toward 0.0 based on the configured failure update rate.
func (r *Reputation) Failure() {
	r.Reputation = (1 - r.config.FailureUpdateRate) * r.Reputation
}

// Forgive applies time-based drift toward the neutral forgiveness target.
// Only increases scores that are below the target - scores above the target are unchanged.
// This should be called on all reputations before making selection decisions,
// so that entities which haven't been interacted with recently have their
// low scores recover rather than remaining at their old values.
func (r *Reputation) Forgive(now time.Time) {
	if r.LastUpdatedAt.IsZero() {
		r.LastUpdatedAt = now
		return
	}

	// Only forgive if score is below the forgiveness target
	if r.Reputation >= r.config.ForgivenessTarget {
		return
	}

	elapsed := now.Sub(r.LastUpdatedAt).Seconds()
	if elapsed <= 0 {
		return
	}

	rate := math.Log(2) / r.config.ForgivenessHalfLife.Seconds()
	recoveryFraction := 1 - math.Exp(-rate*elapsed)

	r.Reputation = (1-recoveryFraction)*r.Reputation + recoveryFraction*r.config.ForgivenessTarget
	r.LastUpdatedAt = now
}
