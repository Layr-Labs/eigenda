package table

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

// initializeRandom initializes the random number generator. Prints the seed so that the test can be rerun
// deterministically. Replace a call to this method with a call to initializeRandomWithSeed to rerun a test
// with a specific seed.
func initializeRandom() {
	rand.Seed(uint64(time.Now().UnixNano()))
	seed := rand.Uint64()
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// initializeRandomWithSeed initializes the random number generator with a specific seed.
func initializeRandomWithSeed(seed uint64) {
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// randomMetadata generates a random BlobMetadata instance.
func randomMetadata(permits int) *BlobMetadata {
	key := make([]byte, 32)
	batchHeaderHash := make([]byte, 32)
	checksum := [16]byte{}
	_, _ = rand.Read(key)
	_, _ = rand.Read(checksum[:])
	_, _ = rand.Read(batchHeaderHash)
	return NewBlobMetadata(&key, &checksum, 1024, &batchHeaderHash, 0, permits)
}

// TestBasicOperation tests basic operations of the BlobTable. Adds blobs and iterates over them.
func TestBasicOperation(t *testing.T) {
	initializeRandom()

	table := NewBlobTable()
	assert.Equal(t, uint(0), table.Size())

	size := 1024
	expectedMetadata := make([]*BlobMetadata, 0)
	for i := 0; i < size; i++ {
		metadata := randomMetadata(1)
		table.Add(metadata)
		expectedMetadata = append(expectedMetadata, metadata)
		assert.Equal(t, uint(i+1), table.Size())
	}

	for i := 0; i < size; i++ {
		assert.Equal(t, expectedMetadata[i], table.Get(uint(i)))
	}

	// Requesting an index that is out of bounds should return nil.
	assert.Nil(t, table.Get(uint(size)))
}

// TestGetRandomWithRemoval tests getting a random blob data, but where the number of permits per blob is unlimited.
func TestGetRandomNoRemovalByConfiguration(t *testing.T) {
	initializeRandom()

	table := NewBlobTable()
	assert.Equal(t, uint(0), table.Size())

	// Requesting a random element from an empty table should return nil.
	element, _ := table.GetRandom(true)
	assert.Nil(t, element)

	expectedMetadata := make([]*BlobMetadata, 0)
	size := 128
	for i := 0; i < size; i++ {
		metadata := randomMetadata(-1) // -1 == unlimited permits
		table.Add(metadata)
		expectedMetadata = append(expectedMetadata, metadata)
		assert.Equal(t, uint(i+1), table.Size())
	}

	randomIndices := make(map[uint]bool)

	// Query more times than the number of blobs to ensure that blobs are not removed.
	for i := 0; i < size*8; i++ {
		// This parameter will be ignored given that the number of permits is unlimited.
		// But not a bad thing to exercise the code path.
		decrement := rand.Intn(2) == 1

		metadata, removed := table.GetRandom(decrement)
		assert.False(t, removed)
		assert.NotNil(t, metadata)
		assert.Equal(t, expectedMetadata[metadata.index], metadata)

		randomIndices[metadata.index] = true
	}

	// Sanity check: ensure that at least 10 different blobs were returned. This check is attempting to verify
	// that we are actually getting random blobs. The probability of this check failing is extremely low if
	// the random number generator is working correctly.
	assert.GreaterOrEqual(t, len(randomIndices), 10)
}

// TestGetRandomWithRemoval tests getting a random blob data, where the number of permits per blob is limited.
func TestGetRandomWithRemoval(t *testing.T) {
	initializeRandom()

	table := NewBlobTable()
	assert.Equal(t, uint(0), table.Size())

	// Requesting a random element from an empty table should return nil.
	element, _ := table.GetRandom(true)
	assert.Nil(t, element)

	permitCount := 2

	size := 1024
	expectedMetadata := make(map[*[]byte]uint, 0)
	for i := 0; i < size; i++ {
		metadata := randomMetadata(permitCount)
		table.Add(metadata)
		expectedMetadata[metadata.Key()] = 0
		assert.Equal(t, uint(i+1), table.Size())
	}

	// Requesting random elements without decrementing should not remove any elements.
	for i := 0; i < size; i++ {
		metadata, removed := table.GetRandom(false)
		assert.NotNil(t, metadata)
		_, exists := expectedMetadata[metadata.Key()]
		assert.True(t, exists)
		assert.False(t, removed)
	}
	assert.Equal(t, uint(size), table.Size())

	// Requesting elements a number of times equal to the size times the number of permits should completely
	// drain the table and return all elements a number of times equal to the number of permits.
	for i := 0; i < size*permitCount; i++ {
		metadata, removed := table.GetRandom(true)
		assert.NotNil(t, metadata)

		permitsUsed := expectedMetadata[metadata.Key()] + 1
		expectedMetadata[metadata.Key()] = permitsUsed
		assert.LessOrEqual(t, permitsUsed, uint(permitCount))

		if int(permitsUsed) == permitCount {
			assert.True(t, removed)
		} else {
			assert.False(t, removed)
		}
	}

	assert.Equal(t, uint(0), table.Size())
}

// TestAddOrReplace tests adding blobs to a table with a maximum capacity. The table should replace blobs when full.
func TestAddOrReplace(t *testing.T) {
	initializeRandom()

	table := NewBlobTable()
	assert.Equal(t, uint(0), table.Size())

	// Adding data to a table with capacity 0 should be a no-op.
	table.AddOrReplace(randomMetadata(1), 0)
	assert.Equal(t, uint(0), table.Size())

	randomIndices := make(map[uint]bool)

	size := 1024
	for i := 0; i < size*2; i++ {
		metadata := randomMetadata(-1) // -1 == unlimited permits

		initialSize := table.Size()
		table.AddOrReplace(metadata, uint(size))
		resultingSize := table.Size()

		assert.LessOrEqual(t, resultingSize, uint(size))
		if initialSize < uint(size) {
			assert.Equal(t, initialSize+1, resultingSize)
		} else {
			randomIndices[metadata.index] = true
		}

		// Verify that the metadata is in the table.
		assert.Less(t, metadata.index, table.Size())
		assert.Equal(t, metadata, table.Get(metadata.index))
	}

	// Sanity check: ensure that replacements happened at least 10 different indices. This check is attempting to
	// verify that we are actually replacing blobs. The probability of this check failing is extremely low if
	// the random number generator is working correctly.
	assert.GreaterOrEqual(t, len(randomIndices), 10)
}
