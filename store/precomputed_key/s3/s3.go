package s3

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/minio/minio-go/v7"

	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
)

const (
	CredentialTypeStatic  CredentialType = "static"
	CredentialTypeIAM     CredentialType = "iam"
	CredentialTypeUnknown CredentialType = "unknown"
)

func StringToCredentialType(s string) CredentialType {
	switch s {
	case "static":
		return CredentialTypeStatic
	case "iam":
		return CredentialTypeIAM
	default:
		return CredentialTypeUnknown
	}
}

var _ store.PrecomputedKeyStore = (*Store)(nil)

type CredentialType string
type Config struct {
	CredentialType  CredentialType
	Endpoint        string
	EnableTLS       bool
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	Path            string
	Backup          bool
	Timeout         time.Duration
	Profiling       bool
}

type Store struct {
	cfg              Config
	client           *minio.Client
	putObjectOptions minio.PutObjectOptions
	stats            *store.Stats
}

func NewS3(cfg Config) (*Store, error) {
	endpointURL, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	putObjectOptions := minio.PutObjectOptions{}
	if s3utils.IsGoogleEndpoint(*endpointURL) {
		putObjectOptions.DisableContentSha256 = true // Avoid chunk signatures on GCS: https://github.com/minio/minio-go/issues/1922
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  creds(cfg),
		Secure: cfg.EnableTLS,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		cfg:              cfg,
		client:           client,
		putObjectOptions: putObjectOptions,
		stats: &store.Stats{
			Entries: 0,
			Reads:   0,
		},
	}, nil
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	result, err := s.client.GetObject(ctx, s.cfg.Bucket, path.Join(s.cfg.Path, hex.EncodeToString(key)), minio.GetObjectOptions{})
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
		s.stats.Reads++
	}

	return data, nil
}

func (s *Store) Put(ctx context.Context, key []byte, value []byte) error {
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, path.Join(s.cfg.Path, hex.EncodeToString(key)), bytes.NewReader(value), int64(len(value)), s.putObjectOptions)
	if err != nil {
		return err
	}

	if s.cfg.Profiling {
		s.stats.Entries++
	}

	return nil
}

func (s *Store) Verify(key []byte, value []byte) error {
	h := crypto.Keccak256Hash(value)
	if !bytes.Equal(h[:], key) {
		return errors.New("key does not match value")
	}

	return nil
}

func (s *Store) Stats() *store.Stats {
	return s.stats
}

func (s *Store) BackendType() store.BackendType {
	return store.S3BackendType
}

func creds(cfg Config) *credentials.Credentials {
	if cfg.CredentialType == CredentialTypeIAM {
		return credentials.NewIAM("")
	}
	return credentials.NewStaticV4(cfg.AccessKeyID, cfg.AccessKeySecret, "")
}
