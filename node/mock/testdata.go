package mock

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

func MockBatch(t *testing.T) ([]v2.BlobKey, *v2.Batch, []map[core.QuorumID]core.Bundle) {
	commitments := MockCommitment(t)
	bh0 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			ReservationPeriod: 5,
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
			ReservationPeriod: 6,
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
			ReservationPeriod: 7,
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

	certs := []*v2.BlobCertificate{blobCert0, blobCert1, blobCert2}
	tree, err := v2.BuildMerkleTree(certs)
	require.NoError(t, err)
	var root [32]byte
	copy(root[:], tree.Root())
	return []v2.BlobKey{blobKey0, blobKey1, blobKey2}, &v2.Batch{
		BatchHeader: &v2.BatchHeader{
			BatchRoot:            root,
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: certs,
	}, []map[core.QuorumID]core.Bundle{bundles0, bundles1, bundles2}
}

func MockCommitment(t *testing.T) encoding.BlobCommitments {
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
