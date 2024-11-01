package workers

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/mock"
)

var _ clients.DisperserClient = (*MockDisperserClient)(nil)

type MockDisperserClient struct {
	mock mock.Mock
}

func (m *MockDisperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	args := m.mock.Called(data, customQuorums)
	return args.Get(0).(*disperser.BlobStatus), args.Get(1).([]byte), args.Error(2)
}

func (m *MockDisperserClient) DisperseBlobAuthenticated(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	args := m.mock.Called(data, customQuorums)
	return args.Get(0).(*disperser.BlobStatus), args.Get(1).([]byte), args.Error(2)
}

func (m *MockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	args := m.mock.Called(key)
	return args.Get(0).(*disperser_rpc.BlobStatusReply), args.Error(1)
}

func (m *MockDisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	args := m.mock.Called(batchHeaderHash, blobIndex)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDisperserClient) Close() error {
	args := m.mock.Called()
	return args.Error(0)
}
