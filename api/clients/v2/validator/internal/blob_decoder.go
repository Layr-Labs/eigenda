package internal

import (
	"fmt"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

// BlobDecoder is responsible for decoding blobs from chunk data.
type BlobDecoder interface {

	// DecodeBlob decodes a blob from the given chunk data.
	DecodeBlob(
		blobKey v2.BlobKey,
		chunks []*encoding.Frame,
		indices []encoding.ChunkNumber,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error)
}

// BlobDecoderFactory is a function that creates a new BlobDecoder instance.
type BlobDecoderFactory func(
	encoder *rs.Encoder,
) BlobDecoder

var _ BlobDecoder = &blobDecoder{}

// blobDecoder is a standard implementation of the BlobDecoder interface.
type blobDecoder struct {
	encoder *rs.Encoder
}

var _ BlobDecoderFactory = NewBlobDecoder

// NewBlobDecoder creates a new BlobDecoder instance.
func NewBlobDecoder(encoder *rs.Encoder) BlobDecoder {
	return &blobDecoder{
		encoder: encoder,
	}
}

func (d *blobDecoder) DecodeBlob(
	_ v2.BlobKey, // used for unit tests
	chunks []*encoding.Frame,
	indices []encoding.ChunkNumber,
	encodingParams *encoding.EncodingParams,
	blobCommitments *encoding.BlobCommitments,
) ([]byte, error) {
	frames := make([]rs.FrameCoeffs, len(chunks))
	for i := range chunks {
		frames[i] = chunks[i].Coeffs
	}

	blob, err := d.encoder.Decode(
		frames,
		indices,
		uint64(blobCommitments.Length)*encoding.BYTES_PER_SYMBOL,
		*encodingParams,
	)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return blob, nil
}
