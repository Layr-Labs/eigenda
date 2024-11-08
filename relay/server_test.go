package relay

import (
	"context"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestReadWriteBlobs(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	blobStore := buildBlobStore(t, logger)

	// This is the server used to read it back
	config := DefaultConfig()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil /* not used in this test*/)
	require.NoError(t, err)

	expectedData := make(map[v2.BlobKey][]byte)

	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, data := randomBlob(t)

		blobKey, err := header.BlobKey()
		require.NoError(t, err)
		expectedData[blobKey] = data

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
			},
			&encoding.FragmentInfo{})
		require.NoError(t, err)

		err = blobStore.StoreBlob(context.Background(), blobKey, data)
		require.NoError(t, err)
	}

	// Read the blobs back.
	for key, data := range expectedData {
		request := &pb.GetBlobRequest{
			BlobKey: key[:],
		}

		response, err := server.GetBlob(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, data, response.Blob)
	}

	// Read the blobs back again to test caching.
	for key, data := range expectedData {
		request := &pb.GetBlobRequest{
			BlobKey: key[:],
		}

		response, err := server.GetBlob(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, data, response.Blob)
	}
}

func TestReadNonExistentBlob(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	blobStore := buildBlobStore(t, logger)

	// This is the server used to read it back
	config := DefaultConfig()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil /* not used in this test */)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		request := &pb.GetBlobRequest{
			BlobKey: tu.RandomBytes(32),
		}

		response, err := server.GetBlob(context.Background(), request)
		require.Error(t, err)
		require.Nil(t, response)
	}
}

func TestReadWriteBlobsWithSharding(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	blobStore := buildBlobStore(t, logger)

	shardCount := rand.Intn(10) + 10
	shardList := make([]v2.RelayKey, 0)
	shardSet := make(map[v2.RelayKey]struct{})
	for i := 0; i < shardCount; i++ {
		if i%2 == 0 {
			shardList = append(shardList, v2.RelayKey(i))
			shardSet[v2.RelayKey(i)] = struct{}{}
		}
	}

	// This is the server used to read it back
	config := DefaultConfig()
	config.Shards = shardList
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil /* not used in this test*/)
	require.NoError(t, err)

	expectedData := make(map[v2.BlobKey][]byte)
	shardMap := make(map[v2.BlobKey][]v2.RelayKey)

	blobCount := 100
	for i := 0; i < blobCount; i++ {
		header, data := randomBlob(t)

		blobKey, err := header.BlobKey()
		require.NoError(t, err)
		expectedData[blobKey] = data

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
			&encoding.FragmentInfo{})
		require.NoError(t, err)

		err = blobStore.StoreBlob(context.Background(), blobKey, data)
		require.NoError(t, err)
	}

	// Read the blobs back. On average, we expect 25% of the blobs to be assigned to shards we don't have.
	for key, data := range expectedData {
		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		request := &pb.GetBlobRequest{
			BlobKey: key[:],
		}

		response, err := server.GetBlob(context.Background(), request)

		if isBlobInCorrectShard {
			require.NoError(t, err)
			require.Equal(t, data, response.Blob)
		} else {
			require.Error(t, err)
			require.Nil(t, response)
		}
	}

	// Read the blobs back again to test caching.
	for key, data := range expectedData {
		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		request := &pb.GetBlobRequest{
			BlobKey: key[:],
		}

		response, err := server.GetBlob(context.Background(), request)

		if isBlobInCorrectShard {
			require.NoError(t, err)
			require.Equal(t, data, response.Blob)
		} else {
			require.Error(t, err)
			require.Nil(t, response)
		}
	}
}

func TestReadWriteChunks(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	chunkReader, chunkWriter := buildChunkStore(t, logger)

	// This is the server used to read it back
	config := DefaultConfig()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		nil, /* not used in this test*/
		chunkReader)
	require.NoError(t, err)

	expectedData := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, _, chunks := randomBlobChunks(t)

		blobKey, err := header.BlobKey()
		require.NoError(t, err)
		expectedData[blobKey] = chunks

		coeffs, chunkProofs := disassembleFrames(chunks)
		err = chunkWriter.PutChunkProofs(context.Background(), blobKey, chunkProofs)
		require.NoError(t, err)
		fragmentInfo, err := chunkWriter.PutChunkCoefficients(context.Background(), blobKey, coeffs)
		require.NoError(t, err)
		fragmentInfoMap[blobKey] = fragmentInfo

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
				FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Request the entire blob by range
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)
		requestedChunks = append(requestedChunks, &pb.ChunkRequest{
			Request: &pb.ChunkRequest_ByRange{
				ByRange: &pb.ChunkRequestByRange{
					BlobKey:    key[:],
					StartIndex: 0,
					EndIndex:   uint32(len(data)),
				},
			},
		})
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		response, err := server.GetChunks(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		for i, frame := range response.Data[0].Data {
			convertedFrame := encoding.FrameFromProtobuf(frame)
			require.Equal(t, data[i], convertedFrame)
		}
	}

	// Request the entire blob by index
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)

		indices := make([]uint32, len(data))
		for i := range data {
			indices[i] = uint32(i)
		}

		requestedChunks = append(requestedChunks, &pb.ChunkRequest{
			Request: &pb.ChunkRequest_ByIndex{
				ByIndex: &pb.ChunkRequestByIndex{
					BlobKey:      key[:],
					ChunkIndices: indices,
				},
			},
		})
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		response, err := server.GetChunks(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		for i, frame := range response.Data[0].Data {
			convertedFrame := encoding.FrameFromProtobuf(frame)
			require.Equal(t, data[i], convertedFrame)
		}
	}

	// Request part of the blob back by range
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)

		startIndex := rand.Intn(len(data))
		var endIndex int
		if startIndex == len(data)-1 {
			endIndex = len(data)
		} else {
			endIndex = startIndex + rand.Intn(len(data)-startIndex)
		}

		requestedChunks = append(requestedChunks, &pb.ChunkRequest{
			Request: &pb.ChunkRequest_ByRange{
				ByRange: &pb.ChunkRequestByRange{
					BlobKey:    key[:],
					StartIndex: uint32(startIndex),
					EndIndex:   uint32(endIndex),
				},
			},
		})
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		response, err := server.GetChunks(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		for i := startIndex; i < endIndex; i++ {
			convertedFrame := encoding.FrameFromProtobuf(response.Data[0].Data[i-startIndex])
			require.Equal(t, data[i], convertedFrame)
		}
	}

	// Request part of the blob back by index
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)

		indices := make([]uint32, 0)
		for i := range data {
			if i%2 == 0 {
				indices = append(indices, uint32(i))
			}
		}

		requestedChunks = append(requestedChunks, &pb.ChunkRequest{
			Request: &pb.ChunkRequest_ByIndex{
				ByIndex: &pb.ChunkRequestByIndex{
					BlobKey:      key[:],
					ChunkIndices: indices,
				},
			},
		})
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		response, err := server.GetChunks(context.Background(), request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		for i := 0; i < len(indices); i++ {
			if i%2 == 0 {
				convertedFrame := encoding.FrameFromProtobuf(response.Data[0].Data[i/2])
				require.Equal(t, data[indices[i]], convertedFrame)
			}
		}
	}
}
