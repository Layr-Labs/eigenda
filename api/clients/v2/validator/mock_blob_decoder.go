package validator

import (
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
)

var _ BlobDecoder = &MockBlobDecoder{}

// MockBlobDecoder is a mock implementation of the BlobDecoder interface.
type MockBlobDecoder struct {
	// A lambda function to be called when DecodeBlob is called.
	DecodeBlobFunction func(
		blobKey corev2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error)
}

func (m MockBlobDecoder) DecodeBlob(
	blobKey corev2.BlobKey,
	chunks []*encoding.Frame,
	indices []uint,
	encodingParams *encoding.EncodingParams,
	blobCommitments *encoding.BlobCommitments,
) ([]byte, error) {
	if m.DecodeBlobFunction == nil {
		return nil, nil
	}
	return m.DecodeBlobFunction(blobKey, chunks, indices, encodingParams, blobCommitments)
}

// NewMockBlobDecoderFactory creates a new BlobDecoderFactory that returns the provided decoder.
func NewMockBlobDecoderFactory(decoder BlobDecoder) BlobDecoderFactory {
	return func(verifier encoding.Verifier) BlobDecoder {
		return decoder
	}
}
