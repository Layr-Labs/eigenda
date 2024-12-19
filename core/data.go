package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

type AccountID = string

// Security and Quorum Parameters

// QuorumID is a unique identifier for a quorum; initially EigenDA wil support upt to 256 quorums
type QuorumID = uint8

// SecurityParam contains the quorum ID and the adversary threshold for the quorum;
type SecurityParam struct {
	QuorumID QuorumID
	// AdversaryThreshold is the maximum amount of stake that can be controlled by an adversary in the quorum as a percentage of the total stake in the quorum
	AdversaryThreshold uint8
	// ConfirmationThreshold is the amount of stake that must sign a message for it to be considered valid as a percentage of the total stake in the quorum
	ConfirmationThreshold uint8
	// Rate Limit. This is a temporary measure until the node can derive rates on its own using rollup authentication. This is used
	// for restricting the rate at which retrievers are able to download data from the DA node to a multiple of the rate at which the
	// data was posted to the DA node.
	QuorumRate common.RateParam
}

type ChunkEncodingFormat = uint8
type BundleEncodingFormat = uint8

const (
	// We use uint8 to count the number of quorums, so we can have at most 255 quorums,
	// which means the max ID can not be larger than 254 (from 0 to 254, there are 255
	// different IDs).
	MaxQuorumID = 254

	// How many bits for the bundle's header.
	NumBundleHeaderBits = 64
	// How many bits (out of header) for representing the bundle's encoding format.
	NumBundleEncodingFormatBits = 8

	// The list of supported encoding formats for bundle.
	// Values must be in range [0, 255].
	// Note that the IDs here may not be the same as the ChunkEncodingFormat enum in
	// the node.proto file. For example, GobBundleEncodingFormat is 0 here, but in
	// ChunkEncodingFormat the GOB is 2 (and UNKNOWN is 0). The reason is because
	// we need to set GobBundleEncodingFormat to 0 for backward compatibility (and
	// in protobuf, UNKNOWN as 0 is a convention).
	GobBundleEncodingFormat   BundleEncodingFormat = 0
	GnarkBundleEncodingFormat BundleEncodingFormat = 1

	// Similar to bundle encoding format, this describes the encoding format of chunks.
	// The difference is ChunkEncodingFormat is just about chunks, whereas BundleEncodingFormat
	// is also about how multiple chunks of the same bundle are packed into a single byte array.
	GobChunkEncodingFormat   ChunkEncodingFormat = 0
	GnarkChunkEncodingFormat ChunkEncodingFormat = 1
)

type ChunksData struct {
	// Chunks is the encoded bytes of the chunks.
	Chunks [][]byte
	// Format describes how the bytes of the chunks are encoded.
	Format ChunkEncodingFormat
	// The number of symbols in each chunk.
	// Note each chunk of the same blob will always have the same number of symbols.
	ChunkLen int
}

func (cd *ChunksData) Size() uint64 {
	if len(cd.Chunks) == 0 {
		return 0
	}
	// GnarkChunkEncoding will create chunks of equal size.
	if cd.Format == GnarkChunkEncodingFormat {
		return uint64(len(cd.Chunks)) * uint64(len(cd.Chunks[0]))
	}
	// GobChunkEncoding can create chunks of different sizes.
	size := uint64(0)
	for _, c := range cd.Chunks {
		size += uint64(len(c))
	}
	return size
}

func (cd *ChunksData) FromFrames(fr []*encoding.Frame) (*ChunksData, error) {
	if len(fr) == 0 {
		return nil, errors.New("no frame is provided")
	}
	var c ChunksData
	c.Format = GnarkChunkEncodingFormat
	c.ChunkLen = fr[0].Length()
	c.Chunks = make([][]byte, 0, len(fr))
	for _, f := range fr {
		bytes, err := f.SerializeGnark()
		if err != nil {
			return nil, err
		}
		c.Chunks = append(c.Chunks, bytes)
	}
	return &c, nil
}

