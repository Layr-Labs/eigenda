package workers

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
	"time"
)

// BlobReader reads blobs from the DA network at a configured rate.
type BlobReader struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// config contains the configuration for the generator.
	config *config.WorkerConfig

	retriever   clients.RetrievalClient
	chainClient eth.ChainClient

	// blobsToRead blobs we are required to read a certain number of times.
	blobsToRead *table.BlobStore

	// metrics for the blob reader.
	metrics *blobReaderMetrics
}

type blobReaderMetrics struct {
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
	config *config.WorkerConfig,
	retriever clients.RetrievalClient,
	chainClient eth.ChainClient,
	blobStore *table.BlobStore,
	generatorMetrics metrics.Metrics) BlobReader {

	return BlobReader{
		ctx:         ctx,
		waitGroup:   waitGroup,
		logger:      logger,
		config:      config,
		retriever:   retriever,
		chainClient: chainClient,
		blobsToRead: blobStore,
		metrics: &blobReaderMetrics{
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
		},
	}
}

// Start begins a blob reader goroutine.
func (r *BlobReader) Start() {
	r.waitGroup.Add(1)
	ticker := time.NewTicker(r.config.ReadRequestInterval)
	go func() {
		defer r.waitGroup.Done()
		for {
			select {
			case <-(*r.ctx).Done():
				err := (*r.ctx).Err()
				if err != nil {
					r.logger.Info("blob reader context closed", "err:", err)
				}
				return
			case <-ticker.C:
				r.randomRead()
			}
		}
	}()
}

// randomRead reads a random blob.
func (r *BlobReader) randomRead() {
	metadata := r.blobsToRead.GetNext()
	if metadata == nil {
		// There are no blobs that we are required to read.
		return
	}

	r.metrics.requiredReadPoolSizeMetric.Set(float64(r.blobsToRead.Size()))

	ctxTimeout, cancel := context.WithTimeout(*r.ctx, r.config.FetchBatchHeaderTimeout)
	defer cancel()

	start := time.Now()
	batchHeader, err := r.chainClient.FetchBatchHeader(
		ctxTimeout,
		gcommon.HexToAddress(r.config.EigenDAServiceManager),
		metadata.BatchHeaderHash[:],
		big.NewInt(int64(0)),
		nil)
	if err != nil {
		r.logger.Error("failed to get batch header", "err:", err)
		r.metrics.fetchBatchHeaderFailure.Increment()
		return
	}
	r.metrics.fetchBatchHeaderMetric.ReportLatency(time.Since(start))

	r.metrics.fetchBatchHeaderSuccess.Increment()

	ctxTimeout, cancel = context.WithTimeout(*r.ctx, r.config.RetrieveBlobChunksTimeout)
	defer cancel()

	start = time.Now()
	chunks, err := r.retriever.RetrieveBlobChunks(
		ctxTimeout,
		metadata.BatchHeaderHash,
		uint32(metadata.BlobIndex),
		uint(batchHeader.ReferenceBlockNumber),
		batchHeader.BlobHeadersRoot,
		core.QuorumID(0))
	if err != nil {
		r.logger.Error("failed to read chunks", "err:", err)
		r.metrics.readFailureMetric.Increment()
		return
	}
	r.metrics.readLatencyMetric.ReportLatency(time.Since(start))

	r.metrics.readSuccessMetric.Increment()

	assignments := chunks.Assignments

	data, err := r.retriever.CombineChunks(chunks)
	if err != nil {
		r.logger.Error("failed to combine chunks", "err:", err)
		r.metrics.recombinationFailureMetric.Increment()
		return
	}
	r.metrics.recombinationSuccessMetric.Increment()

	r.verifyBlob(metadata, &data)

	indexSet := make(map[encoding.ChunkNumber]bool)
	for index := range chunks.Indices {
		indexSet[chunks.Indices[index]] = true
	}

	for id, assignment := range assignments {
		for index := assignment.StartIndex; index < assignment.StartIndex+assignment.NumChunks; index++ {
			if indexSet[index] {
				r.reportChunk(id)
			} else {
				r.reportMissingChunk(id)
			}
		}
	}
}

// reportChunk reports a successful chunk read.
func (r *BlobReader) reportChunk(operatorId core.OperatorID) {
	metric, exists := r.metrics.operatorSuccessMetrics[operatorId]
	if !exists {
		metric = r.metrics.generatorMetrics.NewCountMetric(fmt.Sprintf("operator_%x_returned_chunk", operatorId))
		r.metrics.operatorSuccessMetrics[operatorId] = metric
	}

	metric.Increment()
}

// reportMissingChunk reports a missing chunk.
func (r *BlobReader) reportMissingChunk(operatorId core.OperatorID) {
	metric, exists := r.metrics.operatorFailureMetrics[operatorId]
	if !exists {
		metric = r.metrics.generatorMetrics.NewCountMetric(fmt.Sprintf("operator_%x_witheld_chunk", operatorId))
		r.metrics.operatorFailureMetrics[operatorId] = metric
	}

	metric.Increment()
}

// verifyBlob performs sanity checks on the blob.
func (r *BlobReader) verifyBlob(metadata *table.BlobMetadata, blob *[]byte) {
	// Trim off the padding.
	truncatedBlob := (*blob)[:metadata.Size]
	recomputedChecksum := md5.Sum(truncatedBlob)

	if metadata.Checksum == recomputedChecksum {
		r.metrics.validBlobMetric.Increment()
	} else {
		r.metrics.invalidBlobMetric.Increment()
	}
}
