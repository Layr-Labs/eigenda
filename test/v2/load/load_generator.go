package load

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
)

// LoadGenerator is a utility for generating read and write load for the target network.
type LoadGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc

	// The configuration for the load generator.
	config *LoadGeneratorConfig
	// The test client to use for the load test.
	client *client.TestClient
	// The time between starting each blob submission.
	submissionPeriod time.Duration
	// The channel to limit the number of parallel blob submissions.
	submissionLimiter chan struct{}
	// The channel to limit the number of parallel blob reads.
	readLimiter chan struct{}
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
	client *client.TestClient) *LoadGenerator {

	bytesPerSecond := config.MBPerSecond * units.MiB
	averageBlobSize := config.AverageBlobSizeMB * units.MiB

	submissionFrequency := bytesPerSecond / averageBlobSize
	submissionPeriod := 1 / submissionFrequency
	submissionPeriodAsDuration := time.Duration(submissionPeriod * float64(time.Second))

	submissionLimiter := make(chan struct{}, config.SubmissionParallelism)
	readLimiter := make(chan struct{}, config.ReadParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	metrics := newLoadGeneratorMetrics(client.GetMetricsRegistry())

	if config.EnablePprof {
		pprofProfiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", config.PprofHttpPort), client.GetLogger())
		go pprofProfiler.Start()
		client.GetLogger().Info("Enabled pprof", "port", config.PprofHttpPort)
	}

	client.SetCertVerifierAddress(client.GetConfig().EigenDACertVerifierAddressQuorums0_1)

	return &LoadGenerator{
		ctx:               ctx,
		cancel:            cancel,
		config:            config,
		client:            client,
		submissionPeriod:  submissionPeriodAsDuration,
		submissionLimiter: submissionLimiter,
		readLimiter:       readLimiter,
		alive:             atomic.Bool{},
		finishedChan:      make(chan struct{}),
		metrics:           metrics,
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
		go l.readAndWriteBlob()
	}
}

func (l *LoadGenerator) readAndWriteBlob() {
	rand := random.NewTestRandomNoPrint()

	l.submissionLimiter <- struct{}{}
	blobKey, payload, eigenDACert, err := l.submitBlob(rand)
	<-l.submissionLimiter
	if err != nil {
		return
	}

	l.readLimiter <- struct{}{}
	l.readBlob(rand, blobKey, payload, eigenDACert)
	<-l.readLimiter
}

// Submits a single blob to the network.
func (l *LoadGenerator) submitBlob(rand *random.TestRandom) (
	blobKey *corev2.BlobKey,
	payload []byte,
	eigenDACert *coretypes.EigenDACert,
	err error) {

	payloadSize := int(rand.BoundedGaussian(
		l.config.AverageBlobSizeMB*units.MiB,
		l.config.BlobSizeStdDev*units.MiB,
		1.0,
		float64(l.client.GetConfig().MaxBlobSize+1)))
	payload = rand.Bytes(payloadSize)

	ctx, cancel := context.WithTimeout(l.ctx, l.config.DispersalTimeout*time.Second)
	l.metrics.startOperation("write")
	defer func() {
		l.metrics.endOperation("write")
		cancel()
	}()

	eigenDACert, err = l.client.DispersePayload(ctx, payload)
	if err != nil {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("failed to disperse blob: %v", err)
		return nil, nil, nil, err
	}

	blobKey, err = eigenDACert.ComputeBlobKey()
	if err != nil {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("failed to compute blob key: %v", err)
		return nil, nil, nil, err
	}

	l.metrics.reportDispersalSuccess()
	return blobKey, payload, eigenDACert, nil
}

// readBlob reads a blob from the network. May read from relays or validators or both. May read multiple times.
func (l *LoadGenerator) readBlob(
	rand *random.TestRandom,
	blobKey *corev2.BlobKey,
	payload []byte,
	eigenDACert *coretypes.EigenDACert) {

	ctx, cancel := context.WithTimeout(l.ctx, l.config.ReadTimeout*time.Second)
	l.metrics.startOperation("read")
	defer func() {
		l.metrics.endOperation("read")
		cancel()
	}()

	blobLengthSymbols := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length

	var relayReadCount int
	if l.config.RelayReadAmplification < 1 {
		if rand.Float64() < l.config.RelayReadAmplification {
			relayReadCount = 1
		}
	} else {
		relayReadCount = int(l.config.RelayReadAmplification)
	}

	for i := 0; i < relayReadCount; i++ {
		err := l.client.ReadBlobFromRelays(
			ctx,
			*blobKey,
			eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys,
			payload,
			blobLengthSymbols)
		if err == nil {
			l.metrics.reportRelayReadSuccess()
		} else {
			l.metrics.reportRelayReadFailure()
			l.client.GetLogger().Errorf("failed to read blob from relays: %v", err)
		}
	}

	blobHeader := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader
	commitment, err := coretypes.BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		l.client.GetLogger().Errorf("failed to bind blob commitments: %v", err)
		return
	}

	var validatorReadCount int
	if l.config.ValidatorReadAmplification < 1 {
		if rand.Float64() < l.config.ValidatorReadAmplification {
			validatorReadCount = 1
		}
	} else {
		validatorReadCount = int(l.config.ValidatorReadAmplification)
	}

	for i := 0; i < validatorReadCount; i++ {
		err = l.client.ReadBlobFromValidators(
			ctx,
			*blobKey,
			blobHeader.Version,
			*commitment,
			eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
			payload)
		if err == nil {
			l.metrics.reportValidatorReadSuccess()
		} else {
			l.metrics.reportValidatorReadFailure()
			l.client.GetLogger().Errorf("failed to read blob from validators: %v", err)
		}
	}
}
