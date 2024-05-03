package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	plasma "github.com/Layr-Labs/op-plasma-eigenda"
)

type S3Store struct {
	bucket string
	client *s3.Client
}

func NewS3Store(ctx context.Context, bucket string) (*S3Store, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &S3Store{
		bucket: bucket,
		client: s3.NewFromConfig(sdkConfig),
	}, nil
}

func (s *S3Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(hex.EncodeToString(key)),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return nil, plasma.ErrNotFound
			}
		}
		return nil, err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *S3Store) Put(ctx context.Context, key []byte, value []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(hex.EncodeToString(key)),
		Body:   bytes.NewReader(value),
	})
	return err
}

func (s *S3Store) PutWithComm(ctx context.Context, key []byte, value []byte) error {
	return s.Put(ctx, key, value)
}

func (s *S3Store) PutWithoutComm(ctx context.Context, value []byte) (key []byte, err error) {
	// make key fingerprint of value
	// this could result in collisions
	hasher := sha256.New()
	hasher.Write(value)
	bs := hasher.Sum(nil)

	if err := s.PutWithComm(ctx, bs, value); err != nil {
		return nil, err
	}

	return bs, err
}
