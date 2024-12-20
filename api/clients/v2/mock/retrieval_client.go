package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockRetrievalClient struct {
	mock.Mock
}

var _ clients.RetrievalClient = (*MockRetrievalClient)(nil)

func NewRetrievalClient() *MockRetrievalClient {
	return &MockRetrievalClient{}
}

func (c *MockRetrievalClient) GetBlob(ctx context.Context, blobHeader *corev2.BlobHeader, referenceBlockNumber uint64, quorumID core.QuorumID) ([]byte, error) {
	args := c.Called()

	result := args.Get(0)
	return result.([]byte), args.Error(1)
}
