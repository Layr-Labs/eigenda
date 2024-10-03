package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type MockRetrievalClient struct {
	mock.Mock
}

var _ clients.RetrievalClient = (*MockRetrievalClient)(nil)

func NewRetrievalClient() *MockRetrievalClient {
	return &MockRetrievalClient{}
}

func (c *MockRetrievalClient) StartIndexingChainState(ctx context.Context) error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockRetrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) ([]byte, error) {
	args := c.Called()

	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (c *MockRetrievalClient) RetrieveBlobChunks(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) (*clients.BlobChunks, error) {

	args := c.Called(batchHeaderHash, blobIndex, referenceBlockNumber, batchRoot, quorumID)
	return args.Get(0).(*clients.BlobChunks), args.Error(1)
}

func (c *MockRetrievalClient) CombineChunks(chunks *clients.BlobChunks) ([]byte, error) {
	args := c.Called(chunks)

	result := args.Get(0)
	return result.([]byte), args.Error(1)
}
