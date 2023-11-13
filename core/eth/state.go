package eth

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

type ChainState struct {
	Client common.EthClient
	Tx     core.Transactor
}

func NewChainState(tx core.Transactor, client common.EthClient) *ChainState {
	return &ChainState{
		Client: client,
		Tx:     tx,
	}
}

var _ core.ChainState = (*ChainState)(nil)

func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	operatorsByQuorum, quorumIds, err := cs.Tx.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, quorumIds, uint32(blockNumber))

}

func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	operatorsByQuorum, err := cs.Tx.GetOperatorStakesForQuorums(ctx, quorums, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, quorums, uint32(blockNumber))

}

func (cs *ChainState) GetCurrentBlockNumber() (uint, error) {

	ctx := context.Background()
	header, err := cs.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return uint(header.Number.Uint64()), nil

}

func getOperatorState(operatorsByQuorum [][]core.OperatorStake, quorumIds []core.QuorumID, blockNumber uint32) (*core.OperatorState, error) {
	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for i, quorum := range operatorsByQuorum {

		totalStake := big.NewInt(0)

		operators[quorumIds[i]] = make(map[core.OperatorID]*core.OperatorInfo)

		for ind, op := range quorum {

			operators[quorumIds[i]][op.OperatorID] = &core.OperatorInfo{
				Stake: op.Stake,
				Index: core.OperatorIndex(ind),
			}
			totalStake.Add(totalStake, op.Stake)

		}

		totals[quorumIds[i]] = &core.OperatorInfo{
			Stake: totalStake,
			Index: core.OperatorIndex(len(quorum)),
		}
	}

	state := &core.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: uint(blockNumber),
	}

	return state, nil

}
