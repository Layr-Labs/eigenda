package mock

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type ChainDataMock struct {
	mock.Mock

	KeyPairs     []*core.KeyPair
	NumOperators core.OperatorIndex

	Stakes []int
}

var _ core.ChainState = (*ChainDataMock)(nil)
var _ core.IndexedChainState = (*ChainDataMock)(nil)

type PrivateOperatorInfo struct {
	*core.OperatorInfo
	*core.IndexedOperatorInfo
	KeyPair       *core.KeyPair
	Host          string
	DispersalPort string
	RetrievalPort string
}

type PrivateOperatorState struct {
	*core.OperatorState
	*core.IndexedOperatorState
	PrivateOperators map[core.OperatorID]*PrivateOperatorInfo
}

func makeOperatorId(id int) core.OperatorID {
	data := [32]byte{}
	copy(data[:], []byte(fmt.Sprintf("%d", id)))
	return data
}

func NewChainDataMock(stakes []int) (*ChainDataMock, error) {
	numOperators := uint(len(stakes))

	keyPairs := make([]*core.KeyPair, numOperators)
	for ind := core.OperatorIndex(0); ind < numOperators; ind++ {
		keyPair, err := core.GenRandomBlsKeys()
		if err != nil {
			return nil, err
		}
		keyPairs[ind] = keyPair
	}

	return &ChainDataMock{
		NumOperators: numOperators,
		KeyPairs:     keyPairs,
		Stakes:       stakes,
	}, nil

}

func MakeChainDataMock(numOperators core.OperatorIndex) (*ChainDataMock, error) {

	stakes := make([]int, numOperators)
	for ind := core.OperatorIndex(0); ind < numOperators; ind++ {
		stakes[ind] = int(ind + 1)
	}
	return NewChainDataMock(stakes)

}

func (d *ChainDataMock) GetTotalOperatorState(ctx context.Context, blockNumber uint) *PrivateOperatorState {
	return d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, []core.QuorumID{})
}

func (d *ChainDataMock) GetTotalOperatorStateWithQuorums(ctx context.Context, blockNumber uint, quorums []core.QuorumID) *PrivateOperatorState {
	indexedOperators := make(map[core.OperatorID]*core.IndexedOperatorInfo, d.NumOperators)
	storedOperators := make(map[core.OperatorID]*core.OperatorInfo)
	privateOperators := make(map[core.OperatorID]*PrivateOperatorInfo, d.NumOperators)

	var aggPubKey *core.G1Point

	quorumStake := 0

	for ind := core.OperatorIndex(0); ind < d.NumOperators; ind++ {
		if ind == 0 {
			key := d.KeyPairs[ind].GetPubKeyG1()
			aggPubKey = key.Deserialize(key.Serialize())
		} else {
			aggPubKey.Add(d.KeyPairs[ind].GetPubKeyG1())
		}

		stake := d.Stakes[ind]
		host := "0.0.0.0"
		dispersalPort := fmt.Sprintf("3%03v", int(2*ind))
		retrievalPort := fmt.Sprintf("3%03v", int(2*ind+1))
		socket := core.MakeOperatorSocket(host, dispersalPort, retrievalPort)

		stored := &core.OperatorInfo{
			Stake: big.NewInt(int64(stake)),
			Index: ind,
		}

		indexed := &core.IndexedOperatorInfo{
			Socket:   string(socket),
			PubkeyG1: d.KeyPairs[ind].GetPubKeyG1(),
			PubkeyG2: d.KeyPairs[ind].GetPubKeyG2(),
		}

		private := &PrivateOperatorInfo{
			OperatorInfo:        stored,
			IndexedOperatorInfo: indexed,
			KeyPair:             d.KeyPairs[ind],
			Host:                host,
			DispersalPort:       dispersalPort,
			RetrievalPort:       retrievalPort,
		}

		id := makeOperatorId(int(ind))
		storedOperators[id] = stored
		indexedOperators[id] = indexed
		privateOperators[id] = private

		quorumStake += int(stake)

	}

	totals := map[core.QuorumID]*core.OperatorInfo{
		0: {
			Stake: big.NewInt(int64(quorumStake)),
			Index: d.NumOperators,
		},
		1: {
			Stake: big.NewInt(int64(quorumStake)),
			Index: d.NumOperators,
		},
		2: {
			Stake: big.NewInt(int64(quorumStake)),
			Index: d.NumOperators,
		},
	}

	if len(quorums) > 0 {
		totals = make(map[core.QuorumID]*core.OperatorInfo)
		for _, id := range quorums {
			totals[id] = &core.OperatorInfo{
				Stake: big.NewInt(int64(quorumStake)),
				Index: d.NumOperators,
			}
		}
	}

	operators := map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo{
		0: storedOperators,
		1: storedOperators,
		2: storedOperators,
	}
	if len(quorums) > 0 {
		operators = make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
		for _, id := range quorums {
			operators[id] = storedOperators
		}
	}

	operatorState := &core.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: blockNumber,
	}

	indexedState := &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: indexedOperators,
		AggKeys: map[core.QuorumID]*core.G1Point{
			0: aggPubKey,
			1: aggPubKey,
			2: aggPubKey,
		},
	}

	privateOperatorState := &PrivateOperatorState{
		OperatorState:        operatorState,
		IndexedOperatorState: indexedState,
		PrivateOperators:     privateOperators,
	}

	return privateOperatorState

}

func (d *ChainDataMock) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {

	state := d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, quorums)

	return state.OperatorState, nil

}

func (d *ChainDataMock) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {

	state := d.GetTotalOperatorState(ctx, blockNumber)

	return state.OperatorState, nil

}

func (d *ChainDataMock) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {

	state := d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, quorums)

	return state.IndexedOperatorState, nil

}

func (d *ChainDataMock) GetCurrentBlockNumber() (uint, error) {
	args := d.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (d *ChainDataMock) Start(context.Context) error {
	return nil
}
