package common

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

type RPCEthClient interface {
	BatchCall(b []rpc.BatchElem) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	Call(result interface{}, method string, args ...interface{}) error
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}
