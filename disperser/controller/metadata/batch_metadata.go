package metadata

import "github.com/Layr-Labs/eigenda/core"

// The metadata required to create a new batch.
type BatchMetadata struct {
	// The eth block number associated with the batch.
	referenceBlockNumber uint64

	// The operator state for the specified block number.
	operatorState *core.IndexedOperatorState
}

// Create a new BatchMetadata instance with the specified reference block number and operator state.
func NewBatchMetadata(
	referenceBlockNumber uint64,
	operatorState *core.IndexedOperatorState,
) *BatchMetadata {
	return &BatchMetadata{
		referenceBlockNumber: referenceBlockNumber,
		operatorState:        operatorState,
	}
}

// Get the reference block number (RBN) for this batch metadata.
func (b *BatchMetadata) ReferenceBlockNumber() uint64 {
	return b.referenceBlockNumber
}

// Get the operator state for this batch metadata.
func (b *BatchMetadata) OperatorState() *core.IndexedOperatorState {
	return b.operatorState
}
