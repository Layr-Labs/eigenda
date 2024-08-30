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

// rotateLeft shifts the bits of x to the left by k bits, with bits that fall off the
// left side wrapping around to the right.
func rotateLeft(x uint64, k uint32) uint64 {
	return (x << k) | (x >> (64 - k))
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
func ComputeShuffleOffset(assignmentSeed uint64, shufflePeriod time.Duration) time.Duration {
	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	return time.Duration(int64(randomInt(assignmentSeed) % uint64(shufflePeriod)))
}

// ComputeShuffleEpoch returns the epoch number of a light node's group chunkGroupAssignment at the current time.
func ComputeShuffleEpoch(
	shufflePeriod time.Duration,
	shuffleOffset time.Duration,
	now time.Time) uint64 {

	// TODO can this be a static constant?
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

	// TODO can this be a static constant?
	unixEpoch := time.Unix(0, 0)

	// The time when the first shuffle epoch for this node begins.
	// Note that this will be before unix epoch.
	genesis := unixEpoch.Add(shuffleOffset - shufflePeriod)

	return genesis.Add(shuffleOffset * time.Duration(currentEpoch))
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
	assignmentSeed uint64,
	shuffleEpoch uint64,
	chunkGroupCount uint32) uint32 {

	return uint32(randomInt(assignmentSeed^shuffleEpoch) % uint64(chunkGroupCount))
}
