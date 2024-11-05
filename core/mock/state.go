package mock

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type ChainDataMock struct {
	mock.Mock

	KeyPairs  map[core.OperatorID]*core.KeyPair
	Operators []core.OperatorID
	Stakes    map[core.QuorumID]map[core.OperatorID]int
}

var _ core.ChainState = (*ChainDataMock)(nil)
var _ core.IndexedChainState = (*ChainDataMock)(nil)

type PrivateOperatorInfo struct {
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
	var data [32]byte
	binary.LittleEndian.PutUint64(data[:8], uint64(id))
	return data
}

func NewChainDataMock(stakes map[core.QuorumID]map[core.OperatorID]int) (*ChainDataMock, error) {

	seenOperators := make(map[core.OperatorID]struct{})
	for _, oprStakes := range stakes {
		for opID := range oprStakes {
			if _, ok := seenOperators[opID]; ok {
				continue
			}
			seenOperators[opID] = struct{}{}
		}
	}

	operators := make([]core.OperatorID, 0, len(seenOperators))
	for opID := range seenOperators {
		operators = append(operators, opID)
	}

	sort.Slice(operators, func(i, j int) bool {
		return operators[i].Hex() < operators[j].Hex()
	})

	keyPairs := make(map[core.OperatorID]*core.KeyPair)
	for _, opID := range operators {
		keyPair, err := core.GenRandomBlsKeys()
		if err != nil {
			return nil, err
		}
		keyPairs[opID] = keyPair
	}

	return &ChainDataMock{
		KeyPairs:  keyPairs,
		Operators: operators,
		Stakes:    stakes,
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

	indexedOperators := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(d.Operators))
	privateOperators := make(map[core.OperatorID]*PrivateOperatorInfo, len(d.Operators))

	aggPubKeys := make(map[core.QuorumID]*core.G1Point)
	for i, id := range d.Operators {

		host := "0.0.0.0"
		dispersalPort := fmt.Sprintf("3%03v", 2*i)
		retrievalPort := fmt.Sprintf("3%03v", 2*i+1)
		socket := core.MakeOperatorSocket(host, dispersalPort, retrievalPort)

		indexed := &core.IndexedOperatorInfo{
			Socket:   string(socket),
			PubkeyG1: d.KeyPairs[id].GetPubKeyG1(),
			PubkeyG2: d.KeyPairs[id].GetPubKeyG2(),
		}

		private := &PrivateOperatorInfo{
			IndexedOperatorInfo: indexed,
			KeyPair:             d.KeyPairs[id],
			Host:                host,
			DispersalPort:       dispersalPort,
			RetrievalPort:       retrievalPort,
		}

		indexedOperators[id] = indexed
		privateOperators[id] = private
	}

	storedOperators := make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo, len(d.Stakes))
	totals := make(map[core.QuorumID]*core.OperatorInfo)

	for _, quorumID := range quorums {

		storedOperators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo, len(d.Stakes[quorumID]))

		index := uint(0)
		for _, opID := range d.Operators {
			stake, ok := d.Stakes[quorumID][opID]
			if !ok {
				continue
			}

			storedOperators[quorumID][opID] = &core.OperatorInfo{
				Stake: big.NewInt(int64(stake)),
				Index: index,
			}
			index++
		}

		quorumStake := 0
		for _, stake := range d.Stakes[quorumID] {
			quorumStake += stake
		}
		totals[quorumID] = &core.OperatorInfo{
			Stake: big.NewInt(int64(quorumStake)),
			Index: uint(len(d.Stakes[quorumID])),
		}
	}

	operatorState := &core.OperatorState{
		Operators:   storedOperators,
		Totals:      totals,
		BlockNumber: blockNumber,
	}

	filteredIndexedOperators := make(map[core.OperatorID]*core.IndexedOperatorInfo, 0)
	for quorumID, operatorsByID := range storedOperators {
		for opID := range operatorsByID {
			if aggPubKeys[quorumID] == nil {
				key := privateOperators[opID].KeyPair.GetPubKeyG1()
				aggPubKeys[quorumID] = key.Clone()
			} else {
				aggPubKeys[quorumID].Add(privateOperators[opID].KeyPair.GetPubKeyG1())
			}
			filteredIndexedOperators[opID] = indexedOperators[opID]
		}
	}

	indexedState := &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: filteredIndexedOperators,
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

func (d *ChainDataMock) GetOperatorSocket(ctx context.Context, blockNumber uint, operator core.OperatorID) (string, error) {

	state := d.GetTotalOperatorState(ctx, blockNumber)

	return state.IndexedOperatorState.IndexedOperators[operator].Socket, nil
}

func (d *ChainDataMock) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {

	state := d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, quorums)

	return state.IndexedOperatorState, nil

}

func (d *ChainDataMock) GetIndexedOperators(ctx context.Context, blockNumber uint) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {
	state := d.GetTotalOperatorState(ctx, blockNumber)

	return state.IndexedOperatorState.IndexedOperators, nil
}

func (d *ChainDataMock) GetCurrentBlockNumber() (uint, error) {
	args := d.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (d *ChainDataMock) Start(context.Context) error {
	return nil
}
