package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockNodeDispersalClient struct {
	mock.Mock
}

var _ node.DispersalClient = (*MockNodeDispersalClient)(nil)

func NewMockDispersalClient() *MockNodeDispersalClient {
	return &MockNodeDispersalClient{}
}

func (m *MockNodeDispersalClient) StoreChunks(ctx context.Context, in *node.StoreChunksRequest, opts ...grpc.CallOption) (*node.StoreChunksReply, error) {
	args := m.Called()
	return args.Get(0).(*node.StoreChunksReply), args.Error(1)
}

func (m *MockNodeDispersalClient) StoreBlobs(ctx context.Context, in *node.StoreBlobsRequest, opts ...grpc.CallOption) (*node.StoreBlobsReply, error) {
	args := m.Called()
	return args.Get(0).(*node.StoreBlobsReply), args.Error(1)
}

func (m *MockNodeDispersalClient) AttestBatch(ctx context.Context, in *node.AttestBatchRequest, opts ...grpc.CallOption) (*node.AttestBatchReply, error) {
	args := m.Called()
	return args.Get(0).(*node.AttestBatchReply), args.Error(1)
}

func (m *MockNodeDispersalClient) NodeInfo(ctx context.Context, in *node.NodeInfoRequest, opts ...grpc.CallOption) (*node.NodeInfoReply, error) {
	args := m.Called()
	return args.Get(0).(*node.NodeInfoReply), args.Error(1)
}
