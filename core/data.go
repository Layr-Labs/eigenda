package core

import (
	"errors"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type AccountID = string

// Security and Quorum Parameters

// QuorumID is a unique identifier for a quorum; initially EigenDA wil support upt to 256 quorums
type QuorumID = uint8

// SecurityParam contains the quorum ID and the adversary threshold for the quorum;
type SecurityParam struct {
	QuorumID QuorumID `json:"quorum_id"`
	// AdversaryThreshold is the maximum amount of stake that can be controlled by an adversary in the quorum as a percentage of the total stake in the quorum
	AdversaryThreshold uint8 `json:"adversary_threshold"`
	// QuorumThreshold is the amount of stake that must sign a message for it to be considered valid as a percentage of the total stake in the quorum
	QuorumThreshold uint8 `json:"quorum_threshold"`
	// Rate Limit. This is a temporary measure until the node can derive rates on its own using rollup authentication. This is used
	// for restricting the rate at which retrievers are able to download data from the DA node to a multiple of the rate at which the
	// data was posted to the DA node.
	QuorumRate common.RateParam `json:"quorum_rate"`
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

type BlobAuthHeader struct {
	// Commitments
	BlobCommitments `json:"commitments"`
	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
	// Nonce
	Nonce uint32 `json:"nonce"`
	// AuthenticationData is the signature of the blob header by the account ID
	AuthenticationData []byte `json:"authentication_data"`
}

// BlobRequestHeader contains the original data size of a blob and the security required
type BlobRequestHeader struct {
	// Commitments
	BlobCommitments `json:"commitments"`
	// For a blob to be accepted by EigenDA, it satisfy the AdversaryThreshold of each quorum contained in SecurityParams
	SecurityParams []*SecurityParam `json:"security_params"`
	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
}

func (h *BlobRequestHeader) Validate() error {
	for _, quorum := range h.SecurityParams {
		if quorum.QuorumThreshold < quorum.AdversaryThreshold+10 {
			return errors.New("invalid request: quorum threshold must be >= 10 + adversary threshold")
		}
		if quorum.QuorumThreshold > 100 {
			return errors.New("invalid request: quorum threshold exceeds 100")
		}
		if quorum.AdversaryThreshold == 0 {
			return errors.New("invalid request: adversary threshold equals 0")
		}
	}
	return nil
}

// BlobQuorumInfo contains the quorum IDs and parameters for a blob specific to a given quorum
type BlobQuorumInfo struct {
	SecurityParam
	// ChunkLength is the number of symbols in a chunk
	ChunkLength uint
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	BlobCommitments
	// QuorumInfos contains the quorum specific parameters for the blob
	QuorumInfos []*BlobQuorumInfo

	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
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

		size += int64(roundUpDivide(b.Length*percentMultiplier*bn254.BYTES_PER_COEFFICIENT, uint(quorum.QuorumThreshold-quorum.AdversaryThreshold)))
	}
	return size
}

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment  *Commitment `json:"commitment"`
	LengthProof *Commitment `json:"length_proof"`
	Length      uint        `json:"length"`
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
type EncodedBlob = map[OperatorID]*BlobMessage

// Chunks

// Chunk is the smallest unit that is distributed to DA nodes, including both data and the associated polynomial opening proofs.
// A chunk corresponds to a set of evaluations of the global polynomial whose coefficients are used to construct the blob Commitment.
type Chunk struct {
	// The Coeffs field contains the coefficients of the polynomial which interolates these evaluations. This is the same as the
	// interpolating polynomial, I(X), used in the KZG multi-reveal (https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html#multiproofs)
	Coeffs []Symbol
	Proof  Proof
}

func (c *Chunk) Length() int {
	return len(c.Coeffs)
}

// Returns the size of chunk in bytes.
func (c *Chunk) Size() int {
	return c.Length() * bn254.BYTES_PER_COEFFICIENT
}

// A Bundle is the collection of chunks associated with a single blob, for a single operator and a single quorum.
type Bundle = []*Chunk

// Bundles is the collection of bundles associated with a single blob and a single operator.
type Bundles map[QuorumID]Bundle

// BlobMessage is the message that is sent to DA nodes. It contains the blob header and the associated chunk bundles.
type BlobMessage struct {
	BlobHeader *BlobHeader
	Bundles    Bundles
}

// Serialize encodes a batch of chunks into a byte array
func (cb Bundles) Serialize() ([][][]byte, error) {
	data := make([][][]byte, len(cb))
	for i, bundle := range cb {
		for _, chunk := range bundle {
			chunkData, err := chunk.Serialize()
			if err != nil {
				return nil, err
			}
			data[i] = append(data[i], chunkData)
		}
	}
	return data, nil
}

// Returns the size of the bundles in bytes.
func (cb Bundles) Size() int64 {
	size := int64(0)
	for _, bundle := range cb {
		for _, chunk := range bundle {
			size += int64(chunk.Size())
		}
	}
	return size
}

// Sample is a chunk with associated metadata used by the Universal Batch Verifier
type Sample struct {
	Commitment      *Commitment
	Chunk           *Chunk
	AssignmentIndex ChunkNumber
	BlobIndex       int
}

// SubBatch is a part of the whole Batch with identical Encoding Parameters, i.e. (ChunkLen, NumChunk)
// Blobs with the same encoding parameters are collected in a single subBatch
type SubBatch struct {
	Samples  []Sample
	NumBlobs int
}
