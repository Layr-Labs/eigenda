package eth

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ChainState implements core.ChainState interface to query on-chain states
// and filter socket updates through event logs to track the socket addresses
// of operators in a map.
type ChainState struct {
	Client common.EthClient
	Reader core.Reader
	logger logging.Logger
	// A cache map of the operator registry, key: operator id, value: socket string
	SocketMap map[core.OperatorID]*string
	// Mutex to access socket map
	socketMu sync.Mutex
	// The previous block number the socket map was updated at, inclusive
	socketPrevBlockNumber atomic.Uint32
}

func NewChainState(reader core.Reader, client common.EthClient, logger logging.Logger) (*ChainState, error) {
	currentBlockNumber, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %v", err)
	}
	cs := &ChainState{
		Client:    client,
		Reader:    reader,
		logger:    logger,
		SocketMap: make(map[core.OperatorID]*string),
	}
	// Set initial block number to current block number
	cs.socketPrevBlockNumber.Store(uint32(currentBlockNumber.Number().Uint64()))
	return cs, nil
}

var _ core.ChainState = (*ChainState)(nil)

// GetOperatorStateByOperator returns the operator state for a given operator at a specific block number.
func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	operatorsByQuorum, _, err := cs.Reader.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for operator %x at block %d: %w", operator, blockNumber, err)
	}

	return cs.getOperatorState(ctx, operatorsByQuorum, uint32(blockNumber))
}

// GetOperatorState returns the operator state for a given quorum at a specific block number.
func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	operatorsByQuorum, err := cs.Reader.GetOperatorStakesForQuorums(ctx, quorums, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for quorums %v at block %d: %v", quorums, blockNumber, err)
	}

	return cs.getOperatorState(ctx, operatorsByQuorum, uint32(blockNumber))
}

// GetCurrentBlockNumber returns the current block number.
func (cs *ChainState) GetCurrentBlockNumber(ctx context.Context) (uint, error) {
	number, err := cs.Client.BlockNumber(ctx)
	if err != nil {
		cs.logger.Warn("failed to get current block number", "error", err)
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}

	return uint(number), nil
}

// GetOperatorSocket returns the socket address for a given operator at the current block number,
// and it takes blockNumber due to the core.ChainState interface.
func (cs *ChainState) GetOperatorSocket(ctx context.Context, blockNumber uint, operator core.OperatorID) (string, error) {
	socket, err := cs.Reader.GetOperatorSocket(ctx, operator)
	if err != nil {
		cs.logger.Warn("failed to get operator socket from Eth Client reader", "operator", operator, "blockNumber", blockNumber, "error", err)
		return "", fmt.Errorf("failed to get socket for operator %x at block %d: %w", operator, blockNumber, err)
	}
	return socket, nil
}

// indexSocketMap filters event logs from the previously checked block number to the current block,
// to identify all socket update events in that block range, and update the socket map accordingly
func (cs *ChainState) indexSocketMap(ctx context.Context) error {
	logs, err := cs.getSocketUpdateEventLogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get logs: %v", err)
	}

	// logs are in order of block number, so we can just iterate through them
	prematureBreak := false
	for _, log := range logs {
		cs.socketPrevBlockNumber.Store(uint32(log.BlockNumber - 1))

		operatorID, socket, err := cs.parseOperatorSocketUpdate(&log)
		if err != nil {
			cs.logger.Warn("failed to get transaction data for operator", "operatorID", operatorID, "error", err)
			continue
		}

		cs.SocketMap[operatorID] = &socket
	}

	// If above loop completed without a premature break, increment prevBlockNum by 1 to ensure we don't handle the same block twice
	if !prematureBreak {
		cs.socketPrevBlockNumber.Add(1)
	}

	return nil
}

func (cs *ChainState) getSocketUpdateEventLogs(ctx context.Context) ([]types.Log, error) {
	currentBlockNumber, err := cs.Reader.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %v", err)
	}

	registryCoordinator, err := cs.Reader.RegistryCoordinator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry coordinator address: %v", err)
	}
	prevBlockNum := cs.socketPrevBlockNumber.Load()
	// The chain hasn't progressed since the last filter, so no logs
	if prevBlockNum >= currentBlockNumber {
		return []types.Log{}, nil
	}
	// Add 1 to prevBlockNum since we already processed that block
	logs, err := cs.Client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(prevBlockNum + 1)),
		ToBlock:   big.NewInt(int64(currentBlockNumber)),
		Addresses: []gcommon.Address{registryCoordinator},
		Topics: [][]gcommon.Hash{
			{common.OperatorSocketUpdateEventSigHash},
		},
	})
	if err != nil {
		cs.logger.Warn("failed to filter logs from block", "FromBlock", prevBlockNum+1, "ToBlock", currentBlockNumber, "error", err)
		return nil, fmt.Errorf("failed to filter logs from block %d to %d: %v", prevBlockNum+1, currentBlockNumber, err)
	}
	if len(logs) == 0 {
		return []types.Log{}, nil
	}
	return logs, nil
}

