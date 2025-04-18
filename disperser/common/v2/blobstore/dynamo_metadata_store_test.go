package blobstore_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
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

func checkAttestationsAsc(t *testing.T, items []*corev2.Attestation) {
	if len(items) > 1 {
		for i := 1; i < len(items); i++ {
			assert.Less(t,
				items[i-1].AttestedAt, // previous should be less
				items[i].AttestedAt,   // than current
				"attestations should be in ascending order",
			)
		}
	}
}

func checkAttestationsDesc(t *testing.T, items []*corev2.Attestation) {
	for i := 1; i < len(items); i++ {
		assert.Greater(t,
			items[i-1].AttestedAt, // previous should be greater
			items[i].AttestedAt,   // than current
			"attestations should be in descending order",
		)
	}
}

func checkDispersalsAsc(t *testing.T, items []*corev2.DispersalRequest) {
	if len(items) > 1 {
		for i := 1; i < len(items); i++ {
			assert.Less(
				t,
				items[i-1].DispersedAt, // previous should be less
				items[i].DispersedAt,   // than current
				"DispersalRequests should be in ascending order",
			)

		}
	}
}

func checkDispersalsDesc(t *testing.T, items []*corev2.DispersalRequest) {
	for i := 1; i < len(items); i++ {
		assert.Greater(
			t,
			items[i-1].DispersedAt, // previous should be greater
			items[i].DispersedAt,   // than current
			"DispersalRequests should be in descending order",
		)
	}
}

func checkBlobsAsc(t *testing.T, items []*v2.BlobMetadata) {
	if len(items) > 1 {
		for i := 1; i < len(items); i++ {
			assert.Less(t,
				items[i-1].RequestedAt, // previous should be less
				items[i].RequestedAt,   // than current
				"blobs should be in ascending order",
			)
		}
	}
}

func checkBlobsDesc(t *testing.T, items []*v2.BlobMetadata) {
	for i := 1; i < len(items); i++ {
		assert.Greater(t,
			items[i-1].RequestedAt, // previous should be greater
			items[i].RequestedAt,   // than current
			"blobs should be in descending order",
		)
	}
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
		Signature:  []byte{1, 2, 3},
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	metadata2 := &v2.BlobMetadata{
		BlobHeader: blobHeader2,
		Signature:  []byte{4, 5, 6},
		BlobStatus: v2.Complete,
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

	complete, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Complete, 0)
	assert.NoError(t, err)
	assert.Len(t, complete, 1)
	assert.Equal(t, metadata2, complete[0])

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

func TestBlobMetadataStoreGetBlobMetadataByRequestedAtForwardWithIdenticalTimestamp(t *testing.T) {
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
			Signature:   []byte{1, 2, 3},
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

	metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
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
		RequestedAt: requestedAts[4],
		BlobKey:     nil,
	}

	// Test with different end cursors
	testCases := []struct {
		endBlobKey *corev2.BlobKey
		expectLen  int
		expectLast int
	}{
		{nil, 1, 2},
		{&keys[3], 1, 2}, // keys[2] will be retrieved
		{&keys[4], 2, 3}, // keys[2], keys[3] will be retrieved
	}

	for _, tc := range testCases {
		endCursor.BlobKey = tc.endBlobKey
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
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

func TestBlobMetadataStoreGetBlobMetadataByRequestedAtForwardWithDynamoPagination(t *testing.T) {
	ctx := context.Background()

	// Make all blobs happen in 120s
	numBlobs := 1200
	nanoSecsPerBlob := uint64(1e8) // 10 blob per second

	now := uint64(time.Now().UnixNano())
	firstBlobTime := now - uint64(10*time.Minute.Nanoseconds())
	// Adjust "now" so all blobs will deterministically fall in just one
	// bucket.
	startBucket, endBucket := blobstore.GetRequestedAtBucketIDRange(firstBlobTime-1, now)
	if startBucket < endBucket {
		now -= uint64(11 * time.Minute.Nanoseconds())
		firstBlobTime = now - uint64(10*time.Minute.Nanoseconds())
	}
	startBucket, endBucket = blobstore.GetAttestedAtBucketIDRange(firstBlobTime-1, now)
	require.Equal(t, startBucket, endBucket)

	// Create blobs for testing
	// The num of blobs here are large enough to make it more than 1MB (the max response
	// size of DyanamoDB) so it will have to use DynamoDB's pagination to get all desired
	// results.
	keys := make([]corev2.BlobKey, numBlobs)
	dynamoKeys := make([]commondynamodb.Key, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobKey, blobHeader := newBlob(t)
		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader:  blobHeader,
			Signature:   []byte{1, 2, 3},
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

	startCursor := blobstore.BlobFeedCursor{
		RequestedAt: firstBlobTime,
		BlobKey:     nil,
	}
	endCursor := blobstore.BlobFeedCursor{
		RequestedAt: now + 1,
		BlobKey:     nil,
	}
	blobs, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
	require.NoError(t, err)
	require.Equal(t, numBlobs, len(blobs))
	require.NotNil(t, lastProcessedCursor)
	assert.Equal(t, firstBlobTime+nanoSecsPerBlob*uint64(numBlobs-1), lastProcessedCursor.RequestedAt)
	assert.Equal(t, keys[numBlobs-1], *lastProcessedCursor.BlobKey)
}

func TestBlobMetadataStoreGetBlobMetadataByRequestedAtForward(t *testing.T) {
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
			Signature:   []byte{1, 2, 3},
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
		_, _, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, startCursor, 10)
		assert.Error(t, err)
		assert.Equal(t, "after cursor must be less than before cursor", err.Error())

		// Test empty range
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 10)
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
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*102, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[102], *lastProcessedCursor.BlobKey)

		// Test with limit
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 32)
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
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)

		// Test exclusive end
		endCursor.BlobKey = &keys[1]
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		require.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)

		endCursor.RequestedAt = firstBlobTime + nanoSecsPerBlob + 1 // pass the time of second blob
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		require.Equal(t, 1, len(metadata))
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob, metadata[0].RequestedAt)
		checkBlobKeyEqual(t, keys[1], metadata[0].BlobHeader)
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[1], *lastProcessedCursor.BlobKey)

		// Test nil start blob key, so it should return the first blob
		startCursor.BlobKey = nil
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
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

		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*102, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[102], *lastProcessedCursor.BlobKey)

		// Test future start time
		startCursor.RequestedAt = uint64(time.Now().UnixNano()) + 3600*1e9
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 0)
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
			metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtForward(ctx, startCursor, endCursor, 1)
			require.NoError(t, err)
			require.Equal(t, 1, len(metadata))
			checkBlobKeyEqual(t, keys[i], metadata[0].BlobHeader)
			require.NotNil(t, lastProcessedCursor)
			assert.Equal(t, keys[i], *lastProcessedCursor.BlobKey)
			startCursor = *lastProcessedCursor
		}
	})
}

