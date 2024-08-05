package table

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
)

// randomMetadata generates a random BlobMetadata instance.
func randomMetadata(t *testing.T, permits int) *BlobMetadata {
	key := make([]byte, 32)
	checksum := [16]byte{}
	_, _ = rand.Read(key)
	_, _ = rand.Read(checksum[:])
	metadata, err := NewBlobMetadata(key, checksum, 1024, 0, permits)
	assert.Nil(t, err)

	return metadata
}

// TestBasicOperation tests basic operations of the BlobTable. Adds blobs and iterates over them.
func TestBasicOperation(t *testing.T) {
	tu.InitializeRandom()

	store := NewBlobStore()
	assert.Equal(t, uint(0), store.Size())

	size := 1024
	expectedMetadata := make([]*BlobMetadata, 0)
	for i := 0; i < size; i++ {
		metadata := randomMetadata(t, 1)
		store.Add(metadata)
		expectedMetadata = append(expectedMetadata, metadata)
		assert.Equal(t, uint(i+1), store.Size())
	}

	for i := 0; i < size; i++ {
		assert.Equal(t, expectedMetadata[i], store.blobs[uint64(i)])
	}
}

// TestGetRandomWithRemoval tests getting a random blob data, but where the number of permits per blob is unlimited.
func TestGetRandomNoRemoval(t *testing.T) {
	tu.InitializeRandom()

	table := NewBlobStore()
	assert.Equal(t, uint(0), table.Size())

	// Requesting a random element from an empty table should return nil.
	element := table.GetNext()
	assert.Nil(t, element)

	expectedMetadata := make(map[string]*BlobMetadata)
	size := 128
	for i := 0; i < size; i++ {
		metadata := randomMetadata(t, -1) // -1 == unlimited permits
		table.Add(metadata)
		expectedMetadata[string(metadata.Key)] = metadata
		assert.Equal(t, uint(i+1), table.Size())
	}

	randomIndices := make(map[string]bool)

	// Query more times than the number of blobs to ensure that blobs are not removed.
	for i := 0; i < size*8; i++ {
		metadata := table.GetNext()
		assert.NotNil(t, metadata)
		assert.Equal(t, expectedMetadata[string(metadata.Key)], metadata)
		randomIndices[string(metadata.Key)] = true
	}

	// Sanity check: ensure that at least 10 different blobs were returned. This check is attempting to verify
	// that we are actually getting random blobs. The probability of this check failing is extremely low if
	// the random number generator is working correctly.
	assert.GreaterOrEqual(t, len(randomIndices), 10)
}

// TestGetRandomWithRemoval tests getting a random blob data, where the number of permits per blob is limited.
func TestGetRandomWithRemoval(t *testing.T) {
	tu.InitializeRandom()

	table := NewBlobStore()
	assert.Equal(t, uint(0), table.Size())

	// Requesting a random element from an empty table should return nil.
	element := table.GetNext()
	assert.Nil(t, element)

	permitCount := 2

	size := 1024
	expectedMetadata := make(map[string]uint)
	for i := 0; i < size; i++ {
		metadata := randomMetadata(t, permitCount)
		table.Add(metadata)
		expectedMetadata[string(metadata.Key)] = 0
		assert.Equal(t, uint(i+1), table.Size())
	}

	// Requesting elements a number of times equal to the size times the number of permits should completely
	// drain the table and return all elements a number of times equal to the number of permits.
	for i := 0; i < size*permitCount; i++ {
		metadata := table.GetNext()
		assert.NotNil(t, metadata)

		k := string(metadata.Key)
		permitsUsed := expectedMetadata[k] + 1
		expectedMetadata[k] = permitsUsed
		assert.LessOrEqual(t, permitsUsed, uint(permitCount))
	}

	assert.Equal(t, uint(0), table.Size())
}
