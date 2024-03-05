package eth

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gcommon "github.com/ethereum/go-ethereum/common"
)

type ChainClient interface {
	FetchBatchHeader(ctx context.Context, serviceManagerAddress gcommon.Address, batchHeaderHash []byte) (*binding.IEigenDAServiceManagerBatchHeader, error)
}

type chainClient struct {
	ethClient common.EthClient
	logger    common.Logger
}

func NewChainClient(ethClient common.EthClient, logger common.Logger) *chainClient {
	return &chainClient{
		ethClient: ethClient,
		logger:    logger,
	}
}

// FetchBatchHeader fetches batch header from chain given a service manager contract address and batch header hash.
// It filters logs by the batch header hashes which are logged as events by the service manager contract.
// From those logs, it identifies corresponding confirmBatch transaction and decodes batch header from the calldata
func (c *chainClient) FetchBatchHeader(ctx context.Context, serviceManagerAddress gcommon.Address, batchHeaderHash []byte) (*binding.IEigenDAServiceManagerBatchHeader, error) {
	logs, err := c.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []gcommon.Address{serviceManagerAddress},
		Topics: [][]gcommon.Hash{
			{common.BatchConfirmedEventSigHash},
			{gcommon.BytesToHash(batchHeaderHash)},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return nil, fmt.Errorf("could not find confirmBatch events for batch header %s", string(batchHeaderHash))
	}

	if len(logs) > 1 {
		c.logger.Error("found more than 1 confirmBatch events", "batchHeader", string(batchHeaderHash))
	}

	txnLog := logs[0]
	tx, isPending, err := c.ethClient.TransactionByHash(ctx, txnLog.TxHash)
	if err != nil {
		return nil, err
	}
	if isPending {
		return nil, fmt.Errorf("confirmBatch transaction pending for batch header %s", string(batchHeaderHash))
	}

	calldata := tx.Data()

	smAbi, err := abi.JSON(bytes.NewReader(common.ServiceManagerAbi))
	if err != nil {
		return nil, err
	}
	methodSig := calldata[:4]
	method, err := smAbi.MethodById(methodSig)
	if err != nil {
		return nil, err
	}

	inputs, err := method.Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, err
	}
	batchHeaderInput := inputs[0].(struct {
		BlobHeadersRoot                  [32]byte "json:\"blobHeadersRoot\""
		QuorumNumbers                    []byte   "json:\"quorumNumbers\""
		ConfirmationThresholdPercentages []byte   "json:\"quorumThresholdPercentages\""
		ReferenceBlockNumber             uint32   "json:\"referenceBlockNumber\""
	})

	return (*binding.IEigenDAServiceManagerBatchHeader)(&batchHeaderInput), nil
}
