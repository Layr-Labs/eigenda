package test

import (
	"context"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// mockRetrievalClient is a mock implementation of the clients.RetrievalClient interface.
type mockRetrievalClient struct {
	t *testing.T

	lock *sync.Mutex

	// Since it isn't being used during this test, blob index field is used
	// as a convenient unique identifier for the blob.

	// A map from blob index to the blob data.
	blobData map[uint]*[]byte

	// A map from blob index to the blob metadata.
	blobMetadata map[uint]*table.BlobMetadata

	// A map from blob index to the blob chunks corresponding to that blob.
	blobChunks map[uint]*clients.BlobChunks

	RetrieveBlobChunksCount uint
	CombineChunksCount      uint
}

func newMockRetrievalClient(t *testing.T, lock *sync.Mutex) *mockRetrievalClient {
	return &mockRetrievalClient{
		t:            t,
		lock:         lock,
		blobData:     make(map[uint]*[]byte),
		blobMetadata: make(map[uint]*table.BlobMetadata),
		blobChunks:   make(map[uint]*clients.BlobChunks),
	}
}

// AddBlob adds a blob to the mock retrieval client. Once added, the retrieval client will act as if
// it is able to retrieve the blob.
func (m *mockRetrievalClient) AddBlob(metadata *table.BlobMetadata, data []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.blobData[metadata.BlobIndex] = &data
	m.blobMetadata[metadata.BlobIndex] = metadata

	// The blob index is used in this test as a convenient unique identifier for the blob.

	m.blobChunks[metadata.BlobIndex] = &clients.BlobChunks{
		// Since it isn't otherwise used in this field, we can use it to store the unique identifier for the blob.
		BlobHeaderLength: metadata.BlobIndex,
	}
}

func (m *mockRetrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) ([]byte, error) {
	panic("this method should not be called during this test")
}

func (m *mockRetrievalClient) RetrieveBlobChunks(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) (*clients.BlobChunks, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.RetrieveBlobChunksCount++

	chunks, ok := m.blobChunks[uint(blobIndex)]
	assert.True(m.t, ok, "blob not found")

	metadata := m.blobMetadata[uint(blobIndex)]
	assert.Equal(m.t, metadata.BlobIndex, uint(blobIndex))
	assert.Equal(m.t, metadata.BatchHeaderHash[:32], batchHeaderHash[:32])

	return chunks, nil

}

func (m *mockRetrievalClient) CombineChunks(chunks *clients.BlobChunks) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.CombineChunksCount++

	blobIndex := chunks.BlobHeaderLength
	data, ok := m.blobData[blobIndex]
	assert.True(m.t, ok, "blob not found")

	return *data, nil
}
