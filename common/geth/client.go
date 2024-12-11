package geth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	FallbackGasTipCap       = big.NewInt(15000000000)
	ErrCannotGetECDSAPubKey = errors.New("ErrCannotGetECDSAPubKey")
	ErrTransactionFailed    = errors.New("ErrTransactionFailed")
)

type EthClient struct {
	*ethclient.Client
	RPCURL           string
	privateKey       *ecdsa.PrivateKey
	chainID          *big.Int
	AccountAddress   gethcommon.Address
	Contracts        map[gethcommon.Address]*bind.BoundContract
	Logger           logging.Logger
	numConfirmations int
}

var _ common.EthClient = (*EthClient)(nil)

// NewClient creates a new Ethereum client.
// If PrivateKeyString in the config is empty, the client will not be able to send transactions, and it will use the senderAddress to create transactions.
// If PrivateKeyString in the config is not empty, the client will be able to send transactions, and the senderAddress is ignored.
func NewClient(config EthClientConfig, senderAddress gethcommon.Address, rpcIndex int, _logger logging.Logger) (*EthClient, error) {
	if rpcIndex >= len(config.RPCURLs) {
		return nil, fmt.Errorf("NewClient: index out of bound, array size is %v, requested is %v", len(config.RPCURLs), rpcIndex)
	}
	logger := _logger.With("component", "EthClient")

	rpcUrl := config.RPCURLs[rpcIndex]
	chainClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("NewClient: cannot connect to provider: %w", err)
	}
	var privateKey *ecdsa.PrivateKey

	accountAddress := senderAddress
	if len(config.PrivateKeyString) != 0 {
		privateKey, err = crypto.HexToECDSA(config.PrivateKeyString)
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot parse private key: %w", err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

		if !ok {
			logger.Error("cannot get publicKeyECDSA")
			return nil, ErrCannotGetECDSAPubKey
		}
		accountAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	}

	chainIDBigInt, err := chainClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("NewClient: cannot get chainId: %w", err)
	}

	c := &EthClient{
		RPCURL:           rpcUrl,
		privateKey:       privateKey,
		chainID:          chainIDBigInt,
		AccountAddress:   accountAddress,
		Client:           chainClient,
		Contracts:        make(map[gethcommon.Address]*bind.BoundContract),
		Logger:           logger,
		numConfirmations: config.NumConfirmations,
	}

	return c, err
}

func (c *EthClient) GetAccountAddress() gethcommon.Address {
	return c.AccountAddress
}

func NoopSigner(addr gethcommon.Address, tx *types.Transaction) (*types.Transaction, error) {
	return tx, nil
}

func (c *EthClient) GetNoSendTransactOpts() (*bind.TransactOpts, error) {
	if c.privateKey != nil {
		opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot create NoSendTransactOpts: %w", err)
		}
		opts.NoSend = true

		return opts, nil
	}

	if c.AccountAddress.Cmp(gethcommon.Address{}) != 0 {
		return &bind.TransactOpts{
			From:   c.AccountAddress,
			Signer: NoopSigner,
			NoSend: true,
		}, nil
	}

	return nil, errors.New("NewClient: cannot create NoSendTransactOpts: private key and account address are both empty")
}

func (c *EthClient) GetLatestGasCaps(ctx context.Context) (gasTipCap, gasFeeCap *big.Int, err error) {
	gasTipCap, err = c.SuggestGasTipCap(ctx)
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

	// pay 25% more than suggested
	extraTip := big.NewInt(0).Quo(gasTipCap, big.NewInt(4))
	// at least pay extra 2 wei
	if extraTip.Cmp(big.NewInt(2)) == -1 {
		extraTip = big.NewInt(2)
	}
	gasTipCap.Add(gasTipCap, extraTip)

	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	gasFeeCap = getGasFeeCap(gasTipCap, header.BaseFee)
	return
}

