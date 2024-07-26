package store

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	Bucket          string
	Path            string
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Profiling       bool
}

type S3Store struct {
	cfg    S3Config
	client *minio.Client
	stats  *Stats
}

func NewS3(cfg S3Config) (*S3Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewIAM(""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &S3Store{
		cfg:    cfg,
		client: client,
		stats: &Stats{
			Entries: 0,
			Reads:   0,
		},
	}, nil
}

func (s *S3Store) Get(ctx context.Context, key []byte) ([]byte, error) {

	result, err := s.client.GetObject(ctx, s.cfg.Bucket, s.cfg.Path+hex.EncodeToString(key), minio.GetObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, errors.New("value not found in s3 bucket")
		}
		return nil, err
	}
	defer result.Close()
	data, err := io.ReadAll(result)
	if err != nil {
		return nil, err
	}

	if s.cfg.Profiling {
		s.stats.Reads += 1
	}

	return data, nil
}

func (s *S3Store) Put(ctx context.Context, key []byte, value []byte) error {
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, s.cfg.Path+hex.EncodeToString(key), bytes.NewReader(value), int64(len(value)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	if s.cfg.Profiling {
		s.stats.Entries += 1
	}

	return nil
}

func (s *S3Store) Stats() *Stats {
	return s.stats
}