func TestBlobMetadataStoreGetBlobMetadataByRequestedAtBackward(t *testing.T) {
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
			Signature:   []byte{1, 2, 3},
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
		beforeCursor := blobstore.BlobFeedCursor{
			RequestedAt: now + 10*1e9,
			BlobKey:     nil,
		}
		afterCursor := blobstore.BlobFeedCursor{
			RequestedAt: now,
			BlobKey:     nil,
		}

		// Test equal cursors error
		_, _, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, beforeCursor, 10)
		assert.Error(t, err)
		assert.Equal(t, "after cursor must be less than before cursor", err.Error())

		// Test empty range
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)
	})

	// Test full range query
	t.Run("full range", func(t *testing.T) {
		beforeCursor := blobstore.BlobFeedCursor{
			RequestedAt: now,
			BlobKey:     nil,
		}
		afterCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime,
			BlobKey:     nil,
		}

		// Test without limit
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[0], *lastProcessedCursor.BlobKey)

		// Test with limit
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 32)
		require.NoError(t, err)
		assert.Equal(t, 32, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob*71, lastProcessedCursor.RequestedAt) // numBlobs-32
		assert.Equal(t, keys[71], *lastProcessedCursor.BlobKey)
	})

	t.Run("cursor boundaries", func(t *testing.T) {
		beforeCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime + nanoSecsPerBlob, // time of blob[1]
			BlobKey:     &keys[1],                        // exclusive
		}
		afterCursor := blobstore.BlobFeedCursor{
			RequestedAt: firstBlobTime, // time of blob[0]
			BlobKey:     &keys[0],      // exclusive
		}

		// Test exclusive before, exclusive after
		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(
			ctx,
			beforeCursor, // blob[1] excluded
			afterCursor,  // blob[0] excluded
			0,
		)
		require.NoError(t, err)
		require.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)

		// Test the effects of blob key in before cursor
		beforeCursor.RequestedAt = firstBlobTime + nanoSecsPerBlob*2 // time of blob[2]
		beforeCursor.BlobKey = &keys[2]                              // exclusive of blob[2]
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtBackward(
			ctx,
			beforeCursor, // excludes blob[2]
			afterCursor,  // excludes blob[0]
			0,
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(metadata))
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob, metadata[0].RequestedAt) // blob[1]
		checkBlobKeyEqual(t, keys[1], metadata[0].BlobHeader)
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[1], *lastProcessedCursor.BlobKey)

		// Test when removing blob key from after cursor
		afterCursor.BlobKey = nil // makes after cursor point to before blob[0]
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtBackward(
			ctx,
			beforeCursor, // excludes blob[2]
			afterCursor,  // now points to before blob[0], so blob[0] will be included
			0,
		)
		require.NoError(t, err)
		require.Equal(t, 2, len(metadata))
		assert.Equal(t, firstBlobTime+nanoSecsPerBlob, metadata[0].RequestedAt) // blob[1]
		assert.Equal(t, firstBlobTime, metadata[1].RequestedAt)                 // blob[0]
		checkBlobKeyEqual(t, keys[1], metadata[0].BlobHeader)
		checkBlobKeyEqual(t, keys[0], metadata[1].BlobHeader)
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, keys[0], *lastProcessedCursor.BlobKey)
	})

	// Test min/max timestamp range
	t.Run("min/max timestamp range", func(t *testing.T) {
		beforeCursor := blobstore.BlobFeedCursor{
			RequestedAt: math.MaxUint64,
			BlobKey:     nil,
		}
		afterCursor := blobstore.BlobFeedCursor{
			RequestedAt: 0,
			BlobKey:     nil,
		}

		metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, numBlobs, len(metadata))
		require.NotNil(t, lastProcessedCursor)
		assert.Equal(t, firstBlobTime, lastProcessedCursor.RequestedAt)
		assert.Equal(t, keys[0], *lastProcessedCursor.BlobKey)

		// Test past `after` time
		afterCursor.RequestedAt = uint64(time.Now().UnixNano()) + 3600*1e9
		metadata, lastProcessedCursor, err = blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, len(metadata))
		assert.Nil(t, lastProcessedCursor)
	})

	// Test pagination
	t.Run("pagination", func(t *testing.T) {
		beforeCursor := blobstore.BlobFeedCursor{
			RequestedAt: math.MaxUint64,
			BlobKey:     nil,
		}
		afterCursor := blobstore.BlobFeedCursor{
			RequestedAt: 0,
			BlobKey:     nil,
		}

		for i := numBlobs - 1; i >= 0; i-- {
			metadata, lastProcessedCursor, err := blobMetadataStore.GetBlobMetadataByRequestedAtBackward(ctx, beforeCursor, afterCursor, 1)
			require.NoError(t, err)
			assert.Equal(t, 1, len(metadata))
			checkBlobKeyEqual(t, keys[i], metadata[0].BlobHeader)
			require.NotNil(t, lastProcessedCursor)
			assert.Equal(t, keys[i], *lastProcessedCursor.BlobKey)
			beforeCursor = *lastProcessedCursor
		}
	})
}

