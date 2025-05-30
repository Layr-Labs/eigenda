package bench

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OperationType represents different metadata store operations
type OperationType string

const (
	OpUpdateBlobStatus     OperationType = "UpdateBlobStatus"
	OpPutBlobMetadata      OperationType = "PutBlobMetadata"
	OpGetBlobMetadata      OperationType = "GetBlobMetadata"
	OpPutBlobCertificate   OperationType = "PutBlobCertificate"
	OpGetBlobCertificate   OperationType = "GetBlobCertificate"
	OpPutBatch             OperationType = "PutBatch"
	OpGetBatch             OperationType = "GetBatch"
	OpPutDispersalRequest  OperationType = "PutDispersalRequest"
	OpGetDispersalRequest  OperationType = "GetDispersalRequest"
	OpPutAttestation       OperationType = "PutAttestation"
	OpGetAttestation       OperationType = "GetAttestation"
	OpPutBlobInclusionInfo OperationType = "PutBlobInclusionInfo"
	OpGetBlobInclusionInfo OperationType = "GetBlobInclusionInfo"
)

// StoreType represents the type of metadata store
type StoreType string

const (
	DynamoDB   StoreType = "dynamodb"
	PostgreSQL StoreType = "postgresql"
)

// OperationConfig defines the configuration for a specific operation
type OperationConfig struct {
	Type       OperationType
	RatePerSec int           // Target operations per second
	Duration   time.Duration // How long to run this operation
	Workers    int           // Number of concurrent workers
}

// BenchmarkConfig defines the overall benchmark configuration
type BenchmarkConfig struct {
	StoreType      StoreType
	Operations     []OperationConfig
	WarmupTime     time.Duration // Time to warm up before collecting metrics
	ReportInterval time.Duration // How often to report intermediate results

	// Store-specific configurations
	DynamoDBConfig *blobstore.DynamoDBConfig
	PostgresConfig *blobstore.PostgresConfig
}

// OperationResult stores the result of a single operation
type OperationResult struct {
	Operation OperationType
	StartTime time.Time
	Duration  time.Duration
	Success   bool
	Error     error
}

// MetricsCollector collects and aggregates benchmark metrics
type MetricsCollector struct {
	mu            sync.RWMutex
	results       map[OperationType][]OperationResult
	startTime     time.Time
	totalOps      map[OperationType]*atomic.Int64
	successfulOps map[OperationType]*atomic.Int64
	failedOps     map[OperationType]*atomic.Int64
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		results:       make(map[OperationType][]OperationResult),
		totalOps:      make(map[OperationType]*atomic.Int64),
		successfulOps: make(map[OperationType]*atomic.Int64),
		failedOps:     make(map[OperationType]*atomic.Int64),
		startTime:     time.Now(),
	}
}

// RecordOperation records the result of an operation
func (mc *MetricsCollector) RecordOperation(result OperationResult) {
	mc.mu.Lock()
	mc.results[result.Operation] = append(mc.results[result.Operation], result)
	mc.mu.Unlock()

	if mc.totalOps[result.Operation] == nil {
		mc.totalOps[result.Operation] = &atomic.Int64{}
		mc.successfulOps[result.Operation] = &atomic.Int64{}
		mc.failedOps[result.Operation] = &atomic.Int64{}
	}

	mc.totalOps[result.Operation].Add(1)
	if result.Success {
		mc.successfulOps[result.Operation].Add(1)
	} else {
		mc.failedOps[result.Operation].Add(1)
	}
}