func (cs *ChainState) getTransaction(ctx context.Context, txHash gcommon.Hash) (*types.Transaction, error) {
	transaction, isPending, err := cs.Client.TransactionByHash(ctx, txHash)
	// Don't continue filtering through logs if we fail to get the transaction or the transaction is pending
	if err != nil {
		cs.logger.Warn("failed to get transaction", "txHash", txHash.Hex(), "error", err)
		return nil, fmt.Errorf("failed to get transaction %s: %v", txHash.Hex(), err)
	}
	if isPending {
		cs.logger.Warn("transaction is still pending for operator socket update event", "txHash", txHash.Hex())
		return nil, fmt.Errorf("transaction %s is still pending for operator socket update event", txHash.Hex())
	}
	return transaction, nil
}

func (cs *ChainState) parseOperatorSocketUpdate(log *types.Log) (core.OperatorID, string, error) {
	operatorID, err := cs.parseOperatorIDFromEventLog(log)
	if err != nil {
		cs.logger.Warn("failed to parse operator ID from log. skipping malformed log",
			"txHash", log.TxHash.Hex(),
			"operatorId", operatorID,
			"error", err)
	}

	socket, err := cs.parseSocketFromEventLog(log.Data)
	if err != nil {
		cs.logger.Warn("failed to parse socket update event. skipping malformed log",
			"txHash", log.TxHash.Hex(),
			"operatorId", operatorID,
			"error", err)
	}
	return operatorID, socket, nil
}

// parseOperatorIDFromEventLog parses the operator ID from a log and returns the operator ID.
func (cs *ChainState) parseOperatorIDFromEventLog(log *types.Log) (core.OperatorID, error) {
	if len(log.Topics) < 2 {
		return core.OperatorID{}, fmt.Errorf("log topics too short: expected at least 2 topics, got %d", len(log.Topics))
	}
	if len(log.Topics[1].Bytes()) != 32 {
		return core.OperatorID{}, fmt.Errorf("operatorID is expecting 32 bytes, got %d", len(log.Topics[1].Bytes()))
	}
	operatorID := core.OperatorID(log.Topics[1].Bytes())
	return operatorID, nil
}

// parseSocketFromEventLog parses the socket string directly from the event log data
func (cs *ChainState) parseSocketFromEventLog(logData []byte) (string, error) {
	rcAbi, err := abi.JSON(bytes.NewReader(common.RegistryCoordinatorAbi))
	if err != nil {
		return "", err
	}

	event := rcAbi.Events["OperatorSocketUpdate"]
	// Used NonIndexed() to specifically get the non-indexed parameters
	values, err := event.Inputs.NonIndexed().Unpack(logData)
	if err != nil {
		return "", fmt.Errorf("failed to unpack event data: %w", err)
	}
	// Expect exactly one value since there's only one non-indexed parameter (the socket)
	if len(values) != 1 {
		return "", fmt.Errorf("unexpected number of values in event data")
	}

	socket, ok := values[0].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert socket to string")
	}

	return socket, nil
}

// refreshSocketMap refresh the socket map for the given operators by quorums at the current block.
func (cs *ChainState) refreshSocketMap(ctx context.Context, operatorsByQuorum core.OperatorStakes) error {
	for _, quorum := range operatorsByQuorum {
		for _, operator := range quorum {
			_, ok := cs.SocketMap[operator.OperatorID]

			if !ok {
				socket, err := cs.Reader.GetOperatorSocket(ctx, operator.OperatorID)
				if err != nil {
					return fmt.Errorf("failed to get socket for operator %x: %w", operator.OperatorID, err)
				}
				cs.SocketMap[operator.OperatorID] = &socket
			}
		}
	}

	// Index for recent socket updates
	if err := cs.indexSocketMap(ctx); err != nil {
		return err
	}
	return nil
}

// getOperatorState returns the current operator state for a given operatorsByQuorum.
// It ensures the socket map is refreshed before returning the state.
// This fucntion locks the socket map upon entry, refreshes the cache, reads
// required values, and then returns the lock.
func (cs *ChainState) getOperatorState(ctx context.Context, operatorsByQuorum core.OperatorStakes, blockNumber uint32) (*core.OperatorState, error) {
	cs.socketMu.Lock()
	defer cs.socketMu.Unlock()
	// Ensure socket map is refreshed before getting operator state
	if err := cs.refreshSocketMap(ctx, operatorsByQuorum); err != nil {
		return nil, fmt.Errorf("failed to refresh socket map: %w", err)
	}

	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)
		for ind, op := range quorum {
			socket, ok := cs.SocketMap[op.OperatorID]
			if !ok || socket == nil {
				return nil, fmt.Errorf("socket not found for operator %x", op.OperatorID)
			}

			operators[quorumID][op.OperatorID] = &core.OperatorInfo{
				Stake:  op.Stake,
				Index:  ind,
				Socket: *socket,
			}
			totalStake.Add(totalStake, op.Stake)
		}

		totals[quorumID] = &core.OperatorInfo{
			Stake: totalStake,
			Index: core.OperatorIndex(len(quorum)),
			// no socket for the total
			Socket: "",
		}
	}

	state := &core.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: uint(blockNumber),
	}

	return state, nil
}
