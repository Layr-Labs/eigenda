package node_test

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	nodemock "github.com/Layr-Labs/eigenda/node/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreBatchV2(t *testing.T) {
	_, batch, bundles := nodemock.MockBatch(t)

	blobShards := make([]*corev2.BlobShard, len(batch.BlobCertificates))
	rawBundles := make([]*node.RawBundles, len(batch.BlobCertificates))
	for i, cert := range batch.BlobCertificates {
		blobShards[i] = &corev2.BlobShard{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID]core.Bundle),
		}
		rawBundles[i] = &node.RawBundles{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID][]byte),
		}

		for quorum, bundle := range bundles[i] {
			blobShards[i].Bundles[quorum] = bundle
			bundleBytes, err := bundle.Serialize()
			assert.NoError(t, err)
			rawBundles[i].Bundles[quorum] = bundleBytes
		}
	}

	s, db := createStoreV2(t)
	defer func() {
		_ = db.Shutdown()
	}()
	keys, _, err := s.StoreBatch(batch, rawBundles)
	require.NoError(t, err)
	require.Len(t, keys, 7)

	tables := db.GetTables()
	require.ElementsMatch(t, []string{node.BatchHeaderTableName, node.BlobCertificateTableName, node.BundleTableName}, tables)

	// Check batch header
	bhh, err := batch.BatchHeader.Hash()
	require.NoError(t, err)
	batchHeaderKeyBuilder, err := db.GetKeyBuilder(node.BatchHeaderTableName)
	require.NoError(t, err)
	bhhBytes, err := db.Get(batchHeaderKeyBuilder.Key(bhh[:]))
	require.NoError(t, err)
	assert.NotNil(t, bhhBytes)
	deserializedBatchHeader, err := corev2.DeserializeBatchHeader(bhhBytes)
	require.NoError(t, err)
	assert.Equal(t, batch.BatchHeader, deserializedBatchHeader)

	// Check bundles
	bundleKeyBuilder, err := db.GetKeyBuilder(node.BundleTableName)
	require.NoError(t, err)
	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		require.NoError(t, err)
		for quorum, bundle := range bundles.Bundles {
			k, err := node.BundleKey(blobKey, quorum)
			require.NoError(t, err)
			bundleBytes, err := db.Get(bundleKeyBuilder.Key(k))
			require.NoError(t, err)
			assert.NotNil(t, bundleBytes)
			require.Equal(t, bundle, bundleBytes)
		}
	}

	// Try to store the same batch again
	_, _, err = s.StoreBatch(batch, rawBundles)
	require.ErrorIs(t, err, node.ErrBatchAlreadyExist)

	// Check deletion
	err = s.DeleteKeys(keys)
	require.NoError(t, err)

	bhhBytes, err = db.Get(batchHeaderKeyBuilder.Key(bhh[:]))
	require.Error(t, err)
	require.Empty(t, bhhBytes)

	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		require.NoError(t, err)
		for quorum := range bundles.Bundles {
			k, err := node.BundleKey(blobKey, quorum)
			require.NoError(t, err)
			bundleBytes, err := db.Get(bundleKeyBuilder.Key(k))
			require.Error(t, err)
			require.Empty(t, bundleBytes)
		}
	}
}

func TestGetChunks(t *testing.T) {
	blobKeys, batch, bundles := nodemock.MockBatch(t)

	blobShards := make([]*corev2.BlobShard, len(batch.BlobCertificates))
	rawBundles := make([]*node.RawBundles, len(batch.BlobCertificates))
	for i, cert := range batch.BlobCertificates {
		blobShards[i] = &corev2.BlobShard{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID]core.Bundle),
		}
		rawBundles[i] = &node.RawBundles{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID][]byte),
		}

		for quorum, bundle := range bundles[i] {
			blobShards[i].Bundles[quorum] = bundle
			bundleBytes, err := bundle.Serialize()
			assert.NoError(t, err)
			rawBundles[i].Bundles[quorum] = bundleBytes
		}
	}

	s, db := createStoreV2(t)
	defer func() {
		_ = db.Shutdown()
	}()
	_, _, err := s.StoreBatch(batch, rawBundles)
	require.NoError(t, err)

	chunks, err := s.GetChunks(blobKeys[0], 0)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[0][0]))

	chunks, err = s.GetChunks(blobKeys[0], 1)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[0][1]))

	chunks, err = s.GetChunks(blobKeys[1], 0)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[1][0]))

	chunks, err = s.GetChunks(blobKeys[1], 1)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[1][1]))

	chunks, err = s.GetChunks(blobKeys[2], 1)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[2][1]))

	chunks, err = s.GetChunks(blobKeys[2], 2)
	require.NoError(t, err)
	require.Len(t, chunks, len(bundles[2][2]))

	// wrong quorum
	_, err = s.GetChunks(blobKeys[0], 2)
	require.Error(t, err)
}

func createStoreV2(t *testing.T) (node.StoreV2, kvstore.TableStore) {
	logger := logging.NewNoopLogger()
	config := tablestore.DefaultLevelDBConfig(t.TempDir())
	config.Schema = []string{node.BatchHeaderTableName, node.BlobCertificateTableName, node.BundleTableName}
	tStore, err := tablestore.Start(logger, config)
	require.NoError(t, err)
	s := node.NewLevelDBStoreV2(tStore, logger, 10*time.Second)
	return s, tStore
}
