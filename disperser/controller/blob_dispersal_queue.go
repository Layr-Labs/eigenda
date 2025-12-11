package controller

import (
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

// BlobDispersalQueue acquires and provides information about blobs that are ready for immediate dispersal. This object
// forms the controller's interface to the encoder->controller pipeline.
type BlobDispersalQueue interface {

	// GetBlobChannel returns a channel that yields blobs ready for dispersal.
	//
	// Due to some tech debt with the dynamoDB implementation, assume that this channel may return
	// the same blob multiple times, and that the caller is responsible for deduplicating them.
	GetBlobChannel() <-chan *v2.BlobMetadata
}
