package v2_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	prommock "github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus/mock"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	subgraphmock "github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph/mock"
	serverv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	//go:embed testdata/prometheus-resp-avg-throughput.json
	mockPrometheusRespAvgThroughput string

	//go:embed testdata/prometheus-response-network-signing-rate.json
	mockPrometheusResponseNetworkSigningRate string

	UUID                = uuid.New()
	metadataTableName   = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	blobMetadataStore   *blobstorev2.BlobMetadataStore
	testDataApiServerV2 *serverv2.ServerV2

	logger = testutils.GetLogger()

	// Local stack
	localStackPort     = "4566"
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	deployLocalStack   bool

	dynamoClient dynamodb.Client

	serverVersion     = uint(2)
	mockLogger        = testutils.GetLogger()
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
	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
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
	_, err := test_utils.CreateTable(context.Background(), cfg, metadataTableName, blobstorev2.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	// Create BlobMetadataStore
	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}
	blobMetadataStore = blobstorev2.NewBlobMetadataStore(dynamoClient, logger, metadataTableName)

	mockTx.On("GetCurrentBlockNumber").Return(uint32(1), nil)
	mockTx.On("GetQuorumCount").Return(uint8(2), nil)

	metrics := dataapi.NewMetrics(serverVersion, nil, blobMetadataStore, "9001", mockLogger)
	testDataApiServerV2, err = serverv2.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, metrics)
	if err != nil {
		teardown()
		panic("failed to create v2 server: " + err.Error())
	}
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
	accountID := gethcommon.HexToAddress(hex.EncodeToString(accountBytes))
	timestamp, err := rand.Int(rand.Reader, big.NewInt(int64(time.Now().Nanosecond())))
	require.NoError(t, err)
	cumulativePayment, err := rand.Int(rand.Reader, big.NewInt(int64(time.Now().Nanosecond())))
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
			Timestamp:         timestamp.Int64(),
			CumulativePayment: cumulativePayment,
		},
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
	defer core.CloseLogOnError(body, "response body", mockLogger)
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

func checkOperatorSigningInfoEqual(t *testing.T, actual, expected *serverv2.OperatorSigningInfo) {
	assert.Equal(t, expected.OperatorId, actual.OperatorId)
	assert.Equal(t, expected.OperatorAddress, actual.OperatorAddress)
	assert.Equal(t, expected.QuorumId, actual.QuorumId)
	assert.Equal(t, expected.TotalUnsignedBatches, actual.TotalUnsignedBatches)
	assert.Equal(t, expected.TotalResponsibleBatches, actual.TotalResponsibleBatches)
	assert.Equal(t, expected.TotalBatches, actual.TotalBatches)
}

func checkCursor(t *testing.T, token string, requestedAt uint64, blobKey corev2.BlobKey) {
	cursor, err := new(blobstorev2.BlobFeedCursor).FromCursorKey(token)
	require.NoError(t, err)
	assert.True(t, cursor.Equal(requestedAt, &blobKey))
}

func deleteItems(t *testing.T, keys []dynamodb.Key) {
	failed, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
	assert.Len(t, failed, 0)
}

func TestFetchBlob(t *testing.T) {
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

	r.GET("/v2/blobs/:blob_key", testDataApiServerV2.FetchBlob)

	w := executeRequest(t, r, http.MethodGet, "/v2/blobs/"+blobKey.Hex())
	response := decodeResponseBody[serverv2.BlobResponse](t, w)

	assert.Equal(t, "Queued", response.Status)
	assert.Equal(t, uint16(0), response.BlobHeader.BlobVersion)
	assert.Equal(t, blobHeader.PaymentMetadata.AccountID, response.BlobHeader.PaymentMetadata.AccountID)
	assert.Equal(t, blobHeader.PaymentMetadata.Timestamp, response.BlobHeader.PaymentMetadata.Timestamp)
	assert.Equal(t, blobHeader.PaymentMetadata.CumulativePayment, response.BlobHeader.PaymentMetadata.CumulativePayment)
}

