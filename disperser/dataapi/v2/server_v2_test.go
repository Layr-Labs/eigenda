package v2_test

import (
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	prommock "github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus/mock"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	subgraphmock "github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph/mock"
	serverv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/prometheus/common/model"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	//go:embed testdata/prometheus-response-sample.json
	mockPrometheusResponse string

	//go:embed testdata/prometheus-resp-avg-throughput.json
	mockPrometheusRespAvgThroughput string

	blobMetadataStore   *blobstorev2.BlobMetadataStore
	testDataApiServerV2 *serverv2.ServerV2

	logger = testutils.GetLogger()

	// Local stack
	localStackPort     = "4566"
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	deployLocalStack   bool

	mockLogger        = testutils.GetLogger()
	blobstore         = inmem.NewBlobStore()
	mockPrometheusApi = &prommock.MockPrometheusApi{}
	prometheusClient  = dataapi.NewPrometheusClient(mockPrometheusApi, "test-cluster")
	mockSubgraphApi   = &subgraphmock.MockSubgraphApi{}
	subgraphClient    = dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger)

	config = dataapi.Config{ServerMode: "test", SocketAddr: ":8080", AllowOrigins: []string{"*"}, DisperserHostname: "localhost:32007", ChurnerHostname: "localhost:32009"}

	mockTx            = &coremock.MockWriter{}
	opId0, _          = core.OperatorIDFromHex("e22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311")
	opId1, _          = core.OperatorIDFromHex("e23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312")
	mockChainState, _ = coremock.NewChainDataMock(map[uint8]map[core.OperatorID]int{
		0: {
			opId0: 1,
			opId1: 1,
		},
		1: {
			opId0: 1,
			opId1: 3,
		},
	})
	mockIndexedChainState, _ = coremock.MakeChainDataMock(map[uint8]int{
		0: 10,
		1: 10,
		2: 10,
	})
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil, nil)

	operatorInfo = &subgraph.IndexedOperatorInfo{
		Id:         "0xa96bfb4a7ca981ad365220f336dc5a3de0816ebd5130b79bbc85aca94bc9b6ac",
		PubkeyG1_X: "1336192159512049190945679273141887248666932624338963482128432381981287252980",
		PubkeyG1_Y: "25195175002875833468883745675063986308012687914999552116603423331534089122704",
		PubkeyG2_X: []graphql.String{
			"31597023645215426396093421944506635812143308313031252511177204078669540440732",
			"21405255666568400552575831267661419473985517916677491029848981743882451844775",
		},
		PubkeyG2_Y: []graphql.String{
			"8416989242565286095121881312760798075882411191579108217086927390793923664442",
			"23612061731370453436662267863740141021994163834412349567410746669651828926551",
		},
		SocketUpdates: []subgraph.SocketUpdates{
			{
				Socket: "23.93.76.1:32005;32006",
			},
		},
	}
)

type MockSubgraphClient struct {
	mock.Mock
}

type MockGRPCConnection struct{}

type MockHttpClient struct {
	ShouldSucceed bool
}

func (mc *MockGRPCConnection) Dial(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Here, return a mock connection. How you implement this depends on your testing framework
	// and what aspects of the gRPC connection you wish to mock.
	// For a simple approach, you might not even need to return a real *grpc.ClientConn
	// but rather a mock or stub that satisfies the interface.
	return &grpc.ClientConn{}, nil // Simplified, consider using a more sophisticated mock.
}

type MockGRPNilConnection struct{}

func (mc *MockGRPNilConnection) Dial(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Here, return a mock connection. How you implement this depends on your testing framework
	// and what aspects of the gRPC connection you wish to mock.
	// For a simple approach, you might not even need to return a real *grpc.ClientConn
	// but rather a mock or stub that satisfies the interface.
	return nil, nil // Simplified, consider using a more sophisticated mock.
}

type MockHealthCheckService struct {
	ResponseMap map[string]*grpc_health_v1.HealthCheckResponse
}

func TestMain(m *testing.M) {
	setup(m)
	m.Run()
	teardown()
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func setup(m *testing.M) {
	// Start localstack
	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}
	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	// Create DynamoDB table
	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	metadataTableName := fmt.Sprintf("test-BlobMetadata-%v", uuid.New())
	_, err := test_utils.CreateTable(context.Background(), cfg, metadataTableName, blobstorev2.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	// Create BlobMetadataStore
	dynamoClient, err := dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}
	blobMetadataStore = blobstorev2.NewBlobMetadataStore(dynamoClient, logger, metadataTableName)
	testDataApiServerV2 = serverv2.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger))
}

