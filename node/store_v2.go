package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	BatchHeaderTableName     = "batch_headers"
	BlobCertificateTableName = "blob_certificates"
	BundleTableName          = "bundles"
)

type StoreV2 interface {

	// StoreBatch stores a batch and its raw bundles in the database. Returns the keys of the stored data
	// and the size of the stored data, in bytes.
	//
	// All modifications to the database within this method are performed atomically.
	StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, uint64, error)

	// DeleteKeys deletes the keys from local storage.
	//
	// All modifications to the database within this method are performed atomically.
	DeleteKeys(keys []kvstore.Key) error

	// GetChunks returns the chunks of a blob with the given blob key and quorum.
	GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error)
}

type storeV2 struct {
	levelDB    kvstore.TableStore
	littDB     litt.DB
	chunkTable litt.Table
	logger     logging.Logger
	ttl        time.Duration
}

var _ StoreV2 = &storeV2{}

func NewStoreV2(
	config *Config,
	ttl time.Duration,
	logger logging.Logger,
	registry *prometheus.Registry) (StoreV2, error) {

	if config.LittDBEnabled {
		littDBPath := config.DbPath + "/chunk_v2_litt"
		littConfig, err := littbuilder.DefaultConfig(littDBPath)
		littConfig.ShardingFactor = 1
		littConfig.MetricsEnabled = true
		littConfig.MetricsRegistry = registry
		littConfig.MetricsNamespace = "node_littdb"
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt config: %w", err)
		}

		littDB, err := littConfig.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt store: %w", err)
		}

		chunkTable, err := littDB.GetTable("chunks")
		if err != nil {
			return nil, fmt.Errorf("failed to get chunks table: %w", err)
		}

		err = chunkTable.SetTTL(ttl)
		if err != nil {
			return nil, fmt.Errorf("failed to set TTL for chunks table: %w", err)
		}

		return &storeV2{
			littDB:     littDB,
			chunkTable: chunkTable,
			logger:     logger,
			ttl:        ttl,
		}, nil
	} else {
		levelDBPath := config.DbPath + "/chunk_v2"
		levelDB, err := tablestore.Start(logger, &tablestore.Config{
			Type:                          tablestore.LevelDB,
			Path:                          &levelDBPath,
			GarbageCollectionEnabled:      true,
			GarbageCollectionInterval:     time.Duration(config.ExpirationPollIntervalSec) * time.Second,
			GarbageCollectionBatchSize:    1024,
			Schema:                        []string{BatchHeaderTableName, BlobCertificateTableName, BundleTableName},
			MetricsRegistry:               registry,
			LevelDBDisableSeeksCompaction: config.LevelDBDisableSeeksCompactionV2,
			LevelDBSyncWrites:             config.LevelDBSyncWritesV2,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create new tablestore: %w", err)
		}

		return &storeV2{
			levelDB: levelDB,
			logger:  logger,

			ttl: ttl,
		}, nil
	}

}

// NewLevelDBStoreV2 creates a new StoreV2 backed by a LevelDB database.
func NewLevelDBStoreV2(db kvstore.TableStore, logger logging.Logger, ttl time.Duration) StoreV2 {
	return &storeV2{
		levelDB: db,
		logger:  logger,

		ttl: ttl,
	}
}

// NewLittDBStoreV2 creates a new StoreV2 backed by a Litt database.
func NewLittDBStoreV2(db litt.DB, logger logging.Logger, ttl time.Duration) (StoreV2, error) {
	chunkTable, err := db.GetTable("chunks")
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks table: %w", err)
	}

	err = chunkTable.SetTTL(ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to set TTL for chunks table: %w", err)
	}

	return &storeV2{
		chunkTable: chunkTable,
		logger:     logger,
		ttl:        ttl,
	}, nil
}

func (s *storeV2) StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, uint64, error) {
	if len(rawBundles) == 0 {
		return nil, 0, fmt.Errorf("no raw bundles")
	}
	if len(rawBundles) != len(batch.BlobCertificates) {
		return nil, 0, fmt.Errorf("mismatch between raw bundles (%d) and blob certificates (%d)", len(rawBundles), len(batch.BlobCertificates))
	}

	if s.levelDB != nil {
		return s.storeBatchLevelDB(batch, rawBundles)
	} else {
		size, err := s.storeBatchLittDB(batch, rawBundles)
		return nil, size, err
	}
}

