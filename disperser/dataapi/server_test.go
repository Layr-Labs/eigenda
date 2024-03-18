package dataapi_test

import (
	"context"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	prommock "github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus/mock"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	subgraphmock "github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	//go:embed testdata/prometheus-response-sample.json
	mockPrometheusResponse string
	//go:embed testdata/prometheus-resp-avg-throughput.json
	mockPrometheusRespAvgThroughput string

	expectedBlobCommitment *encoding.BlobCommitments
	mockLogger             = logging.NewNoopLogger()
	blobstore              = inmem.NewBlobStore()
	mockPrometheusApi      = &prommock.MockPrometheusApi{}
	prometheusClient       = dataapi.NewPrometheusClient(mockPrometheusApi, "test-cluster")
	mockSubgraphApi        = &subgraphmock.MockSubgraphApi{}
	subgraphClient         = dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger)

	config = dataapi.Config{ServerMode: "test", SocketAddr: ":8080", AllowOrigins: []string{"*"}, DisperserHostname: "localhost:32007", ChurnerHostname: "localhost:32009"}

	mockTx                          = &coremock.MockTransactor{}
	mockChainState, _               = coremock.MakeChainDataMock(core.OperatorIndex(1))
	testDataApiServer               = dataapi.NewServer(config, blobstore, prometheusClient, subgraphClient, mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)
	expectedBatchHeaderHash         = [32]byte{1, 2, 3}
	expectedBlobIndex               = uint32(1)
	expectedRequestedAt             = uint64(5567830000000000000)
	expectedDataLength              = 32
	expectedBatchId                 = uint32(99)
	expectedBatchRoot               = []byte("hello")
	expectedReferenceBlockNumber    = uint32(132)
	expectedConfirmationBlockNumber = uint32(150)
	expectedSignatoryRecordHash     = [32]byte{0}
	expectedFee                     = []byte{0}
	expectedInclusionProof          = []byte{1, 2, 3, 4, 5}
	gettysburgAddressBytes          = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

type MockSubgraphClient struct {
	mock.Mock
}

type MockGRPCConnection struct{}

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

func NewMockHealthCheckService() *MockHealthCheckService {
	return &MockHealthCheckService{
		ResponseMap: make(map[string]*grpc_health_v1.HealthCheckResponse),
	}
}

func (m *MockHealthCheckService) CheckHealth(ctx context.Context, serviceName string) (*grpc_health_v1.HealthCheckResponse, error) {
	response, exists := m.ResponseMap[serviceName]
	if !exists {
		// Simulate an unsupported service error or return a default response
		return nil, fmt.Errorf("unsupported service: %s", serviceName)
	}
	return response, nil
}

func (m *MockHealthCheckService) CloseConnections() error {
	// Close any open connections or resources
	return nil
}

func (m *MockHealthCheckService) AddResponse(serviceName string, response *grpc_health_v1.HealthCheckResponse) {
	m.ResponseMap[serviceName] = response
}

