package eigenda

import (
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

// EigenDA finalized blob certificate
type Cert struct {
	BatchHeaderHash      []byte
	BlobIndex            uint32
	ReferenceBlockNumber uint32
	QuorumIDs            []uint32
	// Used for kzg verification when reading blob data from DA
	BlobCommitment *common.G1Commitment
}

// BlobCommitmentFields transforms commitment byte arrays to bn254 field elements
func (c *Cert) BlobCommitmentFields() (*fp.Element, *fp.Element) {
	xBytes, yBytes := c.BlobCommitment.X, c.BlobCommitment.Y
	xElement, yElement := &fp.Element{}, &fp.Element{}

	xElement.Unmarshal(xBytes)
	yElement.Unmarshal(yBytes)

	return xElement, yElement
}
