package blobstore_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func checkBlobKeyEqual(t *testing.T, blobKey corev2.BlobKey, blobHeader *corev2.BlobHeader) {
	bk, err := blobHeader.BlobKey()
	assert.Nil(t, err)
	assert.Equal(t, blobKey, bk)
}

func TestBlobFeedCursor_Equal(t *testing.T) {
	bk1 := corev2.BlobKey([32]byte{1, 2, 3})
	bk2 := corev2.BlobKey([32]byte{2, 3, 4})
	tests := []struct {
		cursor      *blobstore.BlobFeedCursor
		requestedAt uint64
		blobKey     *corev2.BlobKey
		expected    bool
	}{
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			requestedAt: 1,
			blobKey:     &bk1,
			expected:    true,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: nil},
			requestedAt: 1,
			blobKey:     nil,
			expected:    true,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			requestedAt: 2,
			blobKey:     &bk1,
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			requestedAt: 1,
			blobKey:     nil,
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: nil},
			requestedAt: 1,
			blobKey:     &bk1,
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			requestedAt: 1,
			blobKey:     &bk2,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run("Equal", func(t *testing.T) {
			result := tt.cursor.Equal(tt.requestedAt, tt.blobKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBlobFeedCursor_LessThan(t *testing.T) {
	bk1 := corev2.BlobKey([32]byte{1, 2, 3})
	bk2 := corev2.BlobKey([32]byte{2, 3, 4})
	tests := []struct {
		cursor      *blobstore.BlobFeedCursor
		otherCursor *blobstore.BlobFeedCursor
		expected    bool
	}{
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 2, BlobKey: &bk1},
			expected:    true,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 2, BlobKey: &bk1},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: nil},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			expected:    true,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: nil},
			expected:    false,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk2},
			expected:    true,
		},
		{
			cursor:      &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk2},
			otherCursor: &blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk1},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run("LessThan", func(t *testing.T) {
			result := tt.cursor.LessThan(tt.otherCursor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBlobFeedCursor_CursorKeyCodec(t *testing.T) {
	bk := corev2.BlobKey([32]byte{1, 2, 3})
	cursors := []*blobstore.BlobFeedCursor{
		&blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: nil},
		&blobstore.BlobFeedCursor{RequestedAt: 1, BlobKey: &bk},
	}
	for _, cursor := range cursors {
		encoded := cursor.ToCursorKey()
		c, err := new(blobstore.BlobFeedCursor).FromCursorKey(encoded)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), c.RequestedAt)
		assert.Equal(t, cursor.BlobKey, c.BlobKey)
	}
}

func TestBlobFeedCursor_OrderPreserving(t *testing.T) {
	bk1 := corev2.BlobKey([32]byte{1, 2, 3})
	bk2 := corev2.BlobKey([32]byte{2, 3, 4})
	cursors := []*blobstore.BlobFeedCursor{
		{RequestedAt: 100, BlobKey: nil},
		{RequestedAt: 100, BlobKey: &bk1},
		{RequestedAt: 100, BlobKey: &bk2},
		{RequestedAt: 101, BlobKey: nil},
		{RequestedAt: 101, BlobKey: &bk1},
	}

	// Test that ordering is consistent between LessThan and ToCursorKey
	for i := 0; i < len(cursors); i++ {
		for j := 0; j < len(cursors); j++ {
			if i != j {
				cursorLessThan := cursors[i].LessThan(cursors[j])
				encodedLessThan := cursors[i].ToCursorKey() < cursors[j].ToCursorKey()
				assert.Equal(t, encodedLessThan, cursorLessThan)
			}
		}
	}
}

