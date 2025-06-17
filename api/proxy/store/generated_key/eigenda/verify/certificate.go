package verify

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
)

// G1Point struct to represent G1Point in Solidity
type G1Point struct {
	X *big.Int
	Y *big.Int
}

// QuorumBlobParam struct to represent QuorumBlobParam in Solidity
type QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// BlobHeader struct to represent BlobHeader in Solidity
type BlobHeader struct {
	Commitment       G1Point
	DataLength       uint32
	QuorumBlobParams []QuorumBlobParam
}

type Certificate disperser.BlobInfo

// NoNilFields ... checks if any referenced fields in the certificate
// are nil and returns an error if so
func (c *Certificate) NoNilFields() error {
	if c.BlobVerificationProof == nil {
		return fmt.Errorf("BlobVerificationProof is nil")
	}

	if c.BlobVerificationProof.BatchMetadata == nil {
		return fmt.Errorf("BlobVerificationProof.BatchMetadata is nil")
	}

	if c.BlobVerificationProof.BatchMetadata.BatchHeader == nil {
		return fmt.Errorf("BlobVerificationProof.BatchMetadata.BatchHeader is nil")
	}

	if c.BlobHeader == nil {
		return fmt.Errorf("BlobHeader is nil")
	}

	if c.BlobHeader.Commitment == nil {
		return fmt.Errorf("BlobHeader.Commitment is nil")
	}

	return nil
}

// ValidFieldLengths ... enforces length invariance on certificate fields which are expected
// to be size constrained but are read as unfixed byte arrays from the disperser.
// This is necessary to remove a trust assumption and grieving vector where the
// disperser can intentionally increase the data sizes and cause a rollup to incur higher
// operating costs when publishing certificates to some batcher inbox
func (c *Certificate) ValidFieldLengths() error {
	bvp := c.BlobVerificationProof
	bh := c.BlobHeader

	// 1 - necessary invariants to remove disperser trust assumption

	// 1.a necessary since only first 32 bytes of header hash are checked
	//     in verification equivalence check which could allow data padding at end
	if hashLen := len(bvp.BatchMetadata.BatchHeaderHash); hashLen != 32 {
		return fmt.Errorf("BlobVerification.BatchMetadata.BatchHeaderHash is not 32 bytes, got %d", hashLen)
	}

	// 1.b necessary since commitment verification parses the byte field byte arrays
	//     into a field element representation which disregards 0x0 padded bytes
	if xLen := len(bh.Commitment.X); xLen != 32 {
		return fmt.Errorf("BlobHeader.Commitment.X is not 32 bytes, got %d", xLen)
	}

	if yLen := len(bh.Commitment.Y); yLen != 32 {
		return fmt.Errorf("BlobHeader.Commitment.Y is not 32 bytes, got %d", yLen)
	}

	// 2 - unnecessary but preemptive checks that would trigger failure in downstream
	//     verification since these values are used as input for batch metadata hash
	//     recomputation. Capturing here is more efficient!

	if hashLen := len(bvp.BatchMetadata.SignatoryRecordHash); hashLen != 32 {
		return fmt.Errorf("BlobVerification.BatchMetadata.SignatoryRecordHash is not 32 bytes, got %d", hashLen)
	}

	if hashLen := len(bvp.BatchMetadata.BatchHeader.BatchRoot); hashLen != 32 {
		return fmt.Errorf("BlobVerification.BatchMetadata.BatchHeader.BatchRoot is not 32 bytes, got %d", hashLen)
	}

	return nil
}

func (c *Certificate) BlobIndex() uint32 {
	return c.BlobVerificationProof.BlobIndex
}

func (c *Certificate) BatchHeaderRoot() []byte {
	return c.BlobVerificationProof.BatchMetadata.BatchHeader.BatchRoot
}

func (c *Certificate) ReadBlobHeader() BlobHeader {
	// parse quorum params

	qps := make([]QuorumBlobParam, len(c.BlobHeader.BlobQuorumParams))
	for i, qp := range c.BlobHeader.BlobQuorumParams {
		qps[i] = QuorumBlobParam{
			QuorumNumber:                    uint8(qp.QuorumNumber),                    // #nosec G115
			AdversaryThresholdPercentage:    uint8(qp.AdversaryThresholdPercentage),    // #nosec G115
			ConfirmationThresholdPercentage: uint8(qp.ConfirmationThresholdPercentage), // #nosec G115
			ChunkLength:                     qp.ChunkLength,
		}
	}

	return BlobHeader{
		Commitment: G1Point{
			X: new(big.Int).SetBytes(c.BlobHeader.Commitment.X),
			Y: new(big.Int).SetBytes(c.BlobHeader.Commitment.Y),
		},
		DataLength:       c.BlobHeader.DataLength,
		QuorumBlobParams: qps,
	}
}

func (c *Certificate) Proof() *disperser.BlobVerificationProof {
	return c.BlobVerificationProof
}
