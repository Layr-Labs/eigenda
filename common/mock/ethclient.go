package mock

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"

	dacommon "github.com/Layr-Labs/eigenda/common"
)

type MockEthClient struct {
	mock.Mock
}

var _ dacommon.EthClient = (*MockEthClient)(nil)

func (mock *MockEthClient) GetAccountAddress() common.Address {
	args := mock.Called()
	result := args.Get(0)
	return result.(common.Address)
}

func (mock *MockEthClient) GetNoSendTransactOpts() (*bind.TransactOpts, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*bind.TransactOpts), args.Error(1)
}

func (mock *MockEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*types.Block), args.Error(1)
}

func (mock *MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*types.Block), args.Error(1)
}

func (mock *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), nil
}

func (mock *MockEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockEthClient) FeeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock *big.Int,
	rewardPercentiles []float64,
) (*ethereum.FeeHistory, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*ethereum.FeeHistory), args.Error(1)
}

func (mock *MockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	args := mock.Called(q)
	result := args.Get(0)
	return result.([]types.Log), args.Error(1)
}

func (mock *MockEthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*types.Header), args.Error(1)
}

func (mock *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*types.Header), args.Error(1)
}

func (mock *MockEthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockEthClient) PeerCount(ctx context.Context) (uint64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockEthClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (mock *MockEthClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint), args.Error(1)
}

func (mock *MockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockEthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *MockEthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(ethereum.Subscription), args.Error(1)
}

func (mock *MockEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(ethereum.Subscription), args.Error(1)
}

func (mock *MockEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (mock *MockEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*ethereum.SyncProgress), args.Error(1)
}

func (mock *MockEthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	args := mock.Called(hash)
	result1 := args.Get(0)
	result2 := args.Get(1)
	return result1.(*types.Transaction), result2.(bool), args.Error(2)
}

func (mock *MockEthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(uint), args.Error(1)
}

func (mock *MockEthClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*types.Transaction), args.Error(1)
}

func (mock *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := mock.Called()
	var result *types.Receipt
	if args.Get(0) != nil {
		result = args.Get(0).(*types.Receipt)
	}

	return result, args.Error(1)
}

func (mock *MockEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(common.Address), args.Error(1)
}

func (mock *MockEthClient) GetLatestGasCaps(ctx context.Context) (gasTipCap, gasFeeCap *big.Int, err error) {
	args := mock.Called()
	result1 := args.Get(0)
	result2 := args.Get(1)
	return result1.(*big.Int), result2.(*big.Int), args.Error(2)
}

func (mock *MockEthClient) EstimateGasPriceAndLimitAndSendTx(ctx context.Context, tx *types.Transaction, tag string, value *big.Int) (*types.Receipt, error) {
	args := mock.Called()
	var result *types.Receipt
	if args.Get(0) != nil {
		result = args.Get(0).(*types.Receipt)
	}

	return result, args.Error(1)
}

func (mock *MockEthClient) UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	args := mock.Called()
	var newTx *types.Transaction
	if args.Get(0) != nil {
		newTx = args.Get(0).(*types.Transaction)
	}
	return newTx, args.Error(1)
}

func (mock *MockEthClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	args := mock.Called()
	var result *types.Receipt
	if args.Get(0) != nil {
		result = args.Get(0).(*types.Receipt)
	}

	return result, args.Error(1)
}

func (mock *MockEthClient) EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error) {
	args := mock.Called()
	var result *types.Receipt
	if args.Get(0) != nil {
		result = args.Get(0).(*types.Receipt)
	}

	return result, args.Error(1)
}
