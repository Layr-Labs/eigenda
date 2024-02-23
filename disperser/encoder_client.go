package disperser

import (
	"context"

	"github.com/Layr-Labs/eigenda/encoding"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, []*encoding.Frame, error)
}
