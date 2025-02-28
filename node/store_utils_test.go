package node_test

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"testing"

	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/assert"
)

func TestDecodeHashSlice(t *testing.T) {
	hash0 := [32]byte{0, 1}
	hash1 := [32]byte{0, 1, 2, 3, 4}
	hash2 := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	input := make([]byte, 0)
	input = append(input, hash0[:]...)
	input = append(input, hash1[:]...)
	input = append(input, hash2[:]...)

	hashes, err := node.DecodeHashSlice(input)
	assert.NoError(t, err)
	assert.Len(t, hashes, 3)
	assert.Equal(t, hash0, hashes[0])
	assert.Equal(t, hash1, hashes[1])
	assert.Equal(t, hash2, hashes[2])
}

func TestEncodeDecodeBatchMappingExpirationKey(t *testing.T) {
	expirationTime := int64(1234567890)
	batchHeaderHash := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	key := node.EncodeBatchMappingExpirationKey(expirationTime, batchHeaderHash)
	decodedExpirationTime, err := node.DecodeBatchMappingExpirationKey(key)
	assert.NoError(t, err)
	assert.Equal(t, expirationTime, decodedExpirationTime)
}

func TestBatchMappingExpirationKeyOrdering(t *testing.T) {
	dbPath := t.TempDir()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	db, err := leveldb.NewStore(logger, dbPath)
	defer func() {
		err = db.Destroy()
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)

	batch := db.NewBatch()

	// test ordering using expiration time
	expirationTime := int64(1111111111)
	batchHeaderHash := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	key := node.EncodeBatchMappingExpirationKey(expirationTime, batchHeaderHash)
	batch.Put(key, []byte("value"))

	expirationTime = int64(2222222222)
	batchHeaderHash = [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	key = node.EncodeBatchMappingExpirationKey(expirationTime, batchHeaderHash)
	batch.Put(key, []byte("value"))

	expirationTime = int64(3333333333)
	batchHeaderHash = [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	key = node.EncodeBatchMappingExpirationKey(expirationTime, batchHeaderHash)
	batch.Put(key, []byte("value"))

	err = batch.Apply()
	assert.NoError(t, err)

	iter, err := db.NewIterator(node.EncodeBatchMappingExpirationKeyPrefix())
	assert.NoError(t, err)
	defer iter.Release()
	i := 0
	expectedExpirationTimes := []int64{1111111111, 2222222222, 3333333333}
	for iter.Next() {
		ts, err := node.DecodeBatchMappingExpirationKey(iter.Key())
		assert.NoError(t, err)
		assert.Equal(t, expectedExpirationTimes[i], ts)
		i++
	}
	assert.Equal(t, 3, i)
}
