package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// percentage multiplier for gas price. It needs to be >= 10 to properly replace existing transaction
// e.g. 10 means 10% increase
var (
	gasPricePercentageMultiplier = big.NewInt(10)
	hundred                      = big.NewInt(100)
	maxSendTransactionRetry      = 3
	queryTickerDuration          = 3 * time.Second
	ErrTransactionNotBroadcasted = errors.New("transaction not broadcasted")
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

type transaction struct {
	*types.Transaction
	TxID        walletsdk.TxID
	requestedAt time.Time
}

type TxnRequest struct {
	Tx       *types.Transaction
	Tag      string
	Value    *big.Int
	Metadata interface{}

	requestedAt time.Time
	// txAttempts are the transactions that have been attempted to be mined for this request.
	// If a transaction hasn't been confirmed within the timeout and a replacement transaction is sent,
	// the original transaction hash will be kept in this slice
	txAttempts []*transaction
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

	ethClient        EthClient
	wallet           walletsdk.Wallet
	numConfirmations int
	requestChan      chan *TxnRequest
	logger           logging.Logger

	receiptChan         chan *ReceiptOrErr
	queueSize           int
	txnBroadcastTimeout time.Duration
	txnRefreshInterval  time.Duration
	metrics             *TxnManagerMetrics
}

var _ TxnManager = (*txnManager)(nil)

func NewTxnManager(ethClient EthClient, wallet walletsdk.Wallet, numConfirmations, queueSize int, txnBroadcastTimeout time.Duration, txnRefreshInterval time.Duration, logger logging.Logger, metrics *TxnManagerMetrics) TxnManager {
	return &txnManager{
		ethClient:           ethClient,
		wallet:              wallet,
		numConfirmations:    numConfirmations,
		requestChan:         make(chan *TxnRequest, queueSize),
		logger:              logger.With("component", "TxnManager"),
		receiptChan:         make(chan *ReceiptOrErr, queueSize),
		queueSize:           queueSize,
		txnBroadcastTimeout: txnBroadcastTimeout,
		txnRefreshInterval:  txnRefreshInterval,
		metrics:             metrics,
	}
}

func NewTxnRequest(tx *types.Transaction, tag string, value *big.Int, metadata interface{}) *TxnRequest {
	return &TxnRequest{
		Tx:       tx,
		Tag:      tag,
		Value:    value,
		Metadata: metadata,

		requestedAt: time.Now(),
		txAttempts:  make([]*transaction, 0),
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
				t.metrics.ObserveLatency("total", float64(time.Since(req.requestedAt).Milliseconds()))
			}
		}
	}()
	t.logger.Info("started TxnManager")
}

// ProcessTransaction sends the transaction and queues the transaction for monitoring.
// It returns an error if the transaction fails to be confirmed for reasons other than timeouts.
// TxnManager monitors the transaction and resends it with a higher gas price if it is not mined without a timeout until the transaction is confirmed or failed.
func (t *txnManager) ProcessTransaction(ctx context.Context, req *TxnRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Debug("new transaction", "tag", req.Tag, "nonce", req.Tx.Nonce(), "gasFeeCap", req.Tx.GasFeeCap(), "gasTipCap", req.Tx.GasTipCap())

	var txn *types.Transaction
	var txID walletsdk.TxID
	var err error
	retryFromFailure := 0
	for retryFromFailure < maxSendTransactionRetry {
		gasTipCap, gasFeeCap, err := t.ethClient.GetLatestGasCaps(ctx)
		if err != nil {
			return fmt.Errorf("failed to get latest gas caps: %w", err)
		}

		txn, err = t.ethClient.UpdateGas(ctx, req.Tx, req.Value, gasTipCap, gasFeeCap)
		if err != nil {
			return fmt.Errorf("failed to update gas price: %w", err)
		}
		txID, err = t.wallet.SendTransaction(ctx, txn)
		var urlErr *url.Error
		didTimeout := false
		if errors.As(err, &urlErr) {
			didTimeout = urlErr.Timeout()
		}
		if didTimeout || errors.Is(err, context.DeadlineExceeded) {
			t.logger.Warn("failed to send txn due to timeout", "tag", req.Tag, "hash", txn.Hash().Hex(), "numRetries", retryFromFailure, "maxRetry", maxSendTransactionRetry, "err", err)
			retryFromFailure++
			continue
		} else if err != nil {
			return fmt.Errorf("failed to send txn (%s) %s: %w", req.Tag, txn.Hash().Hex(), err)
		} else {
			t.logger.Debug("successfully sent txn", "tag", req.Tag, "txID", txID, "txHash", txn.Hash().Hex())
			break
		}
	}

	if txn == nil || txID == "" {
		return fmt.Errorf("failed to send txn (%s) %s: %w", req.Tag, req.Tx.Hash().Hex(), err)
	}

	req.Tx = txn
	req.txAttempts = append(req.txAttempts, &transaction{
		TxID:        txID,
		Transaction: txn,
		requestedAt: time.Now(),
	})

	t.requestChan <- req
	t.metrics.UpdateTxQueue(len(t.requestChan))
	return nil
}

