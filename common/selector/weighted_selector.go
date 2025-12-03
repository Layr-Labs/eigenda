package selector

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Performs weighted random selection with configurable filtering.
type WeightedSelector[T any] struct {
	logger        logging.Logger
	config        *WeightedSelectorConfig
	random        *rand.Rand
	scoreFunction func(T) float64
}

func NewWeightedSelector[T any](
	logger logging.Logger,
	config *WeightedSelectorConfig,
	random *rand.Rand,
	// Function to extract score from candidate. Score must be in range [0.0, 1.0), and is used for weighted selection.
	scoreFunction func(T) float64,
) (*WeightedSelector[T], error) {
	if random == nil {
		return nil, fmt.Errorf("random must not be nil")
	}
	if scoreFunction == nil {
		return nil, fmt.Errorf("scoreFunction must not be nil")
	}
	return &WeightedSelector[T]{
		logger:        logger,
		config:        config,
		random:        random,
		scoreFunction: scoreFunction,
	}, nil
}

// Chooses one item from the provided candidates using weighted random selection.
// Returns an error if candidates is empty.
func (ws *WeightedSelector[T]) Select(candidates []T) (T, error) {
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
func (ws *WeightedSelector[T]) filterLowPerformers(candidates []T) []T {
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
		ws.logger.Errorf("all %d candidates were filtered out, which means there is a bug in the filtering logic.",
			len(candidates))
		// fall back to using all candidates
		filtered = candidates
	}

	return filtered
}

// Performs weighted random selection based on scores.
func (ws *WeightedSelector[T]) weightedRandomSelect(candidates []T) (T, error) {
	scores := make([]float64, len(candidates))
	var totalWeight float64
	for i, candidate := range candidates {
		score := ws.scoreFunction(candidate)
		// The weighted selection algorithm doesn't handle zero values well, so we use a small positive value instead
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