func TestBlobMetadataStoreOperations(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t)
	blobKey2, blobHeader2 := newBlob(t)
	now := time.Now()
	metadata1 := &v2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	metadata2 := &v2.BlobMetadata{
		BlobHeader: blobHeader2,
		BlobStatus: v2.Certified,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	queued, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, 0)
	assert.NoError(t, err)
	assert.Len(t, queued, 1)
	assert.Equal(t, metadata1, queued[0])
	// query to get newer blobs should result in 0 results
	queued, err = blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, metadata1.UpdatedAt+100)
	assert.NoError(t, err)
	assert.Len(t, queued, 0)

	certified, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Certified, 0)
	assert.NoError(t, err)
	assert.Len(t, certified, 1)
	assert.Equal(t, metadata2, certified[0])

	queuedCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, v2.Queued)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), queuedCount)

	// attempt to put metadata with the same key should fail
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey1.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey2.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func TestBlobMetadataStoreGetBlobMetadataByRequestedAtWithIdenticalTimestamp(t *testing.T) {
	ctx := context.Background()
	now := uint64(time.Now().UnixNano())
	firstBlobTime := now - uint64(time.Hour.Nanoseconds())
	numBlobs := 5
	dynamoKeys := make([]commondynamodb.Key, numBlobs)

	// Create blobs: first 3 blobs have the same requestedAt, and last 2 blobs have the same requestedAt
	for i := 0; i < numBlobs; i++ {
		blobKey, blobHeader := newBlob(t)
		requestedAt := firstBlobTime
		if i >= 3 {
			requestedAt += 1
		}
		metadata := &v2.BlobMetadata{
			BlobHeader:  blobHeader,
			BlobStatus:  v2.Encoded,
			Expiry:      uint64(time.Now().Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   now,
			RequestedAt: requestedAt,
		}

		err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	keys := make([]corev2.BlobKey, numBlobs)
	requestedAts := make([]uint64, numBlobs)

	// Test blobs are returned in cursor order, i.e. <requestedAt, blobKey>
	startCursor := blobstore.BlobFeedCursor{
		RequestedAt: firstBlobTime - 1,
		BlobKey:     nil,
	}
	endCursor := blobstore.BlobFeedCursor{
		RequestedAt: now,
		BlobKey:     nil,
	}

	metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
	require.NoError(t, err)
	assert.Equal(t, len(metadata), 5)
	require.NotNil(t, lastProcessedCursor)

	// Verify ordering
	for i := 0; i < len(metadata); i++ {
		keys[i], err = metadata[i].BlobHeader.BlobKey()
		require.NoError(t, err)
		requestedAts[i] = metadata[i].RequestedAt
		if i > 0 {
			if metadata[i].RequestedAt != metadata[i-1].RequestedAt {
				assert.True(t, metadata[i].RequestedAt > metadata[i-1].RequestedAt)
			} else {
				assert.True(t, keys[i].Hex() > keys[i-1].Hex())
			}
		}
	}

	// The first 3 blobs have same requestedAt
	assert.Equal(t, requestedAts[0], requestedAts[1])
	assert.Equal(t, requestedAts[0], requestedAts[2])
	// The last 2 blobs have same requestedAt
	assert.Equal(t, requestedAts[3], requestedAts[4])

	// Test iteration from the middle of same-timestamp blobs
	startCursor = blobstore.BlobFeedCursor{
		RequestedAt: requestedAts[1],
		BlobKey:     &keys[1],
	}
	endCursor = blobstore.BlobFeedCursor{
		RequestedAt: requestedAts[3],
		BlobKey:     nil,
	}

	// Test with different end cursors
	testCases := []struct {
		endBlobKey *corev2.BlobKey
		expectLen  int
		expectLast int
	}{
		{nil, 1, 2},
		{&keys[3], 2, 3},
		{&keys[4], 3, 4},
	}

	for _, tc := range testCases {
		endCursor.BlobKey = tc.endBlobKey
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, tc.expectLen, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[tc.expectLast], *lastProcessedCursor.BlobKey)

		// Verify first blob is always keys[2]
		checkBlobKeyEqual(t, keys[2], metadata[0].BlobHeader)

		// Verify remaining blobs if present
		for i := 1; i < len(metadata); i++ {
			checkBlobKeyEqual(t, keys[i+2], metadata[i].BlobHeader)
		}
	}
}

func TestBlobMetadataStoreGetBlobMetadataByRequestedAt(t *testing.T) {
	ctx := context.Background()
	numBlobs := 103
	now := uint64(time.Now().UnixNano())
	firstBlobTime := now - uint64(24*time.Hour.Nanoseconds())
	nanoSecsPerBlob := uint64(60 * 1e9) // 1 blob per minute

	// Create blobs for testing
	keys := make([]corev2.BlobKey, numBlobs)
	dynamoKeys := make([]commondynamodb.Key, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobKey, blobHeader := newBlob(t)
		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader:  blobHeader,
			BlobStatus:  v2.Encoded,
			Expiry:      uint64(now.Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   uint64(now.UnixNano()),
			RequestedAt: firstBlobTime + nanoSecsPerBlob*uint64(i),
		}

		err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		keys[i] = blobKey
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	// Test empty range
	t.Run("empty range", func(t *testing.T) {
		startCursor := blobstore.BlobFeedCursor{
			RequestedAt: now,
			BlobKey:     nil,
		}
		endCursor := blobstore.BlobFeedCursor{
			RequestedAt: now + 10*1e9,
			BlobKey:     nil,
		}

		// Test equal cursors error
		_, _, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, startCursor, 10)
		assert.Error(t, err)
		assert.Equal(t, "start cursor is expected to be less than end cursor", err.Error())

		// Test empty range
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)
	})

	// Test full range query
	t.Run("full range", func(t *testing.T) {
		startCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime,
			BlobKey:     nil,
		}
		endCursor := blobstore.BlobFeedCursor{
			RequestedAt: now,
			BlobKey:     nil,
		}

		// Test without limit
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*102, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[102], *lastProcessedCursor.BlobKey)

		// Test with limit
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 32)
		require.NoError(t, err)
		assert.Equal(t, 32, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*31, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[31], *lastProcessedCursor.BlobKey)
	})

	// Test cursor range boundaries
	t.Run("cursor boundaries", func(t *testing.T) {
		startCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime,
			BlobKey:     &keys[0],
		}
		endCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime + nanoSecsPerBlob,
			BlobKey:     nil,
		}

		// Test exclusive start
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)

		// Test inclusive end
		endCursor.BlobKey = &keys[1]
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(metadata))
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob, metadata[0].RequestedAt)
		checkBlobKeyEqual(t, keys[1], metadata[0].BlobHeader)
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[1], *lastProcessedCursor.BlobKey)

		// Test nil start blob key, so it should return the first blob
		startCursor.BlobKey = nil
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(metadata))
		assert.Equal(t, firstBlobTime, metadata[0].RequestedAt)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob, metadata[1].RequestedAt)
		checkBlobKeyEqual(t, keys[0], metadata[0].BlobHeader)
		checkBlobKeyEqual(t, keys[1], metadata[1].BlobHeader)
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[1], *lastProcessedCursor.BlobKey)
	})

	// Test min/max timestamp range
	t.Run("min/max timestamp range", func(t *testing.T) {
		startCursor := blobstore.BlobFeedCursor{
			RequestedAt: 0,
			BlobKey:     nil,
		}
		endCursor := blobstore.BlobFeedCursor{
			RequestedAt: math.MaxUint64,
			BlobKey:     nil,
		}

		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*102, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[102], *lastProcessedCursor.BlobKey)

		// Test future start time
		startCursor.RequestedAt = uint64(time.Now().UnixNano()) + 3600*1e9
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)
	})

	// Test pagination
	t.Run("pagination", func(t *testing.T) {
		startCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime,
			BlobKey:     nil,
		}
		endCursor := blobstore.BlobFeedCursor{
			RequestedAt: math.MaxUint64,
			BlobKey:     nil,
		}

		for i := 0; i < numBlobs; i++ {
			metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAt(ctx, startCursor, endCursor, 1)
			require.NoError(t, err)
			assert.Equal(t, 1, len(metadata))
			checkBlobKeyEqual(t, keys[i], metadata[0].BlobHeader)
			require.NotNil(t, lastProcessedCursor)
			assert.Equal(t, keys[i], *lastProcessedCursor.BlobKey)
			startCursor = *lastProcessedCursor
		}
	})
}

