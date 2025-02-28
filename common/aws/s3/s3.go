package s3

import "context"

// Client encapsulates the functionality of an S3 client.
type Client interface {

	// DownloadObject downloads an object from S3.
	DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error)

	// HeadObject retrieves the size of an object in S3. Returns error if the object does not exist.
	HeadObject(ctx context.Context, bucket string, key string) (*int64, error)

	// UploadObject uploads an object to S3.
	UploadObject(ctx context.Context, bucket string, key string, data []byte) error

	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, bucket string, key string) error

	// ListObjects lists all objects in a bucket with the given prefix. Note that this method may return
	// file fragments if the bucket contains files uploaded via FragmentedUploadObject.
	ListObjects(ctx context.Context, bucket string, prefix string) ([]Object, error)

	// CreateBucket creates a bucket in S3.
	CreateBucket(ctx context.Context, bucket string) error

	// FragmentedUploadObject uploads a file to S3. The fragmentSize parameter specifies the maximum size of each
	// file uploaded to S3. If the file is larger than fragmentSize then it will be broken into
	// smaller parts and uploaded in parallel. The file will be reassembled on download.
	//
	// Note: if a file is uploaded with this method, only the FragmentedDownloadObject method should be used to
	// download the file. It is not advised to use DeleteObject on files uploaded with this method (if such
	// functionality is required, a new method to do so should be added to this interface).
	//
	// Note: if this operation fails partway through, some file fragments may have made it to S3 and others may not.
	// In order to prevent long term accumulation of fragments, it is suggested to use this method in conjunction with
	// a bucket configured to have a TTL.
	FragmentedUploadObject(
		ctx context.Context,
		bucket string,
		key string,
		data []byte,
		fragmentSize int) error

	// FragmentedDownloadObject downloads a file from S3, as written by Upload. The fileSize (in bytes) and fragmentSize
	// must be the same as the values used in the FragmentedUploadObject call.
	//
	// Note: this method can only be used to download files that were uploaded with the FragmentedUploadObject method.
	FragmentedDownloadObject(
		ctx context.Context,
		bucket string,
		key string,
		fileSize int,
		fragmentSize int) ([]byte, error)
}
