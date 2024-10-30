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

// ChunkWriter writes chunks that can be read by ChunkReader.
type ChunkWriter interface {
	// PutChunks stores chunks for a given blob key.
	PutChunks(ctx context.Context, blobKey disperser.BlobKey, chunks []*encoding.Frame) error
}

var _ ChunkWriter = (*chunkWriter)(nil)

type chunkWriter struct {
	logger             logging.Logger
	chunkMetadataStore *ChunkMetadataStore
	s3Client           s3.Client
	bucketName         string
}

func NewChunkWriter(
	logger logging.Logger,
	chunkMetadataStore *ChunkMetadataStore,
	s3Client s3.Client,
	bucketName string) ChunkWriter {

	return &chunkWriter{
		logger:             logger,
		chunkMetadataStore: chunkMetadataStore,
		s3Client:           s3Client,
		bucketName:         bucketName,
	}
}

func (c *chunkWriter) PutChunks(ctx context.Context, blobKey disperser.BlobKey, chunks []*encoding.Frame) error {
	var bundle core.Bundle = chunks
	bytes, err := bundle.Serialize()
	if err != nil {
		c.logger.Error("Failed to serialize bundle: %v", err)
		return fmt.Errorf("failed to serialize bundle: %w", err)
	}

	s3Key := blobKey.String()
	err = c.s3Client.UploadObject(ctx, c.bucketName, s3Key, bytes)
	if err != nil {
		c.logger.Error("Failed to upload chunks to S3: %v", err)
		return fmt.Errorf("failed to upload chunks to S3: %w", err)
	}

	return nil
}
