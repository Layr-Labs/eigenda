package lightnode

import (
	"fmt"
	"math/rand"
	"time"
)

// ChunkGroupMap keeps track of light nodes and their chunk group assignments.
type ChunkGroupMap struct {
	// A map from light node ID to light node data.
	lightNodes map[uint64]*chunkGroupAssignment

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
		lightNodes:      make(map[uint64]*chunkGroupAssignment),
		chunkGroupCount: chunkGroupCount,
		genesis:         genesis,
		shufflePeriod:   shufflePeriod,
	}
}

// Add adds a light node to be tracked by the map.
func (cgm *ChunkGroupMap) Add(now time.Time, registration *Registration) {
	shuffleOffset := ComputeShuffleOffset(registration.Seed(), cgm.shufflePeriod)
	epoch := ComputeShuffleEpoch(cgm.genesis, cgm.shufflePeriod, shuffleOffset, now)
	chunkGroup := ComputeChunkGroup(registration.Seed(), epoch, cgm.chunkGroupCount)
	nextShuffleTime := ComputeEndOfShuffleEpoch(cgm.genesis, cgm.shufflePeriod, shuffleOffset, epoch)

	assignment := newChunkGroupAssignment(registration, shuffleOffset, chunkGroup, nextShuffleTime)

	cgm.addToChunkGroupMap(chunkGroup, registration)
	cgm.lightNodes[registration.ID()] = assignment
	cgm.shuffleQueue.Push(assignment)
}

// Remove removes a light node from the map. This is a no-op if the light node is not being tracked.
func (cgm *ChunkGroupMap) Remove(lightNodeID uint64) {
	assignment, ok := cgm.lightNodes[lightNodeID]
	if !ok {
		return
	}

	cgm.shuffleQueue.Remove(lightNodeID)
	delete(cgm.lightNodes, lightNodeID)
	cgm.removeFromChunkGroupMap(assignment.chunkGroup, assignment.registration)
}

// Get returns the light node registration with the given ID, or nil if no such light node is being tracked.
func (cgm *ChunkGroupMap) Get(lightNodeID uint64) *Registration {
	assignment, ok := cgm.lightNodes[lightNodeID]
	if !ok {
		return nil
	}
	return assignment.registration
}

// GetChunkGroup returns the current chunk group of the light node with the given ID.
func (cgm *ChunkGroupMap) getChunkGroup(now time.Time, lightNodeID uint64) (uint64, error) {
	cgm.shuffle(now)

	assignment, ok := cgm.lightNodes[lightNodeID]
	if !ok {
		return 0, fmt.Errorf("light node with ID %d not found", lightNodeID)
	}

	return assignment.chunkGroup, nil
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
	if cgm.Size() == 0 {
		return
	}

	// As a sanity check, ensure that we don't shuffle each light node more than once during this call.
	shufflesRemaining := cgm.Size()

	for {
		shufflesRemaining--
		if shufflesRemaining < 0 {
			panic("too many shuffles")
		}

		next := cgm.shuffleQueue.Peek()
		if next.nextShuffleTime.After(now) {
			// The next light node is not yet ready to be shuffled.
			break
		}
		cgm.shuffleQueue.Pop()

		newEpoch := ComputeShuffleEpoch(cgm.genesis, cgm.shufflePeriod, next.shuffleOffset, now)
		newChunkGroup := ComputeChunkGroup(next.registration.Seed(), newEpoch, cgm.chunkGroupCount)
		nextShuffleTime := ComputeEndOfShuffleEpoch(cgm.genesis, cgm.shufflePeriod, next.shuffleOffset, newEpoch)

		previousChunkGroup := next.chunkGroup
		next.chunkGroup = newChunkGroup
		next.nextShuffleTime = nextShuffleTime

		if newChunkGroup != previousChunkGroup {
			cgm.removeFromChunkGroupMap(previousChunkGroup, next.registration)
			cgm.addToChunkGroupMap(newChunkGroup, next.registration)
		}

		cgm.shuffleQueue.Push(next)
	}
}

// addToChunkGroupMap adds a light node to the given chunk group.
func (cgm *ChunkGroupMap) addToChunkGroupMap(chunkGroup uint64, registration *Registration) {
	oldGroup := cgm.chunkGroups[chunkGroup]
	newGroup := append(oldGroup, registration)
	cgm.chunkGroups[chunkGroup] = newGroup
}

// removeFromChunkGroupMap removes a light node from the given chunk group.
func (cgm *ChunkGroupMap) removeFromChunkGroupMap(chunkGroup uint64, registration *Registration) {
	// TODO this is not efficient, refactor to do this in O(1) time

	oldGroup := cgm.chunkGroups[chunkGroup]
	for i, node := range oldGroup {
		if node.ID() == registration.ID() {
			newGroup := append(oldGroup[:i], oldGroup[i+1:]...)
			cgm.chunkGroups[chunkGroup] = newGroup
			return
		}
	}
}
