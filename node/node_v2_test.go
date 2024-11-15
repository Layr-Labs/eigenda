package node_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	nodemock "github.com/Layr-Labs/eigenda/node/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDownloadBundles(t *testing.T) {
	c := newComponents(t)
	ctx := context.Background()
	blobKeys, batch, bundles := nodemock.MockBatch(t)
	blobCerts := batch.BlobCertificates

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
	state, err := c.node.ChainState.GetOperatorState(ctx, uint(10), []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch, state)
	require.NoError(t, err)
	require.Len(t, blobShards, 3)
	require.Equal(t, blobCerts[0], blobShards[0].BlobCertificate)
	require.Equal(t, blobCerts[1], blobShards[1].BlobCertificate)
	require.Equal(t, blobCerts[2], blobShards[2].BlobCertificate)
	require.Contains(t, blobShards[0].Bundles, core.QuorumID(0))
	require.Contains(t, blobShards[0].Bundles, core.QuorumID(1))
	require.Contains(t, blobShards[1].Bundles, core.QuorumID(0))
	require.Contains(t, blobShards[1].Bundles, core.QuorumID(1))
	require.Contains(t, blobShards[2].Bundles, core.QuorumID(1))
	require.Contains(t, blobShards[2].Bundles, core.QuorumID(2))
	bundleEqual(t, bundles[0][0], blobShards[0].Bundles[0])
	bundleEqual(t, bundles[0][1], blobShards[0].Bundles[1])
	bundleEqual(t, bundles[1][0], blobShards[1].Bundles[0])
	bundleEqual(t, bundles[1][1], blobShards[1].Bundles[1])
	bundleEqual(t, bundles[2][1], blobShards[2].Bundles[1])
	bundleEqual(t, bundles[2][2], blobShards[2].Bundles[2])

	require.Len(t, rawBundles, 3)
	require.Equal(t, blobCerts[0], rawBundles[0].BlobCertificate)
	require.Equal(t, blobCerts[1], rawBundles[1].BlobCertificate)
	require.Equal(t, blobCerts[2], rawBundles[2].BlobCertificate)
	require.Contains(t, rawBundles[0].Bundles, core.QuorumID(0))
	require.Contains(t, rawBundles[0].Bundles, core.QuorumID(1))
	require.Contains(t, rawBundles[1].Bundles, core.QuorumID(0))
	require.Contains(t, rawBundles[1].Bundles, core.QuorumID(1))
	require.Contains(t, rawBundles[2].Bundles, core.QuorumID(1))
	require.Contains(t, rawBundles[2].Bundles, core.QuorumID(2))

	require.Equal(t, bundles00Bytes, rawBundles[0].Bundles[0])
	require.Equal(t, bundles01Bytes, rawBundles[0].Bundles[1])
	require.Equal(t, bundles10Bytes, rawBundles[1].Bundles[0])
	require.Equal(t, bundles11Bytes, rawBundles[1].Bundles[1])
	require.Equal(t, bundles21Bytes, rawBundles[2].Bundles[1])
	require.Equal(t, bundles22Bytes, rawBundles[2].Bundles[2])
}

func TestDownloadBundlesFail(t *testing.T) {
	c := newComponents(t)
	ctx := context.Background()
	blobKeys, batch, bundles := nodemock.MockBatch(t)

	bundles00Bytes, err := bundles[0][0].Serialize()
	require.NoError(t, err)
	bundles01Bytes, err := bundles[0][1].Serialize()
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
	relayServerError := fmt.Errorf("relay server error")
	c.relayClient.On("GetChunksByRange", mock.Anything, v2.RelayKey(1), mock.Anything).Return(nil, relayServerError).Run(func(args mock.Arguments) {
		requests := args.Get(2).([]*clients.ChunkRequestByRange)
		require.Len(t, requests, 2)
		require.Equal(t, blobKeys[1], requests[0].BlobKey)
		require.Equal(t, blobKeys[1], requests[1].BlobKey)
	})
	state, err := c.node.ChainState.GetOperatorState(ctx, uint(10), []core.QuorumID{0, 1, 2})
	require.NoError(t, err)
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch, state)
	require.Error(t, err)
	require.Nil(t, blobShards)
	require.Nil(t, rawBundles)
}

func bundleEqual(t *testing.T, expected, actual core.Bundle) {
	for i := range expected {
		frameEqual(t, expected[i], actual[i])
	}
}

func frameEqual(t *testing.T, expected, actual *encoding.Frame) {
	require.Equal(t, expected.Proof.Bytes(), actual.Proof.Bytes())
	require.Equal(t, expected.Coeffs, actual.Coeffs)
}