func TestBlobMetadataStoreGetBlobMetadataByAccountID(t *testing.T) {
	ctx := context.Background()

	// Make all blobs happen in 12s
	numBlobs := 120
	nanoSecsPerBlob := uint64(1e8) // 10 blobs per second

	now := uint64(time.Now().UnixNano())
	firstBlobTime := now - uint64(10*time.Minute.Nanoseconds())

	accountId := gethcommon.HexToAddress(fmt.Sprintf("0x000000000000000000000000000000000000000%d", 5))

	// Create blobs for testing
	keys := make([]corev2.BlobKey, numBlobs)
	requestedAt := make([]uint64, numBlobs)
	dynamoKeys := make([]commondynamodb.Key, numBlobs)
	for i := 0; i < numBlobs; i++ {
		_, blobHeader := newBlob(t)
		blobHeader.PaymentMetadata.AccountID = accountId
		blobKey, err := blobHeader.BlobKey()
		require.NoError(t, err)
		requestedAt[i] = firstBlobTime + nanoSecsPerBlob*uint64(i)
		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader:  blobHeader,
			Signature:   []byte{1, 2, 3},
			BlobStatus:  v2.Encoded,
			Expiry:      uint64(now.Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   uint64(now.UnixNano()),
			RequestedAt: requestedAt[i],
		}
		err = blobMetadataStore.PutBlobMetadata(ctx, metadata)
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
		// Test invalid time range
		_, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, 1, 1, 0, true)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 1)", err.Error())

		_, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, 1, 2, 0, true)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 2)", err.Error())

		// Test empty range
		blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, now, now+1024, 0, true)
		require.NoError(t, err)
		assert.Equal(t, 0, len(blobs))
	})

	// Test full range query
	t.Run("ascending full range", func(t *testing.T) {
		// Test without limit
		blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsAsc(t, blobs)

		// Test with limit
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, now, 10, true)
		require.NoError(t, err)
		require.Equal(t, 10, len(blobs))
		checkBlobsAsc(t, blobs)

		// Test min/max timestamp range
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, 0, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsAsc(t, blobs)
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, math.MaxInt64, 0, true)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsAsc(t, blobs)
	})

	// Test full range query
	t.Run("descending full range", func(t *testing.T) {
		// Test without limit
		blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsDesc(t, blobs)

		// Test with limit
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, now, 10, false)
		require.NoError(t, err)
		require.Equal(t, 10, len(blobs))
		checkBlobsDesc(t, blobs)

		// Test min/max timestamp range
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, 0, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsDesc(t, blobs)
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, math.MaxInt64, 0, false)
		require.NoError(t, err)
		require.Equal(t, numBlobs, len(blobs))
		checkBlobsDesc(t, blobs)
	})

	// Test range boundaries
	t.Run("ascending range boundaries", func(t *testing.T) {
		// Test exclusive start
		blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numBlobs-1, len(blobs))
		assert.Equal(t, requestedAt[1], blobs[0].RequestedAt)
		assert.Equal(t, requestedAt[numBlobs-1], blobs[numBlobs-2].RequestedAt)
		checkBlobsAsc(t, blobs)

		// Test exclusive end
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, requestedAt[4], 0, true)
		require.NoError(t, err)
		require.Equal(t, 4, len(blobs))
		assert.Equal(t, requestedAt[0], blobs[0].RequestedAt)
		assert.Equal(t, requestedAt[3], blobs[3].RequestedAt)
		checkBlobsAsc(t, blobs)
	})

	// Test range boundaries
	t.Run("descending range boundaries", func(t *testing.T) {
		// Test exclusive start
		blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numBlobs-1, len(blobs))
		assert.Equal(t, requestedAt[numBlobs-1], blobs[0].RequestedAt)
		assert.Equal(t, requestedAt[1], blobs[numBlobs-2].RequestedAt)
		checkBlobsDesc(t, blobs)

		// Test exclusive end
		blobs, err = blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, firstBlobTime-1, requestedAt[4], 0, false)
		require.NoError(t, err)
		require.Equal(t, 4, len(blobs))
		assert.Equal(t, requestedAt[3], blobs[0].RequestedAt)
		assert.Equal(t, requestedAt[0], blobs[3].RequestedAt)
		checkBlobsDesc(t, blobs)
	})

	// Test pagination
	t.Run("pagination", func(t *testing.T) {
		for i := 1; i < numBlobs; i++ {
			blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, requestedAt[i-1], requestedAt[i]+1, 0, true)
			require.NoError(t, err)
			require.Equal(t, 1, len(blobs))
			assert.Equal(t, requestedAt[i], blobs[0].RequestedAt)
		}

		for i := 1; i < numBlobs; i++ {
			blobs, err := blobMetadataStore.GetBlobMetadataByAccountID(ctx, accountId, requestedAt[i-1], requestedAt[i]+1, 0, false)
			require.NoError(t, err)
			require.Equal(t, 1, len(blobs))
			assert.Equal(t, requestedAt[i], blobs[0].RequestedAt)
		}
	})
}

