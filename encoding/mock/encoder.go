package encoding

import (
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/mock"
)

type MockEncoder struct {
	mock.Mock

	Delay time.Duration
}

var _ encoding.Prover = &MockEncoder{}

var _ encoding.Verifier = &MockEncoder{}

func (e *MockEncoder) EncodeAndProve(data []byte, params encoding.EncodingParams) (encoding.BlobCommitments, []*encoding.Frame, error) {
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

func (e *MockEncoder) VerifyFrames(chunks []*encoding.Frame, indices []encoding.ChunkNumber, commitments encoding.BlobCommitments, params encoding.EncodingParams) error {
	args := e.Called(chunks, indices, commitments, params)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) UniversalVerifySubBatch(params encoding.EncodingParams, samples []encoding.Sample, numBlobs int) error {
	args := e.Called(params, samples, numBlobs)
	time.Sleep(e.Delay)
	return args.Error(0)
}
func (e *MockEncoder) VerifyCommitEquivalenceBatch(commitments []encoding.BlobCommitments) error {
	args := e.Called(commitments)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) VerifyBlobLength(commitments encoding.BlobCommitments) error {

	args := e.Called(commitments)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) Decode(chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64) ([]byte, error) {
	args := e.Called(chunks, indices, params, maxInputSize)
	time.Sleep(e.Delay)
	return args.Get(0).([]byte), args.Error(1)
}
