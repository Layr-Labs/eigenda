package workers

import (
	"context"
	"crypto/md5"
	"fmt"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

// TODO create test util package maybe

// initializeRandom initializes the random number generator. Prints the seed so that the test can be rerun
// deterministically. Replace a call to this method with a call to initializeRandomWithSeed to rerun a test
// with a specific seed.
func initializeRandom() {
	rand.Seed(uint64(time.Now().UnixNano()))
	seed := rand.Uint64()
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// initializeRandomWithSeed initializes the random number generator with a specific seed.
func initializeRandomWithSeed(seed uint64) {
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// assertEventuallyTrue asserts that a condition is true within a given duration. Repeatably checks the condition.
func assertEventuallyTrue(t *testing.T, condition func() bool, duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		if condition() {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	assert.True(t, condition(), "Condition did not become true within the given duration")
}

// executeWithTimeout executes a function with a timeout.
// Panics if the function does not complete within the given duration.
func executeWithTimeout(f func(), duration time.Duration) {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(duration):
		panic("function did not complete within the given duration")
	}
}

// MockUnconfirmedKeyHandler is a stand-in for the blob verifier's UnconfirmedKeyHandler.
type MockUnconfirmedKeyHandler struct {
	t *testing.T

	// TODO rename
	ProvidedKey      []byte
	ProvidedChecksum [16]byte
	ProvidedSize     uint

	// Incremented each time AddUnconfirmedKey is called.
	Count uint

	lock *sync.Mutex
}

func NewMockUnconfirmedKeyHandler(t *testing.T, lock *sync.Mutex) *MockUnconfirmedKeyHandler {
	return &MockUnconfirmedKeyHandler{
		t:    t,
		lock: lock,
	}
}

func (m *MockUnconfirmedKeyHandler) AddUnconfirmedKey(key *[]byte, checksum *[16]byte, size uint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.ProvidedKey = *key
	m.ProvidedChecksum = *checksum
	m.ProvidedSize = size

	m.Count++
}

type MockDisperserClient struct {
	t *testing.T
	// if true, DisperseBlobAuthenticated is expected to be used, otherwise DisperseBlob is expected to be used
	authenticated bool

	// The next status, key, and error to return from DisperseBlob or DisperseBlobAuthenticated
	StatusToReturn disperser.BlobStatus
	KeyToReturn    []byte
	ErrorToReturn  error

	// The previous values passed to DisperseBlob or DisperseBlobAuthenticated
	ProvidedData   []byte
	ProvidedQuorum []uint8

	// Incremented each time DisperseBlob or DisperseBlobAuthenticated is called.
	Count uint

	lock *sync.Mutex
}

func NewMockDisperserClient(t *testing.T, lock *sync.Mutex, authenticated bool) *MockDisperserClient {
	return &MockDisperserClient{
		t:             t,
		lock:          lock,
		authenticated: authenticated,
	}
}

func (m *MockDisperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	assert.False(m.t, m.authenticated, "writer configured to use non-authenticated disperser method")
	m.ProvidedData = data
	m.ProvidedQuorum = customQuorums
	m.Count++
	return &m.StatusToReturn, m.KeyToReturn, m.ErrorToReturn
}

func (m *MockDisperserClient) DisperseBlobAuthenticated(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	assert.True(m.t, m.authenticated, "writer configured to use authenticated disperser method")
	m.ProvidedData = data
	m.ProvidedQuorum = customQuorums
	m.Count++
	return &m.StatusToReturn, m.KeyToReturn, m.ErrorToReturn
}

func (m *MockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	panic("this method should not be called in this test")
}

func (m *MockDisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	panic("this method should not be called in this test")
}

// TestBasicBehavior tests the basic behavior of the BlobWriter.
func TestBasicBehavior(t *testing.T) {
	initializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := NewMockTicker(startTime)

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

	config := &Config{
		DataSize:         dataSize,
		SignerPrivateKey: signerPrivateKey,
		RandomizeBlobs:   randomizeBlobs,
		CustomQuorums:    customQuorum,
	}

	lock := sync.Mutex{}

	disperserClient := NewMockDisperserClient(t, &lock, authenticated)
	unconfirmedKeyHandler := NewMockUnconfirmedKeyHandler(t, &lock)

	generatorMetrics := metrics.NewMockMetrics()

	writer := NewBlobWriter(
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
		assertEventuallyTrue(t, func() bool {
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
	executeWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