func TestBlobMetadataStoreGetBlobMetadataByStatusPaginated(t *testing.T) {
	ctx := context.Background()
	numBlobs := 103
	pageSize := 10
	keys := make([]corev2.BlobKey, numBlobs)
	headers := make([]*corev2.BlobHeader, numBlobs)
	metadataList := make([]*v2.BlobMetadata, numBlobs)
	dynamoKeys := make([]commondynamodb.Key, numBlobs)
	expectedCursors := make([]*blobstore.StatusIndexCursor, 0)
	for i := 0; i < numBlobs; i++ {
		blobKey, blobHeader := newBlob(t)
		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader: blobHeader,
			BlobStatus: v2.Encoded,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()),
		}

		err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		keys[i] = blobKey
		headers[i] = blobHeader
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
		metadataList[i] = metadata
		if (i+1)%pageSize == 0 {
			expectedCursors = append(expectedCursors, &blobstore.StatusIndexCursor{
				BlobKey:   &blobKey,
				UpdatedAt: metadata.UpdatedAt,
			})
		}
	}

	// Querying blobs in Queued status should return 0 results
	cursor := &blobstore.StatusIndexCursor{
		BlobKey:   nil,
		UpdatedAt: 0,
	}
	metadata, newCursor, err := blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Queued, cursor, 10)
	require.NoError(t, err)
	require.Len(t, metadata, 0)
	require.Equal(t, cursor, newCursor)

	// Querying blobs in Encoded status should return results
	cursor = &blobstore.StatusIndexCursor{
		BlobKey:   nil,
		UpdatedAt: 0,
	}
	i := 0
	numIterations := (numBlobs + pageSize - 1) / pageSize
	for i < numIterations {
		metadata, cursor, err = blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, cursor, int32(pageSize))
		require.NoError(t, err)
		if i < len(expectedCursors) {
			require.Len(t, metadata, pageSize)
			require.NotNil(t, cursor)
			require.Equal(t, cursor.BlobKey, expectedCursors[i].BlobKey)
			require.Equal(t, cursor.UpdatedAt, expectedCursors[i].UpdatedAt)
		} else {
			require.Len(t, metadata, numBlobs%pageSize)
			require.Equal(t, cursor.BlobKey, &keys[numBlobs-1])
			require.Equal(t, cursor.UpdatedAt, metadataList[numBlobs-1].UpdatedAt)
		}
		i++
	}
	lastCursor := cursor
	metadata, cursor, err = blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, cursor, int32(pageSize))
	require.NoError(t, err)
	require.Len(t, metadata, 0)
	require.Equal(t, cursor, lastCursor)

	deleteItems(t, dynamoKeys)
}

