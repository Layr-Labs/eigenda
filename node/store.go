package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node/leveldb"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/protobuf/proto"
)

const (
	// How many batches to delete in one atomic operation during the expiration
	// garbage collection.
	numBatchesToDeleteAtomically = 8
)

var ErrBatchAlreadyExist = errors.New("batch already exists")

// Store is a key-value database to store blob data (blob header, blob chunks etc).
type Store struct {
	db     DB
	logger logging.Logger

	blockStaleMeasure   uint32
	storeDurationBlocks uint32

	// The DA Node's metrics.
	metrics *Metrics
}

// NewLevelDBStore creates a new Store object with a db at the provided path and the given logger.
// TODO(jianoaix): parameterize this so we can switch between different database backends.
func NewLevelDBStore(path string, logger logging.Logger, metrics *Metrics, blockStaleMeasure, storeDurationBlocks uint32) (*Store, error) {
	// Create the db at the path. This is currently hardcoded to use
	// levelDB.
	db, err := leveldb.NewLevelDBStore(path)
	if err != nil {
		logger.Error("Could not create leveldb database", "err", err)
		return nil, err
	}

	return &Store{
		db:                  db,
		logger:              logger.With("component", "NodeStore"),
		blockStaleMeasure:   blockStaleMeasure,
		storeDurationBlocks: storeDurationBlocks,
		metrics:             metrics,
	}, nil
}

// Delete expired entries in the store.
// An entry is expired if its expiry <= currentTimeUnixSec, where expiry and
// currentTimeUnixSec are time since Unix epoch (in seconds).
// The deletion of a batch is done atomically, i.e. either all or none entries of a batch will be deleted.
// The function will exit with deadline exceeded error if it cannot finish after timeLimitSec seconds.
// The function returns the number of batches deleted and the status of deletion. Note that the
// number of batches deleted can be positive even if the status is error (e.g. the error happened
// after it had successfully deleted some batches).
func (s *Store) DeleteExpiredEntries(currentTimeUnixSec int64, timeLimitSec uint64) (numBatchesDeleted int, numMappingsDeleted int, numBlobsDeleted int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitSec)*time.Second)
	defer cancel()

	totalBatchesDeleted := 0
	totalMappingsDeleted := 0
	totalBlobsDeleted := 0
	for {
		select {
		case <-ctx.Done():
			return totalBatchesDeleted, totalMappingsDeleted, totalBlobsDeleted, ctx.Err()
		default:
			blobsDeleted, err := s.deleteExpiredBlobs(currentTimeUnixSec, numBatchesToDeleteAtomically)
			if err != nil {
				return totalBatchesDeleted, totalMappingsDeleted, totalBlobsDeleted, err
			}
			totalBlobsDeleted += blobsDeleted

			batchesDeleted, err := s.deleteNBatches(currentTimeUnixSec, numBatchesToDeleteAtomically)
			if err != nil {
				return totalBatchesDeleted, totalMappingsDeleted, totalBlobsDeleted, err
			}
			totalBatchesDeleted += batchesDeleted

			mappingsDeleted, batchesDeleted, err := s.deleteExpiredBatchMapping(currentTimeUnixSec, numBatchesToDeleteAtomically)
			if err != nil {
				return totalBatchesDeleted, totalMappingsDeleted, totalBlobsDeleted, err
			}
			totalMappingsDeleted += mappingsDeleted
			totalBatchesDeleted += batchesDeleted
			// When there is no error and we didn't delete any batch, it means we have
			// no obsolete batches to delete, so we can return.
			if blobsDeleted == 0 && batchesDeleted == 0 && mappingsDeleted == 0 {
				return totalBatchesDeleted, totalMappingsDeleted, totalBlobsDeleted, nil
			}
		}
	}
}

