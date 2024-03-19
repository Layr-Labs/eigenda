package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/node/leveldb"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
		logger:              logger,
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
func (s *Store) DeleteExpiredEntries(currentTimeUnixSec int64, timeLimitSec uint64) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitSec)*time.Second)
	defer cancel()

	numBatchesDeleted := 0
	for {
		select {
		case <-ctx.Done():
			return numBatchesDeleted, ctx.Err()
		default:
			numDeleted, err := s.deleteNBatches(currentTimeUnixSec, numBatchesToDeleteAtomically)
			if err != nil {
				return numBatchesDeleted, err
			}
			// When there is no error and we didn't delete any batch, it means we have
			// no obsolete batches to delete, so we can return.
			if numDeleted == 0 {
				return numBatchesDeleted, nil
			}
			numBatchesDeleted += numDeleted
		}
	}
}

// Returns the number of batches we deleted and the status of deletion. The number
// is set to -1 (invalid value) if the deletion status is an error.
func (s *Store) deleteNBatches(currentTimeUnixSec int64, numBatches int) (int, error) {
	// Scan for expired batches.
	iter := s.db.NewIterator(EncodeBatchExpirationKeyPrefix())
	expiredKeys := make([][]byte, 0)
	expiredBatches := make([][]byte, 0)
	for iter.Next() {
		ts, err := DecodeBatchExpirationKey(iter.Key())
		if err != nil {
			s.logger.Error("Could not decode the expiration key", "key:", iter.Key(), "error:", err)
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
		s.logger.Error("Failed to delete the expired keys in batch", "keys:", expiredKeys, "error:", err)
		return -1, err
	}

	// Update the current live batch metric.
	s.metrics.RemoveNCurrentBatch(len(expiredBatches), size)

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
	expirationTime := curr + int64(timeToExpire)
	expirationKey := EncodeBatchExpirationKey(expirationTime)
	keys = append(keys, expirationKey)
	values = append(values, batchHeaderHash[:])

	// Generate key/value pairs for all blob headers and blob chunks .
	size := int64(0)
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

		// blob chunks
		for quorumID, bundle := range blob.Bundles {
			key, err := EncodeBlobKey(batchHeaderHash, idx, quorumID)
			if err != nil {
				log.Error("Cannot generate the key for storing blob:", "err", err)
				return nil, err
			}

			bundleRaw := make([][]byte, len(bundle))
			for i, chunk := range bundle {
				bundleRaw[i], err = chunk.Serialize()
				if err != nil {
					log.Error("Cannot serialize chunk:", "err", err)
					return nil, err
				}
			}
			chunkBytes, err := encodeChunks(bundleRaw)
			if err != nil {
				return nil, err
			}
			size += int64(len(chunkBytes))

			keys = append(keys, key)
			values = append(values, chunkBytes)

		}
	}

	// Write all the key/value pairs to the local database atomically.
	err = s.db.WriteBatch(keys, values)
	if err != nil {
		log.Error("Failed to write the batch into local database:", "err", err)
		return nil, err
	}

	return &keys, nil
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

// GetChunks returns the list of byte arrays stored for given blobKey along with a boolean
// indicating if the read was unsuccessful or the chunks were serialized correctly
func (s *Store) GetChunks(ctx context.Context, batchHeaderHash [32]byte, blobIndex int, quorumID core.QuorumID) ([][]byte, bool) {
	log := s.logger

	blobKey, err := EncodeBlobKey(batchHeaderHash, blobIndex, quorumID)
	if err != nil {
		return nil, false
	}
	data, err := s.db.Get(blobKey)
	if err != nil {
		return nil, false
	}
	log.Debug("Retrieved chunk", "blobKey", hexutil.Encode(blobKey), "length", len(data))

	chunks, err := decodeChunks(data)
	if err != nil {
		return nil, false
	}
	return chunks, true
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
// encodeChunks(chunks) = (len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)
func encodeChunks(chunks [][]byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, chunk := range chunks {
		if err := binary.Write(buf, binary.LittleEndian, uint64(len(chunk))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(chunk); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// Converts a flattened array of chunks into an array of its constituent chunks,
// throwing an error in case the chunks were not serialized correctly
//
// decodeChunks((len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)) = chunks
func decodeChunks(data []byte) ([][]byte, error) {
	buf := bytes.NewReader(data)
	chunks := make([][]byte, 0)

	for {
		var length uint64
		err := binary.Read(buf, binary.LittleEndian, &length)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk := make([]byte, length)
		_, err = buf.Read(chunk)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, chunk)
		if buf.Len() < 8 {
			break
		}
	}

	return chunks, nil
}

func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
