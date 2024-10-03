package disperser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/api/grpc/node"
	gcommon "github.com/ethereum/go-ethereum/common"
)

type BlobStatus uint

// WARNING: THESE VALUES BECOME PART OF PERSISTENT SYSTEM STATE;
// ALWAYS INSERT NEW ENUM VALUES AS THE LAST ELEMENT TO MAINTAIN COMPATIBILITY
const (
	Processing BlobStatus = iota
	Confirmed
	Failed
	Finalized
	InsufficientSignatures
	Dispersing
)

var enumStrings = map[BlobStatus]string{
	Processing:             "Processing",
	Confirmed:              "Confirmed",
	Failed:                 "Failed",
	Finalized:              "Finalized",
	InsufficientSignatures: "InsufficientSignatures",
	Dispersing:             "Dispersing",
}

func (bs BlobStatus) String() string {
	if str, ok := enumStrings[bs]; ok {
		return str
	}
	return "Unknown value"
}

type BlobHash = string
type MetadataHash = string

type BlobKey struct {
	BlobHash     BlobHash
	MetadataHash MetadataHash
}

func (mk BlobKey) String() string {
	return fmt.Sprintf("%s-%s", mk.BlobHash, mk.MetadataHash)
}

func ParseBlobKey(key string) (BlobKey, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return BlobKey{}, fmt.Errorf("invalid metadata key: %s", key)
	}
	return BlobKey{
		BlobHash:     parts[0],
		MetadataHash: parts[1],
	}, nil
}

type BlobMetadata struct {
	BlobHash     BlobHash     `json:"blob_hash"`
	MetadataHash MetadataHash `json:"metadata_hash"`
	BlobStatus   BlobStatus   `json:"blob_status"`
	// Expiry is unix epoch time in seconds at which the blob will expire
	Expiry uint64 `json:"expiry"`
	// NumRetries is the number of times the blob has been retried
	// After few failed attempts, the blob will be marked as failed
	NumRetries uint `json:"num_retries"`
	// RequestMetadata is the request metadata of the blob when it was requested
	// This field is omitted when marshalling to DynamoDB attributevalue as this field will be flattened
	RequestMetadata *RequestMetadata `json:"request_metadata" dynamodbav:"-"`
	// ConfirmationInfo is the confirmation metadata of the blob when it was confirmed
	// This field is nil if the blob has not been confirmed
	// This field is omitted when marshalling to DynamoDB attributevalue as this field will be flattened
	ConfirmationInfo *ConfirmationInfo `json:"blob_confirmation_info" dynamodbav:"-"`
}

func (m *BlobMetadata) GetBlobKey() BlobKey {
	return BlobKey{
		BlobHash:     m.BlobHash,
		MetadataHash: m.MetadataHash,
	}
}

func (m *BlobMetadata) IsConfirmed() (bool, error) {
	if m.BlobStatus != Confirmed && m.BlobStatus != Finalized {
		return false, nil
	}

	if m.ConfirmationInfo == nil {
		return false, fmt.Errorf("blob status is confirmed but missing confirmation info: %s", m.GetBlobKey().String())
	}
	return true, nil
}

type RequestMetadata struct {
	core.BlobRequestHeader
	BlobSize    uint   `json:"blob_size"`
	RequestedAt uint64 `json:"requested_at"`
}

type ConfirmationInfo struct {
	BatchHeaderHash         [32]byte                             `json:"batch_header_hash"`
	BlobIndex               uint32                               `json:"blob_index"`
	BlobCount               uint32                               `json:"blob_count"`
	SignatoryRecordHash     [32]byte                             `json:"signatory_record_hash"`
	ReferenceBlockNumber    uint32                               `json:"reference_block_number"`
	BatchRoot               []byte                               `json:"batch_root"`
	BlobInclusionProof      []byte                               `json:"blob_inclusion_proof"`
	BlobCommitment          *encoding.BlobCommitments            `json:"blob_commitment"`
	BatchID                 uint32                               `json:"batch_id"`
	ConfirmationTxnHash     gcommon.Hash                         `json:"confirmation_txn_hash"`
	ConfirmationBlockNumber uint32                               `json:"confirmation_block_number"`
	Fee                     []byte                               `json:"fee"`
	QuorumResults           map[core.QuorumID]*core.QuorumResult `json:"quorum_results"`
	BlobQuorumInfos         []*core.BlobQuorumInfo               `json:"blob_quorum_infos"`
}

type BlobStoreExclusiveStartKey struct {
	BlobHash     BlobHash
	MetadataHash MetadataHash
	BlobStatus   int32 // BlobStatus is an integer
	Expiry       int64 // Expiry is epoch time in seconds
}

type BatchIndexExclusiveStartKey struct {
	BlobHash        BlobHash
	MetadataHash    MetadataHash
	BatchHeaderHash []byte
	BlobIndex       uint32
}

