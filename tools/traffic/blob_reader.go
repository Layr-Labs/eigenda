package traffic

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	contractEigenDAServiceManager "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	gcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
	"time"
)

// TODO for all of these new types, decide if variables need to be pointers or not

// BlobReader reads blobs from a disperser at a configured rate.
type BlobReader struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// TODO use code from this class
	retriever   clients.RetrievalClient
	chainClient eth.ChainClient

	// table of blobs to read from.
	table *BlobTable

	metrics *Metrics

	fetchBatchHeaderMetric  LatencyMetric
	fetchBatchHeaderSuccess CountMetric
	fetchBatchHeaderFailure CountMetric
	readLatencyMetric       LatencyMetric
	readSuccessMetric       CountMetric
	readFailureMetric       CountMetric
	recombinationSuccess    CountMetric
	recombinationFailure    CountMetric
	operatorSuccessMetrics  map[core.OperatorID]CountMetric
	operatorFailureMetrics  map[core.OperatorID]CountMetric
	candidatePoolSize       GaugeMetric
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	retriever clients.RetrievalClient,
	chainClient eth.ChainClient,
	table *BlobTable,
	metrics *Metrics) BlobReader {

	return BlobReader{
		ctx:                     ctx,
		waitGroup:               waitGroup,
		retriever:               retriever,
		chainClient:             chainClient,
		table:                   table,
		metrics:                 metrics,
		fetchBatchHeaderMetric:  metrics.NewLatencyMetric("fetch_batch_header"),
		fetchBatchHeaderSuccess: metrics.NewCountMetric("fetch_batch_header_success"),
		fetchBatchHeaderFailure: metrics.NewCountMetric("fetch_batch_header_failure"),
		recombinationSuccess:    metrics.NewCountMetric("recombination_success"),
		recombinationFailure:    metrics.NewCountMetric("recombination_failure"),
		readLatencyMetric:       metrics.NewLatencyMetric("read"),
		readSuccessMetric:       metrics.NewCountMetric("read_success"),
		readFailureMetric:       metrics.NewCountMetric("read_failure"),
		operatorSuccessMetrics:  make(map[core.OperatorID]CountMetric),
		operatorFailureMetrics:  make(map[core.OperatorID]CountMetric),
		candidatePoolSize:       metrics.NewGaugeMetric("candidate_pool_size"),
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
	ticker := time.NewTicker(time.Second) // TODO setting
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

	reader.candidatePoolSize.Set(float64(reader.table.Size()))

	metadata := reader.table.GetRandom(true)
	if metadata == nil {
		// There are no blobs to read, do nothing.
		return
	}

	// TODO add timeout config
	ctxTimeout, cancel := context.WithTimeout(*reader.ctx, time.Second*5)
	batchHeader, err := InvokeAndReportLatency(&reader.fetchBatchHeaderMetric,
		func() (*contractEigenDAServiceManager.IEigenDAServiceManagerBatchHeader, error) {
			return reader.chainClient.FetchBatchHeader(
				ctxTimeout,
				gcommon.HexToAddress("0x851356ae760d987E095750cCeb3bC6014560891C"),
				*metadata.batchHeaderHash,
				big.NewInt(int64(0)),
				nil)
		})
	cancel()
	if err != nil {
		// TODO log
		reader.fetchBatchHeaderFailure.Increment()
		return
	}
	reader.fetchBatchHeaderSuccess.Increment()

	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], *metadata.batchHeaderHash)

	// TODO add timeout config
	ctxTimeout, cancel = context.WithTimeout(*reader.ctx, time.Second*5)
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
		// TODO log
		reader.readFailureMetric.Increment()
		return
	}
	reader.readSuccessMetric.Increment()

	chunkCount := chunks.AssignmentInfo.TotalChunks

	var assignments map[core.OperatorID]core.Assignment
	assignments = chunks.Assignments

	data, err := reader.retriever.CombineChunks(chunks)

	if err != nil {
		fmt.Println("Error combining chunks:", err) // TODO
		reader.recombinationFailure.Increment()
		return
	}
	reader.recombinationSuccess.Increment()

	// TODO verify blob data

	fmt.Printf("=====================================\nRead blob. Total chunk count = %d\nRetrieved chunk count = %d\nData length = %d\n",
		chunkCount, len(chunks.Chunks), len(data)) // TODO

	indexSet := make(map[encoding.ChunkNumber]bool)
	for index := range chunks.Indices {
		indexSet[chunks.Indices[index]] = true
	}

	for id, assignment := range assignments {
		fmt.Printf("  - Operator ID: %d, Start Index: %d, Num Chunks: %d\n", id, assignment.StartIndex, assignment.NumChunks)
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
