package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/mock"
)

type MockEncoderClient struct {
	mock.Mock
}

var _ disperser.EncoderClient = (*MockEncoderClient)(nil)

func NewMockEncoderClient() *MockEncoderClient {
	return &MockEncoderClient{}
}

func (m *MockEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams core.EncodingParams) (*core.BlobCommitments, []*core.Chunk, error) {
	args := m.Called(ctx, data, encodingParams)
	var commitments *core.BlobCommitments
	if args.Get(0) != nil {
		commitments = args.Get(0).(*core.BlobCommitments)
	}
	var chunks []*core.Chunk
	if args.Get(1) != nil {
		chunks = args.Get(1).([]*core.Chunk)
	}
	return commitments, chunks, args.Error(2)
}