func TestBlobMetadataStoreCerts(t *testing.T) {
	ctx := context.Background()
	blobKey, blobHeader := newBlob(t)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}
	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}
	err := blobMetadataStore.PutBlobCertificate(ctx, blobCert, fragmentInfo)
	assert.NoError(t, err)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, blobCert, fetchedCert)
	assert.Equal(t, fragmentInfo, fetchedFragmentInfo)

	// blob cert with the same key should fail
	blobCert1 := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0},
	}
	err = blobMetadataStore.PutBlobCertificate(ctx, blobCert1, fragmentInfo)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// get multiple certs
	numCerts := 100
	keys := make([]corev2.BlobKey, numCerts)
	for i := 0; i < numCerts; i++ {
		blobCert := &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				QuorumNumbers:   []core.QuorumID{0},
				BlobCommitments: mockCommitment,
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         "0x123",
					ReservationPeriod: uint32(i),
					CumulativePayment: big.NewInt(321),
				},
				Signature: []byte("signature"),
			},
			RelayKeys: []corev2.RelayKey{0},
		}
		blobKey, err := blobCert.BlobHeader.BlobKey()
		assert.NoError(t, err)
		keys[i] = blobKey
		err = blobMetadataStore.PutBlobCertificate(ctx, blobCert, fragmentInfo)
		assert.NoError(t, err)
	}

	certs, fragmentInfos, err := blobMetadataStore.GetBlobCertificates(ctx, keys)
	assert.NoError(t, err)
	assert.Len(t, certs, numCerts)
	assert.Len(t, fragmentInfos, numCerts)
	reservationPeriodes := make(map[uint32]struct{})
	for i := 0; i < numCerts; i++ {
		assert.Equal(t, fragmentInfos[i], fragmentInfo)
		reservationPeriodes[certs[i].BlobHeader.PaymentMetadata.ReservationPeriod] = struct{}{}
	}
	assert.Len(t, reservationPeriodes, numCerts)
	for i := 0; i < numCerts; i++ {
		assert.Contains(t, reservationPeriodes, uint32(i))
	}

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobCertificate"},
		},
	})
}

