package internal

import (
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
	verifier *rs.Encoder,
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

	return d.encoder.Decode(
		frames,
		toUint64Array(indices),
		uint64(blobCommitments.Length)*encoding.BYTES_PER_SYMBOL,
		*encodingParams,
	)
}

// TODO(samlaf): this is dumb, we shouldn't have to allocate like this just to fit into functions signatures...
// we should standardize all uses of ChunkNumber to be uint64, not some places uint, others uint64.
func toUint64Array(chunkIndices []encoding.ChunkNumber) []uint64 {
	res := make([]uint64, len(chunkIndices))
	for i, d := range chunkIndices {
		res[i] = uint64(d)
	}
	return res
}
