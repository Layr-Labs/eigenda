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
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/common/pprof"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigenda/test/random"
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
	// The channel to limit the number of parallel gas estimation operations.
	gasEstimationLimiter chan struct{}
	// The channel to limit the number of blobs in all phases of the read/write lifecycle.
	lifecycleLimiter chan struct{}
	// if true, the load generator is running.
	alive atomic.Bool
	// The channel to signal when the load generator is finished.
	finishedChan chan struct{}
	// Pool of random number generators
	randPool *sync.Pool
	// The time when the load generator started.
	startTime time.Time
	// The size of the payload that will result in an encoded blob of the target size.
	payloadSize uint32
}

// ReadConfigFile loads a LoadGeneratorConfig from a file.
func ReadConfigFile(filePath string) (*LoadGeneratorConfig, error) {
	configFile, err := util.SanitizePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize path: %w", err)
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
	client *client.TestClient) (*LoadGenerator, error) {

	bytesPerSecond := config.MBPerSecond * units.MiB

	// The size of the blob we want to send.
	targetBlobSize := uint64(config.BlobSizeMB * units.MiB)
	// The target blob size must be a power of 2.
	targetBlobSize = math.NextPowOf2u64(targetBlobSize)

	// The size of the payload necessary to create a blob of the target size.
	payloadSize, err := codec.BlobSizeToMaxPayloadSize(uint32(targetBlobSize))
	if err != nil {
		return nil, fmt.Errorf("failed to compute payload size for target blob size %d: %w", targetBlobSize, err)
	}

	submissionFrequency := bytesPerSecond / float64(targetBlobSize)

	client.GetLogger().Infof("Target blob size: %s bytes, submission frequency: %f hz",
		common.PrettyPrintBytes(targetBlobSize), submissionFrequency)

	submissionLimiter := make(chan struct{}, config.SubmissionParallelism)
	relayReadLimiter := make(chan struct{}, config.RelayReadParallelism)
	validatorReadLimiter := make(chan struct{}, config.ValidatorReadParallelism)
	gasEstimationLimiter := make(chan struct{}, config.GasEstimationParallelism)
	lifecycleLimiter := make(chan struct{},
		config.SubmissionParallelism+
			config.RelayReadParallelism+
			config.ValidatorReadParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if config.EnablePprof {
		pprofProfiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", config.PprofHttpPort), client.GetLogger())
		go pprofProfiler.Start()
		client.GetLogger().Info("Enabled pprof", "port", config.PprofHttpPort)
	}

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
		gasEstimationLimiter: gasEstimationLimiter,
		lifecycleLimiter:     lifecycleLimiter,
		validatorReadLimiter: validatorReadLimiter,
		alive:                atomic.Bool{},
		finishedChan:         make(chan struct{}),
		randPool:             randPool,
		startTime:            time.Now(),
		payloadSize:          payloadSize,
	}, nil
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
	ticker, err := common.NewVariableTickerWithFrequency(l.ctx, l.client.GetLogger(), 0)
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
			if l.config.UseProxy {
				l.readAndWriteBlobWithProxy()
			} else {
				l.readAndWriteBlob()
			}
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
		l.client.GetLogger().Errorf("failed to disperse blob: %w", err)
		return
	}

	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
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

// Submits a single blob to the network using the GRPC clients.
func (l *LoadGenerator) disperseBlob(rand *random.TestRandom) (
	blobKey corev2.BlobKey,
	payload []byte,
	eigenDACert coretypes.EigenDACert,
	err error) {

	payload = rand.Bytes(int(l.payloadSize))

	timeout := time.Duration(l.config.DispersalTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	eigenDACert, err = l.client.DispersePayload(ctx, payload)
	if err != nil {
		l.client.GetLogger().Errorf("failed to disperse blob: %v", err)
		return corev2.BlobKey{}, nil, nil, fmt.Errorf("failed to disperse blob: %w", err)
	}

	// Ensure the eigenDACert is of type EigenDACertV3
	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
		l.client.GetLogger().Errorf("expected EigenDACertV3, got %T", eigenDACert)
		return corev2.BlobKey{}, nil, nil, fmt.Errorf("expected EigenDACertV3, got %T", eigenDACert)
	}

	blobKey, err = eigenDAV3Cert.ComputeBlobKey()
	if err != nil {
		l.client.GetLogger().Errorf("failed to compute blob key: %v", err)
		return corev2.BlobKey{}, nil, nil, fmt.Errorf("failed to compute blob key: %w", err)
	}

	// Estimate gas for CheckDACert call
	go l.estimateAndReportGasCheckDACert(eigenDAV3Cert)

	return blobKey, payload, eigenDACert, nil
}

