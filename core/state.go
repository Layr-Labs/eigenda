package core

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"slices"
	"strings"
)

// Operators

// OperatorSocket is formatted as "host:dispersalPort;retrievalPort;v2DispersalPort;v2RetrievalPort"
type OperatorSocket struct {
	host            string
	v1DispersalPort string
	v1RetrievalPort string
	v2DispersalPort string
	v2RetrievalPort string
}

func (s *OperatorSocket) String() string {
	if s.v2DispersalPort == "" && s.v2RetrievalPort == "" {
		return fmt.Sprintf("%s:%s;%s", s.host, s.v1DispersalPort, s.v2RetrievalPort)
	}
	return fmt.Sprintf("%s:%s;%s;%s;%s", s.host, s.v1DispersalPort, s.v1RetrievalPort, s.v2DispersalPort, s.v2RetrievalPort)
}

func NewOperatorSocket(nodeIP, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort string) OperatorSocket {
	return OperatorSocket{
		host:            nodeIP,
		v1DispersalPort: dispersalPort,
		v1RetrievalPort: retrievalPort,
		v2DispersalPort: v2DispersalPort,
		v2RetrievalPort: v2RetrievalPort,
	}
}

type StakeAmount = *big.Int

func ParseOperatorSocket(opSocketStr string) (*OperatorSocket, error) {

	operatorSocket := OperatorSocket{}
	s := strings.Split(opSocketStr, ";")

	host, v1DispersalPort, err := net.SplitHostPort(s[0])
	if err != nil {
		return nil, fmt.Errorf("invalid host address format in %s: it must specify valid IP or host name (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
	}
	if _, err = net.LookupHost(host); err != nil {
		return nil, fmt.Errorf("invalid host address format in %s: it must specify valid IP or host name (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
	}
	if err = ValidatePort(v1DispersalPort); err != nil {
		return nil, fmt.Errorf("invalid v1 dispersal port format in %s: it must specify valid v1 dispersal port (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
	}
	operatorSocket.host = host
	operatorSocket.v1DispersalPort = v1DispersalPort

	switch len(s) {
	case 4:
		if err = ValidatePort(s[2]); err != nil {
			return nil, fmt.Errorf("invalid v2 dispersal port format in %s: it must specify valid v2 dispersal port (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
		}

		if err = ValidatePort(s[3]); err != nil {
			return nil, fmt.Errorf("invalid v2 retrieval port format in %s: it must specify valid v2 retrieval port (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
		}
		operatorSocket.v2DispersalPort = s[2]
		operatorSocket.v2RetrievalPort = s[3]
		fallthrough
	case 2:
		// V1 Parsing
		//parsedSocket.v1RetrievalPort = s[1]
		if err = ValidatePort(s[1]); err != nil {
			return nil, fmt.Errorf("invalid v1 retrieval port format in %s: it must specify valid v1 retrieval port (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
		}
		operatorSocket.v1DispersalPort = s[1]
		return &operatorSocket, nil
	default:
		return nil, fmt.Errorf("invalid opSocketStr address format %s: it must specify v1 dispersal/retrieval ports, or v2 dispersal/retrieval ports (ex. 0.0.0.0:32004;32005;32006;32007)", opSocketStr)
	}
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
	GetCurrentBlockNumber(ctx context.Context) (uint, error)
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
