package chunkgroup

import (
	"github.com/Layr-Labs/eigenda/lightnode"
	"time"
)

// assignmentKey uniquely identifies a light node chunkGroupAssignment.
type assignmentKey struct {
	lightNodeID     uint64
	assignmentIndex uint32
}

// chunkGroupAssignment is a struct that holds a registration and the chunk group it is currently assigned to.
type chunkGroupAssignment struct {

	// registration contains publicly known information about a light node that is registered on-chain.
	registration *lightnode.Registration

	// assignmentIndex describes which of a light node's multiple groups this struct represents.
	// The first of a light node's groups has an chunkGroupAssignment index of 0, the second has
	// an index of 1, and so on.
	assignmentIndex uint32

	// assignmentKey is the key that uniquely identifies this chunkGroupAssignment. Note that this information
	// is also stored in the registration, but we cache this object here for convenience.
	key assignmentKey

	// assignmentSeed is the seed used for this group chunkGroupAssignment.
	//
	// This value is deterministic and does not change, so we cache it here.
	assignmentSeed uint64

	// shuffleOffset is the offset at which this group chunkGroupAssignment should be shuffled into a
	// new chunk group relative the beginning of each shuffle interval.
	//
	// This value is deterministic and does not change, so we cache it here.
	shuffleOffset time.Duration

	// chunkGroup is the chunk group currently associated with this chunkGroupAssignment index.
	chunkGroup uint32

	// startOfEpoch is the start of the current shuffle epoch,
	// i.e. the time when this chunkGroupAssignment index was last shuffled into the current chunk group.
	startOfEpoch time.Time

	// endOfEpoch is the end of the current shuffle epoch,
	// i.e. the next time when this chunkGroupAssignment index will be shuffled into a new chunk group.
	endOfEpoch time.Time
}

// newChunkGroupAssignment creates a new chunkGroupAssignment.
func newChunkGroupAssignment(
	registration *lightnode.Registration,
	assignmentIndex uint32,
	assignmentSeed uint64,
	shuffleOffset time.Duration,
	chunkGroup uint32,
	startOfEpoch time.Time,
	endOfEpoch time.Time) *chunkGroupAssignment {

	return &chunkGroupAssignment{
		registration:    registration,
		assignmentIndex: assignmentIndex,
		key: assignmentKey{
			lightNodeID:     registration.ID(),
			assignmentIndex: assignmentIndex,
		},
		assignmentSeed: assignmentSeed,
		shuffleOffset:  shuffleOffset,
		chunkGroup:     chunkGroup,
		startOfEpoch:   startOfEpoch,
		endOfEpoch:     endOfEpoch,
	}
}
