package blobstore

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	StatusIndexName            = "StatusIndex"
	OperatorDispersalIndexName = "OperatorDispersalIndex"
	OperatorResponseIndexName  = "OperatorResponseIndex"
	RequestedAtIndexName       = "RequestedAtIndex"
	AttestedAtIndexName        = "AttestedAtAIndex"

	blobKeyPrefix             = "BlobKey#"
	dispersalKeyPrefix        = "Dispersal#"
	batchHeaderKeyPrefix      = "BatchHeader#"
	blobMetadataSK            = "BlobMetadata"
	blobCertSK                = "BlobCertificate"
	dispersalRequestSKPrefix  = "DispersalRequest#"
	dispersalResponseSKPrefix = "DispersalResponse#"
	batchHeaderSK             = "BatchHeader"
	attestationSK             = "Attestation"

	// The number of nanoseconds for a requestedAt bucket (1h).
	// The rationales are:
	// - 1h would be a good estimate for blob feed request (e.g. fetch blobs in past hour can be a common use case)
	// - at 100 blobs/s, it'll be 360,000 blobs in a bucket, which is reasonable
	// - and then it'll be 336 buckets in total (24 buckets/day * 14 days), which is also reasonable
	requestedAtBucketSizeNano = uint64(time.Hour / time.Nanosecond)
	// 14 days with 1 hour margin of safety.
	maxBlobAgeInNano = uint64((14*24*time.Hour + 1*time.Hour) / time.Nanosecond)

	// The number of nanoseconds for an attestedAt bucket (1d)
	// - 1d would be a good estimate for attestation needs (e.g. signing rate over past 24h is a common use case)
	// - even at 1 attesation/s, it'll be 86,400 attestations in a bucket, which is reasonable
	attestedAtBucketSizeNano = uint64(24 * time.Hour / time.Nanosecond)
)

var (
	statusUpdatePrecondition = map[v2.BlobStatus][]v2.BlobStatus{
		v2.Queued:              {},
		v2.Encoded:             {v2.Queued},
		v2.GatheringSignatures: {v2.Encoded},
		// TODO: when GatheringSignatures is fully supported, remove v2.Encoded from below
		v2.Complete: {v2.Encoded, v2.GatheringSignatures},
		v2.Failed:   {v2.Queued, v2.Encoded, v2.GatheringSignatures},
	}
	ErrInvalidStateTransition = errors.New("invalid state transition")
)

type StatusIndexCursor struct {
	BlobKey   *corev2.BlobKey
	UpdatedAt uint64
}

// BlobFeedCursor represents a position in the blob feed, which contains all blobs
// accepted by Disperser, ordered by (requestedAt, blobKey).
type BlobFeedCursor struct {
	RequestedAt uint64

	// The BlobKey can be nil, and a nil BlobKey is treated as equal to another nil BlobKey
	BlobKey *corev2.BlobKey
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

// BlobMetadataStore is a blob metadata storage backed by DynamoDB
type BlobMetadataStore struct {
	dynamoDBClient commondynamodb.Client
	logger         logging.Logger
	tableName      string
}

func NewBlobMetadataStore(dynamoDBClient commondynamodb.Client, logger logging.Logger, tableName string) *BlobMetadataStore {
	logger.Debugf("creating blob metadata store v2 with table %s", tableName)
	return &BlobMetadataStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger.With("component", "blobMetadataStoreV2"),
		tableName:      tableName,
	}
}

