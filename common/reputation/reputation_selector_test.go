package reputation

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

func createTestSelector(t *testing.T, config ReputationSelectorConfig) *ReputationSelector[testItem] {
	selector, err := NewReputationSelector(
		common.TestLogger(t),
		&config,
		random.NewTestRandom().Rand,
		func(item testItem) float64 { return item.score },
	)
	require.NoError(t, err)
	return selector
}

func TestReputationSelector_EmptyCandidates(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

	_, err := selector.Select([]testItem{})
	require.Error(t, err)
}

func TestReputationSelector_SingleCandidate(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

	candidates := []testItem{{id: "a", score: 0.5}}
	result, err := selector.Select(candidates)
	require.NoError(t, err)
	require.Equal(t, "a", result.id)
}

func TestReputationSelector_EqualWeights(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

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

func TestReputationSelector_ZeroScores(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

	candidates := []testItem{
		{id: "zeroA", score: 0.0},
		{id: "zeroB", score: 0.0},
	}

	_, err := selector.Select(candidates)
	require.NoError(t, err)
}

func TestReputationSelector_Filtering(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

	candidates := []testItem{
		{id: "a", score: 0.1},  // Bottom 50% AND below threshold -> filtered
		{id: "b", score: 0.11}, // Bottom 50% AND below threshold -> filtered
		{id: "c", score: 0.12}, // Not in bottom 50%, but below threshold -> included
		{id: "d", score: 1.0},  // Not in bottom 50%, and above threshold -> included
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

func TestReputationSelector_ThresholdPreservation(t *testing.T) {
	selector := createTestSelector(t, DefaultReputationSelectorConfig())

	candidates := []testItem{
		{id: "a", score: 0.3},  // Bottom 50% AND below threshold -> filtered
		{id: "b", score: 0.51}, // Bottom 50% BUT above threshold -> KEPT
		{id: "c", score: 0.75}, // Not in bottom 50% -> included
		{id: "d", score: 1.0},  // Not in bottom 50% -> included
	}

	selections := make(map[string]int)
	for range 2000 {
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
