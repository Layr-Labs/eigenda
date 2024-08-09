package node

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/Layr-Labs/eigenda/core"
)

const (
	// Caution: the change to these prefixes needs to handle the backward compatibility,
	// making sure the new code work with old data in DA Node store.
	blobHeaderPrefix      = "_BLOB_HEADER_"  // The prefix of the blob header key.
	batchHeaderPrefix     = "_BATCH_HEADER_" // The prefix of the batch header key.
	batchExpirationPrefix = "_EXPIRATION_"   // The prefix of the batch expiration key.
	// blobExpirationPrefix is the prefix of the blob and blob header expiration key.
	// The blobs/blob headers expired by this prefix are those that are not associated with any batch.
	// All blobs/blob headers in a batch are expired by the batch expiration key above.
	blobExpirationPrefix = "_BLOBEXPIRATION_"
	// batchMappingExpirationPrefix is the prefix of the batch mapping expiration key.
	// This key is used to expire the batch to blob index mapping used to identify blob index in a full batch.
	batchMappingExpirationPrefix = "_BATCHEXPIRATION_"
	blobPrefix                   = "_BLOB_"      // The prefix of the blob key.
	blobIndexPrefix              = "_BLOB_INDEX" // The prefix of the blob index key.
)

// EncodeBlobKey returns an encoded key as blob identification.
func EncodeBlobKey(batchHeaderHash [32]byte, blobIndex int, quorumID core.QuorumID) ([]byte, error) {
	buf := bytes.NewBuffer(batchHeaderHash[:])
	err := binary.Write(buf, binary.LittleEndian, int32(blobIndex))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, quorumID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeBlobKeyByHash(blobHeaderHash [32]byte, quorumID core.QuorumID) ([]byte, error) {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, blobHeaderHash[:]...))
	err := binary.Write(buf, binary.LittleEndian, quorumID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeBlobKeyByHashPrefix(blobHeaderHash [32]byte) []byte {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, blobHeaderHash[:]...))
	return buf.Bytes()
}

// EncodeBlobHeaderKey returns an encoded key as blob header identification.
func EncodeBlobHeaderKey(batchHeaderHash [32]byte, blobIndex int) ([]byte, error) {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	err := binary.Write(buf, binary.LittleEndian, int32(blobIndex))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Returns an encoded prefix of blob header key.
func EncodeBlobHeaderKeyPrefix(batchHeaderHash [32]byte) []byte {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	return buf.Bytes()
}

func EncodeBlobHeaderKeyByHash(blobHeaderHash [32]byte) []byte {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, blobHeaderHash[:]...))
	return buf.Bytes()
}

// EncodeBatchHeaderKey returns an encoded key as batch header identification.
func EncodeBatchHeaderKey(batchHeaderHash [32]byte) []byte {
	prefix := []byte(batchHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	return buf.Bytes()
}

func EncodeBlobIndexKey(batchHeaderHash [32]byte, blobIndex int) []byte {
	prefix := []byte(blobIndexPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	err := binary.Write(buf, binary.LittleEndian, int32(blobIndex))
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func EncodeBlobIndexKeyPrefix(batchHeaderHash [32]byte) []byte {
	prefix := []byte(blobIndexPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	return buf.Bytes()
}

// EncodeBatchExpirationKeyPrefix returns the encoded prefix for batch expiration key.
func EncodeBatchExpirationKeyPrefix() []byte {
	return []byte(batchExpirationPrefix)
}

// EncodeBlobExpirationKeyPrefix returns the encoded prefix for blob expiration key.
func EncodeBlobExpirationKeyPrefix() []byte {
	return []byte(blobExpirationPrefix)
}

// EncodeBatchMappingExpirationKeyPrefix returns the encoded prefix for the expiration key of the batch to blob index mapping.
func EncodeBatchMappingExpirationKeyPrefix() []byte {
	return []byte(batchMappingExpirationPrefix)
}

// Returns an encoded key for expration time.
// Note: the encoded key will preserve the order of expiration time, that is,
// expirationTime1 < expirationTime2 <=>
// EncodeBatchExpirationKey(expirationTime1) < EncodeBatchExpirationKey(expirationTime2)
func EncodeBatchExpirationKey(expirationTime int64) []byte {
	prefix := []byte(batchExpirationPrefix)
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts[0:8], uint64(expirationTime))
	buf := bytes.NewBuffer(append(prefix, ts[:]...))
	return buf.Bytes()
}

// EncodeBlobExpirationKey returns an encoded key for expration time for blob header hashes.
// Note: the encoded key will preserve the order of expiration time, that is,
// expirationTime1 < expirationTime2 <=>
// EncodeBlobExpirationKey(expirationTime1) < EncodeBlobExpirationKey(expirationTime2)
func EncodeBlobExpirationKey(expirationTime int64) []byte {
	prefix := []byte(blobExpirationPrefix)
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts[0:8], uint64(expirationTime))
	buf := bytes.NewBuffer(append(prefix, ts[:]...))
	return buf.Bytes()
}

// EncodeBatchMappingExpirationKeyPrefix returns an encoded key for expration time for the batch to blob index mapping.
// Encodes the expiration time and the batch header hash into a single key.
// Note: the encoded key will preserve the order of expiration time, that is,
// expirationTime1 < expirationTime2 <=>
// EncodeBatchMappingExpirationKeyPrefix(expirationTime1) < EncodeBatchMappingExpirationKeyPrefix(expirationTime2)
func EncodeBatchMappingExpirationKey(expirationTime int64, batchHeaderHash [32]byte) []byte {
	prefix := []byte(batchMappingExpirationPrefix)
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts[0:8], uint64(expirationTime))
	buf := bytes.NewBuffer(append(prefix, ts[:]...))
	buf.Write(batchHeaderHash[:])
	return buf.Bytes()
}

// DecodeBatchExpirationKey returns the expiration timestamp encoded in the key.
func DecodeBatchExpirationKey(key []byte) (int64, error) {
	if len(key) != len(batchExpirationPrefix)+8 {
		return 0, errors.New("the expiration key is invalid")
	}
	ts := int64(binary.BigEndian.Uint64(key[len(key)-8:]))
	return ts, nil
}

// Returns the expiration timestamp encoded in the key.
func DecodeBlobExpirationKey(key []byte) (int64, error) {
	if len(key) != len(blobExpirationPrefix)+8 {
		return 0, errors.New("the expiration key is invalid")
	}
	ts := int64(binary.BigEndian.Uint64(key[len(key)-8:]))
	return ts, nil
}

// DecodeBatchMappingExpirationKey returns the expiration timestamp encoded in the key.
func DecodeBatchMappingExpirationKey(key []byte) (int64, error) {
	if len(key) != len(batchMappingExpirationPrefix)+8+32 {
		return 0, errors.New("the expiration key is invalid")
	}
	ts := int64(binary.BigEndian.Uint64(key[len(key)-8-32 : len(key)-32]))
	return ts, nil
}

func DecodeHashSlice(input []byte) ([][32]byte, error) {
	if len(input)%32 != 0 {
		return nil, errors.New("input length is not a multiple of 32")
	}
	numHashes := len(input) / 32

	result := make([][32]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(result[i][:], input[i*32:(i+1)*32])
	}

	return result, nil
}
