package common

import (
	"context"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
)

// BackendType ... Storage backend type
type BackendType uint8

const (
	EigenDABackendType BackendType = iota
	EigenDAV2BackendType
	MemstoreV1BackendType
	MemstoreV2BackendType
	S3BackendType

	UnknownBackendType
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
	case "unknown":
		fallthrough
	default:
		return UnknownBackendType
	}
}

// GETOpts defines the options for the Get method of a Store.
// The values in here are optional query params for the cert GET routes,
// are parsed in the handlers and passed down to the Store.Get method.
type GETOpts struct {
	// L1 block number at which the cert was included in the rollup batcher inbox.
	// This is optional, and should be set to 0 to mean to skip the RBN recency check.
	// It is impossible for a batch inbox tx to have been included in the genesis block,
	// so we are free to give this special meaning to the zero value.
	//
	// Used to determine the validity of the eigenDA batch.
	// The eigenDA cert contains a reference block number (RBN) which is used
	// to lookup the stake of the eigenda operators before verifying signature thresholds.
	// The rollup commitment containing the eigenDA cert is only valid if it was included
	// within a certain number of blocks after the RBN.
	// validity condition is: certRBN < L1InclusionBlockNum <= RBN + RBNRecencyWindowSize
	L1InclusionBlockNum uint64

	// When true, the Get method will return the encoded_payload without decoding
	// it. This is useful when clients need to decode the encoded_payload themselves,
	// such as inside an fpvm to prove that a decoding fails and can thus be discarded.
	ReturnEncodedPayload bool
}

type Store interface {
	// BackendType returns the backend type provider of the store.
	BackendType() BackendType
}

// EigenDAV1Store is the interface for an EigenDA V1 data store as well as V1 memstore.
type EigenDAV1Store interface {
	Store
	// Put inserts the given value into the key-value (serializedCert-payload) data store.
	Put(ctx context.Context, payload []byte) (serializedCert []byte, err error)
	// Get retrieves the given key if it's present in the key-value (serializedCert-payload) data store.
	Get(ctx context.Context, serializedCert []byte) (payload []byte, err error)
	// Verify verifies the cert and that the payload (after encoding) matches the kzg commitment in the cert.
	Verify(ctx context.Context, serializedCert []byte, payload []byte) error
}

// EigenDAV2Store is the interface for an EigenDA V2 data store as well as V2 memstore.
type EigenDAV2Store interface {
	Store
	// Put inserts the given value into the key-value (serializedCert-payload) data store.
	Put(ctx context.Context, payload []byte) (serializedCert []byte, err error)
	// Get retrieves the given key if it's present in the key-value (serializedCert-payload) data store.
	// If returnEncodedPayload is true, the payload is returned without decoding.
	Get(ctx context.Context,
		versionedCert certs.VersionedCert,
		returnEncodedPayload bool,
	) (payloadOrEncodedPayload []byte, err error)
	// VerifyCert verifies the cert validity and rbn recency.
	VerifyCert(ctx context.Context, versionedCert certs.VersionedCert, l1InclusionBlockNum uint64) error
}

// SecondaryStore is the interface for a key-value data store that uses keccak(value) as the key.
// It is used for Optimism altda keccak commitments, as well as for caching EigenDAStore entries.
type SecondaryStore interface {
	Store
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, key []byte, value []byte) error
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Verify verifies the given key-value pair.
	Verify(ctx context.Context, key []byte, value []byte) error
}
