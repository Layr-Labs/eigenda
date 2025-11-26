package grpc_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/common/version"
	"github.com/Layr-Labs/eigenda/core/eth/operatorstate"
	"github.com/Layr-Labs/eigenda/test/random"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	coremockv2 "github.com/Layr-Labs/eigenda/core/mock/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigenda/node/grpc"
	nodemock "github.com/Layr-Labs/eigenda/node/mock"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	blobParams = &core.BlobVersionParameters{
		NumChunks:       8192,
		CodingRate:      8,
		MaxNumOperators: 2048,
	}
	blobParamsMap = map[v2.BlobVersion]*core.BlobVersionParameters{
		0: blobParams,
	}
)

type testComponents struct {
	server        *grpc.ServerV2
	node          *node.Node
	store         *nodemock.MockStoreV2
	validator     *coremockv2.MockShardValidator
	relayClient   *clientsmock.MockRelayClient
	disperserKey  *ecdsa.PrivateKey
	disperserAddr gethcommon.Address
}

func newTestComponents(t *testing.T, config *node.Config) *testComponents {
	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	require.NoError(t, err)
	signer, err := blssigner.NewSigner(blssignerTypes.SignerConfig{
		SignerType: blssignerTypes.PrivateKey,
		PrivateKey: keyPair.PrivKey.String(),
	})
	require.NoError(t, err)
	opID = [32]byte{0}
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	require.NoError(t, err)
	err = os.MkdirAll(config.DbPath, os.ModePerm)
	require.NoError(t, err)
	noopMetrics := metrics.NewNoopMetrics()
	reg := prometheus.NewRegistry()
	tx := &coremock.MockWriter{}

	rand := random.NewTestRandom()
	disperserAddr, disperserKey, err := rand.EthAccount()
	require.NoError(t, err)

	// Set up mock for disperser address lookup (used by authentication)
	tx.On("GetDisperserAddress", mock.Anything, mock.Anything).Return(disperserAddr, nil)
	// Set up mock for quorum count (used by blob validation)
	tx.On("GetQuorumCount", mock.Anything, mock.Anything).Return(uint8(3), nil)

	ratelimiter := &commonmock.NoopRatelimiter{}

	val := coremockv2.NewMockShardValidator()
	metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090", opID, -1, tx, chainState)

	operatorStateCache := operatorstate.NewMockOperatorStateCache()
	operatorState, err := chainState.GetOperatorState(t.Context(), 100, []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	operatorStateCache.SetOperatorState(t.Context(), 100, operatorState)

	s := nodemock.NewMockStoreV2()
	relay := clientsmock.NewRelayClient()
	var atomicRelayClient atomic.Value
	atomicRelayClient.Store(relay)
	node := &node.Node{
		Config:             config,
		Logger:             logger,
		KeyPair:            keyPair,
		BLSSigner:          signer,
		Metrics:            metrics,
		ValidatorStore:     s,
		ChainState:         chainState,
		ValidatorV2:        val,
		RelayClient:        atomicRelayClient,
		DownloadPool:       workerpool.New(1),
		ValidationPool:     workerpool.New(1),
		OperatorStateCache: operatorStateCache,
	}
	node.BlobVersionParams.Store(v2.NewBlobVersionParameterMap(blobParamsMap))
	// Set quorum count for validation
	node.QuorumCount.Store(3)

	// Create listeners with OS-allocated ports for testing
	v2DispersalListener, err := net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err)
	v2RetrievalListener, err := net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err)

	server, err := grpc.NewServerV2(
		context.Background(),
		config,
		node,
		logger,
		ratelimiter,
		prometheus.NewRegistry(),
		tx,
		version.DefaultVersion(),
		v2DispersalListener,
		v2RetrievalListener)

	require.NoError(t, err)
	return &testComponents{
		server:        server,
		node:          node,
		store:         s,
		validator:     val,
		relayClient:   relay,
		disperserKey:  disperserKey,
		disperserAddr: disperserAddr,
	}
}

// Signs a StoreChunksRequest with the test disperser key and sets the timestamp
func (c *testComponents) signRequest(t *testing.T, request *validator.StoreChunksRequest) {
	request.Timestamp = uint32(time.Now().Unix())
	signature, err := auth.SignStoreChunksRequest(c.disperserKey, request)
	require.NoError(t, err)
	request.Signature = signature
}

func TestV2NodeInfoRequest(t *testing.T) {
	c := newTestComponents(t, makeConfig(t))
	resp, err := c.server.GetNodeInfo(context.Background(), &validator.GetNodeInfoRequest{})
	require.NoError(t, err)
	require.Equal(t, resp.GetSemver(), version.DefaultVersion().String())
}

func TestV2ServerWithoutV2(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = false
	c := newTestComponents(t, config)
	_, err := c.server.StoreChunks(context.Background(), &validator.StoreChunksRequest{})
	requireErrorStatus(t, err, codes.InvalidArgument)

	_, err = c.server.GetChunks(context.Background(), &validator.GetChunksRequest{})
	requireErrorStatus(t, err, codes.InvalidArgument)
}

