package reputation

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

var _ config.VerifiableConfig = (*ReputationConfig)(nil)

type ReputationConfig struct {
	// How strongly to adjust the score after success.
	SuccessUpdateRate float64
	// How strongly to adjust the score after failure.
	FailureUpdateRate float64
	// How long it takes for a score to drift halfway back to the neutral point.
	ForgivenessHalfLife time.Duration
	// The score that reputation drifts toward over time when there are no interactions.
	ForgivenessTarget float64
}

func DefaultConfig() ReputationConfig {
	return ReputationConfig{
		SuccessUpdateRate:   0.05,
		FailureUpdateRate:   0.20,
		ForgivenessHalfLife: 24 * time.Hour,
		ForgivenessTarget:   0.5,
	}
}

// Verify implements [config.VerifiableConfig].
func (c *ReputationConfig) Verify() error {
	if c.SuccessUpdateRate < 0 || c.SuccessUpdateRate > 1 {
		return fmt.Errorf("SuccessUpdateRate must be between 0 and 1, got %f", c.SuccessUpdateRate)
	}
	if c.FailureUpdateRate < 0 || c.FailureUpdateRate > 1 {
		return fmt.Errorf("FailureUpdateRate must be between 0 and 1, got %f", c.FailureUpdateRate)
	}
	if c.ForgivenessHalfLife <= 0 {
		return fmt.Errorf("ForgivenessHalfLife must be positive, got %v", c.ForgivenessHalfLife)
	}
	if c.ForgivenessTarget <= 0 || c.ForgivenessTarget > 1 {
		return fmt.Errorf(
			"ForgivenessTarget must be between 0 (exclusive) and 1 (inclusive), got %f", c.ForgivenessTarget)
	}
	return nil
}
