package chunkgroup

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/lightnode"
	"golang.org/x/exp/rand"
	"time"
)

// Map keeps track of light nodes and their chunk group assignments.
type Map struct {
	// A map from light node ID to light node data.
	lightNodes map[uint64]*assignment

	// A map from chunk group ID to a list of light nodes in that chunk group.
	chunkGroups map[uint64][]*lightnode.Registration

	// Light node registrations are stored in this queue. The next light node to be shuffled is always at the front.
	shuffleQueue *assignmentQueue

	// The number of chunk groups.
	chunkGroupCount uint64

	// The genesis time of the protocol.
	genesis time.Time

	// The shuffle period of the protocol.
	shufflePeriod time.Duration
}

// NewMap creates a new Map.
func NewMap(
	chunkGroupCount uint64,
	genesis time.Time,
	shufflePeriod time.Duration) Map {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	return Map{
		lightNodes:      make(map[uint64]*assignment),
		chunkGroups:     make(map[uint64][]*lightnode.Registration),
		shuffleQueue:    newAssignmentQueue(),
		chunkGroupCount: chunkGroupCount,
		genesis:         genesis,
		shufflePeriod:   shufflePeriod,
	}
}

// Add adds a light node to be tracked by the map.
func (m *Map) Add(now time.Time, registration *lightnode.Registration) {
	shuffleOffset := ComputeShuffleOffset(registration.Seed(), m.shufflePeriod)
	epoch := ComputeShuffleEpoch(m.genesis, m.shufflePeriod, shuffleOffset, now)
	chunkGroup := ComputeChunkGroup(registration.Seed(), epoch, m.chunkGroupCount)
	startOfEpoch := ComputeStartOfShuffleEpoch(m.genesis, m.shufflePeriod, shuffleOffset, epoch)
	endOfEpoch := ComputeEndOfShuffleEpoch(m.genesis, m.shufflePeriod, shuffleOffset, epoch)

	entry := newAssignment(registration, shuffleOffset, chunkGroup, startOfEpoch, endOfEpoch)

	m.addToChunkGroupMap(chunkGroup, registration)
	m.lightNodes[registration.ID()] = entry
	m.shuffleQueue.Push(entry)
}

// Remove removes a light node from the map. This is a no-op if the light node is not being tracked.
func (m *Map) Remove(lightNodeID uint64) {
	entry, ok := m.lightNodes[lightNodeID]
	if !ok {
		return
	}

	m.shuffleQueue.Remove(lightNodeID)
	delete(m.lightNodes, lightNodeID)
	m.removeFromChunkGroupMap(entry.chunkGroup, entry.registration)
}

// Get returns the light node registration with the given ID, or nil if no such light node is being tracked.
func (m *Map) Get(lightNodeID uint64) *lightnode.Registration {
	entry, ok := m.lightNodes[lightNodeID]
	if !ok {
		return nil
	}
	return entry.registration
}

// GetChunkGroup returns the current chunk group of the light node with the given ID. The second return value
// is true if the light node is being tracked, and false otherwise.
func (m *Map) GetChunkGroup(now time.Time, lightNodeID uint64) (uint64, bool) {
	m.shuffle(now)

	entry, ok := m.lightNodes[lightNodeID]
	if !ok {
		return 0, false
	}

	return entry.chunkGroup, true
}

// Size returns the number of light nodes in the map.
func (m *Map) Size() uint {
	return uint(len(m.lightNodes))
}

// GetNodesInChunkGroup returns all light nodes in the given chunk group.
func (m *Map) GetNodesInChunkGroup(
	now time.Time,
	chunkGroup uint64) []*lightnode.Registration {

	m.shuffle(now)

	nodes := m.chunkGroups[chunkGroup]
	nodesCopy := make([]*lightnode.Registration, len(nodes))
	copy(nodesCopy, nodes)

	return nodesCopy
}

// GetRandomNode returns a random light node in the given chunk group. If minimumTimeInGroup is
// non-zero, the light node must have been in the chunk group for at least that amount of time. Returns nil
// if no light node is found that satisfies the constraints.
func (m *Map) GetRandomNode(
	now time.Time,
	rand *rand.Rand,
	chunkGroup uint64,
	minimumTimeInGroup time.Duration) (*lightnode.Registration, bool) {

	if chunkGroup >= m.chunkGroupCount {
		panic(fmt.Sprintf("chunk group %d is out of bounds, there are only %d chunk groups",
			chunkGroup, m.chunkGroupCount))
	}

	m.shuffle(now)

	nodes := m.chunkGroups[chunkGroup]
	var filteredNodes []*lightnode.Registration

	if minimumTimeInGroup == 0 {
		filteredNodes = nodes
	} else {
		filteredNodes = make([]*lightnode.Registration, 0, len(nodes))
		for _, node := range nodes {
			entry := m.lightNodes[node.ID()]
			timeInGroup := now.Sub(entry.startOfEpoch)
			if timeInGroup >= minimumTimeInGroup {
				filteredNodes = append(filteredNodes, node)
			}
		}
	}

	if len(filteredNodes) == 0 {
		return nil, false
	}

	index := rand.Intn(len(nodes))
	node := nodes[index]
	return node, true
}

// shuffle shuffles the light nodes into new chunk groups given the current time.
func (m *Map) shuffle(now time.Time) {
	if m.Size() == 0 {
		return
	}

	// As a sanity check, ensure that we don't shuffle each light node more than once during this call.
	shufflesRemaining := int(m.Size()) + 1

	for {
		shufflesRemaining--
		if shufflesRemaining < 0 {
			panic("too many shuffles")
		}

		next := m.shuffleQueue.Peek()
		if next.endOfEpoch.After(now) {
			// The next light node is not yet ready to be shuffled.
			break
		}
		m.shuffleQueue.Pop()

		newEpoch := ComputeShuffleEpoch(m.genesis, m.shufflePeriod, next.shuffleOffset, now)
		newChunkGroup := ComputeChunkGroup(next.registration.Seed(), newEpoch, m.chunkGroupCount)
		startOfEpoch := ComputeStartOfShuffleEpoch(m.genesis, m.shufflePeriod, next.shuffleOffset, newEpoch)
		endOfEpoch := ComputeEndOfShuffleEpoch(m.genesis, m.shufflePeriod, next.shuffleOffset, newEpoch)

		previousChunkGroup := next.chunkGroup
		next.chunkGroup = newChunkGroup
		next.startOfEpoch = startOfEpoch
		next.endOfEpoch = endOfEpoch

		if newChunkGroup != previousChunkGroup {
			m.removeFromChunkGroupMap(previousChunkGroup, next.registration)
			m.addToChunkGroupMap(newChunkGroup, next.registration)
		}

		m.shuffleQueue.Push(next)
	}
}

// addToChunkGroupMap adds a light node to the given chunk group.
func (m *Map) addToChunkGroupMap(chunkGroup uint64, registration *lightnode.Registration) {
	oldGroup := m.chunkGroups[chunkGroup]
	newGroup := append(oldGroup, registration)
	m.chunkGroups[chunkGroup] = newGroup
}

// removeFromChunkGroupMap removes a light node from the given chunk group.
func (m *Map) removeFromChunkGroupMap(chunkGroup uint64, registration *lightnode.Registration) {
	// TODO this is not efficient, refactor to do this in O(1) time

	oldGroup := m.chunkGroups[chunkGroup]
	for i, node := range oldGroup {
		if node.ID() == registration.ID() {
			newGroup := append(oldGroup[:i], oldGroup[i+1:]...)
			m.chunkGroups[chunkGroup] = newGroup
			return
		}
	}
}