func TestV2StoreChunksInputValidation(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)
	_, batch, _ := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)

	req := &validator.StoreChunksRequest{
		DisperserID: 0,
	}
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "signature must be 65 bytes")

	req = &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch:       &pbcommon.Batch{},
	}
	c.signRequest(t, req)
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "failed to deserialize batch")

	req = &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch: &pbcommon.Batch{
			Header:           &pbcommon.BatchHeader{},
			BlobCertificates: batchProto.GetBlobCertificates(),
		},
	}
	c.signRequest(t, req)
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "failed to deserialize batch")

	req = &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch: &pbcommon.Batch{
			Header:           batchProto.GetHeader(),
			BlobCertificates: []*pbcommon.BlobCertificate{},
		},
	}
	c.signRequest(t, req)
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "failed to deserialize batch")
}

func TestV2StoreChunksSuccess(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)

	blobKeys, batch, bundles := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)

	c.validator.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	c.validator.On("ValidateBatchHeader", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)
	c.relayClient.On(
		"GetChunksByRange",
		mock.Anything,
		v2.RelayKey(0),
		mock.Anything,
	).Return([][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	c.relayClient.On(
		"GetChunksByRange",
		mock.Anything,
		v2.RelayKey(1),
		mock.Anything,
	).Return([][]byte{bundles10Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})
	c.store.On("StoreBatch", mock.Anything, mock.Anything).Return(nil, nil)
	request := &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch:       batchProto,
	}
	c.signRequest(t, request)
	reply, err := c.server.StoreChunks(t.Context(), request)
	require.NoError(t, err)
	require.NotNil(t, reply.GetSignature())
	sigBytes := reply.GetSignature()
	point, err := new(core.Signature).Deserialize(sigBytes)
	require.NoError(t, err)
	sig := &core.Signature{G1Point: point}
	bhh, err := batch.BatchHeader.Hash()
	require.NoError(t, err)
	require.True(t, sig.Verify(c.node.KeyPair.GetPubKeyG2(), bhh))
}

func TestV2StoreChunksDownloadFailure(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)

	_, batch, _ := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)
	c.validator.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	c.validator.On("ValidateBatchHeader", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	relayErr := errors.New("error")
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{}, relayErr)
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{}, relayErr)
	request := &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch:       batchProto,
	}
	c.signRequest(t, request)
	reply, err := c.server.StoreChunks(t.Context(), request)
	require.Nil(t, reply.GetSignature())
	requireErrorStatus(t, err, codes.Internal)
}

func TestV2StoreChunksStorageFailure(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)

	blobKeys, batch, bundles := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)

	c.validator.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	c.validator.On("ValidateBatchHeader", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)
	c.relayClient.On(
		"GetChunksByRange",
		mock.Anything,
		v2.RelayKey(0),
		mock.Anything,
	).Return([][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	c.relayClient.On(
		"GetChunksByRange",
		mock.Anything,
		v2.RelayKey(1),
		mock.Anything,
	).Return([][]byte{bundles10Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})
	c.store.On("StoreBatch", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	request := &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch:       batchProto,
	}
	c.signRequest(t, request)
	reply, err := c.server.StoreChunks(t.Context(), request)
	require.Nil(t, reply.GetSignature())
	requireErrorStatusAndMsg(t, err, codes.Internal, "failed to store batch")
}

func TestV2StoreChunksLevelDBValidationFailure(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)

	blobKeys, batch, bundles := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)

	c.validator.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(
		errors.New("error"))
	c.validator.On("ValidateBatchHeader", mock.Anything, mock.Anything, mock.Anything).Return(
		nil)
	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(0), mock.Anything).Return(
		[][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return(
		[][]byte{bundles10Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByRange)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})
	c.store.On("StoreBatch", mock.Anything, mock.Anything).Return([]kvstore.Key{mockKey{}}, nil)
	c.store.On("DeleteKeys", mock.Anything, mock.Anything).Return(nil)
	request := &validator.StoreChunksRequest{
		DisperserID: 0,
		Batch:       batchProto,
	}
	c.signRequest(t, request)
	reply, err := c.server.StoreChunks(context.Background(), request)
	require.Nil(t, reply.GetSignature())
	requireErrorStatus(t, err, codes.Internal)
}

func TestV2GetChunksInputValidation(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	c := newTestComponents(t, config)
	ctx := context.Background()
	req := &validator.GetChunksRequest{
		BlobKey: []byte{0},
	}
	_, err := c.server.GetChunks(ctx, req)
	requireErrorStatus(t, err, codes.InvalidArgument)

	bk := [32]byte{0}
	maxUInt32 := uint32(0xFFFFFFFF)
	req = &validator.GetChunksRequest{
		BlobKey:  bk[:],
		QuorumId: maxUInt32,
	}
	_, err = c.server.GetChunks(ctx, req)
	requireErrorStatus(t, err, codes.InvalidArgument)
}

func requireErrorStatus(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, s.Code(), code)
}

func requireErrorStatusAndMsg(t *testing.T, err error, code codes.Code, substring string) {
	requireErrorStatus(t, err, code)
	assert.True(t, strings.Contains(err.Error(), substring))
}

type mockKey struct{}
type mockKeyBuilder struct{}

var _ kvstore.Key = mockKey{}
var _ kvstore.KeyBuilder = mockKeyBuilder{}

func (mockKey) Bytes() []byte {
	return []byte{0}
}

func (mockKey) Raw() []byte {
	return []byte{0}
}

func (mockKey) Builder() kvstore.KeyBuilder {
	return &mockKeyBuilder{}
}

func (mockKeyBuilder) TableName() string {
	return "tableName"
}

func (mockKeyBuilder) Key(data []byte) kvstore.Key {
	return mockKey{}
}