// deleteExpiredBlobs returns the number of blobs deleted and the status of deletion.
// The number is set to -1 (invalid value) if the deletion status is an error.
// Note that the blobs/blob headers expired by this method are those that are not associated with any batch.
// All blobs & blob headers in a batch are expired by deleteNBatches method.
func (s *Store) deleteExpiredBlobs(currentTimeUnixSec int64, numBlobs int) (int, error) {
	// Scan for expired batches.
	iter := s.db.NewIterator(EncodeBlobExpirationKeyPrefix())
	expiredKeys := make([][]byte, 0)
	expiredBlobHeaders := make([][32]byte, 0)
	for iter.Next() {
		ts, err := DecodeBlobExpirationKey(iter.Key())
		if err != nil {
			s.logger.Error("Could not decode the expiration key", "key", iter.Key(), "error", err)
			continue
		}
		// No more rows expired up to current time.
		if currentTimeUnixSec < ts {
			break
		}
		expiredKeys = append(expiredKeys, copyBytes(iter.Key()))
		blobHeaderBytes := copyBytes(iter.Value())
		blobHeaders, err := DecodeHashSlice(blobHeaderBytes)
		if err != nil {
			s.logger.Error("Could not decode the blob header hashes", "error", err)
			continue
		}
		expiredBlobHeaders = append(expiredBlobHeaders, blobHeaders...)
		if len(expiredKeys) == numBlobs {
			break
		}
	}
	iter.Release()

	// No expired batch found.
	if len(expiredKeys) == 0 {
		return 0, nil
	}

	// Calculate the num of bytes (for chunks) that will be purged from the database.
	size := int64(0)
	// Scan for the batch header, blob headers and chunks of each expired batch.
	for _, blobHeaderHash := range expiredBlobHeaders {
		// Blob headers.
		blobHeaderIter := s.db.NewIterator(EncodeBlobHeaderKeyByHash(blobHeaderHash))
		for blobHeaderIter.Next() {
			expiredKeys = append(expiredKeys, copyBytes(blobHeaderIter.Key()))
		}
		blobHeaderIter.Release()

		// Blob chunks.
		blobIter := s.db.NewIterator(EncodeBlobKeyByHashPrefix(blobHeaderHash))
		for blobIter.Next() {
			expiredKeys = append(expiredKeys, copyBytes(blobIter.Key()))
			size += int64(len(blobIter.Value()))
		}
		blobIter.Release()
	}

	// Perform the removal.
	err := s.db.DeleteBatch(expiredKeys)
	if err != nil {
		return -1, fmt.Errorf("failed to delete the expired keys in batch: %w", err)
	}

	// Update the current live metric.
	s.metrics.RemoveNBlobs(len(expiredBlobHeaders), size)

	return len(expiredBlobHeaders), nil
}

// deleteExpiredBatchMapping returns the number of batch to blob index mapping entries deleted and the status of deletion.
// The first return value is the number of batch to blob index mapping entries deleted.
// The second return value is the number of batch header entries deleted.
func (s *Store) deleteExpiredBatchMapping(currentTimeUnixSec int64, numBatches int) (numExpiredMappings int, numExpiredBatches int, err error) {
	// Scan for expired batches.
	iter := s.db.NewIterator(EncodeBatchMappingExpirationKeyPrefix())
	expiredKeys := make([][]byte, 0)
	expiredBatches := make([][]byte, 0)
	for iter.Next() {
		ts, err := DecodeBatchMappingExpirationKey(iter.Key())
		if err != nil {
			s.logger.Error("Could not decode the batch mapping expiration key", "key", iter.Key(), "error", err)
			continue
		}
		// No more rows expired up to current time.
		if currentTimeUnixSec < ts {
			break
		}
		expiredKeys = append(expiredKeys, copyBytes(iter.Key()))
		expiredBatches = append(expiredBatches, copyBytes(iter.Value()))
		if len(expiredKeys) == numBatches {
			break
		}
	}
	iter.Release()

	// No expired batch found.
	if len(expiredKeys) == 0 {
		return 0, 0, nil
	}

	numMappings := 0
	// Scan for the batch header, blob headers and chunks of each expired batch.
	for _, hash := range expiredBatches {
		var batchHeaderHash [32]byte
		copy(batchHeaderHash[:], hash)

		// Batch header.
		expiredKeys = append(expiredKeys, EncodeBatchHeaderKey(batchHeaderHash))

		// Blob index mapping.
		blobIndexIter := s.db.NewIterator(EncodeBlobIndexKeyPrefix(batchHeaderHash))
		for blobIndexIter.Next() {
			expiredKeys = append(expiredKeys, copyBytes(blobIndexIter.Key()))
			numMappings++
		}
		blobIndexIter.Release()
	}

	// Perform the removal.
	err = s.db.DeleteBatch(expiredKeys)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to delete the expired keys in batch: %w", err)
	}

	s.logger.Info("Deleted expired batch mapping", "numBatches", len(expiredBatches), "numMappings", numMappings)
	numExpiredMappings = numMappings
	numExpiredBatches = len(expiredBatches)
	return numExpiredMappings, numExpiredBatches, nil
}

