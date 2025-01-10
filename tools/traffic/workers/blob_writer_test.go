package workers

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/rand"
)

func TestBlobWriter(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)

	dataSize := rand.Uint64()%1024 + 64

	randomizeBlobs := rand.Intn(2) == 0
	useCustomQuorum := rand.Intn(2) == 0
	var customQuorum []uint8
	if useCustomQuorum {
		customQuorum = []uint8{1, 2, 3}
	}

	config := &config.BlobWriterConfig{
		DataSize:       dataSize,
		RandomizeBlobs: randomizeBlobs,
		CustomQuorums:  customQuorum,
	}

	disperserClient := &MockDisperserClient{}
	generatorMetrics := metrics.NewMockMetrics()

	writer := NewBlobWriter(
		"test-writer",
		&ctx,
		config,
		&waitGroup,
		logger,
		disperserClient,
		generatorMetrics)

	errorCount := 0

	var previousData []byte

	for i := 0; i < 100; i++ {
		var errorToReturn error
		if i%10 == 0 {
			errorToReturn = fmt.Errorf("intentional error for testing purposes")
			errorCount++
		} else {
			errorToReturn = nil
		}

		// This is the Key that will be assigned to the next blob.
		var keyToReturn corev2.BlobKey
		_, err = rand.Read(keyToReturn[:])
		assert.Nil(t, err)

		status := dispv2.Queued
		disperserClient.mock = mock.Mock{} // reset mock state
		disperserClient.mock.On("DisperseBlob",
			mock.Anything,
			mock.AnythingOfType("[]uint8"),
			mock.AnythingOfType("uint16"),
			mock.AnythingOfType("[]uint8"),
			mock.AnythingOfType("uint32"),
		).Return(&status, keyToReturn, errorToReturn)

		// Simulate the advancement of time (i.e. allow the writer to write the next blob).
		writer.writeNextBlob()

		disperserClient.mock.AssertNumberOfCalls(t, "DisperseBlob", 1)

		if errorToReturn == nil {
			dataSentToDisperser := disperserClient.mock.Calls[0].Arguments.Get(1).([]byte)
			assert.NotNil(t, dataSentToDisperser)

			// Strip away the extra encoding bytes. We should have data of the expected Size.
			decodedData := codec.RemoveEmptyByteFromPaddedBytes(dataSentToDisperser)
			assert.Equal(t, dataSize, uint64(len(decodedData)))

			// Verify that data has the proper amount of randomness.
			if previousData != nil {
				if randomizeBlobs {
					// We expect each blob to be different.
					assert.NotEqual(t, previousData, dataSentToDisperser)
				} else {
					// We expect each blob to be the same.
					assert.Equal(t, previousData, dataSentToDisperser)
				}
			}
			previousData = dataSentToDisperser
		}

		// Verify metrics.
		assert.Equal(t, float64(i+1-errorCount), generatorMetrics.GetCount("write_success"))
		assert.Equal(t, float64(errorCount), generatorMetrics.GetCount("write_failure"))
	}

	cancel()
}
