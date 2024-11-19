package relay

import (
	"context"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestGetNonExistentBlob(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()
	metadataStore := buildMetadataStore(t)

	server, err := newMetadataProvider(context.Background(), logger, metadataStore, 1024*1024, 32, nil)
	require.NoError(t, err)

	// Try to fetch a non-existent blobs
	for i := 0; i < 10; i++ {
		_, err := server.GetMetadataForBlobs([]v2.BlobKey{v2.BlobKey(tu.RandomBytes(32))})
		require.Error(t, err)
	}
}

func TestFetchingIndividualMetadata(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()
	metadataStore := buildMetadataStore(t)

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, _ := randomBlob(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		totalChunkSizeBytes := uint32(rand.Intn(1024 * 1024 * 1024))
		fragmentSizeBytes := uint32(rand.Intn(1024 * 1024))

		totalChunkSizeMap[blobKey] = totalChunkSizeBytes
		fragmentSizeMap[blobKey] = fragmentSizeBytes

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: totalChunkSizeBytes,
				FragmentSizeBytes:   fragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Sanity check, make sure the metadata is in the low level store
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		cert, fragmentInfo, err := metadataStore.GetBlobCertificate(context.Background(), blobKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.NotNil(t, fragmentInfo)
		require.Equal(t, totalChunkSizeBytes, fragmentInfo.TotalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], fragmentInfo.FragmentSizeBytes)
	}

	server, err := newMetadataProvider(context.Background(), logger, metadataStore, 1024*1024, 32, nil)
	require.NoError(t, err)

	// Fetch the metadata from the server.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		mMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})
		require.NoError(t, err)
		require.Equal(t, 1, len(mMap))
		metadata := mMap[blobKey]
		require.NotNil(t, metadata)
		require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
	}

	// Read it back again. This uses a different code pathway due to the cache.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		mMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})
		require.NoError(t, err)
		require.Equal(t, 1, len(mMap))
		metadata := mMap[blobKey]
		require.NotNil(t, metadata)
		require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
	}
}

func TestBatchedFetch(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()
	metadataStore := buildMetadataStore(t)

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, _ := randomBlob(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		totalChunkSizeBytes := uint32(rand.Intn(1024 * 1024 * 1024))
		fragmentSizeBytes := uint32(rand.Intn(1024 * 1024))

		totalChunkSizeMap[blobKey] = totalChunkSizeBytes
		fragmentSizeMap[blobKey] = fragmentSizeBytes

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: totalChunkSizeBytes,
				FragmentSizeBytes:   fragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Sanity check, make sure the metadata is in the low level store
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		cert, fragmentInfo, err := metadataStore.GetBlobCertificate(context.Background(), blobKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.NotNil(t, fragmentInfo)
		require.Equal(t, totalChunkSizeBytes, fragmentInfo.TotalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], fragmentInfo.FragmentSizeBytes)
	}

	server, err := newMetadataProvider(context.Background(), logger, metadataStore, 1024*1024, 32, nil)
	require.NoError(t, err)

	// Each iteration, choose a random subset of the keys to fetch
	for i := 0; i < 10; i++ {
		keyCount := rand.Intn(blobCount) + 1
		keys := make([]v2.BlobKey, 0, keyCount)
		for key := range totalChunkSizeMap {
			keys = append(keys, key)
			if len(keys) == keyCount {
				break
			}
		}

		mMap, err := server.GetMetadataForBlobs(keys)
		require.NoError(t, err)

		assert.Equal(t, keyCount, len(mMap))
		for _, key := range keys {
			metadata := mMap[key]
			require.NotNil(t, metadata)
			require.Equal(t, totalChunkSizeMap[key], metadata.totalChunkSizeBytes)
			require.Equal(t, fragmentSizeMap[key], metadata.fragmentSizeBytes)
		}
	}
}

