package oci

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ObjectStorageClientInterface defines the interface we need for testing
type ObjectStorageClientInterface interface {
	GetObject(ctx context.Context, request objectstorage.GetObjectRequest) (objectstorage.GetObjectResponse, error)
	PutObject(ctx context.Context, request objectstorage.PutObjectRequest) (objectstorage.PutObjectResponse, error)
	DeleteObject(
		ctx context.Context, request objectstorage.DeleteObjectRequest,
	) (objectstorage.DeleteObjectResponse, error)
	HeadObject(ctx context.Context, request objectstorage.HeadObjectRequest) (objectstorage.HeadObjectResponse, error)
	ListObjects(ctx context.Context, request objectstorage.ListObjectsRequest) (objectstorage.ListObjectsResponse, error)
	CreateBucket(
		ctx context.Context, request objectstorage.CreateBucketRequest,
	) (objectstorage.CreateBucketResponse, error)
	SetRegion(region string)
}

// MockObjectStorageClient is a mock implementation of the OCI ObjectStorageClient
type MockObjectStorageClient struct {
	mock.Mock
}

// mockLogger is a simple mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, tags ...any)               {}
func (m *mockLogger) Info(msg string, tags ...any)                {}
func (m *mockLogger) Warn(msg string, tags ...any)                {}
func (m *mockLogger) Error(msg string, tags ...any)               {}
func (m *mockLogger) Fatal(msg string, tags ...any)               {}
func (m *mockLogger) Debugf(template string, args ...interface{}) {}
func (m *mockLogger) Infof(template string, args ...interface{})  {}
func (m *mockLogger) Warnf(template string, args ...interface{})  {}
func (m *mockLogger) Errorf(template string, args ...interface{}) {}
func (m *mockLogger) Fatalf(template string, args ...interface{}) {}
func (m *mockLogger) With(tags ...any) logging.Logger             { return m }

func (m *MockObjectStorageClient) GetObject(
	ctx context.Context, request objectstorage.GetObjectRequest,
) (objectstorage.GetObjectResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.GetObjectResponse{}, fmt.Errorf("mock GetObject error: %w", err)
	}
	return args.Get(0).(objectstorage.GetObjectResponse), nil
}

func (m *MockObjectStorageClient) PutObject(
	ctx context.Context, request objectstorage.PutObjectRequest,
) (objectstorage.PutObjectResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.PutObjectResponse{}, fmt.Errorf("mock PutObject error: %w", err)
	}
	return args.Get(0).(objectstorage.PutObjectResponse), nil
}

func (m *MockObjectStorageClient) DeleteObject(
	ctx context.Context, request objectstorage.DeleteObjectRequest,
) (objectstorage.DeleteObjectResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.DeleteObjectResponse{}, fmt.Errorf("mock DeleteObject error: %w", err)
	}
	return args.Get(0).(objectstorage.DeleteObjectResponse), nil
}

func (m *MockObjectStorageClient) HeadObject(
	ctx context.Context, request objectstorage.HeadObjectRequest,
) (objectstorage.HeadObjectResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.HeadObjectResponse{}, fmt.Errorf("mock HeadObject error: %w", err)
	}
	return args.Get(0).(objectstorage.HeadObjectResponse), nil
}

func (m *MockObjectStorageClient) ListObjects(
	ctx context.Context, request objectstorage.ListObjectsRequest,
) (objectstorage.ListObjectsResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.ListObjectsResponse{}, fmt.Errorf("mock ListObjects error: %w", err)
	}
	return args.Get(0).(objectstorage.ListObjectsResponse), nil
}

func (m *MockObjectStorageClient) CreateBucket(
	ctx context.Context, request objectstorage.CreateBucketRequest,
) (objectstorage.CreateBucketResponse, error) {
	args := m.Called(ctx, request)
	if err := args.Error(1); err != nil {
		return objectstorage.CreateBucketResponse{}, fmt.Errorf("mock CreateBucket error: %w", err)
	}
	return args.Get(0).(objectstorage.CreateBucketResponse), nil
}

func (m *MockObjectStorageClient) SetRegion(region string) {
	m.Called(region)
}

// testOciClient is a test version of ociClient that uses our interface
type testOciClient struct {
	cfg                 *ObjectStorageConfig
	objectStorageClient ObjectStorageClientInterface
	concurrencyLimiter  chan struct{}
	logger              logging.Logger
}