func (s *BlobMetadataStore) PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	s.logger.Debug("store put blob metadata", "blobMetadata", blobMetadata)
	item, err := MarshalBlobMetadata(blobMetadata)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) UpdateBlobStatus(ctx context.Context, blobKey corev2.BlobKey, status v2.BlobStatus) error {
	validStatuses := statusUpdatePrecondition[status]
	if len(validStatuses) == 0 {
		return fmt.Errorf("%w: invalid status transition to %s", ErrInvalidStateTransition, status.String())
	}

	expValues := make([]expression.OperandBuilder, len(validStatuses))
	for i, validStatus := range validStatuses {
		expValues[i] = expression.Value(int(validStatus))
	}
	condition := expression.Name("BlobStatus").In(expValues[0], expValues[1:]...)
	_, err := s.dynamoDBClient.UpdateItemWithCondition(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	}, map[string]types.AttributeValue{
		"BlobStatus": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		"UpdatedAt": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().UnixNano(), 10),
		},
	}, condition)

	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		blob, err := s.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			return fmt.Errorf("failed to get blob metadata for key %s: %v", blobKey.Hex(), err)
		}

		if blob.BlobStatus == status {
			return fmt.Errorf("%w: blob already in status %s", common.ErrAlreadyExists, status.String())
		}

		return fmt.Errorf("%w: invalid status transition to %s", ErrInvalidStateTransition, status.String())
	}

	return err
}

func (s *BlobMetadataStore) DeleteBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) error {
	err := s.dynamoDBClient.DeleteItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	})

	return err
}

func (s *BlobMetadataStore) GetBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobMetadata, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	})

	if item == nil {
		return nil, fmt.Errorf("%w: metadata not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	if err != nil {
		return nil, err
	}

	metadata, err := UnmarshalBlobMetadata(item)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetBlobMetadataByStatus returns all the metadata with the given status that were updated after lastUpdatedAt
// Because this function scans the entire index, it should only be used for status with a limited number of items.
// Results are ordered by UpdatedAt in ascending order.
func (s *BlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, StatusIndexName, "BlobStatus = :status AND UpdatedAt > :updatedAt", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":updatedAt": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(int64(lastUpdatedAt), 10),
		}})
	if err != nil {
		return nil, err
	}

	metadata := make([]*v2.BlobMetadata, len(items))
	for i, item := range items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

// queryBucketBlobMetadata returns blobs (as metadata) within range [startKey, endKey] from a single bucket.
// Results are ordered by <RequestedAt, Bobkey> in ascending order.
//
// The function handles DynamoDB's 1MB response size limitation by performing multiple queries if necessary.
func (s *BlobMetadataStore) queryBucketBlobMetadata(
	ctx context.Context,
	bucket uint64,
	startKey string,
	endKey string,
) ([]*v2.BlobMetadata, error) {
	metadata := make([]*v2.BlobMetadata, 0)
	var lastEvaledKey map[string]types.AttributeValue

	for {
		start := startKey
		if lastEvaledKey != nil {
			requestedAtBlobkey, err := UnmarshalRequestedAtBlobKey(lastEvaledKey)
			if err != nil {
				return nil, fmt.Errorf("failed to parse the RequestedAtBlobkey from the LastEvaluatedKey: %w", err)
			}
			start = requestedAtBlobkey
		}

		res, err := s.dynamoDBClient.QueryIndexWithPagination(
			ctx,
			s.tableName,
			RequestedAtIndexName,
			"RequestedAtBucket = :pk AND RequestedAtBlobKey BETWEEN :start AND :end",
			commondynamodb.ExpressionValues{
				":pk":    &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", bucket)},
				":start": &types.AttributeValueMemberS{Value: start},
				":end":   &types.AttributeValueMemberS{Value: endKey},
			},
			0, // no limit within a bucket
			lastEvaledKey,
		)
		if err != nil {
			return nil, fmt.Errorf("query failed for bucket %d: %w", bucket, err)
		}

		// Collect results
		for _, item := range res.Items {
			bm, err := UnmarshalBlobMetadata(item)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal blob metadata: %w", err)
			}
			metadata = append(metadata, bm)
		}

		// Exhausted all items already
		if res.LastEvaluatedKey == nil {
			break
		}
		// For next iteration
		lastEvaledKey = res.LastEvaluatedKey
	}

	return metadata, nil
}

