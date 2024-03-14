package geth

import (
	"context"
	"fmt"
	"math/big"

	dacommon "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MultiHomingClient struct {
	RPCs     []*EthClient
	rpcUrls  []string
	NumRetry int
	Logger   logging.Logger
	*FailoverController
}

var _ dacommon.EthClient = (*MultiHomingClient)(nil)

func NewMultiHomingClient(config EthClientConfig, senderAddress gethcommon.Address, logger logging.Logger) (*MultiHomingClient, error) {
	rpcUrls := config.RPCURLs

	controller := NewFailoverController(len(rpcUrls), logger)

	rpcs := make([]*EthClient, len(rpcUrls))
	for i := 0; i < len(rpcUrls); i++ {
		rpcurl := rpcUrls[i]
		rpc, err := NewClient(config, senderAddress, rpcurl, logger)
		if err != nil {
			logger.Info("cannot connect to rpc at start", "url", rpcUrls[i])
			return nil, err
		}
		rpcs[i] = rpc
	}

	return &MultiHomingClient{
		RPCs:               rpcs,
		rpcUrls:            rpcUrls,
		NumRetry:           config.NumRetries,
		FailoverController: controller,
		Logger:             logger,
	}, nil
}

func (m *MultiHomingClient) GetRPCInstance() *EthClient {
	index := m.GetTotalNumberFault() % uint64(len(m.RPCs))
	return m.RPCs[index]
}

func (m *MultiHomingClient) GetAccountAddress() gethcommon.Address {
	return m.GetRPCInstance().GetAccountAddress()
}

func (m *MultiHomingClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().SuggestGasTipCap(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().HeaderByNumber(ctx, number)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().EstimateGas(ctx, msg)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

func (m *MultiHomingClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		err := m.GetRPCInstance().SendTransaction(ctx, tx)
		if err == nil {
			return nil
		}
		m.ProcessError(err)
	}
	return err
}

func (m *MultiHomingClient) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().TransactionReceipt(ctx, txHash)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) BlockNumber(ctx context.Context) (uint64, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().BlockNumber(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

// rest is just inherited
func (m *MultiHomingClient) BalanceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().BalanceAt(ctx, account, blockNumber)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().BlockByHash(ctx, hash)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().BlockByNumber(ctx, number)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().CallContract(ctx, call, blockNumber)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) CallContractAtHash(
	ctx context.Context,
	msg ethereum.CallMsg,
	blockHash gethcommon.Hash,
) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().CallContractAtHash(ctx, msg, blockHash)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) CodeAt(
	ctx context.Context,
	contract gethcommon.Address,
	blockNumber *big.Int,
) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().CodeAt(ctx, contract, blockNumber)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) FeeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock *big.Int,
	rewardPercentiles []float64,
) (*ethereum.FeeHistory, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().FilterLogs(ctx, q)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) HeaderByHash(ctx context.Context, hash gethcommon.Hash) (*types.Header, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().HeaderByHash(ctx, hash)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) NetworkID(ctx context.Context) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().NetworkID(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) NonceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (uint64, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().NonceAt(ctx, account, blockNumber)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

func (m *MultiHomingClient) PeerCount(ctx context.Context) (uint64, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PeerCount(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

func (m *MultiHomingClient) PendingBalanceAt(ctx context.Context, account gethcommon.Address) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingBalanceAt(ctx, account)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingCallContract(ctx, msg)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingCodeAt(ctx, account)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingNonceAt(ctx, account)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}
func (m *MultiHomingClient) PendingStorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingStorageAt(ctx, account, key)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}
func (m *MultiHomingClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().PendingTransactionCount(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

func (m *MultiHomingClient) StorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().StorageAt(ctx, account, key, blockNumber)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	var err error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetry; i++ {
		result, err = m.GetRPCInstance().SubscribeFilterLogs(ctx, q, ch)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return result, err
}

func (m *MultiHomingClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	var err error
	var result ethereum.Subscription
	for i := 0; i < m.NumRetry; i++ {
		result, err = m.GetRPCInstance().SubscribeNewHead(ctx, ch)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return result, err
}

func (m *MultiHomingClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().SuggestGasPrice(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().SyncProgress(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (*types.Transaction, bool, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		tx, isPending, err := m.GetRPCInstance().TransactionByHash(ctx, hash)
		if err == nil {
			return tx, isPending, nil
		}
		m.ProcessError(err)
	}
	return nil, true, err
}

func (m *MultiHomingClient) TransactionCount(ctx context.Context, blockHash gethcommon.Hash) (uint, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().TransactionCount(ctx, blockHash)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return 0, err
}

func (m *MultiHomingClient) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*types.Transaction, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().TransactionInBlock(ctx, blockHash, index)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) TransactionSender(ctx context.Context, tx *types.Transaction, block gethcommon.Hash, index uint) (gethcommon.Address, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().TransactionSender(ctx, tx, block, index)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return gethcommon.Address{}, err
}

func (m *MultiHomingClient) ChainID(ctx context.Context) (*big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().ChainID(ctx)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) GetLatestGasCaps(ctx context.Context) (*big.Int, *big.Int, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		gasTipCap, gasFeeCap, err := m.GetRPCInstance().GetLatestGasCaps(ctx)
		if err == nil {
			return gasTipCap, gasFeeCap, nil
		}
		m.ProcessError(err)
	}
	return nil, nil, err
}

func (m *MultiHomingClient) EstimateGasPriceAndLimitAndSendTx(ctx context.Context, tx *types.Transaction, tag string, value *big.Int) (*types.Receipt, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().EstimateGasPriceAndLimitAndSendTx(ctx, tx, tag, value)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().UpdateGas(ctx, tx, value, gasTipCap, gasFeeCap)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}
func (m *MultiHomingClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().EnsureTransactionEvaled(ctx, tx, tag)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error) {
	var err error
	for i := 0; i < m.NumRetry; i++ {
		result, err := m.GetRPCInstance().EnsureAnyTransactionEvaled(ctx, txs, tag)
		if err == nil {
			return result, nil
		}
		m.ProcessError(err)
	}
	return nil, err
}

func (m *MultiHomingClient) GetNoSendTransactOpts() (*bind.TransactOpts, error) {
	instance := m.GetRPCInstance()
	opts, err := bind.NewKeyedTransactorWithChainID(instance.privateKey, instance.chainID)
	if err != nil {
		return nil, fmt.Errorf("NewClient: cannot create NoSendTransactOpts: %w", err)
	}
	opts.NoSend = true

	return opts, nil
}
