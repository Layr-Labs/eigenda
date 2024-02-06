package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	cm "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/indexer/eth"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	ttfMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	logger            = &cm.Logger{}
	blockNumber int64 = 17320293
)

func TestHeaderService_PullNewHeaders(t *testing.T) {
	ctx := context.Background()

	pullNewHeaders := func(
		input indexer.Header,
		expected []indexer.Header,
		expecIsHead bool,
		expecErr error,
		prepare func() indexer.HeaderService) func(t *testing.T) {
		return func(t *testing.T) {
			srv := prepare()
			got, isHead, err := srv.PullNewHeaders(&input)
			if expecErr != nil {
				require.NotNil(t, err)
				assert.EqualError(t, err, expecErr.Error())
				return
			}
			require.Nil(t, err, "Error should be nil")
			require.NotNil(t, got, "Got should not be nil")
			assert.Equal(t, len(expected), len(got), "Length of expected and got should be equal")
			assert.Equal(t, expected[0].Number, got[0].Number, "Number not equal to expected")
			assert.Equal(t, expected[0].Finalized, got[0].Finalized, "Finalized not equal to expected")
			assert.Equal(t, expecIsHead, isHead, "isHead not equal to expected")
		}
	}

	t.Run("Pull new headers successfully",
		pullNewHeaders(
			indexer.Header{Number: uint64(blockNumber - 1)},
			[]indexer.Header{
				{
					Number:    uint64(blockNumber),
					Finalized: false,
				},
			},
			false,
			nil,
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Once().Return(nil)

				batchElems := make([]rpc.BatchElem, 0, 1)
				batchElems = append(batchElems, rpc.BatchElem{
					Method: "eth_getBlockByNumber",
					Args:   []interface{}{hexutil.EncodeBig(big.NewInt(blockNumber)), false},
					Result: new(types.Header),
				})

				mockRPCEthClient.On("BatchCallContext", ctx, batchElems).
					Run(func(args ttfMock.Arguments) {
						args[1].([]rpc.BatchElem)[0].Result = &types.Header{
							Number: big.NewInt(blockNumber),
						}
					}).Once().Return(nil)

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull new headers with errors at getting latest header",
		pullNewHeaders(
			indexer.Header{},
			[]indexer.Header{},
			false,
			errors.New("fake error"),
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Once().Return(errors.New("fake error"))

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull new headers returning latest header",
		pullNewHeaders(
			indexer.Header{Number: uint64(blockNumber)},
			[]indexer.Header{
				{
					Number:    uint64(blockNumber),
					Finalized: false,
				},
			},
			true,
			nil,
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Once().Return(nil)

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull new headers with errors at batch call",
		pullNewHeaders(
			indexer.Header{Number: uint64(blockNumber - 1)},
			[]indexer.Header{},
			false,
			errors.New("fake error"),
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Once().Return(nil)

				batchElems := make([]rpc.BatchElem, 0, 1)
				batchElems = append(batchElems, rpc.BatchElem{
					Method: "eth_getBlockByNumber",
					Args:   []interface{}{hexutil.EncodeBig(big.NewInt(blockNumber)), false},
					Result: new(types.Header),
				})

				mockRPCEthClient.On("BatchCallContext", ctx, batchElems).Once().Return(errors.New("fake error"))

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull new headers with errors at batch elems",
		pullNewHeaders(
			indexer.Header{Number: uint64(blockNumber - 1)},
			[]indexer.Header{},
			false,
			errors.New("fake error"),
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Once().Return(nil)

				batchElems := make([]rpc.BatchElem, 0, 1)
				batchElems = append(batchElems, rpc.BatchElem{
					Method: "eth_getBlockByNumber",
					Args:   []interface{}{hexutil.EncodeBig(big.NewInt(blockNumber)), false},
					Result: new(types.Header),
				})

				mockRPCEthClient.On("BatchCallContext", ctx, batchElems).
					Run(func(args ttfMock.Arguments) {
						args[1].([]rpc.BatchElem)[0].Error = errors.New("fake error")
					}).Once().Return(nil)

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

}

func TestHeaderService_PullLatestHeader(t *testing.T) {
	ctx := context.Background()

	pullLatestHeader := func(
		input bool,
		expected indexer.Header,
		expecErr error,
		prepare func() indexer.HeaderService) func(t *testing.T) {
		return func(t *testing.T) {
			srv := prepare()
			got, err := srv.PullLatestHeader(input)
			if expecErr != nil {
				require.NotNil(t, err)
				assert.EqualError(t, err, expecErr.Error())
				return
			}
			require.Nil(t, err, "Error should be nil")
			require.NotNil(t, got, "Got should not be nil")
			assert.Equal(t, expected.Number, got.Number, "Number not equal to expected")
			assert.Equal(t, expected.Finalized, got.Finalized, "Finalized not equal to expected")
		}
	}

	t.Run("Pull latest header successfully",
		pullLatestHeader(
			false,
			indexer.Header{
				Number:    uint64(blockNumber),
				Finalized: false,
			},
			nil,
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)

				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Once().Return(nil)

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull latest header with errors at getting latest header",
		pullLatestHeader(
			false,
			indexer.Header{},
			errors.New("fake error"),
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)

				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Return(errors.New("fake error")).Once()

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull latest finalized header successfully",
		pullLatestHeader(
			true,
			indexer.Header{
				Number:    uint64(blockNumber - eth.DistanceFromHead),
				Finalized: true,
			},
			nil,
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)

				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = big.NewInt(blockNumber)
					}).Return(nil).Once()

				blockNumBig := big.NewInt(blockNumber - eth.DistanceFromHead)
				blockEncoded := hexutil.EncodeBig(blockNumBig)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", blockEncoded, false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = blockNumBig
					}).
					Return(nil).Once()

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))

	t.Run("Pull latest header with errors at getting latest finalized header",
		pullLatestHeader(
			true,
			indexer.Header{},
			errors.New("fake error"),
			func() indexer.HeaderService {
				mockRPCEthClient := new(cm.MockRPCEthClient)

				blockNumBig := big.NewInt(blockNumber)
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", "latest", false).
					Run(func(args ttfMock.Arguments) {
						args[1].(*types.Header).Number = blockNumBig
					}).Return(nil).Once()

				blockEncoded := hexutil.EncodeBig(big.NewInt(blockNumBig.Int64() - eth.DistanceFromHead))
				mockRPCEthClient.On("CallContext", ctx, &types.Header{}, "eth_getBlockByNumber", blockEncoded, false).
					Return(errors.New("fake error")).Once()

				return eth.NewHeaderService(logger, mockRPCEthClient)
			},
		))
}
