package chunkgroup

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/lightnode"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sort"
	"testing"
	"time"
)

func randomRegistration() *lightnode.Registration {
	return lightnode.NewRegistration(rand.Uint64(), rand.Uint64(), tu.RandomTime())
}

// TODO test that deletes things and ensures there is no garbage left in maps
// TODO take into account when a light node was registered

func TestAddRemoveGetOneAssignment(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(1)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(startTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint32(i+1), cgMap.Size())
	}

	// Removing non-existent elements should be a no-op.
	for i := 0; i < 10; i++ {
		cgMap.Remove(rand.Uint64())
		assert.Equal(t, uint32(elementsToAdd), cgMap.Size())
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedRegistrations {
		assert.Equal(t, registration, cgMap.Get(id))
	}

	// Remove all nodes that have an ID divisible by 2.
	removalCount := 0
	for id, registration := range expectedRegistrations {
		if id%2 == 0 {
			assert.Equal(t, registration, cgMap.Get(id))
			cgMap.Remove(id)
			assert.Nil(t, cgMap.Get(id))
			removalCount++
			assert.Equal(t, uint32(elementsToAdd-removalCount), cgMap.Size())
		}
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedRegistrations {
		if id%2 == 0 {
			assert.Nil(t, cgMap.Get(id))
		} else {
			assert.Equal(t, registration, cgMap.Get(id))
		}
	}
}

func TestAddRemoveGetMultipleAssignments(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(rand.Intn(3) + 2)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(startTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint32(i+1), cgMap.Size())
	}

	// Removing non-existent elements should be a no-op.
	for i := 0; i < 10; i++ {
		cgMap.Remove(rand.Uint64())
		assert.Equal(t, uint32(elementsToAdd), cgMap.Size())
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedRegistrations {
		assert.Equal(t, registration, cgMap.Get(id))
	}

	// Remove all nodes that have an ID divisible by 2.
	removalCount := 0
	for id, registration := range expectedRegistrations {
		if id%2 == 0 {
			assert.Equal(t, registration, cgMap.Get(id))
			cgMap.Remove(id)
			assert.Nil(t, cgMap.Get(id))
			removalCount++
			assert.Equal(t, uint32(elementsToAdd-removalCount), cgMap.Size())
		}
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedRegistrations {
		if id%2 == 0 {
			assert.Nil(t, cgMap.Get(id))
		} else {
			assert.Equal(t, registration, cgMap.Get(id))
		}
	}
}

func TestChunkGroupCalculationsSingleAssignment(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(1)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration
		cgMap.Add(startTime, registration)
	}

	now := startTime
	steps := 100
	for step := 0; step < steps; step++ {
		if rand.Float64() < (1.0 / 3.0) {
			// Add less than a full shuffle period.
			now = now.Add(shufflePeriod * time.Duration(rand.Float64()))
		} else if rand.Float64() < (2.0 / 3.0) {
			// Add exactly one shuffle period.
			now = now.Add(shufflePeriod)
		} else {
			// Add several shuffle periods.
			now = now.Add(shufflePeriod * time.Duration(rand.Intn(10)+2))
		}

		// Add a few elements.
		numberToAdd := rand.Intn(10)
		count += numberToAdd
		for i := 0; i < numberToAdd; i++ {
			registration := randomRegistration()
			expectedRegistrations[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedRegistrations {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedRegistrations, key)
			numberToRemove--
		}

		// Verify the chunk group for each element.
		for id, registration := range expectedRegistrations {
			chunkGroups, ok := cgMap.GetChunkGroups(now, id)

			assert.Equal(t, 1, len(chunkGroups))
			chunkGroup := chunkGroups[0]

			assert.True(t, ok)
			offset := ComputeShuffleOffset(registration.Seed(), 0, shufflePeriod)
			epoch := ComputeShuffleEpoch(shufflePeriod, offset, now)
			expectedChunkGroup := ComputeChunkGroup(registration.Seed(), 0, epoch, chunkGroupCount)

			assert.Equal(t, expectedChunkGroup, chunkGroup)
		}

		// Query for full chunk groups.
		nodesReported := 0
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)
			nodesReported += len(chunk)

			for _, nodeId := range chunk {
				chunkGroups, ok := cgMap.GetChunkGroups(now, nodeId)
				assert.True(t, ok)
				assert.Equal(t, 1, len(chunkGroups))
				chunkGroup := chunkGroups[0]

				assert.Equal(t, chunkIndex, chunkGroup)
			}
		}

		assert.Equal(t, count, nodesReported)
	}
}

