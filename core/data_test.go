package core_test

import (
	"bytes"
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

func createChunksData(t *testing.T, seed int) (core.Bundle, *core.ChunksData, *core.ChunksData) {
	bundle := createBundle(t, 64, 64, seed)
	gobChunks := make([][]byte, len(bundle))
	gnarkChunks := make([][]byte, len(bundle))
	for i, frame := range bundle {
		gobChunk, err := frame.Serialize()
		assert.Nil(t, err)
		gobChunks[i] = gobChunk

		gnarkChunk, err := frame.SerializeGnark()
		assert.Nil(t, err)
		gnarkChunks[i] = gnarkChunk
	}
	gob := &core.ChunksData{
		Chunks:   gobChunks,
		Format:   core.GobChunkEncodingFormat,
		ChunkLen: 64,
	}
	gnark := &core.ChunksData{
		Chunks:   gnarkChunks,
		Format:   core.GnarkChunkEncodingFormat,
		ChunkLen: 64,
	}
	return bundle, gob, gnark
}

func TestChunksData(t *testing.T) {
	numTrials := 16
	for i := 0; i < numTrials; i++ {
		bundle, gob, gnark := createChunksData(t, i)
		assert.Equal(t, len(gob.Chunks), 64)
		assert.Equal(t, len(gnark.Chunks), 64)
		assert.Equal(t, gnark.Size(), uint64(64*(32+64*encoding.BYTES_PER_SYMBOL)))
		// ToGobFormat
		convertedGob, err := gob.ToGobFormat()
		assert.Nil(t, err)
		assert.Equal(t, convertedGob, gob)
		convertedGob, err = gnark.ToGobFormat()
		assert.Nil(t, err)
		assert.Equal(t, len(gob.Chunks), len(convertedGob.Chunks))
		for i := 0; i < len(gob.Chunks); i++ {
			assert.True(t, bytes.Equal(gob.Chunks[i], convertedGob.Chunks[i]))
		}
		// ToGnarkFormat
		convertedGnark, err := gnark.ToGnarkFormat()
		assert.Nil(t, err)
		assert.Equal(t, convertedGnark, gnark)
		convertedGnark, err = gob.ToGnarkFormat()
		assert.Nil(t, err)
		assert.Equal(t, len(gnark.Chunks), len(convertedGnark.Chunks))
		for i := 0; i < len(gnark.Chunks); i++ {
			assert.True(t, bytes.Equal(gnark.Chunks[i], convertedGnark.Chunks[i]))
		}
		// FlattenToBundle
		bytesFromChunksData, err := gnark.FlattenToBundle()
		assert.Nil(t, err)
		bytesFromBundle, err := bundle.Serialize()
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(bytesFromChunksData, bytesFromBundle))
		// Invalid cases
		gnark.Chunks[0] = gnark.Chunks[0][1:]
		_, err = gnark.FlattenToBundle()
		assert.EqualError(t, err, "all chunks must be of same size")
		_, err = gob.FlattenToBundle()
		assert.EqualError(t, err, "unsupported chunk encoding format to flatten: 0")
		gob.Format = core.ChunkEncodingFormat(3)
		_, err = gob.ToGobFormat()
		assert.EqualError(t, err, "unsupported chunk encoding format: 3")
		_, err = gob.ToGnarkFormat()
		assert.EqualError(t, err, "unsupported chunk encoding format: 3")
	}
}
