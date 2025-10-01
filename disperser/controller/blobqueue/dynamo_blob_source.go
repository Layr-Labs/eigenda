package blobqueue

import (
	"context"
	"fmt"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
)

var _ BlobSource = (*dynamoBlobSource)(nil)

// Provides blobs that need to be put into batches for dispersal from DynamoDB.
type dynamoBlobSource struct {

	// used to iterate through blobs in Dynamo
	cursor *blobstore.StatusIndexCursor

	// a handle for interacting with dynamoDB
	metadataStore blobstore.MetadataStore

	// The maximum number of blobs to return in a single dynamo query.
	maxQuerySize uint64
}

// NewDynamoBlobSource creates a new BlobSource built on top of DynamoDB.
func NewDynamoBlobSource(
	metadataStore blobstore.MetadataStore,
	maxQuerySize uint64,
) BlobSource {

	if maxQuerySize == 0 {
		maxQuerySize = 1
	}

	return &dynamoBlobSource{
		metadataStore: metadataStore,
		maxQuerySize:  maxQuerySize,
	}
}

func (d *dynamoBlobSource) GetBlobsToDisperse(ctx context.Context) ([]*v2.BlobMetadata, error) {
	blobMetadata, cursor, err := d.metadataStore.GetBlobMetadataByStatusPaginated(
		ctx,
		v2.Encoded,
		d.cursor,
		int32(d.maxQuerySize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata: %w", err)
	}

	d.cursor = cursor
	return blobMetadata, nil
}