func TestBlobMetadataStoreUpdateBlobStatus(t *testing.T) {
	ctx := context.Background()
	blobKey, blobHeader := newBlob(t)

	now := time.Now()
	metadata := &v2.BlobMetadata{
		BlobHeader: blobHeader,
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
	assert.NoError(t, err)

	// Update the blob status to invalid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Certified)
	assert.ErrorIs(t, err, blobstore.ErrInvalidStateTransition)

	// Update the blob status to a valid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Encoded)
	assert.NoError(t, err)

	// Update the blob status to same status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Encoded)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, fetchedMetadata.BlobStatus, v2.Encoded)
	assert.Greater(t, fetchedMetadata.UpdatedAt, metadata.UpdatedAt)

	// Update the blob status to a valid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
	assert.NoError(t, err)

	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, fetchedMetadata.BlobStatus, v2.Failed)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func TestBlobMetadataStoreDispersals(t *testing.T) {
	ctx := context.Background()
	opID := core.OperatorID{0, 1}
	dispersalRequest := &corev2.DispersalRequest{
		OperatorID:      opID,
		OperatorAddress: gethcommon.HexToAddress("0x1234567"),
		Socket:          "socket",
		DispersedAt:     uint64(time.Now().UnixNano()),

		BatchHeader: corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: 100,
		},
	}

	err := blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest)
	assert.NoError(t, err)

	bhh, err := dispersalRequest.BatchHeader.Hash()
	assert.NoError(t, err)

	fetchedRequest, err := blobMetadataStore.GetDispersalRequest(ctx, bhh, dispersalRequest.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersalRequest, fetchedRequest)

	// attempt to put dispersal request with the same key should fail
	err = blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	dispersalResponse := &corev2.DispersalResponse{
		DispersalRequest: dispersalRequest,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        [32]byte{1, 1, 1},
		Error:            "error",
	}

	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
	assert.NoError(t, err)

	fetchedResponse, err := blobMetadataStore.GetDispersalResponse(ctx, bhh, dispersalRequest.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersalResponse, fetchedResponse)

	// attempt to put dispersal response with the same key should fail
	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// the other operator's response for the same batch
	opID2 := core.OperatorID{2, 3}
	dispersalRequest2 := &corev2.DispersalRequest{
		OperatorID:      opID2,
		OperatorAddress: gethcommon.HexToAddress("0x2234567"),
		Socket:          "socket",
		DispersedAt:     uint64(time.Now().UnixNano()),
		BatchHeader: corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: 100,
		},
	}
	err = blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest2)
	assert.NoError(t, err)
	dispersalResponse2 := &corev2.DispersalResponse{
		DispersalRequest: dispersalRequest2,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        [32]byte{1, 1, 1},
		Error:            "",
	}
	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse2)
	assert.NoError(t, err)

	responses, err := blobMetadataStore.GetDispersalResponses(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responses))
	assert.Equal(t, dispersalResponse, responses[0])
	assert.Equal(t, dispersalResponse2, responses[1])

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalRequest#" + opID.Hex()},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalRequest#" + opID2.Hex()},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalResponse#" + opID.Hex()},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalResponse#" + opID2.Hex()},
		},
	})
}