func TestChunkGroupCalculationsMultipleAssignments(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(rand.Intn(3) + 2)
	//assignmentCount := uint32(rand.Intn(3) + 2)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration
		cgMap.Add(startTime, registration)
	}

	now := startTime
	steps := 100
	for step := 0; step < steps; step++ {
		if rand.Float64() < (1.0 / 3.0) {
			// Add less than a full shuffle period.
			now = now.Add(shufflePeriod * time.Duration(rand.Float64()))
		} else if rand.Float64() < (2.0 / 3.0) {
			// Add exactly one shuffle period.
			now = now.Add(shufflePeriod)
		} else {
			// Add several shuffle periods.
			now = now.Add(shufflePeriod * time.Duration(rand.Intn(10)+2))
		}

		// Add a few elements.
		numberToAdd := rand.Intn(10)
		count += numberToAdd
		for i := 0; i < numberToAdd; i++ {
			registration := randomRegistration()
			expectedRegistrations[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedRegistrations {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedRegistrations, key)
			numberToRemove--
		}

		// A map from chunk group ID to set of nodes expected to be in that chunk group.
		expectedChunkMembership := make(map[uint32]map[uint64]bool)

		// Verify the chunk groups for each element.
		for id, registration := range expectedRegistrations {
			chunkGroups, _ := cgMap.GetChunkGroups(now, id)

			// Note that a light node may be assigned to the same chunk group multiple times,
			// resulting in a smaller number of chunk groups than the number of assignments.
			groupCount := len(chunkGroups)
			assert.True(t, groupCount > 0 && groupCount <= int(assignmentCount))

			uniqueExpectedChunkGroups := make(map[uint32]bool)
			for assignmentIndex := uint32(0); assignmentIndex < assignmentCount; assignmentIndex++ {
				shuffleOffset := ComputeShuffleOffset(registration.Seed(), assignmentIndex, shufflePeriod)
				epoch := ComputeShuffleEpoch(shufflePeriod, shuffleOffset, now)
				expectedChunkGroup := ComputeChunkGroup(registration.Seed(), assignmentIndex, epoch, chunkGroupCount)
				uniqueExpectedChunkGroups[expectedChunkGroup] = true
			}
			expectedChunkGroups := make([]uint32, 0, assignmentCount)
			for g := range uniqueExpectedChunkGroups {
				expectedChunkGroups = append(expectedChunkGroups, g)

				if _, ok := expectedChunkMembership[g]; !ok {
					expectedChunkMembership[g] = make(map[uint64]bool)
				}
				expectedChunkMembership[g][id] = true
			}

			// sort the lists to make comparison easier
			sort.Slice(chunkGroups, func(i, j int) bool {
				return chunkGroups[i] < chunkGroups[j]
			})
			sort.Slice(expectedChunkGroups, func(i, j int) bool {
				return expectedChunkGroups[i] < expectedChunkGroups[j]
			})
			assert.Equal(t, expectedChunkGroups, chunkGroups)
		}

		// Query GetNodesInChunkGroup()
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			nodes := cgMap.GetNodesInChunkGroup(now, chunkIndex)

			expectedNodes := make([]uint64, 0)
			for node := range expectedChunkMembership[chunkIndex] {
				expectedNodes = append(expectedNodes, node)
			}

			// Sort for easier comparison
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i] < nodes[j]
			})
			sort.Slice(expectedNodes, func(i, j int) bool {
				return expectedNodes[i] < expectedNodes[j]
			})

			assert.Equal(t, expectedNodes, nodes)
		}
	}
}

