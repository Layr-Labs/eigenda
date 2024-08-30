package chunkgroup

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/lightnode"
	"time"
)

// assignmentKey uniquely identifies a light node chunkGroupAssignment.
type assignmentKey struct {
	lightNodeID     uint64
	assignmentIndex uint32
}

// A set of light node assignments.
type assignmentSet map[assignmentKey]*chunkGroupAssignment

// Map keeps track of light nodes and their chunk group assignments.
type Map struct {
	// A map from light node ID to a list of assignments.
	assignmentMap map[uint64][]*chunkGroupAssignment

	// A map from light node ID to the light node registration.
	lightNodes map[uint64]*lightnode.Registration

	// A map from chunk group ID to a set of light nodes in that chunk group.
	chunkGroups map[uint32]assignmentSet

	// Light node registrations are stored in this queue. The next light node to be shuffled is always at the front.
	shuffleQueue *assignmentQueue

	// The number of chunk groups.
	chunkGroupCount uint32

	// The number of chunk groups each light node is assigned to.
	assignmentCount uint32

	// The shuffle period of the protocol.
	shufflePeriod time.Duration
}

// NewMap creates a new Map.
func NewMap(
	chunkGroupCount uint32,
	assignmentCount uint32,
	shufflePeriod time.Duration) Map {

	if shufflePeriod <= 0 {
		panic(fmt.Sprintf("shuffle period must be positive, got %s", shufflePeriod))
	}

	return Map{
		assignmentMap:   make(map[uint64][]*chunkGroupAssignment),
		lightNodes:      make(map[uint64]*lightnode.Registration),
		chunkGroups:     make(map[uint32]assignmentSet),
		shuffleQueue:    newAssignmentQueue(),
		chunkGroupCount: chunkGroupCount,
		assignmentCount: assignmentCount,
		shufflePeriod:   shufflePeriod,
	}
}

// Add adds a light node to be tracked by the map.
func (m *Map) Add(now time.Time, registration *lightnode.Registration) {
	m.lightNodes[registration.ID()] = registration
	assignments := make([]*chunkGroupAssignment, 0, m.assignmentCount)

	for assignmentIndex := uint32(0); assignmentIndex < m.assignmentCount; assignmentIndex++ {
		// TODO can some of this logic be de-duplicated?

		assignmentSeed := rotateLeft(registration.Seed(), assignmentIndex)
		shuffleOffset := ComputeShuffleOffset(assignmentSeed, m.shufflePeriod)
		epoch := ComputeShuffleEpoch(m.shufflePeriod, shuffleOffset, now)
		chunkGroup := ComputeChunkGroup(assignmentSeed, epoch, m.chunkGroupCount)
		startOfEpoch := ComputeStartOfShuffleEpoch(m.shufflePeriod, shuffleOffset, epoch)
		endOfEpoch := ComputeEndOfShuffleEpoch(m.shufflePeriod, shuffleOffset, epoch)

		entry := newChunkGroupAssignment(
			registration,
			assignmentIndex,
			assignmentSeed,
			shuffleOffset,
			chunkGroup,
			startOfEpoch,
			endOfEpoch)

		m.addToChunkGroupMap(chunkGroup, registration.ID(), assignmentIndex)
		assignments = append(assignments, entry)
		m.shuffleQueue.Push(entry)
	}

	m.assignmentMap[registration.ID()] = assignments
}

// Remove removes a light node from the map. This is a no-op if the light node is not being tracked.
func (m *Map) Remove(lightNodeID uint64) {
	assignments, ok := m.assignmentMap[lightNodeID]
	if !ok {
		return
	}

	m.shuffleQueue.Remove(lightNodeID)
	delete(m.assignmentMap, lightNodeID)
	delete(m.lightNodes, lightNodeID)

	for _, entry := range assignments {
		m.removeFromChunkGroupMap(entry.chunkGroup, entry.registration.ID(), entry.assignmentIndex)
	}
}

// Get returns the light node registration with the given ID, or nil if no such light node is being tracked.
func (m *Map) Get(lightNodeID uint64) *lightnode.Registration {
	registration, ok := m.lightNodes[lightNodeID]
	if !ok {
		return nil
	}
	return registration
}

