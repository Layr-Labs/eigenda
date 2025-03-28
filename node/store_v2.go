package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	BatchHeaderTableName     = "batch_headers"
	BlobCertificateTableName = "blob_certificates"
	BundleTableName          = "bundles"
	MigrationTableName       = "migration"

	// LevelDBPath is he path where the levelDB database is stored.
	LevelDBPath = "chunk_v2"
	// LevelDBDeletionPath is the path where the levelDB database is stored while it is being deleted (for atomicity).
	LevelDBDeletionPath = "chunk_v2_deleted"
	// LittDBPath is the path where the littDB database is stored.
	LittDBPath = "chunk_v2_litt"
)

// TODO consider renaming to ValidatorStore

// StoreV2 encapsulates the database for storing batches of chunk data for the V2 validator node.
type StoreV2 interface {

	// StoreBatch stores a batch and its raw bundles in the database. Returns the keys of the stored data
	// and the size of the stored data, in bytes.
	StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, uint64, error)

	// DeleteKeys deletes the keys from local storage.
	//
	// All modifications to the database within this method are performed atomically.
	DeleteKeys(keys []kvstore.Key) error

	// GetChunks returns the chunks of a blob with the given blob key and quorum.
	GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error)
}

type storeV2 struct {
	logger     logging.Logger
	timeSource func() time.Time

	// The levelDB database for storing chunk data. If nil, then the store is backed by a littDB database.
	levelDB kvstore.TableStore

	// The path to the levelDB database.
	levelDBPath string

	// The path to the levelDB database while it is being deleted.
	levelDBDeletionPath string

	// The littDB database for storing chunk data. If nil, then the store has not yet been migrated to littDB.
	littDB litt.DB

	// The table where chunks are stored in the littDB database.
	chunkTable litt.Table

	// The table where batch headers are stored in the littDB database.
	headerTable litt.Table

	// The length of time to store data in the database.
	ttl time.Duration

	// If a migration is in progress, this is the timestamp when the migration is considered to be complete.
	// The migration is completed once all data in levelDB has outlived its TTL.
	migrationCompleteTime time.Time

	// Used to make migration thread safe.
	migrationLock sync.RWMutex

	// A lock used to prevent concurrent requests from storing the same data multiple times.
	duplicateRequestLock *common.IndexLock

	// The salt used to prevent an attacker from causing hash collisions in the duplicate request lock.
	duplicateRequestSalt uint32
}

var _ StoreV2 = &storeV2{}

