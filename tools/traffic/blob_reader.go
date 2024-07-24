package traffic

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
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
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	retriever clients.RetrievalClient,
	chainClient eth.ChainClient,
	table *BlobTable) BlobReader {

	return BlobReader{
		ctx:         ctx,
		waitGroup:   waitGroup,
		retriever:   retriever,
		chainClient: chainClient,
		table:       table,
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

	metadata := reader.table.GetRandom(true)
	if metadata == nil {
		// There are no blobs to read, do nothing.
		return
	}

	// TODO use NodeClient

	batchHeader, err := reader.chainClient.FetchBatchHeader(
		*reader.ctx,
		gcommon.HexToAddress("0x851356ae760d987E095750cCeb3bC6014560891C"),
		*metadata.batchHeaderHash,
		big.NewInt(int64(0)),
		nil)
	if err != nil {
		fmt.Println("Error fetching batch header:", err) // TODO
		return
	}

	// TODO is this correct?
	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], *metadata.batchHeaderHash)

	fmt.Printf("attempting to read blob, header hash: %x, index: %d\n, batch header: %+v\n", *metadata.batchHeaderHash, metadata.blobIndex, batchHeader)
	chunks, err := reader.retriever.RetrieveBlobChunks(
		*reader.ctx,
		batchHeaderHash,
		metadata.blobIndex,
		uint(batchHeader.ReferenceBlockNumber),
		batchHeader.BlobHeadersRoot,
		core.QuorumID(0))

	if err != nil {
		fmt.Println("Error reading blob:", err) // TODO
		return
	}

	chunkCount := chunks.AssignmentInfo.TotalChunks

	var assignments map[core.OperatorID]core.Assignment
	assignments = chunks.Assignments

	data, err := reader.retriever.CombineChunks(chunks)

	if err != nil {
		fmt.Println("Error combining chunks:", err) // TODO
		return
	}

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
			if !indexSet[index] {
				fmt.Printf("    - Missing chunk: %d\n", index)
			} else {
				fmt.Printf("    - Chunk found: %d\n", index)
			}
		}
	}

}