// GetBlobMetadataByRequestedAt returns blobs (as BlobMetadata) that are in cursor position
// range (start, end]. Blobs returned are in cursor order.
//
// If limit > 0, returns at most that many blobs. If limit <= 0, returns all blobs in range.
// Also returns the cursor of the last processed blob, or nil if no blobs were processed.
func (s *BlobMetadataStore) GetBlobMetadataByRequestedAt(
	ctx context.Context,
	start BlobFeedCursor,
	end BlobFeedCursor,
	limit int,
) ([]*v2.BlobMetadata, *BlobFeedCursor, error) {
	if !start.LessThan(&end) {
		return nil, nil, errors.New("start cursor is expected to be less than end cursor")
	}

	startBucket, endBucket := GetRequestedAtBucketIDRange(start.RequestedAt, end.RequestedAt)
	startKey := start.ToCursorKey()
	endKey := end.ToCursorKey()

	result := make([]*v2.BlobMetadata, 0)
	var lastProcessedCursor *BlobFeedCursor

	for bucket := startBucket; bucket <= endBucket; bucket++ {
		if limit > 0 && len(result) >= limit {
			break
		}

		bucketMetadata, err := s.queryBucketBlobMetadata(ctx, bucket, startKey, endKey)
		if err != nil {
			return nil, nil, err
		}

		// Process bucket results
		for _, bm := range bucketMetadata {
			blobKey, err := bm.BlobHeader.BlobKey()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get blob key: %w", err)
			}

			// Skip the start cursor's blob
			if start.Equal(bm.RequestedAt, &blobKey) {
				continue
			}

			result = append(result, bm)
			lastProcessedCursor = &BlobFeedCursor{
				RequestedAt: bm.RequestedAt,
				BlobKey:     &blobKey,
			}

			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	return result, lastProcessedCursor, nil
}

// queryBucketAttestation returns attestations within a single bucket of time range [start, end]. Results are ordered by AttestedAt in
// ascending order.
//
// The function handles DynamoDB's 1MB response size limitation by performing multiple queries  if necessary.
// If there are more than numToReturn attestations in the bucket, returns numToReturn attestations; otherwise returns all attestations in bucket.
func (s *BlobMetadataStore) queryBucketAttestation(ctx context.Context, bucket, start, end uint64, numToReturn int) ([]*corev2.Attestation, error) {
	attestations := make([]*corev2.Attestation, 0)
	var lastEvaledKey map[string]types.AttributeValue

	// Iteratively fetch results from the bucket until we get desired number of items
	// or exhaust the entire bucket.
	// This needs to be processed in a loop because DynamoDb has a limit on the response
	// size of a query (1MB) and we may have more data than that.
	for {
		startTime := start
		if lastEvaledKey != nil {
			at, err := UnmarshalAttestedAt(lastEvaledKey)
			if err != nil {
				return nil, fmt.Errorf("failed to parse the AttestedAt from the LastEvaluatedKey: %w", err)
			}
			startTime = at
		}
		res, err := s.dynamoDBClient.QueryIndexWithPagination(
			ctx,
			s.tableName,
			AttestedAtIndexName,
			"AttestedAtBucket = :pk AND AttestedAt BETWEEN :start AND :end",
			commondynamodb.ExpressionValues{
				":pk":    &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", bucket)},
				":start": &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(startTime), 10)},
				":end":   &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(end), 10)},
			},
			0, // no limit within a bucket
			lastEvaledKey,
		)
		if err != nil {
			return nil, fmt.Errorf("query failed for bucket %d: %w", bucket, err)
		}

		// Collect results
		for _, item := range res.Items {
			at, err := UnmarshalAttestation(item)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal attestation: %w", err)
			}
			attestations = append(attestations, at)

			// Desired number of items collected
			if len(attestations) >= numToReturn {
				return attestations, nil
			}
		}

		// Exhausted all items already
		if res.LastEvaluatedKey == nil {
			break
		}
		// For next iteration
		lastEvaledKey = res.LastEvaluatedKey
	}

	return attestations, nil
}

