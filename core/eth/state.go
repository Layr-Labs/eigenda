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

func NewChainState(reader core.Reader, client common.EthClient) (*ChainState, error) {
	currentBlockNumber, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	cs := &ChainState{
		Client:    client,
		Reader:    reader,
		SocketMap: make(map[core.OperatorID]*string),
	}
	// Set initial block number to current block number
	cs.socketPrevBlockNumber.Store(uint32(currentBlockNumber.Number().Uint64()))
	return cs, nil
}

var _ core.ChainState = (*ChainState)(nil)

// GetOperatorStateByOperator returns the operator state for a given operator at a specific block number.
func (cs *ChainState) GetOperatorStateByOperator(creader context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	operatorsByQuorum, _, err := cs.Reader.GetOperatorStakes(creader, operator, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for operator %x at block %d: %w", operator, blockNumber, err)
	}

	return cs.getOperatorState(creader, operatorsByQuorum, uint32(blockNumber))
}

// GetOperatorState returns the operator state for a given quorum at a specific block number.
func (cs *ChainState) GetOperatorState(creader context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	operatorsByQuorum, err := cs.Reader.GetOperatorStakesForQuorums(creader, quorums, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for quorums %v at block %d: %w", quorums, blockNumber, err)
	}

	return cs.getOperatorState(creader, operatorsByQuorum, uint32(blockNumber))
}

// GetCurrentBlockNumber returns the current block number.
func (cs *ChainState) GetCurrentBlockNumber(creader context.Context) (uint, error) {
	number, err := cs.Client.BlockNumber(creader)
	if err != nil {
		cs.logger.Warn("failed to get current block number: %w", "error", err)
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}

	return uint(number), nil
}

// GetOperatorSocket returns the socket address for a given operator at the current block number,
// and it takes blockNumber due to the core.ChainState interface.
func (cs *ChainState) GetOperatorSocket(creader context.Context, blockNumber uint, operator core.OperatorID) (string, error) {
	socket, err := cs.Reader.GetOperatorSocket(creader, operator)
	if err != nil {
		cs.logger.Warn("failed to get socket for operator %x at block %d: %w", "operator", operator, "blockNumber", blockNumber, "error", err)
		return "", fmt.Errorf("failed to get socket for operator %x at block %d: %w", operator, blockNumber, err)
	}
	return socket, nil
}

// indexSocketMap filters event logs from the previously checked block number to the current block,
// to identify all socket update events in that block range, and update the socket map accordingly
func (cs *ChainState) indexSocketMap(creader context.Context) error {
	currentBlockNumber, err := cs.Reader.GetCurrentBlockNumber(creader)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	registryCoordinator, err := cs.Reader.RegistryCoordinator(creader)
	if err != nil {
		return fmt.Errorf("failed to get registry coordinator address: %w", err)
	}
	prevBlockNum := cs.socketPrevBlockNumber.Load()
	// The chain hasn't progressed since the last filter, so no need to filter logs
	if prevBlockNum >= currentBlockNumber {
		return nil
	}
	// Add 1 to prevBlockNum since we already processed that block
	logs, err := cs.Client.FilterLogs(creader, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(prevBlockNum + 1)),
		ToBlock:   big.NewInt(int64(currentBlockNumber)),
		Addresses: []gcommon.Address{registryCoordinator},
		Topics: [][]gcommon.Hash{
			{common.OperatorSocketUpdateEventSigHash},
		},
	})
	if err != nil {
		cs.logger.Warn("failed to filter logs from block %d to %d: %w", prevBlockNum+1, currentBlockNumber, err)
		return fmt.Errorf("failed to filter logs from block %d to %d: %w", prevBlockNum+1, currentBlockNumber, err)
	}
	if len(logs) == 0 {
		return nil
	}

	var socketUpdates []*socketUpdateParams

	// logs are in order of block number, so we can just iterate through them
	for _, log := range logs {
		socketUpdate, err := cs.parseSocketUpdateEvent(creader, &log)
		if err != nil {
			cs.logger.Warn("failed to parse socket update event; skipping",
				"error", err,
				"txHash", log.TxHash.Hex(),
				"operatorId", log.Topics[1].Hex(),
				"socket", log.Data)
			continue
		}
		socketUpdates = append(socketUpdates, socketUpdate)
	}

	cs.socketMu.Lock()
	for _, socketUpdate := range socketUpdates {
		cs.SocketMap[socketUpdate.OperatorID] = &socketUpdate.Socket
	}
	cs.socketPrevBlockNumber.Store(currentBlockNumber)
	cs.socketMu.Unlock()

	return nil
}

