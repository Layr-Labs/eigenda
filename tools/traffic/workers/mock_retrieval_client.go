package workers

import (
	"context"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

// mockRetrievalClient is a mock implementation of the clients.RetrievalClient interface.
type mockRetrievalClient struct {
	mock mock.Mock
}

func (m *mockRetrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) ([]byte, error) {
	args := m.mock.Called(batchHeaderHash, blobIndex, referenceBlockNumber, batchRoot, quorumID)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockRetrievalClient) RetrieveBlobChunks(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) (*clients.BlobChunks, error) {

	args := m.mock.Called(batchHeaderHash, blobIndex, referenceBlockNumber, batchRoot, quorumID)
	return args.Get(0).(*clients.BlobChunks), args.Error(1)
}

func (m *mockRetrievalClient) CombineChunks(chunks *clients.BlobChunks) ([]byte, error) {
	args := m.mock.Called(chunks)
	return args.Get(0).([]byte), args.Error(1)
}
