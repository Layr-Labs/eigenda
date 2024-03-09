package common

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	GetAccountAddress() common.Address
	GetNoSendTransactOpts() (*bind.TransactOpts, error)
	ChainID(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	FeeHistory(
		ctx context.Context,
		blockCount uint64,
		lastBlock *big.Int,
		rewardPercentiles []float64,
	) (*ethereum.FeeHistory, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	PeerCount(ctx context.Context) (uint64, error)
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	GetLatestGasCaps(ctx context.Context) (gasTipCap, gasFeeCap *big.Int, err error)
	EstimateGasPriceAndLimitAndSendTx(ctx context.Context, tx *types.Transaction, tag string, value *big.Int) (*types.Receipt, error)
	UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error)
	EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error)
	EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error)
}
