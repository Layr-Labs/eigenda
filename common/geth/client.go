package geth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
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
	RPCURL             string
	privateKey         *ecdsa.PrivateKey
	AccountAddress     gethcommon.Address
	NoSendTransactOpts *bind.TransactOpts
	Contracts          map[gethcommon.Address]*bind.BoundContract
	Logger             common.Logger
}

var _ common.EthClient = (*EthClient)(nil)

func NewClient(config EthClientConfig, logger common.Logger) (*EthClient, error) {
	chainClient, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("NewClient: cannot connect to provider: %w", err)
	}
	var accountAddress gethcommon.Address
	var privateKey *ecdsa.PrivateKey
	var opts *bind.TransactOpts

	if len(config.PrivateKeyString) != 0 {
		privateKey, err = crypto.HexToECDSA(config.PrivateKeyString)
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot parse private key: %w", err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

		if !ok {
			logger.Error("NewClient: cannot get publicKeyECDSA")
			return nil, ErrCannotGetECDSAPubKey
		}
		accountAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

		chainIDBigInt, err := chainClient.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot get chainId: %w", err)
		}

		// generate and memoize NoSendTransactOpts
		opts, err = bind.NewKeyedTransactorWithChainID(privateKey, chainIDBigInt)
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot create NoSendTransactOpts: %w", err)
		}
		opts.NoSend = true
	}

	c := &EthClient{
		RPCURL:         config.RPCURL,
		privateKey:     privateKey,
		AccountAddress: accountAddress,
		Client:         chainClient,
		Contracts:      make(map[gethcommon.Address]*bind.BoundContract),
		Logger:         logger,
	}

	c.NoSendTransactOpts = opts

	return c, err
}

func (c *EthClient) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	bn, err := c.Client.BlockNumber(ctx)
	return uint32(bn), err
}

func (c *EthClient) GetAccountAddress() gethcommon.Address {
	return c.AccountAddress
}

func (c *EthClient) GetNoSendTransactOpts() *bind.TransactOpts {
	return c.NoSendTransactOpts
}

// EstimateGasPriceAndLimitAndSendTx sends and returns an otherwise identical txn
// to the one provided but with updated gas prices sampled from the existing network
// conditions and an accurate gasLimit
//
// Note: tx must be a to a contract, not an EOA
//
// Slightly modified from: https://github.com/ethereum-optimism/optimism/blob/ec266098641820c50c39c31048aa4e953bece464/batch-submitter/drivers/sequencer/driver.go#L314
func (c *EthClient) EstimateGasPriceAndLimitAndSendTx(
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

	// pay 25% more than suggested
	extraTip := big.NewInt(0).Quo(gasTipCap, big.NewInt(4))
	// at least pay extra 2 wei
	if extraTip.Cmp(big.NewInt(2)) == -1 {
		extraTip = big.NewInt(2)
	}
	gasTipCap.Add(gasTipCap, extraTip)

	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	gasFeeCap := new(big.Int).Add(header.BaseFee, gasTipCap)

	// The estimated gas limits performed by RawTransact fail semi-regularly
	// with out of gas exceptions. To remedy this we extract the internal calls
	// to perform gas price/gas limit estimation here and add a buffer to
	// account for any network variability.
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

	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, tx.ChainId())
	if err != nil {
		return nil, fmt.Errorf("EstimateGasPriceAndLimitAndSendTx: cannot create transactOpts: %w", err)
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

func (c *EthClient) EnsureTransactionEvaled(ctx context.Context, tx *types.Transaction, tag string) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("EnsureTransactionEvaled: failed to wait for transaction (%s) to mine: %w", tag, err)
	}
	if receipt.Status != 1 {
		c.Logger.Error("Transaction Failed", "tag", tag, "txHash", tx.Hash().Hex(), "status", receipt.Status, "GasUsed", receipt.GasUsed)
		return nil, ErrTransactionFailed
	}
	c.Logger.Trace("successfully submitted transaction", "txHash", tx.Hash().Hex(), "tag", tag, "gasUsed", receipt.GasUsed)
	return receipt, nil
}

func addGasBuffer(gasLimit uint64) uint64 {
	return 6 * gasLimit / 5 // add 20% buffer to gas limit
}
