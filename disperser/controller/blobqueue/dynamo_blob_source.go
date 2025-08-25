package blobqueue

import (
	"context"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

var _ BlobSource = (*dynamoBlobSource)(nil)

// Provides blobs that need to be put into batches for dispersal from DynamoDB.
type dynamoBlobSource struct {
}

// NewDynamoBlobSource creates a new BlobSource built on top of DynamoDB.
func NewDynamoBlobSource() (BlobSource, error) {
	return &dynamoBlobSource{}, nil
}

func (d *dynamoBlobSource) GetBlobsToDisperse(ctx context.Context) ([]*v2.BlobMetadata, error) {
	//TODO implement me
	panic("implement me")
}
