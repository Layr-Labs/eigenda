package disperser

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, encodingParams core.EncodingParams) (*core.BlobCommitments, []*core.Chunk, error)
}
