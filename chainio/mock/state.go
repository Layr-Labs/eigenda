package mock

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/stretchr/testify/mock"
)

type ChainDataMock struct {
	mock.Mock

	KeyPairs  map[chainio.OperatorID]*bn254.KeyPair
	Operators []chainio.OperatorID
	Stakes    map[chainio.QuorumID]map[chainio.OperatorID]int
}

var _ chainio.ChainState = (*ChainDataMock)(nil)
var _ chainio.IndexedChainState = (*ChainDataMock)(nil)

type PrivateOperatorInfo struct {
	*chainio.IndexedOperatorInfo
	KeyPair       *bn254.KeyPair
	Host          string
	DispersalPort string
	RetrievalPort string
}

type PrivateOperatorState struct {
	*chainio.OperatorState
	*chainio.IndexedOperatorState
	PrivateOperators map[chainio.OperatorID]*PrivateOperatorInfo
}

func MakeOperatorId(id int) chainio.OperatorID {
	data := [32]byte{uint8(id)}
	return data
}

func NewChainDataMock(stakes map[chainio.QuorumID]map[chainio.OperatorID]int) (*ChainDataMock, error) {

	seenOperators := make(map[chainio.OperatorID]struct{})
	for _, oprStakes := range stakes {
		for opID := range oprStakes {
			if _, ok := seenOperators[opID]; ok {
				continue
			}
			seenOperators[opID] = struct{}{}
		}
	}

	operators := make([]chainio.OperatorID, 0, len(seenOperators))
	for opID := range seenOperators {
		operators = append(operators, opID)
	}

	sort.Slice(operators, func(i, j int) bool {
		return chainio.GetOperatorHex(operators[i]) < chainio.GetOperatorHex(operators[j])
	})

	keyPairs := make(map[chainio.OperatorID]*bn254.KeyPair)
	for _, opID := range operators {
		keyPair, err := bn254.GenRandomBlsKeys()
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
func MakeChainDataMock(numOperatorsPerQuorum map[chainio.QuorumID]int) (*ChainDataMock, error) {
	stakes := make(map[chainio.QuorumID]map[chainio.OperatorID]int)
	for quorumID, numOpr := range numOperatorsPerQuorum {
		stakes[quorumID] = make(map[chainio.OperatorID]int)
		for i := 0; i < numOpr; i++ {
			id := MakeOperatorId(i)
			stakes[quorumID][id] = int(i + 1)
		}
	}

	return NewChainDataMock(stakes)
}

func (d *ChainDataMock) GetTotalOperatorState(ctx context.Context, blockNumber uint) *PrivateOperatorState {
	return d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, []chainio.QuorumID{})
}

func (d *ChainDataMock) GetTotalOperatorStateWithQuorums(ctx context.Context, blockNumber uint, filterQuorums []chainio.QuorumID) *PrivateOperatorState {
	quorums := filterQuorums
	if len(quorums) == 0 {
		for quorumID := range d.Stakes {
			quorums = append(quorums, quorumID)
		}
	}

	indexedOperators := make(map[chainio.OperatorID]*chainio.IndexedOperatorInfo, len(d.Operators))
	privateOperators := make(map[chainio.OperatorID]*PrivateOperatorInfo, len(d.Operators))

	aggPubKeys := make(map[chainio.QuorumID]*bn254.G1Point)
	for i, id := range d.Operators {

		host := "0.0.0.0"
		dispersalPort := fmt.Sprintf("3%03v", 2*i)
		retrievalPort := fmt.Sprintf("3%03v", 2*i+1)
		socket := chainio.MakeOperatorSocket(host, dispersalPort, retrievalPort)

		indexed := &chainio.IndexedOperatorInfo{
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

	storedOperators := make(map[chainio.QuorumID]map[chainio.OperatorID]*chainio.OperatorInfo, len(d.Stakes))
	totals := make(map[chainio.QuorumID]*chainio.OperatorInfo)

	for _, quorumID := range quorums {

		storedOperators[quorumID] = make(map[chainio.OperatorID]*chainio.OperatorInfo, len(d.Stakes[quorumID]))

		index := uint32(0)
		for _, opID := range d.Operators {
			stake, ok := d.Stakes[quorumID][opID]
			if !ok {
				continue
			}

			storedOperators[quorumID][opID] = &chainio.OperatorInfo{
				Stake: big.NewInt(int64(stake)),
				Index: index,
			}
			index++
		}

		quorumStake := 0
		for _, stake := range d.Stakes[quorumID] {
			quorumStake += stake
		}
		totals[quorumID] = &chainio.OperatorInfo{
			Stake: big.NewInt(int64(quorumStake)),
			Index: uint32(len(d.Stakes[quorumID])),
		}
	}

	operatorState := &chainio.OperatorState{
		Operators:   storedOperators,
		Totals:      totals,
		BlockNumber: blockNumber,
	}

	filteredIndexedOperators := make(map[chainio.OperatorID]*chainio.IndexedOperatorInfo, 0)
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

	indexedState := &chainio.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: filteredIndexedOperators,
		AggKeys:          make(map[chainio.QuorumID]*bn254.G1Point),
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

func (d *ChainDataMock) GetOperatorState(ctx context.Context, blockNumber uint, quorums []chainio.QuorumID) (*chainio.OperatorState, error) {
	state := d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, quorums)

	return state.OperatorState, nil
}

func (d *ChainDataMock) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator chainio.OperatorID) (*chainio.OperatorState, error) {

	state := d.GetTotalOperatorState(ctx, blockNumber)

	return state.OperatorState, nil

}

func (d *ChainDataMock) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []chainio.QuorumID) (*chainio.IndexedOperatorState, error) {

	state := d.GetTotalOperatorStateWithQuorums(ctx, blockNumber, quorums)

	return state.IndexedOperatorState, nil

}

func (d *ChainDataMock) GetIndexedOperators(ctx context.Context, blockNumber uint) (map[chainio.OperatorID]*chainio.IndexedOperatorInfo, error) {
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