// deleteNBatches returns the number of batches we deleted and the status of deletion. The number
// is set to -1 (invalid value) if the deletion status is an error.
func (s *Store) deleteNBatches(currentTimeUnixSec int64, numBatches int) (int, error) {
	// Scan for expired batches.
	iter := s.db.NewIterator(EncodeBatchExpirationKeyPrefix())
	expiredKeys := make([][]byte, 0)
	expiredBatches := make([][]byte, 0)
	for iter.Next() {
		ts, err := DecodeBatchExpirationKey(iter.Key())
		if err != nil {
			s.logger.Error("Could not decode the expiration key", "key:", iter.Key(), "error", err)
			continue
		}
		// No more rows expired up to current time.
		if currentTimeUnixSec < ts {
			break
		}
		expiredKeys = append(expiredKeys, copyBytes(iter.Key()))
		expiredBatches = append(expiredBatches, copyBytes(iter.Value()))
		if len(expiredKeys) == numBatches {
			break
		}
	}
	iter.Release()

	// No expired batch found.
	if len(expiredKeys) == 0 {
		return 0, nil
	}

	// Calculate the num of bytes (for chunks) that will be purged from the database.
	size := int64(0)
	numBlobs := 0
	// Scan for the batch header, blob headers and chunks of each expired batch.
	for _, hash := range expiredBatches {
		var batchHeaderHash [32]byte
		copy(batchHeaderHash[:], hash)

		// Batch header.
		expiredKeys = append(expiredKeys, EncodeBatchHeaderKey(batchHeaderHash))

		// Blob headers.
		blobHeaderIter := s.db.NewIterator(EncodeBlobHeaderKeyPrefix(batchHeaderHash))
		for blobHeaderIter.Next() {
			expiredKeys = append(expiredKeys, copyBytes(blobHeaderIter.Key()))
			numBlobs++
		}
		blobHeaderIter.Release()

		// Blob chunks.
		blobIter := s.db.NewIterator(bytes.NewBuffer(hash).Bytes())
		for blobIter.Next() {
			expiredKeys = append(expiredKeys, copyBytes(blobIter.Key()))
			size += int64(len(blobIter.Value()))
		}
		blobIter.Release()
	}

	// Perform the removal.
	err := s.db.DeleteBatch(expiredKeys)
	if err != nil {
		s.logger.Error("Failed to delete the expired keys in batch", "keys", expiredKeys, "error", err)
		return -1, err
	}

	// Update the current live batch metric.
	s.metrics.RemoveNCurrentBatch(len(expiredBatches), size)
	s.metrics.RemoveNBlobs(numBlobs, 0)

	return len(expiredBatches), nil
}