func NewStoreV2(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	timeSource func() time.Time,
	ttl time.Duration,
	registry *prometheus.Registry) (StoreV2, error) {

	if !config.LittDBEnabled {
		logger.Warn("WARNING: This node is running with littDB disabled. " +
			"This is a deprecated mode of operation, and will not be supported in future versions.")
	}

	littDBPath := path.Join(config.DbPath, LittDBPath)
	levelDBPath := path.Join(config.DbPath, LevelDBPath)
	levelDBDeletionPath := path.Join(config.DbPath, LevelDBDeletionPath)

	// If we previously made an attempt at deleting the levelDB database but it was interrupted, delete it now.
	_, err := os.Stat(levelDBDeletionPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat path %s: %v", levelDBDeletionPath, err)
	}
	if err == nil {
		// The previous attempt at deleting the levelDB database was interrupted.
		logger.Warnf("partial deletion of levelDB database detected at %s. Deleting.", levelDBDeletionPath)
		err = os.RemoveAll(levelDBDeletionPath)
		if err != nil {
			return nil, fmt.Errorf("failed to delete path %s: %v", levelDBDeletionPath, err)
		}
	}

	// Check to see which DBs currently have data on disk.
	_, err = os.Stat(levelDBPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat path %s: %v", levelDBPath, err)
	}
	levelDBExists := err == nil

	_, err = os.Stat(littDBPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat path %s: %v", littDBPath, err)
	}
	littDBExists := err == nil
	if littDBExists && !config.LittDBEnabled {
		return nil, fmt.Errorf("Unable to do backwards migration. Once enabled, littDB cannot be disabled.")
	}

	var littDB litt.DB
	var chunkTable litt.Table
	var headerTable litt.Table
	var levelDB kvstore.TableStore

	// If we are still running with levelDB, start it up.
	if !config.LittDBEnabled || levelDBExists {
		logger.Infof("Using levelDB at %s", levelDBPath)

		levelDB, err = tablestore.Start(logger, &tablestore.Config{
			Type:                       tablestore.LevelDB,
			Path:                       &levelDBPath,
			GarbageCollectionEnabled:   true,
			GarbageCollectionInterval:  time.Duration(config.ExpirationPollIntervalSec) * time.Second,
			GarbageCollectionBatchSize: 1024,
			Schema: []string{
				BatchHeaderTableName,
				BlobCertificateTableName,
				BundleTableName,
				MigrationTableName},
			MetricsRegistry:               registry,
			LevelDBDisableSeeksCompaction: config.LevelDBDisableSeeksCompactionV2,
			LevelDBSyncWrites:             config.LevelDBSyncWritesV2,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create new tablestore: %w", err)
		}
	}

	// Set up migration if necessary.
	var migrationComplete time.Time
	if config.LittDBEnabled && levelDBExists {
		// Both DBs are in play, meaning we are either about to start a migration or already in the middle of one.

		migrationKeyBuilder, err := levelDB.GetKeyBuilder(MigrationTableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get key builder for migration: %v", err)
		}
		migrationKey := migrationKeyBuilder.Key([]byte("migrationCompleteTime"))

		if littDBExists {
			// This is the first time we are starting up with littDB enabled and there is potentially data still
			// in the levelDB. Start the data migration from levelDB to littDB.

			migrationComplete = timeSource().Add(ttl)
			logger.Infof("Begining data migration from levelDB to littDB. Migration will be completed at %s",
				migrationComplete)

			err = levelDB.Put(migrationKey, []byte(fmt.Sprintf("%d", migrationComplete.Unix())))
			if err != nil {
				return nil, fmt.Errorf("failed to put migration key: %v", err)
			}

		} else {
			// A data migration from levelDB to littDB is currently in progress.

			migrationCompleteString, err := levelDB.Get(migrationKey)
			if err != nil {
				return nil, fmt.Errorf("failed to get migration complete time: %v", err)
			}
			migrationCompleteUnix, err := binary.ReadVarint(bytes.NewReader(migrationCompleteString))
			if err != nil {
				return nil, fmt.Errorf("failed to read migration complete time: %v", err)
			}
			migrationComplete = time.Unix(migrationCompleteUnix, 0)

			logger.Infof(
				"Data migration from levelDB to littDB is in progress. Migration will be completed at %s",
				migrationComplete)
		}
	}

	// Start littDB.
	// The ordering of this step is important. By starting littDB, we will create the littDB directory, which
	// will cause the littDBExists variable to be true for all future runs. It's important to have written the
	// migration complete timestamp prior to this happening so that a crash during startup does not leave the
	// migration in a broken state.
	if config.LittDBEnabled {
		logger.Infof("Using littDB at %s", littDBPath)

		littConfig, err := littbuilder.DefaultConfig(littDBPath)
		littConfig.ShardingFactor = 1
		littConfig.MetricsEnabled = true
		littConfig.MetricsRegistry = registry
		littConfig.MetricsNamespace = "node_littdb"
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt config: %w", err)
		}

		littDB, err = littConfig.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt store: %w", err)
		}

		chunkTable, err = littDB.GetTable("chunks")
		if err != nil {
			return nil, fmt.Errorf("failed to get chunks table: %w", err)
		}

		headerTable, err = littDB.GetTable("headers")
		if err != nil {
			return nil, fmt.Errorf("failed to get headers table: %w", err)
		}

		err = chunkTable.SetTTL(ttl)
		if err != nil {
			return nil, fmt.Errorf("failed to set TTL for chunks table: %w", err)
		}
		// We store headers to prevent duplicate insertions. Store the headers for 2x the TTL to ensure that
		// we don't delete the headers before the corresponding chunks.
		err = headerTable.SetTTL(2 * ttl)
		if err != nil {
			return nil, fmt.Errorf("failed to set TTL for headers table: %w", err)
		}
	}

	store := &storeV2{
		logger:                logger,
		timeSource:            timeSource,
		levelDB:               levelDB,
		levelDBPath:           levelDBPath,
		levelDBDeletionPath:   levelDBDeletionPath,
		littDB:                littDB,
		chunkTable:            chunkTable,
		ttl:                   ttl,
		migrationCompleteTime: migrationComplete,
		duplicateRequestLock:  common.NewIndexLock(128),
		duplicateRequestSalt:  rand.Uint32(),
	}

	if config.LittDBEnabled && levelDBExists {
		// This sleeps until the migration is complete then deletes the levelDB database.
		go store.finalizeMigration(ctx)
	}

	return store, nil
}

