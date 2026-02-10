package store

import (
	"context"

	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
)

// Store is the interface for persisting and querying indexed chain state.
// This abstraction allows for different storage backends (memory, database, etc.).
type Store interface {
	// Operator operations

	// SaveOperator saves or updates an operator in the store.
	SaveOperator(ctx context.Context, op *types.Operator) error

	// GetOperator retrieves an operator by ID.
	GetOperator(ctx context.Context, id core.OperatorID) (*types.Operator, error)

	// ListOperators retrieves a paginated list of operators matching the filter.
	ListOperators(ctx context.Context, filter types.OperatorFilter, limit, offset int) ([]*types.Operator, error)

	// UpdateOperatorSocket updates the socket address for an operator.
	UpdateOperatorSocket(ctx context.Context, id core.OperatorID, socket string, blockNum uint64) error

	// DeregisterOperator marks an operator as deregistered.
	DeregisterOperator(ctx context.Context, id core.OperatorID, blockNum uint64, txHash common.Hash) error

	// Quorum APK operations

	// SaveQuorumAPK saves a quorum aggregate public key snapshot.
	SaveQuorumAPK(ctx context.Context, apk *types.QuorumAPK) error

	// GetQuorumAPK retrieves the aggregate public key for a quorum at a specific block.
	GetQuorumAPK(ctx context.Context, quorumID core.QuorumID, blockNum uint64) (*types.QuorumAPK, error)

	// ListQuorumAPKs retrieves quorum APK snapshots matching the filter.
	ListQuorumAPKs(ctx context.Context, filter types.QuorumAPKFilter) ([]*types.QuorumAPK, error)

	// Ejection operations

	// SaveEjection records an operator ejection event.
	SaveEjection(ctx context.Context, ejection *types.OperatorEjection) error

	// ListEjections retrieves ejection events, optionally filtered by operator.
	ListEjections(ctx context.Context, operatorID *core.OperatorID, limit, offset int) ([]*types.OperatorEjection, error)

	// Socket update operations

	// SaveSocketUpdate records an operator socket update event.
	SaveSocketUpdate(ctx context.Context, update *types.OperatorSocketUpdate) error

	// ListSocketUpdates retrieves socket update events for an operator.
	ListSocketUpdates(ctx context.Context, operatorID core.OperatorID, limit, offset int) ([]*types.OperatorSocketUpdate, error)

	// Block tracking

	// GetLastIndexedBlock returns the last block number that was indexed.
	GetLastIndexedBlock(ctx context.Context) (uint64, error)

	// SetLastIndexedBlock updates the last indexed block number.
	SetLastIndexedBlock(ctx context.Context, blockNum uint64) error

	// Persistence operations

	// Snapshot serializes the entire store state to bytes for persistence.
	Snapshot(ctx context.Context) ([]byte, error)

	// Restore deserializes and loads store state from bytes.
	Restore(ctx context.Context, data []byte) error
}
