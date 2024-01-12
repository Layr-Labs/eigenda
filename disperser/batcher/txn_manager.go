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

// TxnManager receives transactions from the caller, sends them to the chain, and monitors their status.
// It also handles the case where a transaction is not mined within a certain time. In this case, it will
// resend the transaction with a higher gas price. It is assumed that all transactions originate from the
// same account.
type TxnManager interface {
	Start(ctx context.Context)
	ProcessTransaction(ctx context.Context, req *TxnRequest) error
	ReceiptChan() chan *ReceiptOrErr
}

type TxnRequest struct {
	Tx       *types.Transaction
	Tag      string
	Value    *big.Int
	Metadata interface{}

	requestedAt time.Time
}

// ReceiptOrErr is a wrapper for a transaction receipt or an error.
// Receipt should be nil if there is an error, and non-nil if there is no error.
// Metadata is the metadata passed in with the transaction request.
type ReceiptOrErr struct {
	Receipt  *types.Receipt
	Metadata interface{}
	Err      error
}

type txnManager struct {
	mu sync.Mutex

	ethClient   common.EthClient
	requestChan chan *TxnRequest
	logger      common.Logger

	receiptChan        chan *ReceiptOrErr
	queueSize          int
	txnRefreshInterval time.Duration
	metrics            *TxnManagerMetrics
}

var _ TxnManager = (*txnManager)(nil)

func NewTxnManager(ethClient common.EthClient, queueSize int, txnRefreshInterval time.Duration, logger common.Logger, metrics *TxnManagerMetrics) TxnManager {
	return &txnManager{
		ethClient:          ethClient,
		requestChan:        make(chan *TxnRequest, queueSize),
		logger:             logger,
		receiptChan:        make(chan *ReceiptOrErr, queueSize),
		queueSize:          queueSize,
		txnRefreshInterval: txnRefreshInterval,
		metrics:            metrics,
	}
}

func NewTxnRequest(tx *types.Transaction, tag string, value *big.Int, metadata interface{}) *TxnRequest {
	return &TxnRequest{
		Tx:       tx,
		Tag:      tag,
		Value:    value,
		Metadata: metadata,

		requestedAt: time.Now(),
	}
}

func (t *txnManager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-t.requestChan:
				receipt, err := t.monitorTransaction(ctx, req)
				if err != nil {
					t.receiptChan <- &ReceiptOrErr{
						Receipt:  nil,
						Metadata: req.Metadata,
						Err:      err,
					}
				} else {
					t.receiptChan <- &ReceiptOrErr{
						Receipt:  receipt,
						Metadata: req.Metadata,
						Err:      nil,
					}
					if receipt.GasUsed > 0 {
						t.metrics.UpdateGasUsed(receipt.GasUsed)
					}
				}
				t.metrics.ObserveLatency(float64(time.Since(req.requestedAt).Milliseconds()))
			}
		}
	}()
	t.logger.Info("started TxnManager")
}

// ProcessTransaction sends the transaction and queues the transaction for monitoring.
// It returns an error if the transaction fails to be sent for reasons other than timeouts.
// TxnManager monitors the transaction and resends it with a higher gas price if it is not mined without a timeout.
func (t *txnManager) ProcessTransaction(ctx context.Context, req *TxnRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Debug("[TxnManager] new transaction", "tag", req.Tag, "nonce", req.Tx.Nonce(), "gasFeeCap", req.Tx.GasFeeCap(), "gasTipCap", req.Tx.GasTipCap())
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
	} else {
		t.logger.Debug("[TxnManager] successfully sent txn", "tag", req.Tag, "txn", txn.Hash().Hex())
	}
	req.Tx = txn

	t.requestChan <- req
	t.metrics.UpdateTxQueue(len(t.requestChan))
	return nil
}

func (t *txnManager) ReceiptChan() chan *ReceiptOrErr {
	return t.receiptChan
}

// monitorTransaction monitors the transaction and resends it with a higher gas price if it is not mined without a timeout.
// It returns an error if the transaction fails to be sent for reasons other than timeouts.
func (t *txnManager) monitorTransaction(ctx context.Context, req *TxnRequest) (*types.Receipt, error) {
	numSpeedUps := 0
	for {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, t.txnRefreshInterval)
		defer cancel()

		t.logger.Debug("[TxnManager] monitoring transaction", "txHash", req.Tx.Hash().Hex(), "tag", req.Tag, "nonce", req.Tx.Nonce())
		receipt, err := t.ethClient.EnsureTransactionEvaled(
			ctxWithTimeout,
			req.Tx,
			req.Tag,
		)
		if err == nil {
			t.metrics.UpdateSpeedUps(numSpeedUps)
			t.metrics.IncrementTxnCount("success")
			return receipt, nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			if receipt != nil {
				t.logger.Warn("[TxnManager] transaction has been mined, but hasn't accumulated the required number of confirmations", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
				continue
			}
			t.logger.Warn("[TxnManager] transaction not mined within timeout, resending with higher gas price", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
			newTx, err := t.speedUpTxn(ctx, req.Tx, req.Tag)
			if err != nil {
				t.logger.Error("[TxnManager] failed to speed up transaction", "err", err)
				t.metrics.IncrementTxnCount("failure")
				return nil, err
			}
			err = t.ethClient.SendTransaction(ctx, newTx)
			if err != nil {
				t.logger.Error("[TxnManager] failed to send txn", "tag", req.Tag, "txn", req.Tx.Hash().Hex(), "err", err)
				t.metrics.IncrementTxnCount("failure")
				return nil, err
			} else {
				t.logger.Debug("[TxnManager] successfully sent txn", "tag", req.Tag, "txn", newTx.Hash().Hex())
			}
			req.Tx = newTx
			numSpeedUps++
		} else {
			t.logger.Error("[TxnManager] transaction failed", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "err", err)
			t.metrics.IncrementTxnCount("failure")
			return nil, err
		}
	}
}

// speedUpTxn increases the gas price of the existing transaction by specified percentage.
// It makes sure the new gas price is not lower than the current gas price.
func (t *txnManager) speedUpTxn(ctx context.Context, tx *types.Transaction, tag string) (*types.Transaction, error) {
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

	t.logger.Info("[TxnManager] increasing gas price", "tag", tag, "txHash", tx.Hash().Hex(), "nonce", tx.Nonce(), "prevGasTipCap", prevGasTipCap, "prevGasFeeCap", prevGasFeeCap, "newGasTipCap", newGasTipCap, "newGasFeeCap", newGasFeeCap)
	return t.ethClient.UpdateGas(ctx, tx, tx.Value(), newGasTipCap, newGasFeeCap)
}

// increaseGasPrice increases the gas price by specified percentage.
// i.e. gasPrice + ((gasPrice * gasPricePercentageMultiplier + 99) / 100)
func increaseGasPrice(gasPrice *big.Int) *big.Int {
	if gasPrice == nil {
		return nil
	}
	bump := new(big.Int).Mul(gasPrice, gasPricePercentageMultiplier)
	bump = roundUpDivideBig(bump, hundred)
	return new(big.Int).Add(gasPrice, bump)
}

func roundUpDivideBig(a, b *big.Int) *big.Int {
	if a == nil || b == nil || b.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	one := new(big.Int).SetUint64(1)
	num := new(big.Int).Sub(new(big.Int).Add(a, b), one) // a + b - 1
	res := new(big.Int).Div(num, b)                      // (a + b - 1) / b
	return res
}
