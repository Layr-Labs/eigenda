package storeutil

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/kvstore/mapstore"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestExpiryKeyParsing(t *testing.T) {
	tu.InitializeRandom()

	for i := 0; i < 1000; i++ {
		expiryTime := tu.RandomTime()
		expiryKey := buildExpiryKey(expiryTime)
		parsedExpiryTime, err := parseExpiryKey(expiryKey)
		assert.NoError(t, err)
		assert.Equal(t, expiryTime, parsedExpiryTime)
	}

	// Try a very large key.
	expiryTime := time.Unix(0, 1<<62-1)
	expiryKey := buildExpiryKey(expiryTime)
	parsedExpiryTime, err := parseExpiryKey(expiryKey)
	assert.NoError(t, err)
	assert.Equal(t, expiryTime, parsedExpiryTime)
}

func TestExpiryKeyOrdering(t *testing.T) {
	tu.InitializeRandom()

	expiryKeys := make([][]byte, 0)

	for i := 0; i < 1000; i++ {
		expiryTime := tu.RandomTime()
		expiryKey := buildExpiryKey(expiryTime)
		expiryKeys = append(expiryKeys, expiryKey)
	}

	// Add some keys with very large expiry times.
	for i := 0; i < 1000; i++ {
		expiryTime := tu.RandomTime().Add(time.Duration(1<<62 - 1))
		expiryKey := buildExpiryKey(expiryTime)
		expiryKeys = append(expiryKeys, expiryKey)
	}

	// Sort the keys.
	sort.Slice(expiryKeys, func(i, j int) bool {
		return string(expiryKeys[i]) < string(expiryKeys[j])
	})

	// Check that the keys are sorted.
	for i := 1; i < len(expiryKeys)-1; i++ {
		a := expiryKeys[i-1]
		b := expiryKeys[i]

		aTime, err := parseExpiryKey(a)
		assert.NoError(t, err)
		bTime, err := parseExpiryKey(b)
		assert.NoError(t, err)

		assert.True(t, aTime.Before(bTime))
	}
}

func TestRandomDataExpired(t *testing.T) {
	tu.InitializeRandom(0) // TODO

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
	endingTime := startingTime.Add(time.Duration(simulatedSeconds) * time.Second)

	// Generate some random data
	for i := 0; i < 1000; i++ { // TODO 1000
		key := tu.RandomBytes(10)
		stringifiedKey := string(key)
		value := tu.RandomBytes(10)
		expiryTime := startingTime.Add(time.Duration(rand.Intn(simulatedSeconds)) * time.Second)

		data[stringifiedKey] = value
		expiryTimes[stringifiedKey] = expiryTime

		err := store.PutWithExpiration(key, value, expiryTime)
		assert.NoError(t, err)
	}

	currentTime := startingTime

	// Simulate time passing
	for currentTime.Before(endingTime) {

		elapsedSeconds := rand.Intn(simulatedSeconds / 10)
		currentTime = currentTime.Add(time.Duration(elapsedSeconds) * time.Second)

		fmt.Printf(">>>>>>>>>>>>>> expiring keys until %v\n", currentTime)

		err := store.expireKeys(currentTime)
		assert.NoError(t, err)

		for key := range data {
			keyExpirationTime := expiryTimes[key]
			expired := currentTime.After(keyExpirationTime)

			if expired {

				fmt.Printf(">>>>>>>>>>>>>> key %s expired at %v\n", key, keyExpirationTime)

				value, err := store.Get([]byte(key))
				assert.Error(t, err)
				assert.Nil(t, value)
			} else {
				value, err := store.Get([]byte(key))
				assert.NoError(t, err)

				expectedValue := data[key]
				assert.Equal(t, expectedValue, value)
			}
		}
	}

}
