package eth

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
		return nil, fmt.Errorf("failed to get current block number: %w", err)
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
		return nil, fmt.Errorf("failed to get operator stakes for quorums %v at block %d: %w", quorums, blockNumber, err)
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
	currentBlockNumber, err := cs.Client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	updates, err := cs.Reader.GetSocketUpdates(ctx, uint64(cs.socketPrevBlockNumber.Load()), uint64(currentBlockNumber))
	if err != nil {
		return fmt.Errorf("failed to get socket updates: %w", err)
	}

	for _, update := range updates {
		cs.SocketMap[core.OperatorID(update.OperatorId)] = &update.Socket
	}

	// Store current block as a processed block
	cs.socketPrevBlockNumber.Store(uint32(currentBlockNumber))

	return nil
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
