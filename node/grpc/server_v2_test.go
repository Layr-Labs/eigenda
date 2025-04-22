package grpc_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/gammazero/workerpool"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
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
		MaxNumOperators: 3537,
	}
	blobParamsMap = map[v2.BlobVersion]*core.BlobVersionParameters{
		0: blobParams,
	}
	ecdsaSig = []byte{
		0x2a, 0x3b, 0x4c, 0x5d, 0x6e, 0x7f, 0x8a, 0x9b,
		0x0c, 0x1d, 0x2e, 0x3f, 0x4a, 0x5b, 0x6c, 0x7d,
		0x8e, 0x9f, 0x0a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f,
		0x6a, 0x7b, 0x8c, 0x9d, 0x0e, 0x1f, 0x2a, 0x3b,
		0x4c, 0x5d, 0x6e, 0x7f, 0x8a, 0x9b, 0x0c, 0x1d,
		0x2e, 0x3f, 0x4a, 0x5b, 0x6c, 0x7d, 0x8e, 0x9f,
		0x0a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b,
		0x8c, 0x9d, 0x0e, 0x1f, 0x2a, 0x3b, 0x4c, 0x5d,
		0x66,
	}
)

type testComponents struct {
	server      *grpc.ServerV2
	node        *node.Node
	store       *nodemock.MockStoreV2
	validator   *coremockv2.MockShardValidator
	relayClient *clientsmock.MockRelayClient
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

	ratelimiter := &commonmock.NoopRatelimiter{}

	val := coremockv2.NewMockShardValidator()
	metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090", opID, -1, tx, chainState)

	s := nodemock.NewMockStoreV2()
	relay := clientsmock.NewRelayClient()
	var atomicRelayClient atomic.Value
	atomicRelayClient.Store(relay)
	node := &node.Node{
		Config:         config,
		Logger:         logger,
		KeyPair:        keyPair,
		BLSSigner:      signer,
		Metrics:        metrics,
		ValidatorStore: s,
		ChainState:     chainState,
		ValidatorV2:    val,
		RelayClient:    atomicRelayClient,
		DownloadPool:   workerpool.New(1),
	}
	node.BlobVersionParams.Store(v2.NewBlobVersionParameterMap(blobParamsMap))

	// The eth client is only utilized for StoreChunks validation, which is disabled in these tests
	var reader *coreeth.Reader

	server, err := grpc.NewServerV2(
		context.Background(),
		config,
		node,
		logger,
		ratelimiter,
		prometheus.NewRegistry(),
		reader)

	require.NoError(t, err)
	return &testComponents{
		server:      server,
		node:        node,
		store:       s,
		validator:   val,
		relayClient: relay,
	}
}

func TestV2NodeInfoRequest(t *testing.T) {
	c := newTestComponents(t, makeConfig(t))
	resp, err := c.server.GetNodeInfo(context.Background(), &validator.GetNodeInfoRequest{})
	assert.True(t, resp.Semver == "0.0.0")
	assert.True(t, err == nil)
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
		Signature:   ecdsaSig,
		Batch:       &pbcommon.Batch{},
	}
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "failed to deserialize batch")

	req = &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch: &pbcommon.Batch{
			Header:           &pbcommon.BatchHeader{},
			BlobCertificates: batchProto.BlobCertificates,
		},
	}
	_, err = c.server.StoreChunks(context.Background(), req)
	requireErrorStatusAndMsg(t, err, codes.InvalidArgument, "failed to deserialize batch")

	req = &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch: &pbcommon.Batch{
			Header:           batchProto.Header,
			BlobCertificates: []*pbcommon.BlobCertificate{},
		},
	}
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
	bundles01Bytes, err := bundles[0][1].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles11Bytes, err := bundles[1][1].Serialize()
	require.NoError(t, err)
	bundles21Bytes, err := bundles[2][1].Serialize()
	require.NoError(t, err)
	bundles22Bytes, err := bundles[2][2].Serialize()
	require.NoError(t, err)
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles01Bytes, bundles21Bytes, bundles22Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 4)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[0], requests[1].BlobKey)
		require.Equal(t, blobKeys[2], requests[2].BlobKey)
		require.Equal(t, blobKeys[2], requests[3].BlobKey)
	})
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{bundles10Bytes, bundles11Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
		require.Equal(t, blobKeys[1], requests[1].BlobKey)
	})
	c.store.On("StoreBatch", mock.Anything, mock.Anything).Return(nil, nil)
	reply, err := c.server.StoreChunks(context.Background(), &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch:       batchProto,
	})
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
	reply, err := c.server.StoreChunks(context.Background(), &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch:       batchProto,
	})
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
	bundles01Bytes, err := bundles[0][1].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles11Bytes, err := bundles[1][1].Serialize()
	require.NoError(t, err)
	bundles21Bytes, err := bundles[2][1].Serialize()
	require.NoError(t, err)
	bundles22Bytes, err := bundles[2][2].Serialize()
	require.NoError(t, err)
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles01Bytes, bundles21Bytes, bundles22Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 4)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[0], requests[1].BlobKey)
		require.Equal(t, blobKeys[2], requests[2].BlobKey)
		require.Equal(t, blobKeys[2], requests[3].BlobKey)
	})
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{bundles10Bytes, bundles11Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
		require.Equal(t, blobKeys[1], requests[1].BlobKey)
	})
	c.store.On("StoreBatch", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	reply, err := c.server.StoreChunks(context.Background(), &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch:       batchProto,
	})
	require.Nil(t, reply.GetSignature())
	requireErrorStatusAndMsg(t, err, codes.Internal, "failed to store batch")
}

func TestV2StoreChunksValidationFailure(t *testing.T) {
	config := makeConfig(t)
	config.EnableV2 = true
	config.LittDBEnabled = true
	c := newTestComponents(t, config)

	blobKeys, batch, bundles := nodemock.MockBatch(t)
	batchProto, err := batch.ToProtobuf()
	require.NoError(t, err)

	c.validator.On("ValidateBlobs", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
	c.validator.On("ValidateBatchHeader", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles01Bytes, err := bundles[0][1].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles11Bytes, err := bundles[1][1].Serialize()
	require.NoError(t, err)
	bundles21Bytes, err := bundles[2][1].Serialize()
	require.NoError(t, err)
	bundles22Bytes, err := bundles[2][2].Serialize()
	require.NoError(t, err)
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles01Bytes, bundles21Bytes, bundles22Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 4)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[0], requests[1].BlobKey)
		require.Equal(t, blobKeys[2], requests[2].BlobKey)
		require.Equal(t, blobKeys[2], requests[3].BlobKey)
	})
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{bundles10Bytes, bundles11Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
		require.Equal(t, blobKeys[1], requests[1].BlobKey)
	})
	c.store.On("StoreBatch", batch, mock.Anything).Return([]kvstore.Key{mockKey{}}, nil)
	c.store.On("DeleteKeys", mock.Anything, mock.Anything).Return(nil)
	reply, err := c.server.StoreChunks(context.Background(), &validator.StoreChunksRequest{
		DisperserID: 0,
		Signature:   ecdsaSig,
		Batch:       batchProto,
	})
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
