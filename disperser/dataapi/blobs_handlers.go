package dataapi

import (
	"context"
	"encoding/hex"
	"sort"

	"github.com/Layr-Labs/eigenda/disperser"
)

func (s *server) getBlob(ctx context.Context, key string) (*BlobMetadataResponse, error) {
	s.logger.Info("Calling get blob", "key", key)
	blobKey, err := disperser.ParseBlobKey(string(key))
	if err != nil {
		return nil, err
	}
	metadata, err := s.blobstore.GetBlobMetadata(ctx, blobKey)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Got blob metadata", "metadata", metadata)
	return convertMetadataToBlobMetadataResponse(metadata)
}

func (s *server) getBlobs(ctx context.Context, limit int) ([]*BlobMetadataResponse, error) {
	_, blobMetadatas, err := s.getBlobMetadataByBatchesWithLimit(ctx, limit)
	if err != nil {
		return nil, err
	}
	if len(blobMetadatas) == 0 {
		return nil, errNotFound
	}

	return s.convertBlobMetadatasToBlobMetadataResponse(ctx, blobMetadatas)
}

func (s *server) getBlobsFromBatchHeaderHash(ctx context.Context, batcherHeaderHash [32]byte, limit int, exclusiveStartKey *disperser.BatchIndexExclusiveStartKey) ([]*BlobMetadataResponse, *disperser.BatchIndexExclusiveStartKey, error) {
	blobMetadatas, newExclusiveStartKey, err := s.getBlobMetadataByBatchHeaderHashWithLimit(ctx, batcherHeaderHash, int32(limit), exclusiveStartKey)
	if err != nil {
		return nil, nil, err
	}
	if len(blobMetadatas) == 0 {
		return nil, nil, errNotFound
	}

	responses, err := s.convertBlobMetadatasToBlobMetadataResponse(ctx, blobMetadatas)
	if err != nil {
		return nil, nil, err
	}

	return responses, newExclusiveStartKey, nil
}

func (s *server) convertBlobMetadatasToBlobMetadataResponse(ctx context.Context, metadatas []*disperser.BlobMetadata) ([]*BlobMetadataResponse, error) {
	var (
		err               error
		responseMetadatas = make([]*BlobMetadataResponse, len(metadatas))
	)

	sort.SliceStable(metadatas, func(i, j int) bool {
		// We may have unconfirmed blobs to fetch, which will not have the ConfirmationInfo.
		// In such case, we order them by request timestamp.
		if metadatas[i].ConfirmationInfo == nil || metadatas[j].ConfirmationInfo == nil {
			return metadatas[i].RequestMetadata.RequestedAt < metadatas[j].RequestMetadata.RequestedAt
		}
		if metadatas[i].ConfirmationInfo.BatchID != metadatas[j].ConfirmationInfo.BatchID {
			return metadatas[i].ConfirmationInfo.BatchID < metadatas[j].ConfirmationInfo.BatchID
		}
		return metadatas[i].ConfirmationInfo.BlobIndex < metadatas[j].ConfirmationInfo.BlobIndex
	})

	for i := range metadatas {
		responseMetadatas[i], err = convertMetadataToBlobMetadataResponse(metadatas[i])
		if err != nil {
			return nil, err
		}
	}

	return responseMetadatas, nil
}

func convertMetadataToBlobMetadataResponse(metadata *disperser.BlobMetadata) (*BlobMetadataResponse, error) {
	// If the blob is not confirmed or finalized, return the metadata without the confirmation info
	isConfirmed, err := metadata.IsConfirmed()
	if err != nil {
		return nil, err
	}
	if !isConfirmed {
		return &BlobMetadataResponse{
			BlobKey:        metadata.GetBlobKey().String(),
			SecurityParams: metadata.RequestMetadata.SecurityParams,
			RequestAt:      ConvertNanosecondToSecond(metadata.RequestMetadata.RequestedAt),
			BlobStatus:     metadata.BlobStatus,
		}, nil
	}

	return &BlobMetadataResponse{
		BlobKey:                 metadata.GetBlobKey().String(),
		BatchHeaderHash:         hex.EncodeToString(metadata.ConfirmationInfo.BatchHeaderHash[:]),
		BlobIndex:               metadata.ConfirmationInfo.BlobIndex,
		SignatoryRecordHash:     hex.EncodeToString(metadata.ConfirmationInfo.SignatoryRecordHash[:]),
		ReferenceBlockNumber:    metadata.ConfirmationInfo.ReferenceBlockNumber,
		BatchRoot:               hex.EncodeToString(metadata.ConfirmationInfo.BatchRoot),
		BlobInclusionProof:      hex.EncodeToString(metadata.ConfirmationInfo.BlobInclusionProof),
		BlobCommitment:          metadata.ConfirmationInfo.BlobCommitment,
		BatchId:                 metadata.ConfirmationInfo.BatchID,
		ConfirmationBlockNumber: metadata.ConfirmationInfo.ConfirmationBlockNumber,
		ConfirmationTxnHash:     metadata.ConfirmationInfo.ConfirmationTxnHash.String(),
		Fee:                     hex.EncodeToString(metadata.ConfirmationInfo.Fee),
		SecurityParams:          metadata.RequestMetadata.SecurityParams,
		RequestAt:               ConvertNanosecondToSecond(metadata.RequestMetadata.RequestedAt),
		BlobStatus:              metadata.BlobStatus,
	}, nil
}

