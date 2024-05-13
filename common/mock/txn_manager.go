package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/mock"
)

type MockTxnManager struct {
	mock.Mock

	Requests []*common.TxnRequest
}

var _ common.TxnManager = (*MockTxnManager)(nil)

func NewTxnManager() *MockTxnManager {
	return &MockTxnManager{}
}

func (b *MockTxnManager) Start(ctx context.Context) {}

func (b *MockTxnManager) ProcessTransaction(ctx context.Context, req *common.TxnRequest) error {
	args := b.Called()
	b.Requests = append(b.Requests, req)
	return args.Error(0)
}

func (b *MockTxnManager) ReceiptChan() chan *common.ReceiptOrErr {
	args := b.Called()
	return args.Get(0).(chan *common.ReceiptOrErr)
}
