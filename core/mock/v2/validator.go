package v2

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

// MockShardValidator is a mock implementation of ShardValidator
type MockShardValidator struct {
	mock.Mock
}

var _ corev2.ShardValidator = (*MockShardValidator)(nil)

func NewMockShardValidator() *MockShardValidator {
	return &MockShardValidator{}
}

func (v *MockShardValidator) ValidateBatchHeader(ctx context.Context, header *corev2.BatchHeader, blobCerts []*corev2.BlobCertificate) error {
	args := v.Called()
	return args.Error(0)
}

func (v *MockShardValidator) ValidateBlobs(ctx context.Context, blobs []*corev2.BlobShard, pool common.WorkerPool, state *core.OperatorState) error {
	args := v.Called()
	return args.Error(0)
}