// GetMetrics returns aggregated metrics for an operation
func (mc *MetricsCollector) GetMetrics(op OperationType) OperationMetrics {
	mc.mu.RLock()
	results := mc.results[op]
	mc.mu.RUnlock()

	if len(results) == 0 {
		return OperationMetrics{Operation: op}
	}

	// Calculate latencies
	latencies := make([]float64, 0, len(results))
	for _, r := range results {
		if r.Success {
			latencies = append(latencies, float64(r.Duration.Nanoseconds())/1_000_000)
		}
	}

	sort.Float64s(latencies)

	metrics := OperationMetrics{
		Operation:     op,
		TotalOps:      mc.totalOps[op].Load(),
		SuccessfulOps: mc.successfulOps[op].Load(),
		FailedOps:     mc.failedOps[op].Load(),
		Duration:      time.Since(mc.startTime),
	}

	if len(latencies) > 0 {
		metrics.MinLatencyMs = latencies[0]
		metrics.MaxLatencyMs = latencies[len(latencies)-1]
		metrics.AvgLatencyMs = average(latencies)
		metrics.P50LatencyMs = percentile(latencies, 0.50)
		metrics.P90LatencyMs = percentile(latencies, 0.90)
		metrics.P95LatencyMs = percentile(latencies, 0.95)
		metrics.P99LatencyMs = percentile(latencies, 0.99)
	}

	metrics.Throughput = float64(metrics.TotalOps) / metrics.Duration.Seconds()
	if metrics.TotalOps > 0 {
		metrics.ErrorRate = float64(metrics.FailedOps) / float64(metrics.TotalOps) * 100
	}

	return metrics
}

// OperationMetrics contains aggregated metrics for an operation type
type OperationMetrics struct {
	Operation     OperationType
	TotalOps      int64
	SuccessfulOps int64
	FailedOps     int64
	Throughput    float64 // ops/sec
	ErrorRate     float64 // percentage
	Duration      time.Duration

	// Latency percentiles in milliseconds
	MinLatencyMs float64
	MaxLatencyMs float64
	AvgLatencyMs float64
	P50LatencyMs float64
	P90LatencyMs float64
	P95LatencyMs float64
	P99LatencyMs float64
}

// BenchmarkRunner runs the benchmark
type BenchmarkRunner struct {
	config    BenchmarkConfig
	store     blobstore.MetadataStore
	collector *MetricsCollector
	logger    logging.Logger

	// Test data
	testBlobs     []corev2.BlobKey
	testBatches   [][32]byte
	testOperators []core.OperatorID
	
	// Blob state tracking for efficient UpdateBlobStatus benchmarking
	blobStates    map[corev2.BlobKey]v2.BlobStatus
	mu            sync.RWMutex
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner(config BenchmarkConfig, store blobstore.MetadataStore, logger logging.Logger) *BenchmarkRunner {
	return &BenchmarkRunner{
		config:     config,
		store:      store,
		collector:  NewMetricsCollector(),
		logger:     logger,
		blobStates: make(map[corev2.BlobKey]v2.BlobStatus),
	}
}

// Run executes the benchmark
func (br *BenchmarkRunner) Run(ctx context.Context) error {
	br.logger.Info("Starting benchmark", "store_type", br.config.StoreType)

	// Setup test data
	if err := br.setupTestData(ctx); err != nil {
		return fmt.Errorf("failed to setup test data: %w", err)
	}

	// Warmup phase
	if br.config.WarmupTime > 0 {
		br.logger.Info("Starting warmup phase", "duration", br.config.WarmupTime)
		warmupCtx, cancel := context.WithTimeout(ctx, br.config.WarmupTime)
		br.runOperations(warmupCtx, false) // Don't collect metrics during warmup
		cancel()

		// Reset metrics after warmup
		br.collector = NewMetricsCollector()
	}

	// Main benchmark phase
	br.logger.Info("Starting main benchmark phase")

	// Start report ticker
	reportTicker := time.NewTicker(br.config.ReportInterval)
	defer reportTicker.Stop()

	// Start operations
	opCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, opConfig := range br.config.Operations {
		wg.Add(1)
		go func(cfg OperationConfig) {
			defer wg.Done()
			br.runOperation(opCtx, cfg, true)
		}(opConfig)
	}

	// Report intermediate results
	go func() {
		for {
			select {
			case <-reportTicker.C:
				br.printIntermediateResults()
			case <-opCtx.Done():
				return
			}
		}
	}()

	// Wait for all operations to complete
	wg.Wait()

	// Print final results
	br.printFinalResults()

	return nil
}

// setupTestData creates test data for the benchmark
func (br *BenchmarkRunner) setupTestData(ctx context.Context) error {
	br.logger.Info("Setting up test data")

	numTestItems := 10000 // Adjust based on needs

	// Generate test blob keys
	br.testBlobs = make([]corev2.BlobKey, numTestItems)
	for i := 0; i < numTestItems; i++ {
		rand.Read(br.testBlobs[i][:])
	}

	// Generate test batch hashes
	br.testBatches = make([][32]byte, numTestItems)
	for i := 0; i < numTestItems; i++ {
		rand.Read(br.testBatches[i][:])
	}

	// Generate test operator IDs
	br.testOperators = make([]core.OperatorID, 100)
	for i := 0; i < 100; i++ {
		rand.Read(br.testOperators[i][:])
	}

	// Pre-populate some data for read operations with various statuses
	// Also track their states in memory for efficient UpdateBlobStatus benchmarking
	for i := 0; i < min(1000, numTestItems); i++ {
		metadata := br.generateTestBlobMetadataWithRandomStatus(br.testBlobs[i])
		if err := br.store.PutBlobMetadata(ctx, metadata); err != nil {
			br.logger.Warn("Failed to pre-populate blob metadata", "error", err)
			continue
		}
		
		// Get the actual blob key from the header (this might be different from br.testBlobs[i])
		actualBlobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			br.logger.Warn("Failed to get blob key from header", "error", err)
			continue
		}
		
		// Track the blob state in memory using the actual blob key
		br.mu.Lock()
		br.blobStates[actualBlobKey] = metadata.BlobStatus
		br.mu.Unlock()
		
		// Update our test blobs array to use the actual keys
		br.testBlobs[i] = actualBlobKey
	}

	return nil
}

