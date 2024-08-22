package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/mock"
)

type MockEncoderClient struct {
	mock.Mock
}

var _ disperser.EncoderClient = (*MockEncoderClient)(nil)

func NewMockEncoderClient() *MockEncoderClient {
	return &MockEncoderClient{}
}

func (m *MockEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, *core.ChunksData, error) {
	args := m.Called(ctx, data, encodingParams)
	var commitments *encoding.BlobCommitments
	if args.Get(0) != nil {
		commitments = args.Get(0).(*encoding.BlobCommitments)
	}
	var chunks *core.ChunksData
	if args.Get(1) != nil {
		chunks = args.Get(1).(*core.ChunksData)
	}
	return commitments, chunks, args.Error(2)
}
