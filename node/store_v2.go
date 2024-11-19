package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

const (
	BatchHeaderTableName     = "batch_headers"
	BlobCertificateTableName = "blob_certificates"
	BundleTableName          = "bundles"
)

type StoreV2 interface {
	StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, error)
	DeleteKeys(keys []kvstore.Key) error
	GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error)
}

type storeV2 struct {
	db     kvstore.TableStore
	logger logging.Logger

	ttl time.Duration
}

var _ StoreV2 = &storeV2{}

func NewLevelDBStoreV2(db kvstore.TableStore, logger logging.Logger) *storeV2 {
	return &storeV2{
		db:     db,
		logger: logger,
	}
}

func (s *storeV2) StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, error) {
	dbBatch := s.db.NewTTLBatch()
	keys := make([]kvstore.Key, 0)

	batchHeaderKeyBuilder, err := s.db.GetKeyBuilder(BatchHeaderTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder for batch header: %v", err)
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %v", err)
	}

	// Store batch header
	batchHeaderKey := batchHeaderKeyBuilder.Key(batchHeaderHash[:])
	if _, err = s.db.Get(batchHeaderKey); err == nil {
		return nil, ErrBatchAlreadyExist
	}
	batchHeaderBytes, err := batch.BatchHeader.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize batch header: %v", err)
	}

	keys = append(keys, batchHeaderKey)
	dbBatch.PutWithTTL(batchHeaderKey, batchHeaderBytes, s.ttl)

	// Store blob shards
	for _, bundles := range rawBundles {
		// Store blob certificate
		blobCertificateKeyBuilder, err := s.db.GetKeyBuilder(BlobCertificateTableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get key builder for blob certificate: %v", err)
		}
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %v", err)
		}
		blobCertificateKey := blobCertificateKeyBuilder.Key(blobKey[:])
		blobCertificateBytes, err := bundles.BlobCertificate.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize blob certificate: %v", err)
		}
		keys = append(keys, blobCertificateKey)
		dbBatch.PutWithTTL(blobCertificateKey, blobCertificateBytes, s.ttl)

		// Store bundles
		for quorum, bundle := range bundles.Bundles {
			bundlesKeyBuilder, err := s.db.GetKeyBuilder(BundleTableName)
			if err != nil {
				return nil, fmt.Errorf("failed to get key builder for bundles: %v", err)
			}

			k, err := BundleKey(blobKey, quorum)
			if err != nil {
				return nil, fmt.Errorf("failed to get key for bundles: %v", err)
			}

			keys = append(keys, bundlesKeyBuilder.Key(k))
			dbBatch.PutWithTTL(bundlesKeyBuilder.Key(k), bundle, s.ttl)
		}
	}

	if err := dbBatch.Apply(); err != nil {
		return nil, fmt.Errorf("failed to apply batch: %v", err)
	}

	return keys, nil
}

func (s *storeV2) DeleteKeys(keys []kvstore.Key) error {
	dbBatch := s.db.NewTTLBatch()
	for _, key := range keys {
		dbBatch.Delete(key)
	}
	return dbBatch.Apply()
}

func (s *storeV2) GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error) {
	bundlesKeyBuilder, err := s.db.GetKeyBuilder(BundleTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder for bundles: %v", err)
	}

	k, err := BundleKey(blobKey, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get key for bundles: %v", err)
	}

	bundle, err := s.db.Get(bundlesKeyBuilder.Key(k))
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %v", err)
	}

	chunks, _, err := DecodeChunks(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to decode chunks: %v", err)
	}

	return chunks, nil
}

func BundleKey(blobKey corev2.BlobKey, quorumID core.QuorumID) ([]byte, error) {
	buf := bytes.NewBuffer(blobKey[:])
	err := binary.Write(buf, binary.LittleEndian, quorumID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
