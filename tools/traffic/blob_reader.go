package traffic

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
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
	retriever clients.RetrievalClient

	// table of blobs to read from.
	table *BlobTable
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	retriever clients.RetrievalClient,
	table *BlobTable) BlobReader {

	return BlobReader{
		ctx:       ctx,
		waitGroup: waitGroup,
		retriever: retriever,
		table:     table,
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

	// TODO
	//data, err := reader.generator.DisperserClient.RetrieveBlob(*reader.ctx, *metadata.batchHeaderHash, metadata.blobIndex)
	//if err != nil {
	//	fmt.Println("Error reading blob:", err) // TODO
	//	return
	//}

	// TODO use NodeClient

	fmt.Println("attempting to read blob")
	// TODO arguments probably not right
	data, err := reader.retriever.RetrieveBlob(
		*reader.ctx,
		[32]byte(*metadata.batchHeaderHash),
		metadata.blobIndex,
		0,
		[32]byte{},
		0)

	if err != nil {
		fmt.Println("Error reading blob:", err) // TODO
		return
	}

	// TODO verify blob data

	fmt.Println("Read blob:", data) // TODO
}