// makeCommitment returns a test hardcoded BlobCommitments
func makeCommitment(t *testing.T) encoding.BlobCommitments {
	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	require.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	require.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	require.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	require.NoError(t, err)

	var lengthProof bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	return encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: *new(fp.Element).SetBigInt(big.NewInt(1)),
			Y: *new(fp.Element).SetBigInt(big.NewInt(2)),
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthProof),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           16,
	}
}

// makeBlobHeaderV2 returns a test hardcoded V2 BlobHeader
func makeBlobHeaderV2(t *testing.T) *corev2.BlobHeader {
	accountBytes := make([]byte, 32)
	_, err := rand.Read(accountBytes)
	require.NoError(t, err)
	accountID := hex.EncodeToString(accountBytes)
	reservationPeriod, err := rand.Int(rand.Reader, big.NewInt(42))
	require.NoError(t, err)
	salt, err := rand.Int(rand.Reader, big.NewInt(1000))
	require.NoError(t, err)
	cumulativePayment, err := rand.Int(rand.Reader, big.NewInt(123))
	require.NoError(t, err)
	sig := make([]byte, 32)
	_, err = rand.Read(sig)
	require.NoError(t, err)
	return &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0, 1},
		BlobCommitments: makeCommitment(t),
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			ReservationPeriod: uint32(reservationPeriod.Int64()),
			CumulativePayment: cumulativePayment,
			Salt:              uint32(salt.Int64()),
		},
		Signature: sig,
	}
}

func setUpRouter() *gin.Engine {
	return gin.Default()
}

func executeRequest(t *testing.T, router *gin.Engine, method, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	return w
}

func decodeResponseBody[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	body := w.Result().Body
	defer body.Close()
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response T
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	return response
}

func checkBlobKeyEqual(t *testing.T, blobKey corev2.BlobKey, blobHeader *corev2.BlobHeader) {
	bk, err := blobHeader.BlobKey()
	assert.Nil(t, err)
	assert.Equal(t, blobKey, bk)
}

func checkPaginationToken(t *testing.T, token string, requestedAt uint64, blobKey corev2.BlobKey) {
	cursor, err := new(blobstorev2.BlobFeedCursor).FromCursorKey(token)
	require.NoError(t, err)
	assert.True(t, cursor.Equal(requestedAt, &blobKey))
}

func TestFetchBlobHandlerV2(t *testing.T) {
	r := setUpRouter()

	// Set up blob metadata in metadata store
	now := time.Now()
	blobHeader := makeBlobHeaderV2(t)
	metadata := &commonv2.BlobMetadata{
		BlobHeader: blobHeader,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(context.Background(), metadata)
	require.NoError(t, err)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)
	require.NoError(t, err)

	r.GET("/v2/blobs/:blob_key", testDataApiServerV2.FetchBlobHandler)

	w := executeRequest(t, r, http.MethodGet, "/v2/blobs/"+blobKey.Hex())
	response := decodeResponseBody[serverv2.BlobResponse](t, w)

	assert.Equal(t, "Queued", response.Status)
	assert.Equal(t, uint16(0), response.BlobHeader.BlobVersion)
	assert.Equal(t, blobHeader.Signature, response.BlobHeader.Signature)
	assert.Equal(t, blobHeader.PaymentMetadata.AccountID, response.BlobHeader.PaymentMetadata.AccountID)
	assert.Equal(t, blobHeader.PaymentMetadata.ReservationPeriod, response.BlobHeader.PaymentMetadata.ReservationPeriod)
	assert.Equal(t, blobHeader.PaymentMetadata.CumulativePayment, response.BlobHeader.PaymentMetadata.CumulativePayment)
}

func TestFetchBlobCertificateHandler(t *testing.T) {
	r := setUpRouter()

	// Set up blob certificate in metadata store
	blobHeader := makeBlobHeaderV2(t)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}
	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}
	err = blobMetadataStore.PutBlobCertificate(context.Background(), blobCert, fragmentInfo)
	require.NoError(t, err)

	r.GET("/v2/blobs/:blob_key/certificate", testDataApiServerV2.FetchBlobCertificateHandler)

	w := executeRequest(t, r, http.MethodGet, "/v2/blobs/"+blobKey.Hex()+"/certificate")
	response := decodeResponseBody[serverv2.BlobCertificateResponse](t, w)

	assert.Equal(t, blobCert.RelayKeys, response.Certificate.RelayKeys)
	assert.Equal(t, uint16(0), response.Certificate.BlobHeader.BlobVersion)
	assert.Equal(t, blobHeader.Signature, response.Certificate.BlobHeader.Signature)
}

