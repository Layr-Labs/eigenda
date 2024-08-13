package chunkgroup

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/lightnode"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func randomRegistration() *lightnode.Registration {
	return lightnode.NewRegistration(rand.Uint64(), rand.Uint64(), tu.RandomTime())
}

func TestAddRemoveGet(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	genesisTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)

	cgMap := NewMap(chunkGroupCount, genesisTime, shufflePeriod)
	assert.Equal(t, uint(0), cgMap.Size())

	expectedMap := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedMap[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(genesisTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint(i+1), cgMap.Size())
	}

	// Removing non-existent elements should be a no-op.
	for i := 0; i < 10; i++ {
		cgMap.Remove(rand.Uint64())
		assert.Equal(t, uint(elementsToAdd), cgMap.Size())
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedMap {
		assert.Equal(t, registration, cgMap.Get(id))
	}

	// Remove all nodes that have an ID divisible by 2.
	removalCount := 0
	for id, registration := range expectedMap {
		if id%2 == 0 {
			assert.Equal(t, registration, cgMap.Get(id))
			cgMap.Remove(id)
			assert.Nil(t, cgMap.Get(id))
			removalCount++
			assert.Equal(t, uint(elementsToAdd-removalCount), cgMap.Size())
		}
	}

	// Verify that get returns the correct registrations.
	for id, registration := range expectedMap {
		if id%2 == 0 {
			assert.Nil(t, cgMap.Get(id))
		} else {
			assert.Equal(t, registration, cgMap.Get(id))
		}
	}
}

func TestChunkGroupCalculations(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	genesisTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)

	cgMap := NewMap(chunkGroupCount, genesisTime, shufflePeriod)
	assert.Equal(t, uint(0), cgMap.Size())

	expectedMap := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedMap[registration.ID()] = registration
		cgMap.Add(genesisTime, registration)
	}

	now := genesisTime
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
			expectedMap[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedMap {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedMap, key)
			numberToRemove--
		}

		// Verify the chunk group for each element.
		for id, registration := range expectedMap {
			chunkGroup, ok := cgMap.GetChunkGroup(now, id)

			assert.True(t, ok)
			offset := ComputeShuffleOffset(registration.Seed(), shufflePeriod)
			epoch := ComputeShuffleEpoch(genesisTime, shufflePeriod, offset, now)
			expectedChunkGroup := ComputeChunkGroup(registration.Seed(), epoch, chunkGroupCount)

			assert.Equal(t, expectedChunkGroup, chunkGroup)
		}

		// Query for full chunk groups.
		nodesReported := 0
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)
			nodesReported += len(chunk)

			for _, registration := range chunk {
				chunkGroup, ok := cgMap.GetChunkGroup(now, registration.ID())
				assert.True(t, ok)
				assert.Equal(t, chunkIndex, chunkGroup)
			}
		}

		assert.Equal(t, count, nodesReported)
	}
}

func TestGetRandomNode(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(rand.Intn(100) + 1)
	genesisTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)

	cgMap := NewMap(chunkGroupCount, genesisTime, shufflePeriod)
	assert.Equal(t, uint(0), cgMap.Size())

	expectedMap := make(map[uint64]*lightnode.Registration)

	// Add elements
	elementsToAdd := 1_000
	for i := 0; i < elementsToAdd; i++ {
		registration := randomRegistration()
		expectedMap[registration.ID()] = registration

		assert.Nil(t, cgMap.Get(registration.ID()))
		cgMap.Add(genesisTime, registration)
		assert.Equal(t, registration, cgMap.Get(registration.ID()))

		assert.Equal(t, uint(i+1), cgMap.Size())
	}

	now := genesisTime.Add(shufflePeriod * time.Duration(rand.Float64()*1000))

	for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {

		chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)

		if len(chunk) == 0 {
			_, ok := cgMap.GetRandomNode(now, chunkIndex, 0)
			assert.False(t, ok)
			continue
		}

		for i := 0; i < 10; i++ {

			var minimumTimeInGroup time.Duration
			if rand.Float64() < 0.1 {
				minimumTimeInGroup = 0
			} else {
				minimumTimeInGroup = shufflePeriod / time.Duration(rand.Intn(5)+1)
			}

			randomNode, ok := cgMap.GetRandomNode(now, chunkIndex, minimumTimeInGroup)

			if ok {
				assert.NotNil(t, randomNode)
				assert.Contains(t, chunk, randomNode)
			} else {
				// there shouldn't be any nodes in the chunk group for the minimum time
				for _, registration := range chunk {
					offset := ComputeShuffleOffset(registration.Seed(), shufflePeriod)
					epoch := ComputeShuffleEpoch(genesisTime, shufflePeriod, offset, genesisTime)
					epochBeginning := ComputeStartOfShuffleEpoch(genesisTime, shufflePeriod, offset, epoch)
					timeInGroup := now.Sub(epochBeginning)
					assert.True(t, timeInGroup >= minimumTimeInGroup)
				}
			}
		}
	}
}

func TestSingleChunkGroup(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint32(1)
	genesisTime := tu.RandomTime()
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)

	cgMap := NewMap(chunkGroupCount, genesisTime, shufflePeriod)
	assert.Equal(t, uint(0), cgMap.Size())

	expectedMap := make(map[uint64]*lightnode.Registration)

	// Add elements
	count := 1_000
	for i := 0; i < count; i++ {
		registration := randomRegistration()
		expectedMap[registration.ID()] = registration
		cgMap.Add(genesisTime, registration)
	}

	now := genesisTime
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
			expectedMap[registration.ID()] = registration
			cgMap.Add(now, registration)
		}

		// Remove a few elements.
		numberToRemove := rand.Intn(10)
		count -= numberToRemove
		for key := range expectedMap {
			if numberToRemove == 0 {
				break
			}
			cgMap.Remove(key)
			delete(expectedMap, key)
			numberToRemove--
		}

		// Verify the chunk group for each element.
		for id, registration := range expectedMap {
			chunkGroup, ok := cgMap.GetChunkGroup(now, id)

			assert.True(t, ok)
			offset := ComputeShuffleOffset(registration.Seed(), shufflePeriod)
			epoch := ComputeShuffleEpoch(genesisTime, shufflePeriod, offset, now)
			expectedChunkGroup := ComputeChunkGroup(registration.Seed(), epoch, chunkGroupCount)

			assert.Equal(t, expectedChunkGroup, chunkGroup)
		}

		// Query for full chunk groups.
		nodesReported := 0
		for chunkIndex := uint32(0); chunkIndex < chunkGroupCount; chunkIndex++ {
			chunk := cgMap.GetNodesInChunkGroup(now, chunkIndex)
			nodesReported += len(chunk)

			for _, registration := range chunk {
				chunkGroup, ok := cgMap.GetChunkGroup(now, registration.ID())
				assert.True(t, ok)
				assert.Equal(t, chunkIndex, chunkGroup)
			}
		}

		assert.Equal(t, count, nodesReported)
	}
}
