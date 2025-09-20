package encoding_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerDeserGnark(t *testing.T) {
	var XCoord, YCoord fp.Element
	_, err := XCoord.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = YCoord.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)

	numCoeffs := 64
	var f encoding.Frame
	f.Proof = encoding.Proof{
		X: XCoord,
		Y: YCoord,
	}
	for i := 0; i < numCoeffs; i++ {
		f.Coeffs = append(f.Coeffs, fr.NewElement(uint64(i)))
	}

	gnark, err := f.SerializeGnark()
	assert.Nil(t, err)
	// The gnark encoding via f.Serialize() will generate less bytes
	// than gob.
	assert.Equal(t, 32*(1+numCoeffs), len(gnark))
	gob, err := f.SerializeGob()
	assert.Nil(t, err)
	// 2080 with gnark v.s. 2574 with gob
	assert.Equal(t, 2574, len(gob))

	// Verify the deserialization can get back original data
	c, err := new(encoding.Frame).DeserializeGnark(gnark)
	assert.Nil(t, err)
	assert.True(t, f.Proof.Equal(&c.Proof))
	assert.Equal(t, len(f.Coeffs), len(c.Coeffs))
	for i := 0; i < len(f.Coeffs); i++ {
		assert.True(t, f.Coeffs[i].Equal(&c.Coeffs[i]))
	}

	// invalid length should return error
	_, err = new(encoding.Frame).DeserializeGnark([]byte{1, 2, 3})
	assert.ErrorContains(t, err, "chunk length must be at least")
}

func createFrames(b *testing.B, numFrames int) []encoding.Frame {
	var XCoord, YCoord fp.Element
	_, err := XCoord.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(b, err)
	_, err = YCoord.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(b, err)
	r := rand.New(rand.NewSource(2024))
	numCoeffs := 64
	frames := make([]encoding.Frame, numFrames)
	for n := 0; n < numFrames; n++ {
		frames[n].Proof = encoding.Proof{
			X: XCoord,
			Y: YCoord,
		}
		for i := 0; i < numCoeffs; i++ {
			frames[n].Coeffs = append(frames[n].Coeffs, fr.NewElement(r.Uint64()))
		}
	}
	return frames
}

// randomG1 generates a random G1 point. There is no direct way to generate a random G1 point in the bn254 library,
// but we can generate a random BLS key and steal the public key.
func randomG1() (*bn254.G1Point, error) {
	key, err := bn254.GenRandomBlsKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random BLS keys: %w", err)
	}
	return key.PubKey, nil
}

func TestSerializeFrameProof(t *testing.T) {
	g1, err := randomG1()
	require.NoError(t, err)

	proof := g1.G1Affine

	bytes := make([]byte, encoding.SerializedProofLength)
	err = encoding.SerializeFrameProof(proof, bytes)
	require.NoError(t, err)

	proof2, err := encoding.DeserializeFrameProof(bytes)
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

	bytes, err := encoding.SerializeFrameProofs(proofs)
	require.NoError(t, err)
	proofs2, err := encoding.DeserializeFrameProofs(bytes)
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

	bytes, err := encoding.SerializeFrameProofs(proofs)
	require.NoError(t, err)
	splitBytes, err := encoding.SplitSerializedFrameProofs(bytes)
	require.NoError(t, err)

	require.Equal(t, len(proofs), len(splitBytes))
	for i := 0; i < len(proofs); i++ {
		proof := &encoding.Proof{}
		err := proof.Unmarshal(splitBytes[i])
		require.NoError(t, err)
		require.True(t, proofs[i].Equal(proof))
	}
}

func BenchmarkFrameGobSerialization(b *testing.B) {
	numSamples := 64
	frames := createFrames(b, numSamples)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = frames[i%numSamples].SerializeGob()
	}
}

func BenchmarkFrameGnarkSerialization(b *testing.B) {
	numSamples := 64
	frames := createFrames(b, numSamples)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = frames[i%numSamples].SerializeGnark()
	}
}

func BenchmarkFrameGobDeserialization(b *testing.B) {
	numSamples := 64
	frames := createFrames(b, numSamples)
	bytes := make([][]byte, numSamples)
	for n := 0; n < numSamples; n++ {
		gob, _ := frames[n].SerializeGob()
		bytes[n] = gob
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = new(encoding.Frame).DeserializeGob(bytes[i%numSamples])
	}
}

func BenchmarkFrameGnarkDeserialization(b *testing.B) {
	numSamples := 64
	frames := createFrames(b, numSamples)
	bytes := make([][]byte, numSamples)
	for n := 0; n < numSamples; n++ {
		gnark, _ := frames[n].SerializeGnark()
		bytes[n] = gnark
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = new(encoding.Frame).DeserializeGnark(bytes[i%numSamples])
	}
}