func (c *EthClient) UpdateGas(ctx context.Context, tx *types.Transaction, value, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	gasLimit, err := c.Client.EstimateGas(ctx, ethereum.CallMsg{
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

	opts, err := c.GetNoSendTransactOpts()
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.GasTipCap = gasTipCap
	opts.GasFeeCap = gasFeeCap
	opts.GasLimit = addGasBuffer(gasLimit)
	opts.Value = value

	contract := c.Contracts[*tx.To()]
	// if the contract has not been cached
	if contract == nil {
		// create a dummy bound contract tied to the `to` address of the transaction
		contract = bind.NewBoundContract(*tx.To(), abi.ABI{}, c.Client, c.Client, c.Client)
		// cache the contract for later use
		c.Contracts[*tx.To()] = contract
	}
	return contract.RawTransact(opts, tx.Data())
}

// EstimateGasPriceAndLimitAndSendTx sends and returns a transaction receipt.
//
// Note: tx must be a to a contract, not an EOA
func (c *EthClient) EstimateGasPriceAndLimitAndSendTx(
	ctx context.Context,
	tx *types.Transaction,
	tag string,
	value *big.Int,
) (*types.Receipt, error) {
	gasTipCap, gasFeeCap, err := c.GetLatestGasCaps(ctx)
	if err != nil {
		return nil, fmt.Errorf("EstimateGasPriceAndLimitAndSendTx: failed to get gas price for txn (%s): %w", tag, err)
	}

	tx, err = c.UpdateGas(ctx, tx, value, gasTipCap, gasFeeCap)
	if err != nil {
		return nil, fmt.Errorf("EstimateGasPriceAndLimitAndSendTx: failed to update gas for txn (%s): %w", tag, err)
	}

	err = c.SendTransaction(ctx, tx)
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

// EnsureTransactionEvaled waits for tx to be mined on the blockchain and returns the receipt.
// If the context times out but the receipt is available, it returns both receipt and error, noting that the transaction is confirmed but has not accumulated the required number of confirmations.
func (c *EthClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	receipt, err := c.waitMined(ctx, []*types.Transaction{tx})
	if err != nil {
		return receipt, fmt.Errorf("failed to wait for transaction (%s) to mine: %w", tag, err)
	}
	if receipt.Status != 1 {
		c.Logger.Error("Transaction Failed", "tag", tag, "txHash", tx.Hash().Hex(), "status", receipt.Status, "GasUsed", receipt.GasUsed)
		return nil, ErrTransactionFailed
	}
	c.Logger.Debug("transaction confirmed", "txHash", tx.Hash().Hex(), "tag", tag, "gasUsed", receipt.GasUsed, "blockNumber", receipt.BlockNumber)
	return receipt, nil
}

// EnsureAnyTransactionEvaled takes multiple transactions and waits for any of them to be mined on the blockchain and returns the receipt.
// If the context times out but the receipt is available, it returns both receipt and error, noting that the transaction is confirmed but has not accumulated the required number of confirmations.
func (c *EthClient) EnsureAnyTransactionEvaled(ctx context.Context, txs []*types.Transaction, tag string) (*types.Receipt, error) {
	receipt, err := c.waitMined(ctx, txs)
	if err != nil {
		return receipt, fmt.Errorf("EnsureTransactionEvaled: failed to wait for transaction (%s) to mine: %w", tag, err)
	}
	if receipt.Status != 1 {
		c.Logger.Error("Transaction Failed", "tag", tag, "txHash", receipt.TxHash.Hex(), "status", receipt.Status, "GasUsed", receipt.GasUsed)
		return nil, ErrTransactionFailed
	}
	c.Logger.Debug("transaction confirmed", "txHash", receipt.TxHash.Hex(), "tag", tag, "gasUsed", receipt.GasUsed)
	return receipt, nil
}

// waitMined takes multiple transactions and waits for any of them to be mined on the blockchain and returns the receipt.
// If the context times out but the receipt is available, it returns both receipt and error, noting that the transaction is confirmed but has not accumulated the required number of confirmations.
// Taken from https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/bind/util.go#L32,
// but added a check for number of confirmations.
func (c *EthClient) waitMined(ctx context.Context, txs []*types.Transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(3 * time.Second)
	defer queryTicker.Stop()
	var receipt *types.Receipt
	var err error
	for {
		for _, tx := range txs {
			receipt, err = c.TransactionReceipt(ctx, tx.Hash())
			if err == nil {
				chainTip, err := c.BlockNumber(ctx)
				if err == nil {
					if receipt.BlockNumber.Uint64()+uint64(c.numConfirmations) > chainTip {
						c.Logger.Debug("transaction has been mined but doesn't have enough confirmations at current chain head", "txnBlockNumber", receipt.BlockNumber.Uint64(), "numConfirmations", c.numConfirmations, "chainTip", chainTip)
						break
					} else {
						return receipt, nil
					}
				} else {
					c.Logger.Debug("failed to query block height while waiting for transaction to mine", "err", err)
				}
			}

			if errors.Is(err, ethereum.NotFound) {
				c.Logger.Debug("Transaction not yet mined", "txHash", tx.Hash().Hex())
			} else if err != nil {
				c.Logger.Debug("Transaction receipt retrieval failed", "err", err)
			}
		}
		// Wait for the next round.
		select {
		case <-ctx.Done():
			return receipt, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

// getGasFeeCap returns the gas fee cap for a transaction, calculated as:
// gasFeeCap = 2 * baseFee + gasTipCap
// Rationale: https://www.blocknative.com/blog/eip-1559-fees
func getGasFeeCap(gasTipCap *big.Int, baseFee *big.Int) *big.Int {
	return new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)
}

func addGasBuffer(gasLimit uint64) uint64 {
	return 6 * gasLimit / 5 // add 20% buffer to gas limit
}
