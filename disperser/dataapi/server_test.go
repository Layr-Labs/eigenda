package dataapi_test

import (
	"context"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	commock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	prommock "github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus/mock"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	subgraphmock "github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph/mock"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

var (
	//go:embed testdata/prometheus-response-sample.json
	mockPrometheusResponse string
	//go:embed testdata/prometheus-resp-avg-throughput.json
	mockPrometheusRespAvgThroughput string

	expectedBlobCommitment *core.BlobCommitments
	mockLogger             = &commock.Logger{}
	blobstore              = inmem.NewBlobStore()
	mockPrometheusApi      = &prommock.MockPrometheusApi{}
	prometheusClient       = dataapi.NewPrometheusClient(mockPrometheusApi, "test-cluster")
	mockSubgraphApi        = &subgraphmock.MockSubgraphApi{}
	subgraphClient         = dataapi.NewSubgraphClient(mockSubgraphApi, mockLogger)
	config                 = dataapi.Config{ServerMode: "test", SocketAddr: ":8080"}

	mockTx                          = &coremock.MockTransactor{}
	mockChainState, _               = coremock.MakeChainDataMock(core.OperatorIndex(1))
	testDataApiServer               = dataapi.NewServer(config, blobstore, prometheusClient, subgraphClient, mockTx, mockChainState, mockLogger, dataapi.NewMetrics(nil, "9001", mockLogger))
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
	mockSubgraphApi.On("QueryBatches").Return(subgraphBatches, nil)
	mockSubgraphApi.On("QueryOperators").Return(subgraphOperatorRegistereds, nil)
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

	mockSubgraphApi.On("QueryBatches").Return(subgraphBatches, nil)

	nonSigning := struct {
		NonSigners []struct {
			OperatorId graphql.String `graphql:"operatorId"`
		} `graphql:"nonSigners"`
	}{
		NonSigners: []struct {
			OperatorId graphql.String `graphql:"operatorId"`
		}{
			{OperatorId: "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568310"},
			{OperatorId: "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568312"},
		},
	}
	batchNonSigningOperatorIds := []*subgraph.BatchNonSigningOperatorIds{
		{
			NonSigning: nonSigning,
		},
	}

	mockSubgraphApi.On("QueryBatchNonSigningOperatorIdsInInterval").Return(batchNonSigningOperatorIds, nil).Once()
	mockSubgraphApi.On("QueryRegisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorRegistereds, nil)
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryBatchesByBlockTimestampRange").Return(subgraphBatches, nil)

	r.GET("/v1/metrics/operators_nonsigning_percentage", testDataApiServer.FetchOperatorsNonsigningPercentageHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/metrics/operators_nonsigning_percentage", nil)
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

	responseData := response.Data[0]
	operatorId := responseData.OperatorId
	assert.Equal(t, 2, response.Meta.Size)
	assert.Equal(t, 3, responseData.TotalBatches)
	assert.Equal(t, 1, responseData.TotalUnsignedBatches)
	assert.Equal(t, float64(33.33), responseData.Percentage)
	assert.Equal(t, "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568310", operatorId)
	assert.Equal(t, 2, len(response.Data))
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
	commitment := &core.G1Commitment{
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
		BlobCommitment: &core.BlobCommitments{
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