func TestBlobMetadataStoreGetAttestationByAttestedAtForward(t *testing.T) {
	ctx := context.Background()
	numBatches := 72
	now := uint64(time.Now().UnixNano())
	firstBatchTs := now - uint64((72+2)*time.Hour.Nanoseconds())
	nanoSecsPerBatch := uint64(time.Hour.Nanoseconds()) // 1 batch per hour

	// Create attestations for testing
	attestedAt := make([]uint64, numBatches)
	batchHeaders := make([]*corev2.BatchHeader, numBatches)
	dynamoKeys := make([]commondynamodb.Key, numBatches)
	for i := 0; i < numBatches; i++ {
		batchHeaders[i] = &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, byte(i)},
			ReferenceBlockNumber: uint64(i + 1),
		}
		bhh, err := batchHeaders[i].Hash()
		assert.NoError(t, err)
		keyPair, err := core.GenRandomBlsKeys()
		assert.NoError(t, err)
		apk := keyPair.GetPubKeyG2()
		attestedAt[i] = firstBatchTs + uint64(i)*nanoSecsPerBatch
		attestation := &corev2.Attestation{
			BatchHeader: batchHeaders[i],
			AttestedAt:  attestedAt[i],
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
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	// Test empty range
	t.Run("empty range", func(t *testing.T) {
		// Test invalid time range
		_, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, 1, 1, 0)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 1)", err.Error())

		_, err = blobMetadataStore.GetAttestationByAttestedAtForward(ctx, 1, 2, 0)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 2)", err.Error())

		// Test empty range
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, now, now+uint64(240*time.Hour.Nanoseconds()), 0)
		require.NoError(t, err)
		assert.Equal(t, 0, len(attestations))
	})

	// Test full range query
	t.Run("full range", func(t *testing.T) {
		// Test without limit
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs-1, now, 0)
		require.NoError(t, err)
		require.Equal(t, numBatches, len(attestations))
		checkAttestationsAsc(t, attestations)

		// Test with limit
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs, now, 10)
		require.NoError(t, err)
		require.Equal(t, 10, len(attestations))
		checkAttestationsAsc(t, attestations)

		// Test min/max timestamp range
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtForward(ctx, 0, now, 0)
		require.NoError(t, err)
		require.Equal(t, numBatches, len(attestations))
		checkAttestationsAsc(t, attestations)
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs-1, math.MaxInt64, 0)
		require.NoError(t, err)
		require.Equal(t, numBatches, len(attestations))
		checkAttestationsAsc(t, attestations)
	})

	// Test range boundaries
	t.Run("range boundaries", func(t *testing.T) {
		// Test exclusive start
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs, now+1, 0)
		require.NoError(t, err)
		require.Equal(t, numBatches-1, len(attestations))
		checkAttestationsAsc(t, attestations)
		assert.Equal(t, attestedAt[1], attestations[0].AttestedAt)
		assert.Equal(t, batchHeaders[1].BatchRoot, attestations[0].BatchRoot)
		assert.Equal(t, attestedAt[numBatches-1], attestations[numBatches-2].AttestedAt)
		assert.Equal(t, batchHeaders[numBatches-1].BatchRoot, attestations[numBatches-2].BatchRoot)

		// Test exclusive end
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs-1, attestedAt[4], 0)
		require.NoError(t, err)
		require.Equal(t, 4, len(attestations))
		checkAttestationsAsc(t, attestations)
		assert.Equal(t, attestedAt[0], attestations[0].AttestedAt)
		assert.Equal(t, batchHeaders[0].BatchRoot, attestations[0].BatchRoot)
		assert.Equal(t, attestedAt[3], attestations[3].AttestedAt)
		assert.Equal(t, batchHeaders[3].BatchRoot, attestations[3].BatchRoot)
	})

	// Test pagination
	t.Run("pagination", func(t *testing.T) {
		for i := 1; i < numBatches; i++ {
			attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, attestedAt[i-1], attestedAt[i]+1, 1)
			require.NoError(t, err)
			require.Equal(t, 1, len(attestations))
			assert.Equal(t, attestedAt[i], attestations[0].AttestedAt)
			assert.Equal(t, batchHeaders[i].BatchRoot, attestations[0].BatchRoot)
		}
	})
}