// Store the batch into the store.
//
// The batch will be itemized into multiple entries when it's stored:
//   - Batch header: keyed by <batchHeaderPrefix, batchHeaderHash>
//   - Batch expiry: keyed by <batchExprationPrefix, expirationTime>
//   - The header of each blob in the batch: one entry to each blob header, keyed by <blobHeaderPrefix, batchHeaderHash, blobIdx>
//   - The chunks of each blob in the batch: one entry for each blob chunks, keyed by <batchHeaderHash, blobIdx, quorumID>
//
// These entries will be stored atomically, i.e. either all or none entries will be stored.
func (s *Store) StoreBatch(ctx context.Context, header *core.BatchHeader, blobs []*core.BlobMessage, blobsProto []*node.Blob) (*[][]byte, error) {
	storeBatchStart := time.Now()

	log := s.logger
	batchHeaderHash, err := header.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}

	// The key/value pairs that need to be written to the local database.
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	// Generate the key/value pair for batch header.
	batchHeaderKey := EncodeBatchHeaderKey(batchHeaderHash)
	batchHeaderBytes, err := header.Serialize()
	if err != nil {
		log.Error("Cannot serialize the batch header:", "err", err)
		return nil, err
	}

	// If the batch header exists already in store, we know that all data items associated
	// with this batch should be in the store already (because they are written atomically).
	// In this case, we do nothing and just return.
	if s.HasKey(ctx, batchHeaderKey) {
		return nil, ErrBatchAlreadyExist
	}

	keys = append(keys, batchHeaderKey)
	values = append(values, batchHeaderBytes)

	expirationTime := s.expirationTime()
	expirationKey := EncodeBatchExpirationKey(expirationTime)
	keys = append(keys, expirationKey)
	values = append(values, batchHeaderHash[:])

	// Generate key/value pairs for all blob headers and blob chunks .
	size := int64(0)
	var serializationDuration, encodingDuration time.Duration
	for idx, blob := range blobs {
		// blob header
		blobHeaderKey, err := EncodeBlobHeaderKey(batchHeaderHash, idx)
		if err != nil {
			log.Error("Cannot generate the key for storing blob header:", "err", err)
			return nil, err
		}
		blobHeaderBytes, err := proto.Marshal(blobsProto[idx].GetHeader())
		if err != nil {
			log.Error("Cannot serialize the blob header proto:", "err", err)
			return nil, err
		}
		keys = append(keys, blobHeaderKey)
		values = append(values, blobHeaderBytes)

		// Get raw chunks
		start := time.Now()
		rawBlob := blobsProto[idx]
		if len(rawBlob.GetBundles()) != len(blob.Bundles) {
			return nil, errors.New("internal error: the number of bundles in parsed blob must be the same as in raw blob")
		}
		format := GetBundleEncodingFormat(rawBlob)
		rawBundles := make(map[core.QuorumID][]byte)
		rawChunks := make(map[core.QuorumID][][]byte)
		for i, bundle := range rawBlob.GetBundles() {
			quorumID := uint8(rawBlob.GetHeader().GetQuorumHeaders()[i].GetQuorumId())
			if format == core.GnarkBundleEncodingFormat {
				if len(bundle.GetChunks()) > 0 && len(bundle.GetChunks()[0]) > 0 {
					return nil, errors.New("chunks of a bundle are encoded together already")
				}
				rawBundles[quorumID] = bundle.GetBundle()
			} else {
				rawChunks[quorumID] = make([][]byte, len(bundle.GetChunks()))
				for j, chunk := range bundle.GetChunks() {
					rawChunks[quorumID][j] = chunk
				}
			}
		}
		serializationDuration += time.Since(start)
		start = time.Now()
		// blob chunks
		for quorumID, bundle := range blob.Bundles {
			key, err := EncodeBlobKey(batchHeaderHash, idx, quorumID)
			if err != nil {
				log.Error("Cannot generate the key for storing blob:", "err", err)
				return nil, err
			}

			if format == core.GnarkBundleEncodingFormat {
				rawBundle, ok := rawBundles[quorumID]
				if ok {
					size += int64(len(rawBundle))
					keys = append(keys, key)
					values = append(values, rawBundle)
				}
			} else if format == core.GobBundleEncodingFormat {
				if len(rawChunks[quorumID]) != len(bundle) {
					return nil, errors.New("internal error: the number of chunks in parsed blob bundle must be the same as in raw blob bundle")
				}
				chunksBytes, ok := rawChunks[quorumID]
				if ok {

					bundleRaw := make([][]byte, len(bundle))
					for i := 0; i < len(bundle); i++ {
						bundleRaw[i] = chunksBytes[i]
					}
					chunkBytes, err := EncodeChunks(bundleRaw)
					if err != nil {
						return nil, err
					}
					size += int64(len(chunkBytes))
					keys = append(keys, key)
					values = append(values, chunkBytes)
				}
			} else {
				return nil, fmt.Errorf("invalid bundle encoding format: %d", format)
			}
		}
		encodingDuration += time.Since(start)
	}

	start := time.Now()
	// Write all the key/value pairs to the local database atomically.
	err = s.db.WriteBatch(keys, values)
	if err != nil {
		log.Error("Failed to write the batch into local database:", "err", err)
		return nil, err
	}
	throughput := float64(size) / time.Since(start).Seconds()
	s.metrics.DBWriteThroughput.Set(throughput)
	log.Debug("StoreBatch succeeded", "chunk serialization duration", serializationDuration, "bytes encoding duration", encodingDuration, "num blobs", len(blobs), "num of key-value pair entries", len(keys), "write batch duration", time.Since(start), "write throughput (MB/s)", throughput/1000_000, "total store batch duration", time.Since(storeBatchStart), "total bytes", size)

	return &keys, nil
}

