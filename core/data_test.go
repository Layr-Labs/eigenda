package core_test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
)

func createBundle(t *testing.T, numFrames, numCoeffs, seed int) core.Bundle {
	var XCoord, YCoord fp.Element
	_, err := XCoord.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = YCoord.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	r := rand.New(rand.NewSource(int64(seed)))
	frames := make([]*encoding.Frame, numFrames)
	for n := 0; n < numFrames; n++ {
		frames[n] = new(encoding.Frame)
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

func TestInvalidBundleSer(t *testing.T) {
	b1 := createBundle(t, 1, 0, 0)
	_, err := b1.Serialize()
	assert.EqualError(t, err, "invalid bundle: the coeffs length is zero")

	b2 := createBundle(t, 1, 1, 0)
	b3 := createBundle(t, 1, 2, 0)
	b3 = append(b3, b2...)
	_, err = b3.Serialize()
	assert.EqualError(t, err, "invalid bundle: all chunks should have the same length")
}

func TestInvalidBundleDeser(t *testing.T) {
	tooSmallBytes := []byte{byte(0b01000000)}
	_, err := new(core.Bundle).Deserialize(tooSmallBytes)
	assert.EqualError(t, err, "bundle data must have at least 8 bytes")

	invalidFormat := make([]byte, 0, 8)
	for i := 0; i < 7; i++ {
		invalidFormat = append(invalidFormat, byte(0))
	}
	invalidFormat = append(invalidFormat, byte(0b01000000))
	_, err = new(core.Bundle).Deserialize(invalidFormat)
	assert.EqualError(t, err, "invalid bundle data encoding format")

	invliadChunkLen := make([]byte, 0, 8)
	for i := 0; i < 7; i++ {
		invliadChunkLen = append(invliadChunkLen, byte(0))
	}
	invliadChunkLen = append(invliadChunkLen, byte(1))
	_, err = new(core.Bundle).Deserialize(invliadChunkLen)
	assert.EqualError(t, err, "chunk length must be greater than zero")

	data := make([]byte, 0, 9)
	for i := 0; i < 6; i++ {
		data = append(data, byte(0))
	}
	data = append(data, byte(0b00100000))
	data = append(data, byte(1))
	data = append(data, byte(5))
	data = append(data, byte(0b01000000))
	_, err = new(core.Bundle).Deserialize(data)
	assert.EqualError(t, err, "bundle data is invalid")
}

func TestBundleEncoding(t *testing.T) {
	numTrials := 16
	for i := 0; i < numTrials; i++ {
		bundle := createBundle(t, 64, 64, i)
		bytes, err := bundle.Serialize()
		assert.Nil(t, err)
		decoded, err := new(core.Bundle).Deserialize(bytes)
		assert.Nil(t, err)
		assert.Equal(t, len(bundle), len(decoded))
		for i := 0; i < len(bundle); i++ {
			assert.True(t, bundle[i].Proof.Equal(&decoded[i].Proof))
			assert.Equal(t, len(bundle[i].Coeffs), len(decoded[i].Coeffs))
			for j := 0; j < len(bundle[i].Coeffs); j++ {
				assert.True(t, bundle[i].Coeffs[j].Equal(&decoded[i].Coeffs[j]))
			}
		}
	}
}