func TestBlobMetadataStoreGetAttestationByAttestedAtBackward(t *testing.T) {
	ctx := context.Background()
	numBatches := 72
	now := uint64(time.Now().UnixNano())
	firstBatchTs := now - uint64((72+2)*time.Hour.Nanoseconds())
	nanoSecsPerBatch := uint64(time.Hour.Nanoseconds()) // 1 batch per hour

	// Create attestations for testing
	attestedAt := make([]uint64, numBatches)
	batchHeaders := make([]*corev2.BatchHeader, numBatches)
	dynamoKeys := make([]commondynamodb.Key, numBatches)
	for i := 0; i < numBatches; i++ {
		batchHeaders[i] = &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, byte(i)},
			ReferenceBlockNumber: uint64(i + 1),
		}
		bhh, err := batchHeaders[i].Hash()
		assert.NoError(t, err)
		keyPair, err := core.GenRandomBlsKeys()
		assert.NoError(t, err)
		apk := keyPair.GetPubKeyG2()
		attestedAt[i] = firstBatchTs + uint64(i)*nanoSecsPerBatch
		attestation := &corev2.Attestation{
			BatchHeader: batchHeaders[i],
			AttestedAt:  attestedAt[i],
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
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	t.Run("empty range", func(t *testing.T) {
		// Test invalid time range
		_, err := blobMetadataStore.GetAttestationByAttestedAtBackward(ctx, 1, 1, 0)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 1)", err.Error())

		_, err = blobMetadataStore.GetAttestationByAttestedAtBackward(ctx, 2, 1, 0)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 2)", err.Error())

		// Test empty range
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx,
			now-uint64(240*time.Hour.Nanoseconds()), // before
			now-uint64(241*time.Hour.Nanoseconds()), // after
			0,
		)
		require.NoError(t, err)
		assert.Equal(t, 0, len(attestations))
	})

	t.Run("full range", func(t *testing.T) {
		// Test without limit - traverse from now back to firstBatchTs
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx,
			now+1,          // before (exclusive)
			firstBatchTs-1, // after (inclusive)
			0,
		)
		require.NoError(t, err)
		require.Equal(t, numBatches, len(attestations))
		checkAttestationsDesc(t, attestations)

		// Test with limit
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx,
			now+1,          // before
			firstBatchTs-1, // after
			10,
		)
		require.NoError(t, err)
		require.Equal(t, 10, len(attestations))
		checkAttestationsDesc(t, attestations)
	})

	t.Run("range boundaries", func(t *testing.T) {
		// Test exclusive before - should skip the newest item
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx,
			attestedAt[numBatches-1], // before (exclusive)
			firstBatchTs,             // after (exclusive)
			0,
		)
		require.NoError(t, err)
		require.Equal(t, numBatches-2, len(attestations))
		// The first one returned is not "before" (as "before" is exclusive)
		assert.Equal(t, attestedAt[numBatches-2], attestations[0].AttestedAt)
		// The last one returned is the second batch (as "after" is exclusive)
		assert.Equal(t, attestedAt[1], attestations[numBatches-3].AttestedAt)
		checkAttestationsDesc(t, attestations)

		// Test exclusive after - should not include the oldest item
		attestations, err = blobMetadataStore.GetAttestationByAttestedAtBackward(
			ctx,
			attestedAt[4]+1, // before: just after 4th item (so this batch should be included)
			attestedAt[0],   // after: oldest item (should not be included)
			0,
		)
		require.NoError(t, err)
		require.Equal(t, 4, len(attestations))
		assert.Equal(t, attestedAt[4], attestations[0].AttestedAt)
		assert.Equal(t, attestedAt[1], attestations[3].AttestedAt)
		checkAttestationsDesc(t, attestations)
	})

	t.Run("pagination", func(t *testing.T) {
		for i := numBatches - 1; i > 0; i-- {
			attestations, err := blobMetadataStore.GetAttestationByAttestedAtBackward(
				ctx,
				attestedAt[i]+1, // before: just after current item
				attestedAt[i-1], // after: previous item (included)
				1,
			)
			require.NoError(t, err)
			require.Equal(t, 1, len(attestations))
			assert.Equal(t, attestedAt[i], attestations[0].AttestedAt)
		}
	})
}