func (cd *ChunksData) ToFrames() ([]*encoding.Frame, error) {
	frames := make([]*encoding.Frame, 0, len(cd.Chunks))
	switch cd.Format {
	case GobChunkEncodingFormat:
		for _, data := range cd.Chunks {
			fr, err := new(encoding.Frame).Deserialize(data)
			if err != nil {
				return nil, err
			}
			frames = append(frames, fr)
		}
	case GnarkChunkEncodingFormat:
		for _, data := range cd.Chunks {
			fr, err := new(encoding.Frame).DeserializeGnark(data)
			if err != nil {
				return nil, err
			}
			frames = append(frames, fr)
		}
	default:
		return nil, fmt.Errorf("invalid chunk encoding format: %v", cd.Format)
	}
	return frames, nil
}

func (cd *ChunksData) FlattenToBundle() ([]byte, error) {
	// Only Gnark coded chunks are dispersed as a byte array.
	// Gob coded chunks are not flattened.
	if cd.Format != GnarkChunkEncodingFormat {
		return nil, fmt.Errorf("unsupported chunk encoding format to flatten: %v", cd.Format)
	}
	result := make([]byte, cd.Size()+8)
	buf := result
	metadata := (uint64(cd.Format) << (NumBundleHeaderBits - NumBundleEncodingFormatBits)) | uint64(cd.ChunkLen)
	binary.LittleEndian.PutUint64(buf, metadata)
	buf = buf[8:]
	for _, c := range cd.Chunks {
		if len(c) != len(cd.Chunks[0]) {
			return nil, errors.New("all chunks must be of same size")
		}
		copy(buf, c)
		buf = buf[len(c):]
	}
	return result, nil
}

func (cd *ChunksData) ToGobFormat() (*ChunksData, error) {
	if cd.Format == GobChunkEncodingFormat {
		return cd, nil
	}
	if cd.Format != GnarkChunkEncodingFormat {
		return nil, fmt.Errorf("unsupported chunk encoding format: %d", cd.Format)
	}
	gobChunks := make([][]byte, 0, len(cd.Chunks))
	for _, chunk := range cd.Chunks {
		c, err := new(encoding.Frame).DeserializeGnark(chunk)
		if err != nil {
			return nil, err
		}
		gob, err := c.Serialize()
		if err != nil {
			return nil, err
		}
		gobChunks = append(gobChunks, gob)
	}
	return &ChunksData{
		Chunks:   gobChunks,
		Format:   GobChunkEncodingFormat,
		ChunkLen: cd.ChunkLen,
	}, nil
}

func (cd *ChunksData) ToGnarkFormat() (*ChunksData, error) {
	if cd.Format == GnarkChunkEncodingFormat {
		return cd, nil
	}
	if cd.Format != GobChunkEncodingFormat {
		return nil, fmt.Errorf("unsupported chunk encoding format: %d", cd.Format)
	}
	gnarkChunks := make([][]byte, 0, len(cd.Chunks))
	for _, chunk := range cd.Chunks {
		c, err := new(encoding.Frame).Deserialize(chunk)
		if err != nil {
			return nil, err
		}
		gnark, err := c.SerializeGnark()
		if err != nil {
			return nil, err
		}
		gnarkChunks = append(gnarkChunks, gnark)
	}
	return &ChunksData{
		Chunks:   gnarkChunks,
		Format:   GnarkChunkEncodingFormat,
		ChunkLen: cd.ChunkLen,
	}, nil
}

func (s *SecurityParam) String() string {
	return fmt.Sprintf("QuorumID: %d, AdversaryThreshold: %d, ConfirmationThreshold: %d", s.QuorumID, s.AdversaryThreshold, s.ConfirmationThreshold)
}

// QuorumResult contains the quorum ID and the amount signed for the quorum
type QuorumResult struct {
	QuorumID QuorumID
	// PercentSigned is percentage of the total stake for the quorum that signed for a particular batch.
	PercentSigned uint8
}

// Blob stores the data and header of a single data blob. Blobs are the fundamental unit of data posted to EigenDA by users.
type Blob struct {
	RequestHeader BlobRequestHeader
	Data          []byte
}

func (b *Blob) GetQuorumNumbers() []uint8 {
	quorumNumbers := make([]uint8, 0, len(b.RequestHeader.SecurityParams))
	for _, sp := range b.RequestHeader.SecurityParams {
		quorumNumbers = append(quorumNumbers, sp.QuorumID)
	}
	return quorumNumbers
}

