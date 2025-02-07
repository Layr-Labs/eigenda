package load

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
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
		RelayReadAmplification:     3,
		ValidatorReadAmplification: 3,
		MaxParallelism:             10,
		DispersalTimeout:           5 * time.Minute,
		Quorums:                    []core.QuorumID{0, 1},
	}
}

type LoadGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc

	// The configuration for the load generator.
	config *LoadGeneratorConfig
	// The test client to use for the load test.
	client *client.TestClient
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
	// The metrics for the load generator.
	metrics *loadGeneratorMetrics
}

// NewLoadGenerator creates a new LoadGenerator.
func NewLoadGenerator(
	config *LoadGeneratorConfig,
	client *client.TestClient,
	rand *random.TestRandom) *LoadGenerator {

	submissionFrequency := config.BytesPerSecond / config.AverageBlobSize
	submissionPeriod := time.Second / time.Duration(submissionFrequency)

	parallelismLimiter := make(chan struct{}, config.MaxParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	metrics := newLoadGeneratorMetrics(client.MetricsRegistry)

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
		metrics:            metrics,
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
	l.metrics.startOperation()
	defer func() {
		<-l.parallelismLimiter
		l.metrics.endOperation()
		cancel()
	}()

	// TODO: failure metrics

	payloadSize := int(l.rand.BoundedGaussian(
		float64(l.config.AverageBlobSize),
		float64(l.config.BlobSizeStdDev),
		1.0,
		float64(l.client.Config.MaxBlobSize+1)))
	payload := l.rand.Bytes(payloadSize)

	eigenDACert, err := l.client.DispersePayload(ctx, payload, rand.Uint32())
	if err != nil {
		fmt.Printf("failed to disperse blob: %v\n", err)
		return
	}

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		fmt.Printf("failed to compute blob key: %v\n", err)
		return
	}

	// Read the blob from the relays and validators
	for i := uint64(0); i < l.config.RelayReadAmplification; i++ {
		err = l.client.ReadBlobFromRelays(
			ctx,
			*blobKey,
			eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys,
			payload)
		if err != nil {
			fmt.Printf("failed to read blob from relays: %v\n", err)
		}
	}

	blobHeader := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader
	commitment, err := verification.BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		fmt.Printf("failed to compute blob commitment: %v\n", err)
	}

	for i := uint64(0); i < l.config.ValidatorReadAmplification; i++ {
		err = l.client.ReadBlobFromValidators(
			ctx,
			*blobKey,
			blobHeader.Version,
			*commitment,
			eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
			payload)
		if err != nil {
			fmt.Printf("failed to read blob from validators: %v\n", err)
		}
	}
}
