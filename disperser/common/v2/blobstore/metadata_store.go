package blobstore

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// BlobFeedCursor represents a position in the blob feed, which contains all blobs
// accepted by Disperser, ordered by (requestedAt, blobKey).
type BlobFeedCursor struct {
	RequestedAt uint64

	// The BlobKey can be nil, and a nil BlobKey is treated as equal to another nil BlobKey
	BlobKey *corev2.BlobKey
}

// StatusIndexCursor represents a cursor for paginated queries by blob status
type StatusIndexCursor struct {
	BlobKey   *corev2.BlobKey
	UpdatedAt uint64
}

// MetadataStore defines the interface for a blob metadata storage system
type MetadataStore interface {
	// Blob Metadata Operations
	// These methods manage the core blob metadata in the system
	CheckBlobExists(ctx context.Context, blobKey corev2.BlobKey) (bool, error)
	GetBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobMetadata, error)
	PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error
	UpdateBlobStatus(ctx context.Context, key corev2.BlobKey, status v2.BlobStatus) error
	DeleteBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) error // Only used in testing

	// Blob Query Operations
	// These methods provide various ways to query blobs based on different criteria
	GetBlobMetadataByAccountID(
		ctx context.Context,
		accountId gethcommon.Address,
		start uint64,
		end uint64,
		limit int,
		ascending bool,
	) ([]*v2.BlobMetadata, error)
	GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error)
	GetBlobMetadataByStatusPaginated(
		ctx context.Context,
		status v2.BlobStatus,
		exclusiveStartKey *StatusIndexCursor,
		limit int32,
	) ([]*v2.BlobMetadata, *StatusIndexCursor, error)
	GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error)

	// Blob Feed Operations
	// These methods support retrieving blobs in chronological order for feed-like functionality
	GetBlobMetadataByRequestedAtForward(
		ctx context.Context,
		after BlobFeedCursor,
		before BlobFeedCursor,
		limit int,
	) ([]*v2.BlobMetadata, *BlobFeedCursor, error)
	GetBlobMetadataByRequestedAtBackward(
		ctx context.Context,
		before BlobFeedCursor,
		after BlobFeedCursor,
		limit int,
	) ([]*v2.BlobMetadata, *BlobFeedCursor, error)

	// Blob Certificate Operations
	// These methods handle blob certificates which contain cryptographic proofs
	PutBlobCertificate(ctx context.Context, blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) error
	DeleteBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) error
	GetBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) (*corev2.BlobCertificate, *encoding.FragmentInfo, error)
	GetBlobCertificates(ctx context.Context, blobKeys []corev2.BlobKey) ([]*corev2.BlobCertificate, []*encoding.FragmentInfo, error)

	// Batch Operations
	// These methods manage batches of blobs that are processed together
	PutBatch(ctx context.Context, batch *corev2.Batch) error
	GetBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Batch, error)
	PutBatchHeader(ctx context.Context, batchHeader *corev2.BatchHeader) error
	DeleteBatchHeader(ctx context.Context, batchHeaderHash [32]byte) error
	GetBatchHeader(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, error)

	// Dispersal Operations
	// These methods handle the distribution of blobs to operators
	PutDispersalRequest(ctx context.Context, req *corev2.DispersalRequest) error
	GetDispersalRequest(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalRequest, error)
	PutDispersalResponse(ctx context.Context, res *corev2.DispersalResponse) error
	GetDispersalResponse(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalResponse, error)
	GetDispersalResponses(ctx context.Context, batchHeaderHash [32]byte) ([]*corev2.DispersalResponse, error)
	GetDispersalsByRespondedAt(
		ctx context.Context,
		operatorId core.OperatorID,
		start uint64,
		end uint64,
		limit int,
		ascending bool,
	) ([]*corev2.DispersalResponse, error)

	// Attestation Operations
	// These methods handle cryptographic attestations of batches
	PutAttestation(ctx context.Context, attestation *corev2.Attestation) error
	GetAttestation(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Attestation, error)
	GetAttestationByAttestedAtForward(
		ctx context.Context,
		after uint64,
		before uint64,
		limit int,
	) ([]*corev2.Attestation, error)
	GetAttestationByAttestedAtBackward(
		ctx context.Context,
		before uint64,
		after uint64,
		limit int,
	) ([]*corev2.Attestation, error)

	// Blob Inclusion Operations
	// These methods handle information about blob inclusion in batches
	PutBlobInclusionInfo(ctx context.Context, inclusionInfo *corev2.BlobInclusionInfo) error
	PutBlobInclusionInfos(ctx context.Context, inclusionInfos []*corev2.BlobInclusionInfo) error
	GetBlobInclusionInfo(ctx context.Context, blobKey corev2.BlobKey, batchHeaderHash [32]byte) (*corev2.BlobInclusionInfo, error)
	GetBlobInclusionInfos(ctx context.Context, blobKey corev2.BlobKey) ([]*corev2.BlobInclusionInfo, error)
	GetBlobAttestationInfo(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobAttestationInfo, error)

	// Combined Operations
	// These methods provide convenient access to related data in a single call
	GetSignedBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, *corev2.Attestation, error)
}

// Equal returns true if the cursor is equal to the given <requestedAt, blobKey>
func (cursor *BlobFeedCursor) Equal(requestedAt uint64, blobKey *corev2.BlobKey) bool {
	if cursor.RequestedAt != requestedAt {
		return false
	}

	// Both nil
	if cursor.BlobKey == nil && blobKey == nil {
		return true
	}

	// One nil
	if cursor.BlobKey == nil || blobKey == nil {
		return false
	}

	return cursor.BlobKey.Hex() == blobKey.Hex()
}

