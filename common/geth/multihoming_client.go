package geth

import (
	"context"
	"math/big"
	"sync"

	dacommon "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MultiHomingClient struct {
	RPCs         []dacommon.EthClient
	rpcUrls      []string
	NumRetries   int
	Logger       logging.Logger
	lastRPCIndex uint64
	*FailoverController
	mu sync.Mutex
}

var _ dacommon.EthClient = (*MultiHomingClient)(nil)

// NewMultiHomingClient is an EthClient that automatically handles RPC failures and retries by cycling through
// multiple RPC clients. All EthClients underneath maintain active connections throughout the life time. The
// MultiHomingClient keeps using the same EthClient for a new RPC invocation until it encounters a connection
// error (i.e. any Non EVM error). Then the next EthClient is chosen in a round robin fashion, and the same rpc call
// can be retried. The total number of retry is configured through cli argument. When the rpc call has used up all
// the retry opportunity, the rpc would fail and return error. The MultiHomingClient assumes a single private key.
func NewMultiHomingClient(config EthClientConfig, senderAddress gethcommon.Address, logger logging.Logger) (*MultiHomingClient, error) {
	rpcUrls := config.RPCURLs

	if len(config.RPCURLs) > 1 {
		logger.Info("Fallback chain RPC enabled")
	} else {
		logger.Info("Fallback chain RPC not available")
	}

	FailoverController, err := NewFailoverController(logger, rpcUrls)
	if err != nil {
		return nil, err
	}

	client := &MultiHomingClient{
		rpcUrls:            rpcUrls,
		NumRetries:         config.NumRetries,
		FailoverController: FailoverController,
		lastRPCIndex:       0,
		Logger:             logger.With("component", "MultiHomingClient"),
		mu:                 sync.Mutex{},
	}

	for i := 0; i < len(rpcUrls); i++ {
		rpc, err := NewClient(config, senderAddress, i, logger)
		if err != nil {
			logger.Info("cannot connect to rpc at start", "url", rpcUrls[i])
			return nil, err
		}
		client.RPCs = append(client.RPCs, rpc)
	}

	return client, nil
}

func (m *MultiHomingClient) GetRPCInstance() (int, dacommon.EthClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	index := m.GetTotalNumberRpcFault() % uint64(len(m.RPCs))
	if index != m.lastRPCIndex {
		m.Logger.Info("[MultiHomingClient] Switch RPC", "new index", index, "old index", m.lastRPCIndex)
		m.lastRPCIndex = index
	}
	return int(index), m.RPCs[index]
}

func (m *MultiHomingClient) GetAccountAddress() gethcommon.Address {
	_, instance := m.GetRPCInstance()
	return instance.GetAccountAddress()
}

func (m *MultiHomingClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()
		result, err := instance.SuggestGasTipCap(ctx)
		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SuggestGasTipCap") {
			break
		}
	}
	return nil, errLast
}

func (m *MultiHomingClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.HeaderByNumber(ctx, number)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "HeaderByNumber") {
			break
		}
	}
	return nil, errLast
}

func (m *MultiHomingClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.EstimateGas(ctx, msg)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "EstimateGas") {
			break
		}

	}
	return 0, errLast
}

func (m *MultiHomingClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		err := instance.SendTransaction(ctx, tx)

		if err == nil {
			return nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SendTransaction") {
			break
		}

	}
	return errLast
}

func (m *MultiHomingClient) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.TransactionReceipt(ctx, txHash)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "TransactionReceipt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockNumber(ctx context.Context) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.BlockNumber(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "BlockNumber") {
			break
		}

	}
	return 0, errLast
}

// rest is just inherited
func (m *MultiHomingClient) BalanceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.BalanceAt(ctx, account, blockNumber)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "BalanceAt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.BlockByHash(ctx, hash)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "BlockByHash") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.BlockByNumber(ctx, number)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "BlockByNumber") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.CallContract(ctx, call, blockNumber)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "CallContract") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) CallContractAtHash(
	ctx context.Context,
	msg ethereum.CallMsg,
	blockHash gethcommon.Hash,
) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.CallContractAtHash(ctx, msg, blockHash)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "CallContractAtHash") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) CodeAt(
	ctx context.Context,
	contract gethcommon.Address,
	blockNumber *big.Int,
) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.CodeAt(ctx, contract, blockNumber)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "CodeAt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) FeeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock *big.Int,
	rewardPercentiles []float64,
) (*ethereum.FeeHistory, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "FeeHistory") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.FilterLogs(ctx, q)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "FilterLogs") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) HeaderByHash(ctx context.Context, hash gethcommon.Hash) (*types.Header, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.HeaderByHash(ctx, hash)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "HeaderByHash") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) NetworkID(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.NetworkID(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "NetworkID") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) NonceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.NonceAt(ctx, account, blockNumber)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "NonceAt") {
			break
		}

	}
	return 0, errLast
}

