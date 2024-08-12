package lightnode

import (
	"fmt"
	"time"
)

// Registration describes publicly known information about a light node.
// This information is registered on-chain.
type Registration struct {
	// The ID of the light node.
	id uint64

	// A seed assigned to the light node when it was registered.
	// Used for deterministically random operations involving this light node.
	seed uint64

	// The time at which the light node was initially registered.
	registrationTime time.Time

	// FUTURE WORK: add public key, payment address, etc.
}

// NewRegistration creates a new Registration instance.
func NewRegistration(id uint64, seed uint64, registrationTime time.Time) *Registration {
	return &Registration{
		id:               id,
		seed:             seed,
		registrationTime: registrationTime,
	}
}

// ID returns the ID of the light node.
func (registration *Registration) ID() uint64 {
	return registration.id
}

// Seed returns the seed of the light node.
func (registration *Registration) Seed() uint64 {
	return registration.seed
}

// RegistrationTime returns the time at which the light node was registered.
func (registration *Registration) RegistrationTime() time.Time {
	return registration.registrationTime
}

// String returns a string representation of the light node.
func (registration *Registration) String() string {
	return fmt.Sprintf("LightNode{ID: %d, Seed: %v, Registration time: %s}", registration.id, registration.seed, registration.registrationTime)
}

// GetChunkGroup returns the chunk group that the light node is in.
func (registration *Registration) GetChunkGroup(
	now time.Time,
	genesis time.Time,
	shufflePeriod time.Duration,
	chunkGroupCount uint64) uint64 {

	shuffleOffset := ComputeShuffleOffset(registration.seed, shufflePeriod)
	epoch := ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, now)
	return ComputeChunkGroup(registration.seed, epoch, chunkGroupCount)
}
