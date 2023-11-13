package s3

import "context"

type Client interface {
	DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error)
	UploadObject(ctx context.Context, bucket string, key string, data []byte) error
	DeleteObject(ctx context.Context, bucket string, key string) error
	ListObjects(ctx context.Context, bucket string, prefix string) ([]Object, error)
}
