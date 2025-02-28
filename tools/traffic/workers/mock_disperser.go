package workers

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/stretchr/testify/mock"
)

var _ clients.DisperserClient = (*MockDisperserClient)(nil)

type MockDisperserClient struct {
	mock mock.Mock
}

func (m *MockDisperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	blobVersion corev2.BlobVersion,
	quorums []core.QuorumID,
) (*dispv2.BlobStatus, corev2.BlobKey, error) {

	args := m.mock.Called(ctx, data, blobVersion, quorums)
	return args.Get(0).(*dispv2.BlobStatus), args.Get(1).(corev2.BlobKey), args.Error(2)
}

func (m *MockDisperserClient) GetBlobStatus(ctx context.Context, blobKey corev2.BlobKey) (*disperser_rpc.BlobStatusReply, error) {
	args := m.mock.Called(blobKey)
	return args.Get(0).(*disperser_rpc.BlobStatusReply), args.Error(1)
}

func (m *MockDisperserClient) GetBlobCommitment(ctx context.Context, data []byte) (*disperser_rpc.BlobCommitmentReply, error) {
	args := m.mock.Called(data)
	return args.Get(0).(*disperser_rpc.BlobCommitmentReply), args.Error(1)
}

func (m *MockDisperserClient) Close() error {
	args := m.mock.Called()
	return args.Error(0)
}
