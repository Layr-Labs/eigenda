package node

import (
	"bytes"
	"context"
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

type StoreV2 struct {
	db     kvstore.TableStore
	logger logging.Logger

	ttl time.Duration
}

func NewLevelDBStoreV2(db kvstore.TableStore, logger logging.Logger) *StoreV2 {
	return &StoreV2{
		db:     db,
		logger: logger,
	}
}

func (s *StoreV2) StoreBatch(ctx context.Context, batch *corev2.Batch, rawBundles []*RawBundles) error {
	dbBatch := s.db.NewTTLBatch()

	batchHeaderKeyBuilder, err := s.db.GetKeyBuilder(BatchHeaderTableName)
	if err != nil {
		return fmt.Errorf("failed to get key builder for batch header: %v", err)
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash batch header: %v", err)
	}

	// Store batch header
	batchHeaderBytes, err := batch.BatchHeader.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize batch header: %v", err)
	}

	batchHeaderKey := batchHeaderKeyBuilder.Key(batchHeaderHash[:])
	dbBatch.PutWithTTL(batchHeaderKey, batchHeaderBytes, s.ttl)

	// Store blob shards
	for _, bundles := range rawBundles {
		// Store blob certificate
		blobCertificateKeyBuilder, err := s.db.GetKeyBuilder(BlobCertificateTableName)
		if err != nil {
			return fmt.Errorf("failed to get key builder for blob certificate: %v", err)
		}
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return fmt.Errorf("failed to get blob key: %v", err)
		}
		blobCertificateKey := blobCertificateKeyBuilder.Key(blobKey[:])
		blobCertificateBytes, err := bundles.BlobCertificate.Serialize()
		if err != nil {
			return fmt.Errorf("failed to serialize blob certificate: %v", err)
		}
		dbBatch.PutWithTTL(blobCertificateKey, blobCertificateBytes, s.ttl)

		// Store bundles
		for quorum, bundle := range bundles.Bundles {
			bundlesKeyBuilder, err := s.db.GetKeyBuilder(BundleTableName)
			if err != nil {
				return fmt.Errorf("failed to get key builder for bundles: %v", err)
			}

			k, err := BundleKey(blobKey, quorum)
			if err != nil {
				return fmt.Errorf("failed to get key for bundles: %v", err)
			}

			dbBatch.PutWithTTL(bundlesKeyBuilder.Key(k), bundle, s.ttl)
		}
	}

	return dbBatch.Apply()
}

func BundleKey(blobKey corev2.BlobKey, quorumID core.QuorumID) ([]byte, error) {
	buf := bytes.NewBuffer(blobKey[:])
	err := binary.Write(buf, binary.LittleEndian, quorumID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
