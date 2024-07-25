package traffic

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	contractEigenDAServiceManager "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
	"time"
)

// BlobReader reads blobs from a disperser at a configured rate.
type BlobReader struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// config contains the configuration for the generator.
	config *Config

	retriever   clients.RetrievalClient
	chainClient eth.ChainClient

	// requiredReads blobs we are required to read a certain number of times.
	requiredReads *BlobTable

	// optionalReads blobs we are not required to read, but can choose to read if we want.
	optionalReads *BlobTable

	metrics                    *Metrics
	fetchBatchHeaderMetric     LatencyMetric
	fetchBatchHeaderSuccess    CountMetric
	fetchBatchHeaderFailure    CountMetric
	readLatencyMetric          LatencyMetric
	readSuccessMetric          CountMetric
	readFailureMetric          CountMetric
	recombinationSuccessMetric CountMetric
	recombinationFailureMetric CountMetric
	validBlobMetric            CountMetric
	invalidBlobMetric          CountMetric
	operatorSuccessMetrics     map[core.OperatorID]CountMetric
	operatorFailureMetrics     map[core.OperatorID]CountMetric
	requiredReadPoolSizeMetric GaugeMetric
	optionalReadPoolSizeMetric GaugeMetric
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	logger logging.Logger,
	config *Config,
	retriever clients.RetrievalClient,
	chainClient eth.ChainClient,
	table *BlobTable,
	metrics *Metrics) BlobReader {

	optionalReads := NewBlobTable()

	return BlobReader{
		ctx:                        ctx,
		waitGroup:                  waitGroup,
		logger:                     logger,
		config:                     config,
		retriever:                  retriever,
		chainClient:                chainClient,
		requiredReads:              table,
		optionalReads:              &optionalReads,
		metrics:                    metrics,
		fetchBatchHeaderMetric:     metrics.NewLatencyMetric("fetch_batch_header"),
		fetchBatchHeaderSuccess:    metrics.NewCountMetric("fetch_batch_header_success"),
		fetchBatchHeaderFailure:    metrics.NewCountMetric("fetch_batch_header_failure"),
		recombinationSuccessMetric: metrics.NewCountMetric("recombination_success"),
		recombinationFailureMetric: metrics.NewCountMetric("recombination_failure"),
		readLatencyMetric:          metrics.NewLatencyMetric("read"),
		validBlobMetric:            metrics.NewCountMetric("valid_blob"),
		invalidBlobMetric:          metrics.NewCountMetric("invalid_blob"),
		readSuccessMetric:          metrics.NewCountMetric("read_success"),
		readFailureMetric:          metrics.NewCountMetric("read_failure"),
		operatorSuccessMetrics:     make(map[core.OperatorID]CountMetric),
		operatorFailureMetrics:     make(map[core.OperatorID]CountMetric),
		requiredReadPoolSizeMetric: metrics.NewGaugeMetric("required_read_pool_size"),
		optionalReadPoolSizeMetric: metrics.NewGaugeMetric("optional_read_pool_size"),
	}
}

// Start begins a blob reader goroutine.
func (reader *BlobReader) Start() {
	reader.waitGroup.Add(1)
	go func() {
		defer reader.waitGroup.Done()
		reader.run()
	}()
}

// run periodically performs reads on blobs.
func (reader *BlobReader) run() {
	ticker := time.NewTicker(reader.config.ReadRequestInterval)
	for {
		select {
		case <-(*reader.ctx).Done():
			return
		case <-ticker.C:
			reader.randomRead()
		}
	}
}

