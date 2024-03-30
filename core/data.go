package core

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
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

		size += int64(roundUpDivide(b.Length*percentMultiplier*encoding.BYTES_PER_COEFFICIENT, uint(quorum.ConfirmationThreshold-quorum.AdversaryThreshold)))
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
