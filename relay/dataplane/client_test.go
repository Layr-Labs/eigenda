package dataplane

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
)

const (
	localstackPort = "4566"
	localstackHost = "http://0.0.0.0:4566"
)

type clientBuilder struct {
	// This method is called at the beginning of the test.
	start func() error
	// This method is called to build a new client.
	build func() (S3Client, error)
	// This method is called at the end of the test when all operations are done.
	finish func() error
}

var clientBuilders = []*clientBuilder{
	{
		start: func() error {
			return nil
		},
		build: func() (S3Client, error) {
			return NewLocalClient(mapstore.NewStore()), nil
		},
		finish: func() error {
			return nil
		},
	},
	{
		start: func() error {
			return setupLocalstack()
		},
		build: func() (S3Client, error) {

			config := DefaultS3Config()
			config.AWSConfig.Endpoint = aws.String(localstackHost)
			config.AWSConfig.S3ForcePathStyle = aws.Bool(true)
			config.AWSConfig.WithRegion("us-east-1")

			err := os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
			if err != nil {
				return nil, err
			}
			err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
			if err != nil {
				return nil, err
			}

			config.Bucket = "this-is-a-test-bucket"
			config.AutoCreateBucket = true

			client, err := NewS3Client(context.Background(), config)
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
	var err error
	dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
	if err != nil {
		teardownLocalstack()
		return err
	}
	return nil
}

func teardownLocalstack() {
	deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
}

func RandomOperationsTest(t *testing.T, client S3Client) {
	numberToWrite := 100
	expectedData := make(map[string][]byte)

	fragmentSize := rand.Intn(1000) + 1000

	for i := 0; i < numberToWrite; i++ {
		key := tu.RandomString(10)
		fragmentMultiple := rand.Float64() * 10
		dataSize := int(fragmentMultiple*float64(fragmentSize)) + 1
		data := tu.RandomBytes(dataSize)
		expectedData[key] = data

		err := client.Upload(key, data, fragmentSize, time.Hour)
		assert.NoError(t, err)
	}

	// Read back the data
	for key, expected := range expectedData {
		data, err := client.Download(key, len(expected), fragmentSize)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	}
}

func TestRandomOperations(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		assert.NoError(t, err)

		client, err := builder.build()
		assert.NoError(t, err)
		RandomOperationsTest(t, client)
		err = client.Close()
		assert.NoError(t, err)

		err = builder.finish()
		assert.NoError(t, err)
	}
}

func ReadNonExistentValueTest(t *testing.T, client S3Client) {
	_, err := client.Download("nonexistent", 1000, 1000)
	assert.Error(t, err)
	randomKey := tu.RandomString(10)
	_, err = client.Download(randomKey, 0, 0)
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
		err = client.Close()
		assert.NoError(t, err)

		err = builder.finish()
		assert.NoError(t, err)
	}
}
