package verification

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/stretchr/testify/require"
)

func TestAttestationProtoToBinding(t *testing.T) {
	var X0, Y0, X1, Y1 fp.Element
	_, err := X0.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	require.NoError(t, err)
	_, err = Y0.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	require.NoError(t, err)
	_, err = X1.SetString("18730744272503541936633286178165146673834730535090946570310418711896464442549")
	require.NoError(t, err)
	_, err = Y1.SetString("15356431458378126778840641829778151778222945686256112821552210070627093656047")
	require.NoError(t, err)

	pt0 := &core.G1Point{
		G1Affine: &bn254.G1Affine{
			X: X0,
			Y: Y0,
		},
	}
	pt1 := &core.G1Point{
		G1Affine: &bn254.G1Affine{
			X: X1,
			Y: Y1,
		},
	}

	var e0, e1, e2, e3 fp.Element
	_, err = e0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	require.NoError(t, err)
	_, err = e1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	require.NoError(t, err)
	_, err = e2.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	require.NoError(t, err)
	_, err = e3.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	require.NoError(t, err)

	var apk bn254.G2Affine
	apk.X.A0 = e0
	apk.X.A1 = e1
	apk.Y.A0 = e2
	apk.Y.A1 = e3

	inputAttestation := &v2.Attestation{
		NonSignerPubKeys: []*core.G1Point{pt0, pt1},
		APKG2: &core.G2Point{
			G2Affine: &apk,
		},
		QuorumAPKs: map[uint8]*core.G1Point{
			0: pt0,
			3: pt0,
			2: pt1,
		},
		Sigma: &core.Signature{
			G1Point: pt0,
		},
		QuorumNumbers: []core.QuorumID{3, 0, 2},
		QuorumResults: map[uint8]uint8{
			0: 100,
			3: 50,
			2: 25,
		},
	}
	attestationProtobuf, err := inputAttestation.ToProtobuf()
	require.NoError(t, err)

	bindingAttestation, err := attestationProtoToBinding(attestationProtobuf)
	require.NoError(t, err)

	require.Equal(t, len(inputAttestation.NonSignerPubKeys), len(bindingAttestation.NonSignerPubkeys))
	for i := range inputAttestation.NonSignerPubKeys {
		require.Equal(t, inputAttestation.NonSignerPubKeys[i].G1Affine.X.BigInt(new(big.Int)).Bytes(), bindingAttestation.NonSignerPubkeys[i].X.Bytes())
		require.Equal(t, inputAttestation.NonSignerPubKeys[i].G1Affine.Y.BigInt(new(big.Int)).Bytes(), bindingAttestation.NonSignerPubkeys[i].Y.Bytes())
	}

	require.Equal(t, inputAttestation.APKG2.G2Affine.X.A0.BigInt(new(big.Int)).Bytes(), bindingAttestation.ApkG2.X[1].Bytes())
	require.Equal(t, inputAttestation.APKG2.G2Affine.X.A1.BigInt(new(big.Int)).Bytes(), bindingAttestation.ApkG2.X[0].Bytes())
	require.Equal(t, inputAttestation.APKG2.G2Affine.Y.A0.BigInt(new(big.Int)).Bytes(), bindingAttestation.ApkG2.Y[1].Bytes())
	require.Equal(t, inputAttestation.APKG2.G2Affine.Y.A1.BigInt(new(big.Int)).Bytes(), bindingAttestation.ApkG2.Y[0].Bytes())

	require.Equal(t, len(inputAttestation.QuorumAPKs), len(bindingAttestation.QuorumApks))
	require.Equal(t, bindingAttestation.QuorumApks[0].X.Bytes(), pt0.G1Affine.X.BigInt(new(big.Int)).Bytes())
	require.Equal(t, bindingAttestation.QuorumApks[0].Y.Bytes(), pt0.G1Affine.Y.BigInt(new(big.Int)).Bytes())
	require.Equal(t, bindingAttestation.QuorumApks[1].X.Bytes(), pt1.G1Affine.X.BigInt(new(big.Int)).Bytes())
	require.Equal(t, bindingAttestation.QuorumApks[1].Y.Bytes(), pt1.G1Affine.Y.BigInt(new(big.Int)).Bytes())
	require.Equal(t, bindingAttestation.QuorumApks[2].X.Bytes(), pt0.G1Affine.X.BigInt(new(big.Int)).Bytes())
	require.Equal(t, bindingAttestation.QuorumApks[2].Y.Bytes(), pt0.G1Affine.Y.BigInt(new(big.Int)).Bytes())

	require.Equal(t, inputAttestation.Sigma.G1Point.G1Affine.X.BigInt(new(big.Int)).Bytes(), bindingAttestation.Sigma.X.Bytes())
	require.Equal(t, inputAttestation.Sigma.G1Point.G1Affine.Y.BigInt(new(big.Int)).Bytes(), bindingAttestation.Sigma.Y.Bytes())

	require.Equal(t, bindingAttestation.QuorumNumbers, []uint32{0, 2, 3})
}