// parseSocketUpdateEvent parses the socket update event from a log and returns the operator ID and socket address.
func (cs *ChainState) parseSocketUpdateEvent(creader context.Context, log *types.Log) (*socketUpdateParams, error) {
	reader, isPending, err := cs.Client.TransactionByHash(creader, log.TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction %s: %w", log.TxHash.Hex(), err)
	}
	if isPending {
		return nil, fmt.Errorf("transaction %s is still pending for operator socket update event", log.TxHash.Hex())
	}

	calldata := reader.Data()
	// Add length check for method name and input data
	if len(calldata) <= 4 {
		return nil, fmt.Errorf("calldata too short: expected more than 4 bytes for method name and input data, got %d bytes", len(calldata))
	}

	rcAbi, err := abi.JSON(bytes.NewReader(common.RegistryCoordinatorAbi))
	if err != nil {
		return nil, err
	}
	methodSig := calldata[:4]
	method, err := rcAbi.MethodById(methodSig)
	if err != nil {
		return nil, err
	}

	inputs, err := method.Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, err
	}

	var socket string
	if (method.Name == "registerOperator" || method.Name == "registerOperatorWithChurn") && len(inputs) >= 2 {
		socket = inputs[1].(string)
	} else if method.Name == "updateSocket" && len(inputs) >= 1 {
		socket = inputs[0].(string)
	} else {
		// this should never happen; we are going to return nil, so it will be skipped
		return nil, fmt.Errorf("method and input length mismatch for socket update event: %s", method.Name)
	}
	if len(log.Topics) < 2 {
		return nil, fmt.Errorf("log topics too short: expected at least 2 topics, got %d", len(log.Topics))
	}
	if len(log.Topics[1].Bytes()) != 32 {
		return nil, fmt.Errorf("operatorID is expecting 32 bytes, got %d", len(log.Topics[1].Bytes()))
	}
	operatorID := core.OperatorID(log.Topics[1].Bytes())
	return &socketUpdateParams{
		Socket:     socket,
		OperatorID: operatorID,
	}, nil
}

// refreshSocketMap refresh the socket map for the given operators by quorums at the current block.
func (cs *ChainState) refreshSocketMap(creader context.Context, operatorsByQuorum core.OperatorStakes) error {
	for _, quorum := range operatorsByQuorum {
		for _, operator := range quorum {
			cs.socketMu.Lock()
			_, ok := cs.SocketMap[operator.OperatorID]
			cs.socketMu.Unlock()

			if !ok {
				socket, err := cs.Reader.GetOperatorSocket(creader, operator.OperatorID)
				if err != nil {
					return fmt.Errorf("failed to get socket for operator %x: %w", operator.OperatorID, err)
				}
				cs.socketMu.Lock()
				cs.SocketMap[operator.OperatorID] = &socket
				cs.socketMu.Unlock()
			}
		}
	}

	// Index for recent socket updates
	if err := cs.indexSocketMap(creader); err != nil {
		return err
	}
	return nil
}

// getOperatorState returns the current operator state for a given operatorsByQuorum.
// It ensures the socket map is refreshed before returning the state.
func (cs *ChainState) getOperatorState(creader context.Context, operatorsByQuorum core.OperatorStakes, blockNumber uint32) (*core.OperatorState, error) {
	// Ensure socket map is refreshed before getting operator state
	if err := cs.refreshSocketMap(creader, operatorsByQuorum); err != nil {
		return nil, fmt.Errorf("failed to refresh socket map: %w", err)
	}

	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)
		for ind, op := range quorum {
			cs.socketMu.Lock()
			socket, ok := cs.SocketMap[op.OperatorID]
			cs.socketMu.Unlock()
			if !ok || socket == nil {
				return nil, fmt.Errorf("socket not found for operator %x", op.OperatorID)
			}

			operators[quorumID][op.OperatorID] = &core.OperatorInfo{
				Stake:  op.Stake,
				Index:  core.OperatorIndex(ind),
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

type socketUpdateParams struct {
	OperatorID core.OperatorID
	Socket     string
}