func (s *Store) StoreBlobs(ctx context.Context, blobs []*core.BlobMessage, blobsProto []*node.Blob) (*[][]byte, error) {
	storeBlobsStart := time.Now()
	// The key/value pairs that need to be written to the local database.
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	expirationTime := s.expirationTime()
	expirationKey := EncodeBlobExpirationKey(expirationTime)
	// expirationValue is a list of blob header hashes that are expired.
	expirationValue := make([]byte, 0)
	var err error
	// If there is already an expiration key in the store, we need to get the value and append to it.
	if s.HasKey(ctx, expirationKey) {
		expirationValue, err = s.db.Get(expirationKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get the expiration value: %w", err)
		}
	}

	// Generate key/value pairs for all blob headers and blob chunks .
	size := int64(0)
	var serializationDuration, encodingDuration time.Duration
	for idx, blob := range blobs {
		rawBlob := blobsProto[idx]
		if len(rawBlob.GetBundles()) != len(blob.Bundles) {
			return nil, fmt.Errorf("internal error: the number of bundles in parsed blob (%d) must be the same as in raw blob (%d)", len(rawBlob.GetBundles()), len(blob.Bundles))
		}

		// blob header
		blobHeaderHash, err := blob.BlobHeader.GetBlobHeaderHash()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob header hash: %w", err)
		}
		blobHeaderKey := EncodeBlobHeaderKeyByHash(blobHeaderHash)
		if s.HasKey(ctx, blobHeaderKey) {
			s.logger.Warn("Blob already exists", "blobHeaderHash", hexutil.Encode(blobHeaderHash[:]))
			continue
		}

		blobHeaderBytes, err := proto.Marshal(blobsProto[idx].GetHeader())
		if err != nil {
			return nil, fmt.Errorf("failed to serialize the blob header proto: %w", err)
		}
		keys = append(keys, blobHeaderKey)
		values = append(values, blobHeaderBytes)
		expirationValue = append(expirationValue, blobHeaderHash[:]...)

		// Get raw chunks
		start := time.Now()
		format := GetBundleEncodingFormat(rawBlob)
		rawBundles := make(map[core.QuorumID][]byte)
		rawChunks := make(map[core.QuorumID][][]byte)
		for i, bundle := range rawBlob.GetBundles() {
			quorumID := uint8(rawBlob.GetHeader().GetQuorumHeaders()[i].GetQuorumId())
			if format == core.GnarkBundleEncodingFormat {
				if len(bundle.GetChunks()) > 0 && len(bundle.GetChunks()[0]) > 0 {
					return nil, errors.New("chunks of a bundle are encoded together already")
				}
				rawBundles[quorumID] = bundle.GetBundle()
			} else {
				rawChunks[quorumID] = make([][]byte, len(bundle.GetChunks()))
				for j, chunk := range bundle.GetChunks() {
					rawChunks[quorumID][j] = chunk
				}
			}
		}
		serializationDuration += time.Since(start)
		start = time.Now()
		// blob chunks
		for quorumID, bundle := range blob.Bundles {
			key, err := EncodeBlobKeyByHash(blobHeaderHash, quorumID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate the key for storing blob: %w", err)
			}

			if format == core.GnarkBundleEncodingFormat {
				rawBundle, ok := rawBundles[quorumID]
				if ok {
					size += int64(len(rawBundle))
					keys = append(keys, key)
					values = append(values, rawBundle)
				}
			} else if format == core.GobBundleEncodingFormat {
				if len(rawChunks[quorumID]) != len(bundle) {
					return nil, errors.New("internal error: the number of chunks in parsed blob bundle must be the same as in raw blob bundle")
				}
				chunksBytes, ok := rawChunks[quorumID]
				if ok {

					bundleRaw := make([][]byte, len(bundle))
					for i := 0; i < len(bundle); i++ {
						bundleRaw[i] = chunksBytes[i]
					}
					chunkBytes, err := EncodeChunks(bundleRaw)
					if err != nil {
						return nil, err
					}
					size += int64(len(chunkBytes))
					keys = append(keys, key)
					values = append(values, chunkBytes)
				}
			} else {
				return nil, fmt.Errorf("invalid bundle encoding format: %d", format)
			}
		}
		encodingDuration += time.Since(start)
	}

	keys = append(keys, expirationKey)
	values = append(values, expirationValue)

	start := time.Now()
	// Write all the key/value pairs to the local database atomically.
	err = s.db.WriteBatch(keys, values)
	if err != nil {
		return nil, fmt.Errorf("failed to write the batch into local database: %w", err)
	}
	throughput := float64(size) / time.Since(start).Seconds()
	s.metrics.DBWriteThroughput.Set(throughput)
	s.logger.Debug("StoreBlobs succeeded", "chunk serialization duration", serializationDuration, "bytes encoding duration", encodingDuration, "num blobs", len(blobs), "num of key-value pair entries", len(keys), "write batch duration", time.Since(start), "write throughput (MB/s)", throughput/1000_000, "total store batch duration", time.Since(storeBlobsStart), "total bytes", size)

	return &keys, nil
}

