package lightnode

import (
	"fmt"
	"time"
)

// LightNodeData describes publicly known information about a light node.
// This information is registered on-chain.
type LightNodeData struct {
	// The ID of the light node.
	id uint64

	// A seed assigned to the light node when it was registered.
	// Used for deterministically random operations involving this light node.
	seed uint64

	// The time at which the light node was initially registered.
	registrationTime time.Time

	// TODO public key
	// TODO payment address
}

// NewLightNodeData creates a new LightNodeData instance.
func NewLightNodeData(id uint64, seed uint64, registrationTime time.Time) LightNodeData {
	return LightNodeData{
		id:               id,
		seed:             seed,
		registrationTime: registrationTime,
	}
}

// ID returns the ID of the light node.
func (ln *LightNodeData) ID() uint64 {
	return ln.id
}

// Seed returns the seed of the light node.
func (ln *LightNodeData) Seed() uint64 {
	return ln.seed
}

// RegistrationTime returns the time at which the light node was registered.
func (ln *LightNodeData) RegistrationTime() time.Time {
	return ln.registrationTime
}

// String returns a string representation of the light node.
func (ln *LightNodeData) String() string {
	return fmt.Sprintf("LightNode{ID: %d, Seed: %v, Registration time: %s}", ln.id, ln.seed, ln.registrationTime)
}
