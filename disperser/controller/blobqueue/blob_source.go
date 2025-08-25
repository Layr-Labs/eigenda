package blobqueue

import (
	"context"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

// Provides blobs that need to be put into batches for dispersal.
type BlobSource interface {

	// GetBlobsToDisperse returns a list of blobs that need to be put into batches for dispersal.
	// Some implementations (i.e. dynamo) may return duplicates depending on timing.
	GetBlobsToDisperse(ctx context.Context) ([]*v2.BlobMetadata, error)
}
