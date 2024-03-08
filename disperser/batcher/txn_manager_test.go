package batcher_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil)
	ethClient.On("EnsureAnyTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil).Once()

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	assert.NoError(t, err)
	receiptOrErr := <-txnManager.ReceiptChan()
	assert.NoError(t, receiptOrErr.Err)
	assert.Equal(t, uint64(1), receiptOrErr.Receipt.BlockNumber.Uint64())
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 1)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 1)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 1)
	ethClient.AssertNumberOfCalls(t, "EnsureAnyTransactionEvaled", 1)

	// now test the case where the transaction fails
	randomErr := errors.New("random error")
	ethClient.On("EnsureAnyTransactionEvaled").Return(nil, randomErr)
	err = txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	receiptOrErr = <-txnManager.ReceiptChan()
	assert.Error(t, receiptOrErr.Err, randomErr)
	assert.Nil(t, receiptOrErr.Receipt)
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 2)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 2)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 2)
	ethClient.AssertNumberOfCalls(t, "EnsureAnyTransactionEvaled", 2)
}

func TestReplaceGasFee(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil)
	// assume that the transaction is not mined within the timeout
	ethClient.On("EnsureAnyTransactionEvaled").Return(nil, context.DeadlineExceeded).Once()
	ethClient.On("EnsureAnyTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil)

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 2)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 2)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 2)
	ethClient.AssertNumberOfCalls(t, "EnsureAnyTransactionEvaled", 2)
}

func TestTransactionFailure(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil).Once()
	// now assume that the transaction fails on retry
	speedUpFailure := errors.New("speed up failure")
	ethClient.On("UpdateGas").Return(nil, speedUpFailure).Once()
	ethClient.On("SendTransaction").Return(nil)
	// assume that the transaction is not mined within the timeout
	ethClient.On("EnsureAnyTransactionEvaled").Return(nil, context.DeadlineExceeded).Once()
	ethClient.On("EnsureAnyTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil)

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	res := <-txnManager.ReceiptChan()
	assert.Error(t, res.Err, speedUpFailure)
}

func TestSendTransactionRetry(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil).Once()
	// assume that it fails to send the replacement transaction once
	ethClient.On("SendTransaction").Return(errors.New("send txn failure")).Once()
	// assume that the transaction is not mined within the timeout
	ethClient.On("EnsureAnyTransactionEvaled").Return(nil, context.DeadlineExceeded).Once()
	ethClient.On("EnsureAnyTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil)

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	res := <-txnManager.ReceiptChan()
	assert.NoError(t, res.Err)
	assert.Equal(t, uint64(1), res.Receipt.BlockNumber.Uint64())
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 2)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 2)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 2)
	ethClient.AssertNumberOfCalls(t, "EnsureAnyTransactionEvaled", 2)
}

func TestSendTransactionRetryFailure(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil).Once()
	// assume that it keeps failing to send the replacement transaction
	sendErr := errors.New("send txn failure")
	ethClient.On("SendTransaction").Return(sendErr)
	// assume that the transaction is not mined within the timeout
	ethClient.On("EnsureAnyTransactionEvaled").Return(nil, context.DeadlineExceeded)

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	res := <-txnManager.ReceiptChan()
	assert.Error(t, res.Err, sendErr)
	assert.Nil(t, res.Receipt)
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 5)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 5)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 5)
	ethClient.AssertNumberOfCalls(t, "EnsureAnyTransactionEvaled", 4)
}