// randomRead reads a random blob.
func (reader *BlobReader) randomRead() {

	reader.requiredReadPoolSizeMetric.Set(float64(reader.requiredReads.Size()))

	metadata, removed := reader.requiredReads.GetRandom(true)
	if metadata == nil {
		// There are no blobs that we are required to read. Get a random blob from the optionalReads.
		metadata, _ = reader.optionalReads.GetRandom(false)
		if metadata == nil {
			// No blobs to read.
			return
		}
	} else if removed {
		// We have removed a blob from the requiredReads. Add it to the optionalReads.
		reader.optionalReads.AddOrReplace(metadata, reader.config.ReadOverflowTableSize)
		reader.optionalReadPoolSizeMetric.Set(float64(reader.optionalReads.Size()))
	}

	ctxTimeout, cancel := context.WithTimeout(*reader.ctx, reader.config.FetchBatchHeaderTimeout)
	batchHeader, err := InvokeAndReportLatency(&reader.fetchBatchHeaderMetric,
		func() (*contractEigenDAServiceManager.IEigenDAServiceManagerBatchHeader, error) {
			return reader.chainClient.FetchBatchHeader(
				ctxTimeout,
				gcommon.HexToAddress(reader.config.EigenDAServiceManager),
				*metadata.batchHeaderHash,
				big.NewInt(int64(0)),
				nil)
		})
	cancel()
	if err != nil {
		reader.logger.Error("failed to get batch header", "err:", err)
		reader.fetchBatchHeaderFailure.Increment()
		return
	}
	reader.fetchBatchHeaderSuccess.Increment()

	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], *metadata.batchHeaderHash)

	ctxTimeout, cancel = context.WithTimeout(*reader.ctx, reader.config.RetrieveBlobChunksTimeout)
	chunks, err := InvokeAndReportLatency(&reader.readLatencyMetric, func() (*clients.BlobChunks, error) {
		return reader.retriever.RetrieveBlobChunks(
			ctxTimeout,
			batchHeaderHash,
			metadata.blobIndex,
			uint(batchHeader.ReferenceBlockNumber),
			batchHeader.BlobHeadersRoot,
			core.QuorumID(0))
	})
	cancel()
	if err != nil {
		reader.logger.Error("failed to read chunks", "err:", err)
		reader.readFailureMetric.Increment()
		return
	}
	reader.readSuccessMetric.Increment()

	var assignments map[core.OperatorID]core.Assignment
	assignments = chunks.Assignments

	data, err := reader.retriever.CombineChunks(chunks)
	if err != nil {
		reader.logger.Error("failed to combine chunks", "err:", err)
		reader.recombinationFailureMetric.Increment()
		return
	}
	reader.recombinationSuccessMetric.Increment()

	reader.verifyBlob(metadata, &data)

	indexSet := make(map[encoding.ChunkNumber]bool)
	for index := range chunks.Indices {
		indexSet[chunks.Indices[index]] = true
	}

	for id, assignment := range assignments {
		for index := assignment.StartIndex; index < assignment.StartIndex+assignment.NumChunks; index++ {
			if indexSet[index] {
				reader.reportChunk(id)
			} else {
				reader.reportMissingChunk(id)
			}
		}
	}
}

// reportChunk reports a successful chunk read.
func (reader *BlobReader) reportChunk(operatorId core.OperatorID) {
	metric, exists := reader.operatorSuccessMetrics[operatorId]
	if !exists {
		metric = reader.metrics.NewCountMetric(fmt.Sprintf("operator_%x_returned_chunk", operatorId))
		reader.operatorSuccessMetrics[operatorId] = metric
	}

	metric.Increment()
}

// reportMissingChunk reports a missing chunk.
func (reader *BlobReader) reportMissingChunk(operatorId core.OperatorID) {
	metric, exists := reader.operatorFailureMetrics[operatorId]
	if !exists {
		metric = reader.metrics.NewCountMetric(fmt.Sprintf("operator_%x_witheld_chunk", operatorId))
		reader.operatorFailureMetrics[operatorId] = metric
	}

	metric.Increment()
}

// verifyBlob performs sanity checks on the blob.
func (reader *BlobReader) verifyBlob(metadata *BlobMetadata, blob *[]byte) {
	// Trim off the padding.
	truncatedBlob := (*blob)[:metadata.size]
	recomputedChecksum := md5.Sum(truncatedBlob)

	if *metadata.checksum == recomputedChecksum {
		reader.validBlobMetric.Increment()
	} else {
		reader.invalidBlobMetric.Increment()
	}
}
