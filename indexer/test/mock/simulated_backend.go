package mock

import (
	"context"
	"errors"
	"math/big"
	"time"

	cm "github.com/Layr-Labs/eigenda/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type (
	SimulatedBackend interface {
		AdjustTime(adjustment time.Duration) error
		BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
		BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
		BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
		Blockchain() *core.BlockChain
		CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
		Close() error
		CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
		Commit() common.Hash
		EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
		FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
		Fork(ctx context.Context, parent common.Hash) error
		HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
		HeaderByNumber(ctx context.Context, block *big.Int) (*types.Header, error)
		NonceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (uint64, error)
		PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
		PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error)
		PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
		Rollback()
		SendTransaction(ctx context.Context, tx *types.Transaction) error
		StorageAt(ctx context.Context, contract common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
		SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
		SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
		SuggestGasPrice(ctx context.Context) (*big.Int, error)
		SuggestGasTipCap(ctx context.Context) (*big.Int, error)
		TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
		TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
		TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
		TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

		cm.RPCEthClient
	}

	simulatedBackend struct {
		*backends.SimulatedBackend
	}
)

func (sb *simulatedBackend) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "eth_getBlockByNumber":
		number := args[0].(string)
		h := result.(*types.Header)
		return sb.getBlockByNumber(ctx, h, number)
	default:
		return errors.New("method not found")
	}
}

func (sb *simulatedBackend) Call(result interface{}, method string, args ...interface{}) error {
	return sb.CallContext(context.Background(), result, method, args...)
}

func (sb *simulatedBackend) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	for _, elem := range b {
		if err := sb.CallContext(ctx, elem.Result, elem.Method, elem.Args...); err != nil {
			return err
		}
	}
	return nil
}

func (sb *simulatedBackend) BatchCall(b []rpc.BatchElem) error {
	return sb.BatchCallContext(context.Background(), b)
}

func (sb *simulatedBackend) getBlockByNumber(ctx context.Context, result *types.Header, blockNum string) error {
	var blockNumBigInt *big.Int

	if blockNum == "latest" {
		blockNumBigInt = nil
	} else {
		bn, err := hexutil.DecodeBig(blockNum)
		if err != nil {
			return err
		}
		blockNumBigInt = bn
	}

	header, err := sb.HeaderByNumber(ctx, blockNumBigInt)
	if err != nil || header == nil {
		return err
	}

	*result = *header
	return nil
}