func TestFetchBlobFeedHandler(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()

	// Create a timeline of test blobs:
	// - Total of 103 blobs
	// - First 3 blobs share the same timestamp (firstBlobTime)
	// - The last blob has timestamp "now"
	// - Remaining blobs are spaced 1 minute apart
	// - Timeline spans roughly 100 minutes into the past from now
	numBlobs := 103
	now := uint64(time.Now().UnixNano())
	nanoSecsPerBlob := uint64(60 * 1e9) // 1 blob per minute
	firstBlobTime := now - uint64(numBlobs-3)*nanoSecsPerBlob
	keys := make([]corev2.BlobKey, numBlobs)
	requestedAt := make([]uint64, numBlobs)

	// Actually create blobs
	firstBlobKeys := make([][32]byte, 3)
	for i := 0; i < numBlobs; i++ {
		blobHeader := makeBlobHeaderV2(t)
		blobKey, err := blobHeader.BlobKey()
		require.NoError(t, err)
		keys[i] = blobKey
		if i < 3 {
			requestedAt[i] = firstBlobTime
			firstBlobKeys[i] = keys[i]
		} else {
			requestedAt[i] = firstBlobTime + nanoSecsPerBlob*uint64(i-2)
		}

		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader:  blobHeader,
			BlobStatus:  v2.Encoded,
			Expiry:      uint64(now.Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   uint64(now.UnixNano()),
			RequestedAt: requestedAt[i],
		}
		err = blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
	}
	sort.Slice(firstBlobKeys, func(i, j int) bool {
		return bytes.Compare(firstBlobKeys[i][:], firstBlobKeys[j][:]) < 0
	})

	r.GET("/v2/blobs/feed", testDataApiServerV2.FetchBlobFeedHandler)

	t.Run("invalid params", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v2/blobs/feed?pagination_token=abc", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/v2/blobs/feed?limit=abc", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/v2/blobs/feed?interval=abc", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/v2/blobs/feed?end=2006-01-02T15:04:05", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("default params", func(t *testing.T) {
		// Default query returns:
		// - Most recent 1 hour of blobs (60 blobs total available, keys[43], ..., keys[102])
		// - Limited to 20 results (the default "limit")
		// - Starting from blob[43] through blob[62]
		w := executeRequest(t, r, http.MethodGet, "/v2/blobs/feed")
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[62], keys[62])
	})

	t.Run("various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// Returns keys[43] through keys[102] (60 blobs)
		w := executeRequest(t, r, http.MethodGet, "/v2/blobs/feed?limit=0")
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[102], keys[102])

		// Test 2: 2-hour window captures all test blobs
		// Verifies correct ordering of timestamp-colliding blobs
		w = executeRequest(t, r, http.MethodGet, "/v2/blobs/feed?interval=7200&limit=-1")
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, numBlobs, len(response.Blobs))
		// First 3 blobs ordered by key due to same timestamp
		checkBlobKeyEqual(t, firstBlobKeys[0], response.Blobs[0].BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[1], response.Blobs[1].BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[2], response.Blobs[2].BlobHeader)
		for i := 3; i < numBlobs; i++ {
			checkBlobKeyEqual(t, keys[i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[102], keys[102])

		// Test 3: Custom end time with 1-hour window
		// Retrieves keys[41] through keys[100]
		tm := time.Unix(0, int64(requestedAt[100])+1).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl := fmt.Sprintf("/v2/blobs/feed?end=%s&limit=-1", endTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[41+i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[41+i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[100], keys[100])
	})

	t.Run("pagination", func(t *testing.T) {
		// Test pagination behavior:
		// 1. First page: blobs in past 1h limited to 20, returns keys[43] through keys[62]
		// 2. Second page: the next 20 blobs, returns keys[63] through keys[82]
		// Verifies:
		// - Correct sequencing across pages
		// - Proper token handling
		tm := time.Unix(0, time.Now().UnixNano()).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/blobs/feed?end=%s&limit=20", endTime)
		w := executeRequest(t, r, http.MethodGet, reqUrl)
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[62], keys[62])

		// Request next page using pagination token
		reqUrl = fmt.Sprintf("/v2/blobs/feed?end=%s&limit=20&pagination_token=%s", endTime, response.PaginationToken)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[63+i], response.Blobs[i].BlobHeader)
			assert.Equal(t, requestedAt[63+i], response.Blobs[i].RequestedAt)
		}
		assert.True(t, len(response.PaginationToken) > 0)
		checkPaginationToken(t, response.PaginationToken, requestedAt[82], keys[82])
	})
}