// runOperation runs a specific operation type at the specified rate
func (br *BenchmarkRunner) runOperation(ctx context.Context, config OperationConfig, collectMetrics bool) {
	rateLimiter := time.NewTicker(time.Second / time.Duration(config.RatePerSec))
	defer rateLimiter.Stop()

	opCtx, cancel := context.WithTimeout(ctx, config.Duration)
	defer cancel()

	// Create worker pool
	workChan := make(chan struct{}, config.Workers)
	var wg sync.WaitGroup

	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-workChan:
					result := br.executeOperation(opCtx, config.Type)
					if collectMetrics {
						br.collector.RecordOperation(result)
					}
				case <-opCtx.Done():
					return
				}
			}
		}()
	}

	// Generate load at specified rate
	for {
		select {
		case <-rateLimiter.C:
			select {
			case workChan <- struct{}{}:
			default:
				// Workers are busy, skip this tick
				br.logger.Debug("Workers busy, skipping tick", "operation", config.Type)
			}
		case <-opCtx.Done():
			close(workChan)
			wg.Wait()
			return
		}
	}
}

// executeOperation executes a single operation and returns the result
func (br *BenchmarkRunner) executeOperation(ctx context.Context, opType OperationType) OperationResult {
	start := time.Now()
	var err error

	switch opType {
	case OpUpdateBlobStatus:
		err = br.executeUpdateBlobStatus(ctx)
	default:
		err = fmt.Errorf("unknown operation type: %s", opType)
	}

	return OperationResult{
		Operation: opType,
		StartTime: start,
		Duration:  time.Since(start),
		Success:   err == nil,
		Error:     err,
	}
}

// Individual operation implementations
func (br *BenchmarkRunner) executeUpdateBlobStatus(ctx context.Context) error {
	// Pick a random blob from our tracked blobs
	blobKey := br.getRandomTrackedBlobKey()
	if blobKey == nil {
		return fmt.Errorf("no tracked blobs available")
	}

	// Get current status from our in-memory tracking (no DB read needed)
	br.mu.RLock()
	currentStatus, exists := br.blobStates[*blobKey]
	br.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("blob not found in tracking map")
	}

	// Determine valid next status based on current status
	nextStatus, err := br.getValidNextStatus(currentStatus)
	if err != nil {
		return err
	}

	// If the blob is already in a terminal state and we're trying to set it to the same state,
	// that's a successful no-op
	if currentStatus == nextStatus {
		return nil
	}

	// Perform the actual UpdateBlobStatus operation (this is what we're benchmarking)
	err = br.store.UpdateBlobStatus(ctx, *blobKey, nextStatus)
	if err == nil {
		// Update our in-memory tracking on success
		br.mu.Lock()
		br.blobStates[*blobKey] = nextStatus
		br.mu.Unlock()
	} else {
		// Log the specific error for debugging
		br.logger.Debug("UpdateBlobStatus failed", "error", err, "blobKey", blobKey.Hex(), "currentStatus", currentStatus.String(), "nextStatus", nextStatus.String())
	}
	
	return err
}

