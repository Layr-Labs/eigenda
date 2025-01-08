package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockNodeClient struct {
	mock.Mock
}

var _ clients.NodeClient = (*MockNodeClient)(nil)

func NewNodeClient() *MockNodeClient {
	return &MockNodeClient{}
}

func (c *MockNodeClient) StoreChunks(ctx context.Context, batch *corev2.Batch) (*core.Signature, error) {
	args := c.Called()
	var signature *core.Signature
	if args.Get(0) != nil {
		signature = (args.Get(0)).(*core.Signature)
	}
	return signature, args.Error(1)
}

func (c *MockNodeClient) Close() error {
	args := c.Called()
	return args.Error(0)
}
