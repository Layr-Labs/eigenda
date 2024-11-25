package eth_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	damock "github.com/Layr-Labs/eigenda/common/mock"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestFetchBatchHeader(t *testing.T) {
	ethClient := &damock.MockEthClient{}
	logger := logging.NewNoopLogger()
	serviceManagerAddress := gcommon.HexToAddress("0x0000000000000000000000000000000000000000")
	batchHeaderHash := []byte("hashhash")
	chainClient := eth.NewChainClient(ethClient, logger)
	topics := [][]gcommon.Hash{
		{common.BatchConfirmedEventSigHash},
		{gcommon.BytesToHash(batchHeaderHash)},
	}
	txHash := gcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	refBlock := 86
	ethClient.On("FilterLogs", ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(refBlock)),
		ToBlock:   nil,
		Addresses: []gcommon.Address{serviceManagerAddress},
		Topics:    topics,
	}).Return([]types.Log{
		{
			Address: serviceManagerAddress,
			Topics: []gcommon.Hash{
				topics[0][0], topics[1][0],
			},
			Data:        []byte{},
			BlockHash:   gcommon.HexToHash("0x0"),
			BlockNumber: 123,
			TxHash:      txHash,
			TxIndex:     0,
			Index:       0,
		},
	}, nil)
	expectedHeader := binding.BatchHeader{
		BlobHeadersRoot:       [32]byte{0},
		QuorumNumbers:         []byte{0},
		SignedStakeForQuorums: []byte{100},
		ReferenceBlockNumber:  uint32(refBlock),
	}
	calldata, err := hex.DecodeString("7794965a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000560000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001c01b4136a161225e9cebe4e2c561148043b2fde423fc5b64e01d897d0fb7970a142d5474fb609bda1b747bdb5c47375d5819000e3c5cbc75baf55b19849410a2610de9c40eb95b49aca940e0bec6ae8b2868855a6324d04d864cbfa61128cf06a51c069e5a0c490c5a359086b0a3660c2ea2e4fb50722bec1ef593c5245413e4cd0a3c7e490348fb279ccb58f91a3bd494511c2ab0321e3922a0cd26012ef3133c043acb758e735db805d360196f3fc89a6395a4b174c19b981afb7f64c2b1193e0000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000026000000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001170c867415fef7db6d88e37598228f43de085616a25939dacbb6b5900f680c7f1d582c9ea38023afb08f368ea93692d17946619d9cf5f3c4d7b3c0cff1a92dff0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000")

	assert.Nil(t, err)
	r, ok := new(big.Int).SetString("8ad2b300a012fb0e90dceb8b66fa564717a2d218ca0fd25f11a1875e0153d1d8", 16)
	assert.True(t, ok)
	s, ok := new(big.Int).SetString("1accb1e1c69fa07bd4237d92143275960b24eec780862a673d54ffaaa5e77f9b", 16)
	assert.True(t, ok)
	ethClient.On("TransactionByHash", txHash).Return(
		types.NewTx(&types.DynamicFeeTx{
			ChainID:    big.NewInt(1),
			Nonce:      1,
			GasTipCap:  big.NewInt(1_000_000),
			GasFeeCap:  big.NewInt(1_000_000),
			Gas:        298617,
			To:         &serviceManagerAddress,
			Value:      big.NewInt(0),
			Data:       calldata,
			AccessList: types.AccessList{},
			V:          big.NewInt(0x1),
			R:          r,
			S:          s,
		}), false, nil)
	batchHeader, err := chainClient.FetchBatchHeader(context.Background(), serviceManagerAddress, batchHeaderHash, big.NewInt(int64(refBlock)), nil)
	assert.Nil(t, err)
	assert.Equal(t, batchHeader.BlobHeadersRoot, expectedHeader.BlobHeadersRoot)
	assert.Equal(t, batchHeader.QuorumNumbers, expectedHeader.QuorumNumbers)
	assert.Equal(t, batchHeader.SignedStakeForQuorums, expectedHeader.SignedStakeForQuorums)
	assert.Equal(t, batchHeader.ReferenceBlockNumber, expectedHeader.ReferenceBlockNumber)
}
