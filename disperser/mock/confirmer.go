package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type MockBatchConfirmer struct {
	mock.Mock
}

var _ disperser.BatchConfirmer = (*MockBatchConfirmer)(nil)

func NewBatchConfirmer() *MockBatchConfirmer {
	return &MockBatchConfirmer{}
}

func (b *MockBatchConfirmer) ConfirmBatch(ctx context.Context, header *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, sig *core.SignatureAggregation) (*types.Receipt, error) {
	args := b.Called()
	var receipt *types.Receipt
	if args.Get(0) != nil {
		receipt = args.Get(0).(*types.Receipt)
	}
	return receipt, args.Error(1)
}
