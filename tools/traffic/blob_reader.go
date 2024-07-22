package traffic

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BlobReader reads blobs from a disperser at a configured rate.
type BlobReader struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// TODO this type should be refactored maybe
	generator *TrafficGenerator

	// table of blobs to read from.
	table *BlobTable
}

// NewBlobReader creates a new BlobReader instance.
func NewBlobReader(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	generator *TrafficGenerator,
	table *BlobTable) BlobReader {

	return BlobReader{
		ctx:       ctx,
		waitGroup: waitGroup,
		generator: generator,
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
	data, err := reader.generator.EigenDAClient.GetBlob(*reader.ctx, *metadata.batchHeaderHash, metadata.blobIndex)
	if err != nil {
		fmt.Println("Error reading blob:", err) // TODO
		return
	}

	// TODO verify blob data

	fmt.Println("Read blob:", data) // TODO
}
