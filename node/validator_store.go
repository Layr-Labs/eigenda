package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
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
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	BatchHeaderTableName     = "batch_headers"
	BlobCertificateTableName = "blob_certificates"
	BundleTableName          = "bundles"
	MigrationTableName       = "migration"

	// LevelDBPath is the path where the levelDB database is stored.
	LevelDBPath = "chunk_v2"
	// LevelDBDeletionPath is the path where the levelDB database is stored while it is being deleted (for atomicity).
	LevelDBDeletionPath = "chunk_v2_deleted"
	// LittDBPath is the path where the littDB database is stored.
	LittDBPath = "chunk_v2_litt"

	// The name of the littDB table containing chunk data.
	chunksTableName = "chunks"
	// A legacy littDB table, this will exist until all old implementations are migrated and delete this table.
	headersTableName = "headers"
	// The metrics prefix for littDB.
	littDBMetricsPrefix = "node_littdb"
)

// BundleToStore is a struct that holds the bundle key and the bundle bytes.
type BundleToStore struct {
	// A bundle key, as encoded by BundleKey()
	BundleKey []byte
	// The binary bundle bytes.
	BundleBytes []byte
}

// ValidatorStore encapsulates the database for storing batches of chunk data for the V2 validator node.
type ValidatorStore interface {

	// StoreBatch stores a batch and its raw bundles in the database. Returns the keys of the stored data
	// and the size of the stored data, in bytes.
	StoreBatch(batchHeaderHash []byte, batchData []*BundleToStore) ([]kvstore.Key, uint64, error)

	// DeleteKeys deletes the keys from local storage.
	//
	// All modifications to the database within this method are performed atomically.
	DeleteKeys(keys []kvstore.Key) error

	// GetBundleData returns the chunks of a blob with the given bundle key.
	// The returned chunks are encoded in bundle format.
	GetBundleData(bundleKey []byte) ([]byte, error)

	// Stop stops the store.
	Stop() error
}

type validatorStore struct {
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

	// A flag indicating whether the migration is complete. Used to prevent a double migration race condition
	// (which is possible only in a unit test).
	migrationComplete bool
}

var _ ValidatorStore = &validatorStore{}

