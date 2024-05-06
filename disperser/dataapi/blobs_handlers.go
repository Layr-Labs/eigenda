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

func (s *server) getBlobCountByAccountId(ctx context.Context, accountID string) (*BlobCountForAccountIdResponse, error) {
	s.logger.Info("Calling get blob", "AccountId", accountID)

	metadataCount, err := s.blobstore.GetBlobMetadataCountByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Got blob metadata count for AccountId", "AccountId", accountID, "metadataCount", metadataCount)
	return &BlobCountForAccountIdResponse{
		Count:     metadataCount,
		AccountId: accountID,
	}, nil
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