func TestFetchBlobVerificationInfoHandler(t *testing.T) {
	r := setUpRouter()

	// Set up blob verification info in metadata store
	blobHeader := makeBlobHeaderV2(t)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)

	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	batchHeaderHash, err := batchHeader.Hash()
	require.NoError(t, err)

	ctx := context.Background()
	err = blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	require.NoError(t, err)
	verificationInfo := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      123,
		InclusionProof: []byte("inclusion proof"),
	}
	err = blobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo)
	require.NoError(t, err)

	r.GET("/v2/blobs/:blob_key/verification-info", testDataApiServerV2.FetchBlobVerificationInfoHandler)

	reqStr := fmt.Sprintf("/v2/blobs/%s/verification-info?batch_header_hash=%s", blobKey.Hex(), hex.EncodeToString(batchHeaderHash[:]))
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[serverv2.BlobVerificationInfoResponse](t, w)

	assert.Equal(t, verificationInfo.InclusionProof, response.VerificationInfo.InclusionProof)
}

func TestFetchBatchHandlerV2(t *testing.T) {
	r := setUpRouter()

	// Set up batch header in metadata store
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 0, 2, 4},
		ReferenceBlockNumber: 1024,
	}
	err := blobMetadataStore.PutBatchHeader(context.Background(), batchHeader)
	require.NoError(t, err)
	batchHeaderHashBytes, err := batchHeader.Hash()
	require.NoError(t, err)
	batchHeaderHash := hex.EncodeToString(batchHeaderHashBytes[:])

	// Set up attestation in metadata store
	commitment := makeCommitment(t)
	attestation := &corev2.Attestation{
		BatchHeader: batchHeader,
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(1), big.NewInt(0)),
			core.NewG1Point(big.NewInt(2), big.NewInt(4)),
		},
		APKG2: &core.G2Point{
			G2Affine: &bn254.G2Affine{
				X: commitment.LengthCommitment.X,
				Y: commitment.LengthCommitment.Y,
			},
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(2), big.NewInt(0)),
		},
	}
	err = blobMetadataStore.PutAttestation(context.Background(), attestation)
	require.NoError(t, err)

	r.GET("/v2/batches/:batch_header_hash", testDataApiServerV2.FetchBatchHandler)

	w := executeRequest(t, r, http.MethodGet, "/v2/batches/"+batchHeaderHash)
	response := decodeResponseBody[serverv2.BatchResponse](t, w)

	assert.Equal(t, batchHeaderHash, response.BatchHeaderHash)
	assert.Equal(t, batchHeader.BatchRoot, response.SignedBatch.BatchHeader.BatchRoot)
	assert.Equal(t, batchHeader.ReferenceBlockNumber, response.SignedBatch.BatchHeader.ReferenceBlockNumber)
	assert.Equal(t, attestation.AttestedAt, response.SignedBatch.Attestation.AttestedAt)
	assert.Equal(t, attestation.QuorumNumbers, response.SignedBatch.Attestation.QuorumNumbers)
}

