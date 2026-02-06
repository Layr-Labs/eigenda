package types

import (
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
)

// Operator represents an EigenDA operator with their registration information.
type Operator struct {
	// Unique operator identifier
	ID core.OperatorID

	// Ethereum address of the operator
	Address common.Address

	// BLS public key (G1 point)
	BLSPubKeyG1 *core.G1Point

	// BLS public key (G2 point)
	BLSPubKeyG2 *core.G2Point

	// Socket address for connecting to the operator
	Socket string

	// Block number when the operator was registered
	RegisteredAtBlockNumber uint64

	// Block number when the operator was deregistered (nil if still registered)
	DeregisteredAtBlockNumber *uint64

	// List of quorum IDs the operator is part of
	QuorumIDs []core.QuorumID

	// Transaction hash of the registration
	RegisteredTxHash common.Hash

	// Transaction hash of the deregistration (nil if still registered)
	DeregisteredTxHash *common.Hash
}

// IsRegistered returns true if the operator is currently registered.
func (o *Operator) IsRegistered() bool {
	return o.DeregisteredAtBlockNumber == nil
}

// QuorumAPK represents the aggregate public key for a quorum at a specific block.
type QuorumAPK struct {
	// Quorum identifier
	QuorumID core.QuorumID

	// Block number for this snapshot
	BlockNumber uint64

	// Aggregate public key for all operators in the quorum
	APK *core.G1Point

	// Total stake in the quorum
	TotalStake *big.Int

	// Timestamp when this was recorded
	UpdatedAt time.Time
}

// OperatorSocketUpdate represents a socket address update for an operator.
type OperatorSocketUpdate struct {
	// Operator identifier
	OperatorID core.OperatorID

	// New socket address
	Socket string

	// Block number of the update
	BlockNumber uint64

	// Transaction hash of the update
	TxHash common.Hash

	// Timestamp when this was recorded
	UpdatedAt time.Time
}

// OperatorEjection represents an operator ejection event.
type OperatorEjection struct {
	// Operator identifier
	OperatorID core.OperatorID

	// Quorums from which the operator was ejected
	QuorumIDs []core.QuorumID

	// Block number of the ejection
	BlockNumber uint64

	// Transaction hash of the ejection
	TxHash common.Hash

	// Timestamp when this was recorded
	EjectedAt time.Time
}

// OperatorFilter is used to filter operators when querying.
type OperatorFilter struct {
	// Only return registered operators
	RegisteredOnly bool

	// Only return deregistered operators
	DeregisteredOnly bool

	// Filter by specific quorum ID (nil for all)
	QuorumID *core.QuorumID

	// Minimum block number (inclusive)
	MinBlock uint64

	// Maximum block number (inclusive, 0 for no limit)
	MaxBlock uint64
}

// QuorumAPKFilter is used to filter quorum APK snapshots when querying.
type QuorumAPKFilter struct {
	// Quorum identifier
	QuorumID core.QuorumID

	// Block number for the snapshot (0 for latest)
	BlockNumber uint64

	// Get all snapshots after this block (inclusive)
	MinBlock uint64

	// Get all snapshots before this block (inclusive, 0 for no limit)
	MaxBlock uint64
}
