package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockRelayClient struct {
	mock.Mock
}

var _ clients.RelayClient = (*MockRelayClient)(nil)

func NewRelayClient() *MockRelayClient {
	return &MockRelayClient{}
}

func (c *MockRelayClient) GetBlob(ctx context.Context, relayKey corev2.RelayKey, blobKey corev2.BlobKey) ([]byte, error) {
	args := c.Called(blobKey)
	return args.Get(0).([]byte), args.Error(1)
}

func (c *MockRelayClient) GetChunksByRange(ctx context.Context, relayKey corev2.RelayKey, requests []*clients.ChunkRequestByRange) ([][]byte, error) {
	args := c.Called()
	return args.Get(0).([][]byte), args.Error(1)
}

func (c *MockRelayClient) GetChunksByIndex(ctx context.Context, relayKey corev2.RelayKey, requests []*clients.ChunkRequestByIndex) ([][]byte, error) {
	args := c.Called()
	return args.Get(0).([][]byte), args.Error(1)
}

func (c *MockRelayClient) Close() error {
	args := c.Called()
	return args.Error(0)
}