func TestBlobMetadataStoreGetAttestationByAttestedAtForwardWithDynamoPagination(t *testing.T) {
	ctx := context.Background()

	now := uint64(time.Now().UnixNano())
	firstBatchTs := now - uint64(5*time.Minute.Nanoseconds())
	// Adjust "now" so all attestations will deterministically fall in just one
	// bucket.
	startBucket, endBucket := blobstore.GetAttestedAtBucketIDRange(firstBatchTs-1, now)
	if startBucket < endBucket {
		now -= uint64(time.Hour.Nanoseconds())
		firstBatchTs = now - uint64(5*time.Minute.Nanoseconds())
	}
	startBucket, endBucket = blobstore.GetAttestedAtBucketIDRange(firstBatchTs-1, now)
	require.Equal(t, startBucket, endBucket)

	numBatches := 240
	nanoSecsPerBatch := uint64(time.Second.Nanoseconds()) // 1 batch per second

	// Create attestations for testing
	attestedAt := make([]uint64, numBatches)
	batchHeaders := make([]*corev2.BatchHeader, numBatches)
	dynamoKeys := make([]commondynamodb.Key, numBatches)
	for i := 0; i < numBatches; i++ {
		batchHeaders[i] = &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, byte(i)},
			ReferenceBlockNumber: uint64(i + 1),
		}
		bhh, err := batchHeaders[i].Hash()
		assert.NoError(t, err)
		keyPair, err := core.GenRandomBlsKeys()
		assert.NoError(t, err)
		apk := keyPair.GetPubKeyG2()
		attestedAt[i] = firstBatchTs + uint64(i)*nanoSecsPerBatch
		// Create a sizable nonsigners so the attestation message is big
		nonsigners := make([]*core.G1Point, 0)
		for i := 0; i < 200; i++ {
			nonsigners = append(nonsigners, core.NewG1Point(big.NewInt(int64(i)), big.NewInt(int64(i+1))))
		}
		attestation := &corev2.Attestation{
			BatchHeader:      batchHeaders[i],
			AttestedAt:       attestedAt[i],
			NonSignerPubKeys: nonsigners,
			APKG2:            apk,
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
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		}
	}
	// The total bytes written to the bucket will be greater than 1MB, so if a query tries to
	// fetch all results in the bucket, it has to use pagination.
	// Each attestation has 200 nonsigners and the G1 point has 32 bytes, so we have
	// 32*3200*numBatches bytes just for nonsigners (attestations' size must be greater).
	assert.True(t, 32*200*numBatches > 1*1024*1024)

	defer deleteItems(t, dynamoKeys)

	// Test the query can fetch all attestations in a bucket
	t.Run("full range", func(t *testing.T) {
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs-1, now, 0)
		require.NoError(t, err)
		require.Equal(t, numBatches, len(attestations))
		checkAttestationsAsc(t, attestations)
	})

	// Test the query returns after getting desired num of attestations in a bucket
	t.Run("return after getting desired num of items", func(t *testing.T) {
		attestations, err := blobMetadataStore.GetAttestationByAttestedAtForward(ctx, firstBatchTs-1, now, 125)
		require.NoError(t, err)
		require.Equal(t, 125, len(attestations))
		checkAttestationsAsc(t, attestations)
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
			require.Nil(t, cursor)
		}
		i++
	}

	for i := 0; i < numBlobs; i++ {
		err = blobMetadataStore.UpdateBlobStatus(ctx, keys[i], v2.GatheringSignatures)
		require.NoError(t, err)
	}

	metadata, cursor, err = blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, cursor, int32(pageSize))
	require.NoError(t, err)
	require.Len(t, metadata, 0)
	require.Nil(t, cursor)

	deleteItems(t, dynamoKeys)
}

