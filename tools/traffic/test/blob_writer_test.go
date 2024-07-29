package test

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

func TestBlobWriter(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := newMockTicker(startTime)

	dataSize := rand.Uint64()%1024 + 64

	authenticated := rand.Intn(2) == 0
	var signerPrivateKey string
	if authenticated {
		signerPrivateKey = "asdf"
	}

	randomizeBlobs := rand.Intn(2) == 0

	useCustomQuorum := rand.Intn(2) == 0
	var customQuorum []uint8
	if useCustomQuorum {
		customQuorum = []uint8{1, 2, 3}
	}

	config := &workers.Config{
		DataSize:         dataSize,
		SignerPrivateKey: signerPrivateKey,
		RandomizeBlobs:   randomizeBlobs,
		CustomQuorums:    customQuorum,
	}

	lock := sync.Mutex{}

	disperserClient := newMockDisperserClient(t, &lock, authenticated)
	unconfirmedKeyHandler := newMockKeyHandler(t, &lock)

	generatorMetrics := metrics.NewMockMetrics()

	writer := workers.NewBlobWriter(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		disperserClient,
		unconfirmedKeyHandler,
		generatorMetrics)
	writer.Start()

	errorProbability := 0.1
	errorCount := 0

	var previousData []byte

	for i := 0; i < 100; i++ {
		if rand.Float64() < errorProbability {
			disperserClient.ErrorToReturn = fmt.Errorf("intentional error for testing purposes")
			errorCount++
		} else {
			disperserClient.ErrorToReturn = nil
		}

		// This is the key that will be assigned to the next blob.
		disperserClient.KeyToReturn = make([]byte, 32)
		_, err = rand.Read(disperserClient.KeyToReturn)
		assert.Nil(t, err)

		// Move time forward, allowing the writer to attempt to send a blob.
		ticker.Tick(1 * time.Second)

		// Wait until the writer finishes its work.
		tu.AssertEventuallyTrue(t, func() bool {
			lock.Lock()
			defer lock.Unlock()
			return int(disperserClient.Count) > i && int(unconfirmedKeyHandler.Count)+errorCount > i
		}, time.Second)

		// These methods should be called exactly once per tick if there are no errors.
		// In the presence of errors, nothing should be passed to the unconfirmed key handler.
		assert.Equal(t, uint(i+1), disperserClient.Count)
		assert.Equal(t, uint(i+1-errorCount), unconfirmedKeyHandler.Count)

		if disperserClient.ErrorToReturn == nil {
			assert.NotNil(t, disperserClient.ProvidedData)
			assert.Equal(t, customQuorum, disperserClient.ProvidedQuorum)

			// Strip away the extra encoding bytes. We should have data of the expected size.
			decodedData := codec.RemoveEmptyByteFromPaddedBytes(disperserClient.ProvidedData)
			assert.Equal(t, dataSize, uint64(len(decodedData)))

			// Verify that the proper data was sent to the unconfirmed key handler.
			assert.Equal(t, uint(len(disperserClient.ProvidedData)), unconfirmedKeyHandler.ProvidedSize)
			checksum := md5.Sum(disperserClient.ProvidedData)
			assert.Equal(t, checksum, unconfirmedKeyHandler.ProvidedChecksum)
			assert.Equal(t, disperserClient.KeyToReturn, unconfirmedKeyHandler.ProvidedKey)

			// Verify that data has the proper amount of randomness.
			if previousData != nil {
				if randomizeBlobs {
					// We expect each blob to be different.
					assert.NotEqual(t, previousData, disperserClient.ProvidedData)
				} else {
					// We expect each blob to be the same.
					assert.Equal(t, previousData, disperserClient.ProvidedData)
				}
			}
			previousData = disperserClient.ProvidedData
		}

		// Verify metrics.
		assert.Equal(t, float64(i+1-errorCount), generatorMetrics.GetCount("write_success"))
		assert.Equal(t, float64(errorCount), generatorMetrics.GetCount("write_failure"))
	}

	cancel()
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
