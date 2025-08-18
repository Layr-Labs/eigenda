package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/internal"
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

var _ internal.ValidatorGRPCManager = (*MockValidatorGRPCManager)(nil)

// MockValidatorGRPCManager is a mock implementation of the ValidatorGRPCManager interface.
type MockValidatorGRPCManager struct {
	// A lambda function to be called when DownloadChunks is called.
	DownloadChunksFunction func(ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error)
}

func (m *MockValidatorGRPCManager) DownloadChunks(
	ctx context.Context,
	key v2.BlobKey,
	operatorID core.OperatorID,
) (*grpcnode.GetChunksReply, error) {
	if m.DownloadChunksFunction == nil {
		return nil, nil
	}
	return m.DownloadChunksFunction(ctx, key, operatorID)
}

// NewMockValidatorGRPCManager creates a new ValidatorGRPCManager instance with the provided download function.
func NewMockValidatorGRPCManager(
	downloadChunksFunction func(ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error),
) internal.ValidatorGRPCManager {
	return &MockValidatorGRPCManager{
		DownloadChunksFunction: downloadChunksFunction,
	}
}
