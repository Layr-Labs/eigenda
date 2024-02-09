package encoding

import (
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type MockEncoder struct {
	mock.Mock

	Delay time.Duration
}

var _ core.Encoder = &MockEncoder{}

func (e *MockEncoder) Encode(data []byte, params core.EncodingParams) (core.BlobCommitments, []*core.Chunk, error) {
	args := e.Called(data, params)
	time.Sleep(e.Delay)
	return args.Get(0).(core.BlobCommitments), args.Get(1).([]*core.Chunk), args.Error(2)
}

func (e *MockEncoder) VerifyChunks(chunks []*core.Chunk, indices []core.ChunkNumber, commitments core.BlobCommitments, params core.EncodingParams) error {
	args := e.Called(chunks, indices, commitments, params)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) UniversalVerifySubBatch(params core.EncodingParams, samples []core.Sample, numBlobs int) error {
	args := e.Called(params, samples, numBlobs)
	time.Sleep(e.Delay)
	return args.Error(0)
}
func (e *MockEncoder) VerifyCommitEquivalenceBatch(commitments []core.BlobCommitments) error {
	args := e.Called(commitments)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) VerifyBlobLength(commitments core.BlobCommitments) error {

	args := e.Called(commitments)
	time.Sleep(e.Delay)
	return args.Error(0)
}

func (e *MockEncoder) Decode(chunks []*core.Chunk, indices []core.ChunkNumber, params core.EncodingParams, maxInputSize uint64) ([]byte, error) {
	args := e.Called(chunks, indices, params, maxInputSize)
	time.Sleep(e.Delay)
	return args.Get(0).([]byte), args.Error(1)
}
