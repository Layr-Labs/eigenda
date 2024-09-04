package workers

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	contractEigenDAServiceManager "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	config2 "github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
)

// BlobReader reads blobs from a disperser at a configured rate.
type BlobReader struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// ticker is used to control the rate at which blobs are read.
	ticker InterceptableTicker

	// config contains the configuration for the generator.
	config *config2.WorkerConfig

	retriever   clients.RetrievalClient
	chainClient eth.ChainClient

	// blobsToRead blobs we are required to read a certain number of times.
	blobsToRead *table.BlobStore

	generatorMetrics           metrics.Metrics
	fetchBatchHeaderMetric     metrics.LatencyMetric
	fetchBatchHeaderSuccess    metrics.CountMetric
	fetchBatchHeaderFailure    metrics.CountMetric
	readLatencyMetric          metrics.LatencyMetric
	readSuccessMetric          metrics.CountMetric
	readFailureMetric          metrics.CountMetric
	recombinationSuccessMetric metrics.CountMetric
	recombinationFailureMetric metrics.CountMetric
	validBlobMetric            metrics.CountMetric
	invalidBlobMetric          metrics.CountMetric
	operatorSuccessMetrics     map[core.OperatorID]metrics.CountMetric
	operatorFailureMetrics     map[core.OperatorID]metrics.CountMetric
	requiredReadPoolSizeMetric metrics.GaugeMetric
	optionalReadPoolSizeMetric metrics.GaugeMetric
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	logger logging.Logger,
	ticker InterceptableTicker,
	config *config2.WorkerConfig,
	retriever clients.RetrievalClient,
	chainClient eth.ChainClient,
	blobStore *table.BlobStore,
	generatorMetrics metrics.Metrics) BlobReader {

	return BlobReader{
		ctx:                        ctx,
		waitGroup:                  waitGroup,
		logger:                     logger,
		ticker:                     ticker,
		config:                     config,
		retriever:                  retriever,
		chainClient:                chainClient,
		blobsToRead:                blobStore,
		generatorMetrics:           generatorMetrics,
		fetchBatchHeaderMetric:     generatorMetrics.NewLatencyMetric("fetch_batch_header"),
		fetchBatchHeaderSuccess:    generatorMetrics.NewCountMetric("fetch_batch_header_success"),
		fetchBatchHeaderFailure:    generatorMetrics.NewCountMetric("fetch_batch_header_failure"),
		recombinationSuccessMetric: generatorMetrics.NewCountMetric("recombination_success"),
		recombinationFailureMetric: generatorMetrics.NewCountMetric("recombination_failure"),
		readLatencyMetric:          generatorMetrics.NewLatencyMetric("read"),
		validBlobMetric:            generatorMetrics.NewCountMetric("valid_blob"),
		invalidBlobMetric:          generatorMetrics.NewCountMetric("invalid_blob"),
		readSuccessMetric:          generatorMetrics.NewCountMetric("read_success"),
		readFailureMetric:          generatorMetrics.NewCountMetric("read_failure"),
		operatorSuccessMetrics:     make(map[core.OperatorID]metrics.CountMetric),
		operatorFailureMetrics:     make(map[core.OperatorID]metrics.CountMetric),
		requiredReadPoolSizeMetric: generatorMetrics.NewGaugeMetric("required_read_pool_size"),
		optionalReadPoolSizeMetric: generatorMetrics.NewGaugeMetric("optional_read_pool_size"),
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
	ticker := reader.ticker.GetTimeChannel()
	for {
		select {
		case <-(*reader.ctx).Done():
			return
		case <-ticker:
			reader.randomRead()
		}
	}
}

// randomRead reads a random blob.
func (reader *BlobReader) randomRead() {
	metadata := reader.blobsToRead.GetNext()
	if metadata == nil {
		// There are no blobs that we are required to read.
		return
	}

	reader.requiredReadPoolSizeMetric.Set(float64(reader.blobsToRead.Size()))

	ctxTimeout, cancel := context.WithTimeout(*reader.ctx, reader.config.FetchBatchHeaderTimeout)
	batchHeader, err := metrics.InvokeAndReportLatency(reader.fetchBatchHeaderMetric,
		func() (*contractEigenDAServiceManager.IEigenDAServiceManagerBatchHeader, error) {
			return reader.chainClient.FetchBatchHeader(
				ctxTimeout,
				gcommon.HexToAddress(reader.config.EigenDAServiceManager),
				metadata.BatchHeaderHash[:],
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
	copy(batchHeaderHash[:], metadata.BatchHeaderHash[:])

	ctxTimeout, cancel = context.WithTimeout(*reader.ctx, reader.config.RetrieveBlobChunksTimeout)
	chunks, err := metrics.InvokeAndReportLatency(reader.readLatencyMetric, func() (*clients.BlobChunks, error) {
		return reader.retriever.RetrieveBlobChunks(
			ctxTimeout,
			batchHeaderHash,
			uint32(metadata.BlobIndex),
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

	assignments := chunks.Assignments

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
		metric = reader.generatorMetrics.NewCountMetric(fmt.Sprintf("operator_%x_returned_chunk", operatorId))
		reader.operatorSuccessMetrics[operatorId] = metric
	}

	metric.Increment()
}

// reportMissingChunk reports a missing chunk.
func (reader *BlobReader) reportMissingChunk(operatorId core.OperatorID) {
	metric, exists := reader.operatorFailureMetrics[operatorId]
	if !exists {
		metric = reader.generatorMetrics.NewCountMetric(fmt.Sprintf("operator_%x_witheld_chunk", operatorId))
		reader.operatorFailureMetrics[operatorId] = metric
	}

	metric.Increment()
}

// verifyBlob performs sanity checks on the blob.
func (reader *BlobReader) verifyBlob(metadata *table.BlobMetadata, blob *[]byte) {
	// Trim off the padding.
	truncatedBlob := (*blob)[:metadata.Size]
	recomputedChecksum := md5.Sum(truncatedBlob)

	if metadata.Checksum == recomputedChecksum {
		reader.validBlobMetric.Increment()
	} else {
		reader.invalidBlobMetric.Increment()
	}
}
