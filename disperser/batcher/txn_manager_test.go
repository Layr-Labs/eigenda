package batcher_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	sdkmock "github.com/Layr-Labs/eigensdk-go/chainio/clients/mocks"
	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProcessTransaction(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, 100*time.Millisecond, 100*time.Millisecond, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txID := "1234"
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)
	gomock.InOrder(
		w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil),
		w.EXPECT().GetTransactionReceipt(gomock.Any(), gomock.Any()).Return(&types.Receipt{
			BlockNumber: new(big.Int).SetUint64(1),
		}, nil).Times(2),
	)

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	assert.NoError(t, err)
	receiptOrErr := <-txnManager.ReceiptChan()
	assert.NoError(t, receiptOrErr.Err)
	assert.Equal(t, uint64(1), receiptOrErr.Receipt.BlockNumber.Uint64())

	// now test the case where the replacement transaction fails
	randomErr := errors.New("random error")
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil)
	w.EXPECT().GetTransactionReceipt(gomock.Any(), gomock.Any()).Return(nil, randomErr).AnyTimes()
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return("", randomErr).AnyTimes()

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
}

func TestReplaceGasFee(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, 100*time.Millisecond, 100*time.Millisecond, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)

	// assume that the transaction is not mined within the timeout
	badTxID := "1234"
	validTxID := "4321"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(badTxID, nil)
	w.EXPECT().GetTransactionReceipt(gomock.Any(), badTxID).Return(nil, walletsdk.ErrReceiptNotYetAvailable).AnyTimes()
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(validTxID, nil)
	w.EXPECT().GetTransactionReceipt(gomock.Any(), validTxID).Return(&types.Receipt{
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
}

func TestTransactionReplacementFailure(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, time.Second, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil).Once()
	// now assume that the transaction fails on retry
	speedUpFailure := errors.New("speed up failure")
	ethClient.On("UpdateGas").Return(nil, speedUpFailure).Once()

	// assume that the transaction is not mined within the timeout
	badTxID := "1234"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(badTxID, nil)
	w.EXPECT().GetTransactionReceipt(gomock.Any(), badTxID).Return(nil, errors.New("blah")).AnyTimes()

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

func TestSendTransactionReceiptRetry(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, time.Second, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)
	txID := "1234"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil)
	// assume that the transaction is not mined within the timeout
	w.EXPECT().GetTransactionReceipt(gomock.Any(), txID).Return(nil, walletsdk.ErrReceiptNotYetAvailable).Times(3)
	// assume that it fails to send the replacement transaction once
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return("", errors.New("send txn failure"))
	w.EXPECT().GetTransactionReceipt(gomock.Any(), txID).Return(&types.Receipt{
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
}

func TestSendTransactionRetrySuccess(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, time.Second, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)
	txID := "1234"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil)
	// assume that the transaction is not mined within the timeout
	w.EXPECT().GetTransactionReceipt(gomock.Any(), txID).Return(nil, walletsdk.ErrReceiptNotYetAvailable).AnyTimes()

	// assume that it fails to send the replacement transaction once
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return("", errors.New("send txn failure"))
	newTxID := "4321"
	// second try succeeds
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(newTxID, nil)
	w.EXPECT().GetTransactionReceipt(gomock.Any(), newTxID).Return(&types.Receipt{
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
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 3)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 3)
}

func TestSendTransactionRetryFailure(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, time.Second, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)
	txID := "1234"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil)

	// assume that it keeps failing to send the replacement transaction
	sendErr := errors.New("send txn failure")
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return("", sendErr).Times(4)

	// assume that the transaction is not mined within the timeout
	w.EXPECT().GetTransactionReceipt(gomock.Any(), txID).Return(nil, walletsdk.ErrReceiptNotYetAvailable).AnyTimes()

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
}

func TestTransactionNotBroadcasted(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	ctrl := gomock.NewController(t)
	w := sdkmock.NewMockWallet(ctrl)
	logger := logging.NewNoopLogger()
	metrics := batcher.NewMetrics("9100", logger)
	txnManager := batcher.NewTxnManager(ethClient, w, 0, 5, 100*time.Millisecond, 48*time.Second, logger, metrics.TxnManagerMetrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("BlockNumber").Return(uint64(123), nil)
	txID := "1234"
	w.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(txID, nil)

	// assume that the transaction does not get broadcasted to the network
	w.EXPECT().GetTransactionReceipt(gomock.Any(), txID).Return(nil, walletsdk.ErrNotYetBroadcasted).AnyTimes()

	err := txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
	})
	<-ctx.Done()
	assert.NoError(t, err)
	res := <-txnManager.ReceiptChan()
	assert.ErrorAs(t, res.Err, &batcher.ErrTransactionNotBroadcasted)
	assert.Nil(t, res.Receipt)
}
