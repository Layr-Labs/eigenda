package common

import (
	"crypto/sha256"
	"encoding/binary"
)

// DeterministicRand is a deterministic random number generator that uses SHA-256
// to produce a stream of random bytes from an initial seed
type DeterministicRand struct {
	currentHash []byte
	buffer      []byte
	bufferPos   int
}

// NewDeterministicRand creates a new deterministic random number generator from a seed
func NewDeterministicRand(seed []byte) *DeterministicRand {
	// Initialize with the seed
	hasher := sha256.New()
	hasher.Write(seed)
	initialHash := hasher.Sum(nil)

	return &DeterministicRand{
		currentHash: initialHash,
		buffer:      initialHash,
		bufferPos:   0,
	}
}

// getNextBytes generates the next batch of random bytes
func (r *DeterministicRand) getNextBytes() {
	// Use the current hash as input to generate the next hash
	hasher := sha256.New()
	hasher.Write(r.currentHash)
	r.currentHash = hasher.Sum(nil)
	r.buffer = r.currentHash
	r.bufferPos = 0
}

// Uint32 returns a random uint32
func (r *DeterministicRand) Uint32() uint32 {
	// If we don't have enough bytes left in the buffer, generate more
	if r.bufferPos+4 > len(r.buffer) {
		r.getNextBytes()
	}

	// Extract 4 bytes from the buffer
	value := binary.BigEndian.Uint32(r.buffer[r.bufferPos : r.bufferPos+4])
	r.bufferPos += 4
	return value
}

// Uint32N returns a random uint32 in the range [0, n)
func (r *DeterministicRand) Uint32N(n uint32) uint32 {
	if n == 0 {
		return 0
	}

	// To avoid modulo bias, we need to create a mask of bits that will work
	// for our range and reject values outside of it

	// Find the smallest power of 2 greater than n
	mask := uint32(1)
	for mask < n {
		mask <<= 1
	}
	mask-- // Create a mask with all bits set up to the required number of bits

	// Keep generating random numbers until we get one in the range [0, n)
	for {
		// Generate a random number and apply the mask
		val := r.Uint32() & mask

		// If the value is less than n, return it
		if val < n {
			return val
		}
		// Otherwise, try again
	}
}