func (s *server) getBlobMetadataByBatchesWithLimit(ctx context.Context, limit int) ([]*Batch, []*disperser.BlobMetadata, error) {
	var (
		blobMetadatas   = make([]*disperser.BlobMetadata, 0)
		batches         = make([]*Batch, 0)
		blobKeyPresence = make(map[string]struct{})
		batchPresence   = make(map[string]struct{})
	)

	for skip := 0; len(blobMetadatas) < limit && skip < limit; skip += maxQueryBatchesLimit {
		batchesWithLimit, err := s.subgraphClient.QueryBatchesWithLimit(ctx, maxQueryBatchesLimit, skip)
		if err != nil {
			s.logger.Error("Failed to query batches", "error", err)
			return nil, nil, err
		}

		if len(batchesWithLimit) == 0 {
			break
		}

		for i := range batchesWithLimit {
			s.logger.Debug("Getting blob metadata", "batchHeaderHash", batchesWithLimit[i].BatchHeaderHash)
			var (
				batch = batchesWithLimit[i]
			)
			if batch == nil {
				continue
			}
			batchHeaderHash, err := ConvertHexadecimalToBytes(batch.BatchHeaderHash)
			if err != nil {
				s.logger.Error("Failed to convert batch header hash to hex string", "error", err)
				continue
			}
			batchKey := string(batchHeaderHash[:])
			if _, found := batchPresence[batchKey]; !found {
				batchPresence[batchKey] = struct{}{}
			} else {
				// The batch has processed, skip it.
				s.logger.Error("Getting duplicate batch from the graph", "batch header hash", batchKey)
				continue
			}

			metadatas, err := s.blobstore.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
			if err != nil {
				s.logger.Error("Failed to get blob metadata", "error", err)
				continue
			}
			for _, bm := range metadatas {
				blobKey := bm.GetBlobKey().String()
				if _, found := blobKeyPresence[blobKey]; !found {
					blobKeyPresence[blobKey] = struct{}{}
					blobMetadatas = append(blobMetadatas, bm)
				} else {
					s.logger.Error("Getting duplicate blob key from the blobstore", "blobkey", blobKey)
				}
			}
			batches = append(batches, batch)
			if len(blobMetadatas) >= limit {
				break
			}
		}
	}

	if len(blobMetadatas) >= limit {
		blobMetadatas = blobMetadatas[:limit]
	}

	return batches, blobMetadatas, nil
}

func (s *server) getBlobMetadataByBatchHeaderHashWithLimit(ctx context.Context, batchHeaderHash [32]byte, limit int32, exclusiveStartKey *disperser.BatchIndexExclusiveStartKey) ([]*disperser.BlobMetadata, *disperser.BatchIndexExclusiveStartKey, error) {
	var allMetadata []*disperser.BlobMetadata
	var nextKey *disperser.BatchIndexExclusiveStartKey = exclusiveStartKey

	const maxLimit int32 = 1000
	remainingLimit := min(limit, maxLimit)

	s.logger.Debug("Getting blob metadata by batch header hash", "batchHeaderHash", batchHeaderHash, "remainingLimit", remainingLimit, "nextKey", nextKey)
	for int32(len(allMetadata)) < remainingLimit {
		metadatas, newNextKey, err := s.blobstore.GetAllBlobMetadataByBatchWithPagination(ctx, batchHeaderHash, remainingLimit-int32(len(allMetadata)), nextKey)
		if err != nil {
			s.logger.Error("Failed to get blob metadata", "error", err)
			return nil, nil, err
		}

		allMetadata = append(allMetadata, metadatas...)

		if newNextKey == nil {
			// No more data to fetch
			return allMetadata, nil, nil
		}

		nextKey = newNextKey

		if int32(len(allMetadata)) == remainingLimit {
			// We've reached the limit
			break
		}
	}

	return allMetadata, nextKey, nil
}
