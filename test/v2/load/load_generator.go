package load

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
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
	// The frequency at which blobs are submitted, in HZ.
	submissionFrequency float64
	// The channel to limit the number of parallel blob submissions.
	submissionLimiter chan struct{}
	// The channel to limit the number of parallel blob reads sent to the relays.
	relayReadLimiter chan struct{}
	// The channel to limit the number of parallel blob reads sent to the validators.
	validatorReadLimiter chan struct{}
	// The channel to limit the number of blobs in all phases of the read/write lifecycle.
	lifecycleLimiter chan struct{}
	// if true, the load generator is running.
	alive atomic.Bool
	// The channel to signal when the load generator is finished.
	finishedChan chan struct{}
	// The metrics for the load generator.
	metrics *loadGeneratorMetrics
	// Pool of random number generators
	randPool *sync.Pool
	// The time when the load generator started.
	startTime time.Time
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

	config := DefaultLoadGeneratorConfig()
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

	submissionLimiter := make(chan struct{}, config.SubmissionParallelism)
	relayReadLimiter := make(chan struct{}, config.RelayReadParallelism)
	validatorReadLimiter := make(chan struct{}, config.ValidatorReadParallelism)
	lifecycleLimiter := make(chan struct{},
		config.SubmissionParallelism+
			config.RelayReadParallelism+
			config.ValidatorReadParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	metrics := newLoadGeneratorMetrics(client.GetMetricsRegistry())

	if config.EnablePprof {
		pprofProfiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", config.PprofHttpPort), client.GetLogger())
		go pprofProfiler.Start()
		client.GetLogger().Info("Enabled pprof", "port", config.PprofHttpPort)
	}

	client.SetCertVerifierAddress(client.GetConfig().EigenDACertVerifierAddressQuorums0_1)

	// Initialize a pool for random number generators
	randPool := &sync.Pool{
		New: func() interface{} {
			return random.NewTestRandomNoPrint()
		},
	}

	return &LoadGenerator{
		ctx:                  ctx,
		cancel:               cancel,
		config:               config,
		client:               client,
		submissionFrequency:  submissionFrequency,
		submissionLimiter:    submissionLimiter,
		relayReadLimiter:     relayReadLimiter,
		lifecycleLimiter:     lifecycleLimiter,
		validatorReadLimiter: validatorReadLimiter,
		alive:                atomic.Bool{},
		finishedChan:         make(chan struct{}),
		randPool:             randPool,
		metrics:              metrics,
		startTime:            time.Now(),
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

	// Start with frequency 0.
	ticker, err := common.NewVariableTickerWithFrequency(l.ctx, 0)
	if err != nil {
		// Not possible, error is only returned with invalid arguments, and 0hz is a valid frequency.
		panic(fmt.Errorf("failed to create variable ticker: %w", err))
	}

	defer ticker.Close()
	// Set acceleration prior to setting target frequency, since acceleration 0 allows "infinite" acceleration.
	err = ticker.SetAcceleration(l.config.FrequencyAcceleration)
	if err != nil {
		// load generator configuration error, no way to recover
		panic(fmt.Errorf("failed to set acceleration: %w", err))
	}
	err = ticker.SetTargetFrequency(l.submissionFrequency)
	if err != nil {
		// load generator configuration error, no way to recover
		panic(fmt.Errorf("failed to set target frequency: %w", err))
	}

	for l.alive.Load() {
		<-ticker.Tick()

		l.lifecycleLimiter <- struct{}{}
		go func() {
			l.readAndWriteBlob()
			<-l.lifecycleLimiter
		}()
	}
}

func (l *LoadGenerator) readAndWriteBlob() {
	// Get a random generator from the pool
	randObj := l.randPool.Get()
	rand := randObj.(*random.TestRandom)
	defer l.randPool.Put(randObj) // Return to pool when done

	l.submissionLimiter <- struct{}{}
	blobKey, payload, eigenDACert, err := l.disperseBlob(rand)
	<-l.submissionLimiter
	if err != nil {
		return
	}

	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("expected EigenDACertV3, got %T", eigenDACert)
		return
	}

	l.relayReadLimiter <- struct{}{}
	l.readBlobFromRelays(rand, blobKey, payload, eigenDAV3Cert)
	<-l.relayReadLimiter

	l.validatorReadLimiter <- struct{}{}
	l.readBlobFromValidators(rand, payload, eigenDAV3Cert)
	<-l.validatorReadLimiter
}

