package dataapi_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	blobMetadataStore   *blobstorev2.BlobMetadataStore
	testDataApiServerV2 *dataapi.ServerV2

	logger = logging.NewNoopLogger()

	// Local stack
	localStackPort     = "4566"
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	deployLocalStack   bool
)

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
	testDataApiServerV2 = dataapi.NewServerV2(config, blobMetadataStore, prometheusClient, subgraphClient, mockTx, mockChainState, mockIndexedChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger))
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
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/blobs/"+blobKey.Hex(), nil)
	r.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.BlobResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "Queued", response.Status)
	assert.Equal(t, uint16(0), response.BlobHeader.BlobVersion)
	assert.Equal(t, blobHeader.Signature, response.BlobHeader.Signature)
	assert.Equal(t, blobHeader.PaymentMetadata.AccountID, response.BlobHeader.PaymentMetadata.AccountID)
	assert.Equal(t, blobHeader.PaymentMetadata.ReservationPeriod, response.BlobHeader.PaymentMetadata.ReservationPeriod)
	assert.Equal(t, blobHeader.PaymentMetadata.CumulativePayment, response.BlobHeader.PaymentMetadata.CumulativePayment)
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
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/batches/"+batchHeaderHash, nil)
	r.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.BatchResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
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

	operator_id := "0xa96bfb4a7ca981ad365220f336dc5a3de0816ebd5130b79bbc85aca94bc9b6ab"
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(operatorInfo, nil)

	r.GET("/v2/operators/reachability", testDataApiServerV2.CheckOperatorsReachability)

	w := httptest.NewRecorder()
	reqStr := fmt.Sprintf("/v2/operators/reachability?operator_id=%v", operator_id)
	req := httptest.NewRequest(http.MethodGet, reqStr, nil)
	ctxWithDeadline, cancel := context.WithTimeout(req.Context(), 500*time.Microsecond)
	defer cancel()
	req = req.WithContext(ctxWithDeadline)
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.OperatorPortCheckResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, "23.93.76.1:32005", response.DispersalSocket)
	assert.Equal(t, false, response.DispersalOnline)
	assert.Equal(t, "23.93.76.1:32006", response.RetrievalSocket)
	assert.Equal(t, false, response.RetrievalOnline)

	mockSubgraphApi.ExpectedCalls = nil
	mockSubgraphApi.Calls = nil
}

func TestFetchOperatorsStake(t *testing.T) {
	r := setUpRouter()

	mockIndexedChainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	r.GET("/v2/operators/stake", testDataApiServerV2.FetchOperatorsStake)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/operators/stake", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.OperatorsStakeResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

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

	req := httptest.NewRequest(http.MethodGet, "/v2/metrics/summary", nil)
	req.Close = true
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response dataapi.MetricSummary
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, http.StatusOK, res.StatusCode)
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

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/metrics/timeseries/throughput", nil)
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
	assert.Equal(t, 3361, len(response))
	assert.Equal(t, float64(12000), response[0].Throughput)
	assert.Equal(t, uint64(1701292920), response[0].Timestamp)
	assert.Equal(t, float64(3.503022666666651e+07), totalThroughput)
}
