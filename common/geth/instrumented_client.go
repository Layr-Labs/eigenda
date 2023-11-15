package geth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// InstrumentedEthClient is a wrapper around our EthClient that instruments all underlying json-rpc calls.
// It counts each eth_ call made to it, as well as its duration, and exposes them as prometheus metrics
//
// TODO: This client is a temporary hack. Ideally this should be done at the geth rpcclient level,
// not the ethclient level, which would be much cleaner... but geth implemented the gethclient
// using an rpcClient struct instead of interface... see https://github.com/ethereum/go-ethereum/issues/28267
// to track progress on this
type InstrumentedEthClient struct {
	*EthClient
	rpcCallsCollector *rpccalls.Collector
	clientAndVersion  string
}

var _ common.EthClient = (*InstrumentedEthClient)(nil)

func NewInstrumentedEthClient(config EthClientConfig, rpcCallsCollector *rpccalls.Collector, logger common.Logger) (*InstrumentedEthClient, error) {
	ethClient, err := NewClient(config, logger)
	if err != nil {
		return nil, err
	}
	c := &InstrumentedEthClient{
		EthClient:         ethClient,
		rpcCallsCollector: rpcCallsCollector,
		clientAndVersion:  getClientAndVersion(ethClient),
	}

	return c, err
}

func (iec *InstrumentedEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	chainID := func() (*big.Int, error) { return iec.Client.ChainID(ctx) }
	id, err := instrumentFunction[*big.Int](chainID, "eth_chainId", iec)
	return id, err
}

