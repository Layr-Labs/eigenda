package validator

import (
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
)

// BlobDecoder is responsible for decoding blobs from chunk data.
type BlobDecoder interface {

	// DecodeBlob decodes a blob from the given chunk data.
	DecodeBlob(
		blobKey v2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error)
}

// BlobDecoderFactory is a function that creates a new BlobDecoder instance.
type BlobDecoderFactory func(
	verifier encoding.Verifier,
) BlobDecoder

var _ BlobDecoder = &blobDecoder{}

// blobDecoder is a standard implementation of the BlobDecoder interface.
type blobDecoder struct {
	verifier encoding.Verifier
}

var _ BlobDecoderFactory = NewBlobDecoder

// NewBlobDecoder creates a new BlobDecoder instance.
func NewBlobDecoder(verifier encoding.Verifier) BlobDecoder {
	return &blobDecoder{
		verifier: verifier,
	}
}

func (d *blobDecoder) DecodeBlob(
	_ v2.BlobKey, // used for unit tests
	chunks []*encoding.Frame,
	indices []uint,
	encodingParams *encoding.EncodingParams,
	blobCommitments *encoding.BlobCommitments,
) ([]byte, error) {

	return d.verifier.Decode(
		chunks,
		indices,
		*encodingParams,
		uint64(blobCommitments.Length)*encoding.BYTES_PER_SYMBOL)
}
