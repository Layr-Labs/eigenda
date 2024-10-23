package v2

import (
	"math"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
)

var (
	// TODO(mooselumph): Put these parameters on chain and add on-chain checks to ensure that the number of operators does not
	// conflict with the existing on-chain limits
	ParametersMap = map[uint8]BlobVersionParameters{
		0: {CodingRate: 8, ReconstructionThreshold: 0.22, NumChunks: 8192},
	}
)

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

	// PaymentHeader contains the payment information for the blob
	core.PaymentMetadata

	// AuthenticationData is the signature of the blob header by the account ID
	AuthenticationData []byte `json:"authentication_data"`
}

func (b *BlobHeader) GetEncodingParams() (encoding.EncodingParams, error) {

	params := ParametersMap[b.Version]

	length, err := GetChunkLength(b.Version, uint32(b.Length))
	if err != nil {
		return encoding.EncodingParams{}, err
	}

	return encoding.EncodingParams{
		NumChunks:   uint64(params.NumChunks),
		ChunkLength: uint64(length),
	}, nil

}

type BlobCertificate struct {
	BlobHeader

	// ReferenceBlockNumber is the block number of the block at which the operator state will be referenced
	ReferenceBlockNumber uint64

	// RelayKeys
	RelayKeys []uint16
}

type BlobVersionParameters struct {
	CodingRate              uint32
	ReconstructionThreshold float64
	NumChunks               uint32
}

func (p BlobVersionParameters) MaxNumOperators() uint32 {

	return uint32(math.Floor(float64(p.NumChunks) * (1 - 1/(p.ReconstructionThreshold*float64(p.CodingRate)))))

}

const (
	// We use uint8 to count the number of quorums, so we can have at most 255 quorums,
	// which means the max ID can not be larger than 254 (from 0 to 254, there are 255
	// different IDs).
	MaxQuorumID = 254
)
