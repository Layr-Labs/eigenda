package oci

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"

	s3common "github.com/Layr-Labs/eigenda/common/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
	oraclecommon "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

// ObjectStorageConfig holds configuration for OCI Object Storage
type ObjectStorageConfig struct {
	Namespace                   string
	Region                      string
	CompartmentID               string
	BucketName                  string
	FragmentParallelismConstant int
	FragmentParallelismFactor   int
}

// ociS3Client implements the S3 Client interface using OCI Object Storage
type ociS3Client struct {
	cfg                 *ObjectStorageConfig
	objectStorageClient objectstorage.ObjectStorageClient

	// concurrencyLimiter is a channel that limits the number of concurrent operations.
	concurrencyLimiter chan struct{}

	logger logging.Logger
}

var _ s3common.S3Client = (*ociS3Client)(nil)

// NewOciS3Client creates a new OCI Object Storage client that implements the S3 Client interface
func NewOciS3Client(
	ctx context.Context,
	cfg ObjectStorageConfig,
	logger logging.Logger) (s3common.S3Client, error) {

	// Create OCI configuration provider using workload identity
	configProvider, err := auth.OkeWorkloadIdentityConfigurationProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI Object Storage client: %w", err)
	}

	// Create Object Storage client
	objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI Object Storage client: %w", err)
	}

	// Get namespace dynamically if not provided in config
	finalCfg := cfg
	if finalCfg.Namespace == "" {
		namespaceReq := objectstorage.GetNamespaceRequest{}
		namespaceResp, err := objectStorageClient.GetNamespace(ctx, namespaceReq)
		if err != nil {
			return nil, fmt.Errorf("failed to get OCI namespace: %w", err)
		}
		finalCfg.Namespace = *namespaceResp.Value
		logger.Info("Retrieved OCI namespace dynamically", "namespace", finalCfg.Namespace)
	}

	// Set region
	if finalCfg.Region != "" {
		objectStorageClient.SetRegion(finalCfg.Region)
	}

	// Calculate workers for concurrency
	workers := 0
	if cfg.FragmentParallelismConstant > 0 {
		workers = cfg.FragmentParallelismConstant
	}
	if cfg.FragmentParallelismFactor > 0 {
		workers = cfg.FragmentParallelismFactor * runtime.NumCPU()
	}

	if workers == 0 {
		workers = 1
	}

	// Initialize concurrency limiter with tokens
	limiter := make(chan struct{}, workers)
	for i := 0; i < workers; i++ {
		limiter <- struct{}{}
	}
	return &ociS3Client{
		cfg:                 &finalCfg,
		objectStorageClient: objectStorageClient,
		concurrencyLimiter:  limiter,
		logger:              logger.With("component", "OCIObjectStorageClient"),
	}, nil
}

// NOTE: The methods below have 0% test coverage because they all require live OCI credentials
// and network access to Oracle Cloud. We could refactor to use dependency injection with
// interfaces, but that adds complexity for minimal benefit since these are just thin wrappers
// around the OCI SDK. The utility functions (GetFragmentCount, RecombineFragments) and
// config processing in NewObjectStorageClient have good coverage where it matters.

func (c *ociS3Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, bool, error) {
	getObjectRequest := objectstorage.GetObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	response, err := c.objectStorageClient.GetObject(ctx, getObjectRequest)
	if err != nil {
		if response.RawResponse != nil && response.RawResponse.StatusCode == 404 {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get object from OCI: %w", err)
	}
	defer func() {
		if closeErr := response.Content.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	data, err := io.ReadAll(response.Content)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read object content: %w", err)
	}

	if len(data) == 0 {
		return nil, false, nil
	}

	return data, true, nil
}

func (c *ociS3Client) DownloadPartialObject(
	ctx context.Context,
	bucket string,
	key string,
	startIndex int64,
	endIndex int64,
) ([]byte, bool, error) {

	if startIndex < 0 || endIndex <= startIndex {
		return nil, false, fmt.Errorf("invalid startIndex (%d) or endIndex (%d)", startIndex, endIndex)
	}

	rangeString := fmt.Sprintf("bytes=%d-%d", startIndex, endIndex-1)

	getObjectRequest := objectstorage.GetObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
		Range:         oraclecommon.String(rangeString),
	}

	response, err := c.objectStorageClient.GetObject(ctx, getObjectRequest)
	if err != nil {
		if response.RawResponse != nil && response.RawResponse.StatusCode == 404 {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get object from OCI: %w", err)
	}
	defer func() {
		if closeErr := response.Content.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	data, err := io.ReadAll(response.Content)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read object content: %w", err)
	}

	if len(data) == 0 {
		return nil, false, nil
	}

	return data, true, nil
}

func (c *ociS3Client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	headObjectRequest := objectstorage.HeadObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	response, err := c.objectStorageClient.HeadObject(ctx, headObjectRequest)
	if err != nil {
		// Check if it's a 404 error
		if response.RawResponse != nil && response.RawResponse.StatusCode == 404 {
			return nil, s3common.ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to head object: %w", err)
	}

	return response.ContentLength, nil
}

func (c *ociS3Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	putObjectRequest := objectstorage.PutObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
		PutObjectBody: io.NopCloser(bytes.NewReader(data)),
		ContentLength: oraclecommon.Int64(int64(len(data))),
	}

	_, err := c.objectStorageClient.PutObject(ctx, putObjectRequest)
	if err != nil {
		return fmt.Errorf("failed to put object to OCI: %w", err)
	}

	return nil
}

func (c *ociS3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	deleteObjectRequest := objectstorage.DeleteObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	_, err := c.objectStorageClient.DeleteObject(ctx, deleteObjectRequest)
	if err != nil {
		return fmt.Errorf("failed to delete object from OCI: %w", err)
	}

	return nil
}

func (c *ociS3Client) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3common.ListedObject, error) {
	listObjectsRequest := objectstorage.ListObjectsRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		Prefix:        oraclecommon.String(prefix),
		Limit:         oraclecommon.Int(1000), // Match S3 behavior of up to 1000 items
	}

	response, err := c.objectStorageClient.ListObjects(ctx, listObjectsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects from OCI: %w", err)
	}

	objects := make([]s3common.ListedObject, 0, len(response.Objects))
	for _, object := range response.Objects {
		var size int64 = 0
		if object.Size != nil {
			size = *object.Size
		}
		var key string
		if object.Name != nil {
			key = *object.Name
		}
		objects = append(objects, s3common.ListedObject{
			Key:  key,
			Size: size,
		})
	}

	return objects, nil
}

func (c *ociS3Client) CreateBucket(ctx context.Context, bucket string) error {
	createBucketRequest := objectstorage.CreateBucketRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:             oraclecommon.String(bucket),
			CompartmentId:    oraclecommon.String(c.cfg.CompartmentID),
			PublicAccessType: objectstorage.CreateBucketDetailsPublicAccessTypeNopublicaccess,
		},
	}

	_, err := c.objectStorageClient.CreateBucket(ctx, createBucketRequest)
	if err != nil {
		return fmt.Errorf("failed to create bucket in OCI: %w", err)
	}

	return nil
}
