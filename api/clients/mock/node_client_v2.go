package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockNodeClientV2 struct {
	mock.Mock
}

var _ clients.NodeClientV2 = (*MockNodeClientV2)(nil)

func NewNodeClientV2() *MockNodeClientV2 {
	return &MockNodeClientV2{}
}

func (c *MockNodeClientV2) StoreChunks(ctx context.Context, batch *corev2.Batch) (*core.Signature, error) {
	args := c.Called()
	var signature *core.Signature
	if args.Get(0) != nil {
		signature = (args.Get(0)).(*core.Signature)
	}
	return signature, args.Error(1)
}

func (c *MockNodeClientV2) Close() error {
	args := c.Called()
	return args.Error(0)
}
