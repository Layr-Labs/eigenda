package storeutil

import (
	"context"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/kvstore/mapstore"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRandomDataExpired(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	baseStore := mapstore.NewStore()
	store := ttlStore{
		store:  baseStore,
		ctx:    context.Background(),
		logger: logger,
	}

	data := make(map[string][]byte)
	expiryTimes := make(map[string]time.Time)

	startingTime := tu.RandomTime()
	simulatedSeconds := 1000
	endingTime := startingTime.Add(time.Duration(simulatedSeconds))

	// Generate some random data
	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(10)
		stringifiedKey := string(key)
		value := tu.RandomBytes(10)
		expiryTime := startingTime.Add(time.Duration(rand.Intn(simulatedSeconds)))

		data[stringifiedKey] = value
		expiryTimes[stringifiedKey] = expiryTime

		err := store.PutWithExpiration(key, value, expiryTime)
		assert.NoError(t, err)
	}

	currentTime := startingTime

	// Simulate time passing
	for currentTime.Before(endingTime) {
		currentTime = currentTime.Add(time.Duration(rand.Intn(simulatedSeconds / 10)))

		for key := range data {
			err := store.expireKeys(currentTime)
			assert.NoError(t, err)

			keyExpirationTime := expiryTimes[key]
			expired := currentTime.After(keyExpirationTime)

			if expired {
				value, err := store.Get([]byte(key))
				assert.Nil(t, value)
				assert.Error(t, err)
			} else {
				value, err := store.Get([]byte(key))
				assert.NoError(t, err)

				expectedValue := data[key]
				assert.Equal(t, expectedValue, value)
			}
		}
	}

}
