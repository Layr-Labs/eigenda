package workers

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
)

func TestBlobWriter(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)

	dataSize := rand.Uint64()%1024 + 64

	authenticated := rand.Intn(2) == 0
	var signerPrivateKey string
	if authenticated {
		signerPrivateKey = "asdf"
	}
	var functionName string
	if authenticated {
		functionName = "DisperseBlobAuthenticated"
	} else {
		functionName = "DisperseBlob"
	}

	randomizeBlobs := rand.Intn(2) == 0

	useCustomQuorum := rand.Intn(2) == 0
	var customQuorum []uint8
	if useCustomQuorum {
		customQuorum = []uint8{1, 2, 3}
	}

	config := &config.WorkerConfig{
		DataSize:         dataSize,
		SignerPrivateKey: signerPrivateKey,
		RandomizeBlobs:   randomizeBlobs,
		CustomQuorums:    customQuorum,
	}

	disperserClient := &MockDisperserClient{}
	unconfirmedKeyHandler := &MockKeyHandler{}
	unconfirmedKeyHandler.mock.On(
		"AddUnconfirmedKey", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	generatorMetrics := metrics.NewMockMetrics()

	writer := NewBlobWriter(
		&ctx,
		&waitGroup,
		logger,
		config,
		disperserClient,
		unconfirmedKeyHandler,
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

		// This is the key that will be assigned to the next blob.
		keyToReturn := make([]byte, 32)
		_, err = rand.Read(keyToReturn)
		assert.Nil(t, err)

		status := disperser.Processing
		disperserClient.mock = mock.Mock{} // reset mock state
		disperserClient.mock.On(functionName, mock.Anything, customQuorum).Return(&status, keyToReturn, errorToReturn)

		// Simulate the advancement of time (i.e. allow the writer to write the next blob).
		writer.writeNextBlob()

		disperserClient.mock.AssertNumberOfCalls(t, functionName, 1)
		unconfirmedKeyHandler.mock.AssertNumberOfCalls(t, "AddUnconfirmedKey", i+1-errorCount)

		if errorToReturn == nil {

			dataSentToDisperser := disperserClient.mock.Calls[0].Arguments.Get(0).([]byte)
			assert.NotNil(t, dataSentToDisperser)

			// Strip away the extra encoding bytes. We should have data of the expected size.
			decodedData := codec.RemoveEmptyByteFromPaddedBytes(dataSentToDisperser)
			assert.Equal(t, dataSize, uint64(len(decodedData)))

			// Verify that the proper data was sent to the unconfirmed key handler.
			checksum := md5.Sum(dataSentToDisperser)

			unconfirmedKeyHandler.mock.AssertCalled(t, "AddUnconfirmedKey", keyToReturn, checksum, uint(len(dataSentToDisperser)))

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