// Test data generation helpers
func (br *BenchmarkRunner) generateRandomBlobKey() corev2.BlobKey {
	var key corev2.BlobKey
	rand.Read(key[:])
	return key
}

func (br *BenchmarkRunner) getRandomBlobKey() corev2.BlobKey {
	br.mu.RLock()
	defer br.mu.RUnlock()
	if len(br.testBlobs) == 0 {
		return br.generateRandomBlobKey()
	}
	return br.testBlobs[rand.Intn(len(br.testBlobs))]
}

func (br *BenchmarkRunner) getRandomTrackedBlobKey() *corev2.BlobKey {
	br.mu.RLock()
	defer br.mu.RUnlock()
	
	if len(br.blobStates) == 0 {
		return nil
	}
	
	// Convert map keys to slice for random selection
	keys := make([]corev2.BlobKey, 0, len(br.blobStates))
	for key := range br.blobStates {
		keys = append(keys, key)
	}
	
	selectedKey := keys[rand.Intn(len(keys))]
	return &selectedKey
}

// getValidNextStatus returns a valid next status for state transition
func (br *BenchmarkRunner) getValidNextStatus(currentStatus v2.BlobStatus) (v2.BlobStatus, error) {
	switch currentStatus {
	case v2.Queued:
		// From Queued, can go to Encoded or Failed
		nextStates := []v2.BlobStatus{v2.Encoded, v2.Failed}
		return nextStates[rand.Intn(len(nextStates))], nil
	case v2.Encoded:
		// From Encoded, can go to GatheringSignatures or Failed
		nextStates := []v2.BlobStatus{v2.GatheringSignatures, v2.Failed}
		return nextStates[rand.Intn(len(nextStates))], nil
	case v2.GatheringSignatures:
		// From GatheringSignatures, can go to Complete or Failed
		nextStates := []v2.BlobStatus{v2.Complete, v2.Failed}
		return nextStates[rand.Intn(len(nextStates))], nil
	case v2.Complete, v2.Failed:
		// Terminal states - return the same status (idempotent)
		return currentStatus, nil
	default:
		return 0, fmt.Errorf("unknown blob status: %v", currentStatus)
	}
}

func (br *BenchmarkRunner) generateTestBlobMetadataWithRandomStatus(blobKey corev2.BlobKey) *v2.BlobMetadata {
	metadata := br.generateTestBlobMetadata(blobKey)
	
	// Assign a random status with weights heavily favoring earlier states
	// This ensures we have many blobs that can still transition for better benchmarking
	statusWeights := []struct {
		status v2.BlobStatus
		weight int
	}{
		{v2.Queued, 50},              // 50% - can transition to Encoded or Failed
		{v2.Encoded, 35},             // 35% - can transition to GatheringSignatures or Failed
		{v2.GatheringSignatures, 13}, // 13% - can transition to Complete or Failed
		{v2.Complete, 1},             // 1% - terminal state
		{v2.Failed, 1},               // 1% - terminal state
	}
	
	totalWeight := 0
	for _, sw := range statusWeights {
		totalWeight += sw.weight
	}
	
	randomWeight := rand.Intn(totalWeight)
	currentWeight := 0
	for _, sw := range statusWeights {
		currentWeight += sw.weight
		if randomWeight < currentWeight {
			metadata.BlobStatus = sw.status
			break
		}
	}
	
	return metadata
}

