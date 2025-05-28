package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoogleEndpoint_StorageGoogleapis(t *testing.T) {
	endpoint := "storage.googleapis.com"
	result := isGoogleEndpoint(endpoint)
	assert.True(t, result, "Expected true for Google Cloud Storage endpoint")
}

func TestIsGoogleEndpoint_HttpsStorageGoogleapis(t *testing.T) {
	endpoint := "https://storage.googleapis.com"
	result := isGoogleEndpoint(endpoint)
	assert.True(t, result, "Expected true for Google Cloud Storage endpoint")
}

func TestIsGoogleEndpoint_False(t *testing.T) {
	endpoint := "https://s3.amazonaws.com/my-bucket"
	result := isGoogleEndpoint(endpoint)
	assert.False(t, result, "Expected false for non-Google endpoint")
}

func TestIsGoogleEndpoint_Empty(t *testing.T) {
	endpoint := ""
	result := isGoogleEndpoint(endpoint)
	assert.False(t, result, "Expected false for empty endpoint")
}
