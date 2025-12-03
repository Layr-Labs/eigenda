package selector

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

type testItem struct {
	id    string
	score float64
}

func createTestSelector(t *testing.T, config WeightedSelectorConfig) *WeightedSelector[testItem] {
	selector, err := NewWeightedSelector(
		common.TestLogger(t),
		&config,
		random.NewTestRandom().Rand,
		func(item testItem) float64 { return item.score },
	)
	require.NoError(t, err)
	return selector
}

func TestWeightedSelector_EmptyCandidates(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	_, err := selector.Select([]testItem{})
	require.Error(t, err)
}

func TestWeightedSelector_SingleCandidate(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	candidates := []testItem{{id: "a", score: 0.5}}
	result, err := selector.Select(candidates)
	require.NoError(t, err)
	require.Equal(t, "a", result.id)
}

func TestWeightedSelector_EqualWeights(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	candidates := []testItem{
		{id: "a", score: 0.5},
		{id: "b", score: 0.5},
		{id: "c", score: 0.5},
	}

	selections := make(map[string]int)
	for range 1000 {
		result, err := selector.Select(candidates)
		require.NoError(t, err)
		selections[result.id]++
	}

	// With equal weights, all should be selected roughly equally
	// Allow for some randomness (within 20% of expected value)
	expected := 1000 / 3
	for id, count := range selections {
		require.Greater(t, count, expected-expected/5, "item %s selected too few times", id)
		require.Less(t, count, expected+expected/5, "item %s selected too many times", id)
	}
}

func TestWeightedSelector_ZeroScores(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	candidates := []testItem{
		{id: "zero", score: 0.0},
		{id: "nonzero", score: 0.1},
	}

	selections := make(map[string]int)
	for range 1000 {
		result, err := selector.Select(candidates)
		require.NoError(t, err)
		selections[result.id]++
	}

	require.Greater(t, selections["zero"], 0, "zero score item should be selected at least once")
	require.Greater(t, selections["nonzero"], selections["zero"], "nonzero should be selected more than zero")
}

func TestWeightedSelector_Filtering(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	candidates := []testItem{
		{id: "a", score: 0.1}, // Bottom 50% AND below threshold -> filtered
		{id: "b", score: 0.2}, // Bottom 50% AND below threshold -> filtered
		{id: "c", score: 0.3}, // Not in bottom 50%, but below threshold -> included
		{id: "d", score: 0.9}, // Not in bottom 50%, and above threshold -> included
	}

	selections := make(map[string]int)
	for range 1000 {
		result, err := selector.Select(candidates)
		require.NoError(t, err)
		selections[result.id]++
	}

	// Items a and b should be filtered out
	require.Equal(t, 0, selections["a"], "item a should be filtered out")
	require.Equal(t, 0, selections["b"], "item b should be filtered out")
	// Items c and d should be selected
	require.Greater(t, selections["c"], 0, "item c should be selected")
	require.Greater(t, selections["d"], selections["c"], "item d should be selected more than item c")
}

func TestWeightedSelector_ThresholdPreservation(t *testing.T) {
	selector := createTestSelector(t, DefaultWeightedSelectorConfig())

	candidates := []testItem{
		{id: "a", score: 0.3}, // Bottom 50% AND below threshold -> filtered
		{id: "b", score: 0.6}, // Bottom 50% BUT above threshold -> KEPT
		{id: "c", score: 0.7}, // Not in bottom 50% -> included
		{id: "d", score: 0.9}, // Not in bottom 50% -> included
	}

	selections := make(map[string]int)
	for range 1000 {
		result, err := selector.Select(candidates)
		require.NoError(t, err)
		selections[result.id]++
	}

	// Item a should be filtered out
	require.Equal(t, 0, selections["a"], "item a should be filtered out")
	// Items b, c, d should all be selected (b is preserved by threshold)
	require.Greater(t, selections["b"], 0, "item b should be preserved by threshold")
	require.Greater(t, selections["c"], selections["b"], "item c should be selected more than item b")
	require.Greater(t, selections["d"], selections["c"], "item d should be selected more than item c")
}

func TestWeightedSelectorConfig_Validation(t *testing.T) {
	// Test invalid LowPerformerFraction
	config := WeightedSelectorConfig{LowPerformerFraction: -0.1, ScoreThreshold: 0.4}
	require.Error(t, config.Verify())

	config = WeightedSelectorConfig{LowPerformerFraction: 1.1, ScoreThreshold: 0.4}
	require.Error(t, config.Verify())

	// Test invalid ScoreThreshold
	config = WeightedSelectorConfig{LowPerformerFraction: 0.5, ScoreThreshold: -0.1}
	require.Error(t, config.Verify())

	config = WeightedSelectorConfig{LowPerformerFraction: 0.5, ScoreThreshold: 1.1}
	require.Error(t, config.Verify())

	// Test valid configs
	config = WeightedSelectorConfig{LowPerformerFraction: 0.5, ScoreThreshold: 0.4}
	require.NoError(t, config.Verify())

	config = WeightedSelectorConfig{LowPerformerFraction: 0, ScoreThreshold: 0}
	require.NoError(t, config.Verify())

	config = WeightedSelectorConfig{LowPerformerFraction: 1, ScoreThreshold: 1}
	require.NoError(t, config.Verify())
}
