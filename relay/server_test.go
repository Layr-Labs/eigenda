package relay

import (
	"context"
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/relay/auth"
	"github.com/Layr-Labs/eigenda/relay/mock"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/relay/limiter"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func defaultConfig() *Config {
	return &Config{
		GRPCPort:                   50051,
		MaxGRPCMessageSize:         1024 * 1024 * 300,
		MetadataCacheSize:          1024 * 1024,
		MetadataMaxConcurrency:     32,
		BlobCacheBytes:             1024 * 1024,
		BlobMaxConcurrency:         32,
		ChunkCacheSize:             1024 * 1024,
		ChunkMaxConcurrency:        32,
		MaxKeysPerGetChunksRequest: 1024,
		AuthenticationKeyCacheSize: 1024,
		AuthenticationDisabled:     false,
		RateLimits: limiter.Config{
			MaxGetBlobOpsPerSecond:          1024,
			GetBlobOpsBurstiness:            1024,
			MaxGetBlobBytesPerSecond:        20 * 1024 * 1024,
			GetBlobBytesBurstiness:          20 * 1024 * 1024,
			MaxConcurrentGetBlobOps:         1024,
			MaxGetChunkOpsPerSecond:         1024,
			GetChunkOpsBurstiness:           1024,
			MaxGetChunkBytesPerSecond:       20 * 1024 * 1024,
			GetChunkBytesBurstiness:         20 * 1024 * 1024,
			MaxConcurrentGetChunkOps:        1024,
			MaxGetChunkOpsPerSecondClient:   8,
			GetChunkOpsBurstinessClient:     8,
			MaxGetChunkBytesPerSecondClient: 2 * 1024 * 1024,
			GetChunkBytesBurstinessClient:   2 * 1024 * 1024,
			MaxConcurrentGetChunkOpsClient:  1,
		},
		Timeouts: TimeoutConfig{
			GetBlobTimeout:                 10 * time.Second,
			GetChunksTimeout:               10 * time.Second,
			InternalGetMetadataTimeout:     10 * time.Second,
			InternalGetBlobTimeout:         10 * time.Second,
			InternalGetProofsTimeout:       10 * time.Second,
			InternalGetCoefficientsTimeout: 10 * time.Second,
		},
		MetricsPort: 9101,
	}
}

func getBlob(t *testing.T, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("0.0.0.0:50051", opts...)
	require.NoError(t, err)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := pb.NewRelayClient(conn)
	response, err := client.GetBlob(context.Background(), request)
	return response, err
}

func getChunks(
	t *testing.T,
	random *random.TestRandom,
	operatorKeys map[uint32]*core.KeyPair,
	request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// Choose a random operator to send this request as. Operator IDs are expected to be sequential starting at 0.
	operatorID := random.Uint32() % uint32(len(operatorKeys))
	operatorIDBytes := make([]byte, 32)
	binary.BigEndian.PutUint32(operatorIDBytes[24:], operatorID)
	request.OperatorId = operatorIDBytes
	signature := auth.SignGetChunksRequest(operatorKeys[operatorID], request)
	request.OperatorSignature = signature

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("0.0.0.0:50051", opts...)
	require.NoError(t, err)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := pb.NewRelayClient(conn)
	response, err := client.GetChunks(context.Background(), request)
	return response, err
}

func TestReadWriteBlobs(t *testing.T) {
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	blobStore := buildBlobStore(t, logger)
	chainReader := newMockChainReader()

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil, /* not used in this test*/
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

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

		response, err := getBlob(t, request)
		require.NoError(t, err)

		require.Equal(t, data, response.Blob)
	}

	// Read the blobs back again to test caching.
	for key, data := range expectedData {
		request := &pb.GetBlobRequest{
			BlobKey: key[:],
		}

		response, err := getBlob(t, request)
		require.NoError(t, err)

		require.Equal(t, data, response.Blob)
	}
}

func TestReadNonExistentBlob(t *testing.T) {
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	blobStore := buildBlobStore(t, logger)

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil, /* not used in this test */
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

	for i := 0; i < 10; i++ {
		request := &pb.GetBlobRequest{
			BlobKey: tu.RandomBytes(32),
		}

		response, err := getBlob(t, request)
		require.Error(t, err)
		require.Nil(t, response)
	}
}

