package lightnode

import (
	"fmt"
	"time"
)

// Registration describes publicly known information about a light node.
// This information is registered on-chain.
type Registration struct {
	// The ID of the light node.
	id int64

	// A seed assigned to the light node when it was registered.
	// Used for deterministically random operations involving this light node.
	seed int64

	// The time at which the light node was initially registered.
	registrationTime time.Time

	// FUTURE WORK: add public key, payment address, etc.
}

// NewRegistration creates a new Registration instance.
func NewRegistration(id int64, seed int64, registrationTime time.Time) *Registration {
	return &Registration{
		id:               id,
		seed:             seed,
		registrationTime: registrationTime,
	}
}

// ID returns the ID of the light node.
func (registration *Registration) ID() int64 {
	return registration.id
}

// Seed returns the seed of the light node.
func (registration *Registration) Seed() int64 {
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
