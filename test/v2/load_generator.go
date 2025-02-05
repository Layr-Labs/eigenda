package v2

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/docker/go-units"
	"math/rand"
	"sync/atomic"
	"time"
)

// LoadGeneratorConfig is the configuration for the load generator.
type LoadGeneratorConfig struct {
	// The desired number of bytes per second to write.
	BytesPerSecond uint64
	// The average size of the blobs to write.
	AverageBlobSize uint64
	// The standard deviation of the blob size.
	BlobSizeStdDev uint64
	// By default, this utility reads each blob back from each relay once. The number of
	// reads per relay is multiplied by this factor. For example, If this is set to 3,
	// then each blob is read back from each relay 3 times.
	RelayReadAmplification uint64 // TODO use this
	// By default, this utility reads chunks once. The number of chunk reads is multiplied
	// by this factor. If this is set to 3, then chunks are read back 3 times.
	ValidatorReadAmplification uint64 // TODO use this
	// The maximum number of parallel blobs in flight.
	MaxParallelism uint64
	// The timeout for each blob dispersal.
	DispersalTimeout time.Duration
	// The quorums to use for the load test.
	Quorums []core.QuorumID
}

// DefaultLoadGeneratorConfig returns the default configuration for the load generator.
func DefaultLoadGeneratorConfig() *LoadGeneratorConfig {
	return &LoadGeneratorConfig{
		BytesPerSecond:             10 * units.MiB,
		AverageBlobSize:            1 * units.MiB,
		BlobSizeStdDev:             0.5 * units.MiB,
		RelayReadAmplification:     1,
		ValidatorReadAmplification: 1,
		MaxParallelism:             1000,
		DispersalTimeout:           5 * time.Minute,
		Quorums:                    []core.QuorumID{0, 1},
	}
}

type LoadGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc

	// The configuration for the load generator.
	config *LoadGeneratorConfig
	// The  test client to use for the load test.
	client *TestClient
	// The random number generator to use for the load test.
	rand *random.TestRandom
	// The time between starting each blob submission.
	submissionPeriod time.Duration
	// The channel to limit the number of parallel blob submissions.
	parallelismLimiter chan struct{}
	// if true, the load generator is running.
	alive atomic.Bool
	// The channel to signal when the load generator is finished.
	finishedChan chan struct{}
}

// NewLoadGenerator creates a new LoadGenerator.
func NewLoadGenerator(
	config *LoadGeneratorConfig,
	client *TestClient,
	rand *random.TestRandom) *LoadGenerator {

	submissionFrequency := config.BytesPerSecond / config.AverageBlobSize
	submissionPeriod := time.Second / time.Duration(submissionFrequency)

	parallelismLimiter := make(chan struct{}, config.MaxParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &LoadGenerator{
		ctx:                ctx,
		cancel:             cancel,
		config:             config,
		client:             client,
		rand:               rand,
		submissionPeriod:   submissionPeriod,
		parallelismLimiter: parallelismLimiter,
		alive:              atomic.Bool{},
		finishedChan:       make(chan struct{}),
	}
}

// Start starts the load generator. If block is true, this function will block until Stop() or
// the load generator crashes. If block is false, this function will return immediately.
func (l *LoadGenerator) Start(block bool) {
	l.alive.Store(true)
	l.run()
	if block {
		<-l.finishedChan
	}
}

// Stop stops the load generator.
func (l *LoadGenerator) Stop() {
	l.finishedChan <- struct{}{}
	l.alive.Store(false)
	l.client.Stop()
	l.cancel()
}

// run runs the load generator.
func (l *LoadGenerator) run() {
	ticker := time.NewTicker(l.submissionPeriod)
	for l.alive.Load() {
		<-ticker.C
		l.parallelismLimiter <- struct{}{}
		go l.submitBlob()
	}
}

// Submits a single blob to the network. This function does not return until it reads the blob back
// from the network, which may take tens of seconds.
func (l *LoadGenerator) submitBlob() {
	ctx, cancel := context.WithTimeout(l.ctx, l.config.DispersalTimeout)
	defer func() {
		<-l.parallelismLimiter
		cancel()
	}()

	payloadSize := int(l.rand.BoundedGaussian(
		float64(l.config.AverageBlobSize),
		float64(l.config.BlobSizeStdDev),
		1.0,
		float64(l.client.Config.MaxBlobSize+1)))
	payload := l.rand.Bytes(payloadSize)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	if uint64(len(paddedPayload)) > l.client.Config.MaxBlobSize {
		paddedPayload = paddedPayload[:l.client.Config.MaxBlobSize]
	}

	key, err := l.client.DispersePayload(ctx, paddedPayload, l.config.Quorums, rand.Uint32())
	if err != nil {
		fmt.Printf("failed to disperse blob: %v\n", err)
	}
	blobCert := l.client.WaitForCertification(ctx, *key, l.config.Quorums)

	// Unpad the payload
	unpaddedPayload := codec.RemoveEmptyByteFromPaddedBytes(paddedPayload)

	// Read the blob from the relays and validators
	for i := uint64(0); i < l.config.RelayReadAmplification; i++ {
		l.client.ReadBlobFromRelays(ctx, *key, blobCert, unpaddedPayload)
	}
	for i := uint64(0); i < l.config.ValidatorReadAmplification; i++ {
		l.client.ReadBlobFromValidators(ctx, blobCert, l.config.Quorums, unpaddedPayload)
	}
}
