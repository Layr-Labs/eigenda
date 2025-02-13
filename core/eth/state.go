package eth

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gcommon "github.com/ethereum/go-ethereum/common"
)

type ChainState struct {
	Client common.EthClient
	Tx     core.Reader
	// A cache map of the operator registry, key: operator id, value: socket string
	SocketMap map[core.OperatorID]*string
	// Mutex to access socket map
	socketMu sync.Mutex
	// The previous block number the socket map was updated at, inclusive
	socketPrevBlockNumber uint32
}

func NewChainState(tx core.Reader, client common.EthClient) (*ChainState, error) {
	currentBlockNumber, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	return &ChainState{
		Client: client,
		Tx:     tx,
		// only start filtering for socket updates from current block
		socketPrevBlockNumber: uint32(currentBlockNumber.Number().Uint64()),
		SocketMap:             make(map[core.OperatorID]*string),
	}, nil
}

var _ core.ChainState = (*ChainState)(nil)

func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	operatorsByQuorum, _, err := cs.Tx.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for operator %x at block %d: %w", operator, blockNumber, err)
	}

	err = cs.refreshSocketMap(ctx, operatorsByQuorum)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh socket map for operator %x at block %d: %w", operator, blockNumber, err)
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber), cs.SocketMap)
}

func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	operatorsByQuorum, err := cs.Tx.GetOperatorStakesForQuorums(ctx, quorums, uint32(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes for quorums %v at block %d: %w", quorums, blockNumber, err)
	}

	err = cs.refreshSocketMap(ctx, operatorsByQuorum)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh socket map for quorums %v at block %d: %w", quorums, blockNumber, err)
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber), cs.SocketMap)
}

func (cs *ChainState) GetCurrentBlockNumber() (uint, error) {
	ctx := context.Background()
	header, err := cs.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block header: %w", err)
	}

	return uint(header.Number.Uint64()), nil
}

func (cs *ChainState) GetOperatorSocket(ctx context.Context, blockNumber uint, operator core.OperatorID) (string, error) {
	socket, err := cs.Tx.GetOperatorSocket(ctx, operator)
	if err != nil {
		return "", fmt.Errorf("failed to get socket for operator %x at block %d: %w", operator, blockNumber, err)
	}
	return socket, nil
}

// buildSocketMap returns a map from operatorID to socket address for the operators in the operatorsByQuorum
func (cs *ChainState) buildSocketMap(ctx context.Context, operatorIds []core.OperatorID) error {
	socketMap := make(map[core.OperatorID]*string)
	for _, operatorID := range operatorIds {
		// if the socket is already in the map, skip
		if _, ok := socketMap[operatorID]; ok {
			continue
		}
		socket, err := cs.Tx.GetOperatorSocket(ctx, operatorID)
		if err != nil {
			return fmt.Errorf("failed to get socket for operator %x: %w", operatorID, err)
		}
		socketMap[operatorID] = &socket
	}

	cs.socketMu.Lock()
	for operatorID, socket := range socketMap {
		cs.SocketMap[operatorID] = socket
	}
	cs.socketMu.Unlock()

	return nil
}

// indexSocketMap indexes a block range of socket update events
func (cs *ChainState) indexSocketMap(ctx context.Context) error {
	currentBlockNumber, err := cs.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	registryCoordinator, err := cs.Tx.RegistryCoordinator(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registry coordinator address: %w", err)
	}

	logs, err := cs.Client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(cs.socketPrevBlockNumber)),
		ToBlock:   big.NewInt(int64(currentBlockNumber)),
		Addresses: []gcommon.Address{registryCoordinator},
		Topics: [][]gcommon.Hash{
			{common.OperatorSocketUpdateEventSigHash},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to filter logs from block %d to %d: %w", cs.socketPrevBlockNumber, currentBlockNumber, err)
	}
	if len(logs) == 0 {
		return nil
	}

	var socketUpdates []*socketUpdateParams

	// logs are in order of block number, so we can just iterate through them
	for _, log := range logs {
		tx, isPending, err := cs.Client.TransactionByHash(ctx, log.TxHash)
		if err != nil {
			return fmt.Errorf("failed to get transaction %s: %w", log.TxHash.Hex(), err)
		}
		if isPending {
			return fmt.Errorf("transaction %s is still pending for operator socket update event", log.TxHash.Hex())
		}

		calldata := tx.Data()
		rcAbi, err := abi.JSON(bytes.NewReader(common.RegistryCoordinatorAbi))
		if err != nil {
			return err
		}
		methodSig := calldata[:4]
		method, err := rcAbi.MethodById(methodSig)
		if err != nil {
			return err
		}
		inputs, err := method.Inputs.Unpack(calldata[4:])
		if err != nil {
			return err
		}

		var socket string
		if method.Name == "registerOperator" || method.Name == "registerOperatorWithChurn" {
			socket = inputs[1].(string)
		} else if method.Name == "updateOperatorSocket" {
			socket = inputs[0].(string)
		} else {
			return fmt.Errorf("unknown method filtered for socket update event: %s", method.Name)
		}
		operatorID := core.OperatorID(log.Topics[1].Bytes())
		socketUpdates = append(socketUpdates, &socketUpdateParams{
			Socket:     socket,
			OperatorID: operatorID,
		})
	}

	cs.socketMu.Lock()
	for _, socketUpdate := range socketUpdates {
		cs.SocketMap[socketUpdate.OperatorID] = &socketUpdate.Socket
	}
	cs.socketPrevBlockNumber = uint32(currentBlockNumber)
	cs.socketMu.Unlock()

	return nil
}

// refreshSocketMap refresh the socket map for the given operators by quorums at the current block.
func (cs *ChainState) refreshSocketMap(ctx context.Context, operatorsByQuorum core.OperatorStakes) error {
	// for all operators in operatorsByQuorum, check if the socket is in the map
	missingOperatorIds := make([]core.OperatorID, 0)
	for _, quorum := range operatorsByQuorum {
		for _, operator := range quorum {
			missingOperatorIds = append(missingOperatorIds, operator.OperatorID)
		}
	}

	if err := cs.buildSocketMap(ctx, missingOperatorIds); err != nil {
		return err
	}

	// Index for recent socket updates
	if err := cs.indexSocketMap(ctx); err != nil {
		return err
	}
	return nil
}

func getOperatorState(operatorsByQuorum core.OperatorStakes, blockNumber uint32, socketMap map[core.OperatorID]*string) (*core.OperatorState, error) {
	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)

		for ind, op := range quorum {
			operators[quorumID][op.OperatorID] = &core.OperatorInfo{
				Stake:  op.Stake,
				Index:  core.OperatorIndex(ind),
				Socket: *socketMap[op.OperatorID],
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
