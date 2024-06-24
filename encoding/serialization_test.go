package encoding_test

import (
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

	var f encoding.Frame
	f.Proof = encoding.Proof{
		X: XCoord,
		Y: YCoord,
	}
	for i := 0; i < 3; i++ {
		f.Coeffs = append(f.Coeffs, fr.NewElement(uint64(i)))
	}

	gnark, err := f.SerializeGnark()
	assert.Nil(t, err)
	// The gob encoding via f.Serialize() will generate 318 bytes
	// whereas gnark only 128 bytes
	assert.Equal(t, 128, len(gnark))
	gob, err := f.Serialize()
	assert.Nil(t, err)
	assert.Equal(t, 318, len(gob))

	// Verify the deserialization can get back original data
	c, err := new(encoding.Frame).DeserializeGnark(gnark)
	assert.Nil(t, err)
	assert.True(t, f.Proof.Equal(&c.Proof))
	assert.Equal(t, len(f.Coeffs), len(c.Coeffs))
	for i := 0; i < len(f.Coeffs); i++ {
		assert.True(t, f.Coeffs[i].Equal(&c.Coeffs[i]))
	}
}
