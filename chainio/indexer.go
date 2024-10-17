package chainio

import (
	"context"

	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
)

// IndexedOperatorInfo contains information about an operator which is contained in events from the EigenDA smart contracts. Note that
// this information does not depend on the quorum.
type IndexedOperatorInfo struct {
	// PubKeyG1 and PubKeyG2 are the public keys of the operator, which are retrieved from the EigenDAPubKeyCompendium smart contract
	PubkeyG1 *bn254.G1Point
	PubkeyG2 *bn254.G2Point
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
	AggKeys map[QuorumID]*bn254.G1Point
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