// TODO get rid of this constructor

// NewLevelDBStoreV2 creates a new StoreV2 backed by a LevelDB database.
func NewLevelDBStoreV2(db kvstore.TableStore, logger logging.Logger, ttl time.Duration) StoreV2 {
	return &storeV2{
		levelDB: db,
		logger:  logger,

		ttl: ttl,
	}
}

func (s *storeV2) StoreBatch(batch *corev2.Batch, rawBundles []*RawBundles) ([]kvstore.Key, uint64, error) {
	if len(rawBundles) == 0 {
		return nil, 0, fmt.Errorf("no raw bundles")
	}
	if len(rawBundles) != len(batch.BlobCertificates) {
		return nil, 0, fmt.Errorf("mismatch between raw bundles (%d) and blob certificates (%d)",
			len(rawBundles), len(batch.BlobCertificates))
	}

	if s.littDB != nil {
		size, err := s.storeBatchLittDB(batch, rawBundles)
		return nil, size, err
	} else {
		return s.storeBatchLevelDB(batch, rawBundles)
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

// TODO do we really need to be returning size?

// storeBatchHeader stores the batch header in the database, returning an error if it is already present.
// This method is guaranteed to only return nil exactly once if called multiple times with the same batch header hash.
func (s *storeV2) storeBatchHeader(batchHeader *corev2.BatchHeader) (uint64, error) {

	batchHeaderHash, err := batchHeader.Hash()
	if err != nil {
		return 0, fmt.Errorf("failed to hash batch header: %v", err)
	}
	batchHeaderBytes, err := batchHeader.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize batch header: %v", err)
	}

	// Grab a lock that is mutually exclusive for this batch header hash. This prevents us from storing
	// data twice if we receive concurrent requests to store the same batch.
	lockIndex := uint64(util.HashKey(batchHeaderHash[:], s.duplicateRequestSalt))
	s.duplicateRequestLock.Lock(lockIndex)
	defer s.duplicateRequestLock.Unlock(lockIndex)

	ok, err := s.headerTable.Exists(batchHeaderHash[:])
	if err != nil {
		return 0, fmt.Errorf("failed to check batch header existence: %v", err)
	}
	if ok {
		return 0, ErrBatchAlreadyExist
	}

	// Store batch header.
	err = s.headerTable.Put(batchHeaderHash[:], batchHeaderBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to put batch header: %v", err)
	}

	return uint64(len(batchHeaderBytes)), nil
}

func (s *storeV2) storeBatchLittDB(batch *corev2.Batch, rawBundles []*RawBundles) (uint64, error) {
	var size uint64

	// Store the batch header to prevent duplicate insertions.
	batchHeaderSize, err := s.storeBatchHeader(batch.BatchHeader)
	if err != nil {
		return 0, fmt.Errorf("failed to store batch header: %v", err)
	}
	size += batchHeaderSize
	err = s.headerTable.Flush()
	if err != nil {
		return 0, fmt.Errorf("failed to flush header table: %v", err)
	}

	// Now that the batch header is durable on disk, this validator will reject all future requests to store this batch.
	// If the validator crashes between now and when it returns the availability signature, it will never
	// return the availability signature for this batch. Although the validator could, in theory, reprocess this
	// batch after starting back up, the performance and complexity overhead of doing so is cost prohibitive.
	// Even if the validator was capable of reprocessing the batch after a crash, it's highly likely that the time
	// window for returning the availability signature for this batch would have already passed -- thus negating
	// any potential benefit of being able to reprocess the batch. It is, after all, not unreasonable to expect that
	// a validator may not sign for some batches if it crashes.

	// Store chunk data.
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

// TODO this should not be a public method!!
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
	if s.littDB == nil {
		// We haven't thrown the switch for littDB yet. Just look in levelDB.
		return s.getChunksLevelDB(blobKey, quorum)
	}

	// Regardless of migration status, always check littDB first.
	data, ok, err := s.getChunksLittDB(blobKey, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks: %v", err)
	}

	if !ok {
		// The data wasn't found in littDB.

		if s.levelDB == nil {
			// There is no data in levelDB.
			return nil, fmt.Errorf("failed to get chunks: not found")
		}

		s.migrationLock.RLock()
		defer s.migrationLock.RUnlock()

		if s.levelDB == nil {
			// The migration completed while we were waiting for the lock.
			return nil, fmt.Errorf("failed to get chunks: not found")
		}

		return s.getChunksLevelDB(blobKey, quorum)
	}

	return data, nil
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

func (s *storeV2) getChunksLittDB(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, bool, error) {
	k, err := BundleKey(blobKey, quorum)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get key for bundles: %v", err)
	}

	bundle, ok, err := s.chunkTable.Get(k)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get bundle: %v", err)
	}
	if !ok {
		return nil, false, nil
	}

	chunks, _, err := DecodeChunks(bundle)
	if err != nil {
		return nil, false, fmt.Errorf("failed to decode chunks: %v", err)
	}

	return chunks, true, nil
}

