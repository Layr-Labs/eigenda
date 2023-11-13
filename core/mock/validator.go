package mock

import (
	"errors"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrInvalidHeader       = errors.New("invalid header")
)

// MockChunkValidator is a mock implementation of ChunkValidator
type MockChunkValidator struct {
	mock.Mock
}

var _ core.ChunkValidator = (*MockChunkValidator)(nil)

func NewMockChunkValidator() *MockChunkValidator {
	return &MockChunkValidator{}
}

func (v *MockChunkValidator) ValidateBlob(blob *core.BlobMessage, operatorState *core.OperatorState) error {
	args := v.Called(blob, operatorState)
	return args.Error(0)
}

func (v *MockChunkValidator) UpdateOperatorID(operatorID core.OperatorID) {
	v.Called(operatorID)
}
