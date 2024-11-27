package core

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"
)

// Operators

type OperatorSocket string

func (s OperatorSocket) String() string {
	return string(s)
}

func MakeOperatorSocket(nodeIP, dispersalPort, retrievalPort string) OperatorSocket {
	return OperatorSocket(fmt.Sprintf("%s:%s;%s", nodeIP, dispersalPort, retrievalPort))
}

type StakeAmount = *big.Int

func ParseOperatorSocket(socket string) (host string, dispersalPort string, retrievalPort string, err error) {
	s := strings.Split(socket, ";")
	if len(s) != 2 {
		err = fmt.Errorf("invalid socket address format, missing retrieval port: %s", socket)
		return
	}
	retrievalPort = s[1]

	s = strings.Split(s[0], ":")
	if len(s) != 2 {
		err = fmt.Errorf("invalid socket address format: %s", socket)
		return
	}
	host = s[0]
	dispersalPort = s[1]

	return
}

// OperatorInfo contains information about an operator which is stored on the blockchain state,
// corresponding to a particular quorum
type OperatorInfo struct {
	// Stake is the amount of stake held by the operator in the quorum
	Stake StakeAmount
	// Index is the index of the operator within the quorum
	Index OperatorIndex
}

// OperatorState contains information about the current state of operators which is stored in the blockchain state
type OperatorState struct {
	// Operators is a map from quorum ID to a map from the operators in that quourm to their StoredOperatorInfo. Membership
	// in the map implies membership in the quorum.
	Operators map[QuorumID]map[OperatorID]*OperatorInfo
	// Totals is a map from quorum ID to the total stake (Stake) and total count (Index) of all operators in that quorum
	Totals map[QuorumID]*OperatorInfo
	// BlockNumber is the block number at which this state was retrieved
	BlockNumber uint
}

func (s *OperatorState) Hash() (map[QuorumID][16]byte, error) {
	res := make(map[QuorumID][16]byte)
	type operatorInfoWithID struct {
		OperatorID string
		Stake      string
		Index      uint
	}
	for quorumID, opInfos := range s.Operators {
		marshalable := struct {
			Operators   []operatorInfoWithID
			Totals      OperatorInfo
			BlockNumber uint
		}{
			Operators:   make([]operatorInfoWithID, 0, len(opInfos)),
			Totals:      OperatorInfo{},
			BlockNumber: s.BlockNumber,
		}

		for opID, opInfo := range opInfos {
			marshalable.Operators = append(marshalable.Operators, operatorInfoWithID{
				OperatorID: opID.Hex(),
				Stake:      opInfo.Stake.String(),
				Index:      uint(opInfo.Index),
			})
		}
		slices.SortStableFunc(marshalable.Operators, func(a, b operatorInfoWithID) int {
			return strings.Compare(a.OperatorID, b.OperatorID)
		})

		marshalable.Totals = *s.Totals[quorumID]
		data, err := json.Marshal(marshalable)
		if err != nil {
			return nil, err
		}
		res[quorumID] = md5.Sum(data)
	}

	return res, nil
}

// IndexedOperatorInfo contains information about an operator which is contained in events from the EigenDA smart contracts. Note that
// this information does not depend on the quorum.
type IndexedOperatorInfo struct {
	// PubKeyG1 and PubKeyG2 are the public keys of the operator, which are retrieved from the EigenDAPubKeyCompendium smart contract
	PubkeyG1 *G1Point
	PubkeyG2 *G2Point
	// Socket is the socket address of the operator, in the form "host:port"
	Socket string
}

// IndexedOperatorState contains information about the current state of operators which is contained in events from the EigenDA smart contracts,
// in addition to the information contained in OperatorState
type IndexedOperatorState struct {
	*OperatorState
	// IndexedOperators is a map from operator ID to the IndexedOperatorInfo for that operator.
	IndexedOperators map[OperatorID]*IndexedOperatorInfo
	// AggKeys is a map from quorum ID to the aggregate public key of the operators in that quorum
	AggKeys map[QuorumID]*G1Point
}

// ChainState is an interface for getting information about the current chain state.
type ChainState interface {
	GetCurrentBlockNumber() (uint, error)
	GetOperatorState(ctx context.Context, blockNumber uint, quorums []QuorumID) (*OperatorState, error)
	GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator OperatorID) (*OperatorState, error)
	GetOperatorSocket(ctx context.Context, blockNumber uint, operator OperatorID) (string, error)
}

// ChainState is an interface for getting information about the current chain state.
type IndexedChainState interface {
	ChainState
	// GetIndexedOperatorState returns the IndexedOperatorState for the given block number and quorums
	// If the quorum is not found, the quorum will be ignored and the IndexedOperatorState will be returned for the remaining quorums
	GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []QuorumID) (*IndexedOperatorState, error)
	GetIndexedOperators(ctx context.Context, blockNumber uint) (map[OperatorID]*IndexedOperatorInfo, error)
	Start(context context.Context) error
}
