package node

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/memory"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

const (
	// The name of the littDB table containing chunk data.
	chunksTableName = "chunks"
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
	StoreBatch(batchData []*BundleToStore) (uint64, error)

	// GetBundleData returns the chunks of a blob with the given bundle key.
	// The returned chunks are encoded in bundle format.
	GetBundleData(bundleKey []byte) ([]byte, error)

	// Stop stops the store.
	Stop() error
}

type validatorStore struct {
	logger     logging.Logger
	timeSource func() time.Time

	// The littDB database for storing chunk data. If nil, then the store has not yet been migrated to littDB.
	littDB litt.DB

	// The table where chunks are stored in the littDB database.
	chunkTable litt.Table

	// The length of time to store data in the database.
	ttl time.Duration

	// A lock used to prevent concurrent requests from storing the same data multiple times.
	duplicateRequestLock *common.IndexLock

	// The salt used to prevent an attacker from causing hash collisions in the duplicate request lock.
	duplicateRequestSalt [16]byte

	// limits the frequency of hot reads (i.e. reads that hit the cache)
	hotReadRateLimiter *rate.Limiter

	// limits the frequency of cold reads (i.e. reads that miss the cache)
	coldReadRateLimiter *rate.Limiter
}

var _ ValidatorStore = &validatorStore{}

func NewValidatorStore(
	logger logging.Logger,
	config *Config,
	timeSource func() time.Time,
	ttl time.Duration,
	registry *prometheus.Registry) (ValidatorStore, error) {

	if len(config.LittDBStoragePaths) == 0 {
		return nil, fmt.Errorf("no littDB paths provided")
	}

	stringBuilder := strings.Builder{}
	if len(config.LittDBStoragePaths) > 1 {
		stringBuilder.WriteString("s")
	}
	for i, path := range config.LittDBStoragePaths {
		stringBuilder.WriteString(" ")
		stringBuilder.WriteString(path)
		if i < len(config.LittDBStoragePaths)-1 {
			stringBuilder.WriteString(",")
		}
	}
	logger.Infof("Using littDB at path%s", stringBuilder.String())

	littConfig, err := litt.DefaultConfig(config.LittDBStoragePaths...)
	littConfig.ShardingFactor = 1
	littConfig.MetricsEnabled = true
	littConfig.MetricsRegistry = registry
	littConfig.MetricsNamespace = littDBMetricsPrefix
	littConfig.Logger = logger
	littConfig.DoubleWriteProtection = config.LittDBDoubleWriteProtection
	if err != nil {
		return nil, fmt.Errorf("failed to create new litt config: %w", err)
	}

	littDB, err := littbuilder.NewDB(littConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new litt store: %w", err)
	}

	chunkTable, err := littDB.GetTable(chunksTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks table: %w", err)
	}

	maxMemory, err := memory.GetMaximumAvailableMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get maximum available memory: %w", err)
	}

	writeCacheSize := uint64(0)
	if config.LittDBWriteCacheSizeGB > 0 {
		writeCacheSize = uint64(config.LittDBWriteCacheSizeGB * units.GiB)
		logger.Infof("LittDB write cache size configured to use %.2f GB.\n", config.LittDBWriteCacheSizeGB)
	} else {
		writeCacheSize = uint64(config.LittDBWriteCacheSizeFraction * float64(maxMemory))
		logger.Infof("LittDB write cache is configured to use %.1f%% of %.2f GB available (%.2f GB).",
			config.LittDBWriteCacheSizeFraction*100.0,
			float64(maxMemory)/float64(units.GiB),
			float64(writeCacheSize)/float64(units.GiB))
	}

	readCacheSize := uint64(0)
	if config.LittDBReadCacheSizeGB > 0 {
		readCacheSize = uint64(config.LittDBReadCacheSizeGB * units.GiB)
		logger.Infof("LittDB read cache size configured to use %.2f GB.\n", config.LittDBReadCacheSizeGB)
	} else {
		readCacheSize = uint64(config.LittDBReadCacheSizeFraction * float64(maxMemory))
		logger.Infof("LittDB read cache is configured to use %.1f%% of %.2f GB available (%.2f GB).",
			config.LittDBReadCacheSizeFraction*100.0,
			float64(maxMemory)/float64(units.GiB),
			float64(readCacheSize)/float64(units.GiB))
	}

	if writeCacheSize+readCacheSize >= maxMemory {
		return nil, fmt.Errorf("Write cache size + read cache size must be less than max memory. "+
			"Write cache size: %d, read cache size: %d, max memory: %d", writeCacheSize, readCacheSize, maxMemory)
	}

	err = chunkTable.SetWriteCacheSize(writeCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to set write cache size for chunks table: %w", err)
	}

	err = chunkTable.SetReadCacheSize(readCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to set read cache size for chunks table: %w", err)
	}

	err = chunkTable.SetTTL(ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to set TTL for chunks table: %w", err)
	}

	salt := [16]byte{}
	_, err = rand.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %v", err)
	}

	hotReadRateLimiter := rate.NewLimiter(
		rate.Limit(config.GetChunksHotCacheReadLimitMB*units.MiB),
		int(config.GetChunksHotBurstLimitMB*units.MiB))
	coldReadRateLimiter := rate.NewLimiter(
		rate.Limit(config.GetChunksColdCacheReadLimitMB*units.MiB),
		int(config.GetChunksColdBurstLimitMB*units.MiB))

	store := &validatorStore{
		logger:               logger,
		timeSource:           timeSource,
		littDB:               littDB,
		chunkTable:           chunkTable,
		ttl:                  ttl,
		duplicateRequestLock: common.NewIndexLock(1024),
		duplicateRequestSalt: salt,
		hotReadRateLimiter:   hotReadRateLimiter,
		coldReadRateLimiter:  coldReadRateLimiter,
	}

	return store, nil
}

