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

// MockDataValidator is a mock implementation of DataValidator
type MockDataValidator struct {
	mock.Mock
}

var _ core.DataValidator = (*MockDataValidator)(nil)

func NewMockDataValidator() *MockDataValidator {
	return &MockDataValidator{}
}

func (v *MockDataValidator) ValidateBatch(batchHeader *core.BatchHeader, blobs []*core.BlobMessage, operatorState *core.OperatorState, pool common.WorkerPool) error {
	args := v.Called(blobs, operatorState, pool)
	return args.Error(0)
}

func (v *MockDataValidator) UpdateOperatorID(operatorID core.OperatorID) {
	v.Called(operatorID)
}