func TestFetchBlobHandler(t *testing.T) {
	r := setUpRouter()

	blob := makeTestBlob(0, 80)
	key := queueBlob(t, &blob, blobstore)
	markBlobConfirmed(t, &blob, key, expectedBatchHeaderHash, blobstore)
	blobKey := key.String()
	r.GET("/v1/feed/blobs/:blob_key", testDataApiServer.FetchBlobHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/feed/blobs/"+blobKey, nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.BlobMetadataResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, hex.EncodeToString(expectedBatchHeaderHash[:]), response.BatchHeaderHash)
	assert.Equal(t, expectedBlobIndex, uint32(response.BlobIndex))
	assert.Equal(t, hex.EncodeToString(expectedSignatoryRecordHash[:]), response.SignatoryRecordHash)
	assert.Equal(t, expectedReferenceBlockNumber, uint32(response.ReferenceBlockNumber))
	assert.Equal(t, hex.EncodeToString(expectedBatchRoot), response.BatchRoot)
	assert.Equal(t, hex.EncodeToString(expectedInclusionProof), response.BlobInclusionProof)
	assert.Equal(t, expectedBlobCommitment, response.BlobCommitment)
	assert.Equal(t, expectedBatchId, uint32(response.BatchId))
	assert.Equal(t, expectedConfirmationBlockNumber, uint32(response.ConfirmationBlockNumber))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000123", response.ConfirmationTxnHash)
	assert.Equal(t, hex.EncodeToString(expectedFee), response.Fee)
	assert.Equal(t, blob.RequestHeader.SecurityParams, response.SecurityParams)
	assert.Equal(t, uint64(5567830000), response.RequestAt)
}

func TestFetchBlobsHandler(t *testing.T) {
	defer goleak.VerifyNone(t)

	r := setUpRouter()
	blob := makeTestBlob(0, 10)

	for _, batch := range subgraphBatches {
		var (
			key = queueBlob(t, &blob, blobstore)
		)
		// Convert the string to a byte slice
		batchHeaderHashBytes := []byte(batch.BatchHeaderHash)
		batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes(batchHeaderHashBytes)
		assert.NoError(t, err)
		markBlobConfirmed(t, &blob, key, batchHeaderHash, blobstore)
	}

	mockSubgraphApi.On("QueryBatches").Return(subgraphBatches, nil)

	r.GET("/v1/feed/blobs", testDataApiServer.FetchBlobsHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/feed/blobs?limit=2", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.BlobsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))
}

func TestFetchMetricsHandler(t *testing.T) {
	defer goleak.VerifyNone(t)

	r := setUpRouter()

	blob := makeTestBlob(0, 10)
	for _, batch := range subgraphBatches {
		var (
			key = queueBlob(t, &blob, blobstore)
		)

		batchHeaderHashBytes := []byte(batch.BatchHeaderHash)
		batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes(batchHeaderHashBytes)
		assert.NoError(t, err)

		markBlobConfirmed(t, &blob, key, batchHeaderHash, blobstore)
	}

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusResponse))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockTx.On("GetCurrentBlockNumber").Return(uint32(1), nil)
	mockTx.On("GetQuorumCount").Return(uint8(2), nil)
	mockSubgraphApi.On("QueryBatches").Return(subgraphBatches, nil)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v1/metrics", testDataApiServer.FetchMetricsHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/metrics", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.Metric
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 16555.555555555555, response.Throughput)
	assert.Equal(t, float64(85.14485344239945), response.CostInGas)
	assert.Equal(t, uint64(1), response.TotalStake)
	assert.Len(t, response.TotalStakePerQuorum, 2)
	assert.Equal(t, uint64(1), response.TotalStakePerQuorum[0])
	assert.Equal(t, uint64(1), response.TotalStakePerQuorum[1])
}

func TestFetchMetricsTroughputHandler(t *testing.T) {
	r := setUpRouter()

	s := new(model.SampleStream)
	err := s.UnmarshalJSON([]byte(mockPrometheusRespAvgThroughput))
	assert.NoError(t, err)

	matrix := make(model.Matrix, 0)
	matrix = append(matrix, s)
	mockPrometheusApi.On("QueryRange").Return(matrix, nil, nil).Once()

	r.GET("/v1/metrics/throughput", testDataApiServer.FetchMetricsTroughputHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/metrics/throughput", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response []*dataapi.Throughput
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	var totalThroughput float64
	for _, v := range response {
		totalThroughput += v.Throughput
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 3481, len(response))
	assert.Equal(t, float64(11666.666666666666), response[0].Throughput)
	assert.Equal(t, uint64(1701292800), response[0].Timestamp)
	assert.Equal(t, float64(3.599722666666646e+07), totalThroughput)
}

func TestFetchUnsignedBatchesHandler(t *testing.T) {
	r := setUpRouter()

	mockSubgraphApi.On("QueryBatchNonSigningInfo").Return(batchNonSigningInfo, nil)
	addr1 := gethcommon.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa")
	addr2 := gethcommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	mockTx.On("BatchOperatorIDToAddress").Return([]gethcommon.Address{addr1, addr2}, nil)
	mockTx.On("GetQuorumBitmapForOperatorsAtBlockNumber").Return([]*big.Int{big.NewInt(3), big.NewInt(0)}, nil)
	mockSubgraphApi.On("QueryOperatorAddedToQuorum").Return(operatorAddedToQuorum, nil)
	mockSubgraphApi.On("QueryOperatorRemovedFromQuorum").Return(operatorRemovedFromQuorum, nil)

	r.GET("/v1/metrics/operator-nonsigning-percentage", testDataApiServer.FetchOperatorsNonsigningPercentageHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/metrics/operator-nonsigning-percentage", nil)
	ctxWithDeadline, cancel := context.WithTimeout(req.Context(), 500*time.Microsecond)
	defer cancel()

	req = req.WithContext(ctxWithDeadline)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.OperatorsNonsigningPercentage
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	responseData := response.Data[0]
	operatorId := responseData.OperatorId
	assert.Equal(t, 1, responseData.TotalBatches)
	assert.Equal(t, 1, responseData.TotalUnsignedBatches)
	assert.Equal(t, uint8(0), responseData.QuorumId)
	assert.Equal(t, float64(100), responseData.Percentage)
	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operatorId)

	responseData = response.Data[1]
	operatorId = responseData.OperatorId
	assert.Equal(t, 2, responseData.TotalBatches)
	assert.Equal(t, 2, responseData.TotalUnsignedBatches)
	assert.Equal(t, uint8(1), responseData.QuorumId)
	assert.Equal(t, float64(100), responseData.Percentage)
	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operatorId)
}