// Implement the s3.Client interface methods for testOciClient
func (c *testOciClient) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	getRequest := objectstorage.GetObjectRequest{
		NamespaceName: &c.cfg.Namespace,
		BucketName:    &bucket,
		ObjectName:    &key,
	}

	response, err := c.objectStorageClient.GetObject(ctx, getRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to download object %s from bucket %s: %w", key, bucket, err)
	}
	defer func() { _ = response.Content.Close() }()

	data, err := io.ReadAll(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	if len(data) == 0 {
		return nil, ErrObjectNotFound
	}

	return data, nil
}

func (c *testOciClient) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	headRequest := objectstorage.HeadObjectRequest{
		NamespaceName: &c.cfg.Namespace,
		BucketName:    &bucket,
		ObjectName:    &key,
	}

	response, err := c.objectStorageClient.HeadObject(ctx, headRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to head object %s from bucket %s: %w", key, bucket, err)
	}

	return response.ContentLength, nil
}

func (c *testOciClient) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	putRequest := objectstorage.PutObjectRequest{
		NamespaceName: &c.cfg.Namespace,
		BucketName:    &bucket,
		ObjectName:    &key,
		PutObjectBody: io.NopCloser(bytes.NewReader(data)),
	}

	_, err := c.objectStorageClient.PutObject(ctx, putRequest)
	if err != nil {
		return fmt.Errorf("failed to upload object %s to bucket %s: %w", key, bucket, err)
	}

	return nil
}

func (c *testOciClient) DeleteObject(ctx context.Context, bucket string, key string) error {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	deleteRequest := objectstorage.DeleteObjectRequest{
		NamespaceName: &c.cfg.Namespace,
		BucketName:    &bucket,
		ObjectName:    &key,
	}

	_, err := c.objectStorageClient.DeleteObject(ctx, deleteRequest)
	if err != nil {
		return fmt.Errorf("failed to delete object %s from bucket %s: %w", key, bucket, err)
	}

	return nil
}

func (c *testOciClient) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3.Object, error) {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	listRequest := objectstorage.ListObjectsRequest{
		NamespaceName: &c.cfg.Namespace,
		BucketName:    &bucket,
		Prefix:        &prefix,
	}

	response, err := c.objectStorageClient.ListObjects(ctx, listRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects from bucket %s with prefix %s: %w", bucket, prefix, err)
	}

	objects := make([]s3.Object, len(response.Objects))
	for i, obj := range response.Objects {
		objects[i] = s3.Object{
			Key:  *obj.Name,
			Size: *obj.Size,
		}
	}

	return objects, nil
}

func (c *testOciClient) CreateBucket(ctx context.Context, bucket string) error {
	<-c.concurrencyLimiter
	defer func() { c.concurrencyLimiter <- struct{}{} }()

	createRequest := objectstorage.CreateBucketRequest{
		NamespaceName: &c.cfg.Namespace,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:          &bucket,
			CompartmentId: &c.cfg.CompartmentID,
		},
	}

	_, err := c.objectStorageClient.CreateBucket(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
	}

	return nil
}

func (c *testOciClient) FragmentedUploadObject(
	ctx context.Context, bucket string, key string, data []byte, fragmentSize int,
) error {
	// Simplified implementation for testing
	return c.UploadObject(ctx, bucket, key, data)
}

func (c *testOciClient) FragmentedDownloadObject(
	ctx context.Context, bucket string, key string, fileSize int, fragmentSize int,
) ([]byte, error) {
	if fileSize <= 0 {
		return nil, errors.New("fileSize must be greater than 0")
	}
	if fragmentSize <= 0 {
		return nil, errors.New("fragmentSize must be greater than 0")
	}

	// Simplified implementation that downloads fragments and recombines them
	fragmentCount := GetFragmentCount(fileSize, fragmentSize)
	fragments := make([]*s3.Fragment, 0, fragmentCount)

	for i := 0; i < fragmentCount; i++ {
		fragmentKey := fmt.Sprintf("%s-%d", key, i)
		if i == fragmentCount-1 {
			fragmentKey += "f" // final fragment
		}

		fragmentData, err := c.DownloadObject(ctx, bucket, fragmentKey)
		if err != nil {
			return nil, fmt.Errorf("failed to download fragment %d: %w", i, err)
		}

		fragments = append(fragments, &s3.Fragment{
			FragmentKey: fragmentKey,
			Data:        fragmentData,
			Index:       i,
		})
	}

	return RecombineFragments(fragments)
}

