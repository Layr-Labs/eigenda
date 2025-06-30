package node_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval/test"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	nodemock "github.com/Layr-Labs/eigenda/node/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDownloadBundles(t *testing.T) {
	c := newComponents(t, op0)
	c.node.RelayClient.Store(c.relayClient)
	ctx := context.Background()
	blobKeys, batch, bundles := nodemock.MockBatch(t)
	blobCerts := batch.BlobCertificates

	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)

	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{bundles10Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})
	state, err := c.node.ChainState.GetOperatorStateByOperator(ctx, uint(10), op0)
	require.NoError(t, err)
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch, state, nil)
	require.NoError(t, err)
	require.Len(t, blobShards, 3)
	require.Equal(t, blobCerts[0], blobShards[0].BlobCertificate)
	require.Equal(t, blobCerts[1], blobShards[1].BlobCertificate)
	require.Equal(t, blobCerts[2], blobShards[2].BlobCertificate)

	require.Len(t, rawBundles, 3)
	require.Equal(t, blobCerts[0], rawBundles[0].BlobCertificate)
	require.Equal(t, blobCerts[1], rawBundles[1].BlobCertificate)
	require.Equal(t, blobCerts[2], rawBundles[2].BlobCertificate)
}

func TestDownloadBundlesFail(t *testing.T) {
	c := newComponents(t, op0)
	c.node.RelayClient.Store(c.relayClient)
	ctx := context.Background()
	blobKeys, batch, bundles := nodemock.MockBatch(t)

	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)
	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	relayServerError := fmt.Errorf("relay server error")
	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(1), mock.Anything).Return(nil, relayServerError).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})
	state, err := c.node.ChainState.GetOperatorState(ctx, uint(10), []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch, state, nil)
	require.Error(t, err)
	require.Nil(t, blobShards)
	require.Nil(t, rawBundles)
}

func TestDownloadBundlesOnlyParticipatingQuorums(t *testing.T) {
	// Operator 3 is not participating in quorum 2, so it should only download bundles for quorums 0 and 1
	c := newComponents(t, op3)
	c.node.RelayClient.Store(c.relayClient)
	ctx := context.Background()
	blobKeys, batch, bundles := nodemock.MockBatch(t)
	blobCerts := batch.BlobCertificates

	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles10Bytes, err := bundles[1][0].Serialize()
	require.NoError(t, err)
	bundles20Bytes, err := bundles[2][0].Serialize()
	require.NoError(t, err)
	// there shouldn't be a request to quorum 2 for blobKeys[2]
	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(0), mock.Anything).Return([][]byte{bundles00Bytes, bundles20Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[0], requests[0].BlobKey)
		require.Equal(t, blobKeys[2], requests[1].BlobKey)
	})
	c.relayClient.On("GetChunksByIndex", mock.Anything, v2.RelayKey(1), mock.Anything).Return([][]byte{bundles10Bytes}, nil).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*relay.ChunkRequestByIndex)
		require.Len(t, requests, 1)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
	})

	state, err := c.node.ChainState.GetOperatorState(ctx, uint(10), []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch, state, nil)
	require.NoError(t, err)
	require.Len(t, blobShards, 3)
	require.Equal(t, blobCerts[0], blobShards[0].BlobCertificate)
	require.Equal(t, blobCerts[1], blobShards[1].BlobCertificate)
	require.Equal(t, blobCerts[2], blobShards[2].BlobCertificate)

	require.Len(t, rawBundles, 3)
	require.Equal(t, blobCerts[0], rawBundles[0].BlobCertificate)
	require.Equal(t, blobCerts[1], rawBundles[1].BlobCertificate)
	require.Equal(t, blobCerts[2], rawBundles[2].BlobCertificate)
}

