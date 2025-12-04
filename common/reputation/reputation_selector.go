package reputation

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Performs weighted random selection with configurable filtering of low performers.
//
// Selection is a two-stage process:
//  1. Filtering: Candidates that are in the bottom LowPerformerFraction AND have scores below ScoreThreshold
//     are excluded.
//  2. Weighted Selection: From remaining candidates, one is chosen randomly with probability proportional to score.
//
// The score function must return values >= 0. Higher scores increase selection probability.
// Zero scores are treated as 0.001 to ensure all candidates that aren't filtered have non-zero selection probability.
type ReputationSelector[T any] struct {
	config        *ReputationSelectorConfig
	random        *rand.Rand
	scoreFunction func(T) float64
}

func NewReputationSelector[T any](
	logger logging.Logger,
	config *ReputationSelectorConfig,
	random *rand.Rand,
	// Function to extract score from candidate. Score must be >= 0, and is used for weighted selection.
	scoreFunction func(T) float64,
) (*ReputationSelector[T], error) {
	err := config.Verify()
	if err != nil {
		return nil, fmt.Errorf("invalid reputation selector config: %w", err)
	}

	if random == nil {
		return nil, fmt.Errorf("random must not be nil")
	}
	if scoreFunction == nil {
		return nil, fmt.Errorf("scoreFunction must not be nil")
	}
	return &ReputationSelector[T]{
		config:        config,
		random:        random,
		scoreFunction: scoreFunction,
	}, nil
}

// Chooses one item from the provided candidates using weighted random selection.
// Returns an error if candidates is empty.
func (ws *ReputationSelector[T]) Select(candidates []T) (T, error) {
	var zero T

	if len(candidates) == 0 {
		return zero, fmt.Errorf("no candidates provided for selection")
	}

	// Sort candidates by score (ascending)
	slices.SortFunc(candidates, func(a, b T) int {
		scoreA := ws.scoreFunction(a)
		scoreB := ws.scoreFunction(b)
		if scoreA < scoreB {
			return -1
		} else if scoreA > scoreB {
			return 1
		}
		return 0
	})

	filteredCandidates := ws.filterLowPerformers(candidates)
	return ws.weightedRandomSelect(filteredCandidates)
}

// Filters out low performers based on config.
func (ws *ReputationSelector[T]) filterLowPerformers(candidates []T) []T {
	// Calculate how many candidates are in the low performer fraction. Round down to ensure we don't exclude all
	// candidates in cases where there are few eligible candidates.
	lowPerformerCount := int(math.Floor(float64(len(candidates)) * ws.config.LowPerformerFraction))

	// Filter out low performers
	filtered := make([]T, 0, len(candidates))
	for i, candidate := range candidates {
		score := ws.scoreFunction(candidate)
		// Exclude if in bottom percentile AND below threshold
		if i < lowPerformerCount && score < ws.config.ScoreThreshold {
			continue
		}
		filtered = append(filtered, candidate)
	}

	if len(filtered) == 0 {
		// fall back to using all candidates
		filtered = candidates
	}

	return filtered
}

// Performs weighted random selection based on scores.
func (ws *ReputationSelector[T]) weightedRandomSelect(candidates []T) (T, error) {
	scores := make([]float64, len(candidates))
	var totalWeight float64
	for i, candidate := range candidates {
		score := ws.scoreFunction(candidate)
		// Items with zero score would never be selected, so we use a small positive value instead
		if score == 0 {
			score = 0.001
		}
		scores[i] = score
		totalWeight += score
	}

	// Generate random number in [0, totalWeight)
	target := ws.random.Float64() * totalWeight

	// Walk through candidates, accumulating weight until we exceed target
	var accumulated float64
	for i, score := range scores {
		accumulated += score
		if accumulated >= target {
			return candidates[i], nil
		}
	}

	// We should never reach here, but return last candidate just in case
	return candidates[len(candidates)-1], nil
}