func (iec *InstrumentedEthClient) BalanceAt(
	ctx context.Context,
	account gethcommon.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	balanceAt := func() (*big.Int, error) { return iec.Client.BalanceAt(ctx, account, blockNumber) }
	balance, err := instrumentFunction[*big.Int](balanceAt, "eth_getBalance", iec)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (iec *InstrumentedEthClient) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	blockByHash := func() (*types.Block, error) { return iec.Client.BlockByHash(ctx, hash) }
	block, err := instrumentFunction[*types.Block](blockByHash, "eth_getBlockByHash", iec)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (iec *InstrumentedEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	blockByNumber := func() (*types.Block, error) { return iec.Client.BlockByNumber(ctx, number) }
	block, err := instrumentFunction[*types.Block](
		blockByNumber,
		"eth_getBlockByNumber",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (iec *InstrumentedEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber := func() (uint64, error) { return iec.Client.BlockNumber(ctx) }
	number, err := instrumentFunction[uint64](blockNumber, "eth_blockNumber", iec)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (iec *InstrumentedEthClient) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	callContract := func() ([]byte, error) { return iec.Client.CallContract(ctx, call, blockNumber) }
	bytes, err := instrumentFunction[[]byte](callContract, "eth_call", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) CallContractAtHash(
	ctx context.Context,
	msg ethereum.CallMsg,
	blockHash gethcommon.Hash,
) ([]byte, error) {
	callContractAtHash := func() ([]byte, error) { return iec.Client.CallContractAtHash(ctx, msg, blockHash) }
	bytes, err := instrumentFunction[[]byte](callContractAtHash, "eth_call", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) CodeAt(
	ctx context.Context,
	contract gethcommon.Address,
	blockNumber *big.Int,
) ([]byte, error) {
	call := func() ([]byte, error) { return iec.Client.CodeAt(ctx, contract, blockNumber) }
	bytes, err := instrumentFunction[[]byte](call, "eth_getCode", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	estimateGas := func() (uint64, error) { return iec.Client.EstimateGas(ctx, call) }
	gas, err := instrumentFunction[uint64](estimateGas, "eth_estimateGas", iec)
	if err != nil {
		return 0, err
	}
	return gas, nil
}

func (iec *InstrumentedEthClient) FeeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock *big.Int,
	rewardPercentiles []float64,
) (*ethereum.FeeHistory, error) {
	feeHistory := func() (*ethereum.FeeHistory, error) {
		return iec.Client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
	}
	history, err := instrumentFunction[*ethereum.FeeHistory](
		feeHistory,
		"eth_feeHistory",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (iec *InstrumentedEthClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	filterLogs := func() ([]types.Log, error) { return iec.Client.FilterLogs(ctx, query) }
	logs, err := instrumentFunction[[]types.Log](filterLogs, "eth_getLogs", iec)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (iec *InstrumentedEthClient) HeaderByHash(ctx context.Context, hash gethcommon.Hash) (*types.Header, error) {
	headerByHash := func() (*types.Header, error) { return iec.Client.HeaderByHash(ctx, hash) }
	header, err := instrumentFunction[*types.Header](
		headerByHash,
		"eth_getBlockByHash",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (iec *InstrumentedEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	headerByNumber := func() (*types.Header, error) { return iec.Client.HeaderByNumber(ctx, number) }
	header, err := instrumentFunction[*types.Header](
		headerByNumber,
		"eth_getBlockByNumber",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (iec *InstrumentedEthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	networkID := func() (*big.Int, error) { return iec.Client.NetworkID(ctx) }
	id, err := instrumentFunction[*big.Int](networkID, "net_version", iec)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (iec *InstrumentedEthClient) NonceAt(
	ctx context.Context,
	account gethcommon.Address,
	blockNumber *big.Int,
) (uint64, error) {
	nonceAt := func() (uint64, error) { return iec.Client.NonceAt(ctx, account, blockNumber) }
	nonce, err := instrumentFunction[uint64](nonceAt, "eth_getTransactionCount", iec)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (iec *InstrumentedEthClient) PeerCount(ctx context.Context) (uint64, error) {
	peerCount := func() (uint64, error) { return iec.Client.PeerCount(ctx) }
	count, err := instrumentFunction[uint64](peerCount, "net_peerCount", iec)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (iec *InstrumentedEthClient) PendingBalanceAt(ctx context.Context, account gethcommon.Address) (*big.Int, error) {
	pendingBalanceAt := func() (*big.Int, error) { return iec.Client.PendingBalanceAt(ctx, account) }
	balance, err := instrumentFunction[*big.Int](pendingBalanceAt, "eth_getBalance", iec)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (iec *InstrumentedEthClient) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	pendingCallContract := func() ([]byte, error) { return iec.Client.PendingCallContract(ctx, call) }
	bytes, err := instrumentFunction[[]byte](pendingCallContract, "eth_call", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	pendingCodeAt := func() ([]byte, error) { return iec.Client.PendingCodeAt(ctx, account) }
	bytes, err := instrumentFunction[[]byte](pendingCodeAt, "eth_getCode", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	pendingNonceAt := func() (uint64, error) { return iec.Client.PendingNonceAt(ctx, account) }
	nonce, err := instrumentFunction[uint64](
		pendingNonceAt,
		"eth_getTransactionCount",
		iec,
	)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (iec *InstrumentedEthClient) PendingStorageAt(
	ctx context.Context,
	account gethcommon.Address,
	key gethcommon.Hash,
) ([]byte, error) {
	pendingStorageAt := func() ([]byte, error) { return iec.Client.PendingStorageAt(ctx, account, key) }
	bytes, err := instrumentFunction[[]byte](pendingStorageAt, "eth_getStorageAt", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	pendingTransactionCount := func() (uint, error) { return iec.Client.PendingTransactionCount(ctx) }
	count, err := instrumentFunction[uint](
		pendingTransactionCount,
		"eth_getBlockTransactionCountByNumber",
		iec,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (iec *InstrumentedEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// instrumentFunction takes a function that returns a value and an error
	// so we just wrap the SendTransaction method in a function that returns 0 as its value,
	// which we throw out below
	sendTransaction := func() (int, error) { return 0, iec.Client.SendTransaction(ctx, tx) }
	_, err := instrumentFunction[int](sendTransaction, "eth_sendRawTransaction", iec)
	return err
}

func (iec *InstrumentedEthClient) StorageAt(
	ctx context.Context,
	account gethcommon.Address,
	key gethcommon.Hash,
	blockNumber *big.Int,
) ([]byte, error) {
	storageAt := func() ([]byte, error) { return iec.Client.StorageAt(ctx, account, key, blockNumber) }
	bytes, err := instrumentFunction[[]byte](storageAt, "eth_getStorageAt", iec)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (iec *InstrumentedEthClient) SubscribeFilterLogs(
	ctx context.Context,
	query ethereum.FilterQuery,
	ch chan<- types.Log,
) (ethereum.Subscription, error) {
	subscribeFilterLogs := func() (ethereum.Subscription, error) { return iec.Client.SubscribeFilterLogs(ctx, query, ch) }
	subscription, err := instrumentFunction[ethereum.Subscription](
		subscribeFilterLogs,
		"eth_subscribe",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (iec *InstrumentedEthClient) SubscribeNewHead(
	ctx context.Context,
	ch chan<- *types.Header,
) (ethereum.Subscription, error) {
	subscribeNewHead := func() (ethereum.Subscription, error) { return iec.Client.SubscribeNewHead(ctx, ch) }
	subscription, err := instrumentFunction[ethereum.Subscription](
		subscribeNewHead,
		"eth_subscribe",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (iec *InstrumentedEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	suggestGasPrice := func() (*big.Int, error) { return iec.Client.SuggestGasPrice(ctx) }
	gasPrice, err := instrumentFunction[*big.Int](suggestGasPrice, "eth_gasPrice", iec)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (iec *InstrumentedEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	suggestGasTipCap := func() (*big.Int, error) { return iec.Client.SuggestGasTipCap(ctx) }
	gasTipCap, err := instrumentFunction[*big.Int](
		suggestGasTipCap,
		"eth_maxPriorityFeePerGas",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return gasTipCap, nil
}

func (iec *InstrumentedEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	syncProgress := func() (*ethereum.SyncProgress, error) { return iec.Client.SyncProgress(ctx) }
	progress, err := instrumentFunction[*ethereum.SyncProgress](
		syncProgress,
		"eth_syncing",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// We write the instrumentation of this function directly because instrumentFunction[] generic fct only takes a single
// return value
func (iec *InstrumentedEthClient) TransactionByHash(
	ctx context.Context,
	hash gethcommon.Hash,
) (tx *types.Transaction, isPending bool, err error) {
	start := time.Now()
	tx, isPending, err = iec.Client.TransactionByHash(ctx, hash)
	// we count both successful and erroring calls (even though this is not well defined in the spec)
	iec.rpcCallsCollector.AddRPCRequestTotal("eth_getTransactionByHash", iec.clientAndVersion)
	if err != nil {
		return nil, false, err
	}
	rpcRequestDuration := time.Since(start)
	// we only observe the duration of successful calls (even though this is not well defined in the spec)
	iec.rpcCallsCollector.ObserveRPCRequestDurationSeconds(
		float64(rpcRequestDuration),
		"eth_getTransactionByHash",
		iec.clientAndVersion,
	)

	return tx, isPending, nil
}

func (iec *InstrumentedEthClient) TransactionCount(ctx context.Context, blockHash gethcommon.Hash) (uint, error) {
	transactionCount := func() (uint, error) { return iec.Client.TransactionCount(ctx, blockHash) }
	count, err := instrumentFunction[uint](
		transactionCount,
		"eth_getBlockTransactionCountByHash",
		iec,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (iec *InstrumentedEthClient) TransactionInBlock(
	ctx context.Context,
	blockHash gethcommon.Hash,
	index uint,
) (*types.Transaction, error) {
	transactionInBlock := func() (*types.Transaction, error) { return iec.Client.TransactionInBlock(ctx, blockHash, index) }
	tx, err := instrumentFunction[*types.Transaction](
		transactionInBlock,
		"eth_getTransactionByBlockHashAndIndex",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (iec *InstrumentedEthClient) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	transactionReceipt := func() (*types.Receipt, error) { return iec.Client.TransactionReceipt(ctx, txHash) }
	receipt, err := instrumentFunction[*types.Receipt](
		transactionReceipt,
		"eth_getTransactionReceipt",
		iec,
	)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (iec *InstrumentedEthClient) TransactionSender(
	ctx context.Context,
	tx *types.Transaction,
	block gethcommon.Hash,
	index uint,
) (gethcommon.Address, error) {
	transactionSender := func() (gethcommon.Address, error) { return iec.Client.TransactionSender(ctx, tx, block, index) }
	address, err := instrumentFunction[gethcommon.Address](
		transactionSender,
		"eth_getSender",
		iec,
	)
	if err != nil {
		return gethcommon.Address{}, err
	}
	return address, nil
}

// Copied from ethclient.go so make sure to change this implementation if the other one changes!
// We need to do this because this method makes a bunch of internal eth_ calls so copying them
// here forces them to use the instrumented versions instead of ethClient's non instrumented versions
// eg: c.HeaderByNumber(ctx, nil) below calls the instrumented HeaderByNumber implemented in this file.
// if we didn't overwrite EstimateGasPriceAndLimitAndSendTx it would be calling the non instrumented version
// which would be equivalent to having all calls here be c.Client.HeaderByNumber instead of c.HeaderByNumber
//
// EstimateGasPriceAndLimitAndSendTx sends and returns an otherwise identical txn
// to the one provided but with updated gas prices sampled from the existing network
// conditions and an accurate gasLimit
//
// Note: tx must be a to a contract, not an EOA
//
// Slightly modified from: https://github.com/ethereum-optimism/optimism/blob/ec266098641820c50c39c31048aa4e953bece464/batch-submitter/drivers/sequencer/driver.go#L314
func (c *InstrumentedEthClient) EstimateGasPriceAndLimitAndSendTx(
	ctx context.Context,
	tx *types.Transaction,
	tag string,
	value *big.Int,
) (*types.Receipt, error) {
	gasTipCap, err := c.SuggestGasTipCap(ctx)
	if err != nil {
		// If the transaction failed because the backend does not support
		// eth_maxPriorityFeePerGas, fallback to using the default constant.
		// Currently Alchemy is the only backend provider that exposes this
		// method, so in the event their API is unreachable we can fallback to a
		// degraded mode of operation. This also applies to our test
		// environments, as hardhat doesn't support the query either.
		c.Logger.Info("eth_maxPriorityFeePerGas is unsupported by current backend, using fallback gasTipCap")
		gasTipCap = FallbackGasTipCap
	}

	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	gasFeeCap := new(big.Int).Add(header.BaseFee, gasTipCap)

	// The estimated gas limits performed by RawTransact fail semi-regularly
	// with out of gas exceptions. To remedy this we extract the internal calls
	// to perform gas price/gas limit estimation here and add a buffer to
	// account for any network variability.
	gasLimit, err := c.EstimateGas(ctx, ethereum.CallMsg{
		From:      c.AccountAddress,
		To:        tx.To(),
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Value:     value,
		Data:      tx.Data(),
	})

	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, tx.ChainId())
	if err != nil {
		return nil, fmt.Errorf("EstimateGasPriceAndLimitAndSendTx: cannot create transactOpts: %w", err)
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.GasTipCap = gasTipCap
	opts.GasFeeCap = gasFeeCap
	opts.GasLimit = addGasBuffer(gasLimit)

	contract := c.Contracts[*tx.To()]
	// if the contract has not been cached
	if contract == nil {
		// create a dummy bound contract tied to the `to` address of the transaction
		contract = bind.NewBoundContract(*tx.To(), abi.ABI{}, c.Client, c.Client, c.Client)
		// cache the contract for later use
		c.Contracts[*tx.To()] = contract
	}

	tx, err = contract.RawTransact(opts, tx.Data())
	if err != nil {
		return nil, fmt.Errorf("EstimateGasPriceAndLimitAndSendTx: failed to send txn (%s): %w", tag, err)
	}

	receipt, err := c.EnsureTransactionEvaled(
		ctx,
		tx,
		tag,
	)
	if err != nil {
		return nil, err
	}

	return receipt, err
}

// Generic function used to instrument all the eth calls that we make below
func instrumentFunction[T any](
	rpcCall func() (T, error),
	rpcMethodName string,
	iec *InstrumentedEthClient,
) (value T, err error) {
	start := time.Now()
	result, err := rpcCall()
	// we count both successful and erroring calls (even though this is not well defined in the spec)
	iec.rpcCallsCollector.AddRPCRequestTotal(rpcMethodName, iec.clientAndVersion)
	if err != nil {
		return value, err
	}
	rpcRequestDuration := time.Since(start)
	// we only observe the duration of successful calls (even though this is not well defined in the spec)
	iec.rpcCallsCollector.ObserveRPCRequestDurationSeconds(
		float64(rpcRequestDuration),
		rpcMethodName,
		iec.clientAndVersion,
	)
	return result, nil
}

// Not sure why this method is not exposed in the ethclient itself...
// but it is needed to comply with the rpc metrics defined in avs-node spec
// https://eigen.nethermind.io/docs/metrics/metrics-prom-spec
func getClientAndVersion(client *EthClient) string {
	var clientVersion string
	err := client.Client.Client().Call(&clientVersion, "web3_clientVersion")
	if err != nil {
		return "unavailable"
	}
	return clientVersion
}