func TestRefreshOnchainStateFailure(t *testing.T) {
	c := newComponents(t, op0)
	c.node.Config.EnableV2 = true
	c.node.RelayClient.Store(c.relayClient)
	c.node.Config.OnchainStateRefreshInterval = time.Millisecond
	ctx := context.Background()
	bp, ok := c.node.BlobVersionParams.Load().Get(0)
	require.True(t, ok)
	require.Equal(t, bp, blobParams)
	_, ok = c.node.BlobVersionParams.Load().Get(1)
	require.False(t, ok)
	relayClient, ok := c.node.RelayClient.Load().(relay.RelayClient)
	require.True(t, ok)
	require.NotNil(t, relayClient)

	// Both updates fail
	newCtx, cancel := context.WithTimeout(ctx, c.node.Config.OnchainStateRefreshInterval*2)
	defer cancel()

	c.tx.On("GetAllVersionedBlobParams", mock.Anything).Return(nil, assert.AnError)
	c.relayClient.On("GetSockets").Return(nil)
	c.tx.On("GetRelayURLs", mock.Anything).Return(nil, assert.AnError)
	c.tx.On("GetCurrentBlockNumber", mock.Anything).Return(uint32(10), nil)
	c.tx.On("GetQuorumCount", mock.Anything).Return(uint8(2), nil)

	err := c.node.RefreshOnchainState(newCtx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	bp, ok = c.node.BlobVersionParams.Load().Get(0)
	require.True(t, ok)
	require.Equal(t, bp, blobParams)
	_, ok = c.node.BlobVersionParams.Load().Get(1)
	require.False(t, ok)
	newRelayClient := c.node.RelayClient.Load().(relay.RelayClient)
	require.Same(t, relayClient, newRelayClient)
	quorumCount := c.node.QuorumCount.Load()
	require.Equal(t, quorumCount, uint32(2))

	// Same relay URLs shouldn't trigger update
	newCtx1, cancel1 := context.WithTimeout(ctx, c.node.Config.OnchainStateRefreshInterval*2)
	defer cancel1()

	c.tx.On("GetAllVersionedBlobParams", mock.Anything).Return(nil, assert.AnError)
	relayURLs := map[v2.RelayKey]string{
		0: "http://localhost:8080",
	}
	c.relayClient.On("GetSockets").Return(relayURLs).Once()
	c.tx.On("GetRelayURLs", mock.Anything).Return(relayURLs, nil)
	c.tx.On("GetCurrentBlockNumber", mock.Anything).Return(uint32(10), nil)
	c.tx.On("GetQuorumCount", mock.Anything).Return(uint8(3), nil)

	err = c.node.RefreshOnchainState(newCtx1)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	newRelayClient = c.node.RelayClient.Load().(relay.RelayClient)
	require.Same(t, relayClient, newRelayClient)
	quorumCount = c.node.QuorumCount.Load()
	require.Equal(t, quorumCount, uint32(2))
}

func TestRefreshOnchainStateSuccess(t *testing.T) {
	c := newComponents(t, op0)
	c.node.Config.EnableV2 = true
	c.node.Config.OnchainStateRefreshInterval = time.Millisecond

	relayUrlProvider := test.NewTestRelayUrlProvider()
	relayUrlProvider.StoreRelayUrl(0, "http://localhost:8080")

	messageSigner := func(ctx context.Context, data [32]byte) (*core.Signature, error) {
		return nil, nil
	}

	relayClientConfig := &relay.RelayClientConfig{
		OperatorID:         &c.node.Config.ID,
		MessageSigner:      messageSigner,
		MaxGRPCMessageSize: units.GiB,
	}

	relayClient, err := relay.NewRelayClient(relayClientConfig, c.node.Logger, relayUrlProvider)
	require.NoError(t, err)
	// set up non-mock client
	c.node.RelayClient.Store(relayClient)
	ctx := context.Background()
	bp, ok := c.node.BlobVersionParams.Load().Get(0)
	require.True(t, ok)
	require.Equal(t, bp, blobParams)
	_, ok = c.node.BlobVersionParams.Load().Get(1)
	require.False(t, ok)

	// Blob params updated successfully
	newCtx, cancel := context.WithTimeout(ctx, c.node.Config.OnchainStateRefreshInterval*2)
	defer cancel()

	blobParams2 := &core.BlobVersionParameters{
		NumChunks:       111,
		CodingRate:      1,
		MaxNumOperators: 2048,
	}
	c.tx.On("GetAllVersionedBlobParams", mock.Anything).Return(map[v2.BlobVersion]*core.BlobVersionParameters{
		0: blobParams,
		1: blobParams2,
	}, nil)
	c.tx.On("GetCurrentBlockNumber", mock.Anything).Return(uint32(10), nil)
	c.tx.On("GetQuorumCount", mock.Anything).Return(uint8(2), nil)

	err = c.node.RefreshOnchainState(newCtx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	bp, ok = c.node.BlobVersionParams.Load().Get(0)
	require.True(t, ok)
	require.Equal(t, bp, blobParams)
	bp, ok = c.node.BlobVersionParams.Load().Get(1)
	require.True(t, ok)
	require.Equal(t, bp, blobParams2)
	quorumCount := c.node.QuorumCount.Load()
	require.Equal(t, quorumCount, uint32(2))
}
