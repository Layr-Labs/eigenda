package mock

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
)

type MockRPCEthClient struct {
	mock.Mock
}

func (mock *MockRPCEthClient) BatchCall(b []rpc.BatchElem) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRPCEthClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	args := mock.Called(ctx, b)
	return args.Error(0)
}

func (mock *MockRPCEthClient) Call(result interface{}, method string, args ...interface{}) error {
	mokcArgs := mock.Called()
	return mokcArgs.Error(0)
}

func (mock *MockRPCEthClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	args = append([]interface{}{ctx, result, method}, args...)
	mokcArgs := mock.Called(args...)
	return mokcArgs.Error(0)
}
