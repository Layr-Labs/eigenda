package batcher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// percentage multiplier for gas price. It needs to be >= 10 to properly replace existing transaction
// e.g. 10 means 10% increase
var (
	gasPricePercentageMultiplier = big.NewInt(10)
	hundred                      = big.NewInt(100)
)

type TxnRequest struct {
	Tx    *types.Transaction
	Tag   string
	Value *big.Int
}

// ReceiptOrErr is a wrapper for a transaction receipt or an error.
// Receipt should be nil if there is an error, and non-nil if there is no error.
type ReceiptOrErr struct {
	Receipt *types.Receipt
	Err     error
}

// TxnManager receives transactions from the caller, sends them to the chain, and monitors their status.
// It also handles the case where a transaction is not mined within a certain time. In this case, it will
// resend the transaction with a higher gas price. It is assumed that all transactions originate from the
// same account.
type TxnManager struct {
	mu sync.Mutex

	ReceiptChan chan *ReceiptOrErr

	ethClient   common.EthClient
	requestChan chan *TxnRequest
	logger      common.Logger

	queueSize          int
	txnRefreshInterval time.Duration
}

func NewTxnManager(ethClient common.EthClient, queueSize int, txnRefreshInterval time.Duration, logger common.Logger) *TxnManager {
	return &TxnManager{
		ReceiptChan:        make(chan *ReceiptOrErr, queueSize),
		ethClient:          ethClient,
		requestChan:        make(chan *TxnRequest, queueSize),
		logger:             logger,
		queueSize:          queueSize,
		txnRefreshInterval: txnRefreshInterval,
	}
}

func (t *TxnManager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-t.requestChan:
				receipt, err := t.monitorTransaction(ctx, req)
				if err != nil {
					t.ReceiptChan <- &ReceiptOrErr{
						Receipt: nil,
						Err:     err,
					}
				} else {
					t.ReceiptChan <- &ReceiptOrErr{
						Receipt: receipt,
						Err:     nil,
					}
				}
			}
		}
	}()
	t.logger.Info("started TxnManager")
}

// ProcessTransaction sends the transaction and queues the transaction for monitoring.
// It returns an error if the transaction fails to be sent for reasons other than timeouts.
// TxnManager monitors the transaction and resends it with a higher gas price if it is not mined without a timeout.
func (t *TxnManager) ProcessTransaction(ctx context.Context, req *TxnRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Debug("[ProcessTransaction] new transaction", "tag", req.Tag, "nonce", req.Tx.Nonce(), "gasFeeCap", req.Tx.GasFeeCap(), "gasTipCap", req.Tx.GasTipCap())
	gasTipCap, gasFeeCap, err := t.ethClient.GetLatestGasCaps(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest gas caps: %w", err)
	}

	txn, err := t.ethClient.UpdateGas(ctx, req.Tx, req.Value, gasTipCap, gasFeeCap)
	if err != nil {
		return fmt.Errorf("failed to update gas price: %w", err)
	}
	err = t.ethClient.SendTransaction(ctx, txn)
	if err != nil {
		return fmt.Errorf("failed to send txn (%s) %s: %w", req.Tag, req.Tx.Hash().Hex(), err)
	}

	t.requestChan <- req
	return nil
}

// monitorTransaction monitors the transaction and resends it with a higher gas price if it is not mined without a timeout.
// It returns an error if the transaction fails to be sent for reasons other than timeouts.
func (t *TxnManager) monitorTransaction(ctx context.Context, req *TxnRequest) (*types.Receipt, error) {
	for {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, t.txnRefreshInterval)
		defer cancel()

		receipt, err := t.ethClient.EnsureTransactionEvaled(
			ctxWithTimeout,
			req.Tx,
			req.Tag,
		)
		if err == nil {
			return receipt, nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			if receipt != nil {
				t.logger.Warn("transaction has been mined, but hasn't accumulated the required number of confirmations", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
				continue
			}
			t.logger.Warn("transaction not mined within timeout, resending with higher gas price", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
			req.Tx, err = t.speedUpTxn(ctx, req.Tx, req.Tag)
			if err != nil {
				t.logger.Error("failed to speed up transaction", "err", err)
				continue
			}
			err = t.ethClient.SendTransaction(ctx, req.Tx)
			if err != nil {
				t.logger.Error("failed to send txn", "tag", req.Tag, "txn", req.Tx.Hash().Hex(), "err", err)
				continue
			}
		} else {
			t.logger.Error("transaction failed", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "err", err)
			return nil, err
		}
	}
}

// speedUpTxn increases the gas price of the existing transaction by specified percentage.
// It makes sure the new gas price is not lower than the current gas price.
func (t *TxnManager) speedUpTxn(ctx context.Context, tx *types.Transaction, tag string) (*types.Transaction, error) {
	prevGasTipCap := tx.GasTipCap()
	prevGasFeeCap := tx.GasFeeCap()
	// get the gas tip cap and gas fee cap based on current network condition
	currentGasTipCap, currentGasFeeCap, err := t.ethClient.GetLatestGasCaps(ctx)
	if err != nil {
		return nil, err
	}
	increasedGasTipCap := increaseGasPrice(prevGasTipCap)
	increasedGasFeeCap := increaseGasPrice(prevGasFeeCap)
	// make sure increased gas prices are not lower than current gas prices
	var newGasTipCap, newGasFeeCap *big.Int
	if currentGasTipCap.Cmp(increasedGasTipCap) > 0 {
		newGasTipCap = currentGasTipCap
	} else {
		newGasTipCap = increasedGasTipCap
	}
	if currentGasFeeCap.Cmp(increasedGasFeeCap) > 0 {
		newGasFeeCap = currentGasFeeCap
	} else {
		newGasFeeCap = increasedGasFeeCap
	}

	t.logger.Debug("[speedUpTxn] increasing gas price", "tag", tag, "txHash", tx.Hash().Hex(), "nonce", tx.Nonce(), "prevGasTipCap", prevGasTipCap, "prevGasFeeCap", prevGasFeeCap, "newGasTipCap", newGasTipCap, "newGasFeeCap", newGasFeeCap)
	return t.ethClient.UpdateGas(ctx, tx, tx.Value(), newGasTipCap, newGasFeeCap)
}

// increaseGasPrice increases the gas price by specified percentage.
// i.e. gasPrice + (gasPrice * gasPricePercentageMultiplier / 100)
func increaseGasPrice(gasPrice *big.Int) *big.Int {
	bump := new(big.Int).Mul(gasPrice, gasPricePercentageMultiplier)
	bump.Div(bump, hundred)
	return new(big.Int).Add(gasPrice, bump)
}
