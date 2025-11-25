package reputation

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	testRandom := random.NewTestRandom()
	now := testRandom.Time()
	reputation := NewReputation(DefaultConfig(), now)

	for range 100 {
		reputation.Success(now)
	}
	require.Greater(t, reputation.Score(now), 0.99)
}

func TestFailure(t *testing.T) {
	testRandom := random.NewTestRandom()
	now := testRandom.Time()
	reputation := NewReputation(DefaultConfig(), now)

	for range 100 {
		reputation.Failure(now)
	}
	require.Less(t, reputation.Score(now), 0.01)
}

func TestForgive(t *testing.T) {
	t.Run("score above target unchanged", func(t *testing.T) {
		testRandom := random.NewTestRandom()
		startTime := testRandom.Time()
		reputation := NewReputation(DefaultConfig(), startTime)

		// lots of successes will result in high reputation
		for range 50 {
			reputation.Success(startTime)
		}
		scoreBeforeForgive := reputation.Score(startTime)

		// calling Score() after time has elapsed triggers forgiveness
		require.Equal(t, scoreBeforeForgive, reputation.Score(startTime.Add(1*time.Minute)),
			"forgiveness should only be applied to scores below the target")
	})

	t.Run("forgiveness converges to target", func(t *testing.T) {
		config := DefaultConfig()

		testRandom := random.NewTestRandom()
		startTime := testRandom.Time()
		reputation := NewReputation(config, startTime)

		// lots of failures will result in low reputation
		for range 50 {
			reputation.Failure(startTime)
		}

		// calling Score() after time has elapsed triggers forgiveness
		require.InDelta(t, config.ForgivenessTarget, reputation.Score(startTime.Add(100*config.ForgivenessHalfLife)), 0.0001,
			"forgiveness after a long time period should converge to the target level")
	})
}
