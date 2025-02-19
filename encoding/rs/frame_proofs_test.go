package rs_test

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/require"
	"testing"
)

// randomG1 generates a random G1 point. There is no direct way to generate a random G1 point in the bn254 library,
// but we can generate a random BLS key and steal the public key.
func randomG1() (*bn254.G1Point, error) {
	key, err := bn254.GenRandomBlsKeys()
	if err != nil {
		return nil, err
	}
	return key.PubKey, nil
}

func TestSerializeFrameProof(t *testing.T) {
	g1, err := randomG1()
	require.NoError(t, err)

	proof := g1.G1Affine

	bytes := make([]byte, rs.SerializedProofLength)
	err = rs.SerializeFrameProof(proof, bytes)
	require.NoError(t, err)

	proof2, err := rs.DeserializeFrameProof(bytes)
	require.NoError(t, err)

	require.True(t, proof.Equal(proof2))
}

func TestSerializeFrameProofs(t *testing.T) {
	rand := random.NewTestRandom()

	count := 10 + rand.Intn(10)
	proofs := make([]*encoding.Proof, count)

	for i := 0; i < count; i++ {
		g1, err := randomG1()
		require.NoError(t, err)
		proofs[i] = g1.G1Affine
	}

	bytes, err := rs.SerializeFrameProofs(proofs)
	require.NoError(t, err)
	proofs2, err := rs.DeserializeFrameProofs(bytes)
	require.NoError(t, err)

	require.Equal(t, len(proofs), len(proofs2))
	for i := 0; i < len(proofs); i++ {
		require.True(t, proofs[i].Equal(proofs2[i]))
	}
}

func TestSplitSerializedFrameProofs(t *testing.T) {
	rand := random.NewTestRandom()

	count := 10 + rand.Intn(10)
	proofs := make([]*encoding.Proof, count)

	for i := 0; i < count; i++ {
		g1, err := randomG1()
		require.NoError(t, err)
		proofs[i] = g1.G1Affine
	}

	bytes, err := rs.SerializeFrameProofs(proofs)
	require.NoError(t, err)
	splitBytes, err := rs.SplitSerializedFrameProofs(bytes)
	require.NoError(t, err)

	require.Equal(t, len(proofs), len(splitBytes))
	for i := 0; i < len(proofs); i++ {
		proof := &encoding.Proof{}
		err := proof.Unmarshal(splitBytes[i])
		require.NoError(t, err)
		require.True(t, proofs[i].Equal(proof))
	}
}
