package disperser

import (
	"context"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
)

type EncoderClientV2 interface {
	EncodeBlob(ctx context.Context, blobKey v2.BlobKey, encodingParams encoding.EncodingParams) (*encoding.FragmentInfo, error)
}
