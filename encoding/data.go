package encoding

import (
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	framepb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
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

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment       *G1Commitment `json:"commitment"`
	LengthCommitment *G2Commitment `json:"length_commitment"`
	LengthProof      *LengthProof  `json:"length_proof"`
	Length           uint          `json:"length"`
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

func BlobCommitmentsFromProtobuf(c *pbcommon.BlobCommitment) (*BlobCommitments, error) {
	commitment, err := new(G1Commitment).Deserialize(c.Commitment)
	if err != nil {
		return nil, err
	}

	lengthCommitment, err := new(G2Commitment).Deserialize(c.LengthCommitment)
	if err != nil {
		return nil, err
	}

	lengthProof, err := new(G2Commitment).Deserialize(c.LengthProof)
	if err != nil {
		return nil, err
	}

	return &BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           uint(c.Length),
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

// ToProtobuf converts the FragmentInfo to protobuf format
func (f *Frame) ToProtobuf() *framepb.Frame {
	proof := &framepb.Proof{
		X: fpElementToProtobuf(&f.Proof.X),
		Y: fpElementToProtobuf(&f.Proof.Y),
	}

	coeffs := make([]*framepb.Element, len(f.Coeffs))
	for i, c := range f.Coeffs {
		coeffs[i] = frElementToProtobuf(&c)
	}

	return &framepb.Frame{
		Proof:  proof,
		Coeffs: coeffs,
	}
}

// fpElementToProtobuf converts an fp.Element to protobuf format
func fpElementToProtobuf(e *fp.Element) *framepb.Element {
	return &framepb.Element{
		C0: e[0],
		C1: e[1],
		C2: e[2],
		C3: e[3],
	}
}

// frElementToProtobuf converts an fr.Element to protobuf format
func frElementToProtobuf(e *fr.Element) *framepb.Element {
	return &framepb.Element{
		C0: e[0],
		C1: e[1],
		C2: e[2],
		C3: e[3],
	}
}

// fpElementFromProtobuf converts a protobuf element to an fp.Element
func fpElementFromProtobuf(e *framepb.Element) fp.Element {
	return fp.Element{
		e.C0,
		e.C1,
		e.C2,
		e.C3,
	}
}

// frElementFromProtobuf converts a protobuf element to an fr.Element
func frElementFromProtobuf(e *framepb.Element) fr.Element {
	return fr.Element{
		e.C0,
		e.C1,
		e.C2,
		e.C3,
	}
}

// FrameFromProtobuf converts a protobuf frame to a Frame.
func FrameFromProtobuf(f *framepb.Frame) *Frame {
	proof := Proof{
		X: fpElementFromProtobuf(f.Proof.X),
		Y: fpElementFromProtobuf(f.Proof.Y),
	}

	coeffs := make([]Symbol, len(f.Coeffs))
	for i, c := range f.Coeffs {
		coeffs[i] = frElementFromProtobuf(c)
	}

	return &Frame{
		Proof:  proof,
		Coeffs: coeffs,
	}
}

// FramesToProtobuf converts a slice of Frames to protobuf format
func FramesToProtobuf(frames []*Frame) []*framepb.Frame {
	protobufFrames := make([]*framepb.Frame, len(frames))
	for i, f := range frames {
		protobufFrames[i] = f.ToProtobuf()
	}
	return protobufFrames
}

// FramesFromProtobuf converts a slice of protobuf Frames to Frames.
func FramesFromProtobuf(frames []*framepb.Frame) []*Frame {
	protobufFrames := make([]*Frame, len(frames))
	for i, f := range frames {
		protobufFrames[i] = FrameFromProtobuf(f)
	}
	return protobufFrames
}