func TestFetchDeregisteredOperatorNoSocketInfoOneOperatorHandler(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfoNoSocketInfo

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfoNoSocketInfo, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", response.Data[0].OperatorId)
	assert.Equal(t, "failed to convert operator info gql to indexed operator info at blocknumber: 22 for operator 0x3078653232646165313261303037346632306238666339366130343839333736", response.Data[0].OperatorProcessError)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredMultipleOperatorsOneWithNoSocketInfoHandler(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfoNoSocketInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfoNoSocketInfo, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	// Start test server for Operator
	closeServer, err := startTestGRPCServer("localhost:32009") // Let the OS assign a free port
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer() // Ensure the server is closed after the test

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]

	responseJson := string(data)
	fmt.Printf("Response: %v\n", responseJson)

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)
	assert.Equal(t, "failed to convert operator info gql to indexed operator info at blocknumber: 22 for operator 0x3078653232646165313261303037346632306238666339366130343839333736", operator1Data.OperatorProcessError)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, true, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorInfoInvalidTimeStampHandler(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfoInvalidTimeStamp

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregisteredInvalidTimeStamp, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 0, response.Meta.Size)
	assert.Equal(t, 0, len(response.Data))

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorInfoInvalidTimeStampTwoOperatorsHandler(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfoInvalidTimeStamp
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregisteredInvalidTimeStampTwoOperator, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	operator1Data := response.Data[0]

	responseJson := string(data)
	fmt.Printf("Response: %v\n", responseJson)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator1Data.OperatorId)
	assert.Equal(t, uint(24), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchMetricsDeregisteredOperatorHandler(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	// Start the test server for Operator 2
	closeServer, err := startTestGRPCServer("localhost:32009")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer()

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]

	responseJson := string(data)
	fmt.Printf("Response: %v\n", responseJson)

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, true, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorOffline(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorState := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorState[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil)
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	for _, data := range response.Data {
		fmt.Printf("Data: %v\n", data)
		assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", data.OperatorId)
		assert.Equal(t, uint(22), data.BlockNumber)
		assert.Equal(t, "localhost:32007", data.Socket)
	}

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorsWithoutDaysQueryParam(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators/", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators/", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]
	fmt.Printf("Operator1Data: %v\n", operator1Data)
	fmt.Printf("Operator2Data: %v\n", operator2Data)

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, false, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorInvalidDaysQueryParam(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil)
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=ten", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	fmt.Printf("Response: %v\n", res)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Assert the response body
	var responseBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error unmarshaling response body: %v", err)
	}
	expectedErrorMessage := "Invalid 'days' parameter"
	assert.Equal(t, expectedErrorMessage, responseBody["error"])

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorQueryDaysGreaterThan30(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorState := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorState[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil)
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=31", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	fmt.Printf("Response: %v\n", res)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Assert the response body
	var responseBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error unmarshaling response body: %v", err)
	}
	expectedErrorMessage := "Invalid 'days' parameter. Max value is 30"
	assert.Equal(t, expectedErrorMessage, responseBody["error"])

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorsMultipleOffline(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	fmt.Printf("Response: %v\n", response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]
	fmt.Printf("Operator1Data: %v\n", operator1Data)
	fmt.Printf("Operator2Data: %v\n", operator2Data)

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, false, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorOnline(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorState := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorState[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil)
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorState, nil)

	// Start test server for Operator
	closeServer, err := startTestGRPCServer("localhost:32007") // Let the OS assign a free port
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer() // Ensure the server is closed after the test

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, true, response.Data[0].IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorsMultipleOfflineOnline(t *testing.T) {
	// Skipping this test as repported being flaky but could not reproduce it locally
	t.Skip("Skipping testing in CI environment")

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)

	// Set up the mock calls for the two operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	// Start the test server for Operator 2
	closeServer, err := startTestGRPCServer("localhost:32009")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer()

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]

	responseJson := string(data)
	fmt.Printf("Response: %v\n", responseJson)

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, true, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorsMultipleOnline(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphTwoOperatorsDeregistered, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	// Start test server for Operator 1
	closeServer1, err := startTestGRPCServer("localhost:32007") // Let the OS assign a free port
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer1() // Ensure the server is closed after the test

	// Start test server for Operator 2
	closeServer2, err := startTestGRPCServer("localhost:32009") // Let the OS assign a free port
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer closeServer2() // Ensure the server is closed after the test

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	operator1Data := response.Data[0]
	operator2Data := response.Data[1]

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, true, operator1Data.IsOnline)

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, true, operator2Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchDeregisteredOperatorsMultipleOfflineSameBlock(t *testing.T) {

	defer goleak.VerifyNone(t)

	r := setUpRouter()

	indexedOperatorStates := make(map[core.OperatorID]*subgraph.DeregisteredOperatorInfo)
	indexedOperatorStates[core.OperatorID{0}] = subgraphDeregisteredOperatorInfo
	indexedOperatorStates[core.OperatorID{1}] = subgraphDeregisteredOperatorInfo2
	indexedOperatorStates[core.OperatorID{2}] = subgraphDeregisteredOperatorInfo3

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphThreeOperatorsDeregistered, nil)

	// Set up the mock calls for the three operators
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo1, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo2, nil).Once()
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfo3, nil).Once()
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, nil)

	mockSubgraphApi.On("QueryIndexedDeregisteredOperatorsForTimeWindow").Return(indexedOperatorStates, nil)

	r.GET("/v1/operators-info/deregistered-operators", testDataApiServer.FetchDeregisteredOperators)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/operators-info/deregistered-operators?days=14", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.DeregisteredOperatorsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 3, response.Meta.Size)
	assert.Equal(t, 3, len(response.Data))

	operator1Data := response.Data[0]

	assert.Equal(t, "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311", operator1Data.OperatorId)
	assert.Equal(t, uint(22), operator1Data.BlockNumber)
	assert.Equal(t, "localhost:32007", operator1Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	operator2Data := getOperatorData(response.Data, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312")
	operator3Data := getOperatorData(response.Data, "0xe24cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568313")

	assert.Equal(t, "0xe23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312", operator2Data.OperatorId)
	assert.Equal(t, uint(24), operator2Data.BlockNumber)
	assert.Equal(t, "localhost:32009", operator2Data.Socket)
	assert.Equal(t, false, operator1Data.IsOnline)

	assert.Equal(t, "0xe24cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568313", operator3Data.OperatorId)
	assert.Equal(t, uint(24), operator3Data.BlockNumber)
	assert.Equal(t, "localhost:32011", operator3Data.Socket)
	assert.Equal(t, false, operator3Data.IsOnline)

	// Reset the mock
	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestGetServiceAvailability(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	// Assuming "Churner" service is also expected to be SERVING for this test
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)
	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	service1Data := response.Data[0]
	service2Data := response.Data[1]

	assert.Equal(t, "Disperser", service1Data.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING.String(), service1Data.ServiceStatus)

	assert.Equal(t, "Churner", service2Data.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING.String(), service2Data.ServiceStatus)
}

func TestGetServiceAvailability_QueryDisperser(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)

	// Initialize the gRPC client pools
	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability?service-name=disperser", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	fmt.Printf("Response: %v\n", response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	serviceData := response.Data[0]
	assert.Equal(t, "Disperser", serviceData.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING.String(), serviceData.ServiceStatus)
}

