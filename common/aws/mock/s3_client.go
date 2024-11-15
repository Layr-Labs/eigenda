package mock

import (
	"context"
	"errors"
	"strings"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
)

type S3Client struct {
	bucket map[string][]byte
	Called map[string]int
}

var _ s3.Client = (*S3Client)(nil)

func NewS3Client() *S3Client {
	return &S3Client{
		bucket: make(map[string][]byte),
		Called: map[string]int{
			"DownloadObject":           0,
			"HeadObject":               0,
			"UploadObject":             0,
			"DeleteObject":             0,
			"ListObjects":              0,
			"CreateBucket":             0,
			"FragmentedUploadObject":   0,
			"FragmentedDownloadObject": 0,
		},
	}
}

func (s *S3Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	s.Called["DownloadObject"]++
	data, ok := s.bucket[key]
	if !ok {
		return []byte{}, s3.ErrObjectNotFound
	}
	return data, nil
}

func (s *S3Client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	s.Called["HeadObject"]++
	data, ok := s.bucket[key]
	if !ok {
		return nil, s3.ErrObjectNotFound
	}
	size := int64(len(data))
	return &size, nil
}

func (s *S3Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	s.Called["UploadObject"]++
	s.bucket[key] = data
	return nil
}

func (s *S3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	s.Called["DeleteObject"]++
	delete(s.bucket, key)
	return nil
}

func (s *S3Client) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3.Object, error) {
	s.Called["ListObjects"]++
	objects := make([]s3.Object, 0, 1000)
	for k, v := range s.bucket {
		if strings.HasPrefix(k, prefix) {
			objects = append(objects, s3.Object{Key: k, Size: int64(len(v))})
		}
	}
	return objects, nil
}

func (s *S3Client) CreateBucket(ctx context.Context, bucket string) error {
	s.Called["CreateBucket"]++
	return nil
}

func (s *S3Client) FragmentedUploadObject(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	fragmentSize int) error {
	s.Called["FragmentedUploadObject"]++
	fragments, err := s3.BreakIntoFragments(key, data, fragmentSize)
	if err != nil {
		return err
	}
	for _, fragment := range fragments {
		s.bucket[fragment.FragmentKey] = fragment.Data
	}
	return nil
}

func (s *S3Client) FragmentedDownloadObject(
	ctx context.Context,
	bucket string,
	key string,
	fileSize int,
	fragmentSize int) ([]byte, error) {
	s.Called["FragmentedDownloadObject"]++
	if fileSize <= 0 {
		return nil, errors.New("fileSize must be greater than 0")
	}
	if fragmentSize <= 0 {
		return nil, errors.New("fragmentSize must be greater than 0")
	}

	count := 0
	if fileSize < fragmentSize {
		count = 1
	} else if fileSize%fragmentSize == 0 {
		count = fileSize / fragmentSize
	} else {
		count = fileSize/fragmentSize + 1
	}
	fragmentKeys, err := s3.GetFragmentKeys(key, count)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 0, fileSize)
	for _, fragmentKey := range fragmentKeys {
		fragmentData, ok := s.bucket[fragmentKey]
		if !ok {
			return nil, s3.ErrObjectNotFound
		}
		data = append(data, fragmentData...)
	}
	return data, nil
}
