package geth

import (
	"context"
	"math/big"
	"time"

	dacommon "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MultiHomingClient struct {
	RPCs       []dacommon.EthClient
	rpcUrls    []string
	NumRetries int
	Logger     logging.Logger
	*FailoverController
	Timeout time.Duration // Network timeout is injected in additional to parent context timeout
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

	controller := NewFailoverController(len(rpcUrls), logger)

	client := &MultiHomingClient{
		rpcUrls:            rpcUrls,
		NumRetries:         config.NumRetries,
		FailoverController: controller,
		Logger:             logger,
		Timeout:            config.NetworkTimeout,
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
	index := m.GetTotalNumberFault() % uint64(len(m.RPCs))
	return int(index), m.RPCs[index]
}

func (m *MultiHomingClient) GetAccountAddress() gethcommon.Address {
	_, instance := m.GetRPCInstance()
	return instance.GetAccountAddress()
}

func (m *MultiHomingClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.SuggestGasTipCap(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.HeaderByNumber(instanceCtx, number)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.EstimateGas(instanceCtx, msg)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

func (m *MultiHomingClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		err := instance.SendTransaction(instanceCtx, tx)
		instanceCtxCancel()
		if err == nil {
			return nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return errLast
}

func (m *MultiHomingClient) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.TransactionReceipt(instanceCtx, txHash)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockNumber(ctx context.Context) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.BlockNumber(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

// rest is just inherited
func (m *MultiHomingClient) BalanceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.BalanceAt(instanceCtx, account, blockNumber)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.BlockByHash(instanceCtx, hash)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.BlockByNumber(instanceCtx, number)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
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
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.CallContract(instanceCtx, call, blockNumber)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
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
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.CallContractAtHash(instanceCtx, msg, blockHash)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
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
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.CodeAt(instanceCtx, contract, blockNumber)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
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
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.FeeHistory(instanceCtx, blockCount, lastBlock, rewardPercentiles)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.FilterLogs(instanceCtx, q)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) HeaderByHash(ctx context.Context, hash gethcommon.Hash) (*types.Header, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.HeaderByHash(instanceCtx, hash)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) NetworkID(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.NetworkID(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) NonceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.NonceAt(instanceCtx, account, blockNumber)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

func (m *MultiHomingClient) PeerCount(ctx context.Context) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PeerCount(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

func (m *MultiHomingClient) PendingBalanceAt(ctx context.Context, account gethcommon.Address) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingBalanceAt(instanceCtx, account)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingCallContract(instanceCtx, msg)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingCodeAt(instanceCtx, account)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingNonceAt(instanceCtx, account)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}
func (m *MultiHomingClient) PendingStorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingStorageAt(instanceCtx, account, key)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}
func (m *MultiHomingClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.PendingTransactionCount(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

func (m *MultiHomingClient) StorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.StorageAt(instanceCtx, account, key, blockNumber)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	var errLast error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.SubscribeFilterLogs(instanceCtx, q, ch)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return result, errLast
}

func (m *MultiHomingClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	var errLast error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.SubscribeNewHead(instanceCtx, ch)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return result, errLast
}

func (m *MultiHomingClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.SuggestGasPrice(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.SyncProgress(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (*types.Transaction, bool, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		tx, isPending, err := instance.TransactionByHash(instanceCtx, hash)
		instanceCtxCancel()
		if err == nil {
			return tx, isPending, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, true, errLast
}

func (m *MultiHomingClient) TransactionCount(ctx context.Context, blockHash gethcommon.Hash) (uint, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.TransactionCount(instanceCtx, blockHash)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return 0, errLast
}

func (m *MultiHomingClient) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*types.Transaction, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.TransactionInBlock(instanceCtx, blockHash, index)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) TransactionSender(ctx context.Context, tx *types.Transaction, block gethcommon.Hash, index uint) (gethcommon.Address, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.TransactionSender(instanceCtx, tx, block, index)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return gethcommon.Address{}, errLast
}

func (m *MultiHomingClient) ChainID(ctx context.Context) (*big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.ChainID(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) GetLatestGasCaps(ctx context.Context) (*big.Int, *big.Int, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		gasTipCap, gasFeeCap, err := instance.GetLatestGasCaps(instanceCtx)
		instanceCtxCancel()
		if err == nil {
			return gasTipCap, gasFeeCap, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, nil, errLast
}

func (m *MultiHomingClient) EstimateGasPriceAndLimitAndSendTx(ctx context.Context, tx *types.Transaction, tag string, value *big.Int) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.EstimateGasPriceAndLimitAndSendTx(instanceCtx, tx, tag, value)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.UpdateGas(instanceCtx, tx, value, gasTipCap, gasFeeCap)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}
func (m *MultiHomingClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.EnsureTransactionEvaled(instanceCtx, tx, tag)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error) {
	var errLast error
	for i := 0; i < m.NumRetries+1; i++ {
		_, instance := m.GetRPCInstance()
		instanceCtx, instanceCtxCancel := context.WithTimeout(ctx, m.Timeout)
		result, err := instance.EnsureAnyTransactionEvaled(instanceCtx, txs, tag)
		instanceCtxCancel()
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
		errLast = err
	}
	return nil, errLast
}

func (m *MultiHomingClient) GetNoSendTransactOpts() (*bind.TransactOpts, error) {
	_, instance := m.GetRPCInstance()
	return instance.GetNoSendTransactOpts()
}
