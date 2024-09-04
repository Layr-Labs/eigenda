package chunkgroup

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"time"
)

// uint64ToBytes writes a uint64 to a byte array in big-endian order.
func uint64ToBytes(value uint64, bytes []byte) {
	for i := 0; i < 8; i++ {
		bytes[i] = byte(value >> (56 - 8*i))
	}
}

// bytesToUint64 converts a byte array in big-endian order to a uint64.
func bytesToUint64(bytes []byte) uint64 {
	var value uint64
	for i := 0; i < 8; i++ {
		value |= uint64(bytes[i]) << (56 - 8*i)
	}
	return value
}

// hash hashes a seed using the Keccak-256 algorithm. The returned value is an integer formed from the
// first 8 bytes of the hash.
func randomInt(seeds ...uint64) uint64 {
	seedBytes := make([]byte, 8*len(seeds))

	for i, seed := range seeds {
		uint64ToBytes(seed, seedBytes[i*8:(i+1)*8])
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(seedBytes)
	value := hasher.Sum(nil)
	return bytesToUint64(value)
}

// ComputeShuffleOffset returns the offset at which a light node should be shuffled into a new chunk group,
// relative to the beginning each shuffle interval.
func ComputeShuffleOffset(
	seed uint64,
	assignmentIndex uint64,
	shufflePeriod time.Duration) time.Duration {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if assignmentIndex >= 64 {
		panic(fmt.Sprintf("assignment index must be between 0 and 63, got %d", assignmentIndex))
	}

	return time.Duration(int64(randomInt(seed, assignmentIndex) % uint64(shufflePeriod)))
}

// ComputeShuffleEpoch returns the epoch number of a light node's group chunkGroupAssignment at the current time.
func ComputeShuffleEpoch(
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	now time.Time) uint64 {

	unixEpoch := time.Unix(0, 0)

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}
	if now.Before(unixEpoch) {
		panic(fmt.Sprintf("provided time %s is before epoch", now))
	}

	// The time when the first shuffle epoch for this node begins.
	// Note that this will be before unix epoch.
	genesis := unixEpoch.Add(shuffleOffset - shufflePeriod)

	timeSinceGenesis := now.Sub(genesis)
	return uint64(timeSinceGenesis / shufflePeriod)
}

// ComputeStartOfShuffleEpoch returns the time when a shuffle epoch begins.
func ComputeStartOfShuffleEpoch(
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	currentEpoch uint64) time.Time {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}
	if shuffleOffset < 0 {
		panic(fmt.Sprintf("shuffle offset must be non-negative, got %s", shuffleOffset))
	}

	unixEpoch := time.Unix(0, 0)

	// The time when the first shuffle epoch for this node begins.
	// Note that this will be before unix epoch.
	genesis := unixEpoch.Add(shuffleOffset - shufflePeriod)

	return genesis.Add(shufflePeriod * time.Duration(currentEpoch))
}

// ComputeEndOfShuffleEpoch given an epoch, return the time when that epoch will end and the next epoch will begin.
func ComputeEndOfShuffleEpoch(
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	currentEpoch uint64) time.Time {

	return ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, currentEpoch+1)
}

// ComputeChunkGroup returns the chunk group of a light node given its current shuffle epoch.
func ComputeChunkGroup(
	seed uint64,
	assignmentIndex uint64,
	shuffleEpoch uint64,
	chunkGroupCount uint64) uint64 {

	if assignmentIndex >= 64 {
		panic(fmt.Sprintf("assignment index must be between 0 and 63, got %d", assignmentIndex))
	}

	return randomInt(seed, assignmentIndex, shuffleEpoch) % chunkGroupCount
}
