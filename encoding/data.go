package encoding

import (
	"bytes"
	"fmt"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Commitment is a polynomial commitment (e.g. a kzg commitment)
type G1Commitment bn254.G1Affine

// Commitment is a polynomial commitment (e.g. a kzg commitment)
type G2Commitment bn254.G2Affine

// LengthProof is a polynomial commitment on G2 (e.g. a kzg commitment) used for low degree proof
type LengthProof = G2Commitment

// Proof is used to open a commitment. In the case of Kzg, this is also a kzg commitment, and is different from a Commitment only semantically.
type Proof = bn254.G1Affine

// Symbol is a symbol in the field used for polynomial commitments
type Symbol = fr.Element

// BlobCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment       *G1Commitment `json:"commitment"`
	LengthCommitment *G2Commitment `json:"length_commitment"`
	LengthProof      *LengthProof  `json:"length_proof"`
	// this is the length in SYMBOLS (32 byte field elements) of the blob. it must be a power of 2
	Length uint `json:"length"`
}

// ToProfobuf converts the BlobCommitments to protobuf format
func (c *BlobCommitments) ToProtobuf() (*pbcommon.BlobCommitment, error) {
	commitData, err := c.Commitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthCommitData, err := c.LengthCommitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthProofData, err := c.LengthProof.Serialize()
	if err != nil {
		return nil, err
	}

	return &pbcommon.BlobCommitment{
		Commitment:       commitData,
		LengthCommitment: lengthCommitData,
		LengthProof:      lengthProofData,
		Length:           uint32(c.Length),
	}, nil
}

// Equal checks if two BlobCommitments are equal, and returns an error if not.
// TODO(samlaf): should return structured errors to diffentiate 400 from 500 errors
// Any error returned here is currently returned as a 400 to users, but failing to Serialize a commitment
// should return a 500.
func (c *BlobCommitments) Equal(c1 *BlobCommitments) error {
	if c.Length != c1.Length {
		return fmt.Errorf("lengths are different: %d vs %d", c.Length, c1.Length)
	}

	cCommitment, err := c.Commitment.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize commitment: %w", err)
	}
	c1Commitment, err := c1.Commitment.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize c1 commitment: %w", err)
	}
	if !bytes.Equal(cCommitment, c1Commitment) {
		return fmt.Errorf("commitments are different: %v vs %v", cCommitment, c1Commitment)
	}

	cLengthCommitment, err := c.LengthCommitment.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize length commitment: %w", err)
	}
	c1LengthCommitment, err := c1.LengthCommitment.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize c1 length commitment: %w", err)
	}
	if !bytes.Equal(cLengthCommitment, c1LengthCommitment) {
		return fmt.Errorf("length commitments are different: %v vs %v", cLengthCommitment, c1LengthCommitment)
	}

	cLengthProof, err := c.LengthProof.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize length proof: %w", err)
	}
	c1LengthProof, err := c1.LengthProof.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize c1 length proof: %w", err)
	}
	if !bytes.Equal(cLengthProof, c1LengthProof) {
		return fmt.Errorf("length proofs are different: %v vs %v", cLengthProof, c1LengthProof)
	}

	return nil
}

func BlobCommitmentsFromProtobuf(c *pbcommon.BlobCommitment) (*BlobCommitments, error) {
	commitment, err := new(G1Commitment).Deserialize(c.GetCommitment())
	if err != nil {
		return nil, err
	}

	lengthCommitment, err := new(G2Commitment).Deserialize(c.GetLengthCommitment())
	if err != nil {
		return nil, err
	}

	lengthProof, err := new(G2Commitment).Deserialize(c.GetLengthProof())
	if err != nil {
		return nil, err
	}

	return &BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           uint(c.GetLength()),
	}, nil
}

// Frame is a chunk of data with the associated multi-reveal proof
type Frame struct {
	// Proof is the multireveal proof corresponding to the chunk
	Proof Proof
	// Coeffs contains the coefficients of the interpolating polynomial of the chunk
	Coeffs []Symbol
}

func (f *Frame) Length() int {
	return len(f.Coeffs)
}

// Size return the size of chunks in bytes.
func (f *Frame) Size() uint64 {
	return uint64(f.Length() * BYTES_PER_SYMBOL)
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

// FragmentInfo contains metadata about how chunk coefficients file is stored.
type FragmentInfo struct {
	// TotalChunkSizeBytes is the total size of the file containing all chunk coefficients for the blob.
	TotalChunkSizeBytes uint32
	// FragmentSizeBytes is the maximum fragment size used to store the chunk coefficients.
	FragmentSizeBytes uint32
}