func (t *txnManager) ReceiptChan() chan *ReceiptOrErr {
	return t.receiptChan
}

// ensureAnyTransactionBroadcasted waits until all given transactions are broadcasted to the network.
func (t *txnManager) ensureAnyTransactionBroadcasted(ctx context.Context, txs []*transaction) error {
	queryTicker := time.NewTicker(queryTickerDuration)
	defer queryTicker.Stop()

	for {
		for _, tx := range txs {
			_, err := t.wallet.GetTransactionReceipt(ctx, tx.TxID)
			if err == nil || errors.Is(err, walletsdk.ErrReceiptNotYetAvailable) {
				t.metrics.ObserveLatency("broadcasted", float64(time.Since(tx.requestedAt).Milliseconds()))
				return nil
			}
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func (t *txnManager) ensureAnyTransactionEvaled(ctx context.Context, txs []*transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(queryTickerDuration)
	defer queryTicker.Stop()
	var receipt *types.Receipt
	var err error
	// transactions that need to be queried. Some transactions will be removed from this map depending on their status.
	txnsToQuery := make(map[walletsdk.TxID]*types.Transaction, len(txs))
	for _, tx := range txs {
		txnsToQuery[tx.TxID] = tx.Transaction
	}

	for {
		for txID, tx := range txnsToQuery {
			receipt, err = t.wallet.GetTransactionReceipt(ctx, txID)
			if err == nil {
				chainTip, err := t.ethClient.BlockNumber(ctx)
				if err == nil {
					if receipt.BlockNumber.Uint64()+uint64(t.numConfirmations) > chainTip {
						t.logger.Debug("transaction has been mined but don't have enough confirmations at current chain tip", "txnBlockNumber", receipt.BlockNumber.Uint64(), "numConfirmations", t.numConfirmations, "chainTip", chainTip)
						break
					} else {
						return receipt, nil
					}
				} else {
					t.logger.Debug("failed to get chain tip while waiting for transaction to mine", "err", err)
				}
			}

			if errors.Is(err, ethereum.NotFound) || errors.Is(err, walletsdk.ErrReceiptNotYetAvailable) {
				t.logger.Debug("Transaction not yet mined", "txID", txID, "txHash", tx.Hash().Hex(), "err", err)
			} else if errors.Is(err, walletsdk.ErrTransactionFailed) {
				t.logger.Debug("Transaction failed", "txID", txID, "txHash", tx.Hash().Hex(), "err", err)
				delete(txnsToQuery, txID)
			} else if errors.Is(err, walletsdk.ErrNotYetBroadcasted) {
				t.logger.Error("Transaction has not been broadcasted to network but attempted to retrieve receipt", "err", err)
			} else {
				t.logger.Debug("Transaction receipt retrieval failed", "err", err)
			}
		}

		if len(txnsToQuery) == 0 {
			return nil, fmt.Errorf("all transactions failed")
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return receipt, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

// monitorTransaction waits until the transaction is confirmed (or failed) and resends it with a higher gas price if it is not mined without a timeout.
// It returns the receipt once the transaction has been confirmed.
// It returns an error if the transaction fails to be sent for reasons other than timeouts.
func (t *txnManager) monitorTransaction(ctx context.Context, req *TxnRequest) (*types.Receipt, error) {
	numSpeedUps := 0
	retryFromFailure := 0

	var receipt *types.Receipt
	var err error

	rpcCallAttempt := func() error {
		t.logger.Debug("monitoring transaction", "txHash", req.Tx.Hash().Hex(), "tag", req.Tag, "nonce", req.Tx.Nonce())

		ctxWithTimeout, cancelBroadcastTimeout := context.WithTimeout(ctx, t.txnBroadcastTimeout)
		defer cancelBroadcastTimeout()

		// Ensure transactions are broadcasted to the network before querying the receipt.
		// This is to avoid querying the receipt of a transaction that hasn't been broadcasted yet.
		// For example, when Fireblocks wallet is used, there may be delays in broadcasting the transaction due to latency from cosigning and MPC operations.
		err = t.ensureAnyTransactionBroadcasted(ctxWithTimeout, req.txAttempts)
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			t.logger.Warn("transaction not broadcasted within timeout", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
			fireblocksWallet, ok := t.wallet.(interface {
				CancelTransactionBroadcast(ctx context.Context, txID walletsdk.TxID) (bool, error)
			})
			if ok {
				// Consider these transactions failed as they haven't been broadcasted within timeout.
				// Cancel these transactions to avoid blocking the next transactions.
				for _, tx := range req.txAttempts {
					cancelled, err := fireblocksWallet.CancelTransactionBroadcast(ctx, tx.TxID)
					if err != nil {
						t.logger.Warn("failed to cancel Fireblocks transaction broadcast", "txID", tx.TxID, "err", err)
					} else if cancelled {
						t.logger.Info("cancelled Fireblocks transaction broadcast because it didn't get broadcasted within timeout", "txID", tx.TxID, "timeout", t.txnBroadcastTimeout.String())
					}
				}
			}
			return ErrTransactionNotBroadcasted
		} else if err != nil {
			t.logger.Error("unexpected error while waiting for Fireblocks transaction to broadcast", "txHash", req.Tx.Hash().Hex(), "err", err)
			return err
		}

		ctxWithTimeout, cancelEvaluationTimeout := context.WithTimeout(ctx, t.txnRefreshInterval)
		defer cancelEvaluationTimeout()
		receipt, err = t.ensureAnyTransactionEvaled(
			ctxWithTimeout,
			req.txAttempts,
		)
		return err
	}

	for {
		err = rpcCallAttempt()
		if err == nil {
			t.metrics.UpdateSpeedUps(numSpeedUps)
			t.metrics.IncrementTxnCount("success")
			return receipt, nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			if receipt != nil {
				t.logger.Warn("transaction has been mined, but hasn't accumulated the required number of confirmations", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
				continue
			}
			t.logger.Warn("transaction not mined within timeout, resending with higher gas price", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "nonce", req.Tx.Nonce())
			newTx, err := t.speedUpTxn(ctx, req.Tx, req.Tag)
			if err != nil {
				t.logger.Error("failed to speed up transaction", "err", err)
				t.metrics.IncrementTxnCount("failure")
				return nil, err
			}
			txID, err := t.wallet.SendTransaction(ctx, newTx)
			if err != nil {
				if retryFromFailure >= maxSendTransactionRetry {
					t.logger.Warn("failed to send txn - retries exhausted", "tag", req.Tag, "txn", req.Tx.Hash().Hex(), "attempt", retryFromFailure, "maxRetry", maxSendTransactionRetry, "err", err)
					t.metrics.IncrementTxnCount("failure")
					return nil, err
				} else {
					t.logger.Warn("failed to send txn - retrying", "tag", req.Tag, "txn", req.Tx.Hash().Hex(), "attempt", retryFromFailure, "maxRetry", maxSendTransactionRetry, "err", err)
				}
				retryFromFailure++
				continue
			}

			t.logger.Debug("successfully sent txn", "tag", req.Tag, "txID", txID, "txHash", newTx.Hash().Hex())
			req.Tx = newTx
			req.txAttempts = append(req.txAttempts, &transaction{
				TxID:        txID,
				Transaction: newTx,
			})
			numSpeedUps++
		} else {
			t.logger.Error("transaction failed", "tag", req.Tag, "txHash", req.Tx.Hash().Hex(), "err", err)
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

	t.logger.Info("increasing gas price", "tag", tag, "txHash", tx.Hash().Hex(), "nonce", tx.Nonce(), "prevGasTipCap", prevGasTipCap, "prevGasFeeCap", prevGasFeeCap, "newGasTipCap", newGasTipCap, "newGasFeeCap", newGasFeeCap)
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
