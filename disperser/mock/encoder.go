// nolint: wrapcheck
package mock

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
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

type MockEncoder struct {
	mock.Mock

	Delay time.Duration
}

var _ encoder.Prover = &MockEncoder{}

func (e *MockEncoder) Decode(
	chunks []*encoding.Frame, indices []encoding.ChunkNumber,
	params encoding.EncodingParams, maxInputSize uint64,
) ([]byte, error) {
	args := e.Called(chunks, indices, params, maxInputSize)
	time.Sleep(e.Delay)
	return args.Get(0).([]byte), args.Error(1)
}

func (e *MockEncoder) EncodeAndProve(
	data []byte, params encoding.EncodingParams,
) (encoding.BlobCommitments, []*encoding.Frame, error) {
	args := e.Called(data, params)
	time.Sleep(e.Delay)
	return args.Get(0).(encoding.BlobCommitments), args.Get(1).([]*encoding.Frame), args.Error(2)
}

func (e *MockEncoder) GetCommitmentsForPaddedLength(data []byte) (encoding.BlobCommitments, error) {
	args := e.Called(data)
	time.Sleep(e.Delay)
	return args.Get(0).(encoding.BlobCommitments), args.Error(1)
}

func (e *MockEncoder) GetFrames(data []byte, params encoding.EncodingParams) ([]*encoding.Frame, error) {
	args := e.Called(data, params)
	time.Sleep(e.Delay)
	return args.Get(0).([]*encoding.Frame), args.Error(1)
}

func (e *MockEncoder) GetMultiFrameProofs(data []byte, params encoding.EncodingParams) ([]encoding.Proof, error) {
	args := e.Called(data, params)
	time.Sleep(e.Delay)
	return args.Get(0).([]encoding.Proof), args.Error(1)
}

func (e *MockEncoder) GetSRSOrder() uint64 {
	args := e.Called()
	return args.Get(0).(uint64)
}
