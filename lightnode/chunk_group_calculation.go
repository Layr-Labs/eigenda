package lightnode

import (
	"fmt"
	"math/rand"
	"time"
)

// Consider the timeline below, with time moving from left to right.
//
//                               The "+" marks represent                    The 7th time this
//    The genesis time,          the time when a particular                 light node is shuffled.
//    i.e. protocol start        light node is shuffled.                              |
//           |                            |                                           ↓
//           ↓      1          2          ↓          4          5          6          7          8          9
//           |------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|
//           \          /          \      /                                \         /\         /\         /
//            \        /            \    /                                  \       /  \       /  \       /
//             \      /              \  /                                    \     /    \     /    \     /
//              \    /                \/                                      \   /      \   /      \   /
//               \  /                The "shuffle offset".                     \ /        \ /        \ /
//                \/                 Each light node has a                   epoch 6     epoch 7    epoch 8
//   A "shuffle period". Each node   random offset assigned
//   changes chunk groups once per   at registration time.
//   shuffle period. Each shuffle
//   period is marked with a "|".
//
// The algorithm for determining which chunk group a particular light node in is as follows:
// 1. Using the node's seed and a CSPRNG, determine the node's shuffle offset.
// 2. Define the genesis time to be "epoch 0".
// 3. Moving left to right over the timeline, add one to the epoch number for each time the clock is equal to
//    (genesis time + shuffle offset + X * shuffle period) for all integer values of X.
// 4. Using a CSPRNG, use the node epoch number and the node's seed to determine the node's chunk group.

// ComputeShuffleOffset returns the offset at which a light node should be shuffled into a new chunk group,
// relative to the beginning each shuffle interval.
func ComputeShuffleOffset(seed uint64, shufflePeriod time.Duration) time.Duration {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	rng := rand.New(rand.NewSource(int64(seed)))

	// TODO the algorithm used to determine this floating point value must be part of the spec
	multiple := rng.Float64()

	return time.Duration(float64(shufflePeriod) * multiple)
}

// ComputeShuffleEpoch returns the epoch number of a light node at the current time.
func ComputeShuffleEpoch(
	genesis time.Time,
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	now time.Time) uint64 {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}
	if now.Before(genesis) {
		panic(fmt.Sprintf("provided time %s is before genesis time %s", now, genesis))
	}

	// The time when the first epoch for this node begins.
	// Note that this will be before the genesis time unless the shuffle offset is exactly zero.
	epochGenesis := genesis.Add(shuffleOffset - shufflePeriod)

	timeSinceEpochGenesis := now.Sub(epochGenesis)
	return uint64(timeSinceEpochGenesis / shufflePeriod)
}

// ComputeEndOfShuffleEpoch given an epoch, return the time when that epoch will end and the next epoch will begin.
func ComputeEndOfShuffleEpoch(
	genesis time.Time,
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	currentEpoch uint64) time.Time {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}

	return genesis.Add(shuffleOffset).Add(shufflePeriod * time.Duration(currentEpoch))
}

// ComputeChunkGroup returns the chunk group of a light node given its current shuffle epoch.
func ComputeChunkGroup(
	seed uint64,
	shuffleEpoch uint64,
	chunkGroupCount uint64) uint64 {

	// TODO is adding the seed to the epoch sufficient?
	rng := rand.New(rand.NewSource(int64(seed + shuffleEpoch)))

	// TODO the algorithm used to determine this random value must be part of the spec
	return rng.Uint64() % chunkGroupCount
}
