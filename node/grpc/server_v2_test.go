package grpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/grpc"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestServerV2(t *testing.T, mockValidator bool) *grpc.ServerV2 {
	return newTestServerV2WithConfig(t, mockValidator, makeConfig(t))
}

func newTestServerV2WithConfig(t *testing.T, mockValidator bool, config *node.Config) *grpc.ServerV2 {
	var err error
	keyPair, err = core.GenRandomBlsKeys()
	if err != nil {
		panic("failed to create a BLS Key")
	}
	opID = [32]byte{}
	copy(opID[:], []byte(fmt.Sprintf("%d", 3)))
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic("failed to create a logger")
	}

	err = os.MkdirAll(config.DbPath, os.ModePerm)
	if err != nil {
		panic("failed to create a directory for db")
	}
	noopMetrics := metrics.NewNoopMetrics()
	reg := prometheus.NewRegistry()
	tx := &coremock.MockWriter{}

	ratelimiter := &commonmock.NoopRatelimiter{}

	var val core.ShardValidator

	if mockValidator {
		mockVal := coremock.NewMockShardValidator()
		mockVal.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockVal.On("ValidateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		val = mockVal
	} else {

		_, v, err := makeTestComponents()
		if err != nil {
			panic("failed to create test encoder")
		}

		asn := &core.StdAssignmentCoordinator{}

		cst, err := coremock.MakeChainDataMock(map[uint8]int{
			0: 10,
			1: 10,
			2: 10,
		})
		if err != nil {
			panic("failed to create test encoder")
		}

		val = core.NewShardValidator(v, asn, cst, opID)
	}

	metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090", opID, -1, tx, chainState)
	store, err := node.NewLevelDBStore(config.DbPath, logger, metrics, 1e9, 1e9)
	if err != nil {
		panic("failed to create a new levelDB store")
	}

	node := &node.Node{
		Config:     config,
		Logger:     logger,
		KeyPair:    keyPair,
		Metrics:    metrics,
		Store:      store,
		ChainState: chainState,
		Validator:  val,
	}
	return grpc.NewServerV2(config, node, logger, ratelimiter)
}

func TestV2NodeInfoRequest(t *testing.T) {
	server := newTestServerV2(t, true)
	resp, err := server.NodeInfo(context.Background(), &pbv2.NodeInfoRequest{})
	assert.True(t, resp.Semver == "0.0.0")
	assert.True(t, err == nil)
}

func TestV2StoreChunks(t *testing.T) {
	server := newTestServerV2(t, true)
	_, err := server.StoreChunks(context.Background(), &pbv2.StoreChunksRequest{
		BlobCertificates: []*commonpb.BlobCertificate{},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func TestV2GetChunks(t *testing.T) {
	server := newTestServerV2(t, true)

	_, err := server.GetChunks(context.Background(), &pbv2.GetChunksRequest{
		BlobKey: []byte{0},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func GetV2BlobCertificate(t *testing.T) {
	server := newTestServerV2(t, true)

	_, err := server.GetBlobCertificate(context.Background(), &pbv2.GetBlobCertificateRequest{
		BlobKey: []byte{0},
	})
	assert.ErrorContains(t, err, "not implemented")
}