func (s *storeV2) storeBatchLevelDB(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, uint64, error) {
	dbBatch := s.levelDB.NewTTLBatch()
	var size uint64

	keys := make([]kvstore.Key, 0)

	batchHeaderKeyBuilder, err := s.levelDB.GetKeyBuilder(BatchHeaderTableName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get key builder for batch header: %v", err)
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to hash batch header: %v", err)
	}

	// Store batch header
	batchHeaderKey := batchHeaderKeyBuilder.Key(batchHeaderHash[:])
	if _, err = s.levelDB.Get(batchHeaderKey); err == nil {
		return nil, 0, ErrBatchAlreadyExist
	}
	batchHeaderBytes, err := batch.BatchHeader.Serialize()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize batch header: %v", err)
	}

	dbBatch.PutWithTTL(batchHeaderKey, batchHeaderBytes, s.ttl)
	keys = append(keys, batchHeaderKey)
	size += uint64(len(batchHeaderBytes))

	// Store blob shards
	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get blob key: %v", err)
		}

		// Store bundles
		for quorum, bundle := range bundles.Bundles {
			bundlesKeyBuilder, err := s.levelDB.GetKeyBuilder(BundleTableName)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to get key builder for bundles: %v", err)
			}

			k, err := BundleKey(blobKey, quorum)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to get key for bundles: %v", err)
			}

			keys = append(keys, bundlesKeyBuilder.Key(k))
			dbBatch.PutWithTTL(bundlesKeyBuilder.Key(k), bundle, s.ttl)
			size += uint64(len(bundle))
		}
	}

	if err := dbBatch.Apply(); err != nil {
		return nil, 0, fmt.Errorf("failed to apply batch: %v", err)
	}

	return keys, size, nil
}

func (s *storeV2) storeBatchLittDB(batch *corev2.Batch, rawBundles []*RawBundles) (uint64, error) {
	var size uint64

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return 0, fmt.Errorf("failed to hash batch header: %v", err)
	}

	// Don't store duplicate requests
	_, ok, err := s.chunkTable.Get(batchHeaderHash[:])
	if err != nil {
		return 0, fmt.Errorf("failed to check batch header: %v", err)
	}
	if ok {
		return 0, ErrBatchAlreadyExist
	}

	// Store batch header
	batchHeaderBytes, err := batch.BatchHeader.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize batch header: %v", err)
	}
	err = s.chunkTable.Put(batchHeaderHash[:], batchHeaderBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to put batch header: %v", err)
	}
	size += uint64(len(batchHeaderBytes))

	// Store blob shards
	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return 0, fmt.Errorf("failed to get blob key: %v", err)
		}

		// Store bundles
		for quorum, bundle := range bundles.Bundles {
			k, err := BundleKey(blobKey, quorum)
			if err != nil {
				return 0, fmt.Errorf("failed to get key for bundles: %v", err)
			}

			err = s.chunkTable.Put(k, bundle)
			if err != nil {
				return 0, fmt.Errorf("failed to put bundle: %v", err)
			}

			size += uint64(len(bundle))
		}
	}

	err = s.chunkTable.Flush()
	if err != nil {
		return 0, fmt.Errorf("failed to flush chunk table: %v", err)
	}

	return size, nil
}

func (s *storeV2) DeleteKeys(keys []kvstore.Key) error {
	if s.levelDB == nil {
		return fmt.Errorf("littDB does not support deletion")
	}

	dbBatch := s.levelDB.NewTTLBatch()
	for _, key := range keys {
		dbBatch.Delete(key)
	}
	return dbBatch.Apply()
}

func (s *storeV2) GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error) {
	if s.levelDB != nil {
		return s.getChunksLevelDB(blobKey, quorum)
	} else {
		return s.getChunksLittDB(blobKey, quorum)
	}
}

func (s *storeV2) getChunksLevelDB(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error) {
	bundlesKeyBuilder, err := s.levelDB.GetKeyBuilder(BundleTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder for bundles: %v", err)
	}

	k, err := BundleKey(blobKey, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get key for bundles: %v", err)
	}

	bundle, err := s.levelDB.Get(bundlesKeyBuilder.Key(k))
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %v", err)
	}

	chunks, _, err := DecodeChunks(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to decode chunks: %v", err)
	}

	return chunks, nil
}

func (s *storeV2) getChunksLittDB(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error) {
	k, err := BundleKey(blobKey, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get key for bundles: %v", err)
	}

	bundle, ok, err := s.chunkTable.Get(k)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("bundle not found")
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
