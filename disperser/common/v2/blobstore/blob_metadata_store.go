package blobstore

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

type BlobMetadataStore interface {
	PutBlobMetadata(ctx context.Context, metadata *v2.BlobMetadata) error
	GetBlobMetadata(ctx context.Context, blobKey core.BlobKey) (*v2.BlobMetadata, error)
	GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus) ([]*v2.BlobMetadata, error)
	GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error)
}
