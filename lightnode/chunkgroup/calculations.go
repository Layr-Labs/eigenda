package chunkgroup

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"time"
)

// uint64ToBytes converts a uint64 to a byte array in big-endian order.
func uint64ToBytes(value uint64) []byte {
	bytes := [8]byte{}
	for i := 0; i < 8; i++ {
		bytes[i] = byte(value >> (56 - 8*i))
	}
	return bytes[:]
}

// bytesToUint64 converts a byte array to a uint64 in big-endian order.
func bytesToUint64(bytes []byte) uint64 {
	var value uint64
	for i := 0; i < 8; i++ {
		value |= uint64(bytes[i]) << (56 - 8*i)
	}
	return value
}

// hash hashes a seed using the Keccak-256 algorithm. The returned value is an integer formed from the
// first 8 bytes of the hash.
func randomInt(seed uint64) uint64 {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(uint64ToBytes(seed))
	value := hasher.Sum(nil)
	return bytesToUint64(value)
}

// ComputeShuffleOffset returns the offset at which a light node should be shuffled into a new chunk group,
// relative to the beginning each shuffle interval.
func ComputeShuffleOffset(seed uint64, shufflePeriod time.Duration) time.Duration {
	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	return time.Duration(int64(randomInt(seed) % uint64(shufflePeriod)))
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

// ComputeStartOfShuffleEpoch returns the time when a shuffle epoch begins.
func ComputeStartOfShuffleEpoch(
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
	chunkGroupCount uint32) uint32 {

	return uint32(randomInt(seed^shuffleEpoch) % uint64(chunkGroupCount))
}
