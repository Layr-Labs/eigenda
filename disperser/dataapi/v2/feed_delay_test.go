package v2_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	serverv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// TestBlobFeedDelay verifies that blobs younger than feedDelay are not returned
func TestBlobFeedDelay(t *testing.T) {
	if blobMetadataStore == nil {
		t.Skip("Blob metadata store is not initialized")
	}

	now := time.Now()
	feedDelay := 10 * time.Second
	router := testDataApiServerV2

	// Test that beforeTime is properly capped
	_, err := http.NewRequest("GET", "/v2/blobs/feed", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(rr)
	router.FetchBlobFeed(ginContext)

	// Verify response
	require.Equal(t, http.StatusOK, rr.Code)

	var response serverv2.BlobFeedResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

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
		t.Skip("Blob metadata store is not initialized")
	}

	now := time.Now()
	feedDelay := 10 * time.Second

	_, err := http.NewRequest("GET", "/v2/batches/feed", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := testDataApiServerV2
	ginContext, _ := gin.CreateTestContext(rr)
	router.FetchBatchFeed(ginContext)

	// Verify response
	require.Equal(t, http.StatusOK, rr.Code)

	var response serverv2.BatchFeedResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify that all returned batches have AttestedAt timestamp older than (now - feedDelay)
	minTime := now.Add(-feedDelay)
	for _, batch := range response.Batches {
		attestedAt := time.Unix(0, int64(batch.AttestedAt))
		require.True(t, attestedAt.Before(minTime) || attestedAt.Equal(minTime),
			"Batch %s has AttestedAt=%v which is newer than allowed minTime=%v",
			batch.BatchHeaderHash, attestedAt, minTime)
	}
}
