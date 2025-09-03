package test

import (
	"context"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/mock"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/stretchr/testify/assert"
)

const (
	bucket = "eigen-test"
)

var (
	// Shared LocalStack container and client for all tests
	sharedLocalStackContainer *testbed.LocalStackContainer
	sharedLocalStackClient    s3.Client
)

// TestMain sets up and tears down shared resources for all tests
func TestMain(m *testing.M) {
	deployLocalStack := os.Getenv("DEPLOY_LOCALSTACK") != "false"

	if deployLocalStack {
		// Setup shared LocalStack container
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		lsConfig := testbed.LocalStackConfig{
			Enabled:  true,
			Services: []string{"s3"},
			Region:   "us-east-1",
			Debug:    false,
		}

		container, err := testbed.NewLocalStackContainer(ctx, lsConfig)
		if err != nil {
			log.Fatalf("Failed to start shared LocalStack: %v", err)
		}
		sharedLocalStackContainer = container

		// Setup shared S3 client
		logger, err := common.NewLogger(common.DefaultLoggerConfig())
		if err != nil {
			container.Terminate(context.Background())
			log.Fatalf("Failed to create logger: %v", err)
		}

		config := *aws.DefaultClientConfig()
		config.EndpointURL = container.Endpoint()
		config.Region = "us-east-1"

		os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")

		client, err := s3.NewClient(context.Background(), config, logger)
		if err != nil {
			container.Terminate(context.Background())
			log.Fatalf("Failed to create S3 client: %v", err)
		}

		err = client.CreateBucket(context.Background(), bucket)
		if err != nil {
			container.Terminate(context.Background())
			log.Fatalf("Failed to create bucket: %v", err)
		}

		sharedLocalStackClient = client
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup
	if sharedLocalStackContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		sharedLocalStackContainer.Terminate(ctx)
	}

	os.Exit(exitCode)
}

// testClientSetup represents a test client with optional cleanup
type testClientSetup struct {
	name   string
	client s3.Client
}

// setupMockClient creates a mock S3 client for testing
func setupMockClient() *testClientSetup {
	return &testClientSetup{
		name:   "mock",
		client: mock.NewS3Client(),
	}
}

// setupLocalStackClient returns the shared LocalStack S3 client for testing
func setupLocalStackClient(_ *testing.T) *testClientSetup {
	if sharedLocalStackClient == nil {
		return nil // LocalStack is disabled or not initialized
	}

	return &testClientSetup{
		name:   "localstack",
		client: sharedLocalStackClient,
	}
}

// getAllTestClients returns all available test clients
func getAllTestClients(t *testing.T) []*testClientSetup {
	clients := []*testClientSetup{
		setupMockClient(),
	}

	if localStackClient := setupLocalStackClient(t); localStackClient != nil {
		clients = append(clients, localStackClient)
	}

	return clients
}

func TestRandomOperations(t *testing.T) {
	tu.InitializeRandom()

	for _, setup := range getAllTestClients(t) {
		t.Run(setup.name, func(t *testing.T) {
			numberToWrite := 100
			expectedData := make(map[string][]byte)

			fragmentSize := rand.Intn(1000) + 1000
			for i := 0; i < numberToWrite; i++ {
				key := tu.RandomString(10)
				fragmentMultiple := rand.Float64() * 10
				dataSize := int(fragmentMultiple*float64(fragmentSize)) + 1
				data := tu.RandomBytes(dataSize)
				expectedData[key] = data
				err := setup.client.FragmentedUploadObject(context.Background(), bucket, key, data, fragmentSize)
				assert.NoError(t, err)
			}

			// Read back the data
			for key, expected := range expectedData {
				data, err := setup.client.FragmentedDownloadObject(context.Background(), bucket, key, len(expected), fragmentSize)
				assert.NoError(t, err)
				assert.Equal(t, expected, data)

				// List the objects
				objects, err := setup.client.ListObjects(context.Background(), bucket, key)
				assert.NoError(t, err)
				numFragments := math.Ceil(float64(len(expected)) / float64(fragmentSize))
				assert.Len(t, objects, int(numFragments))
				totalSize := int64(0)
				for _, object := range objects {
					totalSize += object.Size
				}
				assert.Equal(t, int64(len(expected)), totalSize)
			}

			// Attempt to list non-existent objects
			objects, err := setup.client.ListObjects(context.Background(), bucket, "nonexistent")
			assert.NoError(t, err)
			assert.Len(t, objects, 0)
		})
	}
}

func TestReadNonExistentValue(t *testing.T) {
	tu.InitializeRandom()

	for _, setup := range getAllTestClients(t) {
		t.Run(setup.name, func(t *testing.T) {
			_, err := setup.client.FragmentedDownloadObject(context.Background(), bucket, "nonexistent", 1000, 1000)
			assert.Error(t, err)
			randomKey := tu.RandomString(10)
			_, err = setup.client.FragmentedDownloadObject(context.Background(), bucket, randomKey, 0, 0)
			assert.Error(t, err)
		})
	}
}

func TestHeadObject(t *testing.T) {
	tu.InitializeRandom()

	for _, setup := range getAllTestClients(t) {
		t.Run(setup.name, func(t *testing.T) {
			key := tu.RandomString(10)
			err := setup.client.UploadObject(context.Background(), bucket, key, []byte("test"))
			assert.NoError(t, err)
			size, err := setup.client.HeadObject(context.Background(), bucket, key)
			assert.NoError(t, err)
			assert.NotNil(t, size)
			assert.Equal(t, int64(4), *size)

			size, err = setup.client.HeadObject(context.Background(), bucket, "nonexistent")
			assert.ErrorIs(t, err, s3.ErrObjectNotFound)
			assert.Nil(t, size)
		})
	}
}