// GetAttestationByAttestedAt returns attestations within a specified time range (start, end],
// ordered by their AttestedAt timestamp in ascending order.
//
// The function splits the time range into buckets and queries each bucket sequentially.
// Results from all buckets are combined while maintaining the ordering.
//
// If limit > 0, returns at most that many attestations. If limit <= 0, returns all attestations
// in the time range.
func (s *BlobMetadataStore) GetAttestationByAttestedAt(
	ctx context.Context,
	start uint64,
	end uint64,
	limit int,
) ([]*corev2.Attestation, error) {
	if start >= end {
		return nil, errors.New("start must be less than end")
	}

	startBucket, endBucket := GetAttestedAtBucketIDRange(start, end)

	result := make([]*corev2.Attestation, 0)
	for bucket := startBucket; bucket <= endBucket; bucket++ {
		if limit > 0 && len(result) >= limit {
			break
		}

		remaining := math.MaxInt
		if limit > 0 {
			remaining = limit - len(result)
		}
		bucketAttestation, err := s.queryBucketAttestation(ctx, bucket, start+1, end, remaining)
		if err != nil {
			return nil, err
		}

		for _, ba := range bucketAttestation {
			result = append(result, ba)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

// GetBlobMetadataByStatusPaginated returns all the metadata with the given status that were updated after the given cursor.
// It also returns a new cursor (last evaluated key) to be used for the next page
// even when there are no more results or there are no results at all.
// This cursor can be used to get new set of results when they become available.
// Therefore, it's possible to get an empty result from a request with exclusive start key returned from previous response.
func (s *BlobMetadataStore) GetBlobMetadataByStatusPaginated(
	ctx context.Context,
	status v2.BlobStatus,
	exclusiveStartKey *StatusIndexCursor,
	limit int32,
) ([]*v2.BlobMetadata, *StatusIndexCursor, error) {
	var cursor map[string]types.AttributeValue
	if exclusiveStartKey != nil {
		pk := blobKeyPrefix
		if exclusiveStartKey.BlobKey != nil && len(exclusiveStartKey.BlobKey) == 32 {
			pk = blobKeyPrefix + exclusiveStartKey.BlobKey.Hex()
		}
		cursor = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: pk,
			},
			"SK": &types.AttributeValueMemberS{
				Value: blobMetadataSK,
			},
			"UpdatedAt": &types.AttributeValueMemberN{
				Value: strconv.FormatUint(exclusiveStartKey.UpdatedAt, 10),
			},
			"BlobStatus": &types.AttributeValueMemberN{
				Value: strconv.Itoa(int(status)),
			},
		}
	} else {
		cursor = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: blobKeyPrefix,
			},
			"SK": &types.AttributeValueMemberS{
				Value: blobMetadataSK,
			},
			"UpdatedAt": &types.AttributeValueMemberN{
				Value: "0",
			},
			"BlobStatus": &types.AttributeValueMemberN{
				Value: strconv.Itoa(int(status)),
			},
		}
	}
	res, err := s.dynamoDBClient.QueryIndexWithPagination(ctx, s.tableName, StatusIndexName, "BlobStatus = :status", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	}, limit, cursor)
	if err != nil {
		return nil, nil, err
	}

	// No results
	if len(res.Items) == 0 && res.LastEvaluatedKey == nil {
		// return the same cursor
		return nil, exclusiveStartKey, nil
	}

	metadata := make([]*v2.BlobMetadata, 0, len(res.Items))
	for _, item := range res.Items {
		m, err := UnmarshalBlobMetadata(item)
		// Skip invalid/corrupt items
		if err != nil {
			s.logger.Errorf("failed to unmarshal blob metadata: %v", err)
			continue
		}
		metadata = append(metadata, m)
	}

	lastEvaludatedKey := res.LastEvaluatedKey
	if lastEvaludatedKey == nil {
		lastItem := res.Items[len(res.Items)-1]
		u, ok := lastItem["UpdatedAt"].(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, fmt.Errorf("expected *types.AttributeValueMemberN for UpdatedAt, got %T", u)
		}
		updatedAt, err := strconv.ParseUint(u.Value, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		bk, err := UnmarshalBlobKey(lastItem)
		if err != nil {
			return nil, nil, err
		}
		return metadata, &StatusIndexCursor{
			BlobKey:   &bk,
			UpdatedAt: updatedAt,
		}, nil
	}

	newCursor := StatusIndexCursor{}
	err = attributevalue.UnmarshalMap(lastEvaludatedKey, &newCursor)
	if err != nil {
		return nil, nil, err
	}
	bk, err := UnmarshalBlobKey(lastEvaludatedKey)
	if err != nil {
		return nil, nil, err
	}
	newCursor.BlobKey = &bk

	return metadata, &newCursor, nil
}

