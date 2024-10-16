package tablestore

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestExpiryKeyParsing(t *testing.T) {
	tu.InitializeRandom()

	for i := 0; i < 1000; i++ {
		key := tu.RandomBytes(rand.Intn(100))
		expiryTime := tu.RandomTime()
		expiryKey := prependTimestamp(expiryTime, key)
		parsedExpiryTime, parsedKey := parsePrependedTimestamp(expiryKey)
		assert.Equal(t, key, parsedKey)
		assert.Equal(t, expiryTime, parsedExpiryTime)
	}

	// Try a very large key.
	key := tu.RandomBytes(100)
	expiryTime := time.Unix(0, 1<<62-1)
	expiryKey := prependTimestamp(expiryTime, key)
	parsedExpiryTime, parsedKey := parsePrependedTimestamp(expiryKey)
	assert.Equal(t, key, parsedKey)
	assert.Equal(t, expiryTime, parsedExpiryTime)
}

func TestExpiryKeyOrdering(t *testing.T) {
	tu.InitializeRandom()

	expiryKeys := make([][]byte, 0)

	for i := 0; i < 1000; i++ {
		expiryTime := tu.RandomTime()
		expiryKey := prependTimestamp(expiryTime, tu.RandomBytes(10))
		expiryKeys = append(expiryKeys, expiryKey)
	}

	// Add some keys with very large expiry times.
	for i := 0; i < 1000; i++ {
		expiryTime := tu.RandomTime().Add(time.Duration(1<<62 - 1))
		expiryKey := prependTimestamp(expiryTime, tu.RandomBytes(10))
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

		aTime, _ := parsePrependedTimestamp(a)
		bTime, _ := parsePrependedTimestamp(b)

		assert.True(t, aTime.Before(bTime))
	}
}

//func TestRandomDataExpired(t *testing.T) {
//	tu.InitializeRandom()
//
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//
//	baseStore := mapstore.NewStore()
//	//store := ttlStore{
//	//	store:  baseStore,
//	//	ctx:    context.Background(),
//	//	logger: logger,
//	//}
//	var store tableStore
//
//	data := make(map[string][]byte)
//	expiryTimes := make(map[string]time.Time)
//
//	startingTime := tu.RandomTime()
//	simulatedSeconds := 1000
//	endingTime := startingTime.Add(time.Duration(simulatedSeconds) * time.Second)
//
//	// Generate some random data
//	for i := 0; i < 1000; i++ {
//		key := tu.RandomBytes(10)
//		stringifiedKey := string(key)
//		value := tu.RandomBytes(10)
//		expiryTime := startingTime.Add(time.Duration(rand.Intn(simulatedSeconds)) * time.Second)
//
//		data[stringifiedKey] = value
//		expiryTimes[stringifiedKey] = expiryTime
//
//		err := store.PutWithExpiration(key, value, expiryTime)
//		assert.NoError(t, err)
//	}
//
//	currentTime := startingTime
//
//	// Simulate time passing
//	for currentTime.Before(endingTime) {
//
//		elapsedSeconds := rand.Intn(simulatedSeconds / 10)
//		currentTime = currentTime.Add(time.Duration(elapsedSeconds) * time.Second)
//
//		err = store.expireKeys(currentTime)
//		assert.NoError(t, err)
//
//		for key := range data {
//			keyExpirationTime := expiryTimes[key]
//			expired := !currentTime.Before(keyExpirationTime)
//
//			if expired {
//				value, err := store.Get([]byte(key))
//				assert.Error(t, err)
//				assert.Nil(t, value)
//			} else {
//				value, err := store.Get([]byte(key))
//				assert.NoError(t, err)
//				expectedValue := data[key]
//				assert.Equal(t, expectedValue, value)
//			}
//		}
//	}
//}

//func TestBatchRandomDataExpired(t *testing.T) {
//	tu.InitializeRandom()
//
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//
//	baseStore := mapstore.NewStore()
//	store := ttlStore{
//		store:  baseStore,
//		ctx:    context.Background(),
//		logger: logger,
//	}
//
//	data := make(map[string][]byte)
//	expiryTimes := make(map[string]time.Time)
//
//	startingTime := tu.RandomTime()
//	simulatedSeconds := 1000
//	endingTime := startingTime.Add(time.Duration(simulatedSeconds) * time.Second)
//
//	// Generate some random data
//	for i := 0; i < 100; i++ {
//
//		expiryTime := startingTime.Add(time.Duration(rand.Intn(simulatedSeconds)) * time.Second)
//
//		keys := make([][]byte, 0)
//		values := make([][]byte, 0)
//
//		// Generate a batch of random data
//		for j := 0; j < 10; j++ {
//			key := tu.RandomBytes(10)
//			keys = append(keys, key)
//			stringifiedKey := string(key)
//			value := tu.RandomBytes(10)
//			values = append(values, value)
//
//			data[stringifiedKey] = value
//			expiryTimes[stringifiedKey] = expiryTime
//		}
//
//		err := store.PutBatchWithExpiration(keys, values, expiryTime)
//		assert.NoError(t, err)
//	}
//
//	currentTime := startingTime
//
//	// Simulate time passing
//	for currentTime.Before(endingTime) {
//
//		elapsedSeconds := rand.Intn(simulatedSeconds / 10)
//		currentTime = currentTime.Add(time.Duration(elapsedSeconds) * time.Second)
//
//		err = store.expireKeys(currentTime)
//		assert.NoError(t, err)
//
//		for key := range data {
//			keyExpirationTime := expiryTimes[key]
//			expired := !currentTime.Before(keyExpirationTime)
//
//			if expired {
//				value, err := store.Get([]byte(key))
//				assert.Error(t, err)
//				assert.Nil(t, value)
//			} else {
//				value, err := store.Get([]byte(key))
//				assert.NoError(t, err)
//				expectedValue := data[key]
//				assert.Equal(t, expectedValue, value)
//			}
//		}
//	}
//}
//
//func TestBigBatchOfDeletions(t *testing.T) {
//	tu.InitializeRandom()
//
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//
//	baseStore := mapstore.NewStore()
//	store := ttlStore{
//		store:  baseStore,
//		ctx:    context.Background(),
//		logger: logger,
//	}
//
//	data := make(map[string][]byte)
//	expiryTimes := make(map[string]time.Time)
//
//	startingTime := tu.RandomTime()
//	simulatedSeconds := 1000
//
//	// Generate some random data
//	for i := 0; i < 2345; i++ {
//		key := tu.RandomBytes(10)
//		stringifiedKey := string(key)
//		value := tu.RandomBytes(10)
//		expiryTime := startingTime.Add(time.Duration(rand.Intn(simulatedSeconds)) * time.Second)
//
//		data[stringifiedKey] = value
//		expiryTimes[stringifiedKey] = expiryTime
//
//		err = store.PutWithExpiration(key, value, expiryTime)
//		assert.NoError(t, err)
//	}
//
//	// Move time forward by one large step
//	elapsedSeconds := simulatedSeconds * 2
//	currentTime := startingTime.Add(time.Duration(elapsedSeconds) * time.Second)
//
//	err = store.expireKeys(currentTime)
//	assert.NoError(t, err)
//
//	// All keys should be expired
//	for key := range data {
//		value, err := store.Get([]byte(key))
//		assert.Error(t, err)
//		assert.Nil(t, value)
//	}
//}
