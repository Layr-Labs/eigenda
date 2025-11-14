package s3

import (
	"context"
	"errors"
)

var (
	// ErrObjectNotFound is returned when an object is not found in the storage backend
	ErrObjectNotFound = errors.New("object not found")
)

// S3Client encapsulates the functionality of talking to AWS S3 (or an S3 mimic service).
type S3Client interface {

	// HeadObject retrieves the size of an object in S3. Returns error if the object does not exist.
	HeadObject(ctx context.Context, bucket string, key string) (*int64, error)

	// UploadObject uploads an object to S3.
	UploadObject(ctx context.Context, bucket string, key string, data []byte) error

	// DownloadObject downloads an object from S3. The returned boolean indicates whether the object was found.
	DownloadObject(ctx context.Context, bucket string, key string) ([]byte, bool, error)

	// Download part of the object, specified by startIndex (inclusive) and endIndex (exclusive).
	// The returned boolean indicates whether the object was found.
	DownloadPartialObject(
		ctx context.Context,
		bucket string,
		key string,
		// inclusive
		startIndex int64,
		// exclusive
		endIndex int64,
	) ([]byte, bool, error)

	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, bucket string, key string) error

	// ListObjects lists all objects in a bucket with the given prefix. Note that this method may return
	// file fragments if the bucket contains files uploaded via FragmentedUploadObject.
	ListObjects(ctx context.Context, bucket string, prefix string) ([]ListedObject, error)

	// CreateBucket creates a bucket in S3.
	CreateBucket(ctx context.Context, bucket string) error
}

type ListedObject struct {
	Key  string
	Size int64
}