// BlobAuthHeader contains the data that a user must sign to authenticate a blob request.
// Signing the combination of the Nonce and the BlobCommitments prohibits the disperser from
// using the signature to charge the user for a different blob or for dispersing the same blob
// multiple times (Replay attack).
type BlobAuthHeader struct {
	// Commitments
	encoding.BlobCommitments `json:"commitments"`
	// AccountID is the account that is paying for the blob to be stored. AccountID is hexadecimal representation of the ECDSA public key
	AccountID AccountID `json:"account_id"`
	// Nonce
	Nonce uint32 `json:"nonce"`
	// AuthenticationData is the signature of the blob header by the account ID
	AuthenticationData []byte `json:"authentication_data"`
}

// BlobRequestHeader contains the original data size of a blob and the security required
type BlobRequestHeader struct {
	// BlobAuthHeader
	BlobAuthHeader `json:"blob_auth_header"`
	// For a blob to be accepted by EigenDA, it satisfy the AdversaryThreshold of each quorum contained in SecurityParams
	SecurityParams []*SecurityParam `json:"security_params"`
}

func ValidateSecurityParam(confirmationThreshold, adversaryThreshold uint32) error {
	if confirmationThreshold > 100 {
		return errors.New("confimration threshold exceeds 100")
	}
	if adversaryThreshold == 0 {
		return errors.New("adversary threshold equals 0")
	}
	if confirmationThreshold < adversaryThreshold || confirmationThreshold-adversaryThreshold < 10 {
		return errors.New("confirmation threshold must be >= 10 + adversary threshold")
	}
	return nil
}

func (sp *SecurityParam) Validate() error {
	return ValidateSecurityParam(uint32(sp.ConfirmationThreshold), uint32(sp.AdversaryThreshold))
}

// BlobQuorumInfo contains the quorum IDs and parameters for a blob specific to a given quorum
type BlobQuorumInfo struct {
	SecurityParam
	// ChunkLength is the number of symbols in a chunk
	ChunkLength uint
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	encoding.BlobCommitments
	// QuorumInfos contains the quorum specific parameters for the blob
	QuorumInfos []*BlobQuorumInfo

	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID
}

func (b *BlobHeader) GetQuorumInfo(quorum QuorumID) *BlobQuorumInfo {
	for _, quorumInfo := range b.QuorumInfos {
		if quorumInfo.QuorumID == quorum {
			return quorumInfo
		}
	}
	return nil
}

// Returns the total encoded size in bytes of the blob across all quorums.
func (b *BlobHeader) EncodedSizeAllQuorums() int64 {
	size := int64(0)
	for _, quorum := range b.QuorumInfos {

		size += int64(RoundUpDivide(b.Length*percentMultiplier*encoding.BYTES_PER_SYMBOL, uint(quorum.ConfirmationThreshold-quorum.AdversaryThreshold)))
	}
	return size
}

// Batch
// A batch is a collection of blobs. DA nodes receive and attest to the blobs in a batch together to amortize signature verification costs

// BatchHeader contains the metadata associated with a Batch for which DA nodes must attest; DA nodes sign on the hash of the batch header
type BatchHeader struct {
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte
}

// EncodedBlob contains the messages to be sent to a group of DA nodes corresponding to a single blob
type EncodedBlob struct {
	BlobHeader        *BlobHeader
	BundlesByOperator map[OperatorID]Bundles
	// EncodedBundlesByOperator is bundles in encoded format (not deserialized)
	EncodedBundlesByOperator map[OperatorID]EncodedBundles
}

// A Bundle is the collection of chunks associated with a single blob, for a single operator and a single quorum.
type Bundle []*encoding.Frame

// Bundles is the collection of bundles associated with a single blob and a single operator.
type Bundles map[QuorumID]Bundle

// This is similar to Bundle, but tracks chunks in encoded format (i.e. not deserialized).
type EncodedBundles map[QuorumID]*ChunksData

// BlobMessage is the message that is sent to DA nodes. It contains the blob header and the associated chunk bundles.
type BlobMessage struct {
	BlobHeader *BlobHeader
	Bundles    Bundles
}

// This is similar to BlobMessage, but keep the commitments and chunks in encoded format
// (i.e. not deserialized)
type EncodedBlobMessage struct {
	// TODO(jianoaix): Change the commitments to encoded format.
	BlobHeader     *BlobHeader
	EncodedBundles map[QuorumID]*ChunksData
}

func (b Bundle) Size() uint64 {
	size := uint64(0)
	for _, chunk := range b {
		size += chunk.Size()
	}
	return size
}