// GetBlobMetadataCountByStatus returns the count of all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
func (s *BlobMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error) {
	count, err := s.dynamoDBClient.QueryIndexCount(ctx, s.tableName, StatusIndexName, "BlobStatus = :status", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BlobMetadataStore) PutBlobCertificate(ctx context.Context, blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) error {
	item, err := MarshalBlobCertificate(blobCert, fragmentInfo)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) DeleteBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) error {
	err := s.dynamoDBClient.DeleteItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobCertSK,
		},
	})

	return err
}

func (s *BlobMetadataStore) GetBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobCertSK,
		},
	})

	if err != nil {
		return nil, nil, err
	}

	if item == nil {
		return nil, nil, fmt.Errorf("%w: certificate not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	cert, fragmentInfo, err := UnmarshalBlobCertificate(item)
	if err != nil {
		return nil, nil, err
	}

	return cert, fragmentInfo, nil
}

// GetBlobCertificates returns the certificates for the given blob keys
// Note: the returned certificates are NOT necessarily ordered by the order of the input blob keys
func (s *BlobMetadataStore) GetBlobCertificates(ctx context.Context, blobKeys []corev2.BlobKey) ([]*corev2.BlobCertificate, []*encoding.FragmentInfo, error) {
	keys := make([]map[string]types.AttributeValue, len(blobKeys))
	for i, blobKey := range blobKeys {
		keys[i] = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: blobKeyPrefix + blobKey.Hex(),
			},
			"SK": &types.AttributeValueMemberS{
				Value: blobCertSK,
			},
		}
	}

	items, err := s.dynamoDBClient.GetItems(ctx, s.tableName, keys)
	if err != nil {
		return nil, nil, err
	}

	certs := make([]*corev2.BlobCertificate, len(items))
	fragmentInfos := make([]*encoding.FragmentInfo, len(items))
	for i, item := range items {
		cert, fragmentInfo, err := UnmarshalBlobCertificate(item)
		if err != nil {
			return nil, nil, err
		}
		certs[i] = cert
		fragmentInfos[i] = fragmentInfo
	}

	return certs, fragmentInfos, nil
}

