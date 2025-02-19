package load

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/test/v2/client"
)

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

// ReadConfigFile loads a LoadGeneratorConfig from a file.
func ReadConfigFile(filePath string) (*LoadGeneratorConfig, error) {
	configFile, err := client.ResolveTildeInPath(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &LoadGeneratorConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

// NewLoadGenerator creates a new LoadGenerator.
func NewLoadGenerator(
	config *LoadGeneratorConfig,
	client *client.TestClient,
	rand *random.TestRandom) *LoadGenerator {

	bytesPerSecond := config.MBPerSecond * units.MiB
	averageBlobSize := config.AverageBlobSizeMB * units.MiB

	submissionFrequency := bytesPerSecond / averageBlobSize
	submissionPeriod := 1 / submissionFrequency
	submissionPeriodAsDuration := time.Duration(submissionPeriod * float64(time.Second))

	parallelismLimiter := make(chan struct{}, config.MaxParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	metrics := newLoadGeneratorMetrics(client.GetMetricsRegistry())

	return &LoadGenerator{
		ctx:                ctx,
		cancel:             cancel,
		config:             config,
		client:             client,
		rand:               rand,
		submissionPeriod:   submissionPeriodAsDuration,
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
	ctx, cancel := context.WithTimeout(l.ctx, l.config.DispersalTimeout*time.Second)
	l.metrics.startOperation()
	defer func() {
		<-l.parallelismLimiter
		l.metrics.endOperation()
		cancel()
	}()

	// TODO: failure metrics

	payloadSize := int(l.rand.BoundedGaussian(
		l.config.AverageBlobSizeMB*units.MiB,
		l.config.BlobSizeStdDev*units.MiB,
		1.0,
		float64(l.client.GetConfig().MaxBlobSize+1)))
	payload := l.rand.Bytes(payloadSize)

	eigenDACert, err := l.client.DispersePayload(
		ctx,
		l.client.GetConfig().EigenDACertVerifierAddress,
		l.config.Quorums,
		payload)
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