func NewValidatorStore(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	timeSource func() time.Time,
	ttl time.Duration,
	registry *prometheus.Registry) (ValidatorStore, error) {

	if !config.LittDBEnabled {
		logger.Warn("WARNING: This node is running with littDB disabled. " +
			"This is a deprecated mode of operation, and will not be supported in future versions.")
	}

	littDBPath := path.Join(config.DbPath, LittDBPath)
	levelDBPath := path.Join(config.DbPath, LevelDBPath)
	levelDBDeletionPath := path.Join(config.DbPath, LevelDBDeletionPath)

	// If we previously made an attempt at deleting the levelDB database but it was interrupted, delete it now.
	exists, err := util.Exists(levelDBDeletionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path %s: %v", levelDBDeletionPath, err)
	}
	if exists {
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

	// Set up migration if necessary. The migrationComplete time is the time when all data in the old levelDB instance
	// has exceeded its TTL. If the TTL is 2 weeks, then the migrationComplete time will be 2 weeks from now. At that
	// moment, it becomes safe to permanently stop and delete the levelDB database.
	var migrationComplete time.Time
	if config.LittDBEnabled && levelDBExists {
		// Both DBs are in play, meaning we are either about to start a migration or already in the middle of one.

		migrationKeyBuilder, err := levelDB.GetKeyBuilder(MigrationTableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get key builder for migration: %v", err)
		}
		migrationKey := migrationKeyBuilder.Key([]byte("migrationCompleteTime"))

		if !littDBExists {
			// This is the first time we are starting up with littDB enabled and there is potentially data still
			// in the levelDB. Start the data migration from levelDB to littDB.

			migrationComplete = timeSource().Add(ttl)
			logger.Infof("Beginning data migration from levelDB to littDB. Migration will be completed at %s",
				migrationComplete)

			migrationCompleteUnix := migrationComplete.Unix()
			err = levelDB.Put(migrationKey, []byte(fmt.Sprintf("%d", migrationCompleteUnix)))
			if err != nil {
				return nil, fmt.Errorf("failed to put migration key: %v", err)
			}

		} else {
			// A data migration from levelDB to littDB is currently in progress.

			migrationCompleteString, err := levelDB.Get(migrationKey)
			if err != nil {
				return nil, fmt.Errorf("failed to get migration complete time: %v", err)
			}
			migrationCompleteUnix, err := strconv.ParseUint(string(migrationCompleteString), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to read migration complete time: %v", err)
			}
			migrationComplete = time.Unix(int64(migrationCompleteUnix), 0)

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

		littConfig, err := litt.DefaultConfig(littDBPath)
		littConfig.ShardingFactor = 1
		littConfig.MetricsEnabled = true
		littConfig.MetricsRegistry = registry
		littConfig.MetricsNamespace = littDBMetricsPrefix
		littConfig.Logger = logger
		littConfig.DoubleWriteProtection = config.LittDBDoubleWriteProtection
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt config: %w", err)
		}

		littDB, err = littbuilder.NewDB(littConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create new litt store: %w", err)
		}

		chunkTable, err = littDB.GetTable(chunksTableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get chunks table: %w", err)
		}

		err = chunkTable.SetCacheSize(uint64(config.LittDBChunkCacheSizeGB * units.GiB))
		if err != nil {
			return nil, fmt.Errorf("failed to set cache size for chunks table: %w", err)
		}

		// A prior implementation stored data here. Delete it if it exists.
		// This is safe to delete once all old validators have been migrated to the new version.
		err = littDB.DropTable(headersTableName)
		if err != nil {
			return nil, fmt.Errorf("failed to drop headers table: %w", err)
		}

		err = chunkTable.SetTTL(ttl)
		if err != nil {
			return nil, fmt.Errorf("failed to set TTL for chunks table: %w", err)
		}
	}

	store := &validatorStore{
		logger:                logger,
		timeSource:            timeSource,
		levelDB:               levelDB,
		levelDBPath:           levelDBPath,
		levelDBDeletionPath:   levelDBDeletionPath,
		littDB:                littDB,
		chunkTable:            chunkTable,
		headerTable:           headerTable,
		ttl:                   ttl,
		migrationCompleteTime: migrationComplete,
		duplicateRequestLock:  common.NewIndexLock(1024),
		duplicateRequestSalt:  rand.Uint32(),
	}

	if config.LittDBEnabled && levelDBExists {
		// This sleeps until the migration is complete then deletes the levelDB database.
		go store.finalizeMigration(ctx)
	}

	return store, nil
}

func (s *validatorStore) StoreBatch(batchHeaderHash []byte, batchData []*BundleToStore) ([]kvstore.Key, uint64, error) {
	if len(batchData) == 0 {
		return nil, 0, fmt.Errorf("no batch data")
	}

	if s.littDB != nil {
		size, err := s.storeBatchLittDB(batchData)
		return nil, size, err
	} else {
		return s.storeBatchLevelDB(batchHeaderHash, batchData)
	}
}

func (s *validatorStore) storeBatchLevelDB(batchHeaderHash []byte, batchData []*BundleToStore) ([]kvstore.Key, uint64, error) {
	dbBatch := s.levelDB.NewTTLBatch()
	var size uint64

	keys := make([]kvstore.Key, 0, len(batchData))

	batchHeaderKeyBuilder, err := s.levelDB.GetKeyBuilder(BatchHeaderTableName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get key builder for batch header: %v", err)
	}

	// Store batch header
	batchHeaderKey := batchHeaderKeyBuilder.Key(batchHeaderHash[:])
	if _, err = s.levelDB.Get(batchHeaderKey); err == nil {
		return nil, 0, ErrBatchAlreadyExist
	}

	dbBatch.PutWithTTL(batchHeaderKey, []byte{}, s.ttl)
	keys = append(keys, batchHeaderKey)
	size += uint64(len(batchHeaderKey.Raw()))

	// Store blob shards
	for _, batchDatum := range batchData {

		bundleKeyBytes := batchDatum.BundleKey
		bundleData := batchDatum.BundleBytes

		bundleKeyBuilder, err := s.levelDB.GetKeyBuilder(BundleTableName)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get key builder for bundles: %v", err)
		}
		bundleKey := bundleKeyBuilder.Key(bundleKeyBytes)

		keys = append(keys, bundleKey)
		dbBatch.PutWithTTL(bundleKey, bundleData, s.ttl)
		size += uint64(len(bundleData) + len(bundleKey.Raw()))

	}

	if err := dbBatch.Apply(); err != nil {
		return nil, 0, fmt.Errorf("failed to apply batch: %v", err)
	}

	return keys, size, nil
}

