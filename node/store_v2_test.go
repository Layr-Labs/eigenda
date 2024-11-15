package node_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreBatchV2(t *testing.T) {
	ctx := context.Background()
	_, batch, bundles := mockBatch(t)

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
	err := s.StoreBatch(ctx, batch, rawBundles)
	require.NoError(t, err)

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

	// Check blob certificates
	blobCertKeyBuilder, err := db.GetKeyBuilder(node.BlobCertificateTableName)
	require.NoError(t, err)
	for _, cert := range batch.BlobCertificates {
		blobKey, err := cert.BlobHeader.BlobKey()
		require.NoError(t, err)
		blobCertKey := blobCertKeyBuilder.Key(blobKey[:])
		blobCertBytes, err := db.Get(blobCertKey)
		require.NoError(t, err)
		assert.NotNil(t, blobCertBytes)
		deserializedBlobCert, err := corev2.DeserializeBlobCertificate(blobCertBytes)
		require.NoError(t, err)
		assert.Equal(t, cert, deserializedBlobCert)
	}

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
}

func createStoreV2(t *testing.T) (*node.StoreV2, kvstore.TableStore) {
	logger := logging.NewNoopLogger()
	config := tablestore.DefaultLevelDBConfig(t.TempDir())
	config.Schema = []string{node.BatchHeaderTableName, node.BlobCertificateTableName, node.BundleTableName}
	tStore, err := tablestore.Start(logger, config)
	require.NoError(t, err)
	s := node.NewLevelDBStoreV2(tStore, logger)
	return s, tStore
}
