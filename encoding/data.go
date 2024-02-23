package encoding

import "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"

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

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment       *G1Commitment `json:"commitment"`
	LengthCommitment *G2Commitment `json:"length_commitment"`
	LengthProof      *LengthProof  `json:"length_proof"`
	Length           uint          `json:"length"`
}

// Frame is a chunk of data with the associated multi-reveal proof
type Frame struct {
	// Proof is the multireveal proof corresponding to the chunk
	Proof Proof
	// Coeffs contains the coefficience of the interpolating polynomial of the chunk
	Coeffs []Symbol
}

func (f *Frame) Length() int {
	return len(f.Coeffs)
}

// Returns the size of chunk in bytes.
func (f *Frame) Size() int {
	return f.Length() * bn254.BYTES_PER_COEFFICIENT
}

// Sample is a chunk with associated metadata used by the Universal Batch Verifier
type Sample struct {
	Commitment      *G1Commitment
	Chunk           *Frame
	AssignmentIndex ChunkNumber
	BlobIndex       int
}

// SubBatch is a part of the whole Batch with identical Encoding Parameters, i.e. (ChunkLength, NumChunk)
// Blobs with the same encoding parameters are collected in a single subBatch
type SubBatch struct {
	Samples  []Sample
	NumBlobs int
}

type ChunkNumber = uint
