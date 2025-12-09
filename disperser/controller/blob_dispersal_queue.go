package controller

import (
	"context"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

// BlobDispersalQueue acquires and provides information about blobs that are ready for immediate dispersal. This object
// forms the controller's interface to the encoder->controller pipeline.
type BlobDispersalQueue interface {

	// GetNextBlobForDispersal returns a channel that yields blobs ready for dispersal.
	GetNextBlobForDispersal(ctx context.Context) <-chan *v2.BlobMetadata
}