// Serialize returns the serialized bytes of the bundle.
//
// The bytes are packed in this format:
// <8 bytes header><chunk 1 bytes>chunk 2 bytes>...
//
// The header format:
//   - First byte: describes the encoding format. Currently, only GnarkBundleEncodingFormat (1)
//     is supported.
//   - Remaining 7 bytes: describes the information about chunks.
//
// The chunk format will depend on the encoding format. With the GnarkBundleEncodingFormat,
// each chunk is formated as <32 bytes proof><32 bytes coeff>...<32 bytes coefff>, where the
// proof and coeffs are all encoded with Gnark.
func (b Bundle) Serialize() ([]byte, error) {
	if len(b) == 0 {
		return []byte{}, nil
	}
	if len(b[0].Coeffs) == 0 {
		return nil, errors.New("invalid bundle: the coeffs length is zero")
	}
	size := 0
	for _, f := range b {
		if len(f.Coeffs) != len(b[0].Coeffs) {
			return nil, errors.New("invalid bundle: all chunks should have the same length")
		}
		size += bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*len(f.Coeffs)
	}
	result := make([]byte, size+8)
	buf := result
	metadata := (uint64(GnarkBundleEncodingFormat) << (NumBundleHeaderBits - NumBundleEncodingFormatBits)) | uint64(len(b[0].Coeffs))
	binary.LittleEndian.PutUint64(buf, metadata)
	buf = buf[8:]
	for _, f := range b {
		chunk, err := f.SerializeGnark()
		if err != nil {
			return nil, err
		}
		copy(buf, chunk)
		buf = buf[len(chunk):]
	}
	return result, nil
}

func (b Bundle) Deserialize(data []byte) (Bundle, error) {
	if len(data) < 8 {
		return nil, errors.New("bundle data must have at least 8 bytes")
	}
	// Parse metadata
	meta := binary.LittleEndian.Uint64(data)
	if (meta >> (NumBundleHeaderBits - NumBundleEncodingFormatBits)) != uint64(GnarkBundleEncodingFormat) {
		return nil, errors.New("invalid bundle data encoding format")
	}
	chunkLen := (meta << NumBundleEncodingFormatBits) >> NumBundleEncodingFormatBits
	if chunkLen == 0 {
		return nil, errors.New("chunk length must be greater than zero")
	}
	chunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*int(chunkLen)
	if (len(data)-8)%chunkSize != 0 {
		return nil, errors.New("bundle data is invalid")
	}
	// Decode
	bundle := make([]*encoding.Frame, 0, (len(data)-8)/chunkSize)
	buf := data[8:]
	for len(buf) > 0 {
		if len(buf) < chunkSize {
			return nil, errors.New("bundle data is invalid")
		}
		f, err := new(encoding.Frame).DeserializeGnark(buf[:chunkSize])
		if err != nil {
			return nil, err
		}
		bundle = append(bundle, f)
		buf = buf[chunkSize:]
	}
	return bundle, nil
}

// Serialize encodes a batch of chunks into a byte array
func (cb Bundles) Serialize() (map[uint32][][]byte, error) {
	data := make(map[uint32][][]byte, len(cb))
	for quorumID, bundle := range cb {
		for _, chunk := range bundle {
			chunkData, err := chunk.Serialize()
			if err != nil {
				return nil, err
			}
			data[uint32(quorumID)] = append(data[uint32(quorumID)], chunkData)
		}
	}
	return data, nil
}

// Returns the size of the bundles in bytes.
func (cb Bundles) Size() uint64 {
	size := uint64(0)
	for _, bundle := range cb {
		size += bundle.Size()
	}
	return size
}

func (cb Bundles) ToEncodedBundles() (EncodedBundles, error) {
	eb := make(EncodedBundles)
	for quorum, bundle := range cb {
		cd, err := new(ChunksData).FromFrames(bundle)
		if err != nil {
			return nil, err
		}
		eb[quorum] = cd
	}
	return eb, nil
}

func (cb Bundles) FromEncodedBundles(eb EncodedBundles) (Bundles, error) {
	c := make(Bundles)
	for quorum, chunkData := range eb {
		fr, err := chunkData.ToFrames()
		if err != nil {
			return nil, err
		}
		c[quorum] = fr
	}
	return c, nil
}

