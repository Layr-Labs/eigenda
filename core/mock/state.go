package mock

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type ChainDataMock struct {
	mock.Mock

	KeyPairs     map[core.OperatorID]*core.KeyPair
	NumOperators uint8
	Stakes       map[core.QuorumID]map[core.OperatorID]int
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

func MakeOperatorId(id int) core.OperatorID {
	data := [32]byte{uint8(id)}
	return data
}

func NewChainDataMock(stakes map[core.QuorumID]map[core.OperatorID]int) (*ChainDataMock, error) {
	numOperators := 0
	keyPairs := make(map[core.OperatorID]*core.KeyPair)
	for _, oprStakes := range stakes {
		if len(oprStakes) > 255 {
			return nil, errors.New("too many operators")
		}
		if len(oprStakes) > numOperators {
			numOperators = len(oprStakes)
		}
		for opID := range oprStakes {
			if _, ok := keyPairs[opID]; ok {
				continue
			}
			keyPair, err := core.GenRandomBlsKeys()
			if err != nil {
				return nil, err
			}
			keyPairs[opID] = keyPair
		}
	}

	return &ChainDataMock{
		KeyPairs:     keyPairs,
		NumOperators: uint8(numOperators),
		Stakes:       stakes,
	}, nil
}

// MakeChainDataMock creates a ChainDataMock with a given number of operators per quorum
// For example, given
//
//	numOperatorsPerQuorum = map[core.QuorumID]int{
//		 0: 2,
//		 1: 3,
//	}
//
// It will create a ChainDataMock with 2 operators in quorum 0 and 3 operators in quorum 1
// with stakes distributed as
//
//	map[core.QuorumID]map[core.OperatorID]int{
//	  0: {
//		   core.OperatorID{0}: 1,
//		   core.OperatorID{1}: 2,
//	  },
//	  1: {
//		   core.OperatorID{0}: 1,
//		   core.OperatorID{1}: 2,
//		   core.OperatorID{2}: 3,
//	  },
//	}
func MakeChainDataMock(numOperatorsPerQuorum map[core.QuorumID]int) (*ChainDataMock, error) {
	stakes := make(map[core.QuorumID]map[core.OperatorID]int)
	for quorumID, numOpr := range numOperatorsPerQuorum {
		stakes[quorumID] = make(map[core.OperatorID]int)
		for i := 0; i < numOpr; i++ {
			id := MakeOperatorId(i)
			stakes[quorumID][id] = int(i + 1)
		}
	}

	return NewChainDataMock(stakes)
}

func (d *ChainDataMock) GetTotalOperatorState(ctx context.Context, blockNumber uint) *PrivateOperatorState {
	return d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, []core.QuorumID{})
}

func (d *ChainDataMock) GetTotalOperatorStateWithQuorums(ctx context.Context, blockNumber uint, filterQuorums []core.QuorumID) *PrivateOperatorState {
	quorums := filterQuorums
	if len(quorums) == 0 {
		for quorumID := range d.Stakes {
			quorums = append(quorums, quorumID)
		}
	}

	indexedOperators := make(map[core.OperatorID]*core.IndexedOperatorInfo, d.NumOperators)
	storedOperators := make(map[core.OperatorID]*core.OperatorInfo)
	privateOperators := make(map[core.OperatorID]*PrivateOperatorInfo, d.NumOperators)

	aggPubKeys := make(map[core.QuorumID]*core.G1Point)
	for i := 0; i < int(d.NumOperators); i++ {
		id := MakeOperatorId(i)
		stake := 0
		for _, stakesByOp := range d.Stakes {
			if s, ok := stakesByOp[id]; ok {
				stake = s
				break
			}
		}
		host := "0.0.0.0"
		dispersalPort := fmt.Sprintf("3%03v", 2*i)
		retrievalPort := fmt.Sprintf("3%03v", 2*i+1)
		socket := core.MakeOperatorSocket(host, dispersalPort, retrievalPort)

		indexed := &core.IndexedOperatorInfo{
			Socket:   string(socket),
			PubkeyG1: d.KeyPairs[id].GetPubKeyG1(),
			PubkeyG2: d.KeyPairs[id].GetPubKeyG2(),
		}

		stored := &core.OperatorInfo{
			Stake: big.NewInt(int64(stake)),
			Index: uint(i),
		}

		private := &PrivateOperatorInfo{
			OperatorInfo:        stored,
			IndexedOperatorInfo: indexed,
			KeyPair:             d.KeyPairs[id],
			Host:                host,
			DispersalPort:       dispersalPort,
			RetrievalPort:       retrievalPort,
		}

		storedOperators[id] = stored
		indexedOperators[id] = indexed
		privateOperators[id] = private
	}

	totals := make(map[core.QuorumID]*core.OperatorInfo)
	for _, quorumID := range quorums {
		stakesByOp := d.Stakes[quorumID]
		quorumStake := 0
		for _, stake := range stakesByOp {
			quorumStake += stake
		}
		totals[quorumID] = &core.OperatorInfo{
			Stake: big.NewInt(int64(quorumStake)),
			Index: uint(len(stakesByOp)),
		}
	}

	operators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo)
	for _, quorumID := range quorums {
		stakesByOp := d.Stakes[quorumID]

		quorumOperators := make(map[core.OperatorID]*core.OperatorInfo)
		for oprID := range stakesByOp {
			quorumOperators[oprID] = storedOperators[oprID]
		}
		operators[quorumID] = quorumOperators
	}

	operatorState := &core.OperatorState{
		Operators:   operators,
		Totals:      totals,
		BlockNumber: blockNumber,
	}

	for quorumID, operatorsByID := range operators {
		for opID := range operatorsByID {
			if aggPubKeys[quorumID] == nil {
				key := privateOperators[opID].KeyPair.GetPubKeyG1()
				aggPubKeys[quorumID] = key.Deserialize(key.Serialize())
			} else {
				aggPubKeys[quorumID].Add(privateOperators[opID].KeyPair.GetPubKeyG1())
			}
		}
	}

	indexedState := &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: indexedOperators,
		AggKeys:          make(map[core.QuorumID]*core.G1Point),
	}
	for quorumID, apk := range aggPubKeys {
		indexedState.AggKeys[quorumID] = apk
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
