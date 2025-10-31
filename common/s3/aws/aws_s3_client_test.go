package aws

import (
	"context"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	s3common "github.com/Layr-Labs/eigenda/common/s3"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/stretchr/testify/require"
)

var (
	logger = test.GetLogger()
)

const (
	bucket         = "eigen-test"
	localstackPort = "4578"
	localstackHost = "http://0.0.0.0:4578"
)

func setupLocalStackTest(t *testing.T) s3common.S3Client {
	t.Helper()

	ctx := t.Context()

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start LocalStack container")

	t.Cleanup(func() {
		logger.Info("Stopping LocalStack container")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = localstackContainer.Terminate(ctx)
	})

	config := aws.DefaultClientConfig()
	config.EndpointURL = localstackHost
	config.Region = "us-east-1"

	err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
	require.NoError(t, err, "failed to set AWS_ACCESS_KEY_ID")
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
	require.NoError(t, err, "failed to set AWS_SECRET_ACCESS_KEY")

	client, err := NewAwsS3Client(
		ctx,
		logger,
		config.EndpointURL,
		config.Region,
		config.FragmentParallelismFactor,
		config.FragmentParallelismConstant,
		config.AccessKey,
		config.SecretAccessKey,
	)
	require.NoError(t, err, "failed to create S3 client")

	err = client.CreateBucket(ctx, bucket)
	require.NoError(t, err, "failed to create S3 bucket")

	return client
}

func runRandomOperationsTest(t *testing.T, client s3common.S3Client) {
	t.Helper()
	ctx := t.Context()
	numberToWrite := 100
	expectedData := make(map[string][]byte)

	fragmentSize := rand.Intn(1000) + 1000
	for i := 0; i < numberToWrite; i++ {
		key := random.RandomString(10)
		fragmentMultiple := rand.Float64() * 10
		dataSize := int(fragmentMultiple*float64(fragmentSize)) + 1
		data := random.RandomBytes(dataSize)
		expectedData[key] = data
		err := client.FragmentedUploadObject(ctx, bucket, key, data, fragmentSize)
		require.NoError(t, err, "failed to upload fragmented object for key %s", key)
	}

	// Read back the data
	for key, expected := range expectedData {
		data, err := client.FragmentedDownloadObject(ctx, bucket, key, len(expected), fragmentSize)
		require.NoError(t, err, "failed to download fragmented object for key %s", key)
		require.Equal(t, expected, data, "downloaded data should match uploaded data for key %s", key)

		// List the objects
		objects, err := client.ListObjects(ctx, bucket, key)
		require.NoError(t, err, "failed to list objects for key %s", key)
		numFragments := math.Ceil(float64(len(expected)) / float64(fragmentSize))
		require.Len(t, objects, int(numFragments), "should have correct number of fragments for key %s", key)
		totalSize := int64(0)
		for _, object := range objects {
			totalSize += object.Size
		}
		require.Equal(t, int64(len(expected)), totalSize,
			"total fragment size should match original data size for key %s", key)
	}

	// Attempt to list non-existent objects
	objects, err := client.ListObjects(ctx, bucket, "nonexistent")
	require.NoError(t, err, "failed to list non-existent objects")
	require.Len(t, objects, 0, "should return empty list for non-existent objects")
}

func TestRandomOperations(t *testing.T) {
	random.InitializeRandom()

	t.Run("mock_client", func(t *testing.T) {
		client := s3common.NewMockS3Client()
		runRandomOperationsTest(t, client)
	})

	t.Run("localstack_client", func(t *testing.T) {
		client := setupLocalStackTest(t)
		runRandomOperationsTest(t, client)
	})
}

func TestReadNonExistentValue(t *testing.T) {
	random.InitializeRandom()

	t.Run("mock_client", func(t *testing.T) {
		client := s3common.NewMockS3Client()
		runReadNonExistentValueTest(t, client)
	})

	t.Run("localstack_client", func(t *testing.T) {
		client := setupLocalStackTest(t)
		runReadNonExistentValueTest(t, client)
	})
}

func runReadNonExistentValueTest(t *testing.T, client s3common.S3Client) {
	t.Helper()
	ctx := t.Context()

	_, err := client.FragmentedDownloadObject(ctx, bucket, "nonexistent", 1000, 1000)
	require.Error(t, err, "should fail to download non-existent object")

	randomKey := random.RandomString(10)
	_, err = client.FragmentedDownloadObject(ctx, bucket, randomKey, 0, 0)
	require.Error(t, err, "should fail to download random non-existent object")
}

func TestHeadObject(t *testing.T) {
	random.InitializeRandom()

	t.Run("mock_client", func(t *testing.T) {
		client := s3common.NewMockS3Client()
		runHeadObjectTest(t, client)
	})

	t.Run("localstack_client", func(t *testing.T) {
		client := setupLocalStackTest(t)
		runHeadObjectTest(t, client)
	})
}

func runHeadObjectTest(t *testing.T, client s3common.S3Client) {
	t.Helper()
	ctx := t.Context()

	key := random.RandomString(10)
	err := client.UploadObject(ctx, bucket, key, []byte("test"))
	require.NoError(t, err, "failed to upload test object")

	size, err := client.HeadObject(ctx, bucket, key)
	require.NoError(t, err, "failed to get head object for existing key")
	require.NotNil(t, size, "size should not be nil for existing object")
	require.Equal(t, int64(4), *size, "size should match uploaded data")

	size, err = client.HeadObject(ctx, bucket, "nonexistent")
	require.Error(t, err, "should fail to get head object for non-existent key")
	require.Nil(t, size, "size should be nil for non-existent object")
}