func (m *MultiHomingClient) PeerCount(ctx context.Context) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PeerCount(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PeerCount") {
			break
		}

	}
	return 0, errLast
}

func (m *MultiHomingClient) PendingBalanceAt(ctx context.Context, account gethcommon.Address) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingBalanceAt(ctx, account)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingBalanceAt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingCallContract(ctx, msg)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingCallContract") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingCodeAt(ctx, account)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingCodeAt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingNonceAt(ctx, account)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingNonceAt") {
			break
		}

	}
	return 0, errLast
}
func (m *MultiHomingClient) PendingStorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingStorageAt(ctx, account, key)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingStorageAt") {
			break
		}

	}
	return nil, errLast
}
func (m *MultiHomingClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.PendingTransactionCount(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "PendingTransactionCount") {
			break
		}

	}
	return 0, errLast
}

func (m *MultiHomingClient) StorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.StorageAt(ctx, account, key, blockNumber)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "StorageAt") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	var errLast error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.SubscribeFilterLogs(ctx, q, ch)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SubscribeFilterLogs") {
			break
		}

	}
	return result, errLast
}

func (m *MultiHomingClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	var errLast error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.SubscribeNewHead(ctx, ch)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SubscribeNewHead") {
			break
		}

	}
	return result, errLast
}

func (m *MultiHomingClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.SuggestGasPrice(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SuggestGasPrice") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.SyncProgress(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "SyncProgress") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (*types.Transaction, bool, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		tx, isPending, err := instance.TransactionByHash(ctx, hash)

		if err == nil {
			return tx, isPending, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "TransactionByHash") {
			break
		}

	}
	return nil, true, errLast
}

func (m *MultiHomingClient) TransactionCount(ctx context.Context, blockHash gethcommon.Hash) (uint, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.TransactionCount(ctx, blockHash)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "TransactionCount") {
			break
		}

	}
	return 0, errLast
}

func (m *MultiHomingClient) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*types.Transaction, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.TransactionInBlock(ctx, blockHash, index)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "TransactionInBlock") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) TransactionSender(ctx context.Context, tx *types.Transaction, block gethcommon.Hash, index uint) (gethcommon.Address, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.TransactionSender(ctx, tx, block, index)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "TransactionSender") {
			break
		}

	}
	return gethcommon.Address{}, errLast
}

func (m *MultiHomingClient) ChainID(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.ChainID(ctx)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "ChainID") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) GetLatestGasCaps(ctx context.Context) (*big.Int, *big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		gasTipCap, gasFeeCap, err := instance.GetLatestGasCaps(ctx)

		if err == nil {
			return gasTipCap, gasFeeCap, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "GetLatestGasCaps") {
			break
		}

	}
	return nil, nil, errLast
}

func (m *MultiHomingClient) EstimateGasPriceAndLimitAndSendTx(ctx context.Context, tx *types.Transaction, tag string, value *big.Int) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.EstimateGasPriceAndLimitAndSendTx(ctx, tx, tag, value)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "EstimateGasPriceAndLimitAndSendTx") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.UpdateGas(ctx, tx, value, gasTipCap, gasFeeCap)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "UpdateGas") {
			break
		}

	}
	return nil, errLast
}
func (m *MultiHomingClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.EnsureTransactionEvaled(ctx, tx, tag)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "EnsureTransactionEvaled") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		rpcIndex, instance := m.GetRPCInstance()

		result, err := instance.EnsureAnyTransactionEvaled(ctx, txs, tag)

		if err == nil {
			return result, nil
		}
		errLast = err
		if m.ProcessError(err, rpcIndex, "EnsureAnyTransactionEvaled") {
			break
		}

	}
	return nil, errLast
}

func (m *MultiHomingClient) GetNoSendTransactOpts() (*bind.TransactOpts, error) {
	_, instance := m.GetRPCInstance()
	return instance.GetNoSendTransactOpts()
}
