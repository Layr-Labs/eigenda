package s3

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/minio/minio-go/v7"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	CredentialTypeStatic  CredentialType = "static"
	CredentialTypeIAM     CredentialType = "iam"
	CredentialTypePublic  CredentialType = "public"
	CredentialTypeUnknown CredentialType = "unknown"
)

func StringToCredentialType(s string) CredentialType {
	switch s {
	case "static":
		return CredentialTypeStatic
	case "iam":
		return CredentialTypeIAM
	case "public":
		return CredentialTypePublic
	default:
		return CredentialTypeUnknown
	}
}

var _ common.PrecomputedKeyStore = (*Store)(nil)

type CredentialType string
type Config struct {
	CredentialType  CredentialType
	Endpoint        string
	EnableTLS       bool
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	Path            string
}

// Custom MarshalJSON function to control what gets included in the JSON output
// TODO: Probably best would be to separate config from secrets everywhere.
// Then we could just log the config and not worry about secrets.
func (c Config) MarshalJSON() ([]byte, error) {
	type Alias Config // Use an alias to avoid recursion with MarshalJSON
	aux := (Alias)(c)
	// Conditionally include a masked password if it is set
	if aux.AccessKeySecret != "" {
		aux.AccessKeySecret = "*****"
	}
	return json.Marshal(aux)
}

// Store ... S3 store
// client safe for concurrent use: https://github.com/minio/minio-go/issues/598#issuecomment-569457863
type Store struct {
	cfg              Config
	client           *minio.Client
	putObjectOptions minio.PutObjectOptions
}

func isGoogleEndpoint(endpoint string) bool {
	return strings.Contains(endpoint, "storage.googleapis.com")
}

func NewStore(cfg Config) (*Store, error) {
	putObjectOptions := minio.PutObjectOptions{}
	if isGoogleEndpoint(cfg.Endpoint) {
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

	return data, nil
}

func (s *Store) Put(ctx context.Context, key []byte, value []byte) error {
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, path.Join(s.cfg.Path, hex.EncodeToString(key)), bytes.NewReader(value), int64(len(value)), s.putObjectOptions)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Verify(_ context.Context, key []byte, value []byte) error {
	h := crypto.Keccak256Hash(value)
	if !bytes.Equal(h[:], key) {
		return fmt.Errorf("key does not match value, expected: %s got: %s", hex.EncodeToString(key), h.Hex())
	}

	return nil
}

func (s *Store) BackendType() common.BackendType {
	return common.S3BackendType
}

func creds(cfg Config) *credentials.Credentials {
	if cfg.CredentialType == CredentialTypeIAM {
		return credentials.NewIAM("")
	}
	if cfg.CredentialType == CredentialTypePublic {
		return nil
	}
	return credentials.NewStaticV4(cfg.AccessKeyID, cfg.AccessKeySecret, "")
}