func (s *Store) StoreBatchBlobMapping(ctx context.Context, batchHeader *core.BatchHeader, blobHeaderHashes [][32]byte) error {
	start := time.Now()
	// The key/value pairs that need to be written to the local database.
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
	if err != nil {
		return fmt.Errorf("failed to get the batch header hash: %w", err)
	}

	expirationTime := s.expirationTime()
	expirationKey := EncodeBatchMappingExpirationKey(expirationTime, batchHeaderHash)
	keys = append(keys, expirationKey)
	values = append(values, batchHeaderHash[:])

	for blobIndex, blobHeaderHash := range blobHeaderHashes {
		blobIndexKey := EncodeBlobIndexKey(batchHeaderHash, blobIndex)
		keys = append(keys, blobIndexKey)
		values = append(values, copyBytes(blobHeaderHash[:]))
	}

	// Generate the key/value pair for batch header.
	batchHeaderKey := EncodeBatchHeaderKey(batchHeaderHash)
	batchHeaderBytes, err := batchHeader.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize the batch header: %w", err)
	}
	keys = append(keys, batchHeaderKey)
	values = append(values, batchHeaderBytes)

	err = s.db.WriteBatch(keys, values)
	if err != nil {
		return fmt.Errorf("failed to write the blob index mappings into local database: %w", err)
	}
	s.logger.Debug("StoreBatchBlobMapping succeeded", "duration", time.Since(start))
	return nil
}

