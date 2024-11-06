package v2

import (
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"strings"

	pb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
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

func NewBlobHeader(proto *pb.BlobHeader) (*BlobHeader, error) {
	commitment, err := new(encoding.G1Commitment).Deserialize(proto.GetCommitment().GetCommitment())
	if err != nil {
		return nil, err
	}
	lengthCommitment, err := new(encoding.G2Commitment).Deserialize(proto.GetCommitment().GetLengthCommitment())
	if err != nil {
		return nil, err
	}
	lengthProof, err := new(encoding.LengthProof).Deserialize(proto.GetCommitment().GetLengthProof())
	if err != nil {
		return nil, err
	}

	if !(*bn254.G1Affine)(commitment).IsInSubGroup() {
		return nil, errors.New("commitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(lengthCommitment).IsInSubGroup() {
		return nil, errors.New("lengthCommitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(lengthProof).IsInSubGroup() {
		return nil, errors.New("lengthProof is not in the subgroup")
	}

	quorumNumbers := make([]core.QuorumID, len(proto.QuorumNumbers))
	for i, q := range proto.GetQuorumNumbers() {
		if q > MaxQuorumID {
			return nil, errors.New("quorum number exceeds maximum allowed")
		}
		quorumNumbers[i] = core.QuorumID(q)
	}

	paymentMetadata := core.PaymentMetadata{
		AccountID:         proto.GetPaymentHeader().GetAccountId(),
		BinIndex:          proto.GetPaymentHeader().GetBinIndex(),
		CumulativePayment: new(big.Int).SetBytes(proto.GetPaymentHeader().GetCumulativePayment()),
	}

	return &BlobHeader{
		BlobVersion: BlobVersion(proto.GetVersion()),
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: lengthCommitment,
			LengthProof:      lengthProof,
			Length:           uint(proto.GetCommitment().GetLength()),
		},
		QuorumNumbers:   quorumNumbers,
		PaymentMetadata: paymentMetadata,
		Signature:       proto.GetSignature(),
	}, nil
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

type RelayKey uint16

type BlobCertificate struct {
	BlobHeader *BlobHeader

	// RelayKeys
	RelayKeys []RelayKey
}

type BatchHeader struct {
	BatchRoot            [32]byte
	ReferenceBlockNumber uint64
}

// GetBatchHeaderHash returns the hash of the batch header
func (h BatchHeader) Hash() ([32]byte, error) {
	var headerHash [32]byte

	// The order here has to match the field ordering of ReducedBatchHeader defined in IEigenDAServiceManager.sol
	// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43
	batchHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobHeadersRoot",
			Type: "bytes32",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})
	if err != nil {
		return headerHash, err
	}

	arguments := abi.Arguments{
		{
			Type: batchHeaderType,
		},
	}

	s := struct {
		BlobHeadersRoot      [32]byte
		ReferenceBlockNumber uint32
	}{
		BlobHeadersRoot:      h.BatchRoot,
		ReferenceBlockNumber: uint32(h.ReferenceBlockNumber),
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return headerHash, err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

type Batch struct {
	BatchHeader      *BatchHeader
	BlobCertificates []*BlobCertificate
}

type BlobVerificationInfo struct {
	BlobCertificate *BlobCertificate
	BlobIndex       uint32
	InclusionProof  []byte
}

type BlobVersionParameters struct {
	CodingRate              uint32
	ReconstructionThreshold float64
	NumChunks               uint32
}

func (p BlobVersionParameters) MaxNumOperators() uint32 {
	return uint32(math.Floor(float64(p.NumChunks) * (1 - 1/(p.ReconstructionThreshold*float64(p.CodingRate)))))
}

// DispersalRequest is a request to disperse a batch to a specific operator
type DispersalRequest struct {
	core.OperatorID `dynamodbav:"-"`
	OperatorAddress gethcommon.Address
	Socket          string
	DispersedAt     uint64

	BatchHeader
}

// DispersalResponse is a response to a dispersal request
type DispersalResponse struct {
	*DispersalRequest

	RespondedAt uint64
	// Signature is the signature of the response by the operator
	Signature [32]byte
	// Error is the error message if the dispersal failed
	Error string
}

const (
	// We use uint8 to count the number of quorums, so we can have at most 255 quorums,
	// which means the max ID can not be larger than 254 (from 0 to 254, there are 255
	// different IDs).
	MaxQuorumID = 254
)
