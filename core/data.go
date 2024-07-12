package core

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
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
	GobBundleEncodingFormat   = 0
	GnarkBundleEncodingFormat = 1
)

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

		size += int64(roundUpDivide(b.Length*percentMultiplier*encoding.BYTES_PER_SYMBOL, uint(quorum.ConfirmationThreshold-quorum.AdversaryThreshold)))
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
}

// A Bundle is the collection of chunks associated with a single blob, for a single operator and a single quorum.
type Bundle []*encoding.Frame

// Bundles is the collection of bundles associated with a single blob and a single operator.
type Bundles map[QuorumID]Bundle

// BlobMessage is the message that is sent to DA nodes. It contains the blob header and the associated chunk bundles.
type BlobMessage struct {
	BlobHeader *BlobHeader
	Bundles    Bundles
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
	if (meta >> (NumBundleHeaderBits - NumBundleEncodingFormatBits)) != GnarkBundleEncodingFormat {
		return nil, errors.New("invalid bundle data encoding format")
	}
	chunkLen := meta << NumBundleEncodingFormatBits >> NumBundleEncodingFormatBits
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