func (s *Store) expirationTime() int64 {
	// Setting the expiration time for the batch.
	curr := time.Now().Unix()
	timeToExpire := (s.blockStaleMeasure + s.storeDurationBlocks) * 12 // 12s per block
	// Why this expiration time is safe?
	//
	// The batch must be confirmed before referenceBlockNumber+blockStaleMeasure, otherwise
	// it's stale and won't be accepted onchain. This means the blob's lifecycle will end
	// before referenceBlockNumber+blockStaleMeasure+storeDurationBlocks.
	// Since time@referenceBlockNumber < time.Now() (we always use a reference block that's
	// already onchain), we have
	// time@(referenceBlockNumber+blockStaleMeasure+storeDurationBlocks)
	// = time@referenceBlockNumber + 12*(blockStaleMeasure+storeDurationBlocks)
	// < time.Now() + 12*(blockStaleMeasure+storeDurationBlocks).
	//
	// Note if a batch is unconfirmed, it could be removed even earlier; here we treat its
	// lifecycle the same as confirmed batches for simplicity.
	return curr + int64(timeToExpire)
}

// GetBatchHeader returns the batch header for the given batchHeaderHash.
func (s *Store) GetBatchHeader(ctx context.Context, batchHeaderHash [32]byte) ([]byte, error) {
	batchHeaderKey := EncodeBatchHeaderKey(batchHeaderHash)
	data, err := s.db.Get(batchHeaderKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
}

// GetBlobHeader returns the blob header for the given batchHeaderHash, blob index.
func (s *Store) GetBlobHeader(ctx context.Context, batchHeaderHash [32]byte, blobIndex int) ([]byte, error) {
	blobHeaderKey, err := EncodeBlobHeaderKey(batchHeaderHash, blobIndex)
	if err != nil {
		return nil, err
	}
	data, err := s.db.Get(blobHeaderKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
}

// GetBlobHeaderByHeaderHash returns the blob header for the given blobHeaderHash.
func (s *Store) GetBlobHeaderByHeaderHash(ctx context.Context, blobHeaderHash [32]byte) ([]byte, error) {
	blobHeaderKey := EncodeBlobHeaderKeyByHash(blobHeaderHash)
	data, err := s.db.Get(blobHeaderKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
}

// GetChunks returns the list of byte arrays stored for given blobKey along with the encoding
// format of the bytes.
func (s *Store) GetChunks(ctx context.Context, batchHeaderHash [32]byte, blobIndex int, quorumID core.QuorumID) ([][]byte, node.ChunkEncodingFormat, error) {
	log := s.logger

	blobKey, err := EncodeBlobKey(batchHeaderHash, blobIndex, quorumID)
	if err != nil {
		return nil, node.ChunkEncodingFormat_UNKNOWN, err
	}
	data, err := s.db.Get(blobKey)
	if err != nil {
		return nil, node.ChunkEncodingFormat_UNKNOWN, err
	}

	chunks, format, err := DecodeChunks(data)
	if err != nil {
		return nil, format, err
	}
	log.Debug("Retrieved chunk", "blobKey", hexutil.Encode(blobKey), "length", len(data), "chunk encoding format", format)

	return chunks, format, nil
}

func (s *Store) GetBlobHeaderHashAtIndex(ctx context.Context, batchHeaderHash [32]byte, blobIndex int) ([]byte, error) {
	blobIndexKey := EncodeBlobIndexKey(batchHeaderHash, blobIndex)
	data, err := s.db.Get(blobIndexKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
}

// HasKey returns if a given key has been stored.
func (s *Store) HasKey(ctx context.Context, key []byte) bool {
	_, err := s.db.Get(key)
	return err == nil
}

// DeleteKeys removes a list of keys from the store atomically.
//
// Note: caller should ensure these keys are exactly all the data items for a single batch
// to maintain the integrity of the store.
func (s *Store) DeleteKeys(ctx context.Context, keys *[][]byte) error {
	return s.db.DeleteBatch(*keys)
}

// Flattens an array of byte arrays (chunks) into a single byte array
//
// EncodeChunks(chunks) = (len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)
func EncodeChunks(chunks [][]byte) ([]byte, error) {
	totalSize := 0
	for _, chunk := range chunks {
		totalSize += len(chunk) + 8 // Add size of uint64 for length
	}
	result := make([]byte, totalSize)
	buf := result
	for _, chunk := range chunks {
		binary.LittleEndian.PutUint64(buf, uint64(len(chunk)))
		buf = buf[8:]
		copy(buf, chunk)
		buf = buf[len(chunk):]
	}
	return result, nil
}

func DecodeGnarkChunks(data []byte) ([][]byte, error) {
	format, chunkLen, err := parseHeader(data)
	if err != nil {
		return nil, err
	}
	if format != core.GnarkBundleEncodingFormat {
		return nil, errors.New("invalid bundle data encoding format")
	}
	if chunkLen == 0 {
		return nil, errors.New("chunk length must be greater than zero")
	}
	chunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*int(chunkLen)
	chunks := make([][]byte, 0)
	buf := data[8:]
	for len(buf) > 0 {
		if len(buf) < chunkSize {
			return nil, errors.New("invalid data to decode")
		}
		chunks = append(chunks, buf[:chunkSize])
		buf = buf[chunkSize:]
	}
	return chunks, nil
}

// DecodeChunks((len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)) = chunks
func DecodeGobChunks(data []byte) ([][]byte, error) {
	format, chunkLen, err := parseHeader(data)
	if err != nil {
		return nil, err
	}
	if format != core.GobBundleEncodingFormat {
		return nil, errors.New("invalid bundle data encoding format")
	}
	if chunkLen == 0 {
		return nil, errors.New("chunk length must be greater than zero")
	}
	chunks := make([][]byte, 0)
	buf := data
	for len(buf) > 0 {
		if len(buf) < 8 {
			return nil, errors.New("invalid data to decode")
		}
		chunkSize := binary.LittleEndian.Uint64(buf)
		buf = buf[8:]

		if len(buf) < int(chunkSize) {
			return nil, errors.New("invalid data to decode")
		}
		chunks = append(chunks, buf[:chunkSize])
		buf = buf[chunkSize:]
	}
	return chunks, nil
}

// parseHeader parses the header and returns the encoding format and the chunk length.
func parseHeader(data []byte) (core.BundleEncodingFormat, uint64, error) {
	if len(data) < 8 {
		return 0, 0, errors.New("no header found, the data size is less 8 bytes")
	}
	meta := binary.LittleEndian.Uint64(data)
	format := binary.LittleEndian.Uint64(data) >> (core.NumBundleHeaderBits - core.NumBundleEncodingFormatBits)
	chunkLen := (meta << core.NumBundleEncodingFormatBits) >> core.NumBundleEncodingFormatBits
	return uint8(format), chunkLen, nil
}

// DecodeChunks converts a flattened array of chunks into an array of its constituent chunks,
// throwing an error in case the chunks were not serialized correctly.
func DecodeChunks(data []byte) ([][]byte, node.ChunkEncodingFormat, error) {
	// Empty chunk is valid, but there is nothing to decode.
	if len(data) == 0 {
		return [][]byte{}, node.ChunkEncodingFormat_UNKNOWN, nil
	}
	format, _, err := parseHeader(data)
	if err != nil {
		return nil, node.ChunkEncodingFormat_UNKNOWN, err
	}

	// Note: the encoding format IDs may not be the same as the field ID in protobuf.
	// For example, GobBundleEncodingFormat is 1 but node.ChunkEncodingFormat_GOB has proto
	// field ID 2.
	switch format {
	case 0:
		chunks, err := DecodeGobChunks(data)
		return chunks, node.ChunkEncodingFormat_GOB, err
	case 1:
		chunks, err := DecodeGnarkChunks(data)
		return chunks, node.ChunkEncodingFormat_GNARK, err
	default:
		return nil, node.ChunkEncodingFormat_UNKNOWN, errors.New("invalid data encoding format")
	}
}

func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