func (s *validatorStore) StoreBatch(batchData []*BundleToStore) (uint64, error) {
	if len(batchData) == 0 {
		return 0, fmt.Errorf("no batch data")
	}

	var size uint64

	writeCompleteChan := make(chan error, len(batchData))
	for _, batchDatum := range batchData {
		bundleKeyBytes := batchDatum.BundleKey
		bundleData := batchDatum.BundleBytes

		go func() {
			// Grab a lock on the hash of the blob. This protects against duplicate writes of the same blob.
			hash := util.HashKey(bundleKeyBytes[:], s.duplicateRequestSalt)
			lockIndex := uint64(hash)
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

func (s *validatorStore) GetBundleData(bundleKey []byte) ([]byte, error) {

	// Regardless of migration status, always check littDB first.
	data, exists, err := s.getChunksLittDB(bundleKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks: %v", err)
	}

	if !exists {
		return nil, fmt.Errorf("failed to get chunks: not found")
	}

	return data, nil
}

func (s *validatorStore) getChunksLittDB(bundleKey []byte) (data []byte, exists bool, err error) {

	hotReadsExhausted := s.hotReadRateLimiter.Tokens() <= 0
	if hotReadsExhausted {
		// If hot reads are exhausted we do not allow cold reads either.
		return nil, false, fmt.Errorf("read rate limit exhausted")
	}

	coldReadsExhausted := s.coldReadRateLimiter.Tokens() <= 0

	bundle, exists, hot, err := s.chunkTable.CacheAwareGet(bundleKey, coldReadsExhausted)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get bundle: %v", err)
	}
	if exists && bundle == nil {
		// This can happen when the data is on disk but we've exhausted the cold read rate
		return nil, false, fmt.Errorf("cold read rate limit exhausted")
	}
	if !exists {
		return nil, false, nil
	}

	// If we read the value, debit the rate limiters. This may cause us to exceed the rate limit, in which
	// case the number of tokens will be negative. When this happens, we will not be able to read until
	// we accumulate enough tokens to "pay off the debt".
	if hot {
		s.hotReadRateLimiter.ReserveN(time.Now(), len(bundleKey)+len(bundle))
	} else {
		s.coldReadRateLimiter.ReserveN(time.Now(), len(bundleKey)+len(bundle))
	}

	return bundle, true, nil
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

	return nil
}
