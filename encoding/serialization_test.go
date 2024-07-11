package encoding_test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
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
	gob, err := f.Serialize()
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

func BenchmarkFrameGobSerialization(b *testing.B) {
	numSamples := 64
	frames := createFrames(b, numSamples)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = frames[i%numSamples].Serialize()
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
		gob, _ := frames[n].Serialize()
		bytes[n] = gob
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = new(encoding.Frame).Deserialize(bytes[i%numSamples])
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
