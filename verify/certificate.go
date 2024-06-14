package verify

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
)

var (
	ErrInvalidDomainType = fmt.Errorf("invalid domain type")
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
			QuorumNumber:                    uint8(qp.QuorumNumber),
			AdversaryThresholdPercentage:    uint8(qp.AdversaryThresholdPercentage),
			ConfirmationThresholdPercentage: uint8(qp.ConfirmationThresholdPercentage),
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