func TestGetServiceAvailability_QueryInvalidService(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	// Assuming "Churner" service is also expected to be SERVING for this test
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)

	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability?service-name=encoder", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 0, response.Meta.Size)
	assert.Equal(t, 0, len(response.Data))
}

func TestGetServiceAvailability_HealthCheckError(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	})
	// Assuming "Churner" service is also expected to be SERVING for this test
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	})
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)

	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability?service-name=disperser", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	serviceData := response.Data[0]
	assert.Equal(t, "Disperser", serviceData.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING.String(), serviceData.ServiceStatus)
}

func TestGetServiceAvailability_HealthyUnHealthyService(t *testing.T) {
	r := setUpRouter()
	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	})
	// Assuming "Churner" service is also expected to be SERVING for this test
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	})
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)

	// Initialize the gRPC client pools
	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 2, len(response.Data))

	service1Data := response.Data[0]
	service2Data := response.Data[1]

	assert.Equal(t, "Disperser", service1Data.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING.String(), service1Data.ServiceStatus)

	assert.Equal(t, "Churner", service2Data.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING.String(), service2Data.ServiceStatus)
}

func TestGetServiceAvailability_QueryDisperser_MultipleRequests(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	// Assuming "Churner" service is also expected to be SERVING for this test
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPCConnection{}, mockHealthCheckService)

	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	var concurrentRequests sync.WaitGroup
	responses := make(chan *http.Response, 12) // Channel to collect responses

	for i := 0; i < 12; i++ {
		concurrentRequests.Add(1)
		go func() {
			defer concurrentRequests.Done()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability?service-name=disperser", nil)
			r.ServeHTTP(w, req)
			responses <- w.Result()
		}()
	}

	concurrentRequests.Wait() // Wait for all requests to be processed
	close(responses)          // Close the channel after all goroutines complete

	// Process responses
	for res := range responses {
		processResponse(t, res)
	}
}