// GetChunkGroups returns the current chunk groups of the light node with the given ID. The second return value
// is true if the light node is being tracked, and false otherwise.
func (m *Map) GetChunkGroups(now time.Time, lightNodeID uint64) ([]uint32, bool) {
	m.shuffle(now)

	assignments, ok := m.assignmentMap[lightNodeID]
	if !ok {
		return nil, false
	}

	chunkGroupSet := make(map[uint32]bool)
	for _, entry := range assignments {
		chunkGroupSet[entry.chunkGroup] = true
	}

	chunkGroupList := make([]uint32, 0, m.assignmentCount)
	for chunkGroup := range chunkGroupSet {
		chunkGroupList = append(chunkGroupList, chunkGroup)
	}

	return chunkGroupList, true
}

// Size returns the number of light nodes in the map.
func (m *Map) Size() uint {
	return uint(len(m.lightNodes))
}

// GetNodesInChunkGroup returns the IDs of the light nodes in the given chunk group.
func (m *Map) GetNodesInChunkGroup(
	now time.Time,
	chunkGroup uint32) []uint64 {

	if chunkGroup >= m.chunkGroupCount {
		panic(fmt.Sprintf("chunk group %d is out of bounds, there are only %d chunk groups",
			chunkGroup, m.chunkGroupCount))
	}

	m.shuffle(now)

	assignments := m.chunkGroups[chunkGroup]
	uniqueNodes := make(map[uint64]bool)
	for key := range assignments {
		uniqueNodes[key.lightNodeID] = true
	}

	nodeList := make([]uint64, 0, len(uniqueNodes))
	for node := range uniqueNodes {
		nodeList = append(nodeList, node)
	}

	return nodeList
}

// GetRandomNode returns a random light node in the given chunk group. If minimumTimeInGroup is
// non-zero, the light node must have been in the chunk group for at least that amount of time. Returns nil
// if no light node is found that satisfies the constraints.
func (m *Map) GetRandomNode(
	now time.Time,
	chunkGroup uint32,
	minimumTimeInGroup time.Duration) (*lightnode.Registration, bool) {

	if chunkGroup >= m.chunkGroupCount {
		panic(fmt.Sprintf("chunk group %d is out of bounds, there are only %d chunk groups",
			chunkGroup, m.chunkGroupCount))
	}

	m.shuffle(now)

	assignments := m.chunkGroups[chunkGroup]

	// Collect unique assignments that have been in the group for at least minimumTimeInGroup.
	qualifiedAssignments := map[uint64] /* node ID */ *chunkGroupAssignment{}
	for key := range assignments {
		assignment := assignments[key]
		notYetPresent := qualifiedAssignments[key.lightNodeID] == nil
		meetsTimeRequirement := minimumTimeInGroup == 0 || now.Sub(assignment.startOfEpoch) >= minimumTimeInGroup
		if notYetPresent && meetsTimeRequirement {
			qualifiedAssignments[key.lightNodeID] = assignment
		}
	}

	for entry := range qualifiedAssignments {
		// golang map iteration starts at a random position, so we can return the first node we find
		return qualifiedAssignments[entry].registration, true
	}
	return nil, false
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

		newEpoch := ComputeShuffleEpoch(m.shufflePeriod, next.shuffleOffset, now)
		newChunkGroup := ComputeChunkGroup(next.assignmentSeed, newEpoch, m.chunkGroupCount)
		startOfEpoch := ComputeStartOfShuffleEpoch(m.shufflePeriod, next.shuffleOffset, newEpoch)
		endOfEpoch := ComputeEndOfShuffleEpoch(m.shufflePeriod, next.shuffleOffset, newEpoch)

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
func (m *Map) addToChunkGroupMap(chunkGroup uint32, lightNodeId uint64, assignmentIndex uint32) {
	group := m.chunkGroups[chunkGroup]
	if group == nil {
		group = make(assignmentSet)
		m.chunkGroups[chunkGroup] = group
	}
	key := assignmentKey{
		lightNodeID:     lightNodeId,
		assignmentIndex: assignmentIndex,
	}
	group[key] = true
}

// removeFromChunkGroupMap removes a light node from the given chunk group.
func (m *Map) removeFromChunkGroupMap(chunkGroup uint32, lightNodeId uint64, assignmentIndex uint32) {
	group := m.chunkGroups[chunkGroup]
	key := assignmentKey{
		lightNodeID:     lightNodeId,
		assignmentIndex: assignmentIndex,
	}
	delete(group, key)
}