func TestGetRandomNodeSingleAssignment(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(1)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(startTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint32(i+1), cgMap.Size())
	}

	now := startTime.Add(shufflePeriod * time.Duration(rand.Float64()*1000))

	for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {

		chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)

		if len(chunk) == 0 {
			// There shouldn't be any nodes in the chunk group, so GetRandomNode shouldn't return anything.
			node := cgMap.GetRandomNode(now, chunkIndex, 0)
			assert.Nil(t, node)
			continue
		}

		for i := 0; i < 10; i++ {

			var minimumTimeInGroup time.Duration
			if rand.Float64() < 0.1 {
				minimumTimeInGroup = 0
			} else {
				minimumTimeInGroup = shufflePeriod / time.Duration(rand.Intn(5)+1)
			}

			randomNode := cgMap.GetRandomNode(now, chunkIndex, minimumTimeInGroup)

			if randomNode != nil {
				assert.Contains(t, chunk, randomNode.ID())
			} else {
				// there shouldn't be any nodes in the chunk group for the minimum time
				for _, nodeId := range chunk {
					registration := cgMap.Get(nodeId)

					offset := ComputeShuffleOffset(registration.Seed(), 0, shufflePeriod)
					epoch := ComputeShuffleEpoch(shufflePeriod, offset, startTime)
					epochBeginning := ComputeStartOfShuffleEpoch(shufflePeriod, offset, epoch)
					timeInGroup := now.Sub(epochBeginning)
					assert.True(t, timeInGroup >= minimumTimeInGroup)
				}
			}
		}
	}
}

func TestGetRandomNodeMultipleAssignments(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(rand.Intn(3) + 2)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(startTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint32(i+1), cgMap.Size())
	}

	now := startTime.Add(shufflePeriod * time.Duration(rand.Float64()*1000))

	for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {

		chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)

		if len(chunk) == 0 {
			// There shouldn't be any nodes in the chunk group, so GetRandomNode shouldn't return anything.
			node := cgMap.GetRandomNode(now, chunkIndex, 0)
			assert.Nil(t, node)
			continue
		}

		for i := 0; i < 10; i++ {

			var minimumTimeInGroup time.Duration
			if rand.Float64() < 0.1 {
				minimumTimeInGroup = 0
			} else {
				minimumTimeInGroup = shufflePeriod / time.Duration(rand.Intn(5)+1)
			}

			randomNode := cgMap.GetRandomNode(now, chunkIndex, minimumTimeInGroup)

			if randomNode != nil {
				assert.Contains(t, chunk, randomNode.ID())
			} else {
				// there shouldn't be any nodes in the chunk group for the minimum time
				for _, nodeId := range chunk {
					registration := cgMap.Get(nodeId)

					// We don't know which assignment index this corresponds to, so we just check all of them.
					for assignmentIndex := uint32(0); assignmentIndex < assignmentCount; assignmentIndex++ {
						offset := ComputeShuffleOffset(registration.Seed(), assignmentIndex, shufflePeriod)
						epoch := ComputeShuffleEpoch(shufflePeriod, offset, startTime)

						group := ComputeChunkGroup(registration.Seed(), assignmentIndex, epoch, chunkGroupCount)
						if group != chunkIndex {
							continue
						}

						epochBeginning := ComputeStartOfShuffleEpoch(shufflePeriod, offset, epoch)
						timeInGroup := now.Sub(epochBeginning)
						assert.True(t, timeInGroup >= minimumTimeInGroup)
					}
				}
			}
		}
	}
}