func TestGetServiceAvailability_HealthCheckerNilConnection(t *testing.T) {
	r := setUpRouter()

	mockHealthCheckService := NewMockHealthCheckService()
	mockHealthCheckService.AddResponse("Disperser", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	mockHealthCheckService.AddResponse("Churner", &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	testDataApiServer = dataapi.NewServer(config, blobstore, prometheusClient, dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger), mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger), &MockGRPNilConnection{}, nil)

	// Initialize the gRPC client pools
	r.GET("/v1/eigenda/service-availability", testDataApiServer.GetEigenDAServiceAvailability)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/eigenda/service-availability?service-name=disperser", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	fmt.Printf("Response: %v\n", response)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, 1, response.Meta.Size)
	assert.Equal(t, 1, len(response.Data))

	serviceData := response.Data[0]
	assert.Equal(t, "Disperser", serviceData.ServiceName)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_UNKNOWN.String(), serviceData.ServiceStatus)
}

func setUpRouter() *gin.Engine {
	return gin.Default()
}

func queueBlob(t *testing.T, blob *core.Blob, queue disperser.BlobStore) disperser.BlobKey {
	key, err := queue.StoreBlob(context.Background(), blob, expectedRequestedAt)
	assert.NoError(t, err)
	return key
}

func markBlobConfirmed(t *testing.T, blob *core.Blob, key disperser.BlobKey, batchHeaderHash [32]byte, queue disperser.BlobStore) {
	// simulate blob confirmation
	var commitX, commitY fp.Element
	_, err := commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}

	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:      batchHeaderHash,
		BlobIndex:            expectedBlobIndex,
		SignatoryRecordHash:  expectedSignatoryRecordHash,
		ReferenceBlockNumber: expectedReferenceBlockNumber,
		BatchRoot:            expectedBatchRoot,
		BlobInclusionProof:   expectedInclusionProof,
		BlobCommitment: &encoding.BlobCommitments{
			Commitment: commitment,
			Length:     uint(expectedDataLength),
		},
		BatchID:                 expectedBatchId,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: expectedConfirmationBlockNumber,
		Fee:                     expectedFee,
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     key.BlobHash,
		MetadataHash: key.MetadataHash,
		BlobStatus:   disperser.Confirmed,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: blob.RequestHeader.SecurityParams,
			},
			RequestedAt: expectedRequestedAt,
			BlobSize:    uint(len(blob.Data)),
		},
	}

	expectedBlobCommitment = confirmationInfo.BlobCommitment
	updated, err := queue.MarkBlobConfirmed(context.Background(), metadata, confirmationInfo)
	assert.NoError(t, err)
	assert.Equal(t, disperser.Confirmed, updated.BlobStatus)
}

func makeTestBlob(quorumID core.QuorumID, adversityThreshold uint8) core.Blob {
	blob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: []*core.SecurityParam{
				{
					QuorumID:           quorumID,
					AdversaryThreshold: adversityThreshold,
				},
			},
		},
		Data: gettysburgAddressBytes,
	}
	return blob
}

// startTestGRPCServer starts a gRPC server on a specified address.
// It returns a function to stop the server.
func startTestGRPCServer(address string) (stopFunc func(), err error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	stopFunc = func() {
		grpcServer.Stop()
		lis.Close()
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return stopFunc, nil
}

// Helper to get OperatorData from response
func getOperatorData(operatorMetadtas []*dataapi.DeregisteredOperatorMetadata, operatorId string) dataapi.DeregisteredOperatorMetadata {

	for _, operatorMetadata := range operatorMetadtas {
		if operatorMetadata.OperatorId == operatorId {
			return *operatorMetadata
		}
	}
	return dataapi.DeregisteredOperatorMetadata{}

}

// processResponse processes a single http.Response and closes its body.
func processResponse(t *testing.T, res *http.Response) {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.ServiceAvailabilityResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.GreaterOrEqual(t, response.Meta.Size, 1)
	assert.GreaterOrEqual(t, len(response.Data), 1)

	if len(response.Data) > 0 {
		serviceData := response.Data[0]
		assert.Equal(t, "Disperser", serviceData.ServiceName)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING.String(), serviceData.ServiceStatus)
	}
}
