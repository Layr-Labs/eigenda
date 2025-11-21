package s3

import (
	"context"
	"errors"
	"strings"
)

type MockS3Client struct {
	bucket map[string][]byte
	Called map[string]int
}

var _ S3Client = (*MockS3Client)(nil)

func NewMockS3Client() *MockS3Client {
	return &MockS3Client{
		bucket: make(map[string][]byte),
		Called: map[string]int{
			"DownloadObject":        0,
			"HeadObject":            0,
			"UploadObject":          0,
			"DeleteObject":          0,
			"ListObjects":           0,
			"CreateBucket":          0,
			"DownloadPartialObject": 0,
		},
	}
}

func (s *MockS3Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, bool, error) {
	s.Called["DownloadObject"]++
	data, ok := s.bucket[key]
	if !ok {
		return []byte{}, false, nil
	}
	return data, true, nil
}

func (s *MockS3Client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	s.Called["HeadObject"]++
	data, ok := s.bucket[key]
	if !ok {
		return nil, ErrObjectNotFound
	}
	size := int64(len(data))
	return &size, nil
}

func (s *MockS3Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	s.Called["UploadObject"]++
	s.bucket[key] = data
	return nil
}

func (s *MockS3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	s.Called["DeleteObject"]++
	delete(s.bucket, key)
	return nil
}

func (s *MockS3Client) ListObjects(
	ctx context.Context,
	bucket string,
	prefix string,
) ([]ListedObject, error) {

	s.Called["ListObjects"]++
	objects := make([]ListedObject, 0, 1000)
	for k, v := range s.bucket {
		if strings.HasPrefix(k, prefix) {
			objects = append(objects, ListedObject{Key: k, Size: int64(len(v))})
		}
	}
	return objects, nil
}

func (s *MockS3Client) CreateBucket(ctx context.Context, bucket string) error {
	s.Called["CreateBucket"]++
	return nil
}

func (s *MockS3Client) DownloadPartialObject(
	ctx context.Context,
	bucket string,
	key string,
	startIndex int64,
	endIndex int64,
) ([]byte, bool, error) {
	s.Called["DownloadPartialObject"]++
	data, ok := s.bucket[key]
	if !ok {
		return []byte{}, false, nil
	}
	if startIndex < 0 || endIndex > int64(len(data)) || startIndex >= endIndex {
		return []byte{}, false, errors.New("invalid startIndex or endIndex")
	}
	return data[startIndex:endIndex], true, nil
}
