package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockRelayClient struct {
	mock.Mock
}

var _ relay.RelayClient = (*MockRelayClient)(nil)

func NewRelayClient() *MockRelayClient {
	return &MockRelayClient{}
}

//nolint:wrapcheck // mock code intentionally returns unwrapped errors
func (c *MockRelayClient) GetBlob(ctx context.Context, cert coretypes.EigenDACert) (*coretypes.Blob, error) {
	args := c.Called(ctx, cert)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*coretypes.Blob), args.Error(1)
}

func (c *MockRelayClient) GetChunksByRange(ctx context.Context, relayKey corev2.RelayKey, requests []*relay.ChunkRequestByRange) ([][]byte, error) {
	args := c.Called(ctx, relayKey, requests)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([][]byte), args.Error(1)
}

func (c *MockRelayClient) GetChunksByIndex(ctx context.Context, relayKey corev2.RelayKey, requests []*relay.ChunkRequestByIndex) ([][]byte, error) {
	args := c.Called(ctx, relayKey, requests)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([][]byte), args.Error(1)
}

func (c *MockRelayClient) Close() error {
	args := c.Called()
	return args.Error(0)
}
