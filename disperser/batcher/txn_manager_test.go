package batcher_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil)
	ethClient.On("EnsureTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil)

	successful := false
	err = txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
		HandleSuccess: func(ctx context.Context, receipt *types.Receipt) {
			fmt.Println("success function called", successful)
			successful = true
		},
		HandleFailure: func(ctx context.Context, err error) {
			assert.Fail(t, "should not be called")
		},
	})
	<-ctx.Done()
	assert.NoError(t, err)
	assert.True(t, successful)
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 1)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 1)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 1)
}

func TestReplaceGasFee(t *testing.T) {
	ethClient := &mock.MockEthClient{}
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	txnManager := batcher.NewTxnManager(ethClient, 5, 48*time.Second, logger)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	txnManager.Start(ctx)
	txn := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1e18), 100000, big.NewInt(1e9), []byte{})
	ethClient.On("GetLatestGasCaps").Return(big.NewInt(1e9), big.NewInt(1e9), nil)
	ethClient.On("UpdateGas").Return(txn, nil)
	ethClient.On("SendTransaction").Return(nil)
	// assume that the transaction is not mined within the timeout
	ethClient.On("EnsureTransactionEvaled").Return(nil, context.DeadlineExceeded).Once()
	ethClient.On("EnsureTransactionEvaled").Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1),
	}, nil)

	successful := false
	err = txnManager.ProcessTransaction(ctx, &batcher.TxnRequest{
		Tx:    txn,
		Tag:   "test transaction",
		Value: nil,
		HandleSuccess: func(ctx context.Context, receipt *types.Receipt) {
			fmt.Println("success function called", successful)
			successful = true
		},
		HandleFailure: func(ctx context.Context, err error) {
			assert.Fail(t, "should not be called")
		},
	})
	<-ctx.Done()
	assert.NoError(t, err)
	assert.True(t, successful)
	ethClient.AssertNumberOfCalls(t, "GetLatestGasCaps", 2)
	ethClient.AssertNumberOfCalls(t, "UpdateGas", 2)
	ethClient.AssertNumberOfCalls(t, "SendTransaction", 2)
}
