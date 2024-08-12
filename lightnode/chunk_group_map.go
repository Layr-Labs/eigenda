package lightnode

import (
	"math/rand"
	"time"
)

// ChunkGroupMap keeps track of light nodes and their chunk group assignments.
type ChunkGroupMap struct {
	// A map from light node ID to light node data.
	lightNodes map[uint64]*Registration

	// A map from chunk group ID to a list of light nodes in that chunk group.
	chunkGroups map[uint64][]*Registration

	// Light node registrations are stored in this queue. The next light node to be shuffled is always at the front.
	shuffleQueue *assignmentQueue

	// The number of chunk groups.
	chunkGroupCount uint64

	// The genesis time of the protocol.
	genesis time.Time

	// The shuffle period of the protocol.
	shufflePeriod time.Duration
}

// NewChunkGroupMap creates a new ChunkGroupMap.
func NewChunkGroupMap(
	chunkGroupCount uint64,
	genesis time.Time,
	shufflePeriod time.Duration) ChunkGroupMap {

	return ChunkGroupMap{
		lightNodes:      make(map[uint64]*Registration),
		chunkGroupCount: chunkGroupCount,
		genesis:         genesis,
		shufflePeriod:   shufflePeriod,
	}
}

// Add adds a light node to be tracked by the map.
func (cgm *ChunkGroupMap) Add(now time.Time, registration *Registration) {
	cgm.lightNodes[registration.ID()] = registration

	shuffleOffset := ComputeShuffleOffset(registration.Seed(), cgm.shufflePeriod)
	epoch := ComputeShuffleEpoch(cgm.genesis, cgm.shufflePeriod, shuffleOffset, now)

	chunkGroup := ComputeChunkGroup(registration.Seed(), epoch, cgm.chunkGroupCount)
	cgm.chunkGroups[chunkGroup] = append(cgm.chunkGroups[chunkGroup], registration)

	nextShuffleTime := ComputeNextShuffleTime(cgm.genesis, cgm.shufflePeriod, shuffleOffset, epoch)

	assignment := newChunkGroupAssignment(registration, shuffleOffset, chunkGroup, nextShuffleTime)
	cgm.shuffleQueue.Push(assignment)

}

// Remove removes a light node from the map.
func (cgm *ChunkGroupMap) Remove(lightNodeID uint64) {
	cgm.shuffleQueue.Remove(lightNodeID)
	// TODO remove from the chunkGroups map
	delete(cgm.lightNodes, lightNodeID)
}

// Get returns the light node with the given ID.
func (cgm *ChunkGroupMap) Get(lightNodeID uint64) (*Registration, bool) {
	registration, ok := cgm.lightNodes[lightNodeID]
	return registration, ok
}

// Size returns the number of light nodes in the map.
func (cgm *ChunkGroupMap) Size() int {
	return len(cgm.lightNodes)
}

// GetLightNodesInChunkGroup returns all light nodes in the given chunk group.
func (cgm *ChunkGroupMap) GetLightNodesInChunkGroup(
	now time.Time,
	chunkGroup uint64) []*Registration {

	cgm.shuffle(now)

	nodes := cgm.chunkGroups[chunkGroup]
	nodesCopy := make([]*Registration, len(nodes))
	copy(nodesCopy, nodes)

	return nodesCopy
}

// GetRandomNode returns a random light node in the given chunk group. If minimumTimeInGroup is
// non-zero, the light node must have been in the chunk group for at least that amount of time. Returns nil
// if no light node is found that satisfies the constraints.
func (cgm *ChunkGroupMap) GetRandomNode(
	now time.Time,
	rand *rand.Rand,
	chunkGroup uint64,
	minimumTimeInGroup time.Duration) (Registration, bool) {

	cgm.shuffle(now)

	return Registration{}, false // TODO
}

// shuffle shuffles the light nodes into new chunk groups given the current time.
func (cgm *ChunkGroupMap) shuffle(now time.Time) {
	// TODO
}
