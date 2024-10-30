package chunkstore

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// ChunkReader reads chunks written by ChunkWriter.
type ChunkReader interface {
	// GetChunks retrieves chunks for a given blob key.
	GetChunks(ctx context.Context, blobKey disperser.BlobKey) ([]*encoding.Frame, error)
}

var _ ChunkReader = (*chunkReader)(nil)

type chunkReader struct {
	logger             logging.Logger
	chunkMetadataStore *ChunkMetadataStore
	client             s3.Client
	bucket             string
	shards             []uint32
}

// NewChunkReader creates a new ChunkReader.
func NewChunkReader(
	logger logging.Logger,
	chunkMetadataStore *ChunkMetadataStore,
	s3Client s3.Client,
	bucketName string,
	shards []uint32) ChunkReader {

	return &chunkReader{
		logger:             logger,
		chunkMetadataStore: chunkMetadataStore,
		client:             s3Client,
		bucket:             bucketName,
		shards:             shards,
	}
}

func (r *chunkReader) GetChunks(
	ctx context.Context,
	blobKey disperser.BlobKey) ([]*encoding.Frame, error) {

	s3Key := blobKey.String()

	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3Key)
	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	bundle, err := core.Bundle{}.Deserialize(bytes)
	if err != nil {
		r.logger.Error("Failed to deserialize bundle: %v", err)
		return nil, fmt.Errorf("failed to deserialize bundle: %w", err)
	}

	return bundle, nil
}