func TestReadWriteBlobsWithSharding(t *testing.T) {
	rand := random.NewTestRandom(t)

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

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	config.RelayIDs = shardList
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		blobStore,
		nil, /* not used in this test*/
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

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

		response, err := getBlob(t, request)

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

		response, err := getBlob(t, request)

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
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	chunkReader, chunkWriter := buildChunkStore(t, logger)

	operatorCount := rand.Intn(3) + 1
	operatorKeys := make(map[uint32]*core.KeyPair)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	for i := 0; i < operatorCount; i++ {
		keypair := rand.BLS()
		operatorKeys[uint32(i)] = keypair

		var operatorID core.OperatorID
		binary.BigEndian.PutUint32(operatorID[24:], uint32(i))
		operatorInfo[operatorID] = &core.IndexedOperatorInfo{
			PubkeyG1: keypair.GetPubKeyG1(),
			PubkeyG2: keypair.GetPubKeyG2(),
		}
	}

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	config.RateLimits.MaxGetChunkOpsPerSecond = 1000
	config.RateLimits.GetChunkOpsBurstiness = 1000
	config.RateLimits.MaxGetChunkOpsPerSecondClient = 1000
	config.RateLimits.GetChunkOpsBurstinessClient = 1000
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		nil, /* not used in this test*/
		chunkReader,
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

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

		response, err := getChunks(t, rand, operatorKeys, request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		bundle, err := core.Bundle{}.Deserialize(response.Data[0])
		require.NoError(t, err)

		for i, frame := range bundle {
			require.Equal(t, data[i], frame)
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

		response, err := getChunks(t, rand, operatorKeys, request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		bundle, err := core.Bundle{}.Deserialize(response.Data[0])
		require.NoError(t, err)

		for i, frame := range bundle {
			require.Equal(t, data[i], frame)
		}
	}

	// Request part of the blob back by range
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)

		startIndex := rand.Intn(len(data) - 1)
		endIndex := startIndex + rand.Intn(len(data)-startIndex-1) + 1

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

		response, err := getChunks(t, rand, operatorKeys, request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		bundle, err := core.Bundle{}.Deserialize(response.Data[0])
		require.NoError(t, err)

		for i := startIndex; i < endIndex; i++ {
			require.Equal(t, data[i], bundle[i-startIndex])
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

		response, err := getChunks(t, rand, operatorKeys, request)
		require.NoError(t, err)

		require.Equal(t, 1, len(response.Data))

		bundle, err := core.Bundle{}.Deserialize(response.Data[0])
		require.NoError(t, err)

		for i := 0; i < len(indices); i++ {
			if i%2 == 0 {
				require.Equal(t, data[indices[i]], bundle[i/2])
			}
		}
	}
}

func TestBatchedReadWriteChunks(t *testing.T) {
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	chunkReader, chunkWriter := buildChunkStore(t, logger)

	operatorCount := rand.Intn(3) + 1
	operatorKeys := make(map[uint32]*core.KeyPair)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	for i := 0; i < operatorCount; i++ {
		keypair := rand.BLS()
		operatorKeys[uint32(i)] = keypair

		var operatorID core.OperatorID
		binary.BigEndian.PutUint32(operatorID[24:], uint32(i))
		operatorInfo[operatorID] = &core.IndexedOperatorInfo{
			PubkeyG1: keypair.GetPubKeyG1(),
			PubkeyG2: keypair.GetPubKeyG2(),
		}
	}

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		nil, /* not used in this test */
		chunkReader,
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

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

	keyCount := 3

	for i := 0; i < 10; i++ {
		keys := make([]v2.BlobKey, 0, keyCount)
		for key := range expectedData {
			keys = append(keys, key)
			if len(keys) == keyCount {
				break
			}
		}

		requestedChunks := make([]*pb.ChunkRequest, 0)
		for _, key := range keys {

			boundKey := key
			request := &pb.ChunkRequest{
				Request: &pb.ChunkRequest_ByRange{
					ByRange: &pb.ChunkRequestByRange{
						BlobKey:    boundKey[:],
						StartIndex: 0,
						EndIndex:   uint32(len(expectedData[key])),
					},
				},
			}

			requestedChunks = append(requestedChunks, request)
		}
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		response, err := getChunks(t, rand, operatorKeys, request)
		require.NoError(t, err)

		require.Equal(t, keyCount, len(response.Data))

		for keyIndex, key := range keys {
			data := expectedData[key]

			bundle, err := core.Bundle{}.Deserialize(response.Data[keyIndex])
			require.NoError(t, err)

			for frameIndex, frame := range bundle {
				require.Equal(t, data[frameIndex], frame)
			}
		}
	}
}

func TestReadWriteChunksWithSharding(t *testing.T) {
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	chunkReader, chunkWriter := buildChunkStore(t, logger)

	shardCount := rand.Intn(10) + 10
	shardList := make([]v2.RelayKey, 0)
	shardSet := make(map[v2.RelayKey]struct{})
	for i := 0; i < shardCount; i++ {
		if i%2 == 0 {
			shardList = append(shardList, v2.RelayKey(i))
			shardSet[v2.RelayKey(i)] = struct{}{}
		}
	}
	shardMap := make(map[v2.BlobKey][]v2.RelayKey)

	operatorCount := rand.Intn(3) + 1
	operatorKeys := make(map[uint32]*core.KeyPair)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	for i := 0; i < operatorCount; i++ {
		keypair := rand.BLS()
		operatorKeys[uint32(i)] = keypair

		var operatorID core.OperatorID
		binary.BigEndian.PutUint32(operatorID[24:], uint32(i))
		operatorInfo[operatorID] = &core.IndexedOperatorInfo{
			PubkeyG1: keypair.GetPubKeyG1(),
			PubkeyG2: keypair.GetPubKeyG2(),
		}
	}

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	config.RelayIDs = shardList
	config.RateLimits.MaxGetChunkOpsPerSecond = 1000
	config.RateLimits.GetChunkOpsBurstiness = 1000
	config.RateLimits.MaxGetChunkOpsPerSecondClient = 1000
	config.RateLimits.GetChunkOpsBurstinessClient = 1000
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		nil, /* not used in this test*/
		chunkReader,
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

	expectedData := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, _, chunks := randomBlobChunks(t)

		blobKey, err := header.BlobKey()
		require.NoError(t, err)
		expectedData[blobKey] = chunks

		// Assign two shards to each blob.
		shard1 := v2.RelayKey(rand.Intn(shardCount))
		shard2 := v2.RelayKey(rand.Intn(shardCount))
		shards := []v2.RelayKey{shard1, shard2}
		shardMap[blobKey] = shards

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
				RelayKeys:  shards,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
				FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Request the entire blob by range. 25% of the blobs will be assigned to shards we don't have.
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

		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		response, err := getChunks(t, rand, operatorKeys, request)

		if isBlobInCorrectShard {
			require.NoError(t, err)

			require.Equal(t, 1, len(response.Data))

			bundle, err := core.Bundle{}.Deserialize(response.Data[0])
			require.NoError(t, err)

			for i, frame := range bundle {
				require.Equal(t, data[i], frame)
			}
		} else {
			require.Error(t, err)
			require.Nil(t, response)
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

		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		if isBlobInCorrectShard {
			response, err := getChunks(t, rand, operatorKeys, request)
			require.NoError(t, err)

			require.Equal(t, 1, len(response.Data))

			bundle, err := core.Bundle{}.Deserialize(response.Data[0])
			require.NoError(t, err)

			for i, frame := range bundle {
				require.Equal(t, data[i], frame)
			}
		} else {
			response, err := getChunks(t, rand, operatorKeys, request)
			require.Error(t, err)
			require.Nil(t, response)
		}
	}

	// Request part of the blob back by range
	for key, data := range expectedData {
		requestedChunks := make([]*pb.ChunkRequest, 0)

		startIndex := rand.Intn(len(data) - 1)
		endIndex := startIndex + rand.Intn(len(data)-startIndex-1) + 1

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

		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		if isBlobInCorrectShard {
			response, err := getChunks(t, rand, operatorKeys, request)
			require.NoError(t, err)

			require.Equal(t, 1, len(response.Data))

			bundle, err := core.Bundle{}.Deserialize(response.Data[0])
			require.NoError(t, err)

			for i := startIndex; i < endIndex; i++ {
				require.Equal(t, data[i], bundle[i-startIndex])
			}
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

		isBlobInCorrectShard := false
		blobShards := shardMap[key]
		for _, shard := range blobShards {
			if _, ok := shardSet[shard]; ok {
				isBlobInCorrectShard = true
				break
			}
		}

		if isBlobInCorrectShard {
			response, err := getChunks(t, rand, operatorKeys, request)
			require.NoError(t, err)

			require.Equal(t, 1, len(response.Data))

			bundle, err := core.Bundle{}.Deserialize(response.Data[0])
			require.NoError(t, err)

			for i := 0; i < len(indices); i++ {
				if i%2 == 0 {
					require.Equal(t, data[indices[i]], bundle[i/2])
				}
			}
		} else {
			response, err := getChunks(t, rand, operatorKeys, request)
			require.Error(t, err)
			require.Nil(t, response)
		}
	}
}

func TestBatchedReadWriteChunksWithSharding(t *testing.T) {
	rand := random.NewTestRandom(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	// These are used to write data to S3/dynamoDB
	metadataStore := buildMetadataStore(t)
	chunkReader, chunkWriter := buildChunkStore(t, logger)

	shardCount := rand.Intn(10) + 10
	shardList := make([]v2.RelayKey, 0)
	shardSet := make(map[v2.RelayKey]struct{})
	for i := 0; i < shardCount; i++ {
		if i%2 == 0 {
			shardList = append(shardList, v2.RelayKey(i))
			shardSet[v2.RelayKey(i)] = struct{}{}
		}
	}
	shardMap := make(map[v2.BlobKey][]v2.RelayKey)

	operatorCount := rand.Intn(3) + 1
	operatorKeys := make(map[uint32]*core.KeyPair)
	operatorInfo := make(map[core.OperatorID]*core.IndexedOperatorInfo)
	for i := 0; i < operatorCount; i++ {
		keypair := rand.BLS()
		operatorKeys[uint32(i)] = keypair

		var operatorID core.OperatorID
		binary.BigEndian.PutUint32(operatorID[24:], uint32(i))
		operatorInfo[operatorID] = &core.IndexedOperatorInfo{
			PubkeyG1: keypair.GetPubKeyG1(),
			PubkeyG2: keypair.GetPubKeyG2(),
		}
	}

	ics := &mock.IndexedChainState{}
	blockNumber := uint(rand.Uint32())
	ics.Mock.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	ics.Mock.On("GetIndexedOperators", blockNumber).Return(operatorInfo, nil)

	// This is the server used to read it back
	config := defaultConfig()
	config.RelayIDs = shardList
	config.RateLimits.MaxGetChunkOpsPerSecond = 1000
	config.RateLimits.GetChunkOpsBurstiness = 1000
	config.RateLimits.MaxGetChunkOpsPerSecondClient = 1000
	config.RateLimits.GetChunkOpsBurstinessClient = 1000
	chainReader := newMockChainReader()
	server, err := NewServer(
		context.Background(),
		logger,
		config,
		metadataStore,
		nil, /* not used in this test */
		chunkReader,
		chainReader,
		ics)
	require.NoError(t, err)

	go func() {
		err = server.Start(context.Background())
		require.NoError(t, err)
	}()
	defer func() {
		err = server.Stop()
		require.NoError(t, err)
	}()

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
				TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
				FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	keyCount := 2

	// Read the blobs back. On average, we expect 25% of the blobs to be assigned to shards we don't have.
	for i := 0; i < 10; i++ {
		keys := make([]v2.BlobKey, 0, keyCount)
		for key := range expectedData {
			keys = append(keys, key)
			if len(keys) == keyCount {
				break
			}
		}

		requestedChunks := make([]*pb.ChunkRequest, 0)
		for _, key := range keys {

			boundKey := key
			request := &pb.ChunkRequest{
				Request: &pb.ChunkRequest_ByRange{
					ByRange: &pb.ChunkRequestByRange{
						BlobKey:    boundKey[:],
						StartIndex: 0,
						EndIndex:   uint32(len(expectedData[key])),
					},
				},
			}

			requestedChunks = append(requestedChunks, request)
		}
		// Add a request for duplicate key with different index range
		requestedChunks = append(requestedChunks, &pb.ChunkRequest{
			Request: &pb.ChunkRequest_ByRange{
				ByRange: &pb.ChunkRequestByRange{
					BlobKey:    keys[0][:],
					StartIndex: uint32(len(expectedData[keys[0]]) / 2),
					EndIndex:   uint32(len(expectedData[keys[0]])),
				},
			},
		})
		request := &pb.GetChunksRequest{
			ChunkRequests: requestedChunks,
		}

		allInCorrectShard := true
		for _, key := range keys {
			isBlobInCorrectShard := false
			blobShards := shardMap[key]
			for _, shard := range blobShards {
				if _, ok := shardSet[shard]; ok {
					isBlobInCorrectShard = true
					break
				}
			}
			if !isBlobInCorrectShard {
				allInCorrectShard = false
				break
			}
		}

		response, err := getChunks(t, rand, operatorKeys, request)
		if allInCorrectShard {
			require.NoError(t, err)

			require.Equal(t, keyCount+1, len(response.Data))

			for keyIndex, key := range keys {
				data := expectedData[key]

				bundle, err := core.Bundle{}.Deserialize(response.Data[keyIndex])
				require.NoError(t, err)

				for frameIndex, frame := range bundle {
					require.Equal(t, data[frameIndex], frame)
				}
			}

			// Check the duplicate key
			key := keys[0]
			data := expectedData[key][len(expectedData[key])/2:]

			bundle, err := core.Bundle{}.Deserialize(response.Data[keyCount])
			require.NoError(t, err)

			for frameIndex, frame := range bundle {
				require.Equal(t, data[frameIndex], frame)
			}
		} else {
			require.Error(t, err)
			require.Nil(t, response)
		}
	}
}