func (s *BlobMetadataStore) PutDispersalRequest(ctx context.Context, req *corev2.DispersalRequest) error {
	item, err := MarshalDispersalRequest(req)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetDispersalRequest(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalRequest, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: dispersalKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s%s", dispersalRequestSKPrefix, operatorID.Hex()),
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: dispersal request not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
	}

	req, err := UnmarshalDispersalRequest(item)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *BlobMetadataStore) PutDispersalResponse(ctx context.Context, res *corev2.DispersalResponse) error {
	item, err := MarshalDispersalResponse(res)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetDispersalResponse(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalResponse, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: dispersalKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s%s", dispersalResponseSKPrefix, operatorID.Hex()),
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: dispersal response not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
	}

	res, err := UnmarshalDispersalResponse(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *BlobMetadataStore) GetDispersalResponses(ctx context.Context, batchHeaderHash [32]byte) ([]*corev2.DispersalResponse, error) {
	items, err := s.dynamoDBClient.Query(ctx, s.tableName, "PK = :pk AND begins_with(SK, :prefix)", commondynamodb.ExpressionValues{
		":pk": &types.AttributeValueMemberS{
			Value: dispersalKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: dispersalResponseSKPrefix,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("%w: dispersal responses not found for batch header hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	responses := make([]*corev2.DispersalResponse, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersalResponse(item)
		if err != nil {
			return nil, err
		}
	}

	return responses, nil
}

func (s *BlobMetadataStore) PutBatchHeader(ctx context.Context, batchHeader *corev2.BatchHeader) error {
	item, err := MarshalBatchHeader(batchHeader)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) DeleteBatchHeader(ctx context.Context, batchHeaderHash [32]byte) error {
	err := s.dynamoDBClient.DeleteItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchHeaderSK,
		},
	})

	return err
}

func (s *BlobMetadataStore) GetBatchHeader(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchHeaderSK,
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: batch header not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	header, err := UnmarshalBatchHeader(item)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (s *BlobMetadataStore) PutAttestation(ctx context.Context, attestation *corev2.Attestation) error {
	item, err := MarshalAttestation(attestation)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetAttestation(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Attestation, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: attestationSK,
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: attestation not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	attestation, err := UnmarshalAttestation(item)
	if err != nil {
		return nil, err
	}

	return attestation, nil
}

func (s *BlobMetadataStore) PutBlobInclusionInfo(ctx context.Context, inclusionInfo *corev2.BlobInclusionInfo) error {
	item, err := MarshalBlobInclusionInfo(inclusionInfo)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

// PutBlobInclusionInfos puts multiple inclusion infos into the store
// It retries failed items up to 2 times
func (s *BlobMetadataStore) PutBlobInclusionInfos(ctx context.Context, inclusionInfos []*corev2.BlobInclusionInfo) error {
	items := make([]commondynamodb.Item, len(inclusionInfos))
	for i, info := range inclusionInfos {
		item, err := MarshalBlobInclusionInfo(info)
		if err != nil {
			return err
		}
		items[i] = item
	}

	numRetries := 3
	for i := 0; i < numRetries; i++ {
		failedItems, err := s.dynamoDBClient.PutItems(ctx, s.tableName, items)
		if err != nil {
			return err
		}

		if len(failedItems) > 0 {
			s.logger.Warnf("failed to put inclusion infos, retrying: %v", failedItems)
			items = failedItems
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
		} else {
			return nil
		}
	}

	return nil
}

func (s *BlobMetadataStore) GetBlobInclusionInfo(ctx context.Context, blobKey corev2.BlobKey, batchHeaderHash [32]byte) (*corev2.BlobInclusionInfo, error) {
	bhh := hex.EncodeToString(batchHeaderHash[:])
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix + bhh,
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: inclusion info not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	info, err := UnmarshalBlobInclusionInfo(item)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *BlobMetadataStore) GetBlobAttestationInfo(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobAttestationInfo, error) {
	blobInclusionInfos, err := s.GetBlobInclusionInfos(ctx, blobKey)
	if err != nil {
		s.logger.Error("failed to get blob inclusion info for blob", "err", err, "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob inclusion info: %s", err.Error()))
	}

	if len(blobInclusionInfos) == 0 {
		s.logger.Error("no blob inclusion info found for blob", "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal("no blob inclusion info found")
	}

	if len(blobInclusionInfos) > 1 {
		s.logger.Warn("multiple inclusion info found for blob", "blobKey", blobKey.Hex())
	}

	for _, inclusionInfo := range blobInclusionInfos {
		// get the signed batch from this inclusion info
		batchHeaderHash, err := inclusionInfo.BatchHeader.Hash()
		if err != nil {
			s.logger.Error("failed to get batch header hash from blob inclusion info", "err", err, "blobKey", blobKey.Hex())
			continue
		}
		_, attestation, err := s.GetSignedBatch(ctx, batchHeaderHash)
		if err != nil {
			s.logger.Error("failed to get signed batch", "err", err, "blobKey", blobKey.Hex())
			continue
		}

		return &v2.BlobAttestationInfo{
			InclusionInfo: inclusionInfo,
			Attestation:   attestation,
		}, nil
	}

	return nil, fmt.Errorf("no attestation info found for blobkey: %s", blobKey.Hex())
}

func (s *BlobMetadataStore) GetBlobInclusionInfos(ctx context.Context, blobKey corev2.BlobKey) ([]*corev2.BlobInclusionInfo, error) {
	items, err := s.dynamoDBClient.Query(ctx, s.tableName, "PK = :pk AND begins_with(SK, :prefix)", commondynamodb.ExpressionValues{
		":pk": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("%w: inclusion info not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	responses := make([]*corev2.BlobInclusionInfo, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalBlobInclusionInfo(item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal inclusion info: %w", err)
		}
	}

	return responses, nil
}

func (s *BlobMetadataStore) GetSignedBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, *corev2.Attestation, error) {
	items, err := s.dynamoDBClient.Query(ctx, s.tableName, "PK = :pk", commondynamodb.ExpressionValues{
		":pk": &types.AttributeValueMemberS{
			Value: batchHeaderKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
	})

	if err != nil {
		return nil, nil, err
	}

	if len(items) == 0 {
		return nil, nil, fmt.Errorf("%w: no records found for batch header hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	var header *corev2.BatchHeader
	var attestation *corev2.Attestation
	for _, item := range items {
		sk, ok := item["SK"].(*types.AttributeValueMemberS)
		if !ok {
			return nil, nil, fmt.Errorf("expected *types.AttributeValueMemberS for SK, got %T", item["SK"])
		}
		if strings.HasPrefix(sk.Value, batchHeaderSK) {
			header, err = UnmarshalBatchHeader(item)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal batch header: %w", err)
			}
		} else if strings.HasPrefix(sk.Value, attestationSK) {
			attestation, err = UnmarshalAttestation(item)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal attestation: %w", err)
			}
		}
	}

	if header == nil {
		return nil, nil, fmt.Errorf("%w: batch header not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	if attestation == nil {
		return nil, nil, fmt.Errorf("%w: attestation not found for hash %x", common.ErrMetadataNotFound, batchHeaderHash)
	}

	return header, attestation, nil
}

func GenerateTableSchema(tableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			// PK is the composite partition key
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			// SK is the composite sort key
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BlobStatus"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("UpdatedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("OperatorID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("DispersedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("RespondedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("RequestedAtBucket"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("RequestedAtBlobKey"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("AttestedAtBucket"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("AttestedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(tableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(StatusIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BlobStatus"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("UpdatedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
			{
				IndexName: aws.String(OperatorDispersalIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("DispersedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
			{
				IndexName: aws.String(OperatorResponseIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RespondedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
			{
				IndexName: aws.String(RequestedAtIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("RequestedAtBucket"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RequestedAtBlobKey"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
			{
				IndexName: aws.String(AttestedAtIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("AttestedAtBucket"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("AttestedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		},
	}
}

func MarshalBlobMetadata(metadata *v2.BlobMetadata) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob metadata: %w", err)
	}

	// Add PK and SK fields
	blobKey, err := metadata.BlobHeader.BlobKey()
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + blobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobMetadataSK}
	fields["RequestedAtBucket"] = &types.AttributeValueMemberS{Value: computeRequestedAtBucket(metadata.RequestedAt)}
	fields["RequestedAtBlobKey"] = &types.AttributeValueMemberS{Value: encodeBlobFeedCursorKey(metadata.RequestedAt, &blobKey)}
	return fields, nil
}

func UnmarshalBlobKey(item commondynamodb.Item) (corev2.BlobKey, error) {
	type Blob struct {
		PK string
	}

	blob := Blob{}
	err := attributevalue.UnmarshalMap(item, &blob)
	if err != nil {
		return corev2.BlobKey{}, err
	}

	bk := strings.TrimPrefix(blob.PK, blobKeyPrefix)
	return corev2.HexToBlobKey(bk)
}

func UnmarshalBlobMetadata(item commondynamodb.Item) (*v2.BlobMetadata, error) {
	metadata := v2.BlobMetadata{}
	err := attributevalue.UnmarshalMap(item, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func MarshalBlobCertificate(blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(blobCert)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob certificate: %w", err)
	}

	// merge fragment info
	fragmentInfoFields, err := attributevalue.MarshalMap(fragmentInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fragment info: %w", err)
	}
	for k, v := range fragmentInfoFields {
		fields[k] = v
	}

	// Add PK and SK fields
	blobKey, err := blobCert.BlobHeader.BlobKey()
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + blobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobCertSK}

	return fields, nil
}

func UnmarshalBlobCertificate(item commondynamodb.Item) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	cert := corev2.BlobCertificate{}
	err := attributevalue.UnmarshalMap(item, &cert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal blob certificate: %w", err)
	}
	fragmentInfo := encoding.FragmentInfo{}
	err = attributevalue.UnmarshalMap(item, &fragmentInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal fragment info: %w", err)
	}
	return &cert, &fragmentInfo, nil
}

func UnmarshalBatchHeaderHash(item commondynamodb.Item) ([32]byte, error) {
	type Object struct {
		PK string
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return [32]byte{}, err
	}

	root := strings.TrimPrefix(obj.PK, dispersalKeyPrefix)
	return hexToHash(root)
}

func UnmarshalRequestedAtBlobKey(item commondynamodb.Item) (string, error) {
	type Object struct {
		RequestedAtBlobKey string
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return "", err
	}

	return obj.RequestedAtBlobKey, nil
}

func UnmarshalAttestedAt(item commondynamodb.Item) (uint64, error) {
	type Object struct {
		AttestedAt uint64
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return 0, err
	}

	return obj.AttestedAt, nil
}

func UnmarshalOperatorID(item commondynamodb.Item) (*core.OperatorID, error) {
	type Object struct {
		OperatorID string
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return nil, err
	}

	operatorID, err := core.OperatorIDFromHex(obj.OperatorID)
	if err != nil {
		return nil, err
	}

	return &operatorID, nil
}

func MarshalDispersalRequest(req *corev2.DispersalRequest) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dispersal request: %w", err)
	}

	batchHeaderHash, err := req.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(batchHeaderHash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: dispersalKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", dispersalRequestSKPrefix, req.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: req.OperatorID.Hex()}

	return fields, nil
}

func UnmarshalDispersalRequest(item commondynamodb.Item) (*corev2.DispersalRequest, error) {
	req := corev2.DispersalRequest{}
	err := attributevalue.UnmarshalMap(item, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal request: %w", err)
	}

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	req.OperatorID = *operatorID

	return &req, nil
}

func MarshalDispersalResponse(res *corev2.DispersalResponse) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(res)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dispersal response: %w", err)
	}

	batchHeaderHash, err := res.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(batchHeaderHash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: dispersalKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", dispersalResponseSKPrefix, res.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: res.OperatorID.Hex()}

	return fields, nil
}

func UnmarshalDispersalResponse(item commondynamodb.Item) (*corev2.DispersalResponse, error) {
	res := corev2.DispersalResponse{}
	err := attributevalue.UnmarshalMap(item, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal response: %w", err)
	}

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	res.OperatorID = *operatorID

	return &res, nil
}

func MarshalBatchHeader(batchHeader *corev2.BatchHeader) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(batchHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch header: %w", err)
	}

	hash, err := batchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(hash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: batchHeaderKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: batchHeaderSK}

	return fields, nil
}

func UnmarshalBatchHeader(item commondynamodb.Item) (*corev2.BatchHeader, error) {
	header := corev2.BatchHeader{}
	err := attributevalue.UnmarshalMap(item, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch header: %w", err)
	}

	return &header, nil
}

func MarshalBlobInclusionInfo(inclusionInfo *corev2.BlobInclusionInfo) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(inclusionInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob inclusion info: %w", err)
	}

	bhh, err := inclusionInfo.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(bhh[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + inclusionInfo.BlobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: batchHeaderKeyPrefix + hashstr}

	return fields, nil
}

func UnmarshalBlobInclusionInfo(item commondynamodb.Item) (*corev2.BlobInclusionInfo, error) {
	inclusionInfo := corev2.BlobInclusionInfo{}
	err := attributevalue.UnmarshalMap(item, &inclusionInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal blob inclusion info: %w", err)
	}

	return &inclusionInfo, nil
}

func MarshalAttestation(attestation *corev2.Attestation) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal attestation: %w", err)
	}

	hash, err := attestation.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(hash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: batchHeaderKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: attestationSK}
	fields["AttestedAtBucket"] = &types.AttributeValueMemberS{Value: computeAttestedAtBucket(attestation.AttestedAt)}
	return fields, nil
}

func UnmarshalAttestation(item commondynamodb.Item) (*corev2.Attestation, error) {
	attestation := corev2.Attestation{}
	err := attributevalue.UnmarshalMap(item, &attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal attestation: %w", err)
	}

	return &attestation, nil
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