func TestBlobMetadataStoreCerts(t *testing.T) {
	ctx := context.Background()
	blobKey, blobHeader := newBlob(t)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		Signature:  []byte("signature"),
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
					AccountID:         gethcommon.HexToAddress("0x123"),
					Timestamp:         int64(i),
					CumulativePayment: big.NewInt(321),
				},
			},
			Signature: []byte("signature"),
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
	timestamps := make(map[int64]struct{})
	for i := 0; i < numCerts; i++ {
		assert.Equal(t, fragmentInfos[i], fragmentInfo)
		timestamps[certs[i].BlobHeader.PaymentMetadata.Timestamp] = struct{}{}
	}
	assert.Len(t, timestamps, numCerts)
	for i := 0; i < numCerts; i++ {
		assert.Contains(t, timestamps, int64(i))
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
		Signature:  []byte("signature"),
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
	assert.NoError(t, err)

	// Update the blob status to invalid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Complete)
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

func TestBlobMetadataStoreDispersalsByDispersedAt(t *testing.T) {
	ctx := context.Background()

	numRequests := 60
	opID := core.OperatorID{16, 32}
	now := uint64(time.Now().UnixNano())
	firstRequestTs := now - uint64(int64(numRequests)*time.Second.Nanoseconds())
	nanoSecsPerRequest := uint64(time.Second.Nanoseconds()) // 1 batch/s

	dispersedAt := make([]uint64, numRequests)
	dynamoKeys := make([]commondynamodb.Key, numRequests)
	for i := 0; i < numRequests; i++ {
		dispersedAt[i] = firstRequestTs + uint64(i)*nanoSecsPerRequest
		dispersalRequest := &corev2.DispersalRequest{
			OperatorID:      opID,
			OperatorAddress: gethcommon.HexToAddress("0x1234567"),
			Socket:          "socket",
			DispersedAt:     dispersedAt[i],
			BatchHeader: corev2.BatchHeader{
				BatchRoot:            [32]byte{1, 2, 3},
				ReferenceBlockNumber: uint64(i + 100),
			},
		}

		err := blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest)
		require.NoError(t, err)

		bhh, err := dispersalRequest.BatchHeader.Hash()
		require.NoError(t, err)
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalRequest#" + opID.Hex()},
		}
	}
	defer deleteItems(t, dynamoKeys)

	// Test empty range
	t.Run("empty range", func(t *testing.T) {
		// Test invalid time range
		_, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, 1, 1, 0, true)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 1)", err.Error())

		_, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, 1, 2, 0, true)
		require.Error(t, err)
		assert.Equal(t, "no time point in exclusive time range (1, 2)", err.Error())

		// Test empty range
		dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, now, now+1024, 0, true)
		require.NoError(t, err)
		assert.Equal(t, 0, len(dispersals))
	})

	// Test full range query
	t.Run("ascending full range", func(t *testing.T) {
		// Test without limit
		dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsAsc(t, dispersals)

		// Test with limit
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, now, 10, true)
		require.NoError(t, err)
		require.Equal(t, 10, len(dispersals))
		checkDispersalsAsc(t, dispersals)

		// Test min/max timestamp range
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, 0, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsAsc(t, dispersals)
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, math.MaxInt64, 0, true)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsAsc(t, dispersals)
	})

	// Test full range query
	t.Run("descending full range", func(t *testing.T) {
		// Test without limit
		dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsDesc(t, dispersals)

		// Test with limit
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs, now, 10, false)
		require.NoError(t, err)
		require.Equal(t, 10, len(dispersals))
		checkDispersalsDesc(t, dispersals)

		// Test min/max timestamp range
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, 0, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsDesc(t, dispersals)
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, math.MaxInt64, 0, false)
		require.NoError(t, err)
		require.Equal(t, numRequests, len(dispersals))
		checkDispersalsDesc(t, dispersals)
	})

	// Test range boundaries
	t.Run("ascending range boundaries", func(t *testing.T) {
		// Test exclusive start
		dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs, now, 0, true)
		require.NoError(t, err)
		require.Equal(t, numRequests-1, len(dispersals))
		assert.Equal(t, dispersedAt[1], dispersals[0].DispersedAt)
		assert.Equal(t, dispersedAt[numRequests-1], dispersals[numRequests-2].DispersedAt)
		checkDispersalsAsc(t, dispersals)

		// Test exclusive end
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, dispersedAt[4], 0, true)
		require.NoError(t, err)
		require.Equal(t, 4, len(dispersals))
		assert.Equal(t, dispersedAt[0], dispersals[0].DispersedAt)
		assert.Equal(t, dispersedAt[3], dispersals[3].DispersedAt)
		checkDispersalsAsc(t, dispersals)
	})

	// Test range boundaries
	t.Run("descending range boundaries", func(t *testing.T) {
		// Test exclusive start
		dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs, now, 0, false)
		require.NoError(t, err)
		require.Equal(t, numRequests-1, len(dispersals))
		assert.Equal(t, dispersedAt[numRequests-1], dispersals[0].DispersedAt)
		assert.Equal(t, dispersedAt[1], dispersals[numRequests-2].DispersedAt)
		checkDispersalsDesc(t, dispersals)

		// Test exclusive end
		dispersals, err = blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, firstRequestTs-1, dispersedAt[4], 0, false)
		require.NoError(t, err)
		require.Equal(t, 4, len(dispersals))
		assert.Equal(t, dispersedAt[3], dispersals[0].DispersedAt)
		assert.Equal(t, dispersedAt[0], dispersals[3].DispersedAt)
		checkDispersalsDesc(t, dispersals)
	})

	// Test pagination
	t.Run("pagination", func(t *testing.T) {
		for i := 1; i < numRequests; i++ {
			dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, dispersedAt[i-1], dispersedAt[i]+1, 0, true)
			require.NoError(t, err)
			require.Equal(t, 1, len(dispersals))
			assert.Equal(t, dispersedAt[i], dispersals[0].DispersedAt)
		}

		for i := 1; i < numRequests; i++ {
			dispersals, err := blobMetadataStore.GetDispersalRequestByDispersedAt(ctx, opID, dispersedAt[i-1], dispersedAt[i]+1, 0, false)
			require.NoError(t, err)
			require.Equal(t, 1, len(dispersals))
			assert.Equal(t, dispersedAt[i], dispersals[0].DispersedAt)
		}
	})
}

func TestBlobMetadataStoreBatch(t *testing.T) {
	ctx := context.Background()
	_, blobHeader := newBlob(t)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		Signature:  []byte("signature"),
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}

	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 1024,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)

	batch := &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: []*corev2.BlobCertificate{blobCert},
	}
	err = blobMetadataStore.PutBatch(ctx, batch)
	require.NoError(t, err)

	b, err := blobMetadataStore.GetBatch(ctx, bhh)
	require.NoError(t, err)
	assert.Equal(t, batch, b)
}

func TestBlobMetadataStoreBlobAttestationInfo(t *testing.T) {
	ctx := context.Background()
	blobKey := corev2.BlobKey{1, 1, 1}
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 1024,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)
	err = blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	assert.NoError(t, err)

	inclusionInfo := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      10,
		InclusionProof: []byte("proof"),
	}
	err = blobMetadataStore.PutBlobInclusionInfo(ctx, inclusionInfo)
	assert.NoError(t, err)

	// Test 1: the batch isn't signed yet, so there is no attestation info
	_, err = blobMetadataStore.GetBlobAttestationInfo(ctx, blobKey)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "no attestation info found"))

	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	apk := keyPair.GetPubKeyG2()
	attestation := &corev2.Attestation{
		BatchHeader: batchHeader,
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

	// Test 2: the batch is signed, so we can fetch blob's attestation info
	blobAttestationInfo, err := blobMetadataStore.GetBlobAttestationInfo(ctx, blobKey)
	require.NoError(t, err)
	assert.Equal(t, inclusionInfo, blobAttestationInfo.InclusionInfo)
	assert.Equal(t, attestation, blobAttestationInfo.Attestation)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
		},
	})
}

