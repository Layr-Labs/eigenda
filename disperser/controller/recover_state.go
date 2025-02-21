package controller

import (
	"context"
	"fmt"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// RecoverState checks for blobs in the GatheringSignatures state and updates their status to Failed.
func RecoverState(
	ctx context.Context,
	blobStore *blobstore.BlobMetadataStore,
	logger logging.Logger,
) error {
	logger.Info("recovering state...")

	metadata, err := blobStore.GetBlobMetadataByStatus(ctx, v2.GatheringSignatures, 0)
	if err != nil {
		return fmt.Errorf("failed to get blobs in gathering signatures state: %w", err)
	}

	if len(metadata) == 0 {
		logger.Info("no blobs in gathering signatures state")
		return nil
	}

	logger.Info("found blobs in gathering signatures state", "count", len(metadata))

	for _, blob := range metadata {
		key, err := blob.BlobHeader.BlobKey()
		if err != nil {
			logger.Error("failed to get blob key", "err", err)
			continue
		}

		logger.Debug("updating blob status", "key", key, "status", v2.Failed)
		if err := blobStore.UpdateBlobStatus(ctx, key, v2.Failed); err != nil {
			logger.Error("failed to update blob status", "blobKey", key.Hex(), "err", err)
		}
	}
	logger.Info("recovered state successfully")
	return nil
}