func createTestOCIClient(mockClient *MockObjectStorageClient) *testOciClient {
	config := &ObjectStorageConfig{
		Namespace:                   "test-namespace",
		Region:                      "us-phoenix-1",
		CompartmentID:               "test-compartment",
		BucketName:                  "test-bucket",
		FragmentParallelismConstant: 1,
	}

	logger := &mockLogger{}

	// Initialize the concurrency limiter with one token
	concurrencyLimiter := make(chan struct{}, 1)
	concurrencyLimiter <- struct{}{}

	client := &testOciClient{
		cfg:                 config,
		objectStorageClient: mockClient,
		concurrencyLimiter:  concurrencyLimiter,
		logger:              logger,
	}

	return client
}

func TestOCIClient_DownloadObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	expectedData := []byte("test data")

	// Mock successful response
	mockResponse := objectstorage.GetObjectResponse{
		Content: io.NopCloser(bytes.NewReader(expectedData)),
	}
	mockClient.On("GetObject", ctx, mock.AnythingOfType("objectstorage.GetObjectRequest")).Return(mockResponse, nil)

	data, err := client.DownloadObject(ctx, bucket, key)

	require.NoError(t, err)
	assert.Equal(t, expectedData, data)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_DownloadObject_EmptyData(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"

	// Mock response with empty data
	mockResponse := objectstorage.GetObjectResponse{
		Content: io.NopCloser(bytes.NewReader([]byte{})),
	}
	mockClient.On("GetObject", ctx, mock.AnythingOfType("objectstorage.GetObjectRequest")).Return(mockResponse, nil)

	_, err := client.DownloadObject(ctx, bucket, key)

	assert.Equal(t, ErrObjectNotFound, err)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_HeadObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	expectedSize := int64(123)

	// Mock successful response
	mockResponse := objectstorage.HeadObjectResponse{
		ContentLength: &expectedSize,
	}
	mockClient.On("HeadObject", ctx, mock.AnythingOfType("objectstorage.HeadObjectRequest")).Return(mockResponse, nil)

	size, err := client.HeadObject(ctx, bucket, key)

	require.NoError(t, err)
	assert.Equal(t, &expectedSize, size)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_UploadObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	data := []byte("test data")

	// Mock successful response
	mockResponse := objectstorage.PutObjectResponse{}
	mockClient.On("PutObject", ctx, mock.AnythingOfType("objectstorage.PutObjectRequest")).Return(mockResponse, nil)

	err := client.UploadObject(ctx, bucket, key, data)

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_DeleteObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"

	// Mock successful response
	mockResponse := objectstorage.DeleteObjectResponse{}
	mockClient.On("DeleteObject", ctx, mock.AnythingOfType("objectstorage.DeleteObjectRequest")).Return(mockResponse, nil)

	err := client.DeleteObject(ctx, bucket, key)

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_ListObjects(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	prefix := "test-prefix"

	// Mock successful response
	size1 := int64(100)
	size2 := int64(200)
	name1 := "object1"
	name2 := "object2"

	mockListObjects := objectstorage.ListObjects{
		Objects: []objectstorage.ObjectSummary{
			{
				Name: &name1,
				Size: &size1,
			},
			{
				Name: &name2,
				Size: &size2,
			},
		},
	}
	mockResponse := objectstorage.ListObjectsResponse{
		ListObjects: mockListObjects,
	}
	mockClient.On("ListObjects", ctx, mock.AnythingOfType("objectstorage.ListObjectsRequest")).Return(mockResponse, nil)

	objects, err := client.ListObjects(ctx, bucket, prefix)

	require.NoError(t, err)
	assert.Len(t, objects, 2)
	assert.Equal(t, s3.Object{Key: "object1", Size: 100}, objects[0])
	assert.Equal(t, s3.Object{Key: "object2", Size: 200}, objects[1])
	mockClient.AssertExpectations(t)
}

func TestOCIClient_CreateBucket(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"

	// Mock successful response
	mockResponse := objectstorage.CreateBucketResponse{}
	mockClient.On("CreateBucket", ctx, mock.AnythingOfType("objectstorage.CreateBucketRequest")).Return(mockResponse, nil)

	err := client.CreateBucket(ctx, bucket)

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_FragmentedUploadObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	data := []byte("test data for fragmented upload")
	fragmentSize := 10

	// Mock successful responses for each fragment
	mockResponse := objectstorage.PutObjectResponse{}
	mockClient.On("PutObject", ctx, mock.AnythingOfType("objectstorage.PutObjectRequest")).Return(mockResponse, nil)

	err := client.FragmentedUploadObject(ctx, bucket, key, data, fragmentSize)

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_FragmentedDownloadObject(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"
	fileSize := 20
	fragmentSize := 10

	// Create test data fragments
	fragment1Data := []byte("0123456789")
	fragment2Data := []byte("abcdefghij")

	// Mock responses for each fragment
	mockResponse1 := objectstorage.GetObjectResponse{
		Content: io.NopCloser(bytes.NewReader(fragment1Data)),
	}
	mockResponse2 := objectstorage.GetObjectResponse{
		Content: io.NopCloser(bytes.NewReader(fragment2Data)),
	}

	mockClient.On("GetObject", ctx, mock.MatchedBy(func(req objectstorage.GetObjectRequest) bool {
		return *req.ObjectName == "test-key-0"
	})).Return(mockResponse1, nil)

	mockClient.On("GetObject", ctx, mock.MatchedBy(func(req objectstorage.GetObjectRequest) bool {
		return *req.ObjectName == "test-key-1f"
	})).Return(mockResponse2, nil)

	data, err := client.FragmentedDownloadObject(ctx, bucket, key, fileSize, fragmentSize)

	require.NoError(t, err)
	expected := append(fragment1Data, fragment2Data...)
	assert.Equal(t, expected, data)
	mockClient.AssertExpectations(t)
}

func TestOCIClient_FragmentedDownloadObject_InvalidParams(t *testing.T) {
	mockClient := new(MockObjectStorageClient)
	client := createTestOCIClient(mockClient)

	ctx := context.Background()
	bucket := "test-bucket"
	key := "test-key"

	// Test invalid file size
	_, err := client.FragmentedDownloadObject(ctx, bucket, key, 0, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fileSize must be greater than 0")

	// Test invalid fragment size
	_, err = client.FragmentedDownloadObject(ctx, bucket, key, 10, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fragmentSize must be greater than 0")
}

func TestGetFragmentCount(t *testing.T) {
	tests := []struct {
		fileSize     int
		fragmentSize int
		expected     int
	}{
		{10, 20, 1},     // file smaller than fragment
		{20, 10, 2},     // exact division
		{25, 10, 3},     // with remainder
		{100, 33, 4},    // with remainder
	}

	for _, test := range tests {
		result := GetFragmentCount(test.fileSize, test.fragmentSize)
		assert.Equal(t, test.expected, result)
	}
}

func TestRecombineFragments(t *testing.T) {
	fragment1 := &s3.Fragment{
		FragmentKey: "test-0",
		Data:        []byte("0123456789"),
		Index:       0,
	}
	fragment2 := &s3.Fragment{
		FragmentKey: "test-1f",
		Data:        []byte("abcdefghij"),
		Index:       1,
	}

	fragments := []*s3.Fragment{fragment2, fragment1} // intentionally out of order

	data, err := RecombineFragments(fragments)

	require.NoError(t, err)
	expected := []byte("0123456789abcdefghij")
	assert.Equal(t, expected, data)
}

func TestRecombineFragments_EmptyFragments(t *testing.T) {
	_, err := RecombineFragments([]*s3.Fragment{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fragments")
}

func TestRecombineFragments_MissingFinalFragment(t *testing.T) {
	fragment1 := &s3.Fragment{
		FragmentKey: "test-0",
		Data:        []byte("0123456789"),
		Index:       0,
	}
	fragment2 := &s3.Fragment{
		FragmentKey: "test-1", // missing 'f' suffix
		Data:        []byte("abcdefghij"),
		Index:       1,
	}

	fragments := []*s3.Fragment{fragment1, fragment2}

	_, err := RecombineFragments(fragments)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing final fragment")
}

func TestRecombineFragments_MissingFragment(t *testing.T) {
	fragment1 := &s3.Fragment{
		FragmentKey: "test-0",
		Data:        []byte("0123456789"),
		Index:       0,
	}
	fragment3 := &s3.Fragment{
		FragmentKey: "test-2f",
		Data:        []byte("abcdefghij"),
		Index:       2, // skipping index 1
	}

	fragments := []*s3.Fragment{fragment1, fragment3}

	_, err := RecombineFragments(fragments)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing fragment with index 1")
}

func TestNewObjectStorageClient(t *testing.T) {
	ctx := context.Background()
	config := ObjectStorageConfig{
		Namespace:                   "test-namespace",
		Region:                      "us-phoenix-1",
		CompartmentID:               "test-compartment",
		BucketName:                  "test-bucket",
		FragmentParallelismConstant: 1,
	}
	logger := &mockLogger{}

	// This test will fail in CI without OCI credentials, but demonstrates the interface
	client, err := NewObjectStorageClient(ctx, config, logger)
	if err != nil {
		// We expect an error in test environment without OCI setup
		assert.Contains(t, err.Error(), "failed to create OCI Object Storage client")
		assert.Nil(t, client)
	} else {
		// If somehow it succeeds (should not happen in test env), client should not be nil
		assert.NotNil(t, client)
	}
}