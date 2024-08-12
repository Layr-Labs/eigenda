package chunkgroup

import (
	"fmt"
	"math/rand"
	"time"
)

// TODO use a regular uint for chunk index

// ComputeShuffleOffset returns the offset at which a light node should be shuffled into a new chunk group,
// relative to the beginning each shuffle interval.
func ComputeShuffleOffset(seed int64, shufflePeriod time.Duration) time.Duration {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	rng := rand.New(rand.NewSource(seed))

	return time.Duration(rng.Int63() % int64(shufflePeriod))
}

// ComputeShuffleEpoch returns the epoch number of a light node at the current time.
func ComputeShuffleEpoch(
	genesis time.Time,
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	now time.Time) int64 {

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
	return int64(timeSinceEpochGenesis / shufflePeriod)
}

// ComputeStartOfShuffleEpoch returns the time when a shuffle epoch begins.
func ComputeStartOfShuffleEpoch(
	genesis time.Time,
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	currentEpoch int64) time.Time {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}
	if currentEpoch < 0 {
		panic(fmt.Sprintf("current epoch must be non-negative, got %d", currentEpoch))
	}

	if currentEpoch == 0 {
		return genesis
	}
	return genesis.Add(shuffleOffset).Add(shufflePeriod * time.Duration(currentEpoch-1))
}

// ComputeEndOfShuffleEpoch given an epoch, return the time when that epoch will end and the next epoch will begin.
func ComputeEndOfShuffleEpoch(
	genesis time.Time,
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	currentEpoch int64) time.Time {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}
	if currentEpoch < 0 {
		panic(fmt.Sprintf("current epoch must be non-negative, got %d", currentEpoch))
	}

	return genesis.Add(shuffleOffset).Add(shufflePeriod * time.Duration(currentEpoch))
}

// ComputeChunkGroup returns the chunk group of a light node given its current shuffle epoch.
func ComputeChunkGroup(
	seed int64,
	shuffleEpoch int64,
	chunkGroupCount uint) uint {

	rng := rand.New(rand.NewSource(seed + shuffleEpoch))
	return uint(rng.Uint64()) % chunkGroupCount
}