func TestSingleChunkGroupSingleAssignment(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(1)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration
		cgMap.Add(startTime, registration)
	}

	now := startTime
	steps := 10
	for step := 0; step < steps; step++ {
		if rand.Float64() < (1.0 / 3.0) {
			// Add less than a full shuffle period.
			now = now.Add(shufflePeriod * time.Duration(rand.Float64()))
		} else if rand.Float64() < (2.0 / 3.0) {
			// Add exactly one shuffle period.
			now = now.Add(shufflePeriod)
		} else {
			// Add several shuffle periods.
			now = now.Add(shufflePeriod * time.Duration(rand.Intn(10)+2))
		}

		// Add a few elements.
		numberToAdd := rand.Intn(10)
		count += numberToAdd
		for i := 0; i < numberToAdd; i++ {
			registration := randomRegistration()
			expectedRegistrations[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedRegistrations {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedRegistrations, key)
			numberToRemove--
		}

		// Verify the chunk group for each element.
		for id, registration := range expectedRegistrations {
			chunkGroups, ok := cgMap.GetChunkGroups(now, id)
			assert.True(t, ok)
			assert.Equal(t, 1, len(chunkGroups))
			chunkGroup := chunkGroups[0]

			offset := ComputeShuffleOffset(registration.Seed(), 0, shufflePeriod)
			epoch := ComputeShuffleEpoch(shufflePeriod, offset, now)
			expectedChunkGroup := ComputeChunkGroup(registration.Seed(), 0, epoch, chunkGroupCount)

			assert.Equal(t, expectedChunkGroup, chunkGroup)
		}

		// Query for full chunk groups.
		nodesReported := 0
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)
			nodesReported += len(chunk)

			for _, nodeID := range chunk {
				chunkGroups, ok := cgMap.GetChunkGroups(now, nodeID)
				assert.True(t, ok)
				assert.Equal(t, 1, len(chunkGroups))
				chunkGroup := chunkGroups[0]

				assert.Equal(t, chunkIndex, chunkGroup)
			}
		}

		assert.Equal(t, count, nodesReported)
	}
}

func TestSingleChunkGroupMultipleAssignments(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(1)
	startTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	assignmentCount := uint32(rand.Intn(3) + 2)

	cgMap := NewMap(chunkGroupCount, assignmentCount, shufflePeriod)
	assert.Equal(t, uint32(0), cgMap.Size())

	expectedRegistrations := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedRegistrations[registration.ID()] = registration
		cgMap.Add(startTime, registration)
	}

	now := startTime
	steps := 10
	for step := 0; step < steps; step++ {
		if rand.Float64() < (1.0 / 3.0) {
			// Add less than a full shuffle period.
			now = now.Add(shufflePeriod * time.Duration(rand.Float64()))
		} else if rand.Float64() < (2.0 / 3.0) {
			// Add exactly one shuffle period.
			now = now.Add(shufflePeriod)
		} else {
			// Add several shuffle periods.
			now = now.Add(shufflePeriod * time.Duration(rand.Intn(10)+2))
		}

		// Add a few elements.
		numberToAdd := rand.Intn(10)
		count += numberToAdd
		for i := 0; i < numberToAdd; i++ {
			registration := randomRegistration()
			expectedRegistrations[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedRegistrations {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedRegistrations, key)
			numberToRemove--
		}

		// Verify the chunk group for each element.
		for id, registration := range expectedRegistrations {
			chunkGroups, ok := cgMap.GetChunkGroups(now, id)
			assert.True(t, ok)
			assert.Equal(t, 1, len(chunkGroups))
			chunkGroup := chunkGroups[0]

			offset := ComputeShuffleOffset(registration.Seed(), 0, shufflePeriod)
			epoch := ComputeShuffleEpoch(shufflePeriod, offset, now)
			expectedChunkGroup := ComputeChunkGroup(registration.Seed(), 0, epoch, chunkGroupCount)

			assert.Equal(t, expectedChunkGroup, chunkGroup)
		}

		// Query for full chunk groups.
		nodesReported := 0
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)
			nodesReported += len(chunk)

			for _, nodeID := range chunk {
				chunkGroups, ok := cgMap.GetChunkGroups(now, nodeID)
				assert.True(t, ok)
				assert.Equal(t, 1, len(chunkGroups))
				chunkGroup := chunkGroups[0]

				assert.Equal(t, chunkIndex, chunkGroup)
			}
		}

		assert.Equal(t, count, nodesReported)
	}
}