func TestCheckOperatorsReachability(t *testing.T) {
	r := setUpRouter()

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil

	operatorId := "0xa96bfb4a7ca981ad365220f336dc5a3de0816ebd5130b79bbc85aca94bc9b6ab"
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(operatorInfo, nil)

	r.GET("/v2/operators/reachability", testDataApiServerV2.CheckOperatorsReachability)

	reqStr := fmt.Sprintf("/v2/operators/reachability?operator_id=%v", operatorId)
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[dataapi.OperatorPortCheckResponse](t, w)

	assert.Equal(t, "23.93.76.1:32005", response.DispersalSocket)
	assert.Equal(t, false, response.DispersalOnline)
	assert.Equal(t, "23.93.76.1:32006", response.RetrievalSocket)
	assert.Equal(t, false, response.RetrievalOnline)

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchOperatorResponses(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()
	// Set up batch header in metadata store
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 0, 2, 4},
		ReferenceBlockNumber: 1024,
	}
	batchHeaderHashBytes, err := batchHeader.Hash()
	require.NoError(t, err)
	batchHeaderHash := hex.EncodeToString(batchHeaderHashBytes[:])

	// Set up dispersal response in metadata store
	operatorId := core.OperatorID{0, 1}
	dispersalRequest := &corev2.DispersalRequest{
		OperatorID:      operatorId,
		OperatorAddress: gethcommon.HexToAddress("0x1234567"),
		Socket:          "socket",
		DispersedAt:     uint64(time.Now().UnixNano()),
		BatchHeader:     *batchHeader,
	}
	dispersalResponse := &corev2.DispersalResponse{
		DispersalRequest: dispersalRequest,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        [32]byte{1, 1, 1},
		Error:            "error",
	}
	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
	assert.NoError(t, err)

	// Set up the other dispersal response in metadata store
	operatorId2 := core.OperatorID{2, 3}
	dispersalRequest2 := &corev2.DispersalRequest{
		OperatorID:      operatorId2,
		OperatorAddress: gethcommon.HexToAddress("0x1234567"),
		Socket:          "socket",
		DispersedAt:     uint64(time.Now().UnixNano()),
		BatchHeader:     *batchHeader,
	}
	dispersalResponse2 := &corev2.DispersalResponse{
		DispersalRequest: dispersalRequest2,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        [32]byte{1, 1, 1},
		Error:            "",
	}
	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse2)
	assert.NoError(t, err)

	r.GET("/v2/operators/response/:batch_header_hash", testDataApiServerV2.FetchOperatorsResponses)

	// Fetch response of a specific operator
	reqStr := fmt.Sprintf("/v2/operators/response/%s?operator_id=%v", batchHeaderHash, operatorId.Hex())
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[serverv2.OperatorDispersalResponses](t, w)
	require.Equal(t, 1, len(response.Responses))
	require.Equal(t, dispersalResponse, response.Responses[0])

	// Fetch all operators' responses for a batch
	reqStr2 := fmt.Sprintf("/v2/operators/response/%s", batchHeaderHash)
	w2 := executeRequest(t, r, http.MethodGet, reqStr2)
	response2 := decodeResponseBody[serverv2.OperatorDispersalResponses](t, w2)

	require.Equal(t, 2, len(response2.Responses))
	require.Equal(t, response2.Responses[0], dispersalResponse)
	require.Equal(t, response2.Responses[1], dispersalResponse2)
}

func TestFetchOperatorsStake(t *testing.T) {
	r := setUpRouter()

	mockIndexedChainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	r.GET("/v2/operators/stake", testDataApiServerV2.FetchOperatorsStake)

	w := executeRequest(t, r, http.MethodGet, "/v2/operators/stake")
	response := decodeResponseBody[dataapi.OperatorsStakeResponse](t, w)

	// The quorums and the operators in the quorum are defined in "mockChainState"
	// There are 3 quorums (0, 1) and a "total" entry for TotalQuorumStake
	assert.Equal(t, 3, len(response.StakeRankedOperators))
	// Quorum 0
	ops, ok := response.StakeRankedOperators["0"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(ops))
	assert.Equal(t, opId0.Hex(), ops[0].OperatorId)
	assert.Equal(t, opId1.Hex(), ops[1].OperatorId)
	// Quorum 1
	ops, ok = response.StakeRankedOperators["1"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(ops))
	assert.Equal(t, opId1.Hex(), ops[0].OperatorId)
	assert.Equal(t, opId0.Hex(), ops[1].OperatorId)
	// "total"
	ops, ok = response.StakeRankedOperators["total"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(ops))
	assert.Equal(t, opId1.Hex(), ops[0].OperatorId)
	assert.Equal(t, opId0.Hex(), ops[1].OperatorId)
}

func TestFetchMetricsSummaryHandler(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusResponse))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v2/metrics/summary", testDataApiServerV2.FetchMetricsSummaryHandler)

	w := executeRequest(t, r, http.MethodGet, "/v2/metrics/summary")
	response := decodeResponseBody[serverv2.MetricSummary](t, w)

	assert.Equal(t, 16555.555555555555, response.AvgThroughput)
}

func TestFetchMetricsThroughputTimeseriesHandler(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusRespAvgThroughput))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v2/metrics/timeseries/throughput", testDataApiServer.FetchMetricsThroughputHandler)

	w := executeRequest(t, r, http.MethodGet, "/v2/metrics/timeseries/throughput")
	response := decodeResponseBody[[]*dataapi.Throughput](t, w)

	var totalThroughput float64
	for _, v := range response {
		totalThroughput += v.Throughput
	}

	assert.Equal(t, 3361, len(response))
	assert.Equal(t, float64(12000), response[0].Throughput)
	assert.Equal(t, uint64(1701292920), response[0].Timestamp)
	assert.Equal(t, float64(3.503022666666651e+07), totalThroughput)
}
