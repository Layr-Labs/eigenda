package oci

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
	oraclecommon "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

var (
	ErrObjectNotFound = errors.New("object not found")
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

// ociClient implements the S3 Client interface using OCI Object Storage
type ociClient struct {
	cfg                 *ObjectStorageConfig
	objectStorageClient objectstorage.ObjectStorageClient

	// concurrencyLimiter is a channel that limits the number of concurrent operations.
	concurrencyLimiter chan struct{}

	logger logging.Logger
}

var _ s3.Client = (*ociClient)(nil)

// NewObjectStorageClient creates a new OCI Object Storage client that implements the S3 Client interface
func NewObjectStorageClient(
	ctx context.Context,
	cfg ObjectStorageConfig,
	logger logging.Logger) (s3.Client, error) {

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

	// Fall back to standard OCI environment variables if application config is empty
	finalCfg := cfg
	if finalCfg.Region == "" {
		if region := os.Getenv("OCI_REGION"); region != "" {
			finalCfg.Region = region
		}
	}
	if finalCfg.CompartmentID == "" {
		if compartmentID := os.Getenv("OCI_COMPARTMENT_ID"); compartmentID != "" {
			finalCfg.CompartmentID = compartmentID
		}
	}
	if finalCfg.Namespace == "" {
		if namespace := os.Getenv("OCI_NAMESPACE"); namespace != "" {
			finalCfg.Namespace = namespace
		} else {
			// Get namespace dynamically if not provided (like in the working example)
			namespaceReq := objectstorage.GetNamespaceRequest{}
			namespaceResp, err := objectStorageClient.GetNamespace(ctx, namespaceReq)
			if err != nil {
				return nil, fmt.Errorf("failed to get OCI namespace: %w", err)
			}
			finalCfg.Namespace = *namespaceResp.Value
			logger.Info("Retrieved OCI namespace dynamically", "namespace", finalCfg.Namespace)
		}
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
	return &ociClient{
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

func (c *ociClient) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	getObjectRequest := objectstorage.GetObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	response, err := c.objectStorageClient.GetObject(ctx, getObjectRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from OCI: %w", err)
	}
	defer func() {
		if closeErr := response.Content.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	data, err := io.ReadAll(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	if len(data) == 0 {
		return nil, ErrObjectNotFound
	}

	return data, nil
}

func (c *ociClient) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	headObjectRequest := objectstorage.HeadObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	response, err := c.objectStorageClient.HeadObject(ctx, headObjectRequest)
	if err != nil {
		// Check if it's a 404 error
		if response.RawResponse != nil && response.RawResponse.StatusCode == 404 {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to head object: %w", err)
	}

	return response.ContentLength, nil
}

func (c *ociClient) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
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

func (c *ociClient) DeleteObject(ctx context.Context, bucket string, key string) error {
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

func (c *ociClient) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3.Object, error) {
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

	objects := make([]s3.Object, 0, len(response.Objects))
	for _, object := range response.Objects {
		var size int64 = 0
		if object.Size != nil {
			size = *object.Size
		}
		var key string
		if object.Name != nil {
			key = *object.Name
		}
		objects = append(objects, s3.Object{
			Key:  key,
			Size: size,
		})
	}

	return objects, nil
}

func (c *ociClient) CreateBucket(ctx context.Context, bucket string) error {
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

func (c *ociClient) FragmentedUploadObject(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	fragmentSize int) error {

	fragments, err := s3.BreakIntoFragments(key, data, fragmentSize)
	if err != nil {
		return fmt.Errorf("failed to break data into fragments: %w", err)
	}
	resultChannel := make(chan error, len(fragments))

	for _, fragment := range fragments {
		fragmentCapture := fragment
		c.concurrencyLimiter <- struct{}{}
		go func() {
			defer func() {
				<-c.concurrencyLimiter
			}()
			c.fragmentedWriteTask(ctx, resultChannel, fragmentCapture, bucket)
		}()
	}

	for range fragments {
		err = <-resultChannel
		if err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error during fragmented upload: %w", err)
	}
	return nil
}

// fragmentedWriteTask writes a single fragment to OCI Object Storage.
func (c *ociClient) fragmentedWriteTask(
	ctx context.Context,
	resultChannel chan error,
	fragment *s3.Fragment,
	bucket string) {

	putObjectRequest := objectstorage.PutObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(fragment.FragmentKey),
		PutObjectBody: io.NopCloser(bytes.NewReader(fragment.Data)),
		ContentLength: oraclecommon.Int64(int64(len(fragment.Data))),
	}

	_, err := c.objectStorageClient.PutObject(ctx, putObjectRequest)
	resultChannel <- err
}

func (c *ociClient) FragmentedDownloadObject(
	ctx context.Context,
	bucket string,
	key string,
	fileSize int,
	fragmentSize int) ([]byte, error) {

	if fileSize <= 0 {
		return nil, errors.New("fileSize must be greater than 0")
	}

	if fragmentSize <= 0 {
		return nil, errors.New("fragmentSize must be greater than 0")
	}

	fragmentKeys, err := s3.GetFragmentKeys(key, GetFragmentCount(fileSize, fragmentSize))
	if err != nil {
		return nil, fmt.Errorf("failed to get fragment keys: %w", err)
	}
	resultChannel := make(chan *readResult, len(fragmentKeys))

	for i, fragmentKey := range fragmentKeys {
		boundFragmentKey := fragmentKey
		boundI := i
		c.concurrencyLimiter <- struct{}{}
		go func() {
			defer func() {
				<-c.concurrencyLimiter
			}()
			c.readTask(ctx, resultChannel, bucket, boundFragmentKey, boundI)
		}()
	}

	fragments := make([]*s3.Fragment, len(fragmentKeys))
	for i := 0; i < len(fragmentKeys); i++ {
		result := <-resultChannel
		if result.err != nil {
			return nil, result.err
		}
		fragments[result.fragment.Index] = result.fragment
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error during fragmented download: %w", err)
	}

	return RecombineFragments(fragments)
}

// readResult is the result of a read task.
type readResult struct {
	fragment *s3.Fragment
	err      error
}

// readTask reads a single fragment from OCI Object Storage.
func (c *ociClient) readTask(
	ctx context.Context,
	resultChannel chan *readResult,
	bucket string,
	key string,
	index int) {

	result := &readResult{}
	defer func() {
		resultChannel <- result
	}()

	getObjectRequest := objectstorage.GetObjectRequest{
		NamespaceName: oraclecommon.String(c.cfg.Namespace),
		BucketName:    oraclecommon.String(bucket),
		ObjectName:    oraclecommon.String(key),
	}

	response, err := c.objectStorageClient.GetObject(ctx, getObjectRequest)
	if err != nil {
		result.err = err
		return
	}
	defer func() {
		if closeErr := response.Content.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	data, err := io.ReadAll(response.Content)
	if err != nil {
		result.err = err
		return
	}

	result.fragment = &s3.Fragment{
		FragmentKey: key,
		Data:        data,
		Index:       index,
	}
}

// Helper functions copied from s3 package (unexported)

// GetFragmentCount returns the number of fragments that a file of the given size will be broken into.
func GetFragmentCount(fileSize int, fragmentSize int) int {
	if fileSize < fragmentSize {
		return 1
	} else if fileSize%fragmentSize == 0 {
		return fileSize / fragmentSize
	} else {
		return fileSize/fragmentSize + 1
	}
}

// recombineFragments recombines fragments into a single file.
// Returns an error if any fragments are missing.
func RecombineFragments(fragments []*s3.Fragment) ([]byte, error) {
	if len(fragments) == 0 {
		return nil, fmt.Errorf("no fragments")
	}

	// Sort the fragments by index
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].Index < fragments[j].Index
	})

	// Make sure there aren't any gaps in the fragment indices
	dataSize := 0
	for i, fragment := range fragments {
		if fragment.Index != i {
			return nil, fmt.Errorf("missing fragment with index %d", i)
		}
		dataSize += len(fragment.Data)
	}

	// Make sure we have the last fragment
	if !strings.HasSuffix(fragments[len(fragments)-1].FragmentKey, "f") {
		return nil, fmt.Errorf("missing final fragment")
	}

	fragmentSize := len(fragments[0].Data)

	// Concatenate the data
	result := make([]byte, dataSize)
	for _, fragment := range fragments {
		copy(result[fragment.Index*fragmentSize:], fragment.Data)
	}

	return result, nil
}