// finalizeMigration sleeps until the migration is complete, then deletes the levelDB database.
func (s *storeV2) finalizeMigration(ctx context.Context) {
	timeUntilMigrationComplete := s.migrationCompleteTime.Sub(s.timeSource())

	select {
	case <-ctx.Done():
		s.logger.Info("context cancelled, migration finalization aborted")
		return
	case <-time.After(timeUntilMigrationComplete):
		s.migrationLock.Lock()
		defer s.migrationLock.Unlock()

		s.logger.Infof("migration to littDB complete, deleting levelDB at %s", s.levelDBPath)

		err := s.levelDB.Shutdown()
		if err != nil {
			s.logger.Errorf("failed to stop levelDB: %v", err)
			return
		}

		// In order to make levelDB deletion atomic, first rename it.
		err = os.Rename(s.levelDBPath, s.levelDBDeletionPath)
		if err != nil {
			s.logger.Errorf("failed to rename levelDB: %v", err)
		}

		// Now, delete the levelDB database.
		err = os.RemoveAll(s.levelDBDeletionPath)
		if err != nil {
			s.logger.Errorf("failed to delete levelDB: %v", err)
			return
		}

		s.levelDB = nil
		s.logger.Infof("levelDB has been deleted")
	}
}

func BundleKey(blobKey corev2.BlobKey, quorumID core.QuorumID) ([]byte, error) {
	buf := bytes.NewBuffer(blobKey[:])
	err := binary.Write(buf, binary.LittleEndian, quorumID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