// Submits a single blob to the network.
func (l *LoadGenerator) disperseBlob(rand *random.TestRandom) (
	blobKey *corev2.BlobKey,
	payload []byte,
	eigenDACert coretypes.EigenDACert,
	err error) {

	payloadSize := int(rand.BoundedGaussian(
		l.config.AverageBlobSizeMB*units.MiB,
		l.config.BlobSizeStdDev*units.MiB,
		1.0,
		float64(l.client.GetConfig().MaxBlobSize+1)))
	payload = rand.Bytes(payloadSize)

	timeout := time.Duration(l.config.DispersalTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	l.metrics.startOperation("dispersal")
	defer func() {
		l.metrics.endOperation("dispersal")
		cancel()
	}()

	eigenDACert, err = l.client.DispersePayload(ctx, payload)
	if err != nil {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("failed to disperse blob: %v", err)
		return nil, nil, nil, err
	}

	// Ensure the eigenDACert is of type EigenDACertV3
	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("expected EigenDACertV3, got %T", eigenDACert)
		return nil, nil, nil, fmt.Errorf("expected EigenDACertV3, got %T", eigenDACert)
	}

	blobKey, err = eigenDAV3Cert.ComputeBlobKey()
	if err != nil {
		l.metrics.reportDispersalFailure()
		l.client.GetLogger().Errorf("failed to compute blob key: %v", err)
		return nil, nil, nil, err
	}

	l.metrics.reportDispersalSuccess()
	return blobKey, payload, eigenDACert, nil
}

// readBlobFromRelays reads a blob from the relays.
func (l *LoadGenerator) readBlobFromRelays(
	rand *random.TestRandom,
	blobKey *corev2.BlobKey,
	payload []byte,
	eigenDACert *coretypes.EigenDACertV3,
) {

	timeout := time.Duration(l.config.RelayReadTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	var relayReadCount int
	if l.config.RelayReadAmplification < 1 {
		if rand.Float64() < l.config.RelayReadAmplification {
			relayReadCount = 1
		} else {
			return
		}
	} else {
		relayReadCount = int(l.config.RelayReadAmplification)
	}

	l.metrics.startOperation("relay_read")
	defer func() {
		l.metrics.endOperation("relay_read")
	}()

	blobLengthSymbols := uint32(eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length)
	relayKeys := eigenDACert.RelayKeys()
	readStartIndex := rand.Int32Range(0, int32(len(relayKeys)))

	for i := 0; i < relayReadCount; i++ {
		err := l.client.ReadBlobFromRelay(
			ctx,
			*blobKey,
			relayKeys[(int(readStartIndex)+i)%len(relayKeys)],
			payload,
			blobLengthSymbols,
			0)
		if err == nil {
			l.metrics.reportRelayReadSuccess()
		} else {
			l.metrics.reportRelayReadFailure()
			l.client.GetLogger().Errorf("failed to read blob from relays: %v", err)
		}
	}
}

// readBlobFromValidators reads a blob from the validators using the validator retrieval client.
func (l *LoadGenerator) readBlobFromValidators(
	rand *random.TestRandom,
	payload []byte,
	eigenDACert *coretypes.EigenDACertV3) {

	timeout := time.Duration(l.config.ValidatorReadTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	var validatorReadCount int
	if l.config.ValidatorReadAmplification < 1 {
		if rand.Float64() < l.config.ValidatorReadAmplification {
			validatorReadCount = 1
		} else {
			return
		}
	} else {
		validatorReadCount = int(l.config.ValidatorReadAmplification)
	}

	l.metrics.startOperation("validator_read")
	defer func() {
		l.metrics.endOperation("validator_read")
	}()

	blobHeader, err := eigenDACert.BlobHeader()
	if err != nil {
		l.metrics.reportValidatorReadFailure()
		l.client.GetLogger().Errorf("failed to get blob header: %v", err)
		return
	}

	for i := 0; i < validatorReadCount; i++ {
		validateAndDecode := rand.Float64() < l.config.ValidatorVerificationFraction

		err = l.client.ReadBlobFromValidators(
			ctx,
			blobHeader,
			uint32(eigenDACert.ReferenceBlockNumber()),
			payload,
			0,
			validateAndDecode)
		if err == nil {
			l.metrics.reportValidatorReadSuccess()
		} else {
			l.metrics.reportValidatorReadFailure()
			l.client.GetLogger().Errorf("failed to read blob from validators: %v", err)
		}
	}
}
