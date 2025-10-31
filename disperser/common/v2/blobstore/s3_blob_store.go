package blobstore

import (
	"context"
	"fmt"

	s3common "github.com/Layr-Labs/eigenda/common/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/pkg/errors"
)

type BlobStore struct {
	bucketName string
	s3Client   s3common.S3Client
	logger     logging.Logger
}

func NewBlobStore(s3BucketName string, s3Client s3common.S3Client, logger logging.Logger) *BlobStore {
	return &BlobStore{
		bucketName: s3BucketName,
		s3Client:   s3Client,
		logger:     logger,
	}
}

// StoreBlob adds a blob to the blob store
func (b *BlobStore) StoreBlob(ctx context.Context, key corev2.BlobKey, data []byte) error {
	_, err := b.s3Client.HeadObject(ctx, b.bucketName, s3common.ScopedBlobKey(key))
	if err == nil {
		b.logger.Warnf("blob already exists in bucket %s: %s", b.bucketName, key)
		return ErrAlreadyExists
	}

	err = b.s3Client.UploadObject(ctx, b.bucketName, s3common.ScopedBlobKey(key), data)
	if err != nil {
		b.logger.Errorf("failed to upload blob in bucket %s: %w", b.bucketName, err)
		return err
	}
	return nil
}

// GetBlob retrieves a blob from the blob store
func (b *BlobStore) GetBlob(ctx context.Context, key corev2.BlobKey) ([]byte, error) {
	data, err := b.s3Client.DownloadObject(ctx, b.bucketName, s3common.ScopedBlobKey(key))
	if errors.Is(err, s3common.ErrObjectNotFound) {
		b.logger.Warnf("blob not found in bucket %s: %s", b.bucketName, key)
		return nil, ErrBlobNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%s, bucket: %s,: %w", ErrBlobNotFound.Error(), b.bucketName, err)
	}
	return data, nil
}
