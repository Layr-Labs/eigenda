package v2

import (
	"encoding/hex"
	"math"
	"math/big"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

var (
	// TODO(mooselumph): Put these parameters on chain and add on-chain checks to ensure that the number of operators does not
	// conflict with the existing on-chain limits
	ParametersMap = map[BlobVersion]BlobVersionParameters{
		0: {CodingRate: 8, ReconstructionThreshold: 0.22, NumChunks: 8192},
	}
)

type BlobVersion uint8

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

type BlobKey [32]byte

func (b BlobKey) Hex() string {
	return hex.EncodeToString(b[:])
}

func HexToBlobKey(h string) (BlobKey, error) {
	s := strings.TrimPrefix(h, "0x")
	s = strings.TrimPrefix(s, "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return BlobKey{}, err
	}
	return BlobKey(b), nil
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	BlobVersion BlobVersion

	BlobCommitments encoding.BlobCommitments

	// QuorumNumbers contains the quorums the blob is dispersed to
	QuorumNumbers []core.QuorumID

	// PaymentMetadata contains the payment information for the blob
	PaymentMetadata core.PaymentMetadata

	// Signature is the signature of the blob header by the account ID
	Signature []byte
}

func (b *BlobHeader) GetEncodingParams() (encoding.EncodingParams, error) {
	params := ParametersMap[b.BlobVersion]

	length, err := GetChunkLength(b.BlobVersion, uint32(b.BlobCommitments.Length))
	if err != nil {
		return encoding.EncodingParams{}, err
	}

	return encoding.EncodingParams{
		NumChunks:   uint64(params.NumChunks),
		ChunkLength: uint64(length),
	}, nil
}

func (b *BlobHeader) BlobKey() (BlobKey, error) {
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobVersion",
			Type: "uint8",
		},
		{
			Name: "blobCommitments",
			Type: "tuple",
			Components: []abi.ArgumentMarshaling{
				{
					Name: "commitment",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256",
						},
						{
							Name: "Y",
							Type: "uint256",
						},
					},
				},
				{
					Name: "lengthCommitment",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256[2]",
						},
						{
							Name: "Y",
							Type: "uint256[2]",
						},
					},
				},
				{
					Name: "lengthProof",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256[2]",
						},
						{
							Name: "Y",
							Type: "uint256[2]",
						},
					},
				},
				{
					Name: "length",
					Type: "uint32",
				},
			},
		},
		{
			Name: "quorumNumbers",
			Type: "bytes",
		},
		{
			Name: "paymentMetadataHash",
			Type: "bytes32",
		},
	})
	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: blobHeaderType,
		},
	}

	type g1Commit struct {
		X *big.Int
		Y *big.Int
	}
	type g2Commit struct {
		X [2]*big.Int
		Y [2]*big.Int
	}
	type blobCommitments struct {
		Commitment       g1Commit
		LengthCommitment g2Commit
		LengthProof      g2Commit
		Length           uint32
	}

	paymentHash, err := b.PaymentMetadata.Hash()
	if err != nil {
		return [32]byte{}, err
	}
	s := struct {
		BlobVersion         uint8
		BlobCommitments     blobCommitments
		QuorumNumbers       []byte
		PaymentMetadataHash [32]byte
	}{
		BlobVersion: uint8(b.BlobVersion),
		BlobCommitments: blobCommitments{
			Commitment: g1Commit{
				X: b.BlobCommitments.Commitment.X.BigInt(new(big.Int)),
				Y: b.BlobCommitments.Commitment.Y.BigInt(new(big.Int)),
			},
			LengthCommitment: g2Commit{
				X: [2]*big.Int{
					b.BlobCommitments.LengthCommitment.X.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthCommitment.X.A1.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					b.BlobCommitments.LengthCommitment.Y.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthCommitment.Y.A1.BigInt(new(big.Int)),
				},
			},
			LengthProof: g2Commit{
				X: [2]*big.Int{
					b.BlobCommitments.LengthProof.X.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthProof.X.A1.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					b.BlobCommitments.LengthProof.Y.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthProof.Y.A1.BigInt(new(big.Int)),
				},
			},
			Length: uint32(b.BlobCommitments.Length),
		},
		QuorumNumbers:       b.QuorumNumbers,
		PaymentMetadataHash: paymentHash,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
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