func (s *validatorStore) storeBatchLittDB(batchData []*BundleToStore) (uint64, error) {
	var size uint64

	writeCompleteChan := make(chan error, len(batchData))
	for _, batchDatum := range batchData {
		bundleKeyBytes := batchDatum.BundleKey
		bundleData := batchDatum.BundleBytes

		go func() {
			// Grab a lock on the hash of the blob. This protects against duplicate writes of the same blob.
			lockIndex := uint64(util.HashKey(bundleKeyBytes[:], s.duplicateRequestSalt))
			s.duplicateRequestLock.Lock(lockIndex)
			defer s.duplicateRequestLock.Unlock(lockIndex)

			exists, err := s.chunkTable.Exists(bundleKeyBytes[:])
			if err != nil {
				writeCompleteChan <- fmt.Errorf("failed to check existence: %v", err)
				return
			}

			if exists {
				// Data is already present, no need to write it again.
				writeCompleteChan <- nil
				return
			}

			err = s.chunkTable.Put(bundleKeyBytes, bundleData)
			if err != nil {
				writeCompleteChan <- fmt.Errorf("failed to put data: %v", err)
				return
			}

			writeCompleteChan <- nil
		}()

		size += uint64(len(bundleKeyBytes) + len(bundleData))
	}

	var failedToWrite bool
	for i := 0; i < len(batchData); i++ {
		err := <-writeCompleteChan
		if err != nil {
			failedToWrite = true
			s.logger.Errorf("failed to write data: %v", err)
		}
	}
	if failedToWrite {
		return 0, fmt.Errorf("failed to write data")
	}

	err := s.chunkTable.Flush()
	if err != nil {
		return 0, fmt.Errorf("failed to flush chunk table: %v", err)
	}

	return size, nil
}

func (s *validatorStore) DeleteKeys(keys []kvstore.Key) error {
	if s.littDB != nil {
		return fmt.Errorf("littDB does not support deletion")
	}

	dbBatch := s.levelDB.NewTTLBatch()
	for _, key := range keys {
		dbBatch.Delete(key)
	}
	return dbBatch.Apply()
}

func (s *validatorStore) GetBundleData(bundleKey []byte) ([]byte, error) {
	if s.littDB == nil {
		// We haven't thrown the switch for littDB yet. Just look in levelDB.
		return s.getChunksLevelDB(bundleKey)
	}

	// Regardless of migration status, always check littDB first.
	data, ok, err := s.getChunksLittDB(bundleKey)
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

		return s.getChunksLevelDB(bundleKey)
	}

	return data, nil
}

func (s *validatorStore) getChunksLevelDB(bundleKey []byte) ([]byte, error) {
	bundlesKeyBuilder, err := s.levelDB.GetKeyBuilder(BundleTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder for bundles: %v", err)
	}

	bundle, err := s.levelDB.Get(bundlesKeyBuilder.Key(bundleKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %v", err)
	}

	return bundle, nil
}

func (s *validatorStore) getChunksLittDB(bundleKey []byte) ([]byte, bool, error) {
	bundle, ok, err := s.chunkTable.Get(bundleKey)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get bundle: %v", err)
	}
	if !ok {
		return nil, false, nil
	}

	return bundle, true, nil
}

// finalizeMigration sleeps until the migration is complete, then deletes the levelDB database.
func (s *validatorStore) finalizeMigration(ctx context.Context) {
	timeUntilMigrationComplete := s.migrationCompleteTime.Sub(s.timeSource())

	select {
	case <-ctx.Done():
		s.logger.Info("context cancelled, migration finalization aborted")
		return
	case <-time.After(timeUntilMigrationComplete):
		s.migrationLock.Lock()
		defer s.migrationLock.Unlock()

		if s.migrationComplete {
			s.logger.Info("migration already completed, nothing to do")
			return
		}
		s.migrationComplete = true

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

func (s *validatorStore) Stop() error {
	if s.littDB != nil {
		err := s.littDB.Close()
		if err != nil {
			return fmt.Errorf("failed to close littDB: %v", err)
		}
	}
	if s.levelDB != nil {
		err := s.levelDB.Shutdown()
		if err != nil {
			return fmt.Errorf("failed to close levelDB: %v", err)
		}
	}

	return nil
}
