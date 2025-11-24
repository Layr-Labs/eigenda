package v2_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	serverv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

// TestBlobFeedDelay verifies that blobs younger than feedDelay are not returned
func TestBlobFeedDelay(t *testing.T) {
	if blobMetadataStore == nil {
		t.Skip("Test dependencies are not initialized")
	}

	ctx := context.Background()
	now := time.Now()
	feedDelay := 10 * time.Second

	// Create test blobs: some old (should appear) and some recent (should not appear)
	oldBlobTime := now.Add(-feedDelay - time.Second) // 11 seconds ago
	recentBlobTime := now.Add(-feedDelay/2)          // 5 seconds ago

	oldBlobHeader := makeBlobHeaderV2(t)
	recentBlobHeader := makeBlobHeaderV2(t)

	oldBlobKey, err := oldBlobHeader.BlobKey()
	require.NoError(t, err)
	recentBlobKey, err := recentBlobHeader.BlobKey()
	require.NoError(t, err)

	oldBlobMetadata := &commonv2.BlobMetadata{
		BlobHeader:  oldBlobHeader,
		Signature:   []byte{0, 1, 2, 3, 4},
		BlobStatus:  commonv2.Encoded,
		Expiry:      uint64(now.Add(time.Hour).Unix()),
		NumRetries:  0,
		UpdatedAt:   uint64(oldBlobTime.UnixNano()),
		RequestedAt: uint64(oldBlobTime.UnixNano()),
	}

	recentBlobMetadata := &commonv2.BlobMetadata{
		BlobHeader:  recentBlobHeader,
		Signature:   []byte{5, 6, 7, 8, 9},
		BlobStatus:  commonv2.Encoded,
		Expiry:      uint64(now.Add(time.Hour).Unix()),
		NumRetries:  0,
		UpdatedAt:   uint64(recentBlobTime.UnixNano()),
		RequestedAt: uint64(recentBlobTime.UnixNano()),
	}

	err = blobMetadataStore.PutBlobMetadata(ctx, oldBlobMetadata)
	require.NoError(t, err)
	err = blobMetadataStore.PutBlobMetadata(ctx, recentBlobMetadata)
	require.NoError(t, err)

	// Clean up after test
	defer deleteItems(t, []dynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + oldBlobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + recentBlobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})

	// Create a server with the feed delay configured
	configWithDelay := dataapi.Config{
		ServerMode:        "test",
		SocketAddr:        ":8080",
		AllowOrigins:      []string{"*"},
		DisperserHostname: "localhost:32007",
		ChurnerHostname:   "localhost:32009",
		FeedDelay:         feedDelay,
	}

	metrics := dataapi.NewMetrics(serverVersion, prometheus.NewRegistry(), blobMetadataStore, "9001", logger)
	testServer, err := serverv2.NewServerV2(
		configWithDelay, blobMetadataStore, prometheusClient, subgraphClient,
		mockTx, mockChainState, mockIndexedChainState, logger, metrics)
	require.NoError(t, err)

	// Set up router and request
	r := gin.Default()
	r.GET("/v2/blobs/feed", testServer.FetchBlobFeed)

	req := httptest.NewRequest("GET", "/v2/blobs/feed", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Verify response
	require.Equal(t, http.StatusOK, rr.Code)

	var response serverv2.BlobFeedResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify that only the old blob is returned
	require.Len(t, response.Blobs, 1, "Expected exactly 1 blob (the old one)")
	require.Equal(t, oldBlobKey.Hex(), response.Blobs[0].BlobKey, "Expected the old blob to be returned")

	// Verify that all returned blobs have RequestedAt timestamp older than (now - feedDelay)
	minTime := now.Add(-feedDelay)
	for _, blob := range response.Blobs {
		requestedAt := time.Unix(0, int64(blob.BlobMetadata.RequestedAt))
		require.True(t, requestedAt.Before(minTime) || requestedAt.Equal(minTime),
			"Blob %s has RequestedAt=%v which is newer than allowed minTime=%v",
			blob.BlobKey, requestedAt, minTime)
	}
}

// TestBatchFeedDelay verifies that batches younger than feedDelay are not returned
func TestBatchFeedDelay(t *testing.T) {
	if blobMetadataStore == nil {
		t.Skip("Test dependencies are not initialized")
	}

	ctx := context.Background()
	now := time.Now()
	feedDelay := 10 * time.Second

	// Create test batches: some old (should appear) and some recent (should not appear)
	oldBatchTime := now.Add(-feedDelay - time.Second) // 11 seconds ago
	recentBatchTime := now.Add(-feedDelay/2)          // 5 seconds ago

	oldBatchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{7, 8, 9},
		ReferenceBlockNumber: 100,
	}
	recentBatchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{10, 11, 12},
		ReferenceBlockNumber: 101,
	}

	oldBatchHash, err := oldBatchHeader.Hash()
	require.NoError(t, err)
	recentBatchHash, err := recentBatchHeader.Hash()
	require.NoError(t, err)

	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	apk := keyPair.GetPubKeyG2()

	oldAttestation := &corev2.Attestation{
		BatchHeader: oldBatchHeader,
		AttestedAt:  uint64(oldBatchTime.UnixNano()),
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
			1: 80,
		},
	}

	recentAttestation := &corev2.Attestation{
		BatchHeader: recentBatchHeader,
		AttestedAt:  uint64(recentBatchTime.UnixNano()),
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		},
		APKG2: apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(11), big.NewInt(12)),
			1: core.NewG1Point(big.NewInt(13), big.NewInt(14)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(15), big.NewInt(16)),
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}

	err = blobMetadataStore.PutAttestation(ctx, oldAttestation)
	require.NoError(t, err)
	err = blobMetadataStore.PutAttestation(ctx, recentAttestation)
	require.NoError(t, err)

	// Clean up after test
	defer deleteItems(t, []dynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(oldBatchHash[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(recentBatchHash[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
	})

	// Create a server with the feed delay configured
	configWithDelay := dataapi.Config{
		ServerMode:        "test",
		SocketAddr:        ":8080",
		AllowOrigins:      []string{"*"},
		DisperserHostname: "localhost:32007",
		ChurnerHostname:   "localhost:32009",
		FeedDelay:         feedDelay,
	}

	metrics := dataapi.NewMetrics(serverVersion, prometheus.NewRegistry(), blobMetadataStore, "9001", logger)
	testServer, err := serverv2.NewServerV2(
		configWithDelay, blobMetadataStore, prometheusClient, subgraphClient,
		mockTx, mockChainState, mockIndexedChainState, logger, metrics)
	require.NoError(t, err)

	// Set up router and request
	r := gin.Default()
	r.GET("/v2/batches/feed", testServer.FetchBatchFeed)

	req := httptest.NewRequest("GET", "/v2/batches/feed", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Verify response
	require.Equal(t, http.StatusOK, rr.Code)

	var response serverv2.BatchFeedResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify that only the old batch is returned
	require.Len(t, response.Batches, 1, "Expected exactly 1 batch (the old one)")

	// Verify that all returned batches have AttestedAt timestamp older than (now - feedDelay)
	minTime := now.Add(-feedDelay)
	for _, batch := range response.Batches {
		attestedAt := time.Unix(0, int64(batch.AttestedAt))
		require.True(t, attestedAt.Before(minTime) || attestedAt.Equal(minTime),
			"Batch %s has AttestedAt=%v which is newer than allowed minTime=%v",
			batch.BatchHeaderHash, attestedAt, minTime)
	}
}
