package store

import (
	"context"
	"fmt"
	"strings"
)

type BackendType uint8

const (
	EigenDA BackendType = iota
	Memory
	S3
	Redis

	Unknown
)

var (
	ErrProxyOversizedBlob   = fmt.Errorf("encoded blob is larger than max blob size")
	ErrEigenDAOversizedBlob = fmt.Errorf("blob size cannot exceed")
)

func (b BackendType) String() string {
	switch b {
	case EigenDA:
		return "EigenDA"
	case Memory:
		return "Memory"
	case S3:
		return "S3"
	case Redis:
		return "Redis"
	case Unknown:
		fallthrough
	default:
		return "Unknown"
	}
}

func StringToBackendType(s string) BackendType {
	lower := strings.ToLower(s)

	switch lower {
	case "eigenda":
		return EigenDA
	case "memory":
		return Memory
	case "s3":
		return S3
	case "redis":
		return Redis
	case "unknown":
		fallthrough
	default:
		return Unknown
	}
}

// Used for E2E tests
type Stats struct {
	Entries int
	Reads   int
}

type Store interface {
	// Stats returns the current usage metrics of the key-value data store.
	Stats() *Stats
	// Backend returns the backend type provider of the store.
	BackendType() BackendType
	// Verify verifies the given key-value pair.
	Verify(key []byte, value []byte) error
}

type KeyGeneratedStore interface {
	Store
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, value []byte) (key []byte, err error)
}

type PrecomputedKeyStore interface {
	Store
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, key []byte, value []byte) error
}
