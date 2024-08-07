package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common/pubip"
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
	blobPrefix           = "_BLOB_" // The prefix of the blob key.
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

// Returns the encoded prefix for batch expiration key.
func EncodeBatchExpirationKeyPrefix() []byte {
	return []byte(batchExpirationPrefix)
}

// Returns the encoded prefix for blob expiration key.
func EncodeBlobExpirationKeyPrefix() []byte {
	return []byte(blobExpirationPrefix)
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

// Returns an encoded key for expration time for blob header hashes.
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

// Returns the expiration timestamp encoded in the key.
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

func SocketAddress(ctx context.Context, provider pubip.Provider, dispersalPort string, retrievalPort string) (string, error) {
	ip, err := provider.PublicIPAddress(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get public ip address from IP provider: %w", err)
	}
	socket := core.MakeOperatorSocket(ip, dispersalPort, retrievalPort)
	return socket.String(), nil
}

func GetBundleEncodingFormat(blob *pb.Blob) core.BundleEncodingFormat {
	// We expect all the bundles of the blob are either using combined bundle
	// (with all chunks in a single byte array) or separate chunks, no mixed
	// use.
	for _, bundle := range blob.GetBundles() {
		// If the blob is using combined bundle encoding, there must be at least
		// one non-empty bundle (i.e. the node is in at least one quorum otherwise
		// it shouldn't have received this blob).
		if len(bundle.GetBundle()) > 0 {
			return core.GnarkBundleEncodingFormat
		}
	}
	return core.GobBundleEncodingFormat
}
