package test

import (
	"context"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/mock"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
	bucket         = "eigen-test"
)

type clientBuilder struct {
	// This method is called at the beginning of the test.
	start func() error
	// This method is called to build a new client.
	build func() (s3.Client, error)
	// This method is called at the end of the test when all operations are done.
	finish func() error
}

var clientBuilders = []*clientBuilder{
	{
		start: func() error {
			return nil
		},
		build: func() (s3.Client, error) {
			return mock.NewS3Client(), nil
		},
		finish: func() error {
			return nil
		},
	},
	{
		start: func() error {
			return setupLocalstack()
		},
		build: func() (s3.Client, error) {

			logger, err := common.NewLogger(common.DefaultLoggerConfig())
			if err != nil {
				return nil, err
			}

			config := aws.DefaultClientConfig()
			config.EndpointURL = localstackHost
			config.Region = "us-east-1"

			err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
			if err != nil {
				return nil, err
			}
			err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
			if err != nil {
				return nil, err
			}

			client, err := s3.NewClient(context.Background(), *config, logger)
			if err != nil {
				return nil, err
			}

			err = client.CreateBucket(context.Background(), bucket)
			if err != nil {
				return nil, err
			}

			return client, nil
		},
		finish: func() error {
			teardownLocalstack()
			return nil
		},
	},
}

func setupLocalstack() error {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
		if err != nil && err.Error() == "container already exists" {
			teardownLocalstack()
			return err
		}
	}
	return nil
}

func teardownLocalstack() {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func RandomOperationsTest(t *testing.T, client s3.Client) {
	numberToWrite := 100
	expectedData := make(map[string][]byte)

	fragmentSize := rand.Intn(1000) + 1000
	for i := 0; i < numberToWrite; i++ {
		key := tu.RandomString(10)
		fragmentMultiple := rand.Float64() * 10
		dataSize := int(fragmentMultiple*float64(fragmentSize)) + 1
		data := tu.RandomBytes(dataSize)
		expectedData[key] = data
		err := client.FragmentedUploadObject(context.Background(), bucket, key, data, fragmentSize)
		assert.NoError(t, err)
	}

	// Read back the data
	for key, expected := range expectedData {
		data, err := client.FragmentedDownloadObject(context.Background(), bucket, key, len(expected), fragmentSize)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)

		// List the objects
		objects, err := client.ListObjects(context.Background(), bucket, key)
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
	objects, err := client.ListObjects(context.Background(), bucket, "nonexistent")
	assert.NoError(t, err)
	assert.Len(t, objects, 0)
}

func TestRandomOperations(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		assert.NoError(t, err)

		client, err := builder.build()
		assert.NoError(t, err)
		RandomOperationsTest(t, client)

		err = builder.finish()
		assert.NoError(t, err)
	}
}

func ReadNonExistentValueTest(t *testing.T, client s3.Client) {
	_, err := client.FragmentedDownloadObject(context.Background(), bucket, "nonexistent", 1000, 1000)
	assert.Error(t, err)
	randomKey := tu.RandomString(10)
	_, err = client.FragmentedDownloadObject(context.Background(), bucket, randomKey, 0, 0)
	assert.Error(t, err)
}

func TestReadNonExistentValue(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		assert.NoError(t, err)

		client, err := builder.build()
		assert.NoError(t, err)
		ReadNonExistentValueTest(t, client)

		err = builder.finish()
		assert.NoError(t, err)
	}
}

func TestHeadObject(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		assert.NoError(t, err)

		client, err := builder.build()
		assert.NoError(t, err)

		key := tu.RandomString(10)
		err = client.UploadObject(context.Background(), bucket, key, []byte("test"))
		assert.NoError(t, err)
		size, err := client.HeadObject(context.Background(), bucket, key)
		assert.NoError(t, err)
		assert.NotNil(t, size)
		assert.Equal(t, int64(4), *size)

		size, err = client.HeadObject(context.Background(), bucket, "nonexistent")
		assert.ErrorIs(t, err, s3.ErrObjectNotFound)
		assert.Nil(t, size)

		err = builder.finish()
		assert.NoError(t, err)
	}
}
