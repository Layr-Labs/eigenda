package chunkstore

import (
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"time"
)

// ChunkMetadataStore is an interface for storing and retrieving metadata about chunks.
type ChunkMetadataStore interface {
	// WriteChunkMetadata stores metadata about a chunk. This metadata is required to read the chunk via ChunkReader.
	WriteChunkMetadata(key *disperser.BlobKey, size int, fragmentSize int) error
	// ReadChunkMetadata retrieves metadata about a chunk. This metadata is required to read the chunk via ChunkReader.
	ReadChunkMetadata(key *disperser.BlobKey) (size int, fragmentSize int, err error)
}

var _ ChunkMetadataStore = (*chunkMetadataStore)(nil)

// chunkMetadataStore is currently just a placeholder. A real implementation will follow in a future PR.
type chunkMetadataStore struct {
}

// NewChunkMetadataStore creates a new ChunkMetadataStore.
func NewChunkMetadataStore(
	logger logging.Logger,
	dynamoDBClient *commondynamodb.Client,
	tableName string,
	ttl time.Duration) ChunkMetadataStore {

	return &chunkMetadataStore{}
}

func (c *chunkMetadataStore) WriteChunkMetadata(key *disperser.BlobKey, size int, fragmentSize int) error {
	// TODO this is an intentional place holder
	return nil
}

func (c *chunkMetadataStore) ReadChunkMetadata(key *disperser.BlobKey) (size int, fragmentSize int, err error) {
	//TODO this is an intentional place holder
	return 0, 0, nil
}