func TestBlobMetadataStoreVerificationInfo(t *testing.T) {
	ctx := context.Background()
	blobKey := corev2.BlobKey{1, 1, 1}
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)
	verificationInfo := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      10,
		InclusionProof: []byte("proof"),
	}

	err = blobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo)
	assert.NoError(t, err)

	fetchedInfo, err := blobMetadataStore.GetBlobVerificationInfo(ctx, blobKey, bhh)
	assert.NoError(t, err)
	assert.Equal(t, verificationInfo, fetchedInfo)

	// attempt to put verification info with the same key should fail
	err = blobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// put multiple verification infos
	blobKey1 := corev2.BlobKey{2, 2, 2}
	verificationInfo1 := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey1,
		BlobIndex:      12,
		InclusionProof: []byte("proof 1"),
	}
	blobKey2 := corev2.BlobKey{3, 3, 3}
	verificationInfo2 := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey2,
		BlobIndex:      14,
		InclusionProof: []byte("proof 2"),
	}
	err = blobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.NoError(t, err)

	// test retries
	nonTransientError := errors.New("non transient error")
	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, nonTransientError).Once()
	err = mockedBlobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.ErrorIs(t, err, nonTransientError)

	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).Return([]dynamodb.Item{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey1.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
		},
	}, nil).Run(func(args mock.Arguments) {
		items := args.Get(2).([]dynamodb.Item)
		assert.Len(t, items, 2)
	}).Once()
	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Run(func(args mock.Arguments) {
			items := args.Get(2).([]dynamodb.Item)
			assert.Len(t, items, 1)
		}).
		Once()
	err = mockedBlobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.NoError(t, err)
	mockDynamoClient.AssertNumberOfCalls(t, "PutItems", 3)
}

func TestBlobMetadataStoreBatchAttestation(t *testing.T) {
	ctx := context.Background()
	h := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	bhh, err := h.Hash()
	assert.NoError(t, err)

	err = blobMetadataStore.PutBatchHeader(ctx, h)
	assert.NoError(t, err)

	fetchedHeader, err := blobMetadataStore.GetBatchHeader(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)

	// attempt to put batch header with the same key should fail
	err = blobMetadataStore.PutBatchHeader(ctx, h)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)

	apk := keyPair.GetPubKeyG2()
	attestation := &corev2.Attestation{
		BatchHeader: h,
		AttestedAt:  uint64(time.Now().UnixNano()),
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(1), big.NewInt(2)),
			core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		},
		APKG2: apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
			1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}

	err = blobMetadataStore.PutAttestation(ctx, attestation)
	assert.NoError(t, err)

	fetchedAttestation, err := blobMetadataStore.GetAttestation(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, attestation, fetchedAttestation)

	// attempt to put attestation with the same key should fail
	err = blobMetadataStore.PutAttestation(ctx, attestation)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// attempt to retrieve batch header and attestation at the same time
	fetchedHeader, fetchedAttestation, err = blobMetadataStore.GetSignedBatch(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)
	assert.Equal(t, attestation, fetchedAttestation)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
	})
}

func deleteItems(t *testing.T, keys []commondynamodb.Key) {
	failed, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
	assert.Len(t, failed, 0)
}

func newBlob(t *testing.T) (corev2.BlobKey, *corev2.BlobHeader) {
	accountBytes := make([]byte, 32)
	_, err := rand.Read(accountBytes)
	require.NoError(t, err)
	accountID := hex.EncodeToString(accountBytes)
	reservationPeriod, err := rand.Int(rand.Reader, big.NewInt(256))
	require.NoError(t, err)
	cumulativePayment, err := rand.Int(rand.Reader, big.NewInt(1024))
	require.NoError(t, err)
	sig := make([]byte, 32)
	_, err = rand.Read(sig)
	require.NoError(t, err)
	bh := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			ReservationPeriod: uint32(reservationPeriod.Int64()),
			CumulativePayment: cumulativePayment,
		},
		Signature: sig,
	}
	bk, err := bh.BlobKey()
	require.NoError(t, err)
	return bk, bh
}
