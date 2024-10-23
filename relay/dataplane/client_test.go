package dataplane

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var clientBuilders = []func() (S3Client, error){
	func() (S3Client, error) {
		return NewLocalClient(mapstore.NewStore()), nil
	},
	func() (S3Client, error) {

		config := DefaultS3Config()
		config.Bucket = "eigen-cody-test"
		config.AutoCreateBucket = true

		client, err := NewS3Client(context.Background(), config)
		if err != nil {
			return nil, err
		}

		return client, nil
	},
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
	for _, clientBuilder := range clientBuilders {
		client, err := clientBuilder()
		assert.NoError(t, err)
		RandomOperationsTest(t, client)
		err = client.Close()
		assert.NoError(t, err)
	}
}

func ReadNonExistentValueTest(t *testing.T, client S3Client) {
	_, err := client.Download("nonexistent", 0, 0)
	assert.Error(t, err)
	randomKey := tu.RandomString(10)
	_, err = client.Download(randomKey, 0, 0)
	assert.Error(t, err)
}

func TestReadNonExistentValue(t *testing.T) {
	tu.InitializeRandom()
	for _, clientBuilder := range clientBuilders {
		client, err := clientBuilder()
		assert.NoError(t, err)
		ReadNonExistentValueTest(t, client)
		err = client.Close()
		assert.NoError(t, err)
	}
}

// TODO:
//  - test bucket creation
//  - test a store that already has a bucket created
