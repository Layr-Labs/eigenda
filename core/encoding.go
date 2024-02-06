package core

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Commitments

// Commitment is a polynomial commitment (e.g. a kzg commitment)
type G1Commitment bn254.G1Point

// Commitment is a polynomial commitment (e.g. a kzg commitment)
type G2Commitment bn254.G2Point

// LengthProof is a polynomial commitment on G2 (e.g. a kzg commitment) used for low degree proof
type LengthProof = G2Commitment

// The proof used to open a commitment. In the case of Kzg, this is also a kzg commitment, and is different from a Commitment only semantically.
type Proof = bn254.G1Point

// Symbol is a symbol in the field used for polynomial commitments
type Symbol = bn254.Fr

// Encoding

// EncodingParams contains the encoding parameters that the encoder must satisfy.
type EncodingParams struct {
	ChunkLength uint // ChunkSize is the length of the chunk in symbols
	NumChunks   uint
}

// Encoder is responsible for encoding, decoding, and chunk verification
type Encoder interface {
	// Encode takes in a blob and returns the commitments and encoded chunks. The encoding will satisfy the property that
	// for any number M such that M*params.ChunkLength > BlobCommitments.Length, then any set of M chunks will be sufficient to
	// reconstruct the blob.
	Encode(data []byte, params EncodingParams) (BlobCommitments, []*Chunk, error)

	// VerifyChunks takes in the chunks, indices, commitments, and encoding parameters and returns an error if the chunks are invalid.
	VerifyChunks(chunks []*Chunk, indices []ChunkNumber, commitments BlobCommitments, params EncodingParams) error

	// VerifyBatch takes in the encoding parameters, samples and the number of blobs and returns an error if a chunk in any sample is invalid.
	UniversalVerifySubBatch(params EncodingParams, samples []Sample, numBlobs int) error

	// VerifyBlobLength takes in the commitments and returns an error if the blob length is invalid.
	VerifyBlobLength(commitments BlobCommitments) error

	// VerifyCommitEquivalence takes in a list of commitments and returns an error if the commitment of G1 and G2 are inconsistent
	VerifyCommitEquivalenceBatch(commitments []BlobCommitments) error

	// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
	Decode(chunks []*Chunk, indices []ChunkNumber, params EncodingParams, inputSize uint64) ([]byte, error)
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	symSize := uint(bn254.BYTES_PER_COEFFICIENT)
	return (blobSize + symSize - 1) / symSize
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * bn254.BYTES_PER_COEFFICIENT
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint, quorumThreshold, advThreshold uint8) uint {
	return roundUpDivide(blobLength*100, uint(quorumThreshold)-uint(advThreshold))
}

// GetEncodingParams takes in the minimum chunk length and the minimum number of chunks and returns the encoding parameters.
// Both the ChunkLength and NumChunks must be powers of 2, and the ChunkLength returned here should be used in constructing the BlobHeader.
func GetEncodingParams(minChunkLength, minNumChunks uint) (EncodingParams, error) {
	return EncodingParams{
		ChunkLength: uint(encoder.NextPowerOf2(uint64(minChunkLength))),
		NumChunks:   uint(encoder.NextPowerOf2(uint64(minNumChunks))),
	}, nil
}

// ValidateEncodingParams takes in the encoding parameters and returns an error if they are invalid.
func ValidateEncodingParams(params EncodingParams, blobLength, SRSOrder int) error {

	if int(params.ChunkLength*params.NumChunks) > SRSOrder {
		return fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLength, params.NumChunks, SRSOrder)
	}

	if int(params.ChunkLength*params.NumChunks) < blobLength {
		return fmt.Errorf("the supplied encoding parameters are not sufficient for the size of the data input")
	}

	return nil

}