func TestBlobMetadataStoreInclusionInfo(t *testing.T) {
	ctx := context.Background()
	blobKey := corev2.BlobKey{1, 1, 1}
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)
	inclusionInfo := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      10,
		InclusionProof: []byte("proof"),
	}

	err = blobMetadataStore.PutBlobInclusionInfo(ctx, inclusionInfo)
	assert.NoError(t, err)

	fetchedInfo, err := blobMetadataStore.GetBlobInclusionInfo(ctx, blobKey, bhh)
	assert.NoError(t, err)
	assert.Equal(t, inclusionInfo, fetchedInfo)

	// attempt to put inclusion info with the same key should fail
	err = blobMetadataStore.PutBlobInclusionInfo(ctx, inclusionInfo)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// put multiple inclusion infos
	blobKey1 := corev2.BlobKey{2, 2, 2}
	inclusionInfo1 := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey1,
		BlobIndex:      12,
		InclusionProof: []byte("proof 1"),
	}
	blobKey2 := corev2.BlobKey{3, 3, 3}
	inclusionInfo2 := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey2,
		BlobIndex:      14,
		InclusionProof: []byte("proof 2"),
	}
	err = blobMetadataStore.PutBlobInclusionInfos(ctx, []*corev2.BlobInclusionInfo{inclusionInfo1, inclusionInfo2})
	assert.NoError(t, err)

	// test retries
	nonTransientError := errors.New("non transient error")
	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, nonTransientError).Once()
	err = mockedBlobMetadataStore.PutBlobInclusionInfos(ctx, []*corev2.BlobInclusionInfo{inclusionInfo1, inclusionInfo2})
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
	err = mockedBlobMetadataStore.PutBlobInclusionInfos(ctx, []*corev2.BlobInclusionInfo{inclusionInfo1, inclusionInfo2})
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

	// attempt to retrieve batch header and attestation at the same time
	fetchedHeader, fetchedAttestation, err = blobMetadataStore.GetSignedBatch(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)
	assert.Equal(t, attestation, fetchedAttestation)

	// overwrite existing attestation
	updatedAttestation := &corev2.Attestation{
		BatchHeader: h,
		AttestedAt:  uint64(time.Now().UnixNano()),
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(1), big.NewInt(2)),
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
			1: 90,
		},
	}

	err = blobMetadataStore.PutAttestation(ctx, updatedAttestation)
	assert.NoError(t, err)
	fetchedAttestation, err = blobMetadataStore.GetAttestation(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, updatedAttestation, fetchedAttestation)

	fetchedHeader, fetchedAttestation, err = blobMetadataStore.GetSignedBatch(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)
	assert.Equal(t, updatedAttestation, fetchedAttestation)

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
	accountID := gethcommon.HexToAddress(hex.EncodeToString(accountBytes))
	timestamp := time.Now().UnixNano()
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
			Timestamp:         timestamp,
			CumulativePayment: cumulativePayment,
		},
	}
	bk, err := bh.BlobKey()
	require.NoError(t, err)
	return bk, bh
}

func TestCheckBlobExists(t *testing.T) {
	ctx := context.Background()
	// Create a test blob
	blobKey, blobHeader := newBlob(t)

	// Check that the blob does not exist initially
	exists, err := blobMetadataStore.CheckBlobExists(ctx, blobKey)
	require.NoError(t, err)
	require.False(t, exists, "Blob should not exist before being added")

	// Create blob metadata
	blobMetadata := &v2.BlobMetadata{
		BlobHeader:  blobHeader,
		Signature:   []byte("test-signature"),
		BlobStatus:  v2.Queued,
		Expiry:      uint64(time.Now().Add(time.Hour).Unix()),
		NumRetries:  0,
		BlobSize:    1024,
		RequestedAt: uint64(time.Now().UnixNano()),
		UpdatedAt:   uint64(time.Now().UnixNano()),
	}

	// Store the blob metadata
	err = blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
	require.NoError(t, err)

	// Check that the blob now exists
	exists, err = blobMetadataStore.CheckBlobExists(ctx, blobKey)
	require.NoError(t, err)
	require.True(t, exists, "Blob should exist after being added")

	// Delete the blob metadata
	err = blobMetadataStore.DeleteBlobMetadata(ctx, blobKey)
	require.NoError(t, err)

	// Check that the blob no longer exists
	exists, err = blobMetadataStore.CheckBlobExists(ctx, blobKey)
	require.NoError(t, err)
	require.False(t, exists, "Blob should not exist after being deleted")

	// Test with non-existent blob key
	randomKey := corev2.BlobKey{}
	_, err = rand.Read(randomKey[:])
	require.NoError(t, err)

	exists, err = blobMetadataStore.CheckBlobExists(ctx, randomKey)
	require.NoError(t, err)
	require.False(t, exists, "Random blob key should not exist")
}
