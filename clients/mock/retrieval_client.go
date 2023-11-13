package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/clients"
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
