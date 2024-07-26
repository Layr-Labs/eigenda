package workers

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
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

// MockUnconfirmedKeyHandler is a stand-in for the blob verifier's UnconfirmedKeyHandler.
type MockUnconfirmedKeyHandler struct {
	t        *testing.T
	key      *[]byte
	checksum *[16]byte
	size     uint
}

func (m *MockUnconfirmedKeyHandler) AddUnconfirmedKey(key *[]byte, checksum *[16]byte, size uint) {
	// Ensure that we have already verified the previous key.
	assert.Nil(m.t, m.key)

	m.key = key
	m.checksum = checksum
	m.size = size
}

// GetAndClear returns the unconfirmed key and clears the internal state. Must be called after each
// AddUnconfirmedKey call.
func (m *MockUnconfirmedKeyHandler) GetAndCLear() (*[]byte, *[16]byte, uint) {
	defer func() {
		m.key = nil
		m.checksum = nil
		m.size = 0
	}()
	return m.key, m.checksum, m.size
}

// TestBasicBehavior tests the basic behavior of the BlobWriter with no special cases.
func TestBasicBehavior(t *testing.T) {
	initializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	// TODO use deterministic start time
	ticker := NewMockTicker(time.Now())

	dataSize := rand.Uint64()%1024 + 64

	config := &Config{
		DataSize: dataSize,
	}
	var disperser *clients.DisperserClient // TODO create mock

	unconfirmedKeyHandler := &MockUnconfirmedKeyHandler{
		t: t,
	}

	generatorMetrics := metrics.NewMockMetrics()

	writer := NewBlobWriter(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		disperser,
		unconfirmedKeyHandler,
		generatorMetrics)
	writer.Start()

	ticker.Tick(1 * time.Second)
	_, _, _ = unconfirmedKeyHandler.GetAndCLear() // TODO

	// ...

	cancel()
	// TODO add timeout
	waitGroup.Wait()
}
