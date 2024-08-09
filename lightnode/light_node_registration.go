package lightnode

import (
	"fmt"
	"time"
)

// LightNodeRegistration describes publicly known information about a light node.
// This information is registered on-chain.
type LightNodeRegistration struct {
	// The ID of the light node.
	ID uint64

	// A seed assigned to the light node when it was registered.
	// Used for deterministically random operations involving this light node.
	Seed uint64

	// TODO join time
	// TODO public key
	// TODO payment address
}

// String returns a string representation of the light node.
func (ln *LightNodeRegistration) String() string {
	return fmt.Sprintf("LightNode{ID: %d, Seed: %v}", ln.ID, ln.Seed)
}

// CurrentChunkGroup returns the chunk group of the light node is currently assigned to.
func (ln *LightNodeRegistration) CurrentChunkGroup(
	genesis time.Time,
	shufflePeriod time.Duration,
	now time.Time,
	chunkGroupCount uint64) (uint64, error) {

	// TODO can we cache this value?
	offset, err := ComputeShuffleOffset(ln.Seed, shufflePeriod)
	if err != nil {
		return 0, err
	}

	epoch, err := ComputeShuffleEpoch(genesis, shufflePeriod, offset, now)
	if err != nil {
		return 0, err
	}

	return ComputeChunkGroup(ln.Seed, epoch, chunkGroupCount), nil
}