type BlobStore interface {
	// StoreBlob adds a blob to the queue and returns a key that can be used to retrieve the blob later
	StoreBlob(ctx context.Context, blob *core.Blob, requestedAt uint64) (BlobKey, error)
	// GetBlobContent retrieves a blob's content
	GetBlobContent(ctx context.Context, blobHash BlobHash) ([]byte, error)
	// MarkBlobConfirmed updates blob metadata to Confirmed status with confirmation info
	// Returns the updated metadata and error
	MarkBlobConfirmed(ctx context.Context, existingMetadata *BlobMetadata, confirmationInfo *ConfirmationInfo) (*BlobMetadata, error)
	// MarkBlobDispersing updates blob metadata to Dispersing status
	MarkBlobDispersing(ctx context.Context, blobKey BlobKey) error
	// MarkBlobInsufficientSignatures updates blob metadata to InsufficientSignatures status with confirmation info
	// Returns the updated metadata and error
	MarkBlobInsufficientSignatures(ctx context.Context, existingMetadata *BlobMetadata, confirmationInfo *ConfirmationInfo) (*BlobMetadata, error)
	// MarkBlobFinalized marks a blob as finalized
	MarkBlobFinalized(ctx context.Context, blobKey BlobKey) error
	// MarkBlobProcessing marks a blob as processing
	MarkBlobProcessing(ctx context.Context, blobKey BlobKey) error
	// MarkBlobFailed marks a blob as failed
	MarkBlobFailed(ctx context.Context, blobKey BlobKey) error
	// IncrementBlobRetryCount increments the retry count of a blob
	IncrementBlobRetryCount(ctx context.Context, existingMetadata *BlobMetadata) error
	// UpdateConfirmationBlockNumber updates the confirmation block number of a blob
	UpdateConfirmationBlockNumber(ctx context.Context, existingMetadata *BlobMetadata, confirmationBlockNumber uint32) error
	// GetBlobsByMetadata retrieves a list of blobs given a list of metadata
	GetBlobsByMetadata(ctx context.Context, metadata []*BlobMetadata) (map[BlobKey]*core.Blob, error)
	// GetBlobMetadataByStatus returns a list of blob metadata for blobs with the given status
	GetBlobMetadataByStatus(ctx context.Context, blobStatus BlobStatus) ([]*BlobMetadata, error)
	// GetMetadataInBatch returns the metadata in a given batch at given index.
	GetMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*BlobMetadata, error)
	// GetBlobMetadataByStatusWithPagination returns a list of blob metadata for blobs with the given status
	// Results are limited to the given limit and the pagination token is returned
	GetBlobMetadataByStatusWithPagination(ctx context.Context, blobStatus BlobStatus, limit int32, exclusiveStartKey *BlobStoreExclusiveStartKey) ([]*BlobMetadata, *BlobStoreExclusiveStartKey, error)
	// GetAllBlobMetadataByBatch returns the metadata of all the blobs in the batch.
	GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*BlobMetadata, error)
	// GetAllBlobMetadataByBatchWithPagination returns all the blobs in the batch using pagination
	GetAllBlobMetadataByBatchWithPagination(ctx context.Context, batchHeaderHash [32]byte, limit int32, exclusiveStartKey *BatchIndexExclusiveStartKey) ([]*BlobMetadata, *BatchIndexExclusiveStartKey, error)
	// GetBlobMetadata returns a blob metadata given a metadata key
	GetBlobMetadata(ctx context.Context, blobKey BlobKey) (*BlobMetadata, error)
	// GetBulkBlobMetadata returns a list of blob metadata given a list of blob keys
	GetBulkBlobMetadata(ctx context.Context, blobKeys []BlobKey) ([]*BlobMetadata, error)
	// HandleBlobFailure handles a blob failure by either incrementing the retry count or marking the blob as failed
	// Returns a boolean indicating whether the blob should be retried and an error
	HandleBlobFailure(ctx context.Context, metadata *BlobMetadata, maxRetry uint) (bool, error)
}

type Dispatcher interface {
	DisperseBatch(context.Context, *core.IndexedOperatorState, []core.EncodedBlob, *core.BatchHeader) chan core.SigningMessage
	SendBlobsToOperator(ctx context.Context, blobs []*core.EncodedBlobMessage, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) ([]*core.Signature, error)
	AttestBatch(ctx context.Context, state *core.IndexedOperatorState, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader) (chan core.SigningMessage, error)
	SendAttestBatchRequest(ctx context.Context, nodeDispersalClient node.DispersalClient, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error)
}

// GenerateReverseIndexKey returns the key used to store the blob key in the reverse index
func GenerateReverseIndexKey(batchHeaderHash [32]byte, blobIndex uint32) (string, error) {
	blobIndexHash, err := common.Hash[uint32](blobIndex)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, 0, len(batchHeaderHash)+len(blobIndexHash))
	bytes = append(bytes, batchHeaderHash[:]...)
	bytes = append(bytes, blobIndexHash...)
	return hex.EncodeToString(sha256.New().Sum(bytes)), nil
}

func FromBlobStatusProto(status disperser_rpc.BlobStatus) (*BlobStatus, error) {
	var res BlobStatus
	switch status {
	case disperser_rpc.BlobStatus_UNKNOWN:
		return nil, errors.New("unexpected blob status BlobStatus_UNKNOWN")
	case disperser_rpc.BlobStatus_PROCESSING:
		res = Processing
		return &res, nil
	case disperser_rpc.BlobStatus_CONFIRMED:
		res = Confirmed
		return &res, nil
	case disperser_rpc.BlobStatus_FAILED:
		res = Failed
		return &res, nil
	case disperser_rpc.BlobStatus_FINALIZED:
		res = Finalized
		return &res, nil
	case disperser_rpc.BlobStatus_INSUFFICIENT_SIGNATURES:
		res = InsufficientSignatures
		return &res, nil
	case disperser_rpc.BlobStatus_DISPERSING:
		res = Dispersing
		return &res, nil
	}

	return nil, fmt.Errorf("unknown blob status: %v", status)
}
