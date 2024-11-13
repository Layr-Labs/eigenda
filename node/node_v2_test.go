package node_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mockBatch(t *testing.T) ([]v2.BlobKey, *v2.Batch, []map[core.QuorumID]core.Bundle) {
	commitments := mockCommitment(t)
	bh0 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			BinIndex:          5,
			CumulativePayment: big.NewInt(100),
		},
		Signature: []byte{1, 2, 3},
	}
	bh1 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x456",
			BinIndex:          6,
			CumulativePayment: big.NewInt(200),
		},
		Signature: []byte{1, 2, 3},
	}
	bh2 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{1, 2},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x789",
			BinIndex:          7,
			CumulativePayment: big.NewInt(300),
		},
		Signature: []byte{1, 2, 3},
	}
	blobKey0, err := bh0.BlobKey()
	require.NoError(t, err)
	blobKey1, err := bh1.BlobKey()
	require.NoError(t, err)
	blobKey2, err := bh2.BlobKey()
	require.NoError(t, err)

	// blobCert 0 and blobCert 2 will be downloaded from relay 0
	// blobCert 1 will be downloaded from relay 1
	blobCert0 := &v2.BlobCertificate{
		BlobHeader: bh0,
		RelayKeys:  []v2.RelayKey{0},
	}
	blobCert1 := &v2.BlobCertificate{
		BlobHeader: bh1,
		RelayKeys:  []v2.RelayKey{1},
	}
	blobCert2 := &v2.BlobCertificate{
		BlobHeader: bh2,
		RelayKeys:  []v2.RelayKey{0},
	}

	bundles0 := map[core.QuorumID]core.Bundle{
		0: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(1), big.NewInt(2)).G1Affine),
				Coeffs: []fr.Element{
					{1, 2, 3, 4},
					{5, 6, 7, 8},
				},
			},
		},
		1: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(3), big.NewInt(4)).G1Affine),
				Coeffs: []fr.Element{
					{9, 10, 11, 12},
					{13, 14, 15, 16},
				},
			},
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(5), big.NewInt(6)).G1Affine),
				Coeffs: []fr.Element{
					{17, 18, 19, 20},
					{21, 22, 23, 24},
				},
			},
		},
	}
	bundles1 := map[core.QuorumID]core.Bundle{
		0: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(7), big.NewInt(8)).G1Affine),
				Coeffs: []fr.Element{
					{25, 26, 27, 28},
					{29, 30, 31, 32},
				},
			},
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(9), big.NewInt(10)).G1Affine),
				Coeffs: []fr.Element{
					{33, 34, 35, 36},
					{37, 38, 39, 40},
				},
			},
		},
		1: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(11), big.NewInt(12)).G1Affine),
				Coeffs: []fr.Element{
					{41, 42, 43, 44},
					{45, 46, 47, 48},
				},
			},
		},
	}
	bundles2 := map[core.QuorumID]core.Bundle{
		1: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(13), big.NewInt(14)).G1Affine),
				Coeffs: []fr.Element{
					{49, 50, 51, 52},
					{53, 54, 55, 56},
				},
			},
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(15), big.NewInt(16)).G1Affine),
				Coeffs: []fr.Element{
					{57, 58, 59, 60},
					{61, 62, 63, 64},
				},
			},
		},
		2: {
			{
				Proof: encoding.Proof(*core.NewG1Point(big.NewInt(17), big.NewInt(18)).G1Affine),
				Coeffs: []fr.Element{
					{65, 66, 67, 68},
					{69, 70, 71, 72},
				},
			},
		},
	}

	return []v2.BlobKey{blobKey0, blobKey1, blobKey2}, &v2.Batch{
		BatchHeader: &v2.BatchHeader{
			BatchRoot:            [32]byte{1, 1, 1},
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: []*v2.BlobCertificate{blobCert0, blobCert1, blobCert2},
	}, []map[core.QuorumID]core.Bundle{bundles0, bundles1, bundles2}
}

func TestDownloadBundles(t *testing.T) {
	c := newComponents(t)
	ctx := context.Background()
	blobKeys, batch, bundles := mockBatch(t)
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
	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch)
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
	blobKeys, batch, bundles := mockBatch(t)

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

	blobShards, rawBundles, err := c.node.DownloadBundles(ctx, batch)
	require.Error(t, err)
	require.Nil(t, blobShards)
	require.Nil(t, rawBundles)
}

func mockCommitment(t *testing.T) encoding.BlobCommitments {
	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	require.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	require.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	require.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	require.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	return encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           10,
	}
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