// LessThan returns true if the current cursor is less than the other cursor
// in the ordering defined by (requestedAt, blobKey).
func (cursor *BlobFeedCursor) LessThan(other *BlobFeedCursor) bool {
	if other == nil {
		return false
	}

	// First, compare the RequestedAt timestamps
	if cursor.RequestedAt != other.RequestedAt {
		return cursor.RequestedAt < other.RequestedAt
	}

	// If RequestedAt is the same, compare BlobKey
	if cursor.BlobKey != nil && other.BlobKey != nil {
		return cursor.BlobKey.Hex() < other.BlobKey.Hex()
	}

	// Handle cases where BlobKey might be nil
	if cursor.BlobKey == nil && other.BlobKey != nil {
		return true // cursor.BlobKey is nil, so it comes first
	}
	if cursor.BlobKey != nil && other.BlobKey == nil {
		return false // other.BlobKey is nil, so "other" comes first
	}

	// If both RequestedAt and BlobKey are equal, return false (because they are equal)
	return false
}

// ToCursorKey encodes the cursor into a string that preserves ordering.
// For any two cursors A and B:
// - A < B if and only if A.ToCursorKey() < B.ToCursorKey()
// - A == B if and only if A.ToCursorKey() == B.ToCursorKey()
func (cursor *BlobFeedCursor) ToCursorKey() string {
	return encodeBlobFeedCursorKey(cursor.RequestedAt, cursor.BlobKey)
}

// FromCursorKey decodes the cursor key string back to the cursor.
func (cursor *BlobFeedCursor) FromCursorKey(encoded string) (*BlobFeedCursor, error) {
	requestedAt, blobKey, err := decodeBlobFeedCursorKey(encoded)
	if err != nil {
		return nil, err
	}
	return &BlobFeedCursor{
		RequestedAt: requestedAt,
		BlobKey:     blobKey,
	}, nil
}

// GetRequestedAtBucketIDRange returns the adjusted start and end bucket IDs based on
// the allowed time range for blobs.
func GetRequestedAtBucketIDRange(startTime, endTime uint64) (uint64, uint64) {
	now := uint64(time.Now().UnixNano())
	oldestAllowed := now - maxBlobAgeInNano

	startBucket := computeBucketID(startTime, requestedAtBucketSizeNano)
	if startTime < oldestAllowed {
		startBucket = computeBucketID(oldestAllowed, requestedAtBucketSizeNano)
	}

	endBucket := computeBucketID(endTime, requestedAtBucketSizeNano)
	if endTime > now {
		endBucket = computeBucketID(now, requestedAtBucketSizeNano)
	}

	return startBucket, endBucket
}

// GetAttestedAtBucketIDRange returns the adjusted start and end bucket IDs based on
// the allowed time range for blobs.
func GetAttestedAtBucketIDRange(startTime, endTime uint64) (uint64, uint64) {
	now := uint64(time.Now().UnixNano())
	oldestAllowed := now - maxBlobAgeInNano

	startBucket := computeBucketID(startTime, attestedAtBucketSizeNano)
	if startTime < oldestAllowed {
		startBucket = computeBucketID(oldestAllowed, attestedAtBucketSizeNano)
	}

	endBucket := computeBucketID(endTime, attestedAtBucketSizeNano)
	if endTime > now {
		endBucket = computeBucketID(now, attestedAtBucketSizeNano)
	}

	return startBucket, endBucket
}

// encodeBlobFeedCursorKey encodes <requestedAt, blobKey> into string which
// preserves the order.
func encodeBlobFeedCursorKey(requestedAt uint64, blobKey *corev2.BlobKey) string {
	result := make([]byte, 40) // 8 bytes for timestamp + 32 bytes for blobKey

	// Write timestamp
	binary.BigEndian.PutUint64(result[:8], requestedAt)

	if blobKey != nil {
		copy(result[8:], blobKey[:])
	}
	// Use hex encoding to preserve byte ordering
	return hex.EncodeToString(result)
}

// decodeBlobFeedCursorKey decodes the cursor key back to <requestedAt, blobKey>.
func decodeBlobFeedCursorKey(encoded string) (uint64, *corev2.BlobKey, error) {
	// Decode hex string
	bytes, err := hex.DecodeString(encoded)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid hex encoding: %w", err)
	}

	// Check length
	if len(bytes) != 40 { // 8 bytes timestamp + 32 bytes blobKey
		return 0, nil, fmt.Errorf("invalid length: expected 40 bytes, got %d", len(bytes))
	}

	// Get timestamp
	requestedAt := binary.BigEndian.Uint64(bytes[:8])

	// Check if the remaining bytes are all zeros
	allZeros := true
	for i := 8; i < len(bytes); i++ {
		if bytes[i] != 0 {
			allZeros = false
			break
		}
	}

	if allZeros {
		return requestedAt, nil, nil
	}
	var bk corev2.BlobKey
	copy(bk[:], bytes[8:])
	return requestedAt, &bk, nil
}

func hexToHash(h string) ([32]byte, error) {
	s := strings.TrimPrefix(h, "0x")
	s = strings.TrimPrefix(s, "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	return [32]byte(b), nil
}

// computeBucketID maps a given timestamp to a time bucket.
// Note each bucket represents a time range [start, end) (i.e. inclusive start, exclusive end).
func computeBucketID(timestamp, bucketSizeNano uint64) uint64 {
	return timestamp / bucketSizeNano
}

func computeRequestedAtBucket(requestedAt uint64) string {
	id := computeBucketID(requestedAt, requestedAtBucketSizeNano)
	return fmt.Sprintf("%d", id)
}

func computeAttestedAtBucket(attestedAt uint64) string {
	id := computeBucketID(attestedAt, attestedAtBucketSizeNano)
	return fmt.Sprintf("%d", id)
}
