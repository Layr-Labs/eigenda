package eth

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
)

type ChainState struct {
	client common.EthClient
	reader corev2.Reader
}

func NewChainState(reader corev2.Reader, client common.EthClient) *ChainState {
	return &ChainState{
		client: client,
		reader: reader,
	}
}

var _ corev2.ChainState = (*ChainState)(nil)

func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator corev2.OperatorID) (*corev2.OperatorState, error) {
	operatorsByQuorum, _, err := cs.reader.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber))

}

func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []corev2.QuorumID) (*corev2.OperatorState, error) {
	operatorsByQuorum, err := cs.reader.GetOperatorStakesForQuorums(ctx, quorums, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber))
}

func (cs *ChainState) GetCurrentBlockNumber() (uint, error) {
	ctx := context.Background()
	header, err := cs.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return uint(header.Number.Uint64()), nil
}

func getOperatorState(operatorsByQuorum corev2.OperatorStakes, blockNumber uint32) (*corev2.OperatorState, error) {
	operators := make(map[corev2.QuorumID]map[corev2.OperatorID]*corev2.OperatorInfo)
	totals := make(map[corev2.QuorumID]*corev2.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[corev2.OperatorID]*corev2.OperatorInfo)

		for ind, op := range quorum {
			operators[quorumID][op.OperatorID] = &corev2.OperatorInfo{
				Stake: op.Stake,
				Index: ind,
			}
			totalStake.Add(totalStake, op.Stake)
		}

		totals[quorumID] = &corev2.OperatorInfo{
			Stake: totalStake,
			Index: uint32(len(quorum)),
		}
	}

	state := &corev2.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: uint(blockNumber),
	}

	return state, nil
}
