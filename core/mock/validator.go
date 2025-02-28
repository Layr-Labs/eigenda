package mock

import (
	"errors"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
)

// MockShardValidator is a mock implementation of ShardValidator
type MockShardValidator struct {
	mock.Mock
}

var _ core.ShardValidator = (*MockShardValidator)(nil)

func NewMockShardValidator() *MockShardValidator {
	return &MockShardValidator{}
}

func (v *MockShardValidator) ValidateBatch(batchHeader *core.BatchHeader, blobs []*core.BlobMessage, operatorState *core.OperatorState, pool common.WorkerPool) error {
	args := v.Called(blobs, operatorState, pool)
	return args.Error(0)
}

func (v *MockShardValidator) ValidateBlobs(blobs []*core.BlobMessage, operatorState *core.OperatorState, pool common.WorkerPool) error {
	args := v.Called(blobs, operatorState, pool)
	return args.Error(0)
}

func (v *MockShardValidator) UpdateOperatorID(operatorID core.OperatorID) {
	v.Called(operatorID)
}
