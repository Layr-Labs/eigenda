package lightnode

import "time"

// chunkGroupAssignment is a struct that holds a registration and the chunk group it is currently assigned to.
type chunkGroupAssignment struct {

	// registration contains publicly known information about a light node that is registered on-chain.
	registration *Registration

	// shuffleOffset is the offset at which a light node should be shuffled into a new chunk group relative
	// the beginning of each shuffle interval. This is a function of the light node's seed and the shuffle period
	// and does not change, so we cache it here.
	shuffleOffset time.Duration

	// chunkGroup is the chunk group that the light node is currently assigned to.
	chunkGroup uint64

	// startOfEpoch is the start of the current shuffle epoch,
	// i.e. the time when this light node was last shuffled into the current chunk group.
	startOfEpoch time.Time

	// endOfEpoch is the end of the current shuffle epoch,
	// i.e. the next time when this light node will be shuffled into a new chunk group.
	endOfEpoch time.Time
}

// newChunkGroupAssignment creates a new chunkGroupAssignment.
func newChunkGroupAssignment(
	registration *Registration,
	shuffleOffset time.Duration,
	chunkGroup uint64,
	startOfEpoch time.Time,
	endOfEpoch time.Time) *chunkGroupAssignment {

	return &chunkGroupAssignment{
		registration:  registration,
		shuffleOffset: shuffleOffset,
		chunkGroup:    chunkGroup,
		startOfEpoch:  startOfEpoch,
		endOfEpoch:    endOfEpoch,
	}
}