// PaymentMetadata represents the header information for a blob
type PaymentMetadata struct {
	// AccountID is the ETH account address for the payer
	AccountID string `json:"account_id"`

	// ReservationPeriod represents the range of time at which the dispersal is made
	ReservationPeriod uint32 `json:"reservation_period"`
	// TODO: we are thinking the contract can use uint128 for cumulative payment,
	// but the definition on v2 uses uint64. Double check with team.
	CumulativePayment *big.Int `json:"cumulative_payment"`
	// Allow same blob to be dispersed multiple times within the same reservation period
	Salt uint32 `json:"salt"`
}

// Hash returns the Keccak256 hash of the PaymentMetadata
func (pm *PaymentMetadata) Hash() ([32]byte, error) {
	if pm == nil {
		return [32]byte{}, errors.New("payment metadata is nil")
	}
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "accountID",
			Type: "string",
		},
		{
			Name: "reservationPeriod",
			Type: "uint32",
		},
		{
			Name: "cumulativePayment",
			Type: "uint256",
		},
		{
			Name: "salt",
			Type: "uint32",
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

	bytes, err := arguments.Pack(pm)
	if err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(hash[:], hasher.Sum(nil)[:32])

	return hash, nil
}

func (pm *PaymentMetadata) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if pm == nil {
		return nil, errors.New("payment metadata is nil")
	}

	return &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"AccountID":         &types.AttributeValueMemberS{Value: pm.AccountID},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", pm.ReservationPeriod)},
			"CumulativePayment": &types.AttributeValueMemberN{
				Value: pm.CumulativePayment.String(),
			},
			"Salt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", pm.Salt)},
		},
	}, nil
}

func (pm *PaymentMetadata) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberM, got %T", av)
	}
	pm.AccountID = m.Value["AccountID"].(*types.AttributeValueMemberS).Value
	reservationPeriod, err := strconv.ParseUint(m.Value["ReservationPeriod"].(*types.AttributeValueMemberN).Value, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse ReservationPeriod: %w", err)
	}
	pm.ReservationPeriod = uint32(reservationPeriod)
	pm.CumulativePayment, _ = new(big.Int).SetString(m.Value["CumulativePayment"].(*types.AttributeValueMemberN).Value, 10)
	salt, err := strconv.ParseUint(m.Value["Salt"].(*types.AttributeValueMemberN).Value, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse Salt: %w", err)
	}
	pm.Salt = uint32(salt)
	return nil
}

func (pm *PaymentMetadata) ToProtobuf() *commonpb.PaymentHeader {
	if pm == nil {
		return nil
	}
	return &commonpb.PaymentHeader{
		AccountId:         pm.AccountID,
		ReservationPeriod: pm.ReservationPeriod,
		CumulativePayment: pm.CumulativePayment.Bytes(),
		Salt:              pm.Salt,
	}
}

// ConvertToProtoPaymentHeader converts a PaymentMetadata to a protobuf payment header
func ConvertToPaymentMetadata(ph *commonpb.PaymentHeader) *PaymentMetadata {
	if ph == nil {
		return nil
	}

	return &PaymentMetadata{
		AccountID:         ph.AccountId,
		ReservationPeriod: ph.ReservationPeriod,
		CumulativePayment: new(big.Int).SetBytes(ph.CumulativePayment),
		Salt:              ph.Salt,
	}
}

// ReservedPayment contains information the onchain state about a reserved payment
type ReservedPayment struct {
	// reserve number of symbols per second
	SymbolsPerSecond uint64
	// reservation activation timestamp
	StartTimestamp uint64
	// reservation expiration timestamp
	EndTimestamp uint64

	// allowed quorums
	QuorumNumbers []uint8
	// ordered mapping of quorum number to payment split; on-chain validation should ensure split <= 100
	QuorumSplits []byte
}

type OnDemandPayment struct {
	// Total amount deposited by the user
	CumulativePayment *big.Int
}

type BlobVersionParameters struct {
	CodingRate      uint32
	MaxNumOperators uint32
	NumChunks       uint32
}

// IsActive returns true if the reservation is active at the given timestamp
func (ar *ReservedPayment) IsActive(currentTimestamp uint64) bool {
	return ar.StartTimestamp <= currentTimestamp && ar.EndTimestamp >= currentTimestamp
}
