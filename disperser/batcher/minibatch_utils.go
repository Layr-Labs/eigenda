package batcher

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
)

// DispersedBlobs returns the blobs that were successfully dispersed and the blobs that failed to be dispersed.
// Really simple implementation that only considers minibatches successful if all operators have responded.
func DispersedBlobs(operators map[core.OperatorID]*core.IndexedOperatorInfo, dispersals []*MinibatchDispersal, blobMinibatchMappings []*BlobMinibatchMapping) (successful []*BlobMinibatchMapping, failed []*BlobMinibatchMapping, err error) {
	if len(dispersals) == 0 {
		return nil, nil, fmt.Errorf("no dispersals")
	}
	if len(blobMinibatchMappings) == 0 {
		return nil, nil, fmt.Errorf("no blob minibatch mappings")
	}

	received := make(map[uint]map[core.OperatorID]struct{})
	failedMinibatches := make(map[uint]struct{})
	for _, d := range dispersals {
		// If the response is pending or errored, we consider the minibatch failed
		if d.RespondedAt.IsZero() || d.Error != nil {
			failedMinibatches[d.MinibatchIndex] = struct{}{}
			continue
		}
		if _, ok := failedMinibatches[d.MinibatchIndex]; ok {
			continue
		}

		if _, ok := received[d.MinibatchIndex]; !ok {
			received[d.MinibatchIndex] = make(map[core.OperatorID]struct{})
		}
		received[d.MinibatchIndex][d.OperatorID] = struct{}{}
	}

	for minibatchIndex, opIDs := range received {
		for opID := range operators {
			if _, ok := opIDs[opID]; !ok {
				failedMinibatches[minibatchIndex] = struct{}{}
				break
			}
		}
	}

	batchID := blobMinibatchMappings[0].BatchID
	successful = make([]*BlobMinibatchMapping, 0)
	failed = make([]*BlobMinibatchMapping, 0)
	for _, blobMinibatchMapping := range blobMinibatchMappings {
		if blobMinibatchMapping.BatchID != batchID {
			return nil, nil, fmt.Errorf("batch ID mismatch")
		}
		if _, ok := failedMinibatches[blobMinibatchMapping.MinibatchIndex]; ok {
			failed = append(failed, blobMinibatchMapping)
		} else {
			successful = append(successful, blobMinibatchMapping)
		}
	}

	return successful, failed, nil
}
