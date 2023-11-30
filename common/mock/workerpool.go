package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/mock"
)

type MockWorkerpool struct {
	mock.Mock
}

var _ common.WorkerPool = (*MockWorkerpool)(nil)

func (mock *MockWorkerpool) Size() int {
	args := mock.Called()
	result := args.Get(0)
	return result.(int)
}

func (mock *MockWorkerpool) Stop() {
	mock.Called()
}

func (mock *MockWorkerpool) StopWait() {
	mock.Called()
}

func (mock *MockWorkerpool) Stopped() bool {
	args := mock.Called()
	result := args.Get(0)
	return result.(bool)
}

func (mock *MockWorkerpool) Submit(task func()) {
	mock.Called(task)
}

func (mock *MockWorkerpool) SubmitWait(task func()) {
	mock.Called(task)
}

func (mock *MockWorkerpool) WaitingQueueSize() int {
	args := mock.Called()
	result := args.Get(0)
	return result.(int)
}

func (mock *MockWorkerpool) Pause(ctx context.Context) {
	mock.Called(ctx)
}
