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

// EncodeBatchHeaderKey returns an encoded key as batch header identification.
func EncodeBatchHeaderKey(batchHeaderHash [32]byte) []byte {
	prefix := []byte(batchHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, batchHeaderHash[:]...))
	return buf.Bytes()
}

// Returns the encoded prefix for batch expiration key.
func EncodeBatchExpirationKeyPrefix() []byte {
	return []byte(batchExpirationPrefix)
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

// Returns the expiration timestamp encoded in the key.
func DecodeBatchExpirationKey(key []byte) (int64, error) {
	if len(key) != len(batchExpirationPrefix)+8 {
		return 0, errors.New("the expiration key is invalid")
	}
	ts := int64(binary.BigEndian.Uint64(key[len(key)-8:]))
	return ts, nil
}
