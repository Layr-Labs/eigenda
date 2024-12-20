package eth

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

type ChainState struct {
	Client common.EthClient
	Tx     core.Reader
}

func NewChainState(tx core.Reader, client common.EthClient) *ChainState {
	return &ChainState{
		Client: client,
		Tx:     tx,
	}
}

var _ core.ChainState = (*ChainState)(nil)

func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	operatorsByQuorum, _, err := cs.Tx.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	socketMap := make(map[core.OperatorID]string)
	socket, err := cs.Tx.GetOperatorSocket(ctx, operator)
	if err != nil {
		return nil, err
	}
	socketMap[operator] = socket

	return getOperatorState(operatorsByQuorum, uint32(blockNumber), socketMap)

}

func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	operatorsByQuorum, err := cs.Tx.GetOperatorStakesForQuorums(ctx, quorums, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	socketMap, err := cs.buildSocketMap(ctx, operatorsByQuorum)
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber), socketMap)
}

func (cs *ChainState) GetCurrentBlockNumber() (uint, error) {
	ctx := context.Background()
	header, err := cs.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return uint(header.Number.Uint64()), nil
}

func (cs *ChainState) GetOperatorSocket(ctx context.Context, blockNumber uint, operator core.OperatorID) (string, error) {
	socket, err := cs.Tx.GetOperatorSocket(ctx, operator)
	if err != nil {
		return "", err
	}
	return socket, nil
}

// buildSocketMap returns a map from operatorID to socket address for the operators in the operatorsByQuorum
func (cs *ChainState) buildSocketMap(ctx context.Context, operatorsByQuorum core.OperatorStakes) (map[core.OperatorID]string, error) {
	socketMap := make(map[core.OperatorID]string)
	for _, quorum := range operatorsByQuorum {
		for _, op := range quorum {
			// if the socket is already in the map, skip
			if _, ok := socketMap[op.OperatorID]; ok {
				continue
			}
			socket, err := cs.Tx.GetOperatorSocket(ctx, op.OperatorID)
			if err != nil {
				return nil, err
			}
			socketMap[op.OperatorID] = socket
		}
	}
	return socketMap, nil
}

func getOperatorState(operatorsByQuorum core.OperatorStakes, blockNumber uint32, socketMap map[core.OperatorID]string) (*core.OperatorState, error) {
	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)

		for ind, op := range quorum {
			operators[quorumID][op.OperatorID] = &core.OperatorInfo{
				Stake:  op.Stake,
				Index:  core.OperatorIndex(ind),
				Socket: socketMap[op.OperatorID],
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