func TestIndividualFetchWithSharding(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()
	metadataStore := buildMetadataStore(t)

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)
	shardMap := make(map[v2.BlobKey][]v2.RelayKey)

	shardCount := rand.Intn(10) + 10
	shardList := make([]v2.RelayKey, 0)
	shardSet := make(map[v2.RelayKey]struct{})
	for i := 0; i < shardCount; i++ {
		if i%2 == 0 {
			shardList = append(shardList, v2.RelayKey(i))
			shardSet[v2.RelayKey(i)] = struct{}{}
		}
	}

	// Write some metadata
	blobCount := 100
	for i := 0; i < blobCount; i++ {
		header, _ := randomBlob(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		totalChunkSizeBytes := uint32(rand.Intn(1024 * 1024 * 1024))
		fragmentSizeBytes := uint32(rand.Intn(1024 * 1024))

		totalChunkSizeMap[blobKey] = totalChunkSizeBytes
		fragmentSizeMap[blobKey] = fragmentSizeBytes

		// Assign two shards to each blob.
		shard1 := v2.RelayKey(rand.Intn(shardCount))
		shard2 := v2.RelayKey(rand.Intn(shardCount))
		shards := []v2.RelayKey{shard1, shard2}
		shardMap[blobKey] = shards

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
				RelayKeys:  shards,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: totalChunkSizeBytes,
				FragmentSizeBytes:   fragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Sanity check, make sure the metadata is in the low level store
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		cert, fragmentInfo, err := metadataStore.GetBlobCertificate(context.Background(), blobKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.NotNil(t, fragmentInfo)
		require.Equal(t, totalChunkSizeBytes, fragmentInfo.TotalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], fragmentInfo.FragmentSizeBytes)
	}

	server, err := newMetadataProvider(context.Background(), logger, metadataStore, 1024*1024, 32, shardList)
	require.NoError(t, err)

	// Fetch the metadata from the server.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		isBlobInCorrectShard := false
		blobShards := shardMap[blobKey]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		mMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})

		if isBlobInCorrectShard {
			// The blob is in the relay's shard, should be returned like normal
			require.NoError(t, err)
			require.Equal(t, 1, len(mMap))
			metadata := mMap[blobKey]
			require.NotNil(t, metadata)
			require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
			require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
		} else {
			// the blob is not in the relay's shard, should return an error
			require.Error(t, err)
		}
	}

	// Read it back again. This uses a different code pathway due to the cache.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		isBlobInCorrectShard := false
		blobShards := shardMap[blobKey]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		mMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})

		if isBlobInCorrectShard {
			// The blob is in the relay's shard, should be returned like normal
			require.NoError(t, err)
			require.Equal(t, 1, len(mMap))
			metadata := mMap[blobKey]
			require.NotNil(t, metadata)
			require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
			require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
		} else {
			// the blob is not in the relay's shard, should return an error
			require.Error(t, err)
		}
	}
}

func TestBatchedFetchWithSharding(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()
	metadataStore := buildMetadataStore(t)

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)
	shardMap := make(map[v2.BlobKey][]v2.RelayKey)

	shardCount := rand.Intn(10) + 10
	shardList := make([]v2.RelayKey, 0)
	shardSet := make(map[v2.RelayKey]struct{})
	for i := 0; i < shardCount; i++ {
		if i%2 == 0 {
			shardList = append(shardList, v2.RelayKey(i))
			shardSet[v2.RelayKey(i)] = struct{}{}
		}
	}

	// Write some metadata
	blobCount := 100
	for i := 0; i < blobCount; i++ {
		header, _ := randomBlob(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		totalChunkSizeBytes := uint32(rand.Intn(1024 * 1024 * 1024))
		fragmentSizeBytes := uint32(rand.Intn(1024 * 1024))

		totalChunkSizeMap[blobKey] = totalChunkSizeBytes
		fragmentSizeMap[blobKey] = fragmentSizeBytes

		// Assign two shards to each blob.
		shard1 := v2.RelayKey(rand.Intn(shardCount))
		shard2 := v2.RelayKey(rand.Intn(shardCount))
		shards := []v2.RelayKey{shard1, shard2}
		shardMap[blobKey] = shards

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
				RelayKeys:  shards,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: totalChunkSizeBytes,
				FragmentSizeBytes:   fragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Sanity check, make sure the metadata is in the low level store
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		cert, fragmentInfo, err := metadataStore.GetBlobCertificate(context.Background(), blobKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.NotNil(t, fragmentInfo)
		require.Equal(t, totalChunkSizeBytes, fragmentInfo.TotalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], fragmentInfo.FragmentSizeBytes)
	}

	server, err := newMetadataProvider(context.Background(), logger, metadataStore, 1024*1024, 32, shardList)
	require.NoError(t, err)

	// Each iteration, choose two random keys to fetch. There will be a 25% chance that both blobs map to valid shards.
	for i := 0; i < 100; i++ {

		keyCount := 2
		keys := make([]v2.BlobKey, 0, keyCount)
		areKeysInCorrectShard := true
		for key := range totalChunkSizeMap {
			keys = append(keys, key)

			keyShards := shardMap[key]
			keyIsInShard := false
			for _, shard := range keyShards {
				if _, ok := shardSet[shard]; ok {
					keyIsInShard = true
					break
				}
			}
			if !keyIsInShard {
				// If both keys are not in the shard, we expect an error.
				areKeysInCorrectShard = false
			}

			if len(keys) == keyCount {
				break
			}
		}

		mMap, err := server.GetMetadataForBlobs(keys)
		if areKeysInCorrectShard {
			require.NoError(t, err)
			assert.Equal(t, keyCount, len(mMap))
			for _, key := range keys {
				metadata := mMap[key]
				require.NotNil(t, metadata)
				require.Equal(t, totalChunkSizeMap[key], metadata.totalChunkSizeBytes)
				require.Equal(t, fragmentSizeMap[key], metadata.fragmentSizeBytes)
			}
		} else {
			require.Error(t, err)
		}
	}
}
