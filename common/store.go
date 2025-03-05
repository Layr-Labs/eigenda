package common

import (
	"context"
	"fmt"
	"strings"
)

// BackendType ... Storage backend type
type BackendType uint8

const (
	EigenDABackendType BackendType = iota
	EigenDAV2BackendType
	MemstoreV1BackendType
	MemstoreV2BackendType
	S3BackendType
	RedisBackendType

	UnknownBackendType
)

var (
	ErrProxyOversizedBlob = fmt.Errorf("encoded blob is larger than max blob size")
)

func (b BackendType) String() string {
	switch b {
	case EigenDABackendType:
		return "EigenDA"
	case EigenDAV2BackendType:
		return "EigenDAV2"
	case MemstoreV1BackendType:
		return "EigenDAV1Memstore"
	case MemstoreV2BackendType:
		return "EigenDAV2Memstore"
	case S3BackendType:
		return "S3"
	case RedisBackendType:
		return "Redis"
	case UnknownBackendType:
		fallthrough
	default:
		return "Unknown"
	}
}

func StringToBackendType(s string) BackendType {
	lower := strings.ToLower(s)

	switch lower {
	case "eigenda":
		return EigenDABackendType
	case "eigenda_v2":
		return EigenDAV2BackendType
	case "memory_v1":
		return MemstoreV1BackendType
	case "memory_v2":
		return MemstoreV2BackendType
	case "s3":
		return S3BackendType
	case "redis":
		return RedisBackendType
	case "unknown":
		fallthrough
	default:
		return UnknownBackendType
	}
}

type Store interface {
	// BackendType returns the backend type provider of the store.
	BackendType() BackendType
	// Verify verifies the given key-value pair.
	Verify(ctx context.Context, key []byte, value []byte) error
}

type GeneratedKeyStore interface {
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
