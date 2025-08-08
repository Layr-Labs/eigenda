package controller

import "github.com/Layr-Labs/eigenda/core"

// When creating a new batch, the controller needs two important pieces of metadata:
// 1. a reference block number (or RBN)
// 2. for each quorum number associated with the batch, the operator state for that quorum
//
// When it comes time to create a new batch, the primary goal of this cache is to avoid having to perform a
// blocking RPC call in order to retrieve the required metadata.
type BatchMetadataCache struct {
}

// NewBatchMetadataCache creates a new BatchMetadataCache instance.
func NewBatchMetadataCache() *BatchMetadataCache {
	return nil
}

func (c *BatchMetadataCache) GetBatchMetadata(quorums []core.QuorumID) (
	rbn uint64,
	operatorState *core.IndexedOperatorState,
	err error,
) {

	return 0, nil, nil

}