func (br *BenchmarkRunner) generateTestBlobMetadata(blobKey corev2.BlobKey) *v2.BlobMetadata {
	// Generate a test blob header

	// Create a random G1 commitment
	commitment := &encoding.G1Commitment{}
	// Generate random values for X and Y using big.Int
	x := new(big.Int)
	y := new(big.Int)
	x.SetInt64(rand.Int63())
	y.SetInt64(rand.Int63())
	commitment.X.SetBigInt(x)
	commitment.Y.SetBigInt(y)

	// Create a random G2 commitment
	lengthCommitment := &encoding.G2Commitment{}
	lengthProof := &encoding.G2Commitment{}
	// Generate random values for X and Y using big.Int
	x = new(big.Int)
	y = new(big.Int)
	x.SetInt64(rand.Int63())
	y.SetInt64(rand.Int63())
	lengthCommitment.X.A0.SetBigInt(x)
	lengthCommitment.X.A1.SetBigInt(y)
	lengthCommitment.Y.A0.SetBigInt(x)
	lengthCommitment.Y.A1.SetBigInt(y)
	lengthProof.X.A0.SetBigInt(x)
	lengthProof.X.A1.SetBigInt(y)
	lengthProof.Y.A0.SetBigInt(x)
	lengthProof.Y.A1.SetBigInt(y)

	// Generate random payment amount
	cumulativePayment := new(big.Int)
	cumulativePayment.SetInt64(rand.Int63n(1000))

	header := &corev2.BlobHeader{
		BlobVersion:   0,
		QuorumNumbers: []core.QuorumID{0, 1}, // Add some quorum numbers
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: lengthCommitment,
			LengthProof:      lengthProof,
		},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x1234567890123456789012345678901234567890"),
			CumulativePayment: cumulativePayment,
		},
	}

	signature := make([]byte, 65)
	rand.Read(signature)

	return &v2.BlobMetadata{
		BlobHeader:  header,
		BlobStatus:  v2.Queued,
		Expiry:      uint64(time.Now().Add(24 * time.Hour).Unix()),
		NumRetries:  0,
		RequestedAt: uint64(time.Now().UnixNano()),
		UpdatedAt:   uint64(time.Now().UnixNano()),
		BlobSize:    1024,
		Signature:   signature,
	}
}

// Reporting functions
func (br *BenchmarkRunner) printIntermediateResults() {
	fmt.Println("\n=== Intermediate Results ===")
	fmt.Printf("Elapsed Time: %s\n", time.Since(br.collector.startTime))
	fmt.Println()

	for _, opConfig := range br.config.Operations {
		metrics := br.collector.GetMetrics(opConfig.Type)
		if metrics.TotalOps > 0 {
			br.printOperationMetrics(metrics)
		}
	}
}

func (br *BenchmarkRunner) printFinalResults() {
	fmt.Println("\n=== Final Benchmark Results ===")
	fmt.Printf("Store Type: %s\n", br.config.StoreType)
	fmt.Printf("Total Duration: %s\n", time.Since(br.collector.startTime))
	fmt.Println()

	for _, opConfig := range br.config.Operations {
		metrics := br.collector.GetMetrics(opConfig.Type)
		if metrics.TotalOps > 0 {
			fmt.Printf("--- %s ---\n", opConfig.Type)
			fmt.Printf("Target Rate: %d ops/sec\n", opConfig.RatePerSec)
			br.printOperationMetrics(metrics)
		}
	}
}

func (br *BenchmarkRunner) printOperationMetrics(m OperationMetrics) {
	fmt.Printf("Total Operations: %d (Success: %d, Failed: %d)\n", m.TotalOps, m.SuccessfulOps, m.FailedOps)
	fmt.Printf("Actual Throughput: %.2f ops/sec\n", m.Throughput)
	fmt.Printf("Error Rate: %.2f%%\n", m.ErrorRate)

	if m.SuccessfulOps > 0 {
		fmt.Printf("Latency (ms):\n")
		fmt.Printf("  Min: %.2f\n", m.MinLatencyMs)
		fmt.Printf("  Avg: %.2f\n", m.AvgLatencyMs)
		fmt.Printf("  P50: %.2f\n", m.P50LatencyMs)
		fmt.Printf("  P90: %.2f\n", m.P90LatencyMs)
		fmt.Printf("  P95: %.2f\n", m.P95LatencyMs)
		fmt.Printf("  P99: %.2f\n", m.P99LatencyMs)
		fmt.Printf("  Max: %.2f\n", m.MaxLatencyMs)
	}
	fmt.Println()
}

// runOperations runs all configured operations (used for warmup)
func (br *BenchmarkRunner) runOperations(ctx context.Context, collectMetrics bool) {
	var wg sync.WaitGroup
	for _, opConfig := range br.config.Operations {
		wg.Add(1)
		go func(cfg OperationConfig) {
			defer wg.Done()
			// For warmup, use a shorter duration
			warmupConfig := cfg
			warmupConfig.Duration = br.config.WarmupTime
			br.runOperation(ctx, warmupConfig, collectMetrics)
		}(opConfig)
	}
	wg.Wait()
}

// Helper functions
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	index := int(math.Ceil(p*float64(len(values)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
