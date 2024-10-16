package eth

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/common"
)

type ChainState struct {
	client common.EthClient
	reader chainio.Reader
}

func NewChainState(reader chainio.Reader, client common.EthClient) *ChainState {
	return &ChainState{
		client: client,
		reader: reader,
	}
}

var _ chainio.ChainState = (*ChainState)(nil)

func (cs *ChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator chainio.OperatorID) (*chainio.OperatorState, error) {
	operatorsByQuorum, _, err := cs.reader.GetOperatorStakes(ctx, operator, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	return getOperatorState(operatorsByQuorum, uint32(blockNumber))

}

func (cs *ChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []chainio.QuorumID) (*chainio.OperatorState, error) {
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

func getOperatorState(operatorsByQuorum chainio.OperatorStakes, blockNumber uint32) (*chainio.OperatorState, error) {
	operators := make(map[chainio.QuorumID]map[chainio.OperatorID]*chainio.OperatorInfo)
	totals := make(map[chainio.QuorumID]*chainio.OperatorInfo)

	for quorumID, quorum := range operatorsByQuorum {
		totalStake := big.NewInt(0)
		operators[quorumID] = make(map[chainio.OperatorID]*chainio.OperatorInfo)

		for ind, op := range quorum {
			operators[quorumID][op.OperatorID] = &chainio.OperatorInfo{
				Stake: op.Stake,
				Index: ind,
			}
			totalStake.Add(totalStake, op.Stake)
		}

		totals[quorumID] = &chainio.OperatorInfo{
			Stake: totalStake,
			Index: uint32(len(quorum)),
		}
	}

	state := &chainio.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: uint(blockNumber),
	}

	return state, nil
}
