package relay

import (
	"context"
	"fmt"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbcommonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	p "github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	prover             *p.Prover
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
)

func setup(t *testing.T) {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "..")
	changeDirectory(filepath.Join(rootPath, "inabox"))

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
		require.NoError(t, err)
	}

	// Only set up the prover once, it's expensive
	if prover == nil {
		config := &kzg.KzgConfig{
			G1Path:          "./resources/kzg/g1.point.300000",
			G2Path:          "./resources/kzg/g2.point.300000",
			CacheDir:        "./resources/kzg/SRSTables",
			SRSOrder:        8192,
			SRSNumberToLoad: 8192,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		}
		var err error
		prover, err = p.NewProver(config, true)
		require.NoError(t, err)
	}
}

func changeDirectory(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("Failed to change directories. Error: %s", err)
	}

	newDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get working directory. Error: %s", err)
	}
	log.Printf("Current Working Directory: %s\n", newDir)
}

func teardown() {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func buildMetadataStore(t *testing.T) *blobstore.BlobMetadataStore {
	setup(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
	require.NoError(t, err)
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
	require.NoError(t, err)

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     localstackHost,
	}

	_, err = test_utils.CreateTable(
		context.Background(),
		cfg,
		metadataTableName,
		blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	require.NoError(t, err)

	dynamoClient, err := dynamodb.NewClient(cfg, logger)
	require.NoError(t, err)

	return blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		metadataTableName)
}

func randomBlobHeader(t *testing.T) *v2.BlobHeader {

	data := tu.RandomBytes(128)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitments(data)
	require.NoError(t, err)
	require.NoError(t, err)
	commitmentProto, err := commitments.ToProfobuf()
	require.NoError(t, err)

	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         tu.RandomString(10),
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := v2.NewBlobHeader(blobHeaderProto)
	require.NoError(t, err)

	return blobHeader
}

// TODO verify blob size once it is added to metadata

func TestFetchingIndividualMetadata(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardown()
	}()

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header := randomBlobHeader(t)
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

	server, err := NewMetadataServer(context.Background(), logger, metadataStore, 1024*1024, 32, nil)
	require.NoError(t, err)

	// Fetch the metadata from the server.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		metadataMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})
		require.NoError(t, err)
		require.Equal(t, 1, len(*metadataMap))
		metadata := (*metadataMap)[blobKey]
		require.NotNil(t, metadata)
		require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
	}

	// Read it back again. This uses a different code pathway due to the cache.
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		metadataMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})
		require.NoError(t, err)
		require.Equal(t, 1, len(*metadataMap))
		metadata := (*metadataMap)[blobKey]
		require.NotNil(t, metadata)
		require.Equal(t, totalChunkSizeBytes, metadata.totalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], metadata.fragmentSizeBytes)
	}
}

func TestBatchedFetch(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardown()
	}()

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header := randomBlobHeader(t)
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

	server, err := NewMetadataServer(context.Background(), logger, metadataStore, 1024*1024, 32, nil)
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

		metadataMap, err := server.GetMetadataForBlobs(keys)
		require.NoError(t, err)

		assert.Equal(t, keyCount, len(*metadataMap))
		for _, key := range keys {
			metadata := (*metadataMap)[key]
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

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardown()
	}()

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
		header := randomBlobHeader(t)
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

	server, err := NewMetadataServer(context.Background(), logger, metadataStore, 1024*1024, 32, shardList)
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

		metadataMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})

		if isBlobInCorrectShard {
			// The blob is in the relay's shard, should be returned like normal
			require.NoError(t, err)
			require.Equal(t, 1, len(*metadataMap))
			metadata := (*metadataMap)[blobKey]
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

		metadataMap, err := server.GetMetadataForBlobs([]v2.BlobKey{blobKey})

		if isBlobInCorrectShard {
			// The blob is in the relay's shard, should be returned like normal
			require.NoError(t, err)
			require.Equal(t, 1, len(*metadataMap))
			metadata := (*metadataMap)[blobKey]
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

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardown()
	}()

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
		header := randomBlobHeader(t)
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

	server, err := NewMetadataServer(context.Background(), logger, metadataStore, 1024*1024, 32, shardList)
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

		metadataMap, err := server.GetMetadataForBlobs(keys)
		if areKeysInCorrectShard {
			require.NoError(t, err)
			assert.Equal(t, keyCount, len(*metadataMap))
			for _, key := range keys {
				metadata := (*metadataMap)[key]
				require.NotNil(t, metadata)
				require.Equal(t, totalChunkSizeMap[key], metadata.totalChunkSizeBytes)
				require.Equal(t, fragmentSizeMap[key], metadata.fragmentSizeBytes)
			}
		} else {
			require.Error(t, err)
		}
	}
}
