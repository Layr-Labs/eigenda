package mock

import (
	"context"
	"strings"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
)

type S3Client struct {
	bucket map[string][]byte
}

var _ s3.Client = (*S3Client)(nil)

func NewS3Client() *S3Client {
	return &S3Client{bucket: make(map[string][]byte)}
}

func (s *S3Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	data, ok := s.bucket[key]
	if !ok {
		return []byte{}, s3.ErrObjectNotFound
	}
	return data, nil
}

func (s *S3Client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	data, ok := s.bucket[key]
	if !ok {
		return nil, s3.ErrObjectNotFound
	}
	size := int64(len(data))
	return &size, nil
}

func (s *S3Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	s.bucket[key] = data
	return nil
}

func (s *S3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	delete(s.bucket, key)
	return nil
}

func (s *S3Client) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3.Object, error) {
	objects := make([]s3.Object, 0, 5)
	for k, v := range s.bucket {
		if strings.HasPrefix(k, prefix) {
			objects = append(objects, s3.Object{Key: k, Size: int64(len(v))})
		}
	}
	return objects, nil
}

func (s *S3Client) CreateBucket(ctx context.Context, bucket string) error {
	return nil
}

func (s *S3Client) FragmentedUploadObject(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	fragmentSize int) error {
	s.bucket[key] = data
	return nil
}

func (s *S3Client) FragmentedDownloadObject(
	ctx context.Context,
	bucket string,
	key string,
	fileSize int,
	fragmentSize int) ([]byte, error) {
	data, ok := s.bucket[key]
	if !ok {
		return []byte{}, s3.ErrObjectNotFound
	}
	return data, nil
}
