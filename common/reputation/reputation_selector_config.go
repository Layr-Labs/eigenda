package reputation

import "fmt"

// Configuration for the [ReputationSelector]
type ReputationSelectorConfig struct {
	// The fraction of candidates (sorted by score) to consider as "low performers", which may potentially be
	// excluded from selection.
	LowPerformerFraction float64
	// Candidates with a score higher than this will always be considered for selection, even if they fall within
	// the low performer fraction.
	ScoreThreshold float64
}

func DefaultReputationSelectorConfig() ReputationSelectorConfig {
	return ReputationSelectorConfig{
		LowPerformerFraction: 0.5,
		ScoreThreshold:       0.4,
	}
}

func (c *ReputationSelectorConfig) Verify() error {
	if c.LowPerformerFraction < 0 || c.LowPerformerFraction > 1 {
		return fmt.Errorf("LowPerformerFraction must be between 0 and 1, got %f", c.LowPerformerFraction)
	}
	if c.ScoreThreshold < 0 {
		return fmt.Errorf("ScoreThreshold must be >= 0, got %f", c.ScoreThreshold)
	}
	return nil
}
