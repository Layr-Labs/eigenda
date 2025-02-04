package v2

import (
	"github.com/docker/go-units"
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
	RelayReadAmplification uint64
	// By default, this utility reads chunks once. The number of chunk reads is multiplied
	// by this factor. If this is set to 3, then chunks are read back 3 times.
	ValidatorReadAmplification uint64
	// The maximum number of parallel blobs in flight.
	MaxParallelism uint64
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
	}
}

type LoadGenerator struct {
	// The configuration for the load generator.
	config *LoadGeneratorConfig
	// The  test client to use for the load test.
	client *TestClient
	// The time between starting each blob submission.
	submissionPeriod time.Duration
	// if true, the load generator is running.
	alive atomic.Bool
	// The channel to signal when the load generator is finished.
	finishedChan chan struct{}
}

// NewLoadGenerator creates a new LoadGenerator.
func NewLoadGenerator(config *LoadGeneratorConfig, client *TestClient) *LoadGenerator {
	submissionFrequency := time.Duration(config.BytesPerSecond/config.AverageBlobSize) * time.Second
	submissionPeriod := 1.0 / submissionFrequency

	return &LoadGenerator{
		config:           config,
		client:           client,
		submissionPeriod: submissionPeriod,
		alive:            atomic.Bool{},
		finishedChan:     make(chan struct{}),
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
	// unblock Start()
	l.finishedChan <- struct{}{}
	l.alive.Store(false)
}

// run runs the load generator.
func (l *LoadGenerator) run() {
	ticker := time.NewTicker(l.submissionPeriod)
	for l.alive.Load() {
		<-ticker.C
		// TODO limit parallelism
		go l.submitBlob()
	}
}

// Submits a single blob to the network. This function does not return until it reads the blob back
// from the network, which may take tens of seconds.
func (l *LoadGenerator) submitBlob() {

}