// estimateAndReportGasCheckDACert performs gas estimation and reports it as a metric.
// Make sure to call this in a separate goroutine to avoid blocking blob dispersal.
func (l *LoadGenerator) estimateAndReportGasCheckDACert(eigenDAV3Cert *coretypes.EigenDACertV3) {
	l.gasEstimationLimiter <- struct{}{}
	defer func() {
		<-l.gasEstimationLimiter
	}()

	gasTimeout := time.Duration(l.config.GasEstimationTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, gasTimeout)
	defer cancel()

	_, err := l.client.EstimateGasAndReportCheckDACert(ctx, eigenDAV3Cert)
	if err != nil {
		l.client.GetLogger().Errorf("failed to estimate gas for CheckDACert call: %v", err)
	}
}

func (l *LoadGenerator) readAndWriteBlobWithProxy() {
	// Get a random generator from the pool
	randObj := l.randPool.Get()
	rand := randObj.(*random.TestRandom)
	defer l.randPool.Put(randObj) // Return to pool when done

	l.submissionLimiter <- struct{}{}

	cert, payload, err := l.dispersePayloadWithProxy(rand)
	<-l.submissionLimiter
	if err != nil {
		l.client.GetLogger().Errorf("failed to disperse blob: %w", err)
		return
	}

	l.relayReadLimiter <- struct{}{}
	err = l.doReadsWithProxy(rand, cert, payload)
	<-l.relayReadLimiter
	if err != nil {
		l.client.GetLogger().Errorf("failed to read blob from proxy: %w", err)
	}
}

// Disperses a blob using the proxy (as opposed to using the GRPC clients directly). Returns the blob cert in byte
// form since this is how the proxy forces the user to interact with it.
func (l *LoadGenerator) dispersePayloadWithProxy(rand *random.TestRandom) (
	cert []byte,
	payload []byte,
	err error,
) {

	payload = rand.Bytes(int(l.payloadSize))

	timeout := time.Duration(l.config.DispersalTimeout) * time.Second
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	cert, err = l.client.DispersePayloadWithProxy(ctx, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to disperse blob with proxy: %w", err)
	}

	return cert, payload, nil
}

// Reads the blob using the proxy client. The proxy may in theory read the blob from the relays or validators, but
// unless the relays are malfunctioning it will always read from the relays.
func (l *LoadGenerator) doReadsWithProxy(
	rand *random.TestRandom,
	cert []byte,
	expectedPayload []byte,
) error {

	var readCount int
	if l.config.RelayReadAmplification < 1 {
		if rand.Float64() < l.config.RelayReadAmplification {
			readCount = 1
		} else {
			return nil // Skip reading this time
		}
	} else {
		readCount = int(l.config.RelayReadAmplification)
	}

	for i := 0; i < readCount; i++ {
		_, err := l.client.ReadPayloadWithProxy(l.ctx, cert, expectedPayload, 0)
		if err != nil {
			return fmt.Errorf("failed to read blob from proxy: %w", err)
		}
	}

	return nil
}

// readBlobFromRelays reads a blob from the relays using the GRPC clients.
func (l *LoadGenerator) readBlobFromRelays(
	rand *random.TestRandom,
	blobKey corev2.BlobKey,
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

	blobLengthSymbols := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length
	relayKeys := eigenDACert.RelayKeys()
	readStartIndex := rand.Int32Range(0, int32(len(relayKeys)))

	for i := 0; i < relayReadCount; i++ {
		err := l.client.ReadBlobFromRelay(
			ctx,
			blobKey,
			relayKeys[(int(readStartIndex)+i)%len(relayKeys)],
			payload,
			blobLengthSymbols,
			0)
		if err != nil {
			l.client.GetLogger().Errorf("failed to read blob from relays: %v", err)
		}
	}
}

// readBlobFromValidators reads a blob from the validators using the validator GRPC client.
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

	blobHeader, err := eigenDACert.BlobHeader()
	if err != nil {
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
		if err != nil {
			l.client.GetLogger().Errorf("failed to read blob from validators: %v", err)
		}
	}
}
