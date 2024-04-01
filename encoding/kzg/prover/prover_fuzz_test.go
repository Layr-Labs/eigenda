package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
)

func FuzzOnlySystematic(f *testing.F) {
	f.Add(gettysburgAddressBytesIFFT)
	f.Fuzz(func(t *testing.T, input []byte) {

		group, _ := prover.NewProver(kzgConfig, true)

		params := encoding.ParamsFromSysPar(10, 3, uint64(len(input)))
		enc, err := group.GetKzgEncoder(params)
		if err != nil {
			t.Errorf("Error making rs: %q", err)
		}

		inputFr := rs.ToFrArrayWith254Bits(input)

		//encode the data
		_, _, _, frames, _, err := enc.Encode(inputFr)

		for _, frame := range frames {
			assert.NotEqual(t, len(frame.Coeffs), 0)
		}

		if err != nil {
			t.Errorf("Error Encoding:\n Data:\n %q \n Err: %q", input, err)
		}

		//sample the correct systematic frames
		samples, indices := sampleFrames(frames, uint64(len(frames)))

		data, err := enc.Decode(samples, indices, uint64(len(input)))
		if err != nil {
			t.Errorf("Error Decoding:\n Data:\n %q \n Err: %q", input, err)
		}
		assert.Equal(t, gettysburgAddressBytes, data[:len(gettysburgAddressBytes)], "Input data was not equal to the decoded data")
	})
}
