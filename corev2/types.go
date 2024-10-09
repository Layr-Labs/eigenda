package corev2

import (
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
)

var (
	ParametersMap = map[uint8]BlobVersionParameters{
		0: {CodingRate: 8, ReconstructionThreshold: 0.22, NumChunks: 8192},
	}
)

type QuorumID = uint8

type OperatorID = [32]byte

// Assignment contains information about the set of chunks that a specific node will receive
type Assignment struct {
	StartIndex uint32
	NumChunks  uint32
}

// GetIndices generates the list of ChunkIndices associated with a given assignment
func (c *Assignment) GetIndices() []uint32 {
	indices := make([]uint32, c.NumChunks)
	for ind := range indices {
		indices[ind] = c.StartIndex + uint32(ind)
	}
	return indices
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	Version uint8

	encoding.BlobCommitments

	// QuorumInfos contains the quorum specific parameters for the blob
	QuorumNumbers []uint8

	// ReferenceBlockNumber is the block number of the block at which the operator state will be referenced
	ReferenceBlockNumber uint64
}

type PaymentHeader struct {
	// BlobKey is the hash of the blob header
	BlobKey [32]byte

	// AccountID is the account that is paying for the blob to be stored. AccountID is hexadecimal representation of the ECDSA public key
	AccountID string

	// Cumulative Payment
	CumulativePayment uint64

	BinIndex uint64

	// AuthenticationData is the signature of the blob header by the account ID
	AuthenticationData []byte `json:"authentication_data"`
}

type BlobVersionParameters struct {
	CodingRate              uint32
	ReconstructionThreshold float64
	NumChunks               uint32
}

func (p BlobVersionParameters) MaxNumOperators() uint32 {

	return uint32(math.Floor(float64(p.NumChunks) * (1 - 1/(p.ReconstructionThreshold*float64(p.CodingRate)))))

}