func TestFetchOperatorDispersalFeed(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()

	numRequests := 60
	opID := core.OperatorID{16, 32}
	now := uint64(time.Now().UnixNano())
	firstRequestTs := now - uint64(int64(numRequests)*time.Minute.Nanoseconds())
	nanoSecsPerRequest := uint64(time.Minute.Nanoseconds()) // 1 batch/min

	dispersedAt := make([]uint64, numRequests)
	batchHeaders := make([]*corev2.BatchHeader, numRequests)
	signatures := make([][32]byte, numRequests)
	dynamoKeys := make([]dynamodb.Key, numRequests)
	for i := 0; i < numRequests; i++ {
		dispersedAt[i] = firstRequestTs + uint64(i)*nanoSecsPerRequest
		batchHeaders[i] = &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: uint64(i + 100),
		}
		dispersalRequest := &corev2.DispersalRequest{
			OperatorID:      opID,
			OperatorAddress: gethcommon.HexToAddress("0x1234567"),
			Socket:          "socket",
			DispersedAt:     dispersedAt[i],
			BatchHeader:     *batchHeaders[i],
		}
		signatures[i] = [32]byte{}
		if i%2 == 0 {
			signatures[i] = [32]byte{1, 1, uint8(i)}
		}
		dispersalResponse := &corev2.DispersalResponse{
			DispersalRequest: dispersalRequest,
			RespondedAt:      dispersedAt[i],
			Signature:        signatures[i],
			Error:            "",
		}

		err := blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
		require.NoError(t, err)

		bhh, err := dispersalRequest.BatchHeader.Hash() // go:nolint QF1008
		require.NoError(t, err)
		dynamoKeys[i] = dynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalResponse#" + opID.Hex()},
		}
	}
	defer deleteItems(t, dynamoKeys)

	r.GET("/v2/operators/:operator_id/dispersals", testDataApiServerV2.FetchOperatorDispersalFeed)
	baseUrl := fmt.Sprintf("/v2/operators/%s/dispersals", opID.Hex())

	t.Run("invalid params", func(t *testing.T) {
		now := time.Now()

		tests := []struct {
			name        string
			queryParams map[string]string
			wantError   string // expected error message
		}{
			// Invalid direction
			{
				name:        "invalid direction",
				queryParams: map[string]string{"direction": "abc"},
				wantError:   "`direction` must be either \"forward\" or \"backward\", found: \"abc\"",
			},

			// Invalid time formats
			{
				name:        "invalid before format",
				queryParams: map[string]string{"before": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid before value",
				queryParams: map[string]string{"before": "abc"},
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid after format",
				queryParams: map[string]string{"after": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "invalid after value",
				queryParams: map[string]string{"after": "abc"},
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "after in future",
				queryParams: map[string]string{"after": "3025-01-02T15:04:05Z"},
				wantError:   "`after` must be before current time",
			},

			// Invalid time ranges
			{
				name: "after >= before",
				queryParams: map[string]string{
					"after":  now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.999999999Z"),
					"before": now.Add(-time.Hour).UTC().Format("2006-01-02T15:04:05.999999999Z"),
				},
				wantError: "must be earlier than `before` timestamp",
			},
			{
				name: "before too old",
				queryParams: map[string]string{
					"before": "2020-01-02T15:04:05Z",
				},
				wantError: "`before` time cannot be more than 14 days in the past",
			},

			// Invalid limit
			{
				name:        "invalid limit format",
				queryParams: map[string]string{"limit": "abc"},
				wantError:   "failed to parse `limit` param",
			},
		}

		for _, tt := range tests {
			params := url.Values{}
			for k, v := range tt.queryParams {
				params.Add(k, v)
			}

			url := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

			var errResp serverv2.ErrorResponse
			require.NoError(t, json.NewDecoder(w.Body).Decode(&errResp))
			assert.Contains(t, errResp.Error, tt.wantError)
		}
	})

	t.Run("nonexistent operatorid", func(t *testing.T) {
		otherID := core.OperatorID{4, 16}
		url := fmt.Sprintf("/v2/operators/%s/dispersals", otherID.Hex())
		w := executeRequest(t, r, http.MethodGet, url)
		response := decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 0, len(response.Dispersals))
	})

	t.Run("default params", func(t *testing.T) {
		// Default query returns:
		// - Most recent 1 hour of dispersals include all of dispersals[1] through dispersals[59]
		// - Limited to 20 results (the default "limit")
		// - Result will first 20 dispersals
		w := executeRequest(t, r, http.MethodGet, baseUrl)
		response := decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 20, len(response.Dispersals))
		for i := 0; i < 20; i++ {
			assert.Equal(t, dispersedAt[1+i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[1+i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[1+i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
			if (1+i)%2 == 0 {
				assert.Equal(t, hex.EncodeToString(signatures[1+i][:]), response.Dispersals[i].Signature)
			} else {
				assert.Equal(t, "", response.Dispersals[i].Signature)
			}
		}
	})

	t.Run("forward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves dispersals[1] through batch[59] (59 batches)
		w := executeRequest(t, r, http.MethodGet, baseUrl+"?limit=0")
		response := decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 59, len(response.Dispersals))
		for i := 0; i < 59; i++ {
			assert.Equal(t, dispersedAt[1+i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[1+i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[1+i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}

		// Test 2: 2-hour window captures all test batches
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("%s?limit=-1&after=%s", baseUrl, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 60, len(response.Dispersals))
		for i := 0; i < 60; i++ {
			assert.Equal(t, dispersedAt[i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}

		// Teste 3: custom end time
		after := time.Unix(0, int64(dispersedAt[20])).UTC()
		afterTime = after.Format("2006-01-02T15:04:05.999999999Z")
		before := time.Unix(0, int64(dispersedAt[50])).UTC()
		beforeTime := before.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("%s?before=%s&after=%s&limit=-1", baseUrl, beforeTime, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 29, len(response.Dispersals))
		for i := 0; i < 29; i++ {
			assert.Equal(t, dispersedAt[21+i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[21+i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[21+i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}
	})

	t.Run("backward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves dispersals[59] through batch[1] (59 batches)
		w := executeRequest(t, r, http.MethodGet, baseUrl+"?limit=0&direction=backward")
		response := decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 59, len(response.Dispersals))
		for i := 0; i < 59; i++ {
			assert.Equal(t, dispersedAt[59-i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[59-i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[59-i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}

		// Test 2: 2-hour window captures all test batches
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("%s?limit=-1&after=%s&direction=backward", baseUrl, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 60, len(response.Dispersals))
		for i := 0; i < 60; i++ {
			assert.Equal(t, dispersedAt[59-i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[59-i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[59-i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}

		// Teste 3: custom end time
		after := time.Unix(0, int64(dispersedAt[20])).UTC()
		afterTime = after.Format("2006-01-02T15:04:05.999999999Z")
		before := time.Unix(0, int64(dispersedAt[50])).UTC()
		beforeTime := before.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("%s?before=%s&after=%s&limit=-1&direction=backward", baseUrl, beforeTime, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.OperatorDispersalFeedResponse](t, w)
		require.Equal(t, 29, len(response.Dispersals))
		for i := 0; i < 29; i++ {
			assert.Equal(t, dispersedAt[49-i], response.Dispersals[i].DispersedAt)
			assert.Equal(t, batchHeaders[49-i].ReferenceBlockNumber, response.Dispersals[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[49-i].BatchRoot[:]), response.Dispersals[i].BatchHeader.BatchRoot)
		}
	})

}

func TestFetchBlobCertificate(t *testing.T) {
	r := setUpRouter()

	// Set up blob certificate in metadata store
	blobHeader := makeBlobHeaderV2(t)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		Signature:  []byte{0, 1, 2, 3, 4},
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}
	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}
	err = blobMetadataStore.PutBlobCertificate(context.Background(), blobCert, fragmentInfo)
	require.NoError(t, err)

	r.GET("/v2/blobs/:blob_key/certificate", testDataApiServerV2.FetchBlobCertificate)

	w := executeRequest(t, r, http.MethodGet, "/v2/blobs/"+blobKey.Hex()+"/certificate")
	response := decodeResponseBody[serverv2.BlobCertificateResponse](t, w)

	assert.Equal(t, blobCert.RelayKeys, response.Certificate.RelayKeys)
	assert.Equal(t, uint16(0), response.Certificate.BlobHeader.BlobVersion)
	assert.Equal(t, blobCert.Signature, response.Certificate.Signature)
}

func TestFetchBlobFeed(t *testing.T) {
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
	dynamoKeys := make([]dynamodb.Key, numBlobs)
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
		metadata := &commonv2.BlobMetadata{
			BlobHeader:  blobHeader,
			Signature:   []byte{0, 1, 2, 3, 4},
			BlobStatus:  commonv2.Encoded,
			Expiry:      uint64(now.Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   uint64(now.UnixNano()),
			RequestedAt: requestedAt[i],
		}
		err = blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		dynamoKeys[i] = dynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
	}
	sort.Slice(firstBlobKeys, func(i, j int) bool {
		return bytes.Compare(firstBlobKeys[i][:], firstBlobKeys[j][:]) < 0
	})

	defer deleteItems(t, dynamoKeys)

	r.GET("/v2/blobs/feed", testDataApiServerV2.FetchBlobFeed)

	t.Run("invalid params", func(t *testing.T) {
		now := time.Now()

		tests := []struct {
			name        string
			queryParams map[string]string
			wantError   string // expected error message
		}{
			// Invalid direction
			{
				name:        "invalid direction",
				queryParams: map[string]string{"direction": "abc"},
				wantError:   "`direction` must be either \"forward\" or \"backward\", found: \"abc\"",
			},

			// Invalid time formats
			{
				name:        "invalid before format",
				queryParams: map[string]string{"before": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid before value",
				queryParams: map[string]string{"before": "abc"},
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid after format",
				queryParams: map[string]string{"after": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "invalid after value",
				queryParams: map[string]string{"after": "abc"},
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "after in future",
				queryParams: map[string]string{"after": "3025-01-02T15:04:05Z"},
				wantError:   "`after` must be before current time",
			},

			// Invalid time ranges
			{
				name: "after >= before",
				queryParams: map[string]string{
					"after":  now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.999999999Z"),
					"before": now.Add(-time.Hour).UTC().Format("2006-01-02T15:04:05.999999999Z"),
				},
				wantError: "must be earlier than `before` timestamp",
			},
			{
				name: "before too old",
				queryParams: map[string]string{
					"before": "2020-01-02T15:04:05Z",
				},
				wantError: "`before` time cannot be more than 14 days in the past",
			},

			// Invalid cursor
			{
				name:        "invalid cursor format",
				queryParams: map[string]string{"cursor": "not-a-valid-cursor"},
				wantError:   "failed to parse the `cursor`",
			},

			// Invalid limit
			{
				name:        "invalid limit format",
				queryParams: map[string]string{"limit": "abc"},
				wantError:   "failed to parse `limit` param",
			},
		}

		for _, tt := range tests {
			params := url.Values{}
			for k, v := range tt.queryParams {
				params.Add(k, v)
			}
			url := fmt.Sprintf("/v2/blobs/feed?%s", params.Encode())

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

			var errResp serverv2.ErrorResponse
			require.NoError(t, json.NewDecoder(w.Body).Decode(&errResp))
			assert.Contains(t, errResp.Error, tt.wantError)
		}
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
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[62], keys[62])
	})

	t.Run("forward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// Returns keys[43] through keys[102] (60 blobs)
		w := executeRequest(t, r, http.MethodGet, "/v2/blobs/feed?limit=0")
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[102], keys[102])

		// Test 2: 2-hour window captures all test blobs
		// Verifies correct ordering of timestamp-colliding blobs
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/blobs/feed?after=%s&limit=-1", afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, numBlobs, len(response.Blobs))
		// First 3 blobs ordered by key due to same timestamp
		checkBlobKeyEqual(t, firstBlobKeys[0], response.Blobs[0].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[1], response.Blobs[1].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[2], response.Blobs[2].BlobMetadata.BlobHeader)
		for i := 3; i < numBlobs; i++ {
			checkBlobKeyEqual(t, keys[i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[102], keys[102])

		// Test 3: Custom end time with 1-hour window
		// Retrieves keys[41] through keys[100]
		tm := time.Unix(0, int64(requestedAt[100])+1).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("/v2/blobs/feed?before=%s&limit=-1", endTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[41+i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[41+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[100], keys[100])
	})

	t.Run("backward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// Returns keys[102] through keys[43] (60 blobs in descending order of time)
		w := executeRequest(t, r, http.MethodGet, "/v2/blobs/feed?direction=backward&limit=0")
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[102-i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[102-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[43], keys[43])

		// Test 2: 2-hour window captures all test blobs
		// Verifies correct ordering of timestamp-colliding blobs
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/blobs/feed?direction=backward&after=%s&limit=-1", afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, numBlobs, len(response.Blobs))
		// The last 3 blobs ordered by key due to same timestamp
		checkBlobKeyEqual(t, firstBlobKeys[2], response.Blobs[numBlobs-3].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[1], response.Blobs[numBlobs-2].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[0], response.Blobs[numBlobs-1].BlobMetadata.BlobHeader)
		for i := 3; i < numBlobs; i++ {
			checkBlobKeyEqual(t, keys[i], response.Blobs[numBlobs-i-1].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[i], response.Blobs[numBlobs-i-1].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[0], firstBlobKeys[0])

		// Test 3: Custom end time with 1-hour window
		// Retrieves keys[100] through keys[41]
		tm := time.Unix(0, int64(requestedAt[100])+1).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("/v2/blobs/feed?direction=backward&before=%s&limit=-1", endTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			checkBlobKeyEqual(t, keys[100-i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[100-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[41], keys[41])
	})

	t.Run("forward pagination", func(t *testing.T) {
		// Test pagination behavior:
		// 1. First page: blobs in past 1h limited to 20, returns keys[43] through keys[62]
		// 2. Second page: the next 20 blobs, returns keys[63] through keys[82]
		// Verifies:
		// - Correct sequencing across pages
		// - Proper token handling
		tm := time.Unix(0, time.Now().UnixNano()).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/blobs/feed?before=%s&limit=20", endTime)
		w := executeRequest(t, r, http.MethodGet, reqUrl)
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[43+i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[43+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[62], keys[62])

		// Request next page using pagination cursor
		reqUrl = fmt.Sprintf("/v2/blobs/feed?before=%s&limit=20&cursor=%s", endTime, response.Cursor)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[63+i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[63+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[82], keys[82])
	})

	t.Run("backward pagination", func(t *testing.T) {
		// Test backward pagination behavior:
		// 1. First page: the most recent 20 blobs, keys[102] through keys[83]
		// 2. Second page: requesting the next 20 blobs, but only 3 blobs due to "after" time bound
		// Verifies:
		// - Correct sequencing across pages
		// - Proper token handling (cursor is exclusive)
		tm := time.Unix(0, int64(requestedAt[80])).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/blobs/feed?direction=backward&after=%s&limit=20", endTime)
		w := executeRequest(t, r, http.MethodGet, reqUrl)
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			checkBlobKeyEqual(t, keys[102-i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[102-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[83], keys[83])

		// Request next page using pagination cursor
		reqUrl = fmt.Sprintf("/v2/blobs/feed?direction=backward&after=%s&limit=20&cursor=%s", endTime, response.Cursor)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 3, len(response.Blobs))
		for i := 0; i < 3; i++ {
			checkBlobKeyEqual(t, keys[82-i], response.Blobs[i].BlobMetadata.BlobHeader)
			assert.Equal(t, requestedAt[82-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[80], keys[80])
	})

	t.Run("pagination over same-timestamp blobs", func(t *testing.T) {
		// Test pagination behavior in case of same-timestamp blobs
		// - We have 3 blobs with identical timestamp (firstBlobTime): firstBlobKeys[0,1,2]
		// - These are followed by sequential blobs: keys[3,4] with different timestamps
		// - End time is set to requestedAt[5]
		tm := time.Unix(0, int64(requestedAt[5])).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z") // nano precision format

		// First page: fetch 2 blobs, which have same requestedAt timestamp
		reqUrl := fmt.Sprintf("/v2/blobs/feed?before=%s&limit=2", endTime)
		w := executeRequest(t, r, http.MethodGet, reqUrl)
		response := decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		require.Equal(t, 2, len(response.Blobs))
		checkBlobKeyEqual(t, firstBlobKeys[0], response.Blobs[0].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, firstBlobKeys[1], response.Blobs[1].BlobMetadata.BlobHeader)
		assert.Equal(t, firstBlobTime, response.Blobs[0].BlobMetadata.RequestedAt)
		assert.Equal(t, firstBlobTime, response.Blobs[1].BlobMetadata.RequestedAt)
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[1], firstBlobKeys[1])

		// Second page: fetch remaining blobs (limit=0 means no limit, hence reach the last blob)
		reqUrl = fmt.Sprintf("/v2/blobs/feed?before=%s&limit=0&cursor=%s", endTime, response.Cursor)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BlobFeedResponse](t, w)
		// Verify second page contains:
		// 1. Last same-timestamp blob
		// 2. Following blobs with sequential timestamps
		require.Equal(t, 3, len(response.Blobs))
		checkBlobKeyEqual(t, firstBlobKeys[2], response.Blobs[0].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, keys[3], response.Blobs[1].BlobMetadata.BlobHeader)
		checkBlobKeyEqual(t, keys[4], response.Blobs[2].BlobMetadata.BlobHeader)
		assert.Equal(t, firstBlobTime, response.Blobs[0].BlobMetadata.RequestedAt)
		assert.Equal(t, requestedAt[3], response.Blobs[1].BlobMetadata.RequestedAt)
		assert.Equal(t, requestedAt[4], response.Blobs[2].BlobMetadata.RequestedAt)
		assert.True(t, len(response.Cursor) > 0)
		checkCursor(t, response.Cursor, requestedAt[4], keys[4])
	})
}

func TestFetchBlobAttestationInfo(t *testing.T) {
	ctx := context.Background()
	r := setUpRouter()

	// Set up blob inclusion info
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
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 10,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)
	err = blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	require.NoError(t, err)
	inclusionInfo := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      123,
		InclusionProof: []byte("inclusion proof"),
	}
	err = blobMetadataStore.PutBlobInclusionInfo(ctx, inclusionInfo)
	require.NoError(t, err)

	r.GET("/v2/blobs/:blob_key/attestation-info", testDataApiServerV2.FetchBlobAttestationInfo)

	t.Run("no attestation found", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqStr := fmt.Sprintf("/v2/blobs/%s/attestation-info", blobKey.Hex())
		req := httptest.NewRequest(http.MethodGet, reqStr, nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	operatorPubKeys := []*core.G1Point{
		core.NewG1Point(big.NewInt(1), big.NewInt(2)),
		core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		core.NewG1Point(big.NewInt(4), big.NewInt(5)),
		core.NewG1Point(big.NewInt(5), big.NewInt(6)),
	}
	operatorAddresses := []gethcommon.Address{
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fb"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fc"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fd"),
	}
	operatorIDToAddr := make(map[string]gethcommon.Address)
	for i := 0; i < len(operatorPubKeys); i++ {
		operatorIDToAddr[operatorPubKeys[i].GetOperatorID().Hex()] = operatorAddresses[i]
	}
	mockTx.On("BatchOperatorIDToAddress").Return(
		func(ids []core.OperatorID) []gethcommon.Address {
			result := make([]gethcommon.Address, len(ids))
			for i, id := range ids {
				result[i] = operatorIDToAddr[id.Hex()]
			}
			return result
		},
		nil,
	)

	// Set up attestation
	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	apk := keyPair.GetPubKeyG2()
	nonsignerPubKeys := operatorPubKeys[:2]
	attestation := &corev2.Attestation{
		BatchHeader:      batchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: nonsignerPubKeys,
		APKG2:            apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
			1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}
	err = blobMetadataStore.PutAttestation(ctx, attestation)
	assert.NoError(t, err)

	operatorStakesByBlock := map[uint32]core.OperatorStakes{
		10: core.OperatorStakes{
			0: {
				0: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorPubKeys[1].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				2: {
					OperatorID: operatorPubKeys[2].GetOperatorID(),
					Stake:      big.NewInt(3),
				},
			},
			1: {
				0: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorPubKeys[2].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				2: {
					OperatorID: operatorPubKeys[3].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
			},
			2: {
				1: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
			},
		},
	}
	mockTx.On("GetOperatorStakesForQuorums").Return(
		func(quorums []core.QuorumID, blockNum uint32) core.OperatorStakes {
			return operatorStakesByBlock[blockNum]
		},
		nil,
	)

	t.Run("found attestation info", func(t *testing.T) {
		reqStr := fmt.Sprintf("/v2/blobs/%s/attestation-info", blobKey.Hex())
		w := executeRequest(t, r, http.MethodGet, reqStr)
		response := decodeResponseBody[serverv2.BlobAttestationInfoResponse](t, w)

		assert.Equal(t, blobKey.Hex(), response.BlobKey)
		assert.Equal(t, hex.EncodeToString(bhh[:]), response.BatchHeaderHash)
		assert.Equal(t, hex.EncodeToString(inclusionInfo.InclusionProof[:]), response.InclusionInfo.InclusionProof)
		assert.Equal(t, attestation, response.AttestationInfo.Attestation)

		signers := map[uint8][]serverv2.OperatorIdentity{
			0: []serverv2.OperatorIdentity{
				{
					OperatorId:      operatorPubKeys[2].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[2].Hex(),
				},
			},
			1: []serverv2.OperatorIdentity{
				{
					OperatorId:      operatorPubKeys[2].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[2].Hex(),
				},
				{
					OperatorId:      operatorPubKeys[3].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[3].Hex(),
				},
			},
		}
		nonsigners := map[uint8][]serverv2.OperatorIdentity{
			0: []serverv2.OperatorIdentity{
				{
					OperatorId:      operatorPubKeys[0].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[0].Hex(),
				},
				{
					OperatorId:      operatorPubKeys[1].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[1].Hex(),
				},
			},
			1: []serverv2.OperatorIdentity{
				{
					OperatorId:      operatorPubKeys[0].GetOperatorID().Hex(),
					OperatorAddress: operatorAddresses[0].Hex(),
				},
			},
		}
		for key, expectedSigners := range signers {
			actualSigners, exists := response.AttestationInfo.Signers[key]
			require.True(t, exists)
			assert.ElementsMatch(t, expectedSigners, actualSigners)
		}
		for key, expectedNonsigners := range nonsigners {
			actualNonsigners, exists := response.AttestationInfo.Nonsigners[key]
			require.True(t, exists)
			assert.ElementsMatch(t, expectedNonsigners, actualNonsigners)
		}
	})

	mockTx.ExpectedCalls = nil
	mockTx.Calls = nil
	deleteItems(t, []dynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func TestFetchBatch(t *testing.T) {
	r := setUpRouter()

	operatorPubKeys := []*core.G1Point{
		core.NewG1Point(big.NewInt(1), big.NewInt(2)),
		core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		core.NewG1Point(big.NewInt(4), big.NewInt(5)),
		core.NewG1Point(big.NewInt(5), big.NewInt(6)),
	}
	operatorAddresses := []gethcommon.Address{
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fb"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fc"),
		gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fd"),
	}
	operatorIDToAddr := make(map[string]gethcommon.Address)
	for i := 0; i < len(operatorPubKeys); i++ {
		operatorIDToAddr[operatorPubKeys[i].GetOperatorID().Hex()] = operatorAddresses[i]
	}

	// Set up batch header in metadata store
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 10,
	}
	err := blobMetadataStore.PutBatchHeader(context.Background(), batchHeader)
	require.NoError(t, err)
	batchHeaderHashBytes, err := batchHeader.Hash()
	require.NoError(t, err)
	batchHeaderHash := hex.EncodeToString(batchHeaderHashBytes[:])

	// Set up batch in metadata store
	blobHeader := makeBlobHeaderV2(t)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		Signature:  []byte{0, 1, 2, 3, 4},
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}
	batch := &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: []*corev2.BlobCertificate{blobCert},
	}
	err = blobMetadataStore.PutBatch(context.Background(), batch)
	require.NoError(t, err)

	// Set up attestation in metadata store
	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	apk := keyPair.GetPubKeyG2()
	nonsignerPubKeys := operatorPubKeys[:2]
	attestation := &corev2.Attestation{
		BatchHeader:      batchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: nonsignerPubKeys,
		APKG2:            apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
			1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}
	err = blobMetadataStore.PutAttestation(context.Background(), attestation)
	require.NoError(t, err)

	mockTx.On("BatchOperatorIDToAddress").Return(
		func(ids []core.OperatorID) []gethcommon.Address {
			result := make([]gethcommon.Address, len(ids))
			for i, id := range ids {
				result[i] = operatorIDToAddr[id.Hex()]
			}
			return result
		},
		nil,
	)

	operatorStakesByBlock := map[uint32]core.OperatorStakes{
		10: core.OperatorStakes{
			0: {
				0: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorPubKeys[1].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				2: {
					OperatorID: operatorPubKeys[2].GetOperatorID(),
					Stake:      big.NewInt(3),
				},
			},
			1: {
				0: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorPubKeys[2].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
				2: {
					OperatorID: operatorPubKeys[3].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
			},
			2: {
				1: {
					OperatorID: operatorPubKeys[0].GetOperatorID(),
					Stake:      big.NewInt(2),
				},
			},
		},
	}
	mockTx.On("GetOperatorStakesForQuorums").Return(
		func(quorums []core.QuorumID, blockNum uint32) core.OperatorStakes {
			return operatorStakesByBlock[blockNum]
		},
		nil,
	)

	r.GET("/v2/batches/:batch_header_hash", testDataApiServerV2.FetchBatch)

	w := executeRequest(t, r, http.MethodGet, "/v2/batches/"+batchHeaderHash)
	response := decodeResponseBody[serverv2.BatchResponse](t, w)

	assert.Equal(t, batchHeaderHash, response.BatchHeaderHash)
	assert.Equal(t, hex.EncodeToString(batchHeader.BatchRoot[:]), response.SignedBatch.BatchHeader.BatchRoot)
	assert.Equal(t, batchHeader.ReferenceBlockNumber, response.SignedBatch.BatchHeader.ReferenceBlockNumber)
	assert.Equal(t, attestation.AttestedAt, response.SignedBatch.AttestationInfo.Attestation.AttestedAt)
	assert.Equal(t, attestation.QuorumNumbers, response.SignedBatch.AttestationInfo.Attestation.QuorumNumbers)
	assert.Equal(t, 1, len(response.BlobKeys))
	assert.Equal(t, blobKey.Hex(), response.BlobKeys[0])
	assert.Equal(t, 1, len(response.BlobCertificates))
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, response.BlobCertificates[0].Signature)

	signers := map[uint8][]serverv2.OperatorIdentity{
		0: []serverv2.OperatorIdentity{
			{
				OperatorId:      operatorPubKeys[2].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[2].Hex(),
			},
		},
		1: []serverv2.OperatorIdentity{
			{
				OperatorId:      operatorPubKeys[2].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[2].Hex(),
			},
			{
				OperatorId:      operatorPubKeys[3].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[3].Hex(),
			},
		},
	}
	nonsigners := map[uint8][]serverv2.OperatorIdentity{
		0: []serverv2.OperatorIdentity{
			{
				OperatorId:      operatorPubKeys[0].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[0].Hex(),
			},
			{
				OperatorId:      operatorPubKeys[1].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[1].Hex(),
			},
		},
		1: []serverv2.OperatorIdentity{
			{
				OperatorId:      operatorPubKeys[0].GetOperatorID().Hex(),
				OperatorAddress: operatorAddresses[0].Hex(),
			},
		},
	}
	for key, expectedSigners := range signers {
		actualSigners, exists := response.SignedBatch.AttestationInfo.Signers[key]
		require.True(t, exists)
		assert.ElementsMatch(t, expectedSigners, actualSigners)
	}
	for key, expectedNonsigners := range nonsigners {
		actualNonsigners, exists := response.SignedBatch.AttestationInfo.Nonsigners[key]
		require.True(t, exists)
		assert.ElementsMatch(t, expectedNonsigners, actualNonsigners)
	}

	mockTx.ExpectedCalls = nil
	mockTx.Calls = nil
	deleteItems(t, []dynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + batchHeaderHash},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + batchHeaderHash},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + batchHeaderHash},
			"SK": &types.AttributeValueMemberS{Value: "BatchInfo"},
		},
	})
}

func TestFetchBatchFeed(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()

	// Create a timeline of test batches
	numBatches := 72
	now := uint64(time.Now().UnixNano())
	firstBatchTs := now - uint64(72*time.Minute.Nanoseconds())
	nanoSecsPerBatch := uint64(time.Minute.Nanoseconds()) // 1 batch per minute
	attestedAt := make([]uint64, numBatches)
	batchHeaders := make([]*corev2.BatchHeader, numBatches)
	dynamoKeys := make([]dynamodb.Key, numBatches)
	for i := 0; i < numBatches; i++ {
		batchHeaders[i] = &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, byte(i)},
			ReferenceBlockNumber: uint64(i + 1),
		}
		bhh, err := batchHeaders[i].Hash()
		require.NoError(t, err)
		keyPair, err := core.GenRandomBlsKeys()
		assert.NoError(t, err)
		apk := keyPair.GetPubKeyG2()
		attestedAt[i] = firstBatchTs + uint64(i)*nanoSecsPerBatch
		attestation := &corev2.Attestation{
			BatchHeader: batchHeaders[i],
			AttestedAt:  attestedAt[i],
			NonSignerPubKeys: []*core.G1Point{
				core.NewG1Point(big.NewInt(1), big.NewInt(2)),
				core.NewG1Point(big.NewInt(3), big.NewInt(4)),
			},
			APKG2: apk,
			QuorumAPKs: map[uint8]*core.G1Point{
				0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
				1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
			},
			Sigma: &core.Signature{
				G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
			},
			QuorumNumbers: []core.QuorumID{0, 1},
			QuorumResults: map[uint8]uint8{
				0: 100,
				1: 80,
			},
		}
		err = blobMetadataStore.PutAttestation(ctx, attestation)
		require.NoError(t, err)
		dynamoKeys[i] = dynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	mockTx.On("GetCurrentBlockNumber").Return(uint32(1), nil)
	mockTx.On("GetQuorumCount").Return(uint8(2), nil)

	// Create a local server so the internal state (e.g. cache) will be re-created.
	// This is needed because /v2/operators/signing-info API shares the cache state with
	// /v2/batches/feed API.
	testDataApiServerV2, err := serverv2.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, dataapi.NewMetrics(serverVersion, nil, nil, "9001", mockLogger))
	require.NoError(t, err)

	r.GET("/v2/batches/feed", testDataApiServerV2.FetchBatchFeed)

	t.Run("invalid params", func(t *testing.T) {
		now := time.Now()

		tests := []struct {
			name        string
			queryParams map[string]string
			wantError   string // expected error message
		}{
			// Invalid direction
			{
				name:        "invalid direction",
				queryParams: map[string]string{"direction": "abc"},
				wantError:   "`direction` must be either \"forward\" or \"backward\", found: \"abc\"",
			},

			// Invalid time formats
			{
				name:        "invalid before format",
				queryParams: map[string]string{"before": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid before value",
				queryParams: map[string]string{"before": "abc"},
				wantError:   "failed to parse `before` param",
			},
			{
				name:        "invalid after format",
				queryParams: map[string]string{"after": "2006-01-02T15:04:05"}, // missing Z
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "invalid after value",
				queryParams: map[string]string{"after": "abc"},
				wantError:   "failed to parse `after` param",
			},
			{
				name:        "after in future",
				queryParams: map[string]string{"after": "3025-01-02T15:04:05Z"},
				wantError:   "`after` must be before current time",
			},

			// Invalid time ranges
			{
				name: "after >= before",
				queryParams: map[string]string{
					"after":  now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.999999999Z"),
					"before": now.Add(-time.Hour).UTC().Format("2006-01-02T15:04:05.999999999Z"),
				},
				wantError: "must be earlier than `before` timestamp",
			},
			{
				name: "before too old",
				queryParams: map[string]string{
					"before": "2020-01-02T15:04:05Z",
				},
				wantError: "`before` time cannot be more than 14 days in the past",
			},

			// Invalid limit
			{
				name:        "invalid limit format",
				queryParams: map[string]string{"limit": "abc"},
				wantError:   "failed to parse `limit` param",
			},
		}

		for _, tt := range tests {
			params := url.Values{}
			for k, v := range tt.queryParams {
				params.Add(k, v)
			}
			url := fmt.Sprintf("/v2/batches/feed?%s", params.Encode())

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

			var errResp serverv2.ErrorResponse
			require.NoError(t, json.NewDecoder(w.Body).Decode(&errResp))
			assert.Contains(t, errResp.Error, tt.wantError)
		}

	})

	t.Run("default params", func(t *testing.T) {
		// Default query returns:
		// - Most recent 1 hour of batches attested (batch[13], ..., batch[71])
		// - Limited to 20 results (the default "limit")
		// - Result will first 20 batches: batch[13], ..., batch[42]
		w := executeRequest(t, r, http.MethodGet, "/v2/batches/feed")
		response := decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 20, len(response.Batches))
		for i := 0; i < 20; i++ {
			assert.Equal(t, attestedAt[13+i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[13+i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[13+i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}
	})

	t.Run("forward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves batch[13] through batch[71] (59 batches)
		w := executeRequest(t, r, http.MethodGet, "/v2/batches/feed?limit=0")
		response := decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 59, len(response.Batches))
		for i := 0; i < 59; i++ {
			assert.Equal(t, attestedAt[13+i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[13+i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[13+i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}

		// Test 2: 2-hour window captures all test batches
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/batches/feed?limit=-1&after=%s", afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 72, len(response.Batches))
		for i := 0; i < 72; i++ {
			assert.Equal(t, attestedAt[i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}

		// Test 3: Custom end time with 1-hour window
		// With 1h ending time at attestedAt[66], this retrieves batch[7] throught batch[65] (59 batches, as the `before` is exclusive)
		tm := time.Unix(0, int64(attestedAt[66])).UTC()
		beforeTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("/v2/batches/feed?before=%s&limit=-1", beforeTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 59, len(response.Batches))
		for i := 0; i < 59; i++ {
			assert.Equal(t, attestedAt[7+i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[7+i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[7+i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}
	})

	t.Run("backward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves batch[71] through batch[13] (59 batches)
		w := executeRequest(t, r, http.MethodGet, "/v2/batches/feed?direction=backward&limit=0")
		response := decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 59, len(response.Batches))
		for i := 0; i < 59; i++ {
			assert.Equal(t, attestedAt[71-i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[71-i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[71-i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}

		// Test 2: 2-hour window captures all test batches
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("/v2/batches/feed?direction=backward&limit=-1&after=%s", afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 72, len(response.Batches))
		for i := 0; i < 72; i++ {
			assert.Equal(t, attestedAt[71-i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[71-i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[71-i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}

		// Test 3: Custom end time with 1-hour window
		// With 1h ending time at attestedAt[66], this retrieves batch[65] throught batch[7] (59 batches,
		// as the `before` is exclusive)
		tm := time.Unix(0, int64(attestedAt[66])).UTC()
		beforeTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("/v2/batches/feed?direction=backward&before=%s&limit=-1", beforeTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.BatchFeedResponse](t, w)
		require.Equal(t, 59, len(response.Batches))
		for i := 0; i < 59; i++ {
			assert.Equal(t, attestedAt[65-i], response.Batches[i].AttestedAt)
			assert.Equal(t, batchHeaders[65-i].ReferenceBlockNumber, response.Batches[i].BatchHeader.ReferenceBlockNumber)
			assert.Equal(t, hex.EncodeToString(batchHeaders[65-i].BatchRoot[:]), response.Batches[i].BatchHeader.BatchRoot)
		}
	})

}

func TestFetchOperatorSigningInfo(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()

	/*
		Test data setup

		Column definitions:
		- Batch:            Batch number
		- AttestedAt:       Timestamp of attestation (sortkey of this table)
		- RefBlockNum:      Reference block number
		- Quorums:          Quorum numbers used by the batch
		- Nonsigners:       Operators that didn't sign for the batch
		- Active operators: Mapping of operator ID to their quorum assignments at the block

		Data:
		+-------+------------+-------------+---------+------------+------------------------+
		| Batch | AttestedAt | RefBlockNum | Quorums | Nonsigners | Active operators      |
		+-------+------------+-------------+---------+------------+------------------------+
		|     1 |          A |           1 | 0,1     | 3          | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		+-------+------------+-------------+---------+------------+------------------------+
		|     2 |          B |           3 | 1       | 4          | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		|       |            |             |         |            | 4: {0,1}              |
		|       |            |             |         |            | 5: {0}                |
		+-------+------------+-------------+---------+------------+------------------------+
		|     3 |          C |           2 | 0       | 3          | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		|       |            |             |         |            | 4: {0,1}              |
		+-------+------------+-------------+---------+------------+------------------------+
		|     4 |          D |           2 | 0,1     | None       | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		|       |            |             |         |            | 4: {0,1}              |
		+-------+------------+-------------+---------+------------+------------------------+
		|     5 |          E |           4 | 0,1     | 3,5        | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		|       |            |             |         |            | 5: {0}                |
		+-------+------------+-------------+---------+------------+------------------------+
		|     6 |          F |           5 | 0       | 5          | 1: {2}                |
		|       |            |             |         |            | 2: {0,1}              |
		|       |            |             |         |            | 3: {0,1}              |
		|       |            |             |         |            | 5: {0}                |
		+-------+------------+-------------+---------+------------+------------------------+
	*/

	// Create test operators
	// Note: the operator numbered 1-5 in the above tables are corresponding to the
	// operatorIds[0], ..., operatorIds[4] here
	numOperators := 5
	operatorIds := make([]core.OperatorID, numOperators)
	operatorAddresses := make([]gethcommon.Address, numOperators)
	operatorG1s := make([]*core.G1Point, numOperators)
	operatorIDToAddr := make(map[string]gethcommon.Address)
	operatorAddrToID := make(map[string]core.OperatorID)
	for i := 0; i < numOperators; i++ {
		operatorG1s[i] = core.NewG1Point(big.NewInt(int64(i)), big.NewInt(int64(i+1)))
		operatorIds[i] = operatorG1s[i].GetOperatorID()
		privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		require.NoError(t, err)
		publicKey := privateKey.Public().(*ecdsa.PublicKey)
		operatorAddresses[i] = crypto.PubkeyToAddress(*publicKey)

		operatorIDToAddr[operatorIds[i].Hex()] = operatorAddresses[i]
		operatorAddrToID[operatorAddresses[i].Hex()] = operatorIds[i]
	}

	// Mocking using a map function so we can always maintain the ID and address mapping
	// defined above, ie. operatorIds[i] <-> operatorAddresses[i]
	mockTx.On("BatchOperatorIDToAddress").Return(
		func(ids []core.OperatorID) []gethcommon.Address {
			result := make([]gethcommon.Address, len(ids))
			for i, id := range ids {
				result[i] = operatorIDToAddr[id.Hex()]
			}
			return result
		},
		nil,
	)
	mockTx.On("BatchOperatorAddressToID").Return(
		func(addrs []gethcommon.Address) []core.OperatorID {
			result := make([]core.OperatorID, len(addrs))
			for i, addr := range addrs {
				result[i] = operatorAddrToID[addr.Hex()]
			}
			return result
		},
		nil,
	)

	// Mocking using a map function so we can always maintain the ID and address mapping
	// defined above, ie. operatorIds[i] <-> operatorAddresses[i]
	// We prepare data at two blocks (1 and 4) as they will be hit by queries below
	operatorIntialQuorumsByBlock := map[uint32]map[core.OperatorID]*big.Int{
		1: map[core.OperatorID]*big.Int{
			operatorIds[0]: big.NewInt(4), // quorum 2
			operatorIds[1]: big.NewInt(3), // quorum 0,1
			operatorIds[2]: big.NewInt(3), // quorum 0,1
			operatorIds[3]: big.NewInt(0), // no quorum
			operatorIds[4]: big.NewInt(0), // no quorum
		},
		4: map[core.OperatorID]*big.Int{
			operatorIds[0]: big.NewInt(4), // quorum 2
			operatorIds[1]: big.NewInt(3), // quorum 0,1
			operatorIds[2]: big.NewInt(3), // quorum 0,1
			operatorIds[3]: big.NewInt(0), // no quorum
			operatorIds[4]: big.NewInt(1), // quorum 0
		},
	}
	mockTx.On("GetQuorumBitmapForOperatorsAtBlockNumber").Return(
		func(ids []core.OperatorID, blockNum uint32) []*big.Int {
			bitmaps := make([]*big.Int, len(ids))
			for i, id := range ids {
				bitmaps[i] = operatorIntialQuorumsByBlock[blockNum][id]
			}
			return bitmaps
		},
		nil,
	)

	// We prepare data at two blocks (1 and 4) as they will be hit by queries below
	operatorStakesByBlock := map[uint32]core.OperatorStakes{
		1: core.OperatorStakes{
			0: {
				0: {
					OperatorID: operatorIds[1],
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorIds[2],
					Stake:      big.NewInt(2),
				},
			},
			1: {
				0: {
					OperatorID: operatorIds[1],
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorIds[2],
					Stake:      big.NewInt(2),
				},
			},
			2: {
				1: {
					OperatorID: operatorIds[0],
					Stake:      big.NewInt(2),
				},
			},
		},
		4: core.OperatorStakes{
			0: {
				0: {
					OperatorID: operatorIds[1],
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorIds[2],
					Stake:      big.NewInt(2),
				},
				2: {
					OperatorID: operatorIds[4],
					Stake:      big.NewInt(2),
				},
			},
			1: {
				0: {
					OperatorID: operatorIds[1],
					Stake:      big.NewInt(2),
				},
				1: {
					OperatorID: operatorIds[2],
					Stake:      big.NewInt(2),
				},
			},
			2: {
				1: {
					OperatorID: operatorIds[0],
					Stake:      big.NewInt(2),
				},
			},
		},
	}
	mockTx.On("GetOperatorStakesForQuorums").Return(
		func(quorums []core.QuorumID, blockNum uint32) core.OperatorStakes {
			return operatorStakesByBlock[blockNum]
		},
		nil,
	)

	// operatorIds[3], operatorIds[4] were not active at the first block, but were added to
	// quorums after startBlock (see the above table).
	operatorAddedToQuorum := []*subgraph.OperatorQuorum{
		{
			Operator:       graphql.String(operatorAddresses[3].Hex()),
			QuorumNumbers:  "0x0001",
			BlockNumber:    "2",
			BlockTimestamp: "1702666070",
		},
		{
			Operator:       graphql.String(operatorAddresses[4].Hex()),
			QuorumNumbers:  "0x00",
			BlockNumber:    "3",
			BlockTimestamp: "1702666070",
		},
	}
	operatorRemovedFromQuorum := []*subgraph.OperatorQuorum{
		{
			Operator:       graphql.String(operatorAddresses[3].Hex()),
			QuorumNumbers:  "0x0001",
			BlockNumber:    "4",
			BlockTimestamp: "1702666058",
		},
	}
	mockSubgraphApi.On("QueryOperatorAddedToQuorum").Return(operatorAddedToQuorum, nil)
	mockSubgraphApi.On("QueryOperatorRemovedFromQuorum").Return(operatorRemovedFromQuorum, nil)

	// Create a timeline of test batches
	// See the above table for the choices of reference block number, quorums and nonsigners
	// for each batch
	numBatches := 6
	now := uint64(time.Now().UnixNano())
	firstBatchTime := now - uint64(32*time.Minute.Nanoseconds())
	nanoSecsPerBatch := uint64(5 * time.Minute.Nanoseconds()) // 1 batch per 5 minutes
	attestedAt := make([]uint64, numBatches)
	for i := 0; i < numBatches; i++ {
		attestedAt[i] = firstBatchTime + uint64(i)*nanoSecsPerBatch
	}
	referenceBlockNum := []uint64{1, 3, 2, 2, 4, 5}
	quorums := [][]uint8{{0, 1}, {1}, {0}, {0, 1}, {0, 1}, {0}}
	nonsigners := [][]*core.G1Point{
		{operatorG1s[2]},
		{operatorG1s[3]},
		{operatorG1s[2]},
		{},
		{operatorG1s[2], operatorG1s[4]},
		{operatorG1s[4]},
	}
	dynamoKeys := make([]dynamodb.Key, numBatches)
	for i := 0; i < numBatches; i++ {
		attestation := createAttestation(t, referenceBlockNum[i], attestedAt[i], nonsigners[i], quorums[i])
		err := blobMetadataStore.PutAttestation(ctx, attestation)
		require.NoError(t, err)
		bhh, err := attestation.BatchHeader.Hash() // go:nolint QF1008
		require.NoError(t, err)
		dynamoKeys[i] = dynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	/*
		Resulting Operator SigningInfo (for block range [1, 5])

		Column definitions:
		- <operator, quorum>:    Operator ID and quorum pair
		- Total responsible:     Total number of batches the operator was responsible for
		- Total nonsigning:      Number of batches where operator did not sign
		- Signing rate:          Percentage of batches signed by <operator, quorum>

		Data:
		+------------------+-------------------+------------------+--------------+
		| <operator,quorum>| Total responsible | Total nonsigning | Signing rate |
		+------------------+-------------------+------------------+--------------+
		| <2, 0>           |                 5 |                0 |        100%  |
		+------------------+-------------------+------------------+--------------+
		| <2, 1>           |                 4 |                0 |        100%  |
		+------------------+-------------------+------------------+--------------+
		| <3, 0>           |                 5 |                3 |         40%  |
		+------------------+-------------------+------------------+--------------+
		| <3, 1>           |                 4 |                2 |         50%  |
		+------------------+-------------------+------------------+--------------+
		| <4, 0>           |                 2 |                0 |        100%  |
		+------------------+-------------------+------------------+--------------+
		| <4, 1>           |                 2 |                1 |         50%  |
		+------------------+-------------------+------------------+--------------+
		| <5, 0>           |                 2 |                2 |          0%  |
		+------------------+-------------------+------------------+--------------+
	*/

	// Create a local server so the internal state (e.g. cache) will be re-created.
	// This is needed because /v2/operators/signing-info API shares the cache state with
	// /v2/batches/feed API.
	testDataApiServerV2, err := serverv2.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, dataapi.NewMetrics(serverVersion, nil, nil, "9001", mockLogger))
	require.NoError(t, err)

	r.GET("/v2/operators/signing-info", testDataApiServerV2.FetchOperatorSigningInfo)

	t.Run("invalid params", func(t *testing.T) {
		reqUrls := []string{
			"/v2/operators/signing-info?interval=abc",
			"/v2/operators/signing-info?interval=-1",
			"/v2/operators/signing-info?end=2006-01-02T15:04:05",
			"/v2/operators/signing-info?end=2006-01-02T15:04:05Z",
			"/v2/operators/signing-info?quorums=-1",
			"/v2/operators/signing-info?quorums=abc",
			"/v2/operators/signing-info?quorums=10000000",
			"/v2/operators/signing-info?quorums=-1",
			"/v2/operators/signing-info?nonsigner_only=-1",
			"/v2/operators/signing-info?nonsigner_only=deadbeef",
		}
		for _, url := range reqUrls {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)
			require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		}
	})

	t.Run("default params", func(t *testing.T) {
		w := executeRequest(t, r, http.MethodGet, "/v2/operators/signing-info")
		response := decodeResponseBody[serverv2.OperatorsSigningInfoResponse](t, w)
		osi := response.OperatorSigningInfo
		require.Equal(t, 7, len(osi))
		checkOperatorSigningInfoEqual(t, osi[0], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[3].Hex(),
			OperatorAddress:         operatorAddresses[3].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 2,
			TotalBatches:            5,
		})
		checkOperatorSigningInfoEqual(t, osi[1], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[1].Hex(),
			OperatorAddress:         operatorAddresses[1].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 5,
			TotalBatches:            5,
		})
		checkOperatorSigningInfoEqual(t, osi[2], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[1].Hex(),
			OperatorAddress:         operatorAddresses[1].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 4,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[3], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[3].Hex(),
			OperatorAddress:         operatorAddresses[3].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    1,
			TotalResponsibleBatches: 2,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[4], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 4,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[5], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    3,
			TotalResponsibleBatches: 5,
			TotalBatches:            5,
		})
		checkOperatorSigningInfoEqual(t, osi[6], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[4].Hex(),
			OperatorAddress:         operatorAddresses[4].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 2,
			TotalBatches:            5,
		})
	})

	t.Run("nonsigner only", func(t *testing.T) {
		w := executeRequest(t, r, http.MethodGet, "/v2/operators/signing-info?nonsigner_only=true")
		response := decodeResponseBody[serverv2.OperatorsSigningInfoResponse](t, w)
		osi := response.OperatorSigningInfo
		require.Equal(t, 4, len(osi))
		checkOperatorSigningInfoEqual(t, osi[0], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[3].Hex(),
			OperatorAddress:         operatorAddresses[3].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    1,
			TotalResponsibleBatches: 2,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[1], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 4,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[2], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    3,
			TotalResponsibleBatches: 5,
			TotalBatches:            5,
		})
		checkOperatorSigningInfoEqual(t, osi[3], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[4].Hex(),
			OperatorAddress:         operatorAddresses[4].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 2,
			TotalBatches:            5,
		})
	})

	t.Run("quorum 1 only", func(t *testing.T) {
		w := executeRequest(t, r, http.MethodGet, "/v2/operators/signing-info?quorums=1")
		response := decodeResponseBody[serverv2.OperatorsSigningInfoResponse](t, w)
		osi := response.OperatorSigningInfo
		require.Equal(t, 3, len(osi))
		checkOperatorSigningInfoEqual(t, osi[0], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[1].Hex(),
			OperatorAddress:         operatorAddresses[1].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 4,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[1], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[3].Hex(),
			OperatorAddress:         operatorAddresses[3].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    1,
			TotalResponsibleBatches: 2,
			TotalBatches:            4,
		})
		checkOperatorSigningInfoEqual(t, osi[2], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 4,
			TotalBatches:            4,
		})
	})

	t.Run("custom time range", func(t *testing.T) {
		// We query 800 seconds before "now", it should hit the last 2 batches (block 4, 5)
		// in the setup table:
		//
		// +-------+------------+-------------+---------+------------+------------------------+
		// | Batch | AttestedAt | RefBlockNum | Quorums | Nonsigners | Active operators      |
		// +-------+------------+-------------+---------+------------+------------------------+
		// |     5 |          5 |           4 | 0,1     | 3,5        | 1: {2}                |
		// |       |            |             |         |            | 2: {0,1}              |
		// |       |            |             |         |            | 3: {0,1}              |
		// |       |            |             |         |            | 5: {0}                |
		// +-------+------------+-------------+---------+------------+------------------------+
		// |     6 |          6 |           5 | 0       | 5          | 1: {2}                |
		// |       |            |             |         |            | 2: {0,1}              |
		// |       |            |             |         |            | 3: {0,1}              |
		// |       |            |             |         |            | 5: {0}                |
		// +-------+------------+-------------+---------+------------+------------------------+
		//
		// which results in:
		//
		// +------------------+-------------------+------------------+--------------+
		// | <operator,quorum>| Total responsible | Total nonsigning | Signing rate |
		// +------------------+-------------------+------------------+--------------+
		// | <2, 0>           |                 2 |                0 |        100%  |
		// +------------------+-------------------+------------------+--------------+
		// | <2, 1>           |                 1 |                0 |        100%  |
		// +------------------+-------------------+------------------+--------------+
		// | <3, 0>           |                 2 |                1 |         50%  |
		// +------------------+-------------------+------------------+--------------+
		// | <3, 1>           |                 1 |                1 |         0%   |
		// +------------------+-------------------+------------------+--------------+
		// | <5, 0>           |                 2 |                2 |         0%   |
		// +------------------+-------------------+------------------+--------------+

		tm := time.Unix(0, int64(now)+1).UTC()
		endTime := tm.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl := fmt.Sprintf("/v2/operators/signing-info?end=%s&interval=1000", endTime)
		w := executeRequest(t, r, http.MethodGet, reqUrl)
		response := decodeResponseBody[serverv2.OperatorsSigningInfoResponse](t, w)
		osi := response.OperatorSigningInfo
		require.Equal(t, 5, len(osi))
		checkOperatorSigningInfoEqual(t, osi[0], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[1].Hex(),
			OperatorAddress:         operatorAddresses[1].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 2,
			TotalBatches:            2,
		})
		checkOperatorSigningInfoEqual(t, osi[1], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[1].Hex(),
			OperatorAddress:         operatorAddresses[1].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    0,
			TotalResponsibleBatches: 1,
			TotalBatches:            1,
		})
		checkOperatorSigningInfoEqual(t, osi[2], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    1,
			TotalResponsibleBatches: 2,
			TotalBatches:            2,
		})
		checkOperatorSigningInfoEqual(t, osi[3], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[4].Hex(),
			OperatorAddress:         operatorAddresses[4].Hex(),
			QuorumId:                0,
			TotalUnsignedBatches:    2,
			TotalResponsibleBatches: 2,
			TotalBatches:            2,
		})
		checkOperatorSigningInfoEqual(t, osi[4], &serverv2.OperatorSigningInfo{
			OperatorId:              operatorIds[2].Hex(),
			OperatorAddress:         operatorAddresses[2].Hex(),
			QuorumId:                1,
			TotalUnsignedBatches:    1,
			TotalResponsibleBatches: 1,
			TotalBatches:            1,
		})
	})

	mockTx.ExpectedCalls = nil
	mockTx.Calls = nil
}

func TestCheckOperatorsLiveness(t *testing.T) {
	r := setUpRouter()

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil

	mockIndexedChainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	r.GET("/v2/operators/liveness", testDataApiServerV2.CheckOperatorsLiveness)

	operatorId := core.OperatorID{1}.Hex()
	reqStr := fmt.Sprintf("/v2/operators/liveness?operator_id=%v", operatorId)
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[serverv2.OperatorLivenessResponse](t, w)

	assert.Equal(t, 1, len(response.Operators))
	assert.Equal(t, "0.0.0.0:3004", response.Operators[0].DispersalSocket)
	assert.Equal(t, false, response.Operators[0].DispersalOnline)
	assert.Equal(t, "", response.Operators[0].DispersalStatus)
	assert.Equal(t, "0.0.0.0:3005", response.Operators[0].RetrievalSocket)
	assert.Equal(t, false, response.Operators[0].RetrievalOnline)
	assert.Equal(t, "", response.Operators[0].RetrievalStatus)

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestCheckOperatorsLivenessLegacyV1SocketRegistration(t *testing.T) {
	r := setUpRouter()

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil

	operatorId := core.OperatorID{1}
	ios := &core.IndexedOperatorState{
		IndexedOperators: map[core.OperatorID]*core.IndexedOperatorInfo{
			operatorId: &core.IndexedOperatorInfo{
				Socket: "1.2.3.4:3004:3005",
			},
		},
	}

	mockIcs := &coremock.MockIndexedChainState{}

	mockIcs.On("GetCurrentBlockNumber").Return(uint(1), nil)
	mockIcs.On("GetIndexedOperatorState").Return(ios, nil)

	mockTx.On("GetCurrentBlockNumber").Return(uint32(1), nil)
	mockTx.On("GetQuorumCount").Return(uint8(2), nil)

	testDataApiServerV2, err := serverv2.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIcs, mockLogger, dataapi.NewMetrics(serverVersion, nil, nil, "9001", mockLogger))
	require.NoError(t, err)

	r.GET("/v2/operators/liveness", testDataApiServerV2.CheckOperatorsLiveness)

	reqStr := fmt.Sprintf("/v2/operators/liveness?operator_id=%v", operatorId.Hex())
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[serverv2.OperatorLivenessResponse](t, w)

	assert.Equal(t, 1, len(response.Operators))
	assert.Equal(t, "", response.Operators[0].DispersalSocket)
	assert.Equal(t, false, response.Operators[0].DispersalOnline)
	assert.Equal(t, "v2 dispersal port is not registered", response.Operators[0].DispersalStatus)
	assert.Equal(t, "", response.Operators[0].RetrievalSocket)
	assert.Equal(t, false, response.Operators[0].RetrievalOnline)
	assert.Equal(t, "v2 retrieval port is not registered", response.Operators[0].RetrievalStatus)

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchAccountBlobFeed(t *testing.T) {
	r := setUpRouter()
	ctx := context.Background()

	numBlobs := 60
	now := uint64(time.Now().UnixNano())
	firstBlobTime := now - uint64(int64(numBlobs)*time.Minute.Nanoseconds())
	nanoSecsPerBlob := uint64(time.Minute.Nanoseconds()) // 1 blob/min

	accountId := gethcommon.HexToAddress(fmt.Sprintf("0x000000000000000000000000000000000000000%d", 5))

	// Create blobs for testing
	requestedAt := make([]uint64, numBlobs)
	dynamoKeys := make([]dynamodb.Key, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobHeader := makeBlobHeaderV2(t)
		blobHeader.PaymentMetadata.AccountID = accountId
		blobKey, err := blobHeader.BlobKey()
		require.NoError(t, err)
		requestedAt[i] = firstBlobTime + nanoSecsPerBlob*uint64(i)
		now := time.Now()
		metadata := &commonv2.BlobMetadata{
			BlobHeader:  blobHeader,
			Signature:   []byte{1, 2, 3},
			BlobStatus:  commonv2.Encoded,
			Expiry:      uint64(now.Add(time.Hour).Unix()),
			NumRetries:  0,
			UpdatedAt:   uint64(now.UnixNano()),
			RequestedAt: requestedAt[i],
		}
		err = blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		dynamoKeys[i] = dynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
	}
	defer deleteItems(t, dynamoKeys)

	r.GET("/v2/accounts/:account_id/blobs", testDataApiServerV2.FetchAccountBlobFeed)
	baseUrl := fmt.Sprintf("/v2/accounts/%s/blobs", accountId.Hex())

	t.Run("invalid account ID params", func(t *testing.T) {
		tests := []struct {
			name           string
			accountID      string
			expectedStatus int
			expectedError  string
		}{
			// Invalid format cases
			{
				accountID:      "not-a-hex-string",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			{
				accountID:      "0x",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			{
				accountID:      "0xG1234567890123456789012345678901234567",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			{
				accountID:      "0x123",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			{
				accountID:      "0x" + "1234567890123456789012345678901234567890abcdef",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			// Zero address case
			{
				accountID:      "0x0000000000000000000000000000000000000000",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "zero account id is not valid",
			},
			// Empty & whitespace cases
			{
				accountID:      "",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
			{
				accountID:      " ",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "account id is not valid hex",
			},
		}

		for _, tc := range tests {
			url := "/v2/accounts/" + tc.accountID + "/blobs"

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response serverv2.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response.Error, tc.expectedError)
		}
	})

	t.Run("nonexistent account", func(t *testing.T) {
		otherID := gethcommon.HexToAddress(fmt.Sprintf("0x000000000000000000000000000000000000000%d", 6))
		url := fmt.Sprintf("/v2/accounts/%s/blobs", otherID.Hex())
		w := executeRequest(t, r, http.MethodGet, url)
		response := decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		require.Equal(t, 0, len(response.Blobs))
	})

	t.Run("default params", func(t *testing.T) {
		// Default query returns:
		// - Most recent 1 hour of blobs include all of blobs[1] through blobs[59]
		// - Limited to 20 results (the default "limit")
		// - Result will first 20 blobs
		w := executeRequest(t, r, http.MethodGet, baseUrl)
		response := decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		assert.Equal(t, accountId.Hex(), response.AccountId)
		require.Equal(t, 20, len(response.Blobs))
		for i := 0; i < 20; i++ {
			assert.Equal(t, requestedAt[1+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
	})

	t.Run("forward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves blobs[1] through blobs[59] (59 blobs)
		w := executeRequest(t, r, http.MethodGet, baseUrl+"?limit=0")
		response := decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		assert.Equal(t, accountId.Hex(), response.AccountId)
		require.Equal(t, 59, len(response.Blobs))
		for i := 0; i < 59; i++ {
			assert.Equal(t, requestedAt[1+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}

		// Test 2: 2-hour window captures all test blobs
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("%s?limit=-1&after=%s", baseUrl, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		assert.Equal(t, accountId.Hex(), response.AccountId)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			assert.Equal(t, requestedAt[i], response.Blobs[i].BlobMetadata.RequestedAt)
		}

		// Teste 3: custom end time
		after := time.Unix(0, int64(requestedAt[20])).UTC()
		afterTime = after.Format("2006-01-02T15:04:05.999999999Z")
		before := time.Unix(0, int64(requestedAt[50])).UTC()
		beforeTime := before.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("%s?before=%s&after=%s&limit=-1", baseUrl, beforeTime, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		require.Equal(t, 29, len(response.Blobs))
		for i := 0; i < 29; i++ {
			assert.Equal(t, requestedAt[21+i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
	})

	t.Run("backward iteration with various query ranges and limits", func(t *testing.T) {
		// Test 1: Unlimited results in 1-hour window
		// With 1h ending time at now, this retrieves blobs[59] through blobs[1] (59 blobs)
		w := executeRequest(t, r, http.MethodGet, baseUrl+"?limit=0&direction=backward")
		response := decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		assert.Equal(t, accountId.Hex(), response.AccountId)
		require.Equal(t, 59, len(response.Blobs))
		for i := 0; i < 59; i++ {
			assert.Equal(t, requestedAt[59-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}

		// Test 2: 2-hour window captures all test blobs
		afterTime := time.Now().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.999999999Z") // nano precision format
		reqUrl := fmt.Sprintf("%s?limit=-1&after=%s&direction=backward", baseUrl, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		assert.Equal(t, accountId.Hex(), response.AccountId)
		require.Equal(t, 60, len(response.Blobs))
		for i := 0; i < 60; i++ {
			assert.Equal(t, requestedAt[59-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}

		// Teste 3: custom end time
		after := time.Unix(0, int64(requestedAt[20])).UTC()
		afterTime = after.Format("2006-01-02T15:04:05.999999999Z")
		before := time.Unix(0, int64(requestedAt[50])).UTC()
		beforeTime := before.Format("2006-01-02T15:04:05.999999999Z")
		reqUrl = fmt.Sprintf("%s?before=%s&after=%s&limit=-1&direction=backward", baseUrl, beforeTime, afterTime)
		w = executeRequest(t, r, http.MethodGet, reqUrl)
		response = decodeResponseBody[serverv2.AccountBlobFeedResponse](t, w)
		require.Equal(t, 29, len(response.Blobs))
		for i := 0; i < 29; i++ {
			assert.Equal(t, requestedAt[49-i], response.Blobs[i].BlobMetadata.RequestedAt)
		}
	})
}

func TestFetchOperatorDispersalResponse(t *testing.T) {
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

	r.GET("/v2/operators/:operator_id/dispersals/:batch_header_hash/response", testDataApiServerV2.FetchOperatorDispersalResponse)

	// Fetch response of a specific operator
	reqStr := fmt.Sprintf("/v2/operators/%s/dispersals/%s/response", operatorId.Hex(), batchHeaderHash)
	w := executeRequest(t, r, http.MethodGet, reqStr)
	response := decodeResponseBody[serverv2.OperatorDispersalResponse](t, w)
	require.Equal(t, dispersalResponse, response.Response)
}

func TestFetchOperatorsStake(t *testing.T) {
	r := setUpRouter()

	mockIndexedChainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	addr0 := gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa")
	addr1 := gethcommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	mockTx.On("BatchOperatorIDToAddress").Return(
		func(ids []core.OperatorID) []gethcommon.Address {
			result := make([]gethcommon.Address, len(ids))
			for i, id := range ids {
				switch id {
				case opId0:
					result[i] = addr0
				case opId1:
					result[i] = addr1
				default:
					result[i] = gethcommon.Address{}
				}
			}
			return result
		},
		nil,
	)

	r.GET("/v2/operators/stake", testDataApiServerV2.FetchOperatorsStake)

	w := executeRequest(t, r, http.MethodGet, "/v2/operators/stake")
	response := decodeResponseBody[dataapi.OperatorsStakeResponse](t, w)

	// The quorums and the operators in the quorum are defined in "mockChainState"
	// There are 2 quorums (0, 1)
	require.Equal(t, 2, len(response.StakeRankedOperators))
	checkAddress := func(op *dataapi.OperatorStake) {
		if op.OperatorId == opId0.Hex() {
			assert.Equal(t, addr0.Hex(), op.OperatorAddress)
		}
		if op.OperatorId == opId1.Hex() {
			assert.Equal(t, addr1.Hex(), op.OperatorAddress)
		}
	}
	// Quorum 0
	ops, ok := response.StakeRankedOperators["0"]
	require.True(t, ok)
	require.Equal(t, 2, len(ops))
	assert.Equal(t, opId0.Hex(), ops[0].OperatorId)
	assert.Equal(t, opId1.Hex(), ops[1].OperatorId)
	checkAddress(ops[0])
	checkAddress(ops[1])
	// Quorum 1
	ops, ok = response.StakeRankedOperators["1"]
	require.True(t, ok)
	require.Equal(t, 2, len(ops))
	assert.Equal(t, opId1.Hex(), ops[0].OperatorId)
	assert.Equal(t, opId0.Hex(), ops[1].OperatorId)
	checkAddress(ops[0])
	checkAddress(ops[1])
}

func TestFetchMetricsSummary(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusRespAvgThroughput))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v2/metrics/summary", testDataApiServerV2.FetchMetricsSummary)

	w := executeRequest(t, r, http.MethodGet, "/v2/metrics/summary")
	response := decodeResponseBody[serverv2.MetricSummary](t, w)

	assert.Equal(t, 10422.560745809731, response.AverageBytesPerSecond)
}

func TestFetchMetricsThroughputTimeseries(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusRespAvgThroughput))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v2/metrics/timeseries/throughput", testDataApiServerV2.FetchMetricsThroughputTimeseries)

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

func TestFetchMetricsNetworkSigningRateTimeseries(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusResponseNetworkSigningRate))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil)

	r.GET("/v2/metrics/timeseries/network-signing-rate", testDataApiServerV2.FetchNetworkSigningRate)

	w := executeRequest(t, r, http.MethodGet, "/v2/metrics/timeseries/network-signing-rate")
	response := decodeResponseBody[serverv2.NetworkSigningRateResponse](t, w)

	require.Equal(t, 2, len(response.QuorumSigningRates))
	assert.Equal(t, "0", response.QuorumSigningRates[0].QuorumId)
	require.Equal(t, 12, len(response.QuorumSigningRates[0].DataPoints))
	assert.Equal(t, float64(98.1), response.QuorumSigningRates[0].DataPoints[0].SigningRate)
	assert.Equal(t, "1", response.QuorumSigningRates[1].QuorumId)
	assert.Equal(t, 12, len(response.QuorumSigningRates[1].DataPoints))
	assert.Equal(t, float64(98.1), response.QuorumSigningRates[1].DataPoints[0].SigningRate)
}

func createAttestation(
	t *testing.T,
	refBlockNumber uint64,
	attestedAt uint64,
	nonsigners []*core.G1Point,
	quorums []uint8,
) *corev2.Attestation {
	br := make([]byte, 32)
	_, err := rand.Read(br)
	require.NoError(t, err)
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            ([32]byte)(br),
		ReferenceBlockNumber: refBlockNumber,
	}
	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	apk := keyPair.GetPubKeyG2()
	return &corev2.Attestation{
		BatchHeader:      batchHeader,
		AttestedAt:       attestedAt,
		NonSignerPubKeys: nonsigners,
		APKG2:            apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
			1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
		},
		QuorumNumbers: quorums,
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}
}
