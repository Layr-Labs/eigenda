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

// Proof is the multireveal proof
// Coeffs is identical to input data converted into Fr element
type Frame struct {
	Proof  Proof
	Coeffs []Symbol
}

// Sample is a chunk with associated metadata used by the Universal Batch Verifier
type Sample struct {
	Commitment      *G1Commitment
	Chunk           *Frame
	AssignmentIndex ChunkNumber
	BlobIndex       int
}

// SubBatch is a part of the whole Batch with identical Encoding Parameters, i.e. (ChunkLen, NumChunk)
// Blobs with the same encoding parameters are collected in a single subBatch
type SubBatch struct {
	Samples  []Sample
	NumBlobs int
}

type ChunkNumber = uint
